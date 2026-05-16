# Synthetic Telemetry Webhook Sidecar

This starter kit shows how a deployment-owned sidecar can normalize a local
webhook-style batch into Open Transit RT telemetry events.

It is disabled by default, synthetic-only, and dry-run only. It does not listen
on a network port, send telemetry, contact a vendor, certify hardware, prove
vendor compatibility, prove production AVL reliability, or create evidence.

Run:

```bash
go run ./examples/connectors/telemetry-webhook-sidecar
```

The output is a redaction-safe dry-run summary using the shared telemetry SDK
shape.
