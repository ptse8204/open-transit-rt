# Phase 122 -- GTFS-RT Fixture Library And Edge-Case Coverage

## Goal

Expand GTFS/GTFS-RT fixture and edge-case coverage for midnight rollover,
frequency service, canceled trips, stale telemetry, unknown vehicles, and
malformed realtime messages using synthetic/local data only.

Phase 122 is not a consumer acceptance, compliance, production readiness,
vendor compatibility, hardware certification, SLA/uptime, or ETA-quality proof
phase.

## Current Repo Context

- Phase 121 added the offline GTFS-RT conformance library and
  `cmd/gtfsrt-conformance`.
- Existing `testdata/gtfs`, `testdata/replay`, `testdata/telemetry`, and
  adapter-conformance fixtures cover several matching and connector edge cases.
- The new GTFS-RT conformance harness needs a manifest-backed fixture library
  so future generated protobuf cases can stay organized and claim-bounded.

## Scope

- Add or improve synthetic/local GTFS-RT fixture metadata for edge cases.
- Validate fixture metadata with tests so missing cases are caught.
- Cover the requested edge-case families without writing binary evidence,
  changing protected paths, contacting consumers, or claiming production
  outcomes.

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

- Synthetic GTFS-RT fixture library metadata or equivalent fixture suite.
- Tests validating required edge-case coverage and claim boundaries.
- `docs/handoffs/phase-122.md`
- Source-of-truth status updates for Phase 122 closeout.

## Implementation Plan

1. Add this Phase 122 plan and commit checkpoint 000001.
2. Add a synthetic GTFS-RT fixture suite manifest and test-backed manifest
   validation tied to the Phase 121 conformance harness.
3. Run relevant GTFS-RT, connector, claim-boundary, and baseline validation;
   patch repo-caused failures.
4. Close Phase 122 with handoff/status docs and continue immediately to
   Phase 123.

## Checkpoint Plan

- `Phase 122 -- Checkpoint 000001: add gtfs rt fixture library and edge case coverage plan`
- `Phase 122 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 122 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 122 -- Checkpoint 000004: close gtfs rt fixture library and edge case coverage review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 122 -- Checkpoint 000001: add GTFS-RT fixture library and edge-case
coverage plan.

Goal status:
Active. Phase 121 is closed and Phase 122 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-122-gtfs-rt-fixture-library-and-edge-case-coverage.md`.

Validation run:
Initial inspection reviewed the Phase 122 roadmap prompt, existing GTFS,
replay, telemetry, adapter-conformance fixtures, Phase 121 conformance harness,
and validation boundaries.

Blocked checks:
Implementation, tests, connector validation, and closeout validation are
scheduled for later Phase 122 checkpoints.

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
The plan is synthetic/local fixture work. It does not change route auth, CSRF
behavior, credential handling, token handling, public exposure, private payload
handling, or operator command behavior.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 122 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 122 does not plan
visual UX changes.

Master review:
Approved. The plan keeps fixture expansion synthetic/local and
claim-bounded.

Required edits:
Commit checkpoint 000001, then implement the scoped fixture library.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 122 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 122 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 122 implemented a synthetic GTFS-RT fixture suite and manifest
validator.

Sub-agents used or simulated:
The agent thread limit still prevents new real sub-agents. Context / Repo
Truth, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles were simulated by the Master Agent.

Changed files:
`testdata/gtfsrt-conformance/README.md`;
`testdata/gtfsrt-conformance/suite.json`;
`internal/gtfsrtconformance/fixtures.go`;
`internal/gtfsrtconformance/fixtures_test.go`;
`Makefile`;
`docs/phase-122-gtfs-rt-fixture-library-and-edge-case-coverage.md`.

Validation run:
Passed:

- `gofmt -w internal/gtfsrtconformance`
- `go test ./internal/gtfsrtconformance ./cmd/gtfsrt-conformance`
- `make gtfsrt-conformance`
- `python3 -m json.tool testdata/gtfsrt-conformance/suite.json`
- `git diff --check`
- `make check`
- protected-path git status check

Blocked checks:
Full Phase 122 validation, connector validation, and claim audits are
scheduled for checkpoint 000003.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
The fixture suite is synthetic/local, has all-false claim flags, and records
explicit "does not prove" language for every case. It does not claim
compliance, consumer acceptance, production readiness, SLA/uptime, vendor
compatibility, hardware certification, production-grade ETA quality, or
real-world ETA accuracy.

Security/auth status:
Fixture validation rejects evidence-like suite paths and evidence-like
references. No route auth, CSRF, credential handling, token handling, public
exposure, private payload handling, or operator command behavior changed.

Data/migration status:
No migration, schema, durable state, runtime dependency, or Go module change
was added.

Release/publication status:
The public rc1 prerelease remains published. Phase 122 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. The changed
fixture library is current-source hardening and is not part of the published
rc1 tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. CP000002 did not change
visual UI templates.

Master review:
Approved. The fixture suite covers midnight rollover, frequency service,
canceled trips, stale telemetry, unknown vehicles, and malformed realtime
messages across Vehicle Positions, Trip Updates, and Alerts.

Required edits:
Commit checkpoint 000002, then run Phase 122 validation and patch any
repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 validation.

Next checkpoint:
Phase 122 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 122 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Phase 122 fixture library validation passed.

Sub-agents used or simulated:
Context / Repo Truth, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles were simulated by the Master Agent because the
agent thread limit prevents new real sub-agents.

Changed files:
`docs/phase-122-gtfs-rt-fixture-library-and-edge-case-coverage.md`.

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
No Phase 122 validation blocker remains.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. Phase 122 remains bounded to
synthetic/local fixture coverage and makes no stable release readiness,
production readiness, compliance, adoption, agency approval, consumer
acceptance, consumer ingestion/listing/display, final-root readiness, hosted
service, paid support, SLA/uptime, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality, or
real-world ETA accuracy claim.

Security/auth status:
No route auth, CSRF, credential handling, token handling, public exposure,
private payload handling, or operator command behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was
added.

Release/publication status:
The public rc1 prerelease remains published. Phase 122 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed. Phase 122
current-source hardening is not part of the published rc1 tag.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 122 did not change
visual UI templates.

Master review:
Approved. Validation found no required patch after adding the synthetic
GTFS-RT fixture suite and required-case checks.

Required edits:
Commit checkpoint 000003, then close Phase 122 with handoff and status docs.

Decision:
Proceed to checkpoint 000003 commit and Phase 122 closeout.

Next checkpoint:
Phase 122 -- Checkpoint 000004: close GTFS-RT fixture library and edge-case
coverage review.
