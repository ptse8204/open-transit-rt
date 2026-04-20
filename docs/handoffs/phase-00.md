# Phase 0 Handoff

## Phase

Phase 0 — Scaffolding and repo hardening

## Status

- Complete.
- Operational closure audit passed.
- Active phase after this handoff: Phase 1 — Durable telemetry foundation.

## What Was Implemented

- Added `.env.example` with local defaults for DB URLs, service ports, public/feed base URLs, auth placeholders, metrics flags, matcher thresholds, stale telemetry TTLs, and validation settings.
- Added `Taskfile.yml` with optional Task workflows.
- Expanded `Makefile` so it remains independently usable without Task.
- Added `cmd/migrate` as a Goose-backed migration command.
- Added `db/migrations/000001_initial_schema.sql` with a Postgres/PostGIS foundation schema.
- Changed Docker Compose to use `postgis/postgis:16-3.4` with a healthcheck and named volume.
- Mapped local PostGIS to host port `55432` because `5432` was already occupied on this machine.
- Added `scripts/bootstrap-dev.sh`.
- Added deterministic fixture structure under `testdata/`.
- Added `docs/decisions.md`, `docs/backlog.md`, and `docs/open-questions.md`.
- Added plural handoff source of truth under `docs/handoffs/`.
- Removed the retired singular handoff path so `docs/handoffs/latest.md` is the only handoff source of truth.
- Installed Go through Homebrew for this environment.

## Exact Files Changed

- `.env.example`
- `Taskfile.yml`
- `Makefile`
- `README.md`
- `go.mod`
- `go.sum`
- `cmd/migrate/main.go`
- `db/schema.sql`
- `db/migrations/000001_initial_schema.sql`
- `deploy/docker-compose.yml`
- `scripts/bootstrap-dev.sh`
- `testdata/README.md`
- `testdata/expected/README.md`
- `testdata/gtfs/valid-small/agency.txt`
- `testdata/gtfs/valid-small/calendar.txt`
- `testdata/gtfs/valid-small/routes.txt`
- `testdata/gtfs/valid-small/shapes.txt`
- `testdata/gtfs/valid-small/stop_times.txt`
- `testdata/gtfs/valid-small/stops.txt`
- `testdata/gtfs/valid-small/trips.txt`
- `testdata/gtfs/after-midnight/agency.txt`
- `testdata/gtfs/after-midnight/calendar.txt`
- `testdata/gtfs/after-midnight/routes.txt`
- `testdata/gtfs/after-midnight/shapes.txt`
- `testdata/gtfs/after-midnight/stop_times.txt`
- `testdata/gtfs/after-midnight/stops.txt`
- `testdata/gtfs/after-midnight/trips.txt`
- `testdata/gtfs/frequency-based/agency.txt`
- `testdata/gtfs/frequency-based/calendar.txt`
- `testdata/gtfs/frequency-based/frequencies.txt`
- `testdata/gtfs/frequency-based/routes.txt`
- `testdata/gtfs/frequency-based/shapes.txt`
- `testdata/gtfs/frequency-based/stop_times.txt`
- `testdata/gtfs/frequency-based/stops.txt`
- `testdata/gtfs/frequency-based/trips.txt`
- `testdata/gtfs/malformed/README.md`
- `testdata/gtfs/malformed/agency.txt`
- `testdata/gtfs/malformed/routes.txt`
- `testdata/gtfs/malformed/stop_times.txt`
- `testdata/gtfs/malformed/stops.txt`
- `testdata/telemetry/after-midnight.json`
- `testdata/telemetry/frequency-based.json`
- `testdata/telemetry/matched-vehicle.json`
- `testdata/telemetry/stale-vehicle.json`
- `testdata/telemetry/swapped-vehicle.json`
- `testdata/telemetry/unmatched-vehicle.json`
- `docs/backlog.md`
- `docs/current-status.md`
- `docs/decisions.md`
- `docs/dependencies.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-00.md`
- `docs/handoffs/template.md`
- `docs/open-questions.md`
- `docs/phase-plan.md`

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

- Migrations under `db/migrations` are the schema source of truth.
- `db/schema.sql` is deprecated as an executable schema and remains only as a comment-only compatibility pointer.
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

- Installed Go with Homebrew:
  - `go version go1.26.2 darwin/amd64`
- Added Go module requirements:
  - `github.com/jackc/pgx/v5`
  - `github.com/pressly/goose/v3`
- Added `go.sum` from `go mod tidy`.
- Documented `pgx`, Goose, Task, PostGIS, future protobuf/validator tooling, and external adapter boundaries in `docs/dependencies.md`.
- Task is optional. Makefile is the required fallback workflow.

## Migrations Added

- `db/migrations/000001_initial_schema.sql`

The migration creates PostGIS and foundation-only tables/constraints for later phases.

