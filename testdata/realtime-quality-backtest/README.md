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
