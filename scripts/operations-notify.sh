#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/operations-notify/$TIMESTAMP}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
ALLOW_UNIGNORED_SOURCE_DIR="${ALLOW_UNIGNORED_SOURCE_DIR:-false}"
FORCE="${FORCE:-false}"
STRICT_OPERATIONS_NOTIFY="${STRICT_OPERATIONS_NOTIFY:-false}"
MAX_SOURCE_BYTES="${MAX_SOURCE_BYTES:-5242880}"
VALIDATOR_HEALTH_SUMMARY="${VALIDATOR_HEALTH_SUMMARY:-}"
DEPLOYMENT_DOCTOR_SUMMARY="${DEPLOYMENT_DOCTOR_SUMMARY:-}"
NOTIFY_WEBHOOK_URL="${NOTIFY_WEBHOOK_URL:-}"
NOTIFY_EMAIL_TO="${NOTIFY_EMAIL_TO:-}"
DRY_RUN="false"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/operations-notify.sh [--help] [--dry-run]

Environment:
  OUTPUT_DIR                         Default .cache/operations-notify/<timestamp>
  FORCE                              true|false; allow non-empty OUTPUT_DIR reuse
  ALLOW_UNIGNORED_OUTPUT_DIR         true|false; allow output outside .cache except evidence-like paths
  VALIDATOR_HEALTH_SUMMARY           Optional explicit .cache/validator-health/.../summary.json
  DEPLOYMENT_DOCTOR_SUMMARY          Optional explicit .cache/deployment-doctor/.../summary.json
  ALLOW_UNIGNORED_SOURCE_DIR         true|false; allow explicit source summaries outside .cache except evidence-like paths
  STRICT_OPERATIONS_NOTIFY           true|false; fail on missing, malformed, oversized, blocked, or unhealthy inputs
  MAX_SOURCE_BYTES                   Maximum bytes per source summary, default 5242880
  NOTIFY_WEBHOOK_URL                 Presence is recorded as a boolean only; value is never written
  NOTIFY_EMAIL_TO                    Presence is recorded as a boolean only; value is never written

Dry-run:
  --dry-run writes the same five local draft files as default mode. It requires
  no source summaries, network, webhook, email, database, Docker, admin token,
  application process, or running service.

Safety:
  This helper creates a private local notification draft only. It never sends
  notifications, calls admin routes, runs validators, contacts consumers,
  creates evidence packets, writes docs/evidence, edits GTFS, blocks publishing,
  changes consumer statuses, or claims compliance/readiness/adoption.
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
bool_var STRICT_OPERATIONS_NOTIFY "$STRICT_OPERATIONS_NOTIFY"
positive_int MAX_SOURCE_BYTES "$MAX_SOURCE_BYTES"

TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-operations-notify.XXXXXX")"

if ! python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$TMP_DIR" "$TIMESTAMP" "$FORCE" "$ALLOW_UNIGNORED_OUTPUT_DIR" "$ALLOW_UNIGNORED_SOURCE_DIR" "$STRICT_OPERATIONS_NOTIFY" "$MAX_SOURCE_BYTES" "$VALIDATOR_HEALTH_SUMMARY" "$DEPLOYMENT_DOCTOR_SUMMARY" "$DRY_RUN" "$NOTIFY_WEBHOOK_URL" "$NOTIFY_EMAIL_TO" <<'PY'
import json
import os
import pathlib
import re
import shutil
import sys

(
    root_arg,
    output_arg,
    tmp_arg,
    timestamp,
    force_arg,
    allow_output_arg,
    allow_source_arg,
    strict_arg,
    max_source_bytes_arg,
    validator_arg,
    doctor_arg,
    dry_run_arg,
    webhook_arg,
    email_arg,
) = sys.argv[1:15]

ROOT = pathlib.Path(root_arg).resolve()
OUTPUT_RAW = pathlib.Path(output_arg)
TMP = pathlib.Path(tmp_arg).resolve()
FORCE = force_arg == "true"
ALLOW_OUTPUT = allow_output_arg == "true"
ALLOW_SOURCE = allow_source_arg == "true"
STRICT = strict_arg == "true"
MAX_SOURCE_BYTES = int(max_source_bytes_arg)
DRY_RUN = dry_run_arg == "true"

OUTPUT_FILES = ("summary.json", "summary.md", "manifest.json", "manifest.md", "notification.txt")
MAX_NEXT_ACTIONS = 20
MAX_NOTIFICATION_BYTES = 16 * 1024
MAX_SUMMARY_MD_BYTES = 24 * 1024
MAX_TEXT = 220
TIMESTAMP_DIR = re.compile(r"^\d{8}T\d{6}Z$")


