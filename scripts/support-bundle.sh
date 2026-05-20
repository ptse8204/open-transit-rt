#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-http://localhost:8080}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
INCLUDE_ADMIN_READINESS="${INCLUDE_ADMIN_READINESS:-false}"
CONNECT_TIMEOUT_SECONDS="${CONNECT_TIMEOUT_SECONDS:-5}"
REQUEST_TIMEOUT_SECONDS="${REQUEST_TIMEOUT_SECONDS:-30}"
MAX_FEED_BYTES="${MAX_FEED_BYTES:-104857600}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
FORCE="${FORCE:-false}"
TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/support-bundles/$TIMESTAMP}"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/support-bundle.sh
  scripts/support-bundle.sh --self-test-redaction

Environment:
  PUBLIC_BASE_URL                 Public feed root, default http://localhost:8080
  ADMIN_BASE_URL                  Optional admin root; defaults to PUBLIC_BASE_URL only when loopback
  ADMIN_TOKEN                     Optional admin bearer token
  INCLUDE_ADMIN_READINESS         true|false; include authenticated readiness summary when ADMIN_TOKEN is set
  DATABASE_URL                    Optional; migration status is sanitized before writing
  OUTPUT_DIR                      Default .cache/support-bundles/<timestamp>
  ALLOW_UNIGNORED_OUTPUT_DIR      true|false; allow OUTPUT_DIR outside repo .cache
  FORCE                           true|false; allow non-empty OUTPUT_DIR reuse
  CONNECT_TIMEOUT_SECONDS         curl connect timeout, default 5
  REQUEST_TIMEOUT_SECONDS         curl total timeout, default 30
  MAX_FEED_BYTES                  maximum public feed size to checksum, default 104857600

Safety:
  This helper records redaction-safe diagnostics only.
  Cookie auth is not supported.
  Tokens, Authorization headers, cookies, JWTs, CSRF values, raw .env files, and DB passwords are never written.
  The self-test mode exercises the sanitizer and generated-bundle scan without contacting services.
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

need() {
  command -v "$1" >/dev/null 2>&1 || return 1
}

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
PY
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

support_bundle_ref() {
  python3 - "$ROOT_DIR" "${OUT_REAL:-}" "$OUTPUT_DIR" <<'PY'
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out_arg = sys.argv[2] or sys.argv[3]
out = pathlib.Path(out_arg)
if not out.is_absolute():
    out = root / out
out = out.resolve(strict=False)
cache = (root / ".cache" / "support-bundles").resolve(strict=False)
try:
    rel = out.relative_to(cache)
    print(f".cache/support-bundles/{rel.as_posix()}")
except ValueError:
    print(f"<support-bundle-output>/{out.name}")
PY
}

url_path_ref() {
  python3 - "$1" <<'PY'
import sys
from urllib.parse import urlsplit

raw = sys.argv[1].strip()
if not raw:
    print("not_recorded")
    raise SystemExit
parts = urlsplit(raw)
if parts.scheme not in ("http", "https") or not parts.netloc:
    print("<redacted-url>")
    raise SystemExit
path = parts.path or "/"
if parts.query:
    path = f"{path}?<redacted-query>"
print(path)
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
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-support-bundle.XXXXXX")"
}

path_url() {
  base="$1"
  path="$2"
  printf '%s%s\n' "$base" "$path"
}

