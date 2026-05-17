# Architecture Decisions

This document records architecture-significant decisions so later phases do not re-decide core product boundaries.

## ADR-0001 — Keep the backend mostly Go

Open Transit RT should use Go services and internal packages for core backend behavior. Early admin and Studio surfaces should prefer simple server-rendered HTML unless a later phase documents a need for a heavier frontend stack.

## ADR-0002 — Use Postgres/PostGIS as source of truth

Postgres stores agency configuration, GTFS feed versions, telemetry, assignments, overrides, audit logs, validation reports, incidents, feed metadata, and compliance workflow state. PostGIS is required for future nearest-shape and spatial matching behavior.

## ADR-0003 — Use versioned migrations

Schema changes live under `db/migrations` and are applied through `cmd/migrate`. Migrations are the source of truth for the executable database schema.

`db/schema.sql` is deprecated as an executable schema file. It remains only as a short comment-only pointer for readers or tools that still expect the path to exist. It must not contain `CREATE`, `ALTER`, or `DROP` statements, must not be edited as an independent schema definition, and must not be used to apply database changes. If a future phase wants a full schema snapshot, it should generate it from migrations and document that workflow before replacing the pointer file.

## ADR-0004 — Keep Trip Updates pluggable

Trip Updates must stay behind a prediction adapter boundary. Open Transit RT owns GTFS management, telemetry, assignments, audit logs, and Vehicle Positions. Optional predictors such as TheTransitClock may generate ETAs or Trip Updates only behind an adapter.

## ADR-0005 — Publish Vehicle Positions first

Vehicle Positions are the first production-grade public realtime output. Trip Updates and Alerts are architecture-binding but implemented in later phases.

## ADR-0006 — Prefer unknown over false certainty

The matcher must be conservative. Low-confidence or contradictory evidence should produce `unknown` plus incidents/diagnostics instead of a speculative trip descriptor.

## ADR-0007 — Manual overrides take precedence

Operator overrides are part of the core model. Active overrides must beat automatic matching until they expire or are cleared, and privileged actions must be audit logged.

## ADR-0008 — Keep draft and published GTFS separate

GTFS Studio draft data and active published feed versions must not collapse into one model. Import and Studio are two sources that publish through a shared validated feed-version pipeline.

## ADR-0009 — Stable public URLs are product contracts

Public schedule, Vehicle Positions, Trip Updates, and Alerts URLs must stay stable across feed updates and rollback. Version changes happen behind those URLs.

## ADR-0010 — Phase 0 is foundation-only

Phase 0 may design schemas, contracts, and docs for later requirements, but it must not implement later-phase runtime behavior such as durable telemetry, trip matching, GTFS import, protobuf feed generation, Trip Updates, or Alerts.

## ADR-0011 — Persist telemetry through an agency-scoped repository

Telemetry ingest writes must go through a repository backed by Postgres/PostGIS. The repository classifies accepted, duplicate, and out-of-order telemetry inside one transaction protected by a deterministic advisory lock derived from agency and vehicle identity. The lock key is a SHA-256-derived signed 64-bit value; theoretical collisions only serialize unrelated streams and do not merge data because SQL predicates and uniqueness remain authoritative. Canonical accepted telemetry uniqueness is vehicle-scoped by `(agency_id, vehicle_id, observed_at)`; `device_id` is retained for audit/debug but does not define the canonical latest vehicle position.

Invalid JSON and invalid telemetry payloads are rejected before repository storage in Phase 1. The database `rejected` status remains reserved for a later ingest-audit phase that explicitly designs rejected-payload retention.

## ADR-0012 — Persist explicit deterministic assignment outcomes

Phase 2 persists every matcher outcome as a `vehicle_trip_assignment` row, including `unknown`. Unknown results close any previous active row so stale or low-confidence telemetry cannot leave a prior confident trip active. Unknown rows carry `service_date` whenever agency timezone and observed timestamp can be resolved; the column remains nullable only for unresolved cases.

Assignment reasons and degraded state use a small stable taxonomy. `score_details_json` is intentionally loose debug JSON for Phase 2 and is not a stable public API or integration schema. The internal convention is that matcher-generated score details include `score_schema`; candidate-based score details also include `trip_id`, `start_time`, and `observed_local_seconds` when resolvable. Future public or adapter-facing diagnostics should define a separate versioned contract rather than depending on this debug payload.

Phase 2 treats `missing_shape` as both a reason code and a dedicated degraded-state category. Missing shapes reduce confidence but do not automatically prevent a match when other evidence is strong. Route-hint matching is reserved for a future telemetry/input expansion and is not part of the active Phase 2 reason-code taxonomy.

`internal/state.Engine` is the only valid production matcher entry point. It requires schedule and assignment repositories. `NewEngine` returns an error for invalid construction, and `MustNewEngine` is reserved for tests/bootstrap paths that intentionally want panic-on-error behavior. The old placeholder rule-based matcher path was removed so placeholder feed code cannot accidentally look like production matching behavior.

Phase 2 service-day resolution intentionally considers only two service dates for each observation: the observed agency-local date and the immediately previous agency-local date. This covers same-day and typical after-midnight service, including GTFS times greater than `24:00:00`, but later phases must explicitly extend the resolver before assuming broader multi-day post-midnight coverage.

Active manual overrides are absolute in Phase 2 and are evaluated before stale telemetry fallback. When an override references a resolvable active-feed trip, the persisted assignment includes `feed_version_id` and `block_id` so manual rows remain first-class assignment records.

Block-transition scoring requires same block, temporal plausibility, and the nearest plausible next-trip sequencing when start-time identity is available. A later same-block trip does not receive block-transition credit merely because it is later than the previous assignment. Explicit telemetry bearing validity is distinct from numeric truthiness: numeric `bearing: 0` is a valid true-north bearing only when the stored payload explicitly contains a numeric `bearing` field. Null, malformed, or payload-missing zero values are invalid for movement-direction scoring. Persisted shape distance preserves `0` as a valid value.

Repeated identical degraded unknown states reuse the active degraded assignment only when degraded state, reason codes, service date, and telemetry evidence match. Telemetry evidence means matching `telemetry_event_id` when either row has one, with exact `active_from` equality used only as the no-telemetry fallback. Materially new telemetry evidence or a service-day change creates a replacement unknown row and must not leave a previous confident row active.

The Phase 2 handoff is expected to describe actual implemented matcher behavior, not aspirational behavior. After the semantic-closure pass, the handoff and implementation are aligned on constructor behavior, override precedence, degraded-state handling, system-failure taxonomy, batching, block-transition successor rules, bearing validity, and post-midnight service-day limits.

## ADR-0013 — Build Vehicle Positions from one DB-backed snapshot model

