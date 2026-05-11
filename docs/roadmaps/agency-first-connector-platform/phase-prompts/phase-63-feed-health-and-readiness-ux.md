# Phase 63 — Feed Health and Readiness UX

## Goal

Give operators a simple view of feed health and CAL-ITP-style readiness without requiring them to read diagnostic tables.

## Commit Sequence

```text
Phase 63 -- Checkpoint 000001: add feed health and readiness UX plan
Phase 63 -- Checkpoint 000002: implement feed health dashboard
Phase 63 -- Checkpoint 000003: implement readiness checklist v2
Phase 63 -- Checkpoint 000004: close feed health and readiness UX
```

## Required Work

- Plain-language status for `feeds.json`, schedule, Vehicle Positions, Trip Updates, and Alerts.
- Show freshness, last generated/checked time, validator context, and next action.
- Readiness checklist v2 with “what this means / why it matters / what to do next / what it does not prove.”

## Boundaries

- No SLA/uptime proof.
- No compliance claim.
