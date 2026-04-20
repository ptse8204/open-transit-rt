# Phase 1 Handoff

## Phase

Phase 1 — Durable telemetry foundation

## Status

- Complete.
- Closure-polish pass complete; no remaining Phase 1 cleanup items are known.
- Active phase after this handoff: Phase 2 — Deterministic trip matching.

## What Was Implemented

- Added `internal/db` with env-driven `pgxpool` setup, startup ping, readiness ping support, and clean shutdown.
- Added `internal/telemetry` repository contracts and a Postgres implementation for telemetry persistence.
- Replaced in-memory storage in `cmd/telemetry-ingest` with DB-backed persistence.
- Added `POST /v1/telemetry` behavior that persists valid telemetry to Postgres and returns DB-derived ingest status and `received_at`.
- Added `/readyz` DB readiness to telemetry ingest while keeping `/healthz` as process liveness.
- Tightened `GET /v1/events` to a Phase 1 local/debug endpoint that requires `agency_id` and bounded `limit`, returns only that agency's persisted rows, and orders by `received_at DESC, id DESC`.
- Added atomic duplicate and out-of-order classification inside one DB transaction.
- Added deterministic transaction-scoped advisory locking for each `(agency_id, vehicle_id)` telemetry stream.
- Added parsed request payload storage in `telemetry_event.payload_json`.
- Added development agency seeding through `scripts/seed-dev.sql`, `make seed`, and `scripts/bootstrap-dev.sh`.
- Added DB-backed integration tests for telemetry insert/query, duplicates, equal timestamps, different devices, out-of-order events, agency scoping, unknown agencies, payload JSONB persistence, timestamp validation, and advisory-lock key determinism.

## What Was Designed But Intentionally Not Implemented Yet

- Deterministic trip matching.
- Service-day resolver and trip-instance matching.
- Vehicle assignment persistence behavior beyond the Phase 0 foundation table.
- GTFS import or GTFS Studio runtime behavior.
- Protobuf GTFS-RT Vehicle Positions.
- Trip Updates prediction adapters.
- Alerts feed behavior.
- Auth/admin protection for `/v1/events`; it is documented as local/debug only in Phase 1.
- Persistence of invalid JSON or invalid telemetry payloads as rejected rows. `rejected` remains reserved in the database enum for a later ingest-audit phase.

## Schema And Interface Changes

- Added migration `db/migrations/000002_telemetry_ingest_foundation.sql`.
- `telemetry_event.payload_json` is now `JSONB NOT NULL DEFAULT '{}'::jsonb`.
- Removed the old uniqueness constraint on `(agency_id, device_id, vehicle_id, observed_at)`.
- Added accepted-row uniqueness on `(agency_id, vehicle_id, observed_at)` with predicate `WHERE ingest_status = 'accepted'`.
- `device_id` no longer participates in canonical accepted telemetry uniqueness. It remains stored for audit/debug.
- Added indexes for latest accepted telemetry and recent agency-scoped debug listing.
- Added `telemetry.Repository` with:
  - `Store(ctx, event, payload)`
  - `LatestByVehicle(ctx, agencyID, vehicleID)`
  - `ListLatestByAgency(ctx, agencyID, limit)`
  - `ListEvents(ctx, agencyID, limit)`
- `Store` returns `accepted = true` only for `ingest_status = 'accepted'`. Duplicate and out-of-order telemetry persist but return `accepted = false` through the HTTP API.

## Dependency Changes

- No new external dependencies were added.
- Existing `github.com/jackc/pgx/v5` is now used by runtime services through `pgxpool`, not only by migrations.
- Existing Goose dependency is reused by DB-backed integration tests to apply migrations to isolated test databases.

## Migrations Added

- `db/migrations/000002_telemetry_ingest_foundation.sql`

Migration behavior:
- preserves existing `payload_json` rows while ensuring JSONB type, non-null default, and non-null constraint
- replaces telemetry canonical accepted uniqueness with vehicle-scoped accepted-row uniqueness
- keeps `ingest_status = 'rejected'` reserved but unused by Phase 1 runtime behavior