def fail(message):
    raise SystemExit(message)


def rel_to_root(path):
    try:
        return path.resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<redacted-source>"


def path_has_symlink(path):
    probe = pathlib.Path(path)
    if not probe.is_absolute():
        probe = ROOT / probe
    current = probe.anchor
    current_path = pathlib.Path(current) if current else pathlib.Path(".")
    for part in probe.parts[1 if current else 0:]:
        current_path = current_path / part
        if current_path.exists() and current_path.is_symlink():
            return True
    return False


def is_evidence_like(path):
    raw = str(path).replace("\\", "/").lower()
    parts = [p.lower() for p in pathlib.Path(path).parts]
    if "docs/evidence" in raw or "evidence" in parts:
        return True
    return False


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


def source_allowed(path):
    if is_evidence_like(path):
        return False
    if path_has_symlink(path):
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


def discover_latest(kind):
    base = ROOT / ".cache" / kind
    if not base.exists() or path_has_symlink(base):
        return None
    candidates = []
    for child in base.iterdir():
        if not child.is_dir() or child.is_symlink():
            continue
        if not TIMESTAMP_DIR.match(child.name):
            continue
        summary = child / "summary.json"
        if not summary.exists() or summary.is_symlink():
            continue
        candidates.append((child.name, summary))
    if not candidates:
        return None
    return sorted(candidates, key=lambda item: item[0])[-1][1]


def safe_string(value, default=""):
    if value is None:
        return default
    if isinstance(value, (dict, list)):
        return default
    text = str(value).replace("\r", " ").replace("\n", " ").strip()
    text = re.sub(r"\s+", " ", text)
    if len(text) > MAX_TEXT:
        text = text[: MAX_TEXT - 15].rstrip() + " [truncated]"
    return text


def base_source(kind, explicit):
    if explicit:
        path = pathlib.Path(explicit)
        if not path.is_absolute():
            path = ROOT / path
        return path
    return discover_latest(kind)


def read_json_source(name, kind, explicit):
    path = base_source(kind, explicit)
    if path is None:
        return {
            "name": name,
            "kind": "missing_source",
            "status": "missing_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": None,
            "data": None,
            "next_action": f"Run the local {name} diagnostic helper first, then rerun operations-notify.",
        }
    if not source_allowed(path) or path.name != "summary.json":
        return {
            "name": name,
            "kind": "rejected_source",
            "status": "malformed_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": rel_to_root(path),
            "data": None,
            "next_action": f"Use a non-symlink {name} summary.json under .cache, outside docs/evidence.",
        }
    if not path.exists():
        return {
            "name": name,
            "kind": "missing_source",
            "status": "missing_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": rel_to_root(path),
            "data": None,
            "next_action": f"Run the local {name} diagnostic helper first, then rerun operations-notify.",
        }
    try:
        size = path.stat().st_size
    except OSError:
        size = -1
    if size < 0 or size > MAX_SOURCE_BYTES:
        return {
            "name": name,
            "kind": "too_large_source",
            "status": "too_large_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": rel_to_root(path),
            "data": None,
            "next_action": f"Regenerate a bounded {name} summary or raise MAX_SOURCE_BYTES for local diagnostics only.",
        }
    try:
        raw = path.read_bytes()
    except OSError:
        return {
            "name": name,
            "kind": "unreadable_source",
            "status": "malformed_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": rel_to_root(path),
            "data": None,
            "next_action": f"Regenerate the {name} summary because the selected file could not be read.",
        }
    try:
        data = json.loads(raw.decode("utf-8"))
        if not isinstance(data, dict):
            raise ValueError("summary root is not an object")
    except Exception:
        return {
            "name": name,
            "kind": "malformed_source",
            "status": "malformed_source",
            "severity": "blocked" if STRICT else "needs_review",
            "source_file": rel_to_root(path),
            "data": None,
            "next_action": f"Regenerate the {name} summary because the selected JSON is malformed.",
        }
    return {
        "name": name,
        "kind": "loaded",
        "status": "loaded",
        "severity": "unknown",
        "source_file": rel_to_root(path),
        "data": data,
        "next_action": "",
    }


def severity_rank(severity):
    return {"unknown": 0, "info": 1, "needs_review": 2, "blocked": 3}.get(severity, 0)


def aggregate_severity(values):
    selected = "unknown"
    for value in values:
        if severity_rank(value) > severity_rank(selected):
            selected = value
    return selected


