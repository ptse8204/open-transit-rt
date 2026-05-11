# Phase 62 — Guided Setup and Browser GTFS Import

## Goal

Reduce CLI reliance by adding browser-led setup and GTFS import/validation flows.

## Commit Sequence

```text
Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan
Phase 62 -- Checkpoint 000002: implement guided setup wizard v1
Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow
Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import
```

## Required Work

- Guided setup wizard with agency profile, metadata, GTFS, feeds, telemetry, validators, connectors.
- Authenticated GTFS import by URL and/or upload.
- Reuse existing import logic.
- Show validation status and next actions.
- Keep raw files and private operator data out of committed docs/evidence.

## Non-Goals

- No evidence creation.
- No compliance claim.
- No public feed contract changes unless explicitly approved.