resolve_urls() {
  public_info="$(normalize_url "$PUBLIC_BASE_URL" "PUBLIC_BASE_URL")" || fail "$public_info"
  PUBLIC_BASE_NORMALIZED="$(printf '%s\n' "$public_info" | sed -n '1p')"
  PUBLIC_BASE_LOOPBACK="$(printf '%s\n' "$public_info" | sed -n '2p')"
  if [ -z "$ADMIN_BASE_URL" ] && [ "$PUBLIC_BASE_LOOPBACK" = "true" ]; then
    ADMIN_BASE_URL="$PUBLIC_BASE_NORMALIZED"
  fi
  if [ -n "$ADMIN_BASE_URL" ]; then
    if admin_info="$(normalize_url "$ADMIN_BASE_URL" "ADMIN_BASE_URL" 2>"$TMP_DIR/admin-url-error.txt")"; then
      ADMIN_BASE_NORMALIZED="$(printf '%s\n' "$admin_info" | sed -n '1p')"
      ADMIN_URL_BLOCKER=""
    else
      ADMIN_BASE_NORMALIZED=""
      ADMIN_URL_BLOCKER="$(cat "$TMP_DIR/admin-url-error.txt")"
    fi
  else
    ADMIN_BASE_NORMALIZED=""
    ADMIN_URL_BLOCKER="ADMIN_BASE_URL was not explicitly supplied and PUBLIC_BASE_URL is not loopback"
  fi
  if [ -n "$ADMIN_TOKEN" ] && [ -z "$ADMIN_BASE_NORMALIZED" ]; then
    ADMIN_TOKEN_BLOCKER="ADMIN_TOKEN was supplied but safe ADMIN_BASE_URL is not available"
  else
    ADMIN_TOKEN_BLOCKER=""
  fi
}

record_versions() {
  log "Record command and runtime versions"
  mkdir -p "$OUT_REAL/system"
  {
    printf 'generated_at_utc=%s\n' "$TIMESTAMP"
    printf 'uname='
    uname -srm 2>/dev/null || printf 'unavailable\n'
    printf 'go='
    if command -v go >/dev/null 2>&1; then go version; else printf 'unavailable\n'; fi
    printf 'git='
    if command -v git >/dev/null 2>&1; then git --version; else printf 'unavailable\n'; fi
    printf 'curl='
    if command -v curl >/dev/null 2>&1; then curl --version | sed -n '1p'; else printf 'unavailable\n'; fi
    printf 'docker='
    if command -v docker >/dev/null 2>&1; then docker --version; else printf 'unavailable\n'; fi
    printf 'python3='
    if command -v python3 >/dev/null 2>&1; then python3 --version; else printf 'unavailable\n'; fi
  } >"$OUT_REAL/system/versions.txt"
  if command -v git >/dev/null 2>&1; then
    {
      printf 'commit='
      git rev-parse HEAD 2>/dev/null || printf 'unavailable\n'
      if git diff --quiet --ignore-submodules -- 2>/dev/null; then
        printf 'working_tree=clean\n'
      else
        printf 'working_tree=dirty\n'
      fi
    } >"$OUT_REAL/system/git-summary.txt"
  fi
}

record_public_probe() {
  label="$1"
  path="$2"
  mkdir -p "$OUT_REAL/public"
  tmp="$TMP_DIR/$label.tmp"
  meta="$TMP_DIR/$label.meta"
  headers="$TMP_DIR/$label.headers"
  url="$(path_url "$PUBLIC_BASE_NORMALIZED" "$path")"
  curl_exit=0
  set +e
  curl -L -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    --max-filesize "$MAX_FEED_BYTES" \
    -D "$headers" \
    -o "$tmp" \
    -w 'http_code=%{http_code}\nurl_effective=%{url_effective}\nnum_redirects=%{num_redirects}\nsize_download=%{size_download}\ncontent_type=%{content_type}\n' \
    "$url" >"$meta" 2>/dev/null
  curl_exit=$?
  set -e
  status="$(sed -n 's/^http_code=//p' "$meta" | tail -n 1)"
  status="${status:-000}"
  bytes=0
  if [ -f "$tmp" ]; then
    bytes="$(wc -c <"$tmp" | awk '{print $1}')"
  fi
  sha=""
  outcome="unavailable"
  if [ "$bytes" -gt "$MAX_FEED_BYTES" ] || [ "$curl_exit" -eq 63 ]; then
    outcome="too_large"
  elif [ "$curl_exit" -ne 0 ]; then
    outcome="curl_failed"
  elif [ "$status" -lt 200 ] || [ "$status" -ge 300 ]; then
    outcome="bad_status"
  elif [ "$bytes" -le 0 ]; then
    outcome="empty"
  else
    sha="$(sha256_file "$tmp")"
    outcome="ok"
  fi
  printf '%s,%s,%s,%s,%s\n' "$label" "$status" "$bytes" "$sha" "$outcome" >>"$PUBLIC_SUMMARY_CSV"
  {
    printf 'label=%s\n' "$label"
    printf 'path=%s\n' "$path"
    printf 'status=%s\n' "$status"
    printf 'bytes=%s\n' "$bytes"
    printf 'sha256=%s\n' "$sha"
    printf 'outcome=%s\n' "$outcome"
    printf 'curl_exit=%s\n' "$curl_exit"
    effective="$(sed -n 's/^url_effective=//p' "$meta" | tail -n 1)"
    printf 'url_effective_path=%s\n' "$(url_path_ref "$effective")"
    sed -n '/^num_redirects=/p;/^content_type=/p' "$meta"
  } >"$OUT_REAL/public/$label.summary"
  rm -f "$tmp" "$meta" "$headers"
}

