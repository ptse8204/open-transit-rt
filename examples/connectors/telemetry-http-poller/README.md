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

## Adapter Shape

The example uses the shared SDK-style helper in
`examples/connectors/sdk/telemetry` to normalize observations into the
Open Transit RT telemetry event shape. The helper emits dry-run events only
when identity, timestamp, age, quality, and coordinates pass conservative
checks. Unsafe records are withheld with bounded drop reasons such as
`invalid timestamp`, `future timestamp`, `stale observation`, `low quality`,
or `invalid coordinates`.

This example does not create device credentials or send to `/v1/telemetry`.
When adapting it for a deployment-owned sidecar, keep credential lookup,
endpoint configuration, private payload mapping, and send behavior outside the
public example until an operator reviews them.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no production AVL reliability claim
- no evidence writes
