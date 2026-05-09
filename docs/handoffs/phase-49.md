# Phase 49 Handoff — External Predictor Runtime Adapter

## Status

Phase 49 is complete for the approved generic external predictor runtime
adapter scope.

## What Changed

- Added a shared Trip Updates adapter factory/config path in
  `internal/prediction`, used by both `cmd/feed-trip-updates` and
  `cmd/agency-config`.
- Preserved `TRIP_UPDATES_ADAPTER=deterministic` as default and
  `TRIP_UPDATES_ADAPTER=noop` as an explicit fallback.
- Added `TRIP_UPDATES_ADAPTER=external_http` for a generic operator-owned HTTP
  sidecar at fixed path `/v1/predict/trip-updates`.
- Added `TRIP_UPDATES_ADAPTER=external_http_shadow`, which returns
  deterministic public output and records only bounded redacted shadow
  diagnostics.
- Added strict external URL validation: exact allowlisted host, no
  userinfo/query/fragment, no redirects, HTTPS except loopback test stubs, and
  bounded timeout/request/response bytes.
- Added env-name-only bearer token lookup through
  `TRIP_UPDATES_EXTERNAL_HTTP_TOKEN_ENV`; missing referenced tokens fail config
  before calls, and token values are not logged, persisted, or diagnosed.
- Added sanitized external request/response DTOs that do not marshal
  `telemetry.StoredEvent`, `telemetry.Event`, or `state.Assignment` directly.
- Added tests for adapter selection, URL/token validation, DTO redaction,
  response/failure handling, shadow diagnostics, and valid empty Trip Updates
  behavior on external failures.

## Boundaries Preserved

Phase 49 did not add TheTransitClock-specific runtime code, start external
processes, vendor code, Java/Maven/Tomcat/TTC packaging, public/admin routes,
queues, schedulers, daemons, webhooks, evidence writes, consumer status
changes, compliance claims, vendor-compatibility claims, or production-grade
ETA claims.

Vehicle Positions, telemetry ingest, GTFS import, GTFS Studio, Alerts,
assignments, audit state, and publication workflows remain independent of the
external predictor runtime.

## Verification

Master closeout verification passed:

```bash
go test ./internal/prediction
go test ./internal/prediction ./internal/feed/tripupdates ./cmd/feed-trip-updates ./cmd/agency-config ./internal/architecture
go test ./cmd/feed-vehicle-positions ./internal/feed
make validate
make realtime-quality
make test
make smoke
make test-integration
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
exact seven-target prepared-only consumer tracker check
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config
```

The consumer tracker remains unchanged with exactly seven targets, all
`prepared`. No files under `docs/evidence` were edited.

## Next Phase

Phase 50 — Realtime Quality Backtesting should start with a fresh read-only
planning sub-agent pass before implementation.
