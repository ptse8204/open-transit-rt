# Compliance Evidence Checklist

This document is the Phase 11 evidence package for Open Transit RT. It separates what the repository proves locally from what a real agency deployment must prove and what only third-party consumers can confirm.

It uses the truthfulness guardrail in `docs/prompts/calitp-truthfulness.md`: the repo may be described as supporting deployment toward Caltrans/CAL-ITP-style readiness, but it must not be described as fully compliant, consumer-accepted, production ready for all agencies, or equivalent to a marketplace vendor without external evidence.

![Illustrative evidence maturity ladder showing code exists, hosted evidence, prepared packet, submitted, under review, and accepted as separate evidence stages.](assets/evidence-maturity-ladder.png)

## External Reference Points

Use these official sources when discussing California-facing readiness:

- [Caltrans California Transit Data Guidelines, Version 4.0](https://dot.ca.gov/cal-itp/california-transit-data-guidelines)
- [Caltrans California Transit Data Guidelines FAQ, Version 4.0](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs-v4_0)
- [Cal-ITP GTFS overview](https://dot.ca.gov/cal-itp/cal-itp-gtfs)
- [Caltrans Critical GTFS Validation Errors](https://dot.ca.gov/cal-itp/critical-gtfs-validation-errors)
- [Caltrans Website Model Language](https://dot.ca.gov/cal-itp/website-model-language)
- [FTA 2025 NTD Reporting Policy Manual](https://www.transit.dot.gov/ntd/2025-ntd-reporting-policy-manual)

Phase 54 re-checked these official public sources on May 9, 2026. The current
Caltrans Guidelines page identifies itself as Version 4.0, dated December 11,
2024, and the FAQ identifies itself as Version 4.0. This source refresh updates
requirement mapping only; it is not compliance evidence.

Phase 60 completed a final claim review against the retained repository
evidence and this official-source context. The review and audit tooling are not
deployment evidence and do not prove compliance.

The current Caltrans California Transit Data Guidelines describe GTFS Schedule
and GTFS Realtime compliance in terms that include stable public URLs, regular
canonical validation with no errors, open licensing, and ingestion by major trip
planners. For realtime completeness, that context includes all three standard
GTFS Realtime feed types: Trip Updates, Vehicle Positions, and Service Alerts.
The Guidelines also map data availability to provider or regional
source-of-truth website links, a technical contact or online contact workflow,
aggregator availability through Transitland and Mobility Database, and
transparent API-key registration constraints if realtime authentication is
used. The FTA 2025 NTD manual separately requires applicable fixed-route NTD
reporters to maintain a public-domain GTFS dataset and a publicly accessible,
persistent, machine-readable, non-password-protected GTFS ZIP link.

## Evidence Categories

| Area | Implemented in repo | Requires deployment/operator proof | Requires third-party confirmation |
| --- | --- | --- | --- |
| Stable public feed URLs | Stable public paths exist for `/public/gtfs/schedule.zip`, `/public/feeds.json`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, and `/public/gtfsrt/alerts.pb`. `published_feed` stores canonical URLs. | Public HTTPS host, reverse proxy routing, anonymous fetch proof, URL permanence across schedule updates and rollback. If realtime API keys are used, registration must be discoverable, automated, transparent, and nondiscriminatory. | Consumers must confirm they can fetch the deployment URLs if acceptance is claimed. |
| Public publication | Public protobuf and schedule ZIP endpoints are anonymous by design; the Phase 10 demo fetches them through a local public proxy. | Internet-reachable HTTPS deployment, no login wall, uptime evidence, cache/header behavior verified against the live host. | Major consumers must confirm successful automated fetches when claiming ingestion. |
| Open license and contact metadata | `feed_config`, `published_feed`, `/public/feeds.json`, and scorecards carry license/contact fields. | Agency-approved open data license, monitored technical contact, provider or agreed regional website exposing feed links and metadata publicly, and feed contact fields where supported by the active GTFS feed. | Consumers or aggregators may request confirmation that license/contact metadata is acceptable. |
| Static GTFS Schedule | ZIP import, GTFS Studio publish, active feed versions, and `/public/gtfs/schedule.zip` are implemented from database-backed published GTFS. | Current active agency schedule, public fetch evidence, canonical validator result for the deployed schedule, operational publish/rollback procedure. | Trip planners must accept or ingest the static GTFS feed before acceptance can be claimed. |
| Vehicle Positions | DB-backed GTFS-RT Vehicle Positions protobuf generation from latest accepted telemetry plus current assignments is implemented. | Real device telemetry, freshness monitoring, validator result for the deployed feed, evidence that stale/unmatched behavior matches agency policy. | Trip planners must accept or ingest the deployed Vehicle Positions feed if that is claimed. |
| Trip Updates | Stable Trip Updates endpoint and internal deterministic `internal/prediction.Adapter` implementation exist; weak or unsupported cases are withheld. | Real operating data, coverage review, validator result, quality review, and agency approval that conservative schedule-deviation predictions are acceptable for the pilot. | Trip planners must accept or ingest the deployed Trip Updates feed; production-grade ETA quality requires additional evidence beyond repo tests. |
| Alerts | Persisted Service Alerts authoring/lifecycle state and public GTFS-RT Alerts publication are implemented. | Operator workflow proof, live alert lifecycle evidence, validator result, and process for cancellations, disruptions, and expired alerts. | Consumers must accept or ingest the deployed Alerts feed if that is claimed. |
| Canonical validator workflow | Pinned MobilityData static validator install/check path, Docker-backed pinned GTFS-RT validator wrapper, allowlisted `/admin/validation/run`, and normalized `validation_report` records are implemented. | Latest production validation records for schedule, Vehicle Positions, Trip Updates, and Alerts; no-error results before making compliance claims. | Consumer acceptance is separate from validation and must not be inferred from validator success alone. |
| Consumer-ingestion workflow records | `consumer_ingestion` records and admin APIs track workflow status and packet JSON. Default seeded consumers are Google Maps, Apple Maps, Transit App, Bing Maps, and Moovit. | Actual submitted packet, submission dates, rejection/accepted notes, and operator-maintained records. Mobility Database and Transitland availability require retained records before being claimed. | Only the named consumer or aggregator can confirm acceptance, ingestion, listing, display, or production use. |
| Deployment, security, and operations | Production secret checks, admin JWT/cookie auth, CSRF for browser unsafe methods, device token binding, request IDs, readiness checks, and optional `/metrics` output are implemented. | HTTPS/TLS, backups, process supervision, log retention, monitoring/alerting, incident response, key rotation, role assignments, and deployment runbooks. | Third parties do not prove these items except where a consumer requires operational evidence. |
| Readiness workflow | `/admin/operations/readiness` shows CAL-ITP-style readiness rows with status source, current evidence/signal, next action, and claim boundary. | Operators must supply deployment-specific URL, validation, metadata, telemetry, operations, and consumer artifacts before making stronger claims. | Third parties must still confirm submission, review, acceptance, ingestion, listing, or display when those statuses are claimed. |
| Marketplace or vendor-equivalent capability | The repo records marketplace gaps and supports technical workflow evidence. | Service packaging, support runbooks, onboarding templates, SLA/KPI reporting, hardware/BYOD strategy, procurement artifacts, and operations staffing. | Marketplace listing, vendor approval, or consumer partnership status must come from the relevant external program or customer. |

## California Readiness Mapping

Open Transit RT currently supports the technical foundations for the Caltrans/CAL-ITP-style data expectations below:

- Stable URL foundation: implemented as stable paths plus `published_feed.canonical_public_url`; deployment must prove real HTTPS permanence.
- Public GTFS Schedule and GTFS Realtime publication: implemented locally for schedule, Vehicle Positions, Trip Updates, and Alerts; deployment must prove public reachability and no login wall.
- Canonical validation workflow: implemented through pinned tooling and allowlisted validator IDs; deployment must run and store current no-error validation results before any compliant wording is justified.
- Open license/contact metadata: implemented as metadata storage and `/public/feeds.json`; deployment must provide agency-approved values, provider or regional source-of-truth website visibility, and a technical contact or online contact workflow.
- Consumer-ingestion workflow: implemented as records and packet storage; third-party acceptance must come from the specific consumer or aggregator.
- Realtime completeness: all three standard GTFS-RT feed paths exist; deployment must prove live, validator-clean, fresh, and operationally useful feeds.
- Aggregator availability: workflow records and prepared packet docs exist; deployment must retain records before claiming Transitland or Mobility Database availability.

Phase 55 adds a local compliance/readiness packet generator for operator review:

- `make generate-compliance-evidence-packet` writes ignored `.cache`
  blocker/draft packets only.
- `COMPLIANCE_PACKET_DIR=<packet> make audit-compliance-evidence-packet`
  checks packet shape, JSON validity, false claim flags, non-compliance
  statuses, unsafe strings, prepared-only consumer tracker state, README-only
  artifact directories, and misleading claim wording.
- Generated packets summarize configured retained evidence paths and missing
  evidence. They do not create evidence, fetch feeds, contact consumers,
  change consumer statuses, or prove compliance.

Phase 60 adds `make audit-final-claim-review` for bounded public/status claim
review. The audit checks required Phase 60 sections, unsafe private strings,
unsupported positive claims, prepared-only consumer tracker state, and
README-only consumer artifact directories. It is local and read-only.

Truthful wording:

- Allowed: "supports deployment toward Caltrans/CAL-ITP-style readiness."
- Allowed: "supports CAL-ITP-style readiness workflows."
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

## Phase 12 Step 2 Evidence Result

Phase 12 Step 2 collected a dated local evidence packet at `docs/evidence/captured/local-demo/2026-04-22/`.

What it proves:

- local loopback public feed retrieval for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts;
- local protected admin/debug route rejection for anonymous requests;
- local validator workflow execution for schedule plus all three realtime feeds, with failures retained without omission;
- local manual scorecard export;
- local Postgres dump/restore mechanics into an isolated restore database, with restored row counts and public feed fetches against the restored database.

What it does not prove:

- public HTTPS hosting, TLS, or DNS ownership;
- clean validator status;
- production monitoring, alert delivery, or alert lifecycle;
- production backup schedule, retention, or restore operations;
- production publish/rollback URL permanence;
- third-party consumer or aggregator acceptance.

The local packet is useful operator/repo evidence, but it does not support stronger CAL-ITP compliance, production-readiness, or consumer-acceptance claims.

## Phase 12 Step 3 Tooling Result

Phase 12 Step 3 did not collect new hosted evidence. It hardened the validator-evidence path so closure checks are stricter:

- the repo-supported GTFS-RT validator wrapper now drives the pinned MobilityData validator webapp API against local schedule and realtime artifacts instead of passing unsupported CLI flags to the pinned image;
- `make validators-check` now fails if Java is not runnable for the pinned static validator JAR;
- Docker, `curl`, and `python3` are now explicit requirements for the repo-supported GTFS-RT validator wrapper.

This improved future hosted evidence collection, but did not itself move hosted deployment items into the deployment/operator-proof column.

## Phase 12 Hosted OCI Pilot Evidence Result

Phase 12 hosted evidence was collected at `docs/evidence/captured/oci-pilot/2026-04-24/`.

The hosted packet includes public HTTPS feed fetches, public-edge and SSH-tunneled auth-boundary checks, TLS/redirect evidence, clean hosted validator records for schedule plus all three realtime feeds, monitoring and alert lifecycle artifacts, backup/job history artifacts, a deployment data-restore rollback drill, scorecard export job-history evidence, and a final current-live recheck showing active `gtfs-import-3` with `canonical_validation_complete=true`.

The closure command passed:

```bash
EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence
```

This closes Phase 12 for deployment/operator evidence on the OCI pilot only. It does not prove Cal-ITP compliance or third-party consumer acceptance.

## Phase 13 Consumer Submission Evidence Result

Phase 13 adds a truthful consumer-submission evidence layer under `docs/evidence/consumer-submissions/` and `docs/consumer-submission-evidence.md`.

What Phase 13 currently proves:

- a tracker exists for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land;
- each target has a current evidence record and reusable template;
- each current record is `not_started` because no redacted third-party submission, review, acceptance, rejection, or blocker evidence is present in the repository;
- each template requires status timestamp, operator, feed root, exact submitted URLs, submission packet artifact, validation reference, Phase 12 evidence packet reference, correspondence/receipt/ticket/screenshot reference, redaction notes, next action, allowed public wording, and acceptance-scope fields.

What Phase 13 does not prove:

- submission to any consumer or aggregator;
- review by any consumer or aggregator;
- acceptance, ingestion, listing, or production use by any consumer or aggregator;
- CAL-ITP compliance.

Hosted/operator evidence completed in Phase 12 remains separate from consumer submission evidence collected in Phase 13. Validator success and public fetch proof are supporting evidence only, not consumer acceptance.

## Phase 20 Consumer Submission Packet Result

Phase 20 adds target-specific prepared packets under `docs/evidence/consumer-submissions/packets/` and a machine-readable status snapshot at `docs/evidence/consumer-submissions/status.json`.

What Phase 20 currently proves:

- complete packet drafts exist for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land;
- each prepared packet includes all five public feed URLs, license/contact metadata, Phase 12 hosted evidence, validator evidence, redaction notes, next action, and allowed wording;
- each packet includes evidence freshness fields and marks official submission method/contact as `not verified`;
- all seven current records are `prepared` only.

What Phase 20 does not prove:

- submission to any consumer or aggregator;
- review by any consumer or aggregator;
- acceptance, ingestion, listing, display, or production use by any consumer or aggregator;
- agency-owned domain proof, because the OCI pilot DuckDNS domain is pilot evidence and not agency-domain production proof;
- CAL-ITP/Caltrans compliance;
- marketplace/vendor equivalence;
- production-grade ETA quality.

The California readiness summary is `docs/california-readiness-summary.md`. The marketplace/vendor gap review is `docs/marketplace-vendor-gap-review.md`.

## Phase 23 Agency-Owned Deployment Proof Result

Phase 23 closed as blocker-documented only because no agency-owned or
agency-approved final public feed root was available.

What Phase 23 currently proves:

- the final-root evidence gap is explicitly documented;
- the DuckDNS OCI pilot remains labeled as hosted/operator pilot evidence only;
- the future agency-owned domain checklist and migration template are recorded.

What Phase 23 does not prove:

- agency-owned or agency-approved domain/root ownership;
- DNS, TLS, redirect, or public fetch behavior at a final root;
- validator-clean final-root schedule, Vehicle Positions, Trip Updates, or
  Alerts records;
- refreshed prepared packets for a final root;
- consumer submission, review, acceptance, ingestion, listing, display, or
  production use;
- CAL-ITP/Caltrans compliance;
- marketplace/vendor equivalence;
- production-grade ETA quality.

## Post-Phase-32 Final-Root Follow-Up Result

The post-Phase-32 final-root evidence follow-up also closed as
blocker-documented only. No agency-owned or agency-approved final public feed
root was available, no root was used, and no owner or approval evidence was
available.

No final-root DNS proof, TLS certificate metadata, HTTP-to-HTTPS redirect
proof, public fetch proof for the five public feed URLs, validator records,
redacted proxy/config summary, evidence packet README, or checksums were
collected. No final-root evidence packet was created, and
`EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` was intentionally not
run.

This follow-up does not change consumer or aggregator statuses. All target
records remain `prepared` only unless retained, redacted, target-originated
evidence supports a target-specific status transition.
