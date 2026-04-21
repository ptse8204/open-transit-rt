# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 5 — GTFS Studio draft/publish model

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-04.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin Phase 5 GTFS Studio draft/publish model. Do not start Trip Updates, Alerts, rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

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
- GTFS Studio runtime editing flows are not implemented yet.

## First Files Likely To Edit

- `internal/gtfs/`
- `db/migrations/`
- `cmd/*` only for minimal admin/Studio entrypoints if needed
- `testdata/gtfs/`
- `docs/current-status.md`
- `docs/handoffs/phase-05.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md`
- `docs/decisions.md` if architecture-significant draft/publish decisions are made

## Phase 5 Entry Recommendation

Start GTFS Studio draft/publish work without changing Vehicle Positions semantics or weakening the Phase 4 import pipeline:

1. Inspect the existing `gtfs_draft` / `gtfs_draft_record` scaffolding and the Phase 4 import service.
2. Add draft storage and editing boundaries that remain separate from published `feed_version` tables.
3. Add minimal draft CRUD for core GTFS entities needed by the first Studio slice.
4. Publish drafts through the same validation and activation semantics used by GTFS ZIP import.
5. Add tests proving draft data does not leak into active published GTFS until publish.
6. Do not implement Trip Updates, Alerts, rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows in Phase 5.

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
- GTFS times beyond `24:00:00` remain stored as imported text in canonical published GTFS tables.

## Phase 4 Notes For Phase 5

- `cmd/gtfs-import` is a thin CLI wrapper over `internal/gtfs.ImportService`.
- `gtfs_import.feed_version_id` is set only after successful publish and remains `NULL` for failed imports.
- Validation failures create no `feed_version`; publish failures roll back staged rows.
- `validation_report.gtfs_import_id` links schedule validation reports to import attempts.
- `block_id` from `trips.txt` is imported when present and visible through `gtfs.PostgresRepository.ListTripCandidates`.
- `gtfs_shape_line` is built from ordered shape points when a shape has at least two points.
- Optional `shapes.txt` and `frequencies.txt` are accepted.
- Canonical validator tooling remains unwired; internal validation is not a compliance claim.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
