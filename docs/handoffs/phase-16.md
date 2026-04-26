# Phase Handoff

## Phase

Phase 16 — Agency Onboarding Product Packaging

## Status

- Complete for the Phase 16 local packaging, first-run documentation, device onboarding, and operator-friendly command output scope.
- Active phase after this handoff: Phase 17 — Deployment Automation And Pilot Operations.

## What Was Implemented

- Added a full local app packaging path behind `make agency-app-up`.
- Added `scripts/agency-local-app.sh` with `up`, `down`, `logs`, and destructive `reset` subcommands.
- Added a local Docker app image through `deploy/Dockerfile.local`.
- Extended `deploy/docker-compose.yml` with an `app` profile for Postgres/PostGIS, all Go services, GTFS Studio, and a local Caddy reverse proxy.
- Added `deploy/Caddyfile.local` to expose the local demo package at `http://localhost:8080`.
- Added Makefile and Taskfile targets for `agency-app-up`, `agency-app-down`, `agency-app-logs`, and `agency-app-reset`.
- Added `scripts/device-onboarding.sh` for local device token rebind, one-event sample telemetry, simulator-style telemetry, help output, and dry-run payload display.

## What Was Designed But Intentionally Not Implemented Yet

- No backend API contracts were changed.
- No database schema or migration was added.
- No public feed URL paths were changed.
- No consumer-submission status was advanced.
- No hosted SaaS, Kubernetes, external predictor integration, consumer submission API, or major admin UX rewrite was added.
- Validator tooling remains optional for local app startup; validation setup is printed as the next step unless a validation workflow is run explicitly.

## Packaging/Profile/Script Changes

- `make agency-app-up` now starts the full local stack, applies migrations, seeds demo records, imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for readiness, verifies public feed URLs, and prints final operator guidance.
- `make agency-app-up` is idempotent for repeated local runs: it tolerates existing containers and volumes, safely reapplies migrations, reruns seed/import/bootstrap operations, and prints reset guidance when local state exists.
- `make agency-app-reset` is visibly destructive. It states that it removes local containers, the Compose Postgres volume, generated local env files if present, local demo database state, and container logs. `scripts/agency-local-app.sh reset --force` supports automation.
- Compose healthchecks were added for Postgres, Go service health endpoints, and the local reverse proxy.
- The local app image excludes `.cache`, local env files, private keys, logs, generated deploy binaries, and git metadata through `.dockerignore`.

## Documentation Changes

- Added `docs/tutorials/agency-first-run.md`.
- Updated `README.md`, `wiki/README.md`, `docs/README.md`, `docs/tutorials/README.md`, `docs/tutorials/local-quickstart.md`, `docs/tutorials/agency-demo-flow.md`, `docs/tutorials/deploy-with-docker-compose.md`, and `docs/tutorials/production-checklist.md`.
- Updated `docs/dependencies.md` with the local Docker app package dependency boundary.
- Added ADR-0024 to `docs/decisions.md`.
- Updated `docs/current-status.md` and `docs/handoffs/latest.md`.

## Device Onboarding Changes

- `scripts/device-onboarding.sh help` explains commands and security notes.
- `scripts/device-onboarding.sh rebind` calls existing `POST /admin/devices/rebind` and prints the one-time token only because the existing API intentionally returns it.
- `scripts/device-onboarding.sh sample` sends one telemetry event to the existing `/v1/telemetry` API.
- `scripts/device-onboarding.sh simulate` sends a short sequence of telemetry events for local testing when no hardware exists.
- `--dry-run` for `sample` and `simulate` prints the target and payload without sending telemetry.

## Schema And Interface Changes

- No schema changes.
- No backend API contract changes.
- No public feed URL changes.
- New local CLI/script interfaces:
  - `scripts/agency-local-app.sh up|down|logs|reset [--force]`
  - `scripts/device-onboarding.sh help|rebind|sample|simulate`

## Dependency Changes

- Added local app packaging dependency on Docker Compose app profiles and the `caddy:2.8-alpine` image for local demo reverse proxying.
- Added local app image build from `golang:1.23.2-alpine` and `alpine:3.20`.
- No production dependency requirement was added; the local profile is explicitly demo/evaluation packaging.

## Migrations Added

- None.

## Tests Added And Results

- No Go tests were added because this phase changed packaging, scripts, and docs without backend behavior changes.
- Manual script and packaging verification covered the new local app and device helper paths.

## Checks Run And Blocked Checks

- `sh -n scripts/agency-local-app.sh scripts/device-onboarding.sh scripts/demo-agency-flow.sh scripts/bootstrap-dev.sh`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `docker compose -f deploy/docker-compose.yml --profile app config`: passed.
- `make agency-app-up`: passed after fixing local Caddy route ordering.
- `make agency-app-up` repeated against the existing local stack: passed.
- `scripts/device-onboarding.sh sample --dry-run`: passed.
- `scripts/device-onboarding.sh simulate --dry-run`: passed.
- `scripts/device-onboarding.sh sample`: passed.
- `scripts/device-onboarding.sh simulate`: passed.
- `make agency-app-logs`: passed.
- `printf 'reset-local-app\n' | make agency-app-reset`: passed and verified the destructive confirmation text.
- `make validate`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make demo-agency-flow`: passed after tightening the existing database wait logic to verify the host `DATABASE_URL` path before migrations.
- `git diff --check`: passed.
- Blocked commands: none.

## Known Issues

- Local app startup builds a local image and may take time on first run.
- The local app profile uses development defaults and `http://localhost:8080`; it must not be treated as production TLS or admin network policy.
- Validator execution is not part of local app startup success. Operators must run `make validators-install validators-check` and validation workflows separately.
- Repeated `agency-app-up` runs publish another sample GTFS feed version as the current active local feed. This is intentional for idempotent local evaluation; use `make agency-app-reset` for a clean local database.

## Exact Next-Step Recommendation

- First files to read:
  - `docs/current-status.md`
  - `docs/handoffs/phase-16.md`
  - `docs/tutorials/agency-first-run.md`
  - `docs/tutorials/deploy-with-docker-compose.md`
  - `docs/phase-17-deployment-automation-pilot-operations.md`
- First files likely to edit:
  - `scripts/oci-pilot.sh`
  - `deploy/oci/`
  - `deploy/systemd/`
  - `docs/runbooks/`
  - `docs/handoffs/latest.md`
- Commands to run before coding:
  - `make validate`
  - `make test`
  - `make smoke`
  - `make agency-app-up`
  - `docker compose -f deploy/docker-compose.yml --profile app config`
  - `git diff --check`
- Known blockers:
  - Phase 15 found real secrets in ignored local `.cache` files; affected pilot/admin/device/TLS secrets still need rotation or revocation before further real pilot use.
- Recommended first implementation slice:
  - Start Phase 17 by turning the proven local app startup sequence into deployment automation checks for the OCI/systemd path, while preserving the same truthful final-output pattern and secret boundaries.
