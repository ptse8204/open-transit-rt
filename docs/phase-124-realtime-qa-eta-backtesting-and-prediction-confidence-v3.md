# Phase 124 -- Realtime QA ETA Backtesting And Prediction Confidence V3

## Goal

Improve realtime QA, ETA backtesting, and prediction confidence reporting with
conservative aggregate metrics and no real-world ETA-quality claim.

Phase 124 is not a production-grade ETA quality, real-world ETA accuracy,
vendor compatibility, hardware certification, production readiness,
compliance, consumer acceptance, SLA/uptime, hosted-service, adoption, or
release-readiness proof phase.

## Current Repo Context

- `internal/realtimequality` already provides private aggregate backtesting
  over observed events and prediction samples.
- `cmd/realtime-quality-backtest` writes aggregate local diagnostics under
  `.cache/` by default and forbids evidence paths.
- Phase 99 added synthetic conformance signals for prediction backtests.
- Phase 120 improved Vehicle Positions aggregate review summaries.
- Phase 121 and Phase 122 added GTFS-RT conformance and fixture coverage.

## Scope

- Add or reconcile conservative prediction confidence / backtest reporting
  signals that make uncertainty and withholding visible.
- Keep outputs aggregate-only and private/local by default.
- Preserve the Trip Updates adapter boundary and Vehicle Positions
  independence.
- Avoid public API, durable persistence, evidence writes, consumer-status
  changes, or production ETA claims.

## Protected Paths

Do not modify, reformat, delete, stage, or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

The consumer tracker must remain exactly seven targets in order and all
`prepared`.

## Deliverables

- Backtest/prediction confidence reporting improvement or audit.
- Focused tests for conservative aggregate signals.
- `docs/handoffs/phase-124.md`
- Source-of-truth status updates for Phase 124 closeout.

## Implementation Plan

1. Add this Phase 124 plan and commit checkpoint 000001.
2. Inspect `internal/realtimequality`, `internal/prediction`, and
   `cmd/realtime-quality-backtest` for the highest-value conservative V3 gap.
3. Implement the scoped aggregate confidence/reporting improvement and tests.
4. Run relevant backtest, prediction, GTFS-RT, claim-boundary, and baseline
   validation; patch repo-caused failures.
5. Close Phase 124 with handoff/status docs and continue immediately to Phase
   125.

## Checkpoint Plan

- `Phase 124 -- Checkpoint 000001: add realtime qa eta backtesting and prediction confidence v3 plan`
- `Phase 124 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 124 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 124 -- Checkpoint 000004: close realtime qa eta backtesting and prediction confidence v3 review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 124 -- Checkpoint 000001: add realtime QA ETA backtesting and prediction
confidence V3 plan.

Goal status:
Active. Phase 123 is closed and Phase 124 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-124-realtime-qa-eta-backtesting-and-prediction-confidence-v3.md`.

Validation run:
Initial inspection reviewed the Phase 124 prompt, current status and handoff
context, `internal/realtimequality`, `internal/prediction`, and
`cmd/realtime-quality-backtest`.

Blocked checks:
Implementation, tests, realtime quality validation, and closeout validation
are scheduled for later Phase 124 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
compliance, adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, and real-world ETA
accuracy claims.

Security/auth status:
The plan does not change route auth, CSRF behavior, credential handling, token
handling, public exposure, private payload handling, or operator command
behavior.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 124 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 124 does not plan
visual UX changes unless implementation inspection exposes a UX-specific gap.

Master review:
Approved. The plan scopes Phase 124 to aggregate conservative reporting and
claim-boundary-preserving validation.

Required edits:
Commit checkpoint 000001, then implement the scoped realtime QA / prediction
confidence work.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 124 -- Checkpoint 000002: implement or audit primary scoped work.
