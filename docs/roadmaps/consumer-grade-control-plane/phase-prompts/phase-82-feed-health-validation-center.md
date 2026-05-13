# Phase 82 Prompt — Feed Health And Validation Center

## Goal

Unify feed status, validator health, GTFS quality triage, readiness signals, and blocker history into one operator-friendly Validation Center.

## Scope

- Feed status overview for schedule, Vehicle Positions, Trip Updates, Alerts, and `feeds.json`.
- Validation history summary.
- Validator-health drilldowns.
- GTFS Quality issue owner/fix guidance.
- Readiness timeline with `ok`, `needs_review`, `missing`, `blocked`, or `unknown` states.
- Clear “does not prove” labels.

## Important boundaries

- Validator success is a supporting signal only.
- Do not claim CAL-ITP/Caltrans compliance.
- Do not claim consumer acceptance.
- Do not block publish unless existing logic already does.
- Do not expose raw validator stdout/stderr/private paths.

## Validation

Baseline validation plus:

```bash
make validate
make test
make audit-final-claim-review
```

## Commits

```text
Phase 82 -- Checkpoint 000001: add feed health and validation center plan
Phase 82 -- Checkpoint 000002: unify feed status and validator history views
Phase 82 -- Checkpoint 000003: add validation issue drilldowns and fix-owner guidance
Phase 82 -- Checkpoint 000004: add readiness timeline and blocker history
Phase 82 -- Checkpoint 000005: close feed health and validation center review
```
