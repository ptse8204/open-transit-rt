# Phase 4 Handoff

## Phase

Phase 4 — GTFS import and publish pipeline

## Status

- Complete.
- Active phase after this handoff: Phase 5 — GTFS Studio draft/publish model.

## What Was Implemented

- Added `cmd/gtfs-import` for runtime GTFS ZIP import.
- Added `internal/gtfs.ImportService` to parse, validate, report, stage, and atomically publish GTFS ZIPs.
- Added internal validation for the exact Phase 4 required input rule: `agency.txt`, `routes.txt`, `stops.txt`, `trips.txt`, `stop_times.txt`, and at least one usable service source from `calendar.txt` or `calendar_dates.txt`.
- Added support for `agency`, `routes`, `stops`, `trips`, `stop_times`, `calendar`, `calendar_dates`, optional `shapes`, optional `frequencies`, and optional `block_id`.
- Preserved GTFS time strings exactly as imported in canonical published tables, including times beyond `24:00:00`; parsed seconds remain validation/query logic only.
- Preserved `block_id` from `trips.txt` when present and covered downstream visibility through `gtfs.PostgresRepository.ListTripCandidates`.
- Built PostGIS point geometry for stops and shape points, and `gtfs_shape_line` from ordered shape points when a shape has at least two points.
- Implemented transactional publish: create staged feed version, insert GTFS rows, write validation report, retire previous active feed, activate new feed, update import status, and audit successful publish in one transaction.
- Implemented failed-import behavior: validation failures create no `feed_version`; publish failures roll back partial GTFS rows; failed import rows keep `gtfs_import.feed_version_id` `NULL` when the failure report write succeeds.
- Made failure-report write failure explicit: importer/CLI returns a clear nonzero failure and does not claim failure metadata was stored.

## What Was Designed But Intentionally Not Implemented Yet

- GTFS Studio runtime editing flows.
- Static GTFS public ZIP serving/export.
- MobilityData canonical GTFS Validator execution.
- Trip Updates.
- Alerts.
- Compliance dashboard and consumer ingestion workflows.
- Rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Schema And Interface Changes

- Added `gtfs_import` table for durable import attempts, source checksum/size, status, report JSON, actor/notes, and nullable `feed_version_id`.
- Added `validation_report.gtfs_import_id` to link schedule validation reports to import attempts.
- Added `gtfs.ImportService`, `gtfs.ImportOptions`, `gtfs.ImportResult`, `gtfs.ImportReport`, and `gtfs.ImportError`.
- `cmd/gtfs-import` is a thin wrapper over the import service and outputs the service result as JSON.

## Dependency Changes

- No new external dependencies were added.
- Phase 4 uses the Go standard library for ZIP/CSV parsing and the existing pgx/Postgres/PostGIS stack.
- `docs/dependencies.md` now clarifies that Phase 4 internal validation is not a canonical compliance validator.

## Migrations Added

- `db/migrations/000004_gtfs_import_pipeline.sql`.

## Tests Added And Results

- Added parser tests for valid GTFS, after-midnight time preservation, optional missing `shapes.txt`/`frequencies.txt`, required service source enforcement, malformed input, and multi-agency route conflict rejection.
- Added DB-backed import tests for valid import activation, failed import report storage with no staged feed version, active feed switching, downstream `block_id` visibility, and shape-line construction.
- Added CLI smoke tests for required flags, successful JSON result output, and failed-import JSON output when report storage fails.

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
| `make migrate-up` | Passed | Applied `000004_gtfs_import_pipeline.sql`. |
| `make migrate-status` | Passed | Reports migrations 1 through 4 applied. |
| migration down/up smoke for `000004` | Passed | `make migrate-down`, `make migrate-up`, `make migrate-status`. |
| `make test-integration` | Passed | DB-backed telemetry, matcher, Vehicle Positions, and GTFS import tests passed. |
| `make validate` | Passed | Phase 4 smoke only; canonical validators remain unwired. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Task equivalents were not run; Makefile remains independently usable and Task has been optional in prior phases.

## Known Issues

- MobilityData GTFS Validator is still not wired; internal validation blocks obvious unsafe imports but is not a compliance claim.
- Runtime import input is GTFS ZIP only. Directory parsing exists only as a test-fixture convenience that creates ZIPs.
- Original uploaded ZIP bytes are not stored in Postgres and no public static GTFS ZIP endpoint was added.
- Existing current vehicle assignments are not cleared on feed switch; the next matcher pass uses the new active feed.
- Import report rows require agency metadata to be available or bootstrappable from `agency.txt`; if the separate failure-report write fails, the CLI reports that explicitly.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-04.md`
  - `docs/phase-plan.md`
  - `docs/codex-task.md`
  - `docs/requirements-2a-2f.md`
  - `docs/requirements-calitp-compliance.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `internal/gtfs/`
  - `db/migrations/`
  - `cmd/*` only for minimal admin/Studio entrypoints if needed
  - `docs/current-status.md`
  - `docs/handoffs/phase-05.md`
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
  - GTFS Studio draft runtime behavior is not implemented yet.
- Recommended first implementation slice:
  - Begin Phase 5 by adding draft GTFS editing/publish models that remain separate from published `feed_version` tables.
  - Reuse the Phase 4 validation/publish semantics for draft publish.
  - Do not start Trip Updates, Alerts, rider apps, payments, passenger accounts, CAD, or marketplace workflows.
