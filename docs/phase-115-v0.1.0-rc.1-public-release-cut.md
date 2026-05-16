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

## Checkpoint Report -- 000002

Checkpoint:
Phase 115 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. The source archive public-distribution blocker has a proposed safe
fix and is ready for post-commit release gate validation.

Sub-agents used or simulated:
Release/Supply-Chain sub-agent confirmed that a committed root `.gitattributes`
export policy should exclude protected paths from `git archive HEAD` without
editing protected paths. Implementation, QA, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Install Confidence,
Web Design Skill, and GTFS-RT Domain roles were simulated by the Master Agent.

Changed files:
`.gitattributes`; `scripts/test-release-package.sh`; `docs/decisions.md`;
`docs/dependencies.md`; `docs/release-notes-v0.1.0-rc.1-draft.md`;
`docs/phase-115-v0.1.0-rc.1-public-release-cut.md`.

Validation run:
Focused validation is scheduled for checkpoint 000003 because `git archive
HEAD` will not use the new `.gitattributes` policy until this checkpoint is
committed. Pre-commit `make check`, consumer tracker assertion, protected-path
status check, and `git diff --check` are rerun before checkpoint 000002 commit.

Blocked checks:
Final source archive scan, strict package generation, package audit, full
release-candidate gates, GitHub auth recheck, tag push, and GitHub Release
creation are scheduled for checkpoint 000003 after this policy is committed.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
new export policy affects source archive output only.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and must remain `prepared`.

Claim-boundary status:
Release notes remain bounded and explicitly withhold publication/release
claims until Phase 115 gates pass.

Security/auth status:
No route, auth behavior, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or command
execution behavior changed.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
was added.

Release/publication status:
No tag, tag push, GitHub Release, asset upload, or public publication action
was taken in checkpoint 000002.

Install confidence status:
Phase 113 remains the current local install-confidence result.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 115 does not touch UX.

Master review:
Approved. The implementation resolves the known public-distribution blocker at
the archive policy layer while preserving protected paths and claim
boundaries.

Required edits:
Commit checkpoint 000002, then regenerate and audit the package from committed
HEAD before deciding whether to publish.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 release gate
validation.

Next checkpoint:
Phase 115 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 115 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Release gates passed through local validation, strict package audit,
protected archive scan, release-candidate diagnostics, and GitHub publication
preflight; the release notes were patched to remove local-draft wording before
tagging.

Sub-agents used or simulated:
Release/Supply-Chain sub-agent findings were incorporated. QA,
Claim-Boundary, Security/Auth, Data/Migration, Documentation / IA,
Install Confidence, Web Design Skill, GTFS-RT Domain, and Implementation roles
were simulated by the Master Agent for this validation checkpoint.

Changed files:
`docs/release-notes-v0.1.0-rc.1-draft.md`;
`docs/phase-115-v0.1.0-rc.1-public-release-cut.md`.

Validation run:
Passed before the release-notes wording patch:

- `make test-release-package`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- direct `git archive HEAD` protected-path scan: `0`
- strict `v0.1.0-rc.1` package generation from commit
  `10a36dc73a8533ba81ec7f7d6c2b5324b2ee70c5`
