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
BASE = ROOT / ".cache" / "multi-agency-hosting-tests" / str(os.getpid())
SCRIPT = ROOT / "scripts" / "multi-agency-hosting.sh"
EXPECTED_FILES = ["manifest.json", "manifest.md", "summary.json", "summary.md"]
EXPECTED_ROUTES = [
    "/public/agencies/*/feeds.json",
    "/public/agencies/*/gtfs/schedule.zip",
    "/public/agencies/*/gtfsrt/vehicle_positions.pb",
    "/public/agencies/*/gtfsrt/trip_updates.pb",
    "/public/agencies/*/gtfsrt/alerts.pb",
]


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


def load(path):
    return json.loads(pathlib.Path(path).read_text(encoding="utf-8"))


def assert_exact_files(path):
    got = sorted(child.name for child in pathlib.Path(path).iterdir())
    if got != EXPECTED_FILES:
        raise AssertionError(f"{path} files = {got}, want {EXPECTED_FILES}")


def assert_claim_flags_false(data):
    flags = data.get("claim_flags", {})
    if not flags or any(value is not False for value in flags.values()):
        raise AssertionError(f"claim flags not all false: {flags}")


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(SCRIPT), "--help"])

default_out = BASE / "default"
default_proc = run([str(SCRIPT)], env={"OUTPUT_DIR": default_out})
default_path = None
for line in default_proc.stdout.splitlines():
    if "multi-agency hosting diagnostic:" in line:
        default_path = ROOT / line.split(":", 1)[1].strip()
if default_path is None or default_path.resolve() != default_out.resolve():
    raise AssertionError(f"explicit output did not use requested .cache path: {default_proc.stdout}")
assert_exact_files(default_path)

custom = BASE / "custom"
run([str(SCRIPT)], env={"OUTPUT_DIR": custom})
assert_exact_files(custom)
summary = load(custom / "summary.json")
manifest = load(custom / "manifest.json")
if summary["status"] != "passed" or manifest["status"] != "passed":
    raise AssertionError(f"unexpected diagnostic status: {summary['status']} {manifest['status']}")
assert_claim_flags_false(summary)
assert_claim_flags_false(manifest)
if summary["public_route_contract"] != EXPECTED_ROUTES:
    raise AssertionError(summary["public_route_contract"])
if summary["operations_model"]["tenant_restore_into_shared_live_database"] != "blocked":
    raise AssertionError(summary["operations_model"])
for item in summary["caddy_files"]:
    if item["missing_routes"] or item["forbidden_public_edge_matches"]:
        raise AssertionError(item)

run([str(SCRIPT)], env={"OUTPUT_DIR": custom}, ok=False)
run([str(SCRIPT)], env={"OUTPUT_DIR": custom, "FORCE": "true"})

run([str(SCRIPT)], env={"OUTPUT_DIR": "../outside"}, ok=False)
run([str(SCRIPT)], env={"OUTPUT_DIR": ROOT / "docs" / "evidence" / "captured" / "bad"}, ok=False)
run([str(SCRIPT)], env={"OUTPUT_DIR": BASE / "proof" / "bad"}, ok=False)
run([str(SCRIPT)], env={"FORCE": "maybe", "OUTPUT_DIR": BASE / "bad-force"}, ok=False)

symlink_target = BASE / "symlink-target"
symlink_target.mkdir()
symlink_output = BASE / "symlink-output"
symlink_output.symlink_to(symlink_target, target_is_directory=True)
run([str(SCRIPT)], env={"OUTPUT_DIR": symlink_output}, ok=False)

print("multi-agency hosting script tests passed")
PY
