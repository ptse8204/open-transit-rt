# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current repository state

This repository is an early-stage starter for **Open Transit RT**.

The current codebase is not yet a production implementation. It is a scaffold plus requirements/docs.

## What exists now

### Repo guidance and architecture docs
The repo now has:
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

### Starter code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`

These services are scaffolds, not complete implementations.

### Starter database schema
The repo includes an initial SQL schema, but it is not yet a full migration-based production schema.

### Example data
The repo includes basic sample telemetry input and a basic docker-compose setup.

## What does not exist yet

The following are still missing or incomplete unless a later handoff says otherwise:

- durable DB-backed telemetry implementation
- migration command and versioned migrations
- `.env.example`
- one-command bootstrap flow
- integration fixtures under `testdata/`
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

**Active phase:** Phase 0 — Scaffolding and repo hardening

## Current objective

The immediate goal is to make the repo runnable, testable, and properly documented before major feature work begins.

This includes:
- scaffolding
- migration flow
- fixtures
- bootstrap flow
- env template
- task/build flows
- status/handoff docs

## Architecture posture

The codebase must preserve these long-term rules:
- mostly Go backend
- Postgres/PostGIS source of truth
- Vehicle Positions first
- Trip Updates pluggable
- draft GTFS separate from published GTFS
- conservative matching
- external dependencies isolated behind adapters
- no rider app / payments / dispatcher CAD scope

## Known blocking environment constraints

At the time of the current plan:
- `go` may not be on `PATH` in the execution shell
- `task` may not be on `PATH`
- Docker Compose is available
- the folder may not be a git repo

These constraints must be re-checked by the active Codex instance before running commands.

## Next recommended step

Complete Phase 0 before implementing substantive feature work.

That means:
1. add repo scaffolding from `docs/repo-gaps.md`
2. align build/bootstrap/test flows
3. create migration infrastructure
4. add fixtures and status docs
5. run baseline checks and record blocked commands

## What not to do next

Do not:
- jump straight into Trip Updates implementation
- add rider-facing functionality
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths