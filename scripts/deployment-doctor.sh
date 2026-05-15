#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

STRICT_DOCTOR="${STRICT_DOCTOR:-false}"
CONNECT_TIMEOUT_SECONDS="${CONNECT_TIMEOUT_SECONDS:-5}"
REQUEST_TIMEOUT_SECONDS="${REQUEST_TIMEOUT_SECONDS:-30}"
MAX_FEED_BYTES="${MAX_FEED_BYTES:-104857600}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
FORCE="${FORCE:-false}"
TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/deployment-doctor/$TIMESTAMP}"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/deployment-doctor.sh

Environment:
  PUBLIC_BASE_URL                 Public feed root; no default is invented
  ADMIN_BASE_URL                  Admin root; defaults to PUBLIC_BASE_URL only when loopback
  ADMIN_TOKEN                     Optional admin bearer token
  DATABASE_URL                    Optional database URL; value is never printed
  DB_MAX_CONNS                    Optional app DB pool cap for small-host review
  MIGRATIONS_DIR                  Migrations directory, default db/migrations when DATABASE_URL is set
  BACKUP_DIR                      Optional backup destination path
  RESTORE_DATABASE_URL            Optional restore-drill target URL; value is never printed
  RESTORE_DRILL_DATABASE_URL      Optional alias for RESTORE_DATABASE_URL
  RESTORE_BACKUP_FILE             Optional restore-drill backup file path
  RESTORE_DRILL_BACKUP_FILE       Optional alias for RESTORE_BACKUP_FILE
  STRICT_DOCTOR                   true|false; fail on blockers only when true
  OUTPUT_DIR                      Default .cache/deployment-doctor/<timestamp>
  ALLOW_UNIGNORED_OUTPUT_DIR      true|false; allow OUTPUT_DIR outside repo .cache
  FORCE                           true|false; allow non-empty OUTPUT_DIR reuse
  CONNECT_TIMEOUT_SECONDS         curl connect timeout, default 5
  REQUEST_TIMEOUT_SECONDS         curl total timeout, default 30
  MAX_FEED_BYTES                  maximum public feed size, default 104857600

Safety:
  The doctor inspects already-exported environment variables only.
  It does not source private env files, run migrations, create backups, restore
  databases, contact consumers, create evidence packets, or change statuses.
  Secret values, Authorization headers, cookies, database URLs, webhook values,
  and raw env files are never printed or written.
EOF
}

log() {
  printf '\n==> %s\n' "$1"
}

fail() {
  printf '\nERROR: %s\n' "$1" >&2
  exit 1
}

cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT INT TERM

bool_var() {
  name="$1"
  value="$2"
  case "$value" in
    true|false) ;;
    *) fail "$name must be true or false" ;;
  esac
}

positive_int() {
  name="$1"
  value="$2"
  case "$value" in
    ''|*[!0-9]*) fail "$name must be a positive integer" ;;
    0) fail "$name must be greater than zero" ;;
  esac
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

output_realpath() {
  python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$ALLOW_UNIGNORED_OUTPUT_DIR" <<'PY'
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out = pathlib.Path(sys.argv[2])
allow = sys.argv[3] == "true"
if not out.is_absolute():
    out = root / out
out = out.resolve(strict=False)
cache = (root / ".cache").resolve(strict=False)
if not allow:
    try:
        out.relative_to(cache)
    except ValueError:
        raise SystemExit(f"OUTPUT_DIR must resolve under {cache} unless ALLOW_UNIGNORED_OUTPUT_DIR=true")
print(out)
PY
}

prepare_output_dir() {
  if [ -L "$OUTPUT_DIR" ]; then
    fail "OUTPUT_DIR must not be a symlink: $OUTPUT_DIR"
  fi
  OUT_REAL="$(output_realpath)"
  if [ -L "$OUT_REAL" ]; then
    fail "OUTPUT_DIR must not be a symlink: $OUT_REAL"
  fi
  if [ -d "$OUT_REAL" ] && [ "$FORCE" != "true" ] && [ "$(find "$OUT_REAL" -mindepth 1 -maxdepth 1 | sed -n '1p')" ]; then
    fail "OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it"
  fi
  mkdir -p "$OUT_REAL"
  chmod 700 "$OUT_REAL"
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-deployment-doctor.XXXXXX")"
  CHECKS_TSV="$OUT_REAL/checks.tsv"
  printf 'category\tname\tstatus\tdetail\n' >"$CHECKS_TSV"
}

add_check() {
  category="$1"
  name="$2"
  status="$3"
  detail="$4"
  case "$status" in
    passed|blocker|warning|skipped|unavailable) ;;
    *) fail "invalid check status: $status" ;;
  esac
  printf '%s\t%s\t%s\t%s\n' "$category" "$name" "$status" "$detail" >>"$CHECKS_TSV"
}

has_placeholder() {
  value="$1"
  python3 - "$value" <<'PY'
import sys
value = sys.argv[1].strip().lower()
needles = ("replace-with", "change-me", "dev-", "example", "placeholder")
raise SystemExit(0 if (not value or any(n in value for n in needles)) else 1)
PY
}

env_value() {
  key="$1"
  eval "printf '%s' \"\${$key-}\""
}

is_optional_env_key() {
  case "$1" in
    ADMIN_JWT_OLD_SECRETS|ADMIN_TOKEN|NOTIFY_WEBHOOK_URL|NOTIFY_EMAIL_TO|BACKUP_RETENTION_DAYS|CAPTURE_DATE_UTC|RESTORE_DRILL_DATABASE_URL|RESTORE_DRILL_BACKUP_FILE) return 0 ;;
    *) return 1 ;;
  esac
}

reference_env_keys() {
  sed -n 's/^\([A-Z][A-Z0-9_]*\)=.*/\1/p' docs/deployment/oci-reference-env.example
}

record_env_presence() {
  log "Check reference environment variable presence"
  mkdir -p "$OUT_REAL/env"
  env_summary="$OUT_REAL/env/reference-env.summary.tsv"
  printf 'key\tstatus\n' >"$env_summary"
  reference_env_keys | while IFS= read -r key
  do
    value="$(env_value "$key")"
    status="present"
    if [ -z "$value" ]; then
      if is_optional_env_key "$key"; then
        status="optional_empty"
      else
        status="missing"
      fi
    elif has_placeholder "$value"; then
      status="placeholder"
    fi
    printf '%s\t%s\n' "$key" "$status" >>"$env_summary"
    case "$status" in
      present|optional_empty) add_check "env" "$key" "passed" "$status" ;;
      placeholder) add_check "env" "$key" "blocker" "placeholder value configured" ;;
      missing) add_check "env" "$key" "blocker" "missing required reference env" ;;
    esac
  done
}

record_generated_secret_checks() {
  log "Check generated-secret presence"
  mkdir -p "$OUT_REAL/env"
  secret_summary="$OUT_REAL/env/generated-secrets.summary.tsv"
  printf 'key\tstatus\n' >"$secret_summary"
  for key in ADMIN_JWT_SECRET CSRF_SECRET DEVICE_TOKEN_PEPPER
  do
    value="$(env_value "$key")"
    if [ -z "$value" ]; then
      status="missing"
      check_status="blocker"
      detail="missing generated secret"
    elif has_placeholder "$value"; then
      status="placeholder"
      check_status="blocker"
      detail="placeholder generated secret"
    elif [ "${#value}" -lt 32 ]; then
      status="too_short"
      check_status="blocker"
      detail="generated secret shorter than 32 characters"
    else
      status="present"
      check_status="passed"
      detail="generated secret present"
    fi
    printf '%s\t%s\n' "$key" "$status" >>"$secret_summary"
    add_check "generated_secret" "$key" "$check_status" "$detail"
  done

  value="${ADMIN_TOKEN:-}"
  if [ -z "$value" ]; then
    status="skipped"
    check_status="skipped"
    detail="ADMIN_TOKEN optional; authenticated admin checks skipped unless requested"
  elif has_placeholder "$value"; then
    status="placeholder"
    check_status="warning"
    detail="optional ADMIN_TOKEN appears placeholder"
  elif [ "${#value}" -lt 32 ]; then
    status="too_short"
    check_status="warning"
    detail="optional ADMIN_TOKEN shorter than 32 characters"
  else
    status="present"
    check_status="passed"
    detail="optional ADMIN_TOKEN present"
  fi
  printf '%s\t%s\n' "ADMIN_TOKEN" "$status" >>"$secret_summary"
  add_check "generated_secret" "ADMIN_TOKEN" "$check_status" "$detail"
}

