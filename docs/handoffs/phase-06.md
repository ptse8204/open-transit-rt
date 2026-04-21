# Phase 6 Handoff

## Phase

Phase 6 — Trip Updates and Alerts architecture

## Status

- Complete.
- Active phase after this handoff: Phase 7 — Prediction quality and operations workflows.

## What Was Implemented

- Added `internal/prediction.Adapter` as the narrow Trip Updates prediction boundary.
- Added `prediction.NoopAdapter` as the default safe adapter. It emits no Trip Updates and returns explicit `noop` diagnostics.
- Added Trip Updates diagnostics persistence through existing `feed_health_snapshot` rows with `feed_type = 'trip_updates'`.
- Added `internal/feed/tripupdates` with valid GTFS-RT Trip Updates protobuf rendering, JSON debug output, explicit `FeedHeader.timestamp`, deterministic entity ordering, and ordered `stop_time_update` entries.
- Added `cmd/feed-trip-updates` with `/healthz`, `/readyz`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/trip_updates.json`.
- Added exact Vehicle Positions URL config behavior:
  - `VEHICLE_POSITIONS_FEED_URL` is treated as the exact full Phase 3 protobuf URL.
  - otherwise `FEED_BASE_URL` must include `/public` and derives `/public/gtfsrt/vehicle_positions.pb`.
- Added `internal/feed/alerts` with valid empty GTFS-RT Alerts protobuf rendering and JSON debug output.
- Added `cmd/feed-alerts` with `/healthz`, `/readyz`, `/public/gtfsrt/alerts.pb`, and `/public/gtfsrt/alerts.json`.
- Added non-coupling tests proving `cmd/feed-vehicle-positions`, `cmd/telemetry-ingest`, and `cmd/gtfs-studio` do not depend on prediction or Trip Updates packages.
- Preserved Phase 3 Vehicle Positions, Phase 4 GTFS import, and Phase 5 GTFS Studio behavior.

## What Was Designed But Intentionally Not Implemented Yet

- Production stop-level ETA prediction quality.
- TheTransitClock or any other real external predictor adapter.
- Internal deterministic ETA predictor.
- Backtesting and prediction quality metrics.
- Alert authoring UI or persistence.
- Incident/manual-override-to-alert conversion.
- Canceled-trip Trip Updates plus corresponding Service Alerts.
- Compliance dashboard, canonical GTFS-RT validation, consumer ingestion workflows, rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Schema And Interface Changes

- Added `internal/prediction` package:
  - `Adapter`
  - `Request`
  - `Result`
  - `TripUpdate`
  - `StopTimeUpdate`
  - `Diagnostics`
  - `DiagnosticsRepository`
  - `PostgresDiagnosticsRepository`
- Added `internal/feed/tripupdates` package for Trip Updates snapshot/protobuf/debug rendering.
- Added `internal/feed/alerts` package for Alerts snapshot/protobuf/debug rendering.
- Added `cmd/feed-trip-updates`.
- Added `cmd/feed-alerts`.
- Added `VEHICLE_POSITIONS_FEED_URL` and `TRIP_UPDATES_MAX_VEHICLES` to `.env.example`.
- Updated Makefile and Taskfile run/validate targets for Trip Updates and Alerts.
- No database migrations were added. Trip Updates diagnostics reuse `feed_health_snapshot`.

Minimum Trip Updates diagnostics `details_json` fields:
- `adapter_name`
- `diagnostics_status`
- `diagnostics_reason`
- `active_feed_version_id`
- `input_counts`
- `vehicle_positions_url`
- `diagnostics_persistence_outcome`

Alerts diagnostics behavior:
- Alerts deferred status is JSON-only in Phase 6.
- Alerts does not write `feed_health_snapshot` rows in this slice.

## Dependency Changes

- No new external dependencies were added.
- The no-op Trip Updates adapter is internal Go code.
- TheTransitClock remains an optional future prediction backend behind `internal/prediction.Adapter`.

## Migrations Added

- None.

## Tests Added And Results

- Added `internal/prediction` tests for:
  - default no-op adapter behavior
  - DB-backed Trip Updates diagnostics persistence to `feed_health_snapshot`
- Added `internal/feed/tripupdates` tests for:
  - valid empty no-op protobuf output
  - explicit `FeedHeader.timestamp`
  - adapter input contract
  - deterministic entity ordering
  - ordered `stop_time_update` entries
  - missing active-feed behavior
  - adapter error behavior
  - no `SaveAssignment` coupling
- Added `internal/feed/alerts` tests for:
  - valid empty Alerts protobuf output
  - explicit `FeedHeader.timestamp`
  - JSON-only deferred diagnostics
- Added `cmd/feed-trip-updates` handler/config tests for:
  - Vehicle Positions URL derivation and validation
  - protobuf/JSON response headers and bodies
  - `Last-Modified` from snapshot `GeneratedAt`
  - readiness and wrong-method behavior
- Added `cmd/feed-alerts` handler tests for:
  - protobuf/JSON response headers and bodies
  - `Last-Modified` from snapshot `GeneratedAt`
  - readiness and wrong-method behavior
- Added `internal/architecture` non-coupling test.

Results:
- `make test`: passed.
- `make test-integration`: passed.

## Checks Run And Blocked Checks

| Command | Result | Notes |
|---|---|---|
| `command -v go` | Passed | `/usr/local/bin/go`. |
| `go version` | Passed | `go version go1.26.2 darwin/amd64`. |
| `make fmt` | Passed | Ran `gofmt -w ./cmd ./internal`. |
| `make test` | Passed | Unit and non-integration package tests passed. |
| `docker compose -f deploy/docker-compose.yml config` | Passed | Compose file renders successfully. |
| `make db-up` | Passed | PostGIS container running on host port `55432`. |
| `make migrate-status` | Passed | Reports migrations 1 through 5 applied. |
| `make test-integration` | Passed | DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, and Trip Updates diagnostics tests passed. |
| `make validate` | Passed | Phase 6 file smoke only; canonical validators remain unwired. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Task equivalents were not run; Makefile remains independently usable and Task has been optional in prior phases.

## Known Issues

- Trip Updates endpoint intentionally returns an empty feed with no-op diagnostics until a real predictor is implemented.
- Missing active GTFS produces a valid empty Trip Updates feed with explicit diagnostics.
- Alerts endpoint intentionally returns an empty feed with JSON-only deferred diagnostics; no alert authoring/persistence exists yet.
- Canonical GTFS-Realtime validators remain documented but unwired.
- GTFS Studio auth is still minimal/dev-only and not production-grade.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-06.md`
  - `docs/phase-plan.md`
  - `docs/requirements-trip-updates.md`
  - `docs/requirements-calitp-compliance.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `internal/prediction/`
  - `internal/feed/tripupdates/`
  - `internal/feed/alerts/`
  - `internal/state/` only for operation-state inputs needed by prediction
  - `db/migrations/` only if Phase 7 adds alert/prediction persistence beyond existing health snapshots
  - docs and handoff files
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
  - No real predictor is implemented yet.
  - Alerts authoring/persistence is not implemented yet.
  - Canonical GTFS-RT validators are documented but not wired.
- Recommended first implementation slice:
  - Start Phase 7 by adding the first real Trip Updates prediction behavior behind `internal/prediction.Adapter`, using active published GTFS, persisted latest telemetry, and persisted assignments.
  - Preserve the Phase 6 public endpoint shapes and non-coupling guarantees.
  - Add alert persistence/authoring only after deciding whether alerts are operator-authored, incident-derived, or both.
