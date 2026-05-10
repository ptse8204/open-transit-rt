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

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no notification or network sends by default
- no consumer or discovery automation
- no evidence writes

