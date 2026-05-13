# Phase 81 Prompt — Realtime Operations Center

## Goal

Make realtime operations understandable from the browser.

## Scope

- Fleet/vehicle overview.
- Telemetry freshness and device status.
- Vehicle-to-trip assignment explanation.
- Matching confidence and unknown/degraded reasons.
- Vehicle Positions feed status.
- Trip Updates feed status.
- Alerts operational status.
- Issue queue: stale telemetry, unknown device, ambiguous match, withheld Trip Updates, missing alerts.
- Simulator guide/control surface that still uses safe synthetic inputs.

## Feature principles

- Unknown is better than false certainty.
- Vehicle Positions remain independent of external predictors.
- Trip Updates can be withheld when confidence is low.
- Show what was withheld and why.
- Never imply production-grade ETA quality.

## Do not

- Change `/v1/telemetry` semantics.
- Add public fleet map.
- Leak device tokens.
- Claim production AVL reliability.
- Claim vendor compatibility.

## Validation

Baseline validation plus:

```bash
make validate
make test
make telemetry-simulator
```

If simulator requires local app or DB and is blocked, record exact blocker.

## Commits

```text
Phase 81 -- Checkpoint 000001: add realtime operations center plan
Phase 81 -- Checkpoint 000002: add fleet and telemetry freshness overview
Phase 81 -- Checkpoint 000003: add assignment and matching explanation views
Phase 81 -- Checkpoint 000004: add Trip Updates and Alerts status views
Phase 81 -- Checkpoint 000005: add realtime issue queue and simulator guidance
Phase 81 -- Checkpoint 000006: close realtime operations center review
```
