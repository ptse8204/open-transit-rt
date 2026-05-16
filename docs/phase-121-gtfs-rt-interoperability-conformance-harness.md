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