normalize_url() {
  url="$1"
  kind="$2"
  python3 - "$url" "$kind" <<'PY'
import ipaddress
import sys
from urllib.parse import urlsplit, urlunsplit

url, kind = sys.argv[1:3]
parts = urlsplit(url)
if parts.scheme not in ("http", "https"):
    raise SystemExit(f"{kind} must start with http:// or https://")
if not parts.hostname:
    raise SystemExit(f"{kind} must include a host")
if parts.username or parts.password or "@" in parts.netloc:
    raise SystemExit(f"{kind} must not include userinfo or embedded credentials")
if parts.query or parts.fragment:
    raise SystemExit(f"{kind} must not include query strings or fragments")
hostname = parts.hostname
host_lower = hostname.lower()
is_loopback = host_lower in ("localhost", "127.0.0.1", "::1")
if not is_loopback:
    try:
        is_loopback = ipaddress.ip_address(hostname).is_loopback
    except ValueError:
        is_loopback = False
if kind == "ADMIN_BASE_URL" and not is_loopback and parts.scheme != "https":
    raise SystemExit("ADMIN_BASE_URL must use HTTPS unless it is loopback")
path = parts.path.rstrip("/")
print(urlunsplit((parts.scheme, parts.netloc, path, "", "")))
print("true" if is_loopback else "false")
print(parts.scheme)
print(hostname)
PY
}

resolve_urls() {
  PUBLIC_BASE_NORMALIZED=""
  PUBLIC_BASE_LOOPBACK=""
  PUBLIC_BASE_SCHEME=""
  PUBLIC_BASE_HOST=""
  ADMIN_BASE_NORMALIZED=""
  ADMIN_URL_STATUS="skipped"
  if [ -z "${PUBLIC_BASE_URL:-}" ]; then
    add_check "url" "PUBLIC_BASE_URL" "blocker" "missing public base URL"
    return 0
  fi
  if public_info="$(normalize_url "$PUBLIC_BASE_URL" "PUBLIC_BASE_URL" 2>"$TMP_DIR/public-url-error.txt")"; then
    PUBLIC_BASE_NORMALIZED="$(printf '%s\n' "$public_info" | sed -n '1p')"
    PUBLIC_BASE_LOOPBACK="$(printf '%s\n' "$public_info" | sed -n '2p')"
    PUBLIC_BASE_SCHEME="$(printf '%s\n' "$public_info" | sed -n '3p')"
    PUBLIC_BASE_HOST="$(printf '%s\n' "$public_info" | sed -n '4p')"
    add_check "url" "PUBLIC_BASE_URL" "passed" "public URL normalized"
  else
    add_check "url" "PUBLIC_BASE_URL" "blocker" "$(cat "$TMP_DIR/public-url-error.txt")"
    return 0
  fi

  if [ -z "${ADMIN_BASE_URL:-}" ] && [ "$PUBLIC_BASE_LOOPBACK" = "true" ]; then
    ADMIN_BASE_URL="$PUBLIC_BASE_NORMALIZED"
  fi
  if [ -n "${ADMIN_BASE_URL:-}" ]; then
    if admin_info="$(normalize_url "$ADMIN_BASE_URL" "ADMIN_BASE_URL" 2>"$TMP_DIR/admin-url-error.txt")"; then
      ADMIN_BASE_NORMALIZED="$(printf '%s\n' "$admin_info" | sed -n '1p')"
      add_check "url" "ADMIN_BASE_URL" "passed" "admin URL normalized"
    else
      add_check "url" "ADMIN_BASE_URL" "blocker" "$(cat "$TMP_DIR/admin-url-error.txt")"
    fi
  else
    add_check "url" "ADMIN_BASE_URL" "skipped" "not supplied and public URL is not loopback"
  fi
}

path_url() {
  base="$1"
  path="$2"
  printf '%s%s\n' "$base" "$path"
}

record_public_feed() {
  label="$1"
  path="$2"
  mkdir -p "$OUT_REAL/public"
  tmp="$TMP_DIR/$label.body"
  headers="$TMP_DIR/$label.headers"
  meta="$TMP_DIR/$label.meta"
  url="$(path_url "$PUBLIC_BASE_NORMALIZED" "$path")"
  set +e
  curl -L -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    --max-filesize "$MAX_FEED_BYTES" \
    -D "$headers" \
    -o "$tmp" \
    -w 'http_code=%{http_code}\nurl_effective=%{url_effective}\nnum_redirects=%{num_redirects}\nsize_download=%{size_download}\ncontent_type=%{content_type}\n' \
    "$url" >"$meta" 2>"$TMP_DIR/$label.curl-stderr"
  rc=$?
  set -e
  status="$(sed -n 's/^http_code=//p' "$meta" | tail -n 1)"
  status="${status:-000}"
  bytes=0
  if [ -f "$tmp" ]; then
    bytes="$(wc -c <"$tmp" | awk '{print $1}')"
  fi
  sha=""
  outcome="unavailable"
  check_status="unavailable"
  if [ "$bytes" -gt "$MAX_FEED_BYTES" ] || [ "$rc" -eq 63 ]; then
    outcome="too_large"
    check_status="blocker"
  elif [ "$rc" -ne 0 ]; then
    outcome="curl_failed"
    check_status="unavailable"
  elif [ "$status" -lt 200 ] || [ "$status" -ge 300 ]; then
    outcome="bad_status"
    check_status="blocker"
  elif [ "$bytes" -le 0 ]; then
    outcome="empty"
    check_status="blocker"
  else
    sha="$(sha256_file "$tmp")"
    outcome="ok"
    check_status="passed"
  fi
  rm -f "$tmp"
  {
    printf 'label=%s\n' "$label"
    printf 'path=%s\n' "$path"
    printf 'status=%s\n' "$status"
    printf 'bytes=%s\n' "$bytes"
    printf 'sha256=%s\n' "$sha"
    printf 'outcome=%s\n' "$outcome"
    printf 'curl_exit=%s\n' "$rc"
    sed -n '/^url_effective=/p;/^num_redirects=/p;/^content_type=/p' "$meta"
  } >"$OUT_REAL/public/$label.summary"
  printf '%s,%s,%s,%s,%s\n' "$label" "$status" "$bytes" "$sha" "$outcome" >>"$OUT_REAL/public-summary.csv"
  add_check "public_feed_edge" "$label" "$check_status" "$outcome"
}

record_public_feeds() {
  log "Check anonymous public feed edge"
  printf 'label,status,bytes,sha256,outcome\n' >"$OUT_REAL/public-summary.csv"
  if [ -z "$PUBLIC_BASE_NORMALIZED" ]; then
    for label in feeds.json schedule.zip vehicle_positions.pb trip_updates.pb alerts.pb; do
      printf '%s,000,0,,skipped_no_public_base_url\n' "$label" >>"$OUT_REAL/public-summary.csv"
      add_check "public_feed_edge" "$label" "skipped" "PUBLIC_BASE_URL unavailable"
    done
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1; then
    for label in feeds.json schedule.zip vehicle_positions.pb trip_updates.pb alerts.pb; do
      printf '%s,000,0,,unavailable_curl_missing\n' "$label" >>"$OUT_REAL/public-summary.csv"
      add_check "public_feed_edge" "$label" "unavailable" "curl missing"
    done
    return 0
  fi
  record_public_feed "feeds.json" "/public/feeds.json"
  record_public_feed "schedule.zip" "/public/gtfs/schedule.zip"
  record_public_feed "vehicle_positions.pb" "/public/gtfsrt/vehicle_positions.pb"
  record_public_feed "trip_updates.pb" "/public/gtfsrt/trip_updates.pb"
  record_public_feed "alerts.pb" "/public/gtfsrt/alerts.pb"
}

