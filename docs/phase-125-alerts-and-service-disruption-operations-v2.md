# Phase 125 -- Alerts And Service Disruption Operations V2

## Goal

Improve Alerts and service disruption operations, lifecycle review, missing
alert hints, validation guidance, and cancellation linkage.

Phase 125 is not a production operations, compliance, consumer acceptance,
public launch, hosted-service, SLA/uptime, production-readiness, vendor
compatibility, hardware certification, adoption, or release-readiness proof
phase.

## Current Repo Context

- `internal/alerts` models alert records and cancellation reconciliation.
- `internal/feed/alerts` builds GTFS-RT Alerts snapshots and debug summaries.
- `cmd/feed-alerts` serves public Alerts feeds and private Alerts admin /
  console routes.
- Existing tests cover console rendering, reconciliation, private route
  boundaries, agency scope, and mutation role checks.

## Scope

- Add or reconcile service-disruption operational review signals for alerts.
- Improve lifecycle / missing-alert / cancellation-linkage guidance using
  synthetic/local or private aggregate data only.
- Keep Alerts operations authenticated and claim-bounded.
- Avoid public API expansion, evidence writes, consumer-status movement,
  external contact, or production operations claims.

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

- Alerts/service-disruption operations improvement or audit.
- Focused tests for lifecycle or disruption review behavior.
- `docs/handoffs/phase-125.md`
- Source-of-truth status updates for Phase 125 closeout.

## Implementation Plan

1. Add this Phase 125 plan and commit checkpoint 000001.
2. Inspect `internal/alerts`, `internal/feed/alerts`, and `cmd/feed-alerts`
   for the highest-value V2 gap.
3. Implement the scoped alerts operations improvement and tests.
4. Run relevant alerts, GTFS-RT, claim-boundary, and baseline validation;
   patch repo-caused failures.
5. Close Phase 125 with handoff/status docs and continue immediately to Phase
   126.

## Checkpoint Plan

- `Phase 125 -- Checkpoint 000001: add alerts and service disruption operations v2 plan`
- `Phase 125 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 125 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 125 -- Checkpoint 000004: close alerts and service disruption operations v2 review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 125 -- Checkpoint 000001: add alerts and service disruption operations
V2 plan.

Goal status:
Active. Phase 124 is closed and Phase 125 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`docs/phase-125-alerts-and-service-disruption-operations-v2.md`.

Validation run:
Initial inspection reviewed the Phase 125 prompt, current status and handoff
context, `internal/alerts`, `internal/feed/alerts`, `cmd/feed-alerts`, and
existing alerts console tests.

Blocked checks:
Implementation, tests, alerts validation, and closeout validation are
scheduled for later Phase 125 checkpoints.

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
The plan keeps Alerts operations authenticated and does not change credential
handling, public exposure, or evidence/consumer status behavior.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 125 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Phase 118 Web Design Skill artifact remains complete. If Phase 125 changes
visual Alerts Console UX, the Web Design Skill will be used before that UX
work.

Master review:
Approved. The plan scopes Phase 125 to private alerts operations hardening and
claim-boundary-preserving validation.

Required edits:
Commit checkpoint 000001, then implement the scoped alerts operations work.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 125 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 125 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 125 implemented the scoped Alerts/service-disruption operations
review improvement.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, GTFS-RT Domain, Connector, Claim-Boundary,
Security/Auth, Data/Migration, Documentation / IA, Web Design Skill, Release,
and Install Confidence roles are simulated by the Master Agent.

Changed files:
`cmd/feed-alerts/main.go`, `cmd/feed-alerts/main_test.go`, and this phase
report.

Implementation summary:
The private Alerts Console now includes a read-only `Service Disruption
Review` table. It summarizes active published alerts, draft alerts awaiting
review, expired published alerts, published alerts without end windows,
agency-wide/unscoped alerts, alerts without entities, and cancellation
reconciler pairing. The review uses existing private alert records only and
adds no browser-side network action, public route, external contact, evidence
write, or consumer-status mutation.

Validation run:
The Web Design Skill was loaded before editing the Alerts Console UX surface.
`gofmt` passed on touched Go files. `git diff --check` passed. `go test
./cmd/feed-alerts ./internal/alerts ./internal/feed/alerts` passed. `make
gtfsrt-conformance` passed. `scripts/check-consumer-tracker.sh` passed.

Blocked checks:
None for this checkpoint. Full repo validation is scheduled for checkpoint
000003.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
The new review rows explicitly avoid consumer display, agency approval, public
launch, compliance, production readiness, real-world disruption handling
quality, SLA coverage, and consumer ingestion claims.

Security/auth status:
No route auth, CSRF behavior, credential handling, token handling, public
exposure, private payload handling, or operator command behavior was changed.
The new review is rendered only inside the existing authenticated private
Alerts Console.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 125 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
The Web Design Skill at `~/.agents/skills/web-design-engineer` was used for
the Alerts Console UX change. The patch intentionally preserved the existing
dense private admin table style without adding visual decoration, motion, or
new browser interactions.

Master review:
Approved for full validation. The implementation improves service-disruption
review while preserving authenticated, private, bounded-alerts operations.

Required edits:
Run checkpoint 000003 full validation and patch any repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit.

Next checkpoint:
Phase 125 -- Checkpoint 000003: run validation and patch required gaps.
