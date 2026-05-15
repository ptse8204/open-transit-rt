# Realtime Quality Backtest Fixtures

These fixtures are synthetic and public-safe. They exercise the private Phase
50 backtesting workflow without raw telemetry, GTFS-RT payloads, device IDs,
driver IDs, credentials, headers, cookies, database URLs, private paths, or raw
operator logs.

Use:

```sh
go run ./cmd/realtime-quality-backtest \
  --observed testdata/realtime-quality-backtest/observed-events.json \
  --predictions testdata/realtime-quality-backtest/prediction-samples.json
```

The command writes aggregate diagnostics under `.cache/realtime-quality-backtest/`.

The fixture set intentionally includes synthetic rows for:

- after-midnight service using agency-local `start_date` plus `25:xx` GTFS time;
- frequency/headway-style service with unscheduled prediction output;
- stale and missing prediction diagnostics;
- unknown and ambiguous assignment withholding;
- external predictor shadow/fail-closed withholding;
- a prediction-only shadow sample that exercises missing-observation metrics.

These rows are local conformance diagnostics only. They do not prove
real-world ETA accuracy, production-grade ETA quality, vendor compatibility,
consumer display, compliance, public launch, or release readiness.
