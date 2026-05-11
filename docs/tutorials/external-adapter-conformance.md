# External Adapter Conformance

The adapter conformance suite is an offline synthetic check for connector
quality. It validates the manifest contract and representative failure cases
for telemetry sources, prediction sidecars, validator wrappers, and monitoring
exports.

Run the full suite:

```sh
make adapter-conformance
```

Private UI guidance:

```text
/admin/operations/connectors/tests
```

The Operations Console page lists the fixed offline checks as generated
instructions only. It does not execute commands, read commands from connector
manifests, capture output, run validators, start sidecars, write evidence,
contact external parties, or change consumer statuses.

Run one section:

```sh
go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance manifest --suite testdata/adapter-conformance
```

The CLI reads `testdata/adapter-conformance/suite.json` and synthetic fixture
files only. It does not send network traffic, run validators, contact
consumers, automate portals, write evidence, mutate repository state, or start
sidecars.

## Covered Cases

Telemetry cases cover malformed, stale, future-dated, wrong-agency,
unknown-device, low-quality, duplicate, and out-of-order input. Expected
behavior is fail-closed rejection or review, not guessed vehicle state.

Prediction cases cover timeout, malformed output, stale output, wrong-agency
output, and low-confidence output. Vehicle Positions remain independent of
predictor availability, and deterministic prediction remains the default.

Validator cases cover allowlisted validator IDs. The suite does not run
validators and does not accept raw validator commands.

Monitoring cases cover redaction and no-send defaults. The suite does not send
notifications or create SLA/uptime evidence.

Passing conformance is a local quality signal only. It is not CAL-ITP/Caltrans
compliance, consumer acceptance, agency approval, vendor compatibility,
production AVL reliability, production readiness, or production-grade ETA
quality.
