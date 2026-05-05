# Current Status

This document is the short operational summary for the repository.

A fresh Codex instance should be able to read this file and quickly understand:
- what exists
- what does not exist
- what phase is active
- what should happen next

## Current Repository State

This repository is an early-stage starter for **Open Transit RT**.

Phase 0 scaffolding, Phase 1 durable telemetry foundation, Phase 2 deterministic trip matching, Phase 3 Vehicle Positions production feed, Phase 4 GTFS import/publish, and Phase 5 GTFS Studio draft/publish are complete. The repo can format, test, start Postgres/PostGIS, run migrations, seed local agencies, execute the bootstrap flow, import GTFS ZIP files, edit typed GTFS drafts, publish drafts, and run DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, and Trip Updates diagnostics tests.

Phase 6 Trip Updates and Alerts architecture is complete. The repo has a pluggable Trip Updates adapter boundary, default no-op adapter, Trip Updates diagnostics persistence, valid empty Trip Updates protobuf/JSON endpoints, valid empty Alerts protobuf/JSON endpoints, and non-coupling tests that keep prediction packages out of telemetry ingest, Vehicle Positions, and GTFS Studio.

Phase 7 prediction quality and operations workflows are complete for the first conservative production-directed scope. The Trip Updates service now defaults to an internal deterministic predictor behind `internal/prediction.Adapter`, emits non-empty Trip Updates for defensible matched inputs, withholds weak/degraded/deadhead/layover/disrupted cases, persists prediction review items, records audit-backed override workflow operations, emits cancellation Trip Updates with missing-alert linkage signals, and exposes first-class coverage metrics.

Phase 8 compliance and consumer workflow is complete for the first production-directed publication layer. The repo now has persisted Service Alerts authoring/lifecycle state, real GTFS-RT Alerts publication, Alerts-owned canceled-trip reconciliation, stable on-demand public static GTFS ZIP publication, `/public/feeds.json` discoverability metadata, publication/license/contact metadata workflows, consumer ingestion records, marketplace-gap records, compliance scorecard snapshots, and canonical-validator command adapters that normalize passed/warning/failed/not-run results.

Phase 9 production closure is implemented for the current repository surface. Admin and JSON debug routes require JWT/cookie admin auth with DB-backed roles, cookie admin flows require CSRF on unsafe methods, telemetry ingest requires active device Bearer tokens bound to agency/device/vehicle, validator execution uses server-side allowlisted validator IDs with argv-based execution, current assignment writes are serialized and protected by a partial unique index, and production runtime config fails fast without required secrets.

`/admin/validation/run` derives schedule and realtime artifacts itself. Schedule validation uses generated ZIP bytes; realtime validation prefers internally generated Vehicle Positions, Trip Updates, or Alerts protobuf bytes from the service builder boundary and uses configured feed URLs only as a fallback. The endpoint accepts only `validator_id`, `feed_type`, and optional `feed_version_id`; command/path/argv/output/artifact request fields are rejected.

Validator tooling now has a repo-supported pin/install/check workflow. `make validators-install` installs MobilityData GTFS Validator `v7.1.0` with SHA-256 verification and a Docker-backed GTFS-RT validator wrapper pinned to `ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89`. `make validators-check`, `make validate`, and `make smoke` distinguish missing pinned tooling from checksum/digest/path misconfiguration. `VALIDATOR_TOOLING_MODE=stub` is the explicit deterministic stub bypass for targeted tests.

Phase 10 docs, tutorials, deployment, and demo work is complete for the repository surface at that time. It filled the tutorial set under `docs/tutorials/`, added the executable `make demo-agency-flow` agency demo, updated `scripts/bootstrap-dev.sh` output for current services and protected/public surfaces, and added repo-owned docs assets under `docs/assets/`. The demo flow explicitly verifies public `schedule.zip`, `feeds.json`, public realtime protobuf feeds, protected JSON debug/admin access, and protected GTFS Studio access.

Phase 11 compliance evidence and optional external integration review is complete for the selected evidence-only path. The repo now has `docs/compliance-evidence-checklist.md`, which separates repo-proven capability, deployment/operator proof, and third-party confirmation. Dependency docs now explicitly mark wired integrations, workflow-only targets, and deferred optional systems including TheTransitClock, other external predictors, Prometheus/Grafana, OpenTelemetry, consumer submission APIs, Mobility Database, and transit.land.

Phase 12 is closed for the OCI pilot evidence scope. Step 1 (repo-side deployment evidence scaffolding), Step 2 (local demo evidence packet), Step 3 (hosted closure tooling hardening), and the hosted OCI pilot evidence packet are complete. The hosted packet lives at `docs/evidence/captured/oci-pilot/2026-04-24/` and passed `EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence`. A final current-live recheck on April 24, 2026 refreshed the packet with active `gtfs-import-3`, passed schedule/Vehicle Positions/Trip Updates/Alerts validation, and `canonical_validation_complete=true`.

Phase 13 is complete for the initial consumer-submission evidence structure. The tracker lives at `docs/evidence/consumer-submissions/README.md`, with current records and templates for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land. Phase 20 later moved all current target records to `prepared` only because complete packet drafts exist; no repo evidence currently supports submitted, under-review, accepted, rejected, or blocked claims for any target.

Phase 14 is complete for the public launch polish and repo simplification scope. The README is now a concise public front door with a short "what this is / what this is not" block, a single illustrative main visual, quick trial commands, bounded evidence links, quick-action links, and plain-language star/support wording. Public reader guides live under `wiki/`, while `docs/README.md` works as the documentation hub for public guides, practical tutorials, evidence records, architecture references, dependencies, decisions, and maintainer notes. `docs/assets/README.md` records generated-assisted visual specs plus the manual review rule for label accuracy, truthful captions, and useful alt text.

Phase 15 is complete for the targeted public repo hygiene and evidence redaction review scope. The review used `839efd6` (`Phase 14 -- Checkpoint 4 -- Security Cleanup`) as the earlier scrub baseline, reviewed changed files since that point plus tracked high-risk file patterns from `git ls-files`, inventoried committed evidence archives, added `SECURITY.md`, added `docs/evidence/redaction-policy.md`, added `docs/evidence/archive-inventory.md`, expanded `.gitignore`, removed ignored local `.DS_Store` and `.cache` secret artifacts from the working tree, and redacted unnecessary raw public client IP / instance-host detail from OCI operator evidence. The review found real secrets only in ignored local `.cache` files, not in tracked files or history for those `.cache` paths; rotation/revocation is still required before further pilot use.

Phase 16 is complete for the agency onboarding product packaging scope. The repo now has a local Compose `app` profile, `deploy/Dockerfile.local`, `deploy/Caddyfile.local`, and `scripts/agency-local-app.sh` behind `make agency-app-up`, `make agency-app-down`, `make agency-app-logs`, and `make agency-app-reset`. `make agency-app-up` starts the full local stack behind `http://localhost:8080`, applies migrations, seeds demo data, imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for readiness, verifies public feed URLs, and prints public URLs, admin/token instructions, device helper guidance, logs, validation status or next step, and a copy/paste support summary. Device onboarding is clearer through `scripts/device-onboarding.sh` for rebind, sample telemetry, dry-run, and simulator-style telemetry. The local proxy is explicitly demo-only; production still requires HTTPS/TLS and deployment-owned admin network boundaries.

Phase 17 is complete for the deployment automation and pilot operations scope. The repo now has `docs/runbooks/small-agency-pilot-operations.md` with an explicit deployment environment variable matrix, evidence output labels, and naming conventions; `scripts/pilot-ops.sh` with dry-run-capable validator-cycle, backup, restore-drill, feed-monitor, and scorecard-export helpers; and systemd timer examples for validation, backup, feed monitoring, and scorecard export. Hosted evidence refresh now ends with `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence`, and refreshed evidence is not complete unless the audit passes. The Phase 17 work does not change backend API contracts, database schema, public feed URLs, GTFS-RT contracts, consumer-submission statuses, or evidence claims.

Phase 18 is complete for the Admin UX and Agency Operations Console scope. `cmd/agency-config` now serves authenticated server-rendered operations pages under `/admin/operations` for dashboard, feed URL/validation state, telemetry freshness, device rotate/rebind, consumer evidence status, evidence links, and setup checklist views. `cmd/feed-alerts` now has `/admin/alerts/console` for simple alert listing, create/update, publish, and archive flows. GTFS Studio links back to the Operations Console, and the local app output prints `/admin/operations`. The console shows `PUBLICATION_ENVIRONMENT`/feed environment context and section last-updated timestamps where available. It does not add new public feed URLs, protobuf contracts, external consumer APIs, consumer-status claims, or production public-edge admin exposure.

