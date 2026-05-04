# docs/dependencies.md

This document defines the external tools, libraries, and codebases that Open Transit RT may rely on, what role each one plays, how it integrates with the repository, and what to do if it fails or does not fit.

The purpose of this file is to stop the codebase from developing hidden or accidental coupling to outside systems.

---

## Dependency policy

For every external dependency or codebase:

- document its purpose
- pin or declare expected version range when possible
- document how it is started or provisioned
- define the integration boundary
- define failure behavior
- define how to replace or disable it
- do not let its internal types leak into core domain packages unless explicitly required

External integrations must be isolated behind adapters wherever practical.

---

## Dependency classification

Dependencies fall into four groups:

1. **Core runtime infrastructure**
   - required for the application to function normally

2. **Validation and developer tooling**
   - used to validate or build feeds and code, but not part of the runtime request path

3. **Optional prediction backends**
   - not required for Vehicle Positions
   - may be required later for Trip Updates quality

4. **Future optional integrations**
   - planned or possible, but not required for initial implementation

---

## Phase 11 wiring reality

This table is the current external dependency and integration status. It separates repo-wired tooling from optional systems that remain deferred or workflow-only.

| Dependency or external system | Current status | Evidence / boundary |
| --- | --- | --- |
| Postgres | Integrated | Core datastore through repository interfaces; local Compose uses PostgreSQL 16 image with host port `55432`. |
| PostGIS | Integrated | Spatial storage and queries are migration-backed and kept inside repository/matching boundaries. |
| pgx | Integrated | Go PostgreSQL driver and pool used by services and repositories. |
| Goose | Integrated | `cmd/migrate` applies versioned migrations. |
| MobilityData GTFS Validator | Integrated as validation tooling | Pinned `v7.1.0` JAR installed by `make validators-install` and checked by `make validators-check`; invoked through allowlisted validator IDs. |
| MobilityData GTFS Realtime Validator | Integrated as validation tooling | Docker-backed wrapper pinned by image digest; invoked through allowlisted validator IDs with server-derived artifacts. |
| GTFS Realtime protobuf Go bindings | Integrated | Used only at feed serialization boundaries. |
| Go toolchain | Integrated developer/runtime build tool | Version follows `go.mod`; used for build, test, and service commands. |
| Docker / Docker Compose | Integrated local tooling | Provisions local Postgres/PostGIS and supports the Docker-backed GTFS-RT validator wrapper; app containers are not packaged in this repo. |
| Task | Optional local tooling | Mirrors Make targets; Makefile remains independently supported. |
| Local demo tools `curl`, `zip`, `unzip` | Integrated local demo tooling | Used only by `scripts/demo-agency-flow.sh`; not runtime service dependencies. |
| Phase 17 pilot operations helpers | Integrated deployment tooling | `scripts/pilot-ops.sh` provides dry-run-capable validator, backup, restore-drill, feed-monitor, and scorecard-export helpers; systemd timer examples are under `deploy/systemd/`. |
| Internal Prometheus-format `/metrics` | Partially wired | Services can expose internal metrics text when `METRICS_ENABLED=true`; no Prometheus server, Grafana dashboard, alert rules, or SLO deployment assets are integrated. |
| Prometheus / Grafana | Deferred optional integration | Future deployment/observability stack only. |
| OpenTelemetry | Deferred optional integration | Phase 11 repo scan found no OpenTelemetry SDK, collector, exporter, trace propagation, or deployment wiring. |
| TheTransitClock | Deferred optional predictor | Not integrated. Future use must be behind `internal/prediction.Adapter`; Open Transit RT remains source of truth. |
| Other external predictors | Deferred optional predictors | Same adapter boundary as TheTransitClock. |
| AVL / vendor adapter pilot | Synthetic dry-run transform only | Phase 29B adds `internal/avladapter` and `cmd/avl-vendor-adapter` for synthetic fixture transforms into the existing telemetry contract. No named vendor, network send mode, credential, runtime dependency, or real vendor compatibility claim is added. |
| Google Maps, Apple Maps, Transit App, Bing Maps, Moovit | Workflow records and Phase 13 evidence docs only | Default `consumer_ingestion` records can track submission status; Phase 13 docs provide current records and templates. No external API calls or acceptance proof. |
| Mobility Database, transit.land | Workflow targets and Phase 13 evidence docs only | Documented as possible publication/aggregator targets; Phase 13 docs provide current records and templates. No API integration or acceptance proof. |

