# Phase 3 Handoff

## Phase

Phase 3 — Vehicle Positions production feed

## Status

- Complete.
- Active phase after this handoff: Phase 4 — GTFS import and publish pipeline.

## What Was Implemented

- Replaced placeholder Vehicle Positions behavior with DB-backed generation from latest accepted telemetry and current persisted assignments.
- Added official GTFS-RT protobuf Vehicle Positions serialization with `FeedHeader.gtfs_realtime_version = "2.0"`, `FULL_DATASET`, and snapshot-generated timestamps.
- Added stable public protobuf endpoint: `GET /public/gtfsrt/vehicle_positions.pb`.
- Kept JSON debug endpoint: `GET /public/gtfsrt/vehicle_positions.json`, backed by the same snapshot model as protobuf.
- Added `/readyz` DB readiness to `cmd/feed-vehicle-positions` and kept `/healthz`.
- Added startup config validation: `AGENCY_ID` is required, and numeric feed settings must be valid.
- Added `internal/feed.VehiclePositionsSnapshot` as the single immutable per-request model for protobuf and JSON rendering.
- Added vehicle cap behavior through `VEHICLE_POSITIONS_MAX_VEHICLES`, default `2000`, applied before stale/suppression/publication rules.
- Defined `telemetry.Repository.ListLatestByAgency` as newest accepted row per vehicle, ordered by `observed_at DESC, id DESC`.
- Added deterministic stale policy: stale-but-unsuppressed vehicles are included without trip descriptors; suppressed vehicles are omitted from protobuf and visible in JSON debug.
- Added normal successful empty protobuf behavior for no telemetry and all-suppressed snapshots.
- Added JSON debug fields for per-vehicle telemetry age, inclusion, assignment publishability, assignment/telemetry mismatch, trip descriptor publication, and the winning omission reason.
- Added `Last-Modified` derived from `snapshot.generated_at` for both protobuf and JSON endpoints.
- Preserved Phase 2 conservative semantics by omitting trip descriptors for unknown, stale, ambiguous, low-confidence, missing-schedule, stale telemetry, and assignment/telemetry mismatch cases.

## What Was Designed But Intentionally Not Implemented Yet

- Trip Updates.
- Alerts.
- GTFS import runtime logic.
- GTFS Studio runtime logic.
- Canonical GTFS or GTFS-RT validator execution.
- Public metadata/license/discoverability pages.
- Operator UI for manual overrides.
- Production auth and role handling.
- Metrics export and SLO dashboards.

## Schema And Interface Changes

- No database migrations were added.
- Added `state.Repository.ListCurrentAssignments(ctx, agencyID, vehicleIDs)` for bulk active-assignment lookup.
- Extended `state.PostgresRepository` with the new bulk method.
- Hardened the contract of `telemetry.Repository.ListLatestByAgency`: latest accepted vehicle rows first by observed time and id.
- Added `internal/feed.VehiclePositionsConfig`, `VehiclePositionsBuilder`, `VehiclePositionsSnapshot`, and JSON debug DTOs.
- `cmd/feed-vehicle-positions` now depends on repository interfaces and passes concrete Postgres repositories only at command wiring.

## Dependency Changes

- Added `github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs` v1.0.0.
- Upgraded/declared `google.golang.org/protobuf` v1.36.11.
- Updated `docs/dependencies.md` to make GTFS-RT protobuf serialization a Phase 3 runtime dependency isolated to feed boundary code.

## Migrations Added

- None.

## Tests Added And Results

- Added `internal/feed` tests for protobuf validity, header fields, matched entity content, no telemetry, all-suppressed empty feeds, no assignments, stale/suppressed behavior, manual override publication and no-trip omission behavior, truncation-before-publication behavior, non-exact frequency `UNSCHEDULED`, explicit true-north `bearing: 0`, malformed/omitted optional telemetry fields, assignment/telemetry mismatch, telemetry age debug output, and deterministic protobuf bytes.
- Added `cmd/feed-vehicle-positions` handler tests for protobuf headers, JSON debug output, missing config, method handling, readiness, no-telemetry success, and repository errors.
- Added `internal/state` DB-backed coverage for `ListCurrentAssignments`, including active-only behavior, missing vehicles, and agency scoping.
- Added `internal/telemetry` DB-backed coverage for `ListLatestByAgency` ordering by `observed_at DESC, id DESC`.

Results:
- `make test`: passed.
- `make test-integration`: passed.

## Checks Run And Blocked Checks

| Command | Result | Notes |
|---|---|---|
| `go mod tidy` | Passed | Added GTFS-RT protobuf dependencies. |
| `make fmt` | Passed | Ran `gofmt -w ./cmd ./internal`. |
| `make test` | Passed | Unit and non-integration package tests passed. |
| `docker compose -f deploy/docker-compose.yml config` | Passed | Compose file renders successfully. |
| `make db-up` | Passed | PostGIS container running on host port `55432`. |
| `make migrate-status` | Passed | Reports migrations 1, 2, and 3 applied. |
| `make test-integration` | Passed | DB-backed telemetry and matcher tests passed using isolated temporary database setup. |
| `make validate` | Passed | Phase 3 smoke only; canonical GTFS and GTFS-RT validators remain documented but not wired. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Task equivalents were not run; Task has been optional in prior phases and Makefile remains independently usable.

## Known Issues

- Canonical GTFS-RT validator tooling is still not wired; Phase 3 tests prove protobuf marshal/unmarshal and initialization but do not run MobilityData validation.
- GTFS import is still not implemented; feed and matcher tests seed GTFS rows directly where schedule context is needed.
- JSON debug output is diagnostic and not a stable public integration schema.
- Feed generation uses read-committed repository reads rather than one cross-table serializable transaction; automatic assignment telemetry linkage prevents false trip certainty.
- `AGENCY_ID` scopes the feed service to one agency in Phase 3.
- Manual override workflows/UI are not implemented, though persisted manual override assignments publish correctly when current telemetry is fresh.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/phase-plan.md`
  - `docs/codex-task.md`
  - `docs/requirements-2a-2f.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `internal/gtfs/`
  - `cmd/migrate/` only if migration behavior needs extension
  - `db/migrations/`
  - `testdata/gtfs/`
  - `docs/current-status.md`
  - `docs/handoffs/phase-04.md`
  - `docs/handoffs/latest.md`
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
  - Canonical validators are documented but not wired.
  - Existing GTFS tests seed DB tables directly; Phase 4 must build real import/publish runtime behavior rather than extending test-only seed helpers into production.
- Recommended first implementation slice:
  - Begin Phase 4 by adding GTFS ZIP staging and validation around the existing published GTFS tables.
  - Keep draft GTFS and published feed versions separate.
  - Preserve the stable Vehicle Positions endpoint and do not start Trip Updates or Alerts.
