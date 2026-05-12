#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/validate-public-feeds/$TIMESTAMP}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-}"
FORCE="${FORCE:-false}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
FETCH_CONNECT_TIMEOUT="${FETCH_CONNECT_TIMEOUT:-10}"
FETCH_TIMEOUT="${FETCH_TIMEOUT:-90}"
STRICT="${STRICT:-false}"
SKIP_VALIDATORS="${SKIP_VALIDATORS:-false}"
DRY_RUN="false"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/validate-public-feeds.sh [--public-base-url URL] [--output-dir DIR] [--dry-run] [--strict] [--skip-validators]

Environment:
  PUBLIC_BASE_URL              Public root to fetch, for example https://feeds.example.org
  OUTPUT_DIR                   Default .cache/validate-public-feeds/<UTC timestamp>
  FORCE                        true|false; allow non-empty output reuse
  ALLOW_UNIGNORED_OUTPUT_DIR   true|false; allow output outside .cache except evidence-like paths
  FETCH_CONNECT_TIMEOUT        Curl connect timeout seconds, default 10
  FETCH_TIMEOUT                Curl total timeout seconds, default 90
  STRICT                       true|false; fail on failed fetches, failed validators, or missing validator tooling
  SKIP_VALIDATORS              true|false; fetch artifacts without running validators

Safety:
  This helper is an off-host diagnostic path for public feed artifacts. It
  writes under .cache by default, never writes docs/evidence, does not contact
  consumers, does not change consumer statuses, and does not claim compliance,
  consumer acceptance, agency adoption, production readiness, hosted service
  availability, vendor compatibility, SLA coverage, or production-grade ETA
  quality. Validator results are supporting signals only.
EOF
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT INT TERM

bool_var() {
  case "$2" in
    true|false) ;;
    *) fail "$1 must be true or false" ;;
  esac
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --public-base-url)
      shift
      PUBLIC_BASE_URL="${1:-}"
      ;;
    --output-dir)
      shift
      OUTPUT_DIR="${1:-}"
      ;;
    --dry-run)
      DRY_RUN="true"
      ;;
    --strict)
      STRICT="true"
      ;;
    --skip-validators)
      SKIP_VALIDATORS="true"
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

bool_var FORCE "$FORCE"
bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
bool_var STRICT "$STRICT"
bool_var SKIP_VALIDATORS "$SKIP_VALIDATORS"

case "$FETCH_CONNECT_TIMEOUT" in ''|*[!0-9]*) fail "FETCH_CONNECT_TIMEOUT must be a positive whole number" ;; esac
case "$FETCH_TIMEOUT" in ''|*[!0-9]*) fail "FETCH_TIMEOUT must be a positive whole number" ;; esac

if [ -n "$PUBLIC_BASE_URL" ]; then
  case "$PUBLIC_BASE_URL" in
    http://*|https://*) ;;
    *) fail "PUBLIC_BASE_URL must start with http:// or https://" ;;
  esac
  case "$PUBLIC_BASE_URL" in
    *"@"*|*"?"*|*"#"*) fail "PUBLIC_BASE_URL must be a public root without userinfo, query, or fragment" ;;
  esac
  PUBLIC_BASE_URL="${PUBLIC_BASE_URL%/}"
fi

mkdir -p .cache
TMP_DIR="$(mktemp -d ".cache/validate-public-feeds-tmp.XXXXXX")"
chmod 700 "$TMP_DIR"

python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$FORCE" "$ALLOW_UNIGNORED_OUTPUT_DIR" <<'PY' >"$TMP_DIR/output-dir"
import pathlib
import shutil
import sys

root = pathlib.Path(sys.argv[1]).resolve()
raw = pathlib.Path(sys.argv[2])
force = sys.argv[3] == "true"
allow = sys.argv[4] == "true"
out = raw if raw.is_absolute() else root / raw
resolved = out.resolve(strict=False)
cache = (root / ".cache").resolve(strict=False)

def evidence_like(path):
    text = str(path).replace("\\", "/").lower()
    parts = [p.lower() for p in pathlib.Path(path).parts]
    return "docs/evidence" in text or "evidence" in parts or "submission" in parts or "proof" in parts

def has_symlink(path):
    probe = pathlib.Path(path)
    if not probe.is_absolute():
        probe = root / probe
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False

if evidence_like(raw) or evidence_like(resolved):
    raise SystemExit("OUTPUT_DIR must not be evidence-like or under docs/evidence")
if has_symlink(raw):
    raise SystemExit("OUTPUT_DIR must not contain symlink directories")
