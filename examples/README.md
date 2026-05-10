# Examples

This directory contains small, synthetic examples for contributors and
integrators. They are intended to show integration shape, not to prove real
vendor compatibility, production AVL reliability, consumer acceptance, or
CAL-ITP/Caltrans compliance.

## Connector Examples

Run all connector example checks:

```bash
make test-connector-examples
```

Run manifest and conformance checks:

```bash
make external-connection-check
make adapter-conformance
```

| Example | What it demonstrates | What it does not prove |
| --- | --- | --- |
| `connectors/telemetry-http-poller` | A synthetic sidecar that normalizes source observations before telemetry ingest | No real source polling, credentials, network send, or vendor compatibility |
| `connectors/telemetry-csv-replay` | A synthetic CSV replay adapter for local development | No production import workflow or realtime data proof |
| `connectors/predictor-sidecar-stub` | The Trip Updates predictor sidecar boundary | No production-grade ETA quality or named predictor compatibility |
| `connectors/monitoring-export` | Redacted monitoring/export summary shape | No notification delivery, monitoring service, SLA, or evidence creation |

## Fixture Boundary

Examples use synthetic fixtures only. Do not add real device tokens, API keys,
vendor credentials, private identifiers, raw private payloads, unredacted logs,
private hostnames, webhook URLs, database URLs, or private infrastructure
details.

## Related Docs

- [Integration Adapter Kit](../docs/integration-adapter-kit.md)
- [Connector Plugin Contract](../docs/connectors/plugin-contract.md)
- [External Adapter Conformance](../docs/tutorials/external-adapter-conformance.md)
- [Connector Cookbook](../wiki/connector-cookbook.md)