record_public_feeds() {
  log "Record public feed status summaries"
  PUBLIC_SUMMARY_CSV="$OUT_REAL/public-summary.csv"
  printf 'label,status,bytes,sha256,outcome\n' >"$PUBLIC_SUMMARY_CSV"
  if ! command -v curl >/dev/null 2>&1; then
    for label in feeds.json schedule.zip vehicle_positions.pb trip_updates.pb alerts.pb; do
      printf '%s,000,0,,skipped_curl_missing\n' "$label" >>"$PUBLIC_SUMMARY_CSV"
    done
    return 0
  fi
  record_public_probe "feeds.json" "/public/feeds.json"
  record_public_probe "schedule.zip" "/public/gtfs/schedule.zip"
  record_public_probe "vehicle_positions.pb" "/public/gtfsrt/vehicle_positions.pb"
  record_public_probe "trip_updates.pb" "/public/gtfsrt/trip_updates.pb"
  record_public_probe "alerts.pb" "/public/gtfsrt/alerts.pb"
}

http_status_no_redirect() {
  url="$1"
  body="$TMP_DIR/admin-body.tmp"
  if ! command -v curl >/dev/null 2>&1; then
    printf '000'
    return 0
  fi
  curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$body" \
    -w '%{http_code}' \
    "$url" 2>/dev/null || true
}

record_admin_boundary() {
  log "Record unauthenticated admin boundary status"
  mkdir -p "$OUT_REAL/admin"
  status="$(http_status_no_redirect "$(path_url "$PUBLIC_BASE_NORMALIZED" "/admin/operations/readiness")")"
  status="${status:-000}"
  case "$status" in
    401|403|404) ADMIN_BOUNDARY_RESULT="passed:$status" ;;
    2*) ADMIN_BOUNDARY_RESULT="blocker:unexpected_2xx_$status" ;;
    000) ADMIN_BOUNDARY_RESULT="unavailable" ;;
    *) ADMIN_BOUNDARY_RESULT="review:$status" ;;
  esac
  printf 'status=%s\nresult=%s\n' "$status" "$ADMIN_BOUNDARY_RESULT" >"$OUT_REAL/admin/public-readiness.summary"
  rm -f "$TMP_DIR/admin-body.tmp"
}

record_authenticated_readiness() {
  mkdir -p "$OUT_REAL/admin"
  if [ -n "$ADMIN_TOKEN_BLOCKER" ]; then
    AUTH_READINESS_STATUS="blocker:unsafe_or_missing_admin_base_url"
    printf 'status=blocker\nreason=%s\n' "$ADMIN_TOKEN_BLOCKER" >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if [ "$INCLUDE_ADMIN_READINESS" != "true" ]; then
    AUTH_READINESS_STATUS="skipped:INCLUDE_ADMIN_READINESS_not_true"
    printf 'status=skipped\nreason=INCLUDE_ADMIN_READINESS_not_true\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if [ -z "$ADMIN_TOKEN" ]; then
    AUTH_READINESS_STATUS="skipped:no_admin_token"
    printf 'status=skipped\nreason=no_admin_token\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if ! command -v curl >/dev/null 2>&1; then
    AUTH_READINESS_STATUS="unavailable:curl_missing"
    printf 'status=unavailable\nreason=curl_missing\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  log "Record authenticated readiness summary"
  body="$TMP_DIR/auth-readiness.html"
  status="$(curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$body" \
    -w '%{http_code}' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "$(path_url "$ADMIN_BASE_NORMALIZED" "/admin/operations/readiness")" 2>/dev/null || true)"
  status="${status:-000}"
  python3 - "$body" "$status" "$OUT_REAL/admin/auth-readiness.summary" <<'PY'
