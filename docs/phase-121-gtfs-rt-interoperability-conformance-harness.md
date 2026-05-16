# Phase 121 -- GTFS-RT Interoperability Conformance Harness

## Goal

Add GTFS-RT interoperability conformance surfaces for common open-source
consumer expectations and validator differences using synthetic/local data.

Phase 121 is not a consumer acceptance, compliance, production readiness,
vendor compatibility, hardware certification, SLA/uptime, or ETA-quality proof
phase.

## Current Repo Context

- Phase 120 added redaction-safe Vehicle Positions aggregate review summaries.
- Existing validator helpers can run server-owned allowlisted static and
  realtime validators when configured.
- Existing adapter conformance covers connector manifests and synthetic
  telemetry/prediction/validator/monitoring fixtures, but there is no small
  offline GTFS-RT protobuf interoperability harness for local artifacts.

## Scope

- Add a local/synthetic GTFS-RT conformance harness for Vehicle Positions, Trip
  Updates, and Alerts protobuf artifacts.
- Check common interoperability expectations such as parseability, header
  shape, entity IDs, feed-specific entity payloads, and conservative trip
  descriptor identity.
- Keep the harness offline by default; it must not contact consumers, write
  evidence, mutate statuses, or claim validator-clean/compliance outcomes.

## Protected Paths

Do not modify, reformat, delete, stage, or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

The consumer tracker must remain exactly seven targets in order and all
`prepared`.

## Deliverables

- Offline GTFS-RT conformance library/CLI or equivalent harness.
- Tests covering positive and negative synthetic protobuf cases.
- Makefile or script entry point if useful.
- `docs/handoffs/phase-121.md`
- Source-of-truth status updates for Phase 121 closeout.

## Implementation Plan

1. Add this Phase 121 plan and commit checkpoint 000001.
2. Implement a small offline conformance harness for local protobuf artifacts
   and add test coverage.
3. Run relevant GTFS-RT, connector, claim-boundary, and baseline validation;
   patch repo-caused failures.
4. Close Phase 121 with handoff/status docs and continue immediately to
   Phase 122.

## Checkpoint Plan

- `Phase 121 -- Checkpoint 000001: add gtfs rt interoperability conformance harness plan`
- `Phase 121 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 121 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 121 -- Checkpoint 000004: close gtfs rt interoperability conformance harness review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 121 -- Checkpoint 000001: add GTFS-RT interoperability conformance
harness plan.

Goal status:
Active. Phase 120 is closed and Phase 121 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-121-gtfs-rt-interoperability-conformance-harness.md`.

Validation run:
Initial inspection reviewed the Phase 121 roadmap prompt, validation
boundaries, adapter conformance command, validator helper scripts, GTFS-RT
feed packages, Makefile targets, and synthetic fixture layout.

Blocked checks:
Implementation, tests, connector validation, and closeout validation are
scheduled for later Phase 121 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
compliance, adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, and real-world ETA
accuracy claims.

Security/auth status:
The plan is offline/local. It does not change route auth, CSRF behavior,
credential handling, token handling, public exposure, private payload handling,
or operator command behavior.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 121 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 121 does not plan
visual UX changes.

Master review:
Approved. The plan keeps conformance work offline, synthetic/local, and
claim-bounded.

Required edits:
Commit checkpoint 000001, then implement the scoped conformance harness.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 121 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 121 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 121 implemented an offline GTFS-RT interoperability conformance
harness.

Sub-agents used or simulated:
The agent thread limit still prevents new real sub-agents. Context / Repo
Truth, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles were simulated by the Master Agent.

Changed files:
`internal/gtfsrtconformance/conformance.go`;
`internal/gtfsrtconformance/conformance_test.go`;
`cmd/gtfsrt-conformance/main.go`;
`cmd/gtfsrt-conformance/main_test.go`;
`Makefile`;
`docs/phase-121-gtfs-rt-interoperability-conformance-harness.md`.

Validation run:
Passed:

- `gofmt -w internal/gtfsrtconformance cmd/gtfsrt-conformance`
- `go test ./internal/gtfsrtconformance ./cmd/gtfsrt-conformance`
- `make gtfsrt-conformance`
- `git diff --check`
- `make check`
- protected-path git status check

