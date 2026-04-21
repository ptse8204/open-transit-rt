# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 6 — Trip Updates and Alerts architecture

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 GTFS Studio draft/publish model is implemented and complete.
- Phase 6 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-05.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin Phase 6 Trip Updates and Alerts architecture. Do not start rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Exact First Commands

```bash
command -v go
go version
make fmt
make test
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
make test-integration
```

If Task is installed, optional equivalents may be run:

```bash
task fmt
task test
task migrate:status
task test:integration
```

## Known Blockers

- Task is optional and may not be installed; Makefile remains independently usable.
- Docker must be running before DB-backed checks.
- Canonical GTFS and GTFS-RT validators are documented but not wired yet.
- Trip Updates and Alerts are not implemented yet.
- GTFS Studio auth is minimal/dev-only and not production-grade.

## First Files Likely To Edit

- `internal/feed/`
- `internal/state/`
- `internal/gtfs/`
- `cmd/*` only for minimal Phase 6 entrypoints if needed
- `db/migrations/` only if adapter/diagnostics persistence is needed
- `docs/current-status.md`
- `docs/handoffs/phase-06.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md`
- `docs/decisions.md` if architecture-significant prediction/alerts decisions are made

## Phase 6 Entry Recommendation

Start Trip Updates and Alerts architecture without changing Vehicle Positions semantics, weakening GTFS import/Studio publish behavior, or coupling predictor internals into core state:

1. Inspect `internal/feed`, `internal/state`, and `internal/gtfs`.
2. Define a narrow prediction adapter boundary for Trip Updates inputs and outputs.
3. Add a documented no-op or minimal adapter plus diagnostics plumbing.
4. Add Alerts feed model and stable endpoint shape only within Phase 6 scope.
5. Keep Vehicle Positions, telemetry ingest, matching, GTFS import, and GTFS Studio behavior stable.
6. Do not implement rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

## Constraints To Preserve

- Mostly Go.
- Postgres/PostGIS source of truth.
- Vehicle Positions first.
- Trip Updates pluggable.
- Draft GTFS separate from published GTFS.
- Conservative matching.
- Manual overrides take precedence over matching.
- No rider apps, payments, passenger accounts, or dispatcher CAD.
- External integrations stay behind documented adapters.
- Runtime GTFS import input is ZIP; directory parsing is test-fixture setup only.
- GTFS Studio publishes typed draft rows through the shared validation/activation helper directly, not through synthetic ZIP import.
- GTFS times beyond `24:00:00` remain stored as imported text in canonical published GTFS tables.

## Phase 5 Notes For Phase 6

- `cmd/gtfs-studio` is a minimal server-rendered admin UI.
- `internal/gtfs.DraftService` owns draft creation, editing, discard, and publish.
- Typed draft tables cover agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- `gtfs_draft_record` remains unused legacy scaffold.
- Drafts cloned from active published GTFS capture `base_feed_version_id`.
- Drafts in `published` or `discarded` status are read-only by default.
- Entity remove operations only affect rows in the current editable draft and never delete published GTFS rows or publish history.
- Discarded drafts are hidden from the default Studio list and included only with an explicit filter.
- Studio publish uses `feed_version.source_type = 'gtfs_studio'`.
- `gtfs_draft_publish` and `validation_report.gtfs_draft_publish_id` preserve publish traceability.
- Non-editable draft statuses are rejected before draft-to-feed conversion, validation, or shared publish activation.
- Canonical validator tooling remains unwired; internal validation is not a compliance claim.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
