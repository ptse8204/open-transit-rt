# Phase Handoff Template

All future phase handoff files must use this structure unless the phase explicitly documents a reason to diverge.

## Phase

Phase 11 — Compliance Evidence and Optional External Integrations

## Status

- Complete for the selected evidence-only path
- Active phase after this handoff: deployment evidence hardening track, if requested

## What Was Implemented

- Added `docs/compliance-evidence-checklist.md` as the Phase 11 evidence package.
- The checklist separates repo-proven capability, deployment/operator proof, and third-party confirmation.
- The checklist covers stable public URLs, public publication, open license/contact metadata, GTFS Schedule, Vehicle Positions, Trip Updates, Alerts, canonical validator workflow, consumer-ingestion workflow records, deployment/security/ops prerequisites, consumer acceptance limits, and marketplace/vendor-equivalence limits.
- Mapped current repo support to Caltrans/CAL-ITP-style expectations using evidence-bounded wording.
- Updated `docs/dependencies.md` with a Phase 11 wiring reality table for the external tools and repos originally mentioned in project scope.
- Updated README and tutorial wording with links to the evidence checklist and explicit observability/consumer-ingestion limits.
- Updated `docs/current-status.md` and `docs/handoffs/latest.md` to make Phase 11 the closed current state.

## What Was Designed But Intentionally Not Implemented Yet

- No new backend product features were added.
- No TheTransitClock adapter was added.
- No external predictor adapter was added.
- No Prometheus/Grafana deployment assets were added.
- No OpenTelemetry SDK, collector, exporter, or tracing configuration was added.
- No external consumer submission API integration was added.
- No SSO/login UI, marketplace packaging, consumer acceptance workflow automation, or major architecture change was added.

## Schema And Interface Changes

- No database schema changes.
- No public API changes.
- No Go interface changes.
- No environment variable changes.
- `internal/prediction.Adapter` remains the only approved boundary for future external predictor work.

## Dependency Changes

- No runtime dependency changes.
- `docs/dependencies.md` now marks the following as integrated or tool-backed where true: Postgres/PostGIS, pgx, Goose, MobilityData static GTFS validator, Docker-backed MobilityData GTFS-RT validator wrapper, GTFS-RT protobuf bindings, Go toolchain, Docker/Docker Compose, Task, local demo tools, and internal Prometheus-format `/metrics`.
- `docs/dependencies.md` now marks TheTransitClock, other external predictors, Prometheus/Grafana deployment, OpenTelemetry, consumer submission APIs, Mobility Database, and transit.land as deferred or workflow-only.

## Migrations Added

- None.

## Tests Added And Results

- No Go tests were added because Phase 11 is a docs/evidence-only closure pass.
- Existing tests and smoke/demo checks passed after the docs changes.

## Checks Run And Blocked Checks

Pre-edit baseline:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make validators-install`: passed.
- `make validators-check`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make demo-agency-flow`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make validate`: passed.
- `make migrate-status`: passed and reported migrations 1 through 8 applied.
- `make test-integration`: passed.
- `git diff --check`: passed.

Post-edit closure:
- `make validators-check`: passed.
- `make validate`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make demo-agency-flow`: passed.
- `make test-integration`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `git diff --check`: passed.

Blocked commands: none.

## Known Issues

- The repo supports deployment toward Caltrans/CAL-ITP-style readiness but does not prove full compliance by itself.
- Deployment evidence is still required for real HTTPS URLs, public fetch proof, production validator records, monitoring/alerting, backup/operations proof, and live scorecard export.
- Third-party confirmation is still required before claiming Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land, or other consumer acceptance.
- The internal deterministic Trip Updates adapter is conservative schedule-deviation prediction, not production-grade learned ETA quality.
- Prometheus-format `/metrics` exists when enabled, but Prometheus/Grafana deployment and OpenTelemetry tracing are not wired.
- Open Transit RT is not yet a marketplace/vendor-equivalent package.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/compliance-evidence-checklist.md`
  - `docs/dependencies.md`
  - `docs/prompts/calitp-truthfulness.md`
- First files likely to edit:
  - deployment evidence docs or runbooks for a real hosted environment
  - `docs/compliance-evidence-checklist.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - a new handoff file for the next track
- Commands to run before coding:
  - `make validators-check`
  - `make validate`
  - `make test`
  - `make smoke`
  - `make demo-agency-flow`
  - `make test-integration`
  - `docker compose -f deploy/docker-compose.yml config`
  - `git diff --check`
- Known blockers:
  - Docker must be running for DB-backed checks, the GTFS-RT validator wrapper, and the agency demo.
  - Stronger compliance or consumer-ingestion claims require deployment and third-party evidence.
- Recommended first implementation slice:
  - Start deployment evidence hardening: collect public HTTPS fetch proof, production validator records, monitoring/alerting evidence, scorecard export, and real third-party submission or acceptance records where available.
