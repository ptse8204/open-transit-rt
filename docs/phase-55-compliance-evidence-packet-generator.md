# Phase 55 -- Compliance Evidence Packet Generator

## Status

Complete for the approved local packet generator scope. Execution added a
redaction-safe `.cache` packet generator, audit guard, Make targets, local-only
script tests, and status/handoff updates without creating retained evidence,
contacting consumers, changing consumer statuses, or claiming compliance.

## Goal

Generate deployment-specific readiness/compliance review packets from retained
artifacts and current repository mappings, with strict claim gates and explicit
human review requirements.

The packet generator is a summarizer, not a collector. It must report missing
evidence truthfully and must not infer compliance.

## Scope

- Add a local packet generator that defaults to ignored
  `.cache/compliance-evidence-packet/<UTC timestamp>/` output.
- Add an audit guard that fails on unsafe output, unsupported claim flags, fake
  evidence signals, or consumer tracker drift.
- Map retained inputs to RQ-4A through RQ-4G and the Phase 54 official-source
  mappings.
- Summarize current consumer tracker state without changing it.
- Produce blocker-only output when deployment identity or root is missing.
- Add local-only script tests for safety, redaction, claim-gating, and packet
  contract behavior.

## Non-Goals

- No public unauthenticated API.
- No admin route.
- No migration.
- No DB write.
- No runtime behavior change.
- No live feed fetching or proof collection.
- No consumer contact, target selection, portal automation, submission, or
  status transition.
- No final-root evidence creation or fabrication.
- No `docs/evidence/captured` write by default.
- No compliance, consumer acceptance, consumer ingestion, agency adoption,
  hosted SaaS, production-readiness, vendor-compatibility, SLA/uptime,
  marketplace approval, or production-grade ETA claim.

## Files Likely To Change

- `scripts/generate-compliance-evidence-packet.sh`
- `scripts/audit-compliance-evidence-packet.sh`
- `scripts/test-compliance-evidence-packet.sh`
- `Makefile`
- `docs/phase-55-compliance-evidence-packet-generator.md`
- `docs/evidence/README.md`
- `docs/compliance-evidence-checklist.md`
- `docs/california-readiness-summary.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/repo-gaps.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-55.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/captured/`

## Safety Boundaries

Default output:

```text
.cache/compliance-evidence-packet/<UTC timestamp>/
```

The generator and audit must reject:

- output under `docs/evidence/captured`;
- other evidence-like output paths unless explicitly safe and ignored;
- symlink paths;
- path traversal;
- private logs;
- raw telemetry;
- raw correspondence;
- private portal data;
- DB URLs;
- credentials;
- cookies;
- Authorization headers;
- oversized source files.

The generator may reference retained evidence paths and safe summaries. It must
not copy raw ZIP/protobuf artifacts, raw validator reports, private logs, or
private operator artifacts into the generated packet.

`.cache` diagnostics may be summarized only as private local diagnostics, not
retained compliance evidence.

The OCI pilot packet must remain labeled hosted/operator pilot evidence, not
agency-owned final-root proof.

## Evidence And Claim Boundaries

Generated rows may use statuses such as:

- `present`
- `partial`
- `missing`
- `blocked`
- `pilot_only`
- `needs_review`

Generated rows must not use `compliant` as a status value.

Claim flags must remain false:

- `compliance_claimed`
- `consumer_acceptance_claimed`
- `consumer_ingestion_claimed`
- `agency_adoption_claimed`
- `hosted_saas_claimed`
- `production_readiness_claimed`
- `sla_or_uptime_claimed`
- `vendor_compatibility_claimed`
- `production_grade_eta_claimed`

Default generated wording is limited to readiness/support language. Any wording
beyond that must be operator-supplied, redacted, separately reviewed, and
clearly labeled as human-reviewed text. The generator must not infer or endorse
compliance automatically.

## Implementation Details

Add Make targets:

- `generate-compliance-evidence-packet`
- `audit-compliance-evidence-packet`
- `test-compliance-evidence-packet`

Suggested generator environment:

- `COMPLIANCE_PACKET_DEPLOYMENT_NAME`
- `COMPLIANCE_PACKET_ROOT_URL`
- `COMPLIANCE_PACKET_OUTPUT_DIR`
- `COMPLIANCE_PACKET_FORCE`
- `COMPLIANCE_PACKET_MAX_SOURCE_BYTES`
- `COMPLIANCE_PACKET_HUMAN_REVIEW`
- `COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT`
- optional retained evidence path variables for final-root, validation,
  operations, and consumer artifacts

Blocker-only output, used when deployment identity or root is missing, must
write exactly:

- `blocker.json`
- `blocker.md`
- `manifest.json`
- `manifest.md`

Deployment packet output must write exactly:

