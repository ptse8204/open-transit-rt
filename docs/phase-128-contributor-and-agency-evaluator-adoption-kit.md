# Phase 128 -- Contributor And Agency Evaluator Adoption Kit

## Goal

Create contributor and agency evaluator adoption kit: no-claim trial guide,
demo paths, feedback templates, and first contribution map.

Phase 128 is not an adoption proof, agency approval, support commitment,
hosted-service, SLA/uptime, production-readiness, compliance, consumer
acceptance, release-readiness, evidence, or deployment-success proof phase.

## Current Repo Context

- `CONTRIBUTING.md`, `docs/contributor-first-issues.md`, and
  `docs/connectors/contributing-connectors.md` already orient contributors.
- `docs/agency-feedback-template.md`, `docs/tutorials/agency-demo-flow.md`,
  `docs/tutorials/no-cli-agency-first-run.md`, and the wiki already provide
  separate evaluator and demo paths.
- Phase 119 aligned public docs and quickstart with the published rc1 release
  candidate.
- Phase 127 closed small-host deployment and upgrade UX hardening.

## Scope

- Add or reconcile a single no-claim evaluator/contributor kit that links
  agency trial paths, demo paths, feedback templates, and first contribution
  paths.
- Keep the kit public-safe and explicit that evaluation, feedback, and
  contribution do not prove adoption, approval, compliance, consumer
  acceptance, support SLA, hosted availability, or production readiness.
- Improve public docs navigation only where it helps evaluators and
  contributors find the kit.
- Do not collect real agency feedback, contact external parties, move consumer
  statuses, create evidence, or imply maintainer support commitments.

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

- Contributor and agency evaluator adoption kit or reconciled equivalent.
- Public-safe feedback and first-contribution guidance updates.
- `docs/handoffs/phase-128.md`
- Source-of-truth status updates for Phase 128 closeout.

## Implementation Plan

1. Add this Phase 128 plan and commit checkpoint 000001.
2. Inspect current contributor, evaluator, demo, feedback, issue-template, and
   wiki navigation docs.
3. Add a no-claim adoption kit that consolidates evaluator flows, feedback
   templates, demo paths, and first contribution map.
4. Run docs/claim-boundary and baseline validation; patch repo-caused
   failures.
5. Close Phase 128 with handoff/status docs and continue immediately to Phase
   129.

## Checkpoint Plan

- `Phase 128 -- Checkpoint 000001: add contributor and agency evaluator adoption kit plan`
- `Phase 128 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 128 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 128 -- Checkpoint 000004: close contributor and agency evaluator adoption kit review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 128 -- Checkpoint 000001: add contributor and agency evaluator adoption
kit plan.

Goal status:
Active. Phase 127 is closed and Phase 128 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/phase-128-contributor-and-agency-evaluator-adoption-kit.md`.

Validation run:
Initial inspection reviewed the Phase 128 prompt, current contributing docs,
agency/evaluator docs, feedback templates, issue templates, wiki navigation
candidates, and protected-path boundaries.

Blocked checks:
Implementation, docs validation, claim-boundary validation, and closeout
validation are scheduled for later Phase 128 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids adoption proof, agency approval, support SLA,
hosted-service availability, production readiness, compliance, consumer
acceptance, consumer ingestion/listing/display, final-root readiness, stable
release readiness, vendor compatibility, hardware certification, and
production-grade ETA claims.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, or public route behavior is planned.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 128 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000001 because the phase plan is docs/process work and
does not touch a visual UI surface.

Master review:
Approved. The plan scopes Phase 128 to public-safe evaluator and contributor
guidance without collecting evidence or making adoption/support claims.

