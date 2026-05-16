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

ROOT = pathlib.Path(__import__("sys").argv[1]).resolve()
BASE = ROOT / ".cache" / "compliance-evidence-packet-tests" / str(os.getpid())
GENERATE = ROOT / "scripts" / "generate-compliance-evidence-packet.sh"
AUDIT = ROOT / "scripts" / "audit-compliance-evidence-packet.sh"

EXPECTED_BLOCKER = ["blocker.json", "blocker.md", "manifest.json", "manifest.md"]
EXPECTED_DEPLOYMENT = [
    "summary.json",
    "summary.md",
    "readiness-packet.md",
    "evidence-map.json",
    "evidence-map.md",
    "missing-evidence.md",
    "human-review.md",
    "manifest.json",
    "manifest.md",
]
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


def assert_files_exact(path, expected):
    names = sorted(p.name for p in pathlib.Path(path).iterdir())
    if names != sorted(expected):
        raise AssertionError(f"{path} files mismatch: {names} != {sorted(expected)}")


def load(path):
    return json.loads(pathlib.Path(path).read_text(encoding="utf-8"))


def assert_claim_flags_false(data):
    flags = data.get("claim_flags", {})
    if not flags or any(value is not False for value in flags.values()):
        raise AssertionError(f"claim flags are not all false: {flags}")


def write_tracker(path, status="prepared", targets=EXPECTED_TARGETS):
    data = {
        "tracker_last_reviewed_timestamp": "2026-05-09T00:00:00Z",
        "targets": [
            {
                "target": target,
                "status": status,
                "packet_path": f"packets/{target.lower().replace(' ', '-')}/README.md",
            }
            for target in targets
        ],
    }
    pathlib.Path(path).write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def write_artifact_dirs(path, extra=False):
    base = pathlib.Path(path)
    for target in EXPECTED_TARGETS:
        folder = base / target.lower().replace(" ", "-").replace(".", "")
        folder.mkdir(parents=True, exist_ok=True)
        (folder / "README.md").write_text("# Artifact intake\n", encoding="utf-8")
        if extra and target == "Google Maps":
            (folder / "receipt.txt").write_text("not real evidence\n", encoding="utf-8")


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(GENERATE), "--help"])
run([str(AUDIT), "--help"])

blocker = BASE / "blocker"
run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": blocker})
assert_files_exact(blocker, EXPECTED_BLOCKER)
blocker_json = load(blocker / "blocker.json")
assert blocker_json["mode"] == "blocker", blocker_json
assert_claim_flags_false(blocker_json)
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": blocker})

default_out = BASE / "default"
default_proc = run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": default_out})
default_path = None
for line in default_proc.stdout.splitlines():
    if "compliance evidence blocker packet:" in line:
        default_path = line.split(":", 1)[1].strip()
if pathlib.Path(default_path).resolve() != default_out.resolve():
    raise AssertionError(f"explicit output did not use requested .cache path: {default_proc.stdout}")

deployment = BASE / "deployment"
run([str(GENERATE)], env={
    "COMPLIANCE_PACKET_DEPLOYMENT_NAME": "Local Test Deployment",
    "COMPLIANCE_PACKET_ROOT_URL": "http://127.0.0.1:8080",
    "COMPLIANCE_PACKET_OUTPUT_DIR": deployment,
})
assert_files_exact(deployment, EXPECTED_DEPLOYMENT)
summary = load(deployment / "summary.json")
manifest = load(deployment / "manifest.json")
assert_claim_flags_false(summary)
assert_claim_flags_false(manifest)
statuses = {row["requirement"]: row["status"] for row in summary["readiness_rows"]}
if any(status == "compliant" for status in statuses.values()):
    raise AssertionError(statuses)
if statuses.get("RQ-4A") != "pilot_only":
    raise AssertionError(f"OCI pilot not classified pilot_only: {statuses}")
if not any(item["area"] == "final_root" and item["status"] == "missing" for item in summary["missing_evidence"]):
    raise AssertionError("final-root missing evidence not reported")
if not summary["consumer_tracker_summary"]["all_expected_targets_prepared"]:
    raise AssertionError("consumer tracker was not recorded prepared-only")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": deployment})

run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": ROOT / "docs" / "evidence" / "captured" / "bad"}, ok=False)
run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": "../outside"}, ok=False)
run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": BASE / "docs" / "evidence" / "bad"}, ok=False)

symlink_target = BASE / "symlink-target"
symlink_target.mkdir()
symlink_output = BASE / "symlink-output"
symlink_output.symlink_to(symlink_target, target_is_directory=True)
run([str(GENERATE)], env={"COMPLIANCE_PACKET_OUTPUT_DIR": symlink_output}, ok=False)