http_status_head_then_get() {
  url="$1"
  method_used="$2"
  body="$TMP_DIR/route-body.tmp"
  status="$(curl -sS -I \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o /dev/null \
    -w '%{http_code}' \
    "$url" 2>/dev/null || true)"
  status="${status:-000}"
  if [ "$status" = "405" ]; then
    status="$(curl -sS \
      --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
      --max-time "$REQUEST_TIMEOUT_SECONDS" \
      -o "$body" \
      -w '%{http_code}' \
      "$url" 2>/dev/null || true)"
    status="${status:-000}"
    printf 'GET' >"$method_used"
  else
    printf 'HEAD' >"$method_used"
  fi
  rm -f "$body"
  printf '%s\n' "$status"
}

record_private_route_boundaries() {
  log "Check public/private route boundary"
  mkdir -p "$OUT_REAL/admin"
  summary="$OUT_REAL/admin/public-boundary.summary.tsv"
  printf 'path\tmethod\tstatus\tresult\n' >"$summary"
  if [ -z "$PUBLIC_BASE_NORMALIZED" ]; then
    add_check "admin_boundary" "private_routes" "skipped" "PUBLIC_BASE_URL unavailable"
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1; then
    add_check "admin_boundary" "private_routes" "unavailable" "curl missing"
    return 0
  fi
  for path in \
    "/admin/operations" \
    "/admin/operations/readiness" \
    "/admin/operations/validation-health" \
    "/admin/operations/validation-health.json" \
    "/admin/validation/run" \
    "/admin/devices/rebind" \
    "/admin/alerts/console" \
    "/admin/gtfs-studio" \
    "/v1/events" \
    "/metrics"
  do
    method_file="$TMP_DIR/private-route-method"
    status="$(http_status_head_then_get "$(path_url "$PUBLIC_BASE_NORMALIZED" "$path")" "$method_file")"
    method="$(cat "$method_file")"
    case "$status" in
      401|403|404) result="passed"; check_status="passed" ;;
      2*) result="blocker:unexpected_2xx"; check_status="blocker" ;;
      000) result="unavailable"; check_status="unavailable" ;;
      405) result="warning:method_limited"; check_status="warning" ;;
      *) result="warning:review_status_$status"; check_status="warning" ;;
    esac
    printf '%s\t%s\t%s\t%s\n' "$path" "$method" "$status" "$result" >>"$summary"
    add_check "admin_boundary" "$path" "$check_status" "$result"
  done
}

record_validation_health_summary() {
  log "Check private validator health summary when explicitly available"
  mkdir -p "$OUT_REAL/admin"
  if [ -z "${ADMIN_TOKEN:-}" ]; then
    add_check "admin_authenticated" "validation_health" "skipped" "ADMIN_TOKEN not supplied"
    printf '{"status":"skipped","reason":"no_admin_token"}\n' >"$OUT_REAL/admin/validation-health.summary.json"
    return 0
  fi
  if [ -z "$ADMIN_BASE_NORMALIZED" ]; then
    add_check "admin_authenticated" "validation_health" "skipped" "safe ADMIN_BASE_URL unavailable"
    printf '{"status":"skipped","reason":"safe_admin_base_url_unavailable"}\n' >"$OUT_REAL/admin/validation-health.summary.json"
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1; then
    add_check "admin_authenticated" "validation_health" "unavailable" "curl missing"
    printf '{"status":"unavailable","reason":"curl_missing"}\n' >"$OUT_REAL/admin/validation-health.summary.json"
    return 0
  fi
  body="$TMP_DIR/validation-health.body"
  status="$(curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$body" \
    -w '%{http_code}' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "$(path_url "$ADMIN_BASE_NORMALIZED" "/admin/operations/validation-health.json")" 2>/dev/null || true)"
  status="${status:-000}"
  python3 - "$body" "$OUT_REAL/admin/validation-health.summary.json" "$status" <<'PY'
import json
import pathlib
import sys
import sys
import sys

body, out, status = sys.argv[1:4]
summary = {"status": status}
try:
    data = json.loads(pathlib.Path(body).read_text())
except Exception:
    data = {}
for key in (
    "generated_at", "agency_id", "overall_status", "tooling_status",
    "external_evidence_created", "consumer_statuses_changed",
    "compliance_claimed", "production_readiness_claimed",
):
    if key in data:
        summary[key] = data[key]
feeds = []
for row in data.get("feeds", [])[:4]:
    feeds.append({
        "feed_type": row.get("feed_type", ""),
        "validator_id": row.get("validator_id", ""),
        "tooling_status": row.get("tooling_status", ""),
        "artifact_status": row.get("artifact_status", ""),
        "latest_result_status": row.get("latest_result_status", ""),
        "stale_status": row.get("stale_status", ""),
        "health_status": row.get("health_status", ""),
    })
summary["feeds"] = feeds
pathlib.Path(out).write_text(json.dumps(summary, indent=2) + "\n")
PY
  rm -f "$body"
  case "$status" in
    2*) add_check "admin_authenticated" "validation_health" "passed" "private validator health JSON returned 2xx" ;;
    000) add_check "admin_authenticated" "validation_health" "unavailable" "connection failed" ;;
    *) add_check "admin_authenticated" "validation_health" "warning" "unexpected authenticated status $status" ;;
  esac
}

record_authenticated_admin() {
  log "Check authenticated admin readiness when explicitly available"
  mkdir -p "$OUT_REAL/admin"
  if [ -z "${ADMIN_TOKEN:-}" ]; then
    add_check "admin_authenticated" "readiness" "skipped" "ADMIN_TOKEN not supplied"
    printf 'status=skipped\nreason=no_admin_token\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if [ -z "$ADMIN_BASE_NORMALIZED" ]; then
    add_check "admin_authenticated" "readiness" "blocker" "ADMIN_TOKEN supplied but safe ADMIN_BASE_URL unavailable"
    printf 'status=blocker\nreason=safe_admin_base_url_unavailable\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1; then
    add_check "admin_authenticated" "readiness" "unavailable" "curl missing"
    printf 'status=unavailable\nreason=curl_missing\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  body="$TMP_DIR/auth-readiness.body"
  status="$(curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$body" \
    -w '%{http_code}' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "$(path_url "$ADMIN_BASE_NORMALIZED" "/admin/operations/readiness")" 2>/dev/null || true)"
  status="${status:-000}"
  rows="$(python3 - "$body" <<'PY'
import re
import sys
try:
    text = open(sys.argv[1], encoding="utf-8", errors="replace").read(200000)
except OSError:
    text = ""
print(len(re.findall(r"<tr>", text, re.I)))
PY
)"
  rm -f "$body"
  printf 'status=%s\ntable_rows=%s\n' "$status" "$rows" >"$OUT_REAL/admin/auth-readiness.summary"
  case "$status" in
    2*) add_check "admin_authenticated" "readiness" "passed" "authenticated readiness returned 2xx" ;;
    000) add_check "admin_authenticated" "readiness" "unavailable" "connection failed" ;;
    *) add_check "admin_authenticated" "readiness" "warning" "unexpected authenticated status $status" ;;
  esac
}

