# Phase 27 — Multi-Agency Isolation Prototype

## Status

Complete for the Phase 27 prototype scope.

Phase 27 adds tests and documentation for selected agency-isolation boundaries. It does not claim production multi-tenant hosting, hosted SaaS availability, paid support or SLA coverage, agency endorsement, consumer acceptance, CAL-ITP/Caltrans compliance, marketplace/vendor equivalence, or production-grade ETA quality.

## What Phase 27 Proves

Phase 27 adds synthetic two-agency fixtures under `testdata/multi-agency/` and focused tests proving that selected workflows derive agency from trusted server/auth context instead of browser-supplied agency IDs.

Test-proven surfaces include:

- DB-backed auth role loading: a shared subject receives only roles for the JWT claim agency.
- Publication metadata, scorecard, validation, consumer, device-rebind, and Operations Console handler boundaries.
- Public `/public/feeds.json` query routing for synthetic agencies and omission behavior using configured `AGENCY_ID`.
- Public feed discovery JSON shape stays limited to public metadata fields.
- Device credential verification, device binding listing, and device rebind audit rows.
- Telemetry ingest token binding, debug event listing, latest-telemetry summaries, and Operations Console telemetry views.
- Compliance publication metadata, validation status, scorecards, consumer records, and audit rows.
- Alerts admin JSON and Alerts Console list/create/update/publish/archive/reconcile agency boundaries.
- GTFS Studio list/create/draft summary/publish/discard/entity edit boundaries for the current minimal Studio handler.
- Prediction operations records and audit rows for override/review workflows.

## Public Endpoint Scope

Current public endpoint routing is mixed:

| Endpoint | Current scope | Phase 27 status |
| --- | --- | --- |
| `/public/feeds.json` | Query-routed by `agency_id`; omitted query uses configured `AGENCY_ID`. | Handler-tested for query routing and public metadata only. |
| `/public/gtfs/schedule.zip` | Service-instance scoped by configured `AGENCY_ID` through the schedule builder. | Documented as single service-instance scoped. Query params are not per-agency routing. |
| `/public/gtfsrt/vehicle_positions.pb` | Service-instance scoped by configured `AGENCY_ID`. | Public protobuf behavior unchanged. Protected JSON debug now rejects cross-agency principals. |
| `/public/gtfsrt/trip_updates.pb` | Service-instance scoped by configured `AGENCY_ID`. | Public protobuf behavior unchanged. Protected JSON debug now rejects cross-agency principals. |
| `/public/gtfsrt/alerts.pb` | Service-instance scoped by configured `AGENCY_ID`. | Public protobuf behavior unchanged. Protected JSON debug now rejects cross-agency principals. |

Do not infer that one service instance can safely serve multiple public agency roots from Phase 27. Only `/public/feeds.json` has query-routed behavior today. Schedule ZIP and GTFS-RT protobuf feeds remain service-instance scoped unless a future phase adds explicit per-agency public routing, tests, and a migration plan.

## Agency-ID Persistence Audit

This is a current isolation review, not production multi-tenant certification.

Tables currently carrying `agency_id` include:

- `agency_user`, `role_binding`
- `device_credential`
- `feed_config`, `feed_version`, `published_feed`
- published GTFS tables such as `gtfs_agency`, `gtfs_route`, `gtfs_stop`, `gtfs_trip`, `gtfs_stop_time`, `gtfs_calendar`, `gtfs_calendar_date`, `gtfs_shape_point`, and `gtfs_frequency`
- GTFS Studio draft tables such as `gtfs_draft`, `gtfs_draft_agency`, and typed draft entity tables where agency-level scoping is represented through the draft and agency rows
- `gtfs_import`, `gtfs_draft_publish`
- `telemetry_event`
- `manual_override`, `vehicle_trip_assignment`, `incident`
- `validation_report`, `feed_health_snapshot`
- `consumer_ingestion`, `marketplace_gap`, `compliance_scorecard_snapshot`
- `service_alert`, `service_alert_informed_entity`
- `audit_log`

Global or shared objects that still require operational review before any hosted multi-tenant claim include:

- the `agency` root table and `goose_db_version`
- deployment environment variables such as `AGENCY_ID`, public roots, admin secrets, and validator paths
- backup/restore/export tooling
- evidence packet directories and generated operator artifacts
- validator temporary files and external tool outputs
- reverse proxy routing and public feed-root ownership
- any future audit-log read surface

## Remaining Gaps

True hosted multi-tenant service remains deferred. Missing or incomplete areas include:

- explicit per-agency routing for public schedule ZIP and GTFS-RT protobuf feeds;
- tenant-aware backup, restore, export, deletion, and restore-drill evidence;
- a reviewed global-admin model, if ever needed;
- broader GTFS Studio entity isolation with richer multi-draft fixtures;
- end-to-end browser tests for every Operations Console and Alerts Console role combination;
- deployment evidence proving one hosted stack can safely operate multiple agency roots;
- evidence packet and consumer packet generation that proves per-agency redaction and artifact separation.

## Consumer And Evidence Boundary

Runtime DB consumer records remain agency-scoped operational records. They do not override the Phase 20 docs/evidence tracker. The tracker and `docs/evidence/consumer-submissions/status.json` remain prepared-only records unless retained, redacted, target-originated evidence supports a target-specific status change.

Evidence packets for future multi-agency operation must be generated and redacted per agency. One agency's validation, consumer packet, device telemetry, incident, scorecard, or operator evidence must not be copied into another agency packet.
