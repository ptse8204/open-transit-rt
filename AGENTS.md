# AGENTS.md

## Purpose

This repository builds **Open Transit RT**, a modular transit data platform for small agencies.

The product goal is:

- import or author static GTFS
- ingest vehicle telemetry
- perform conservative deterministic trip matching
- publish high-quality GTFS-RT Vehicle Positions first
- keep Trip Updates pluggable
- support eventual compliance with the Caltrans / Cal-ITP transit data expectations
- avoid scope creep into unrelated transit software categories

This file applies to the entire repository unless a deeper `AGENTS.md` overrides it.

---

## Instruction priority

Treat the following as binding, in this order:

1. direct user instructions
2. this `AGENTS.md`
3. repository docs listed below
4. local code conventions inferred from existing code

If instructions conflict, prefer the higher-priority source.

---

## Binding repository docs

Before making changes, read and obey:

- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/dependencies.md`

Interpret `docs/codex-task.md` as the **implementation order**.

Interpret the requirements docs as the **full product contract**.

Do not optimize only for MVP completion. Design schemas, interfaces, and service boundaries so the system can later satisfy all listed requirements without major rewrites.

---

## Product boundaries

This repository is for:

- GTFS import
- GTFS Studio
- telemetry ingest
- deterministic trip matching
- GTFS-RT Vehicle Positions
- GTFS-RT Trip Updates integration boundary
- GTFS-RT Alerts
- feed validation
- monitoring
- admin/operator workflows

Do **not** add:

- rider-facing mobile apps
- fare payments
- CAD/dispatch systems
- passenger account systems
- marketing websites
- unrelated analytics products

If a task drifts toward those areas, stop and keep the implementation inside the defined product boundary.

---

## Core architecture rules

### 1. Mostly Go
Keep the backend codebase mostly in Go.

Prefer:
- Go services
- Go internal packages
- Postgres / PostGIS
- simple HTML or minimal web UI for early admin flows

Do not introduce a heavy frontend stack unless the task clearly requires it.

### 2. Vehicle Positions first
The first production-grade public output is:

- `GTFS-RT Vehicle Positions`

Trip Updates are important, but they are **not** the first implementation target.

### 3. Trip Updates must remain pluggable
Do not hard-wire Trip Updates logic into telemetry ingest or Vehicle Positions publishing.

Use a narrow prediction adapter boundary with:

- input:
  - active GTFS feed version
  - current telemetry
  - current vehicle assignments
  - Vehicle Positions feed URL or feed data
- output:
  - Trip Updates feed
  - optional diagnostics

The codebase must be able to support:
- an internal deterministic predictor
- an external predictor such as TheTransitClock
- a later replacement predictor

### 4. Conservative matching
Prefer `unknown` over false certainty.

Do not emit a trip descriptor unless matching confidence is above the configured threshold.

Manual overrides must take precedence over automatic matching.

### 5. Draft and published GTFS must stay separate
Do not collapse draft GTFS editing and published feed versions into one model.

GTFS Studio and GTFS ZIP import are two sources for the same published feed model, but draft data and active published data must remain distinct.

---

## External dependency rules

This repo may depend on external tools and codebases.

Examples include:
- Postgres
- PostGIS
- GTFS validators
- GTFS Realtime validators
- protobuf toolchains
- TheTransitClock or another predictor

Never assume external codebases “just fit.”

When touching integration points:

1. document the dependency in `docs/dependencies.md`
2. define the adapter/interface explicitly in code and docs
3. make failure modes explicit
4. add tests or stubs for the integration boundary
5. avoid tightly coupling internal logic to external implementation details

If an external tool does not fit cleanly:
- isolate it behind an adapter
- keep the public repo contracts stable
- do not spread tool-specific assumptions throughout the codebase

Do not treat any external codebase as the source of truth for core internal state unless the docs explicitly say so.

---

## Required repo scaffolding

If any of these are missing, add them before major feature work proceeds:

- `.env.example`
- `cmd/migrate`
- one-command bootstrap flow
- integration fixtures in `testdata/`
- `docs/decisions.md`
- `docs/dependencies.md`

Expand the `Makefile` or add `Taskfile.yml` if needed.

---

## Required documentation updates

When you add or change architecture-significant behavior, update the relevant docs.

At minimum:
- update `docs/decisions.md` for significant architectural choices
- update `docs/dependencies.md` for new external integrations
- update requirements docs if implementation reveals missing requirement language
- update bootstrap/test instructions if commands change

---

## Matching and realtime rules

Any implementation affecting trip matching, assignments, or GTFS-RT must preserve support for:

- agency-local service day
- after-midnight trips
- repeated trip instances
- `start_date`
- `start_time`
- `frequencies.txt`
- block continuity
- stale telemetry handling
- low-confidence unknown state
- manual override precedence
- auditability of assignment changes

Do not simplify away these requirements for convenience.

---

## Validation and compliance rules

Design every public feed path so it can eventually satisfy:

- stable public URLs
- validator-clean GTFS
- validator-clean GTFS-Realtime
- open-license publication
- discoverability metadata
- major consumer ingestion workflow
- complete realtime feed set:
  - Trip Updates
  - Vehicle Positions
  - Alerts

Do not assume current MVP scope removes the need to design for these outcomes.

---

## Testing requirements

For every meaningful code change, add or update tests.

Priority test areas:
- telemetry persistence
- GTFS import
- trip matching
- stale/unmatched behavior
- protobuf feed generation
- draft/publish behavior
- adapter boundaries
- rollback and failure cases

Use fixtures whenever possible.

Add integration tests for:
- valid GTFS import
- after-midnight service
- frequency-based service
- matched vehicle
- unmatched vehicle
- stale vehicle
- block transition

---

## Commands and checks

Before finishing work, run the repository’s relevant checks.

At minimum, once available, run:
- formatting
- unit tests
- integration tests
- migration checks
- feed validation checks where implemented

If a required check cannot run, state exactly why.

If you add a new required check, document it here or in the build docs.

---

## Code style

Prefer:
- small focused packages
- explicit interfaces at boundaries
- context-aware DB and HTTP calls
- structured logging
- deterministic tests
- backward-compatible API evolution where possible

Avoid:
- global mutable state
- hidden coupling across services
- speculative abstractions with no immediate use
- embedding external tool assumptions deep inside core domain logic

---

## Implementation posture

Build the smallest thing that is correct and extensible.

Do not build fake demo behavior once real persistence or real feed generation is expected.

Do not leave placeholder sample feed data in production paths.

If a feature is not fully implemented, prefer:
- a documented no-op adapter
- a clear TODO with dependency notes
- a disabled code path

over a misleading partial simulation.
