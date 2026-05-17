#!/usr/bin/env sh
set -eu

ROOT_DIR="${PRODUCT_ACCEPTANCE_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-product-acceptance.sh [--help]

Environment:
  PRODUCT_ACCEPTANCE_ROOT             Optional repository root override.
  PRODUCT_ACCEPTANCE_PUBLIC_DOCS      Optional colon-separated docs to scan.
  PRODUCT_ACCEPTANCE_STATUS_PATH      Optional consumer status JSON path.
  PRODUCT_ACCEPTANCE_SKIP_GIT_STATUS  Set true to skip git protected-path checks.

Safety:
  Local read-only audit. It checks the browser-first product path, public
  docs navigation, capability-versus-evidence wording, unsupported positive
  claim wording, prepared-only consumer tracker state, and protected evidence
  paths. It does not fetch feeds, contact external systems, require Docker, or
  require a running app.
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

DEFAULT_DOCS="README.md:wiki/README.md:wiki/small-agency-quick-start.md:wiki/browser-first-setup.md:wiki/operations-console-tour.md:wiki/can-my-agency-use-this.md:wiki/agency-demo.md:wiki/agency-adoption-checklist.md:wiki/connector-cookbook.md:wiki/calitp-readiness-plain-english.md:wiki/readiness-and-evidence.md:wiki/deployment-guide.md:wiki/how-agencies-can-help.md:wiki/support-and-contribute.md:docs/README.md:docs/index.md:docs/tutorials/small-agency-acceptance-script.md:docs/tutorials/agency-first-run.md:docs/tutorials/reusable-agency-onboarding.md:docs/tutorials/self-hosted-operator-trial.md:docs/tutorials/agency-launchpad.md:docs/release-candidate-readiness.md:docs/requirements-calitp-compliance.md"

python3 - "$ROOT_DIR" \
  "${PRODUCT_ACCEPTANCE_PUBLIC_DOCS:-$DEFAULT_DOCS}" \
  "${PRODUCT_ACCEPTANCE_STATUS_PATH:-docs/evidence/consumer-submissions/status.json}" \
  "${PRODUCT_ACCEPTANCE_SKIP_GIT_STATUS:-false}" <<'PY'
import json
import pathlib
import re
import subprocess
import sys

root_arg, docs_arg, status_arg, skip_git_arg = sys.argv[1:5]
ROOT = pathlib.Path(root_arg).resolve()
DOCS = [p for p in docs_arg.split(":") if p]
STATUS_PATH = pathlib.Path(status_arg)
if not STATUS_PATH.is_absolute():
    STATUS_PATH = ROOT / STATUS_PATH
SKIP_GIT_STATUS = skip_git_arg.lower() in {"1", "true", "yes"}

EXPECTED_CONSUMERS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
PROTECTED_PATHS = [
    "docs/evidence/captured",
    "docs/evidence/consumer-submissions/status.json",
    "docs/evidence/consumer-submissions/current",
    "docs/evidence/consumer-submissions/artifacts",
    "docs/evidence/consumer-submissions/packets",
]
FORBIDDEN_CLAIMS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"\bOpen Transit RT is (?:CAL-?ITP|Caltrans)[-/ ]?(?:style )?compliant\b",
        r"\bCAL-?ITP/Caltrans compliant\b",
        r"\bagency (?:adoption|approval|endorsement|deployment) (?:is |was )?(?:proven|complete|completed|secured|confirmed)\b",
        r"\b(?:an|the) agency (?:adopted|approved|endorsed|deployed) Open Transit RT\b",
        r"\b(?:Google Maps|Apple Maps|Transit App|Bing Maps|Moovit|Mobility Database|transit\.land) (?:has |have )?(?:accepted|ingested|listed|displayed|approved)\b",
        r"\b(?:accepted|ingested|listed|displayed|approved) by (?:Google Maps|Apple Maps|Transit App|Bing Maps|Moovit|Mobility Database|transit\.land)\b",
        r"\bconsumer (?:acceptance|ingestion|listing|display) (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bagency-owned (?:or agency-approved )?final[- ]root (?:is |was )?(?:ready|proven|complete|completed|approved|secured)\b",
        r"\bOpen Transit RT is (?:a )?hosted SaaS\b",
        r"\bhosted (?:service|SaaS) (?:is |was )?(?:available|launched|proven|complete|completed)\b",
        r"\bOpen Transit RT is production[- ]ready\b",
        r"\bproduction[- ]readiness (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bvendor compatibility (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bhardware certification (?:is |was )?(?:proven|complete|completed|confirmed)\b",
        r"\bSLA (?:coverage )?(?:is |was )?(?:available|included|met|proven|guaranteed)\b",
        r"\buptime (?:is |was )?(?:guaranteed|proven)\b",
        r"\bproduction[- ]grade ETA (?:quality )?(?:is |was )?(?:proven|complete|completed|confirmed)\b",
    )
]
NEGATION_CUES = (
    "not ",
    "not-",
    "no ",
    "does not",
    "do not",
    "must not",
    "without",
    "unsupported",
    "forbidden",
    "avoid",
    "what this does not prove",
    "doesn't",
    "cannot",
)
MAX_TEXT_BYTES = 1024 * 1024
failures = []


