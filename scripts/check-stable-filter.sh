#!/usr/bin/env sh
set -eu

usage() {
  cat <<'EOF'
Usage: scripts/check-stable-filter.sh [--source PATH] [--stable-tree PATH] [--skip-ref-check]

Verifies the filtered stable branch contract without contacting external
services. With --stable-tree, checks an already-filtered tree. Without it, the
script builds a temporary filtered tree from --source using
.github/stable-sync-excludes.txt and checks that result.
EOF
}

SOURCE="."
STABLE_TREE=""
CHECK_REFS=1

while [ "$#" -gt 0 ]; do
  case "$1" in
    --source)
      SOURCE="$2"
      shift 2
      ;;
    --stable-tree)
      STABLE_TREE="$2"
      shift 2
      ;;
    --skip-ref-check)
      CHECK_REFS=0
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

SOURCE_ABS="$(python3 -c 'import pathlib,sys; print(pathlib.Path(sys.argv[1]).resolve())' "$SOURCE")"

if [ ! -f "$SOURCE_ABS/.github/stable-sync-excludes.txt" ]; then
  echo "missing stable exclude file under source: $SOURCE_ABS" >&2
  exit 1
fi

TMP_DIR=""
cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT INT TERM

if [ -z "$STABLE_TREE" ]; then
  TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/open-transit-rt-stable-filter.XXXXXX")"
  rsync -a \
    --exclude .git \
    --exclude-from "$SOURCE_ABS/.github/stable-sync-excludes.txt" \
    "$SOURCE_ABS/" "$TMP_DIR/"
  STABLE_TREE="$TMP_DIR"
fi

STABLE_ABS="$(python3 -c 'import pathlib,sys; print(pathlib.Path(sys.argv[1]).resolve())' "$STABLE_TREE")"

python3 - "$SOURCE_ABS" "$STABLE_ABS" "$CHECK_REFS" <<'PY'
from __future__ import annotations

import json
import pathlib
import subprocess
import sys


source = pathlib.Path(sys.argv[1])
stable = pathlib.Path(sys.argv[2])
check_refs = sys.argv[3] == "1"

required_excludes = [
    ".cache/",
    ".DS_Store",
    "AGENTS.md",
    "docs/agent/",
    "docs/handoffs/",
    "docs/prompts/",
    "docs/roadmaps/",
    "docs/codex-task.md",
    "docs/conversation-summary.md",
    "docs/codex-*.md",
    "docs/codex-*.txt",
    "docs/codex-*.yaml",
    "docs/codex-*.yml",
    "docs/phase-*.md",
    "docs/*phase-*.md",
    "docs/**/*phase-*.md",
    "docs/phase-plan.md",
    "docs/phase-plan-production-closure.md",
    "docs/open-transit-rt-master-planner-remaining-work.md",
    "docs/*roadmap*.md",
    "docs/track-b-*.md",
]

must_preserve = [
    ".github/workflows/test.yml",
    ".github/workflows/release-gates.yml",
    ".github/workflows/update-stable.yml",
    ".github/stable-sync-excludes.txt",
    "README.md",
    "Makefile",
    "go.mod",
    "cmd/agency-config/main.go",
    "cmd/feed-vehicle-positions/main.go",
    "cmd/feed-trip-updates/main.go",
    "cmd/feed-alerts/main.go",
    "cmd/gtfs-studio/main.go",
    "internal/admincontrol/model.go",
    "internal/gtfs/importer.go",
    "examples/README.md",
    "testdata/README.md",
    "docs/index.md",
    "docs/tutorials/no-cli-agency-first-run.md",
    "docs/tutorials/small-agency-maintenance-guide.md",
    "docs/connectors/catalog.md",
    "docs/branching-and-release-policy.md",
    "docs/ci.md",
    "docs/release-status-v0.1.0-rc.2.md",
    "site/index.html",
    "site/connectors.html",
    "site/readiness.html",
    "site/video.html",
    "db/migrations/000001_initial_schema.sql",
    "deploy/docker-compose.yml",
    "scripts/check-consumer-tracker.sh",
    "scripts/check-internal-links.sh",
    "scripts/check-stable-filter.sh",
]

failures: list[str] = []