Blocked checks:
Full Phase 121 validation, connector validation, and claim audits are
scheduled for checkpoint 000003.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
The harness reports all-false claim flags and explicitly says offline
conformance diagnostics do not prove compliance, consumer acceptance,
production readiness, SLA/uptime, vendor compatibility, hardware
certification, production-grade ETA quality, or real-world ETA accuracy.

Security/auth status:
The harness is offline and reads local protobuf files only. It rejects
evidence-like paths, enforces a 10 MiB payload cap, does not fetch URLs, does
not contact consumers, does not write files, and does not expose credentials or
private payloads in tests.

Data/migration status:
No migration, schema, durable state, runtime dependency, or Go module change
was added.

Release/publication status:
The public rc1 prerelease remains published. Phase 121 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. The changed
code is post-rc1 current-source hardening and is not part of the published rc1
tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. CP000002 did not change
visual UI templates.

Master review:
Approved. The new harness checks parseability, header/version/timestamp,
FULL_DATASET profile, unique entity IDs, Vehicle Positions payload shape, Trip
Updates trip identity/start-date/stop-time coverage, and Alerts header text,
informed entities, and active periods using synthetic/local data.

Required edits:
Commit checkpoint 000002, then run Phase 121 validation and patch any
repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 validation.

Next checkpoint:
Phase 121 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 121 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Phase 121 conformance harness validation passed.

Sub-agents used or simulated:
Context / Repo Truth, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles were simulated by the Master Agent because the
agent thread limit prevents new real sub-agents.

Changed files:
`docs/phase-121-gtfs-rt-interoperability-conformance-harness.md`.

Validation run:
Passed:

- `git status --short`
- `git diff --check`
- `make check`
- `make gtfsrt-conformance`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json`
- `scripts/check-consumer-tracker.sh`
- protected-path git status check

Blocked checks:
No Phase 121 validation blocker remains.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. Phase 121 remains bounded to
offline/local GTFS-RT conformance diagnostics and makes no stable release
readiness, production readiness, compliance, adoption, agency approval,
consumer acceptance, consumer ingestion/listing/display, final-root readiness,
hosted service, paid support, SLA/uptime, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality, or
real-world ETA accuracy claim.

Security/auth status:
No route auth, CSRF, credential handling, token handling, public exposure,
private payload handling, or operator command behavior changed. The CLI reads
local files only, rejects evidence-like paths, and writes no files.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was
added.

Release/publication status:
The public rc1 prerelease remains published. Phase 121 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. Phase 121
current-source hardening is not part of the published rc1 tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 121 did not change
visual UI templates.

Master review:
Approved. Validation found no required patch after the conformance harness
implementation.

Required edits:
Commit checkpoint 000003, then close Phase 121 with handoff and status docs.

Decision:
Proceed to checkpoint 000003 commit and Phase 121 closeout.

Next checkpoint:
Phase 121 -- Checkpoint 000004: close GTFS-RT interoperability conformance
harness review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 121 -- Checkpoint 000004: close GTFS-RT interoperability conformance
harness review.

Goal status:
Active. Phase 121 is closed and the goal continues to Phase 122.

Sub-agents used or simulated:
Context / Repo Truth, Planning, Implementation, QA, GTFS-RT Domain, Connector,
Claim-Boundary, Security/Auth, Data/Migration, Documentation / IA, Web Design
Skill, Release, and Install Confidence roles were simulated by the Master
Agent because the agent thread limit prevented new real sub-agents.

Changed files:
`docs/handoffs/phase-121.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-121-gtfs-rt-interoperability-conformance-harness.md`.

Validation run:
Full Phase 121 validation passed before closeout docs. Focused closeout
validation passed after closeout docs: `git diff --check`, `make check`,
`make audit-product-acceptance`, `make audit-final-claim-review`,
`scripts/check-consumer-tracker.sh`, and protected-path git status.

Blocked checks:
No Phase 121 check remains blocked.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 121 remains bounded to offline/local GTFS-RT conformance diagnostics and
makes no stronger public claim.

Security/auth status:
No application security behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was added.

Release/publication status:
The public rc1 prerelease remains published. No release action was taken.

Install confidence status:
Public fresh-clone rc1 install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete.

Master review:
Approved. Phase 121 closes with an offline protobuf conformance harness and
full validation passed.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 122.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 122.

Next checkpoint:
Phase 122 -- Checkpoint 000001: add GTFS-RT fixture library and edge-case
coverage plan.