if not allow:
    try:
        resolved.relative_to(cache)
    except ValueError:
        raise SystemExit("OUTPUT_DIR must resolve under repo .cache unless ALLOW_UNIGNORED_OUTPUT_DIR=true")
if resolved.exists() and not resolved.is_dir():
    raise SystemExit("OUTPUT_DIR must be a directory")
if resolved.exists() and any(resolved.iterdir()):
    if not force:
        raise SystemExit("OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it")
    for child in resolved.iterdir():
        if child.is_symlink() or child.is_file():
            child.unlink()
        else:
            shutil.rmtree(child)
resolved.mkdir(parents=True, exist_ok=True)
resolved.chmod(0o700)
print(resolved)
PY

OUT_REAL="$(cat "$TMP_DIR/output-dir")"
ARTIFACT_DIR="$OUT_REAL/artifacts"
HEADER_DIR="$OUT_REAL/headers"
VALIDATOR_DIR="$OUT_REAL/validators"
mkdir -p "$ARTIFACT_DIR" "$HEADER_DIR" "$VALIDATOR_DIR"
ROWS_TSV="$TMP_DIR/rows.tsv"
printf 'id\tlabel\tpublic_path\tconfigured_url\tfetch_status\thttp_status\tbyte_count\tcontent_type\tchecksum\tartifact_file\tvalidator_state\tnext_action\tdoes_not_prove\n' >"$ROWS_TSV"

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

content_type_from_headers() {
  awk 'BEGIN{IGNORECASE=1} /^content-type:/ {sub(/\r$/, "", $0); sub(/^[^:]*:[[:space:]]*/, "", $0); print; exit}' "$1" 2>/dev/null || true
}

rel_output_file() {
  python3 - "$ROOT_DIR" "$1" <<'PY'
import pathlib
import sys
root = pathlib.Path(sys.argv[1]).resolve()
path = pathlib.Path(sys.argv[2]).resolve(strict=False)
try:
    print(path.relative_to(root).as_posix())
except ValueError:
    print(path.name)
PY
}

append_row() {
  printf '%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n' "$@" >>"$ROWS_TSV"
}

fetch_feed() {
  id="$1"
  label="$2"
  public_path="$3"
  file_name="$4"
  validator_kind="$5"
  configured_url=""
  if [ -n "$PUBLIC_BASE_URL" ]; then
    configured_url="${PUBLIC_BASE_URL}${public_path}"
  fi
  artifact="$ARTIFACT_DIR/$file_name"
  headers="$HEADER_DIR/$id.headers"
  fetch_status="not_checked"
  http_status="not checked"
  byte_count="not recorded"
  content_type="not recorded"
  checksum="not recorded"
  artifact_file="not available"
  validator_state="not_applicable"
  next_action="Set PUBLIC_BASE_URL, then rerun this helper from an off-host operator machine."
  if [ "$validator_kind" != "none" ]; then
    validator_state="skipped"
  fi

  if [ "$DRY_RUN" = "true" ]; then
    fetch_status="not_checked"
    next_action="Dry run only. Rerun without --dry-run to fetch and validate this public path."
  elif [ -z "$PUBLIC_BASE_URL" ]; then
    fetch_status="missing_base_url"
  else
    curl_error="$HEADER_DIR/$id.curl-error.txt"
    code="000"
    if code="$(curl -L -sS --connect-timeout "$FETCH_CONNECT_TIMEOUT" --max-time "$FETCH_TIMEOUT" -D "$headers" -o "$artifact" -w '%{http_code}' "$configured_url" 2>"$curl_error")"; then
      http_status="$code"
      if [ "$code" = "200" ] && [ -s "$artifact" ]; then
        fetch_status="fetched"
        byte_count="$(wc -c <"$artifact" | tr -d ' ')"
        checksum="$(sha256_file "$artifact")"
        content_type="$(content_type_from_headers "$headers")"
        if [ -z "$content_type" ]; then
          content_type="not recorded"
        fi
        artifact_file="$(rel_output_file "$artifact")"
        next_action="Review validator state and compare byte count/checksum with the last local check."
      else
        fetch_status="failed_fetch"
        next_action="Confirm the public route, proxy, and service health before relying on this path."
      fi
    else
      fetch_status="failed_fetch"
      http_status="curl failed"
      next_action="Confirm DNS, TLS, proxy routing, and service health before relying on this path."
    fi
  fi

  if [ "$validator_kind" = "static" ]; then
    validator_state="$(run_static_validator "$fetch_status" "$artifact" "$id")"
  elif [ "$validator_kind" = "realtime" ]; then
    validator_state="$(run_realtime_validator "$fetch_status" "$artifact" "$id")"
  fi

  append_row "$id" "$label" "$public_path" "$configured_url" "$fetch_status" "$http_status" "$byte_count" "$content_type" "$checksum" "$artifact_file" "$validator_state" "$next_action" "Supporting diagnostics only; does not prove compliance, consumer acceptance, final-root readiness, agency adoption, or production readiness."
}

