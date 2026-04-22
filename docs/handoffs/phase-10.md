# Phase Handoff Template

All future phase handoff files must use this structure unless the phase explicitly documents a reason to diverge.

## Phase

Phase 10 — Docs, Tutorials, Deployment, and Demo

## Status

- Complete for the current codebase surface
- Active phase after this handoff: Phase 11 — Compliance Evidence and Optional External Integrations

## What Was Implemented

- Rewrote `README.md` to match the Phase 9 runtime surface: public protobuf feeds, protected JSON/admin routes, admin JWT generation, device-token telemetry ingest, pinned validators, local quickstart, deployment path, limitations, and truthful CAL-ITP/Caltrans-aligned wording.
- Added tutorial docs:
  - `docs/tutorials/local-quickstart.md`
  - `docs/tutorials/deploy-with-docker-compose.md`
  - `docs/tutorials/agency-demo-flow.md`
  - `docs/tutorials/production-checklist.md`
  - `docs/tutorials/calitp-readiness-checklist.md`
- Added `scripts/demo-agency-flow.sh`, `make demo-agency-flow`, and `task demo:agency`.
- The demo flow bootstraps the DB, installs/checks validators, imports sample GTFS, starts services, bootstraps publication metadata, ingests authenticated telemetry, fetches and verifies `schedule.zip`, fetches `feeds.json`, fetches public realtime protobuf feeds, verifies protected debug/admin routes including GTFS Studio, creates a demo Service Alert, runs validation flow, and reads scorecard/consumer-ingestion records.
- Updated `scripts/bootstrap-dev.sh` so the bootstrap output lists current services, public feed URLs, protected debug/admin examples, validator commands, and the executable demo target.
- Added repo-owned docs assets and asset notes:
  - `docs/assets/architecture-overview.png`
  - `docs/assets/agency-deployment.png`
  - `docs/assets/quickstart-flow.png`
  - `docs/assets/public-vs-admin-endpoints.png`
  - `docs/assets/README.md`
- Updated `docs/dependencies.md` to document local demo packaging tools.
- Updated `docs/current-status.md`, `docs/backlog.md`, and this handoff to reflect Phase 10.

## What Was Designed But Intentionally Not Implemented Yet

- No new backend product features were added.
- No external predictor integration was added.
- No SSO/login UI, consumer submission API, production app container packaging, Kubernetes manifests, or compliance evidence packet was added.
- The local demo uses a temporary public proxy so `PUBLIC_BASE_URL`/`FEED_BASE_URL` can represent a single feed root; that proxy is demo packaging, not a production reverse proxy implementation.
- Generated bitmap diagram drafts were not checked in because they had label inaccuracies; final repo assets were rendered from exact SVG specs after using the image-generation workflow for draft concepts.

## Schema And Interface Changes

- No database schema changes.
- Added Make target `demo-agency-flow`.
- Added Task target `demo:agency`.
- Added executable script `scripts/demo-agency-flow.sh`.

## Dependency Changes

- No runtime dependency changes.
- Documented local demo packaging tools: `curl`, `zip`, and `unzip`.
- Existing pinned validator tooling remains unchanged.

## Migrations Added

- None.

## Tests Added And Results

- No Go tests were added because Phase 10 was docs/tutorial/demo packaging.
- `scripts/demo-agency-flow.sh` is the executable acceptance path for the new demo flow and passed locally.

## Checks Run And Blocked Checks

- `make validators-install`: passed.
- `make validators-check`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make validate`: passed.
- `make demo-agency-flow`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make test-integration`: passed.
- `git diff --check`: passed.
- Blocked commands: none.

## Known Issues

- The current local Compose file provisions Postgres/PostGIS only. Deployment docs correctly describe app services as Go processes or deployment-owned service units, not as committed app containers.
- The demo starts services on fixed local ports `8081` through `8086` and a public proxy on `8090`; users must free those ports before running it.
- The demo proves local repo operability and feed/public-private boundaries. It does not prove production hosting, consumer acceptance, learned ETA quality, or full CAL-ITP/Caltrans compliance.
- Consumer-ingestion seed rows can have empty/NULL notes from older local DB state, so the demo upserts the default consumer records through the admin API before scorecard reads.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/phase-plan-production-closure.md`
  - `docs/prompts/calitp-truthfulness.md`
  - `docs/dependencies.md`
  - `docs/decisions.md`
- First files likely to edit:
  - a new Phase 11 evidence checklist doc
  - `docs/dependencies.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-11.md`
- Commands to run before coding:
  - `make validators-install`
  - `make validators-check`
  - `make test`
  - `make smoke`
  - `make demo-agency-flow`
- Known blockers:
  - Docker must be running for Postgres/PostGIS, pinned GTFS-RT validator wrapper, and the executable demo.
  - Stronger compliance or consumer-ingestion claims require deployment and third-party evidence, not repo assertions.
- Recommended first implementation slice:
  - Add a Phase 11 evidence checklist that separates code-complete capability, deployment-required proof, and external-consumer-confirmation-required evidence. Keep optional predictors deferred unless they are wired only through `internal/prediction.Adapter`.
