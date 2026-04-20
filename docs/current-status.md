# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current repository state

This repository is an early-stage starter for **Open Transit RT**.

The current codebase is not yet a production implementation. Phase 0 scaffolding is complete enough for the next implementation phase to begin once the Go toolchain is available.

## What exists now

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
- PostGIS-backed Docker Compose configuration
- `scripts/bootstrap-dev.sh`
- deterministic fixtures under `testdata/`
- handoff template and Phase 0 handoff under `docs/handoffs/`

### Starter code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`

These services are scaffolds, not complete implementations.

## What does not exist yet

The following are still missing or incomplete unless a later handoff says otherwise:

- durable DB-backed telemetry runtime implementation
- repository interfaces and DB connection package
- complete GTFS import pipeline
- complete GTFS Studio draft/publish workflow
- deterministic trip matcher with real edge-case handling
- real GTFS-RT Vehicle Positions protobuf feed from persisted data
- Trip Updates adapter implementation
- Alerts feed implementation
- compliance dashboard
- consumer ingestion workflow
- robust auth and role handling
- manual override workflows
- production observability and SLO reporting

## Current phase

**Active phase:** Phase 1 — Durable telemetry foundation

Phase 0 is complete. The next Codex instance should start with `docs/handoffs/latest.md`.

## Architecture posture

The codebase must preserve these long-term rules:
- mostly Go backend
- Postgres/PostGIS source of truth
- Vehicle Positions first
- Trip Updates pluggable
- draft GTFS separate from published GTFS
- conservative matching
- external dependencies isolated behind adapters
- no rider apps, payments, passenger accounts, or dispatcher CAD scope

## Known blocking environment constraints

Checked during Phase 0:
- `go` is not on `PATH`
- `gofmt` is not on `PATH`
- `task` is not on `PATH`
- Docker Compose is available
- Docker Compose config validates

Install or expose the Go toolchain before Phase 1 implementation checks can pass.

## Next recommended step

Begin Phase 1 with the exact entry recommendation in `docs/handoffs/latest.md` and `docs/handoffs/phase-00.md`.

The first implementation slice should be:
1. add a shared DB package using `pgxpool`
2. add telemetry repository interfaces and Postgres implementation
3. wire `cmd/telemetry-ingest` to persist telemetry
4. update health/readiness behavior
5. add DB-backed tests using `testdata/telemetry`

## What not to do next

Do not:
- jump straight into Trip Updates implementation
- add rider-facing functionality
- add payments, passenger accounts, or dispatcher CAD
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths once real feed generation starts