Phase 19 is complete for the Realtime Quality and ETA Improvement measurement-first scope. The repo now has deterministic replay evaluation under `internal/realtimequality`, documented replay fixture schema under `testdata/replay/README.md`, baseline replay fixtures for matched, stale, ambiguous, low-confidence, manual override, canceled-trip, added-trip, short-turn, and detour cases, explicit quality metrics with denominators and `not_applicable` zero-denominator handling, regression guards that keep unknown/ambiguous/stale/withheld/degraded uncertainty visible, and authenticated Operations Console Trip Updates quality summaries from recorded `feed_health_snapshot` diagnostics. Phase 19 did not integrate TheTransitClock or another external predictor, did not claim production-grade ETA quality, and did not change public feed URLs, GTFS-RT protobuf contracts, unauthenticated surfaces, consumer statuses, or evidence claims.

Phase 20 is complete for the Consumer Submission Execution and CAL-ITP Readiness Program docs/evidence scope. The repo now has complete prepared consumer/aggregator packet drafts for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land under `docs/evidence/consumer-submissions/packets/`; a machine-readable tracker snapshot at `docs/evidence/consumer-submissions/status.json`; a California readiness summary at `docs/california-readiness-summary.md`; and a marketplace/vendor gap review at `docs/marketplace-vendor-gap-review.md`. All seven consumer records are `prepared` only. No external portal was contacted, no submission was automated, no submission path was guessed, and no submission, under-review, acceptance, rejection, consumer-ingestion, compliance, marketplace-equivalence, hosted SaaS, agency-endorsement, or production-grade ETA claim was added.

Phase 21 is complete for the Community, Governance, and Multi-Agency Scale docs/process scope. The repo now has contributor guidance, a code of conduct, GitHub issue and PR templates, governance authority docs, release process docs, support boundaries, multi-agency strategy, roadmap/status communication, and teaching visuals under `docs/assets/`. Phase 21 did not change backend behavior, API contracts, database schema, public feed URLs, consumer-submission statuses, external integrations, or evidence claims.

Track A — External Proof And Adoption is complete for the docs-only operator workflow scope. The repo now has an official submission-path verification workflow, pre-submission checklist, evidence intake/status-transition rules, README-only per-target artifact intake directories, and an agency-owned domain readiness checklist. Track A did not contact portals, automate submissions, add backend behavior, add helper scripts, change public feed URLs, change consumer statuses, or introduce submission, review, acceptance, ingestion, compliance, agency-endorsement, hosted-SaaS, vendor-equivalence, or production-grade ETA claims. All seven consumer and aggregator targets remain `prepared` only.

Track B — Agency Productization, Release, And Real-World Adoption has started. Phase 22 — Release And Distribution Hardening is complete for the docs/process scope. The repo now has a changelog, release checklist, release notes template, install/upgrade/rollback guidance, source tag and commit pinning guidance, local Docker build-from-tag guidance, evidence packet version-linkage guidance, and release validation checks that include `make realtime-quality`. Current distribution guidance supports source tags and local Docker builds only; published/versioned production Docker images remain deferred.

Phase 22 did not change backend behavior, API contracts, database schema, public feed URLs, consumer-submission statuses, external integrations, or evidence claims. Track B must preserve Track A truthfulness boundaries: consumer targets remain `prepared` only unless retained, redacted, target-originated evidence supports a target-specific status change.

Phase 23 — Agency-Owned Deployment Proof is closed as blocker-documented only because no agency-owned or agency-approved final feed root is available. No final-root evidence, validator records, or packet refreshes were collected. The DuckDNS OCI pilot remains hosted/operator pilot evidence only, not agency-owned production-domain proof.

Phase 24 — Real Agency Data Onboarding is complete for the docs/process and evidence-template scope. The repo now has a real-agency GTFS onboarding guide, GTFS validation triage guide, metadata approval checklist, publish review checklist, Phase 23-aware final public-feed review guidance, and template-only future real-agency import evidence scaffold. No real agency data, fake evidence, backend behavior, public feed URLs, consumer statuses, final-root evidence, or unsupported readiness claims were added.

Phase 25 — Device And AVL Integration Kit is complete for the docs/process and template-only evidence scope. The repo now has a telemetry API and AVL integration guide, device token lifecycle guide, vendor AVL adapter boundary guidance, simulator/no-hardware testing guidance, clock/timezone/GPS quality expectations, troubleshooting table, and template-only future device/AVL evidence scaffold. No backend API behavior, protobuf contract, prediction logic, public feed URL, consumer status, named vendor dependency, real device data, vendor payload, credential, hardware certification, fake evidence, or production AVL reliability claim was added.

Phase 26 — Admin UX Setup Wizard is complete for the server-rendered Operations Console setup checklist scope. `/admin/operations/setup` now shows a browser-guided checklist with explicit status sources for publication metadata, feed discovery, validation records, device bindings, telemetry repository state, docs/evidence tracker records, and evidence links. Admins can store publication metadata through the existing bootstrap/update repository behavior with agency ID derived from the authenticated principal, and can run validation from the browser by choosing only feed type while the server maps to allowlisted validator IDs. Browser GTFS ZIP upload and manual assignment override/review UI were intentionally deferred. Consumer packet state remains sourced from the Phase 20 docs/evidence tracker and all seven targets remain `prepared` only. Phase 26 did not change public feed URLs, GTFS-RT protobuf contracts, telemetry/device APIs, Trip Updates adapter boundaries, external integrations, consumer statuses, or evidence claims.

Phase 27 — Multi-Agency Isolation Prototype is complete for the test-and-documentation scope. The repo now has synthetic multi-agency fixture notes under `testdata/multi-agency/` and focused tests for DB-backed role loading, protected admin agency conflicts, Operations Console data views, GTFS Studio draft boundaries, Alerts admin/console boundaries, device credential bindings, telemetry ingest/debug listings, compliance publication/validation/scorecard/consumer records, prediction operations/audit rows, and protected realtime JSON debug surfaces. `/public/feeds.json` is tested as query-routed by `agency_id` with omitted query defaulting to configured `AGENCY_ID`; public `schedule.zip` and GTFS-RT protobuf feeds remain service-instance scoped by configured `AGENCY_ID`. Phase 27 is repository-level isolation evidence for selected workflows, not production multi-tenant hosting, hosted SaaS availability, consumer acceptance, agency endorsement, CAL-ITP/Caltrans compliance, or production-grade ETA proof.

Phase 28 — Production Operations Hardening is complete for the docs-first operations hardening scope. The repo now has `docs/runbooks/production-operations-hardening.md`, template-only incident/response records under `docs/runbooks/templates/`, stronger backup/restore cadence and restore-event guidance, monitoring/alerting alert delivery proof pattern, validator failure response guidance, explicit capacity thresholds, secret rotation guidance, operator handover fields, evidence refresh/redaction guidance, and repeated Phase 27 operations-boundary language for deployment/DB-scoped backup/restore/export/evidence workflows. Phase 28 did not change runtime APIs, database schema, public feed URLs, GTFS-RT contracts, consumer statuses, external integrations, systemd/Docker behavior, or evidence claims.

Phase 29 — Real-World Realtime Quality Expansion is complete for the synthetic replay evidence expansion scope. The repo now has richer deterministic replay fixtures for after-midnight service, exact and non-exact frequency trips, block continuity, long layover withholding, sparse telemetry, noisy/off-shape GPS, stale/ambiguous hard patterns, cancellation alert linkage, and manual override before/after expiry. Replay fixture support now includes `frequencies` and optional manual override `expires_at`; replay telemetry snapshots now use latest telemetry per vehicle; replay comparisons now assert already-recorded cancellation alert linkage and unsupported disruption-withheld metrics. Phase 29 expands synthetic replay coverage only. It does not add real-world observed-arrival/departure evidence, real-world ETA accuracy evidence, real route/time-period quality metrics, external predictors, Operations Console changes, public feed URL changes, GTFS-RT contract changes, consumer-status changes, auth-boundary changes, dependency changes, or stronger evidence claims.

Phase 29A — External Predictor Adapter Evaluation is complete for the adapter contract and candidate-only feasibility scope. The repo now documents the external predictor adapter contract, TheTransitClock candidate-only review, Vehicle Positions independence, timeout/failure semantics, and strict wrong-agency/wrong-feed output handling. Trip Updates builder output validation now rejects unsafe adapter output before protobuf serialization, and test-only mock adapter coverage verifies happy-path normalization/diagnostics persistence plus rejection of missing active-feed trips, impossible stops, stale timestamps, wrong agency/feed candidates, unsupported added-trip predictions, and low or missing confidence. Phase 29A does not implement external predictor runtime integration, add runtime external predictor config, start or call TheTransitClock, change public feed URLs, change GTFS-RT protobuf contracts, change consumer statuses, change auth boundaries, add migrations, add runtime dependencies, or support stronger ETA/compliance/vendor-support claims.

Phase 29B — AVL / Vendor Adapter Pilot Implementation is complete for the synthetic dry-run adapter pilot scope. The repo now has a strict mapping-driven `internal/avladapter` package, dry-run-only `cmd/avl-vendor-adapter` CLI, synthetic fixtures under `testdata/avl-vendor/`, stable JSON diagnostics, focused adapter/CLI tests, and updated device/AVL evidence guidance. Phase 29B does not add network send mode, named vendor runtime dependencies, real vendor payloads, credentials, telemetry/device API changes, public feed URL changes, GTFS-RT contract changes, Trip Updates behavior changes, consumer-status changes, or stronger vendor/reliability/compliance claims.

