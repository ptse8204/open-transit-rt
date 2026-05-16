#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

AUDIT_MODE="${AUDIT_MODE:-real}"
BLOCKER_ONLY="false"

usage() {
  cat <<'EOF'
Usage:
  FINAL_ROOT_PACKET_DIR=<packet-dir> [AUDIT_MODE=real|blocker] scripts/audit-final-root-evidence.sh [--help] [--blocker-only]

Environment:
  FINAL_ROOT_PACKET_DIR   Required path to a final-root evidence or blocker packet.
  AUDIT_MODE              real|blocker, default real. --blocker-only sets blocker mode.

Safety:
  Real audit is intentionally conservative. It fails on missing approval,
  root mismatch, placeholders, missing feed/checksum artifacts, missing or
  unavailable validator status, unsafe private strings, checksum mismatches,
  missing redaction notes, or consumer tracker drift. Blocker audit passes only
  for truthful blocker packets with no captured evidence and all claim flags
  false.
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
    --blocker-only)
      BLOCKER_ONLY="true"
      AUDIT_MODE="blocker"
      ;;
    *)
      fail_exit "unknown argument: $1"
      ;;
  esac
  shift
done

if [ -z "${FINAL_ROOT_PACKET_DIR:-}" ]; then
  usage >&2
  fail_exit "missing required environment variable: FINAL_ROOT_PACKET_DIR"
fi

case "$AUDIT_MODE" in
  real|blocker) ;;
  *) fail_exit "AUDIT_MODE must be real or blocker" ;;
esac

python3 - "$ROOT_DIR" "$FINAL_ROOT_PACKET_DIR" "$AUDIT_MODE" "$BLOCKER_ONLY" <<'PY'
import hashlib
import json
import pathlib
import re
import sys
import urllib.parse

root_arg, packet_arg, mode, blocker_flag = sys.argv[1:5]
ROOT = pathlib.Path(root_arg).resolve()
PACKET = (ROOT / packet_arg).resolve(strict=False) if not pathlib.Path(packet_arg).is_absolute() else pathlib.Path(packet_arg).resolve(strict=False)

BLOCKER_FILES = {"blocker.json", "blocker.md", "manifest.json", "manifest.md"}
REAL_TOP = {
    "README.md",
    "approval.md",
    "dns-tls-redirect.md",
    "public-fetches.md",
    "validator-record.md",
    "proxy-config-summary.md",
    "redaction-notes.md",
    "manifest.json",
    "manifest.md",
    "SHA256SUMS.txt",
    "artifacts",
}
REQUIRED_ARTIFACT_DIRS = {
    "artifacts/public",
    "artifacts/tls",
    "artifacts/dns",
    "artifacts/validation",
    "artifacts/operator-supplied",
}
REQUIRED_PUBLIC = {
    "artifacts/public/feeds.json",
    "artifacts/public/schedule.zip",
    "artifacts/public/vehicle_positions.pb",
    "artifacts/public/trip_updates.pb",
    "artifacts/public/alerts.pb",
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
        r"database_url\s*=",
        r"begin [a-z ]*private key",
        r"acme[_ -]?account",
        r"private[_ -]?key",
        r"secret\s*[:=]",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|diagnostic)",
        r"dns[_ -]?provider[_ -]?(token|secret|export)",
    )
]
PLACEHOLDER = re.compile(r"\b(TBD|TODO|placeholder|pending|not collected|not available|unavailable|missing|fake|sample)\b", re.I)

failures = []


def record_failure(message):
    failures.append(message)
    print(f"FAIL: {message}")


def record_pass(message):
    print(f"PASS: {message}")


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


def sha256_file(path):
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def has_symlink(path):
    probe = pathlib.Path(path)
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False


def check_consumer_tracker():
    status_path = ROOT / "docs" / "evidence" / "consumer-submissions" / "status.json"
    if not status_path.exists():
        attrs = ROOT / ".gitattributes"
        archive_export = (
            not (ROOT / ".git").exists()
            and attrs.exists()
            and "/docs/evidence/consumer-submissions/status.json export-ignore" in attrs.read_text(encoding="utf-8", errors="replace")
        )
        if archive_export:
            record_pass("consumer tracker check skipped because protected tracker is export-ignored from source archive")
            return
    data = load_json(status_path)
    records = data.get("targets", [])
    seen = {row.get("target"): row.get("status") for row in records}
    if list(seen) != EXPECTED_CONSUMERS:
        record_failure("consumer tracker target order changed")
        return
    if any(seen[name] != "prepared" for name in EXPECTED_CONSUMERS):
        record_failure("consumer tracker contains non-prepared status")
        return
    record_pass("consumer tracker preserved with seven prepared targets")