record_https_posture() {
  log "Check HTTPS and redirect posture"
  mkdir -p "$OUT_REAL/https"
  if [ -z "$PUBLIC_BASE_NORMALIZED" ]; then
    add_check "https_posture" "public_root" "skipped" "PUBLIC_BASE_URL unavailable"
    return 0
  fi
  if [ "$PUBLIC_BASE_LOOPBACK" = "true" ]; then
    add_check "https_posture" "public_root" "skipped" "loopback local/reference root"
    return 0
  fi
  if [ "$PUBLIC_BASE_SCHEME" != "https" ]; then
    add_check "https_posture" "public_root" "blocker" "non-loopback public root is not HTTPS"
    return 0
  fi
  add_check "https_posture" "public_root" "passed" "non-loopback public root uses HTTPS"
  if command -v python3 >/dev/null 2>&1; then
    cert_meta="$OUT_REAL/https/certificate.summary"
    set +e
    python3 - "$PUBLIC_BASE_HOST" "$CONNECT_TIMEOUT_SECONDS" >"$cert_meta" 2>/dev/null <<'PY'
import socket
import ssl
import sys

host = sys.argv[1]
timeout = float(sys.argv[2])
context = ssl.create_default_context()
with socket.create_connection((host, 443), timeout=timeout) as sock:
    with context.wrap_socket(sock, server_hostname=host) as tls:
        cert = tls.getpeercert()

def name(parts):
    values = []
    for group in parts or []:
        for key, value in group:
            values.append(f"{key}={value}")
    return ", ".join(values)

print(f"subject={name(cert.get('subject'))}")
print(f"issuer={name(cert.get('issuer'))}")
print(f"notBefore={cert.get('notBefore', '')}")
print(f"notAfter={cert.get('notAfter', '')}")
PY
    rc=$?
    set -e
    if [ "$rc" -eq 0 ]; then
      add_check "https_posture" "tls_certificate" "passed" "certificate metadata recorded"
    else
      add_check "https_posture" "tls_certificate" "unavailable" "certificate metadata unavailable"
      rm -f "$cert_meta"
    fi
  else
    add_check "https_posture" "tls_certificate" "unavailable" "python3 missing"
  fi
  if command -v curl >/dev/null 2>&1; then
    http_url="http://$PUBLIC_BASE_HOST"
    meta="$TMP_DIR/http-redirect.meta"
    set +e
    curl -sS -L \
      --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
      --max-time "$REQUEST_TIMEOUT_SECONDS" \
      -o /dev/null \
      -w 'http_code=%{http_code}\nurl_effective=%{url_effective}\nnum_redirects=%{num_redirects}\n' \
      "$http_url" >"$meta" 2>/dev/null
    rc=$?
    set -e
    cp "$meta" "$OUT_REAL/https/http-redirect.summary"
    effective="$(sed -n 's/^url_effective=//p' "$meta" | tail -n 1)"
    redirects="$(sed -n 's/^num_redirects=//p' "$meta" | tail -n 1)"
    if [ "$rc" -ne 0 ]; then
      add_check "https_posture" "http_redirect" "unavailable" "HTTP redirect probe unavailable"
    elif [ "${redirects:-0}" -gt 0 ] && printf '%s' "$effective" | grep -qi '^https://'; then
      add_check "https_posture" "http_redirect" "passed" "HTTP redirects to HTTPS"
    else
      add_check "https_posture" "http_redirect" "warning" "HTTP-to-HTTPS redirect not confirmed"
    fi
  else
    add_check "https_posture" "http_redirect" "unavailable" "curl missing"
  fi
}

record_service_health() {
  log "Check local service health matrix"
  mkdir -p "$OUT_REAL/health"
  summary="$OUT_REAL/health/service-health.summary.tsv"
  printf 'service\tpath\tstatus\tresult\n' >"$summary"
  if ! command -v curl >/dev/null 2>&1; then
    add_check "service_health" "matrix" "unavailable" "curl missing"
    return 0
  fi
  for item in \
    "agency-config 127.0.0.1:8081" \
    "telemetry-ingest 127.0.0.1:8082" \
    "feed-vehicle-positions 127.0.0.1:8083" \
    "feed-trip-updates 127.0.0.1:8084" \
    "feed-alerts 127.0.0.1:8085" \
    "gtfs-studio 127.0.0.1:8086"
  do
    service="$(printf '%s' "$item" | awk '{print $1}')"
    hostport="$(printf '%s' "$item" | awk '{print $2}')"
    for path in /healthz /readyz
    do
      status="$(curl -sS --connect-timeout "$CONNECT_TIMEOUT_SECONDS" --max-time "$REQUEST_TIMEOUT_SECONDS" -o /dev/null -w '%{http_code}' "http://$hostport$path" 2>/dev/null || true)"
      status="${status:-000}"
      case "$status" in
        2*) result="passed"; check_status="passed" ;;
        000) result="unavailable"; check_status="unavailable" ;;
        *) result="warning:status_$status"; check_status="warning" ;;
      esac
      printf '%s\t%s\t%s\t%s\n' "$service" "$path" "$status" "$result" >>"$summary"
      add_check "service_health" "$service$path" "$check_status" "$result"
    done
  done
}

record_small_host_resources() {
  log "Check small-host resource posture"
  mkdir -p "$OUT_REAL/system"
  checks="$TMP_DIR/small-host-resources.checks.tsv"
  python3 - "$OUT_REAL/system/small-host-resources.summary.json" >"$checks" <<'PY'
import json
import os
import shutil
import subprocess
import sys

out = sys.argv[1]

def read_proc_meminfo():
    values = {}
    try:
        with open("/proc/meminfo", encoding="utf-8") as handle:
            for line in handle:
                key, raw = line.split(":", 1)
                parts = raw.strip().split()
                if parts and parts[0].isdigit():
                    values[key] = int(parts[0]) * 1024
    except OSError:
        return {}
    return values

def sysctl_int(name):
    try:
        raw = subprocess.check_output(["sysctl", "-n", name], text=True, stderr=subprocess.DEVNULL).strip()
        return int(raw)
    except Exception:
        return None

meminfo = read_proc_meminfo()
memory_bytes = meminfo.get("MemTotal") or sysctl_int("hw.memsize")
swap_bytes = meminfo.get("SwapTotal")
if swap_bytes is None:
    try:
        raw = subprocess.check_output(["sysctl", "-n", "vm.swapusage"], text=True, stderr=subprocess.DEVNULL)
        # macOS example: total = 2048.00M  used = ...
        token = raw.split("total = ", 1)[1].split()[0]
        swap_bytes = int(float(token.rstrip("M")) * 1024 * 1024)
    except Exception:
        swap_bytes = None

disk = shutil.disk_usage(".")
cpu_count = os.cpu_count() or 0
try:
    load_1m = os.getloadavg()[0]
except OSError:
    load_1m = None

memory_mb = None if memory_bytes is None else int(memory_bytes / 1024 / 1024)
swap_mb = None if swap_bytes is None else int(swap_bytes / 1024 / 1024)
disk_available_mb = int(disk.free / 1024 / 1024)

checks = []
def add(name, status, detail):
    checks.append(("small_host_resources", name, status, detail))

if memory_mb is None:
    add("memory", "unavailable", "total memory unavailable")
elif memory_mb < 1024:
    add("memory", "warning", f"memory={memory_mb}MB below 1024MB tiny-host floor")
elif memory_mb < 2048:
    add("memory", "warning", f"memory={memory_mb}MB; prefer off-host validators and low DB pools")
else:
    add("memory", "passed", f"memory={memory_mb}MB")

if disk_available_mb < 1024:
    add("disk_available", "blocker", f"available_disk={disk_available_mb}MB below 1024MB")
elif disk_available_mb < 5120:
    add("disk_available", "warning", f"available_disk={disk_available_mb}MB below 5120MB review threshold")
else:
    add("disk_available", "passed", f"available_disk={disk_available_mb}MB")

if cpu_count <= 0:
    add("cpu_count", "unavailable", "cpu count unavailable")
elif cpu_count < 2:
    add("cpu_count", "warning", "single CPU host; keep validators off host when possible")
else:
    add("cpu_count", "passed", f"cpu_count={cpu_count}")

if load_1m is None or cpu_count <= 0:
    add("load_average", "unavailable", "load average unavailable")
elif load_1m > cpu_count * 2:
    add("load_average", "warning", f"load_1m={load_1m:.2f} above 2x CPU count")
else:
    add("load_average", "passed", f"load_1m={load_1m:.2f}")

if swap_mb is None:
    add("swap", "unavailable", "swap total unavailable")
elif swap_mb == 0 and (memory_mb or 0) < 2048:
    add("swap", "warning", "no swap detected on small-memory host")
else:
    add("swap", "passed", f"swap={swap_mb}MB")

off_host_recommended = (memory_mb is not None and memory_mb < 2048) or cpu_count < 2
add(
    "validator_off_host_recommendation",
    "warning" if off_host_recommended else "passed",
    "off-host validator workflow recommended for this resource profile" if off_host_recommended else "host resources do not require off-host validators by this heuristic",
)

summary = {
    "memory_mb": memory_mb,
    "swap_mb": swap_mb,
    "disk_available_mb": disk_available_mb,
    "cpu_count": cpu_count,
    "load_1m": load_1m,
    "validator_off_host_recommended": off_host_recommended,
    "boundary": "Private local host resource diagnostic only. It does not prove capacity, uptime, SLA coverage, hosted availability, production readiness, or validator success.",
    "does_not_prove": "Does not prove production capacity, uptime, SLA coverage, hosted service availability, compliance, or consumer acceptance.",
}
with open(out, "w", encoding="utf-8") as handle:
    json.dump(summary, handle, indent=2, sort_keys=True)
    handle.write("\n")
for row in checks:
    print("\t".join(row))
PY
  while IFS="$(printf '\t')" read -r category name status detail
  do
    [ -n "$category" ] || continue
    add_check "$category" "$name" "$status" "$detail"
  done <"$checks"
}

