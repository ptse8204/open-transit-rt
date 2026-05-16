# Phase 129 -- Community Support Feedback And Issue Triage Kit

## Goal

Create community support, feedback, issue triage, release feedback, and support
bundle guidance without claiming adoption or support SLA.

Phase 129 is not a paid support, support SLA, response-time commitment,
hosted-service, production-readiness, compliance, adoption, consumer
acceptance, release-readiness, evidence, or deployment-success proof phase.

## Current Repo Context

- `.github/ISSUE_TEMPLATE/` already has bug, docs, and feature templates with
  secret-safety and claim-boundary checks.
- `docs/support-boundaries.md` defines maintainer/community/operator
  responsibilities.
- `docs/adoption/evaluator-and-contributor-kit.md` now gives evaluator and
  first-contribution entry points.
- `docs/agency-feedback-template.md` gives a structured public-safe agency
  feedback template.

## Scope

- Add or reconcile community support, feedback, release-feedback, issue triage,
  and support-bundle guidance.
- Keep support language community/best-effort and explicit that no SLA, paid
  support, hosted availability, production readiness, adoption, compliance, or
  consumer acceptance is claimed.
- Improve issue-template or docs navigation only where it helps reporters
  provide public-safe, reproducible, scoped reports.
- Do not collect support evidence, contact external parties, move consumer
  statuses, create evidence, or imply response commitments.

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

- Community support and issue triage kit or reconciled equivalent.
- Public-safe release-feedback and support-bundle guidance updates.
- `docs/handoffs/phase-129.md`
- Source-of-truth status updates for Phase 129 closeout.

## Implementation Plan

1. Add this Phase 129 plan and commit checkpoint 000001.
2. Inspect support boundaries, issue templates, feedback docs, release docs,
   support-bundle docs, and wiki support navigation.
3. Add a bounded support/triage kit and update issue or docs links as needed.
4. Run docs/claim-boundary and baseline validation; patch repo-caused
   failures.
5. Close Phase 129 with handoff/status docs and continue immediately to Phase
   130.

## Checkpoint Plan

- `Phase 129 -- Checkpoint 000001: add community support feedback and issue triage kit plan`
- `Phase 129 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 129 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 129 -- Checkpoint 000004: close community support feedback and issue triage kit review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 129 -- Checkpoint 000001: add community support feedback and issue
triage kit plan.

Goal status:
Active. Phase 128 is closed and Phase 129 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/phase-129-community-support-feedback-and-issue-triage-kit.md`.

Validation run:
Initial inspection reviewed the Phase 129 prompt, issue templates, support
boundaries, evaluator kit, and feedback template.

Blocked checks:
Implementation, docs validation, claim-boundary validation, and closeout
validation are scheduled for later Phase 129 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids paid support, support SLA, response-time
commitments, hosted-service availability, production readiness, adoption,
agency approval, compliance, consumer acceptance, final-root readiness, stable
release readiness, vendor compatibility, hardware certification, and
production-grade ETA claims.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, or public route behavior is planned.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 129 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000001 because the phase plan is docs/process work and
does not touch a visual UI surface.

Master review:
Approved. The plan scopes Phase 129 to support and triage guidance without
collecting evidence or making support/adoption claims.

Required edits:
Commit checkpoint 000001, then implement or audit the scoped kit.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 129 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 129 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 129 implemented the scoped community support, feedback, release
feedback, and issue triage kit.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/support/community-support-and-issue-triage-kit.md`,
`.github/ISSUE_TEMPLATE/release_feedback.yml`,
`.github/ISSUE_TEMPLATE/config.yml`, `docs/support-boundaries.md`,
`wiki/support-and-contribute.md`, `docs/README.md`,
`docs/adoption/evaluator-and-contributor-kit.md`, and this phase report.

Implementation summary:
Added a public-safe Community Support And Issue Triage Kit with triage lanes,
reporter and maintainer checklists, release-candidate feedback guidance,
support-bundle sharing boundaries, and bounded public reply patterns. Added a
release-feedback GitHub issue template for `v0.1.0-rc.1` install/download/local
evaluation feedback and linked the kit from issue contact links, support
boundaries, wiki support, docs home, and the evaluator/contributor kit. The
change does not create a support SLA, response-time target, paid support
commitment, hosted-service claim, evidence collection, external contact, or
consumer status movement.

Validation run:
`git diff --check` passed. `scripts/check-consumer-tracker.sh` passed.
Protected-path git status returned no output. All `.github/ISSUE_TEMPLATE/*.yml`
files parsed as YAML. `make audit-product-acceptance` passed. `make
audit-final-claim-review` passed. `make check` passed.

Blocked checks:
None for this checkpoint. Full Phase 129 validation is scheduled for
checkpoint 000003.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
Product acceptance and final claim audits passed. The support kit explicitly
states community support is best-effort and does not claim paid support, SLA,
response-time guarantees, hosted availability, production readiness,
compliance, adoption, agency approval, consumer acceptance, vendor
compatibility, hardware certification, final-root readiness, or ETA quality.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, public route, or browser behavior changed. The issue template and
kit direct suspected vulnerabilities or leaked secrets away from public issues.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 129 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000002 because the implementation was docs and GitHub
issue-template configuration only; it did not touch a visual UI surface.

Master review:
Approved for full validation. The implementation improves triage support
without making support, adoption, release, or production claims.

Required edits:
Run checkpoint 000003 full validation and patch any repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit.

Next checkpoint:
Phase 129 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 129 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Full Phase 129 validation passed with no repo-caused failures.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/phase-129-community-support-feedback-and-issue-triage-kit.md`.

Validation run:
`git status --short` returned clean at validation start. `git diff --check`
passed. `python3 -m json.tool
docs/evidence/consumer-submissions/status.json` passed.
`scripts/check-consumer-tracker.sh` passed. Protected-path git status returned
no output. All `.github/ISSUE_TEMPLATE/*.yml` files parsed as YAML. `make
check` passed. `make audit-product-acceptance` passed. `make
audit-final-claim-review` passed. `make external-connection-check` passed.
`make adapter-conformance` passed. `make gtfsrt-conformance` passed. `make
validate` passed. `make test` passed. `docker compose -f
deploy/docker-compose.yml config` passed. Final `git status --short`, `git
diff --check`, and protected-path git status returned clean.

Blocked checks:
None.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
Product acceptance and final claim audits passed. The kit remains community
best-effort support and triage guidance only.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, public route, or browser behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 129 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000003 because Phase 129 did not touch a visual UI
surface.

Master review:
Approved for Phase 129 closeout.

Required edits:
None.

Decision:
Proceed to checkpoint 000003 commit.

Next checkpoint:
Phase 129 -- Checkpoint 000004: close community support feedback and issue
triage kit review.
