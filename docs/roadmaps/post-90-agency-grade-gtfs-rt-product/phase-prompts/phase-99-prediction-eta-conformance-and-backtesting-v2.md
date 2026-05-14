# Phase 99 — Prediction / ETA Conformance And Backtesting V2

## Goal

Improve predictor conformance, shadow comparison, and backtesting without ETA quality claims.

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

- Keep changes inside Phase 99.
- Preserve protected evidence paths and prepared-only consumer statuses.
- Do not make stronger public claims.
- Add `docs/handoffs/phase-99.md` at closeout.


## Expected deliverables

- `docs/phase-99-prediction-eta-conformance-and-backtesting-v2.md`
- `docs/handoffs/phase-99.md`
- focused tests/docs/code as required
- checkpoint report

## Validation

Minimum baseline from `04-validation-and-claim-boundaries.md`; add heavier
checks when code changes.

## Commit pattern

```text
Phase 99 -- Checkpoint 000001: add prediction / eta conformance and backtesting v2 plan
Phase 99 -- Checkpoint 000002: implement primary scoped work
Phase 99 -- Checkpoint 000003: run validation and patch required gaps
Phase 99 -- Checkpoint 000004: close prediction / eta conformance and backtesting v2 review
```
