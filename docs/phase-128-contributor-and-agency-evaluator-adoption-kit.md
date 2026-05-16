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
