# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 8 — Compliance and consumer workflow

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 GTFS Studio draft/publish model is implemented and complete.
- Phase 6 Trip Updates and Alerts architecture is implemented and complete.
- Phase 7 prediction quality and operations workflows are implemented and complete.
- Phase 8 is ready to start.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-07.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin Phase 8 compliance and consumer workflow. Preserve the Phase 7 prediction adapter, metrics, override, review queue, and cancellation-linkage boundaries. Do not bypass canonical feed boundaries, do not couple validators into core prediction or matching internals, and do not start rider apps, payments, passenger accounts, dispatcher CAD, or marketplace workflows.

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
- Canonical GTFS and GTFS-Realtime validators are documented but not wired yet.
- Alerts authoring and persistence remain deferred.
- Phase 7 predictions are conservative schedule-deviation predictions, not production-grade learned ETAs.
- Prediction review workflow has repository behavior but no full operator UI.
- GTFS Studio auth is minimal/dev-only and not production-grade.

## First Files Likely To Edit

- `internal/gtfs/`
- `internal/feed/tripupdates/`
- `internal/feed/alerts/`
- `internal/prediction/` only if compliance metrics need prediction diagnostics
- `db/migrations/` if validation/report/compliance persistence changes are required
- `docs/current-status.md`
- `docs/handoffs/phase-08.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md`
- `docs/decisions.md` if architecture-significant compliance decisions are made

## Phase 8 Entry Recommendation

Start compliance work behind stable validation and publication boundaries:

1. Inspect Phase 7 Trip Updates metrics, review queue, cancellation linkage, and deterministic adapter behavior.
2. Choose the first compliance slice, likely canonical validation/reporting and public feed metadata.
3. Wire GTFS and GTFS-RT validator execution behind documented validation interfaces.
4. Add public discoverability metadata and stable feed status reporting.
5. Keep compliance workflow separate from prediction internals.
6. Preserve Vehicle Positions, telemetry ingest, GTFS import, GTFS Studio, and Trip Updates endpoint behavior.

## Constraints To Preserve

- Mostly Go.
- Postgres/PostGIS source of truth.
- Vehicle Positions first.
- Trip Updates pluggable.
- Draft GTFS separate from published GTFS.
- Conservative matching and prediction.
- Manual overrides take precedence over matching.
- No rider apps, payments, passenger accounts, or dispatcher CAD.
- External integrations stay behind documented adapters.
- Runtime GTFS import input is ZIP; directory parsing is test-fixture setup only.
- GTFS Studio publishes typed draft rows through the shared validation/activation helper directly, not through synthetic ZIP import.
- GTFS times beyond `24:00:00` remain stored as imported text in canonical published GTFS tables.
- Trip Updates packages must not become dependencies of telemetry ingest, Vehicle Positions, or GTFS Studio.
- Canceled trips are excluded from ETA coverage denominator and tracked separately.
- Prediction review items use `open`, `resolved`, and `deferred` lifecycle states.

## Phase 7 Notes For Phase 8

- `internal/prediction.Adapter` remains the only Trip Updates prediction boundary.
- `prediction.DeterministicAdapter` is the default runtime adapter.
- `prediction.NoopAdapter` remains available through `TRIP_UPDATES_ADAPTER=noop`.
- Trip Updates diagnostics persist to `feed_health_snapshot` with Phase 6 fields plus prediction metrics and adapter details.
- Prediction review items persist as `incident` rows with `incident_type = 'prediction_review'`.
- Override create, replace, clear, and review status updates write audit rows.
- Canceled trips emit conservative `CANCELED` Trip Updates and persist missing-alert linkage signals until Alerts authoring/persistence is implemented.
- Added trips, short turns, detours, deadhead, layover, weak, stale, degraded, and ambiguous cases are withheld with explicit reasons.
- `cmd/feed-trip-updates` still exposes `/public/gtfsrt/trip_updates.pb` and `/public/gtfsrt/trip_updates.json`.
- `cmd/feed-alerts` still exposes valid empty Alerts endpoints; public Alerts authoring/persistence is not implemented.
- No external prediction dependency was added.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
