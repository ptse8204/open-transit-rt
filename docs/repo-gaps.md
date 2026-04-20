# Repo Gaps for Better Codex Performance

This document records what is currently missing from the starter repo in order to improve Codex productivity and reduce ambiguity.

---

## Present in the starter repo
- `Makefile`
- `docs/architecture.md`
- `docs/codex-task.md`
- `docs/conversation-summary.md`
- starter Go service scaffolding
- starter schema file
- example telemetry payload
- docker-compose with Postgres

---

## Missing items to add immediately

### 1. `.env.example`
Add environment variables for:
- `DATABASE_URL`
- `PORT`
- `AGENCY_ID`
- `FEED_BASE_URL`
- `METRICS_ENABLED`
- `LOG_LEVEL`
- `JWT_SECRET` or auth placeholder if auth is added later

### 2. `Taskfile.yml` or a more complete `Makefile`
Provide one-command workflows for:
- install/dev setup
- migrate database
- run all local services
- run tests
- run integration tests
- seed sample GTFS
- fetch dependencies
- lint/format

### 3. `cmd/migrate`
Provide an explicit migration binary that:
- applies versioned migrations
- checks migration status
- can roll back the latest migration
- works in CI and local dev

### 4. One-command dev bootstrap
Add a script or task that can:
- start Postgres/PostGIS
- apply migrations
- seed a sample agency and GTFS feed
- start core services
- print local URLs for feeds and admin UI

Suggested files:
- `scripts/bootstrap-dev.sh`
- or `task dev`

### 5. Integration test fixtures
Add deterministic fixtures for:
- small valid GTFS feed
- after-midnight GTFS feed
- frequency-based GTFS feed
- telemetry traces for matched, unmatched, stale, and swapped vehicles
- expected Vehicle Positions protobuf snapshots
- expected Trip Updates fixtures later

Suggested directory:
- `testdata/`

### 6. `docs/decisions.md`
Record architectural decisions such as:
- why Go is the main backend language
- why Postgres/PostGIS is used
- why Trip Updates is pluggable
- why GTFS Studio and GTFS import share the same published feed model
- why conservative matching is preferred over aggressive guessing

### 7. `docs/dependencies.md`
Document external dependencies and how they are wired in:
- Postgres
- PostGIS
- protobuf generation
- GTFS validator
- GTFS Realtime validator
- TheTransitClock or alternative predictor
- any optional admin UI framework
- Prometheus/Grafana if used

For each dependency, include:
- purpose
- required version
- how it is started
- how the app talks to it
- how to swap or disable it

---

## Why these gaps matter

Without these files and conventions, Codex is more likely to:
- infer the wrong setup flow
- guess missing environment variables
- misunderstand external dependencies
- skip migration or fixture work
- overfit to the starter demo instead of building a real runnable system

---

## Immediate Codex instruction

Add the missing repo scaffolding before major feature work:
1. `.env.example`
2. `cmd/migrate`
3. `docs/decisions.md`
4. `docs/dependencies.md`
5. integration fixtures under `testdata/`
6. one-command bootstrap
7. expanded build/test tasks