Phase 30 — Consumer Submission Execution is closed as Outcome B — blocker-documented closure only. No authorized submission, official-path verification evidence, or target-originated artifact was available. No Phase 30 target was selected, no external portal was contacted, no submission was automated, no submission path was guessed, and no artifact was added. This is a phase-level blocker-documented closure only; no individual target status changed to `blocked` because no target-specific blocker artifact exists. `docs/evidence/consumer-submissions/status.json` and all current target records were left unchanged, artifact directories remain README-only, and tracker/status consistency still shows all seven targets `prepared`.

## What Exists Now

### Repo guidance and architecture docs
The repo has:
- `AGENTS.md`
- `CHANGELOG.md`
- `CONTRIBUTING.md`
- `CODE_OF_CONDUCT.md`
- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/README.md`
- `wiki/README.md`
- `docs/dependencies.md`
- `docs/governance.md`
- `docs/release-process.md`
- `docs/release-checklist.md`
- `docs/upgrade-and-rollback.md`
- `docs/release-notes-template.md`
- `docs/support-boundaries.md`
- `docs/multi-agency-strategy.md`
- `docs/roadmap-status.md`
- `docs/compliance-evidence-checklist.md`
- `docs/agency-owned-domain-readiness.md`
- `docs/phase-plan.md`
- `docs/decisions.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/tutorials/`
- `docs/assets/`
- `docs/evidence/redaction-policy.md`
- `docs/evidence/archive-inventory.md`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/submission-workflow.md`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/california-readiness-summary.md`
- `docs/marketplace-vendor-gap-review.md`
- `docs/track-b-productization-roadmap.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-14.md`
- `docs/handoffs/phase-15.md`
- `docs/handoffs/phase-16.md`
- `docs/handoffs/phase-17.md`
- `docs/handoffs/phase-19.md`
- `docs/handoffs/phase-20.md`
- `docs/handoffs/phase-21.md`
- `docs/handoffs/phase-22.md`
- `docs/handoffs/phase-23.md`
- `docs/handoffs/phase-24.md`
- `docs/handoffs/phase-25.md`
- `docs/handoffs/phase-26.md`
- `docs/handoffs/phase-27.md`
- `docs/handoffs/phase-28.md`
- `docs/handoffs/phase-29.md`
- `docs/handoffs/phase-29a.md`
- `docs/handoffs/phase-29b.md`
- `docs/handoffs/phase-30.md`
- `docs/handoffs/track-a-external-proof.md`
- `docs/handoffs/track-b-roadmap.md`

### Phase 0 scaffolding
The repo now has:
- `.env.example`
- `Taskfile.yml`
- independently usable `Makefile`
- `cmd/migrate`
- versioned migrations under `db/migrations`
- PostGIS-backed Docker Compose configuration on host port `55432`
- local full-stack Compose app profile behind `make agency-app-up`
- `scripts/bootstrap-dev.sh`
- `scripts/agency-local-app.sh`
- `scripts/device-onboarding.sh`
- `scripts/pilot-ops.sh`
- deterministic fixtures under `testdata/`
- handoff template and Phase 0 handoff under `docs/handoffs/`

### Runtime code
The repo includes starter Go services for:
- `agency-config`
- `telemetry-ingest`
- `feed-vehicle-positions`
- `feed-trip-updates`
- `feed-alerts`
- `gtfs-studio`

`cmd/telemetry-ingest` persists valid telemetry to Postgres through a telemetry repository. `cmd/feed-vehicle-positions` serves DB-backed GTFS-RT Vehicle Positions protobuf and JSON debug output from persisted latest accepted telemetry plus persisted current assignments. `cmd/agency-config` serves publication, schedule ZIP, feed discovery, scorecard, validation, consumer-ingestion, and device-rebind workflows.

`cmd/gtfs-studio` serves a minimal server-rendered admin surface for typed GTFS draft editing and draft publishing. It is operational row editing, not a map editor or timetable designer.

`cmd/feed-trip-updates` serves stable Trip Updates endpoints backed by the Phase 7 deterministic prediction adapter by default, with the Phase 6 no-op adapter still selectable as a fallback. It returns valid GTFS-RT Trip Updates protobuf output, JSON diagnostics, prediction metrics, and persisted Trip Updates traceability through `feed_health_snapshot`.

`internal/realtimequality` runs deterministic replay fixtures from `testdata/replay/` to compare matcher assignments, Vehicle Positions publication decisions, Trip Updates behavior, withheld reasons, and quality metrics. `make realtime-quality` runs this focused suite.

`internal/avladapter` and `cmd/avl-vendor-adapter` provide the Phase 29B synthetic, dry-run-only vendor/AVL adapter pilot. The command transforms synthetic payloads from `testdata/avl-vendor/` into existing `telemetry.Event` JSON, prints diagnostics to stderr, and does not send telemetry.

`cmd/feed-alerts` serves DB-backed GTFS-RT Alerts protobuf and JSON output from persisted published Service Alerts. It also exposes minimal JSON admin operations for alert authoring, publish/archive lifecycle, and canceled-trip alert reconciliation.

`cmd/agency-config` now serves publication/compliance workflows: `/public/gtfs/schedule.zip`, `/public/feeds.json`, publication metadata bootstrap, compliance scorecard snapshots, consumer ingestion workflow records, and validator run records.

Admin routes derive actor and agency from auth context. Conflicting request `agency_id` fields or query params are rejected. Scorecard GET reads the latest stored snapshot; scorecard POST recomputes and stores. `/admin/devices/rebind` rotates a device token and binding with audit logging.

Public `.pb` feed endpoints remain anonymous. JSON debug endpoints such as `/public/gtfsrt/vehicle_positions.json`, `/public/gtfsrt/trip_updates.json`, `/public/gtfsrt/alerts.json`, and their `/admin/debug/...` aliases require admin read auth and share the same debug builders.

### Phase 1 telemetry foundation
The repo now has:
- `internal/db` with `pgxpool` connection setup and readiness ping support
- `internal/telemetry` repository interfaces and Postgres implementation
- DB-backed telemetry ingest in `cmd/telemetry-ingest`
- `/healthz` liveness and `/readyz` DB readiness behavior for telemetry ingest
- agency-scoped, bounded `/v1/events` debug listing
- durable parsed request payload storage in `telemetry_event.payload_json`
- active device credential verification before telemetry persistence
- opaque device tokens hashed with `DEVICE_TOKEN_PEPPER`
- device-to-agency/device/vehicle binding checks, including immediate old-token invalidation after rebinding/rotation
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
- per-agency/per-vehicle advisory locking for current assignment writes
- a partial unique index preventing duplicate active assignment rows
- `shape_dist_traveled = 0` is preserved as a valid persisted value, not collapsed to NULL
- repeated identical degraded unknown states reuse the active degraded assignment only when degraded state, reason codes, service date, and telemetry evidence match; telemetry evidence means matching `telemetry_event_id` when present, with `active_from` equality only as the no-telemetry fallback
- batched GTFS schedule detail loading for stop times, shape points, and frequencies under the existing schedule-query boundary
- a small reason-code, degraded-state, and incident taxonomy
- unit and DB-backed integration tests for matcher edge cases

`vehicle_trip_assignment.score_details_json` is intentionally loose debug JSON in Phase 2, not a stable public schema. Matcher-generated score details include `score_schema`; candidate-based details also include `trip_id`, `start_time`, and `observed_local_seconds` when resolvable. Unknown assignment rows carry `service_date` whenever agency timezone and observed timestamp can be resolved; `service_date` is nullable only for truly unresolved cases. Missing shape data uses reason code `missing_shape` and degraded state `missing_shape`. Route-hint matching is reserved for a future input expansion and is not active in Phase 2 because telemetry does not currently carry a route hint.

Phase 2 service-day resolution considers the observed agency-local date and the immediately previous local date. That supports normal same-day service and practical after-midnight GTFS times through the prior service day, but it is not a generalized multi-day lookback for very long service patterns beyond that two-service-day window.

### Phase 3 Vehicle Positions production feed
The repo now has:
- official GTFS-RT protobuf serialization through `github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs`
- `/public/gtfsrt/vehicle_positions.pb` as a stable DB-backed protobuf endpoint
- `/public/gtfsrt/vehicle_positions.json` as DB-backed JSON debug output
- `FeedHeader.gtfs_realtime_version = "2.0"`, `FULL_DATASET`, and snapshot-generated timestamps
- `Last-Modified` derived from the snapshot `generated_at` timestamp
- a single `internal/feed.VehiclePositionsSnapshot` model used by both protobuf and JSON rendering
- a hard `telemetry.Repository.ListLatestByAgency` ordering contract: latest accepted row per vehicle ordered by `observed_at DESC, id DESC`
- `state.Repository.ListCurrentAssignments` for narrow bulk active-assignment reads behind the state repository interface
- configurable vehicle cap, stale TTL, stale suppression TTL, and Vehicle Positions trip publication confidence threshold
- deterministic stale behavior: stale-but-unsuppressed vehicles remain in protobuf without trip descriptors; suppressed vehicles remain visible only in JSON debug
- normal successful empty protobuf feeds when there is no telemetry or all vehicles are suppressed
- JSON debug publication decisions for every snapshot vehicle, including telemetry age, assignment publishability, assignment/telemetry mismatch, trip descriptor publication, and the winning omission reason
- tests for protobuf validity, entity content, no telemetry, no assignments, stale/suppressed behavior, truncation, non-exact frequency mapping, true-north bearing preservation, telemetry mismatch, repository ordering, bulk assignment lookup, and handler headers/status

### Phase 4 GTFS import and publish pipeline
The repo now has:
- `cmd/gtfs-import` as a thin runtime GTFS ZIP import CLI
- `internal/gtfs.ImportService` for GTFS ZIP import, validation, report persistence, staging, and atomic activation
- internal GTFS validation for required files, usable service source availability, route type ranges, core references, service usability, shapes ordering, stop_times references, trips/routes/services consistency, frequencies, blocks, and times beyond `24:00:00`
- exact required runtime input rule: `agency.txt`, `routes.txt`, `stops.txt`, `trips.txt`, `stop_times.txt`, and at least one usable service source from `calendar.txt` or `calendar_dates.txt`
- deterministic service-source validation: usable `calendar.txt` rows must have at least one active weekday; `calendar_dates.txt`-only feeds must include at least one `exception_type=1` addition
- route type validation for the supported GTFS route type domain: base route types `0` through `7` and extended route types `100` through `1702`
- optional `shapes.txt` and `frequencies.txt` handling
- preservation of imported GTFS time text, including values beyond `24:00:00`
- preservation of `block_id` from `trips.txt` when present
- PostGIS point construction for stops and shape points, plus `gtfs_shape_line` construction from ordered shape points when a shape has at least two points
- transactional publish behavior that inserts a new staged `feed_version`, loads published GTFS rows, retires the previous active feed, and activates the new feed atomically
- failed validation behavior that stores `gtfs_import` and `validation_report` rows when possible and creates no staged `feed_version`
- publish/database failure behavior that updates `gtfs_import.report_json` and writes a failed `validation_report` outside the rolled-back publish transaction when possible
- failed publish rollback behavior that leaves no partial GTFS rows and keeps `gtfs_import.feed_version_id` `NULL`
- tests for valid import, invalid import, rollback safety, active feed switching, block visibility to downstream GTFS consumers, shape-line creation, and CLI wrapper behavior

### Phase 5 GTFS Studio draft/publish model
The repo now has:
- `cmd/gtfs-studio` with `/healthz`, `/readyz`, and `/admin/gtfs-studio` routes
- typed draft GTFS tables for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies
- explicit draft traceability fields: status, base feed version, latest publish attempt, latest published feed version, and soft-discard metadata
- `gtfs_draft_publish` attempts linked to schedule `validation_report` rows
- `internal/gtfs.DraftService` for blank draft creation, active-feed cloning, typed entity upsert/remove, soft discard, list/read behavior, and draft publish
- cloned-draft provenance through `gtfs_draft.base_feed_version_id`
- blank draft creation when no active feed exists and explicit blank draft creation when one does exist
- soft discard semantics: discarded drafts keep typed rows and history, are hidden by default, and are read-only/not publishable
- published drafts become read-only by default after successful publish
- entity remove operations affect only rows in the current editable draft and never delete previously published GTFS rows, feed versions, publish attempts, validation reports, or audit history
- draft agency metadata is one row scoped to the draft agency; on successful publish it upserts the canonical `agency` row in the publish transaction
- shared feed-version publishing used directly by both ZIP import and Studio publish; Studio does not generate or re-import a synthetic ZIP
- non-editable draft statuses are rejected before draft-to-feed conversion, validation, or shared publish activation
- minimal server-rendered forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies
- tests for draft CRUD, blank/clone behavior, draft/published separation, publish traceability, read-only published/discarded drafts, discarded list filtering, and summary version visibility

## Schema Source Of Truth

Migrations under `db/migrations` are the source of truth for executable schema changes and are applied through `cmd/migrate`.

`db/schema.sql` is deprecated as an executable schema. It is intentionally a comment-only pointer to the migrations directory and must not be edited independently.

## What Does Not Exist Yet

The following are still missing or incomplete unless a later handoff says otherwise:

- production-grade learned ETA/prediction quality
- real-world observed-arrival/departure ETA accuracy evidence
- real route/time-period realtime quality metrics
- hosted login/SSO and server-side admin JWT `jti` replay tracking
- full operator UI for manual override workflows
- production SLO dashboards and alerting beyond Phase 17 lightweight feed-monitor examples, Phase 18 operator pages, request logs, request IDs, readiness checks, and `/metrics` toggle
- OpenTelemetry tracing/exporter wiring and Prometheus/Grafana deployment assets
- external predictor adapters such as TheTransitClock
- external consumer submission API integrations
- consumer submission, review, acceptance, rejection, or blocker evidence from third parties
- agency-owned or agency-approved stable URL/domain proof for a final public feed root
- real device or vendor AVL integration evidence beyond local simulator/no-hardware examples and templates
- marketplace/vendor-equivalent service packaging and support commitments

## Current Phase

**Active phase:** Phase 30 — Consumer Submission Execution is closed as Outcome B — blocker-documented closure only. No authorized submission, official-path verification evidence, or target-originated artifact was available. Phase 31 — Agency Pilot Program Package is the recommended next implementation phase and must proceed from the prepared-only consumer state. It must not assume submission, review, acceptance, rejection, blocker, ingestion, listing, display, or adoption evidence exists. Track A — External Proof And Adoption is complete for the documented docs-only operator workflow, evidence intake, artifact-directory, and agency-domain readiness scope. Phases 12 through 30 remain closed for their documented scopes.

Phase 12 Step 1 is complete as repo docs/runbooks/evidence-template scaffolding. Phase 12 Step 2 has a partial local evidence packet under `docs/evidence/captured/local-demo/2026-04-22/`. Phase 12 hosted/operator evidence is complete for the OCI pilot under `docs/evidence/captured/oci-pilot/2026-04-24/`.

Phase 13 added documentation-only consumer submission records and templates. Phase 20 added complete prepared packet drafts and moved all seven current records to `prepared` only. Neither phase added runtime/product changes, consumer submission APIs, portal automation, or consumer acceptance claims.

Phase 14 added documentation-only public-facing polish. It did not change backend runtime behavior, API contracts, database schema, public feed URLs, external integrations, evidence claims, or consumer-submission status.

Phase 15 completed targeted public repo hygiene and evidence redaction review. Phase 16 completed local agency onboarding packaging. Phase 17 added deployment/operator automation and documentation only. Phase 18 added authenticated minimal admin UX for existing operational state. Phase 19 added replay measurement, explicit quality metrics, diagnostics, and safe Operations Console quality summaries. Phase 20 added prepared consumer packet docs, `status.json`, California readiness summary, and marketplace/vendor gap review. Phase 21 added community contribution, governance, support, release, multi-agency, roadmap/status, GitHub template, and teaching-visual documentation. It did not add hosted SaaS behavior, Kubernetes, external predictors, consumer submission APIs, public feed URL changes, protobuf changes, portal automation, guessed submission paths, backend behavior changes, or unsupported acceptance/compliance claims.

Track A added the safe operator workflow needed before real consumer adoption steps. It did not verify any target submission path, because no current official target source or operator-retained evidence was added for those paths. It did not change `docs/evidence/consumer-submissions/status.json` or any current target record beyond documentation links.

Track B added repo-native roadmap context for Phase 22 through Phase 32. Phase 22 added release and distribution hardening docs without runtime changes. Phase 23 closed as blocker-documented only because no agency-owned or agency-approved final feed root is available. No final-root evidence, validator records, or packet refreshes were collected. Phase 24 added real-agency GTFS onboarding, validation triage, metadata approval, publish review, and template-only evidence scaffolding without runtime or evidence-claim changes. Phase 25 added device/AVL telemetry onboarding, token lifecycle, vendor-boundary, simulator, troubleshooting, redaction, and template-only evidence guidance without runtime or evidence-claim changes. Phase 26 added browser-guided setup UX without changing public feeds, API contracts, consumer statuses, external integrations, or evidence claims. Phase 27 added selected repository-level multi-agency isolation tests and boundary docs without claiming production multi-tenant operations. Phase 28 added docs-first operations hardening, templates, alert delivery proof, capacity guidance, secret rotation, handover, and evidence refresh guidance without runtime or evidence-claim changes. Phase 29 added synthetic replay quality expansion without claiming real-world ETA accuracy, real route/time-period coverage, production-grade ETA quality, external predictor integration, or evidence-claim changes. Phase 29A documented and tested the external predictor adapter boundary without adding runtime external predictor integration, runtime config, external services, public feed URL changes, GTFS-RT contract changes, consumer-status changes, auth-boundary changes, schema changes, or stronger ETA/compliance/vendor-support claims. Phase 29B added a synthetic dry-run AVL/vendor adapter pilot behind the existing telemetry boundary without network send mode, real vendor data, credentials, external dependencies, public feed URL changes, consumer-status changes, API changes, or stronger vendor/reliability claims. Phase 30 closed as Outcome B — blocker-documented closure only at the phase level; no target was selected, no target-specific blocker artifact exists, no target moved to `blocked`, and all seven targets remain `prepared`. Track B must not advance consumer statuses, change public feed URLs, or introduce stronger readiness claims without the evidence required by Track A, the redaction policy, and the security policy.

The next Codex instance should start with `docs/handoffs/latest.md`.

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

## Phase 24 Closure Audit Results

Checked during Phase 24 closure:
- pre-edit `make validate`: passed.
- pre-edit `make test`: passed.
- pre-edit `make realtime-quality`: passed.
- pre-edit `make smoke`: passed.
- pre-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- pre-edit `git diff --check`: passed.
- post-edit `make validate`: passed.
- post-edit `make test`: passed.
- post-edit `make realtime-quality`: passed.
- post-edit `make smoke`: passed.
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- post-edit `git diff --check`: passed.

Phase 24 implementation results:
- added `docs/tutorials/real-agency-gtfs-onboarding.md` for real-agency GTFS intake, metadata approval, import/publish review, redaction, and Phase 23-aware final public-feed review.
- added `docs/tutorials/gtfs-validation-triage.md` for plain-language importer and validator issue triage, including when to ask for technical help.
- added template-only real-agency GTFS evidence scaffolding under `docs/evidence/real-agency-gtfs/`.
- linked the new onboarding and evidence docs from tutorial, evidence, production checklist, first-run, README, status, and handoff docs.
- did not add real GTFS data, fake validation outputs, fake approvals, fake import evidence, backend behavior, public feed URL changes, consumer status changes, final-root proof, or unsupported compliance/readiness claims.

## Phase 25 Closure Audit Results

Checked during Phase 25 closure:
- pre-edit/planning `make validate`: passed.
- pre-edit/planning `make test`: passed.
- pre-edit/planning `make realtime-quality`: passed.
- pre-edit/planning `make smoke`: passed.
- pre-edit/planning `docker compose -f deploy/docker-compose.yml config`: passed.
- pre-edit/planning `git diff --check`: passed.
- pre-edit/planning `sh -n scripts/device-onboarding.sh`: passed.
- pre-edit/planning `scripts/device-onboarding.sh help`: passed.
- pre-edit/planning `scripts/device-onboarding.sh sample --dry-run`: passed.
- pre-edit/planning `scripts/device-onboarding.sh simulate --dry-run`: passed.
- post-edit `make validate`: passed.
- post-edit `make test`: passed.
- post-edit `make realtime-quality`: passed.
- post-edit `make smoke`: passed.
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- post-edit `git diff --check`: passed.
- post-edit `sh -n scripts/device-onboarding.sh`: passed.
- post-edit `scripts/device-onboarding.sh help`: passed.
- post-edit `scripts/device-onboarding.sh sample --dry-run`: passed.
- post-edit `scripts/device-onboarding.sh simulate --dry-run`: passed.
- post-edit targeted docs secret/example scan: passed.

Phase 25 implementation results:
- added `docs/tutorials/device-avl-integration.md` for telemetry endpoint, payload fields, timestamp/GPS expectations, response behavior confirmed from code/tests, simulator/no-hardware usage, vendor AVL adapter boundaries, troubleshooting, and redaction rules.
- added `docs/tutorials/device-token-lifecycle.md` for bearer token handling, local seeded demo token behavior, rotate/rebind, one-time token display, secure storage, compromise rotation, binding rules, and operator responsibilities.
- added template-only device/AVL evidence scaffolding under `docs/evidence/device-avl/`.
- documented the vendor AVL telemetry adapter boundary in `docs/decisions.md` without adding a named vendor dependency or runtime integration.
- linked the new device/AVL docs from tutorial, evidence, production checklist, first-run, README, status, and handoff docs.
- did not change `scripts/device-onboarding.sh`, backend API behavior, protobuf contracts, prediction logic, public feed URLs, consumer statuses, dependencies, or evidence claims.
- did not add real agency device data, vendor payloads, credentials, hardware certifications, fake evidence, or production AVL reliability claims.

## Phase 23 Closure Audit Results

Checked during Phase 23 closure:
- pre-edit `make validate`: passed.
- pre-edit `make test`: passed.
- pre-edit `make realtime-quality`: passed.
- pre-edit `make smoke`: passed.
- pre-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- pre-edit `git diff --check`: passed.
- post-edit `python3 -m json.tool docs/evidence/consumer-submissions/status.json`: passed.
- post-edit tracker/status consistency check: passed for target name, status, packet path, prepared timestamp, and evidence references.
- post-edit `make validate`: passed.
- post-edit `make test`: passed.
- post-edit `make realtime-quality`: passed.
- post-edit `make smoke`: passed.
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- post-edit `git diff --check`: passed.

Phase 23 implementation results:
- closed Phase 23 as blocker-documented only because no agency-owned or agency-approved final feed root is available.
- added a Phase 23 blocker record and future operator next-actions checklist for agency-owned domain proof.
- kept DuckDNS labeled as hosted/operator pilot evidence only.
- did not create final-root evidence, validator records, evidence packets, migration proof, packet refreshes, consumer status changes, or unsupported readiness/compliance claims.
- did not run `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` because no final-root evidence packet was created.

## Phase 22 Closure Audit Results

Checked during Phase 22 closure:
- pre-edit `make validate`: passed.
- pre-edit `make test`: passed.
- pre-edit `git diff --check`: passed.
- post-edit `make validate`: passed.
- post-edit `make test`: passed.
- post-edit `make realtime-quality`: passed.
- post-edit `make smoke`: passed.
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- post-edit `git diff --check`: passed.

Phase 22 implementation results:
- added `CHANGELOG.md`, `docs/release-checklist.md`, `docs/upgrade-and-rollback.md`, and `docs/release-notes-template.md`.
- expanded `docs/release-process.md` with release-from-main, tag, version verification, artifact, release note, install, upgrade, rollback, and evidence version-linkage guidance.
- documented clean install from a source tag, local app verification, local Docker image builds from tags, required release checks, backup-before-upgrade, migration run order, migration status checks, rollback limits, and restore-procedure links.
- documented version pinning by source tag, commit SHA, local Docker image tag, and artifact checksum when generated.
- explicitly deferred published/versioned production Docker images; current distribution guidance supports source tags and local Docker builds only.
- did not add backend behavior, `/version`, binary `--version`, OCI image labels, migrations, public feed URL changes, consumer status changes, external integrations, production Docker image claims, compliance claims, consumer acceptance claims, hosted-SaaS claims, or vendor-equivalence claims.

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

## Phase 3 Closure Audit Results

Checked during Phase 3 closure:
- `go mod tidy`: passed and added GTFS-RT protobuf dependencies.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed.
- `make migrate-status`: passed and reports migration versions 1, 2, and 3 applied.
- `make test-integration`: passed with DB-backed telemetry and matcher tests using isolated temporary database setup.
- `make validate`: passed Phase 3 scaffold, telemetry, matcher, and Vehicle Positions file validation only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.

Phase 3 implementation results:
- removed placeholder sample Vehicle Positions output from production paths.
- added DB-backed GTFS-RT protobuf Vehicle Positions output.
- added DB-backed JSON debug output from the same snapshot model.
- added snapshot-level cap/truncation behavior and per-vehicle publication decisions.
- preserved stale, suppressed, unknown, no-assignment, no-telemetry, manual override, non-exact frequency, and telemetry-mismatch behavior in tests.
- added official GTFS-RT protobuf Go bindings while keeping protobuf mapping inside `internal/feed`.
- did not add Trip Updates, Alerts, GTFS import, GTFS Studio, rider apps, payments, passenger accounts, CAD, or marketplace workflows.

## Phase 4 Closure Audit Results

Checked during Phase 4 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed; PostGIS container running on host port `55432`.
- `make migrate-up`: passed and applied `000004_gtfs_import_pipeline.sql`.
- `make migrate-status`: passed and reports migration versions 1, 2, 3, and 4 applied.
- migration down/up smoke for `000004_gtfs_import_pipeline.sql`: passed via `make migrate-down`, `make migrate-up`, and `make migrate-status`.
- `make test-integration`: passed with DB-backed telemetry, matcher, Vehicle Positions, and GTFS import tests using isolated temporary database setup.
- `make validate`: passed Phase 4 scaffold, telemetry, matcher, Vehicle Positions, and GTFS import file validation only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.

Phase 4 implementation results:
- added real GTFS ZIP import path through `cmd/gtfs-import` and `internal/gtfs.ImportService`.
- added durable import reports in `gtfs_import` and linked schedule validation reports.
- kept runtime import input as GTFS ZIP; directory handling exists only as test fixture setup that creates ZIPs before invoking importer behavior.
- validates required files, route types, numeric ranges, usable service source availability, core references, service usability, shapes ordering, stop_times references, trips/routes/services consistency, frequencies, agency scoping, and GTFS times beyond `24:00:00`.
- service-source validation now fully matches the Phase 4 contract: mere file or row presence is insufficient; calendar rows with no active weekdays and calendar_dates-only feeds with only removal exceptions are rejected.
- preserves canonical imported GTFS time text in published tables while using parsed seconds only for validation and query logic.
- imports optional `block_id` from `trips.txt` and proves it remains visible through the downstream GTFS repository boundary.
- creates `gtfs_shape_line` rows from ordered shape points when a shape has at least two points.
- publishes atomically by activating a new `feed_version` and retiring the previous active version in one transaction.
- failed validation creates no staged `feed_version`; publish failures roll back partial rows, leave `gtfs_import.feed_version_id` `NULL`, and persist a failed `validation_report` outside the publish transaction when possible.
- did not add GTFS Studio runtime editing, Trip Updates, Alerts, rider apps, payments, passenger accounts, CAD, or marketplace workflows.

## Phase 5 Closure Audit Results

Checked during Phase 5 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed; PostGIS container running on host port `55432`.
- `make migrate-up`: passed and applied `000005_gtfs_studio_drafts.sql`.
- `make migrate-status`: passed and reports migration versions 1, 2, 3, 4, and 5 applied.
- migration down/up smoke for `000005_gtfs_studio_drafts.sql`: passed via `make migrate-down`, `make migrate-up`, and `make migrate-status`.
- `make test-integration`: passed with DB-backed telemetry, matcher, Vehicle Positions, GTFS import, and GTFS Studio tests using isolated temporary database setup.
- `make validate`: passed Phase 5 scaffold, telemetry, matcher, Vehicle Positions, GTFS import, and GTFS Studio file validation only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.

Phase 5 implementation results:
- added typed GTFS Studio draft storage in migration `000005_gtfs_studio_drafts.sql`.
- added `internal/gtfs.DraftService` for blank drafts, active-feed clones, typed draft CRUD, soft discard, list filtering, and draft publish.
- made cloned drafts capture `base_feed_version_id`; blank drafts keep it empty.
- made discarded and published drafts read-only by default.
- made non-editable draft statuses fail before draft-to-feed conversion, validation, or shared publish activation.
- made entity remove operations delete only current editable draft rows, never published GTFS rows or publish history.
- refactored the Phase 4 publish activation into a shared helper used directly by both ZIP import and Studio publish.
- added `cmd/gtfs-studio` as a minimal server-rendered UI with draft summary version visibility and operational row forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- added DB-backed tests for blank/clone behavior, draft/published separation, direct Studio publish, traceability, read-only status behavior, and discarded-draft publish rejection.
- added handler tests for draft list filtering and draft summary version visibility.
- did not add Trip Updates, Alerts, rider apps, payments, passenger accounts, CAD, marketplace workflows, canonical validators, map editing, or timetable designer behavior.

## Phase 6 Closure Audit Results

Checked during Phase 6 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed; PostGIS container running on host port `55432`.
- `make migrate-status`: passed and reports migration versions 1, 2, 3, 4, and 5 applied.
- `make test-integration`: passed with DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, and Trip Updates diagnostics tests using isolated temporary database setup.
- `make validate`: passed Phase 6 file smoke only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.

Phase 6 implementation results:
- added `internal/prediction.Adapter` as the narrow Trip Updates prediction boundary.
- added a default no-op Trip Updates adapter that returns no Trip Updates with explicit diagnostics.
- added Trip Updates diagnostics persistence to existing `feed_health_snapshot` rows with required traceability fields.
- added `internal/feed/tripupdates` with valid empty GTFS-RT Trip Updates protobuf output by default, JSON debug output, explicit `FeedHeader.timestamp`, deterministic entity ordering, and ordered `stop_time_update` entries.
- added `cmd/feed-trip-updates` with `/healthz`, `/readyz`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/trip_updates.json`.
- added exact Vehicle Positions URL derivation: `VEHICLE_POSITIONS_FEED_URL` is an exact full URL, otherwise `FEED_BASE_URL` must include `/public` and derives `/public/gtfsrt/vehicle_positions.pb`.
- added `internal/feed/alerts` and `cmd/feed-alerts` with valid empty GTFS-RT Alerts protobuf output and JSON-only deferred diagnostics.
- added non-coupling tests proving telemetry ingest, Vehicle Positions, and GTFS Studio do not depend on prediction or Trip Updates packages.
- did not add ETA-quality logic, production predictor behavior, alert authoring, alert persistence, incident-to-alert conversion, rider apps, payments, passenger accounts, CAD, marketplace workflows, or canonical validators.