Required edits:
Commit checkpoint 000001, then implement or audit the scoped kit.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 128 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 128 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 128 implemented the public-safe evaluator and contributor kit.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/adoption/evaluator-and-contributor-kit.md`, `README.md`,
`docs/README.md`, `CONTRIBUTING.md`, `wiki/support-and-contribute.md`,
`docs/agency-feedback-template.md`, and this phase report.

Implementation summary:
Added a no-claim Evaluator And Contributor Kit under `docs/adoption/` and
linked it from the README, docs home, CONTRIBUTING guide, wiki support page,
and agency feedback template. The kit consolidates local evaluator paths,
release-candidate install trial guidance, synthetic connector review paths,
feedback template guidance, demo links, first-contribution examples, required
checks, and explicit boundaries. It does not collect feedback, contact
external parties, create evidence, touch consumer statuses, add issue-template
automation, or imply adoption/support/compliance/production claims.

Validation run:
`git diff --check` passed. `scripts/check-consumer-tracker.sh` passed.
Protected-path git status returned no output. `make audit-product-acceptance`
passed. `make audit-final-claim-review` passed. `make check` passed.

Blocked checks:
None for this checkpoint. Full Phase 128 validation is scheduled for
checkpoint 000003.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
Product acceptance and final claim audits passed. The new kit explicitly
states evaluation, feedback, stars, discussions, issues, and PRs are useful
project signals but not proof of adoption, compliance, production readiness,
consumer acceptance, hosted availability, SLA coverage, or vendor/hardware
compatibility.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, public route, or browser behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 128 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000002 because the implementation was public docs and
navigation only; it did not touch a visual UI surface.

Master review:
Approved for full validation. The implementation improves adoption support
without making adoption or support claims.

Required edits:
Run checkpoint 000003 full validation and patch any repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit.

Next checkpoint:
Phase 128 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 128 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Full Phase 128 validation passed with no repo-caused failures.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/phase-128-contributor-and-agency-evaluator-adoption-kit.md`.

Validation run:
`git status --short` returned clean at validation start. `git diff --check`
passed. `python3 -m json.tool
docs/evidence/consumer-submissions/status.json` passed.
`scripts/check-consumer-tracker.sh` passed. Protected-path git status returned
no output. `make check` passed. `make audit-product-acceptance` passed. `make
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
Product acceptance and final claim audits passed. The kit remains no-claim
evaluation and contribution guidance only.

Security/auth status:
No auth, CSRF, credential, token, private payload, support bundle, retained
evidence, public route, or browser behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 128 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000003 because Phase 128 did not touch a visual UI
surface.

Master review:
Approved for Phase 128 closeout.

Required edits:
None.

Decision:
Proceed to checkpoint 000003 commit.

Next checkpoint:
Phase 128 -- Checkpoint 000004: close contributor and agency evaluator
adoption kit review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 128 -- Checkpoint 000004: close contributor and agency evaluator
adoption kit review.

Goal status:
Active. Phase 128 is closed and the goal continues to Phase 129.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release, Install Confidence, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Closeout summary:
Phase 128 added `docs/adoption/evaluator-and-contributor-kit.md` and linked it
from public-facing docs. The kit consolidates no-claim trial paths, demo links,
feedback guidance, and first-contribution guidance without collecting
feedback, contacting external parties, creating evidence, moving consumer
statuses, or implying adoption/support/compliance/production claims.

Changed files:
`docs/handoffs/phase-128.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`; and this phase
report.

Validation run:
Full Phase 128 validation passed before closeout docs. Focused closeout
validation passed after closeout docs: `git diff --check`, `make check`,
`make audit-product-acceptance`, `make audit-final-claim-review`,
`scripts/check-consumer-tracker.sh`, and protected-path git status.

Blocked checks:
No Phase 128 check remains blocked.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 128 remains bounded to no-claim evaluator and contributor guidance and
makes no adoption proof, agency approval, support commitment, stable release
readiness, deployment success, production readiness, compliance, consumer
acceptance, consumer ingestion/listing/display, final-root readiness, hosted
service availability, paid support, SLA/uptime, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality,
real-world ETA accuracy, or consumer display claim.

Security/auth status:
No application security behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was added.

Release/publication status:
The public rc1 prerelease remains published. No release action was taken.

Install confidence status:
Public fresh-clone rc1 install confidence remains passed.

Web design skill status:
Not used; no visual UI changed.

Master review:
Approved. Phase 128 closes with a validated evaluator and contributor kit.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 129.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 129.

Next checkpoint:
Phase 129 -- Checkpoint 000001: add community support feedback and issue
triage kit plan.
