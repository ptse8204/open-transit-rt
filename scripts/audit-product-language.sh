#!/usr/bin/env sh
set -eu

ROOT_DIR="${PRODUCT_LANGUAGE_AUDIT_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-product-language.sh [--help]

Environment:
  PRODUCT_LANGUAGE_AUDIT_ROOT   Optional repository root override.

Safety:
  Read-only local audit. It scans primary public pages, entry docs, and
  Operations Console source for user-facing wording regressions. It does not
  start services, fetch URLs, contact outside systems, write evidence, or
  mutate consumer status records.
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

PRIMARY_FILES = [
    "README.md",
    "docs/index.md",
    "docs/tutorials/no-cli-agency-first-run.md",
    "docs/tutorials/video-recording-guide.md",
    "docs/deployment/oci-reference-deployment.md",
    "wiki/README.md",
]
PRIMARY_GLOBS = [
    "site/*.html",
    "site/assets/*.vtt",
]
CONSOLE_COPY_FILES = [
    "cmd/agency-config/operations.go",
    "cmd/agency-config/operations_navigation.go",
    "cmd/agency-config/operations_route_registry.go",
    "cmd/agency-config/operations_admin.js",
]

BANNED_COPY = [
    "Common Next Action",
    "Common Next Actions",
    "Follow the same path staff use in the console",
    "What stays technical?",
    "site readout",
    "phase closed",
    "checkpoint",
    "AI agent",
    "Codex",
    "claim flags",
]

RAW_FLAG_NAMES = [
    "external_evidence_created",
    "consumer_statuses_changed",
    "production_grade_eta_claimed",
    "hosted_saas_claimed",
    "dynamic_backend_plugin_loading_enabled",
]

ENTRY_DOC_ALLOWED = {
    "docs/index.md": {
        "AI agent",  # allowed only while existing audit-product-acceptance expects an AI context section
        "Codex",
    },
}

failures: list[str] = []


def rel(path: pathlib.Path) -> str:
    try:
        return path.relative_to(ROOT).as_posix()
    except ValueError:
        return str(path)


def read(path: pathlib.Path) -> str:
    try:
        return path.read_text(encoding="utf-8")
    except FileNotFoundError:
        failures.append(f"missing audit target: {rel(path)}")
        return ""


def line_number(text: str, index: int) -> int:
    return text.count("\n", 0, index) + 1


def primary_paths() -> list[pathlib.Path]:
    paths = [ROOT / item for item in PRIMARY_FILES]
    for pattern in PRIMARY_GLOBS:
        paths.extend(ROOT.glob(pattern))
    return sorted(set(path for path in paths if path.is_file() or path.suffix))


def check_phrase(path: pathlib.Path, text: str, phrase: str) -> None:
    allowed = ENTRY_DOC_ALLOWED.get(rel(path), set())
    if phrase in allowed:
        return
    pattern = re.compile(re.escape(phrase), re.IGNORECASE)
    for match in pattern.finditer(text):
        failures.append(f"{rel(path)}:{line_number(text, match.start())}: banned phrase: {phrase}")


def check_raw_flags_in_primary(path: pathlib.Path, text: str) -> None:
    for flag in RAW_FLAG_NAMES:
        pattern = re.compile(re.escape(flag), re.IGNORECASE)
        for match in pattern.finditer(text):
            failures.append(f"{rel(path)}:{line_number(text, match.start())}: raw internal flag visible in primary copy: {flag}")


def extract_template_copy(text: str) -> str:
    # Keep the template string and plain string literals searchable for visible
    # labels. JSON struct tags and test fixtures are intentionally not audited
    # here because companion JSON fields remain part of the product contract.
    return text


for path in primary_paths():
    text = read(path)
    for phrase in BANNED_COPY:
        check_phrase(path, text, phrase)
    check_raw_flags_in_primary(path, text)

for relative in CONSOLE_COPY_FILES:
    path = ROOT / relative
    text = extract_template_copy(read(path))
    for phrase in BANNED_COPY:
        if phrase in {"AI agent", "Codex"}:
            continue
        check_phrase(path, text, phrase)

if failures:
    print("product language audit failed:")
    for failure in failures:
        print(f"  {failure}")
    raise SystemExit(1)

print("product language audit passed")
PY
