# Phase 27 Handoff

## Phase

Phase 27 — Multi-Agency Isolation Prototype

## Status

- Complete for the approved test-and-documentation scope.
- Active phase after this handoff: Phase 28 — Production Operations Hardening.

## What Was Implemented

- Added synthetic multi-agency fixture notes and public-safe fixture metadata under `testdata/multi-agency/`.
- Added repository and handler tests proving selected agency boundaries for auth, publication metadata, feed discovery, validation, scorecards, consumers, devices, telemetry, state assignments/overrides, GTFS imports/feed versions/draft listings, Operations Console, GTFS Studio, Alerts, prediction operations, and audit-row writes.
- Added narrow protected-debug fixes for Vehicle Positions, Trip Updates, and Alerts JSON debug handlers so cross-agency admin principals cannot read another service instance's debug snapshot.
- Documented current public endpoint scope and agency-id persistence audit as current isolation review, not production multi-tenant certification.

## What Was Designed But Intentionally Not Implemented Yet

- No hosted multi-tenant service.
- No global admin model.
- No public feed URL or GTFS-RT protobuf contract changes.
- No consumer status changes, consumer submission APIs, or automated submissions.
- No per-agency routing for public schedule ZIP or GTFS-RT protobuf endpoints.
- No audit-log read endpoint.

## Schema And Interface Changes

- No database migrations.
- No public feed URL changes.
- No GTFS-RT protobuf contract changes.
- Protected JSON debug handlers now require the authenticated principal agency to match the generated service/snapshot agency.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- Fixtures/tests added:
  - `testdata/multi-agency/README.md`
  - `testdata/multi-agency/agencies.json`
  - handler tests in `cmd/agency-config`, `cmd/telemetry-ingest`, `cmd/feed-vehicle-positions`, `cmd/feed-trip-updates`, `cmd/feed-alerts`, and `cmd/gtfs-studio`
  - repository tests in `internal/auth`, `internal/compliance`, `internal/devices`, `internal/telemetry`, `internal/alerts`, and `internal/prediction`
- Auth-boundary tests:
  - DB-backed `PostgresRoleStore` scopes roles by claim agency, even when subject/email exists in multiple agencies.
  - Protected admin paths reject conflicting query/body/form `agency_id` and derive agency from the authenticated principal.
- Data-isolation tests:
  - publication metadata, feed discovery, validation status, scorecards, consumer records, device bindings, telemetry events, state current assignments, state manual overrides, GTFS imported feed versions, GTFS trip candidates, GTFS draft listings, prediction operations, and audit/incident rows are scoped by agency for selected workflows.
- Feed-discovery/public-feed behavior:
  - `/public/feeds.json` is query-routed by `agency_id`; omitted query uses configured `AGENCY_ID`.
  - `/public/gtfs/schedule.zip`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/alerts.pb` remain service-instance scoped by configured `AGENCY_ID`.
- Operations Console isolation behavior:
  - `/admin/operations` and feeds/telemetry/devices/consumers/evidence/setup sections reject conflicting `agency_id`.
  - telemetry and devices views are tested against mixed synthetic agency data.
  - setup publication and validation forms remain principal-derived and conflict-bounded.
- GTFS isolation behavior:
  - Repository-level tests cover active feed lookup by agency, trip candidate lookup by agency/feed version, import row agency IDs, and `DraftService.ListDrafts` agency scoping.
  - GTFS Studio handler tests cover list/create/draft summary/publish/discard/entity edit agency boundaries.
  - `DraftService.GetDraft(ctx, draftID)` remains an ID-only repository method that returns the draft agency for caller enforcement; repository-level denial for cross-agency draft ID lookup is deferred unless a future API adds an agency parameter.
- State isolation behavior:
  - Repository-level tests cover `CurrentAssignment`, `ListCurrentAssignments`, `SaveAssignment`, `ActiveManualOverride`, and incident rows written by assignment saves.
  - `internal/state` does not currently write `audit_log` rows directly; audit coverage for manual override workflows is in prediction operations tests.
- Alerts admin/console isolation behavior:
  - `/admin/alerts`, `/admin/alerts/console`, `/admin/alerts/{id}/publish`, `/admin/alerts/{id}/archive`, and `/admin/alerts/reconcile-cancellations` are covered for list/mutation/reconcile agency boundaries and body/form/query conflicts.
- Device/telemetry isolation behavior:
  - agency A device tokens cannot submit agency B agency/device/vehicle payloads.
  - telemetry debug listings use principal agency and reject conflicting query agency.
  - Postgres telemetry latest/events listings remain agency-scoped.

Focused tests run during implementation:

```bash
go test ./cmd/agency-config ./cmd/feed-alerts ./cmd/feed-trip-updates ./cmd/feed-vehicle-positions ./cmd/gtfs-studio ./cmd/telemetry-ingest
go test ./internal/auth ./internal/compliance ./internal/devices ./internal/telemetry ./internal/alerts ./internal/prediction
```

Both focused test commands passed.

## Checks Run And Blocked Checks

Baseline checks before editing:

```bash
make validate
make test
make test-integration
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

All baseline checks passed after starting the local database for integration tests.

Final checks after editing:

```bash
go test ./cmd/agency-config ./cmd/feed-alerts ./cmd/feed-trip-updates ./cmd/feed-vehicle-positions ./cmd/gtfs-studio ./cmd/telemetry-ingest
go test ./internal/auth ./internal/compliance ./internal/devices ./internal/telemetry ./internal/alerts ./internal/prediction
go test ./internal/state ./internal/gtfs
make validate
make test
make test-integration
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

All final checks passed. `make test-integration` required starting the local Postgres service with `make db-up` after the first attempt found `localhost:55432` unavailable; the rerun passed.

Blocked commands: none.

## Known Issues

- Phase 27 proves repository-level isolation for selected workflows, not production hosted multi-tenant readiness.
- Public schedule ZIP and GTFS-RT protobuf endpoints are still service-instance scoped by configured `AGENCY_ID`.
- Backup, restore, export, deletion, and evidence packet generation are not tenant-safe multi-agency workflows yet.
- Full GTFS Studio entity isolation coverage remains limited to current minimal handler semantics and selected repository tests. `DraftService.GetDraft` is ID-only and returns the agency for handler enforcement; repository-level cross-agency denial for direct draft ID lookup remains deferred.
- There is no global admin model; Phase 27 intentionally did not invent one.
- No audit-log read surface exists; tests inspect written rows directly where needed.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-27.md`
  - `docs/phase-28-production-operations-hardening.md`
  - `docs/multi-agency-strategy.md`
- First files likely to edit:
  - `docs/runbooks/small-agency-pilot-operations.md`
  - `docs/upgrade-and-rollback.md`
  - `scripts/pilot-ops.sh`
  - deployment/systemd docs if Phase 28 changes operations workflows
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make test-integration`
  - `make realtime-quality`
  - `make smoke`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Known blockers:
  - No agency-owned final feed root is available in repo evidence.
  - No production multi-agency public feed routing exists.
  - Consumer targets remain `prepared` only.
- Recommended first implementation slice:
  - Start Phase 28 with backup/restore/upgrade/incident-response hardening that preserves the Phase 27 agency-boundary language and avoids hosted multi-tenant claims.