## Phase 7 Closure Audit Results

Checked during Phase 7 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed; PostGIS container running on host port `55432`.
- `make migrate-up`: passed and applied `000006_prediction_operations.sql`.
- `make migrate-status`: passed and reports migration versions 1 through 6 applied.
- `make test-integration`: passed with DB-backed telemetry, matcher, Vehicle Positions, GTFS import, GTFS Studio, Trip Updates diagnostics, and prediction operations tests using isolated temporary database setup.
- `make validate`: passed Phase 7 file smoke only. Canonical GTFS and GTFS-RT validators remain documented but not wired.
- `git diff --check`: passed.

Phase 7 implementation results:
- added `prediction.DeterministicAdapter` as the first real internal Trip Updates predictor behind `internal/prediction.Adapter`.
- made `cmd/feed-trip-updates` default to the deterministic adapter through `TRIP_UPDATES_ADAPTER=deterministic`, while preserving `TRIP_UPDATES_ADAPTER=noop`.
- generated non-empty Trip Updates for defensible in-service assignments using active published GTFS, latest telemetry, and current assignments.
- kept canceled trips outside the ETA coverage denominator and tracked them separately through canceled-trip and cancellation-alert-linkage metrics.
- persisted canceled-trip missing-alert linkage in prediction review details with `expected_alert_missing=true`.
- added prediction operation repository behavior for override create, replace, clear, expiry reads, review item persistence, review status transitions, and audit logging.
- kept matcher override consumption limited to `trip_assignment` and `service_state`; prediction-only disruption overrides are consumed through `prediction.OperationsRepository`.
- added minimal review queue lifecycle states: `open`, `resolved`, and `deferred`.
- withheld deadhead, layover, weak, stale, degraded, ambiguous, added-trip, short-turn, and detour cases instead of fabricating Trip Updates.
- exposed first-class prediction metrics in diagnostics and `feed_health_snapshot.details_json`.
- preserved Phase 3 Vehicle Positions, Phase 4 GTFS import, Phase 5 GTFS Studio, and Phase 6 public endpoint/non-coupling contracts.

