# Architecture Decisions

This document records architecture-significant decisions so later phases do not re-decide core product boundaries.

## ADR-0001 — Keep the backend mostly Go

Open Transit RT should use Go services and internal packages for core backend behavior. Early admin and Studio surfaces should prefer simple server-rendered HTML unless a later phase documents a need for a heavier frontend stack.

## ADR-0002 — Use Postgres/PostGIS as source of truth

Postgres stores agency configuration, GTFS feed versions, telemetry, assignments, overrides, audit logs, validation reports, incidents, feed metadata, and compliance workflow state. PostGIS is required for future nearest-shape and spatial matching behavior.

## ADR-0003 — Use versioned migrations

Schema changes live under `db/migrations` and are applied through `cmd/migrate`. Migrations are the source of truth for the executable database schema.

`db/schema.sql` is deprecated as an executable schema file. It remains only as a short comment-only pointer for readers or tools that still expect the path to exist. It must not contain `CREATE`, `ALTER`, or `DROP` statements, must not be edited as an independent schema definition, and must not be used to apply database changes. If a future phase wants a full schema snapshot, it should generate it from migrations and document that workflow before replacing the pointer file.

## ADR-0004 — Keep Trip Updates pluggable

Trip Updates must stay behind a prediction adapter boundary. Open Transit RT owns GTFS management, telemetry, assignments, audit logs, and Vehicle Positions. Optional predictors such as TheTransitClock may generate ETAs or Trip Updates only behind an adapter.

## ADR-0005 — Publish Vehicle Positions first

Vehicle Positions are the first production-grade public realtime output. Trip Updates and Alerts are architecture-binding but implemented in later phases.

## ADR-0006 — Prefer unknown over false certainty

The matcher must be conservative. Low-confidence or contradictory evidence should produce `unknown` plus incidents/diagnostics instead of a speculative trip descriptor.

## ADR-0007 — Manual overrides take precedence

Operator overrides are part of the core model. Active overrides must beat automatic matching until they expire or are cleared, and privileged actions must be audit logged.

## ADR-0008 — Keep draft and published GTFS separate

GTFS Studio draft data and active published feed versions must not collapse into one model. Import and Studio are two sources that publish through a shared validated feed-version pipeline.

## ADR-0009 — Stable public URLs are product contracts

Public schedule, Vehicle Positions, Trip Updates, and Alerts URLs must stay stable across feed updates and rollback. Version changes happen behind those URLs.

## ADR-0010 — Phase 0 is foundation-only

Phase 0 may design schemas, contracts, and docs for later requirements, but it must not implement later-phase runtime behavior such as durable telemetry, trip matching, GTFS import, protobuf feed generation, Trip Updates, or Alerts.

## ADR-0011 — Persist telemetry through an agency-scoped repository

Telemetry ingest writes must go through a repository backed by Postgres/PostGIS. The repository classifies accepted, duplicate, and out-of-order telemetry inside one transaction protected by a deterministic advisory lock derived from agency and vehicle identity. The lock key is a SHA-256-derived signed 64-bit value; theoretical collisions only serialize unrelated streams and do not merge data because SQL predicates and uniqueness remain authoritative. Canonical accepted telemetry uniqueness is vehicle-scoped by `(agency_id, vehicle_id, observed_at)`; `device_id` is retained for audit/debug but does not define the canonical latest vehicle position.

Invalid JSON and invalid telemetry payloads are rejected before repository storage in Phase 1. The database `rejected` status remains reserved for a later ingest-audit phase that explicitly designs rejected-payload retention.

## ADR-0012 — Persist explicit deterministic assignment outcomes

Phase 2 persists every matcher outcome as a `vehicle_trip_assignment` row, including `unknown`. Unknown results close any previous active row so stale or low-confidence telemetry cannot leave a prior confident trip active. Unknown rows carry `service_date` whenever agency timezone and observed timestamp can be resolved; the column remains nullable only for unresolved cases.

Assignment reasons and degraded state use a small stable taxonomy. `score_details_json` is intentionally loose debug JSON for Phase 2 and is not a stable public API or integration schema. The internal convention is that matcher-generated score details include `score_schema`; candidate-based score details also include `trip_id`, `start_time`, and `observed_local_seconds` when resolvable. Future public or adapter-facing diagnostics should define a separate versioned contract rather than depending on this debug payload.

Phase 2 treats `missing_shape` as both a reason code and a dedicated degraded-state category. Missing shapes reduce confidence but do not automatically prevent a match when other evidence is strong. Route-hint matching is reserved for a future telemetry/input expansion and is not part of the active Phase 2 reason-code taxonomy.
