#!/usr/bin/env sh
set -eu

ROOT_DIR="${CONSUMER_TRACKER_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"
STATUS_PATH="${CONSUMER_TRACKER_STATUS_PATH:-docs/evidence/consumer-submissions/status.json}"

python3 - "$ROOT_DIR" "$STATUS_PATH" <<'PY'
import json
import pathlib
import sys

root_arg, status_arg = sys.argv[1:3]
ROOT = pathlib.Path(root_arg).resolve()
STATUS = pathlib.Path(status_arg)
if not STATUS.is_absolute():
    STATUS = ROOT / STATUS

EXPECTED = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]


def is_release_archive_without_protected_tracker():
    if STATUS.exists() or (ROOT / ".git").exists():
        return False
    attrs = ROOT / ".gitattributes"
    if not attrs.exists():
        return False
    text = attrs.read_text(encoding="utf-8", errors="replace")
    return "/docs/evidence/consumer-submissions/status.json export-ignore" in text


if not STATUS.exists():
    if is_release_archive_without_protected_tracker():
        print("consumer tracker check skipped: protected tracker is export-ignored from source archive")
        raise SystemExit(0)
    raise SystemExit(f"consumer status JSON missing: {STATUS}")

data = json.loads(STATUS.read_text(encoding="utf-8"))
records = data.get("targets", [])
seen = {row.get("target"): row.get("status") for row in records if isinstance(row, dict)}
if list(seen) != EXPECTED:
    raise SystemExit(f"consumer tracker target order/name drift: {seen}")
if any(seen[name] != "prepared" for name in EXPECTED):
    raise SystemExit(f"consumer tracker status drift: {seen}")
print("consumer tracker has exactly seven prepared-only targets")
PY