## Phase 8 Closure Audit Results

Checked during Phase 8 closure:
- `command -v go`: passed, `/usr/local/bin/go`.
- `go version`: passed, `go version go1.26.2 darwin/amd64`.
- `make fmt`: passed.
- `make test`: passed.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make db-up`: passed; PostGIS container running on host port `55432`.
- `make migrate-up`: passed and applied `000007_phase_8_alerts_compliance.sql`.
- `make migrate-status`: passed and reports migration versions 1 through 7 applied.
- `make test-integration`: passed with DB-backed tests using isolated temporary database setup where supported.
- `make validate`: passed Phase 8 file smoke.
- `git diff --check`: passed.

Phase 8 implementation results:
- added persisted `service_alert`, `service_alert_informed_entity`, and `compliance_scorecard_snapshot` schema.
- added `feed_config.publication_environment` to distinguish dev from production scorecard behavior.
- added DB-backed Alerts authoring, lifecycle, audit logging, and public GTFS-RT Alerts publication.
- added Alerts-owned canceled-trip reconciliation from active canceled-trip overrides and Phase 7 missing-alert review signals.
- added on-demand public GTFS schedule ZIP publication from the active published feed version with deterministic ZIP bytes and stable `Last-Modified`.
- added `/public/feeds.json` with explicit feed metadata, validation, health, license, contact, and readiness fields.
- added publication metadata bootstrap that writes `feed_config`, `published_feed`, `consumer_ingestion`, and `marketplace_gap` records.
- added compliance scorecard snapshot persistence and validator command adapters for static GTFS and GTFS-RT validation.
- kept realtime `published_feed.revision_timestamp` as publication/bootstrap metadata revision; realtime feed generation does not update it.
- kept schedule `published_feed.revision_timestamp` tied to active schedule publication/bootstrap metadata, not request time.

## Phase 9 Closure Audit Results

Checked during Phase 9 closure:
- `gofmt -w ./cmd ./internal`: passed.
- `go mod tidy`: passed.
- `go test ./...`: passed.
- `make validators-install`: passed; installed the pinned static GTFS validator JAR and Docker-backed GTFS-RT validator wrapper.
- `make validators-check`: passed.
- `make validate`: passed with pinned validator tooling checks.
- `make test-integration`: passed with DB-backed tests using isolated temporary databases where supported.
- `make smoke`: passed with pinned validator tooling checks and HTTP/runtime hardening package coverage.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `git diff --check`: passed.

Phase 9 implementation results:
- tightened `/admin/validation/run` to accept only `validator_id`, `feed_type`, and optional `feed_version_id`.
- added handler coverage proving schedule and Vehicle Positions, Trip Updates, and Alerts realtime validation runs return `200`, persist results, normalize status, and record feed type/feed version.
- made realtime validation prefer internal builder-derived protobuf bytes and use configured feed URLs only as fallback.
- added repo-supported validator install/check targets and lock file for pinned static GTFS and GTFS-RT validator tooling.
- added structured request logs, request IDs, redaction rules, and `/metrics` only when `METRICS_ENABLED=true`.
- tightened `/readyz` for `agency-config`, Trip Updates, and Alerts so DB reachability alone is not enough: agency-config also requires an active schedule feed plus complete published feed metadata, and realtime feed services require an active GTFS feed.
- strengthened DB-backed device rebind tests for spoof rejection and immediate old-token invalidation.
- strengthened assignment current-row race tests with a partial-index assertion and higher concurrency.

## Phase 10 Closure Audit Results

Checked during Phase 10 closure:
- `make validators-install`: passed.
- `make validators-check`: passed.
- `make test`: passed.
- `make smoke`: passed.
- `make validate`: passed.
- `make demo-agency-flow`: passed and verified DB bootstrap, validator install/check, sample GTFS import, publication metadata bootstrap, authenticated telemetry ingest, public `schedule.zip`, public `feeds.json`, public realtime protobuf feeds, protected debug/admin routes including GTFS Studio, validation run flow, scorecard, and consumer-ingestion visibility.
- `docker compose -f deploy/docker-compose.yml config`: passed.
- `make test-integration`: passed.
- `git diff --check`: passed.

Phase 10 implementation results:
- rewrote `README.md` to describe the current Phase 9 runtime surface, public/protected endpoints, quickstart, deployment path, limitations, and truthful Caltrans/CAL-ITP-aligned wording.
- added tutorial docs for local quickstart, Docker Compose deployment, agency demo flow, production checklist, and CAL-ITP readiness checklist.
- added `scripts/demo-agency-flow.sh`, `make demo-agency-flow`, and `task demo:agency`.
- updated `scripts/bootstrap-dev.sh` to print current service commands, public feed URLs, protected debug/admin examples, validator setup, and the executable demo target.
- added repo-owned docs assets under `docs/assets/` and documented source specs plus alt text.
- updated `docs/dependencies.md` for local demo packaging tools.

## Phase 11 Closure Audit Results

Checked during Phase 11 closure:
- pre-edit `command -v go`: passed, `/usr/local/bin/go`.
- pre-edit `go version`: passed, `go version go1.26.2 darwin/amd64`.
- pre-edit `make validators-install`: passed.
- pre-edit `make validators-check`: passed.
- pre-edit `make test`: passed.
- pre-edit `make smoke`: passed.
- pre-edit `make demo-agency-flow`: passed.
- pre-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- pre-edit `make validate`: passed.
- pre-edit `make migrate-status`: passed and reports migration versions 1 through 8 applied.
- pre-edit `make test-integration`: passed.
- pre-edit `git diff --check`: passed.
- post-edit `make validators-check`: passed.
- post-edit `make validate`: passed.
- post-edit `make test`: passed.
- post-edit `make smoke`: passed.
- post-edit `make demo-agency-flow`: passed.
- post-edit `make test-integration`: passed.
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed.
- post-edit `git diff --check`: passed.
- Blocked commands: none.

Phase 11 implementation results:
- added `docs/compliance-evidence-checklist.md` as the evidence package separating implemented repo capability, deployment/operator proof, and third-party confirmation.
- mapped current repo support to Caltrans/CAL-ITP-style expectations without claiming full compliance, production readiness, consumer acceptance, or marketplace equivalence.
- updated `docs/dependencies.md` with a Phase 11 wiring reality table for all originally mentioned external tools and repos.
- documented real integrations as wired where code-backed: Postgres/PostGIS, pgx, Goose, MobilityData validators, GTFS-RT protobuf bindings, Docker/Docker Compose, Task, local demo tools, and internal Prometheus-format `/metrics`.
- documented optional/deferred or workflow-only systems truthfully: TheTransitClock, other external predictors, Prometheus/Grafana deployment, OpenTelemetry, consumer submission APIs, Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.
- tightened README and tutorial wording by linking to the evidence checklist and clarifying deployment-owned observability and consumer-ingestion proof limits.

## Phase 12 Step 1 Progress

Phase 12 Step 1 (repo-side docs/runbooks/evidence packaging) is complete:
- added deployment evidence overview and targeted runbooks under `docs/runbooks/`
- added `docs/evidence/` structure with committed templates and operator-owned captured-artifact placeholders
- added lightweight README links to deployment evidence docs
- added Phase 12 Step 1 handoff notes while keeping claim boundaries explicit

## Phase 12 Step 2 Progress

Phase 12 Step 2 produced a real local evidence packet at `docs/evidence/captured/local-demo/2026-04-22/`:
- local loopback public feed fetch proof for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts
- local reverse proxy route map and protected admin/debug boundary checks
- validator records for schedule and all three realtime feeds, all failed and retained without omission
- local request-log and scorecard monitoring evidence, with alert lifecycle explicitly missing
- one local Postgres dump/restore drill into `open_transit_rt_restore_drill_20260422`, including restored row counts and feed fetch checks against the restored database
- manual scorecard export artifacts with checksums

An earlier operator intake packet exists at `docs/evidence/captured/hosted-pending/2026-04-22/`; it remains historical intake material only. The completed hosted proof packet is `docs/evidence/captured/oci-pilot/2026-04-24/`.

Phase 12 Step 3 implemented repo-side closure guardrails but did not collect hosted evidence:
- `scripts/install-validators.sh` now writes a GTFS-RT validator wrapper that drives the pinned MobilityData webapp API against server-derived local artifacts instead of passing unsupported CLI flags to the image.
- `scripts/check-validators.sh` now verifies Java, Docker, `curl`, `python3`, pinned artifacts, and a webapp-API wrapper shape before allowing pinned validator checks to pass. It can use `JAVA_BINARY` or the Homebrew Java 17 path when the macOS `/usr/bin/java` shim is not usable.
- `scripts/duckdns-pilot.sh` can bootstrap a local DuckDNS/Caddy pilot using generated secrets under `.cache/duckdns-pilot/`.
- `docs/dependencies.md` and `README.md` now document the Java and `python3` validator-tooling requirements.

Homebrew Java 17 was installed and the strict repo-side validator gate now passes locally.
The OCI pilot at `https://open-transit-pilot.duckdns.org` now has public HTTPS feed proof, TLS/redirect evidence, clean hosted validator records, public-edge auth-boundary proof, SSH-tunneled admin auth proof, monitoring/alert lifecycle evidence, backup/restore evidence, deployment data-restore rollback proof, and scorecard export job-history proof.

