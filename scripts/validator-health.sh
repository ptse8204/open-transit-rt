#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/validator-health/$TIMESTAMP}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
FORCE="${FORCE:-false}"
RUN_VALIDATORS="${RUN_VALIDATORS:-false}"
STRICT_VALIDATOR_HEALTH="${STRICT_VALIDATOR_HEALTH:-false}"
CONNECT_TIMEOUT_SECONDS="${CONNECT_TIMEOUT_SECONDS:-5}"
REQUEST_TIMEOUT_SECONDS="${REQUEST_TIMEOUT_SECONDS:-30}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-}"
CSRF_TOKEN="${CSRF_TOKEN:-}"
TMP_DIR=""
OUT_REAL=""

usage() {
  cat <<'EOF'
Usage:
  scripts/validator-health.sh [--help] [--dry-run]

Environment:
  ADMIN_BASE_URL                  Private admin root for authenticated calls
  PUBLIC_BASE_URL                 Used as ADMIN_BASE_URL only when it is loopback
  ADMIN_TOKEN                     Optional admin bearer token; value is never printed
  CSRF_TOKEN                      Optional form token; bearer admin requests do not require browser CSRF
  RUN_VALIDATORS                  true|false; POST admin-only run_all when true
  STRICT_VALIDATOR_HEALTH         true|false; fail on blocker statuses when true
  OUTPUT_DIR                      Default .cache/validator-health/<timestamp>
  ALLOW_UNIGNORED_OUTPUT_DIR      true|false; allow output outside .cache except evidence-like paths
  FORCE                           true|false; allow non-empty OUTPUT_DIR reuse
  CONNECT_TIMEOUT_SECONDS         curl connect timeout, default 5
  REQUEST_TIMEOUT_SECONDS         curl total timeout, default 30

Safety:
  This helper writes private diagnostics only. It does not create evidence
  packets, write docs/evidence, contact consumers, submit feeds, edit GTFS,
  change consumer statuses, or change publish state.
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

prepare_output_dir() {
  OUT_REAL="$(python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$ALLOW_UNIGNORED_OUTPUT_DIR" <<'PY'
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out = pathlib.Path(sys.argv[2])
allow = sys.argv[3] == "true"
if not out.is_absolute():
    out = root / out
raw = str(out)
resolved = out.resolve(strict=False)
cache = (root / ".cache").resolve(strict=False)
evidence = (root / "docs" / "evidence").resolve(strict=False)
parts = {p.lower() for p in resolved.parts}
if "evidence" in parts or "docs/evidence" in raw.lower().replace("\\", "/"):
    raise SystemExit("OUTPUT_DIR must not be evidence-like or under docs/evidence")
try:
    resolved.relative_to(evidence)
    raise SystemExit("OUTPUT_DIR must not be under docs/evidence")
except ValueError:
    pass
if not allow:
    try:
        resolved.relative_to(cache)
    except ValueError:
        raise SystemExit(f"OUTPUT_DIR must resolve under {cache} unless ALLOW_UNIGNORED_OUTPUT_DIR=true")
print(resolved)
PY
)"
  if [ -L "$OUTPUT_DIR" ] || [ -L "$OUT_REAL" ]; then
    fail "OUTPUT_DIR must not be a symlink"
  fi
  if [ -d "$OUT_REAL" ] && [ "$FORCE" != "true" ] && [ "$(find "$OUT_REAL" -mindepth 1 -maxdepth 1 | sed -n '1p')" ]; then
    fail "OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it"
  fi
  mkdir -p "$OUT_REAL"
  chmod 700 "$OUT_REAL"
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-validator-health.XXXXXX")"
}

