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
