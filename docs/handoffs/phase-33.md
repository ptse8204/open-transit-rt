# Phase 33 Handoff

## Phase

Phase 33 — Public GTFS Local/Pilot Evidence

## Status

- Complete as Outcome B — attempted public-GTFS run blocked.
- Active phase after this handoff: retry Outcome C evidence collection after
  the post-Phase-33 import fix.

Phase 33 added the public GTFS local/pilot evidence docs and templates, then
attempted the preferred LA Metro Bus GTFS local run. The core run blocked during
repo-supported import/publish because the large LA Metro dataset exceeded the
current importer context while inserting `stop_times.txt`.

Do not call Phase 33 evidence completed. This is an attempted-run blocked
closure only.

Post-Phase-33 implementation note: the large-import blocker was later addressed
in code by adding configurable import timeout handling, bulk `CopyFrom` loading
for large GTFS tables, and fresh-context failure reporting. A local post-fix
verification import published LA Metro Bus GTFS as `gtfs-import-26`, but the
full Outcome C evidence packet was not collected. Phase 33 remains Outcome B
until a future run records published schedule proof, five-path fetch proof, and
claim-boundary summaries.

## What Was Implemented

- Added `docs/phase-33-public-gtfs-local-pilot-evidence.md`.
- Added template-only evidence files under
  `docs/evidence/captured/public-gtfs-local-pilot/templates/`.
- Added a dated Outcome B packet at
  `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/`.
- Updated status/navigation docs for Phase 33 and the Outcome B blocker.

## What Was Designed But Intentionally Not Implemented Yet

- Outcome C evidence was not completed because import/publish blocked.
- No public LA Metro local feed fetch evidence was collected after the import
  blocker.
- No validator evidence was collected for the LA Metro local run.
- No telemetry simulator or dry-run evidence was collected for the LA Metro
  local run.
- No admin/private/debug boundary evidence was collected for the LA Metro local
  run.

## Attempted Run Summary

Catalog facts were checked on `2026-05-06T20:35:48Z`.

- Mobility Database was used as the primary catalog reference for official-feed
  status, route count, service range, producer URL, and dataset size.
- Transitland was used as the secondary catalog reference for current URL, last
  fetch, license URL, and attribution fields.
- Raw LA Metro Bus GTFS was downloaded to ignored `.cache/` storage and not
  committed.
- Source ZIP SHA-256:
  `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94`.
- First import attempt with `demo-agency` failed validation because
  `agency.txt` contains `agency_id=LACMTA`.
- Second import attempt with local `LACMTA` setup reached publish but blocked:
  `context deadline exceeded` while inserting `stop_times.txt`.

The blocked `LACMTA` import parsed these public-safe counts:

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

## Schema And Interface Changes

- None.

## Dependency Changes

- None.

The Phase 33 attempt used existing dependencies only: Docker Compose, Postgres,
the repo importer, `curl`, `unzip`, and `shasum`.

## Migrations Added

- None.

## Tests Added And Results

- No code tests were added because this phase changed documentation/evidence
  files only.
- The attempted LA Metro run exposed an importer/runtime blocker, not a new
  committed test case.

## Checks Run And Blocked Checks

- Planning-pass baseline reportedly passed before implementation:
  `make validate`, `make test`, and `git diff --check`.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `git diff --check` — passed.

Blocked for Outcome C:

- LA Metro import/publish blocked with `context deadline exceeded` while
  inserting `stop_times.txt`.
- Published LA Metro `/public/gtfs/schedule.zip` fetch was not run.
- Proof that fetched schedule is imported public GTFS was not run.
- Five-path fetch was not run.
- Validators were not run for the LA Metro local run.
- Telemetry simulator/dry-run and admin/private/debug checks were not run for
  the LA Metro local run.

## Known Issues

- The repo importer had a fixed command context that was too short for the
  current LA Metro Bus GTFS dataset on this local environment. This was fixed
  after Phase 33 by making the import timeout configurable and increasing the
  default.
- The failed `LACMTA` import left the `gtfs_import` row at `started` because
  the context expired before the failure update completed. This was fixed after
  Phase 33 by recording publish failures with a fresh short context.
- The local Compose app defaults to `AGENCY_ID=demo-agency`; a future Outcome C
  run for public GTFS with a different `agency_id` needs a documented local
  agency setup and service configuration path.

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
real-world ETA accuracy, or production-grade ETA quality from Phase 33.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/phase-33-public-gtfs-local-pilot-evidence.md`
  - `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md`
  - `cmd/gtfs-import/main.go`
  - `internal/gtfs/importer.go`
- First files likely to edit:
  - `cmd/gtfs-import/main.go`
  - `internal/gtfs/importer.go`
  - importer tests under `internal/gtfs/`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `git diff --check`
- Known blockers:
  - Phase 33 still lacks Outcome C evidence collection after the import fix;
  - local app agency configuration defaults to `demo-agency`.
- Recommended first implementation slice:
  - retry Phase 33 Outcome C with a dated packet that proves the published
    schedule is the imported public GTFS and records five-path fetch summaries,
    validators or blockers, telemetry/dry-run status, and admin/private boundary
    status.
