# Phase 19 — Realtime Quality And ETA Improvement

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 19 returns to the hardest original technical problem: trustworthy trip matching and useful Trip Updates. Earlier phases built conservative matching and prediction. This phase improves quality with replay tests, metrics, and optional predictor contracts.

## Scope

1. Replay evaluation framework.
2. Matcher/predictor quality metrics.
3. Better diagnostics for ambiguous or withheld predictions.
4. Improved schedule-deviation prediction where safe.
5. Optional external predictor contract evaluation.
6. Detour/short-turn/cancellation quality review.

## Required Work

### 1) Replay Evaluation

Add a fixture or replay framework that can compare:

- telemetry events;
- expected assignments;
- expected Vehicle Positions output;
- expected Trip Updates behavior;
- withheld/unknown cases.

### 2) Quality Metrics

Track and report:

- assignment confidence distribution;
- unknown rate;
- ambiguous candidate rate;
- stale telemetry rate;
- Trip Updates coverage;
- validator outcomes;
- manual override usage.

### 3) ETA/Prediction Improvements

Improve only where evidence supports it. Keep conservative behavior where data is weak.

### 4) External Predictor Evaluation

If evaluating TheTransitClock or another predictor:

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