Phase 12 is closed for hosted/operator evidence because the OCI pilot packet passed the hosted audit. Third-party consumer confirmation has not been collected and remains outside Phase 12.

## Phase 13 Progress

Phase 13 is complete for the initial consumer submission evidence layer:
- added `docs/consumer-submission-evidence.md` with status definitions, allowed claims by status, tracker requirements, and acceptance-scope rules
- added `docs/evidence/consumer-submissions/README.md` with tracker freshness fields, Phase 12 packet linkage, current target summary, and current OCI pilot feed URLs for future submission packets
- added current evidence records for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land under `docs/evidence/consumer-submissions/current/`
- added reusable templates for all seven targets under `docs/evidence/consumer-submissions/templates/`
- kept all current records at `not_started` because no redacted real submission, review, acceptance, rejection, or blocker evidence is present in the repo
- documented that validator success and public fetch proof are supporting evidence only, not consumer acceptance

## Phase 20 Progress

Phase 20 is complete for the consumer packet preparation and California readiness evidence scope:
- added complete prepared packet drafts for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land under `docs/evidence/consumer-submissions/packets/`
- added packet evidence freshness fields, submission-method fields marked `not verified`, operator warnings, redaction notes, next actions, and allowed wording
- added `docs/evidence/consumer-submissions/status.json` and kept it aligned with the human-readable tracker for target name, status, packet path, prepared timestamp, and evidence references
- updated all current target records to `prepared` only because complete packets exist
- added `docs/california-readiness-summary.md`, including the agency-owned stable URL/domain gap
- added `docs/marketplace-vendor-gap-review.md`
- did not contact external portals, automate submissions, guess submission paths, add backend API behavior, or claim acceptance, compliance, consumer ingestion, marketplace equivalence, hosted SaaS availability, agency endorsement, or production-grade ETA quality

