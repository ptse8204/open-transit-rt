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

Post-Outcome-C note: this file preserves the historical implementation phase
definitions. For the current post-Phase-33 continuation state, use
`docs/handoffs/latest.md` and `docs/future-roadmap-post-outcome-c.md`.

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
| 33 | Public GTFS local/pilot evidence | Attempt real public GTFS local/pilot handling proof without final-root or consumer overclaims |
| 34 | Post-Outcome-C status consistency and evidence readiness | Align post-Outcome-C docs, final-root request guidance, and public-GTFS repeatability guidance without creating external evidence |
| 35 | README and roadmap realignment | Restore the root README as the product front door and make self-hosted agency reuse the default roadmap |
| 36 | OCI/OCL reference deployment productization | Turn the existing OCI/OCL-style pilot server pattern into a repeatable reference deployment path |
| 37 | Agency reusable onboarding flow | Make agency GTFS onboarding and self-hosted first run repeatable without manual DB surgery |
| 38 | Integration adapter kit | Make AVL/device, predictor, validator, monitoring, and consumer integration boundaries reusable |
| 39 | CAL-ITP-style readiness workflow | Make readiness support visible in product workflows without claiming compliance |
| 40 | Guided self-hosted operator trial | Tie reference deployment, reusable onboarding, adapter dry-run, and readiness review into one evidence-bounded operator checklist |

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

## Phase 33 — Public GTFS Local/Pilot Evidence

### Goal
Prove, or prepare a repeatable path to prove, that Open Transit RT can ingest,
publish, validate, and review a real current public agency GTFS dataset in a
local or pilot environment.

### Status
Complete as Outcome C — public-GTFS local/pilot run completed with public-safe
retained summaries.

### Definition of done
Phase 33 must close with exactly one outcome:

- Outcome A: template-only closure.
- Outcome B: attempted public-GTFS run blocked.
- Outcome C: public-GTFS local/pilot run completed with public-safe retained
  summaries.

Only Outcome C may be described as completed public-GTFS local/pilot evidence.
Outcome B records a blocker and does not support stronger public-GTFS handling
claims. The 2026-05-06 packet reached Outcome C after the large-import fix and
retry, but it remains local/pilot evidence only.

### Boundaries
Phase 33 does not prove agency adoption, agency endorsement, agency approval,
official agency feed status, agency-owned final-root proof, consumer
submission/review/acceptance, consumer ingestion/listing/display,
Caltrans/CAL-ITP compliance, hosted SaaS availability, production readiness,
real vendor AVL compatibility, or production-grade ETA quality.

## Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness

### Goal
Make the repo status, roadmap, final-root request, public-GTFS repeatability,
and handoff docs consistent after Phase 33 Outcome C.

### Status
Complete for docs-only status consistency and evidence-readiness scope.

### Definition of done
Phase 34 is done when:

- status and roadmap docs no longer contradict Phase 33 Outcome C;
- `docs/final-root-operator-request.md` exists and is clearly not evidence;
- `docs/tutorials/public-gtfs-local-pilot.md` exists as repeatability guidance
  only;
- the Phase 33 static validator history remains clear: the original static
  validator attempt was blocked because Java was unavailable, and the later
  Homebrew Java 17 retry executed with exit code `0`, system error count `0`,
  and 3 warning notices;
- the later retry is not described as validator-clean, no-warning, compliance,
  final-root, consumer, or production evidence;
- `docs/handoffs/latest.md` identifies the self-hosted continuation path and
  keeps external proof as a future optional track.

### Boundaries
Phase 34 is docs-only. It does not add scripts, Makefile targets, runtime code,
schema changes, migrations, APIs, consumer tracker changes, final-root evidence
packets, target artifacts, OCI pilot final-root wording, or new external
evidence.

## Phase 35 — README And Roadmap Realignment

### Goal
Restore the root README as the Open Transit RT product front door and align the
default roadmap around self-hosted agency reuse instead of external-proof
chasing.

### Status
Complete for docs-only README and roadmap realignment.

### Definition of done
Phase 35 is done when:

- the root `README.md` explains the product, audience, current capabilities,
  three evaluation/deployment paths, integration boundaries, readiness pointers,
  and claim boundaries;