def check_claim_flags_false(flags, allow_final_retained=False):
    ok = True
    for key, value in sorted(flags.items()):
        if key == "final_root_evidence_retained" and allow_final_retained:
            continue
        if value is not False:
            record_failure(f"claim flag is not false: {key}={value!r}")
            ok = False
    if ok:
        record_pass("claim flags are false")


def check_unsafe_strings():
    found = False
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
                found = True
                break
    if not found:
        record_pass("no unsafe private strings detected")


def audit_blocker():
    if not PACKET.is_dir():
        record_failure(f"packet directory missing: {rel(PACKET)}")
        return
    names = {child.name for child in PACKET.iterdir()}
    if names == BLOCKER_FILES:
        record_pass("blocker packet contains exact required files")
    else:
        record_failure(f"blocker packet file set mismatch: {sorted(names)}")
    blocker = load_json(PACKET / "blocker.json")
    manifest = load_json(PACKET / "manifest.json")
    if blocker.get("schema") == "final-root-blocker.v1" and blocker.get("mode") == "blocker":
        record_pass("blocker schema and mode recorded")
    else:
        record_failure("blocker schema or mode is invalid")
    if blocker.get("real_final_root_evidence_retained") is False and blocker.get("captured_evidence_directory_created") is False:
        record_pass("blocker records no retained final-root evidence and no captured directory")
    else:
        record_failure("blocker does not truthfully record absent retained evidence/captured directory")
    if not (blocker.get("final_root_available") and blocker.get("approval_artifact_available")):
        record_pass("blocker records no real root or no approval artifact")
    else:
        record_failure("blocker cannot pass when both real root and approval artifact are recorded available")
    try:
        PACKET.relative_to(ROOT / "docs" / "evidence" / "captured")
        record_failure("blocker packet is under docs/evidence/captured")
    except ValueError:
        record_pass("blocker packet is not under docs/evidence/captured")
    if manifest.get("retained_captured") is False and manifest.get("docs_evidence_captured_changed") is False:
        record_pass("manifest records no captured evidence write")
    else:
        record_failure("manifest indicates captured evidence was retained or changed")
    check_claim_flags_false(blocker.get("claim_flags", {}))
    check_claim_flags_false(manifest.get("claim_flags", {}))
    check_consumer_tracker()
    check_unsafe_strings()


def parse_root(root):
    parsed = urllib.parse.urlparse(root or "")
    if parsed.username or parsed.password or parsed.query or parsed.fragment:
        return None
    host = parsed.hostname or ""
    loopback = host in {"localhost", "127.0.0.1", "::1"} or host.startswith("127.")
    if parsed.scheme == "https" or (parsed.scheme == "http" and loopback):
        return parsed
    return None


