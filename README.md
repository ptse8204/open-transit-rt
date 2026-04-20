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
- shared domain models
- Phase 0 scaffolding for migrations, bootstrap, fixtures, and handoffs
- Phase 1 durable telemetry persistence foundation
- architecture and Codex handoff docs

Not yet implemented:
- protobuf GTFS-RT encoding
- Android client
- trip matching engine
- TheTransitClock integration
- alerts authoring UI
- GTFS Studio interactive editor

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
- `GET /public/gtfsrt/vehicle_positions.json`

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

`make validate` is currently a scaffold and Phase 1 telemetry-file smoke check only. It verifies required migration and fixture paths exist; canonical GTFS and GTFS-Realtime validators are documented but not wired yet.

## Recommended next build order

1. build deterministic trip matcher
2. publish true GTFS-RT protobuf Vehicle Positions
3. integrate TheTransitClock behind a prediction adapter
4. add GTFS import and publish pipeline
5. add GTFS Studio with draft/publish workflow
