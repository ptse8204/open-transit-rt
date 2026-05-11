# CAL-ITP Readiness Plain English

Open Transit RT supports CAL-ITP-style readiness workflows. That means the
software helps operators evaluate and prepare the technical pieces commonly
needed for public GTFS and GTFS Realtime publication.

It does not mean this repository, a local demo, or a private trial is
CAL-ITP/Caltrans compliant.

## What The Software Helps With

- stable public feed paths for static GTFS and GTFS Realtime;
- `feeds.json` discovery metadata;
- Vehicle Positions, Trip Updates, and Alerts feed publication paths;
- license and technical-contact metadata workflows;
- validator workflow records;
- readiness and scorecard views;
- prepared consumer/aggregator packet records;
- local `.cache` readiness gap summaries.

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
| Static GTFS | Browser import, active schedule, public `schedule.zip`, GTFS quality guidance, and validator-health rows. | Agency-approved final dataset, public HTTPS fetch proof, current validator record, license/contact approval, and source-of-truth website listing. |
| Vehicle Positions | Public Vehicle Positions feed path, feed-health row, telemetry freshness, and readiness row. | Real device telemetry, deployed-feed validator result, stale/unmatched handling review, and operations monitoring proof. |
| Trip Updates | Public Trip Updates feed path, prediction adapter boundary, feed-health row, and readiness row. | Real operating-data coverage, deployed-feed validator result, agency review of conservative predictions, and consumer acceptance if claimed. |
| Alerts | Public Alerts feed path, alert lifecycle capability, feed-health row, and readiness row. | Operator alert workflow proof, deployed-feed validator result, live alert lifecycle evidence, and consumer acceptance if claimed. |
| Discovery and consumers | `feeds.json`, metadata workflow, prepared consumer/aggregator records, and readiness status. | Provider or regional website listing, target-specific submission/review records, and target-originated acceptance/ingestion/listing/display proof. |

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
