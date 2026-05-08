#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-http://localhost:8080}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
SKIP_VALIDATORS="${SKIP_VALIDATORS:-false}"
STRICT_VALIDATORS="${STRICT_VALIDATORS:-false}"
CONNECT_TIMEOUT_SECONDS="${CONNECT_TIMEOUT_SECONDS:-5}"
REQUEST_TIMEOUT_SECONDS="${REQUEST_TIMEOUT_SECONDS:-30}"
MAX_FEED_BYTES="${MAX_FEED_BYTES:-104857600}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
FORCE="${FORCE:-false}"
TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/operator-smoke/$TIMESTAMP}"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/operator-smoke.sh

Environment:
  PUBLIC_BASE_URL                 Public feed root, default http://localhost:8080
  ADMIN_BASE_URL                  Admin root for authenticated checks; defaults to PUBLIC_BASE_URL only when loopback
  ADMIN_TOKEN                     Optional admin bearer token
  SKIP_VALIDATORS                 true|false; skips validator API calls, missing tooling is non-fatal
  STRICT_VALIDATORS               true|false; missing/misconfigured tooling is fatal
  OUTPUT_DIR                      Default .cache/operator-smoke/<timestamp>
  ALLOW_UNIGNORED_OUTPUT_DIR      true|false; allow OUTPUT_DIR outside repo .cache
  FORCE                           true|false; allow non-empty OUTPUT_DIR reuse
  CONNECT_TIMEOUT_SECONDS         curl connect timeout, default 5
  REQUEST_TIMEOUT_SECONDS         curl total timeout, default 30
  MAX_FEED_BYTES                  maximum accepted public feed size, default 104857600

Safety:
  Admin requests use Authorization: Bearer "$ADMIN_TOKEN" only.
  Cookie auth is not supported.
  Tokens, Authorization headers, cookies, JWTs, and CSRF values are never printed or written.
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
  command -v "$1" >/dev/null 2>&1 || fail "missing required tool: $1"
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
netloc = parts.netloc
normalized = urlunsplit((parts.scheme, netloc, path, "", ""))
print(normalized)
print("true" if is_loopback else "false")
PY
}

output_realpath() {
  python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$ALLOW_UNIGNORED_OUTPUT_DIR" <<'PY'
import os
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
  if [ "$ALLOW_UNIGNORED_OUTPUT_DIR" != "true" ] && [ "$ALLOW_UNIGNORED_OUTPUT_DIR" != "false" ]; then
    fail "ALLOW_UNIGNORED_OUTPUT_DIR must be true or false"
  fi
  if [ "$FORCE" != "true" ] && [ "$FORCE" != "false" ]; then
    fail "FORCE must be true or false"
  fi
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
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-operator-smoke.XXXXXX")"
}

path_url() {
  base="$1"
  path="$2"
  printf '%s%s\n' "$base" "$path"
}

