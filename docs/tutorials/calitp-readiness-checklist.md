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

The productized readiness workflow is available in the authenticated
Operations Console:

```text
/admin/operations/readiness
```

It starts with a workflow map and then shows detailed readiness cards. Each
area has a status source, current signal, what the local review helps prepare,
next action, and claim boundary. The page supports CAL-ITP-style readiness
workflows; it does not claim CAL-ITP/Caltrans compliance.

| Area | Current repo support | Evidence still needed for a deployment |
| --- | --- | --- |
| Public feed URLs | `/public/feeds.json`, `/public/gtfs/schedule.zip`, Vehicle Positions, Trip Updates, and Alerts URL review | Public HTTPS fetch proof, source-of-truth listing, and target-specific review |
| Static GTFS URL | `/public/gtfs/schedule.zip` from active published GTFS | Public HTTPS fetch proof, current active feed, validator result |
| Vehicle Positions URL | `/public/gtfsrt/vehicle_positions.pb` | Public HTTPS fetch proof, fresh telemetry, validator result |
| Trip Updates URL | `/public/gtfsrt/trip_updates.pb` through prediction adapter | Public HTTPS fetch proof, validation, coverage and quality review |
| Alerts URL | `/public/gtfsrt/alerts.pb` from persisted published alerts | Public HTTPS fetch proof, validation, alert lifecycle operations |
| Discovery metadata | `/public/feeds.json` | Complete license/contact data and stable canonical URLs |
| Validation workflow | `/admin/operations/validation-center` and `/admin/operations/validation-health` with allowlisted validators | Latest canonical validator results for each feed |
| License/contact workflow | `feed_config` and `published_feed` metadata | Agency-approved open license and technical contact |
| Uptime/operations signals | Feed Health, Reliability, Maintenance, and support-bundle guidance | Deployment monitoring, incident response, backup/restore practice, and any SLA/uptime records if claimed |
| Telemetry/device state | Devices, Telemetry, Telemetry Simulator, Realtime Center, and Vehicle Positions review | Real device setup, operational telemetry coverage, and production AVL reliability proof if claimed |
| Consumer workflow records | `consumer_ingestion` records | Actual submissions, responses, and acceptance evidence |
| Scorecard | `/admin/operations/readiness` and `/admin/compliance/scorecard` | Current deployment-specific scorecard with supporting validation records |

## Operations Console Workflow

Use the readiness page after publication metadata and at least one GTFS import
or Studio publish exist.

For a guided local/reference path that prepares those inputs through
`make agency-pilot-up`, verifies the five public paths, handles validators,
runs the synthetic AVL dry-run, and then reviews this page, see
[Self-Hosted Operator Trial](self-hosted-operator-trial.md).

1. Open `/admin/operations/readiness` through the private admin boundary.
2. Review the workflow map for public feed URLs, static GTFS, Vehicle
   Positions, Trip Updates, Alerts, validation, license/contact metadata,
   uptime/operations signals, telemetry/device state, and consumer
   preparedness.
3. Open the detailed readiness cards for the current signal and next action.
4. Run validators or operational helpers from their existing workflows; the
   readiness page itself is read-only.
5. Keep any deployment output private until reviewed against
   `docs/evidence/redaction-policy.md`.

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

The repo stores prepared-only docs tracker targets for:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

It does not call external consumer submission APIs and it does not prove acceptance. Runtime consumer workflow records are local deployment notes only; they do not override the docs tracker prepared state. Record consumer status only when the agency or operator has real evidence for that deployment.

## Marketplace Gap

Open Transit RT is not currently a full California Mobility Marketplace vendor package. Additional non-code work is still needed for vendor-equivalent positioning:

- hardware/BYOD deployment guidance
- implementation plan templates
- support runbooks
- SLA/KPI reporting
- procurement documentation
- third-party journey-planner integration support evidence

Keep those as deployment and service-packaging work unless a later phase explicitly adds them.
