# Phase 113 -- Fresh Clone Install Harness And Release Dry Run

## Goal

Build and run a repeatable install-confidence harness so a non-maintainer path
can be tested from a fresh clone and/or release archive before publication.

Phase 113 must record exact local blockers and must not convert local install
signals into production readiness, compliance, hosted-service, consumer,
adoption, vendor, hardware, SLA/uptime, or ETA-quality claims.

## Current Repo Context

- Phase 112 closed with publication status
  `blocked_public_distribution_review`.
- Install-confidence audit found that GitHub source archives do not include
  `.git`, while the current `make check` assumes a Git worktree because it
  runs `git diff --check`.
- Fresh-clone install confidence can run from a clone outside the active
  checkout without touching protected evidence paths.
- Raw install logs and downloaded feed artifacts must stay under ignored
  `.cache` or temporary directories.

## Scope

- Add a local install-confidence harness that can:
  - clone a source repository into ignored `.cache`;
  - checkout a requested ref;
  - record bounded environment metadata;
  - run `make check`;
  - run `scripts/bootstrap-dev.sh --check`;
  - optionally run a local app startup and five-feed fetch;
  - support archive extraction/replay for later Phase 116/117 use.
- Make lightweight checks archive-aware where needed.
- Add focused tests for the harness boundaries.
- Run the harness from outside the active checkout using the local repository
  as the clone source, because Phase 111/112 commits are not pushed yet.
- Record a non-protected install-confidence summary artifact.

## Boundaries

- Do not tag, push, create a GitHub Release, upload assets, push images, create
  retained evidence, contact external parties, or move consumer statuses.
- Do not modify or generate files under:
  - `docs/evidence/captured/**`
  - `docs/evidence/consumer-submissions/status.json`
  - `docs/evidence/consumer-submissions/current/**`
  - `docs/evidence/consumer-submissions/artifacts/**`
  - `docs/evidence/consumer-submissions/packets/**`
- Do not commit raw install logs, downloaded feed artifacts, credentials,
  private paths, raw external payloads, or `.cache` outputs.
- Do not claim stable release readiness, production readiness, compliance,
  adoption, agency approval, consumer acceptance, consumer
  ingestion/listing/display, final-root readiness, hosted service
  availability, paid support, SLA/uptime, vendor compatibility, hardware
  certification, production AVL reliability, production-grade ETA quality, or
  real-world ETA accuracy.

## Deliverables

- `scripts/install-confidence.sh`
- `scripts/test-install-confidence.sh`
- Make targets for install-confidence checks
- `docs/install-confidence-v0.1.0-rc.1.md`
- `docs/handoffs/phase-113.md`
- Source-of-truth status updates

## Implementation Plan

1. Add this Phase 113 plan and commit checkpoint 000001.
2. Add the install-confidence harness, archive-aware lightweight check
   behavior, tests, Make targets, and a first install-confidence summary.
3. Run focused and baseline validation; patch required gaps.
4. Close Phase 113 with a handoff and continue immediately to Phase 114.

## Checkpoint Plan

- `Phase 113 -- Checkpoint 000001: add fresh clone install harness and release dry run plan`
- `Phase 113 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 113 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 113 -- Checkpoint 000004: close fresh clone install harness and release dry run review`

## Focused Validation Targets

- `git status --short`
- `git diff --check`
- `make check`
- `make test-install-confidence`
- `make install-confidence` with local repository source and startup disabled
- archive-mode harness dry run against a generated local source archive
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Because Phase 113 changes scripts and Makefile behavior, also run:

- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 113 -- Checkpoint 000001: add fresh clone install harness and release
dry run plan.

Goal status:
Active. Phase 112 is closed and Phase 113 has started.

Sub-agents used or simulated:
Install Confidence sub-agent findings from Phase 111 were incorporated.
Planning, Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Release/Supply-Chain, and Web Design Skill
roles are simulated by the Master Agent for this plan checkpoint.

Changed files:
`docs/phase-113-fresh-clone-install-harness-and-release-dry-run.md`.

Validation run:
Initial Phase 113 inspection reviewed the Phase 113 prompt, independent
install-confidence audit prompt, `scripts/bootstrap-dev.sh`,
`scripts/check-validators.sh`, `Makefile`, current git status, and Phase 112
release status. Focused checkpoint validation is scheduled before commit.

Blocked checks:
Implementation, harness tests, fresh-clone dry run, archive dry run, and full
script validation are scheduled for later Phase 113 checkpoints. Release
publication, retained evidence, external contact, consumer status movement,
protected path writes, and stronger claims are out of Phase 113 scope.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
compliance, adoption, consumer acceptance, final-root, hosted-service,
paid-support, SLA/uptime, vendor, hardware, production AVL, and ETA-quality
claims.

Security/auth status:
The plan records only bounded local install metadata and does not change route
auth, CSRF, credentials, token handling, public exposure, private payload
handling, external contact, or command execution behavior.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change is planned.

Release/publication status:
No release action is planned for Phase 113. Publication remains blocked by
Phase 112 source archive public-distribution review.

Install confidence status:
Harness implementation and dry runs are planned for checkpoint 000002.

