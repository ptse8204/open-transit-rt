#!/usr/bin/env sh
set -eu

ROOT_DIR="${UI_LAYOUT_AUDIT_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-ui-layout.sh [--help]

Environment:
  UI_LAYOUT_AUDIT_ROOT   Optional repository root override.

Safety:
  Read-only local audit for static layout regressions. It does not render pages,
  start services, fetch URLs, or write artifacts.
EOF
}

case "${1:-}" in
  --help|-h)
    usage
    exit 0
    ;;
  "")
    ;;
  *)
    printf 'ERROR: unknown argument: %s\n' "$1" >&2
    exit 1
    ;;
esac

cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
from __future__ import annotations

import pathlib
import re
import sys

ROOT = pathlib.Path(sys.argv[1]).resolve()
failures: list[str] = []
warnings: list[str] = []


def rel(path: pathlib.Path) -> str:
    return path.relative_to(ROOT).as_posix()


def read(path: pathlib.Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except FileNotFoundError:
        failures.append(f"missing layout audit target: {rel(path)}")
        return ""


def check_css(path: pathlib.Path) -> None:
    text = read(path)
    for match in re.finditer(r"border-(?:left|right)\s*:\s*([2-9]|\d{2,})px", text, re.I):
        failures.append(f"{rel(path)}:{text.count(chr(10), 0, match.start()) + 1}: avoid decorative side-stripe borders")
    if re.search(r"background-clip\s*:\s*text", text, re.I):
        failures.append(f"{rel(path)}: avoid gradient text/background-clip text patterns")
    if re.search(r"\.card\s+\.card", text):
        failures.append(f"{rel(path)}: nested card styling pattern detected")


def check_site_page(path: pathlib.Path) -> None:
    text = read(path)
    if '<main' not in text or '</main>' not in text:
        failures.append(f"{rel(path)}: page needs a main landmark")
    card_mentions = len(re.findall(r'class="[^"]*(?:card|grid)[^"]*"', text))
    if card_mentions > 10:
        warnings.append(f"{rel(path)}: high card/grid class count ({card_mentions}); review for clutter")


check_css(ROOT / "site/assets/site.css")
check_css(ROOT / "cmd/agency-config/operations_design_system.go")
for path in sorted((ROOT / "site").glob("*.html")):
    check_site_page(path)

if failures:
    print("UI layout audit failed:")
    for failure in failures:
        print(f"  {failure}")
    for warning in warnings:
        print(f"WARN: {warning}")
    raise SystemExit(1)

for warning in warnings:
    print(f"WARN: {warning}")
print("UI layout audit passed")
PY
