#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/oci-reference-check/$TIMESTAMP}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-}"
FORCE="${FORCE:-false}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
STRICT="${STRICT:-false}"
DRY_RUN="false"
SKIP_PUBLIC_FETCH="${SKIP_PUBLIC_FETCH:-false}"
OCI_USER="${OCI_USER:-opc}"
OCI_HOST="${OCI_HOST:-}"
OCI_KEY="${OCI_KEY:-}"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/oci-reference-check.sh [--public-base-url URL] [--output-dir DIR] [--dry-run] [--strict] [--skip-public-fetch]

Environment:
  OUTPUT_DIR                   Default .cache/oci-reference-check/<UTC timestamp>
  PUBLIC_BASE_URL              Public root for five-feed fetches
  FORCE                        true|false; allow non-empty output reuse
  ALLOW_UNIGNORED_OUTPUT_DIR   true|false; allow output outside .cache except evidence-like paths
  STRICT                       true|false; fail on public fetch/validator or SSH loopback blockers
  OCI_HOST, OCI_USER, OCI_KEY  Optional SSH access for loopback health checks

Safety:
  This helper creates private reference-deployment diagnostics only. It does
  not write docs/evidence, does not submit feeds, does not contact consumers,
  does not change consumer statuses, does not require Go on the remote host,
  and does not print secrets, tokens, private keys, database URLs, or populated
  environment values.
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
    --skip-public-fetch)
      SKIP_PUBLIC_FETCH="true"
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
bool_var SKIP_PUBLIC_FETCH "$SKIP_PUBLIC_FETCH"

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
TMP_DIR="$(mktemp -d ".cache/oci-reference-check-tmp.XXXXXX")"
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
ENV_TSV="$TMP_DIR/env.tsv"
SSH_TSV="$TMP_DIR/ssh.tsv"
VALIDATE_SUMMARY="$OUT_REAL/public-feeds/summary.json"
printf 'id\tlabel\tstatus\tcurrent_signal\tnext_action\n' >"$ENV_TSV"
printf 'id\tlabel\tstatus\thttp_status\tnext_action\n' >"$SSH_TSV"

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

presence_row() {
  id="$1"
  label="$2"
  shift 2
  present=""
  missing=""
  for name in "$@"; do
    value="$(eval "printf '%s' \"\${$name:-}\"")"
    if [ -n "$value" ]; then
      present="${present}${present:+, }${name}=configured"
    else
      missing="${missing}${missing:+, }${name}=not configured"
    fi
  done
  if [ -n "$present" ]; then
    printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "configured" "$present; values withheld" "Keep values in the operator environment and out of reports." >>"$ENV_TSV"
  else
    printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "not_configured" "$missing" "Configure only when this reference diagnostic needs that capability." >>"$ENV_TSV"
  fi
}

presence_row "deployment_target" "Deployment helper target" ENVIRONMENT_NAME OCI_HOST OCI_REMOTE_DIR DOMAIN PUBLIC_BASE_URL
presence_row "ssh_access" "Optional SSH loopback access" OCI_HOST OCI_KEY
presence_row "backup_configuration" "Backup configuration" BACKUP_DIR BACKUP_PATH OPEN_TRANSIT_BACKUP_DIR
presence_row "restore_drill_configuration" "Restore-drill configuration" RESTORE_DATABASE_URL RESTORE_BACKUP_FILE RESTORE_DRILL_DATABASE_URL RESTORE_DRILL_BACKUP_FILE RESTORE_DRILL_TARGET OPEN_TRANSIT_RESTORE_DRILL
presence_row "telemetry_credentials" "Telemetry simulator credential presence" DEVICE_TOKEN DEVICE_TOKEN_FILE ADMIN_TOKEN

if [ "$SKIP_PUBLIC_FETCH" = "true" ]; then
  mkdir -p "$OUT_REAL/public-feeds"
  python3 - "$OUT_REAL/public-feeds" "$TIMESTAMP" <<'PY'