Phase 3 Vehicle Positions publishing uses latest accepted telemetry plus current persisted Phase 2 assignments as the source of truth. The protobuf endpoint and JSON debug endpoint both render from the same immutable in-memory snapshot per request, so HTTP handlers do not duplicate publication decisions.

The snapshot caps latest telemetry before assignment lookup and stale/publication evaluation. `ListLatestByAgency` therefore has a hard ordering contract: one latest accepted row per vehicle, ordered by `observed_at DESC, id DESC`. Automatic assignments are only publishable as trip descriptors when linked to the latest telemetry event, which prevents read-committed cross-table timing from producing false trip certainty.

GTFS-RT protobuf types remain isolated to `internal/feed`. Public Vehicle Positions responses set `gtfs_realtime_version = "2.0"` and return normal successful empty `FeedMessage` responses when there is no telemetry or all vehicles are suppressed. JSON debug output carries per-vehicle publication decisions and telemetry age for inspectability, but it is diagnostic rather than a stable public integration contract.

## ADR-0014 — Use transactional feed-version staging for GTFS ZIP import

Phase 4 GTFS ZIP import stages schedule rows by inserting them under a new inactive `feed_version` inside the publish transaction, then retiring the previous active version and activating the new version before commit. Failed validation creates no `feed_version`; publish failures roll back the transaction so no inactive staged version or partial GTFS rows remain. `gtfs_import.feed_version_id` is set only after a successful publish and remains `NULL` for failed imports.

`gtfs_import` and `validation_report` store the normalized internal import report. Validation failures and publish/database failures both write a failed `validation_report` outside the publish transaction when possible; if that report write fails, the importer reports the storage failure and does not claim validation-report persistence. Phase 4 intentionally does not store original ZIP bytes in Postgres and does not wire MobilityData GTFS Validator; canonical validator integration remains a later compliance task.

The internal validator intentionally enforces the Phase 4 contract before activation: required GTFS files, supported `route_type` values (`0`-`7` and extended `100`-`1702`), numeric ranges, usable service sources, references, shape ordering, stop_times references, optional `block_id` preservation, optional shapes/frequencies, and GTFS times beyond `24:00:00` without normalizing away imported time text. A service source is usable only when a calendar row has at least one active weekday or a calendar_dates row adds service with `exception_type=1`.

This staging model began as the GTFS ZIP import publish path. Phase 5 refactored the activation logic into a shared internal publisher used by both ZIP imports and GTFS Studio drafts.

## ADR-0015 — Use typed GTFS Studio draft tables and direct shared publishing

Phase 5 stores GTFS Studio draft data in typed draft tables for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies. The generic `gtfs_draft_record` table remains unused legacy scaffold and is not part of runtime Studio editing.

`gtfs_draft` owns draft metadata and traceability. It records status, cloned-source `base_feed_version_id`, `last_publish_attempt_id`, `last_published_feed_version_id`, and soft-discard fields. Drafts cloned from an active feed capture the active `feed_version` as provenance; explicit blank drafts and drafts created when no active feed exists have no base feed version.

Draft-level discard is soft discard. Discarded drafts retain metadata and typed rows for auditability, are hidden from the default list view, and become read-only and not publishable. Drafts in `published` status also become read-only by default after successful publish. Entity remove operations only delete rows inside the current editable draft and never delete previously published GTFS rows, feed versions, publish attempts, validation reports, or audit history.

Draft agency editing is one row scoped to the draft's agency. On successful draft publish, that draft agency row maps into the canonical `agency` table inside the same publish transaction before the new `feed_version` is activated. Draft agency edits do not mutate published agency metadata before publish.

Studio publish converts typed draft rows into the same internal feed model used by ZIP import, then calls the shared validation and activation helper directly. It does not generate or re-import a synthetic ZIP. Non-editable draft statuses are rejected before draft-to-feed conversion, validation, or shared publish activation begins.

The first Studio UI is intentionally minimal server-rendered HTML from Go stdlib packages. It provides operational row forms, not map editing, timetable design, or a heavy frontend application.

## ADR-0016 — Define Phase 6 Trip Updates and Alerts as pluggable empty-feed architecture

Phase 6 establishes Trip Updates and Alerts feed boundaries without implementing ETA quality or alert authoring. Trip Updates use a narrow `internal/prediction.Adapter` contract that accepts the active published GTFS feed version, persisted latest telemetry, persisted current assignments, and the Vehicle Positions feed URL. The default adapter is an explicit no-op that returns no Trip Updates plus diagnostics; it is not a placeholder prediction algorithm.

Trip Updates diagnostics are persisted to `feed_health_snapshot` with `feed_type = 'trip_updates'` and a normalized `details_json` trace containing adapter name, diagnostics status and reason, active feed version ID, input counts, Vehicle Positions URL, and persistence outcome. This reuses the existing health/traceability schema rather than adding a Phase 6 migration.

Trip Updates and Alerts protobuf endpoints return valid empty GTFS-Realtime `FeedMessage` payloads with `gtfs_realtime_version = "2.0"`, `FULL_DATASET`, and `FeedHeader.timestamp` derived from the same snapshot `GeneratedAt` timestamp used for `Last-Modified`. Non-empty Trip Updates output must use deterministic feed entity ordering and ordered `stop_time_update` entries.

Alerts are architecture-only in Phase 6. The Alerts endpoint returns valid empty protobuf and JSON debug output with deferred status, but it does not write `feed_health_snapshot` rows, persist alert records, or derive public alerts from incidents/manual overrides yet. Phase 7 added canceled-trip missing-alert linkage signals; public Alerts authoring and persistence remain future work.

The Trip Updates packages are intentionally not dependencies of telemetry ingest, Vehicle Positions, or GTFS Studio. A non-coupling test guards that boundary.

## ADR-0017 — Use an internal conservative deterministic predictor for Phase 7

Phase 7 replaces the default Trip Updates no-op runtime path with an internal deterministic adapter behind `internal/prediction.Adapter`. The adapter uses only the active published GTFS feed, latest accepted telemetry, current persisted assignments, and prediction-operation repository interfaces. It does not move matching ownership into the predictor and does not couple predictor internals into telemetry ingest, Vehicle Positions, GTFS import, or GTFS Studio.

The first predictor emits stop-level Trip Updates only when the assignment is in service, current, linked to the active feed, linked to the latest telemetry where required, above the publication confidence threshold, and resolvable to a GTFS trip instance. Prediction times are schedule-deviation projections from the current assigned stop, not production-grade learned ETAs. Weak, stale, degraded, deadhead, layover, ambiguous, added-trip, short-turn, and detour cases are withheld and recorded as prediction review items instead of fabricating Trip Updates.

