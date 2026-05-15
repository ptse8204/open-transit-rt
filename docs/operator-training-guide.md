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

## Demo Scenario Catalog

Use these scenarios for staff training and tabletop review. They reference
committed local/synthetic fixtures only. They are not real agency operations,
vendor proof, consumer evidence, or public launch proof.

| Scenario | Audience | Fixture references | Console path | Training goal | Boundary |
| --- | --- | --- | --- | --- | --- |
| Baseline local agency walkthrough | No-developer evaluator and trainer | `testdata/gtfs/valid-small`, `testdata/telemetry-simulator/on-route.json` | `/admin/operations` | Explain Start Here, schedule import, feed URLs, realtime status, and claim boundaries. | Does not prove agency adoption, final-root readiness, public launch, or consumer display. |
| After-midnight service drill | Daily operator and technical helper | `testdata/gtfs/after-midnight`, `testdata/telemetry-simulator/after-midnight.json`, `testdata/replay/after-midnight-service.json` | `/admin/operations/realtime` | Teach agency-local service day, late-night trip handling, and conservative matching. | Does not prove real-world overnight reliability or ETA accuracy. |
| Frequency-based service drill | Technical helper and integrator | `testdata/gtfs/frequency-based`, `testdata/telemetry/frequency-based.json`, `testdata/replay/frequency-exact-window.json` | `/admin/operations/prediction-lab` | Discuss headways, repeated trip instances, and safe withheld Trip Updates. | Does not prove production-grade ETA quality or real-world headway accuracy. |
| Stale and unknown-device recovery drill | Daily operator | `testdata/telemetry-simulator/stale.json`, `testdata/telemetry-simulator/unknown-device.json`, `testdata/replay/stale-telemetry-withheld.json` | `/admin/operations/devices` | Trace stale telemetry and unknown devices without exposing token values. | Does not prove device fleet reliability, vendor compatibility, or hardware certification. |
| Disruption and alert lifecycle drill | Daily operator and director | `testdata/replay/cancellation-alert-linkage.json`, `testdata/replay/disruption-diagnostics-baseline.json` | `/admin/alerts/console` | Separate alert drafting, lifecycle review, canceled-trip hints, and public-feed usefulness guidance. | Does not prove public alert display, consumer ingestion, or agency approval. |
| Connector conformance tabletop | Integrator and technical helper | `testdata/adapter-conformance/suite.json`, `testdata/connectors/valid/synthetic-telemetry-source-input.json`, `testdata/connectors/valid/synthetic-prediction-input.json` | `/admin/operations/connectors/workbench` | Map a proposed integration to synthetic adapter boundaries before any real credentials exist. | Does not prove named vendor compatibility, real device proof, or hardware certification. |

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

## Trainer Script

Use this 65-minute script for a small staff walkthrough. Keep the session local
and synthetic. Do not record retained evidence, copy secrets, or contact
outside parties.

| Segment | Minutes | Trainer says | Participant does | Boundary |
| --- | --- | --- | --- | --- |
| Opening boundary | 5 | This is private training over local or synthetic data only. | Names one thing the session must not prove. | Creates no evidence, outside contact, or consumer status change. |
| Role path assignment | 10 | Each participant uses the role path closest to their job. | Opens the first console page for that role and reads the first next action. | Does not create support entitlement, adoption proof, or approval. |
| Demo scenario walkthrough | 20 | Use one scenario from the catalog and follow only listed fixture references. | Explains what blocked or missing rows mean. | Uses synthetic/local fixtures only and does not prove field outcomes. |
| Recovery drill | 15 | Pick one common mistake and rehearse the safe first step. | States what must not be copied into notes or tickets. | Does not authorize credential sharing, private data retention, or external contact. |
| Technical-helper handoff | 10 | Handoff notes include page, blocker, owner, intended next action, and authorization need. | Decides whether the next step needs staff, technical helper, maintainer, or evidence-gate review. | Does not start evidence, release, vendor, or consumer workflow. |
| Closeout | 5 | Leave outside statuses unchanged unless a real authorized workflow exists. | Repeats that all consumer targets remain `prepared`. | Does not prove adoption, consumer acceptance, compliance, public launch, hosted service, SLA, or production readiness. |

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

## Technical-Helper Checklist

When a nontechnical operator asks for help, ask for bounded context only.

| Area | Collect | Do not collect | Escalate before |
| --- | --- | --- | --- |
| Startup context | App URL, command attempted, visible status, and local-only environment label. | Secrets, tokens, database URLs, private logs, or screenshots with credentials. | Real deployment, public-root proof, or hosted-service claims. |
| Schedule context | GTFS source type, import status, active feed version, validator row, and fixture path when demo-only. | Real agency private schedule files unless explicitly authorized and redacted. | Source-of-truth listing, agency approval, or compliance proof. |
| Realtime context | Vehicle ID, stale/unmatched state, simulator scenario, device binding status, and withheld reason. | Device tokens, raw private AVL payloads, or unredacted vendor logs. | Real device proof, vendor proof, or ETA-quality claims. |
| Connector context | Chosen recipe, synthetic fixture, conformance case, expected input shape, and no-send posture. | Real credentials, portal tokens, private endpoint URLs, or vendor payload samples. | Live vendor/device integration or network sends. |
| Monitoring context | Health digest row, monitoring export row, no-send channel state, and local summary path if already produced. | Webhook URLs, email credentials, destination values, or private alert recipient lists. | Live notification delivery, hosted monitoring, SLA, or uptime claims. |
| Evidence and claims | Which claim is requested, current prepared-only status, required gate, and missing retained proof. | Protected evidence paths, private portal records, or target-originated artifacts unless an evidence gate authorizes them. | Evidence retention, consumer submission, or status movement. |
| Release or package context | Current local diagnostic status, release note draft path, blocker, and exact command that was run. | Signing keys, registry tokens, private release credentials, or publication secrets. | Tagging, package publication, image push, GitHub Release creation, or release-readiness claims. |

## Claim Boundary

This training guide may describe private browser-first operation, local
diagnostics, staff roles, first-week review, glossary terms, support guidance,
and future authorization gates.

It must not be used to claim CAL-ITP/Caltrans compliance, agency adoption or
approval, consumer submission/review/acceptance/ingestion/listing/display,
final-root readiness, hosted service, paid support, SLA/uptime, production
readiness, vendor compatibility, hardware certification, production-grade ETA
quality, real-world ETA accuracy, or public launch completion.
