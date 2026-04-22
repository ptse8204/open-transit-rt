# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 11 — Compliance Evidence and Optional External Integrations is complete for the selected evidence-only path.

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
- Phase 11 compliance evidence and external-integration reality review is implemented and complete for the evidence-only scope.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-11.md`
4. `docs/compliance-evidence-checklist.md`
5. `docs/dependencies.md`
6. `docs/decisions.md`
7. `docs/prompts/calitp-truthfulness.md`
8. `docs/phase-plan-production-closure.md`
9. `README.md`
10. `docs/tutorials/calitp-readiness-checklist.md`
11. `docs/tutorials/production-checklist.md`

## Current Objective

Start the next hardening track only if requested. Phase 11 closed the evidence layer but did not collect proof from a real public deployment. The next useful work is deployment evidence hardening: HTTPS feed-root proof, production validation records, monitoring/alerting assets, scorecard export, and third-party submission or acceptance records where real evidence exists.

## Exact First Commands

```bash
command -v go
go version
make validators-check
make validate
make test
make smoke
make demo-agency-flow
make test-integration
docker compose -f deploy/docker-compose.yml config
git diff --check
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
- TheTransitClock and other external predictors are not integrated; future work must keep them behind `internal/prediction.Adapter`.
- Prometheus/Grafana deployment assets and OpenTelemetry tracing/exporter wiring are not integrated.
- Consumer-ingestion workflow records exist, but external consumer submission APIs are not integrated.
- Compliance or consumer-acceptance claims need deployment and third-party evidence.

## First Files Likely To Edit

- deployment evidence docs or runbooks for a real hosted environment
- `docs/compliance-evidence-checklist.md` if evidence categories change
- `docs/dependencies.md` if a deferred external integration becomes real
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- a new handoff file for the next track

## Phase 11 Notes For Future Work

- `docs/compliance-evidence-checklist.md` is the Phase 11 evidence package.
- It separates repo-proven capability, deployment/operator proof, and third-party confirmation.
- `docs/dependencies.md` includes the Phase 11 wiring reality table.
- Real integrations are code-backed or tool-backed only where the repo actually wires them.
- TheTransitClock, other external predictors, Prometheus/Grafana deployment, OpenTelemetry, and consumer submission APIs remain deferred or workflow-only.
- Mobility Database and transit.land are documented workflow targets, not seeded default consumers or API integrations.

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
- Do not claim full CAL-ITP/Caltrans compliance, production readiness, marketplace equivalence, or consumer acceptance without actual evidence.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
