# Phase 50 Handoff -- Realtime Quality Backtesting

## Status

Phase 50 is closed for the approved private diagnostics scope.

## Implemented

- Added `internal/realtimequality/backtest.go` with versioned local schemas for
  observed stop events and prediction samples.
- Added aggregate backtest metrics for overall, route, agency-local time
  period, and route plus time period groups.
- Added MAE, median absolute error, p90 absolute error, lead-time aggregates,
  prediction coverage, future stop coverage, stale/missing/withheld counts, and
  diagnostic maturity gates.
- Added `cmd/realtime-quality-backtest`.
- Added synthetic public-safe fixtures under
  `testdata/realtime-quality-backtest`.
- Added `make realtime-quality-backtest`.

## Output Boundary

Default output is `.cache/realtime-quality-backtest/<timestamp>`.

The CLI writes exactly:

- `summary.json`
- `summary.md`
- `metrics.json`
- `metrics.md`
- `manifest.json`

The outputs contain bounded aggregates and redacted manifest metadata only. Raw
observed rows, raw prediction rows, raw telemetry, raw GTFS-RT payloads,
private device IDs, driver IDs, vendor payloads, credentials, headers, cookies,
database URLs, private paths, and raw logs are not copied into the output.

Lead-time diagnostics are aggregate-only and are computed only from matched,
non-stale predictions where `predicted_time` is after `generated_time`.

## Safety Boundary

The workflow rejects `docs/evidence`, evidence-like paths, unsafe traversal,
and symlink ancestors. It adds no database persistence, migrations, Operations
Console route, public API, consumer tracker change, external predictor runtime,
publish gate, evidence packet, or stronger readiness/ETA claim.

Maturity gate labels are limited to:

- `insufficient_data`
- `diagnostic_pass`
- `diagnostic_watch`

## Verification

- `go test ./internal/realtimequality ./cmd/realtime-quality-backtest` -- passed.
- `make realtime-quality` -- passed.
- `make realtime-quality-backtest` -- passed.
- `make validate` -- passed.
- `make test` -- passed.
- `make smoke` -- passed.
- `make test-integration` -- passed with local Postgres/PostGIS started by
  `make db-up` and stopped by `make db-down`.
- `git diff --check` -- passed.
- `docker compose -f deploy/docker-compose.yml config` -- passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` -- passed.
- Exact seven-target prepared-only consumer tracker check -- passed.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` -- passed.
- `git diff --exit-code -- docs/evidence` -- passed.
- Synthetic scale coverage is included in `internal/realtimequality` tests and
  uses indexed joins rather than pairwise prediction/observation scans.

## Next Step

Start Phase 51 planning. Do not reopen Phase 50 unless a concrete
backtesting, evidence-boundary, consumer-tracker, public-route, security, or
claim-boundary regression is found.
