# CAL-ITP Readiness Plain English

Open Transit RT supports CAL-ITP-style readiness workflows. That means the
software helps operators evaluate and prepare the technical pieces commonly
needed for public GTFS and GTFS Realtime publication.

It does not mean this repository, a local demo, or a private trial is
CAL-ITP/Caltrans compliant.

## What The Software Helps With

- public feed URLs for `feeds.json`, static GTFS, Vehicle Positions, Trip
  Updates, and Alerts;
- static GTFS import, review, quality guidance, and validation triage;
- Vehicle Positions review with telemetry/device state and conservative
  assignment signals;
- Trip Updates review through the pluggable prediction boundary;
- Alerts feed and alert lifecycle review;
- license and technical-contact metadata workflows;
- validator workflow records;
- uptime and operations signals through local feed-health, reliability, and
  maintenance views;
- prepared consumer/aggregator packet records;
- connector and support-bundle review paths.

## What A Deployment Still Owns

A real public deployment still needs deployment-specific work:

- hosting, DNS, TLS, and reverse proxy configuration;
- current agency-approved license and contact metadata;
- validator runs against the deployed feeds;
- monitoring, backups, restore practice, and incident handling;
- provider or regional source-of-truth pages;
- consumer or aggregator submissions when the agency chooses that path.

## Readiness Matrix

| Area | UI signal you can review | Missing deployment evidence before stronger claims |
| --- | --- | --- |
| Public feed URLs | Readiness workflow map, Feeds page, Feed Health, and `/public/feeds.json` review. | Final-root/source-of-truth listing, public HTTPS fetch proof, target-specific review, and consumer-originated outcomes. |
| Static GTFS | Browser import, active schedule, public `schedule.zip`, GTFS quality guidance, and validator-health rows. | Agency-approved final dataset, public HTTPS fetch proof, current validator record, license/contact approval, and source-of-truth website listing. |
| Vehicle Positions | Public Vehicle Positions feed path, feed-health row, telemetry freshness, and readiness row. | Real device telemetry, deployed-feed validator result, stale/unmatched handling review, and operations monitoring proof. |
| Trip Updates | Public Trip Updates feed path, prediction adapter boundary, feed-health row, and readiness row. | Real operating-data coverage, deployed-feed validator result, agency review of conservative predictions, and consumer acceptance if claimed. |
| Alerts | Public Alerts feed path, alert lifecycle capability, feed-health row, and readiness row. | Operator alert workflow proof, deployed-feed validator result, live alert lifecycle evidence, and consumer acceptance if claimed. |
| Validation | Validation Health, Validation Center, GTFS Quality, and GTFS-RT conformance checks. | Current deployed validator records and any required external reviewer records. |
| License and contact | Setup, Feeds, and readiness metadata rows. | Agency-approved open license and monitored public technical contact. |
| Uptime and operations | Feed Health, Reliability, Maintenance, and support-bundle guidance. | Deployment monitoring, incident response, backup/restore practice, and any SLA/uptime evidence if claimed. |
| Telemetry and devices | Devices, Telemetry, Telemetry Simulator, Realtime Center, and Vehicle Positions review. | Real device configuration, operational telemetry coverage, and production AVL reliability proof if claimed. |
| Consumer preparedness | Consumers page, prepared tracker, connector catalog, and readiness row. | Provider or regional website listing, target-specific submission/review records, and target-originated acceptance/ingestion/listing/display proof. |

## Claims To Avoid

Do not describe a local run or reference trial as:

- CAL-ITP/Caltrans compliant;
- accepted, listed, ingested, or displayed by a consumer;
- agency approved or agency adopted;
- production ready for all agencies;
- vendor certified or hardware certified;
- production-grade ETA quality.

## How Agencies Can Help

Agencies can help the project by trying the local/reference workflow, testing
with their public GTFS ZIP, contributing connector examples, reviewing
deployment docs, sharing non-sensitive feedback, or sponsoring a later pilot.
Formal agency approval, final feed-root evidence, and consumer acceptance are
not required to use or improve the software; they are future evidence
milestones only for agencies that choose public launch or compliance claims.

Detailed guide: [CAL-ITP Readiness Checklist](../docs/tutorials/calitp-readiness-checklist.md).