unsafe_source = BASE / "unsafe-source.txt"
unsafe_source.write_text("Authorization: Bearer abcdefghijklmnop\n", encoding="utf-8")
run([str(GENERATE)], env={
    "COMPLIANCE_PACKET_DEPLOYMENT_NAME": "Unsafe",
    "COMPLIANCE_PACKET_ROOT_URL": "http://127.0.0.1:8080",
    "COMPLIANCE_PACKET_OUTPUT_DIR": BASE / "unsafe-source-out",
    "COMPLIANCE_PACKET_VALIDATION_EVIDENCE": unsafe_source,
}, ok=False)
run([str(GENERATE)], env={
    "COMPLIANCE_PACKET_DEPLOYMENT_NAME": "Unsafe review",
    "COMPLIANCE_PACKET_ROOT_URL": "http://127.0.0.1:8080",
    "COMPLIANCE_PACKET_OUTPUT_DIR": BASE / "unsafe-review-out",
    "COMPLIANCE_PACKET_HUMAN_REVIEW": "true",
    "COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT": "password: abcdefghijklmnop",
}, ok=False)

bad_files = BASE / "bad-files"
shutil.copytree(deployment, bad_files)
(bad_files / "extra.txt").write_text("extra\n", encoding="utf-8")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": bad_files}, ok=False)

invalid_json = BASE / "invalid-json"
shutil.copytree(deployment, invalid_json)
(invalid_json / "summary.json").write_text("{not json\n", encoding="utf-8")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": invalid_json}, ok=False)

true_flag = BASE / "true-flag"
shutil.copytree(deployment, true_flag)
d = load(true_flag / "summary.json")
d["claim_flags"]["compliance_claimed"] = True
(true_flag / "summary.json").write_text(json.dumps(d, indent=2, sort_keys=True) + "\n", encoding="utf-8")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": true_flag}, ok=False)

compliant = BASE / "compliant"
shutil.copytree(deployment, compliant)
d = load(compliant / "summary.json")
d["readiness_rows"][0]["status"] = "compliant"
(compliant / "summary.json").write_text(json.dumps(d, indent=2, sort_keys=True) + "\n", encoding="utf-8")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": compliant}, ok=False)

misleading = BASE / "misleading"
shutil.copytree(deployment, misleading)
(misleading / "summary.md").write_text("Open Transit RT is CAL-ITP compliant.\n", encoding="utf-8")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": misleading}, ok=False)

large = BASE / "large"
large.mkdir()
(large / "huge.txt").write_bytes(b"x" * (2 * 1024 * 1024))
large_out = BASE / "large-out"
run([str(GENERATE)], env={
    "COMPLIANCE_PACKET_DEPLOYMENT_NAME": "Large Fixture",
    "COMPLIANCE_PACKET_ROOT_URL": "http://127.0.0.1:8080",
    "COMPLIANCE_PACKET_OUTPUT_DIR": large_out,
    "COMPLIANCE_PACKET_OPERATIONS_EVIDENCE": large / "huge.txt",
    "COMPLIANCE_PACKET_MAX_SOURCE_BYTES": "1024",
})
large_summary = load(large_out / "summary.json")
ops = next(source for source in large_summary["source_evidence_path_summaries"] if source["label"] == "operations")
if ops.get("bytes") != 2 * 1024 * 1024 or "sha256" not in ops:
    raise AssertionError(f"large fixture summary is not bounded/stable: {ops}")
run([str(AUDIT)], env={"COMPLIANCE_PACKET_DIR": large_out})

fixture_status = BASE / "status.json"
write_tracker(fixture_status)
fixture_artifacts = BASE / "artifacts-ok"
write_artifact_dirs(fixture_artifacts)
run([str(AUDIT)], env={
    "COMPLIANCE_PACKET_DIR": deployment,
    "COMPLIANCE_CONSUMER_STATUS_PATH": fixture_status,
    "COMPLIANCE_CONSUMER_ARTIFACTS_DIR": fixture_artifacts,
})

drift_status = BASE / "status-drift.json"
write_tracker(drift_status, status="accepted")
run([str(AUDIT)], env={
    "COMPLIANCE_PACKET_DIR": deployment,
    "COMPLIANCE_CONSUMER_STATUS_PATH": drift_status,
    "COMPLIANCE_CONSUMER_ARTIFACTS_DIR": fixture_artifacts,
}, ok=False)

bad_artifacts = BASE / "artifacts-bad"
write_artifact_dirs(bad_artifacts, extra=True)
run([str(AUDIT)], env={
    "COMPLIANCE_PACKET_DIR": deployment,
    "COMPLIANCE_CONSUMER_STATUS_PATH": fixture_status,
    "COMPLIANCE_CONSUMER_ARTIFACTS_DIR": bad_artifacts,
}, ok=False)

print("compliance evidence packet script tests passed")
PY