safe_admin_base() {
  python3 - "$ADMIN_BASE_URL" "$PUBLIC_BASE_URL" <<'PY'
import ipaddress
import sys
from urllib.parse import urlsplit, urlunsplit

admin, public = sys.argv[1], sys.argv[2]
source = "admin"
if not admin and public:
    parts = urlsplit(public)
    host = parts.hostname
    loopback = False
    if host in ("localhost",):
        loopback = True
    else:
        try:
            loopback = ipaddress.ip_address(host or "").is_loopback
        except ValueError:
            loopback = False
    if loopback:
        admin = public
        source = "public_loopback"
if not admin:
    raise SystemExit("missing")
parts = urlsplit(admin)
if parts.scheme not in ("http", "https"):
    raise SystemExit("ADMIN_BASE_URL must be http or https")
if not parts.hostname:
    raise SystemExit("ADMIN_BASE_URL must include a host")
if parts.username or parts.password or "@" in parts.netloc:
    raise SystemExit("ADMIN_BASE_URL must not include embedded credentials")
if parts.query or parts.fragment:
    raise SystemExit("ADMIN_BASE_URL must not include query strings or fragments")
host = parts.hostname
loopback = host == "localhost"
if not loopback:
    try:
        loopback = ipaddress.ip_address(host or "").is_loopback
    except ValueError:
        loopback = False
if not loopback and parts.scheme != "https":
    raise SystemExit("non-loopback ADMIN_BASE_URL must use https")
path = parts.path.rstrip("/")
normalized = urlunsplit((parts.scheme, parts.netloc, path, "", ""))
print(normalized)
print("loopback" if loopback else "non_loopback")
print(source)
PY
}

check_validator_tooling() {
  set +e
  VALIDATOR_TOOLING_MODE="${VALIDATOR_TOOLING_MODE:-pinned}" scripts/check-validators.sh >"$TMP_DIR/check-validators.out" 2>"$TMP_DIR/check-validators.err"
  rc=$?
  set -e
  mode="${VALIDATOR_TOOLING_MODE:-pinned}"
  case "$rc" in
    0)
      if [ "$mode" = "stub" ]; then
        printf '%s\n' "stub"
      else
        printf '%s\n' "configured"
      fi
      ;;
    11) printf '%s\n' "missing_tooling" ;;
    12) printf '%s\n' "misconfigured_tooling" ;;
    *) printf '%s\n' "blocked" ;;
  esac
}

fetch_admin_summary() {
  base="$1"
  method="$2"
  url="$base/admin/operations/validation-health.json"
  if [ "$method" = "POST" ]; then
    url="$base/admin/operations/validation-health"
    curl -sS --fail --connect-timeout "$CONNECT_TIMEOUT_SECONDS" --max-time "$REQUEST_TIMEOUT_SECONDS" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      --data-urlencode "action=run_all" \
      --data-urlencode "csrf_token=$CSRF_TOKEN" \
      "$url" >"$TMP_DIR/admin-response.raw" 2>"$TMP_DIR/admin-response.err"
  else
    curl -sS --fail --connect-timeout "$CONNECT_TIMEOUT_SECONDS" --max-time "$REQUEST_TIMEOUT_SECONDS" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      "$url" >"$TMP_DIR/admin-response.raw" 2>"$TMP_DIR/admin-response.err"
  fi
}

write_summary_from_admin_or_local() {
  tooling_status="$1"
  admin_status="$2"
  python3 - "$TMP_DIR/admin-response.raw" "$OUT_REAL/summary.json" "$tooling_status" "$admin_status" <<'PY'
import json
import pathlib
import sys
from datetime import datetime, timezone

raw_path, out_path, tooling_status, admin_status = sys.argv[1:5]
allowed_top = {
    "generated_at", "agency_id", "overall_status", "tooling_status", "feeds",
    "external_evidence_created", "consumer_statuses_changed", "compliance_claimed",
    "production_readiness_claimed",
}
allowed_row = {
    "feed_type", "validator_id", "validator_name", "tooling_status", "artifact_status",
    "latest_result_status", "latest_result_at", "active_feed_version_id",
    "latest_result_feed_version_id", "stale_status", "health_status", "next_action",
    "claim_boundary",
}
feeds = []
data = None
raw = pathlib.Path(raw_path)
if admin_status == "loaded" and raw.exists() and raw.stat().st_size:
    parsed = json.loads(raw.read_text())
    data = {k: parsed.get(k) for k in allowed_top}
    data["feeds"] = [{k: row.get(k) for k in allowed_row} for row in parsed.get("feeds", [])[:4]]
else:
    now = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    health = "skipped" if admin_status == "skipped" else "blocked"
    for feed_type, validator_id, validator_name in [
        ("schedule", "static-mobilitydata", "mobilitydata-gtfs-validator"),
        ("vehicle_positions", "realtime-mobilitydata", "mobilitydata-gtfs-realtime-validator"),
        ("trip_updates", "realtime-mobilitydata", "mobilitydata-gtfs-realtime-validator"),
        ("alerts", "realtime-mobilitydata", "mobilitydata-gtfs-realtime-validator"),
    ]:
        feeds.append({
            "feed_type": feed_type,
            "validator_id": validator_id,
            "validator_name": validator_name,
            "tooling_status": tooling_status,
            "artifact_status": "unknown",
            "latest_result_status": "not_run",
            "latest_result_at": None,
            "active_feed_version_id": "",
            "latest_result_feed_version_id": "",
            "stale_status": "unknown",
            "health_status": health,
            "next_action": "Supply ADMIN_TOKEN and a safe ADMIN_BASE_URL to check private admin validator health.",
            "claim_boundary": "Private diagnostics only; no evidence packet, consumer status change, compliance claim, acceptance claim, or production readiness claim.",
        })
    data = {
        "generated_at": now,
        "agency_id": "",
        "overall_status": health,
        "tooling_status": tooling_status,
        "feeds": feeds,
        "external_evidence_created": False,
        "consumer_statuses_changed": False,
        "compliance_claimed": False,
        "production_readiness_claimed": False,
    }
data["external_evidence_created"] = False
data["consumer_statuses_changed"] = False
data["compliance_claimed"] = False
data["production_readiness_claimed"] = False
path = pathlib.Path(out_path)
path.write_text(json.dumps(data, indent=2, sort_keys=False) + "\n")
PY
}

