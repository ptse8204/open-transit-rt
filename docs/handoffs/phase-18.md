# Phase Handoff

## Phase

Phase 18 — Admin UX And Agency Operations Console

## Status

- Complete for the approved minimal operator-console scope.
- Active phase after this handoff: Phase 19 — Realtime Quality And ETA Improvement.

## What Was Implemented

- Added authenticated server-rendered Operations Console pages under `/admin/operations`.
- Added dashboard, feeds/validation, telemetry freshness, device rotate/rebind, consumer evidence, evidence links, and setup checklist views.
- Added environment labeling from publication/environment configuration and last-updated timestamps for feeds, telemetry, scorecard, consumers, and evidence links where available.
- Added safe telemetry/assignment diagnostics: assignment state, route/trip IDs, confidence, reason codes, degraded state, assignment source, observed/received/assignment timestamps, and stale status.
- Added safe device binding listing through `internal/devices` without token hashes or token values.
- Added `/admin/alerts/console` for simple alert listing, create/update, publish, archive, and supported route/trip/stop informed-entity fields.
- Added navigation from Operations Console to GTFS Studio, and from GTFS Studio back to Operations Console.
- Updated the local app proxy and output so `/admin/operations` is reachable and printed in local-demo packaging.

## What Was Designed But Intentionally Not Implemented Yet

- No new first-time device credential API was added. The browser surface uses the existing rotate/rebind behavior and documents that limitation.
- No consumer submission API, portal automation, or file-backed evidence parser was added. The console prefers DB `consumer_ingestion` records and links to the docs tracker for targets not present in DB.
- No full setup wizard was added; Phase 18 uses a guided checklist.
- No manual override operator UI was added.
- No hosted login/SSO, production public-edge admin exposure, Prometheus/Grafana dashboards, or OpenTelemetry wiring was added.
- No new hosted evidence packet was collected.

## Dashboard/Pages/Routes Added

- `cmd/agency-config`:
  - `/admin/operations`
  - `/admin/operations/feeds`
  - `/admin/operations/telemetry`
  - `/admin/operations/devices`
  - `/admin/operations/consumers`
  - `/admin/operations/evidence`
  - `/admin/operations/setup`
- `cmd/feed-alerts`:
  - `/admin/alerts/console`
  - `/admin/alerts/console/{id}/publish`
  - `/admin/alerts/console/{id}/archive`

## Auth/CSRF/Security Behavior

- All new console routes require existing admin auth middleware.
- Read-only console pages allow `read_only`, `operator`, `editor`, and `admin`.
- Device rotate/rebind requires `admin`.
- Alert create/update, publish, and archive require `operator` or `admin`.
- Cookie-authenticated unsafe form posts require CSRF through the existing middleware.
- One-time device tokens are displayed only in the immediate POST response from the existing rebind behavior.
- The console does not render raw telemetry payload JSON, full assignment score details, token hashes, device-token pepper, bearer tokens, JWT/CSRF secrets, DB URLs, private URLs, or private debug blobs.
- `deploy/oci/Caddyfile` was not changed; production public edge remains public feed paths only.

## Schema And Interface Changes

- No database migrations were added.
- `internal/devices.Store` now includes `ListBindings(ctx, agencyID)` for safe device/vehicle binding reads.
- Existing JSON admin APIs remain backward-compatible.
- Public feed URLs and GTFS-RT protobuf contracts were not changed.

## Dependency Changes

- No new external dependency was added.
- UI remains server-rendered Go HTML using the standard library.

## Migrations Added

- None.

## Tests Added And Results

- Added handler coverage for Operations Console auth rejection, empty states, demo-like rendering, safe telemetry diagnostics, device token one-time display, admin-only device mutation, cookie CSRF rejection, and consumer status truthfulness.
- Added handler coverage for Alerts Console empty state, unauthenticated rejection, create/update, publish/archive, and role boundaries.

## Checks Run And Blocked Checks

Pre-edit baseline:

- `make validate`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make demo-agency-flow`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `docker compose -f deploy/docker-compose.yml --profile app config`: passed.
- `git diff --check`: passed.

Post-edit targeted checks:

- `go test ./cmd/agency-config ./cmd/feed-alerts ./cmd/gtfs-studio ./cmd/telemetry-ingest ./internal/devices`: passed.

Post-edit full checks:

- `make validate`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `docker compose -f deploy/docker-compose.yml --profile app config`: passed.
- `make smoke`: passed.
- `git diff --check`: passed.
- `make demo-agency-flow`: passed, including protected `/admin/operations` rejection and authenticated render checks.
- `make agency-app-up`: passed and printed `http://localhost:8080/admin/operations`.
- `make agency-app-down`: passed.

Blocked commands:

- None so far.

## Known Issues

- Device credential first-time setup is still represented through the existing rotate/rebind API behavior, not a separately named creation flow.
- Consumer evidence docs are linked from the console but not parsed into runtime DB state.
- The console is intentionally minimal and does not replace all command-line workflows.
- Manual assignment override workflows still need a richer operator UI in a later phase.
- Production SLO dashboards and alerting remain outside this Phase 18 UI.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/phase-19-realtime-quality-eta-improvement.md`
  - `internal/prediction/model.go`
  - `internal/prediction/deterministic.go`
- First files likely to edit:
  - `internal/prediction/`
  - `internal/state/`
  - `internal/feed/tripupdates/`
  - `cmd/feed-trip-updates/`
  - `cmd/agency-config/` if adding quality metrics to the Operations Console
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make smoke`
  - `git diff --check`
- Known blockers:
  - No consumer acceptance evidence exists in the repo.
  - Old ignored `.cache` secrets from Phase 15 must not be reused without rotation/revocation.
- Recommended first implementation slice:
  - Start Phase 19 with replay/evaluation and Trip Updates quality diagnostics behind the existing prediction adapter boundary, then surface only safe summary metrics in the Operations Console.
