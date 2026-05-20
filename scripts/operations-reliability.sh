#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/operations-reliability/$TIMESTAMP}"
FORCE="${FORCE:-false}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
ALLOW_UNIGNORED_SOURCE_DIR="${ALLOW_UNIGNORED_SOURCE_DIR:-false}"
MAX_SOURCE_BYTES="${MAX_SOURCE_BYTES:-5242880}"
VALIDATOR_HEALTH_SUMMARY="${VALIDATOR_HEALTH_SUMMARY:-}"
DEPLOYMENT_DOCTOR_SUMMARY="${DEPLOYMENT_DOCTOR_SUMMARY:-}"
OPERATIONS_NOTIFY_SUMMARY="${OPERATIONS_NOTIFY_SUMMARY:-}"
DIAGNOSTIC_CACHE_ROOT="${DIAGNOSTIC_CACHE_ROOT:-.cache}"
DRY_RUN="false"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/operations-reliability.sh [--help] [--dry-run]

Environment:
  OUTPUT_DIR                         Default .cache/operations-reliability/<UTC timestamp>
  FORCE                              true|false; allow non-empty output reuse
  ALLOW_UNIGNORED_OUTPUT_DIR         true|false; allow output outside .cache except evidence-like paths
  ALLOW_UNIGNORED_SOURCE_DIR         true|false; allow explicit source summaries outside .cache except evidence-like paths
  MAX_SOURCE_BYTES                   Maximum bytes per source summary, default 5242880
  VALIDATOR_HEALTH_SUMMARY           Optional explicit .cache/validator-health/.../summary.json
  DEPLOYMENT_DOCTOR_SUMMARY          Optional explicit .cache/deployment-doctor/.../summary.json
  OPERATIONS_NOTIFY_SUMMARY          Optional explicit .cache/operations-notify/.../summary.json
  DIAGNOSTIC_CACHE_ROOT              Optional source discovery cache root, default .cache

Safety:
  This helper creates private local reliability diagnostics only. It writes
  exactly summary.json, summary.md, manifest.json, manifest.md, and
  reliability-review.txt. It never sends notifications, calls mutating admin
  endpoints, writes docs/evidence, changes consumer statuses, or claims
  compliance, production readiness, SLA coverage, uptime guarantees, hosted
  SaaS availability, agency adoption, consumer acceptance, vendor compatibility,
  or production-grade ETA quality.
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

positive_int() {
  case "$2" in
    ''|*[!0-9]*) fail "$1 must be a positive integer" ;;
    0) fail "$1 must be greater than zero" ;;
  esac
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --dry-run)
      DRY_RUN="true"
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

bool_var FORCE "$FORCE"
bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
bool_var ALLOW_UNIGNORED_SOURCE_DIR "$ALLOW_UNIGNORED_SOURCE_DIR"
positive_int MAX_SOURCE_BYTES "$MAX_SOURCE_BYTES"

TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-operations-reliability.XXXXXX")"

python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$TIMESTAMP" "$FORCE" "$ALLOW_UNIGNORED_OUTPUT_DIR" "$ALLOW_UNIGNORED_SOURCE_DIR" "$MAX_SOURCE_BYTES" "$VALIDATOR_HEALTH_SUMMARY" "$DEPLOYMENT_DOCTOR_SUMMARY" "$OPERATIONS_NOTIFY_SUMMARY" "$DIAGNOSTIC_CACHE_ROOT" "$DRY_RUN" <<'PY'
import json
import os
import pathlib
import re
import shutil
import sys

(
    root_arg,
    output_arg,
    timestamp,
    force_arg,
    allow_output_arg,
    allow_source_arg,
    max_source_bytes_arg,
    validator_arg,
    doctor_arg,
    notify_arg,
    diagnostic_cache_arg,
    dry_run_arg,
) = sys.argv[1:13]

ROOT = pathlib.Path(root_arg).resolve()
OUTPUT_RAW = pathlib.Path(output_arg)
FORCE = force_arg == "true"
ALLOW_OUTPUT = allow_output_arg == "true"
ALLOW_SOURCE = allow_source_arg == "true"
MAX_SOURCE_BYTES = int(max_source_bytes_arg)
DRY_RUN = dry_run_arg == "true"

