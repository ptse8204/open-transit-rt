# Phase 115 -- v0.1.0-rc.1 Public Release Cut

## Goal

Attempt to publish an actual public `v0.1.0-rc.1` release candidate if all
release gates pass and authenticated release tooling is available. If any gate
blocks publication, record exact blocker evidence and continue safe downstream
phases without faking a release.

Phase 115 is authorized for a public release candidate only. It does not
authorize a stable release, production readiness, compliance, adoption,
consumer acceptance, hosted SaaS, SLA/uptime, vendor compatibility, hardware
certification, final-root readiness, or production-grade ETA quality.

## Current Repo Context

- Phase 114 is closed with Web Design Skill UX validation and private
  Operations Console polish.
- Phase 113 recorded bounded local fresh-clone and local source-archive
  install-confidence passes.
- Phase 112 recorded `blocked_public_distribution_review` because the source
  archive contained tracked protected evidence and consumer-submission paths.
- `gh` tooling and repository permission were available during Phase 112, but
  must be rechecked before publication.
- Local release status is documented in
  `docs/release-status-v0.1.0-rc.1.md`.

## Scope

- Re-evaluate the Phase 112 public-distribution blocker from the current
  release commit.
- If safe, add a source-archive export policy that excludes protected evidence
  paths from generated source archives without modifying protected paths.
- Regenerate and audit the local `v0.1.0-rc.1` package from the intended
  release commit.
- Rerun release-candidate, claim, protected-path, and prepared-only consumer
  gates.
- Publish the GitHub Release only if all gates pass and `gh` remains
  authenticated/authorized.
- If blocked, record the exact blocker and continue to Phase 116 without
  pretending publication occurred.

## Protected Paths

Do not modify, reformat, delete, stage, or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

The consumer tracker must remain exactly seven targets in order and all
`prepared`.

## Release Gate Checklist

Publication requires all of the following to pass from the intended release
commit:

- clean worktree;
- `make test-release-package`;
- strict local release package generation;
- local release package audit;
- source archive protected-path scan returns zero protected-path entries;
- release notes are bounded and truthful for the exact commit;
- `make audit-product-acceptance`;
- `make audit-final-claim-review`;
- `make check`;
- `make validate`;
- `make test`;
- `docker compose -f deploy/docker-compose.yml config`;
- `RUN_LOCAL_APP=true RELEASE_PACKAGE_DIR=... RUN_RELEASE_PACKAGE=true make release-candidate-check`;
- protected-path git status is clean;
- exact prepared-only consumer tracker assertion passes;
- `gh` is installed, authenticated, and authorized for the public repository;
- remote tag `v0.1.0-rc.1` and release are absent before publication.

## Publication Commands If Gates Pass

Use actual package asset paths produced by the gate run.

```bash
git tag -a v0.1.0-rc.1 -m "Open Transit RT v0.1.0-rc.1"
git push origin v0.1.0-rc.1
gh release create v0.1.0-rc.1 \
  --repo ptse8204/open-transit-rt \
  --title "Open Transit RT v0.1.0-rc.1" \
  --notes-file docs/release-notes-v0.1.0-rc.1-draft.md \
  <audited-release-assets>
```

Do not run these commands unless every gate passes.

## Deliverables

- Phase 115 release status update in `docs/release-status-v0.1.0-rc.1.md`.
- A regenerated local package record under ignored `.cache` if gates reach
  package generation.
- Published GitHub Release evidence if publication succeeds, or exact blocked
  publication evidence if it does not.
- `docs/handoffs/phase-115.md`.
- Source-of-truth status updates.

## Implementation Plan

1. Add this Phase 115 plan and commit checkpoint 000001.
2. Re-audit source-archive blocker and, if safe, implement a release archive
   export policy without touching protected paths.
3. Run package, release-candidate, claim, consumer, protected-path, and
   code/test gates.
4. Publish public `v0.1.0-rc.1` only if all gates pass; otherwise record exact
   blocker.
5. Close Phase 115 with handoff/status docs and continue immediately to Phase
   116.

## Checkpoint Plan

- `Phase 115 -- Checkpoint 000001: add v0.1.0-rc.1 public release cut plan`
- `Phase 115 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 115 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 115 -- Checkpoint 000001: add v0.1.0-rc.1 public release cut plan.

Goal status:
Active. Phase 114 is closed and Phase 115 has started.

Sub-agents used or simulated:
Release/Supply-Chain sub-agent review was delegated for source archive and
publication blocker audit. Planning, Implementation, QA, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Install Confidence, Web Design
Skill, and GTFS-RT Domain roles are simulated by the Master Agent for this
plan checkpoint.

Changed files:
`docs/phase-115-v0.1.0-rc.1-public-release-cut.md`.

Validation run:
Initial Phase 115 inspection reviewed the Phase 115 prompt, public release
policy, current release status artifact, release package scripts, audit
scripts, and existing release packaging decisions/dependency documentation.
Focused checkpoint validation is scheduled before commit.

Blocked checks:
Package generation, archive review, GitHub tooling recheck, tag/release
publication, release-candidate diagnostics, and full validation are scheduled
for later Phase 115 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release, production readiness, compliance,
adoption, consumer acceptance, final-root readiness, hosted service, paid
support, SLA/uptime, vendor compatibility, hardware certification, production
AVL reliability, and ETA-quality claims.

Security/auth status:
The plan requires `gh` auth and repository permission recheck before any
publication. It does not change route auth, CSRF, credentials, token handling,
public exposure, private payload handling, external contact, or command
execution behavior.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
is planned unless release archive policy docs/scripts require bounded updates.

Release/publication status:
No release action is taken in checkpoint 000001. Publication remains blocked
until Phase 115 gates are re-run and pass.

Install confidence status:
Phase 113 local install-confidence diagnostics remain the current result.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 115 does not touch UX.

Master review:
Approved. The plan preserves the public release authorization while keeping
the source-archive public-distribution gate decisive.

Required edits:
Run checkpoint 000001 validation and commit, then re-audit and, if safe,
address the source archive blocker.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 115 -- Checkpoint 000002: implement or audit primary scoped work.

