# Telemetry CSV Replay Example

Synthetic-only command adapter that replays a local CSV fixture into dry-run
telemetry events. It is intended for adapter development and conformance checks,
not for production import.

## Run

```sh
go run ./examples/connectors/telemetry-csv-replay
```

The command reads `fixtures/replay.csv` and prints normalized JSON events to
stdout. It does not call Open Transit RT or any external service.

## Adapter Shape

The CSV fixture includes `synthetic_only`, agency/device/vehicle identifiers,
`observed_at`, coordinates, and a quality score. The example uses the shared
SDK-style helper in `examples/connectors/sdk/telemetry` so CSV rows follow the
same fail-closed checks as the HTTP poller example. Rows with invalid
timestamps, future or stale timestamps, missing identity, low quality, or
invalid coordinates become bounded dry-run drops instead of guessed vehicle
state.

The output is a dry-run summary with `events` and `drops`. It is useful for
adapter development and conformance checks, not for production import or real
realtime operations.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no production AVL reliability claim
- no evidence writes