def validator_health_status_to_severity(status):
    status = (status or "").lower()
    if status in {"blocked", "failed", "missing_tooling", "misconfigured_tooling", "artifact_unavailable"}:
        return "blocked"
    if status in {"stale", "needs_review", "not_run", "skipped", "unknown"}:
        return "needs_review"
    if status in {"recorded", "configured", "installed", "runnable", "configured_for_tests", "stub", "passed"}:
        return "info"
    return "unknown"


def doctor_status_to_severity(status):
    status = (status or "").lower()
    if status == "blocker":
        return "blocked"
    if status in {"warning", "unavailable", "skipped"}:
        return "needs_review"
    if status == "passed":
        return "info"
    return "unknown"


def add_action(actions, overflow, source, severity, title, action, count=1, overflow_count=0):
    item = {
        "source": source,
        "severity": severity,
        "title": safe_string(title, "Review diagnostic source"),
        "action": safe_string(action, "Review the private local diagnostic summary."),
        "count": int(count or 0),
        "overflow_count": int(overflow_count or 0),
    }
    if len(actions) < MAX_NEXT_ACTIONS:
        actions.append(item)
    else:
        overflow[0] += 1


def summarize_validator(source, actions, overflow):
    if source["kind"] != "loaded":
        add_action(actions, overflow, "validator-health", source["severity"], source["status"], source["next_action"])
        return {
            "source": "validator-health",
            "status": source["status"],
            "severity": source["severity"],
            "loaded": False,
        }
    data = source["data"]
    feeds = data.get("feeds") if isinstance(data.get("feeds"), list) else []
    feed_count = len(feeds)
    row_severities = []
    for row in feeds[:MAX_NEXT_ACTIONS]:
        if not isinstance(row, dict):
            continue
        health = safe_string(row.get("health_status"), "unknown")
        sev = validator_health_status_to_severity(health)
        row_severities.append(sev)
        if sev in {"blocked", "needs_review"}:
            add_action(
                actions,
                overflow,
                "validator-health",
                sev,
                f"{safe_string(row.get('feed_type'), 'feed')} validation health is {health}",
                safe_string(row.get("next_action"), "Review validator-health diagnostics and rerun the validator helper if needed."),
            )
    if feed_count > MAX_NEXT_ACTIONS:
        overflow[0] += feed_count - MAX_NEXT_ACTIONS
    overall = safe_string(data.get("overall_status"), "unknown")
    row_severities.append(validator_health_status_to_severity(overall))
    severity = aggregate_severity(row_severities or ["unknown"])
    if severity == "unknown" and feed_count == 0:
        severity = "needs_review"
        add_action(actions, overflow, "validator-health", severity, "validator-health summary has no feed rows", "Regenerate validator-health before using this draft.")
    return {
        "source": "validator-health",
        "status": overall,
        "severity": severity,
        "loaded": True,
        "counts": {"feeds": feed_count},
    }


def summarize_doctor(source, actions, overflow):
    if source["kind"] != "loaded":
        add_action(actions, overflow, "deployment-doctor", source["severity"], source["status"], source["next_action"])
        return {
            "source": "deployment-doctor",
            "status": source["status"],
            "severity": source["severity"],
            "loaded": False,
        }
    data = source["data"]
    overall = safe_string(data.get("overall_status"), "unknown")
    severity = doctor_status_to_severity(overall)
    counts = data.get("counts") if isinstance(data.get("counts"), dict) else {}
    checks = data.get("checks") if isinstance(data.get("checks"), list) else []
    problem_rows = []
    for row in checks:
        if not isinstance(row, dict):
            continue
        row_status = safe_string(row.get("status"), "unknown")
        row_sev = doctor_status_to_severity(row_status)
        if row_sev in {"blocked", "needs_review"}:
            problem_rows.append((row, row_sev))
    for row, row_sev in problem_rows[:MAX_NEXT_ACTIONS]:
        add_action(
            actions,
            overflow,
            "deployment-doctor",
            row_sev,
            f"{safe_string(row.get('category'), 'deployment')} / {safe_string(row.get('name'), 'check')}: {safe_string(row.get('status'), 'unknown')}",
            safe_string(row.get("detail"), "Review deployment-doctor private summary and resolve the reported blocker or warning."),
        )
    if len(problem_rows) > MAX_NEXT_ACTIONS:
        overflow[0] += len(problem_rows) - MAX_NEXT_ACTIONS
    return {
        "source": "deployment-doctor",
        "status": overall,
        "severity": severity,
        "loaded": True,
        "counts": {
            "checks": len(checks),
            "blockers": int(counts.get("blocker", 0) or 0),
            "warnings": int(counts.get("warning", 0) or 0),
            "unavailable": int(counts.get("unavailable", 0) or 0),
        },
    }


