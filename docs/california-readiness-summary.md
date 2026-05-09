# California Readiness Summary

This Phase 20 summary separates Open Transit RT capability from deployment evidence, prepared submission packets, third-party submission status, and remaining gaps.

It supports California-facing readiness review. It does not claim CAL-ITP/Caltrans compliance, consumer ingestion, consumer acceptance, marketplace/vendor equivalence, hosted SaaS availability, agency endorsement, or production-grade ETA quality.

Agency-owned domain readiness is tracked separately in
`docs/agency-owned-domain-readiness.md`.

Official Caltrans reference points:

- [California Transit Data Guidelines, Version 4.0](https://dot.ca.gov/cal-itp/california-transit-data-guidelines)
- [California Transit Data Guidelines FAQ, Version 4.0](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs-v4_0)
- [Cal-ITP GTFS overview](https://dot.ca.gov/cal-itp/cal-itp-gtfs)
- [Critical GTFS Validation Errors](https://dot.ca.gov/cal-itp/critical-gtfs-validation-errors)
- [Website Model Language](https://dot.ca.gov/cal-itp/website-model-language)
- [FTA 2025 NTD Reporting Policy Manual](https://www.transit.dot.gov/ntd/2025-ntd-reporting-policy-manual)

Phase 54 re-checked these official public sources on May 9, 2026 and refreshed
requirement mappings only. The refresh did not create evidence, change
consumer statuses, or prove compliance.

## Phase 54 Official-Source Mapping

The current Caltrans Guidelines page identifies itself as Version 4.0, dated
December 11, 2024. The FAQ page identifies itself as Version 4.0.

Current official-source mappings recorded by Phase 54:

- GTFS Schedule compliance characteristics include public availability at a
  stable URL, regular canonical validation with no errors, acceptance by major
  trip planners, and an explicit open data license.
- GTFS Realtime compliance characteristics include public availability at
  stable URLs, regular canonical validation with no errors, acceptance by major
  trip planners, and an explicit open data license.
- Complete realtime availability includes all three standard GTFS Realtime
  feed types: Trip Updates, Vehicle Positions, and Service Alerts.
- Data availability expectations include provider or regional website
  source-of-truth links for GTFS Schedule and all three realtime feeds,
  technical contact or online contact routing, optional `feed_info.txt`
  feed-contact fields, and availability through Transitland and Mobility
  Database.
- If a deployment chooses API-key authentication for realtime feeds, the
  registration process must be easy to discover, quick, automated,
  transparent about terms and rate limits, and suitable for HTTPS requests.
- The FTA 2025 NTD Reporting Policy Manual requires applicable fixed-route NTD
  reporters to maintain a public-domain GTFS dataset and a publicly accessible,
  persistent, machine-readable, non-password-protected link for collecting the
  GTFS ZIP.

These mappings are requirements context only. They are not retained deployment
evidence and do not prove agency adoption, final-root readiness, consumer
acceptance, hosted SaaS availability, production readiness,
marketplace/vendor equivalence, SLA/uptime, vendor compatibility, or
production-grade ETA quality.

## Code-Complete Capability

The repository has code and workflow foundations for:

- static GTFS import and GTFS Studio draft publish;
- stable public feed paths for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts;
- GTFS-RT protobuf feed generation for all three realtime feed types;
- license/contact metadata through `feed_config`, `published_feed`, and `/public/feeds.json`;
- canonical validator command adapters and stored validation reports;
- consumer-ingestion workflow records;
- compliance scorecard snapshots;
- deterministic replay quality measurement for current realtime behavior.

These are implementation capabilities, not proof that any specific agency deployment is compliant or accepted by consumers.

## Productized Readiness Workflow

Phase 39 adds an authenticated Operations Console page at
`/admin/operations/readiness`. The page turns the evidence categories in this
summary into a plain-language operator checklist with a status source, current
signal, next action, and claim boundary for:

- stable public URLs;
- static GTFS;
- Vehicle Positions;
- Trip Updates;
- Alerts;
- license/contact metadata;
- validation;
- telemetry freshness;
- operations status;
- consumer packet preparedness.

The page supports CAL-ITP-style readiness workflows. It does not create
external evidence and does not claim CAL-ITP/Caltrans compliance, consumer
acceptance, agency adoption, final-root proof, hosted SaaS availability,
production readiness, vendor compatibility, or production-grade ETA quality.

Phase 55 adds a local packet export for this same evidence-bounded workflow.
`make generate-compliance-evidence-packet` writes ignored `.cache` blocker or
draft packets, and `make audit-compliance-evidence-packet` enforces false claim
flags, non-compliance statuses, redaction checks, prepared-only consumer
tracker state, README-only consumer artifact directories, and misleading-claim
guards. These packets are summaries for human review only. They do not create
retained evidence, contact consumers, fetch live feeds, change consumer
statuses, or prove compliance.

## Deployment-Proven Evidence

The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` provides hosted/operator evidence for the recorded pilot scope:

- anonymous public HTTPS fetches for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts;
- TLS and public-edge/private-admin boundary evidence;
- stable pilot feed URLs through controlled publish and restore drill snapshots;
- hosted validator records where schedule, Vehicle Positions, Trip Updates, and Alerts passed;
- scorecard export evidence with validation and discoverability green and consumer-ingestion red;
- monitoring, alert lifecycle, backup, restore, and job-history artifacts.

This is pilot deployment evidence for `https://open-transit-pilot.duckdns.org`, not agency-domain production proof.

## Phase 23 Agency-Domain Result

Phase 23 closed as blocker-documented only because no agency-owned or
agency-approved final public feed root was available.

No final-root public fetch evidence, DNS proof, TLS/redirect proof, validator
records, evidence packet, migration proof, or prepared packet refreshes were
collected. The OCI pilot remains labeled as hosted/operator pilot evidence only.

## Post-Phase-32 Final-Root Follow-Up Result

The post-Phase-32 final-root evidence follow-up also closed as
blocker-documented only. No agency-owned or agency-approved final public feed
root exists in repo evidence, no root was used, and no owner or approval
artifact was available.

No DNS proof, TLS certificate metadata, HTTP-to-HTTPS redirect proof, public
feed fetch proof, validator records, redacted proxy/config summary, evidence
packet README, or checksums were collected. No final-root evidence packet was
created. Prepared consumer packets and consumer target statuses were not
changed.

The DuckDNS OCI pilot remains pilot evidence only. It must not be described as
agency-owned final-root proof, consumer acceptance, or compliance evidence.

## Prepared Packet Evidence

Phase 20 prepared target-specific packets for:

- Google Maps;
- Apple Maps;
- Transit App;
- Bing Maps;
- Moovit;
- Mobility Database;
- transit.land.

The packet index and completeness checklist live at `docs/evidence/consumer-submissions/packets/README.md`.

The machine-readable status snapshot lives at `docs/evidence/consumer-submissions/status.json`.

All seven targets are `prepared` only. Prepared means a complete reviewable packet exists; it does not mean a submission was made.

The submission workflow lives at
`docs/evidence/consumer-submissions/submission-workflow.md`. It documents how
operators verify official paths, complete pre-submission checks, retain evidence,
and update statuses without overclaiming.

## Submitted Evidence

No submitted evidence is present.

No receipt, ticket, portal screenshot, email correspondence, or target-side artifact exists in the repository for any tracked target.

## Under-Review Evidence

No under-review evidence is present.

No consumer or aggregator has acknowledged review in the retained repo evidence.

## Accepted Evidence

No accepted evidence is present.

No consumer or aggregator acceptance, ingestion, listing, display, or production use may be claimed from the repository evidence.

## Missing Evidence Before Stronger Readiness Language

The following evidence remains missing before stronger CAL-ITP/Caltrans readiness or compliance language would be justified:

- agency-owned stable URL/domain proof; Phase 23 documented this as blocked
  because no agency-owned or provider-approved stable URL root is available,
  and the post-Phase-32 follow-up confirmed the same blocker remains;
- agency-approved identity, license, and contact metadata for any real provider submission;
- current production validation records for the final agency-owned or
  agency-approved URL root;
- provider or agreed regional source-of-truth pages that publish the canonical
  GTFS Schedule and all three GTFS Realtime links;
- public technical contact or online contact workflow for GTFS data quality
  matters, plus feed contact fields where supported;
- Transitland and Mobility Database availability records if discoverability
  through those aggregators is being claimed;
- documented API-key registration terms if realtime authentication is used;
- retained redacted submission receipts or tickets for each named consumer or aggregator;
- retained under-review, rejection, blocker, or acceptance evidence from each named consumer or aggregator when such status is claimed;
- consumer acceptance or ingestion proof for the exact feed scope and URL root being claimed;
- public provider page or discoverability metadata hosted by or approved for the agency;
- ongoing operations evidence for the final deployment environment, including monitoring, backup/restore, incident response, and validation cadence;
- production-grade ETA quality evidence beyond Phase 19 replay measurement.

## Safe Wording

Allowed:

- "Open Transit RT implements technical foundations for California transit data readiness."
- "The OCI pilot has hosted/operator evidence for public feed publication and validation."
- "Consumer submission packets have been prepared for review."

Not supported:

- "Open Transit RT is CAL-ITP compliant."
- "The feeds are accepted by Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, or transit.land."
- "The OCI pilot proves agency-owned production URL compliance."
- "Open Transit RT is marketplace/vendor equivalent."
- "Trip Updates have production-grade ETA quality."
