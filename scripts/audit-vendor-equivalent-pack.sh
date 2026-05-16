#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

VENDOR_EQUIVALENT_PACK_DIR="${VENDOR_EQUIVALENT_PACK_DIR:-docs/vendor-equivalent-pack}"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-vendor-equivalent-pack.sh [--help]

Environment:
  VENDOR_EQUIVALENT_PACK_DIR  Optional template directory, default docs/vendor-equivalent-pack.

Safety:
  Audits template wording and required files only. It does not contact
  marketplaces, consumers, vendors, agencies, procurement systems, or external
  services, and it does not create evidence or change consumer statuses.
EOF
}

if [ "${1:-}" = "--help" ] || [ "${1:-}" = "-h" ]; then
  usage
  exit 0
fi
if [ "$#" -gt 0 ]; then
  printf 'ERROR: unknown argument: %s\n' "$1" >&2
  exit 1
fi

python3 - "$ROOT_DIR" "$VENDOR_EQUIVALENT_PACK_DIR" <<'PY'
import json
import pathlib
import re
import sys

root_arg, pack_arg = sys.argv[1:3]
ROOT = pathlib.Path(root_arg).resolve()
PACK = pathlib.Path(pack_arg)
if not PACK.is_absolute():
    PACK = ROOT / PACK
PACK = PACK.resolve(strict=False)

REQUIRED = [
    "README.md",
    "byod-hardware-intake-template.md",
    "implementation-plan-template.md",
    "support-boundaries-template.md",
    "sla-kpi-template.md",
    "procurement-response-template.md",
]
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
        r"secret\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|telemetry|correspondence)",
    )
]
CLAIM_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"\bis marketplace approved\b",
        r"\bmarketplace approval (?:granted|received|confirmed)\b",
        r"\bis vendor compatible\b",
        r"\bvendor compatibility confirmed\b",
        r"\bcertified hardware\b",
        r"\bhardware certification confirmed\b",
        r"\bpaid support included\b",
        r"\bguaranteed uptime of\b",
        r"\bis production ready\b",
        r"\bis compliant\b",
        r"\bhas consumer acceptance\b",
        r"\bconsumer acceptance confirmed\b",
        r"\bhosted service is available\b",
        r"\bhosted saas is available\b",
        r"\bapproved by (google maps|apple maps|transit app|bing maps|moovit|mobility database|transit\\.land)\b",
    )
]


def fail(message):
    raise SystemExit(f"ERROR: {message}")


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


if not PACK.exists() or not PACK.is_dir():
    fail(f"vendor-equivalent pack directory missing: {rel(PACK)}")
missing = [name for name in REQUIRED if not (PACK / name).is_file()]
if missing:
    fail(f"required vendor-equivalent templates missing: {missing}")

readme = (PACK / "README.md").read_text(encoding="utf-8")
sla = (PACK / "sla-kpi-template.md").read_text(encoding="utf-8")
for phrase, source in (
    ("template material only", readme),
    ("not marketplace approval", readme),
    ("not a service-level commitment", sla),
):
    if phrase.lower() not in source.lower():
        fail(f"required boundary phrase missing: {phrase}")

for name in REQUIRED:
    path = PACK / name
    text = path.read_text(encoding="utf-8", errors="replace")
    unsafe = [pattern.pattern for pattern in UNSAFE_PATTERNS if pattern.search(text)]
    if unsafe:
        fail(f"unsafe private string pattern in {rel(path)}: {unsafe[:3]}")
    claims = [pattern.pattern for pattern in CLAIM_PATTERNS if pattern.search(text)]
    if claims:
        fail(f"unsupported positive claim wording in {rel(path)}: {claims[:3]}")
    if name != "README.md" and ("<" not in text or ">" not in text):
        fail(f"template placeholders missing from {rel(path)}")

tracker_path = ROOT / "docs/evidence/consumer-submissions/status.json"
if not tracker_path.exists():
    attrs = ROOT / ".gitattributes"
    archive_export = (
        not (ROOT / ".git").exists()
        and attrs.exists()
        and "/docs/evidence/consumer-submissions/status.json export-ignore" in attrs.read_text(encoding="utf-8", errors="replace")
    )
    if not archive_export:
        fail(f"consumer tracker missing: {rel(tracker_path)}")
else:
    tracker = json.loads(tracker_path.read_text(encoding="utf-8"))
    records = tracker.get("targets", [])
    seen = {row["target"]: row.get("status") for row in records}
    if list(seen) != EXPECTED_CONSUMERS:
        fail(f"consumer target order drifted: {seen}")
    if any(seen[name] != "prepared" for name in EXPECTED_CONSUMERS):
        fail(f"consumer target status drifted: {seen}")

print(f"vendor-equivalent pack audit passed: {rel(PACK)}")
PY
