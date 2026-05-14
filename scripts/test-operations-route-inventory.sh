#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
import os
import pathlib
import shutil
import subprocess
import sys
import tempfile

ROOT = pathlib.Path(sys.argv[1]).resolve()
AUDIT = ROOT / "scripts" / "audit-operations-route-inventory.sh"


def run(args, env=None, ok=True):
    merged = os.environ.copy()
    if env:
        merged.update({key: str(value) for key, value in env.items()})
    proc = subprocess.run(args, cwd=ROOT, env=merged, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if ok and proc.returncode != 0:
        raise AssertionError(f"command failed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    if not ok and proc.returncode == 0:
        raise AssertionError(f"command unexpectedly passed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    return proc


def copy_required_files(dst):
    files = [
        "cmd/agency-config/main.go",
        "cmd/agency-config/operations.go",
        "cmd/agency-config/operations_navigation.go",
        "cmd/agency-config/operations_admin.js",
        "docs/phase-90-control-plane-final-status.md",
        "README.md",
        "wiki/README.md",
        "wiki/small-agency-quick-start.md",
    ]
    for rel in files:
        target = dst / rel
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(ROOT / rel, target)


def append_doc_routes(dst):
    routes = [
        "/admin/operations/gtfs-workbench",
        "/admin/operations/realtime",
        "/admin/operations/validation-center",
        "/admin/operations/connectors/workbench",
        "/admin/operations/prediction-lab",
        "/admin/operations/access",
        "/admin/operations/audit",
    ]
    for rel in ["README.md", "wiki/README.md", "wiki/small-agency-quick-start.md"]:
        path = dst / rel
        path.write_text(path.read_text(encoding="utf-8") + "\n" + "\n".join(routes) + "\n", encoding="utf-8")


run([str(AUDIT), "--help"])
run([str(AUDIT)])

with tempfile.TemporaryDirectory(prefix="open-transit-route-audit-") as tmp:
    fixture = pathlib.Path(tmp)
    copy_required_files(fixture)
    append_doc_routes(fixture)
    env = {
        "OPERATIONS_ROUTE_AUDIT_ROOT": fixture,
        "OPERATIONS_ROUTE_AUDIT_STRICT_DOCS": "true",
    }
    run([str(AUDIT)], env=env)

    nav = fixture / "cmd" / "agency-config" / "operations_navigation.go"
    nav.write_text(nav.read_text(encoding="utf-8").replace(
        '{Label: "GTFS Workbench", Href: "/admin/operations/gtfs-workbench", Section: "gtfs-workbench"},\n',
        "",
    ), encoding="utf-8")
    run([str(AUDIT)], env=env, ok=False)

with tempfile.TemporaryDirectory(prefix="open-transit-route-doc-audit-") as tmp:
    fixture = pathlib.Path(tmp)
    copy_required_files(fixture)
    append_doc_routes(fixture)
    for rel in ["README.md", "wiki/README.md", "wiki/small-agency-quick-start.md"]:
        path = fixture / rel
        path.write_text(path.read_text(encoding="utf-8").replace("/admin/operations/realtime", ""), encoding="utf-8")
    run([str(AUDIT)], env={
        "OPERATIONS_ROUTE_AUDIT_ROOT": fixture,
        "OPERATIONS_ROUTE_AUDIT_STRICT_DOCS": "true",
    }, ok=False)

print("operations route inventory audit tests passed")
PY