record_service_dependency_review() {
  log "Check service dependency and proxy exposure posture"
  mkdir -p "$OUT_REAL/system"
  checks="$TMP_DIR/service-dependency.checks.tsv"
  python3 - "$OUT_REAL/system/service-dependency.summary.json" >"$checks" <<'PY'
import json
import pathlib
import sys

compose = pathlib.Path("deploy/docker-compose.yml")
local_caddy = pathlib.Path("deploy/Caddyfile.local")
oci_caddy = pathlib.Path("deploy/oci/Caddyfile")
systemd_dir = pathlib.Path("deploy/systemd")

checks = []
def add(category, name, status, detail):
    checks.append((category, name, status, detail))

compose_text = compose.read_text(encoding="utf-8") if compose.exists() else ""
local_text = local_caddy.read_text(encoding="utf-8") if local_caddy.exists() else ""
oci_text = oci_caddy.read_text(encoding="utf-8") if oci_caddy.exists() else ""

if compose.exists() and "depends_on:" in compose_text and "postgres:" in compose_text:
    add("service_dependencies", "compose_dependencies", "passed", "compose dependency graph present")
else:
    add("service_dependencies", "compose_dependencies", "warning", "compose dependency graph needs review")

if "respond \"not found\" 404" in local_text and "respond @local_root" in local_text:
    add("proxy_exposure", "local_caddy_fallback", "passed", "local proxy has explicit root and 404 fallback")
else:
    add("proxy_exposure", "local_caddy_fallback", "blocker", "local proxy fallback is not explicit")

unsafe_oci = any(token in oci_text for token in ("handle /admin", "handle /admin/", "handle /admin*", "handle /v1/events", "handle /admin/debug"))
if oci_caddy.exists() and not unsafe_oci and "handle /public/gtfsrt/vehicle_positions.pb" in oci_text and "respond 404" in oci_text:
    add("proxy_exposure", "oci_public_edge", "passed", "OCI public edge exposes feed paths and unmatched 404 only")
else:
    add("proxy_exposure", "oci_public_edge", "blocker", "OCI public edge may expose unsupported paths")

expected_units = [
    "open-transit-agency-config.service",
    "open-transit-telemetry-ingest.service",
    "open-transit-feed-vehicle-positions.service",
    "open-transit-feed-trip-updates.service",
    "open-transit-feed-alerts.service",
]
missing = [name for name in expected_units if not (systemd_dir / name).exists()]
if missing:
    add("service_dependencies", "systemd_units", "warning", f"missing_units={len(missing)}")
else:
    add("service_dependencies", "systemd_units", "passed", "expected systemd units present")

summary = {
    "compose_dependency_graph_present": "depends_on:" in compose_text,
    "local_proxy_explicit_404": "respond \"not found\" 404" in local_text,
    "oci_public_edge_admin_exposed": unsafe_oci,
    "systemd_expected_units": len(expected_units),
    "systemd_missing_units": len(missing),
    "boundary": "Static dependency and proxy review only. It does not start services, change firewall rules, contact the public edge, or prove production readiness.",
    "does_not_prove": "Does not prove deployment success, public availability, SLA, uptime, hosted service readiness, compliance, or consumer acceptance.",
}
with open(sys.argv[1], "w", encoding="utf-8") as handle:
    json.dump(summary, handle, indent=2, sort_keys=True)
    handle.write("\n")
for row in checks:
    print("\t".join(row))
PY
  while IFS="$(printf '\t')" read -r category name status detail
  do
    [ -n "$category" ] || continue
    add_check "$category" "$name" "$status" "$detail"
  done <"$checks"
}

sanitize_file() {
  src="$1"
  dst="$2"
  python3 - "$src" "$dst" <<'PY'
import re
import sys

src, dst = sys.argv[1:3]
text = open(src, encoding="utf-8", errors="replace").read()
patterns = [
    (re.compile(r"postgres(?:ql)?://[^:\s/]+:[^@\s]+@", re.I), "postgres://<redacted-user>:<redacted-password>@"),
    (re.compile(r"(?i)(DATABASE_URL\s*=\s*)[^\n\r]+"), r"\1<redacted>"),
    (re.compile(r"(?i)(RESTORE_DATABASE_URL\s*=\s*)[^\n\r]+"), r"\1<redacted>"),
    (re.compile(r"(?i)(RESTORE_DRILL_DATABASE_URL\s*=\s*)[^\n\r]+"), r"\1<redacted>"),
]
for pattern, repl in patterns:
    text = pattern.sub(repl, text)
open(dst, "w", encoding="utf-8").write(text)
PY
}

