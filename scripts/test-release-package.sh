#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
import json
import os
import pathlib
import re
import shutil
import subprocess
import sys
import tarfile

ROOT = pathlib.Path(sys.argv[1]).resolve()
BASE = ROOT / ".cache" / "release-package-tests" / str(os.getpid())
GENERATE = ROOT / "scripts" / "release-package.sh"
AUDIT = ROOT / "scripts" / "audit-release-package.sh"


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


def assert_flags_false(data):
    flags = data.get("claim_flags", {})
    if not flags or any(value is not False for value in flags.values()):
        raise AssertionError(f"claim flags are not all false: {flags}")


def protected_archive_entries(archive_path):
    pattern = re.compile(r"(^|/)docs/evidence/(captured|consumer-submissions/(status\.json|current|artifacts|packets))")
    with tarfile.open(archive_path, "r:gz") as archive:
        return [member.name for member in archive.getmembers() if pattern.search(member.name)]


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(GENERATE), "--help"])
run([str(AUDIT), "--help"])

package = BASE / "package"
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": package,
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
})
summary = load(package / "summary.json")
manifest = load(package / "manifest.json")
provenance = load(package / "provenance.json")
sbom = load(package / "sbom.json")
image = load(package / "image.json")
for data in (summary, manifest, provenance, sbom, image):
    assert_flags_false(data)
if summary["version"] != "v0.0.0-test" or summary["source_archive"] != "artifacts/open-transit-rt-v0.0.0-test.source.tar.gz":
    raise AssertionError(summary)
if not (package / "artifacts" / "open-transit-rt-v0.0.0-test.source.tar.gz").exists():
    raise AssertionError("source archive missing")
protected_hits = protected_archive_entries(package / "artifacts" / "open-transit-rt-v0.0.0-test.source.tar.gz")
if protected_hits:
    raise AssertionError(f"source archive contains protected paths: {protected_hits[:10]}")
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": package})

run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": package,
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": package,
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
    "RELEASE_PACKAGE_FORCE": "true",
})
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": package})

run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "bad/version",
    "RELEASE_PACKAGE_OUTPUT_DIR": BASE / "bad-version",
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": "../outside",
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": ROOT / "docs" / "evidence" / "captured" / "release",
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": BASE / "proof" / "release",
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)

symlink_target = BASE / "symlink-target"
symlink_target.mkdir()
symlink_output = BASE / "symlink-output"
symlink_output.symlink_to(symlink_target, target_is_directory=True)
run([str(GENERATE)], env={
    "RELEASE_PACKAGE_VERSION": "v0.0.0-test",
    "RELEASE_PACKAGE_OUTPUT_DIR": symlink_output,
    "RELEASE_PACKAGE_ALLOW_DIRTY": "true",
}, ok=False)

extra = BASE / "extra-file"
shutil.copytree(package, extra)
(extra / "extra.txt").write_text("extra\n", encoding="utf-8")
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": extra}, ok=False)

checksum_drift = BASE / "checksum-drift"
shutil.copytree(package, checksum_drift)
(checksum_drift / "summary.md").write_text((checksum_drift / "summary.md").read_text(encoding="utf-8") + "\nchanged\n", encoding="utf-8")
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": checksum_drift}, ok=False)

true_flag = BASE / "true-flag"
shutil.copytree(package, true_flag)
data = load(true_flag / "summary.json")
data["claim_flags"]["hosted_service_claimed"] = True
(true_flag / "summary.json").write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": true_flag}, ok=False)

unsafe = BASE / "unsafe"
shutil.copytree(package, unsafe)
(unsafe / "summary.md").write_text("Open Transit RT is compliant.\n", encoding="utf-8")
run([str(AUDIT)], env={"RELEASE_PACKAGE_DIR": unsafe}, ok=False)

print("release package script tests passed")
PY
