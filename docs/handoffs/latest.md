# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 9 — Production Closure is complete for the current codebase surface. Continue with Phase 10 from `docs/phase-plan-production-closure.md`.

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
- Phase 9 production closure is implemented for validator execution, validator tooling pins, admin auth/roles, device auth/binding, assignment current-row races, safer config defaults, debug endpoint protection, request logging/request IDs, metrics toggle, stronger feed-service readiness, and smoke coverage.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-09.md`
4. `docs/phase-plan-production-closure.md`
5. `docs/codex-task.md`
6. `docs/requirements-2a-2f.md`
7. `docs/requirements-trip-updates.md`
8. `docs/requirements-calitp-compliance.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Current Objective

Start Phase 10 — Docs, Tutorials, Deployment, and Demo. Preserve stable public protobuf feed URLs, protected JSON debug routes, auth-scoped admin mutations, device-token telemetry ingest, pinned validator setup, Phase 7 prediction boundaries, GTFS import, and GTFS Studio draft/publish behavior.

## Exact First Commands

```bash
command -v go
go version
make fmt
make validators-install
make validators-check
make test
make smoke
docker compose -f deploy/docker-compose.yml config
make db-up
make migrate-status
make test-integration
```

If Task is installed, optional equivalents may be run:

```bash
task fmt
task test
task smoke
task migrate:status
task test:integration
```

## Known Blockers

- Task is optional and may not be installed; Makefile remains independently usable.
- Docker must be running before DB-backed checks.
- Pinned validator tooling is repo-supported through `make validators-install` and `make validators-check`; CI and production automation still need to run or bake those steps into their environment.
- Full hosted login/SSO and server-side `jti` replay tracking are deferred; the current auth contract accepts HS256 admin JWTs plus optional browser cookie sessions.
- Consumer ingestion workflow records exist, but external consumer submission APIs are not integrated.

## First Files Likely To Edit

- `README.md`
- `docs/tutorials/local-quickstart.md`
- `docs/tutorials/deploy-with-docker-compose.md`
- `docs/tutorials/agency-demo-flow.md`
- `docs/tutorials/production-checklist.md`
- `docs/tutorials/calitp-readiness-checklist.md`
- `docs/assets/README.md`

## Phase 8 Notes For Future Work

- `/public/gtfs/schedule.zip` is generated on demand from active published GTFS tables.
- Schedule ZIP bytes are deterministic for unchanged active feed data; ZIP entry modified times and HTTP `Last-Modified` use the active feed revision time.
- Realtime `published_feed.revision_timestamp` is a publication/bootstrap metadata revision and must not change on every feed generation.
- Realtime freshness and generation health belong in `feed_health_snapshot`.
- `/public/feeds.json` reads per-feed data from `published_feed`; license/contact fields resolve from `feed_config` only when per-feed values are empty.
- `feed_config.publication_environment = 'production'` makes missing canonical validator execution red in scorecards. In `dev`, missing validators are yellow/not-run.
- Alerts authoring/persistence is owned by `internal/alerts`; GTFS-RT protobuf rendering is owned by `internal/feed/alerts`.
- Prediction packages must not import Alerts packages. Canceled-trip missing-alert review signals are satisfied by the Alerts-owned reconciler.
- Validator runs are allowlisted by `validator_id`; admin requests may provide only `validator_id`, `feed_type`, and optional `feed_version_id`.
- Validator execution uses server-owned local artifacts/temp files, argv-based `exec.CommandContext`, timeout/output/report caps, output confinement, and redacted argv/path reporting.
- Realtime validation prefers internal builder-derived Vehicle Positions, Trip Updates, and Alerts protobuf bytes; configured feed URLs are fallback only when an internal builder cannot be constructed.
- `/readyz` for `agency-config`, Trip Updates, and Alerts requires DB reachability plus the required active feed/config dependencies; DB ping alone is intentionally insufficient.
- Pinned validator tooling lives in `tools/validators/validators.lock.json`; `VALIDATOR_TOOLING_MODE=stub` is the explicit deterministic stub bypass for targeted tests only.
- Admin JWTs require `sub`, `agency_id`, `iat`, `exp`, `iss`, and `aud`; default TTL is 8h, clock skew allowance is 2m, `ADMIN_JWT_OLD_SECRETS` supports secret rotation, and `jti` replay tracking is deferred.
- Cookie auth is only for browser-admin flows; Bearer auth remains the default for machine/API admin calls.
- Telemetry ingest requires opaque device Bearer tokens. `POST /admin/devices/rebind` rotates token and binding immediately.
- `/public/gtfs/schedule.zip` returns `ETag`, `Last-Modified`, and `X-Checksum-SHA256`, with `SCHEDULE_ZIP_MAX_BYTES` bounding payload size.
- `/metrics` exists only when `METRICS_ENABLED=true` and should be treated as an internal/reverse-proxy-controlled operations surface, not an anonymous public feed.

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
