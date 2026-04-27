# California Readiness Summary

This Phase 20 summary separates Open Transit RT capability from deployment evidence, prepared submission packets, third-party submission status, and remaining gaps.

It supports California-facing readiness review. It does not claim CAL-ITP/Caltrans compliance, consumer ingestion, consumer acceptance, marketplace/vendor equivalence, hosted SaaS availability, agency endorsement, or production-grade ETA quality.

Agency-owned domain readiness is tracked separately in
`docs/agency-owned-domain-readiness.md`.

Official Caltrans reference points:

- [California Transit Data Guidelines v4.0](https://dot.ca.gov/cal-itp/california-minimum-general-transit-feed-specification-gtfs-guidelines)
- [California Transit Data Guidelines FAQ](https://dot.ca.gov/cal-itp/california-transit-data-guidelines-faqs)

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

## Deployment-Proven Evidence

The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` provides hosted/operator evidence for the recorded pilot scope:

- anonymous public HTTPS fetches for `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts;
- TLS and public-edge/private-admin boundary evidence;
- stable pilot feed URLs through controlled publish and restore drill snapshots;
- hosted validator records where schedule, Vehicle Positions, Trip Updates, and Alerts passed;
- scorecard export evidence with validation and discoverability green and consumer-ingestion red;
- monitoring, alert lifecycle, backup, restore, and job-history artifacts.

This is pilot deployment evidence for `https://open-transit-pilot.duckdns.org`, not agency-domain production proof.

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

- agency-owned stable URL/domain proof; the DuckDNS OCI pilot is useful pilot evidence, but an agency-owned domain or provider-approved stable URL remains unproven;
- agency-approved identity, license, and contact metadata for any real provider submission;
- current production validation records for the final agency-owned URL root;
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
