# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 10 — Docs, Tutorials, Deployment, and Demo is complete for the current codebase surface. Continue with Phase 11 from `docs/phase-plan-production-closure.md`.

## Phase Status

- Phase 0 scaffolding is implemented and operationally closed.
- Phase 1 durable telemetry foundation is implemented and operationally closed.
- Phase 2 deterministic trip matching is implemented and semantically closed.
- Phase 3 Vehicle Positions production feed is implemented and complete.
- Phase 4 GTFS import and publish pipeline is implemented and complete.
- Phase 5 GTFS Studio draft/publish model is implemented and complete.
- Phase 6 Trip Updates and Alerts architecture is implemented and complete.
- Phase 7 prediction quality and operations workflows are implemented and complete for the first conservative scope.
- Phase 8 publication/compliance workflow is implemented and complete for the first production-directed layer.
- Phase 9 production closure is implemented for validator execution, validator tooling pins, admin auth/roles, device auth/binding, assignment current-row races, safer config defaults, debug endpoint protection, request logging/request IDs, metrics toggle, stronger feed-service readiness, and smoke coverage.
- Phase 10 docs/tutorial/deployment/demo work is implemented and complete for the current repository surface.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-10.md`
4. `docs/phase-plan-production-closure.md`
5. `docs/prompts/calitp-truthfulness.md`
6. `docs/codex-task.md`
7. `docs/requirements-2a-2f.md`
8. `docs/requirements-trip-updates.md`
9. `docs/requirements-calitp-compliance.md`
10. `docs/dependencies.md`
11. `docs/decisions.md`
12. `docs/tutorials/calitp-readiness-checklist.md`

## Current Objective

Start Phase 11 — Compliance Evidence and Optional External Integrations. Produce truthful evidence mapping before making stronger readiness or compliance claims. Keep optional predictors behind `internal/prediction.Adapter` if any are implemented.

## Exact First Commands

```bash
command -v go
go version
make validators-install
make validators-check
make test
make smoke
make demo-agency-flow
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
make test-integration
```

If Task is installed, optional equivalents may be run:

```bash
task test
task smoke
task demo:agency
task migrate:status
task test:integration
```

## Known Blockers

- Task is optional and may not be installed; Makefile remains independently usable.
- Docker must be running before DB-backed checks, the GTFS-RT validator wrapper, and the agency demo.
- `scripts/demo-agency-flow.sh` uses local ports `8081` through `8086` plus `8090`; free those ports before running it.
- Full hosted login/SSO and server-side `jti` replay tracking are deferred.
- Consumer-ingestion workflow records exist, but external consumer submission APIs are not integrated.
- Compliance or consumer-acceptance claims need deployment and third-party evidence.

## First Files Likely To Edit

- a new Phase 11 evidence checklist doc
- `docs/dependencies.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-11.md`

## Phase 10 Notes For Future Work

- `README.md` now documents current public/protected endpoints, quickstart, deployment path, limitations, and truthful CAL-ITP/Caltrans-aligned wording.
- Tutorials live under `docs/tutorials/`.
- `make demo-agency-flow` is the executable agency demo. It verifies public `schedule.zip`, `feeds.json`, Vehicle Positions / Trip Updates / Alerts protobufs, protected debug/admin routes, protected GTFS Studio root and draft subroutes, validation flow, scorecard, and consumer-ingestion visibility.
- The demo starts a temporary public proxy on `http://localhost:8090` to model a single local public feed root. Production deployments still need a real TLS reverse proxy.
- Docs assets live under `docs/assets/`; final PNGs were rendered from exact SVG specs after the image-generation workflow produced draft concepts with label inaccuracies.
- `docs/tutorials/calitp-readiness-checklist.md` is readiness guidance, not compliance proof.

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