import re
import sys
body, status, out = sys.argv[1:4]
text = ""
try:
    text = open(body, encoding="utf-8", errors="replace").read(200000)
except OSError:
    pass
title = ""
m = re.search(r"<h2>(.*?)</h2>", text, re.I | re.S)
if m:
    title = re.sub(r"\s+", " ", m.group(1)).strip()
rows = len(re.findall(r"<tr>", text, re.I))
with open(out, "w", encoding="utf-8") as f:
    f.write(f"status={status}\n")
    f.write(f"title={title}\n")
    f.write(f"table_rows={rows}\n")
PY
  case "$status" in
    2*) AUTH_READINESS_STATUS="passed:$status" ;;
    *) AUTH_READINESS_STATUS="unavailable:$status" ;;
  esac
  rm -f "$body"
}

record_validator_tooling() {
  log "Record validator tooling status"
  mkdir -p "$OUT_REAL/validators"
  raw="$TMP_DIR/validator-tooling.raw"
  set +e
  ./scripts/check-validators.sh >"$raw" 2>&1
  rc=$?
  set -e
  sanitize_file "$raw" "$OUT_REAL/validators/tooling.txt"
  case "$rc" in
    0) VALIDATOR_TOOLING_STATUS="passed" ;;
    11) VALIDATOR_TOOLING_STATUS="missing" ;;
    12) VALIDATOR_TOOLING_STATUS="misconfigured" ;;
    *) VALIDATOR_TOOLING_STATUS="failed_exit_$rc" ;;
  esac
  printf 'exit_code=%s\nstatus=%s\n' "$rc" "$VALIDATOR_TOOLING_STATUS" >"$OUT_REAL/validators/tooling.summary"
  VALIDATION_API_STATUS="not_run:support_bundle_does_not_store_full_validation_reports"
  printf 'status=not_run\nreason=support_bundle_does_not_store_full_validation_reports\n' >"$OUT_REAL/validators/api.summary"
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
    (re.compile(r"(?i)(DATABASE_URL\s*[:=]\s*)[^\n\r]+"), r"\1<redacted>"),
    (re.compile(r"(?i)(\"(?:token|secret|password|authorization|cookie|api_key|webhook_url|database_url|private_key)\"\s*:\s*\")[^\"]+"), r"\1<redacted>"),
    (re.compile(r"(?i)((?:[A-Z0-9_]*)(?:TOKEN|SECRET|PASSWORD|PRIVATE_KEY|WEBHOOK_URL|API_KEY|DATABASE_URL|DEVICE_TOKEN|ADMIN_TOKEN)(?:[A-Z0-9_]*)\s*[:=]\s*)[^\n\r,}]+"), r"\1<redacted>"),
    (re.compile(r"(?i)(Authorization|Proxy-Authorization):\s*(?:Bearer|Basic)\s+[^\n\r]+"), r"\1: <redacted>"),
    (re.compile(r"(?i)(Cookie|Set-Cookie|X-Api-Key|Api-Key|X-Webhook-Signature):\s*[^\n\r]+"), r"\1: <redacted>"),
    (re.compile(r"(?i)\bBearer\s+[A-Za-z0-9._~+/=-]{8,}"), "Bearer <redacted>"),
    (re.compile(r"(?i)\bBasic\s+[A-Za-z0-9._~+/=-]{8,}"), "Basic <redacted>"),
    (re.compile(r"(?i)(?:postgres|postgresql|mysql|mariadb|mongodb|redis)://[^\s\"'<>]+"), "<redacted-db-url>"),
    (re.compile(r"(?i)https?://[^/\s:@]+:[^@\s]+@[^\s\"'<>]+"), "https://<redacted-user>:<redacted-password>@<redacted-host>"),
    (re.compile(r"(?i)https://[^\s]*(?:hooks|webhook|discord|slack)[^\s]*(?:token|secret|key|/services/)[^\s]*"), "<redacted-webhook-url>"),
    (re.compile(r"(?:(?<=\s)|(?<==)|^)(?:/Users/|/home/|/private/|/var/lib/|/etc/|[A-Za-z]:\\\\Users\\\\)[^\s\"']+"), "<redacted-private-path>"),
]
for pattern, repl in patterns:
    text = pattern.sub(repl, text)
