#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
from __future__ import annotations

import os
import pathlib
import shutil
import subprocess
import tempfile
import sys

root = pathlib.Path(sys.argv[1]).resolve()
audit = root / "scripts" / "audit-product-language.sh"

primary_files = [
    "README.md",
    "docs/index.md",
    "docs/tutorials/no-cli-agency-first-run.md",
    "docs/tutorials/video-recording-guide.md",
    "docs/deployment/oci-reference-deployment.md",
    "wiki/README.md",
    "cmd/agency-config/operations.go",
    "cmd/agency-config/operations_navigation.go",
    "cmd/agency-config/operations_route_registry.go",
    "cmd/agency-config/operations_admin.js",
]


def run(args: list[str], env: dict[str, str] | None = None, ok: bool = True) -> subprocess.CompletedProcess[str]:
    merged = os.environ.copy()
    if env:
        merged.update(env)
    proc = subprocess.run(args, cwd=root, env=merged, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if ok and proc.returncode != 0:
        raise AssertionError(f"command failed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    if not ok and proc.returncode == 0:
        raise AssertionError(f"command unexpectedly passed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    return proc


def copy_fixture(dst: pathlib.Path) -> None:
    for relative in primary_files:
        target = dst / relative
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(root / relative, target)
    for source in (root / "site").glob("*.html"):
        target = dst / "site" / source.name
        target.parent.mkdir(parents=True, exist_ok=True)
        shutil.copy2(source, target)
    assets = dst / "site" / "assets"
    assets.mkdir(parents=True, exist_ok=True)
    for source in (root / "site" / "assets").glob("*.vtt"):
        shutil.copy2(source, assets / source.name)


def write(path: pathlib.Path, text: str) -> None:
    path.write_text(text, encoding="utf-8")


run([str(audit), "--help"])
run([str(audit)])

with tempfile.TemporaryDirectory(prefix="open-transit-product-language-") as tmp:
    fixture = pathlib.Path(tmp)
    copy_fixture(fixture)
    env = {"PRODUCT_LANGUAGE_AUDIT_ROOT": str(fixture)}
    run([str(audit)], env=env)

    index = fixture / "site" / "index.html"
    write(index, index.read_text(encoding="utf-8") + "\n<p>Common Next Actions</p>\n")
    run([str(audit)], env=env, ok=False)

with tempfile.TemporaryDirectory(prefix="open-transit-product-language-flags-") as tmp:
    fixture = pathlib.Path(tmp)
    copy_fixture(fixture)
    env = {"PRODUCT_LANGUAGE_AUDIT_ROOT": str(fixture)}
    tour = fixture / "site" / "ui-tour.html"
    write(tour, tour.read_text(encoding="utf-8") + "\n<p>external_evidence_created</p>\n")
    run([str(audit)], env=env, ok=False)

with tempfile.TemporaryDirectory(prefix="open-transit-product-language-console-") as tmp:
    fixture = pathlib.Path(tmp)
    copy_fixture(fixture)
    env = {"PRODUCT_LANGUAGE_AUDIT_ROOT": str(fixture)}
    registry = fixture / "cmd" / "agency-config" / "operations_route_registry.go"
    write(registry, registry.read_text(encoding="utf-8") + '\nvar badCopy = "What stays technical?"\n')
    run([str(audit)], env=env, ok=False)

print("product language audit tests passed")
PY
