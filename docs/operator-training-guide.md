# Open Transit RT Operator Training Guide

This guide is printable staff training for the private browser-first
Operations Console. It is not evidence, release proof, public launch proof,
consumer status proof, or compliance proof.

Use it to onboard directors, daily operators, technical helpers, integrators,
and no-developer evaluators. Keep credentials, private records, and raw
diagnostics out of notes.

## Roles

| Role | First screen | Main responsibility | Ask for help when |
| --- | --- | --- | --- |
| No-developer evaluator | `/admin/operations` | Decide whether the browser workflow is understandable. | Startup, validator tooling, public-root configuration, or support bundle preparation is needed. |
| Director or manager | `/admin/operations/launchpad` | Assign owners and understand remaining claim boundaries. | A blocker needs staff, technical, maintainer, or authorization review. |
| Daily operator | `/admin/operations/realtime` | Review schedule, realtime, alerts, feed health, validation, and routine maintenance state. | Telemetry stops, validation is blocked, or feed URLs are not configured. |
| Technical helper | `/admin/operations/setup` | Handle setup, validators, device credentials, support review, and deployment diagnostics. | Evidence retention, release packaging, portal work, or schema-changing work is proposed. |
| Integrator | `/admin/operations/connectors/workbench` | Evaluate connectors using synthetic or local inputs first. | Real credentials, real vendor payloads, or network sends are proposed. |

## First Week

| When | Task | Done when | Boundary |
| --- | --- | --- | --- |
| Day 1 | Open Start Here and identify the first missing or blocked setup item. | The next private task is clear. | This does not prove outside readiness. |
| Day 1 | Review Agency Setup and assign owners for profile, feed metadata, schedule, telemetry, validation, and maintenance. | Each setup area has a staff or technical-helper owner. | This does not prove agency approval. |
| Day 2 | Import or review GTFS in the private schedule workflow. | Schedule state and next actions are visible. | This does not prove validator-clean status or source-of-truth listing. |
| Day 3 | Check feed URLs, Feed Health, and Validation Center. | Missing, stale, blocked, or needs-review rows have owners. | This does not prove final-root readiness or consumer action. |
| Day 4 | Review telemetry, devices, Vehicle Positions, Trip Updates, and Alerts. | Unknown, stale, unmatched, and degraded states are understood. | This does not prove ETA quality, field reliability, or target display. |
| Day 5 | Review connectors, simulator paths, and prediction boundaries. | The next connector experiment uses local or synthetic inputs only. | This does not prove vendor compatibility or hardware certification. |
| Day 5 | Review Maintenance, support-bundle guidance, and redaction warnings. | Maintenance owners and cadence are clear. | This does not prove SLA, uptime, managed support, release readiness, or hosted-service availability. |

## Quick Tasks

### Import GTFS

1. Open `/admin/operations/gtfs-workbench`.
2. Review active feed state, import history, required files, quality triage,
   and validation context.
3. Use `/admin/operations/gtfs-import` only for the existing private admin
   import flow.

Stop and ask for help if source schedule data must be corrected or validator
tooling is unavailable.

### Check Feed Health

1. Open `/admin/operations/feeds`.
2. Review exactly five configured feed URLs: feeds.json, schedule, Vehicle
   Positions, Trip Updates, and Alerts.
3. Review metadata, local fetch context, validation context, and off-host
   guidance.

Do not treat a configured URL as final-root proof or consumer action.

### Review Validation

1. Open `/admin/operations/validation-center`.
2. Identify the feed, validator, artifact, latest status, and blocker.
3. Use existing allowlisted validation workflows only.

Validator results are supporting signals, not compliance proof.

### Add Or Simulate Telemetry

1. Open `/admin/operations/devices` and `/admin/operations/telemetry-simulator`.
2. Review device binding, token boundary, simulator scenario, and telemetry
   freshness.
3. Confirm results in `/admin/operations/realtime`.

Do not introduce real credentials or payloads without authorization.

### Review Connectors

1. Open `/admin/operations/connectors/workbench`.
2. Choose the closest recipe.
3. Review synthetic normalization preview, dry-run guidance, conformance
   results, and no-send boundaries.

Synthetic conformance does not prove a real vendor or hardware outcome.

### Prepare For A Technical Helper

Record only:

- current private console page;
- visible blocker text;
- intended next action;
- staff owner;
- whether external authorization is needed.

Do not copy secrets, tokens, database URLs, raw diagnostic output, or private
operator records into issue trackers or public docs.

### Prepare Support Review

1. Open `/admin/operations/maintenance`.
2. Review support-bundle guidance, redaction warnings, backup/restore review,
   upgrade/rollback review, and maintenance cadence.
3. Keep support material local unless a safe sharing path is explicitly
   approved.

Support review does not prove paid support, managed service availability, SLA,
uptime, release readiness, or public launch.

## What Should I Do?

| Situation | Safe first step | Escalate when | Boundary |
| --- | --- | --- | --- |
| Start Here is empty or blocked | Open Agency Setup and follow the first missing private next action. | Startup, environment, validator tooling, or deployment settings are missing. | Clearing a row does not prove outside readiness. |
| GTFS import has errors | Open GTFS Workbench and Schedule Quality. | Source schedule data needs correction. | Import review does not prove validator-clean status. |
| Feed URL is missing | Open Feed Links & Health. | Public root configuration or off-host validation is needed. | A URL does not prove final-root readiness. |
| Validation is blocked | Open Validation Center. | Validator tooling must be installed or run off-host. | Validator success is still only a supporting signal. |
| Vehicles are stale or unmatched | Open Realtime Center, Telemetry, and Devices. | Devices, tokens, or adapters need operational changes. | Fixing staleness does not prove field reliability. |
| Trip Updates are withheld | Open Prediction & ETA Lab. | External predictors or ETA claims are proposed. | Withheld diagnostics do not prove ETA quality. |
| Connector setup is unclear | Open Connector Workbench and Connector Tests. | Real credentials, vendor payloads, or network sends are proposed. | Synthetic checks do not prove vendor compatibility. |
| Consumer packet status is confusing | Open Consumers and Feed Links & Health. | Separate written authorization starts an evidence or consumer track. | Prepared packet visibility does not prove submission, review, listing, display, ingestion, or acceptance. |

## Staff Handoff

Before a staff handoff, confirm:

- the current private console page;
- the role taking ownership;
- the blocker or next action;
- the related docs link;
- whether technical help is needed;
- whether separate written authorization would be required.

Do not mark consumer targets beyond `prepared`, collect retained evidence,
contact external parties, or publish release artifacts from this training
workflow.

## Claim Boundary

This training guide may describe private browser-first operation, local
diagnostics, staff roles, first-week review, glossary terms, support guidance,
and future authorization gates.

It must not be used to claim CAL-ITP/Caltrans compliance, agency adoption or
approval, consumer submission/review/acceptance/ingestion/listing/display,
final-root readiness, hosted service, paid support, SLA/uptime, production
readiness, vendor compatibility, hardware certification, production-grade ETA
quality, real-world ETA accuracy, or public launch completion.