record_database_checks() {
  log "Check database, migrations, and PostGIS"
  mkdir -p "$OUT_REAL/database"
  if [ -z "${DATABASE_URL:-}" ]; then
    status="skipped"
    [ "$STRICT_DOCTOR" = "true" ] && status="blocker"
    add_check "database" "connectivity" "$status" "DATABASE_URL not supplied"
    add_check "migrations" "status" "$status" "DATABASE_URL not supplied"
    add_check "postgis" "availability" "$status" "DATABASE_URL not supplied"
    return 0
  fi
  if ! command -v go >/dev/null 2>&1; then
    add_check "database" "connectivity" "unavailable" "go missing"
    add_check "migrations" "status" "unavailable" "go missing"
  else
    raw="$TMP_DIR/migration-status.raw"
    set +e
    DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="${MIGRATIONS_DIR:-db/migrations}" go run ./cmd/migrate status >"$raw" 2>&1
    rc=$?
    set -e
    sanitize_file "$raw" "$OUT_REAL/database/migration-status.txt"
    if [ "$rc" -eq 0 ]; then
      add_check "database" "connectivity" "passed" "database reachable through migrator"
      add_check "migrations" "status" "passed" "migration status completed"
    else
      add_check "database" "connectivity" "blocker" "migrator could not reach database or status failed"
      add_check "migrations" "status" "blocker" "migration status exited nonzero"
    fi
    printf 'exit_code=%s\n' "$rc" >"$OUT_REAL/database/migration-status.summary"
  fi

  if command -v psql >/dev/null 2>&1; then
    raw="$TMP_DIR/postgis.raw"
    if python3 - "$CONNECT_TIMEOUT_SECONDS" "$raw" <<'PY'
import os
import subprocess
import sys
from urllib.parse import parse_qs, unquote, urlsplit

timeout, out_path = sys.argv[1:3]
url = os.environ.get("DATABASE_URL", "")
parts = urlsplit(url)
if parts.scheme not in ("postgres", "postgresql") or not parts.hostname:
    raise SystemExit(2)
env = os.environ.copy()
env["PGHOST"] = parts.hostname
if parts.port:
    env["PGPORT"] = str(parts.port)
if parts.username:
    env["PGUSER"] = unquote(parts.username)
if parts.password:
    env["PGPASSWORD"] = unquote(parts.password)
if parts.path and parts.path != "/":
    env["PGDATABASE"] = unquote(parts.path.lstrip("/"))
query = parse_qs(parts.query)
if query.get("sslmode"):
    env["PGSSLMODE"] = query["sslmode"][0]
env["PGCONNECT_TIMEOUT"] = timeout
with open(out_path, "w", encoding="utf-8") as out:
    proc = subprocess.run(
        ["psql", "-X", "-q", "-t", "-A", "-c", "SELECT postgis_full_version();"],
        env=env,
        stdout=out,
        stderr=subprocess.STDOUT,
        text=True,
        check=False,
    )
raise SystemExit(proc.returncode)
PY
    then
      rc=0
    else
      rc=$?
    fi
    sanitize_file "$raw" "$OUT_REAL/database/postgis.txt"
    if [ "$rc" -eq 0 ]; then
      add_check "postgis" "availability" "passed" "PostGIS query completed"
    elif [ "$rc" -eq 2 ]; then
      add_check "postgis" "availability" "skipped" "DATABASE_URL could not be parsed for safe psql env"
    else
      add_check "postgis" "availability" "unavailable" "PostGIS query unavailable"
    fi
    printf 'exit_code=%s\n' "$rc" >"$OUT_REAL/database/postgis.summary"
  else
    add_check "postgis" "availability" "skipped" "psql missing"
  fi
}

record_postgres_capacity_review() {
  log "Check Postgres small-host pool guidance"
  mkdir -p "$OUT_REAL/database"
  checks="$TMP_DIR/postgres-capacity.checks.tsv"
  python3 - "$OUT_REAL/database/postgres-capacity.summary.json" "${DB_MAX_CONNS:-}" >"$checks" <<'PY'
import json
import sys

out, raw_max = sys.argv[1:3]
default_pool = 10
service_count = 6
small_host_max_connections = 25
try:
    pool = int(raw_max) if raw_max else default_pool
except ValueError:
    pool = default_pool
configured = bool(raw_max and raw_max.isdigit())
estimated = pool * service_count
if not configured:
    status = "warning"
    detail = "DB_MAX_CONNS unset; default pool may exceed small-host Postgres max_connections=25"
elif estimated >= small_host_max_connections:
    status = "warning"
    detail = f"estimated_pooled_connections={estimated} may exceed small-host max_connections=25"
else:
    status = "passed"
    detail = f"estimated_pooled_connections={estimated} below small-host max_connections=25"
summary = {
    "db_max_conns_configured": configured,
    "per_service_pool": pool,
    "service_count_estimate": service_count,
    "estimated_pooled_connections": estimated,
    "small_host_postgres_max_connections": small_host_max_connections,
    "recommended_db_max_conns": 3,
    "boundary": "Static capacity guidance only. It does not inspect live connection counts, change database settings, or prove production capacity.",
    "does_not_prove": "Does not prove Postgres capacity, uptime, SLA coverage, production readiness, or data safety.",
}
with open(out, "w", encoding="utf-8") as handle:
    json.dump(summary, handle, indent=2, sort_keys=True)
    handle.write("\n")
print("\t".join(("postgres_capacity", "db_pool_budget", status, detail)))
PY
  while IFS="$(printf '\t')" read -r category name status detail
  do
    [ -n "$category" ] || continue
    add_check "$category" "$name" "$status" "$detail"
  done <"$checks"
}

record_upgrade_rollback_review() {
  log "Check upgrade and rollback checklist posture"
  mkdir -p "$OUT_REAL/operations"
  checks="$TMP_DIR/upgrade-rollback.checks.tsv"
  python3 - "$OUT_REAL/operations/upgrade-rollback.summary.json" >"$checks" <<'PY'
import json
import pathlib
import sys

docs = {
    "upgrade": pathlib.Path("docs/upgrade-and-rollback.md"),
    "backup_restore": pathlib.Path("docs/runbooks/backup-and-restore.md"),
    "off_host_validation": pathlib.Path("docs/deployment/off-host-validation.md"),
}
checks = []
for name, path in docs.items():
    if path.exists():
        checks.append(("upgrade_rollback", name, "passed", f"{name} guidance present"))
    else:
        checks.append(("upgrade_rollback", name, "blocker", f"{name} guidance missing"))
summary = {
    "required_review_steps": [
        "record current commit/version without secret values",
        "confirm backup target and restore-drill target before upgrade",
        "run migration status before and after upgrade",
        "run off-host validators when the host is too small for validator tooling",
        "check public feed fetches and private feed health after upgrade",
        "treat rollback as a restore/redeploy decision, not an automatic browser action",
    ],
    "browser_executes_upgrade_or_rollback": False,
    "boundary": "Checklist presence only. This does not execute backup, restore, migration, upgrade, rollback, or validation commands.",
    "does_not_prove": "Does not prove upgrade safety, rollback success, production readiness, compliance, consumer acceptance, SLA coverage, or uptime.",
}
with open(sys.argv[1], "w", encoding="utf-8") as handle:
    json.dump(summary, handle, indent=2, sort_keys=True)
    handle.write("\n")
for row in checks:
    print("\t".join(row))
PY
  while IFS="$(printf '\t')" read -r category name status detail
  do
    [ -n "$category" ] || continue
    add_check "$category" "$name" "$status" "$detail"
  done <"$checks"
}

record_validator_tooling() {
  log "Check pinned validator tooling"
  mkdir -p "$OUT_REAL/validators"
  set +e
  ./scripts/check-validators.sh >"$OUT_REAL/validators/tooling.txt" 2>&1
  rc=$?
  set -e
  case "$rc" in
    0) status="passed"; detail="pinned validator tooling passed" ;;
    11) status="blocker"; detail="pinned validator tooling missing" ;;
    12) status="blocker"; detail="pinned validator tooling misconfigured" ;;
    *) status="unavailable"; detail="validator tooling check exited $rc" ;;
  esac
  printf 'exit_code=%s\nstatus=%s\n' "$rc" "$status" >"$OUT_REAL/validators/tooling.summary"
  add_check "validators" "pinned_tooling" "$status" "$detail"
}

