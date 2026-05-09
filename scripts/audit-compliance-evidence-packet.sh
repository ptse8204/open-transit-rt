#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  COMPLIANCE_PACKET_DIR=<packet-dir> scripts/audit-compliance-evidence-packet.sh [--help]

Environment:
  COMPLIANCE_PACKET_DIR                       Required packet directory.
  COMPLIANCE_CONSUMER_STATUS_PATH             Optional test override for consumer status JSON.
  COMPLIANCE_CONSUMER_ARTIFACTS_DIR           Optional test override for artifact directory.

Safety:
  The audit fails on wrong packet files, invalid JSON, true claim flags,
  `compliant` readiness statuses, unsafe private strings, prepared-only
  consumer tracker drift, non-README consumer artifact files, and misleading
  compliance, consumer, final-root, deployment, operations, vendor, SLA, or ETA
  claims.
EOF
}

fail_exit() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    *)
      fail_exit "unknown argument: $1"
      ;;
  esac
  shift
done

if [ -z "${COMPLIANCE_PACKET_DIR:-}" ]; then
  usage >&2
  fail_exit "missing required environment variable: COMPLIANCE_PACKET_DIR"
fi

python3 - "$ROOT_DIR" "$COMPLIANCE_PACKET_DIR" \
  "${COMPLIANCE_CONSUMER_STATUS_PATH:-docs/evidence/consumer-submissions/status.json}" \
  "${COMPLIANCE_CONSUMER_ARTIFACTS_DIR:-docs/evidence/consumer-submissions/artifacts}" <<'PY'
import json
import pathlib
import re
import sys

root_arg, packet_arg, status_arg, artifacts_arg = sys.argv[1:5]
ROOT = pathlib.Path(root_arg).resolve()
PACKET = pathlib.Path(packet_arg)
if not PACKET.is_absolute():
    PACKET = ROOT / PACKET
PACKET = PACKET.resolve(strict=False)
STATUS_PATH = pathlib.Path(status_arg)
if not STATUS_PATH.is_absolute():
    STATUS_PATH = ROOT / STATUS_PATH
ARTIFACTS_DIR = pathlib.Path(artifacts_arg)
if not ARTIFACTS_DIR.is_absolute():
    ARTIFACTS_DIR = ROOT / ARTIFACTS_DIR

BLOCKER_FILES = {"blocker.json", "blocker.md", "manifest.json", "manifest.md"}
DEPLOYMENT_FILES = {
    "summary.json",
    "summary.md",
    "readiness-packet.md",
    "evidence-map.json",
    "evidence-map.md",
    "missing-evidence.md",
    "human-review.md",
    "manifest.json",
    "manifest.md",
}
JSON_BY_MODE = {
    "blocker": ("blocker.json", "manifest.json"),
    "deployment": ("summary.json", "evidence-map.json", "manifest.json"),
}
EXPECTED_CONSUMERS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
UNSAFE_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"authorization\s*:",
        r"cookie\s*:",
        r"\bbearer\s+[A-Za-z0-9._~+/=-]{8,}",
        r"postgres(?:ql)?://",
        r"database_url\s*[:=]",
        r"begin [a-z ]*private key",
        r"acme[_ -]?account",
        r"private[_ -]?key",
        r"secret\s*[:=]",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|diagnostic|telemetry|correspondence)",
        r"consumer[_ -]?portal",
        r"db[_ -]?url\s*[:=]",
        r"token\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
    )
]
MISLEADING_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"\bis\s+(?:cal-?itp|caltrans)\s+compliant\b",
        r"\bcompliance\s+(?:is\s+)?(?:achieved|proven|certified|complete)\b",
        r"\bconsumer\s+(?:acceptance|ingestion)\s+(?:is\s+)?(?:claimed|proven|complete|accepted)\b",
        r"\baccepted\s+by\s+(?:Google Maps|Apple Maps|Transit App|Bing Maps|Moovit|Mobility Database|transit\.land)\b",
        r"\bagency\s+(?:adoption|endorsement|approval)\s+(?:is\s+)?(?:claimed|proven|complete)\b",
        r"\bfinal[- ]root\s+(?:proof|readiness|evidence)\s+(?:is\s+)?(?:claimed|proven|complete)\b",
        r"\bhosted\s+saas\s+(?:is\s+)?(?:available|claimed|proven)\b",
        r"\bproduction[- ]readiness\s+(?:is\s+)?(?:claimed|proven|complete)\b",
        r"\bSLA\s+(?:is\s+)?(?:met|claimed|proven)\b",
        r"\buptime\s+(?:is\s+)?(?:guaranteed|claimed|proven)\b",
        r"\bvendor\s+compatibility\s+(?:is\s+)?(?:claimed|proven|complete)\b",
        r"\bproduction[- ]grade\s+ETA\s+(?:is\s+)?(?:claimed|proven|complete)\b",
        r"\bmarketplace\s+approval\s+(?:is\s+)?(?:claimed|proven|complete)\b",
    )
]
ALLOWED_STATUSES = {"present", "partial", "missing", "blocked", "pilot_only", "needs_review"}
failures = []


