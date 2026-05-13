# Phase 84 Prompt — Prediction And ETA Lab

## Goal

Improve prediction transparency and experimentation while preserving production-grade ETA claim boundaries.

## Scope

- Deterministic predictor diagnostics.
- Why Trip Updates were emitted or withheld.
- External predictor shadow-mode comparison summary.
- Fail-closed behavior display.
- Backtesting result browser for local/private aggregate diagnostics.
- Fixture guidance and quality caveats.
- Plain-language ETA quality caveats.

## Principles

- Deterministic remains default.
- External predictors are optional.
- Vehicle Positions remain independent.
- Withheld output must be explained.
- Backtesting is local/private diagnostics, not production proof.

## Do not

- Claim production-grade ETA quality.
- Add named predictor compatibility claims.
- Add live external dependency.
- Persist raw observed-arrival records unless separately planned.
- Write evidence.

## Validation

Baseline validation plus:

```bash
make validate
make test
make realtime-quality-backtest
make audit-final-claim-review
```

## Commits

```text
Phase 84 -- Checkpoint 000001: add prediction and ETA lab plan
Phase 84 -- Checkpoint 000002: add deterministic predictor diagnostics view
Phase 84 -- Checkpoint 000003: add external predictor shadow review UI
Phase 84 -- Checkpoint 000004: add backtesting result browser
Phase 84 -- Checkpoint 000005: add ETA quality caveats and withheld explanations
Phase 84 -- Checkpoint 000006: close prediction and ETA lab review
```
