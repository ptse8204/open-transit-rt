# Phase 115 Handoff -- v0.1.0-rc.1 Public Release Cut

## Status

Phase 115 is complete for the authorized public `v0.1.0-rc.1` release-candidate
cut.

The public GitHub prerelease was created after release package, validation,
claim-boundary, protected-path, prepared-only consumer, and GitHub tooling
gates passed:

`https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`

This is a public release candidate for local/self-hosted evaluation. It is not
a stable release and does not prove production readiness, compliance, adoption,
consumer acceptance, final-root readiness, hosted service availability,
SLA/uptime, vendor compatibility, hardware certification, or production-grade
ETA quality.

## Completed Checkpoints

- Phase 115 -- Checkpoint 000001: add v0.1.0-rc.1 public release cut plan.
- Phase 115 -- Checkpoint 000002: implement or audit primary scoped work.
- Phase 115 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review.

## Release Result

- GitHub Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`
- GitHub Release draft: `false`
- GitHub Release prerelease: `true`
- Published at: `2026-05-16T03:09:40Z`
- Annotated tag object: `52e91c7966e0fe1a5a4202277ab32173f8e78e67`
- Tag target commit: `497f99a97baff630af147c83a7e1249bb08e32da`
- Local package directory: `.cache/release-package/v0.1.0-rc.1`
- Local package status: `release_ready`
- Local source archive SHA-256:
  `dedf67537b1ed90c24921db32f0df7770aa42968c2d7cbe4927ec9a49a110e6f`
- Local source archive protected-path hits: `0`

The release status artifact is
`docs/release-status-v0.1.0-rc.1.md`.

## Changed Files

- `.gitattributes`
- `scripts/test-release-package.sh`
- `docs/decisions.md`
- `docs/dependencies.md`
- `docs/release-notes-v0.1.0-rc.1-draft.md`
- `docs/release-status-v0.1.0-rc.1.md`
- `docs/phase-115-v0.1.0-rc.1-public-release-cut.md`
- `docs/handoffs/phase-115.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- `make test-release-package`
- strict `make release-package` for `v0.1.0-rc.1`
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
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `gh auth status`
- `gh repo view ptse8204/open-transit-rt --json nameWithOwner,visibility,viewerPermission`
- remote tag absence check before publication
- GitHub Release absence check before publication
- release metadata verification after publication

Blocked:

- None for Phase 115 publication.
- Published release download replay and GitHub-generated source archive replay
  remain Phase 116 scope.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched.
The generated local release package source archive contained zero protected
entries after the `.gitattributes` `export-ignore` policy.

## Consumer Tracker Status

`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven targets remain present in order and all remain `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim Boundary Status

Phase 115 makes only this release claim: public `v0.1.0-rc.1` release
candidate for local/self-hosted evaluation.

It makes no stable release readiness, production readiness, compliance,
adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy claim.

## Security/Auth Status

GitHub release tooling was authenticated as active account `ptse8204`, and
repository `ptse8204/open-transit-rt` was public with viewer permission
`ADMIN`.

No application route auth, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or operator
command behavior changed.

## Data/Migration Status

No migration, schema, durable state, public feed contract, dependency, runtime
dependency, or Go module change was added.

## Release/Publication Status

Published as a GitHub prerelease. The release is public and draft `false`.

Phase 116 must verify public download replay, including uploaded release
assets and GitHub-generated source archives.

## Install Confidence Status

Phase 113 remains the current bounded local fresh-clone and local source
archive install-confidence result. Public download/archive install confidence
is Phase 116 and Phase 117 scope.

## Web Design Skill Status

Phase 114 Web Design Skill artifact remains complete. Phase 118 remains the
post-release Web Design Skill UX validation pass.

## Commit List

- `b51e22f` -- Phase 115 -- Checkpoint 000001: add v0.1.0-rc.1 public release
  cut plan
- `10a36dc` -- Phase 115 -- Checkpoint 000002: implement or audit primary
  scoped work
- `497f99a` -- Phase 115 -- Checkpoint 000003: run validation and patch
  required gaps
- Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review

## Checkpoint Report

Checkpoint:
Phase 115 -- Checkpoint 000004: close v0.1.0-rc.1 public release cut review.

Goal status:
Active. Phase 115 is closed and the goal continues to Phase 116.

Sub-agents used or simulated:
Release/Supply-Chain sub-agent findings were incorporated. QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Install
Confidence, Web Design Skill, GTFS-RT Domain, Planning, and Implementation
closeout roles were simulated by the Master Agent.

Changed files:
`docs/handoffs/phase-115.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-115-v0.1.0-rc.1-public-release-cut.md`;
`docs/release-status-v0.1.0-rc.1.md`.

Validation run:
Phase 115 release gates passed before publication. After closeout docs were
updated, focused docs/protected-path validation is rerun before the checkpoint
000004 commit.

Blocked checks:
No Phase 115 publication gate is blocked. Phase 116 must perform download
replay and GitHub-generated archive verification.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 115 records a public prerelease only. It makes no stronger production,
compliance, adoption, consumer, hosted-service, vendor, SLA, or ETA-quality
claim.

Security/auth status:
GitHub release tooling was available and authorized. No application security
behavior changed.

Data/migration status:
No migration, schema, durable state, public feed contract, dependency, or Go
module change was added.

Release/publication status:
Public prerelease published.

Install confidence status:
Phase 113 local install-confidence remains current; public download replay is
next.

Web design skill status:
Phase 114 Web Design Skill artifact is complete; Phase 118 remains scheduled.

Master review:
Approved. Phase 115 closes with a real public rc1 prerelease and bounded
release status evidence.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 116.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 116.

Next checkpoint:
Phase 116 -- Checkpoint 000001: add published release verification and
download replay plan.