Canceled trips are not part of the ETA coverage denominator. They are emitted as conservative `CANCELED` Trip Updates when represented by active prediction overrides, and they are tracked by separate cancellation and cancellation-alert-linkage metrics. Because public Alerts authoring remains deferred, canceled-trip review details persist `expected_alert_missing=true` and `cancellation_alert_linkage_status="missing_alert_authoring_deferred"`.

Prediction review workflow uses the existing `incident` table with `incident_type = 'prediction_review'` and a minimal lifecycle of `open`, `resolved`, and `deferred`. Phase 7 extends the incident status check constraint to support `deferred`. Override create, replace, clear, and review-status changes write `audit_log` rows.

The matcher continues to consume only assignment/service-state overrides from `manual_override` (`trip_assignment` and `service_state`). Prediction-only disruption overrides such as canceled trips, added trips, detours, and short turns are consumed through `prediction.OperationsRepository` so they cannot force invalid assignment states into `vehicle_trip_assignment`.

## ADR-0018 — Publish static GTFS ZIPs on demand from active published tables

Phase 8 adds `/public/gtfs/schedule.zip` as the stable public static GTFS URL. The schedule ZIP is generated on demand from the active published `feed_version` tables, not from GTFS Studio draft rows and not from placeholder sample files.

ZIP entries and CSV rows are written in deterministic order. ZIP entry modified times use the active feed revision time, so identical active GTFS data produces stable bytes across requests. The endpoint `Last-Modified` header uses the same active feed revision time. The endpoint does not materialize or cache ZIP bytes in Postgres in Phase 8; a future cache may be added only if it preserves deterministic bytes and stable `Last-Modified` semantics.

For `published_feed`, schedule `revision_timestamp` changes when schedule publication/bootstrap metadata changes or when GTFS import/Studio publish activates a new schedule feed. It does not change merely because `/public/gtfs/schedule.zip` was requested.

## ADR-0019 — Persist Service Alerts separately from feed serialization

Phase 8 stores public Service Alerts in `service_alert` and `service_alert_informed_entity`. `internal/alerts` owns authoring, persistence, lifecycle, audit logging, and canceled-trip reconciliation. `internal/feed/alerts` owns only GTFS-RT protobuf/JSON feed rendering from persisted published alerts.

Canceled-trip Trip Updates from Phase 7 remain prediction-owned, but alert satisfaction is Alerts-owned. The reconciler reads active canceled-trip overrides and open prediction-review incidents with `expected_alert_missing=true`, creates or updates a published cancellation Service Alert, and links the review incident to `service_alert.id`. Prediction packages do not import Alerts packages.

## ADR-0020 — Use feed_config and published_feed as the license/contact contract

Phase 8 makes the metadata contract explicit:
- `feed_config` stores agency-level defaults: `public_base_url`, `feed_base_url`, `technical_contact_email`, `license_name`, `license_url`, `validator_strictness`, and `publication_environment`.
- `published_feed` stores per-feed resolved publication state: `canonical_public_url`, `license_name`, `license_url`, `contact_email`, `revision_timestamp`, `activation_status`, and `active_feed_version_id`.

`/public/feeds.json` reads per-feed values from `published_feed` first. It may resolve empty license/contact fields from `feed_config`, but scorecard readiness still evaluates whether all required values are complete. Response timestamps are RFC3339 UTC JSON timestamps or `null`.

Realtime `published_feed.revision_timestamp` is a publication/bootstrap metadata revision. Vehicle Positions, Trip Updates, and Alerts feed generation must not update it on every request. Feed freshness belongs in `feed_health_snapshot`, not in `published_feed.revision_timestamp`.

## ADR-0021 — Validator-backed scorecards distinguish dev from production

Phase 8 adds canonical validator command adapters for static GTFS and GTFS-RT. Validator results are normalized into `validation_report`. The adapters parse structured JSON from stdout, stderr, or validator output files, count errors/warnings/info notices, preserve the raw parsed report under `report_json.raw_report`, and derive `passed`, `warning`, or `failed` status from the normalized counts plus command exit status. If validator tooling is absent, the system stores `status='not_run'` instead of pretending validation passed.

Production mode is agency-scoped and stored as `feed_config.publication_environment = 'production'`. In production mode, missing canonical validator execution makes scorecard validation red. In dev mode, missing validators are yellow/not-run. `validator_strictness` controls failure handling, but it does not define production mode by itself.

## ADR-0022 — Harden admin, device, validator, and current-assignment boundaries

Post-Phase-8 hardening introduces a service-level `APP_ENV` switch. In `production`, services fail fast without explicit `DATABASE_URL`, admin JWT config, `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`. `BIND_ADDR` defaults to `127.0.0.1`; binding to `0.0.0.0` assumes a TLS-terminating reverse proxy.

Admin identity is an HS256 JWT contract with required `sub`, `agency_id`, `iat`, `exp`, `iss`, and `aud` claims. Bearer auth is the default for machine/API admin calls. Optional `admin_session` cookie auth is for browser-admin flows only and requires CSRF validation on every unsafe cookie-authenticated admin method. Roles are loaded from `agency_user` and `role_binding`, not from token claims. Default token TTL is 8 hours, clock skew allowance is 2 minutes, and secret rotation accepts the active secret plus `ADMIN_JWT_OLD_SECRETS`; server-side `jti` replay tracking is deferred.

Telemetry ingest now verifies opaque device Bearer tokens against peppered HMAC hashes in `device_credential` and enforces active agency/device/vehicle binding before persistence. Device rebinding is an admin-managed operation that rotates the token and immediately invalidates the old token/binding.

Validator runs no longer accept request-supplied commands, paths, argv, or output directories. `/admin/validation/run` accepts only `validator_id`, `feed_type`, and optional `feed_version_id`; the server derives artifacts itself, preferring generated local feed bytes/temp files over URLs whenever possible. Validators run with argv-based `exec.CommandContext`, timeout/output/report caps, and temp/output confinement.

`vehicle_trip_assignment` now has a partial unique index for one current row per `(agency_id, vehicle_id)`, and `SaveAssignment` serializes writes with a per-agency/per-vehicle advisory transaction lock before reading or closing the current row.

## ADR-0023 — Use compiled Go binaries and native services for memory-constrained deployments

Phase 12 introduces `scripts/oci-pilot.sh` and `deploy/oci/` for deploying to an Oracle Cloud VM.Standard.E2.1.Micro instance (1 OCPU, 1 GB RAM). The following decisions apply to this deployment path and any future resource-constrained single-node target:

**No Docker for application services.** The existing `deploy/docker-compose.yml` Postgres image consumes ~150–300 MB RAM by default. On the OCI VM.Standard.E2.1.Micro pilot host, the measured usable RAM is about 503 MiB, so Docker leaves no safe operating margin for five Go services and Caddy. The Oracle Linux 9 pilot installs PostgreSQL 15 + PostGIS natively from PGDG with Oracle EPEL/CodeReady dependencies, uses expanded swap for the package transaction, and tunes PostgreSQL with `shared_buffers=96MB`, `max_connections=25`, and `work_mem=3MB`. The Docker Compose file remains for local development only and must not become the only documented production path.