open(dst, "w", encoding="utf-8").write(text)
PY
}

record_migration_status() {
  mkdir -p "$OUT_REAL/database"
  if [ -z "${DATABASE_URL:-}" ]; then
    MIGRATION_STATUS="skipped:no_database_url"
    printf 'status=skipped\nreason=no_database_url\n' >"$OUT_REAL/database/migration-status.summary"
    return 0
  fi
  if ! command -v go >/dev/null 2>&1; then
    MIGRATION_STATUS="unavailable:go_missing"
    printf 'status=unavailable\nreason=go_missing\n' >"$OUT_REAL/database/migration-status.summary"
    return 0
  fi
  log "Record sanitized migration status"
  raw="$TMP_DIR/migration-status.raw"
  set +e
  DATABASE_URL="$DATABASE_URL" MIGRATIONS_DIR="${MIGRATIONS_DIR:-db/migrations}" go run ./cmd/migrate status >"$raw" 2>&1
  rc=$?
  set -e
  sanitize_file "$raw" "$OUT_REAL/database/migration-status.txt"
  if [ "$rc" -eq 0 ]; then
    MIGRATION_STATUS="available"
  else
    MIGRATION_STATUS="unavailable:exit_$rc"
  fi
  printf 'status=%s\nexit_code=%s\n' "$MIGRATION_STATUS" "$rc" >"$OUT_REAL/database/migration-status.summary"
}

record_local_health() {
  log "Record local health endpoint summaries"
  mkdir -p "$OUT_REAL/health"
  if ! command -v curl >/dev/null 2>&1; then
    printf 'curl missing; health checks skipped\n' >"$OUT_REAL/health/summary.txt"
    return 0
  fi
  : >"$OUT_REAL/health/summary.txt"
  for item in \
    "agency-config http://127.0.0.1:8081/healthz" \
    "telemetry-ingest http://127.0.0.1:8082/healthz" \
    "feed-vehicle-positions http://127.0.0.1:8083/healthz" \
    "feed-trip-updates http://127.0.0.1:8084/healthz" \
    "feed-alerts http://127.0.0.1:8085/healthz" \
    "gtfs-studio http://127.0.0.1:8086/healthz"
  do
    name="$(printf '%s' "$item" | awk '{print $1}')"
    url="$(printf '%s' "$item" | awk '{print $2}')"
    status="$(curl -sS --connect-timeout "$CONNECT_TIMEOUT_SECONDS" --max-time "$REQUEST_TIMEOUT_SECONDS" -o /dev/null -w '%{http_code}' "$url" 2>/dev/null || true)"
    printf '%s status=%s\n' "$name" "${status:-000}" >>"$OUT_REAL/health/summary.txt"
  done
}

run_avl_dry_run() {
  log "Record synthetic AVL dry-run status"
  mkdir -p "$OUT_REAL/avl"
  if ! command -v go >/dev/null 2>&1; then
    AVL_STATUS="unavailable:go_missing"
    printf 'status=%s\n' "$AVL_STATUS" >"$OUT_REAL/avl/summary"
    return 0
  fi
  raw_stdout="$TMP_DIR/avl-telemetry.raw"
  raw_stderr="$TMP_DIR/avl-diagnostics.raw"
  if go run ./cmd/avl-vendor-adapter --dry-run \
      --reference-time 2026-05-04T12:00:00Z \
      --mapping testdata/avl-vendor/mapping.json \
      testdata/avl-vendor/minimal-gps.json \
      >"$raw_stdout" \
      2>"$raw_stderr"; then
    AVL_STATUS="passed"
  else
    AVL_STATUS="failed"
  fi
  sanitize_file "$raw_stdout" "$OUT_REAL/avl/telemetry.json"
  sanitize_file "$raw_stderr" "$OUT_REAL/avl/diagnostics.json"
  printf 'status=%s\nstdout=telemetry.json\nstderr=diagnostics.json\n' "$AVL_STATUS" >"$OUT_REAL/avl/summary"
}