write_markdown_outputs() {
  python3 - "$OUT_REAL/summary.json" "$OUT_REAL/summary.md" "$OUT_REAL/manifest.json" "$OUT_REAL/manifest.md" "$OUTPUT_DIR" <<'PY'
import json
import pathlib
import sys
from datetime import datetime, timezone

summary_path, summary_md_path, manifest_json_path, manifest_md_path, output_dir = sys.argv[1:6]
output_dir_value = output_dir
if pathlib.Path(output_dir).is_absolute():
    output_dir_value = pathlib.Path(output_dir).name
summary = json.loads(pathlib.Path(summary_path).read_text())
lines = [
    "# Validator Health Summary",
    "",
    "Private diagnostics only. This is not evidence, not a consumer submission artifact, not a compliance claim, and not production readiness proof.",
    "",
    f"- generated_at: `{summary.get('generated_at', '')}`",
    f"- agency_id: `{summary.get('agency_id', '')}`",
    f"- overall_status: `{summary.get('overall_status', '')}`",
    f"- tooling_status: `{summary.get('tooling_status', '')}`",
    "",
    "| feed_type | validator_id | tooling_status | artifact_status | latest_result_status | stale_status | health_status |",
    "| --- | --- | --- | --- | --- | --- | --- |",
]
for row in summary.get("feeds", []):
    lines.append("| {feed_type} | {validator_id} | {tooling_status} | {artifact_status} | {latest_result_status} | {stale_status} | {health_status} |".format(**row))
lines.extend([
    "",
    f"- external_evidence_created: `{summary.get('external_evidence_created')}`",
    f"- consumer_statuses_changed: `{summary.get('consumer_statuses_changed')}`",
    f"- compliance_claimed: `{summary.get('compliance_claimed')}`",
    f"- production_readiness_claimed: `{summary.get('production_readiness_claimed')}`",
])
pathlib.Path(summary_md_path).write_text("\n".join(lines) + "\n")
manifest = {
    "created_at": datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z"),
    "output_dir": output_dir_value,
    "included_files": ["summary.json", "summary.md", "manifest.json", "manifest.md"],
    "excluded_categories": [
        "raw_validator_reports", "auth_headers", "cookies", "tokens", "database_urls",
        "private_paths", "evidence_packets", "consumer_submission_artifacts",
    ],
}
pathlib.Path(manifest_json_path).write_text(json.dumps(manifest, indent=2) + "\n")
pathlib.Path(manifest_md_path).write_text("""# Validator Health Manifest

Private diagnostics only.

Included files:
- summary.json
- summary.md
- manifest.json
- manifest.md

Excluded categories:
- raw validator reports
- Authorization headers
- cookies
- tokens
- database URLs
- private paths
- evidence packets
- consumer submission artifacts
""")
PY
}