**Compiled binaries, not `go run`.** `go run` forks the compiler process at startup, consuming an additional ~100–150 MB during compilation. All production commands are pre-compiled with `GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w"` from the developer's local machine and uploaded via scp. This cross-compilation approach requires no Go toolchain on the OCI instance at runtime (though Go is installed by `setup-instance.sh` to support future on-instance rebuilds if needed).

**Caddy via official Debian package.** Caddy is installed from the official Cloudsmith APT repository, creating a systemd-managed service with automatic TLS via Let's Encrypt. No manual certificate management, certbot cron, or Nginx configuration is required.

**Docker remains installed but disabled at boot.** The GTFS-RT MobilityData validator wrapper requires Docker. Docker is installed on the OCI instance but `systemctl disable docker` prevents it from running continuously. Operators start Docker only when running a validation pass, then stop it to reclaim ~80–150 MB of idle RAM.

**4 GB swap on the boot volume.** The OCI boot volume is 47 GB SSD-backed block storage. A 4 GB swapfile with `vm.swappiness=10` is created as a safety net for memory pressure during Java-based static GTFS validation (200–400 MB spike), cold start, or concurrent request bursts. Swap is not a substitute for RAM — production workloads should not rely on it continuously.

**Systemd for service management.** All five application services run as the `open-transit` system user under individual systemd units in `deploy/systemd/`. Units use `ProtectSystem=strict`, `PrivateTmp=yes`, and `NoNewPrivileges=yes` for OS-level hardening. The `{{OCI_REMOTE_DIR}}`, `{{OCI_APP_USER}}`, and `{{DOMAIN}}` placeholders are substituted by `deploy/oci/install-units.sh` at install time.

**DuckDNS for pilot DNS.** The pilot domain `open-transit-pilot.duckdns.org` is updated by `scripts/oci-pilot.sh update-dns` using the DuckDNS API. Custom domains follow the same flow by updating an A record at the registrar. DNS must resolve to the OCI public IP before Caddy can obtain a TLS certificate via the HTTP-01 ACME challenge.

## ADR-0024 — Add a local Compose app profile for agency evaluation only

Phase 16 adds `scripts/agency-local-app.sh`, `deploy/Dockerfile.local`, `deploy/Caddyfile.local`, and a Docker Compose `app` profile so small-agency evaluators can run the full local stack without manually starting each Go service.

The local package is intentionally not the production deployment contract. It uses development defaults, local container networking, and `http://localhost:8080` as a convenience reverse proxy. Production deployments still require deployment-owned HTTPS/TLS, DNS, secret management, admin network boundaries, backup/restore policy, and monitoring.

`make agency-app-up` imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for service readiness, verifies public feed URLs, and prints next actions. The command does not fail solely because validator tooling is missing; validators remain an optional setup step unless a validation workflow is explicitly run.

The local image must not include `.cache`, local env files, generated tokens, private keys, or operator artifacts. Device-token rotation continues to use the existing `/admin/devices/rebind` API contract, which returns a one-time token by design.

## ADR-0025 — Use explicit dry-run-first pilot operations helpers

Phase 17 adds `scripts/pilot-ops.sh` and systemd timer examples for repeatable small-agency pilot operations. This is deployment/operator tooling, not a backend product feature.

The helper owns scheduled validation, backup, restore-drill, feed-monitor, and scorecard-export command sequences. It does not change service APIs, database schema, public feed URLs, GTFS-RT contracts, or consumer-submission statuses.

Every helper subcommand supports `--dry-run`, prints the target environment before doing work, and fails clearly when required target environment variables are missing. State-changing operations require explicit `ENVIRONMENT_NAME` and target paths or URLs. Restore operations are destructive for `RESTORE_DATABASE_URL` and require typed confirmation unless `--force` is passed.

Systemd examples use `EnvironmentFile=` and never inline live secrets. Raw backups, admin tokens, database URLs with passwords, webhook URLs, notification credentials, TLS private material, and unredacted operator artifacts are never public evidence. Missing notification destinations are recorded as `notification not configured`, not as feed failure.

Phase 17 evidence refresh must end with `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence`; refreshed evidence is not complete unless that audit passes. Passing evidence audit remains deployment/operator proof only and does not establish CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, or hosted SaaS availability.

## ADR-0026 — Keep device and AVL vendor integrations at the telemetry adapter boundary

Phase 25 documents device and AVL onboarding without adding a named vendor dependency or runtime adapter implementation. Vendor-specific payloads should be transformed into Open Transit RT telemetry events before calling `/v1/telemetry`.

Vendor credentials, private AVL payloads, private device identifiers, private vehicle identifiers, and private logs must stay outside this public repo unless reviewed and explicitly approved as public-safe. Vendor-specific assumptions must not be embedded into core matching, Vehicle Positions generation, or Trip Updates prediction logic.

Acceptable integration shapes include agency-owned adapter scripts, deployment-owned sidecar services, vendor-owned middleware, or private operator integration processes. These integrations must preserve the existing telemetry contract, validate required fields before forwarding, and avoid claiming certified vendor support or production AVL reliability without retained evidence.

## ADR-0027 — Keep telemetry simulation on the real ingest path

Phase 44 adds `cmd/telemetry-simulator` as local/reference operator tooling.
The simulator must post synthetic events through authenticated
`POST /v1/telemetry` with a real device bearer token; it must not write
directly to telemetry tables or call matcher internals as an ingest bypass.

Optional matcher diagnostics are allowed only after HTTP ingest succeeds. In
that mode, the simulator reads accepted rows back from Postgres, runs the
existing `internal/state.Engine`, and builds a private Vehicle Positions debug
snapshot with the same feed builder used by the service.

Simulator fixtures must remain synthetic-only. Diagnostics under `.cache/` are
private operator artifacts, not evidence packets, not vendor compatibility
proof, not production AVL reliability proof, not real realtime data proof, and
not production-grade ETA or CAL-ITP/Caltrans compliance proof.

## ADR-0028 — Evaluate external predictors only behind the prediction adapter

Phase 29A confirms that external predictors may be evaluated only behind `internal/prediction.Adapter`. The deterministic predictor remains the default runtime Trip Updates adapter, and Vehicle Positions generation remains independent of external predictor availability.

Runtime integration of TheTransitClock or any other external predictor requires a later approved phase, explicit dependency and license review, documented fallback behavior, health/failure semantics, and evidence appropriate to any compatibility or ETA-quality claim. Phase 29A mock adapter tests are contract tests only; they do not prove better ETAs, production-grade ETA quality, real-world predictor compatibility, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or vendor equivalence.

