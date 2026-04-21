# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current Repository State

This repository is an early-stage starter for **Open Transit RT**.

Phase 0 scaffolding, Phase 1 durable telemetry foundation, and Phase 2 deterministic trip matching are complete. The repo can format, test, start Postgres/PostGIS, run migrations, seed local agencies, execute the bootstrap flow, and run DB-backed telemetry plus matcher integration tests.

## What Exists Now

### Repo guidance and architecture docs
The repo has:
- `AGENTS.md`
- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/dependencies.md`
- `docs/phase-plan.md`
- `docs/decisions.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/handoffs/latest.md`

### Phase 0 scaffolding
The repo now has:
- `.env.example`
- `Taskfile.yml`
- independently usable `Makefile`
- `cmd/migrate`
- versioned migrations under `db/migrations`
- PostGIS-backed Docker Compose configuration on host port `55432`
- `scripts/bootstrap-dev.sh`
- deterministic fixtures under `testdata/`
- handoff template and Phase 0 handoff under `docs/handoffs/`

### Runtime code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`

`cmd/telemetry-ingest` now persists valid telemetry to Postgres through a telemetry repository. `agency-config` and `feed-vehicle-positions` remain starter scaffolds; `feed-vehicle-positions` still serves placeholder JSON from sample data and does not read persisted telemetry yet.

### Phase 1 telemetry foundation
The repo now has:
- `internal/db` with `pgxpool` connection setup and readiness ping support
- `internal/telemetry` repository interfaces and Postgres implementation
- DB-backed telemetry ingest in `cmd/telemetry-ingest`
- `/healthz` liveness and `/readyz` DB readiness behavior for telemetry ingest
- agency-scoped, bounded `/v1/events` debug listing
- durable parsed request payload storage in `telemetry_event.payload_json`
- atomic duplicate and out-of-order classification inside a transaction with a deterministic advisory lock
- DB-backed integration tests using `testdata/telemetry`
- development agency seeding through `scripts/seed-dev.sql`

### Phase 2 deterministic trip matching
The repo now has:
- `internal/gtfs` schedule-query boundary over existing published GTFS tables
- agency-local service-day resolution using agency timezone
- GTFS time parsing for times beyond `24:00:00`
- deterministic matcher engine in `internal/state`
- `internal/state.Engine` is the only valid production matcher entry point; legacy placeholder `RuleBasedMatcher` was removed
- `NewEngine` returns an error when schedule or assignment repositories are missing; `MustNewEngine` is available only for tests/bootstrap paths that intentionally want panic-on-error behavior
- conservative candidate scoring using trip hints, shape proximity, movement direction, stop progress, schedule fit, continuity, and block continuity
- time-aware continuity and block-transition scoring using configured windows
- block-transition scoring also requires the nearest plausible next-trip sequencing within the block when start-time identity is available; later same-block trips do not receive block-transition credit just for being later in the block
- explicit telemetry bearing validity is respected, including numeric `bearing: 0` for true north when the stored payload explicitly contains a numeric `bearing` field; malformed or null bearing payload values do not receive movement-direction credit, and non-DB callers without payload evidence treat zero as missing
- exact frequency candidate generation for `exact_times=1`
- conservative frequency-window identity behavior for `exact_times=0`
- non-exact frequency matches are marked as conservative window identities in score details so they are not mistaken for exact scheduled instances
- explicit unknown assignment persistence for stale, ambiguous, low-confidence, or missing-schedule cases
- distinct matcher system-failure reasons for agency lookup, service-day resolution, active-feed lookup, and schedule-query failures
- manual override precedence in matcher logic
- active manual overrides are evaluated before stale-telemetry fallback, so operator state is absolute until cleared or expired
- resolvable manual override assignments populate active `feed_version_id` and trip `block_id`, making override rows first-class persisted assignments alongside automatic matches
- Postgres assignment repository that closes prior active rows and persists assignment confidence, reasons, degraded state, score details, and incident linkage
- `shape_dist_traveled = 0` is preserved as a valid persisted value, not collapsed to NULL
- repeated identical degraded unknown states reuse the active degraded assignment only when degraded state, reason codes, service date, and telemetry evidence match; telemetry evidence means matching `telemetry_event_id` when present, with `active_from` equality only as the no-telemetry fallback
- batched GTFS schedule detail loading for stop times, shape points, and frequencies under the existing schedule-query boundary
- a small reason-code, degraded-state, and incident taxonomy
- unit and DB-backed integration tests for matcher edge cases