record_public_fetch() {
  label="$1"
  path="$2"
  out="$OUT_REAL/public/$label"
  headers="$OUT_REAL/public/$label.headers.txt"
  meta="$OUT_REAL/public/$label.meta"
  tmp="$TMP_DIR/$label.tmp"
  url="$(path_url "$PUBLIC_BASE_NORMALIZED" "$path")"
  mkdir -p "$OUT_REAL/public"
  status="000"
  curl_exit=0
  set +e
  curl -L -sS \
      --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
      --max-time "$REQUEST_TIMEOUT_SECONDS" \
      --max-filesize "$MAX_FEED_BYTES" \
      -D "$headers" \
      -o "$tmp" \
      -w 'http_code=%{http_code}\nurl_effective=%{url_effective}\nnum_redirects=%{num_redirects}\nsize_download=%{size_download}\ncontent_type=%{content_type}\n' \
      "$url" >"$meta" 2>"$OUT_REAL/public/$label.curl-stderr.txt"
  curl_exit=$?
  set -e
  status="$(sed -n 's/^http_code=//p' "$meta" | tail -n 1)"
  status="${status:-000}"
  bytes=0
  if [ -f "$tmp" ]; then
    bytes="$(wc -c <"$tmp" | awk '{print $1}')"
  fi
  outcome="failed"
  sha=""
  if [ "$bytes" -gt "$MAX_FEED_BYTES" ] || [ "$curl_exit" -eq 63 ]; then
    outcome="too_large"
    rm -f "$tmp"
  elif [ "$curl_exit" -ne 0 ]; then
    outcome="curl_failed"
    rm -f "$tmp"
  elif [ "$status" -lt 200 ] || [ "$status" -ge 300 ]; then
    outcome="bad_status"
    rm -f "$tmp"
  elif [ "$bytes" -le 0 ]; then
    outcome="empty"
    rm -f "$tmp"
  else
    mv "$tmp" "$out"
    sha="$(sha256_file "$out")"
    outcome="ok"
  fi
  {
    printf 'label=%s\n' "$label"
    printf 'path=%s\n' "$path"
    printf 'status=%s\n' "$status"
    printf 'bytes=%s\n' "$bytes"
    printf 'sha256=%s\n' "$sha"
    printf 'outcome=%s\n' "$outcome"
    printf 'curl_exit=%s\n' "$curl_exit"
    sed -n '/^url_effective=/p;/^num_redirects=/p;/^content_type=/p' "$meta"
  } >"$OUT_REAL/public/$label.summary"
  printf '%s,%s,%s,%s,%s\n' "$label" "$status" "$bytes" "$sha" "$outcome" >>"$PUBLIC_SUMMARY_CSV"
  if [ "$outcome" != "ok" ]; then
    PUBLIC_FAILURES=$((PUBLIC_FAILURES + 1))
  fi
}

http_status_no_redirect() {
  url="$1"
  curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$TMP_DIR/admin-body.tmp" \
    -w '%{http_code}' \
    "$url" 2>/dev/null || true
}

check_admin_boundary() {
  log "Check unauthenticated admin readiness boundary"
  mkdir -p "$OUT_REAL/admin"
  url="$(path_url "$PUBLIC_BASE_NORMALIZED" "/admin/operations/readiness")"
  status="$(http_status_no_redirect "$url")"
  status="${status:-000}"
  case "$status" in
    401|403|404) ADMIN_BOUNDARY_RESULT="passed:$status" ;;
    2*) ADMIN_BOUNDARY_RESULT="failed:unexpected_2xx_$status"; ADMIN_FAILURES=$((ADMIN_FAILURES + 1)) ;;
    *) ADMIN_BOUNDARY_RESULT="failed:unexpected_status_$status"; ADMIN_FAILURES=$((ADMIN_FAILURES + 1)) ;;
  esac
  printf 'status=%s\nresult=%s\n' "$status" "$ADMIN_BOUNDARY_RESULT" >"$OUT_REAL/admin/public-readiness.summary"
  rm -f "$TMP_DIR/admin-body.tmp"
}

check_authenticated_readiness() {
  mkdir -p "$OUT_REAL/admin"
  if [ -z "$ADMIN_TOKEN" ]; then
    AUTH_READINESS_STATUS="skipped:no_admin_token"
    printf 'status=skipped\nreason=no_admin_token\n' >"$OUT_REAL/admin/auth-readiness.summary"
    return 0
  fi
  if [ -z "${ADMIN_BASE_NORMALIZED:-}" ]; then
    fail "ADMIN_TOKEN was supplied but ADMIN_BASE_URL is missing or unsafe"
  fi
  log "Check authenticated admin readiness"
  url="$(path_url "$ADMIN_BASE_NORMALIZED" "/admin/operations/readiness")"
  body="$TMP_DIR/auth-readiness.html"
  status="$(curl -sS \
    --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
    --max-time "$REQUEST_TIMEOUT_SECONDS" \
    -o "$body" \
    -w '%{http_code}' \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    "$url" 2>/dev/null || true)"
  status="${status:-000}"
  summary="$(python3 - "$body" "$status" <<'PY'