java_binary() {
  if [ -n "${JAVA_BINARY:-}" ] && command -v "$JAVA_BINARY" >/dev/null 2>&1; then
    printf '%s\n' "$JAVA_BINARY"
    return 0
  fi
  if command -v java >/dev/null 2>&1; then
    command -v java
    return 0
  fi
  return 1
}

default_static_validator() {
  if [ -n "${GTFS_VALIDATOR_PATH:-}" ]; then
    printf '%s\n' "$GTFS_VALIDATOR_PATH"
  elif [ -f ".cache/validators/gtfs-validator-7.1.0-cli.jar" ]; then
    printf '%s\n' ".cache/validators/gtfs-validator-7.1.0-cli.jar"
  fi
}

default_realtime_validator() {
  if [ -n "${GTFS_RT_VALIDATOR_PATH:-}" ]; then
    printf '%s\n' "$GTFS_RT_VALIDATOR_PATH"
  elif [ -x ".cache/validators/gtfs-rt-validator-wrapper.sh" ]; then
    printf '%s\n' ".cache/validators/gtfs-rt-validator-wrapper.sh"
  fi
}

run_static_validator() {
  fetch_status="$1"
  artifact="$2"
  id="$3"
  if [ "$SKIP_VALIDATORS" = "true" ]; then
    printf '%s\n' "skipped"
    return
  fi
  if [ "$DRY_RUN" = "true" ]; then
    printf '%s\n' "not_checked"
    return
  fi
  if [ "$fetch_status" != "fetched" ]; then
    printf '%s\n' "artifact_unavailable"
    return
  fi
  validator="$(default_static_validator || true)"
  if [ -z "$validator" ] || [ ! -f "$validator" ]; then
    printf '%s\n' "missing_tooling"
    return
  fi
  java_bin="$(java_binary || true)"
  out_dir="$VALIDATOR_DIR/$id"
  mkdir -p "$out_dir"
  if [ "${validator##*.}" = "jar" ]; then
    if [ -z "$java_bin" ]; then
      printf '%s\n' "missing_tooling"
      return
    fi
    if "$java_bin" -jar "$validator" -i "$artifact" -o "$out_dir/report" >"$out_dir/stdout.log" 2>"$out_dir/stderr.log"; then
      printf '%s\n' "completed"
    else
      printf '%s\n' "validation_failed"
    fi
  else
    if [ ! -x "$validator" ]; then
      printf '%s\n' "missing_tooling"
      return
    fi
    if "$validator" -i "$artifact" -o "$out_dir/report" >"$out_dir/stdout.log" 2>"$out_dir/stderr.log"; then
      printf '%s\n' "completed"
    else
      printf '%s\n' "validation_failed"
    fi
  fi
}

run_realtime_validator() {
  fetch_status="$1"
  artifact="$2"
  id="$3"
  if [ "$SKIP_VALIDATORS" = "true" ]; then
    printf '%s\n' "skipped"
    return
  fi
  if [ "$DRY_RUN" = "true" ]; then
    printf '%s\n' "not_checked"
    return
  fi
  if [ "$fetch_status" != "fetched" ]; then
    printf '%s\n' "artifact_unavailable"
    return
  fi
  schedule="$ARTIFACT_DIR/schedule.zip"
  if [ ! -s "$schedule" ]; then
    printf '%s\n' "artifact_unavailable"
    return
  fi
  validator="$(default_realtime_validator || true)"
  if [ -z "$validator" ] || [ ! -x "$validator" ]; then
    printf '%s\n' "missing_tooling"
    return
  fi
  out_dir="$VALIDATOR_DIR/$id"
  mkdir -p "$out_dir"
  if "$validator" --schedule "$schedule" --realtime "$artifact" --feed_type "$id" --output_dir "$out_dir/report" >"$out_dir/stdout.log" 2>"$out_dir/stderr.log"; then
    printf '%s\n' "completed"
  else
    printf '%s\n' "validation_failed"
  fi
}