`vehicle_trip_assignment.score_details_json` is intentionally loose debug JSON in Phase 2, not a stable public schema. Matcher-generated score details include `score_schema`; candidate-based details also include `trip_id`, `start_time`, and `observed_local_seconds` when resolvable. Unknown assignment rows carry `service_date` whenever agency timezone and observed timestamp can be resolved; `service_date` is nullable only for truly unresolved cases. Missing shape data uses reason code `missing_shape` and degraded state `missing_shape`. Route-hint matching is reserved for a future input expansion and is not active in Phase 2 because telemetry does not currently carry a route hint.

Phase 2 service-day resolution considers the observed agency-local date and the immediately previous local date. That supports normal same-day service and practical after-midnight GTFS times through the prior service day, but it is not a generalized multi-day lookback for very long service patterns beyond that two-service-day window.

## Schema Source Of Truth

Migrations under `db/migrations` are the source of truth for executable schema changes and are applied through `cmd/migrate`.

`db/schema.sql` is deprecated as an executable schema. It is intentionally a comment-only pointer to the migrations directory and must not be edited independently.

## What Does Not Exist Yet

The following are still missing or incomplete unless a later handoff says otherwise:

- complete GTFS import pipeline
- complete GTFS Studio draft/publish workflow
- real GTFS-RT Vehicle Positions protobuf feed from persisted data
- Trip Updates adapter implementation
- Alerts feed implementation
- compliance dashboard
- consumer ingestion workflow
- robust auth and role handling
- manual override workflows
- production observability and SLO reporting

## Current Phase

**Active phase:** Phase 3 — Vehicle Positions production feed

Phase 2 is semantically closed, not just feature-complete. The next Codex instance should start with `docs/handoffs/latest.md`.

## Architecture Posture

The codebase must preserve these long-term rules:
- mostly Go backend
- Postgres/PostGIS source of truth
- Vehicle Positions first
- Trip Updates pluggable
- draft GTFS separate from published GTFS
- conservative matching
- external dependencies isolated behind adapters
- no rider apps, payments, passenger accounts, or dispatcher CAD scope

## Phase 0 Closure Audit Results