def cap_bytes(text, limit):
    data = text.encode("utf-8")
    if len(data) <= limit:
        return text
    suffix = "\n[truncated]\n"
    keep = max(0, limit - len(suffix.encode("utf-8")))
    return data[:keep].decode("utf-8", errors="ignore") + suffix


def output_path_for_manifest(path):
    rel = rel_to_root(path)
    if rel.startswith(".cache/"):
        return rel
    if rel == "<redacted-source>":
        return rel
    return "<redacted-source>"


def scan_text(name, text, strict_names=True):
    patterns = [
        ("authorization", re.compile(r"Authorization\s*:\s*\S+", re.I)),
        ("bearer", re.compile(r"\bBearer\s+[A-Za-z0-9._~+/=-]{8,}", re.I)),
        ("cookie", re.compile(r"\bCookie\s*:\s*\S+|\badmin_session\s*=", re.I)),
        ("csrf", re.compile(r"\bcsrf(?:_token)?\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}", re.I)),
        ("database_url", re.compile(r"\bDATABASE_URL\s*[:=]|postgres(?:ql)?://[^:\s/]+:[^@\s]+@", re.I)),
        ("webhook_url", re.compile(r"https?://[^\s\"']*(?:webhook|hooks)[^\s\"']*(?:token|secret|key|sig|signature)[^\s\"']*", re.I)),
        ("private_key", re.compile(r"BEGIN [A-Z ]*PRIVATE KEY")),
        ("env_secret", re.compile(r"\b(?:TOKEN|SECRET|PASSWORD)\s*=\s*\S+", re.I)),
        ("private_path", re.compile(r"(/Users/|/tmp/|/var/|/etc/|deployment-private)", re.I)),
        ("absolute_private_path", re.compile(r"\b[A-Za-z0-9_.-]*/deployment-private\b", re.I)),
    ]
    if strict_names:
        patterns.append(("raw_field_data", re.compile(r'"(?:raw_report|stdout|stderr|argv)"\s*:', re.I)))
    for label, pattern in patterns:
        if pattern.search(text):
            fail(f"redaction failure in {name}: {label}")