def record_pass(message):
    print(f"PASS: {message}")


def record_failure(message):
    failures.append(message)
    print(f"FAIL: {message}")


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return str(path)


def load_json(path):
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except Exception as exc:
        record_failure(f"invalid JSON {rel(path)}: {exc}")
        return {}


def has_symlink(path):
    probe = pathlib.Path(path)
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False


def check_claim_flags_false(obj, location):
    flags = obj.get("claim_flags", {}) if isinstance(obj, dict) else {}
    if not isinstance(flags, dict):
        record_failure(f"claim_flags missing or invalid in {location}")
        return
    ok = True
    for key, value in sorted(flags.items()):
        if value is not False:
            record_failure(f"claim flag is not false in {location}: {key}={value!r}")
            ok = False
    if ok:
        record_pass(f"claim flags are false in {location}")


def walk_values(obj):
    if isinstance(obj, dict):
        for value in obj.values():
            yield from walk_values(value)
    elif isinstance(obj, list):
        for value in obj:
            yield from walk_values(value)
    else:
        yield obj


def check_no_compliant_status(obj, location):
    bad = False
    if isinstance(obj, dict):
        for key, value in obj.items():
            if key == "status":
                if isinstance(value, str) and value.lower() == "compliant":
                    record_failure(f"disallowed compliant status in {location}")
                    bad = True
                elif isinstance(value, str) and value not in ALLOWED_STATUSES and location.startswith("summary"):
                    if value not in {"prepared"}:
                        pass
            check_no_compliant_status(value, location)
    elif isinstance(obj, list):
        for value in obj:
            check_no_compliant_status(value, location)
    if not bad and location == "summary.json":
        record_pass("no compliant status value in summary.json")


def check_text_scans():
    unsafe_found = False
    misleading_found = False
    compliant_found = False
    for path in PACKET.rglob("*"):
        if not path.is_file():
            continue
        data = path.read_bytes()
        if b"\x00" in data:
            continue
        text = data.decode("utf-8", errors="ignore")
        for pattern in UNSAFE_PATTERNS:
            if pattern.search(text):
                record_failure(f"unsafe private string found in {rel(path)}")
                unsafe_found = True
                break
        for pattern in MISLEADING_PATTERNS:
            if pattern.search(text):
                record_failure(f"misleading claim wording found in {rel(path)}")
                misleading_found = True
                break
        if re.search(r'"status"\s*:\s*"compliant"', text, re.I) or re.search(r"\|\s*`?compliant`?\s*\|", text, re.I):
            record_failure(f"disallowed compliant status text found in {rel(path)}")
            compliant_found = True
    if not unsafe_found:
        record_pass("no unsafe private strings detected")
    if not misleading_found:
        record_pass("no misleading claims detected")
    if not compliant_found:
        record_pass("no compliant status text detected")


def check_consumer_tracker():
    data = load_json(STATUS_PATH)
    records = data.get("targets", [])
    seen = {row.get("target"): row.get("status") for row in records}
    if list(seen) != EXPECTED_CONSUMERS:
        record_failure("consumer tracker target order changed")
        return
    if any(seen[name] != "prepared" for name in EXPECTED_CONSUMERS):
        record_failure("consumer tracker contains non-prepared status")
        return
    record_pass("consumer tracker preserved with seven prepared targets")


