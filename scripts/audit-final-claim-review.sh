#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-final-claim-review.sh [--help]

Environment:
  FINAL_CLAIM_REVIEW_DOCS           Optional colon-separated docs to scan.
  FINAL_CLAIM_REVIEW_PHASE_DOC      Optional Phase 60 doc path.
  FINAL_CLAIM_REVIEW_STATUS_PATH    Optional consumer status JSON path.
  FINAL_CLAIM_REVIEW_ARTIFACTS_DIR  Optional consumer artifact directory path.

Safety:
  Local read-only audit. It scans bounded public/status docs and verifies
  Phase 60 sections, unsupported positive claim wording, unsafe private
  strings, seven prepared-only consumer targets, and README-only consumer
  artifact directories. It does not fetch live feeds, contact agencies,
  contact consumers, create evidence, or change repository state.
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

DEFAULT_DOCS="README.md:docs/README.md:docs/current-status.md:docs/backlog.md:docs/open-questions.md:docs/roadmap-to-calitp-compliance-and-gap-closure.md:docs/roadmap-status.md:docs/public-launch-checklist.md:docs/public-share-copy.md:docs/california-readiness-summary.md:docs/compliance-evidence-checklist.md:docs/handoffs/latest.md"

python3 - "$ROOT_DIR" \
  "${FINAL_CLAIM_REVIEW_DOCS:-$DEFAULT_DOCS}" \
  "${FINAL_CLAIM_REVIEW_PHASE_DOC:-docs/phase-60-final-claim-review-and-public-closeout.md}" \
  "${FINAL_CLAIM_REVIEW_STATUS_PATH:-docs/evidence/consumer-submissions/status.json}" \
  "${FINAL_CLAIM_REVIEW_ARTIFACTS_DIR:-docs/evidence/consumer-submissions/artifacts}" <<'PY'
import json
import pathlib
import re
import sys

root_arg, docs_arg, phase_doc_arg, status_arg, artifacts_arg = sys.argv[1:6]
ROOT = pathlib.Path(root_arg).resolve()


def resolve(path_arg):
    path = pathlib.Path(path_arg)
    if not path.is_absolute():
        path = ROOT / path
    return path.resolve(strict=False)


DOCS = [resolve(p) for p in docs_arg.split(":") if p]
PHASE_DOC = resolve(phase_doc_arg)
STATUS_PATH = resolve(status_arg)
ARTIFACTS_DIR = resolve(artifacts_arg)

EXPECTED_CONSUMERS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
REQUIRED_PHASE60_SECTIONS = [
    "## Final Claim Review",
    "## Claim-To-Evidence Table",
    "## Unsupported Claims",
    "## Official Requirements Context",
    "## Retained Evidence Boundary",
    "## Maintainer Signoff",
    "## Execution Closeout",
]
REQUIRED_PUBLIC_BOUNDARY_PHRASES = [
    "does not claim CAL-ITP/Caltrans compliance",
    "all seven consumer and aggregator targets remain `prepared`",
    "not agency-owned final-root proof",
]
UNSAFE_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"authorization\s*:",
        r"cookie\s*:",
        r"\bbearer\s+[A-Za-z0-9._~+/=-]{8,}",
        r"postgres(?:ql)?://[^ \n`)]*:[^ \n`@)]*@",
        r"database_url\s*[:=]",
        r"begin [a-z ]*private key",
        r"acme[_ -]?account",
        r"private[_ -]?key",
        r"secret\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|diagnostic|telemetry|correspondence)",
        r"consumer[_ -]?portal",
        r"db[_ -]?url\s*[:=]",
        r"token\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
    )
]
CLAIM_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"\bOpen Transit RT is (?:CAL-?ITP|Caltrans) compliant\b",
        r"\bOpen Transit RT has publicly launched\b",
        r"\bpublic launch (?:occurred|is complete|is completed)\b",
        r"\b(?:an|the) agency (?:adopted|endorsed|approved|deployed) Open Transit RT\b",
        r"\bagency (?:adoption|endorsement|approval|deployment) (?:is |was )?(?:proven|complete|completed|secured|confirmed)\b",
        r"\bagency-owned (?:or agency-approved )?final[- ]root (?:is |was )?(?:ready|proven|complete|completed|approved|secured)\b",
        r"\b(?:Google Maps|Apple Maps|Transit App|Bing Maps|Moovit|Mobility Database|transit\.land) (?:has |have )?(?:submitted|accepted|ingested|listed|displayed|adopted|approved)\b",
        r"\b(?:accepted|ingested|listed|displayed|adopted|approved) by (?:Google Maps|Apple Maps|Transit App|Bing Maps|Moovit|Mobility Database|transit\.land)\b",
        r"\bconsumer (?:submission|review|acceptance|ingestion|listing|display|adoption) (?:is |was )?(?:complete|completed|proven|confirmed|accepted)\b",
        r"\bOpen Transit RT is (?:a )?hosted SaaS\b",
        r"\bhosted (?:service|SaaS) (?:is |was )?(?:available|launched|proven|complete|completed)\b",
        r"\bpaid support (?:is |was )?(?:available|included|launched|proven)\b",
        r"\bSLA (?:coverage )?(?:is |was )?(?:available|included|met|proven|guaranteed)\b",
        r"\buptime (?:is |was )?(?:guaranteed|proven)\b",
        r"\bproduction[- ]readiness (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bOpen Transit RT is production[- ]ready\b",
        r"\bproduction multi[- ]tenant hosting (?:is |was )?(?:ready|proven|complete|completed)\b",
        r"\bvendor compatibility (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bhardware certification (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bmarketplace approval (?:is |was )?(?:granted|received|proven|complete|completed|confirmed)\b",
        r"\bproduction[- ]grade ETA (?:quality )?(?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\breal[- ]world ETA accuracy (?:is |was )?(?:proven|complete|completed|confirmed)\b",
    )
]
MAX_TEXT_BYTES = 1024 * 1024
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