## ADR-0029 — Keep AVL vendor adapters outside core telemetry ownership

Phase 29B implements a synthetic, dry-run-only AVL/vendor adapter pilot as an example boundary, not as a runtime vendor integration. Vendor payload identifiers are lookup keys only; a strict mapping file is the authority for emitted Open Transit RT `agency_id`, `device_id`, and `vehicle_id`.

The adapter transforms records into the existing `telemetry.Event` contract and validates the transformed output against `telemetry.Event.Valid()`. It does not introduce a new telemetry request shape, change `/v1/telemetry`, change device token lifecycle, or add network send mode.

Diagnostics are dry-run review output. Duplicate and out-of-order diagnostics from the adapter are batch-level observations only and are not database ingest statuses. Partial stdout from a failed dry run is not submitted telemetry, vendor compatibility proof, production integration evidence, or AVL reliability evidence.

Future real vendor adapters must stay deployment-owned or isolated behind sidecar/middleware boundaries unless a later phase explicitly approves and documents a runtime integration. Vendor credentials, endpoint URLs, private IDs, and real AVL payloads must remain outside the public repo unless reviewed and explicitly approved as public-safe.

## ADR-0030 — Keep AVL adapter send mode on the existing telemetry boundary

Phase 48 adds an optional private send mode to `cmd/avl-vendor-adapter`, but
the only runtime ingest boundary remains authenticated `POST /v1/telemetry`.
The adapter does not introduce vendor-specific APIs, queues, schedulers,
webhooks, public unauthenticated routes, admin routes, or changes to
`/v1/telemetry` payload/auth semantics.

Send mode uses a strict `avl-adapter-send.v1` manifest with environment-only
token references. It preflights mapping, transform, manifest, target URL,
credential, timestamp, warning-gate, and output-path blockers before sending
any records. Stale or future transformed records block the whole send batch;
other warnings send by default unless `AVL_ADAPTER_FAIL_ON_WARNINGS=true`.

Private send diagnostics are written under `.cache/avl-vendor-adapter/` by
default or a safe non-evidence output override. They contain deterministic
`credential_ref` values, per-record outcomes, retry counts, safe parsed success
fields, and response SHA-256 hashes only. They are private operator
diagnostics, not evidence packets, vendor compatibility proof, production AVL
reliability proof, compliance proof, consumer acceptance proof, agency adoption
proof, hosted SaaS proof, or production-grade ETA proof.

## ADR-0031 — Keep generic external predictors behind sanitized HTTP DTOs

Phase 49 adds an optional `external_http` Trip Updates adapter and
`external_http_shadow` mode behind `internal/prediction.Adapter`. The default
adapter remains deterministic, and `noop` remains available.

The runtime endpoint is generic and fixed to `/v1/predict/trip-updates`.
Configuration requires an exact host allowlist, rejects userinfo/query/fragment
URLs, disables redirects, uses HTTPS except loopback test stubs, and enforces
timeout, request-size, and response-size caps. Bearer tokens are referenced by
uppercase environment variable name only; token values are never diagnostics or
persisted data.

External requests use dedicated sanitized DTOs. They never marshal
`telemetry.StoredEvent`, `telemetry.Event`, or `state.Assignment` directly, and
they exclude device IDs, driver IDs, raw payload JSON, score details, override
IDs, audit fields, raw override reasons, credentials, headers, and cookies.

`external_http` call failures degrade to valid empty Trip Updates with
`adapter_error` diagnostics. `external_http_shadow` keeps deterministic output
public and records only bounded redacted shadow counts/status. This phase does
not add TheTransitClock-specific runtime code, start external services, create
evidence, change consumer statuses, or make stronger ETA/compliance claims.

## ADR-0032 — Keep realtime quality backtesting private and aggregate-only

Phase 50 adds local backtesting for observed stop events and prediction samples
under `internal/realtimequality` plus `cmd/realtime-quality-backtest`. The
workflow is a private engineering diagnostic, not a runtime service.

Inputs are versioned local JSON files. Outputs are bounded aggregates and
redacted manifest metadata under `.cache/realtime-quality-backtest/` by
default. The command writes exactly `summary.json`, `summary.md`,
`metrics.json`, `metrics.md`, and `manifest.json`.

The backtest workflow must not persist rows in Postgres, add migrations, add an
Operations Console view, add public routes, block publishing, create evidence
packets, update consumer statuses, or copy raw operator input files. Maturity
gate labels remain diagnostic-only: `insufficient_data`, `diagnostic_pass`, and
`diagnostic_watch`.

## ADR-0033 — Keep operations reliability diagnostics private and existing-table only

Phase 51 adds private operations reliability diagnostics through authenticated
GET-only Operations Console routes and a local `.cache` script helper. Runtime
feed status uses only existing `feed_health_snapshot` rows, incident rollups
use only existing `incident` rows, and Vehicle Positions persists bounded
health snapshots into the existing `feed_health_snapshot` schema.

Missing data remains `missing` or `unknown`; it must not be promoted to `ok`.
Incident summaries are capped and sanitized, and they intentionally omit raw
`details_json`, raw payloads, private text, logs, tokens, hostnames, webhook
values, and backup dumps.

Vehicle Positions health persistence is best-effort and must not change public
feed response status when persistence fails. Phase 51 adds no public route,
migration, monitoring-stack dependency, notification delivery, evidence write,
consumer status change, publish gate, SLA claim, uptime guarantee,
production-readiness claim, compliance claim, hosted SaaS claim, agency
adoption claim, consumer acceptance claim, vendor compatibility claim, or
production-grade ETA claim.

## ADR-0034 — Gate retained final-root evidence behind explicit approval

Phase 52 adds dedicated final public root evidence tooling instead of reusing
hosted-pilot evidence scripts. Final-root evidence has a stricter claim
boundary because a placeholder hosted packet must not look like agency-owned or
agency-approved proof.

`scripts/collect-final-root-evidence.sh` defaults to ignored `.cache` storage
and writes blocker-only packets when no real final root and redacted approval
artifact exist. Writes under `docs/evidence/captured` require explicit
retention opt-in, `ALLOW_CAPTURED_EVIDENCE_WRITE=true`, a valid final root, and
a readable redacted approval artifact.

`scripts/audit-final-root-evidence.sh` has separate blocker and real audit
modes. Real audit must fail on missing approval, root mismatch, placeholders,
missing feed/checksum artifacts, missing or unavailable validator status,
unsafe private strings, checksum mismatch, missing redaction notes, or consumer
tracker drift. Phase 52 adds no runtime public route, consumer contact,
consumer status change, compliance claim, agency adoption claim, hosted SaaS
claim, production-readiness claim, SLA/uptime claim, vendor-compatibility
claim, consumer-acceptance claim, or production-grade ETA claim.

