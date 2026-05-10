# Predictor Sidecar Stub Example

Synthetic-only prediction connector stub. It demonstrates the adapter boundary
without generating production ETAs or replacing the internal prediction path.

## Run

```sh
go run ./examples/connectors/predictor-sidecar-stub
```

The command reads `fixtures/prediction-input.json` and prints a deterministic
no-op response with diagnostics. It does not fetch Vehicle Positions, call a
predictor, or write output files.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no production ETA-quality claim
- no evidence writes

