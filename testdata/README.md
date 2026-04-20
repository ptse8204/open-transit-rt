# Test Fixtures

This directory contains deterministic fixtures for phased implementation.

Phase 0 creates the fixture structure and small seed files only. Later phases should add parser assertions, repository tests, protobuf snapshots, and integration tests against these fixtures.

## Layout

- `gtfs/valid-small/`: minimal valid static GTFS for import smoke tests.
- `gtfs/after-midnight/`: GTFS with service times beyond `24:00:00`.
- `gtfs/frequency-based/`: GTFS with `frequencies.txt` examples.
- `gtfs/malformed/`: intentionally invalid GTFS for validation and rollback tests.
- `telemetry/`: sample telemetry traces for matching and stale/unmatched behavior.
- `expected/`: expected outputs and protobuf decode assertions added by later phases.