import re
import sys
path, status = sys.argv[1:3]
text = ""
try:
    text = open(path, encoding="utf-8", errors="replace").read(200000)
except OSError:
    pass
title = ""
m = re.search(r"<h2>(.*?)</h2>", text, re.I | re.S)
if m:
    title = re.sub(r"\s+", " ", m.group(1)).strip()
rows = len(re.findall(r"<tr>", text, re.I))
print(f"status={status}")
print(f"title={title}")
print(f"table_rows={rows}")
PY
)"
  printf '%s\n' "$summary" >"$OUT_REAL/admin/auth-readiness.summary"
  case "$status" in
    2*) AUTH_READINESS_STATUS="passed:$status" ;;
    *) AUTH_READINESS_STATUS="failed:$status"; ADMIN_FAILURES=$((ADMIN_FAILURES + 1)) ;;
  esac
  rm -f "$body"
}

record_validator_tooling() {
  log "Check pinned validator tooling state"
  mkdir -p "$OUT_REAL/validators"
  set +e
  ./scripts/check-validators.sh >"$OUT_REAL/validators/tooling.txt" 2>&1
  rc=$?
  set -e
  case "$rc" in
    0) VALIDATOR_TOOLING_STATUS="passed" ;;
    11) VALIDATOR_TOOLING_STATUS="missing" ;;
    12) VALIDATOR_TOOLING_STATUS="misconfigured" ;;
    *) VALIDATOR_TOOLING_STATUS="failed_exit_$rc" ;;
  esac
  printf 'exit_code=%s\nstatus=%s\n' "$rc" "$VALIDATOR_TOOLING_STATUS" >"$OUT_REAL/validators/tooling.summary"
  if [ "$STRICT_VALIDATORS" = "true" ] && [ "$rc" -ne 0 ]; then
    fail "validator tooling check failed in strict mode; see $OUT_REAL/validators/tooling.txt"
  fi
}

summarize_validation_json() {
  in="$1"
  out="$2"
  feed_type="$3"
  validator_id="$4"
  python3 - "$in" "$out" "$feed_type" "$validator_id" <<'PY'
import json
import sys

src, dst, feed_type, validator_id = sys.argv[1:5]
summary = {
    "feed_type": feed_type,
    "validator_id": validator_id,
    "status": "unknown",
    "error_count": 0,
    "warning_count": 0,
    "info_count": 0,
    "validator_name": "",
    "validator_version": "",
}
try:
    data = json.load(open(src, encoding="utf-8"))
    for key in ("status", "error_count", "warning_count", "info_count", "validator_name", "validator_version"):
        if key in data:
            summary[key] = data[key]
except Exception as exc:
    summary["status"] = "parse_failed"
    summary["error"] = str(exc)
with open(dst, "w", encoding="utf-8") as f:
    json.dump(summary, f, indent=2, sort_keys=True)
    f.write("\n")
PY
}

