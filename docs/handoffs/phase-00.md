# Phase 0 Handoff

## Phase

Phase 0 — Scaffolding and repo hardening

## Status

- Complete
- Active phase after this handoff: Phase 1 — Durable telemetry foundation

## What Was Implemented

- Added `.env.example` with local defaults for DB URLs, service ports, public/feed base URLs, auth placeholders, metrics flags, matcher thresholds, stale telemetry TTLs, and validation settings.
- Added `Taskfile.yml` with optional Task workflows.
- Expanded `Makefile` so it remains independently usable without Task.
- Added `cmd/migrate` as a Goose-backed migration command.
- Added `db/migrations/000001_initial_schema.sql` with a Postgres/PostGIS foundation schema.
- Changed Docker Compose to use `postgis/postgis:16-3.4` with a healthcheck and named volume.
- Added `scripts/bootstrap-dev.sh`.
- Added deterministic fixture structure under `testdata/`.
- Added `docs/decisions.md`, `docs/backlog.md`, and `docs/open-questions.md`.
- Added plural handoff source of truth under `docs/handoffs/`.
- Retired `docs/handoff/latest.md` as a source of truth; it now points to `docs/handoffs/latest.md`.

## What Was Designed But Intentionally Not Implemented Yet

- Durable telemetry runtime persistence.
- Repository interfaces and DB connection package.
- Deterministic trip matching.
- GTFS import and validation runtime behavior.
- GTFS Studio UI and draft editing runtime behavior.
- GTFS-RT Vehicle Positions protobuf generation.
- Trip Updates predictor adapters.
- Alerts generation.
- Compliance dashboard and consumer ingestion workflows.

Phase 0 intentionally created schemas, contracts, docs, and fixtures so those later requirements can be implemented without major rewrites.

## Schema And Interface Changes

- `db/schema.sql` is now a legacy pointer. The executable schema is under `db/migrations`.
- Added foundation tables for:
  - agency, users, role bindings, device credentials, feed config
  - feed versions and published feed metadata
  - canonical GTFS route/stop/calendar/trip/stop time/shape/frequency data
  - GTFS draft records separate from published feed versions
  - telemetry events
  - manual overrides
  - vehicle trip assignments
  - incidents
  - validation reports
  - feed health snapshots
  - consumer ingestion records
  - marketplace gap tracking
  - audit logs
- Foundation schema is tables, constraints, and indexes only. It does not implement later-phase runtime behavior.
- Added migration command interface:
  - `go run ./cmd/migrate up`
  - `go run ./cmd/migrate down`
  - `go run ./cmd/migrate status`
  - `go run ./cmd/migrate redo`

## Dependency Changes

- Added Go module requirements:
  - `github.com/jackc/pgx/v5`
  - `github.com/pressly/goose/v3`
- Documented `pgx`, Goose, Task, PostGIS, future protobuf/validator tooling, and external adapter boundaries in `docs/dependencies.md`.
- Task is optional. Makefile is the required fallback workflow.

## Migrations Added

- `db/migrations/000001_initial_schema.sql`

The migration creates PostGIS and foundation-only tables/constraints for later phases.

## Tests Added And Results

- Added fixture files under:
  - `testdata/gtfs/valid-small/`
  - `testdata/gtfs/after-midnight/`
  - `testdata/gtfs/frequency-based/`
  - `testdata/gtfs/malformed/`
  - `testdata/telemetry/`
  - `testdata/expected/`
- No Go test files were added in Phase 0 because runtime behavior is intentionally not implemented yet.

## Checks Run And Blocked Checks

- `command -v go`: failed; `go` is not on `PATH`.
- `command -v gofmt`: failed; `gofmt` is not on `PATH`.
- `command -v task`: failed; Task is not on `PATH`.
- `docker compose -f deploy/docker-compose.yml version`: passed; Docker Compose v2.40.3 is available.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make fmt`: blocked because `gofmt` is missing.
- `make test`: blocked because `go` is missing.
- `make test-integration`: blocked because `go` is missing.
- `make migrate-status`: blocked because `go` is missing.
- `make migrate-up`: blocked because `go` is missing.
- `make validate`: passed as a Phase 0 placeholder; validators are documented but not wired yet.
- `make lint`: passed as a no-op fallback; `golangci-lint` is not installed.
- `./scripts/bootstrap-dev.sh`: blocked because `go` is missing.
- `git diff --check`: passed.

## Known Issues

- `go.sum` was not generated because the Go toolchain is unavailable.
- Phase 0 migration syntax has not been executed against Postgres because `cmd/migrate` could not run without Go.
- Runtime services still use starter in-memory/sample behavior.
- `docs/handoff/latest.md` still exists only as a retired-path pointer for compatibility; it must not be treated as source of truth.

## Exact Next-Step Recommendation

### First files to read

1. `AGENTS.md`
2. `docs/phase-plan.md`
3. `docs/current-status.md`
4. `docs/handoffs/latest.md`
5. `docs/dependencies.md`
6. `docs/decisions.md`
7. `docs/repo-gaps.md`

### First files likely to edit

1. `go.mod` and generated `go.sum`
2. `internal/db/` new package for `pgxpool` setup
3. `internal/telemetry/` repository interfaces and Postgres implementation
4. `cmd/telemetry-ingest/main.go`
5. `cmd/migrate/main.go` only if Phase 1 finds migration command issues after Go is available
6. `docs/current-status.md`
7. `docs/handoffs/phase-01.md`
8. `docs/handoffs/latest.md`

### Commands to run before coding

```bash
command -v go
go version
go mod tidy
make fmt
make test
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-up
make migrate-status
```

If Task is installed, the equivalent Task commands may also be run, but Makefile commands must remain supported.

### Known blockers

- Go is currently missing from `PATH`.
- gofmt is currently missing from `PATH`.
- Task is currently missing from `PATH`, but Task is optional.
- Migration execution and Go tests are blocked until Go is available.

### Recommended first implementation slice

Start Phase 1 by making telemetry durable without changing feed or matcher behavior:

1. Add `internal/db` with config loading and `pgxpool` connection setup.
2. Add telemetry repository interfaces in `internal/telemetry`.
3. Add a Postgres telemetry repository that inserts `telemetry_event` rows and queries latest events by agency/vehicle.
4. Update `cmd/telemetry-ingest` to use the repository instead of global process memory.
5. Add readiness behavior that reports DB connectivity.
6. Add DB-backed tests for valid insert/query, duplicate telemetry, and out-of-order telemetry using `testdata/telemetry`.
7. Update Phase 1 handoff docs with checks and any schema adjustments.