## ADR-0035 — Treat Phase 56 multi-agency work as route-boundary hardening only

Phase 56 adds validated `/public/agencies/{agency_id}/...` public feed routes
and a shared `internal/tenant` route parser. Existing single-agency public feed
routes remain service-instance scoped and unchanged.

Agency ID path segments use a conservative ASCII segment-safe contract. Public
per-agency protobuf/static/discovery routes are allowed, but per-agency public
JSON/debug routes are not added. Existing JSON debug routes remain
authenticated, and the OCI public Caddy edge exposes only public feed paths.

`scripts/multi-agency-hosting.sh` is a private `.cache` route/proxy diagnostic
only. It is not a tenant export, not a backup, not retained evidence, and not
production-readiness proof. Tenant restore into a shared live database remains
blocked until a later approved phase defines a safe tested contract.

This decision does not create a hosted SaaS, production multi-tenant hosting,
SLA/uptime, compliance, agency adoption, consumer acceptance, vendor
compatibility, marketplace approval, or production-grade ETA claim.

## ADR-0036 — Keep release packages local and auditable

Phase 57 adds local release package generation and audit tooling. The package
source archive is created from `git archive HEAD` so it does not recursively
copy `.env`, `.cache`, raw logs, private operator artifacts, generated evidence,
or other working-tree-only files.

Release packages include checksums, provenance metadata, Go-module SBOM
metadata, and optional local image metadata when an operator supplies an image
tag. The tooling does not build or push images by default, upload artifacts,
create GitHub releases, contact registries, sign attestations, or create
retained evidence.

Dirty packages are allowed for local diagnostics through the Make target and
are marked not release-ready. Actual release use should run from a clean
checkout with strict settings and an audited checksum manifest.

Phase 115 adds a root `.gitattributes` export policy so `git archive HEAD`
excludes retained/protected evidence and consumer-submission paths from source
release archives without modifying those protected repository paths. The
excluded paths remain tracked in the repository where they already existed;
the policy only controls public source archive distribution.

This decision does not create a hosted service, production image publication,
production-readiness, compliance, agency adoption, consumer acceptance, vendor
compatibility, marketplace approval, SLA/uptime, or production-grade ETA claim.

## ADR-0037 — Keep vendor-equivalent materials template-only

Phase 58 adds optional BYOD, support-boundary, KPI, and procurement templates
for operators who need vendor-style documentation during review.

The templates are not marketplace submissions, marketplace approval, vendor
compatibility, hardware certification, paid support, service commitments,
hosted service availability, production-readiness proof, compliance proof,
agency adoption proof, consumer acceptance proof, or production-grade ETA
proof.

`scripts/audit-vendor-equivalent-pack.sh` scans local template wording and the
prepared-only consumer tracker. It does not contact marketplaces, vendors,
consumers, agencies, procurement systems, or external services, and it does not
create retained evidence.

## ADR-0038 — Use sidecar manifests for connector plugins

Post-60 connector plugins are optional sidecars, command adapters, or
deployment-owned processes described by `open-transit-rt.connector.v1`
manifests. Open Transit RT validates the manifest and synthetic conformance
metadata locally, but it does not load arbitrary Go plugins or execute
connector code from manifests.

The contract covers telemetry sources, prediction sidecars, validator
wrappers, monitoring/export integrations, and consumer/discovery workflows.
All connectors remain adapter-bound, optional, documented, tested, redacted,
and fail-closed. Vehicle Positions must continue without an external
predictor, and deterministic prediction remains the default.

Manifest validation rejects raw secrets, private paths, unsafe URLs,
unsupported positive claims, raw validator commands, notification
send-by-default, consumer submission automation, and status mutation.

This decision adds no DB migration, public feed contract change, telemetry
ingest contract change, Trip Updates hard-coupling, external predictor default
change, evidence write, consumer status change, hosted SaaS claim, compliance
claim, consumer-acceptance claim, agency-adoption claim, vendor-compatibility
claim, production-readiness claim, production AVL reliability claim, or
production-grade ETA-quality claim.

## ADR-0039 — Keep adapter conformance offline and synthetic

Post-60 adapter conformance is an offline synthetic suite under
`testdata/adapter-conformance` with a Go CLI at `cmd/adapter-conformance`.
The suite checks connector manifests and required failure-mode cases for
telemetry sources, prediction sidecars, validator wrappers, and monitoring
exports.

The CLI validates fixture shape, required scenarios, allowlisted-validator
metadata, redaction/no-send assertions, and synthetic-only markers. It does
not start sidecars, send network traffic, run validators, contact consumers,
automate portals, write evidence, mutate repo state, or change runtime
configuration.

Passing conformance is a local connector-quality signal only. It is not
CAL-ITP/Caltrans compliance, consumer acceptance, agency approval, vendor
compatibility, production AVL reliability, production readiness, or
production-grade ETA-quality proof.

## ADR-0040 — Keep generic connector examples synthetic and offline

Post-60 generic connector examples live under `examples/connectors/` and show
the shape of telemetry source, prediction sidecar, and monitoring/export
connectors using `open-transit-rt.connector.v1` manifests, committed synthetic
fixtures, and small Go stdlib stubs.

The examples are included in local manifest, conformance, and `go test`
checks. They do not fetch real vendor endpoints, include real credentials,
include real vendor payloads, submit to consumers, mutate consumer status,
send notifications by default, write evidence, change runtime defaults, or
replace the authenticated `/v1/telemetry` and `internal/prediction.Adapter`
boundaries.

These examples are developer aids only. They are not vendor compatibility,
hardware certification, production AVL reliability, production-grade ETA
quality, consumer acceptance, agency approval, production readiness,
CAL-ITP/Caltrans compliance, hosted SaaS, paid support, or SLA proof.

## ADR-0041 — Keep the agency launchpad private and read-only

Post-60 agency launchpad work adds authenticated GET-only Operations Console
routes at `/admin/operations/launchpad` and
`/admin/operations/launchpad.json`. The launchpad derives setup, GTFS,
metadata, five-feed, telemetry, validator, readiness, connector conformance,
support-bundle, and decision-gate rows from existing private Operations
Console models and repo docs.

The launchpad adds no database table, migration, public route, public feed
contract, telemetry ingest contract, Trip Updates adapter change, POST action,
evidence write, consumer status mutation, external contact, portal automation,
or notification send path.

All launchpad claim flags remain false. The launchpad is a private operator
workflow aid only; it is not approval, agency adoption, CAL-ITP/Caltrans
compliance, consumer acceptance, final-root proof, public launch, hosted SaaS,
paid support, SLA coverage, production readiness, vendor compatibility,
production AVL reliability, or production-grade ETA proof.