---

## 1. Postgres

### Classification
Core runtime infrastructure

### Purpose
Primary relational datastore for:
- agency metadata
- GTFS published feed versions
- GTFS draft data
- telemetry
- assignments
- audit logs
- incidents
- feed metadata

### Expected version
- PostgreSQL 16 preferred for local/containerized and larger production deployments
- PostgreSQL 15 is accepted for the Oracle Linux 9 OCI micro pilot because PGDG packages and PostGIS are available there and the host has about 503 MiB usable RAM

### Integration boundary
- Accessed through Go repository interfaces
- No handler or public-feed code should query SQL directly without going through the data layer
- All schema evolution should go through migrations

### Startup / provisioning
- local dev: Docker or docker-compose
- default local Docker host port: `55432`, mapped to container port `5432`
- CI: disposable service container
- production: managed or self-hosted PostgreSQL

### Failure behavior
- services that require DB access should fail fast on startup if DB is unavailable
- health endpoints should clearly indicate DB unavailability
- public feed generation should not silently serve fake or stale demo data because DB is down

### Replacement strategy
- no replacement planned
- if changed in future, preserve repository interface contracts

---

## 2. PostGIS

### Classification
Core runtime infrastructure

### Purpose
Spatial extension for:
- stop point indexing
- shape point and shape-line storage
- nearest-shape matching
- geometry projections
- efficient spatial queries for trip matching

### Expected version
- compatible with the active PostgreSQL deployment; the OCI micro pilot uses PostGIS 3.4 for PostgreSQL 15 from PGDG

### Integration boundary
- accessed only from repository / matching packages
- do not allow PostGIS-specific SQL to spread across unrelated packages
- keep geometry operations encapsulated behind matching or data-layer helpers

### Failure behavior
- if PostGIS is unavailable or not installed, spatial matching features should fail clearly
- application should not silently downgrade to broken or misleading shape matching

### Replacement strategy
- possible future replacement with non-DB spatial indexing, but not planned
- preserve matching interfaces so internal storage can evolve later

---

## 2A. pgx

### Classification
Core runtime infrastructure

### Purpose
Go PostgreSQL driver and connection pooling layer for repository implementations.

### Expected version
- `github.com/jackc/pgx/v5` v5.x

### Startup / provisioning
- linked through Go modules
- used by `cmd/telemetry-ingest` and repository implementations

### Integration boundary
- database access should go through repository packages
- handlers and feed publishers should not embed SQL directly
- core domain structs should not depend on pgx-specific types

### Input contract
- SQL queries and commands from repository implementations
- `DATABASE_URL` or test database connection strings from configuration

### Output contract
- context-aware query results mapped into internal domain models
- explicit errors returned to callers
- telemetry ingest readiness checks through `pgxpool.Ping`

### Failure behavior
- connection failures should fail startup for DB-required services
- transient query failures should be surfaced; services must not serve fake data

### Replacement strategy
- another Go PostgreSQL driver may replace pgx if repository interfaces remain stable

---

## 2B. Goose

### Classification
Validation and developer tooling

### Purpose
Versioned SQL migration runner for local dev, CI, and production deployment workflows.

### Expected version
- `github.com/pressly/goose/v3` v3.x

### Startup / provisioning
- invoked through `cmd/migrate`
- driven by `DATABASE_URL` and `MIGRATIONS_DIR`

### Integration boundary
- only migration commands should import Goose
- application services should not depend on Goose APIs at runtime

### Input contract
- SQL migration files under `db/migrations`
- command: `up`, `down`, `status`, or `redo`

### Output contract
- migration application status in the Goose schema table
- process exit code indicating success or failure

### Failure behavior
- failed migration commands must exit non-zero and leave the DB transactionally consistent where supported
- bootstrap must stop if migrations fail

### Replacement strategy
- another migration tool may replace Goose if `cmd/migrate` command behavior remains compatible

---

## 3. GTFS static validator

### Classification
Validation and developer tooling

### Purpose
Validate static GTFS before publish and during compliance checks.

### Preferred tooling
- MobilityData GTFS Validator or equivalent canonical validator

