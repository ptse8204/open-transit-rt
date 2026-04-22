# Compliance Evidence Checklist

This document is the Phase 11 evidence package for Open Transit RT. It separates what the repository proves locally from what a real agency deployment must prove and what only third-party consumers can confirm.

It uses the truthfulness guardrail in `docs/prompts/calitp-truthfulness.md`: the repo may be described as supporting deployment toward Caltrans/CAL-ITP-style readiness, but it must not be described as fully compliant, consumer-accepted, production ready for all agencies, or equivalent to a marketplace vendor without external evidence.

## External Reference Points

Use these official sources when discussing California-facing readiness:

- [Caltrans California Transit Data Guidelines v3.1](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-v3_1)
- [Caltrans California Transit Data Guidelines FAQ](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs-v3_0)
- [Cal-ITP GTFS overview](https://dot.ca.gov/cal-itp/cal-itp-gtfs)
- [GTFS Realtime Best Practices](https://gtfs.org/documentation/realtime/realtime-best-practices/)

The Caltrans guidelines describe GTFS Schedule and GTFS Realtime compliance in terms that include stable public URLs, regular canonical validation with no errors, open licensing, and ingestion by major trip planners. For realtime completeness, the guidelines identify all three standard GTFS Realtime feed types: Trip Updates, Vehicle Positions, and Alerts.

## Evidence Categories

| Area | Implemented in repo | Requires deployment/operator proof | Requires third-party confirmation |
| --- | --- | --- | --- |
| Stable public feed URLs | Stable public paths exist for `/public/gtfs/schedule.zip`, `/public/feeds.json`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/alerts.pb`. `published_feed` stores canonical URLs. | Public HTTPS host, reverse proxy routing, anonymous fetch proof, URL permanence across schedule updates and rollback. | Consumers must confirm they can fetch the deployment URLs if acceptance is claimed. |
| Public publication | Public protobuf and schedule ZIP endpoints are anonymous by design; the Phase 10 demo fetches them through a local public proxy. | Internet-reachable HTTPS deployment, no login wall, uptime evidence, cache/header behavior verified against the live host. | Major consumers must confirm successful automated fetches when claiming ingestion. |
| Open license and contact metadata | `feed_config`, `published_feed`, `/public/feeds.json`, and scorecards carry license/contact fields. | Agency-approved open data license, monitored technical contact, provider website or metadata page exposing those values publicly. | Consumers or aggregators may request confirmation that license/contact metadata is acceptable. |
| Static GTFS Schedule | ZIP import, GTFS Studio publish, active feed versions, and `/public/gtfs/schedule.zip` are implemented from database-backed published GTFS. | Current active agency schedule, public fetch evidence, canonical validator result for the deployed schedule, operational publish/rollback procedure. | Trip planners must accept or ingest the static GTFS feed before acceptance can be claimed. |
| Vehicle Positions | DB-backed GTFS-RT Vehicle Positions protobuf generation from latest accepted telemetry plus current assignments is implemented. | Real device telemetry, freshness monitoring, validator result for the deployed feed, evidence that stale/unmatched behavior matches agency policy. | Trip planners must accept or ingest the deployed Vehicle Positions feed if that is claimed. |
| Trip Updates | Stable Trip Updates endpoint and internal deterministic `internal/prediction.Adapter` implementation exist; weak or unsupported cases are withheld. | Real operating data, coverage review, validator result, quality review, and agency approval that conservative schedule-deviation predictions are acceptable for the pilot. | Trip planners must accept or ingest the deployed Trip Updates feed; production-grade ETA quality requires additional evidence beyond repo tests. |
| Alerts | Persisted Service Alerts authoring/lifecycle state and public GTFS-RT Alerts publication are implemented. | Operator workflow proof, live alert lifecycle evidence, validator result, and process for cancellations, disruptions, and expired alerts. | Consumers must accept or ingest the deployed Alerts feed if that is claimed. |
| Canonical validator workflow | Pinned MobilityData static validator install/check path, Docker-backed pinned GTFS-RT validator wrapper, allowlisted `/admin/validation/run`, and normalized `validation_report` records are implemented. | Latest production validation records for schedule, Vehicle Positions, Trip Updates, and Alerts; no-error results before making compliance claims. | Consumer acceptance is separate from validation and must not be inferred from validator success alone. |
| Consumer-ingestion workflow records | `consumer_ingestion` records and admin APIs track workflow status and packet JSON. Default seeded consumers are Google Maps, Apple Maps, Transit App, Bing Maps, and Moovit. | Actual submitted packet, submission dates, rejection/accepted notes, and operator-maintained records. Mobility Database and transit.land may be tracked as workflow records but are not seeded default integrations. | Only the named consumer or aggregator can confirm acceptance, ingestion, or production use. |
| Deployment, security, and operations | Production secret checks, admin JWT/cookie auth, CSRF for browser unsafe methods, device token binding, request IDs, readiness checks, and optional `/metrics` output are implemented. | HTTPS/TLS, backups, process supervision, log retention, monitoring/alerting, incident response, key rotation, role assignments, and deployment runbooks. | Third parties do not prove these items except where a consumer requires operational evidence. |
| Marketplace or vendor-equivalent capability | The repo records marketplace gaps and supports technical workflow evidence. | Service packaging, support runbooks, onboarding templates, SLA/KPI reporting, hardware/BYOD strategy, procurement artifacts, and operations staffing. | Marketplace listing, vendor approval, or consumer partnership status must come from the relevant external program or customer. |

## California Readiness Mapping

Open Transit RT currently supports the technical foundations for the Caltrans/CAL-ITP-style data expectations below:

- Stable URL foundation: implemented as stable paths plus `published_feed.canonical_public_url`; deployment must prove real HTTPS permanence.
- Public GTFS Schedule and GTFS Realtime publication: implemented locally for schedule, Vehicle Positions, Trip Updates, and Alerts; deployment must prove public reachability and no login wall.
- Canonical validation workflow: implemented through pinned tooling and allowlisted validator IDs; deployment must run and store current no-error validation results before any compliant wording is justified.
- Open license/contact metadata: implemented as metadata storage and `/public/feeds.json`; deployment must provide agency-approved values and public visibility.
- Consumer-ingestion workflow: implemented as records and packet storage; third-party acceptance must come from the specific consumer or aggregator.
- Realtime completeness: all three standard GTFS-RT feed paths exist; deployment must prove live, validator-clean, fresh, and operationally useful feeds.

Truthful wording:

- Allowed: "supports deployment toward Caltrans/CAL-ITP-style readiness."
- Allowed: "implements technical foundations for stable URLs, validation workflow, license/contact metadata, and consumer-ingestion records."
- Not supported by repo-only evidence: "is CAL-ITP compliant," "is accepted by Google Maps or Apple Maps," "is production ready for every agency," or "is a marketplace vendor equivalent."

## External Integration Reality

Integrated and testable in the current repo:

- Postgres/PostGIS through Docker Compose or deployment-owned database provisioning.
- pgx repository access and Goose migrations.
- GTFS Realtime protobuf bindings in feed boundary packages.
- MobilityData static GTFS validator pin/install/check workflow.
- Docker-backed MobilityData GTFS-RT validator wrapper pinned by image digest.
- Docker Compose for local Postgres/PostGIS and validator wrapper support.
- Optional Taskfile mirror for Makefile workflows.
- Local demo tools: `curl`, `zip`, and `unzip`.
- Internal Prometheus-format `/metrics` output when `METRICS_ENABLED=true`.

Optional or deferred, not currently integrated as external systems:

- TheTransitClock: deferred. A real integration would require an adapter behind `internal/prediction.Adapter`, setup docs, input/output contract tests, smoke coverage, failure behavior, and replacement strategy.
- Other external predictors: deferred behind the same adapter boundary.
- Prometheus/Grafana deployment: deferred. The repo emits internal metrics text, but it does not provision a Prometheus server, Grafana dashboards, alert rules, or production SLO evidence.
- OpenTelemetry: deferred. Phase 11 repo search found no OpenTelemetry SDK, collector, exporter, trace propagation, or deployment docs.
- Consumer submission APIs for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land: not integrated. The repo stores workflow records only.

## Phase 11 Evidence Result

Phase 11 does not add new backend features or external adapters. It closes the evidence layer by documenting what the repo proves, what a deployment must prove, and what third parties must confirm.

The next hardening track should collect deployment evidence: real HTTPS feed root, production validation records, live scorecard export, monitoring and alerting assets, operations runbooks, and third-party submission or acceptance records.
