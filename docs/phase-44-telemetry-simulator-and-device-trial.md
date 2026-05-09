# Phase 44 — Telemetry Simulator And Device Trial

## Status

Complete for the synthetic-only telemetry simulator and local/reference device
trial scope.

## Goal

Add a safe, repeatable simulator that sends synthetic telemetry through the real
authenticated `POST /v1/telemetry` ingest path so operators can test ingest,
matching, and Vehicle Positions behavior without hardware, vendor payloads, or
private telemetry.

## Implemented

- Added `cmd/telemetry-simulator`, a Go command that loads synthetic scenario
  fixtures, expands deterministic timestamps, sends authenticated HTTP requests
  to `/v1/telemetry`, validates expected HTTP/ingest statuses, and writes
  private diagnostics.
- Added `scripts/telemetry-simulator.sh` and `make telemetry-simulator`.
- Added synthetic scenario fixtures under `testdata/telemetry-simulator/` for
  on-route, stale, out-of-order, unknown-device, low-quality GPS,
  after-midnight, and block-transition cases.
- Added optional `RUN_MATCHER=true` support that reads accepted telemetry back
  from Postgres after HTTP ingest, runs the existing state engine, and writes a
  private Vehicle Positions debug snapshot using the existing feed builder.
- Hardened URL safety so non-dry-run credentialed sends require HTTPS for
  non-loopback targets while preserving loopback HTTP for local app checks and
  non-loopback HTTP validation in dry-run mode.
- Hardened diagnostic output safety so outputs default under `.cache`, custom
  output directories outside repo `.cache` require
  `ALLOW_UNIGNORED_OUTPUT_DIR=true`, `docs/evidence` is always rejected,
  symlink output directories are rejected, and new output directories are
  created with mode `0700`.
- Added final redaction scanning for Authorization headers, Bearer tokens,
  `DEVICE_TOKEN`, `DATABASE_URL`, postgres password URLs, Cookie headers, and
  private-key material.
- Added private timing diagnostics for total run duration, per-event HTTP
  duration, matcher duration, matcher total duration, and Vehicle Positions
  debug duration where applicable.
- Added focused command tests for dry-run redaction, authenticated
  `/v1/telemetry` posting, URL safety, output-directory safety, redaction
  scanning, timing diagnostics, and multi-event count consistency.
- Extended `make validate` to check the simulator command, script, scenarios,
  and docs.

## Operator Entry Points

```bash
make agency-app-up
make telemetry-simulator
RUN_MATCHER=true make telemetry-simulator
make agency-app-down
```

For reference deployments:

```bash
TARGET=https://reference.example.org \
DEVICE_TOKEN=replace-with-private-device-token \
SCENARIO=on-route \
make telemetry-simulator
```

## Boundaries Preserved

- Synthetic data only.
- Real device bearer-token auth.
- No ingest bypass; all sent events use `POST /v1/telemetry`.
- No real vendor payloads.
- No private telemetry.
- No evidence packets.
- No consumer status changes.
- No vendor compatibility claim.
- No production AVL reliability claim.
- No production-grade ETA claim.
- No CAL-ITP/Caltrans compliance claim.

## Notes

The default local scenario uses deterministic fixture time so matching can be
evaluated against `testdata/gtfs/valid-small`. Public feed services still use
their own runtime clock; for deterministic trial snapshots, use
`RUN_MATCHER=true` and review the private `vehicle_positions_debug.json`
diagnostic.
