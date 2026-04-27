# Phase 19 — Realtime Quality And ETA Improvement

## Status

Implemented for the measurement-first Phase 19 scope. See `docs/handoffs/phase-19.md`.

## Purpose

Phase 19 returns to the hardest original technical problem: trustworthy trip matching and useful Trip Updates. Earlier phases built conservative matching and prediction. This phase improves measurable quality with deterministic replay tests, explicit metrics, and safer diagnostics.

## Scope

1. Replay evaluation framework.
2. Matcher/predictor quality metrics.
3. Better diagnostics for ambiguous or withheld predictions.
4. Improved schedule-deviation prediction where safe.
5. Optional external predictor contract evaluation.
6. Detour/short-turn/cancellation quality review.

## Required Work

### 1) Replay Evaluation

Implemented a deterministic fixture replay framework under `internal/realtimequality` with fixtures under `testdata/replay/`. The replay suite compares:

- telemetry events;
- expected assignments;
- expected Vehicle Positions output;
- expected Trip Updates behavior;
- withheld/unknown/degraded cases.

Fixture schema and denominator definitions are documented in `testdata/replay/README.md`.

### 2) Quality Metrics

Track and report:

- assignment confidence distribution;
- unknown rate;
- ambiguous candidate rate;
- stale telemetry rate;
- Trip Updates coverage;
- future-stop coverage;
- withheld-by-reason counts;
- validator outcomes;
- manual override usage.

Rate metrics carry numerator, denominator, status, denominator definition, and an explicit `not_applicable` status when a denominator is zero.

### 3) ETA/Prediction Improvements

No broad ETA algorithm change was made in this phase. The implemented changes preserve conservative behavior and improve evidence/diagnostics around stale, ambiguous, unknown, low-confidence, canceled, added-trip, short-turn, and detour cases.

### 4) External Predictor Evaluation

No TheTransitClock or other external predictor integration was added in this phase. If evaluated later:

- keep it behind `internal/prediction.Adapter`;
- do not make it the source of truth;
- add contract tests;
- document failure behavior.

## Acceptance Criteria

Phase 19 is complete only when:

- quality can be measured repeatably;
- known hard cases are represented in fixtures or replay docs;
- improvements are backed by metrics;
- conservative unknown behavior remains possible;
- optional predictors remain adapter-bound;
- Trip Updates claims remain realistic.

## Implemented Notes

- `make realtime-quality` runs `go test ./internal/realtimequality`.
- Operations Console feed/dashboard views show only safe Trip Updates quality summaries from recorded `feed_health_snapshot` diagnostics.
- If no Trip Updates diagnostics exist, the console says `no Trip Updates diagnostics recorded yet`.
- The console does not synthesize a green or healthy summary without recorded diagnostics.

## Required Checks

```bash
make validate
make test
make smoke
make test-integration
git diff --check
```

## Explicit Non-Goals

Phase 19 does not:

- claim production-grade ETA quality without evidence;
- force an external predictor into the core model;
- reduce uncertainty reporting to make feeds look better;
- ignore manual override workflows.
