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
- simple `telemetry-ingest` HTTP service
- simple `feed-vehicle-positions` HTTP service
- shared domain models
- Phase 0 scaffolding for migrations, bootstrap, fixtures, and handoffs
- architecture and Codex handoff docs

Not yet implemented:
- protobuf GTFS-RT encoding
- Postgres persistence runtime behavior
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
PORT=8082 go run ./cmd/telemetry-ingest
```

Example request:
```bash
curl -X POST http://localhost:8082/v1/telemetry \
  -H 'Content-Type: application/json' \
  --data @examples/telemetry.json
```

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

Useful Phase 0 commands:
```bash
make migrate-status
make test-integration
make validate
```

## Recommended next build order

1. persist telemetry to Postgres/PostGIS
2. build deterministic trip matcher
3. publish true GTFS-RT protobuf Vehicle Positions
4. integrate TheTransitClock behind a prediction adapter
5. add GTFS Studio with draft/publish workflow
