# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 3 — Vehicle Positions production feed

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-02.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/dependencies.md`
9. `docs/decisions.md`

## Current Objective

Begin Phase 3 Vehicle Positions production feed using persisted latest telemetry and persisted Phase 2 assignments. Do not start GTFS import, GTFS Studio, Trip Updates, or Alerts.

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
- GTFS import is not implemented yet, so Phase 3 should use existing published GTFS tables and test fixtures rather than starting Phase 4.
- Canonical GTFS and GTFS-RT validators are documented but not wired yet.

## First Files Likely To Edit

- `cmd/feed-vehicle-positions/`
- `internal/feed/`
- `internal/state/`
- `internal/telemetry/` only for narrow read/query additions if existing latest-telemetry methods are insufficient
- `db/migrations/` only if Phase 3 needs feed-publication metadata changes
- `testdata/expected/`
- `docs/current-status.md`
- `docs/handoffs/phase-03.md`
- `docs/handoffs/latest.md`

## Phase 3 Entry Recommendation

Start Vehicle Positions production feed without changing Trip Updates or GTFS import behavior:

1. Replace placeholder sample data in `cmd/feed-vehicle-positions` with DB-backed latest accepted telemetry.
2. Read persisted current assignments from `internal/state`.
3. Generate valid GTFS-RT Vehicle Positions protobuf output at a stable public endpoint.
4. Keep a JSON debug endpoint for inspection.
5. Preserve Phase 2 conservative behavior: unmatched and stale vehicles must not emit false trip certainty.
6. Add tests for matched, unknown, stale, and assignment-aware Vehicle Positions output.
7. Do not implement Trip Updates, Alerts, GTFS import, or GTFS Studio.

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

## Phase 2 Notes For Phase 3

- `vehicle_trip_assignment.score_details_json` is loose debug JSON, not a stable integration schema.
- Unknown assignment rows are explicit and close previous active rows.
- Unknown rows carry `service_date` whenever agency timezone and observed timestamp are resolvable; the column is nullable for truly unresolved cases.
- Repeated trip instances with the same `trip_id` but different `start_time` must remain distinct.
- `internal/state.Engine` is the only valid production matcher entry point; `NewEngine` returns an error if schedule or assignment repositories are missing, and `MustNewEngine` is reserved for test/bootstrap callers.
- Active manual overrides are absolute and are evaluated before stale-telemetry fallback.
- Continuity and block-transition scoring are time-aware and require configured-window plausibility, not just same trip or same block identity.
- Block-transition scoring also verifies nearest plausible next-trip sequencing within the block when start-time identity is available; later same-block trips do not receive credit solely because they are later.
- Numeric explicit `bearing: 0` is valid true north and can receive movement-direction credit only when the stored payload explicitly contains a numeric `bearing` field; null, malformed, or payload-missing zero values are treated as missing.
- `shape_dist_traveled = 0` is preserved as a valid persisted value.
- Repeated identical degraded unknown states reuse the active degraded assignment only when degraded state, reason codes, service date, and telemetry evidence match. The implementation compares `telemetry_event_id` when present and falls back to exact `active_from` equality only when both rows lack telemetry evidence; materially new evidence or service-day changes replace the unknown row and keep prior confident rows closed.
- Manual override assignments populate active feed and block context when resolvable, so they are not thinner persisted rows than automatic matches.
- Missing shape data uses reason `missing_shape` plus degraded state `missing_shape`; it reduces confidence but does not automatically block a match when other strong evidence exists.
- Non-exact frequency matches use conservative window identity details and must not be treated as exact scheduled instances.
- `no_schedule_candidates` is reserved for successful schedule queries that return no trips. Repository/config/resolution failures use distinct matcher-system-failure reasons.
- Route-hint matching is reserved for future input expansion and is not active in Phase 2 because telemetry does not carry a route hint.
- Service-day resolution checks the observed agency-local date and immediately previous local date only; do not assume broader multi-day post-midnight coverage without extending the resolver.
- The Phase 2 handoff matches the actual implementation after the semantic-closure pass.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
