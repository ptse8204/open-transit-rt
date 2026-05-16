# Phase 126 -- Operator Assistant Safe Command Expansion

## Goal

Expand safe command-backed operator assistance through `internal/admincontrol`
without arbitrary shell execution or destructive browser actions.

Phase 126 is not a production operations, hosted-service, SLA/uptime,
production-readiness, compliance, consumer acceptance, public launch, vendor
compatibility, hardware certification, adoption, release-readiness, or
evidence proof phase.

## Current Repo Context

- `internal/admincontrol` defines bounded command definitions and result
  shapes.
- `validation_health.refresh` is implemented as a read-only private command.
- `validation_health.run_all` is documented as a future/admin private
  diagnostic write definition.
- Browser command requests already reject shell/argv/path/timeout/raw-report
  fields for the implemented validation-health refresh route.

## Scope

- Add or reconcile a safe operator-assistant command catalog in
  `internal/admincontrol`.
- Keep definitions server-owned, role-scoped, claim-bounded, and explicit about
  public-feed impact, private impact, rollback/review paths, and non-claims.
- Do not add arbitrary shell execution, browser-provided command fields,
  destructive actions, public command routes, evidence writes, or consumer
  status changes.

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

- Safe operator-assistant command expansion or audit.
- Focused `internal/admincontrol` tests.
- `docs/handoffs/phase-126.md`
- Source-of-truth status updates for Phase 126 closeout.

## Implementation Plan

1. Add this Phase 126 plan and commit checkpoint 000001.
2. Inspect `internal/admincontrol`, command route tests, and
   `docs/admin-command-model.md`.
3. Implement a bounded command catalog or equivalent safe-command audit and
   tests.
4. Run relevant admincontrol, command-route, claim-boundary, and baseline
   validation; patch repo-caused failures.
5. Close Phase 126 with handoff/status docs and continue immediately to Phase
   127.

## Checkpoint Plan

- `Phase 126 -- Checkpoint 000001: add operator assistant safe command expansion plan`
- `Phase 126 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 126 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 126 -- Checkpoint 000004: close operator assistant safe command expansion review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 126 -- Checkpoint 000001: add operator assistant safe command expansion
plan.

Goal status:
Active. Phase 125 is closed and Phase 126 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-126-operator-assistant-safe-command-expansion.md`.

Validation run:
Initial inspection reviewed the Phase 126 prompt, `internal/admincontrol`,
command route tests, and `docs/admin-command-model.md`.

Blocked checks:
Implementation, tests, command-boundary validation, and closeout validation
are scheduled for later Phase 126 checkpoints.

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
The plan keeps command definitions server-owned and does not add arbitrary
shell execution, browser-supplied argv/path fields, public command routes, or
destructive browser actions.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 126 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. Phase 126 does not plan
visual UX changes.

Master review:
Approved. The plan scopes Phase 126 to safe-command definitions and validation
without expanding browser execution.

Required edits:
Commit checkpoint 000001, then implement the scoped safe-command work.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 126 -- Checkpoint 000002: implement or audit primary scoped work.
