# Phase 8 Handoff

## Phase

Phase 8 — Compliance and consumer workflow

## Status

- Complete for the first publication/compliance layer.
- Active phase after this handoff: no Phase 9 is defined; next work should be focused hardening.

## What Was Implemented

- Added persisted Service Alerts authoring and lifecycle state.
- Published real GTFS-RT Alerts from the existing stable `/public/gtfsrt/alerts.pb` and `/public/gtfsrt/alerts.json` endpoints.
- Added Alerts-owned canceled-trip reconciliation from active canceled-trip overrides and Phase 7 `expected_alert_missing=true` prediction-review incidents.
- Added stable public static GTFS schedule ZIP serving at `/public/gtfs/schedule.zip`.
- Added `/public/feeds.json` discoverability metadata with explicit license, contact, validation, health, revision, and readiness fields.
- Added publication metadata bootstrap that writes `feed_config`, `published_feed`, `consumer_ingestion`, and `marketplace_gap`.
- Added consumer ingestion workflow record updates.
- Added compliance scorecard snapshot persistence.
- Added canonical validator command adapters that store normalized `validation_report` rows and record `not_run` when tooling is missing.
- Expanded architecture tests so Alerts do not couple into prediction, Trip Updates, telemetry ingest, Vehicle Positions, or GTFS Studio.

## What Was Designed But Intentionally Not Implemented Yet

- Pinned canonical validator binary installation/download steps for local dev, CI, or production.
- Production auth and role enforcement around admin mutation endpoints.
- Rich operator UI for Alerts, consumer workflows, and scorecards.
- External consumer submission APIs; Phase 8 stores workflow metadata and packet JSON only.
- Materialized schedule ZIP caching; the schedule ZIP is generated on demand from active published GTFS tables.
- Marketplace/vendor packaging beyond tracked gap records.

## Schema And Interface Changes

- Added migration `000007_phase_8_alerts_compliance.sql`.
- Added `feed_config.publication_environment` with `dev` and `production`.
- Added `service_alert` and `service_alert_informed_entity`.
- Added `compliance_scorecard_snapshot`.
- Added `internal/alerts` for persisted alert authoring, audit logging, lifecycle, and cancellation reconciliation.
- Added `internal/feed/schedule` for deterministic public schedule ZIP generation.
- Replaced the Phase 6 deferred Alerts builder with DB-backed `internal/feed/alerts` protobuf/JSON rendering.
- Added `internal/compliance` for publication metadata, `/public/feeds.json`, consumer ingestion records, scorecards, and validator result storage.
- Added admin/public routes in `cmd/agency-config` and admin alert routes in `cmd/feed-alerts`.

`/public/feeds.json` uses RFC3339 UTC JSON timestamps. Per-feed metadata comes from `published_feed`; license/contact fields resolve from `feed_config` only when the per-feed field is empty. Realtime `published_feed.revision_timestamp` is a publication/bootstrap metadata timestamp and is not updated on every feed generation. Realtime freshness is recorded in `feed_health_snapshot`.

Schedule ZIP serving is on demand from the active published `feed_version`. ZIP entries use deterministic order and the active feed revision timestamp, preserving stable bytes and stable `Last-Modified` semantics for unchanged feed data.

## Dependency Changes

- No new Go module dependencies were added.
- Added runtime/developer configuration:
  - `SCHEDULE_FEED_URL`
  - `TECHNICAL_CONTACT_EMAIL`
  - `FEED_LICENSE_NAME`
  - `FEED_LICENSE_URL`
  - `PUBLICATION_ENVIRONMENT`
  - `GTFS_RT_VALIDATOR_COMMAND`
- `GTFS_VALIDATOR_PATH` is now used by the static GTFS validator command adapter when configured.
- `GTFS_RT_VALIDATOR_COMMAND` is now used by the GTFS-RT validator command-template adapter when configured.

## Migrations Added

- `db/migrations/000007_phase_8_alerts_compliance.sql`

## Tests Added And Results

- Added Alerts feed tests for non-empty GTFS-RT Alert protobuf output, translated text, active period, informed trip selectors, and debug JSON.
- Added `cmd/feed-alerts` handler tests for public endpoints, admin authoring, reconciliation, method rejection, and readiness.
- Added schedule ZIP tests proving deterministic bytes and entry `Modified` timestamp behavior.
- Added compliance tests for `/public/feeds.json` readiness semantics and dev-versus-production validator missing-tool score behavior.
- Expanded non-coupling tests around Alerts and prediction/Trip Updates boundaries.

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
| `make migrate-up` | Passed | Applied `000007_phase_8_alerts_compliance.sql`. |
| `make migrate-status` | Passed | Reports migrations 1 through 7 applied. |
| `make test-integration` | Passed | DB-backed integration tests passed where supported. |
| `make validate` | Passed | Phase 8 file smoke passed. |
| `git diff --check` | Passed | No whitespace errors. |

Blocked checks:
- No required checks were blocked.
- Canonical validator binaries/commands were not configured, so actual external validator execution remains a deployment/configuration step. The app records missing validator execution as `not_run`.

## Known Issues

- Validator command adapters exist, but validator distributions are not pinned or installed by repo automation.
- Production auth and role enforcement are still not implemented.
- Schedule ZIPs are generated on demand; this is correct and deterministic, but a later cache may be useful for large feeds.
- `/public/feeds.json` and scorecards depend on publication bootstrap having been run for each agency.
- Consumer ingestion workflow stores metadata and packets; it does not submit to Google, Apple, Transit App, Bing, or Moovit APIs.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/decisions.md`
  - `docs/dependencies.md`
- First files likely to edit:
  - `.env.example`
  - `Makefile`
  - `internal/compliance/`
  - `cmd/agency-config/`
  - deployment/CI docs or scripts for validator installation
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
  - Exact canonical GTFS and GTFS-RT validator distributions still need to be pinned.
  - Production admin auth remains out of scope so far.
- Recommended first implementation slice:
  - Pin and document canonical validator distributions for local dev and CI, then add production auth/role enforcement around Phase 8 admin mutation endpoints.
