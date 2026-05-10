# Connector Cookbook

Use this cookbook when you want to bring your own GPS, AVL, CSV replay, or
sidecar process to Open Transit RT.

The preferred integration shape is sidecar plus manifest plus conformance
tests. Connectors do not load as dynamic Go plugins, do not bypass auth, and do
not change core database state directly.

## Bring Your Own GPS Source

1. Keep the original source system outside this repository.
2. Transform each observation into the Open Transit RT telemetry event shape.
3. Send validated events to:

   ```text
   POST /v1/telemetry
   Authorization: Bearer <device-token>
   Content-Type: application/json
   ```

4. Use deployment-owned device credentials.
5. Reject malformed, stale, future-dated, wrong-agency, unknown-device,
   low-quality, duplicate, or out-of-order observations instead of guessing.
6. Keep diagnostics redacted and private by default.

Detailed guide: [Device And AVL Integration](../docs/tutorials/device-avl-integration.md).

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

These are local quality checks. They do not prove vendor compatibility,
production AVL reliability, production-grade ETA quality, consumer acceptance,
production readiness, or CAL-ITP/Caltrans compliance.

## Predictor Sidecars

Trip Updates prediction stays behind `internal/prediction.Adapter`. Vehicle
Positions must continue without an external predictor, and deterministic
prediction remains the default.

Use the generic external HTTP adapter only when explicitly configured. Shadow
mode is appropriate for evaluation because it keeps public Trip Updates output
on the deterministic path while recording bounded diagnostics.

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
