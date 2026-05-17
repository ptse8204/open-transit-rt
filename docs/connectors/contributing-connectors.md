# Contributing Connectors

Connector contributions are welcome when they preserve Open Transit RT
boundaries and use synthetic or public-safe examples.

This guide is not vendor certification, hardware certification, production AVL
reliability proof, consumer acceptance, compliance proof, or release proof.

## Good Connector Contributions

- Synthetic telemetry source examples that normalize into `POST /v1/telemetry`.
- CSV replay or polling examples with no real endpoint or credential.
- Prediction sidecar stubs that use the documented adapter boundary.
- Validator allowlist examples that reject raw command injection.
- Monitoring/export examples that redact destinations and default to no-send.
- Consumer/discovery metadata examples that keep public feed URL review
  prepared-only with no submission or status mutation.
- Manifest lint improvements using committed synthetic fixtures.
- Adapter conformance cases for malformed, stale, future, wrong-agency,
  unknown-device, low-quality, duplicate, out-of-order, missing-reference,
  public-mutation, timeout, redaction, no-submit, or status-mutation failures.

## Required Boundaries

- Use synthetic fixtures under `testdata/` or examples under
  `examples/connectors/`.
- Do not commit real credentials, real endpoint URLs, webhook URLs, private
  payloads, private hostnames, private device IDs, vendor secrets, or portal
  records.
- Do not claim compatibility with a named vendor, hardware model, or production
  fleet unless a separate authorized evidence track exists.
- Do not add dynamic backend plugin loading.
- Do not mutate public feeds, consumer status records, or protected evidence
  paths.
- Keep Trip Updates behind the prediction adapter boundary.
- Keep monitoring/export examples no-send by default.

## Expected Files

Connector work usually touches one or more of:

- `examples/connectors/...`
- `examples/connectors/sdk/...`
- `testdata/connectors/...`
- `testdata/adapter-conformance/...`
- `docs/integration-adapter-kit.md`
- `docs/connectors/plugin-contract.md`
- `docs/connectors/catalog.md`
- `docs/connectors/redaction-first-recipes.md`
- `docs/tutorials/external-adapter-conformance.md`
- `wiki/connector-cookbook.md`

Do not touch `docs/evidence/captured/**`,
`docs/evidence/consumer-submissions/status.json`, or consumer submission
artifact directories for connector examples.

## Validation

Run:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
git diff --check
```

If Go code changed, also run:

```bash
make test
```

If public wording changed, also run:

```bash
make audit-product-acceptance
make audit-final-claim-review
```

## PR Description Checklist

In the PR, state:

- connector shape: telemetry, prediction, validator, monitoring/export,
  consumer/discovery, future extension model, or docs-only;
- whether it is synthetic/local only;
- fixtures added or updated;
- checks run;
- no real credentials or private payloads were committed;
- no evidence or consumer status changed;
- what the connector does not prove.
