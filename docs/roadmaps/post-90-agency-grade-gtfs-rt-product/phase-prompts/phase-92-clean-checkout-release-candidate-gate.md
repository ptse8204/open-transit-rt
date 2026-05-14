# Phase 92 — Clean Checkout Release-Candidate Gate

## Goal

Run the clean checkout product gate and record exact release blockers.

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

- Keep changes inside Phase 92.
- Preserve protected evidence paths and prepared-only consumer statuses.
- Do not make stronger public claims.
- Add `docs/handoffs/phase-92.md` at closeout.


## Expected deliverables

- `docs/phase-92-clean-checkout-release-candidate-gate.md`
- `docs/handoffs/phase-92.md`
- focused tests/docs/code as required
- checkpoint report

## Validation

Minimum baseline from `04-validation-and-claim-boundaries.md`; add heavier
checks when code changes.

## Commit pattern

```text
Phase 92 -- Checkpoint 000001: add clean checkout release-candidate gate plan
Phase 92 -- Checkpoint 000002: implement primary scoped work
Phase 92 -- Checkpoint 000003: run validation and patch required gaps
Phase 92 -- Checkpoint 000004: close clean checkout release-candidate gate review
```
