#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
import pathlib
import shutil
import subprocess
import os
import sys

ROOT = pathlib.Path(sys.argv[1]).resolve()
BASE = ROOT / ".cache" / "vendor-equivalent-pack-tests" / str(os.getpid())
AUDIT = ROOT / "scripts" / "audit-vendor-equivalent-pack.sh"
SOURCE = ROOT / "docs" / "vendor-equivalent-pack"


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


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(AUDIT), "--help"])
run([str(AUDIT)])

working = BASE / "working"
shutil.copytree(SOURCE, working)
run([str(AUDIT)], env={"VENDOR_EQUIVALENT_PACK_DIR": working})

missing = BASE / "missing"
shutil.copytree(SOURCE, missing)
(missing / "sla-kpi-template.md").unlink()
run([str(AUDIT)], env={"VENDOR_EQUIVALENT_PACK_DIR": missing}, ok=False)

claim = BASE / "claim"
shutil.copytree(SOURCE, claim)
(claim / "procurement-response-template.md").write_text("Open Transit RT is compliant.\n<placeholder>\n", encoding="utf-8")
run([str(AUDIT)], env={"VENDOR_EQUIVALENT_PACK_DIR": claim}, ok=False)

private = BASE / "private"
shutil.copytree(SOURCE, private)
(private / "support-boundaries-template.md").write_text("Authorization: Bearer abcdefghijklmnop\n<placeholder>\n", encoding="utf-8")
run([str(AUDIT)], env={"VENDOR_EQUIVALENT_PACK_DIR": private}, ok=False)

boundary = BASE / "boundary"
shutil.copytree(SOURCE, boundary)
text = (boundary / "README.md").read_text(encoding="utf-8").replace("template material only", "starter text")
(boundary / "README.md").write_text(text, encoding="utf-8")
run([str(AUDIT)], env={"VENDOR_EQUIVALENT_PACK_DIR": boundary}, ok=False)

print("vendor-equivalent pack script tests passed")
PY
