# Phase 2 Handoff

## Phase

Phase 2 — Deterministic trip matching

## Status

- Complete.
- Active phase after this handoff: Phase 3 — Vehicle Positions production feed.

## What Was Implemented

- Added `internal/gtfs` as a narrow schedule-query boundary over the existing published GTFS tables.
- Added agency-local service-day resolution using agency timezone, including previous-local-date candidate evaluation for after-midnight trips.
- Added GTFS time parsing and formatting for values beyond `24:00:00`.
- Replaced the placeholder-only matcher internals with a deterministic matcher engine in `internal/state`.
- Added candidate scoring for trip hint, shape proximity, movement direction, stop progress, schedule fit, previous-assignment continuity, and block continuity.
- Added exact frequency instance generation for `exact_times=1`.
- Added conservative non-exact frequency-window identity behavior for `exact_times=0`.
- Preserved repeated trip instances with the same `trip_id` but different `start_time`; they are not collapsed into one logical instance.
- Added manual override precedence in matcher logic.
- Added explicit unknown assignment persistence for stale, ambiguous, low-confidence, and missing-schedule cases.
- Added a Postgres assignment repository that closes prior active rows and inserts the new assignment in one transaction.
- Added incident insertion linked to the persisted assignment row.
- Added a small reason-code, degraded-state, and incident taxonomy.
- Updated validation smoke targets to include Phase 2 migration coverage.

## What Was Designed But Intentionally Not Implemented Yet

- GTFS-RT Vehicle Positions protobuf output.
- DB-backed Vehicle Positions public feed behavior.
- GTFS import runtime logic.
- GTFS Studio runtime logic.
- Trip Updates prediction adapters.
- Alerts runtime logic.
- Operator UI for manual overrides.
- Stable public diagnostics schema for matcher scoring.

## Schema And Interface Changes

- Added migration `db/migrations/000003_deterministic_matching.sql`.
- `vehicle_trip_assignment.service_date` is now nullable only for truly unresolved cases. Unknown assignments should still carry a service date whenever agency timezone and observed timestamp can be resolved.
- Added `vehicle_trip_assignment.block_id`.
- Added `vehicle_trip_assignment.telemetry_event_id`.
- Added `vehicle_trip_assignment.degraded_state`.
- Added `vehicle_trip_assignment.score_details_json`.
- Added `incident.vehicle_trip_assignment_id`.
- Added `internal/gtfs.Repository` with:
  - `Agency(ctx, agencyID)`
  - `ActiveFeedVersion(ctx, agencyID)`
  - `ListTripCandidates(ctx, agencyID, feedVersionID, serviceDate)`
- Added `internal/state.Engine` with `MatchEvent(ctx, telemetry.StoredEvent, now)`.
- Added `internal/state.Repository` with:
  - `ActiveManualOverride(ctx, agencyID, vehicleID, at)`
  - `CurrentAssignment(ctx, agencyID, vehicleID)`
  - `SaveAssignment(ctx, assignment, incidents)`
- `score_details_json` is intentionally loose debug JSON for Phase 2, not a stable structured public schema.

## Dependency Changes

- No new external dependencies were added.
- Existing `pgx` and Goose usage was extended to matcher repositories and tests.

## Migrations Added

- `db/migrations/000003_deterministic_matching.sql`

Migration behavior:
- allows unresolved assignment rows to have null `service_date`
- adds assignment block, telemetry-event linkage, degraded-state, and loose score-details fields
- links incidents back to assignment rows
- adds indexes for current assignment lookup, degraded assignment lookup, telemetry-event lookup, and incident-to-assignment lookup

## Tests Added And Results

- Added unit tests under `internal/gtfs` for after-midnight GTFS time parsing and agency-local service-day resolution.
- Added matcher unit tests under `internal/state` for:
  - after-midnight service-date selection
  - exact frequency generated instances
  - non-exact frequency conservative identity
  - explicit unknown row behavior for stale telemetry
  - missing-shape degradation without automatic match rejection
  - manual override precedence
  - block-transition reason recording
  - ambiguous candidates
  - no-schedule unknown behavior
- Added DB-backed matcher integration tests under `internal/state` for:
  - matching persisted latest accepted telemetry
  - stale telemetry unknown assignment plus incident
  - manual override precedence
  - missing-shape degraded match

Test results:
- `make test`: passed.
- `make test-integration`: passed with DB-backed telemetry and matcher tests.

## Checks Run And Blocked Checks

| Command | Result | Notes |
|---|---|---|
| `command -v go` | Passed | `/usr/local/bin/go`. |
| `go version` | Passed | `go version go1.26.2 darwin/amd64`. |
| `make fmt` before coding | Blocked | Plan Mode was active and `make fmt` runs `gofmt -w ./cmd ./internal`; it was run after implementation. |
| `make fmt` | Passed | Ran `gofmt -w ./cmd ./internal`. |
| `make test` | Passed | Unit tests and non-integration package tests passed. |
| `docker compose -f deploy/docker-compose.yml config` | Passed | Compose file renders successfully. |
| `make db-up` | Passed | PostGIS container running on host port `55432`. |
| `make migrate-up` | Passed | Applied `000003_deterministic_matching.sql`. |
| `make migrate-status` | Passed | Reports migrations 1, 2, and 3 applied. |
| `make migrate-down && make migrate-up && make migrate-status` | Passed | Smoke-tested rollback and re-application of migration `000003`. |
| `make test-integration` | Passed | DB-backed telemetry and matcher tests passed using isolated temporary DB setup. |
| `make validate` | Passed | Phase 2 scaffold, telemetry, and matcher-file validation only; canonical validators are still not wired. |
| `git diff --check` | Passed | No whitespace errors. |
| Task equivalents | Not run | `task` is not installed; Makefile remains independently usable. |

## Known Issues

- `cmd/feed-vehicle-positions` still serves placeholder JSON from sample data. Phase 3 must replace this with DB-backed latest telemetry and persisted assignments.
- `score_details_json` is loose debug JSON and must not be consumed as a stable public contract.
- GTFS import is still not implemented. Matcher integration tests seed schedule rows directly through test helpers; this must not evolve into runtime import logic.
- Operator UI for manual overrides is not implemented, although matcher precedence and persistence behavior exist.
- Canonical GTFS and GTFS-RT validators remain documented but not wired.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-02.md`
  - `docs/phase-plan.md`
  - `docs/codex-task.md`
  - `docs/requirements-2a-2f.md`
  - `docs/requirements-trip-updates.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `cmd/feed-vehicle-positions/`
  - `internal/feed/`
  - `internal/state/`
  - `internal/telemetry/` only if a narrow read/query addition is required
  - `testdata/expected/`
  - `docs/current-status.md`
  - `docs/handoffs/phase-03.md`
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
  - GTFS import is not implemented; do not start Phase 4 while implementing Phase 3.
  - GTFS-RT protobuf tooling is documented as a future dependency but not wired yet.
- Recommended first implementation slice:
  - Replace placeholder Vehicle Positions JSON behavior with DB-backed latest accepted telemetry and persisted assignments.
  - Add GTFS-RT protobuf Vehicle Positions serialization and a stable public protobuf endpoint.
  - Preserve explicit unknown and stale behavior from Phase 2.
  - Do not implement Trip Updates, Alerts, GTFS import, or GTFS Studio.
