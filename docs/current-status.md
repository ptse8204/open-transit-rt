# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current Repository State

This repository is an early-stage starter for **Open Transit RT**.

Phase 0 scaffolding and Phase 1 durable telemetry foundation are complete. The repo can format, test, start Postgres/PostGIS, run migrations, seed local agencies, execute the bootstrap flow, and run DB-backed telemetry integration tests.

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

### Runtime code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`

`cmd/telemetry-ingest` now persists valid telemetry to Postgres through a telemetry repository. `agency-config` and `feed-vehicle-positions` remain starter scaffolds; `feed-vehicle-positions` still serves placeholder JSON from sample data and does not read persisted telemetry yet.

### Phase 1 telemetry foundation
The repo now has:
- `internal/db` with `pgxpool` connection setup and readiness ping support
- `internal/telemetry` repository interfaces and Postgres implementation
- DB-backed telemetry ingest in `cmd/telemetry-ingest`
- `/healthz` liveness and `/readyz` DB readiness behavior for telemetry ingest
- agency-scoped, bounded `/v1/events` debug listing
- durable parsed request payload storage in `telemetry_event.payload_json`
- atomic duplicate and out-of-order classification inside a transaction with a deterministic advisory lock
- DB-backed integration tests using `testdata/telemetry`
- development agency seeding through `scripts/seed-dev.sql`

## Schema Source Of Truth

Migrations under `db/migrations` are the source of truth for executable schema changes and are applied through `cmd/migrate`.

`db/schema.sql` is deprecated as an executable schema. It is intentionally a comment-only pointer to the migrations directory and must not be edited independently.

## What Does Not Exist Yet

The following are still missing or incomplete unless a later handoff says otherwise:

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

**Active phase:** Phase 2 — Deterministic trip matching

Phase 1 is operationally closed. The next Codex instance should start with `docs/handoffs/latest.md`.

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
- `make test-integration`: passed; this is currently a Phase 0 integration smoke path that verifies database reachability, migration visibility, and package compilation. There are no DB-backed integration test files yet.
- `scripts/bootstrap-dev.sh`: passed and reports no pending migrations.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make validate`: passed Phase 0 scaffold validation. It checks required migration and fixture scaffolding only; canonical GTFS and GTFS-RT validators are documented but not wired.
- `make lint`: passed optional fallback. `golangci-lint` is not installed, and future CI should make lint required once configured.
- `git diff --check`: passed.
- handoff path audit: passed; repo docs use `docs/handoffs/latest.md` and the retired singular path has been removed.
- Task equivalents were not run because `task` is not installed; Task remains optional because Makefile is independently usable.

## Phase 1 Closure Audit Results

Checked during Phase 1 closure:
- `go mod tidy`: passed.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed.
- `make migrate-up`: passed and applied `000002_telemetry_ingest_foundation.sql`.
- `make migrate-status`: passed and reports migration versions 1 and 2 applied.
- `make test-integration`: passed with DB-backed telemetry tests using an isolated temporary database.
- `scripts/bootstrap-dev.sh`: passed and seeds `demo-agency`, `overnight-agency`, and `freq-agency`.
- `make validate`: passed scaffold validation. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.
- Optional Task equivalents were not run because `task` is not installed.

## Next Recommended Step

Begin Phase 2 using the exact recommendation in `docs/handoffs/latest.md` and `docs/handoffs/phase-01.md`.

The first implementation slice should be:
1. define the GTFS schedule query surface needed by deterministic matching without implementing GTFS import yet
2. add agency-local service-day resolution
3. use persisted latest telemetry from the Phase 1 repository
4. keep low-confidence matches as `unknown`
5. add tests for after-midnight, frequency-based, stale, unmatched, and block-transition fixtures

## What Not To Do Next

Do not:
- jump straight into Trip Updates implementation
- add rider-facing functionality
- add payments, passenger accounts, or dispatcher CAD
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths once real feed generation starts
