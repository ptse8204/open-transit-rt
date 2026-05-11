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

## Adapter Shape

The example uses the SDK-style helper in
`examples/connectors/sdk/prediction` to model a sanitized prediction sidecar
request and a dry-run response. The response keeps the active feed version and
Vehicle Positions reference visible, but it emits no public Trip Updates and
sets `public_trip_updates_mutation` to `false`.

Assignments below the confidence threshold remain `unknown`. Assignments with
no matching telemetry or missing trip descriptor fields are withheld with
bounded diagnostics. Assignments that pass those checks are marked eligible,
but the stub still withholds output because it is an adapter-boundary example,
not a predictor.

Vehicle Positions remain independent of this example. Adapting a real
deployment-owned predictor still requires explicit operator configuration,
sanitized request/response tests, GTFS-Realtime validation, and separate
retained evidence before any stronger ETA-quality or compatibility wording.

## Boundaries

- synthetic fixture data only
- no real credentials
- no real vendor payloads
- no named vendor compatibility claim
- no network sends by default
- no production ETA-quality claim
- no public Trip Updates mutation
- no evidence writes
