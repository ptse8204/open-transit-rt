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
BASE = ROOT / ".cache" / "final-claim-review-tests" / str(os.getpid())
AUDIT = ROOT / "scripts" / "audit-final-claim-review.sh"
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


def write_docs(base, extra_doc_text="", phase_missing_section=None):
    docs = []
    common = """# Public Status

Open Transit RT is an open-source backend project.
It does not claim CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support, SLA coverage, marketplace/vendor equivalence, production-grade ETA quality, or universal production readiness.
The OCI DuckDNS pilot is not agency-owned final-root proof.
All seven consumer and aggregator targets remain `prepared` only unless retained target-originated evidence supports a change.

## Reviewed Text

This section is intentionally neutral.
The audit should still catch direct positive claims here.
This line keeps mutation text outside boundary-context windows.
"""
    for name in ("README.md", "docs-status.md"):
        path = base / name
        path.write_text(common + extra_doc_text, encoding="utf-8")
        docs.append(path)
    sections = [
        "## Final Claim Review",
        "## Claim-To-Evidence Table",
        "## Unsupported Claims",
        "## Official Requirements Context",
        "## Retained Evidence Boundary",
        "## Maintainer Signoff",
        "## Execution Closeout",
    ]
    if phase_missing_section:
        sections.remove(phase_missing_section)
    phase = base / "phase-60.md"
    phase.write_text(
        "# Phase 60 -- Final Claim Review And Public Closeout\n\n"
        "## Status\n\nComplete\n\n"
        + "\n\n".join(f"{section}\n\nReviewed." for section in sections)
        + "\n",
        encoding="utf-8",
    )
    return docs, phase


def write_tracker(path, status="prepared", targets=EXPECTED_TARGETS):
    data = {
        "targets": [
            {
                "target": target,
                "status": status,
                "packet_path": f"docs/evidence/consumer-submissions/packets/{target.lower().replace(' ', '-')}/README.md",
            }
            for target in targets
        ]
    }
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def write_artifacts(path, extra=False):
    for target in EXPECTED_TARGETS:
        folder = path / target.lower().replace(" ", "-").replace(".", "")
        folder.mkdir(parents=True, exist_ok=True)
        (folder / "README.md").write_text("# Artifact intake\n", encoding="utf-8")
        if extra and target == "Google Maps":
            (folder / "receipt.txt").write_text("not retained evidence\n", encoding="utf-8")


def audit_env(docs, phase, status, artifacts):
    return {
        "FINAL_CLAIM_REVIEW_DOCS": ":".join(str(path) for path in docs),
        "FINAL_CLAIM_REVIEW_PHASE_DOC": phase,
        "FINAL_CLAIM_REVIEW_STATUS_PATH": status,
        "FINAL_CLAIM_REVIEW_ARTIFACTS_DIR": artifacts,
    }


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(AUDIT), "--help"])

passing = BASE / "passing"
passing.mkdir()
docs, phase = write_docs(passing)
status = passing / "status.json"
artifacts = passing / "artifacts"
write_tracker(status)
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts))

unsupported = BASE / "unsupported"
unsupported.mkdir()
docs, phase = write_docs(unsupported, "\nOpen Transit RT is CAL-ITP compliant.\n")
status = unsupported / "status.json"
artifacts = unsupported / "artifacts"
write_tracker(status)
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts), ok=False)

unsafe = BASE / "unsafe"
unsafe.mkdir()
docs, phase = write_docs(unsafe, "\nAuthorization: Bearer abcdefghijklmnop\n")
status = unsafe / "status.json"
artifacts = unsafe / "artifacts"
write_tracker(status)
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts), ok=False)

tracker_drift = BASE / "tracker-drift"
tracker_drift.mkdir()
docs, phase = write_docs(tracker_drift)
status = tracker_drift / "status.json"
artifacts = tracker_drift / "artifacts"
write_tracker(status, status="accepted")
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts), ok=False)

missing_section = BASE / "missing-section"
missing_section.mkdir()
docs, phase = write_docs(missing_section, phase_missing_section="## Retained Evidence Boundary")
status = missing_section / "status.json"
artifacts = missing_section / "artifacts"
write_tracker(status)
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts), ok=False)

bad_artifacts = BASE / "bad-artifacts"
bad_artifacts.mkdir()
docs, phase = write_docs(bad_artifacts)
status = bad_artifacts / "status.json"
artifacts = bad_artifacts / "artifacts"
write_tracker(status)
write_artifacts(artifacts, extra=True)
run([str(AUDIT)], env=audit_env(docs, phase, status, artifacts), ok=False)

boundary_missing = BASE / "boundary-missing"
boundary_missing.mkdir()
docs, phase = write_docs(boundary_missing)
docs[0].write_text("# Public Status\n\nOpen Transit RT is an open-source backend project.\n", encoding="utf-8")
status = boundary_missing / "status.json"
artifacts = boundary_missing / "artifacts"
write_tracker(status)
write_artifacts(artifacts)
run([str(AUDIT)], env=audit_env([docs[0]], phase, status, artifacts), ok=False)

print("final claim review script tests passed")
PY