## Tests Added And Results

- Added handler tests under `cmd/telemetry-ingest`.
- Added DB-backed integration tests under `internal/telemetry`.
- Integration tests prefer creating an isolated temporary database from `TEST_DATABASE_URL`.
- If database creation is not available, tests fall back to an isolated temporary schema in the configured test database.
- Tests seed known agencies: `demo-agency`, `overnight-agency`, and `freq-agency`.
- `/readyz` handler behavior is covered for both DB-ready and DB-unavailable responses.
- Advisory-lock key derivation is covered by deterministic unit tests. Repository integration tests exercise telemetry classification through the locked `Store` path; no separate concurrent-ingest stress test was added in Phase 1.

Test results:
- `make test`: passed.
- `make test-integration`: passed with DB-backed telemetry tests.

## Checks Run And Blocked Checks

| Command | Result | Notes |
|---|---|---|
| `go mod tidy` | Passed | No dependency changes beyond existing pgx/Goose usage. |
| `make fmt` | Passed | Ran `gofmt -w ./cmd ./internal`. |
| `make test` | Passed | Unit tests and non-integration package tests passed. |
| `docker compose -f deploy/docker-compose.yml config` | Passed | Compose file renders successfully. |
| `make db-up` | Passed | PostGIS container running on host port `55432`. |
| `make migrate-up` | Passed | Applied `000002_telemetry_ingest_foundation.sql`. |
| `make migrate-down && make migrate-up && make migrate-status` | Passed | Smoke-tested rollback and re-application of migration `000002`. |
| `make migrate-status` | Passed | Reports migrations 1 and 2 applied. |
| `make test-integration` | Passed | DB-backed telemetry tests passed using isolated temporary DB setup. |
| `scripts/bootstrap-dev.sh` | Passed | Applies migrations and seeds development agencies. |
| `make validate` | Passed | Scaffold and durable telemetry file validation only; canonical validators are still not wired. |
| `git diff --check` | Passed | No whitespace errors. |
| Task equivalents | Not run | `task` is not installed; Makefile remains independently usable. |

Intermediate note: one `make test-integration` run failed while developing migration `000002` because Goose needed explicit statement boundaries around a `DO $$` block. The migration was updated with `-- +goose StatementBegin` / `-- +goose StatementEnd`, and the subsequent integration run passed.

## Known Issues

- `/v1/events` has no auth in Phase 1. It is bounded and agency-scoped, but should be disabled or protected before production deployment.
- Invalid JSON and invalid telemetry payloads are rejected without persistence. A later ingest-audit phase can use the reserved `rejected` status.
- `feed-vehicle-positions` still serves placeholder JSON from sample data and does not read persisted telemetry.
- `make validate` is still scaffold and durable telemetry file validation only; GTFS and GTFS-RT validators are documented but not wired.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/phase-plan.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-01.md`
  - `docs/requirements-2a-2f.md`
  - `docs/requirements-trip-updates.md`
- First files likely to edit:
  - `internal/state/`
  - `internal/telemetry/`
  - future schedule-query package under `internal/gtfs/` or equivalent
  - `db/migrations/`
  - `testdata/`
- Commands to run before coding:
  - `command -v go`
  - `go version`
  - `make fmt`
  - `make test`
  - `docker compose -f deploy/docker-compose.yml config`
  - `make db-up`
  - `make migrate-status`
  - `make test-integration`
- Known blockers:
  - Task is not installed, but Task is optional.
  - GTFS import is not implemented, so Phase 2 should use fixtures or narrow test repositories for schedule data rather than starting Phase 4.
- Recommended first implementation slice:
  - Build deterministic trip matching without starting GTFS import or protobuf Vehicle Positions.
  - Add agency-local service-day resolution and a minimal schedule query boundary.
  - Read latest accepted telemetry through the Phase 1 repository.
  - Preserve conservative `unknown` behavior for low-confidence or missing schedule data.
  - Add tests for after-midnight, frequency-based, unmatched, stale, and block-transition scenarios.
