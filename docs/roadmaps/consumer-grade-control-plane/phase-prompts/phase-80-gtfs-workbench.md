# Phase 80 Prompt — GTFS Workbench

## Goal

Turn GTFS import/edit/publish from a developer-centered workflow into an operator workbench.

## Scope

- Active vs draft schedule summary.
- Import diff view.
- Stops/routes/trips/calendar preview tables.
- Optional map previews for stops/routes/shapes if already practical with committed data.
- GTFS Quality issue drilldowns.
- Suggested fix categories and likely owners.
- Draft edit/publish review improvements.
- Rollback/history UX.

## Feature rules

- Preserve existing GTFS import semantics.
- Do not auto-fix GTFS without explicit review.
- Do not publish without confirmation.
- Do not expose raw validator reports.
- Do not rely on external map services unless separately reviewed for privacy/network impact.
- Do not claim validator-clean, compliance, or consumer acceptance.

## Operator usability features

- Plain-language “what this file is” descriptions for GTFS files.
- Examples of common fixes.
- Before/after publish review.
- Undo/rollback guidance.
- Export/download of public-safe summaries only where safe.

## Validation

Baseline validation plus:

```bash
make validate
make test
make audit-final-claim-review
```

## Commits

```text
Phase 80 -- Checkpoint 000001: add GTFS Workbench plan
Phase 80 -- Checkpoint 000002: add import diff and schedule summaries
Phase 80 -- Checkpoint 000003: add GTFS preview tables and filters
Phase 80 -- Checkpoint 000004: improve safe draft publish review
Phase 80 -- Checkpoint 000005: add rollback and schedule history UX
Phase 80 -- Checkpoint 000006: close GTFS Workbench review
```
