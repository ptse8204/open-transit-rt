# Phase 102 — Device / AVL Fleet Onboarding V2

## Goal

Improve fleet device onboarding, mapping, stale-device triage, and support handoff.

## Read first

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-90.md
docs/phase-90-control-plane-final-status.md
docs/roadmap-status.md
docs/evidence/redaction-policy.md
docs/evidence/consumer-submissions/status.json
docs/roadmaps/post-90-agency-grade-gtfs-rt-product/00-CODEX-READ-ME-FIRST.md
docs/roadmaps/post-90-agency-grade-gtfs-rt-product/03-master-subagent-operating-manual.md
docs/roadmaps/post-90-agency-grade-gtfs-rt-product/04-validation-and-claim-boundaries.md
```

## Required workflow

Use the Master/sub-agent flow. Master approval is required before
implementation and before closeout.

## Scope

- Keep changes inside Phase 102.
- Preserve protected evidence paths and prepared-only consumer statuses.
- Do not make stronger public claims.
- Add `docs/handoffs/phase-102.md` at closeout.


## Expected deliverables

- `docs/phase-102-device-avl-fleet-onboarding-v2.md`
- `docs/handoffs/phase-102.md`
- focused tests/docs/code as required
- checkpoint report

## Validation

Minimum baseline from `04-validation-and-claim-boundaries.md`; add heavier
checks when code changes.

## Commit pattern

```text
Phase 102 -- Checkpoint 000001: add device / avl fleet onboarding v2 plan
Phase 102 -- Checkpoint 000002: implement primary scoped work
Phase 102 -- Checkpoint 000003: run validation and patch required gaps
Phase 102 -- Checkpoint 000004: close device / avl fleet onboarding v2 review
```
