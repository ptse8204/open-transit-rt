# Phase 44 Handoff

## Phase

Phase 44 — Telemetry Simulator And Device Trial

## Status

Complete for the synthetic-only local/reference telemetry simulator scope.

## What Was Implemented

- Patched `docs/handoffs/latest.md` before Phase 44 implementation so the
  Current Objective says Phase 43 is complete and Phase 44 is next.
- Added `cmd/telemetry-simulator` for synthetic scenario loading, authenticated
  `POST /v1/telemetry` sends, expectation checks, private diagnostics, and
  optional post-ingest DB-backed matcher/Vehicle Positions debug output.
- Patched the simulator safety closure so non-loopback credentialed sends
  require HTTPS, dry-run can still validate non-loopback HTTP targets, custom
  output directories are confined to repo `.cache` unless explicitly
  overridden, `docs/evidence` is always rejected, symlink output directories
  are rejected, and generated diagnostics receive a final redaction scan.
- Added private timing diagnostics for run duration, per-event HTTP duration,
  matcher duration, matcher total duration, and Vehicle Positions debug
  duration where applicable.
- Added `scripts/telemetry-simulator.sh` and `make telemetry-simulator`.
- Added `testdata/telemetry-simulator/` fixtures for:
  - `on-route`
  - `stale`
  - `out-of-order`
  - `unknown-device`
  - `low-quality-gps`
  - `after-midnight`
  - `block-transition`
- Added `docs/tutorials/telemetry-simulator-and-device-trial.md`.
- Added `docs/phase-44-telemetry-simulator-and-device-trial.md`.
- Updated README/docs navigation, current status, backlog, open questions,
  roadmap, integration adapter kit, and operator tutorials.
- Updated `make validate` to check the simulator script, help/list paths,
  dry-run path, JSON fixtures, and Phase 44 docs.

## Boundaries Preserved

- Synthetic data only.
- Real device bearer-token auth.
- No bypass around `/v1/telemetry`.
- No real vendor payloads.
- No private telemetry.
- No credentialed non-loopback plain HTTP sends.
- No diagnostic output under `docs/evidence`.
- No evidence packets.
- No consumer status changes.
- No vendor compatibility claim.
- No production AVL reliability claim.
- No production-grade ETA claim.
- No CAL-ITP/Caltrans compliance claim.

## Consumer Tracker

`docs/evidence/consumer-submissions/status.json` remains unchanged. All seven
tracked targets remain `prepared`.

## Checks Run

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make agency-app-up` — passed.
- Local simulator check against `http://localhost:8080` with the default
  fixture timestamp and `RUN_MATCHER=true` — passed with `out_of_order`,
  demonstrating repeat-run-safe diagnostics against existing local telemetry.
- Fresh local simulator check against `http://localhost:8080` with
  `REFERENCE_TIME=2026-05-11T15:05:00Z` and `RUN_MATCHER=true` — passed with
  HTTP `201`, ingest status `accepted`, matcher output, and private Vehicle
  Positions debug `trip_descriptor_published=true`.
- Safety patch verification before Phase 45 — passed:
  `make validate`, `make test`, `git diff --check`,
  `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`,
  `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`,
  and `docker compose -f deploy/docker-compose.yml config`.
- Safety patch local simulator check with
  `TARGET=http://localhost:8080 DEVICE_TOKEN=dev-device-token SCENARIO=testdata/telemetry-simulator/on-route.json RUN_MATCHER=true REFERENCE_TIME=2026-05-11T15:05:00Z make telemetry-simulator` — passed with HTTP `202`, ingest status `duplicate`, private timing fields, and no evidence/status changes because the reused local demo database already contained the same accepted synthetic timestamp from earlier Phase 44 verification.
- `make agency-app-down` — passed.

## Next-Step Recommendation

Continue with Phase 45 — GTFS Quality Triage Loop only if the maintainer wants
to proceed with the roadmap. Keep simulator output private unless a future
evidence-specific phase authorizes redaction, retention, and claim mapping.