def rel_paths(root: pathlib.Path) -> list[str]:
    out: list[str] = []
    for path in root.rglob("*"):
        if path.is_file() or path.is_dir():
            out.append(path.relative_to(root).as_posix())
    return out


def is_excluded_agent_path(rel: str) -> bool:
    name = pathlib.PurePosixPath(rel).name.lower()
    if rel == "AGENTS.md":
        return True
    if rel.startswith(("docs/agent/", "docs/handoffs/", "docs/prompts/", "docs/roadmaps/")):
        return True
    if rel in {
        "docs/codex-task.md",
        "docs/conversation-summary.md",
        "docs/phase-plan.md",
        "docs/phase-plan-production-closure.md",
        "docs/open-transit-rt-master-planner-remaining-work.md",
    }:
        return True
    if rel.startswith("docs/codex-") and name.endswith((".md", ".txt", ".yaml", ".yml")):
        return True
    if rel.startswith("docs/") and name.endswith(".md") and "phase-" in name:
        return True
    if rel.startswith("docs/") and name.endswith(".md") and "roadmap" in name:
        return True
    if rel.startswith("docs/track-b-") and name.endswith(".md"):
        return True
    return False


exclude_text = (source / ".github/stable-sync-excludes.txt").read_text(encoding="utf-8")
exclude_lines = {line.strip() for line in exclude_text.splitlines() if line.strip() and not line.strip().startswith("#")}
for pattern in required_excludes:
    if pattern not in exclude_lines:
        failures.append(f"stable exclude list missing required pattern: {pattern}")

if check_refs:
    try:
        subprocess.run(["git", "show-ref", "--verify", "--quiet", "refs/heads/stable"], cwd=source, check=True)
    except subprocess.CalledProcessError:
        failures.append("local stable branch ref is missing")
    try:
        subprocess.run(["git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/stable"], cwd=source, check=True)
    except subprocess.CalledProcessError:
        failures.append("origin/stable remote-tracking ref is missing")

for rel in rel_paths(stable):
    if is_excluded_agent_path(rel):
        failures.append(f"filtered stable tree contains excluded agent/planning path: {rel}")

for rel in must_preserve:
    if not (stable / rel).exists():
        failures.append(f"filtered stable tree is missing preserved product path: {rel}")

status_path = stable / "docs/evidence/consumer-submissions/status.json"
expected_targets = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
if not status_path.is_file():
    failures.append("filtered stable tree is missing consumer submission status tracker")
else:
    data = json.loads(status_path.read_text(encoding="utf-8"))
    targets = data.get("targets", [])
    if len(targets) != 7:
        failures.append(f"consumer tracker target count = {len(targets)}, want 7")
    target_names = [target.get("target") for target in targets if isinstance(target, dict)]
    if target_names != expected_targets:
        failures.append(f"consumer tracker target order/name drift: {target_names}")
    for target in targets:
        if target.get("status") != "prepared":
            failures.append(f"consumer tracker target {target.get('target')} status = {target.get('status')}, want prepared")

workflow = (source / ".github/workflows/update-stable.yml").read_text(encoding="utf-8")
if "dry_run:" not in workflow or "DRY_RUN" not in workflow:
    failures.append("update-stable workflow is missing dry-run support")
if "--force" in workflow or "force-with-lease" in workflow:
    failures.append("update-stable workflow appears to use force-push semantics")
for line in workflow.splitlines():
    stripped = line.strip()
    if stripped.startswith("git push") and "+" in stripped and (
        ":stable" in stripped or ":refs/heads/stable" in stripped
    ):
        failures.append("update-stable workflow appears to use a force refspec for stable")
if "refs/remotes/origin/stable" not in workflow or "STABLE_BASE" not in workflow:
    failures.append("update-stable workflow is missing explicit stable tip guard")
if "check-stable-filter.sh" not in workflow:
    failures.append("update-stable workflow does not run the stable filter checker")

if failures:
    print("stable filter check failed:")
    for failure in failures:
        print(f"  {failure}")
    sys.exit(1)

print("stable filter check passed")
PY

if [ -x "$STABLE_ABS/scripts/check-internal-links.sh" ]; then
  (cd "$STABLE_ABS" && scripts/check-internal-links.sh)
fi