## ADR-0042 — Keep Docker image distribution source/local-only in Phase 66

Phase 66 does not publish versioned production Docker images. The supported
distribution anchors remain source tags, exact commit SHAs, local release
packages, and deployment-owned local Docker builds from a reviewed checkout or
tag.

Maintainers and operators may build local images with `deploy/Dockerfile.local`
for evaluation, release-candidate review, or deployment-owned packaging. Those
local image tags are local metadata only. They are not registry-published
artifacts, hosted service availability, support commitments, or production
readiness proof.

A future image publication track would require an explicit maintainer decision
covering registry ownership, tag policy, supported architectures, build
provenance, vulnerability review, signing/attestation expectations, base-image
update policy, rollback policy, release-note language, and a no-secrets/no-
evidence artifact audit. Until that exists, do not push images, publish
registry tags, tell operators to pull a registry-published app image, or imply
a hosted SaaS or published app image exists.

## ADR-0043 -- Keep browser control behind a private command model

Phase 77 introduces `internal/admincontrol` as the bounded private command
result contract for Operations Console workflows. The model uses explicit
workflow statuses, an action ladder, bounded summaries, bounded next actions,
sanitized errors, and all-false claim flags by default.

The first migrated workflow is `validation_health.refresh`, a private read-only
refresh at `POST /admin/operations/validation-health/refresh.json`. It
recomputes validation-health summaries from existing private records and
server-owned artifact checks. It writes nothing, changes no public feed output,
creates no evidence, and moves no consumer status.

Browser commands must remain authenticated, role-checked, agency-scoped,
CSRF-safe for cookie-auth POSTs, request-capped, allowlisted, auditable, and
free of client-supplied shell commands, argument arrays, paths, validator
binaries, URLs, timeouts, support-bundle destinations, or evidence
destinations. Command responses must not include raw stdout, stderr, raw
validator reports, raw external payloads, private filesystem paths, bearer
tokens, cookies, database URLs, or private hostnames.

The command model is a private operator control-plane boundary only. It is not
a public API, shell runner, plugin loader, evidence collector, release
publisher, consumer submission system, compliance proof, production readiness
proof, hosted SaaS claim, vendor compatibility claim, SLA guarantee, or
production-grade ETA proof.

## ADR-0044 -- Keep Operations Console frontend enhancements buildless and private

Phase 78 adds a small progressive enhancement runtime for the private
Operations Console. The source of truth remains Go server-rendered HTML. The
browser runtime is a committed JavaScript source file embedded into the Go
binary and served only from an allowlisted authenticated route:
`/admin/operations/assets/operations.js`.

The asset route sets `Cache-Control: no-store`, an explicit JavaScript content
type, and `X-Content-Type-Options: nosniff`. It does not use `http.FileServer`,
filesystem path routing, a generated bundle, a package-manager install, a CDN,
or a public route.

Browser fetches are limited to relative private `/admin/operations/*.json`
reads and the approved read-only
`POST /admin/operations/validation-health/refresh.json` command. The runtime
does not fetch `/public/*`, `/v1/events`, external URLs, raw validator routes,
or public-looking GTFS-RT debug JSON. Browser storage may hold only UI
preferences such as filter or sort state; it must not store CSRF tokens, bearer
tokens, cookies, device tokens, raw JSON responses, row contents, commands,
URLs, private paths, hostnames, or credentials.

The runtime is progressive enhancement only. With JavaScript disabled, the
private console must still render useful tables, links, forms, visible
configured values, and native details. The frontend layer is not a public API,
evidence collector, consumer submission system, hosted SaaS surface,
production-readiness proof, vendor compatibility proof, SLA proof, or
production-grade ETA proof.

## ADR-0045 -- Keep the Operations Console information architecture route-driven

The post-rc2 Operations Console uses a central route registry for product
navigation, page titles, grouping, private/no-store posture, and stable route
inventory. Phase 02 reorganizes the console around agency tasks: Start Here,
Setup, GTFS Workbench, Feed Health, Validation, Realtime, Devices / AVL,
Prediction / ETA Lab, Connectors, Alerts, Readiness, Maintenance, Help /
Tutorials, and Support / Troubleshooting.

Route paths remain stable so existing operator links, tests, JSON companions,
and local runbooks do not break. Labels may become more user-facing, but route
inventory and access rules must stay explicit in code and tests.

Each page includes a visible "What to do next" action, and separate private
tools such as GTFS Studio and the Alerts Console are marked as separate tools
instead of hidden diagnostics. The console remains Go server-rendered HTML with
buildless progressive enhancement only. With JavaScript disabled, pages must
still render useful tables, links, forms, and native disclosure controls.

This information architecture does not create evidence, contact outside
systems, change consumer status, prove CAL-ITP/Caltrans compliance, prove
production readiness, prove consumer acceptance, prove vendor compatibility,
provide hosted service availability, or make SLA or production-grade ETA
claims.

## ADR-0046 -- Browser telemetry dry runs are fixture previews only

Phase 03 adds a browser dry-run preview to the private Telemetry Simulator.
The preview reads committed synthetic fixture metadata and renders a redacted
event summary for the selected scenario. It does not execute shell commands,
send telemetry, collect device tokens, read `.cache` diagnostics, write
database rows, or contact external systems.

Intentional telemetry sends remain outside the browser and use the existing
authenticated `/v1/telemetry` boundary from an operator or technical-helper
environment. One-time device tokens may be created or rotated through the
private Devices & Tokens page, but secure storage and installation on devices
or adapters remain operator-owned.

The preview exists to reduce command-line dependence for normal review: a
nontechnical agency user can inspect the synthetic scenario shape and expected
statuses from the browser after startup. It is not a vendor test, hardware
certification, real fleet reliability test, production AVL proof, consumer
acceptance signal, compliance proof, or production-grade ETA quality proof.

## ADR-0047 -- Keep GTFS Workbench read-only while adding staff-facing review

Phase 04 adds Agency Review Summary and Validation Issue Triage to the private
GTFS Workbench. The summary gives staff one browser-first view of required
files, row counts, service dates, core route/stop/trip coverage, import
history, active-vs-previous change signals, and current issue triage. The issue
triage reuses sanitized GTFS quality groups and shows likely owner,
plain-English meaning, suggested fix path, safe next action, and verification
path.

The Workbench remains read-only. It does not import GTFS, edit drafts, publish
schedules, run validators, execute rollback, create evidence, contact external
systems, change consumer status, or turn local review into approval,
compliance, consumer acceptance, production readiness, vendor compatibility, or
ETA-quality proof. Raw validator samples stay out of the Workbench JSON
companion; deeper issue analysis stays on the authenticated GTFS Quality page.

## ADR-0048 -- Keep realtime usefulness review private and read-only