### Integration boundary
- invoked by import/publish workflows, admin validation runs, or CI checks
- Post-Phase-8 hardening replaces request-supplied commands with server-side allowlisted validator IDs
- static validation uses `validator_id=static-mobilitydata`
- local/prod config supplies `GTFS_VALIDATOR_PATH`; if the path ends with `.jar`, the adapter runs `java -jar` through argv-based execution
- Java 17 or newer must be installed and runnable as `java` or configured through `JAVA_BINARY`; `make validators-check` also probes common Homebrew Java locations on macOS and fails when no Java runtime is available
- expected canonical version: MobilityData GTFS Validator `v7.1.0`
- repo-supported installation is `make validators-install`, which downloads `gtfs-validator-7.1.0-cli.jar` into `.cache/validators/` and verifies SHA-256 `52c2785089aaf04e7ba1bb11b2db215692e2622eb0e196b823c194d156d9b58c` from `tools/validators/validators.lock.json`
- CI and production setup should run `make validators-install validators-check` or use a prebuilt runner image whose installed validator paths match `tools/validators/validators.lock.json`
- validation results stored as normalized `validation_report` rows
- validator output should not dictate internal schema design
- Phase 4 implements an internal GTFS import validator/report contract for required files, supported route type ranges, numeric ranges, core references, usable service sources, shape ordering, stop_times references, times beyond `24:00:00`, frequencies, and block preservation. This internal validator is not a substitute for canonical compliance validation.
- validators execute with `exec.CommandContext(binary, args...)`, never `/bin/sh -c`
- admin validation requests may provide only `validator_id`, `feed_type`, and optional `feed_version_id`; the server derives local feed artifacts itself

### Failure behavior
- failed validation should block publish or mark the import unhealthy based on configured strictness
- validation reports must remain visible to operators
- missing validator configuration stores `status='not_run'`; in production scorecards this is red, and in dev scorecards this is yellow
- Phase 4 internal validation failures block activation, store `gtfs_import` and `validation_report` rows when the report write succeeds, and leave `gtfs_import.feed_version_id` `NULL`.
- Phase 4 publish/database failures roll back staged rows, update `gtfs_import.report_json`, and store a failed `validation_report` outside the rolled-back publish transaction when possible.
- If the best-effort failed-import report write also fails, the importer/CLI returns a clear error and must not claim that failure metadata was stored.

### Replacement strategy
- validator implementation can change
- validation report schema and publish gating behavior should remain stable

---

## 4. GTFS Realtime validator

### Classification
Validation and developer tooling

### Purpose
Validate GTFS-RT feeds:
- Vehicle Positions
- Trip Updates
- Alerts

### Preferred tooling
- MobilityData GTFS Realtime validator or equivalent

### Integration boundary
- invoked during CI, smoke tests, scheduled runtime validation, or admin validation runs
- Post-Phase-8 hardening uses `validator_id=realtime-mobilitydata` with server-owned config from `GTFS_RT_VALIDATOR_PATH` and `GTFS_RT_VALIDATOR_ARGS`
- realtime validation should prefer generated local feed bytes/temp files over internal feed URLs whenever the service can build the artifact locally
- `cmd/agency-config` now prefers internal builder-derived protobuf bytes for Vehicle Positions, Trip Updates, and Alerts, then writes those bytes to local temp files before validator execution
- configured feed URLs (`VEHICLE_POSITIONS_FEED_URL`, `TRIP_UPDATES_FEED_URL`, `ALERTS_FEED_URL`, or `REALTIME_VALIDATION_BASE_URL`/`FEED_BASE_URL`) are a fallback only when internal builders cannot be constructed in that runtime context
- the server-owned args may use `{schedule_zip}`, `{realtime_pb}`, `{feed_type}`, and `{output_dir}` placeholders
- repo-supported GTFS-RT installation is Docker-backed: `make validators-install` pulls `ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89` and writes `.cache/validators/gtfs-rt-validator-wrapper.sh`
- the pinned GTFS-RT image is a webapp, so the repo wrapper starts the container, serves the server-derived local schedule/realtime artifacts through a temporary local HTTP server, calls the validator webapp API, and normalizes the resulting monitor data into JSON counts
- `GTFS_RT_VALIDATOR_READY_TIMEOUT_SECONDS` may be set for low-memory hosted deployments where the pinned validator webapp needs longer than the default 20 seconds to become ready after cold start
- the repo-supported GTFS-RT wrapper requires Docker, `curl`, and `python3`; `make validators-check` verifies those runtime dependencies and rejects stale wrappers that do not call the webapp API
- `GTFS_RT_VALIDATOR_PATH` should point to that wrapper for the repo-supported pinned workflow; a direct non-Docker executable is runtime-capable, but `make validators-check` intentionally does not accept it as pinned proof unless this document and `tools/validators/validators.lock.json` are extended with an equivalent checksum/digest contract
- `VALIDATOR_TOOLING_MODE=stub` is the explicit deterministic stub bypass for targeted tests or smoke runs that intentionally do not use pinned canonical tooling
- CI/local/prod must pin the selected MobilityData GTFS Realtime validator distribution by immutable package digest before making compliance claims
- output stored as normalized `validation_report` rows
- does not own business logic; it verifies it

