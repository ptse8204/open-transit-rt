# Test Fixtures

This directory contains deterministic fixtures for importer, telemetry,
matching, realtime feed, connector, and evaluator checks. Fixtures should stay
synthetic or public-safe. Do not add private agency data, real device tokens,
raw vendor payloads, private vehicle identifiers, or retained evidence here.

## Layout

- `gtfs/valid-small/`: minimal valid static GTFS for import smoke tests.
- `gtfs/after-midnight/`: GTFS with service times beyond `24:00:00`.
- `gtfs/frequency-based/`: GTFS with `frequencies.txt` examples.
- `gtfs/malformed/`: intentionally invalid GTFS for validation and rollback tests.
- `telemetry/`: sample telemetry traces for matching and stale/unmatched behavior.
- `telemetry-simulator/`: synthetic scenarios for authenticated telemetry ingest
  trials, including stale, low-quality, after-midnight, and block-transition
  cases.
- `avl-vendor/`: synthetic mapping and GPS-like payloads for adapter dry-runs.
- `connectors/`: connector manifest fixtures for external connection checks.
- `adapter-conformance/`: offline synthetic conformance cases for telemetry,
  prediction, validator, and monitoring/export boundaries.
- `replay/`: replay fixtures for deterministic realtime behavior tests.
- `realtime-quality-backtest/`: observed/predicted sample data for private
  local realtime quality diagnostics.
- `multi-agency/`: fixtures for route and tenant-boundary checks.
- `expected/`: expected outputs and protobuf decode assertions.

## Fixture Rules

- Keep fixtures deterministic and small enough for local tests.
- Prefer synthetic IDs such as `demo-agency`, `vehicle-1`, and `device-1`.
- Keep real credentials, private endpoints, private logs, and raw operator
  data out of committed files.
- Document fixture purpose in a local `README.md` when the directory is not
  self-explanatory.
- Add or update tests when a fixture encodes a behavior that must not regress.