import json
import pathlib
import sys
out = pathlib.Path(sys.argv[1])
generated_at = sys.argv[2]
flags = {
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
summary = {"generated_at": generated_at, "dry_run": True, "skip_validators": True, "rows": [], "counts": {"rows": 0, "fetch_blockers": 0, "validator_blockers": 0}, "claim_flags": flags, "boundary": "Public fetches skipped by operator request.", "strict_failed": False}
(out / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
(out / "summary.md").write_text("# Public Feed Checks Skipped\n")
PY
else
  validate_args="--output-dir"
  if [ "$DRY_RUN" = "true" ]; then
    DRY_ARG="--dry-run"
  else
    DRY_ARG=""
  fi
  if [ "$STRICT" = "true" ]; then
    STRICT_ARG="--strict"
  else
    STRICT_ARG=""
  fi
  PUBLIC_BASE_URL="$PUBLIC_BASE_URL" FORCE=true scripts/validate-public-feeds.sh $validate_args "$OUTPUT_DIR/public-feeds" $DRY_ARG $STRICT_ARG >/dev/null
fi

ssh_run() {
  if [ -n "$OCI_KEY" ]; then
    ssh -o BatchMode=yes -o ConnectTimeout=8 -i "$OCI_KEY" "${OCI_USER}@${OCI_HOST}" "$@"
  else
    ssh -o BatchMode=yes -o ConnectTimeout=8 "${OCI_USER}@${OCI_HOST}" "$@"
  fi
}

loopback_check() {
  id="$1"
  label="$2"
  port="$3"
  if [ "$DRY_RUN" = "true" ]; then
    printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "not_checked" "not checked" "Dry run only. Rerun without --dry-run when SSH access is available." >>"$SSH_TSV"
    return
  fi
  if [ -z "$OCI_HOST" ]; then
    printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "not_configured" "not checked" "Set OCI_HOST and optional OCI_KEY for loopback health checks." >>"$SSH_TSV"
    return
  fi
  code="000"
  if code="$(ssh_run "curl -fsS -o /dev/null -w '%{http_code}' http://127.0.0.1:${port}/healthz" 2>"$OUT_REAL/ssh-$id-error.txt")"; then
    if [ "$code" = "200" ]; then
      printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "healthy" "$code" "No action from this diagnostic row." >>"$SSH_TSV"
    else
      printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "needs_review" "$code" "Review the service unit and loopback listener." >>"$SSH_TSV"
    fi
  else
    printf '%s\t%s\t%s\t%s\t%s\n' "$id" "$label" "blocked" "ssh failed" "Confirm SSH access and service status; no remote Go toolchain is required." >>"$SSH_TSV"
  fi
}

loopback_check "agency_config" "agency-config loopback health" "8081"
loopback_check "telemetry_ingest" "telemetry-ingest loopback health" "8082"
loopback_check "vehicle_positions" "feed-vehicle-positions loopback health" "8083"
loopback_check "trip_updates" "feed-trip-updates loopback health" "8084"
loopback_check "alerts" "feed-alerts loopback health" "8085"

python3 - "$ROOT_DIR" "$OUT_REAL" "$ENV_TSV" "$SSH_TSV" "$VALIDATE_SUMMARY" "$TIMESTAMP" "$DRY_RUN" "$STRICT" <<'PY'
import csv
import json
import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out = pathlib.Path(sys.argv[2]).resolve()
env_rows = list(csv.DictReader(open(sys.argv[3]), delimiter="\t"))
ssh_rows = list(csv.DictReader(open(sys.argv[4]), delimiter="\t"))
feed_summary_path = pathlib.Path(sys.argv[5])
generated_at = sys.argv[6]
dry_run = sys.argv[7] == "true"
strict = sys.argv[8] == "true"
feed_summary = json.loads(feed_summary_path.read_text()) if feed_summary_path.exists() else {"rows": [], "counts": {}, "strict_failed": False, "claim_flags": {}}
flags = {
    "external_evidence_created": False,
    "consumer_statuses_changed": False,
    "compliance_claimed": False,
    "production_readiness_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "public_launch_claimed": False,
    "hosted_saas_claimed": False,
    "vendor_compatibility_claimed": False,
    "hardware_certification_claimed": False,
    "sla_claimed": False,
    "uptime_guarantee_claimed": False,
    "production_grade_eta_claimed": False,
}
ssh_blockers = [r for r in ssh_rows if r["status"] == "blocked"]
env_missing = [r for r in env_rows if r["status"] == "not_configured"]
strict_failed = bool(strict and (feed_summary.get("strict_failed") or ssh_blockers))
summary = {
    "generated_at": generated_at,
    "dry_run": dry_run,
    "strict_mode": strict,
    "boundary": "Private OCI/reference deployment diagnostics only. This is not evidence, final-root proof, consumer submission, compliance, agency adoption, hosted service availability, SLA coverage, vendor compatibility, production readiness, or production-grade ETA proof.",
    "public_feed_summary": feed_summary,
    "deployment_helper_status": env_rows,
    "loopback_health": ssh_rows,
    "backup_restore": [r for r in env_rows if r["id"] in {"backup_configuration", "restore_drill_configuration"}],
    "telemetry_simulator_guidance": {
        "status": "operator_guided",
        "dry_run_command": "scripts/telemetry-simulator.sh --list-scenarios; OUTPUT_DIR=.cache/telemetry-simulator DRY_RUN=true scripts/telemetry-simulator.sh --scenario on-route --dry-run --force",
        "next_action": "Keep device tokens in the operator shell. Use the private Devices and Telemetry Simulator pages to decide whether a technical helper is needed.",
        "does_not_prove": "Does not prove real device reliability, vendor compatibility, or production-grade ETA quality.",
    },
    "counts": {
        "public_rows": len(feed_summary.get("rows", [])),
        "env_missing": len(env_missing),
        "ssh_rows": len(ssh_rows),
        "ssh_blockers": len(ssh_blockers),
    },
    "claim_flags": flags,
    "strict_failed": strict_failed,
}
manifest = {
    "generated_at": generated_at,
    "output_dir": out.relative_to(root).as_posix() if out.is_relative_to(root) else out.name,
    "included_files": ["summary.json", "summary.md", "manifest.json", "manifest.md", "public-feeds/"],
    "excluded_categories": ["docs/evidence", "consumer submissions", "secret values", "raw tokens", "private keys", "database URLs", "remote Go toolchain"],
    "claim_flags": flags,
}
(out / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n")
(out / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n")
lines = [
    "# OCI Reference Check Summary",
    "",
    f"Generated at: {generated_at}",
    "",
    "This is a private reference diagnostic. It creates no evidence and makes no compliance, consumer, adoption, hosted-service, SLA, vendor, production, or ETA-quality claim.",
    "",
    "## Public Feed Rows",
    "",
    "| Feed | Path | Fetch | Validator | Next action |",
    "| --- | --- | --- | --- | --- |",
]
for row in feed_summary.get("rows", []):
    lines.append(f"| {row['label']} | `{row['public_path']}` | {row['fetch_status']} | {row['validator_state']} | {row['next_action']} |")
lines.extend(["", "## Deployment Presence", "", "| Area | Status | Signal | Next action |", "| --- | --- | --- | --- |"])
for row in env_rows:
    lines.append(f"| {row['label']} | {row['status']} | {row['current_signal']} | {row['next_action']} |")
lines.extend(["", "## Loopback Health", "", "| Service | Status | HTTP | Next action |", "| --- | --- | --- | --- |"])
for row in ssh_rows:
    lines.append(f"| {row['label']} | {row['status']} | {row['http_status']} | {row['next_action']} |")
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
  printf 'oci-reference-check summary written to %s\n' "$(rel_output_file "$OUT_REAL")"
else
  printf 'oci-reference-check found strict blockers; summary written to %s\n' "$(rel_output_file "$OUT_REAL")" >&2
  exit 1
fi