### Failure behavior
- validation failure should mark a feed unhealthy
- unhealthy state must be visible in monitoring and admin views
- missing validator binary configuration stores `status='not_run'`; in production scorecards this is red, and in dev scorecards this is yellow
- request-provided shell text, paths, argv, output directories, and validator URLs are not accepted
- validators run with timeout, stdout/stderr caps, report-size caps, and temp-file/output confinement
- `make validate` and `make smoke` report missing pinned tooling separately from checksum/digest/path misconfiguration through `scripts/check-validators.sh`

### Replacement strategy
- validator engine may be swapped
- internal compliance dashboard should not depend on a specific validator output format

---

## 5. GTFS Realtime protobuf tooling

### Classification
Core runtime + developer tooling

### Purpose
Generate and serialize official GTFS Realtime `FeedMessage` protobuf payloads.

### Preferred implementation
- `github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs` v1.0.0
- `google.golang.org/protobuf` v1.36.11

### Integration boundary
- protobuf types should be used at the feed boundary
- internal domain models should remain separate from raw protobuf types where practical
- mapping from internal models to protobuf should happen in feed publisher packages
- Phase 3 keeps protobuf types inside `internal/feed`; telemetry and matcher/domain packages do not expose GTFS-RT types

### Failure behavior
- serialization errors must fail request generation clearly
- do not fall back to placeholder feed data
- empty Vehicle Positions feeds must still return valid protobuf `FeedMessage` responses with populated headers

### Replacement strategy
- protobuf schema version changes should be isolated to feed-mapping packages

---

## 6. TheTransitClock

### Classification
Candidate-only optional prediction backend

### Purpose
Potential external backend for Trip Updates generation and ETA prediction.

### Current Phase 29A status
Candidate-only. TheTransitClock is not vendored, not installed, not required for tests, not required at runtime, and not configured by any production runtime toggle in this repository. Phase 29A reviewed public project sources only and added mock/test-only adapter contract checks; it did not integrate or execute TheTransitClock.

The public repository is GPL-3.0 licensed. Vendoring, linking against, distributing, or tightly packaging TheTransitClock code requires explicit maintainer and license review before any future implementation. A later runtime adapter should prefer a process-level or network-level boundary unless maintainers explicitly approve another shape.

### Important architectural rule
TheTransitClock is **not** the source of truth for:
- telemetry ingestion
- agency configuration
- GTFS source management
- vehicle-trip assignment persistence

If used, it is an optional **prediction engine** behind the prediction adapter.

### Expected role
- consume GTFS and/or Vehicle Positions input from Open Transit RT
- produce Trip Updates output or prediction diagnostics

### Integration boundary
Define a narrow adapter contract:

#### Input
- active published GTFS feed version
- current vehicle assignments
- current telemetry-derived vehicle positions
- canonical Vehicle Positions feed URL or equivalent feed data

#### Output
- Trip Updates feed
- optional diagnostics or health metadata

### Required adapter rules
- no core domain package may depend directly on TheTransitClock internal classes or internal data model
- all integration goes through `prediction-adapter`
- TheTransitClock-specific configuration must live in dependency/config docs and adapter packages only

### Failure behavior
If TheTransitClock is unavailable:
- Vehicle Positions publishing must continue
- telemetry ingest must continue
- assignment persistence must continue
- Trip Updates endpoint may degrade or report unavailable
- admin and monitoring must show degraded prediction status
- no corruption of core internal state is allowed

Phase 29A does not prove runtime compatibility, better ETAs, production-grade ETA quality, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or vendor equivalence.

### Replacement strategy
The codebase must be able to replace TheTransitClock with:
- internal deterministic ETA engine
- another external predictor
- future ML-based ETA predictor

The public Trip Updates endpoint and internal prediction adapter interface must remain stable across replacement.

---

## 7. Other future prediction engines

