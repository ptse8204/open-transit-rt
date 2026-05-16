#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
import json
import os
import pathlib
import shutil
import subprocess
import sys

ROOT = pathlib.Path(sys.argv[1]).resolve()
BASE = ROOT / ".cache" / "product-acceptance-tests" / str(os.getpid())
AUDIT = ROOT / "scripts" / "audit-product-acceptance.sh"
EXPECTED_TARGETS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]


def run(args, env=None, ok=True):
    merged = os.environ.copy()
    if env:
        merged.update({k: str(v) for k, v in env.items()})
    proc = subprocess.run(args, cwd=ROOT, env=merged, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if ok and proc.returncode != 0:
        raise AssertionError(f"command failed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    if not ok and proc.returncode == 0:
        raise AssertionError(f"command unexpectedly passed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    return proc


def write(path, text):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(text, encoding="utf-8")


def write_status(root, status="prepared", targets=EXPECTED_TARGETS):
    path = root / "docs" / "evidence" / "consumer-submissions" / "status.json"
    path.parent.mkdir(parents=True, exist_ok=True)
    data = {"targets": [{"target": target, "status": status} for target in targets]}
    path.write_text(json.dumps(data, indent=2) + "\n", encoding="utf-8")


def make_fixture(name, mutate=None):
    root = BASE / name
    if root.exists():
        shutil.rmtree(root)
    write(
        root / "README.md",
        """# Open Transit RT

Open Transit RT helps small agencies evaluate a self-hosted GTFS and GTFS Realtime path.

## Start The Software UI

```bash
make check
make agency-app-up
```

Open the private Operations Console at `http://localhost:8080/admin/operations`.
This local evaluation does not prove agency approval or consumer acceptance.
""",
    )
    write(
        root / "wiki" / "README.md",
        """# Open Transit RT Wiki

## Start Here

1. [Small Agency Quick Start](small-agency-quick-start.md)
2. [Browser-First Setup](browser-first-setup.md)
3. [Operations Console Tour](operations-console-tour.md)
""",
    )
    write(root / "wiki" / "small-agency-quick-start.md", "# Small Agency Quick Start\n")
    write(root / "wiki" / "browser-first-setup.md", "# Browser-First Setup\n")
    write(root / "wiki" / "operations-console-tour.md", "# Operations Console Tour\n")
    write(
        root / "docs" / "README.md",
        """# Documentation Home

## Public User Docs

Task guides first.

## Maintainer Docs And History

Detailed phase history belongs here.
""",
    )
    write(root / "docs" / "tutorials" / "small-agency-acceptance-script.md", "# Small-Agency Acceptance Script\n")
    write(root / "wiki" / "connector-cookbook.md", "# Connector Cookbook\n\n## Practical Recipes\n\nTelemetry recipe.\n")
    write(
        root / "wiki" / "calitp-readiness-plain-english.md",
        "# CAL-ITP Readiness Plain English\n\n| Area | UI signal you can review | Missing deployment evidence before stronger claims |\n| --- | --- | --- |\n| Static GTFS | Browser import. | Final public proof. |\n",
    )
    write(
        root / "docs" / "requirements-calitp-compliance.md",
        "# Requirements\n\nSoftware capability exists for GTFS import/publication, stable feed paths, Vehicle Positions, Trip Updates, Alerts, validation workflows, feed health, readiness workflows, telemetry ingest, and connector boundaries.\n\nExternal proof tracks are optional and authorization-gated.\n",
    )
    write_status(root)
    for protected in (
        root / "docs" / "evidence" / "captured",
        root / "docs" / "evidence" / "consumer-submissions" / "current",
        root / "docs" / "evidence" / "consumer-submissions" / "artifacts",
        root / "docs" / "evidence" / "consumer-submissions" / "packets",
    ):
        protected.mkdir(parents=True, exist_ok=True)
    if mutate:
        mutate(root)
    return root


def audit_env(root):
    return {
        "PRODUCT_ACCEPTANCE_ROOT": root,
        "PRODUCT_ACCEPTANCE_PUBLIC_DOCS": ":".join(
            [
                "README.md",
                "wiki/README.md",
                "wiki/small-agency-quick-start.md",
                "wiki/browser-first-setup.md",
                "wiki/operations-console-tour.md",
                "wiki/connector-cookbook.md",
                "wiki/calitp-readiness-plain-english.md",
                "docs/README.md",
                "docs/tutorials/small-agency-acceptance-script.md",
                "docs/requirements-calitp-compliance.md",
            ]
        ),
        "PRODUCT_ACCEPTANCE_SKIP_GIT_STATUS": "true",
    }


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(AUDIT), "--help"])
run([str(AUDIT)])

passing = make_fixture("passing")
run([str(AUDIT)], env=audit_env(passing))

missing_quick_start = make_fixture("missing-quick-start", lambda root: (root / "wiki" / "small-agency-quick-start.md").unlink())
run([str(AUDIT)], env=audit_env(missing_quick_start), ok=False)

phase_led = make_fixture(
    "phase-led",
    lambda root: write(root / "README.md", "# Open Transit RT\n\nPhase 69 checkpoint history first.\n\nmake agency-app-up\n\nOperations Console `/admin/operations`.\n"),
)
run([str(AUDIT)], env=audit_env(phase_led), ok=False)

forbidden_claim = make_fixture(
    "forbidden-claim",
    lambda root: write(root / "README.md", "# Open Transit RT\n\nmake agency-app-up\n\nOperations Console `/admin/operations`.\n\nOpen Transit RT is production-ready.\n"),
)
run([str(AUDIT)], env=audit_env(forbidden_claim), ok=False)

bad_status = make_fixture("bad-status", lambda root: write_status(root, status="accepted"))
run([str(AUDIT)], env=audit_env(bad_status), ok=False)

missing_matrix = make_fixture(
    "missing-matrix",
    lambda root: write(root / "wiki" / "calitp-readiness-plain-english.md", "# CAL-ITP Readiness Plain English\n"),
)
run([str(AUDIT)], env=audit_env(missing_matrix), ok=False)

print("product acceptance audit script tests passed")
PY