Web design skill status:
No UX artifact is planned for Phase 113. Web Design Skill artifacts are
scheduled for Phases 114 and 118.

Master review:
Approved. The phase addresses a real install replay friction and keeps all
raw outputs in ignored local directories.

Required edits:
Run checkpoint 000001 validation and commit, then implement the harness and
first summary artifact.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 113 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 113 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Install-confidence harness implementation is complete and ready for
post-commit dry runs.

Sub-agents used or simulated:
Install Confidence sub-agent recommendations were incorporated. Implementation
and QA were performed by the Master Agent. Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Release/Supply-Chain, and Web
Design Skill roles were simulated or deferred according to Phase 113 scope.

Changed files:
`Makefile`; `scripts/install-confidence.sh`;
`scripts/test-install-confidence.sh`;
`docs/phase-113-fresh-clone-install-harness-and-release-dry-run.md`.

Validation run:
`make test-install-confidence` passed; `make check` passed; `git diff --check`
passed. The implementation checkpoint did not run full clone/archive dry runs
because the harness needs to be committed before a fresh clone of the current
repository can include the new script and archive-aware `make check` behavior.

Blocked checks:
Full install-confidence dry runs, `make validate`, `make test`, and compose
config are scheduled for checkpoint 000003. Release publication, retained
evidence, external contact, consumer status movement, protected path writes,
and stronger claims remain out of Phase 113 scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.
The harness enforces output under `.cache/install-confidence/**` and rejects
evidence-like output paths.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets must remain in order and all remain `prepared`.

Claim-boundary status:
The harness summary states that install-confidence diagnostics are local only
and not retained evidence, release publication, production readiness,
compliance proof, consumer acceptance, agency approval, hosted service
availability, vendor compatibility, SLA/uptime, or ETA-quality proof.

Security/auth status:
No route, auth behavior, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or command
execution behavior changed. The harness records bounded local environment
metadata and stores raw logs only under ignored `.cache`.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change was added.

Release/publication status:
No tag, GitHub Release, public package, asset upload, image push, or public
announcement was created. Publication remains blocked by the Phase 112 source
archive public-distribution review.

Install confidence status:
Harness and tests were added. Actual local fresh-clone/archive dry runs are
scheduled for checkpoint 000003.

Web design skill status:
No UX artifact was added in Phase 113.

Master review:
Approved. The implementation addresses archive replay friction by making
`make check` skip `git diff --check` only when outside a Git worktree, while
preserving the normal git diff check in repository clones.

Required edits:
Commit checkpoint 000002, then run fresh-clone and archive-mode install
confidence diagnostics and record the results.

Decision:
Proceed to checkpoint 000002 commit, then checkpoint 000003 validation.

Next checkpoint:
Phase 113 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 113 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Harness validation found and patched a shell portability bug before
fresh-clone closeout diagnostics.

Sub-agents used or simulated:
Install Confidence and QA were simulated by the Master Agent for this
validation checkpoint. Documentation / IA, Claim-Boundary, Security/Auth,
Data/Migration, Release/Supply-Chain, and Web Design Skill roles were
simulated or deferred according to Phase 113 scope.

Changed files:
`scripts/install-confidence.sh`;
`docs/phase-113-fresh-clone-install-harness-and-release-dry-run.md`.

Validation run:
An initial local fresh-clone run reached the harness summary writer after app
startup and feed fetches but failed on a shell `printf` portability issue for
Markdown bullet lines. The script was patched to use `printf --` for bullet
format strings. After the patch, `make agency-app-down` passed; `make
test-install-confidence` passed; `git diff --check` passed; `make check`
passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; and the
protected-path status check returned no output.

Blocked checks:
Full clone/archive dry runs are scheduled for checkpoint 000004 after this
fix is committed, so the cloned checkout includes the patched harness.
`make validate`, `make test`, and compose config remain scheduled before Phase
113 closeout because scripts/Makefile behavior changed. Release publication,
retained evidence, external contact, consumer status movement, protected path
writes, and stronger claims remain out of Phase 113 scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets remain in order and all remain `prepared`.

Claim-boundary status:
The harness wording remains local-diagnostic only and does not claim release
readiness, production readiness, compliance, adoption, consumer acceptance,
final-root readiness, hosted service availability, SLA/uptime, vendor
compatibility, hardware certification, or ETA quality.

Security/auth status:
No route, auth behavior, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or command
execution behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change was added.

Release/publication status:
No tag, GitHub Release, public package, asset upload, image push, or public
announcement was created. Publication remains blocked by the Phase 112 source
archive public-distribution review.

Install confidence status:
Harness tests pass. Full clone/archive diagnostics are deferred until after
this checkpoint commit so the clone contains the patched script.

Web design skill status:
No UX artifact was added in Phase 113.

Master review:
Approved. The bug was caught by using the harness in the intended mode and
patched before install-confidence results were recorded.

Required edits:
Commit checkpoint 000003, rerun fresh-clone and archive diagnostics, then
record results in the Phase 113 closeout.

Decision:
Proceed to checkpoint 000003 commit, then checkpoint 000004 closeout dry runs.

Next checkpoint:
Phase 113 -- Checkpoint 000004: close fresh clone install harness and release
dry run review.
