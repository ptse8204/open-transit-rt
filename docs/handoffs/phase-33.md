# Phase 33 Handoff

## Phase

Phase 33 — Public GTFS Local/Pilot Evidence

## Status

- Complete as Outcome C — public-GTFS local/pilot run completed with
  public-safe retained summaries.
- Active phase after this handoff: no new phase selected.

Phase 33 added the public GTFS local/pilot evidence docs and templates,
attempted the preferred LA Metro Bus GTFS local run, fixed the large-import
timeout exposed by that run, and retried the Outcome C evidence collection.

Do not extend Phase 33 beyond its evidence boundary. It proves local/pilot
handling of a real public static GTFS dataset only.

## What Was Implemented

- Added `docs/phase-33-public-gtfs-local-pilot-evidence.md`.
- Added template-only evidence files under
  `docs/evidence/captured/public-gtfs-local-pilot/templates/`.
- Added a dated Outcome C packet at
  `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/`.
- Updated status/navigation docs for Phase 33 and the Outcome C retry.
- Fixed the large public-GTFS import blocker:
  - `cmd/gtfs-import` now supports configurable timeout handling through
    `-timeout` and `GTFS_IMPORT_TIMEOUT`.
  - Large `gtfs_stop_time` and `gtfs_shape_point` publish inserts use
    `pgx.CopyFrom`.
  - Publish-failure reporting uses a fresh short context.

## Outcome C Evidence Summary

Catalog facts were checked on `2026-05-06T21:15Z`.

- Mobility Database was used as the primary catalog reference for official-feed
  status, route count, service range, producer URL, and dataset size.
- Transitland was used as the secondary catalog reference for current URL, last
  fetch, license URL, and attribution fields.
- Raw LA Metro Bus GTFS was downloaded to ignored `.cache/` storage and not
  committed.
- Source ZIP SHA-256:
  `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94`.
- Import through `cmd/gtfs-import -agency-id LACMTA -timeout 15m` published
  local feed version `gtfs-import-1`.
- Local public root: `http://localhost:19080`.
- Fetched schedule ZIP SHA-256:
  `1819fade012ca53a58d880285bb3ab85a0fce0a1b241d20cf15320e8542503ab`.
- The fetched schedule was unzipped in ignored `.cache/` storage and verified
  as the imported LA Metro public GTFS rather than the repo sample feed.

Imported counts:

| Entity | Count |
| --- | ---: |
| `agency` | 1 |
| `routes` | 114 |
| `stops` | 11884 |
| `trips` | 33642 |
| `stop_times` | 2105503 |
| `shapes` | 343530 |
| `calendar` | 130 |
| `calendar_dates` | 8432 |
| `frequencies` | 0 |

Fetched schedule proof:

- `agency.txt`: `agency_id=LACMTA`, `agency_name=Metro - Los Angeles`,
  timezone `America/Los_Angeles`, URL `https://www.metro.net`.
- `routes.txt`: 114 bus routes.
- Service-date coverage: `20251208` through `20270401`.

Five public paths fetched from the local root:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

Validator results:

- Static GTFS validator attempted but failed to execute because Java runtime
  was unavailable in this local environment.
- Vehicle Positions, Trip Updates, and Alerts GTFS-RT validators passed with 0
  errors, 0 warnings, and 0 info notices against empty valid protobuf feeds.

Additional checks:

- `scripts/device-onboarding.sh simulate --dry-run` printed synthetic payloads
  and sent no telemetry.
- Public root admin/debug paths returned `404`; direct local admin/debug paths
  without auth returned `401`; direct local admin operations with a runtime
  local admin token returned `200`.

## Schema And Interface Changes

- No migrations.
- `cmd/gtfs-import` CLI gained `-timeout`.
- `GTFS_IMPORT_TIMEOUT` can configure the import timeout; `-timeout 0` disables
  it.

## Tests Added And Results

Added focused import timeout and failure-report tests:

- `cmd/gtfs-import/main_test.go`
- `internal/gtfs/importer_test.go`

## Checks Run

- Focused `go test ./cmd/gtfs-import ./internal/gtfs` — passed.
- `make validate` — passed.
- `make test` — passed.
- First `make test-integration` attempt — blocked because Postgres was not
  ready immediately after startup.
- `make test-integration` rerun after `pg_isready` — passed.
- `git diff --check` — passed.
- Final post-Outcome-C-docs `make validate` — passed.
- Final post-Outcome-C-docs `make test` — passed.
- Final post-Outcome-C-docs `git diff --check` — passed.

## Known Issues

- Static GTFS validation could not execute in this local environment because
  Java was unavailable. The schedule validator record is an execution blocker,
  not a data-quality pass.
- The local Compose app defaults to `AGENCY_ID=demo-agency`; the Outcome C run
  used repo service binaries with `AGENCY_ID=LACMTA` and a local public proxy
  rather than `make agency-app-up`, because `make agency-app-up` imports the
  repo sample feed.
- GTFS-RT endpoints were valid empty protobuf publications. They are not real
  LA Metro realtime data.

## Truthfulness And Evidence Boundary

- Final-root blocker unchanged.
- Consumer statuses unchanged.
- No consumer target records changed.
- No target-specific consumer artifacts added.
- No final-root evidence created.
- No OCI pilot final-root wording changed.

Do not claim agency adoption, agency endorsement, agency approval, official
agency feed status, agency-owned final-root proof, consumer submission/review/
acceptance, consumer ingestion/listing/display, Caltrans/CAL-ITP compliance,
hosted SaaS availability, production readiness, real vendor AVL compatibility,
real-world ETA accuracy, production-grade ETA quality, or real LA Metro
realtime data from Phase 33.

## Exact Next-Step Recommendation

The next useful retained-evidence targets are:

- agency-owned/final-root proof if an agency-owned or agency-approved final
  public root becomes available;
- authorized target-specific consumer submission evidence if an operator
  selects a target and retains target-originated artifacts;
- real agency pilot evidence;
- real deployment operations evidence.
