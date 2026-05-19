#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  scripts/product-ui-smoke.sh [--help]

Safety:
  Read-only local product UI smoke. It renders private Operations Console
  routes through Go tests with reference-deployment-style settings and checks
  entry-point docs for local and self-hosted guidance. It does not start
  services, fetch public URLs, contact outside systems, write evidence, or
  change consumer statuses.
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

go test ./cmd/agency-config -run TestOperationsConsoleReferenceDeploymentProductSmoke

python3 - "$ROOT_DIR" <<'PY'
from __future__ import annotations

import pathlib
import sys

root = pathlib.Path(sys.argv[1]).resolve()

checks = {
    "README.md": [
        "self-hosted",
        "https://ptse8204.github.io/open-transit-rt/deploy.html",
        "https://ptse8204.github.io/open-transit-rt/ui-tour.html",
        "/admin/operations",
    ],
    "docs/index.md": [
        "Self-Hosted Reference Path",
        "Operator Workflow Tour",
        "/admin/operations",
    ],
    "docs/deployment/oci-reference-deployment.md": [
        "self-hosted deployment pattern",
        "Operator Workflow Tutorial",
        "/admin/operations",
        "/admin/operations/feed-health",
        "not evidence",
    ],
    "wiki/README.md": [
        "self-hosted path",
        "Operations Console Tour",
        "Self-hosted site guide",
        "/admin/operations",
    ],
}

failures: list[str] = []
for relative, needles in checks.items():
    text = (root / relative).read_text(encoding="utf-8")
    for needle in needles:
        if needle not in text:
            failures.append(f"{relative}: missing {needle!r}")

for relative in ["site/deploy.html", "site/ui-tour.html"]:
    if not (root / relative).is_file():
        failures.append(f"{relative}: missing public site source")

if failures:
    print("product UI smoke failed:")
    for failure in failures:
        print(f"  {failure}")
    raise SystemExit(1)

print("product UI smoke passed")
PY
