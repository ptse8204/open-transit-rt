# Open Transit RT

Open Transit RT is a starter monorepo for a small-agency transit data stack:
- static GTFS authoring/import
- device-based vehicle telemetry ingestion
- deterministic trip matching
- GTFS-RT Vehicle Positions publishing
- a pluggable Trip Updates engine
- alerts, compliance workflows, and hardening foundations

This repo is a phased starter, not a finished product. It focuses on the wedge we discussed: **BYOD or low-cost GPS -> public realtime feeds**, with **Trip Updates** treated as a replaceable module.

## Current status

Implemented now:
- simple `agency-config` HTTP service
- DB-backed `telemetry-ingest` HTTP service
- simple `feed-vehicle-positions` HTTP service
- minimal `gtfs-studio` HTTP service
- shared domain models
- Phase 0 scaffolding for migrations, bootstrap, fixtures, and handoffs
- Phase 1 durable telemetry persistence foundation
- Phase 2 deterministic trip matching foundation
- Phase 3 DB-backed GTFS-RT Vehicle Positions protobuf and JSON debug feed
- Phase 4 GTFS ZIP import/publish pipeline
- Phase 5 GTFS Studio typed draft/publish model
- Phase 6 Trip Updates and Alerts architecture with no-op/default feeds
- Phase 7 deterministic Trip Updates prediction and operations workflows
- Phase 8 Alerts, publication metadata, compliance scorecards, consumer workflow, and validation records
- post-Phase-8 production hardening for admin auth, device auth, validator execution, assignment races, and safer config defaults
- architecture and Codex handoff docs

Not yet implemented:
- Android client
- production Trip Updates prediction quality
- TheTransitClock integration or another real predictor
- GTFS Studio rich map/timetable editor
- full hosted login/identity UI

## Services

Admin/API hardening:
- Public `.pb` feed endpoints remain anonymous and stable for consumers.
- Admin routes and JSON debug endpoints require admin authentication.
- Bearer JWT auth is the default for machine/API admin calls.
- `admin_session` cookie auth is only for browser-admin flows and requires CSRF tokens on unsafe methods.
- Admin JWTs require `sub`, `agency_id`, `iat`, `exp`, `iss`, and `aud`; default local TTL is `8h`, clock skew allowance is `2m`, and secret rotation accepts `ADMIN_JWT_SECRET` plus comma-separated `ADMIN_JWT_OLD_SECRETS`. `jti` is emitted for auditability but server-side replay tracking is deferred.

### agency-config
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" \
AGENCY_ID=demo-agency \
ADMIN_JWT_SECRET=dev-admin-jwt-secret-change-me \
ADMIN_JWT_ISSUER=open-transit-rt-local \
ADMIN_JWT_AUDIENCE=open-transit-rt-admin \
CSRF_SECRET=dev-csrf-secret-change-me \
DEVICE_TOKEN_PEPPER=dev-device-token-pepper-change-me \
PORT=8081 go run ./cmd/agency-config
```

`agency-config` serves `/public/gtfs/schedule.zip`, `/public/feeds.json`, admin publication bootstrap, consumer workflow records, scorecards, device rebinding, and allowlisted validation runs. `/public/gtfs/schedule.zip` emits `ETag`, `Last-Modified`, and `X-Checksum-SHA256`; `SCHEDULE_ZIP_MAX_BYTES` bounds generated payloads.

### telemetry-ingest
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" \
ADMIN_JWT_SECRET=dev-admin-jwt-secret-change-me \
ADMIN_JWT_ISSUER=open-transit-rt-local \
ADMIN_JWT_AUDIENCE=open-transit-rt-admin \
CSRF_SECRET=dev-csrf-secret-change-me \
DEVICE_TOKEN_PEPPER=dev-device-token-pepper-change-me \
PORT=8082 go run ./cmd/telemetry-ingest
```

Example request:
```bash
curl -X POST http://localhost:8082/v1/telemetry \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer dev-device-token' \
  --data @examples/telemetry.json
```

