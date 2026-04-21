# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 7 — Prediction quality and operations workflows

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 GTFS Studio draft/publish model is implemented and complete.
- Phase 6 Trip Updates and Alerts architecture is implemented and complete.
- Phase 7 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-06.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin Phase 7 prediction quality and operations workflows. Do not bypass the Phase 6 prediction adapter boundary, do not couple predictor internals into telemetry ingest, Vehicle Positions, or GTFS Studio, and do not start rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Exact First Commands

```bash
command -v go
go version
make fmt
make test
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
make test-integration
```

If Task is installed, optional equivalents may be run:

```bash
task fmt
task test
task migrate:status
task test:integration
```

## Known Blockers

- Task is optional and may not be installed; Makefile remains independently usable.
- Docker must be running before DB-backed checks.
- Canonical GTFS and GTFS-RT validators are documented but not wired yet.
- Trip Updates currently use a no-op adapter; no ETA-quality predictor is implemented yet.
- Alerts endpoints exist but alert authoring and persistence are intentionally deferred.
- GTFS Studio auth is minimal/dev-only and not production-grade.

## First Files Likely To Edit

- `internal/prediction/`
- `internal/feed/tripupdates/`
- `internal/feed/alerts/`
- `internal/state/` only for operation-state inputs needed by prediction
- `internal/gtfs/` only if additional schedule-query contracts are needed
- `db/migrations/` only if Phase 7 adds alert/prediction persistence beyond existing health snapshots
- `docs/current-status.md`
- `docs/handoffs/phase-07.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md`
- `docs/decisions.md` if architecture-significant prediction/alerts decisions are made

## Phase 7 Entry Recommendation

Start prediction quality and operations work behind the Phase 6 contracts:

1. Inspect `internal/prediction`, `internal/feed/tripupdates`, `internal/feed/alerts`, `internal/state`, and `internal/gtfs`.
2. Choose the first real prediction strategy: internal deterministic ETA logic or an external adapter such as TheTransitClock.
3. Keep public Trip Updates endpoint shape stable while replacing the no-op adapter behavior.
4. Add stop-level predictions only from active published GTFS, persisted latest telemetry, and persisted assignments.
5. Preserve deterministic entity ordering and ordered `stop_time_update` entries.
6. Preserve `FeedHeader.timestamp` and `Last-Modified` alignment from snapshot `GeneratedAt`.
7. Add alert authoring/persistence only after deciding whether alerts are operator-authored, incident-derived, or both.
8. Keep Vehicle Positions, telemetry ingest, GTFS import, and GTFS Studio behavior stable.

## Constraints To Preserve

- Mostly Go.
- Postgres/PostGIS source of truth.
- Vehicle Positions first.
- Trip Updates pluggable.
- Draft GTFS separate from published GTFS.
- Conservative matching.
- Manual overrides take precedence over matching.
- No rider apps, payments, passenger accounts, or dispatcher CAD.
- External integrations stay behind documented adapters.
- Runtime GTFS import input is ZIP; directory parsing is test-fixture setup only.
- GTFS Studio publishes typed draft rows through the shared validation/activation helper directly, not through synthetic ZIP import.
- GTFS times beyond `24:00:00` remain stored as imported text in canonical published GTFS tables.
- Phase 6 Trip Updates packages must not become dependencies of telemetry ingest, Vehicle Positions, or GTFS Studio.

## Phase 6 Notes For Phase 7

- `internal/prediction.Adapter` is the only Trip Updates prediction boundary.
- `prediction.NoopAdapter` is the default adapter and returns explicit no-op diagnostics.
- Trip Updates diagnostics persist to `feed_health_snapshot` with required fields in `details_json`.
- `cmd/feed-trip-updates` exposes `/public/gtfsrt/trip_updates.pb` and `/public/gtfsrt/trip_updates.json`.
- `cmd/feed-alerts` exposes `/public/gtfsrt/alerts.pb` and `/public/gtfsrt/alerts.json`.
- `VEHICLE_POSITIONS_FEED_URL` is an exact full Phase 3 protobuf URL. If unset, `FEED_BASE_URL` must include `/public` and derives `/public/gtfsrt/vehicle_positions.pb`.
- Trip Updates and Alerts responses derive `Last-Modified` from the same snapshot `GeneratedAt` used for `FeedHeader.timestamp`.
- Alerts do not write `feed_health_snapshot` rows in Phase 6; deferred status is JSON-only.
- No database migration was added in Phase 6.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
