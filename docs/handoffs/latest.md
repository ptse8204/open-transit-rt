# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 8 — Compliance and consumer workflow is complete. No Phase 9 is defined in `docs/phase-plan.md`.

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 GTFS Studio draft/publish model is implemented and complete.
- Phase 6 Trip Updates and Alerts architecture is implemented and complete.
- Phase 7 prediction quality and operations workflows are implemented and complete.
- Phase 8 publication/compliance workflow is implemented and complete for the first production-directed layer.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-08.md`
4. `docs/phase-plan.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Begin a focused post-Phase-8 hardening slice. Preserve stable feed URLs, public schedule ZIP behavior, persisted Alerts, Phase 7 prediction boundaries, Vehicle Positions behavior, GTFS import, and GTFS Studio draft/publish behavior.

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
- Canonical validator command adapters exist, but exact validator distributions are not pinned or installed by repo automation.
- Production admin auth and role enforcement are not implemented.
- Consumer ingestion workflow records exist, but external consumer submission APIs are not integrated.

## First Files Likely To Edit

- `.env.example`
- `Makefile`
- `internal/compliance/`
- `cmd/agency-config/`
- deployment/CI docs or scripts for validator installation
- auth/admin docs if production auth is the next slice

## Phase 8 Notes For Future Work

- `/public/gtfs/schedule.zip` is generated on demand from active published GTFS tables.
- Schedule ZIP bytes are deterministic for unchanged active feed data; ZIP entry modified times and HTTP `Last-Modified` use the active feed revision time.
- Realtime `published_feed.revision_timestamp` is a publication/bootstrap metadata revision and must not change on every feed generation.
- Realtime freshness and generation health belong in `feed_health_snapshot`.
- `/public/feeds.json` reads per-feed data from `published_feed`; license/contact fields resolve from `feed_config` only when per-feed values are empty.
- `feed_config.publication_environment = 'production'` makes missing canonical validator execution red in scorecards. In `dev`, missing validators are yellow/not-run.
- Alerts authoring/persistence is owned by `internal/alerts`; GTFS-RT protobuf rendering is owned by `internal/feed/alerts`.
- Prediction packages must not import Alerts packages. Canceled-trip missing-alert review signals are satisfied by the Alerts-owned reconciler.
- Validator adapters parse structured JSON reports from stdout, stderr, or output files, store normalized `validation_report` rows with error/warning/info counts, and record missing validator configuration as `status='not_run'`.

## Constraints To Preserve

- Mostly Go.
- Postgres/PostGIS source of truth.
- Stable public URLs for schedule, Vehicle Positions, Trip Updates, and Alerts.
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
- Trip Updates packages must not become dependencies of telemetry ingest, Vehicle Positions, GTFS Studio, or Alerts.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