### Classification
Optional prediction backend

### Purpose
Allow future support for:
- internal ETA engine
- ML-assisted ETA engine
- alternate open-source predictor

### Integration boundary
Same adapter boundary as TheTransitClock.

### Rule
Do not design internal telemetry, matching, or GTFS storage around one predictor’s assumptions unless the repository docs are updated and the architectural decision is explicitly recorded.

---

## 7A. Phase 6 no-op Trip Updates adapter

### Classification
Core architecture boundary / default runtime behavior

### Purpose
Provide an explicit, safe default Trip Updates adapter while the first real prediction backend remains undecided.

### Integration boundary
- implemented through `internal/prediction.Adapter`
- consumed only by the Trip Updates feed service
- accepts active published GTFS, persisted latest telemetry, persisted current assignments, and the Vehicle Positions feed URL
- returns Trip Updates plus diagnostics
- does not own telemetry ingest, matching, Vehicle Positions, GTFS import, or GTFS Studio state

### Failure behavior
- the no-op adapter returns a valid empty Trip Updates feed with diagnostics status `noop`
- missing active GTFS produces a valid empty Trip Updates feed with explicit diagnostics
- adapter errors produce valid empty Trip Updates feeds with error diagnostics
- Vehicle Positions, telemetry ingest, assignment persistence, and Studio continue independently

### Replacement strategy
Replace the no-op adapter with an internal deterministic predictor, TheTransitClock adapter, or another predictor behind the same `internal/prediction.Adapter` contract.

---

## 7B. Phase 7 internal deterministic Trip Updates adapter

### Classification
Core internal prediction backend

### Purpose
Provide the first real Trip Updates behavior behind `internal/prediction.Adapter` without taking a dependency on an external predictor.

### Integration boundary
- implemented in Go inside `internal/prediction`
- consumed only by `cmd/feed-trip-updates` through `internal/feed/tripupdates`
- reads active published GTFS through `gtfs.Repository`
- reads prediction-affecting overrides and writes review items through `prediction.OperationsRepository`
- does not write telemetry, assignments, Vehicle Positions, or GTFS Studio data
- does not expose GTFS-RT protobuf types outside feed boundary packages

### Input contract
- active published GTFS feed version
- latest accepted telemetry snapshot
- current persisted vehicle-trip assignments
- active prediction-affecting overrides
- configured Vehicle Positions feed URL

### Output contract
- conservative Trip Updates for defensible in-service trip instances
- conservative `CANCELED` Trip Updates for active canceled-trip overrides
- prediction diagnostics, coverage metrics, withheld reasons, and review items

### Failure behavior
- weak, stale, degraded, deadhead, layover, ambiguous, unsupported added-trip, unsupported short-turn, and unsupported detour inputs are withheld
- canceled trips are excluded from ETA coverage metrics and tracked separately
- missing Alerts authoring for cancellations is persisted as `expected_alert_missing=true`
- review item persistence failure is reported in diagnostics without corrupting feed output
- adapter construction or core schedule-query failures surface as Trip Updates adapter errors

### Replacement strategy
The adapter can be replaced by TheTransitClock, another external predictor, or a later higher-quality internal ETA engine if the `internal/prediction.Adapter` and operations repository contracts remain stable.

---

## 7C. Phase 29B synthetic AVL / vendor adapter pilot

### Classification
Synthetic developer/test utility and adapter-boundary example

### Purpose
Demonstrate how a deployment-owned vendor/AVL adapter can transform external-looking payloads into the existing Open Transit RT telemetry event contract before any later private integration calls `/v1/telemetry`.

### Current status
Integrated as dry-run-only Go code:
- `internal/avladapter`
- `cmd/avl-vendor-adapter`
- synthetic fixtures under `testdata/avl-vendor/`

No named vendor, real AVL feed, credential, endpoint URL, network send mode, or runtime external dependency is integrated.

### Integration boundary
- Inputs are synthetic vendor payload fixtures and a strict synthetic mapping file.
- The mapping file is the authority for emitted `agency_id`, `device_id`, and `vehicle_id`; vendor payload IDs are lookup keys only.
- Outputs are transformed `telemetry.Event` JSON records that satisfy the existing telemetry contract.
- Diagnostics are stable JSON-array dry-run review output, not telemetry ingest status.

