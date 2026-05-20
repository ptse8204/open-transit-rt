# Redaction-First Connector Recipes

These templates are for local/synthetic connector review. They help operators
and technical helpers decide what may appear in private diagnostics before any
deployment-owned adapter is connected.

They do not authorize external contact, notification delivery, retained
evidence collection, real credential use, consumer submission, consumer status
movement, vendor compatibility claims, hardware certification claims,
production readiness claims, compliance claims, SLA/uptime claims, or
production-grade ETA claims.

## Decision Tree

| Source shape | Use this boundary | First safe check | Stop if |
| --- | --- | --- | --- |
| CSV vehicle locations | `examples/connectors/telemetry-csv-replay` | `make test-connector-examples` | Rows include private identifiers, unredacted real fleet locations, credentials, or unclear device ownership. |
| GPS polling API | `examples/connectors/telemetry-http-poller` | `make external-connection-check` | A manifest needs live URLs, bearer values, private endpoint text, or source payload bodies. |
| Flat JSON records | `examples/connectors/generic-json-transform` | `go test ./examples/connectors/generic-json-transform` | Field mapping is ambiguous, `send_enabled=true`, private fields are copied into output, or timeout review is missing. |
| AVL source can POST | Deployment-owned receiver before `/v1/telemetry` | `go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance` | Auth, agency mapping, device binding, timestamp handling, or redaction is not reviewed. |
| Synthetic telemetry only | Simulator and committed conformance fixtures | `make telemetry-simulator` | Real credentials, real payloads, retained evidence, or a public claim would be needed. |
| External prediction sidecar | `internal/prediction.Adapter` shadow/fail-closed boundary | `go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance` | Output is stale, malformed, wrong-agency, low-confidence, lacks Vehicle Positions reference, or attempts public mutation. |
| Monitoring/export summaries | Redacted no-send monitoring/export adapter | `go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance` | Destination, contact address, token, private endpoint, or unredacted incident detail appears in output. |
| Off-host validation | Server-owned validator IDs and operator-run tooling | `go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance` | A raw validator command, private artifact path, or unsupported validator ID is required. |
| Consumer/discovery metadata | Prepared-only public feed URL metadata review | `go run ./cmd/adapter-conformance consumer_discovery --suite testdata/adapter-conformance` | Submission automation, consumer status mutation, portal contact, retained evidence write, or acceptance wording appears. |

## Templates

Telemetry source template:

- Allowed fields: `agency_id`, `device_id`, `vehicle_id`, `timestamp`, `lat`,
  `lon`, quality/accuracy signal.
- Redact fields: source record IDs, operator contact, private endpoint,
  authorization header, source payload body.
- Blocked fields: bearer values, API keys, database URLs, private paths,
  unbounded vendor fields.
- Default: `send_enabled=false` and `network_send=false` until deployment
  ownership is reviewed.
- Fail closed: reject before ingest when identity, timestamp, agency,
  coordinate, quality, duplicate, or ordering checks fail.

Prediction sidecar template:

- Allowed fields: active feed version ID, assignment IDs, Vehicle Positions
  reference, synthetic timestamps, confidence.
- Redact fields: source private payload, private URL, operator notes, predictor
  diagnostic detail.
- Blocked fields: public feed mutation flags, consumer submission fields,
  credentials, unvalidated Trip Updates.
- Default: public mutation disabled and Vehicle Positions independent.
- Fail closed: withhold output when predictor response is stale, malformed,
  wrong-agency, low-confidence, or missing Vehicle Positions reference.

Validator/off-host template:

- Allowed fields: validator ID, feed type, artifact reference, run intent,
  blocked-check reason.
- Redact fields: private artifact path, host names, operator workstation path,
  validator command arguments.
- Blocked fields: raw commands, arbitrary argv, private files, unsupported
  validator IDs, evidence paths.
- Default: validator execution remains operator-run or server-owned.
- Fail closed: block checks that require raw commands, unsupported IDs, private
  paths, or evidence writes.

Monitoring/export template:

- Allowed fields: feed freshness bucket, validator count, stale telemetry
  count, redacted incident category.
- Redact fields: notification destination, person name, email address, private
  endpoint, unredacted incident text.
- Blocked fields: webhook tokens, send-enabled defaults, SLA evidence,
  retained evidence paths.
- Default: `send_by_default=false` and destination redacted in examples.
- Fail closed: block export when destination, token, private endpoint, or
  unredacted incident detail appears.

Consumer/discovery template:

- Allowed fields: public feed base URL, static GTFS URL, Vehicle Positions
  URL, Trip Updates URL, Alerts URL, license URL, and technical contact role.
- Redact fields: operator email address, portal URL, private review notes, and
  target-specific instructions.
- Blocked fields: submission status mutation, consumer acceptance flag, portal
  credentials, target-originated evidence write, and automatic contact flag.
- Default: `submit_enabled=false`, `status_mutation=false`,
  `network_send=false`, and `evidence_write=false`.
- Fail closed: block connector behavior when submission automation, consumer
  status mutation, external portal contact, retained evidence write, or
  acceptance wording appears.

## Local Checks

Run:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

The private Connector Workbench at
`/admin/operations/connectors/workbench` renders the same decision tree and
template categories as browser guidance only. It does not run commands, start
sidecars, contact external systems, write evidence, or change consumer
statuses.
