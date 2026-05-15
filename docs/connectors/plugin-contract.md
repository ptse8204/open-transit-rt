# Connector Plugin Contract

Open Transit RT connectors are optional sidecars or command adapters described
by a manifest and checked by conformance tests. A connector is not loaded as a
Go dynamic plugin, cannot change core state directly, and cannot bypass the
existing public/admin/telemetry contracts.

Long-term compatibility, deprecation, and review rules live in
[`docs/extension-governance.md`](../extension-governance.md).

## Contract Model

Every connector manifest uses:

```json
{
  "schema_version": "open-transit-rt.connector.v1",
  "connector_id": "example.telemetry-poller",
  "type": "telemetry_source",
  "display_name": "Example telemetry poller",
  "description": "Synthetic sidecar that transforms observations into /v1/telemetry.",
  "mode": {
    "name": "sidecar_http",
    "disabled_by_default": true
  },
  "input_contracts": [
    {
      "name": "Synthetic source payload",
      "description": "Synthetic payload reviewed before transform.",
      "schema": "open-transit-rt.telemetry.synthetic.v1"
    }
  ],
  "output_contracts": [
    {
      "name": "Telemetry event",
      "description": "Payload for the authenticated ingest boundary.",
      "schema": "open-transit-rt.telemetry.event.v1"
    }
  ],
  "failure_behavior": {
    "timeout_seconds": 10,
    "retry_policy": "bounded retry for retryable failures",
    "degraded_state": "no telemetry emitted",
    "fail_closed": true
  },
  "redaction_policy": {
    "secret_storage": "env_reference_only",
    "redact_fields": ["authorization", "raw_payload", "private_endpoint"]
  },
  "claim_boundary": {
    "positive_claims": ["adapter contract only", "synthetic conformance tested", "disabled by default"],
    "not_claimed": ["vendor compatibility", "consumer acceptance", "compliance", "production readiness"]
  },
  "docs_link": "docs/connectors/plugin-contract.md",
  "conformance_cases": []
}
```

Required fields are connector ID/type, display name, description, mode,
input/output contracts, failure behavior, redaction policy, claim boundary,
docs link, and conformance cases.

The private Connector Workbench at
`/admin/operations/connectors/workbench` reviews committed example manifests
and local/synthetic recipe guidance in the browser. It is not a dynamic plugin
loader, connector installer, sidecar runner, vendor portal client, evidence
collector, or consumer submission tool.

## Connector Types

`telemetry_source`: Transforms external observations into the authenticated
`POST /v1/telemetry` payload. It must fail closed and must not emit fake
vehicle state when upstream data is malformed, stale, future-dated, wrong
agency, unknown device, low quality, duplicate, or out of order.

`prediction`: Runs behind the existing prediction adapter boundary. Vehicle
Positions must remain independent of predictor availability, and deterministic
prediction remains the default.

`validator`: Wraps validator execution only through allowlisted validator IDs
or offline fixture checks. Raw validator commands, arbitrary argv, and
deployment-private paths are not accepted in manifests.

`monitoring_export`: Exports redacted summaries to deployment-owned monitoring.
It must not send by default and must not create evidence or SLA/uptime claims.

`consumer_discovery`: Prepares or exposes public feed metadata for discovery
workflows. It must not submit to consumers, automate portals, or mutate
consumer tracker statuses.

## Safety Rules

Manifest validation rejects:

- raw secrets, bearer tokens, database URLs, private keys, webhook URLs, or
  credential values;
- private local paths, `file://` URLs, or unsafe non-HTTPS URLs outside
  loopback/example documentation hosts;
- unsafe URL/private-endpoint strings across displayable manifest fields,
  including descriptions, mode names, failure behavior, redaction policy,
  claim-boundary text, and conformance descriptions;
- unsupported positive claims about CAL-ITP/Caltrans compliance, agency
  approval/adoption, consumer acceptance, hosted SaaS, paid support/SLA,
  production readiness, vendor compatibility, hardware certification, public
  launch, production AVL reliability, or production-grade ETA quality;
- raw validator command strings such as `java -jar`, `docker run`, or
  direct validator executable paths;
- notification send-by-default flags;
- consumer submission automation or consumer status mutation flags.

The private Connector Workbench also summarizes these lint checks as operator
guidance: secrets/private endpoints, command/plugin boundaries, no status
mutation or send-by-default, positive-claim allowlisting, and synthetic fixture
scope. That summary is read-only guidance over the same local validation
contract; it does not execute manifests or prove external integration.

## Runtime Boundary

Connectors communicate with Open Transit RT through documented interfaces:

- `/v1/telemetry` for telemetry after deployment-owned credential setup;
- `internal/prediction.Adapter` implementations or sidecars selected by
  existing runtime config for Trip Updates;
- server-side validator IDs for validation;
- redacted private diagnostics for monitoring/export;
- public feed URLs and `/public/feeds.json` for discovery.

Connectors do not add DB migrations, public feed contract changes, telemetry
ingest contract changes, Trip Updates hard-coupling, or external predictor
default changes.

## Conformance

Each connector manifest lists synthetic conformance cases. Phase 101 V2
coverage includes additional telemetry missing-field and invalid-coordinate
cases, prediction missing Vehicle Positions reference and public-mutation
attempt cases, validator command-blocking behavior, and monitoring
unredacted-destination blocking. Passing conformance is a local quality signal
only. It is not compliance proof, consumer acceptance, vendor compatibility,
production AVL reliability, or production ETA-quality evidence.

## Compatibility And Deprecation

The current manifest schema is `open-transit-rt.connector.v1`. Maintainers
should keep schema changes additive when possible, update examples and
conformance cases for new connector types or fields, and use the extension
governance deprecation process before removing or renaming operator-facing
contracts.
