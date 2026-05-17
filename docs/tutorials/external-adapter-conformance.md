# External Adapter Conformance

The adapter conformance suite is an offline synthetic check for connector
quality. It validates the manifest contract and representative failure cases
for telemetry sources, prediction sidecars, validator wrappers, monitoring
exports, and consumer/discovery metadata.

Run the full suite:

```sh
make adapter-conformance
```

Private UI guidance:

```text
/admin/operations/connectors/workbench
/admin/operations/connectors/tests
```

The Connector Workbench helps operators choose a local/synthetic recipe,
review committed example manifests, read fixed dry-run guidance, inspect
synthetic telemetry normalization preview rows, and understand offline
conformance coverage. The Connector Tests page lists the fixed offline checks
as generated instructions only. Neither page executes commands, reads commands
from connector manifests, captures output, runs validators, starts sidecars,
writes evidence, contacts external parties, or changes consumer statuses.

Run one section:

```sh
go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance consumer_discovery --suite testdata/adapter-conformance
go run ./cmd/adapter-conformance manifest --suite testdata/adapter-conformance
```

The CLI reads `testdata/adapter-conformance/suite.json` and synthetic fixture
files only. It does not send network traffic, run validators, contact
consumers, automate portals, write evidence, mutate repository state, or start
sidecars.

## Covered Cases

Telemetry cases cover malformed, stale, future-dated, wrong-agency,
unknown-device, low-quality, duplicate, out-of-order, missing required field,
and invalid coordinate input. Expected behavior is fail-closed rejection or
review, not guessed vehicle state.

Prediction cases cover timeout, malformed output, stale output, wrong-agency
output, low-confidence output, missing Vehicle Positions reference, and public
mutation attempts. Vehicle Positions remain independent of predictor
availability, and deterministic prediction remains the default.

Validator cases cover allowlisted validator IDs and command-blocking behavior.
The suite does not run validators and does not accept raw validator commands.

Monitoring cases cover redaction, no-send defaults, and unredacted-destination
blocking. The suite does not send notifications or create SLA/uptime evidence.

Consumer/discovery cases cover public feed URL metadata, status-mutation
blocking, and submission-automation blocking. The suite does not submit to
consumers, contact portals, create retained evidence, or move prepared-only
consumer tracker status.

For operator-facing redaction templates and source-shape decisions, see
[Redaction-First Connector Recipes](../connectors/redaction-first-recipes.md).

Passing conformance is a local quality signal only. It is not CAL-ITP/Caltrans
compliance, consumer acceptance, agency approval, vendor compatibility,
production AVL reliability, production readiness, or production-grade ETA
quality.