redaction_scan() {
  python3 - "$OUT_REAL" <<'PY'
import pathlib
import re
import sys

root = pathlib.Path(sys.argv[1])
patterns = [
    ("authorization_header", re.compile(r"Authorization\s*:\s*", re.I)),
    ("bearer_secret", re.compile(r"Bearer\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("cookie_header", re.compile(r"Cookie\s*:\s*", re.I)),
    ("admin_session", re.compile(r"admin_session\s*=", re.I)),
    ("csrf_secret", re.compile(r"csrf[_-]?(token|secret)?\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}", re.I)),
    ("database_url", re.compile(r"\bDATABASE_URL\b|postgres(?:ql)?://[^\\s\"']+:[^\\s\"']+@", re.I)),
    ("private_key", re.compile(r"BEGIN (?:RSA |EC |OPENSSH |)PRIVATE KEY", re.I)),
    ("private_path", re.compile(r"(/Users/[^\\s\"']+|/tmp/[^\\s\"']+|/var/folders/[^\\s\"']+)", re.I)),
    ("token_env", re.compile(r"\bTOKEN\s*=", re.I)),
    ("secret_env", re.compile(r"\bSECRET\s*=", re.I)),
    ("password_env", re.compile(r"\bPASSWORD\s*=", re.I)),
]
for path in root.rglob("*"):
    if not path.is_file():
        continue
    text = path.read_text(errors="ignore")
    for name, pattern in patterns:
        if pattern.search(text):
            raise SystemExit(f"redaction scan failed: {name} in {path.name}")
PY
}

strict_status_check() {
  python3 - "$OUT_REAL/summary.json" <<'PY'
import json
import sys
bad = {"blocked", "failed", "missing_tooling", "misconfigured_tooling", "artifact_unavailable", "stale"}
summary = json.load(open(sys.argv[1]))
statuses = {summary.get("overall_status"), summary.get("tooling_status")}
statuses.update(row.get("health_status") for row in summary.get("feeds", []))
statuses.update(row.get("tooling_status") for row in summary.get("feeds", []))
raise SystemExit(1 if statuses & bad else 0)
PY
}

main() {
  dry_run="false"
  case "${1:-}" in
    --help|-h) usage; exit 0 ;;
    --dry-run) dry_run="true" ;;
    "") ;;
    *) fail "unknown argument: $1" ;;
  esac
  bool_var RUN_VALIDATORS "$RUN_VALIDATORS"
  bool_var STRICT_VALIDATOR_HEALTH "$STRICT_VALIDATOR_HEALTH"
  bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
  bool_var FORCE "$FORCE"
  positive_int CONNECT_TIMEOUT_SECONDS "$CONNECT_TIMEOUT_SECONDS"
  positive_int REQUEST_TIMEOUT_SECONDS "$REQUEST_TIMEOUT_SECONDS"
  prepare_output_dir
  tooling_status="$(check_validator_tooling)"
  admin_status="skipped"
  if [ "$dry_run" != "true" ] && [ -n "$ADMIN_TOKEN" ]; then
    if ! command -v curl >/dev/null 2>&1; then
      fail "curl is required for authenticated admin requests"
    fi
    safe="$(safe_admin_base)" || fail "$safe"
    admin_base="$(printf '%s\n' "$safe" | sed -n '1p')"
    method="GET"
    if [ "$RUN_VALIDATORS" = "true" ]; then
      method="POST"
    fi
    set +e
    fetch_admin_summary "$admin_base" "$method"
    rc=$?
    set -e
    if [ "$rc" -eq 0 ]; then
      admin_status="loaded"
    else
      admin_status="blocked"
    fi
  elif [ "$dry_run" = "true" ]; then
    admin_status="skipped"
  elif [ "$STRICT_VALIDATOR_HEALTH" = "true" ]; then
    admin_status="blocked"
  fi
  write_summary_from_admin_or_local "$tooling_status" "$admin_status"
  write_markdown_outputs
  python3 -m json.tool "$OUT_REAL/summary.json" >/dev/null
  python3 -m json.tool "$OUT_REAL/manifest.json" >/dev/null
  redaction_scan
  display_output="$(python3 - "$OUTPUT_DIR" <<'PY'
import pathlib
import sys
p = pathlib.Path(sys.argv[1])
print(p.name if p.is_absolute() else sys.argv[1])
PY
)"
  printf 'validator health summary written to %s\n' "$display_output"
  if [ "$STRICT_VALIDATOR_HEALTH" = "true" ]; then
    strict_status_check
  fi
}

main "$@"