def audit_real():
    if not PACKET.is_dir():
        record_failure(f"packet directory missing: {rel(PACKET)}")
        return
    if has_symlink(PACKET):
        record_failure("packet path contains symlink")
    else:
        record_pass("packet path contains no symlink")
    names = {child.name for child in PACKET.iterdir()}
    if names == REAL_TOP:
        record_pass("real packet contains exact top-level files/directories")
    else:
        record_failure(f"real packet top-level set mismatch: {sorted(names)}")
    for artifact_dir in REQUIRED_ARTIFACT_DIRS:
        if (PACKET / artifact_dir).is_dir():
            record_pass(f"{artifact_dir} exists")
        else:
            record_failure(f"{artifact_dir} missing")
    manifest = load_json(PACKET / "manifest.json")
    root = manifest.get("final_root_base_url", "")
    if parse_root(root):
        record_pass("manifest final root is a valid final-root URL shape")
    else:
        record_failure("manifest final root is invalid")
    approval = manifest.get("approval_artifact", "")
    if approval and (PACKET / approval).is_file() and (PACKET / approval).stat().st_size > 0:
        record_pass("approval evidence exists")
    else:
        record_failure("approval evidence missing")
    for required in REQUIRED_PUBLIC:
        path = PACKET / required
        if path.is_file() and path.stat().st_size > 0:
            record_pass(f"public artifact exists: {required}")
        else:
            record_failure(f"public artifact missing or empty: {required}")
    fetches = manifest.get("public_fetches", [])
    expected_paths = {path for _, path, _ in (
        ("feeds_json", "/public/feeds.json", "feeds.json"),
        ("schedule", "/public/gtfs/schedule.zip", "schedule.zip"),
        ("vehicle_positions", "/public/gtfsrt/vehicle_positions.pb", "vehicle_positions.pb"),
        ("trip_updates", "/public/gtfsrt/trip_updates.pb", "trip_updates.pb"),
        ("alerts", "/public/gtfsrt/alerts.pb", "alerts.pb"),
    )}
    seen_paths = {row.get("path") for row in fetches}
    if seen_paths == expected_paths:
        record_pass("all five final feed paths recorded")
    else:
        record_failure("public fetch manifest paths are incomplete")
    for row in fetches:
        if row.get("status") != 200:
            record_failure(f"public fetch status must be exactly 200 for {row.get('feed_type')}: {row.get('status')}")
        if not row.get("sha256"):
            record_failure(f"public fetch checksum missing for {row.get('feed_type')}")
        artifact = pathlib.Path(str(row.get("artifact", "")))
        summary = PACKET / artifact.parent / f"{artifact.name}.summary.json"
        if summary.is_file():
            summary_data = load_json(summary)
            expected_url = root.rstrip("/") + str(row.get("path", ""))
            if summary_data.get("url") != expected_url:
                record_failure(f"public fetch root mismatch for {row.get('feed_type')}")
        else:
            record_failure(f"public fetch summary missing for {row.get('feed_type')}")
    validation_ok = True
    for row in manifest.get("validator_records", []):
        if row.get("status") != "passed":
            record_failure(f"validator status is not passed for {row.get('feed_type')}: {row.get('status')}")
            validation_ok = False
        artifact = PACKET / str(row.get("artifact", ""))
        if not artifact.is_file():
            record_failure(f"validator artifact missing for {row.get('feed_type')}")
            validation_ok = False
    if validation_ok and manifest.get("validator_records"):
        record_pass("validator records are present and passed")
    if (PACKET / "redaction-notes.md").is_file() and (PACKET / "redaction-notes.md").stat().st_size > 0:
        notes = (PACKET / "redaction-notes.md").read_text(encoding="utf-8", errors="ignore").lower()
        if "authorization" in notes and "cookie" in notes and "private" in notes:
            record_pass("required redaction notes exist")
        else:
            record_failure("redaction notes are missing required private-data boundaries")
    else:
        record_failure("redaction notes missing")
    placeholder_found = False
    for path in PACKET.rglob("*.md"):
        text = path.read_text(encoding="utf-8", errors="ignore")
        if PLACEHOLDER.search(text):
            record_failure(f"placeholder or unavailable marker remains in {rel(path)}")
            placeholder_found = True
    if not placeholder_found:
        record_pass("no placeholder/unavailable markers in markdown")
    checksums = PACKET / "SHA256SUMS.txt"
    if checksums.is_file():
        bad = False
        listed = set()
        for line in checksums.read_text(encoding="utf-8").splitlines():
            if not line.strip():
                continue
            parts = line.split(None, 1)
            if len(parts) != 2:
                record_failure(f"malformed checksum line: {line}")
                bad = True
                continue
            expected, name = parts
            name = name.strip()
            if name in listed:
                record_failure(f"duplicate checksum entry: {name}")
                bad = True
            listed.add(name)
            target = PACKET / name
            if not target.is_file():
                record_failure(f"checksum target missing: {name}")
                bad = True
            elif sha256_file(target) != expected:
                record_failure(f"checksum mismatch: {name}")
                bad = True
        expected_files = {
            path.relative_to(PACKET).as_posix()
            for path in PACKET.rglob("*")
            if path.is_file() and path.name != "SHA256SUMS.txt"
        }
        missing_checksums = sorted(expected_files - listed)
        unexpected_checksums = sorted(listed - expected_files)
        if missing_checksums:
            record_failure(f"SHA256SUMS.txt missing entries: {', '.join(missing_checksums[:10])}")
            bad = True
        if unexpected_checksums:
            record_failure(f"SHA256SUMS.txt has unexpected entries: {', '.join(unexpected_checksums[:10])}")
            bad = True
        if not bad:
            record_pass("SHA256SUMS.txt is complete and matches packet files")
    else:
        record_failure("SHA256SUMS.txt missing")
    check_claim_flags_false(manifest.get("claim_flags", {}), allow_final_retained=True)
    check_consumer_tracker()
    check_unsafe_strings()


if mode == "blocker":
    audit_blocker()
else:
    audit_real()

if failures:
    print(f"\nFinal-root evidence audit failed with {len(failures)} blocker(s).", file=sys.stderr)
    sys.exit(1)
print("\nFinal-root evidence audit passed.")
PY
