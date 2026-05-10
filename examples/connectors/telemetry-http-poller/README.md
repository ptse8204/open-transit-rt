# Telemetry HTTP Poller Example

Synthetic-only example for a telemetry source sidecar. It models the shape of a
poller that would normalize observations before Open Transit RT ingest, but it
does not send network requests by default.

## Run

```sh
go run ./examples/connectors/telemetry-http-poller
```

The command reads `fixtures/observations.json`, drops unsafe records, and prints
a dry-run transform summary to stdout.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no evidence writes