- `summary.json`
- `summary.md`
- `readiness-packet.md`
- `evidence-map.json`
- `evidence-map.md`
- `missing-evidence.md`
- `human-review.md`
- `manifest.json`
- `manifest.md`

`summary.json` must include:

- packet version;
- generated timestamp;
- deployment identity and root;
- source evidence path summaries;
- official requirements review date;
- RQ-4A through RQ-4G readiness rows;
- consumer tracker summary;
- missing evidence;
- human review state;
- all-false claim flags.

The audit command must fail when:

- the packet file set is wrong;
- required JSON is invalid;
- any claim flag is true;
- any row status is `compliant`;
- consumer tracker target order or statuses drift;
- artifact directories are no longer README-only in the current repo state;
- unsafe private strings appear in generated files;
- a packet claims final-root, consumer, deployment, operations, vendor, SLA, or
  ETA evidence that is not represented as missing, blocked, pilot-only, or
  needs-review.

Checkpoint strategy:

- `Phase 55 -- Checkpoint 000001: add compliance packet generator plan`
- `Phase 55 -- Checkpoint 000002: add compliance packet generator and audit`
- `Phase 55 -- Checkpoint 000003: add compliance packet tests and validation`
- `Phase 55 -- Checkpoint 000004: close compliance packet handoff and status`

## Tests

Add local-only tests for:

- help output;
- blocker-only exact files;
- deployment packet exact files from local fixtures;
- `.cache` default output;
- rejection of `docs/evidence/captured` output by default;
- symlink, traversal, and evidence-like path rejection;
- redaction scanner failures;
- all claim flags false;
- no `compliant` status value;
- OCI pilot classified as `pilot_only`;
- final-root evidence missing reported truthfully;
- all seven consumers remain `prepared`;
- artifact directories remain README-only;
- bounded large fixture reads and stable manifest generation.

## Performance And Scale Tests

- Cap source summary size.
- Stream checksums for referenced large artifacts when needed.
- Do not copy ZIP, protobuf, raw validation artifacts, raw logs, or dumps into
  the packet.
- Include a large local fixture test proving bounded reads and stable manifest
  generation.

## Docs, Status, And Handoff Updates

Close Phase 55 by updating:

- this phase document;
- `docs/handoffs/phase-55.md`;
- `docs/handoffs/latest.md`;
- `docs/current-status.md`;
- `docs/backlog.md`;
- `docs/open-questions.md`;
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`;
- `docs/repo-gaps.md`;
- `docs/evidence/README.md`;
- relevant readiness/checklist docs.

The handoff must state whether only blocker/draft `.cache` packets were
generated and must explicitly say no retained evidence, consumer statuses,
final-root proof, or compliance claims were created.

## Required Verification Commands

Run and report:

```bash
sh -n scripts/generate-compliance-evidence-packet.sh scripts/audit-compliance-evidence-packet.sh scripts/test-compliance-evidence-packet.sh
./scripts/test-compliance-evidence-packet.sh
make generate-compliance-evidence-packet
COMPLIANCE_PACKET_DIR=.cache/compliance-evidence-packet/<timestamp> make audit-compliance-evidence-packet
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
docker compose -f deploy/docker-compose.yml config
```

The `find` command must print no files for the current Phase 55 state.

## Closure Result

Phase 55 added:

- `scripts/generate-compliance-evidence-packet.sh`
- `scripts/audit-compliance-evidence-packet.sh`
- `scripts/test-compliance-evidence-packet.sh`
- Make targets for generate, audit, and test

The generator defaults to ignored
`.cache/compliance-evidence-packet/<UTC timestamp>/` output. With missing
deployment identity or root URL it writes exactly:

- `blocker.json`
- `blocker.md`
- `manifest.json`
- `manifest.md`

With deployment identity and root URL it writes exactly:

- `summary.json`
- `summary.md`
- `readiness-packet.md`
- `evidence-map.json`
- `evidence-map.md`
- `missing-evidence.md`
- `human-review.md`
- `manifest.json`
- `manifest.md`

The generator is a summarizer only. It reads bounded source metadata, keeps
claim flags false, uses non-compliance statuses only, summarizes the current
consumer tracker as prepared-only, labels OCI evidence as `pilot_only`, and
reports final-root evidence missing unless a real retained final-root path is
configured.

The audit fails on wrong file sets, invalid JSON, true claim flags,
`compliant` status values, unsafe private strings, prepared-only consumer
tracker drift, non-README consumer artifact files, and misleading compliance,
consumer, final-root, deployment, operations, vendor, SLA, marketplace, or ETA
claims.

Phase 55 generated only ignored `.cache` blocker/draft packets during
verification. It did not write retained evidence under `docs/evidence`, did not
change `docs/evidence/captured`, did not change consumer current records or
artifact directories, and did not change
`docs/evidence/consumer-submissions/status.json`.