- `make audit-release-package`
- generated source archive protected-path scan: `0`
- `RUN_LOCAL_APP=true RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true make release-candidate-check`
- `make agency-app-down`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make smoke`
- `git diff --check`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path git status check
- `gh auth status`
- `gh repo view ptse8204/open-transit-rt --json nameWithOwner,visibility,viewerPermission`
- remote tag absence check for `v0.1.0-rc.1`
- GitHub Release absence check for `v0.1.0-rc.1`

The release notes patch is docs-only. Final package, archive, claim, and
publication preflight checks are rerun from this checkpoint commit before any
tag push or GitHub Release creation.

Blocked checks:
No publication gate is currently blocked. Tag push and GitHub Release creation
are pending the final post-commit gate rerun.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.
Source archive protected-path scans returned zero hits after the export-ignore
policy was committed.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. Release notes now describe a
release candidate for local/self-hosted evaluation only and avoid stable,
production, compliance, consumer, hosted-service, vendor, SLA, and ETA-quality
claims.

Security/auth status:
GitHub tooling is available: `gh auth status` shows active account
`ptse8204`; repository view shows `ptse8204/open-transit-rt` is public and
viewer permission is `ADMIN`. No application auth, CSRF, credential, token,
private payload, or route behavior changed.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
was added.

Release/publication status:
No tag, tag push, GitHub Release, asset upload, or public publication action
has been taken yet in this checkpoint. The next step is final post-commit gate
rerun and, if still passed, public rc1 publication.

Install confidence status:
Phase 113 remains the current local install-confidence result.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 115 does not touch UX.

Master review:
Approved. Gates passed before the final release-notes wording patch, and the
remaining action is a final post-commit package/preflight rerun before
publication.

Required edits:
Commit checkpoint 000003, rerun final package/publication gates from the
checkpoint commit, then publish only if they still pass.

Decision:
Proceed to checkpoint 000003 commit and final publication gate execution.

Next checkpoint:
Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review.

Goal status:
Active. Phase 115 published the authorized public `v0.1.0-rc.1` release
candidate after all release gates passed.

Sub-agents used or simulated:
Release/Supply-Chain sub-agent findings from checkpoint 000001 were applied.
QA, Claim-Boundary, Security/Auth, Data/Migration, Documentation / IA,
Install Confidence, Web Design Skill, GTFS-RT Domain, Planning, and
Implementation roles were simulated by the Master Agent for closeout.

Changed files:
`docs/release-status-v0.1.0-rc.1.md`;
`docs/phase-115-v0.1.0-rc.1-public-release-cut.md`;
`docs/handoffs/phase-115.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`.

Validation run:
Passed before publication from release commit
`497f99a97baff630af147c83a7e1249bb08e32da`:

- `make test-release-package`
- strict `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 ... make release-package`
- `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package`
- source archive protected-path scan: `0`
- `make check`
- `make validate`
- `make test`
- `make smoke`
- `docker compose -f deploy/docker-compose.yml config`
- `RUN_LOCAL_APP=true RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true make release-candidate-check`
- `make agency-app-down`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse and exact prepared-only assertion
- protected-path git status check
- `gh auth status`
- `gh repo view ptse8204/open-transit-rt --json nameWithOwner,visibility,viewerPermission`
- remote tag absence check before publication
- GitHub Release absence check before publication
- release metadata verification after publication

Release-candidate diagnostics were written to
`.cache/release-candidate-check/20260516T030728Z`. The helper exited `0` with
36 passed, 0 blocker, 0 `needs_review`, and 3 `not_checked`; the three
not-checked rows were `make validate`, `make test`, and `make smoke`, all run
separately and passed.

Blocked checks:
No Phase 115 publication gate remained blocked. Published release download
replay and GitHub-generated archive replay are Phase 116 scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.
The audited release package source archive contained zero protected-path
entries after the `export-ignore` policy.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. The release status artifact
records a public release candidate for local/self-hosted evaluation only and
explicitly withholds stable, production, compliance, adoption, consumer,
final-root, hosted-service, vendor, SLA, and ETA-quality claims.

Security/auth status:
GitHub tooling was authenticated as active account `ptse8204`; repository
`ptse8204/open-transit-rt` is public and viewer permission was `ADMIN`. No
application route auth, CSRF behavior, credential handling, private payload
handling, or operator command behavior changed.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
was added.

Release/publication status:
Published. The annotated tag `v0.1.0-rc.1` was pushed to origin and
dereferences to `497f99a97baff630af147c83a7e1249bb08e32da`. The GitHub
Release is a public prerelease, draft `false`, published at
`2026-05-16T03:09:40Z`:
`https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`.

Install confidence status:
Phase 113 remains the current bounded local fresh-clone and local
source-archive install-confidence result. Independent public download replay
is Phase 116 and Phase 117 scope.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 115 did not touch UX.

Master review:
Approved. Phase 115 closed with a real public `v0.1.0-rc.1` prerelease after
release, package, claim, protected-path, prepared-only consumer, and GitHub
publication gates passed.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 116 published
release verification and download replay.

Decision:
Phase 115 is complete.

Next checkpoint:
Phase 116 -- Checkpoint 000001: add published release verification and
download replay plan.
