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
- Remaining later work: production-grade ETA quality, backtesting, full operator UI, vehicle swap UI/workflow, and richer detour/short-turn prediction behavior.

## Phase 8 — Compliance And Consumer Workflow

- Complete for the first publication/compliance layer. Persisted Alerts, public Alerts feeds, canceled-trip alert reconciliation, public schedule ZIP publication, public feed metadata, license/contact workflows, consumer ingestion records, marketplace-gap records, validator report execution/recording, and compliance scorecard snapshots are implemented.
- Remaining later work: richer operator UI, production observability/SLO rollups, deployment-specific validator evidence, and external consumer acceptance evidence.

## Phase 9 — Production Closure

- Complete for the current codebase surface. Validator execution hardening, pinned validator install/check workflow, admin JWT/cookie auth, DB-backed roles, device token binding/rebinding, assignment current-row race protection, request IDs/logging, metrics gating, readiness checks, and smoke coverage are implemented.
- Remaining later work: hosted login/SSO, server-side admin JWT `jti` replay tracking, production SLO dashboards, and deployment-specific monitoring/alerting assets.

## Phase 10 — Docs, Tutorials, Deployment, And Demo

- Complete for the current codebase surface. README, local/deployment/demo/checklist tutorials, executable agency demo flow, docs assets, bootstrap output polish, and truthful CAL-ITP/Caltrans-aligned wording are implemented.
- Remaining later work: deployment-specific proof for any stronger readiness or consumer-ingestion claims.

## Phase 11 — Compliance Evidence And Optional External Integrations

- Complete for the selected evidence-only path. The Phase 11 evidence checklist, dependency reality table, README/tutorial truthfulness links, current-status update, and Phase 11 handoff are implemented.
- Remaining later work: deployment evidence hardening, real HTTPS feed proof, production validator records, monitored operations evidence, scorecard export, and third-party submission or acceptance records where they exist.