## Phase 21 Progress

Phase 21 is complete for the community/governance/multi-agency docs and process scope:
- added `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, GitHub issue templates, and a PR template
- added `docs/governance.md`, `docs/release-process.md`, `docs/support-boundaries.md`, `docs/multi-agency-strategy.md`, and `docs/roadmap-status.md`
- documented maintainer authority for PR merges, releases, docs/evidence wording, and competing design decisions
- documented support boundaries without paid support, SLA, hosted SaaS, agency endorsement, vendor equivalence, or production-readiness claims
- documented current single-agency/pilot assumptions and what would need code changes before true multi-tenant hosting
- added teaching visuals for contribution paths, community workflow, single-vs-multi-agency strategy, evidence maturity, and support boundaries
- updated `docs/assets/README.md` with filename, purpose, usage, alt text, generation method, prompt/spec, and truthfulness notes for the new visuals
- did not change backend behavior, API contracts, database schema, public feed URLs, consumer-submission statuses, external integrations, or evidence claims

## Track A Progress

Track A is complete for the external-proof/adoption workflow scope:
- added `docs/evidence/consumer-submissions/submission-workflow.md` for official-path verification, pre-submission checks, evidence intake, and status transition rules
- added README-only target artifact directories under `docs/evidence/consumer-submissions/artifacts/`
- added `docs/agency-owned-domain-readiness.md`
- updated evidence, readiness, roadmap, docs index, status, and handoff docs to point operators to the workflow
- kept all seven consumer and aggregator targets at `prepared` only
- did not add placeholder artifacts, helper scripts, portal automation, backend behavior, public feed URL changes, or external evidence claims

## Phase 23 Progress

Phase 23 is complete as blocker-documented closure only:
- updated `docs/phase-23-agency-owned-deployment-proof.md` with Outcome B blocker status
- updated `docs/agency-owned-domain-readiness.md` with the Phase 23 blocker record and future operator next-actions checklist
- updated `docs/california-readiness-summary.md` and `docs/compliance-evidence-checklist.md` to keep final-root evidence listed as missing before stronger California readiness language
- added `docs/handoffs/phase-23.md`
- did not create a final-root evidence packet
- did not run hosted evidence audit for a final-root packet
- did not refresh prepared consumer packets or update `docs/evidence/consumer-submissions/status.json`
- did not claim compliance, consumer acceptance, agency endorsement, hosted SaaS, marketplace equivalence, or production-grade ETA quality

## Phase 24 Progress

Phase 24 is complete for the docs/process and evidence-template scope:
- added `docs/tutorials/real-agency-gtfs-onboarding.md`
- added `docs/tutorials/gtfs-validation-triage.md`
- added `docs/evidence/real-agency-gtfs/README.md`
- added `docs/evidence/real-agency-gtfs/templates/import-review-template.md`
- updated first-run, production checklist, tutorial index, evidence index, README, phase status, current status, and latest handoff docs
- added `docs/handoffs/phase-24.md`
- kept the evidence scaffold template-only until real agency-approved, public-safe evidence exists
- kept Phase 23 final-root status unchanged: no agency-owned or agency-approved final public feed root is available in repo evidence
- did not add real agency data, placeholder artifacts, fake evidence, backend behavior, public feed URL changes, consumer status changes, or unsupported readiness/compliance claims

## Phase 25 Progress

Phase 25 is complete for the docs/process and template-only evidence scope:
- added `docs/tutorials/device-avl-integration.md`
- added `docs/tutorials/device-token-lifecycle.md`
- added `docs/evidence/device-avl/README.md`
- added `docs/evidence/device-avl/templates/integration-review-template.md`
- updated first-run, production checklist, tutorial index, evidence index, README, phase status, current status, decisions, and latest handoff docs
- added `docs/handoffs/phase-25.md`
- kept the evidence scaffold template-only until real public-safe device or AVL integration evidence exists
- kept response examples limited to telemetry-ingest behavior confirmed by code and tests
- kept `/v1/events` documented as an authenticated admin/debug path, not a public or consumer-facing feed
- kept `docs/dependencies.md` unchanged because no named external vendor, adapter implementation, or dependency status changed
- did not add real device data, private vendor payloads, credentials, hardware certifications, fake evidence, backend behavior, public feed URL changes, consumer status changes, or unsupported AVL reliability/compliance claims

## Phase 29 Progress

Phase 29 is complete for the synthetic replay evidence expansion scope:
- added replay fixtures for after-midnight service, exact and non-exact frequency trips, block continuity, long layover withholding, sparse telemetry, noisy/off-shape GPS, stale/ambiguous hard patterns, cancellation alert linkage, and manual override before/after expiry
- added replay fixture support for `frequencies` and optional manual override `expires_at`
- aligned replay telemetry repository behavior with the production latest-per-vehicle contract for feed snapshots
- strengthened replay assertions for cancellation alert linkage, unsupported disruption-withheld counts, and degraded-by-reason visibility
- added focused realtime-quality tests for Phase 29 scenarios while preserving Phase 19 replay fixtures
- documented that Phase 29 expands synthetic replay coverage only
- documented that real-world observed-arrival/departure evidence remains unavailable
- explicitly deferred real route/time-period quality metrics because no real deployment or observed-arrival evidence exists in the repo
- did not add external predictors, real private telemetry, private agency GTFS, Operations Console changes, public feed URL changes, GTFS-RT contract changes, consumer status changes, auth-boundary changes, dependency changes, production-grade ETA claims, real-world ETA accuracy claims, consumer acceptance claims, agency endorsement claims, hosted SaaS claims, or CAL-ITP/Caltrans compliance claims

## Phase 30 Progress

Phase 30 closed as Outcome B — blocker-documented closure only:
- no authorized submission, official-path verification evidence, or target-originated artifact was available
- no Phase 30 target was selected
- target selection is deferred until an operator is authorized and either official-path verification or target-originated evidence can be retained
- no individual target status changed to `blocked` because no target-specific blocker artifact exists
- `docs/evidence/consumer-submissions/status.json` and all current target records were left unchanged
- tracker/status consistency still shows all seven targets `prepared`
- artifact directories remain README-only; no receipts, screenshots, tickets, correspondence, blocker notes, or placeholder artifacts were added
- Mobility Database and transit.land may be considered as future candidate suggestions once authorized, but neither was selected in Phase 30
- did not contact external portals, automate submissions, guess submission paths, add artifacts, change public feed URLs, change GTFS-RT contracts, change telemetry/device APIs, add consumer submission APIs, or claim submission, review, acceptance, rejection, ingestion, display, compliance, agency endorsement, hosted SaaS availability, marketplace/vendor equivalence, production-grade ETA quality, or consumer adoption

Phase 30 closure audit results:
- pre-edit `make validate`: passed
- pre-edit `make test`: passed
- pre-edit `git diff --check`: passed
- post-edit `make validate`: passed
- post-edit `python3 -m json.tool docs/evidence/consumer-submissions/status.json`: passed
- post-edit tracker/status consistency check: passed; all seven targets remain `prepared`
- post-edit `make test`: passed
- post-edit `make realtime-quality`: passed
- post-edit `make smoke`: passed
- post-edit `make test-integration`: passed
- post-edit `docker compose -f deploy/docker-compose.yml config`: passed
- post-edit `git diff --check`: passed
- post-edit targeted artifact scan: passed; artifact directories contain README files only
- post-edit targeted tracker diff check: passed; `status.json`, current target records, and artifact directories were not edited
- post-edit context-aware forbidden-claim scan: reviewed; matches are negated statements, definitions, transition/future-state wording, or blocker explanations
- post-edit targeted redaction-sensitive term scan: reviewed; matches are security/redaction rules or existing negative boundary wording, not exposed secrets or private artifacts
- blocked commands: none

## Next Recommended Step

Start Phase 31 — Agency Pilot Program Package when maintainers are ready to continue Track B implementation. Phase 31 must proceed from the prepared-only consumer state and must not assume submission, review, acceptance, rejection, blocker, ingestion, listing, display, or adoption evidence exists.

Use the Track A workflow when a human operator is ready to verify an official target path or record real target-originated evidence. If no real third-party artifacts are available, keep every target at `prepared`.

Consumer/aggregator status should still move only when real external artifacts exist:
1. review the prepared target-specific packet under `docs/evidence/consumer-submissions/packets/`
2. verify the target's current official submission path outside the repo
3. confirm agency identity, feed URLs, license/contact metadata, validation status, consumer-specific requirements, and redactions
4. submit through the named consumer or aggregator workflow outside the repo only when authorized
5. add only redacted correspondence, receipt, ticket, portal screenshot, rejection, blocker, or acceptance evidence
6. keep the OCI pilot operator jobs and validator tooling maintained if the pilot remains live

## What Not To Do Next

Do not:
- bypass the prediction adapter boundary
- add rider-facing functionality
- add payments, passenger accounts, or dispatcher CAD
- add a heavy frontend stack
- tightly couple to an external predictor
- merge draft GTFS and published GTFS into one model
- leave placeholder sample feed data in production paths once real feed generation starts
- contact external consumer portals from repo automation
- claim submitted, under-review, accepted, rejected, or blocked consumer status without retained target-originated evidence
