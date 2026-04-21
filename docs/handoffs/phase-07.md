# Phase 7 Handoff

## Phase

Phase 7 — Prediction quality and operations workflows

## Status

- Complete.
- Active phase after this handoff: Phase 8 — Compliance and consumer workflow.

## What Was Implemented

- Added `prediction.DeterministicAdapter` as the first real internal Trip Updates predictor behind `internal/prediction.Adapter`.
- Made `cmd/feed-trip-updates` default to `TRIP_UPDATES_ADAPTER=deterministic`, while preserving `TRIP_UPDATES_ADAPTER=noop`.
- Generated non-empty Trip Updates for defensible in-service assignments using active published GTFS, latest accepted telemetry, and current persisted assignments.
- Added conservative scheduled, exact-frequency, non-exact-frequency, and after-midnight Trip Updates behavior.
- Added deadhead and layover no-prediction behavior with explicit withheld reasons.
- Added first-pass disruption behavior:
  - canceled trips emit conservative `CANCELED` Trip Updates.
  - added trips, short turns, and detours are withheld unless later evidence/contracts make them defensible.
- Added cancellation-to-alert linkage signal in persisted review details:
  - `expected_alert_missing=true`
  - `cancellation_alert_linkage_status="missing_alert_authoring_deferred"`
  - `linked_review_reason="canceled_trip_requires_service_alert"`
- Added first-class prediction metrics in diagnostics and persisted `feed_health_snapshot.details_json`.
- Added prediction operations repository behavior for override create, replace, clear, expiry reads, review item persistence, review status transitions, and audit logging.
- Kept matcher override consumption limited to `trip_assignment` and `service_state`; prediction-only disruption overrides are consumed through `prediction.OperationsRepository`.
- Added prediction review queue lifecycle semantics through `incident.status`: `open`, `resolved`, and `deferred`.
- Preserved Phase 6 public endpoint shapes and non-coupling guarantees.

## What Was Designed But Intentionally Not Implemented Yet

- Production-grade ETA quality, learned travel-time history, and backtesting.
- Public Alerts authoring/persistence and automatic satisfaction of canceled-trip alert linkage.
- Full operator UI for override and prediction review workflows.
- Vehicle swap workflow beyond the existing override schema foundation.
- Detour path modeling and short-turn boundary prediction beyond conservative withholding.
- Canonical GTFS and GTFS-Realtime validators.
- Rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Schema And Interface Changes

- Added migration `000006_prediction_operations.sql` to allow `incident.status = 'deferred'`.
- Extended `prediction.Diagnostics` with metrics and adapter details.
- Added `prediction.Metrics`.
- Added prediction operation interfaces and records:
  - `OperationsRepository`
  - `OverrideRecord`
  - `OverrideInput`
  - `ReviewItem`
  - `ReviewFilter`
  - `ReviewStatus`
- Added `prediction.PostgresOperationsRepository`.
- Extended Trip Updates JSON diagnostics with prediction metrics and optional diagnostics details.
- Added Trip Updates config environment variables:
  - `TRIP_UPDATES_ADAPTER`
  - `TRIP_UPDATES_STALE_TELEMETRY_TTL_SECONDS`
  - `TRIP_UPDATES_ASSIGNMENT_CONFIDENCE_THRESHOLD`
  - `TRIP_UPDATES_MAX_SCHEDULE_DEVIATION_SECONDS`
  - `TRIP_UPDATES_DUPLICATE_CONFIDENCE_GAP`

Coverage denominator semantics:
- `coverage_percent` is ETA coverage only: emitted non-canceled Trip Updates divided by eligible in-service ETA prediction candidates.
- Canceled trips are excluded from that denominator and tracked separately by canceled-trip and cancellation-alert-linkage metrics.

Review queue lifecycle semantics:
- Prediction review items are persisted as `incident` rows with `incident_type = 'prediction_review'`.
- Review items are not append-only; they support minimal states `open`, `resolved`, and `deferred`.
- Status changes write `audit_log` rows.

## Dependency Changes

- No new external dependencies were added.
- The first real predictor is internal Go code.
- Optional external predictors such as TheTransitClock remain future adapter implementations behind `internal/prediction.Adapter`.

## Migrations Added

- `db/migrations/000006_prediction_operations.sql`

## Tests Added And Results

- Added deterministic adapter unit tests for:
  - matched scheduled Trip Updates
  - after-midnight prediction times
  - exact and non-exact frequency behavior
  - deadhead and layover suppression
  - cancellation, added-trip, short-turn, and detour behavior
  - duplicate trip-instance ambiguity
  - prediction metrics and review item creation
- Added DB-backed prediction operations tests for:
  - override create, replace, clear, and expiry semantics
  - audit log coverage for override actions
  - prediction review item persistence
  - deferred review status lifecycle
- Added Trip Updates builder coverage proving deterministic adapter output flows to protobuf and metrics.
- Added command config tests for deterministic default and no-op fallback.
- Existing non-coupling tests still prove Vehicle Positions, telemetry ingest, and GTFS Studio do not depend on prediction or Trip Updates packages.

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
| `make migrate-up` | Passed | Applied `000006_prediction_operations.sql`. |
| `make migrate-status` | Passed | Reports migrations 1 through 6 applied. |
| `make test-integration` | Passed | DB-backed integration tests passed. |
| `make validate` | Passed | Phase 7 file smoke only. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Canonical GTFS and GTFS-Realtime validators remain documented but unwired.

## Known Issues

- Trip Updates predictions are conservative schedule-deviation projections, not production-grade learned ETAs.
- Prediction review rows can accumulate across feed refreshes; later UI/workflow work should decide deduplication and assignment ownership.
- Canceled trips record missing-alert linkage, but public Alerts authoring/persistence remains deferred.
- Added trips, short turns, and detours are conservatively withheld.
- Full operator UI and auth remain future work.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/phase-plan.md`
  - `docs/requirements-calitp-compliance.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `internal/feed/alerts/`
  - `internal/feed/tripupdates/`
  - `internal/prediction/`
  - `internal/gtfs/`
  - `db/migrations/` if compliance/report persistence changes are needed
  - compliance docs and handoff files
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
  - Canonical GTFS and GTFS-Realtime validators are documented but not wired.
  - Alerts authoring/persistence is not implemented yet.
  - Prediction review workflow has repository behavior but no full operator UI.
- Recommended first implementation slice:
  - Begin Phase 8 by wiring canonical validation/reporting and public feed metadata behind stable interfaces, while preserving Phase 7 prediction adapter and review workflow boundaries.