def check_artifacts_readme_only():
    if not ARTIFACTS_DIR.exists():
        record_failure(f"consumer artifact directory missing: {rel(ARTIFACTS_DIR)}")
        return
    bad = []
    for path in ARTIFACTS_DIR.glob("*/*"):
        if path.is_file() and path.name != "README.md":
            bad.append(rel(path))
    if bad:
        record_failure("consumer artifact directories contain non-README files: " + ", ".join(bad[:10]))
    else:
        record_pass("consumer artifact directories are README-only")


def audit_packet():
    if not PACKET.is_dir():
        record_failure(f"packet directory missing: {rel(PACKET)}")
        return
    if has_symlink(PACKET):
        record_failure("packet path contains symlink")
    else:
        record_pass("packet path contains no symlink")
    names = {child.name for child in PACKET.iterdir()}
    if names == BLOCKER_FILES:
        mode = "blocker"
        record_pass("blocker packet contains exact required files")
    elif names == DEPLOYMENT_FILES:
        mode = "deployment"
        record_pass("deployment packet contains exact required files")
    else:
        record_failure(f"packet file set mismatch: {sorted(names)}")
        mode = "unknown"
    loaded = {}
    for filename in JSON_BY_MODE.get(mode, ("manifest.json",)):
        path = PACKET / filename
        if path.is_file():
            loaded[filename] = load_json(path)
        else:
            record_failure(f"required JSON file missing: {filename}")
    for filename, obj in loaded.items():
        check_claim_flags_false(obj, filename)
        check_no_compliant_status(obj, filename)
    if mode == "blocker":
        blocker = loaded.get("blocker.json", {})
        manifest = loaded.get("manifest.json", {})
        if blocker.get("mode") == "blocker" and blocker.get("retained_evidence_created") is False:
            record_pass("blocker records no retained evidence")
        else:
            record_failure("blocker does not truthfully record blocker mode and no retained evidence")
        if manifest.get("packet_type") == "blocker":
            record_pass("manifest records blocker packet type")
        else:
            record_failure("manifest packet type is not blocker")
    if mode == "deployment":
        summary = loaded.get("summary.json", {})
        manifest = loaded.get("manifest.json", {})
        rows = summary.get("readiness_rows", [])
        if rows and all(row.get("status") in ALLOWED_STATUSES for row in rows):
            record_pass("readiness rows use allowed non-compliant statuses")
        else:
            record_failure("readiness rows missing or contain unsupported statuses")
        if manifest.get("retained_evidence_created") is False and manifest.get("consumer_tracker_changed") is False:
            record_pass("manifest records no retained evidence and no consumer tracker change")
        else:
            record_failure("manifest records retained evidence or consumer tracker change")
        if "pilot_only" in {row.get("status") for row in rows if row.get("requirement") in {"RQ-4A", "RQ-4C"}}:
            record_pass("OCI pilot evidence is constrained to pilot_only where represented")
        else:
            record_failure("OCI pilot evidence is not represented as pilot_only")
        missing = summary.get("missing_evidence", [])
        if any(item.get("area") == "final_root" and item.get("status") == "missing" for item in missing):
            record_pass("final-root missing evidence is reported truthfully")
        else:
            sources = summary.get("source_evidence_path_summaries", [])
            final_source = next((source for source in sources if source.get("label") == "final_root"), {})
            if final_source.get("exists") is True:
                record_pass("final-root source exists and missing blocker is not required")
            else:
                record_failure("final-root missing evidence is not reported truthfully")
        consumer = summary.get("consumer_tracker_summary", {})
        if consumer.get("all_expected_targets_prepared") is True:
            record_pass("packet records prepared-only consumer tracker")
        else:
            record_failure("packet does not record prepared-only consumer tracker")


audit_packet()
check_text_scans()
check_consumer_tracker()
check_artifacts_readme_only()

if failures:
    print(f"\nCompliance evidence packet audit failed with {len(failures)} blocker(s).", file=sys.stderr)
    sys.exit(1)
print("\nCompliance evidence packet audit passed.")
PY