fetch_feed "feeds_json" "feeds.json" "/public/feeds.json" "feeds.json" "none"
fetch_feed "schedule" "Static GTFS schedule" "/public/gtfs/schedule.zip" "schedule.zip" "static"
fetch_feed "vehicle_positions" "Vehicle Positions" "/public/gtfsrt/vehicle_positions.pb" "vehicle_positions.pb" "realtime"
fetch_feed "trip_updates" "Trip Updates" "/public/gtfsrt/trip_updates.pb" "trip_updates.pb" "realtime"
fetch_feed "alerts" "Alerts" "/public/gtfsrt/alerts.pb" "alerts.pb" "realtime"

python3 - "$ROOT_DIR" "$OUT_REAL" "$ROWS_TSV" "$TIMESTAMP" "$PUBLIC_BASE_URL" "$DRY_RUN" "$STRICT" "$SKIP_VALIDATORS" <<'PY'
import csv
import json
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out = pathlib.Path(sys.argv[2]).resolve()
rows_path = pathlib.Path(sys.argv[3])
generated_at = sys.argv[4]
public_base_url = sys.argv[5]
dry_run = sys.argv[6] == "true"
strict = sys.argv[7] == "true"
skip_validators = sys.argv[8] == "true"

rows = list(csv.DictReader(rows_path.open(), delimiter="\t"))
claim_flags = {
    "external_evidence_created": False,
    "consumer_statuses_changed": False,
    "compliance_claimed": False,
    "production_readiness_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "public_launch_claimed": False,
    "hosted_saas_claimed": False,
    "vendor_compatibility_claimed": False,
    "sla_claimed": False,
    "uptime_guarantee_claimed": False,
    "production_grade_eta_claimed": False,
}
fetch_blockers = [r for r in rows if r["fetch_status"] in {"failed_fetch", "missing_base_url"}]
validator_blockers = [
    r for r in rows
    if r["validator_state"] in {"missing_tooling", "validation_failed", "artifact_unavailable"} and r["id"] != "feeds_json"
]
summary = {
    "generated_at": generated_at,
    "public_base_url_configured": bool(public_base_url),
    "dry_run": dry_run,
    "strict_mode": strict,
    "skip_validators": skip_validators,
    "rows": rows,
    "counts": {
        "rows": len(rows),
        "fetch_blockers": len(fetch_blockers),
        "validator_blockers": len(validator_blockers),
    },
    "claim_flags": claim_flags,
    "boundary": "Off-host public feed diagnostics only. Validator completion is a supporting signal, not compliance, consumer acceptance, agency adoption, final-root proof, or production readiness.",
    "strict_failed": bool(strict and (fetch_blockers or (validator_blockers and not skip_validators))),
}
manifest = {
    "generated_at": generated_at,
    "output_dir": out.relative_to(root).as_posix() if out.is_relative_to(root) else out.name,
    "included_files": ["summary.json", "summary.md", "manifest.json", "manifest.md", "artifacts/", "headers/", "validators/"],
    "excluded_categories": ["docs/evidence", "consumer submissions", "secret values", "raw tokens", "private keys", "database URLs"],
    "claim_flags": claim_flags,
}
(out / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
(out / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")

lines = [
    "# Validate Public Feeds Summary",
    "",
    f"Generated at: {generated_at}",
    "",
    "This is an off-host diagnostic summary. It is not compliance, consumer acceptance, agency adoption, final-root proof, or production readiness.",
    "",
    "| Feed | Path | Fetch | HTTP | Bytes | Checksum | Validator | Next action |",
    "| --- | --- | --- | --- | --- | --- | --- | --- |",
]
for row in rows:
    lines.append("| {label} | `{public_path}` | {fetch_status} | {http_status} | {byte_count} | `{checksum}` | {validator_state} | {next_action} |".format(**row))
(out / "summary.md").write_text("\n".join(lines) + "\n")
(out / "manifest.md").write_text("# Manifest\n\n" + "\n".join(f"- {name}" for name in manifest["included_files"]) + "\n")
PY

if python3 - "$OUT_REAL/summary.json" <<'PY'
import json
import sys
summary = json.loads(open(sys.argv[1]).read())
raise SystemExit(1 if summary.get("strict_failed") else 0)
PY
then
  printf 'validate-public-feeds summary written to %s\n' "$(rel_output_file "$OUT_REAL")"
else
  printf 'validate-public-feeds found strict blockers; summary written to %s\n' "$(rel_output_file "$OUT_REAL")" >&2
  exit 1
fi