record_backup_restore_readiness() {
  log "Check backup and restore-drill readiness"
  mkdir -p "$OUT_REAL/operations"
  restore_database_url="${RESTORE_DATABASE_URL:-${RESTORE_DRILL_DATABASE_URL:-}}"
  restore_backup_file="${RESTORE_BACKUP_FILE:-${RESTORE_DRILL_BACKUP_FILE:-}}"
  backup_status="blocker"
  restore_url_status="blocker"
  restore_file_status="blocker"
  if [ -z "${BACKUP_DIR:-}" ]; then
    add_check "backup_readiness" "BACKUP_DIR" "blocker" "BACKUP_DIR not supplied"
  elif [ ! -d "$BACKUP_DIR" ]; then
    add_check "backup_readiness" "BACKUP_DIR" "blocker" "BACKUP_DIR does not exist"
  elif [ ! -r "$BACKUP_DIR" ]; then
    add_check "backup_readiness" "BACKUP_DIR" "blocker" "BACKUP_DIR is not readable"
  elif [ ! -w "$BACKUP_DIR" ]; then
    add_check "backup_readiness" "BACKUP_DIR" "warning" "BACKUP_DIR is not writable by current user"
    backup_status="warning"
  else
    add_check "backup_readiness" "BACKUP_DIR" "passed" "BACKUP_DIR exists and is readable/writable"
    backup_status="passed"
  fi

  if [ -z "$restore_database_url" ]; then
    add_check "restore_readiness" "RESTORE_DATABASE_URL" "blocker" "RESTORE_DATABASE_URL or RESTORE_DRILL_DATABASE_URL not supplied"
  else
    add_check "restore_readiness" "RESTORE_DATABASE_URL" "passed" "restore database URL present"
    restore_url_status="passed"
  fi
  if [ -z "$restore_backup_file" ]; then
    add_check "restore_readiness" "RESTORE_BACKUP_FILE" "blocker" "RESTORE_BACKUP_FILE or RESTORE_DRILL_BACKUP_FILE not supplied"
  elif [ ! -f "$restore_backup_file" ]; then
    add_check "restore_readiness" "RESTORE_BACKUP_FILE" "blocker" "RESTORE_BACKUP_FILE does not exist"
  elif [ ! -r "$restore_backup_file" ]; then
    add_check "restore_readiness" "RESTORE_BACKUP_FILE" "blocker" "RESTORE_BACKUP_FILE is not readable"
  else
    add_check "restore_readiness" "RESTORE_BACKUP_FILE" "passed" "RESTORE_BACKUP_FILE exists and is readable"
    restore_file_status="passed"
  fi
  python3 - "$OUT_REAL/operations/backup-restore-readiness.summary.json" "$backup_status" "$restore_url_status" "$restore_file_status" <<'PY'
import json
import sys

out, backup_status, restore_url_status, restore_file_status = sys.argv[1:5]
summary = {
    "backup_dir_status": backup_status,
    "restore_database_url_status": restore_url_status,
    "restore_backup_file_status": restore_file_status,
    "accepted_restore_url_env": ["RESTORE_DATABASE_URL", "RESTORE_DRILL_DATABASE_URL"],
    "accepted_restore_file_env": ["RESTORE_BACKUP_FILE", "RESTORE_DRILL_BACKUP_FILE"],
    "browser_executes_backup_or_restore": False,
    "boundary": "Presence/readability guidance only. The deployment doctor does not create backups or restore databases.",
    "does_not_prove": "Does not prove a backup exists, a restore drill succeeded, disaster recovery coverage, production readiness, SLA, uptime, compliance, or consumer acceptance.",
}
with open(out, "w", encoding="utf-8") as handle:
    json.dump(summary, handle, indent=2, sort_keys=True)
    handle.write("\n")
PY
}

record_release_identity() {
  log "Record release and git identity"
  mkdir -p "$OUT_REAL/system"
  {
    printf 'generated_at_utc=%s\n' "$TIMESTAMP"
    printf 'release_version_status='
    if [ -n "${OPEN_TRANSIT_RT_VERSION:-}" ] || [ -n "${RELEASE_VERSION:-}" ]; then
      printf 'present\n'
    else
      printf 'skipped\n'
    fi
    printf 'git_describe='
    git describe --tags --always --dirty 2>/dev/null || printf 'unavailable\n'
    printf 'git_sha='
    git rev-parse HEAD 2>/dev/null || printf 'unavailable\n'
    printf 'git_branch='
    git rev-parse --abbrev-ref HEAD 2>/dev/null || printf 'unavailable\n'
    if git diff --quiet --ignore-submodules -- 2>/dev/null; then
      printf 'working_tree=clean\n'
    else
      printf 'working_tree=dirty\n'
    fi
  } >"$OUT_REAL/system/release-identity.summary"
  add_check "release_identity" "git" "passed" "git identity recorded when available"
}

record_recent_private_diagnostics() {
  mkdir -p "$OUT_REAL/system"
  {
    printf 'operator_smoke_latest='
    find .cache/operator-smoke -mindepth 1 -maxdepth 1 -type d 2>/dev/null | sort | tail -n 1 || true
    printf 'support_bundle_latest='
    find .cache/support-bundles -mindepth 1 -maxdepth 1 -type d 2>/dev/null | sort | tail -n 1 || true
  } >"$OUT_REAL/system/recent-private-diagnostics.summary"
  add_check "private_diagnostics" "recent_outputs" "skipped" "recent private diagnostics linked only when present"
}

record_consumer_tracker_guard() {
  log "Check consumer tracker remains prepared-only"
  mkdir -p "$OUT_REAL/consumer-tracker"
  if python3 - <<'PY' >"$OUT_REAL/consumer-tracker/status.summary" 2>"$TMP_DIR/consumer-tracker.err"
import json

expected = {"Google Maps", "Apple Maps", "Transit App", "Bing Maps", "Moovit", "Mobility Database", "transit.land"}
with open("docs/evidence/consumer-submissions/status.json", encoding="utf-8") as f:
    data = json.load(f)
seen = {r["target"]: r["status"] for r in data["targets"]}
if set(seen) != expected:
    raise SystemExit(f"targets mismatch: {sorted(seen)}")
bad = {target: status for target, status in seen.items() if status != "prepared"}
if bad:
    raise SystemExit(f"non-prepared statuses: {bad}")
print("status=passed")
print("target_count=7")
print("all_statuses=prepared")
PY
  then
    add_check "consumer_tracker" "prepared_only" "passed" "seven targets remain prepared"
  else
    cp "$TMP_DIR/consumer-tracker.err" "$OUT_REAL/consumer-tracker/status.error"
    add_check "consumer_tracker" "prepared_only" "blocker" "consumer tracker mismatch"
  fi
}

json_string() {
  python3 -c 'import json,sys; print(json.dumps(sys.argv[1]))' "$1"
}