### Failure behavior
- Hard errors reject invalid mapping rows, source mismatches, unknown vendor mappings, malformed payloads, missing coordinates, invalid coordinates, or transformed records that do not satisfy `telemetry.Event.Valid()`.
- Warnings label stale/future timestamps, low GPS accuracy, duplicate dry-run observations, and out-of-order dry-run observations.
- In mixed batches, valid records may still print to stdout while the command exits nonzero for hard errors. That output is dry-run transform output only.

### Replacement strategy
Later real adapters may be agency-owned scripts, sidecars, vendor-owned middleware, or private integration processes if they preserve the `/v1/telemetry` contract and keep vendor credentials outside the public repo.

---

## 8. Go toolchain

### Classification
Validation and developer tooling

### Purpose
Compile, format, test, and run the repository.

### Expected version
- match `go.mod`

### Integration boundary
- standard build/test pipeline
- commands documented in README, Makefile, or Taskfile

### Failure behavior
- if Go is missing from PATH in a shell or CI context, the failure should be explicit
- do not assume environment setup that is undocumented

### Replacement strategy
- none planned

---

## 9. Docker / docker-compose

### Classification
Developer tooling / optional runtime support

### Purpose
Local bootstrapping for:
- Postgres
- PostGIS-enabled DB image if used
- optional validation tools
- optional prediction backend in local development

### Integration boundary
- local orchestration only
- should not become the only documented deployment path
- current local compose file maps PostGIS to host port `55432`

### Failure behavior
- bootstrap scripts should clearly report if required services did not start

### Replacement strategy
- could be replaced by another local dev orchestration method
- keep service env vars and startup assumptions documented

---

## 9A. Task

### Classification
Developer tooling

### Purpose
Optional task runner for local workflows.

### Expected version
- Task v3 compatible syntax

### Startup / provisioning
- optional local install
- `Taskfile.yml` mirrors Makefile workflows

### Integration boundary
- Task is not required for production runtime
- Makefile must remain independently usable when Task is absent

### Input contract
- task names such as `db:up`, `migrate:up`, `test`, and `validate`

### Output contract
- shell commands with the same semantics as the corresponding Makefile targets

### Failure behavior
- if Task is missing, use Makefile targets directly

### Replacement strategy
- Makefile is the fallback and can become the only workflow surface if Task is removed

---

## 9B. GTFS / GTFS-Realtime protobuf and validation tooling

### Classification
Validation and developer tooling; GTFS-RT Vehicle Positions protobuf serialization is core runtime as of Phase 3.

### Purpose
Phase 3 uses official GTFS-Realtime protobuf bindings for Vehicle Positions feed generation. Phase 8 adds canonical validator command adapters for feed validation and compliance checks.

### Expected version
- `github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs` v1.0.0
- `google.golang.org/protobuf` v1.36.11
- MobilityData validators or equivalent canonical validators

### Startup / provisioning
- protobuf bindings are pulled through Go modules
- canonical validator execution runs when `GTFS_VALIDATOR_PATH` or `GTFS_RT_VALIDATOR_PATH` is configured
- canonical validator setup is repo-supported through `make validators-install` and `make validators-check`
- `GTFS_VALIDATOR_PATH` points to the pinned MobilityData static-validator JAR; `GTFS_RT_VALIDATOR_PATH` points to the pinned Docker-backed GTFS-RT wrapper for the repo-supported path

### Integration boundary
- protobuf types may appear in feed boundary packages only
- validators run behind validation adapters and store normalized reports
- validators do not own internal schedule, assignment, or prediction state

### Input contract
- static GTFS ZIP or generated GTFS-RT protobuf payloads

### Output contract
- serialized protobuf feeds from feed publisher packages
- normalized validation reports with errors, warnings, and info counts

### Failure behavior
- serialization failure should fail feed generation clearly
- validation failures should mark feeds unhealthy or block publish when configured

### Replacement strategy
- protobuf version changes should be isolated to feed mapping packages
- validator implementations may change if normalized report contracts remain stable

## 9C. GTFS Studio server-rendered UI

### Classification
Core admin/runtime surface

### Purpose
Phase 5 adds a minimal operator/admin UI for GTFS Studio draft editing and draft publish.

### Expected version
- Go standard library `net/http` and `html/template`

### Startup / provisioning
- command: `go run ./cmd/gtfs-studio`
- Makefile target: `make run-gtfs-studio`
- default local port: `8086` when `PORT=8086`

