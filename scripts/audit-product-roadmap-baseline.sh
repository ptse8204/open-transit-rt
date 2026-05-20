#!/usr/bin/env sh
set -eu

ROOT_DIR="${PRODUCT_ROADMAP_BASELINE_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-product-roadmap-baseline.sh [--help]

Environment:
  PRODUCT_ROADMAP_BASELINE_ROOT   Optional repository root override.

Safety:
  Read-only local audit. It checks that the product-quality roadmap has one
  executable baseline across private Operations Console routes, public site
  source, connector/runtime checks, release gates, and protected consumer
  tracker boundaries. It does not start services, fetch URLs, contact outside
  systems, write evidence, or change consumer statuses.
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

import json
import pathlib
import re
import subprocess
import sys

ROOT = pathlib.Path(sys.argv[1]).resolve()
failures: list[str] = []
warnings: list[str] = []


def read(relative: str) -> str:
    path = ROOT / relative
    try:
        return path.read_text(encoding="utf-8")
    except FileNotFoundError:
        failures.append(f"missing required file: {relative}")
        return ""


def require_text(relative: str, needles: list[str]) -> None:
    text = read(relative)
    for needle in needles:
        if needle not in text:
            failures.append(f"{relative}: missing {needle!r}")


def require_regex(relative: str, pattern: str, label: str) -> None:
    text = read(relative)
    if not re.search(pattern, text, re.M):
        failures.append(f"{relative}: missing {label}")


def require_path(relative: str) -> None:
    if not (ROOT / relative).exists():
        failures.append(f"missing required path: {relative}")


def git(args: list[str]) -> subprocess.CompletedProcess[str]:
    return subprocess.run(["git", *args], cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)


route_text = read("cmd/agency-config/operations_route_registry.go")
required_routes = [
    "/admin/operations",
    "/admin/operations/gtfs-import",
    "/admin/operations/gtfs-workbench",
    "/admin/operations/feed-health",
    "/admin/operations/validation-center",
    "/admin/operations/realtime",
    "/admin/operations/prediction-lab",
    "/admin/operations/devices",
    "/admin/operations/connectors",
    "/admin/operations/connectors/workbench",
    "/admin/operations/connectors/tests",
    "/admin/operations/readiness",
    "/admin/operations/consumers",
    "/admin/operations/maintenance",
    "/admin/operations/reliability",
    "/admin/operations/help",
]
for route in required_routes:
    if route not in route_text:
        failures.append(f"operations route registry missing {route}")

require_text("cmd/agency-config/operations.go", ["IssueCenter", "buildOperationsIssueCenter(page)"])
require_text("cmd/agency-config/operations_operator_issues.go", ["operationsOperatorIssue", "seen", "AdminLink"])

make_targets = [
    "audit-product-language",
    "audit-ui-layout",
    "product-ui-smoke",
    "audit-operations-route-inventory",
    "external-connection-check",
    "adapter-conformance",
    "test-connector-examples",
    "gtfsrt-conformance",
    "release-candidate-check",
    "audit-final-claim-review",
]
makefile = read("Makefile")
for target in make_targets:
    if not re.search(rf"^{re.escape(target)}:", makefile, re.M):
        failures.append(f"Makefile missing target {target}")

site_pages = [
    "site/index.html",
    "site/ui-tour.html",
    "site/check-feeds.html",
    "site/connectors.html",
    "site/connector-support.html",
    "site/readiness.html",
    "site/deploy.html",
]
for page in site_pages:
    require_path(page)
    require_text(page, ["<main", "</main>"])

connector_examples = [
    "examples/connectors/telemetry-csv-replay/connector.json",
    "examples/connectors/telemetry-http-poller/connector.json",
    "examples/connectors/telemetry-webhook-sidecar/connector.json",
    "examples/connectors/predictor-sidecar-stub/connector.json",
    "examples/connectors/validator-allowlist/connector.json",
    "examples/connectors/monitoring-export/connector.json",
    "examples/connectors/consumer-discovery-metadata/connector.json",
]
for manifest in connector_examples:
    require_path(manifest)
    try:
        data = json.loads((ROOT / manifest).read_text(encoding="utf-8"))
    except Exception as exc:
        failures.append(f"{manifest}: invalid JSON: {exc}")
        continue
    if data.get("runtime", {}).get("network", "") == "live-default":
        failures.append(f"{manifest}: live-default connector network mode is not allowed")

require_text("docs/roadmap-status.md", [
    "Operator workflow",
    "GTFS data quality",
    "GTFS-RT usefulness",
    "Connectors",
    "Deployment and observability",
    "Security and redaction",
    "Release gates",
])
require_text("docs/roadmaps/external-connector-runtime-integration/phase-plan.md", [
    "Phase 01 - Runtime Boundary Baseline",
    "Phase 08 - Runtime Roadmap Closeout And Release Gate",
])
require_regex("docs/current-status.md", r"(?i)unsupported claims remain unsupported|does not prove", "bounded claim language")
require_regex("docs/handoffs/latest.md", r"(?i)consumer tracker remain|consumer tracker remains", "consumer tracker boundary")

status_path = ROOT / "docs/evidence/consumer-submissions/status.json"
try:
    status = json.loads(status_path.read_text(encoding="utf-8"))
except Exception as exc:
    failures.append(f"docs/evidence/consumer-submissions/status.json invalid JSON: {exc}")
else:
    targets = status.get("targets", [])
    if len(targets) != 7:
        failures.append(f"consumer tracker target count = {len(targets)}, want 7")
    for target in targets:
        if target.get("status") != "prepared":
            failures.append(f"consumer tracker target {target.get('id') or target.get('name')} status = {target.get('status')!r}, want 'prepared'")

diff = git(["diff", "--name-only", "--", "docs/evidence"])
if diff.returncode == 0 and diff.stdout.strip():
    failures.append("protected evidence paths have tracked diff:\n" + diff.stdout.strip())

branch = git(["branch", "--show-current"])
if branch.returncode == 0 and branch.stdout.strip() not in {"main", "stable"}:
    warnings.append(f"current branch is {branch.stdout.strip()!r}; expected product work to start from main or stable")

gh_pages = git(["rev-parse", "--verify", "gh-pages"])
if gh_pages.returncode == 0:
    listed = git(["ls-tree", "--name-only", "gh-pages"])
    if listed.returncode == 0:
        for item in ["index.html", "ui-tour.html", "connectors.html", "readiness.html"]:
            if item not in listed.stdout.splitlines():
                failures.append(f"gh-pages branch missing {item}")
else:
    warnings.append("local gh-pages branch not present; skipped branch source comparison")

if failures:
    print("product roadmap baseline audit failed:")
    for failure in failures:
        print(f"  {failure}")
    for warning in warnings:
        print(f"WARN: {warning}")
    raise SystemExit(1)

for warning in warnings:
    print(f"WARN: {warning}")
print("product roadmap baseline audit passed")
PY
