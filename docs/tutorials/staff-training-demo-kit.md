# Staff Training Demo Kit

This kit helps a trainer run a local, synthetic staff walkthrough for Open
Transit RT. It is training material only. It does not create retained evidence,
contact outside parties, move consumer status, prove adoption, prove
compliance, or prove production readiness.

## Use This Kit When

- a director wants to understand what staff will do in the private Operations
  Console;
- a daily operator needs to rehearse schedule, realtime, alerts, and recovery
  checks;
- a technical helper needs bounded context without secrets or raw private
  output;
- an integrator wants to map an external data idea to synthetic/local adapter
  boundaries first.

## Before The Session

1. Start from the private local app or a controlled reference environment.
2. Open `/admin/operations/help`.
3. Choose one role path and one demo scenario.
4. Confirm every participant understands that the session is local/synthetic.
5. Keep all seven consumer targets at `prepared`.

Do not collect retained evidence, real agency files, real credentials, raw AVL
payloads, portal records, screenshots with secrets, or private diagnostic
logs.

## Demo Scenario Set

The scenario definitions also live in
`testdata/training-demo/scenarios.json` for local review.

| Scenario | Fixture path examples | Primary console path | What to teach |
| --- | --- | --- | --- |
| `baseline_start` | `testdata/gtfs/valid-small`, `testdata/telemetry-simulator/on-route.json` | `/admin/operations` | Start Here, first missing item, configured feed URLs, and claim boundaries. |
| `after_midnight` | `testdata/gtfs/after-midnight`, `testdata/replay/after-midnight-service.json` | `/admin/operations/realtime` | Agency-local service day, late-night service, and conservative matching. |
| `frequency_service` | `testdata/gtfs/frequency-based`, `testdata/replay/frequency-exact-window.json` | `/admin/operations/prediction-lab` | Headways, repeated trip instances, withheld Trip Updates, and unknown-over-false-certainty. |
| `stale_unknown_device` | `testdata/telemetry-simulator/stale.json`, `testdata/telemetry-simulator/unknown-device.json` | `/admin/operations/devices` | Stale telemetry triage, unknown device handling, and token secrecy. |
| `alerts_disruption` | `testdata/replay/cancellation-alert-linkage.json`, `testdata/replay/disruption-diagnostics-baseline.json` | `/admin/alerts/console` | Alert lifecycle, canceled-trip hints, and public-feed usefulness boundaries. |
| `connector_conformance` | `testdata/adapter-conformance/suite.json`, `testdata/connectors/valid/synthetic-telemetry-source-input.json` | `/admin/operations/connectors/workbench` | Adapter recipe choice, no-send posture, and authorization boundary before real integrations. |

## 65-Minute Trainer Script

1. Opening boundary, 5 minutes: explain that this is local/synthetic training
   only and ask each participant to name one thing it does not prove.
2. Role path assignment, 10 minutes: assign evaluator, director, operator,
   technical helper, or integrator paths from `/admin/operations/help`.
3. Scenario walkthrough, 20 minutes: use one scenario and discuss visible
   missing, blocked, stale, or needs-review rows.
4. Recovery drill, 15 minutes: pick one common mistake and rehearse the safe
   first step.
5. Technical-helper handoff, 10 minutes: write the page, blocker, owner,
   intended next action, fixture or docs reference, and authorization need.
6. Closeout, 5 minutes: confirm no consumer status, evidence state, release
   state, or public claim changed.

## Technical-Helper Handoff

Collect only bounded context:

- current private console page;
- visible blocker text;
- staff owner;
- intended next action;
- relevant fixture path or docs path;
- whether separate authorization is needed.

Do not collect secrets, tokens, database URLs, private endpoint URLs, raw
telemetry payloads, raw validator logs, private screenshots, protected evidence
paths, registry credentials, or portal records.

## Related Docs

- `docs/operator-training-guide.md`
- `docs/tutorials/agency-demo-flow.md`
- `docs/tutorials/telemetry-simulator-and-device-trial.md`
- `docs/tutorials/external-adapter-conformance.md`
- `docs/evidence/redaction-policy.md`
- `docs/evidence/evidence-track-router.md`
