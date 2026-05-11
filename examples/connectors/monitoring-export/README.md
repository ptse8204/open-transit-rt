# Monitoring Export Example

Synthetic-only monitoring export example. It prepares redacted diagnostics for a
deployment-owned monitoring system, but never sends by default and does not
write evidence files.

## Run

```sh
go run ./examples/connectors/monitoring-export
```

The command reads `fixtures/metrics.json`, redacts sensitive fields, and prints
a dry-run export batch to stdout.

## Adapter Shape

The example uses the SDK-style helper in `examples/connectors/sdk/monitoring`
to build a no-send export batch. The helper copies synthetic metrics, keeps
only public incident fields, and marks `send_enabled`, `network_send`,
`status_mutation`, and `evidence_write` as false.

Deployment-owned monitoring can adapt this shape after operator configuration,
but the public example does not carry webhook URLs, credentials, private
operator contacts, raw payloads, notification sends, or retained evidence.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no notification or network sends by default
- no consumer or discovery automation
- no SLA or uptime claim
- no evidence writes
