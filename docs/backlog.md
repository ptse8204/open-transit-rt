# Backlog

This backlog is organized by phase. `docs/phase-plan.md` remains the phase contract.

## Phase 1 — Durable Telemetry Foundation

- Complete. Durable telemetry ingest, DB readiness, duplicate/out-of-order handling, parsed payload JSONB storage, agency-scoped debug listing, and DB-backed tests are implemented.

## Phase 2 — Deterministic Trip Matching

- Complete. The matcher resolves agency-local service days, handles after-midnight and frequency cases, persists explicit unknown rows, respects manual overrides, records reasons/degraded state/incidents, and has unit plus DB-backed tests.

## Phase 3 — Vehicle Positions Production Feed

- Complete. Vehicle Positions protobuf and JSON debug output are DB-backed, generated from the same snapshot model, preserve unknown/stale behavior, and have unit plus DB-backed tests.

## Phase 4 — GTFS Import And Publish

- Complete. GTFS ZIP import, internal validation, transactional feed-version publish, failed-import report storage, rollback-safe activation, block preservation, shape-line construction, and active feed switching tests are implemented.

## Phase 5 — GTFS Studio

- Complete. Typed GTFS draft tables, draft CRUD, minimal server-rendered Studio UI, soft discard, cloned-draft provenance, draft publish traceability, and direct shared validation/activation publishing are implemented.

## Phase 6 — Trip Updates And Alerts Architecture

- Complete. Trip Updates adapter contracts, no-op adapter, diagnostics persistence, stable empty Trip Updates endpoints, stable empty Alerts endpoints, and non-coupling tests are implemented.

## Phase 7 — Prediction Quality And Operations

- Complete for the first conservative operations scope. The deterministic Trip Updates adapter, override lifecycle repository, audit logging, prediction review queue, cancellation linkage signal, deadhead/layover suppression, disruption withholds, and coverage metrics are implemented.
- Remaining later work: production-grade ETA quality, backtesting, full operator UI, canonical realtime validation, public Alerts authoring/persistence, vehicle swap UI/workflow, and richer detour/short-turn prediction behavior.

## Phase 8 — Compliance And Consumer Workflow

- Add compliance dashboard and scorecard.
- Add public metadata and license pages.
- Add consumer ingestion workflow records and export packet generation.
- Track marketplace-vendor-equivalence gaps.