Telemetry timestamps must be RFC 3339 values with a timezone or offset. Device tokens are opaque Bearer tokens bound to agency, device, and vehicle. The development seed creates `device-1` bound to `bus-1` with token `dev-device-token`; rebinding is admin-managed through `POST /admin/devices/rebind` and immediately invalidates the old token/binding.

Admin debug endpoint:
```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  'http://localhost:8082/v1/events?agency_id=demo-agency&limit=25'
```

`/v1/events` and `/admin/debug/telemetry/events` require admin auth and derive agency scope from the token.

### feed-vehicle-positions
```bash
PORT=8083 go run ./cmd/feed-vehicle-positions
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/vehicle_positions.pb`
- `GET /public/gtfsrt/vehicle_positions.json`
- `GET /admin/debug/gtfsrt/vehicle_positions.json`

`feed-vehicle-positions` requires `DATABASE_URL`, `AGENCY_ID`, and admin auth config. The protobuf endpoint is public; JSON debug paths require admin auth and use the same underlying debug representation.

### gtfs-studio
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" PORT=8086 go run ./cmd/gtfs-studio
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /admin/gtfs-studio?agency_id=demo-agency`

GTFS Studio is a minimal server-rendered admin surface for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies. Draft edits are stored separately from published GTFS rows. It uses admin auth and CSRF validation for cookie-authenticated unsafe methods.

### feed-trip-updates
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" AGENCY_ID=demo-agency FEED_BASE_URL=http://localhost:8083/public PORT=8084 go run ./cmd/feed-trip-updates
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/trip_updates.pb`
- `GET /public/gtfsrt/trip_updates.json`
- `GET /admin/debug/gtfsrt/trip_updates.json`

Trip Updates default to the deterministic prediction adapter. The protobuf endpoint is public; JSON debug paths require admin auth and expose diagnostics/traceability from one shared debug builder.

### feed-alerts
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" AGENCY_ID=demo-agency PORT=8085 go run ./cmd/feed-alerts
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/alerts.pb`
- `GET /public/gtfsrt/alerts.json`
- `GET /admin/debug/gtfsrt/alerts.json`
- `GET/POST /admin/alerts`
- `POST /admin/alerts/{id}/publish`
- `POST /admin/alerts/{id}/archive`
- `POST /admin/alerts/reconcile-cancellations`

Alerts publish persisted Service Alerts. The protobuf endpoint is public; JSON debug and admin mutation routes require admin auth. Actor and agency come from the authenticated context.

## Local development

Copy local defaults if needed:
```bash
cp .env.example .env
```

Bring up Postgres/PostGIS and apply migrations:
```bash
make db-up
make migrate-up
```

One-command bootstrap:
```bash
make dev
```

Task is optional. The Makefile remains independently usable when `task` is not installed.

```bash
make build
make test
```

Useful local commands:
```bash
make migrate-status
make test-integration
make smoke
make validate
```

`make test-integration` runs DB-backed telemetry tests. The tests prefer creating an isolated temporary database from `TEST_DATABASE_URL`; if that is not permitted, they fall back to an isolated temporary schema in the configured test database.

`make smoke` runs the hardening HTTP smoke coverage, including unauthenticated admin/debug rejection and telemetry token checks. `make validate` verifies required hardening, validator, migration, and fixture paths exist. Canonical validators run through server-side allowlisted validator IDs when configured.

Production deployment notes:
- Set `APP_ENV=production`; services fail fast without `DATABASE_URL`, admin JWT config, `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`.
- `BIND_ADDR` defaults to `127.0.0.1`. Use `BIND_ADDR=0.0.0.0` only behind a TLS-terminating reverse proxy.
- Public feed URLs should be served over HTTPS by the proxy/domain layer.
- Do not expose `.json` debug paths without admin auth; this repo protects them by default.

## Recommended next build order

1. integrate real hosted identity or SSO in front of the JWT contract
2. add stronger feed SLO dashboards and metrics export
3. install/pin canonical validator distributions in CI and production automation