def write_outputs(out, validator_source, doctor_source):
    actions = []
    overflow = [0]
    validator_summary = summarize_validator(validator_source, actions, overflow)
    doctor_summary = summarize_doctor(doctor_source, actions, overflow)
    overall = aggregate_severity([validator_summary["severity"], doctor_summary["severity"]])
    if overall == "unknown" and not actions:
        overall = "info"
    counts = {
        "sources_loaded": int(validator_summary.get("loaded", False)) + int(doctor_summary.get("loaded", False)),
        "sources_missing_or_skipped": int(not validator_summary.get("loaded", False)) + int(not doctor_summary.get("loaded", False)),
        "next_actions": len(actions),
        "blocked_actions": sum(1 for a in actions if a["severity"] == "blocked"),
        "needs_review_actions": sum(1 for a in actions if a["severity"] == "needs_review"),
    }
    notification = {
        "title": "Self-hosted operations notification draft",
        "severity": overall,
        "not_sent": True,
        "not_sent_to": ["webhook", "email", "consumers", "agency", "public_service"],
        "boundary": "Private local diagnostics only; not evidence, compliance proof, production readiness proof, or consumer acceptance.",
    }
    source_files = []
    for src in (validator_source, doctor_source):
        if src.get("source_file"):
            source_files.append(output_path_for_manifest(ROOT / src["source_file"] if src["source_file"].startswith(".cache/") else pathlib.Path(src["source_file"])))
    summary = {
        "generated_at": timestamp,
        "source_summaries": [validator_summary, doctor_summary],
        "notification": notification,
        "destinations": {
            "webhook_present": bool(webhook_arg),
            "email_present": bool(email_arg),
        },
        "counts": counts,
        "next_actions": actions,
        "overflow_count": overflow[0],
        "external_evidence_created": False,
        "consumer_statuses_changed": False,
        "compliance_claimed": False,
        "production_readiness_claimed": False,
        "hosted_saas_claimed": False,
        "agency_adoption_claimed": False,
        "consumer_acceptance_claimed": False,
        "vendor_compatibility_claimed": False,
        "production_grade_eta_claimed": False,
        "notification_sent": False,
    }
    manifest = {
        "included_files": list(OUTPUT_FILES),
        "excluded_categories": [
            "raw reports",
            "source JSON contents",
            "tokens and secrets",
            "Authorization headers",
            "cookies",
            "database URLs",
            "webhook or email destination values",
            "evidence packets",
            "consumer submissions",
            "recursive source artifacts",
        ],
        "source_files": source_files,
        "output_dir": rel_to_root(out),
        "created_at": timestamp,
    }

    summary_json = json.dumps(summary, indent=2, sort_keys=True) + "\n"
    manifest_json = json.dumps(manifest, indent=2, sort_keys=True) + "\n"
    summary_md = "# Self-Hosted Operations Notification Summary\n\n"
    summary_md += f"- Generated at: `{timestamp}`\n"
    summary_md += f"- Severity: `{overall}`\n"
    summary_md += "- Notification sent: `false`\n"
    summary_md += "- Boundary: private local diagnostics only; not evidence, not a compliance gate, not production health proof, and not consumer acceptance.\n\n"
    summary_md += "## Sources\n\n"
    for src in summary["source_summaries"]:
        summary_md += f"- {src['source']}: `{src['status']}` / `{src['severity']}`\n"
    summary_md += "\n## Next Actions\n\n"
    for action in actions[:MAX_NEXT_ACTIONS]:
        summary_md += f"- `{action['severity']}` {action['source']}: {action['title']} - {action['action']}\n"
    if overflow[0]:
        summary_md += f"- Additional source items omitted: `{overflow[0]}`\n"
    summary_md = cap_bytes(summary_md, MAX_SUMMARY_MD_BYTES)

    manifest_md = "# Operations Notification Manifest\n\n"
    manifest_md += f"- Created at: `{timestamp}`\n"
    manifest_md += f"- Output directory: `{manifest['output_dir']}`\n"
    manifest_md += "- Included files: `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, `notification.txt`\n"
    manifest_md += "- Excluded: raw reports, recursive source artifacts, tokens, secrets, destination values, source JSON contents, consumer submissions, and evidence packets.\n"
    manifest_md += "- This is not evidence and was not sent.\n"

    notification_txt = "DRAFT — NOT SENT\n\n"
    notification_txt += f"Severity: {overall}\n"
    notification_txt += "This local summary was not sent to webhook, email, consumers, agency, or public service.\n"
    notification_txt += "It is not evidence, compliance proof, production readiness proof, or consumer acceptance.\n\n"
    notification_txt += "Next actions:\n"
    if actions:
        for action in actions[:MAX_NEXT_ACTIONS]:
            notification_txt += f"- [{action['severity']}] {action['title']}: {action['action']}\n"
    else:
        notification_txt += "- Review loaded private diagnostics; no send action was taken.\n"
    if overflow[0]:
        notification_txt += f"- {overflow[0]} additional source items omitted from this bounded draft.\n"
    notification_txt = cap_bytes(notification_txt, MAX_NOTIFICATION_BYTES)

    generated = {
        "summary.json": summary_json,
        "summary.md": summary_md,
        "manifest.json": manifest_json,
        "manifest.md": manifest_md,
        "notification.txt": notification_txt,
    }
    for name, text in generated.items():
        scan_text(name, text, strict_names=name in {"summary.json", "summary.md", "notification.txt"})
    for name, text in generated.items():
        tmp_file = TMP / name
        tmp_file.write_text(text, encoding="utf-8")
        os.replace(tmp_file, out / name)
    actual = sorted(p.name for p in out.iterdir())
    if actual != sorted(OUTPUT_FILES):
        fail(f"unexpected output files: {actual}")
    json.loads((out / "summary.json").read_text(encoding="utf-8"))
    json.loads((out / "manifest.json").read_text(encoding="utf-8"))
    return overall


output_dir = resolve_output_dir()
validator_source = read_json_source("validator-health", "validator-health", validator_arg)
doctor_source = read_json_source("deployment-doctor", "deployment-doctor", doctor_arg)
severity = write_outputs(output_dir, validator_source, doctor_source)
print(f"operations notification draft written to {rel_to_root(output_dir)}")

strict_fail = False
if STRICT:
    if severity == "blocked":
        strict_fail = True
    for src in (validator_source, doctor_source):
        if src["kind"] != "loaded":
            strict_fail = True
    for src_summary in json.loads((output_dir / "summary.json").read_text(encoding="utf-8"))["source_summaries"]:
        if src_summary.get("severity") in {"blocked", "needs_review"}:
            strict_fail = True
if strict_fail:
    raise SystemExit(2)
PY
then
  fail "operations notification summary generation failed"
fi
