# Connector Cookbook

Use this cookbook when you want to bring your own GPS, AVL, CSV replay, or
sidecar process to Open Transit RT.

The preferred integration shape is sidecar plus manifest plus conformance
tests. Connectors do not load as dynamic Go plugins, do not bypass auth, and do
not change core database state directly.

For the current product-quality scorecard, see
[Review And Recommendations](../docs/roadmap-status.md#review-and-recommendations).

## Bring Your Own GPS Source

1. Keep the original source system outside this repository.
2. Transform each observation into the Open Transit RT telemetry event shape.
3. Send validated events to:

   ```text
   POST /v1/telemetry
   Bearer device token required
   JSON telemetry payload required
   ```

4. Use deployment-owned device credentials.
5. Reject malformed, stale, future-dated, wrong-agency, unknown-device,
   low-quality, duplicate, or out-of-order observations instead of guessing.
6. Keep diagnostics redacted and private by default.

Detailed guide: [Device And AVL Integration](../docs/tutorials/device-avl-integration.md).

## Start In The UI

Before writing adapter code, open:

```text
/admin/operations/connectors
/admin/operations/connectors/tests
/admin/operations/connectors/workbench
/admin/operations/telemetry-simulator
```

Use those pages to decide whether you are bringing GPS data, replaying a CSV,
testing a prediction sidecar, checking validators, or exporting monitoring
summaries. The browser pages are guidance only; they do not execute connector
commands or contact external systems.

## Practical Recipes

| Recipe | Start with | Verify with | What it does not prove |
| --- | --- | --- | --- |
| I have GPS or AVL observations | Transform to `POST /v1/telemetry` with deployment-owned device tokens | `/admin/operations/telemetry`, `/admin/operations/feed-health`, `make telemetry-simulator` | Vendor compatibility or real AVL reliability |
| I have a CSV replay | Use `examples/connectors/telemetry-csv-replay` with synthetic or redacted fields | `make test-connector-examples` | Production data quality or hardware certification |
| I need predictions | Keep Trip Updates behind `internal/prediction.Adapter` or an external sidecar boundary | `make adapter-conformance` | Production-grade ETA quality |
| I need validation checks | Use server-owned allowlisted validator IDs | `/admin/operations/validation-health` | CAL-ITP/Caltrans compliance |
| I need monitoring/export | Keep redacted summaries local until a separate sharing decision exists | `examples/connectors/monitoring-export` | SLA, uptime, or notification delivery |
| I need a redaction-first starting point | Use `docs/connectors/redaction-first-recipes.md` and the Workbench decision tree | `make external-connection-check` and `make adapter-conformance` | Real integration proof or production readiness |

The release-candidate path should also review feed-consumer URL and metadata
expectations, redaction checks, validator tooling, and monitoring/export
surfaces. Consumer/discovery workflows must not automate submissions or change
target statuses.

## Example Paths

| Example | Purpose |
| --- | --- |
| `examples/connectors/telemetry-http-poller` | Synthetic HTTP polling shape that normalizes observations without sending by default |
| `examples/connectors/telemetry-csv-replay` | Synthetic CSV replay shape for local adapter development |
| `examples/connectors/predictor-sidecar-stub` | Prediction sidecar boundary stub without production ETA claims |
| `examples/connectors/monitoring-export` | Redacted monitoring/export summary that does not send by default |

Run:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

These are local quality checks. Phase 101 V2 conformance also covers missing
telemetry fields, invalid coordinates, missing Vehicle Positions references,
public mutation attempts, validator command blocking, and unredacted monitoring
destinations. These checks do not prove vendor compatibility, production AVL
reliability, production-grade ETA quality, consumer acceptance, production
readiness, or CAL-ITP/Caltrans compliance.

## Predictor Sidecars

Trip Updates prediction stays behind `internal/prediction.Adapter`. Vehicle
Positions must continue without an external predictor, and deterministic
prediction remains the default.

Use the generic external HTTP adapter only when explicitly configured. Shadow
mode is appropriate for evaluation because it keeps public Trip Updates output
on the deterministic path while recording bounded diagnostics.
Fail-closed mode is appropriate only when valid empty/adapted Trip Updates and
diagnostics are acceptable for the review. External predictor output must be
validated before serialization.

## Monitoring And Export

Monitoring remains deployment-owned. Connector examples may prepare redacted
diagnostics, but they must not send notifications, create evidence, or claim
SLA/uptime coverage by default.

## What Not To Commit

Do not commit real device tokens, API keys, vendor credentials, private
identifiers, raw private payloads, unredacted logs, webhook URLs, private
database URLs, or private infrastructure details.

See [Integration Adapter Kit](../docs/integration-adapter-kit.md) and
[Connector Plugin Contract](../docs/connectors/plugin-contract.md) for the
full contract.
