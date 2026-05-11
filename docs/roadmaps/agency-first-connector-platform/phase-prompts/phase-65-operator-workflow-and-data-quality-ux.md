# Phase 65 — Operator Workflow and Data Quality UX

## Goal

Make day-to-day agency operations easier: devices, telemetry simulator, and GTFS quality fixes.

## Commit Sequence

```text
Phase 65 -- Checkpoint 000001: add operator workflow and data quality UX plan
Phase 65 -- Checkpoint 000002: implement device and vehicle onboarding UI
Phase 65 -- Checkpoint 000003: implement telemetry simulator UI
Phase 65 -- Checkpoint 000004: implement GTFS quality fix guidance UI
Phase 65 -- Checkpoint 000005: close operator workflow and data quality UX
```

## Required Work

- Bind/rotate devices through UI.
- Show generated device token once.
- Show latest telemetry per device/vehicle.
- UI to run or guide synthetic telemetry simulator scenarios.
- GTFS quality guidance for common validation/import issues.

## Boundaries

- Do not expose secrets after initial generation.
- Do not claim real AVL reliability from synthetic simulator.
- Do not auto-edit GTFS without explicit operator action.