def resolve(path_arg):
    path = pathlib.Path(path_arg)
    if not path.is_absolute():
        path = ROOT / path
    return path.resolve(strict=False)


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return str(path)


def record_pass(message):
    print(f"PASS: {message}")


def record_failure(message):
    failures.append(message)
    print(f"FAIL: {message}")


def read_text(path_arg):
    path = resolve(path_arg)
    try:
        data = path.read_bytes()
    except FileNotFoundError:
        record_failure(f"missing file: {rel(path)}")
        return ""
    if len(data) > MAX_TEXT_BYTES:
        record_failure(f"file too large for bounded text audit: {rel(path)}")
        return ""
    if b"\x00" in data:
        record_failure(f"file appears binary: {rel(path)}")
        return ""
    try:
        return data.decode("utf-8")
    except UnicodeDecodeError as exc:
        record_failure(f"file is not UTF-8 text: {rel(path)}: {exc}")
        return ""


def require_file(path_arg, label):
    path = resolve(path_arg)
    if path.is_file():
        record_pass(f"{label} exists")
    else:
        record_failure(f"{label} missing: {rel(path)}")
    return path


def require_contains(text, needle, label):
    if needle in text:
        record_pass(label)
    else:
        record_failure(label)


def require_any(text, needles, label):
    if any(needle in text for needle in needles):
        record_pass(label)
    else:
        record_failure(label)


def check_front_doors():
    readme = read_text("README.md")
    require_contains(readme, "make agency-app-up", "README mentions make agency-app-up")
    require_contains(readme, "Operations Console", "README points to the Operations Console")
    require_any(readme, ["http://localhost:8080/admin/operations", "/admin/operations"], "README gives the private UI start path")
    body_start = "\n".join(readme.splitlines()[1:24]).lower()
    if "phase " in body_start or "checkpoint" in body_start:
        record_failure("README leads with phase/checkpoint history")
    else:
        record_pass("README does not lead with phase history")

    wiki_home = read_text("wiki/README.md")
    require_contains(wiki_home, "Small Agency Quick Start", "wiki home links Small Agency Quick Start")
    require_any(
        wiki_home,
        ["Browser-First Setup", "Operations Console Tour"],
        "wiki home links browser-first setup or Operations Console tour",
    )

    docs_home = read_text("docs/README.md")
    require_contains(docs_home, "[Docs Index](index.md)", "docs home points to the role-based docs index")
    docs_index = read_text("docs/index.md")
    for heading in (
        "## New Users",
        "## Agency Staff",
        "## Technical Helpers",
        "## Connector Developers",
        "## Maintainers",
        "## AI Agents",
    ):
        require_contains(docs_index, heading, f"docs index includes {heading}")
    maintainer_index = docs_index.find("## Maintainers")
    ai_index = docs_index.find("## AI Agents")
    if maintainer_index == -1 or ai_index == -1:
        record_failure("docs index is missing maintainer or AI-agent sections")
    else:
        before = docs_index[:maintainer_index].lower()
        after = docs_index[maintainer_index:].lower()
        if "phase" in after and "codex" in after and "phase" not in before:
            record_pass("docs index keeps phase/Codex history out of the new-user path")
        else:
            record_failure("docs index does not separate human-start docs from phase/Codex history")