Migration execution was verified with:
- `make migrate-up`
- `make migrate-status`
- `scripts/bootstrap-dev.sh`

## Tests Added And Results

- Added fixture files under:
  - `testdata/gtfs/valid-small/`
  - `testdata/gtfs/after-midnight/`
  - `testdata/gtfs/frequency-based/`
  - `testdata/gtfs/malformed/`
  - `testdata/telemetry/`
  - `testdata/expected/`
- No Go test files were added in Phase 0 because runtime behavior is intentionally not implemented yet.
- `make test`: passed.
- `make test-integration`: passed; this is currently a Phase 0 integration smoke path that verifies database reachability, migration visibility, and package compilation. There are no DB-backed integration test files yet.

## Checks Run And Blocked Checks

Required closure commands:

| Command | Result | Notes |
|---|---|---|
| `command -v go` | Passed | `/usr/local/bin/go` |
| `command -v gofmt` | Passed | `/usr/local/bin/gofmt` |
| `go version` | Passed | `go version go1.26.2 darwin/amd64` |
| `go mod tidy` | Passed | Generated `go.sum`; resolved pgx to v5.7.4. |
| `make fmt` | Passed | Ran `gofmt -w ./cmd ./internal`. |
| `make test` | Passed | All packages compile; no test files yet. |
| `make db-up` | Passed | Docker daemon is running; PostGIS container starts on host port `55432`. |
| `make migrate-up` | Passed | Applied `000001_initial_schema.sql`. |
| `make migrate-status` | Passed | Reports migration version 1 applied. |
| `make test-integration` | Passed | Verifies database reachability, migration visibility, and package compilation; no DB-backed integration test files yet. |
| `scripts/bootstrap-dev.sh` | Passed | Starts DB, confirms readiness, and reports no pending migrations. |
| Task equivalents | Not run | `task` is not installed; optional because Makefile is independently usable. |

Additional checks:

| Command | Result | Notes |
|---|---|---|
| `docker compose -f deploy/docker-compose.yml version` | Passed | Docker Compose v2.40.3 is available. |
| `docker compose -f deploy/docker-compose.yml config` | Passed | Compose file renders successfully. |
| `make validate` | Passed scaffold validation | Checks required migration and fixture scaffolding only; canonical GTFS and GTFS-RT validators are documented but not wired. |
| `make lint` | Passed optional fallback | `golangci-lint` not installed; future CI should make lint required once configured. |
| `git diff --check` | Passed | No whitespace errors. |
| handoff path audit | Passed | No repo docs reference the retired singular handoff path. |

## Makefile Independent Usability Audit

The Makefile does not require Task and has direct targets for:
- `fmt`
- `test`
- `test-integration`
- `migrate-status`
- `migrate-up`
- `dev` / `bootstrap`
- `db-up`
- `db-down`
- `validate`

Operational result:
- `fmt`, `test`, `test-integration`, `migrate-status`, `migrate-up`, and `bootstrap` all run through Makefile or direct script entrypoints without Task.
- Task remains optional and unavailable in this environment.

## Known Issues

- Runtime services still use starter in-memory/sample behavior.
- `golangci-lint` is not installed; `make lint` clearly reports an optional fallback.
- Static GTFS and GTFS-RT validators are documented but not wired in Phase 0.
- Docker must be running before `make db-up`, migrations, or bootstrap.

## Exact Next-Step Recommendation

### First files to read

1. `AGENTS.md`
2. `docs/phase-plan.md`
3. `docs/current-status.md`
4. `docs/handoffs/latest.md`
5. `docs/dependencies.md`
6. `docs/decisions.md`
7. `docs/repo-gaps.md`

### First files likely to edit for Phase 1

1. `internal/db/`
2. `internal/telemetry/`
3. `cmd/telemetry-ingest/main.go`
4. `docs/current-status.md`
5. `docs/handoffs/phase-01.md`
6. `docs/handoffs/latest.md`

### Commands to run before Phase 1 coding

```bash
command -v go
go version
make fmt
make test
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
```

If Task is installed, optional equivalents may also be run, but Makefile commands must remain supported.

### Known blockers

- No active Phase 0 blockers remain.
- Task is not installed, but Task is optional.

### Recommended first implementation slice

Start Phase 1 by making telemetry durable without changing feed or matcher behavior:

1. Add `internal/db` with config loading and `pgxpool` connection setup.
2. Add telemetry repository interfaces in `internal/telemetry`.
3. Add a Postgres telemetry repository that inserts `telemetry_event` rows and queries latest events by agency/vehicle.
4. Update `cmd/telemetry-ingest` to use the repository instead of global process memory.
5. Add readiness behavior that reports DB connectivity.
6. Add DB-backed tests for valid insert/query, duplicate telemetry, and out-of-order telemetry using `testdata/telemetry`.
7. Update Phase 1 handoff docs with checks and any schema adjustments.
