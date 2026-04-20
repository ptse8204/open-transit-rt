# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current Repository State

This repository is an early-stage starter for **Open Transit RT**.

Phase 0 scaffolding and operational closure are complete. The repo can format, test, start Postgres/PostGIS, run migrations, and execute the bootstrap flow.

## What Exists Now

### Repo guidance and architecture docs
The repo has:
- `AGENTS.md`
- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/dependencies.md`
- `docs/phase-plan.md`
- `docs/decisions.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/handoffs/latest.md`

### Phase 0 scaffolding
The repo now has:
- `.env.example`
- `Taskfile.yml`
- independently usable `Makefile`
- `cmd/migrate`
- versioned migrations under `db/migrations`
- PostGIS-backed Docker Compose configuration on host port `55432`
- `scripts/bootstrap-dev.sh`
- deterministic fixtures under `testdata/`
- handoff template and Phase 0 handoff under `docs/handoffs/`

### Starter code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`

These services are scaffolds, not complete implementations.

## Schema Source Of Truth

Migrations under `db/migrations` are the source of truth for executable schema changes and are applied through `cmd/migrate`.

`db/schema.sql` is deprecated as an executable schema. It remains only as a compatibility pointer to the migrations directory and must not be edited independently.

## What Does Not Exist Yet

The following are still missing or incomplete unless a later handoff says otherwise:

- durable DB-backed telemetry runtime implementation
- repository interfaces and DB connection package
- complete GTFS import pipeline
- complete GTFS Studio draft/publish workflow
- deterministic trip matcher with real edge-case handling
- real GTFS-RT Vehicle Positions protobuf feed from persisted data
- Trip Updates adapter implementation
- Alerts feed implementation
- compliance dashboard
- consumer ingestion workflow
- robust auth and role handling
- manual override workflows
- production observability and SLO reporting

## Current Phase

**Active phase:** Phase 1 — Durable telemetry foundation

Phase 0 is operationally closed. The next Codex instance should start with `docs/handoffs/latest.md`.

## Architecture Posture

The codebase must preserve these long-term rules:
- mostly Go backend
- Postgres/PostGIS source of truth
- Vehicle Positions first
- Trip Updates pluggable
- draft GTFS separate from published GTFS
- conservative matching
- external dependencies isolated behind adapters
- no rider apps, payments, passenger accounts, or dispatcher CAD scope

## Phase 0 Closure Audit Results

Checked during Phase 0 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `command -v gofmt`: passed, `/usr/local/bin/gofmt`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `go mod tidy`: passed and generated `go.sum`.
- `make fmt`: passed.
- `make test`: passed.
- `make db-up`: passed after changing local PostGIS host port to `55432`.
- `make migrate-up`: passed and applied `000001_initial_schema.sql`.
- `make migrate-status`: passed and reports migration version 1 applied.
- `make test-integration`: passed; there are no integration test files yet.
- `scripts/bootstrap-dev.sh`: passed and reports no pending migrations.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make validate`: passed as a Phase 0 placeholder; validators are documented but not wired.
- `make lint`: passed as a no-op fallback; `golangci-lint` is not installed.
- `git diff --check`: passed.
- Task equivalents were not run because `task` is not installed; Task remains optional because Makefile is independently usable.

## Next Recommended Step

Begin Phase 1 using the exact recommendation in `docs/handoffs/latest.md` and `docs/handoffs/phase-00.md`.

The first implementation slice should be:
1. add a shared DB package using `pgxpool`
2. add telemetry repository interfaces and Postgres implementation
3. wire `cmd/telemetry-ingest` to persist telemetry
4. update health/readiness behavior
5. add DB-backed tests using `testdata/telemetry`

## What Not To Do Next

Do not:
- jump straight into Trip Updates implementation
- add rider-facing functionality
- add payments, passenger accounts, or dispatcher CAD
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths once real feed generation starts
