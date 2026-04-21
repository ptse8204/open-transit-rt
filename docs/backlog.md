# Backlog

This backlog is organized by phase. `docs/phase-plan.md` remains the phase contract.

## Phase 1 — Durable Telemetry Foundation

- Complete. Durable telemetry ingest, DB readiness, duplicate/out-of-order handling, parsed payload JSONB storage, agency-scoped debug listing, and DB-backed tests are implemented.

## Phase 2 — Deterministic Trip Matching

- Complete. The matcher resolves agency-local service days, handles after-midnight and frequency cases, persists explicit unknown rows, respects manual overrides, records reasons/degraded state/incidents, and has unit plus DB-backed tests.

## Phase 3 — Vehicle Positions Production Feed

- Add GTFS-RT protobuf bindings.
- Publish `/public/gtfsrt/vehicle_positions.pb` from persisted data.
- Keep JSON debug endpoint.
- Add stale/unmatched behavior and validation tests.

## Phase 4 — GTFS Import And Publish

- Import `gtfs.zip`.
- Validate required files and references.
- Stage imported data and atomically activate feed versions.
- Store import and validation reports.

## Phase 5 — GTFS Studio

- Add draft GTFS CRUD model.
- Add minimal server-rendered admin UI.
- Publish drafts through the same pipeline as ZIP import.

## Phase 6 — Trip Updates And Alerts Architecture

- Define `PredictionAdapter` input/output contracts.
- Add no-op adapter and diagnostics plumbing.
- Add stable endpoint shape for Trip Updates and Alerts.

## Phase 7 — Prediction Quality And Operations

- Add stop-level ETA model and quality metrics.
- Add incident queue workflows.
- Add cancellation, added trip, detour, short turn, and vehicle swap operations.

## Phase 8 — Compliance And Consumer Workflow

- Add compliance dashboard and scorecard.
- Add public metadata and license pages.
- Add consumer ingestion workflow records and export packet generation.
- Track marketplace-vendor-equivalence gaps.
