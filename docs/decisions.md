# Architecture Decisions

This document records architecture-significant decisions so later phases do not re-decide core product boundaries.

## ADR-0001 — Keep the backend mostly Go

Open Transit RT should use Go services and internal packages for core backend behavior. Early admin and Studio surfaces should prefer simple server-rendered HTML unless a later phase documents a need for a heavier frontend stack.

## ADR-0002 — Use Postgres/PostGIS as source of truth

Postgres stores agency configuration, GTFS feed versions, telemetry, assignments, overrides, audit logs, validation reports, incidents, feed metadata, and compliance workflow state. PostGIS is required for future nearest-shape and spatial matching behavior.

## ADR-0003 — Use versioned migrations

Schema changes live under `db/migrations` and are applied through `cmd/migrate`. `db/schema.sql` is retained only as starter/reference context; migrations are the executable schema path.

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
