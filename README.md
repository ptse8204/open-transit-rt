# Open Transit RT

Open Transit RT is a starter monorepo for a small-agency transit data stack:
- static GTFS authoring/import
- device-based vehicle telemetry ingestion
- deterministic trip matching
- GTFS-RT Vehicle Positions publishing
- a pluggable Trip Updates engine
- alerts and monitoring later

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
- architecture and Codex handoff docs

Not yet implemented:
- Android client
- production Trip Updates prediction quality
- TheTransitClock integration or another real predictor
- alerts authoring UI
- GTFS Studio rich map/timetable editor

## Services

### agency-config
```bash
PORT=8081 go run ./cmd/agency-config
```

### telemetry-ingest
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" PORT=8082 go run ./cmd/telemetry-ingest
```

Example request:
```bash
curl -X POST http://localhost:8082/v1/telemetry \
  -H 'Content-Type: application/json' \
  --data @examples/telemetry.json
```

Telemetry timestamps must be RFC 3339 values with a timezone or offset. Unknown agencies are rejected; `make dev` or `make seed` creates the local fixture agencies.

Debug endpoint for local development:
```bash
curl 'http://localhost:8082/v1/events?agency_id=demo-agency&limit=25'
```

`/v1/events` is agency-scoped and bounded, but Phase 1 has no auth layer; production deployments should disable or protect it until admin/auth controls exist.

### feed-vehicle-positions
```bash
PORT=8083 go run ./cmd/feed-vehicle-positions
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/vehicle_positions.pb`
- `GET /public/gtfsrt/vehicle_positions.json`

`feed-vehicle-positions` requires `DATABASE_URL` and `AGENCY_ID`. The protobuf endpoint returns a valid empty `FeedMessage` with normal success headers when no telemetry is available or all vehicles are suppressed as stale.

### gtfs-studio
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" PORT=8086 go run ./cmd/gtfs-studio
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /admin/gtfs-studio?agency_id=demo-agency`

GTFS Studio is a minimal server-rendered admin surface for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies. Draft edits are stored separately from published GTFS rows. Published and discarded drafts are read-only by default.

### feed-trip-updates
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" AGENCY_ID=demo-agency FEED_BASE_URL=http://localhost:8083/public PORT=8084 go run ./cmd/feed-trip-updates
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/trip_updates.pb`
- `GET /public/gtfsrt/trip_updates.json`

Phase 6 Trip Updates use an explicit no-op prediction adapter by default. The protobuf endpoint returns a valid empty `FeedMessage`; the JSON endpoint exposes diagnostics and traceability.

### feed-alerts
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" AGENCY_ID=demo-agency PORT=8085 go run ./cmd/feed-alerts
```

Current endpoints:
- `GET /healthz`
- `GET /readyz`
- `GET /public/gtfsrt/alerts.pb`
- `GET /public/gtfsrt/alerts.json`

Phase 6 Alerts return a valid empty `FeedMessage` plus JSON-only deferred diagnostics. Alert authoring and persistence are not implemented yet.

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
make validate
```

`make test-integration` runs DB-backed telemetry tests. The tests prefer creating an isolated temporary database from `TEST_DATABASE_URL`; if that is not permitted, they fall back to an isolated temporary schema in the configured test database.

`make validate` is currently a scaffold, telemetry, matcher, and Vehicle Positions file smoke check only. It verifies required migration and fixture paths exist; canonical GTFS and GTFS-Realtime validators are documented but not wired yet.

## Recommended next build order

1. add stop-level Trip Updates prediction behind the Phase 6 adapter
2. add operations workflows for prediction repair, cancellations, detours, and alerts
3. add compliance and consumer workflow surfaces
