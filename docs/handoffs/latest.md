# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 1 — Durable telemetry foundation

## Phase Status

- Phase 0 is complete.
- Phase 1 is ready to start once the Go toolchain is available.

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

## Highest-Priority Tasks

- Install or expose Go on `PATH`.
- Run `go mod tidy` to generate `go.sum`.
- Validate and apply Phase 0 migrations.
- Add shared DB connection package using `pgxpool`.
- Add telemetry repository interfaces and Postgres implementation.
- Wire `cmd/telemetry-ingest` to persist telemetry.
- Add DB-backed tests for insert/query, duplicate telemetry, and out-of-order telemetry.

## Exact First Commands

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

## First Files Likely To Edit

- `go.mod`
- `go.sum`
- `internal/db/`
- `internal/telemetry/`
- `cmd/telemetry-ingest/main.go`
- `docs/current-status.md`
- `docs/handoffs/phase-01.md`
- `docs/handoffs/latest.md`

## Known Blockers

- `go` is not on `PATH`.
- `gofmt` is not on `PATH`.
- `task` is not on `PATH`, but Makefile is independently usable.
- Migration execution and tests are blocked until Go is available.

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