Phase 05 adds feed-specific usefulness details and synthetic/local replay
guidance to the private Realtime Center. Vehicle Positions review summarizes
visible vehicle counts, estimated published rows, stale/unmatched/suppressed
vehicles, trip descriptor coverage, and why rows were not published. Trip
Updates review summarizes generated versus withheld output, prediction source,
fallback reason, stale/ambiguous inputs, and low-confidence handling. Alerts
review keeps lifecycle work in the Alerts Console while surfacing active,
stale, cancellation-link, disruption-link, and service-disruption review prompts.

This is a browser-first operator review layer only. It does not change
`/v1/telemetry`, matching, public GTFS-Realtime serialization, Trip Updates
adapter contracts, Alerts authoring, validator execution, evidence records,
consumer status, or release state. Local replay guidance starts with browser
fixture previews and keeps real sends, tokens, generated `.cache` diagnostics,
and private database access under operator or technical-helper control. These
signals do not prove compliance, consumer display, public launch, SLA, uptime,
production readiness, vendor compatibility, hardware certification, production
AVL reliability, production-grade ETA quality, or real-world ETA accuracy.

## ADR-0049 -- Treat connectors as a catalog plus conformance contract

Phase 06 adds one connector catalog across README, docs, private UI, examples,
and static site source. The catalog spells out Vehicle / GPS / AVL,
Prediction, Validator, Monitoring/export, Consumer/discovery, and Future
extension model paths, then maps each one to a copy/adapt starter and first
local check.

Connectors remain manifest-described sidecars or command adapters. Open
Transit RT still does not support arbitrary dynamic backend plugin loading.
Runtime interfaces stay explicit: authenticated `/v1/telemetry`, the
prediction adapter boundary, server-owned validator IDs, redacted monitoring
summaries, and public feed URL metadata.

Consumer/discovery is now covered by a synthetic example and adapter
conformance cases for feed URL metadata, submission blocking, and status
mutation blocking. This closes the gap where the manifest type existed but had
no copyable example or conformance group.

Connector checks remain local quality signals only. They do not contact
external systems, create evidence, move consumer status, prove compliance,
prove production readiness, prove consumer acceptance, prove vendor
compatibility, prove hardware certification, provide hosted service
availability, provide SLA coverage, prove production AVL reliability, or prove
production-grade ETA quality.

## ADR-0050 -- Keep CAL-ITP-style readiness as a browser workflow map

Phase 07 adds a ten-area readiness workflow map to the private Operations
Console. The map organizes public feed URLs, static GTFS, Vehicle Positions,
Trip Updates, Alerts, validation, license/contact metadata, uptime and
operations signals, telemetry/device state, and consumer preparedness before
the detailed readiness cards.

URL readiness and license/contact readiness are intentionally separate. A
configured public URL does not become externally ready merely because it is
present, and complete metadata does not prove final-root ownership, legal
approval, consumer review, or compliance. Feed URL review now depends on
validation or feed-health context before showing a ready status.

Consumer preparedness remains prepared-only. Runtime consumer workflow notes
can be displayed to operators, but they do not override the seven prepared docs
tracker targets or move any protected evidence status. The readiness page is a
private, read-only review surface; it does not contact external systems,
create evidence, change consumer status, prove CAL-ITP/Caltrans compliance,
prove production readiness, prove consumer acceptance, prove hosted service
availability, prove SLA or uptime, prove vendor or hardware compatibility, or
prove ETA quality.

## ADR-0051 -- Make human docs role-based before phase history

Phase 08 makes `README.md`, `docs/index.md`, `docs/README.md`, and the wiki
home point readers to task-based guides before maintainer phase ledgers. The
normal reader path is: understand the product, try it locally, open the private
browser UI, import GTFS, check feed URLs and feed health, connect vehicle data,
review readiness, and understand unsupported claims.

Long phase files, Codex task briefs, handoffs, and roadmap packs remain
discoverable for maintainers and AI agents, but they are no longer presented
as the first path for new users or agency staff. This keeps project history
available without making human readers interpret implementation ledgers before
they can evaluate the product.

## ADR-0052 -- Keep the public site static, local, and claim-safe

Phase 09 adds a static public site source under `site/` with a shared local CSS
file, a concise homepage, generated UI tour captures, connector catalog,
CAL-ITP-style readiness explainer, and video tutorial overview page. The site
uses plain HTML, CSS, and a small local script for role tabs. It adds no
external scripts, tracking, analytics, external fonts, or hosted-service
claims.

Generated browser captures are documentation aids only. They are not retained
evidence and do not prove public deployment, compliance, consumer acceptance,
agency adoption, production readiness, hosted service availability, SLA/uptime,
vendor compatibility, production AVL reliability, or ETA quality.

## ADR-0053 -- Keep tutorial video files out of the repository by default

Phase 10 adds `docs/tutorials/video-recording-guide.md` and expands the static
site video page with public-safe recording storyboards. The guide standardizes
six short tutorial scripts: overview, local setup, browser-first GTFS import,
feed health/readiness review, connector/AVL overview, and maintenance/support
workflow.

The repository stores scripts, checklists, and publication rules only. Raw or
finished video binaries should stay outside git unless a maintainer explicitly
authorizes release assets or another storage path. Recording workflows must use
local/demo or public-safe data, avoid secrets and private records, add captions
or transcripts before publication, and preserve all unsupported-claim limits.

## ADR-0054 -- Index AI-agent docs separately from human docs

Phase 11 adds `docs/agent/` as the explicit hub for Codex continuation context,
including handoffs, roadmap packs, prompt files, and historical phase ledgers.
The canonical historical files stay in their current locations so old links,
scripts, and release records keep working, but the normal reader path now
points agency staff and technical helpers to `README.md`, `docs/index.md`,
tutorials, connector docs, wiki pages, and the static site before agent
history.

This is a documentation separation, not a product or evidence status change.
It does not delete history, move protected evidence, change consumer tracker
state, create evidence packets, or make any compliance, production readiness,
consumer acceptance, hosted service, vendor compatibility, SLA, AVL
reliability, or ETA-quality claim.

## ADR-0055 -- Keep stable as a filtered product branch

Phase 12 introduces a `stable` branch policy and a GitHub Actions workflow that
filters `main` into `stable` while excluding AI-agent-only docs. The workflow
uses `.github/stable-sync-excludes.txt`, supports manual dry runs, commits only
when the filtered tree changes, and pushes without force. If the remote stable
branch diverges, the push should fail for maintainer review.

The branch is a product/user-facing source branch, not a stable release. It
does not publish a tag, contact external systems, create evidence, change
consumer tracker status, prove compliance, prove production readiness, prove
consumer acceptance, provide hosted service availability, prove vendor or
hardware compatibility, provide SLA coverage, or prove AVL/ETA quality.
