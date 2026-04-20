# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 1 — Durable telemetry foundation

## Phase Status

- Phase 0 scaffolding is implemented.
- Phase 0 operational closure audit passed.
- Phase 0 closure-polish pass completed; no remaining Phase 0 cleanup items are known.
- Phase 1 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/phase-plan.md`
3. `docs/current-status.md`
4. `docs/handoffs/phase-00.md`
5. `docs/dependencies.md`
6. `docs/decisions.md`
7. `docs/codex-task.md`

## Current Objective

Replace in-memory telemetry with durable Postgres persistence and create the core DB/repository foundation.

## Exact First Commands

```bash
command -v go
go version
make fmt
make test
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
```

If Task is installed, optional equivalents may be run:

```bash
task fmt
task test
task migrate:status
```

## Known Blockers

- No active Phase 0 blockers remain.
- Task is not installed, but Task is optional and Makefile is independently usable.
- Docker must be running before DB-backed checks.

## First Files Likely To Edit

- `internal/db/`
- `internal/telemetry/`
- `cmd/telemetry-ingest/main.go`
- `docs/current-status.md`
- `docs/handoffs/phase-01.md`
- `docs/handoffs/latest.md`

## Phase 1 Entry Recommendation

Start Phase 1 durable telemetry by making telemetry persistence real without changing matcher or feed behavior:

1. Add `internal/db` with config loading and `pgxpool` connection setup.
2. Add telemetry repository interfaces in `internal/telemetry`.
3. Add a Postgres telemetry repository that inserts `telemetry_event` rows and queries latest events by agency/vehicle.
4. Update `cmd/telemetry-ingest` to use the repository instead of global process memory.
5. Add readiness behavior that reports DB connectivity.
6. Add DB-backed tests for valid insert/query, duplicate telemetry, and out-of-order telemetry using `testdata/telemetry`.
7. Update Phase 1 handoff docs with checks and any schema adjustments.

## Constraints To Preserve

- Mostly Go.
- Vehicle Positions first.
- Trip Updates pluggable.
- Draft GTFS separate from published GTFS.
- Conservative matching.
- Manual overrides take precedence over matching.
- No rider apps, payments, passenger accounts, or dispatcher CAD.
- External integrations stay behind documented adapters.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
