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

Phase 8 adds canonical validator command adapters for static GTFS and GTFS-RT. Validator results are normalized into `validation_report`. If validator tooling is absent, the system stores `status='not_run'` instead of pretending validation passed.

Production mode is agency-scoped and stored as `feed_config.publication_environment = 'production'`. In production mode, missing canonical validator execution makes scorecard validation red. In dev mode, missing validators are yellow/not-run. `validator_strictness` controls failure handling, but it does not define production mode by itself.
