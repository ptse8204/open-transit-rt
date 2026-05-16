#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

RELEASE_PACKAGE_DIR="${RELEASE_PACKAGE_DIR:-}"

usage() {
  cat <<'EOF'
Usage:
  RELEASE_PACKAGE_DIR=.cache/release-package/<version> scripts/audit-release-package.sh

Environment:
  RELEASE_PACKAGE_DIR  Required path to a generated local release package.

Safety:
  Audits local release package structure, JSON, checksums, claim flags, wording,
  and consumer tracker preservation. It does not publish artifacts, contact
  consumers, create evidence, or change repository files.
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
if [ -z "$RELEASE_PACKAGE_DIR" ]; then
  printf 'ERROR: RELEASE_PACKAGE_DIR is required\n' >&2
  exit 1
fi

python3 - "$ROOT_DIR" "$RELEASE_PACKAGE_DIR" <<'PY'
import hashlib
import json
import pathlib
import re
import sys

root_arg, package_arg = sys.argv[1:3]
ROOT = pathlib.Path(root_arg).resolve()
PACKAGE = pathlib.Path(package_arg)
if not PACKAGE.is_absolute():
    PACKAGE = ROOT / PACKAGE
PACKAGE = PACKAGE.resolve(strict=False)

ROOT_EXPECTED = sorted([
    "artifacts",
    "checksums",
    "image.json",
    "manifest.json",
    "manifest.md",
    "provenance.json",
    "provenance.md",
    "sbom.json",
    "summary.json",
    "summary.md",
])
JSON_FILES = ["summary.json", "manifest.json", "provenance.json", "sbom.json", "image.json"]
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
        r"token\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
        r"raw[_-](log|payload|telemetry|correspondence)",
        r"\bis compliant\b",
        r"\bhas consumer acceptance\b",
        r"\bis production ready\b",
        r"\bhosted service is available\b",
        r"\bhosted saas is available\b",
    )
]


def fail(message):
    raise SystemExit(f"ERROR: {message}")


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def is_relative_to(path, base):
    try:
        pathlib.Path(path).resolve(strict=False).relative_to(pathlib.Path(base).resolve(strict=False))
        return True
    except ValueError:
        return False


def is_evidence_like(path):
    raw = str(path).replace("\\", "/").lower()
    parts = [part.lower() for part in pathlib.Path(path).parts]
    return "docs/evidence" in raw or "evidence" in parts or "proof" in parts or "submission" in parts


def sha256_file(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def check_claim_flags(data, label):
    flags = data.get("claim_flags", {})
    if not flags:
        fail(f"{label} missing claim_flags")
    true_flags = [key for key, value in flags.items() if value is not False]
    if true_flags:
        fail(f"{label} has true/nonfalse claim flags: {true_flags}")


def check_consumer_tracker():
    path = ROOT / "docs/evidence/consumer-submissions/status.json"
    if not path.exists():
        attrs = ROOT / ".gitattributes"
        archive_export = (
            not (ROOT / ".git").exists()
            and attrs.exists()
            and "/docs/evidence/consumer-submissions/status.json export-ignore" in attrs.read_text(encoding="utf-8", errors="replace")
        )
        if archive_export:
            return
    data = json.loads(path.read_text(encoding="utf-8"))
    records = data.get("targets", [])
    seen = {row["target"]: row.get("status") for row in records}
    if list(seen) != EXPECTED_CONSUMERS:
        fail(f"consumer target order drifted: {seen}")
    if any(seen[name] != "prepared" for name in EXPECTED_CONSUMERS):
        fail(f"consumer target status drifted: {seen}")


if not PACKAGE.exists() or not PACKAGE.is_dir():
    fail(f"release package directory not found: {rel(PACKAGE)}")
if not is_relative_to(PACKAGE, ROOT / ".cache"):
    fail("release package must be under repo .cache")
if is_evidence_like(PACKAGE):
    fail("release package path must not be evidence-like")

root_names = sorted(child.name for child in PACKAGE.iterdir())
if root_names != ROOT_EXPECTED:
    fail(f"unexpected release package root files: {root_names}")

artifact_files = sorted((PACKAGE / "artifacts").iterdir())
if len(artifact_files) != 1 or not re.fullmatch(r"open-transit-rt-[A-Za-z0-9._-]+\.source\.tar\.gz", artifact_files[0].name):
    fail(f"unexpected artifact files: {[p.name for p in artifact_files]}")
checksum_files = sorted((PACKAGE / "checksums").iterdir())
if [p.name for p in checksum_files] != ["SHA256SUMS.txt"]:
    fail(f"unexpected checksum files: {[p.name for p in checksum_files]}")

for name in JSON_FILES:
    data = json.loads((PACKAGE / name).read_text(encoding="utf-8"))
    check_claim_flags(data, name)

expected = {}
for line in (PACKAGE / "checksums" / "SHA256SUMS.txt").read_text(encoding="utf-8").splitlines():
    if not line.strip():
        continue
    digest, sep, filename = line.partition("  ")
    if sep != "  " or not re.fullmatch(r"[0-9a-f]{64}", digest):
        fail(f"invalid checksum line: {line}")
    target = pathlib.PurePosixPath(filename)
    if target.is_absolute() or ".." in target.parts or filename == "checksums/SHA256SUMS.txt":
        fail(f"unsafe checksum path: {filename}")
    expected[filename] = digest

actual = {}
for path in sorted(PACKAGE.rglob("*")):
    if not path.is_file() or path.name == "SHA256SUMS.txt":
        continue
    actual[path.relative_to(PACKAGE).as_posix()] = sha256_file(path)
if expected != actual:
    fail("checksum manifest does not match package files")

for path in sorted(PACKAGE.rglob("*")):
    if not path.is_file() or path.suffix == ".gz":
        continue
    text = path.read_text(encoding="utf-8", errors="replace")
    hits = [pattern.pattern for pattern in UNSAFE_PATTERNS if pattern.search(text)]
    if hits:
        fail(f"unsafe or unsupported claim text in {rel(path)}: {hits[:3]}")

check_consumer_tracker()
print(f"release package audit passed: {rel(PACKAGE)}")
PY
