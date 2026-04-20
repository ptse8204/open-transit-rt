# Phase Plan

This document defines the phased implementation plan for Open Transit RT.

This file is the source of truth for:
- phase order
- phase goals
- dependencies between phases
- definition of done
- expected handoff behavior between Codex instances

Related binding docs:
- `AGENTS.md`
- `docs/codex-task.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/dependencies.md`

## Implementation policy

- `docs/codex-task.md` defines implementation order.
- The requirements docs define the full product contract.
- Not every requirement must be fully implemented in the first phases, but every phase must preserve the ability to satisfy all binding requirements later without major rewrites.
- Each phase must end with:
  - updated status docs
  - updated handoff docs
  - tests added or updated
  - commands/checks run and reported
  - explicit known limitations

## Global architecture constraints

These constraints apply to every phase:

- keep the backend codebase mostly Go
- use Postgres/PostGIS as the source of truth
- keep GTFS Studio draft data separate from published GTFS feed versions
- keep Vehicle Positions as the first production-grade public output
- keep Trip Updates pluggable behind an adapter boundary
- prefer `unknown` over false certainty in trip matching
- manual overrides must take precedence over automatic matching
- do not add rider apps, passenger accounts, payments, or dispatcher CAD
- do not let external predictor internals leak into core domain packages
- do not treat external tools as the source of truth unless explicitly documented

## Phase overview

| Phase | Name | Goal |
|---|---|---|
| 0 | Scaffolding and repo hardening | Make the repo runnable, testable, and safe for phased development |
| 1 | Durable telemetry foundation | Persist telemetry and basic operational state in Postgres/PostGIS |
| 2 | Deterministic trip matching | Build conservative assignment logic with required edge cases |
| 3 | Vehicle Positions production feed | Publish valid GTFS-RT Vehicle Positions from real persisted data |
| 4 | GTFS import and publish pipeline | Import, validate, stage, and atomically publish GTFS feeds |
| 5 | GTFS Studio draft/publish model | Add draft GTFS editing and publish through the same pipeline |
| 6 | Trip Updates and Alerts architecture | Define pluggable prediction/alerts boundaries and minimal adapters |
| 7 | Prediction quality and operations workflows | Improve ETA quality, overrides, incidents, and realtime coverage |
| 8 | Compliance and consumer workflow | Add compliance scorecard, validation, discoverability, and ingestion workflows |

---

## Phase 0 — Scaffolding and repo hardening

### Goal
Add the missing repo scaffolding so future phases can execute with low ambiguity and strong reproducibility.

### Required work
- add `.env.example`
- add `Taskfile.yml` and/or expand `Makefile`
- add `cmd/migrate`
- add versioned migrations under `db/migrations`
- add `scripts/bootstrap-dev.sh`
- add `testdata/` fixtures
- add `docs/decisions.md`
- add `docs/dependencies.md`
- add and align status/handoff docs
- ensure `AGENTS.md` reflects the current repo contract

### Dependencies
- none

### Definition of done
Phase 0 is done when:
- repo scaffolding exists and is documented
- local bootstrap flow exists
- migration flow exists
- fixture directory exists
- status and handoff docs exist under `docs/handoffs/`
- all docs reference `docs/handoffs/latest.md` as the handoff source of truth
- `docs/handoffs/template.md` is the default required structure for future handoffs unless a phase documents a reason to diverge
- the foundation schema creates only tables, constraints, and indexes for later phases without implementing later runtime behavior
- Makefile workflows remain independently usable even if Task is not installed
- Phase 0 handoff includes exact Phase 1 entry files, commands, blockers, and first implementation slice
- baseline checks can be run or blocked reasons are explicitly documented

---

## Phase 1 — Durable telemetry foundation

### Goal
Replace in-memory telemetry with durable persistence and create the core DB/repository foundation.

### Required work
- add shared DB package using `pgx` / `pgxpool`
- persist telemetry events in Postgres
- add repository interfaces for telemetry, assignments, feed lookup, and agency-scoped access
- add durable health/readiness behavior
- capture raw payload JSON
- handle duplicate and out-of-order telemetry safely
- add DB-backed tests and fixtures

### Dependencies
- Phase 0 complete

### Definition of done
Phase 1 is done when:
- telemetry is no longer stored only in process memory
- DB readiness is checked in health/readiness paths
- repository interfaces are in place
- tests cover telemetry insert/query and basic edge cases
- no placeholder persistence path remains in production code

---

## Phase 2 — Deterministic trip matching

### Goal
Implement conservative rule-based assignment logic with the required operational edge cases.

### Required work
- agency-local service-day resolution
- after-midnight trip handling
- repeated trip-instance handling
- `frequencies.txt` support
- shape proximity and projected progress
- stop-sequence progress
- schedule fit
- continuity from previous assignment
- block transitions
- stale telemetry behavior
- low-confidence `unknown`
- manual override precedence in the data model and logic
- incident generation for degraded or ambiguous cases