OUTPUT_FILES = ("summary.json", "summary.md", "manifest.json", "manifest.md", "reliability-review.txt")
TIMESTAMP_DIR = re.compile(r"^\d{8}T\d{6}Z$")
MAX_TEXT = 220

CLAIM_FLAGS = {
    "external_evidence_created": False,
    "final_root_evidence_created": False,
    "consumer_statuses_changed": False,
    "compliance_claimed": False,
    "production_readiness_claimed": False,
    "sla_claimed": False,
    "uptime_guarantee_claimed": False,
    "hosted_saas_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "vendor_compatibility_claimed": False,
    "production_grade_eta_claimed": False,
}


def fail(message):
    raise SystemExit(message)


def path_has_symlink(path):
    probe = pathlib.Path(path)
    if not probe.is_absolute():
        probe = ROOT / probe
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False


def is_evidence_like(path):
    raw = str(path).replace("\\", "/").lower()
    parts = [p.lower() for p in pathlib.Path(path).parts]
    return "docs/evidence" in raw or "evidence" in parts or "proof" in parts or "submission" in parts


def resolve_diagnostic_cache_root():
    raw = pathlib.Path(diagnostic_cache_arg)
    path = raw if raw.is_absolute() else ROOT / raw
    resolved = path.resolve(strict=False)
    if is_evidence_like(path) or is_evidence_like(resolved):
        fail("DIAGNOSTIC_CACHE_ROOT must not be evidence-like or under docs/evidence")
    if path_has_symlink(path):
        fail("DIAGNOSTIC_CACHE_ROOT must not contain symlink directories")
    return resolved


DIAGNOSTIC_CACHE = resolve_diagnostic_cache_root()