Checked during Phase 0 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `command -v gofmt`: passed, `/usr/local/bin/gofmt`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `go mod tidy`: passed and generated `go.sum`.
- `make fmt`: passed.
- `make test`: passed.
- `make db-up`: passed after changing local PostGIS host port to `55432`.
- `make migrate-up`: passed and applied `000001_initial_schema.sql`.
- `make migrate-status`: passed and reports migration version 1 applied.
- `make test-integration`: passed; this is currently a Phase 0 integration smoke path that verifies database reachability, migration visibility, and package compilation. There are no DB-backed integration test files yet.
- `scripts/bootstrap-dev.sh`: passed and reports no pending migrations.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make validate`: passed Phase 0 scaffold validation. It checks required migration and fixture scaffolding only; canonical GTFS and GTFS-RT validators are documented but not wired.
- `make lint`: passed optional fallback. `golangci-lint` is not installed, and future CI should make lint required once configured.
- `git diff --check`: passed.
- handoff path audit: passed; repo docs use `docs/handoffs/latest.md` and the retired singular path has been removed.
- Task equivalents were not run because `task` is not installed; Task remains optional because Makefile is independently usable.

## Phase 1 Closure Audit Results

Checked during Phase 1 closure:
- `go mod tidy`: passed.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed.
- `make migrate-up`: passed and applied `000002_telemetry_ingest_foundation.sql`.
- `make migrate-status`: passed and reports migration versions 1 and 2 applied.
- `make test-integration`: passed with DB-backed telemetry tests using an isolated temporary database.
- migration down/up smoke for `000002_telemetry_ingest_foundation.sql`: passed via `make migrate-down`, `make migrate-up`, and `make migrate-status`.
- `scripts/bootstrap-dev.sh`: passed and seeds `demo-agency`, `overnight-agency`, and `freq-agency`.
- `/readyz` behavior: covered by handler tests for both DB-ready and DB-unavailable responses.
- advisory-lock behavior: lock-key derivation is covered by deterministic unit tests; repository integration tests exercise classification through the locked `Store` path, but there is no separate concurrent-ingest stress test yet.
- `make validate`: passed scaffold and durable telemetry file validation only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.
- Optional Task equivalents were not run because `task` is not installed.

## Phase 2 Closure Audit Results

Checked during Phase 2 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- Initial pre-coding `make fmt`: blocked while Plan Mode was active because it runs `gofmt -w ./cmd ./internal`; it was run successfully after implementation.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed.
- `make migrate-up`: passed and applied `000003_deterministic_matching.sql`.
- `make migrate-status`: passed and reports migration versions 1, 2, and 3 applied.
- migration down/up smoke for `000003_deterministic_matching.sql`: passed via `make migrate-down`, `make migrate-up`, and `make migrate-status`.
- `make test-integration`: passed with DB-backed telemetry and matcher tests using isolated temporary database setup.
- `make validate`: passed Phase 2 scaffold, telemetry, and matcher-file validation only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.
- Optional Task equivalents were not run because `task` is not installed.

Phase 2 quality-hardening pass results:
- preserved Phase 2 scope only; no Phase 3 runtime work was added.
- made continuity and block-transition scoring require temporal plausibility through configured windows.
- fixed partial matcher config merging so zero fields fall back individually instead of replacing the whole config.
- separated repository/config/resolution failures from true no-schedule-candidate outcomes.
- replaced per-trip GTFS detail queries with batched stop-time, shape-point, and frequency fetches.
- strengthened non-exact frequency score details.
- added DB-backed integration coverage for after-midnight, exact and non-exact frequencies, ambiguous candidates, block transition, and unknown-row replacement.
- removed the legacy placeholder matcher path so the handoff now matches the actual production matcher implementation.
- added the final priority fixes for absolute manual override precedence, true-north bearing validity, zero shape-distance persistence, cleaner `NewEngine` construction, block-transition sequencing, and degraded-state deduplication.
- tightened the final semantic edge cases: degraded dedupe now includes service date and telemetry evidence, block-transition credit is limited to the nearest plausible successor, manual overrides persist feed/block context when resolvable, malformed/null bearings are invalid, and tests cover the two-day service-day boundary plus unknown replacement invariants.
- verified after the semantic-closure pass that the Phase 2 handoff matches the actual implementation.

## Next Recommended Step

Begin Phase 3 using the exact recommendation in `docs/handoffs/latest.md` and `docs/handoffs/phase-02.md`.

The first implementation slice should be:
1. replace placeholder Vehicle Positions behavior with output from persisted latest telemetry and persisted assignments
2. add GTFS-RT protobuf Vehicle Positions serialization
3. expose a stable public protobuf endpoint while keeping JSON debug output
4. preserve stale and unmatched behavior from Phase 2 assignments
5. add feed validation-oriented tests without implementing Trip Updates

## What Not To Do Next

Do not:
- jump straight into Trip Updates implementation
- add rider-facing functionality
- add payments, passenger accounts, or dispatcher CAD
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths once real feed generation starts
