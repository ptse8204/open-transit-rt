# Generic JSON Transform Connector

Synthetic-only telemetry transform example. It maps flat JSON records into
Open Transit RT telemetry events, keeps `send_enabled=false`, and prints a
redaction-safe dry-run batch.

Run the first safe check:

```sh
go run ./examples/connectors/generic-json-transform
go test ./examples/connectors/generic-json-transform
```

The fixture carries an explicit field map, a bounded timeout, a send gate, and
synthetic records. Unmapped fields are not copied into output diagnostics.

This example does not prove vendor compatibility, hardware certification,
production AVL reliability, consumer display, compliance, production readiness,
SLA coverage, or ETA quality.
