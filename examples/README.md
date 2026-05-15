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

The Phase 101 V2 conformance suite includes local synthetic checks for missing
telemetry fields, invalid coordinates, missing Vehicle Positions references,
public mutation attempts, validator command blocking, and unredacted monitoring
destinations. These are fail-closed local checks, not real integration proof.

| Example | What it demonstrates | What it does not prove |
| --- | --- | --- |
| `connectors/telemetry-http-poller` | A synthetic sidecar that normalizes source observations before telemetry ingest | No real source polling, credentials, network send, or vendor compatibility |
| `connectors/telemetry-csv-replay` | A synthetic CSV replay adapter for local development | No production import workflow or realtime data proof |
| `connectors/predictor-sidecar-stub` | The Trip Updates predictor sidecar boundary | No production-grade ETA quality or named predictor compatibility |
| `connectors/validator-allowlist` | Server-owned validator ID allowlist decision | No validator-clean result, compliance proof, or consumer acceptance |
| `connectors/monitoring-export` | Redacted monitoring/export summary shape | No notification delivery, monitoring service, SLA, or evidence creation |

Telemetry examples share a small SDK-style helper under
`connectors/sdk/telemetry` for dry-run normalization, fail-closed drop reasons,
and no-send event output. It is example code for adapter authors, not a
production vendor SDK or compatibility claim.

The prediction sidecar stub uses `connectors/sdk/prediction` for sanitized
dry-run request/response examples, withheld-output diagnostics, and explicit
Vehicle Positions independence. It is not a named predictor integration and
does not prove ETA quality.

The monitoring/export example uses `connectors/sdk/monitoring` for redaction,
no-send output, no status mutation, and no evidence writes. The validator
example uses allowlisted validator ID/feed-type pairings plus safe fixture
references and rejects raw validator command patterns.

## Fixture Boundary

Examples use synthetic fixtures only. Do not add real device tokens, API keys,
vendor credentials, private identifiers, raw private payloads, unredacted logs,
private hostnames, webhook URLs, database URLs, or private infrastructure
details.

## Related Docs

- [Integration Adapter Kit](../docs/integration-adapter-kit.md)
- [Connector Plugin Contract](../docs/connectors/plugin-contract.md)
- [Redaction-First Connector Recipes](../docs/connectors/redaction-first-recipes.md)
- [External Adapter Conformance](../docs/tutorials/external-adapter-conformance.md)
- [Connector Cookbook](../wiki/connector-cookbook.md)