def read_text(path):
    try:
        data = path.read_bytes()
    except FileNotFoundError:
        record_failure(f"missing reviewed file: {rel(path)}")
        return None
    if len(data) > MAX_TEXT_BYTES:
        record_failure(f"reviewed file too large for bounded text scan: {rel(path)}")
        return None
    if b"\x00" in data:
        record_failure(f"reviewed file appears binary: {rel(path)}")
        return None
    try:
        return data.decode("utf-8")
    except UnicodeDecodeError as exc:
        record_failure(f"reviewed file is not UTF-8 text: {rel(path)}: {exc}")
        return None


def check_phase60_doc():
    text = read_text(PHASE_DOC)
    if text is None:
        return
    missing = [section for section in REQUIRED_PHASE60_SECTIONS if section not in text]
    if missing:
        record_failure(f"Phase 60 required sections missing: {missing}")
    else:
        record_pass("Phase 60 required sections present")
    if "Status\n\nComplete" in text or "Status: Complete" in text:
        record_pass("Phase 60 marked complete")
    else:
        record_failure("Phase 60 doc is not marked complete")


def check_doc_scans():
    unsafe_found = False
    claim_found = False
    for path in DOCS + [PHASE_DOC]:
        text = read_text(path)
        if text is None:
            continue
        lines = text.splitlines()
        for pattern in UNSAFE_PATTERNS:
            for index, line in enumerate(lines):
                if not pattern.search(line):
                    continue
                context = " ".join(lines[max(0, index - 4): index + 1]).lower()
                if any(cue in context for cue in (
                    "no secrets",
                    "exclude secrets",
                    "no private",
                    "not include",
                    "do not include",
                    "do not claim",
                    "do not commit",
                    "do not describe",
                    "must not",
                    "reject",
                    "redaction",
                    "unsafe",
                    "without",
                    "no portal",
                    "portal was automated",
                    "portal automation",
                    "automated no portal",
                    "contact external consumer portal",
                    "dev-device-token",
                )):
                    continue
                record_failure(f"unsafe private string found in {rel(path)}")
                unsafe_found = True
                break
            if unsafe_found:
                break
    for path in DOCS:
        text = read_text(path)
        if text is None:
            continue
        lines = text.splitlines()
        for pattern in CLAIM_PATTERNS:
            for index, line in enumerate(lines):
                match = pattern.search(line)
                if not match:
                    continue
                context = " ".join(lines[max(0, index - 4): index + 1]).lower()
                if any(cue in context for cue in (
                    "not supported",
                    "do not use",
                    "do not claim",
                    "do not describe",
                    "unsupported",
                    "forbidden overclaim",
                    "must not",
                    "does not claim",
                    "does not prove",
                    "not claim",
                    "not evidence",
                    "no public launch",
                    "without ",
                )):
                    continue
                record_failure(f"unsupported positive claim wording in {rel(path)}: {match.group(0)!r}")
                claim_found = True
                break
            if claim_found:
                break
    if not unsafe_found:
        record_pass("no unsafe private strings detected in reviewed docs")
    if not claim_found:
        record_pass("no unsupported positive claim wording detected in reviewed docs")


def check_required_public_boundaries():
    combined = "\n".join(read_text(path) or "" for path in DOCS + [PHASE_DOC])
    missing = [phrase for phrase in REQUIRED_PUBLIC_BOUNDARY_PHRASES if phrase.lower() not in combined.lower()]
    if missing:
        record_failure(f"required public boundary phrase missing: {missing}")
    else:
        record_pass("required public boundary phrases present")


def check_consumer_tracker():
    try:
        data = json.loads(STATUS_PATH.read_text(encoding="utf-8"))
    except Exception as exc:
        record_failure(f"invalid consumer status JSON {rel(STATUS_PATH)}: {exc}")
        return
    records = data.get("targets", [])
    seen = {row.get("target"): row.get("status") for row in records if isinstance(row, dict)}
    if list(seen) != EXPECTED_CONSUMERS:
        record_failure(f"consumer tracker target order changed: {seen}")
    elif any(seen[name] != "prepared" for name in EXPECTED_CONSUMERS):
        record_failure(f"consumer tracker contains non-prepared status: {seen}")
    else:
        record_pass("consumer tracker preserved with seven prepared targets")


def check_consumer_artifacts():
    if not ARTIFACTS_DIR.exists() or not ARTIFACTS_DIR.is_dir():
        record_failure(f"consumer artifacts directory missing: {rel(ARTIFACTS_DIR)}")
        return
    bad = []
    for path in ARTIFACTS_DIR.rglob("*"):
        if path.is_file() and path.name != "README.md":
            bad.append(rel(path))
    if bad:
        record_failure("consumer artifact directories contain non-README files: " + ", ".join(bad[:10]))
    else:
        record_pass("consumer artifact directories are README-only")


check_phase60_doc()
check_doc_scans()
check_required_public_boundaries()
check_consumer_tracker()
check_consumer_artifacts()

if failures:
    print("final claim review audit failed")
    for failure in failures:
        print(f"- {failure}")
    raise SystemExit(1)

print("final claim review audit passed")
PY
