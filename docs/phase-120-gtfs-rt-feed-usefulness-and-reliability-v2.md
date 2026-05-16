# Phase 120 -- GTFS-RT Feed Usefulness And Reliability V2

## Goal

Improve GTFS-RT feed usefulness and reliability review for Vehicle Positions,
Trip Updates, and Alerts without making production readiness, SLA, compliance,
consumer acceptance, vendor compatibility, hardware certification, or
production-grade ETA quality claims.

## Current Repo Context

- Phase 115 published public `v0.1.0-rc.1` for local/self-hosted evaluation.
- Phase 116 recorded public download replay and the source-archive `make check`
  limitation for the already-published rc1 archives.
- Phase 117 verified public fresh-clone install confidence from the rc1 tag.
- Phase 118 completed Web Design Skill UX validation.
- Phase 119 aligned public docs and quickstarts to the actual rc1 release and
  install-confidence state.
- Existing feed packages already emit Vehicle Positions, Trip Updates, Alerts,
  and schedule metadata; existing Operations Console routes already show
  private feed health and realtime review state.

## Scope

- Add or improve local/synthetic feed usefulness and reliability diagnostics for
  GTFS-RT outputs.
- Prefer structured helper logic, tests, and bounded Operations Console/readme
  surfaces over ad hoc wording.
- Keep Vehicle Positions first, and keep Trip Updates pluggable and
  fail-closed when unsupported.
- Preserve all protected evidence paths and prepared-only consumer tracker
  statuses.

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

- A bounded implementation or audit artifact that improves GTFS-RT feed
  usefulness/reliability review for local operators.
- Focused tests or validation covering the changed behavior.
- `docs/handoffs/phase-120.md`
- Source-of-truth status updates for Phase 120 closeout.

## Implementation Plan

1. Add this Phase 120 plan and commit checkpoint 000001.
2. Inspect feed/realtimequality/Operations Console code and implement one
   conservative, test-backed usefulness/reliability improvement.
3. Run relevant GTFS-RT, connector, claim-boundary, and baseline validation;
   patch repo-caused failures.
4. Close Phase 120 with handoff/status docs and continue immediately to
   Phase 121.

## Checkpoint Plan

- `Phase 120 -- Checkpoint 000001: add gtfs rt feed usefulness and reliability v2 plan`
- `Phase 120 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 120 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 120 -- Checkpoint 000004: close gtfs rt feed usefulness and reliability v2 review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 120 -- Checkpoint 000001: add GTFS-RT feed usefulness and reliability
v2 plan.

Goal status:
Active. Phase 119 is closed and Phase 120 has started.

Sub-agents used or simulated:
The environment refused an additional GTFS-RT explorer because the agent
thread limit was reached. Context / Repo Truth, Planning, Implementation, QA,
GTFS-RT Domain, Claim-Boundary, Security/Auth, Data/Migration,
Documentation / IA, Web Design Skill, Release, and Install Confidence roles
are simulated by the Master Agent for this checkpoint.

Changed files:
`docs/phase-120-gtfs-rt-feed-usefulness-and-reliability-v2.md`.

Validation run:
Initial inspection reviewed the Phase 120 roadmap prompt, validation
boundaries, master/sub-agent workflow, feed/realtimequality packages,
Operations Console route files, scripts, and Makefile GTFS-RT targets.

Blocked checks:
Implementation, tests, connector validation, and closeout validation are
scheduled for later Phase 120 checkpoints.

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
behavior without focused review.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change is
planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 120 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 120 does not plan
visual UX changes unless implementation inspection shows a tightly scoped
Operations Console review improvement.

Master review:
Approved. The plan keeps Phase 120 inside local/synthetic GTFS-RT usefulness
and reliability review, preserving all evidence and claim boundaries.

Required edits:
Commit checkpoint 000001, then implement or audit the scoped GTFS-RT
improvement.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 120 -- Checkpoint 000002: implement or audit primary scoped work.