### Dependencies
- Phase 1 complete
- GTFS schedule query model available enough for matching

### Definition of done
Phase 2 is done when:
- matcher assigns trips conservatively from real persisted data
- required edge cases are covered by tests
- low-confidence cases resolve to `unknown`
- assignment reasons/confidence are persisted
- incidents or degraded-state markers exist for bad/ambiguous cases

---

## Phase 3 — Vehicle Positions production feed

### Goal
Publish a valid GTFS-RT Vehicle Positions feed from real persisted data.

### Required work
- protobuf-based GTFS-RT serialization
- stable public Vehicle Positions endpoint
- JSON debug endpoint
- feed freshness behavior and stale vehicle handling
- stable entity IDs
- `FeedHeader.timestamp`
- `Last-Modified` and normal HTTP behavior where implemented
- validation tests for feed correctness

### Dependencies
- Phase 1 complete
- Phase 2 sufficiently complete for assignment-aware positions

### Definition of done
Phase 3 is done when:
- `/public/gtfsrt/vehicle_positions.pb` is served from real data
- protobuf output is valid
- stale/unmatched behavior is deterministic and tested
- placeholder sample feed output is removed from production paths

---

## Phase 4 — GTFS import and publish pipeline

### Goal
Support GTFS ZIP ingestion, validation, staging, and atomic publish.

### Required work
- accept `gtfs.zip`
- validate required files and references
- parse times beyond `24:00:00`
- support calendars, calendar dates, shapes, frequencies, and blocks
- stage data before activation
- atomically activate a published feed version
- rollback-safe publish behavior
- import reports with warnings/errors

### Dependencies
- Phase 0 complete
- DB schema and migration system stable

### Definition of done
Phase 4 is done when:
- a GTFS ZIP can be imported and published atomically
- validation reports are stored
- failed imports do not partially activate
- active feed switching is tested

---

## Phase 5 — GTFS Studio draft/publish model

### Goal
Support interactive GTFS editing without collapsing draft and published models.

### Required work
- draft GTFS schema
- minimal admin UI or server-rendered pages
- CRUD for core draft GTFS entities
- publish from draft through the same validation/activation pipeline as ZIP import
- version visibility and publish traceability

### Dependencies
- Phase 4 complete or publish pipeline stable enough to reuse

### Definition of done
Phase 5 is done when:
- draft GTFS can be edited separately from published GTFS
- publish from draft uses the same publish pipeline as ZIP import
- separation between draft and published data is enforced in schema and code

---

## Phase 6 — Trip Updates and Alerts architecture

### Goal
Define the pluggable architecture for Trip Updates and Alerts without overcommitting to a single predictor.

### Required work
- define `PredictionAdapter`
- define input/output contracts
- add a documented no-op adapter or minimal adapter
- add Alerts feed model and stable endpoint shape
- add diagnostics plumbing for prediction status
- document external predictor integration rules

### Dependencies
- Phase 3 complete
- Phase 4 sufficiently complete for published GTFS inputs

### Definition of done
Phase 6 is done when:
- Trip Updates are architecturally pluggable
- Alerts architecture exists
- no external predictor internals leak into core packages
- failure behavior for unavailable predictors is documented and tested at the boundary

---

## Phase 7 — Prediction quality and operations workflows

### Goal
Improve ETA quality and add operational controls.

### Required work
- stop-level ETA logic
- override workflows
- incident queue
- cancellation / added trip / short turn / detour handling
- prediction diagnostics and coverage metrics
- operator-facing workflows for fixing bad assignments or bad predictions

### Dependencies
- Phase 6 complete

### Definition of done
Phase 7 is done when:
- Trip Updates quality is measurable
- operations staff can override or repair bad realtime state
- prediction coverage and degraded cases are visible

---

## Phase 8 — Compliance and consumer workflow

### Goal
Implement the non-core but required compliance and publication workflows.

### Required work
- validation dashboarding
- open-license metadata
- public feed metadata pages
- stable public URLs for all feed types
- consumer ingestion workflow records
- compliance scorecard
- marketplace-gap tracking

### Dependencies
- Phases 3, 4, 6, and 7 sufficiently complete

### Definition of done
Phase 8 is done when:
- technical compliance status is measurable per agency
- discoverability metadata exists
- feed validation posture is visible
- consumer submission workflow is tracked
- gaps between technical compliance and vendor-equivalent packaging are explicit

---

## Handoff rule

At the end of every phase:
- update `docs/current-status.md`
- update `docs/backlog.md`
- update `docs/open-questions.md`
- update the phase-specific handoff file
- update `docs/handoffs/latest.md`

All phase handoff files should use `docs/handoffs/template.md` unless the phase handoff explicitly documents why it diverges.

A fresh Codex instance should be able to continue by reading:
- `AGENTS.md`
- `docs/phase-plan.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`

## Architecture-safe escalation rule

If a phase reveals a conflict with a binding requirement:
- do not silently ignore it
- document it
- update the relevant docs
- choose the smallest implementation that preserves the long-term contract