write_reports() {
  log "Write structured deployment doctor reports"
  python3 - "$OUT_REAL" "$TIMESTAMP" "$STRICT_DOCTOR" <<'PY'
import csv
import json
import pathlib
import sys

out = pathlib.Path(sys.argv[1])
timestamp = sys.argv[2]
strict = sys.argv[3] == "true"
checks = []
with (out / "checks.tsv").open(encoding="utf-8", newline="") as f:
    for row in csv.DictReader(f, delimiter="\t"):
        checks.append(row)
counts = {status: 0 for status in ("passed", "blocker", "warning", "skipped", "unavailable")}
for row in checks:
    counts[row["status"]] = counts.get(row["status"], 0) + 1
if counts["blocker"]:
    overall = "blocker"
elif counts["warning"]:
    overall = "warning"
elif counts["unavailable"]:
    overall = "unavailable"
else:
    overall = "passed"

def category_status(category):
    rows = [r for r in checks if r["category"] == category]
    if not rows:
        return "skipped"
    local_counts = {status: 0 for status in counts}
    for row in rows:
        local_counts[row["status"]] = local_counts.get(row["status"], 0) + 1
    if local_counts["blocker"]:
        return "blocker"
    if local_counts["warning"]:
        return "warning"
    if local_counts["unavailable"]:
        return "unavailable"
    if local_counts["passed"]:
        return "passed"
    return "skipped"

flags = {
    "external_evidence_created": False,
    "final_root_evidence_created": False,
    "consumer_statuses_changed": False,
    "compliance_claimed": False,
    "production_readiness_claimed": False,
    "hosted_saas_claimed": False,
    "sla_claimed": False,
    "uptime_guarantee_claimed": False,
    "vendor_compatibility_claimed": False,
    "hardware_certification_claimed": False,
    "production_grade_eta_claimed": False,
}
summary = {
    "generated_at_utc": timestamp,
    "overall_status": overall,
    "strict_doctor": strict,
    "counts": counts,
    "categories": {
        "public_feed_edge": category_status("public_feed_edge"),
        "admin_boundary": category_status("admin_boundary"),
        "database": category_status("database"),
        "migrations": category_status("migrations"),
        "postgis": category_status("postgis"),
        "validators": category_status("validators"),
        "backup_readiness": category_status("backup_readiness"),
        "restore_readiness": category_status("restore_readiness"),
        "small_host_resources": category_status("small_host_resources"),
        "service_dependencies": category_status("service_dependencies"),
        "proxy_exposure": category_status("proxy_exposure"),
        "postgres_capacity": category_status("postgres_capacity"),
        "upgrade_rollback": category_status("upgrade_rollback"),
    },
    **flags,
    "checks": checks,
}
manifest = {
    "generated_at_utc": timestamp,
    "output_kind": "private_operator_diagnostics",
    "included": [
        "environment key presence statuses without values",
        "generated-secret presence/placeholder/length statuses without values",
        "public feed HTTP metadata and checksums without retained bodies",
        "public/private route boundary statuses",
        "optional authenticated admin readiness summary without token values",
        "service health status matrix",
        "small-host memory, CPU, load, disk, swap, and off-host validator guidance",
        "static service dependency and proxy exposure posture",
        "read-only migration and PostGIS summaries when available",
        "static Postgres pool budget guidance",
        "validator tooling status",
        "backup and restore-drill readiness statuses",
        "upgrade and rollback checklist posture",
        "git/release identity",
        "consumer tracker prepared-only guard",
    ],
    "excluded": [
        "raw environment files",
        "secret values",
        "Authorization headers",
        "cookies",
        "database URLs",
        "webhook values",
        "raw public feed bodies",
        "database dumps",
        "backup file contents",
        "live resource reservations",
        "service start/stop actions",
        "migration execution",
        "backup or restore execution",
        "private keys",
        "consumer submissions",
        "evidence packets",
    ],
    **flags,
}
for name, obj in (("summary.json", summary), ("manifest.json", manifest)):
    with (out / name).open("w", encoding="utf-8") as f:
        json.dump(obj, f, indent=2, sort_keys=True)
        f.write("\n")
with (out / "summary.md").open("w", encoding="utf-8") as f:
    f.write("# Deployment Doctor Summary\n\n")
    f.write(f"- Overall status: `{overall}`\n")
    for key in ("passed", "blocker", "warning", "skipped", "unavailable"):
        f.write(f"- {key}: `{counts[key]}`\n")
    for key, value in flags.items():
        f.write(f"- {key}: `{str(value).lower()}`\n")
    f.write("\n## Category Status\n\n")
    for key, value in summary["categories"].items():
        f.write(f"- {key}: `{value}`\n")
with (out / "manifest.md").open("w", encoding="utf-8") as f:
    f.write("# Deployment Doctor Manifest\n\n")
    f.write("- Output kind: `private_operator_diagnostics`\n")
    f.write("- This is not an evidence packet, final-root proof, compliance claim, or production-readiness claim.\n")
    for key, value in flags.items():
        f.write(f"- {key}: `{str(value).lower()}`\n")
PY
  python3 -m json.tool "$OUT_REAL/summary.json" >/dev/null
  python3 -m json.tool "$OUT_REAL/manifest.json" >/dev/null
}

redaction_scan() {
  log "Run generated output redaction scan"
  python3 - "$OUT_REAL" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
patterns = [
    ("authorization_bearer", re.compile(r"Authorization:\s*Bearer\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("jwt_like_value", re.compile(r"\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b")),
    ("postgres_password_url", re.compile(r"postgres(?:ql)?://[^:\s/]+:[^@\s]+@", re.I)),
    ("cookie_header", re.compile(r"^Cookie:\s*\S+", re.I | re.M)),
    ("private_key", re.compile(r"BEGIN [A-Z ]*PRIVATE KEY")),
    ("env_secret_assignment", re.compile(r"^(?:[A-Z0-9_]*(?:TOKEN|SECRET|PASSWORD|PRIVATE_KEY)[A-Z0-9_]*)\s*=\s*\S+", re.I | re.M)),
    ("webhook_secret_url", re.compile(r"https://[^\s]*(?:hooks|webhook|discord|slack)[^\s]*(?:token|secret|key|/services/)[^\s]*", re.I)),
]
findings = []
for path in root.rglob("*"):
    if not path.is_file():
        continue
    text = path.read_text(encoding="utf-8", errors="replace")
    for name, pattern in patterns:
        if pattern.search(text):
            findings.append(f"{path.relative_to(root)}:{name}")
if findings:
    print("redaction scan failed:")
    for finding in findings:
        print(f"  {finding}")
    raise SystemExit(1)
PY
}

copy_paste_summary() {
  python3 - "$OUT_REAL" <<'PY'
import json
import sys

out = sys.argv[1]
data = json.load(open(f"{out}/summary.json", encoding="utf-8"))
counts = data["counts"]
cat = data["categories"]
print()
print("Deployment doctor copy/paste summary:")
print(f"  output_dir={out}")
print(f"  blocker_count={counts['blocker']}")
print(f"  warning_count={counts['warning']}")
print(f"  unavailable_count={counts['unavailable']}")
for key in ("public_feed_edge", "admin_boundary", "database", "migrations", "postgis", "validators", "backup_readiness", "restore_readiness", "small_host_resources", "service_dependencies", "proxy_exposure", "postgres_capacity", "upgrade_rollback"):
    print(f"  {key}={cat[key]}")
for key in ("external_evidence_created", "final_root_evidence_created", "consumer_statuses_changed", "compliance_claimed", "production_readiness_claimed", "hosted_saas_claimed", "sla_claimed", "uptime_guarantee_claimed"):
    print(f"  {key}={str(data[key]).lower()}")
PY
}

parse_env() {
  bool_var STRICT_DOCTOR "$STRICT_DOCTOR"
  bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
  bool_var FORCE "$FORCE"
  positive_int CONNECT_TIMEOUT_SECONDS "$CONNECT_TIMEOUT_SECONDS"
  positive_int REQUEST_TIMEOUT_SECONDS "$REQUEST_TIMEOUT_SECONDS"
  positive_int MAX_FEED_BYTES "$MAX_FEED_BYTES"
}

main() {
  case "${1:-}" in
    -h|--help|help) usage; exit 0 ;;
    "") ;;
    *) usage; fail "unknown argument: $1" ;;
  esac

  command -v python3 >/dev/null 2>&1 || fail "python3 is required"
  parse_env
  prepare_output_dir

  record_env_presence
  record_generated_secret_checks
  resolve_urls
  record_public_feeds
  record_private_route_boundaries
  record_authenticated_admin
  record_validation_health_summary
  record_https_posture
  record_service_health
  record_small_host_resources
  record_service_dependency_review
  record_database_checks
  record_postgres_capacity_review
  record_validator_tooling
  record_backup_restore_readiness
  record_upgrade_rollback_review
  record_release_identity
  record_recent_private_diagnostics
  record_consumer_tracker_guard
  write_reports
  redaction_scan
  copy_paste_summary

  if [ "$STRICT_DOCTOR" = "true" ]; then
    blockers="$(python3 - "$OUT_REAL/summary.json" <<'PY'
import json
import sys
print(json.load(open(sys.argv[1], encoding="utf-8"))["counts"]["blocker"])
PY
)"
    if [ "$blockers" -gt 0 ]; then
      fail "deployment doctor found $blockers blocker(s) in STRICT_DOCTOR=true mode"
    fi
  fi
}

main "$@"