run_validation_api() {
  mkdir -p "$OUT_REAL/validators"
  if [ "$SKIP_VALIDATORS" = "true" ]; then
    VALIDATION_API_STATUS="skipped:SKIP_VALIDATORS=true"
    printf 'status=skipped\nreason=SKIP_VALIDATORS=true\n' >"$OUT_REAL/validators/api.summary"
    return 0
  fi
  if [ "$VALIDATOR_TOOLING_STATUS" != "passed" ]; then
    VALIDATION_API_STATUS="blocked:validator_tooling_$VALIDATOR_TOOLING_STATUS"
    printf 'status=blocked\nreason=validator_tooling_%s\n' "$VALIDATOR_TOOLING_STATUS" >"$OUT_REAL/validators/api.summary"
    return 0
  fi
  if [ -z "$ADMIN_TOKEN" ]; then
    VALIDATION_API_STATUS="skipped:no_admin_token"
    printf 'status=skipped\nreason=no_admin_token\n' >"$OUT_REAL/validators/api.summary"
    return 0
  fi
  if [ -z "${ADMIN_BASE_NORMALIZED:-}" ]; then
    fail "ADMIN_TOKEN was supplied but ADMIN_BASE_URL is missing or unsafe"
  fi
  log "Run allowlisted validation API calls"
  failures=0
  : >"$OUT_REAL/validators/api.summary"
  for pair in \
    "schedule static-mobilitydata" \
    "vehicle_positions realtime-mobilitydata" \
    "trip_updates realtime-mobilitydata" \
    "alerts realtime-mobilitydata"
  do
    feed_type="$(printf '%s' "$pair" | awk '{print $1}')"
    validator_id="$(printf '%s' "$pair" | awk '{print $2}')"
    body="$TMP_DIR/validation-$feed_type.json"
    summary="$OUT_REAL/validators/validation-$feed_type.summary.json"
    request="$(printf '{"validator_id":"%s","feed_type":"%s"}' "$validator_id" "$feed_type")"
    status="$(curl -sS \
      --connect-timeout "$CONNECT_TIMEOUT_SECONDS" \
      --max-time "$REQUEST_TIMEOUT_SECONDS" \
      -o "$body" \
      -w '%{http_code}' \
      -X POST "$(path_url "$ADMIN_BASE_NORMALIZED" "/admin/validation/run")" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/json" \
      --data "$request" 2>/dev/null || true)"
    summarize_validation_json "$body" "$summary" "$feed_type" "$validator_id"
    result_status="$(python3 - "$summary" <<'PY'
import json
import sys
print(json.load(open(sys.argv[1], encoding="utf-8")).get("status", "unknown"))
PY
)"
    printf '%s http_status=%s result_status=%s summary=%s\n' "$feed_type" "${status:-000}" "$result_status" "$(basename "$summary")" >>"$OUT_REAL/validators/api.summary"
    case "${status:-000}:$result_status" in
      2*:passed|2*:warning) ;;
      *) failures=$((failures + 1)) ;;
    esac
    rm -f "$body"
  done
  if [ "$failures" -gt 0 ]; then
    VALIDATION_API_STATUS="blocked:${failures}_validation_call_failure(s)"
    if [ "$STRICT_VALIDATORS" = "true" ]; then
      fail "$VALIDATION_API_STATUS"
    fi
  else
    VALIDATION_API_STATUS="completed"
  fi
}

run_avl_dry_run() {
  log "Run synthetic AVL dry-run fixture"
  mkdir -p "$OUT_REAL/avl"
  if go run ./cmd/avl-vendor-adapter --dry-run \
      --reference-time 2026-05-04T12:00:00Z \
      --mapping testdata/avl-vendor/mapping.json \
      testdata/avl-vendor/minimal-gps.json \
      >"$OUT_REAL/avl/telemetry.json" \
      2>"$OUT_REAL/avl/diagnostics.json"; then
    AVL_STATUS="passed"
  else
    AVL_STATUS="failed"
    return 1
  fi
  printf 'status=%s\nstdout=telemetry.json\nstderr=diagnostics.json\n' "$AVL_STATUS" >"$OUT_REAL/avl/summary"
}

