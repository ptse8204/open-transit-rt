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

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no evidence writes

