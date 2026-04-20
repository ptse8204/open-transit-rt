# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 2 — Deterministic trip matching

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 1 closure-polish pass completed; no remaining Phase 1 cleanup items are known.
- Phase 2 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/phase-plan.md`
3. `docs/current-status.md`
4. `docs/handoffs/phase-01.md`
5. `docs/requirements-2a-2f.md`
6. `docs/requirements-trip-updates.md`
7. `docs/dependencies.md`
8. `docs/decisions.md`
9. `docs/codex-task.md`

## Current Objective

Begin Phase 2 deterministic trip matching using persisted telemetry from Phase 1. Do not start protobuf Vehicle Positions, GTFS import, GTFS Studio, Trip Updates, or Alerts yet.

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

- Task is not installed, but Task is optional and Makefile is independently usable.
- Docker must be running before DB-backed checks.
- GTFS import is not implemented yet, so Phase 2 should use fixtures or narrow schedule-query test doubles rather than starting Phase 4.

## First Files Likely To Edit

- `internal/state/`
- `internal/telemetry/`
- a narrow schedule-query package under `internal/gtfs/` or equivalent
- `db/migrations/`
- `testdata/`
- `docs/current-status.md`
- `docs/handoffs/phase-02.md`
- `docs/handoffs/latest.md`

## Phase 2 Entry Recommendation

Start deterministic matching without changing feed publication behavior:

1. Define the minimal GTFS schedule query boundary needed by the matcher.
2. Add agency-local service-day resolution.
3. Use latest accepted telemetry from the Phase 1 repository.
4. Preserve conservative `unknown` behavior for missing, ambiguous, or low-confidence matches.
5. Keep manual override precedence in the data model and tests, but do not build the full operator UI yet.
6. Add tests using after-midnight, frequency-based, stale, unmatched, matched, swapped, and block-transition fixtures.
7. Do not implement protobuf Vehicle Positions until Phase 3.

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

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