write_manifest() {
  bundle_ref="$(support_bundle_ref)"
  cat >"$OUT_REAL/manifest.md" <<EOF
# Support Bundle Manifest

- Output reference: \`$bundle_ref\`
- Created at UTC: \`$TIMESTAMP\`
- Included: command/runtime versions, git commit and dirty/clean state, public feed status/size/checksum summaries, unauthenticated admin boundary status, optional authenticated readiness summary, validator tooling status, local health summaries, optional sanitized migration status, synthetic AVL dry-run status.
- Excluded: credential values, admin bearer values, auth headers, cookie values, JWT values, CSRF values, private telemetry records, vendor payload contents, database dumps, environment files, private key material, ACME material, notification credentials, log files, and private operator payload contents.
- Full validation reports are not stored by this support bundle.
- external_evidence_created=false
- consumer_statuses_changed=false
EOF
  cat >"$OUT_REAL/manifest.json" <<EOF
{
  "created_at_utc": "$TIMESTAMP",
  "output_ref": "$bundle_ref",
  "external_evidence_created": false,
  "consumer_statuses_changed": false,
  "included": [
    "system versions",
    "git summary",
    "public feed status and checksum summaries",
    "admin boundary status",
    "optional authenticated readiness summary",
    "validator tooling status",
    "local health summaries",
    "optional sanitized migration status",
    "synthetic AVL dry-run status"
  ],
  "excluded": [
    "credential values",
    "private telemetry records",
    "vendor payload contents",
    "database dumps",
    "environment files",
    "private key material",
    "log files"
  ]
}
EOF
}

redaction_scan() {
  log "Run generated bundle redaction scan"
  scan_root="${1:-$OUT_REAL}"
  python3 - "$scan_root" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
patterns = [
    ("authorization_bearer", re.compile(r"Authorization:\s*Bearer\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("generic_bearer_token", re.compile(r"\bBearer\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("basic_auth_token", re.compile(r"\bBasic\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("jwt_like_value", re.compile(r"\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b")),
    ("db_password_url", re.compile(r"(?:postgres|postgresql|mysql|mariadb|mongodb|redis)://[^:\s/@]+:[^@\s]+@", re.I)),
    ("url_userinfo", re.compile(r"https?://[^/\s:@]+:[^@\s]+@", re.I)),
    ("private_key", re.compile(r"BEGIN [A-Z ]*PRIVATE KEY")),
    ("sensitive_header", re.compile(r"^(?:Cookie|Set-Cookie|Authorization|Proxy-Authorization|X-Api-Key|Api-Key|X-Webhook-Signature):\s*(?!<redacted>)\S+", re.I | re.M)),
    ("env_secret_assignment", re.compile(r"^(?:[A-Z0-9_]*(?:TOKEN|SECRET|PASSWORD|PRIVATE_KEY|WEBHOOK|API_KEY|DATABASE_URL)[A-Z0-9_]*)\s*[:=]\s*(?!<redacted>)\S+", re.I | re.M)),
    ("json_secret_field", re.compile(r'"(?:token|secret|password|authorization|cookie|api_key|webhook_url|database_url|private_key)"\s*:\s*"(?!<redacted>)[^"]+', re.I)),
    ("webhook_secret_url", re.compile(r"https://[^\s]*(?:hooks|webhook|discord|slack)[^\s]*(?:token|secret|key|/services/)[^\s]*", re.I)),
    ("private_path", re.compile(r"(?:(?<=\s)|(?<==)|^)(?:/Users/|/home/|/private/|/var/lib/|/etc/|[A-Za-z]:\\\\Users\\\\)[^\s\"']+")),
    ("raw_payload_label", re.compile(r"\braw[_ -]?(?:payload|telemetry|vendor|log)s?\b", re.I)),
]
findings = []
for path in root.rglob("*"):
    if not path.is_file():
        continue
    try:
        text = path.read_text(encoding="utf-8", errors="replace")
    except OSError:
        continue
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

redaction_self_test() {
  command -v python3 >/dev/null 2>&1 || fail "python3 is required"
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-support-bundle-redaction.XXXXXX")"
  out="$TMP_DIR/out"
  raw="$TMP_DIR/raw.txt"
  mkdir -p "$out"
  cat >"$raw" <<'EOF'
Authorization: Bearer supersecretbearertoken
Proxy-Authorization: Basic Zm9vYmFyYmF6cXV4
Cookie: admin_session=private-cookie
Set-Cookie: admin_session=private-cookie
X-Api-Key: private-api-key
DATABASE_URL=postgres://user:pass@db.internal/open_transit_rt
REDIS_URL=redis://user:pass@cache.internal/0
ADMIN_TOKEN=private-admin-token
DEVICE_TOKEN: private-device-token
webhook_url=https://hooks.slack.com/services/T000/B000/private
private_path=/Users/operator/private/support.log
{"token":"json-token","password":"json-password","authorization":"Bearer json-token","database_url":"postgres://user:pass@db.internal/open_transit_rt","webhook_url":"https://hooks.slack.com/services/T000/B000/private"}
EOF
  sanitize_file "$raw" "$out/sanitized.txt"
  cat >"$out/manifest.md" <<'EOF'
# Redaction Self-Test Manifest

- Output reference: `<support-bundle-output>/self-test`
- external_evidence_created=false
- consumer_statuses_changed=false
EOF
  redaction_scan "$out"
  if grep -E 'supersecret|private-(api-key|admin-token|device-token|cookie)|user:pass|json-token|json-password|hooks\.slack|/Users/operator' "$out/sanitized.txt" >/dev/null 2>&1; then
    fail "redaction self-test left private strings in sanitized output"
  fi
  printf 'support bundle redaction self-test passed\n'
}

copy_paste_summary() {
  public_line="$(awk -F, 'NR>1 {printf "%s:%s/%s ", $1, $5, $2}' "$PUBLIC_SUMMARY_CSV" | sed 's/[[:space:]]*$//')"
  cat <<EOF

Support bundle copy/paste summary:
  output_ref=$(support_bundle_ref)
  public_feed_summary=$public_line
  admin_boundary_result=$ADMIN_BOUNDARY_RESULT
  authenticated_readiness_status=$AUTH_READINESS_STATUS
  validator_tooling_status=$VALIDATOR_TOOLING_STATUS
  validation_api_status=$VALIDATION_API_STATUS
  avl_dry_run_status=$AVL_STATUS
  external_evidence_created=false
  consumer_statuses_changed=false
EOF
}

parse_env() {
  bool_var INCLUDE_ADMIN_READINESS "$INCLUDE_ADMIN_READINESS"
  bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
  bool_var FORCE "$FORCE"
  positive_int CONNECT_TIMEOUT_SECONDS "$CONNECT_TIMEOUT_SECONDS"
  positive_int REQUEST_TIMEOUT_SECONDS "$REQUEST_TIMEOUT_SECONDS"
  positive_int MAX_FEED_BYTES "$MAX_FEED_BYTES"
}

main() {
  case "${1:-}" in
    -h|--help|help) usage; exit 0 ;;
    --self-test-redaction) redaction_self_test; exit 0 ;;
    "") ;;
    *) usage; fail "unknown argument: $1" ;;
  esac

  command -v python3 >/dev/null 2>&1 || fail "python3 is required"
  parse_env
  prepare_output_dir
  TMP_DIR="${TMP_DIR:-$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-support-bundle.XXXXXX")}"
  resolve_urls

  ADMIN_BOUNDARY_RESULT="not_run"
  AUTH_READINESS_STATUS="not_run"
  VALIDATOR_TOOLING_STATUS="not_run"
  VALIDATION_API_STATUS="not_run"
  MIGRATION_STATUS="not_run"
  AVL_STATUS="not_run"

  record_versions
  record_public_feeds
  record_admin_boundary
  record_authenticated_readiness
  record_validator_tooling
  record_migration_status
  record_local_health
  run_avl_dry_run
  write_manifest
  redaction_scan
  copy_paste_summary
}

main "$@"
