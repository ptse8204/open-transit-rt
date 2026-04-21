# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 4 — GTFS import and publish pipeline

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-03.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin Phase 4 GTFS import and publish pipeline. Do not start GTFS Studio, Trip Updates, Alerts, rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

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
- Existing GTFS repository and matcher tests seed schedule rows directly; Phase 4 must build real import/publish runtime behavior rather than promoting test-only seed helpers.

## First Files Likely To Edit

- `internal/gtfs/`
- `db/migrations/`
- `cmd/migrate/` only if migration behavior needs extension
- `testdata/gtfs/`
- `testdata/expected/`
- `docs/current-status.md`
- `docs/handoffs/phase-04.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md`
- `docs/decisions.md` if architecture-significant import/publish decisions are made

## Phase 4 Entry Recommendation

Start GTFS import and publish pipeline without changing Vehicle Positions semantics:

1. Inspect the existing published GTFS tables and schedule-query boundary.
2. Add a staging model for GTFS ZIP imports while keeping draft GTFS separate from published active feed versions.
3. Parse and validate required GTFS files into staged records.
4. Atomically activate a published feed version without changing stable public feed URLs.
5. Add rollback-safe integration tests using `testdata/gtfs/valid-small`, `after-midnight`, `frequency-based`, and `malformed`.
6. Do not implement GTFS Studio, Trip Updates, or Alerts in Phase 4.

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

## Phase 3 Notes For Phase 4

- `cmd/feed-vehicle-positions` now requires `AGENCY_ID` and DB access at startup.
- `/public/gtfsrt/vehicle_positions.pb` is DB-backed and returns valid GTFS-RT protobuf `FeedMessage` responses.
- `/public/gtfsrt/vehicle_positions.json` is diagnostic JSON generated from the same snapshot as protobuf.
- Empty telemetry or all-suppressed snapshots return normal successful empty protobuf feeds with populated headers.
- `Last-Modified` is derived from snapshot `generated_at`.
- Vehicle Positions use `internal/feed.VehiclePositionsSnapshot`; do not duplicate publication business logic in handlers.
- `telemetry.Repository.ListLatestByAgency` ordering is now a hard contract: one latest accepted row per vehicle ordered by `observed_at DESC, id DESC`.
- `state.Repository.ListCurrentAssignments` is the bulk current-assignment read boundary.
- Automatic assignments publish trip descriptors only when linked to the latest telemetry event.
- Non-exact frequency assignments map Vehicle Positions trip descriptors to `UNSCHEDULED`; exact and normal scheduled assignments use `SCHEDULED`.
- JSON debug fields are diagnostic, not a stable public API.
- Canonical GTFS-RT validator tooling remains unwired.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