def check_required_pages():
    require_file("wiki/small-agency-quick-start.md", "small agency quick start page")
    if resolve("wiki/browser-first-setup.md").is_file() or resolve("wiki/operations-console-tour.md").is_file():
        record_pass("browser-first setup or Operations Console tour page exists")
    else:
        record_failure("browser-first setup or Operations Console tour page missing")
    require_file("docs/tutorials/small-agency-acceptance-script.md", "small agency acceptance script page")

    cookbook = read_text("wiki/connector-cookbook.md")
    require_contains(cookbook, "## Practical Recipes", "connector cookbook has practical recipes")

    readiness = read_text("wiki/calitp-readiness-plain-english.md")
    require_any(
        readiness,
        [
            "| UI signal you can review | Missing deployment evidence before stronger claims |",
            "| UI signal you can review | Missing deployment evidence before outside approval, compliance, production, or consumer-acceptance claims |",
        ],
        "plain-English readiness guide distinguishes UI signal from missing deployment evidence",
    )

    requirements = read_text("docs/requirements-calitp-compliance.md")
    require_contains(requirements, "Software capability exists for GTFS import/publication", "requirements doc identifies software capability")
    require_contains(requirements, "External proof tracks are optional and authorization-gated", "requirements doc distinguishes optional evidence tracks")


def check_forbidden_claims():
    found = False
    for doc in DOCS:
        text = read_text(doc)
        lines = text.splitlines()
        for pattern in FORBIDDEN_CLAIMS:
            for index, line in enumerate(lines):
                if not pattern.search(line):
                    continue
                context = " ".join(lines[max(0, index - 4): index + 1]).lower()
                if any(cue in context for cue in NEGATION_CUES):
                    continue
                record_failure(f"unsupported positive claim wording found in {rel(resolve(doc))}: {line.strip()}")
                found = True
    if not found:
        record_pass("no forbidden positive claim phrases detected in public-facing docs")


def check_consumer_tracker():
    if not STATUS_PATH.exists():
        attrs = ROOT / ".gitattributes"
        archive_export = (
            not (ROOT / ".git").exists()
            and attrs.exists()
            and "/docs/evidence/consumer-submissions/status.json export-ignore" in attrs.read_text(encoding="utf-8", errors="replace")
        )
        if archive_export:
            record_pass("consumer tracker check skipped because protected tracker is export-ignored from source archive")
            return
    try:
        data = json.loads(STATUS_PATH.read_text(encoding="utf-8"))
    except FileNotFoundError:
        record_failure(f"consumer status JSON missing: {rel(STATUS_PATH)}")
        return
    except json.JSONDecodeError as exc:
        record_failure(f"consumer status JSON invalid: {rel(STATUS_PATH)}: {exc}")
        return
    records = data.get("targets", [])
    seen = {row.get("target"): row.get("status") for row in records}
    if list(seen) != EXPECTED_CONSUMERS:
        record_failure(f"consumer tracker target order/name drift: {seen}")
    elif all(seen[name] == "prepared" for name in EXPECTED_CONSUMERS):
        record_pass("consumer tracker has exactly seven prepared-only targets")
    else:
        record_failure(f"consumer tracker status drift: {seen}")


def check_protected_paths():
    if SKIP_GIT_STATUS:
        record_pass("protected evidence path git status skipped by test override")
        return
    if not (ROOT / ".git").exists():
        record_pass("protected evidence path git status skipped outside a git checkout")
        return
    diff = subprocess.run(
        ["git", "diff", "--quiet", "--"] + PROTECTED_PATHS,
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if diff.returncode == 0:
        record_pass("protected evidence paths have no tracked diff")
    else:
        record_failure(f"protected evidence paths have tracked changes: {diff.stderr.strip()}")
    status = subprocess.run(
        ["git", "status", "--short", "--"] + PROTECTED_PATHS,
        cwd=ROOT,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if status.returncode != 0:
        record_failure(f"protected evidence path status failed: {status.stderr.strip()}")
    elif status.stdout.strip():
        record_failure(f"protected evidence paths have tracked/untracked status: {status.stdout.strip()}")
    else:
        record_pass("protected evidence paths have no tracked or untracked status")


check_front_doors()
check_required_pages()
check_forbidden_claims()
check_consumer_tracker()
check_protected_paths()

if failures:
    print("\nproduct acceptance audit failed:")
    for failure in failures:
        print(f"- {failure}")
    sys.exit(1)

print("product acceptance audit passed")
PY