- roadmap/status docs identify self-hosted agency reuse and OCI/OCL reference
  deployment productization as the default next path;
- `docs/handoffs/latest.md` recommends Phase 36 — OCI/OCL Reference Deployment
  Productization next;
- `docs/backlog.md` and `docs/open-questions.md` are updated or explicitly
  reviewed in the Phase 35 handoff;
- external-proof docs remain available as future optional tracks, not the
  default roadmap;
- consumer statuses remain unchanged.

### Boundaries
Phase 35 is docs-only. It does not change runtime code, schemas, migrations,
APIs, consumer tracker statuses, final-root evidence, external evidence packets,
validator artifacts, or public feed contracts.

## Phase 36 — OCI/OCL Reference Deployment Productization

### Goal
Make the existing OCI/OCL-style pilot server pattern a repeatable self-hosted
reference deployment path.

### Required work
- add or refine reference deployment docs with placeholder-only environment
  values;
- document install, update, rollback, service supervision, reverse-proxy
  public/private routing, validator install/check, backup/restore, feed monitor,
  scorecard export, and server smoke checks;
- keep the OCI DuckDNS pilot labeled as hosted/operator pilot evidence only.

### Boundaries
Phase 36 must not claim hosted SaaS availability, agency-owned final-root proof,
consumer acceptance, CAL-ITP/Caltrans compliance, production readiness, paid
support/SLA coverage, or production multi-tenant hosting.

## Phase 37 — Agency Reusable Onboarding Flow

### Goal
Make local and server agency onboarding repeatable for a supplied agency ID,
GTFS URL, metadata, and import timeout.

### Required work
- create or refine a guided one-command agency onboarding path;
- download GTFS only to ignored storage;
- import and publish through existing GTFS paths;
- print feed URLs, admin URL, validator status or blocker, device/AVL next
  steps, and support summary;
- add tests for argument validation where practical.

### Boundaries
Phase 37 must not claim agency approval, agency adoption, official feed status,
consumer acceptance, or compliance from tooling alone.

## Phase 38 — Integration Adapter Kit

### Goal
Make integration with existing systems first-class through documented adapters.

### Required work
- expand AVL/device adapter guidance and fixtures without naming unsupported
  vendors;
- document external predictor adapter lifecycle behind
  `internal/prediction.Adapter`;
- keep validators behind pinned/allowlisted tooling;
- clarify monitoring and consumer integration boundaries;
- add conformance tests or fixtures where practical.

### Boundaries
Phase 38 must not add certified vendor compatibility, production AVL
reliability, production-grade ETA quality, or consumer acceptance claims without
retained evidence.

## Phase 39 — CAL-ITP-Style Readiness Workflow

### Goal
Make readiness support visible in product workflows while preserving evidence
boundaries.

### Required work
- expose or refine a readiness checklist for stable URLs, metadata,
  validation, freshness, GTFS-RT completeness, and consumer packet state;
- add plain-language remediation steps;
- keep wording as readiness support or technical foundations unless retained
  evidence supports stronger language.

### Boundaries
Phase 39 must not claim CAL-ITP/Caltrans compliance, consumer acceptance,
agency endorsement, production readiness, hosted SaaS availability, or
marketplace/vendor equivalence.

## Phase 40 — Guided Self-Hosted Operator Trial

### Goal
Give operators one guided local/reference trial that ties the reference
deployment docs, reusable agency onboarding, adapter kit, and readiness
workflow together.

### Required work
- add an operator trial tutorial covering local/reference prep,
  `make agency-pilot-up`, no-external-network fixture setup, five public feed
  checks, `/admin/operations/readiness`, validators, synthetic AVL dry-run,
  next actions, and teardown;
- update docs navigation, status docs, and handoffs;
- add validation checks for the Phase 40 phase doc, tutorial, and handoff.

### Boundaries
Phase 40 must not create external evidence, final-root evidence, consumer
submission APIs, real vendor payloads, consumer status changes, hosted SaaS
claims, CAL-ITP/Caltrans compliance claims, consumer acceptance claims, agency
approval/adoption claims, production-readiness claims, vendor-compatibility
claims, or production-grade ETA-quality claims.