### Integration boundary
- uses `internal/gtfs.DraftService`
- reads and writes typed draft GTFS tables only for editing
- publishes by converting draft rows into the internal GTFS feed model and calling the shared Phase 4 validation/activation helper directly
- does not introduce a frontend build pipeline

### Failure behavior
- DB unavailability is reported through `/readyz`
- non-editable drafts are rejected before draft-to-feed conversion or publish activation
- validation failures create no active feed version and keep editable drafts editable

### Replacement strategy
- a richer admin frontend can replace the HTML surface later if it preserves `DraftService` and publish contracts

---

## 9D. Local demo packaging tools

### Classification
Developer tooling

### Purpose
Support the Phase 10 local agency demo wrapper:
- package `testdata/gtfs/valid-small` into a runtime GTFS ZIP
- verify the fetched public `schedule.zip`
- fetch local HTTP endpoints

### Expected tools
- `zip`
- `unzip`
- `curl`

### Startup / provisioning
- these are invoked by `scripts/demo-agency-flow.sh`
- they are not runtime service dependencies

### Integration boundary
- used only for local/demo packaging and verification
- the production GTFS import path remains `cmd/gtfs-import -zip <gtfs.zip>`
- public feed serving and validation do not depend on these tools

### Failure behavior
- the demo wrapper fails fast if one of these tools is missing
- missing local demo tools must not be interpreted as application runtime failure

### Replacement strategy
- the demo can replace these shell tools with a Go helper later if needed, without changing service contracts

---

## 10. Prometheus / Grafana

### Classification
Future optional integrations

### Current wiring status
- Open Transit RT has an internal Prometheus-format `/metrics` endpoint when `METRICS_ENABLED=true`.
- The repo does not provision a Prometheus server, Grafana dashboards, alert rules, remote write, retention, uptime checks, or production SLO evidence.
- Therefore Prometheus/Grafana are deferred deployment integrations, not current repo-integrated systems.

### Purpose
Observability stack for:
- feed freshness
- assignment quality
- endpoint availability
- validation trends

### Integration boundary
- metrics exposed from services via HTTP
- dashboards optional and external

### Failure behavior
- app should continue if metrics backend is absent
- loss of metrics sink must not break request-path behavior

### Replacement strategy
- any metrics backend is acceptable if service-level metrics contracts remain stable

---

## 10A. OpenTelemetry

### Classification
Future optional integration

### Current wiring status
- Not integrated.
- Phase 11 repo search found no OpenTelemetry SDK, collector, exporter, trace propagation, sampling configuration, or deployment documentation.
- Request IDs and structured logs exist, but they are not OpenTelemetry tracing.

### Purpose
Potential future observability layer for:
- distributed traces across Go services
- structured span attributes for feed generation, validation, ingest, and database operations
- trace export to a collector or hosted observability backend

### Integration boundary
- OpenTelemetry instrumentation should stay in service/server middleware, repository wrappers, or observability packages.
- Core transit domain models must not depend on OpenTelemetry types.
- Exporter configuration must be optional and deployment-owned.

### Failure behavior
- If OpenTelemetry export is unavailable, request handling and feed generation must continue.
- Telemetry export failures must not corrupt application state or block public feed responses.

### Replacement strategy
- Any tracing backend is acceptable if request/operation context remains decoupled from core domain contracts.

---

## 11. Consumer submission targets

### Classification
Workflow and compliance metadata

### Purpose
Operational workflows for:
- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

### Integration boundary
- these are not runtime dependencies
- they are workflow and compliance dependencies
- Phase 8 tracks submission status and packet JSON in `consumer_ingestion`
- default seeded consumer records are Google Maps, Apple Maps, Transit App, Bing Maps, and Moovit
- Mobility Database and transit.land are documented workflow targets; they are not seeded defaults, but operators can track them as `consumer_ingestion.consumer_name` values when needed
- Phase 13 adds documentation-only evidence records and templates under `docs/evidence/consumer-submissions/`
- validator success and public fetch proof are supporting evidence only, not consumer acceptance
- no external submission API is called directly by the app

### Failure behavior
- failed submission or rejection must not break feed generation
- status must remain visible to operators

### Replacement strategy
- none needed; this is workflow metadata rather than runtime coupling

---

## External codebase compatibility rule

If an external codebase does not fit cleanly with Open Transit RT:

1. do not force its internal assumptions into the core model
2. isolate it behind an adapter
3. document the mismatch
4. preserve Open Transit RT’s internal source-of-truth boundaries
5. prefer disabling the integration over introducing hidden coupling