write_manifest() {
  cat >"$OUT_REAL/manifest.md" <<EOF
# Operator Smoke Manifest

- Output directory: \`$OUT_REAL\`
- Created at UTC: \`$TIMESTAMP\`
- Included: public feed bodies, public feed headers, public feed summaries, validator tooling status, validation API summaries when run, readiness status summaries, synthetic AVL dry-run output.
- Excluded: admin bearer tokens, Authorization headers, cookies, JWTs, CSRF values, raw private telemetry, private vendor payloads, raw database dumps, raw environment files, private keys, ACME material, webhook URLs, notification credentials, and unredacted logs.
- external_evidence_created=false
- consumer_statuses_changed=false
EOF
}

copy_paste_summary() {
  public_line="$(awk -F, 'NR>1 {printf "%s:%s/%s ", $1, $5, $2}' "$PUBLIC_SUMMARY_CSV" | sed 's/[[:space:]]*$//')"
  cat <<EOF

Operator smoke copy/paste summary:
  output_dir=$OUT_REAL
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
  bool_var SKIP_VALIDATORS "$SKIP_VALIDATORS"
  bool_var STRICT_VALIDATORS "$STRICT_VALIDATORS"
  bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
  bool_var FORCE "$FORCE"
  positive_int CONNECT_TIMEOUT_SECONDS "$CONNECT_TIMEOUT_SECONDS"
  positive_int REQUEST_TIMEOUT_SECONDS "$REQUEST_TIMEOUT_SECONDS"
  positive_int MAX_FEED_BYTES "$MAX_FEED_BYTES"
  if [ "$SKIP_VALIDATORS" = "true" ] && [ "$STRICT_VALIDATORS" = "true" ]; then
    fail "SKIP_VALIDATORS and STRICT_VALIDATORS cannot both be true"
  fi
}

resolve_urls() {
  public_info="$(normalize_url "$PUBLIC_BASE_URL" "PUBLIC_BASE_URL")" || fail "$public_info"
  PUBLIC_BASE_NORMALIZED="$(printf '%s\n' "$public_info" | sed -n '1p')"
  PUBLIC_BASE_LOOPBACK="$(printf '%s\n' "$public_info" | sed -n '2p')"
  if [ -z "$ADMIN_BASE_URL" ] && [ "$PUBLIC_BASE_LOOPBACK" = "true" ]; then
    ADMIN_BASE_URL="$PUBLIC_BASE_NORMALIZED"
  fi
  if [ -n "$ADMIN_BASE_URL" ]; then
    admin_info="$(normalize_url "$ADMIN_BASE_URL" "ADMIN_BASE_URL")" || fail "$admin_info"
    ADMIN_BASE_NORMALIZED="$(printf '%s\n' "$admin_info" | sed -n '1p')"
  else
    ADMIN_BASE_NORMALIZED=""
  fi
  if [ -n "$ADMIN_TOKEN" ] && [ -z "$ADMIN_BASE_NORMALIZED" ]; then
    fail "ADMIN_TOKEN was supplied, but safe ADMIN_BASE_URL is not available"
  fi
}

main() {
  case "${1:-}" in
    -h|--help|help) usage; exit 0 ;;
    "") ;;
    *) usage; fail "unknown argument: $1" ;;
  esac

  need curl
  need awk
  need sed
  need python3
  need go
  parse_env
  resolve_urls
  prepare_output_dir

  PUBLIC_FAILURES=0
  ADMIN_FAILURES=0
  ADMIN_BOUNDARY_RESULT="not_run"
  AUTH_READINESS_STATUS="not_run"
  VALIDATOR_TOOLING_STATUS="not_run"
  VALIDATION_API_STATUS="not_run"
  AVL_STATUS="not_run"
  PUBLIC_SUMMARY_CSV="$OUT_REAL/public-summary.csv"
  printf 'label,status,bytes,sha256,outcome\n' >"$PUBLIC_SUMMARY_CSV"

  log "Fetch public feed paths"
  record_public_fetch "feeds.json" "/public/feeds.json"
  record_public_fetch "schedule.zip" "/public/gtfs/schedule.zip"
  record_public_fetch "vehicle_positions.pb" "/public/gtfsrt/vehicle_positions.pb"
  record_public_fetch "trip_updates.pb" "/public/gtfsrt/trip_updates.pb"
  record_public_fetch "alerts.pb" "/public/gtfsrt/alerts.pb"

  check_admin_boundary
  check_authenticated_readiness
  record_validator_tooling
  run_validation_api
  run_avl_dry_run || true
  write_manifest
  copy_paste_summary

  if [ "$PUBLIC_FAILURES" -gt 0 ]; then
    fail "$PUBLIC_FAILURES public feed check(s) failed"
  fi
  if [ "$ADMIN_FAILURES" -gt 0 ]; then
    fail "$ADMIN_FAILURES admin readiness/boundary check(s) failed"
  fi
  if [ "$AVL_STATUS" != "passed" ]; then
    fail "synthetic AVL dry-run failed"
  fi
}

main "$@"