def rel_to_root(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<redacted-source>"


def resolve_output_dir():
    out = OUTPUT_RAW if OUTPUT_RAW.is_absolute() else ROOT / OUTPUT_RAW
    resolved = out.resolve(strict=False)
    cache = (ROOT / ".cache").resolve(strict=False)
    if is_evidence_like(out) or is_evidence_like(resolved):
        fail("OUTPUT_DIR must not be evidence-like or under docs/evidence")
    if path_has_symlink(out):
        fail("OUTPUT_DIR must not contain symlink directories")
    if not ALLOW_OUTPUT:
        try:
            resolved.relative_to(cache)
        except ValueError:
            fail("OUTPUT_DIR must resolve under repo .cache unless ALLOW_UNIGNORED_OUTPUT_DIR=true")
    if resolved.exists() and not resolved.is_dir():
        fail("OUTPUT_DIR must be a directory")
    if resolved.exists() and any(resolved.iterdir()):
        if not FORCE:
            fail("OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it")
        for child in resolved.iterdir():
            if child.is_symlink() or child.is_file():
                child.unlink()
            else:
                shutil.rmtree(child)
    resolved.mkdir(parents=True, exist_ok=True)
    os.chmod(resolved, 0o700)
    return resolved


def discover_latest(kind):
    base = DIAGNOSTIC_CACHE / kind
    if not base.exists() or path_has_symlink(base):
        return None
    candidates = []
    for child in base.iterdir():
        if not child.is_dir() or child.is_symlink() or not TIMESTAMP_DIR.match(child.name):
            continue
        summary = child / "summary.json"
        if summary.exists() and not summary.is_symlink():
            candidates.append((child.name, summary))
    if not candidates:
        return None
    return sorted(candidates, key=lambda item: item[0])[-1][1]


def source_allowed(path):
    if is_evidence_like(path) or path_has_symlink(path):
        return False
    resolved = path.resolve(strict=False)
    if is_evidence_like(resolved):
        return False
    if ALLOW_SOURCE:
        return True
    cache = (ROOT / ".cache").resolve(strict=False)
    try:
        resolved.relative_to(cache)
        return True
    except ValueError:
        return False


def unsafe_source_name(path):
    lower = path.name.lower()
    return lower.endswith((".log", ".sql", ".dump", ".bak", ".backup", ".gz", ".zip")) or lower in {"payload.json", "raw.json", "debug.json"}


def unsafe_text(text):
    lower = text.lower()
    patterns = [
        "postgres://",
        "database_url=",
        "restore_database_url=",
        "authorization:",
        "bearer ",
        "cookie:",
        "admin_session",
        "webhook_url=",
        "https://hooks.",
        "secret=",
        "secret:",
        "\"secret\"",
        "password",
        "token_hash",
        "private_key",
        "begin private key",
        "payload_json",
        "details_json",
        "raw_log",
        "backup_dump",
        "/users/",
        "/etc/",
        "/var/lib/",
    ]
    return any(pattern in lower for pattern in patterns)


def safe_string(value, default="unknown"):
    if value is None or isinstance(value, (dict, list)):
        return default
    text = str(value).replace("\r", " ").replace("\n", " ").strip()
    text = re.sub(r"\s+", " ", text)
    if unsafe_text(text):
        return "<redacted>"
    if len(text) > MAX_TEXT:
        text = text[: MAX_TEXT - 15].rstrip() + " [truncated]"
    return text or default


def read_source(kind, explicit):
    path = pathlib.Path(explicit) if explicit else discover_latest(kind)
    if path is None:
        return {
            "kind": kind,
            "status": "missing",
            "source": "not_found",
            "summary": "No safe private summary was found.",
            "next_action": f"Run {kind} first, then rerun operations-reliability.",
        }
    path = path if path.is_absolute() else ROOT / path
    if unsafe_source_name(path):
        fail(f"{kind} source must be a safe summary.json, not raw logs or backup dumps")
    if not source_allowed(path):
        fail(f"{kind} source path is not allowed")
    if not path.exists() or not path.is_file() or path.is_symlink():
        return {
            "kind": kind,
            "status": "missing",
            "source": rel_to_root(path),
            "summary": "The configured safe summary path is absent.",
            "next_action": f"Regenerate {kind} before reviewing reliability.",
        }
    size = path.stat().st_size
    if size > MAX_SOURCE_BYTES:
        fail(f"{kind} source exceeds MAX_SOURCE_BYTES")
    text = path.read_text(encoding="utf-8", errors="replace")
    if unsafe_text(text):
        fail(f"{kind} source contains raw logs, private payloads, secrets, DB URLs, or webhook values")
    try:
        data = json.loads(text)
    except json.JSONDecodeError:
        fail(f"{kind} source is not valid JSON")
    status = map_status(data.get("overall_status") or data.get("status"))
    return {
        "kind": kind,
        "status": status,
        "source": rel_to_root(path),
        "summary": safe_string(data.get("summary") or data.get("overall_status") or data.get("status"), "Safe summary loaded."),
        "next_action": safe_string(data.get("next_action"), "Review this safe private summary."),
        "data": data,
    }


def map_status(value):
    raw = safe_string(value, "unknown").lower()
    if raw in {"ok", "passed", "recorded", "info", "runnable", "configured"}:
        return "ok"
    if raw in {"missing", "not_found", "not_run"}:
        return "missing"
    if raw in {"failed", "blocked", "unhealthy", "error"}:
        return "unhealthy"
    if raw in {"warning", "warnings", "needs_review", "stale", "degraded"}:
        return "needs_review"
    return "unknown"


def worse(a, b):
    rank = {"ok": 0, "unknown": 1, "missing": 2, "needs_review": 3, "unhealthy": 4}
    return b if rank.get(b, 1) > rank.get(a, 1) else a


def section_from_source(source, label):
    return {
        "status": source["status"],
        "source": source["source"],
        "summary": source["summary"],
        "next_action": source["next_action"],
        "section": label,
    }


def notify_bool(data, path):
    cursor = data
    for part in path:
        if not isinstance(cursor, dict):
            return False
        cursor = cursor.get(part)
    return bool(cursor)


def build_monitoring_export_section(source):
    data = source.get("data") if isinstance(source.get("data"), dict) else {}
    digest = data.get("health_digest") if isinstance(data.get("health_digest"), dict) else {}
    channel = data.get("channel_guidance") if isinstance(data.get("channel_guidance"), dict) else {}
    notification = data.get("notification") if isinstance(data.get("notification"), dict) else {}
    return {
        "status": source["status"],
        "source": source["source"],
        "summary": safe_string(digest.get("source_summary"), source["summary"]),
        "action_summary": safe_string(digest.get("action_summary"), "not available"),
        "template": safe_string(digest.get("template"), "severity + source counts + no-send boundary"),
        "next_action": safe_string(digest.get("next_action"), source["next_action"]),
        "not_sent": bool(notification.get("not_sent", True)),
        "webhook_present": notify_bool(channel, ["webhook", "present"]) or notify_bool(data, ["destinations", "webhook_present"]),
        "email_present": notify_bool(channel, ["email", "present"]) or notify_bool(data, ["destinations", "email_present"]),
        "webhook_send_enabled": notify_bool(channel, ["webhook", "send_enabled"]),
        "email_send_enabled": notify_bool(channel, ["email", "send_enabled"]),
        "destination_values_recorded": notify_bool(channel, ["webhook", "destination_value_recorded"]) or notify_bool(channel, ["email", "destination_value_recorded"]),
        "does_not_prove": "Does not prove notification delivery, uptime, SLA coverage, hosted service availability, production readiness, compliance, consumer acceptance, or public launch.",
    }


def build_private_ops_summary(overall, validator, doctor, notify, monitoring_export):
    return {
        "status": overall,
        "summary_json": "summary.json",
        "manifest_json": "manifest.json",
        "scope": "private diagnostic summary only",
        "export_formats": {
            "feed_health": "status, source, next_action rows only",
            "connector_health": "category status, redaction state, no-send blockers only",
            "validator_posture": "validator id, feed type, health status, stale status only",
            "telemetry_freshness": "freshness bucket counts and stale/unmatched summaries only",
            "maintenance_tasks": "cadence, owner category, status, next step only",
        },
        "redaction_boundary": "no endpoint values, tokens, private paths, raw payloads, DB URLs, notification destinations, or retained evidence paths",
        "sources": {
            "validator_health": validator["status"],
            "deployment_doctor": doctor["status"],
            "operations_notify": notify["status"],
        },
        "monitoring_export_status": monitoring_export["status"],
        "notification_not_sent": bool(monitoring_export["not_sent"]),
        "next_action": "Use this private ops summary JSON for local review only; keep live delivery and evidence collection separately authorized.",
        "does_not_prove": "Does not prove hosted monitoring, uptime, SLA coverage, compliance, consumer acceptance, or production readiness.",
    }


out = resolve_output_dir()
validator = read_source("validator-health", validator_arg)
doctor = read_source("deployment-doctor", doctor_arg)
notify = read_source("operations-notify", notify_arg)

feeds = [
    {"feed_type": "schedule", "status": validator["status"], "source": validator["source"], "next_action": "Review validator-health schedule row; missing data is not ok."},
    {"feed_type": "vehicle_positions", "status": validator["status"], "source": validator["source"], "next_action": "Review validator-health and runtime Vehicle Positions health snapshots."},
    {"feed_type": "trip_updates", "status": validator["status"], "source": validator["source"], "next_action": "Review validator-health and Trip Updates diagnostics."},
    {"feed_type": "alerts", "status": validator["status"], "source": validator["source"], "next_action": "Review validator-health and Alerts diagnostics."},
]
incidents = {
    "status": "missing",
    "source": "runtime incident table not read by script",
    "total": 0,
    "counts_by_status": {},
    "counts_by_severity": {},
    "counts_by_type": {},
    "recent": [],
    "recent_limit": 0,
    "next_action": "Use the private admin reliability route for DB-backed incident rollups.",
}
backup_restore = section_from_source(doctor, "backup_restore")
alerting = section_from_source(notify, "alerting")
availability = section_from_source(doctor, "availability_sampling")
long_running = section_from_source(notify, "long_running_operations")
monitoring_export = build_monitoring_export_section(notify)

overall = "unknown"
for row in feeds:
    overall = worse(overall, row["status"])
for section in (incidents, backup_restore, alerting, availability, long_running, monitoring_export):
    overall = worse(overall, section["status"])
private_ops_summary = build_private_ops_summary(overall, validator, doctor, notify, monitoring_export)

summary = {
    "generated_at": timestamp,
    "overall_status": overall,
    "feeds": feeds,
    "incidents": incidents,
    "backup_restore": backup_restore,
    "alerting": alerting,
    "availability_sampling": availability,
    "long_running_operations": long_running,
    "monitoring_export": monitoring_export,
    "private_ops_summary": private_ops_summary,
    "claim_flags": CLAIM_FLAGS,
    "dry_run": DRY_RUN,
}

manifest = {
    "created_at": timestamp,
    "output_dir": rel_to_root(out),
    "included_files": list(OUTPUT_FILES),
    "source_summaries": [
        {"kind": validator["kind"], "status": validator["status"], "source": validator["source"]},
        {"kind": doctor["kind"], "status": doctor["status"], "source": doctor["source"]},
        {"kind": notify["kind"], "status": notify["status"], "source": notify["source"]},
    ],
    "excluded_categories": [
        "docs/evidence",
        "consumer status changes",
        "raw logs",
        "backup dumps",
        "database URLs",
        "webhook values",
        "secrets",
        "private payloads",
        "mutating admin endpoints",
        "notifications",
    ],
    "claim_flags": CLAIM_FLAGS,
}

summary_md = [
    "# Operations Reliability",
    "",
    f"- generated_at: {timestamp}",
    f"- overall_status: {overall}",
    "- scope: private diagnostic sampling only",
    "- claim_flags: all false",
    "",
    "## Feeds",
]
for row in feeds:
    summary_md.append(f"- {row['feed_type']}: {row['status']} ({row['source']})")
summary_md.extend([
    "",
    "## Sections",
    f"- incidents: {incidents['status']}",
    f"- backup_restore: {backup_restore['status']}",
    f"- alerting: {alerting['status']}",
    f"- availability_sampling: {availability['status']}",
    f"- long_running_operations: {long_running['status']}",
    f"- monitoring_export: {monitoring_export['status']} (not_sent={str(monitoring_export['not_sent']).lower()})",
    "",
    "## Private Ops Summary",
    f"- summary_json: {private_ops_summary['summary_json']}",
    f"- monitoring_export_status: {private_ops_summary['monitoring_export_status']}",
    "- notification_not_sent: true",
    "- export_formats: feed_health, connector_health, validator_posture, telemetry_freshness, maintenance_tasks",
    "- redaction_boundary: no endpoint values, tokens, private paths, raw payloads, DB URLs, notification destinations, or retained evidence paths",
])

manifest_md = [
    "# Operations Reliability Manifest",
    "",
    f"- output_dir: {rel_to_root(out)}",
    "- files: summary.json, summary.md, manifest.json, manifest.md, reliability-review.txt",
    "- excluded: docs/evidence, evidence-like paths, symlinked paths, oversized inputs, raw logs, backup dumps, DB URLs, webhook values, secrets, private payloads",
    "- sends_notifications: false",
    "- mutating_admin_endpoints_called: false",
]

review = [
    "Operations reliability review",
    "",
    "Private diagnostic sampling only.",
    "No evidence, compliance, production readiness, SLA, uptime guarantee, hosted SaaS, agency adoption, consumer acceptance, vendor compatibility, or production-grade ETA claim is made.",
    f"Overall status: {overall}",
    f"Monitoring export status: {monitoring_export['status']}",
    "Notification sent: false",
]

(out / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8")
(out / "summary.md").write_text("\n".join(summary_md) + "\n", encoding="utf-8")
(out / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
(out / "manifest.md").write_text("\n".join(manifest_md) + "\n", encoding="utf-8")
(out / "reliability-review.txt").write_text("\n".join(review) + "\n", encoding="utf-8")

actual = sorted(p.name for p in out.iterdir())
if actual != sorted(OUTPUT_FILES):
    fail(f"output contract violation: {actual}")

for name in OUTPUT_FILES:
    text = (out / name).read_text(encoding="utf-8")
    if unsafe_text(text):
        fail(f"{name} contains forbidden private value")
    if len(text.encode("utf-8")) > 65536:
        fail(f"{name} exceeds bounded output size")

print(f"operations reliability diagnostics written to {rel_to_root(out)}")
PY
