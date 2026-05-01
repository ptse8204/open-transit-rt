# Phase 29 — Real-World Realtime Quality Expansion

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Expand beyond Phase 19 baseline replay fixtures toward real-world realtime quality evidence.

Phase 19 made quality measurable. Phase 29 should add richer cases and, where evidence supports it, improve deterministic ETA/matching behavior without hiding uncertainty.

## Scope

1. More replay fixtures.
2. Real-world telemetry patterns if available.
3. Observed arrival/departure comparison if available.
4. ETA accuracy metrics.
5. Route/time-period coverage metrics.
6. Safe deterministic improvements.
7. Optional adapter-bound predictor contract tests, if explicitly approved.

## Required Work

### 1) New Replay Fixtures

Add fixtures for:

- after-midnight service;
- frequency-window trips;
- block continuity;
- long layovers;
- sparse telemetry;
- noisy GPS;
- stale/ambiguous real-world patterns;
- cancellation/alert linkage;
- manual overrides over time.

### 2) ETA Quality Metrics

If observed stop times are available, define:

- ETA error by route/time period;
- coverage by route/trip;
- withheld by reason;
- stale rate;
- unknown/ambiguous rates;
- manual override influence.

### 3) Safe Improvements

Improve only when replay/tests show benefit.

Preserve:

- unknown is better than false certainty;
- withheld reasons remain visible;
- zero-denominator honesty;
- adapter boundary;
- no production-grade claim without evidence.

## Acceptance Criteria

Phase 29 is complete only when:

- richer quality fixtures exist;
- metrics are repeatable;
- any ETA improvement is evidence-backed;
- uncertainty remains visible;
- docs state remaining quality limits clearly.

## Required Checks

```bash
make validate
make realtime-quality
make test
make smoke
make test-integration
git diff --check
```

## Explicit Non-Goals

Phase 29 does not:

- claim production-grade ETA quality from small fixtures;
- force external predictors into core;
- hide unknown/withheld cases;
- change consumer statuses;
- claim compliance or acceptance.

## Likely Files

- `internal/realtimequality/`
- `testdata/replay/`
- `internal/prediction/`
- `internal/state/`
- `internal/feed/tripupdates/`
- `cmd/agency-config/operations.go` only for safe summary display
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-29.md`