---

## Current required additions to the repo

The following repo artifacts should exist to make dependency handling explicit:

- `.env.example`
- `docs/decisions.md`
- `docs/dependencies.md`
- migration command
- local bootstrap workflow
- integration fixtures under `testdata/`

---

## Source-of-truth summary

### Open Transit RT owns
- GTFS import and GTFS Studio draft/publish model
- telemetry ingest
- vehicle state and assignment persistence
- public Vehicle Positions feed
- operator workflows
- audit logs
- compliance tracking

### External predictors may own
- ETA generation
- Trip Updates generation logic, if configured behind the adapter

### Validators own
- validation checks only

### Database owns
- durable storage only

This separation must be preserved as the repository evolves.

---

## Local Docker app package

### Classification
Local development and agency-evaluation packaging

### Purpose
Phase 16 adds a local Compose `app` profile so an evaluator can start the full local stack without manually launching six Go services. It is used by:

- `scripts/agency-local-app.sh`
- `make agency-app-up`
- `make agency-app-down`
- `make agency-app-logs`
- `make agency-app-reset`

### Components
- `deploy/Dockerfile.local` builds local Go service binaries into `open-transit-rt-local:latest`.
- `deploy/docker-compose.yml --profile app` runs the Go services, GTFS Studio, Postgres/PostGIS, and a local Caddy proxy.
- `deploy/Caddyfile.local` routes local `http://localhost:8080` requests to the service containers.

### Integration boundary
- This profile is for local demo packaging only.
- It does not define the production TLS, DNS, or admin network boundary.
- It does not bake generated credentials, private keys, `.cache` validator downloads, or local env files into the image.
- Admin/debug routes may be proxied locally but still require admin auth.
- Validators remain an optional host-side setup step for startup; run `make validators-install validators-check` before validation workflows.

### Failure behavior
- `scripts/agency-local-app.sh up` waits for Postgres, migrations, service readiness, proxy health, and public feed URL fetches before reporting success.
- Docker unavailable, port conflicts, DB failures, service readiness failures, missing validators, and feed fetch failures should print next-action guidance.
- `scripts/agency-local-app.sh reset` is destructive and must state that it removes containers, the Compose volume, local demo DB state, generated local env files if present, and container logs.

---

## Phase 17 pilot operations helpers

### Classification
Deployment/operator tooling

### Purpose
Provide repeatable pilot operations for:

- scheduled canonical validation for static GTFS, Vehicle Positions, Trip Updates, and Alerts;
- Postgres backup and retention cleanup;
- restore-drill command sequence with post-restore public feed checks;
- public feed availability monitoring;
- compliance scorecard export evidence.

### Startup / provisioning
- `scripts/pilot-ops.sh` is the helper entry point.
- Example systemd service/timer files live under `deploy/systemd/`.
- Deployment-owned values should live in a private environment file such as `/opt/open-transit-rt/ops/pilot-ops.env`.

### Integration boundary
- Helpers call existing service endpoints and standard Postgres tools; they do not change backend API contracts, database schema, public feed URLs, GTFS-RT contracts, or consumer-submission statuses.
- All helpers support `--dry-run`.
- State-changing helper runs require explicit `ENVIRONMENT_NAME` and target paths/URLs.
- Restore requires explicit confirmation unless `--force` is passed.
- Systemd examples use `EnvironmentFile=` and must not inline live secrets.

### Evidence outputs
- `validator-cycle-YYYY-MM-DD.json`
- `backup-run-YYYY-MM-DD.txt`
- `restore-drill-YYYY-MM-DD.txt`
- `feed-monitor-YYYY-MM-DD.txt`
- `scorecard-export-YYYY-MM-DD.json`

Raw backups, private env files, admin tokens, database URLs with passwords, webhook URLs, notification credentials, private keys, and unredacted operator artifacts are never public evidence.

### Failure behavior
- Missing required environment variables fail clearly instead of assuming deployment defaults.
- Missing webhook/email notification destinations are recorded as `notification not configured`, not as feed failures.
- Feed monitor exits non-zero only for feed availability failures.
- Validation and scorecard helpers fail on admin/API errors.

### Replacement strategy
The helpers may be replaced by a deployment scheduler, CI runner, managed monitoring service, or external backup system if the replacement preserves evidence outputs, redaction rules, auditability, and truthfulness boundaries.
