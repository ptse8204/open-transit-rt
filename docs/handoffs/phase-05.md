# Phase 5 Handoff

## Phase

Phase 5 — GTFS Studio draft/publish model

## Status

- Complete.
- Active phase after this handoff: Phase 6 — Trip Updates and Alerts architecture.

## What Was Implemented

- Added typed GTFS Studio draft storage for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Added `internal/gtfs.DraftService` for blank draft creation, active-feed cloning, typed draft upsert/remove operations, draft listing, soft discard, and draft publish.
- Added cloned-draft provenance through `gtfs_draft.base_feed_version_id`.
- Added explicit draft traceability for latest publish attempt, latest published feed version, and soft discard metadata.
- Added `gtfs_draft_publish` attempts linked to `validation_report`.
- Added direct shared feed-version publishing used by both GTFS ZIP import and Studio publish. Studio publish converts typed draft rows into the internal GTFS feed model and does not generate or re-import a synthetic ZIP.
- Added minimal server-rendered `cmd/gtfs-studio` admin UI with draft list/create/discard/publish, draft summary version visibility, and operational row forms for all Phase 5 core entities.
- Made discarded drafts hidden by default and visible only through an explicit filter.
- Made discarded and published drafts read-only by default.
- Made non-editable draft statuses reject before draft-to-feed conversion, validation, or shared publish activation.
- Made entity remove operations affect only rows in the current editable draft, never published GTFS rows or publish history.
- Preserved Phase 3 Vehicle Positions behavior and Phase 4 GTFS ZIP import behavior.

## What Was Designed But Intentionally Not Implemented Yet

- Canonical MobilityData GTFS Validator execution.
- Rich map editing for shape points.
- Timetable designer behavior for stop_times.
- Static GTFS public ZIP serving/export.
- Robust auth and role enforcement for Studio.
- Trip Updates.
- Alerts.
- Compliance dashboard and consumer ingestion workflows.
- Rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Schema And Interface Changes

- Added typed `gtfs_draft_*` tables for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Expanded `gtfs_draft` with `last_published_feed_version_id`, `last_publish_attempt_id`, `discarded_at`, `discarded_by`, and `discard_reason`.
- Added `gtfs_draft_publish`.
- Added `validation_report.gtfs_draft_publish_id`.
- Added `internal/gtfs.DraftService`, draft entity structs, draft publish result/error types, and draft CRUD/publish methods.
- Added `cmd/gtfs-studio`.
- Added `make run-gtfs-studio` and `task run:gtfs-studio`.

## Dependency Changes

- No new external dependencies were added.
- GTFS Studio uses Go standard library `net/http` and `html/template`.
- Canonical GTFS and GTFS-Realtime validators remain documented but unwired.

## Migrations Added

- `db/migrations/000005_gtfs_studio_drafts.sql`.

## Tests Added And Results

- Added DB-backed GTFS Studio integration tests for:
  - blank draft creation when no active feed exists
  - explicit blank draft creation when an active feed exists
  - typed CRUD for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies
  - cloned draft source `feed_version` provenance
  - draft edits not mutating published GTFS rows before publish
  - active published GTFS behavior remaining stable before draft publish
  - draft publish creating a new active `gtfs_studio` feed version
  - draft agency metadata mapping into published agency metadata on publish
  - published drafts becoming read-only
  - discarded drafts becoming read-only and not publishable
  - discarded publish rejection before publish attempt/shared publisher invocation
  - draft publish traceability through validation reports, publish attempts, and draft metadata
- Added `cmd/gtfs-studio` handler tests for:
  - default draft list hiding discarded drafts
  - explicit discarded filter including discarded drafts with status labels
  - draft summary showing status, base feed version, latest publish attempt, and published feed version

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
| `make migrate-up` | Passed | Applied `000005_gtfs_studio_drafts.sql`. |
| `make migrate-status` | Passed | Reports migrations 1 through 5 applied. |
| migration down/up smoke for `000005` | Passed | `make migrate-down`, `make migrate-up`, `make migrate-status`. |
| `make test-integration` | Passed | DB-backed telemetry, matcher, Vehicle Positions, GTFS import, and GTFS Studio tests passed. |
| `make validate` | Passed | Phase 5 file smoke only; canonical validators remain unwired. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Task equivalents were not run; Makefile remains independently usable and Task has been optional in prior phases.

## Known Issues

- `cmd/gtfs-studio` has no production auth layer; actor identity is accepted as form input/defaults to `system`.
- Studio UI is intentionally basic row editing. Shape points are not map-edited and stop_times are not managed through a timetable designer.
- Canonical GTFS validation is still not wired; internal validation is not a compliance claim.
- Original uploaded ZIP bytes are still not stored in Postgres and no public static GTFS ZIP endpoint exists.
- Existing current vehicle assignments are not cleared on schedule feed switch; the next matcher pass uses the new active feed, matching the Phase 4 posture.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-05.md`
  - `docs/phase-plan.md`
  - `docs/codex-task.md`
  - `docs/requirements-trip-updates.md`
  - `docs/requirements-calitp-compliance.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - `internal/feed/`
  - `internal/state/`
  - `internal/gtfs/`
  - `cmd/*` for minimal Trip Updates/Alerts architecture entrypoints if needed
  - `db/migrations/` only if Phase 6 needs adapter/diagnostics persistence
  - `docs/current-status.md`
  - `docs/handoffs/phase-06.md`
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
  - Trip Updates and Alerts are not implemented yet.
  - Studio auth is intentionally not production-grade.
- Recommended first implementation slice:
  - Begin Phase 6 by defining the Trip Updates prediction adapter contract and minimal diagnostics/no-op behavior.
  - Keep Trip Updates pluggable and do not couple predictor internals into telemetry ingest, Vehicle Positions, or GTFS Studio.
  - Preserve Phase 3 Vehicle Positions, Phase 4 GTFS import, and Phase 5 GTFS Studio behavior.
