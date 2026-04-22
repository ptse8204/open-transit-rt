# CAL-ITP / Caltrans Readiness Checklist

This checklist uses evidence-bounded language. The repository supports the technical foundations for California transit data readiness, but a specific deployment must provide validation, hosting, metadata, operations, and consumer evidence before stronger compliance claims are appropriate.

The Phase 11 evidence package is [Compliance Evidence Checklist](../compliance-evidence-checklist.md). Use it when deciding whether a claim is repo-proven, deployment-dependent, or dependent on third-party confirmation.

## Current Claim Boundaries

Allowed wording:

- "Open Transit RT supports stable public GTFS and GTFS Realtime feed URLs."
- "Open Transit RT implements technical foundations for CAL-ITP/Caltrans-style readiness."
- "Open Transit RT provides workflows for validation records, license/contact metadata, scorecards, and consumer-ingestion tracking."

Avoid unless backed by deployment evidence:

- "fully CAL-ITP compliant"
- "accepted by Google Maps, Apple Maps, Transit App, or other consumers"
- "production ready for all agencies"
- "complete marketplace vendor equivalent"

## Technical Readiness Areas

| Area | Current repo support | Evidence still needed for a deployment |
| --- | --- | --- |
| Static GTFS URL | `/public/gtfs/schedule.zip` from active published GTFS | Public HTTPS fetch proof, current active feed, validator result |
| Vehicle Positions URL | `/public/gtfsrt/vehicle_positions.pb` | Public HTTPS fetch proof, fresh telemetry, validator result |
| Trip Updates URL | `/public/gtfsrt/trip_updates.pb` through prediction adapter | Public HTTPS fetch proof, validation, coverage and quality review |
| Alerts URL | `/public/gtfsrt/alerts.pb` from persisted published alerts | Public HTTPS fetch proof, validation, alert lifecycle operations |
| Discovery metadata | `/public/feeds.json` | Complete license/contact data and stable canonical URLs |
| Validation workflow | `/admin/validation/run` with allowlisted validators | Latest canonical validator results for each feed |
| License/contact workflow | `feed_config` and `published_feed` metadata | Agency-approved open license and technical contact |
| Consumer workflow records | `consumer_ingestion` records | Actual submissions, responses, and acceptance evidence |
| Scorecard | `/admin/compliance/scorecard` | Current production-mode scorecard with supporting validation records |

## Local Evidence Commands

For local development evidence, run:

```bash
make validators-install
make validators-check
make demo-agency-flow
```

The demo proves the repo flow is runnable locally. It does not prove a public deployment is compliant.

For code checks:

```bash
make validate
make test
make smoke
```

For DB-backed integration coverage:

```bash
make test-integration
```

## Deployment Evidence To Collect

For a real agency deployment, collect:

- Public HTTPS fetch output for `/public/gtfs/schedule.zip`.
- Public HTTPS fetch output for `/public/feeds.json`.
- Public HTTPS fetch output for each GTFS-RT protobuf feed.
- Static GTFS canonical validator result.
- GTFS-RT validator result for Vehicle Positions.
- GTFS-RT validator result for Trip Updates.
- GTFS-RT validator result for Alerts.
- License name and URL approved by the agency.
- Technical contact email monitored by the agency or operator.
- Scorecard JSON from `/admin/compliance/scorecard`.
- Consumer-ingestion records showing each target consumer’s actual status.
- Any external consumer acceptance evidence that will be claimed publicly.

## Consumer Ingestion

The repo stores consumer-ingestion workflow records for:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit

It does not call external consumer submission APIs and it does not prove acceptance. Mobility Database and transit.land may be tracked as workflow records when an operator adds them, but they are not API integrations. Record consumer status only when the agency or operator has real evidence for that deployment.

## Marketplace Gap

Open Transit RT is not currently a full California Mobility Marketplace vendor package. Additional non-code work is still needed for vendor-equivalent positioning:

- hardware/BYOD deployment guidance
- implementation plan templates
- support runbooks
- SLA/KPI reporting
- procurement documentation
- third-party journey-planner integration support evidence

Keep those as deployment and service-packaging work unless a later phase explicitly adds them.
