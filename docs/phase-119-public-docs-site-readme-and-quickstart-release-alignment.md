# Phase 119 -- Public Docs Site README And Quickstart Release Alignment

## Goal

Align README, docs, wiki, quickstart, release notes, and contributor entry
points around the published public `v0.1.0-rc.1` release candidate and the
verified install paths.

Phase 119 is documentation alignment only. It does not claim stable release
readiness, production readiness, compliance, adoption, consumer acceptance,
final-root readiness, hosted service availability, SLA/uptime, vendor
compatibility, hardware certification, or production-grade ETA quality.

## Current Repo Context

- Phase 115 published the public `v0.1.0-rc.1` GitHub prerelease.
- Phase 116 verified public release downloads and recorded the published
  source-archive `make check` blocker.
- Phase 117 verified public fresh-clone install confidence from the rc1 tag.
- Phase 118 completed post-release Web Design Skill UX validation.
- README/docs/wiki still contain pre-publication language such as
  "execute the authorized public release" and "next recommended release
  milestone is `v0.1.0-rc.1`."

## Scope

- Update top-level README release status and quickstart path.
- Update docs home release-candidate section.
- Update wiki home and small-agency quickstart around the public tag and
  fresh-clone path.
- Update release-candidate readiness and release notes language where stale
  pre-publication blockers remain.
- Preserve clear boundaries for source archive replay, stable release,
  production readiness, compliance, adoption, consumer acceptance, hosted
  service, vendor, SLA, and ETA-quality claims.

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

- Public docs and quickstart copy aligned to the actual published rc1 status.
- `docs/handoffs/phase-119.md`
- Source-of-truth status updates for Phase 119 closeout.

## Implementation Plan

1. Add this Phase 119 plan and commit checkpoint 000001.
2. Patch public-facing docs and quickstart entry points to point at the public
   rc1 release and verified fresh-clone path.
3. Run docs/claim/validation checks and patch any stale release wording.
4. Close Phase 119 with handoff/status docs and continue immediately to
   Phase 120.

## Checkpoint Plan

- `Phase 119 -- Checkpoint 000001: add public docs site readme and quickstart release alignment plan`
- `Phase 119 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 119 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 119 -- Checkpoint 000004: close public docs site readme and quickstart release alignment review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 119 -- Checkpoint 000001: add public docs site README and quickstart
release alignment plan.

Goal status:
Active. Phase 118 is closed and Phase 119 has started.

Sub-agents used or simulated:
Documentation / IA, Claim-Boundary, Release, Install Confidence, QA,
Security/Auth, Data/Migration, Web Design Skill, and GTFS-RT Domain roles are
simulated by the Master Agent.

Changed files:
`docs/phase-119-public-docs-site-readme-and-quickstart-release-alignment.md`.

Validation run:
Initial Phase 119 inspection reviewed README, docs home, wiki home,
small-agency quickstart, release-candidate readiness, release notes, and the
Phase 115-118 release/install/UX artifacts. Focused checkpoint validation is
scheduled before commit.

Blocked checks:
Public docs copy patching, stale release wording audit, and full validation
are scheduled for later Phase 119 checkpoints.

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
The plan is documentation-only and does not change route auth, CSRF,
credential handling, token handling, public exposure, private payload
handling, or operator command behavior.

Data/migration status:
No migration, schema, durable state, public feed contract, dependency, or Go
module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 119 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence passed and is the install path
to highlight.

Web design skill status:
Phase 118 Web Design Skill artifact is complete. Phase 119 does not touch UX.

Master review:
Approved. The plan reconciles stale public docs with the actual rc1 outcome
while preserving all claim and evidence boundaries.

Required edits:
Run checkpoint 000001 validation and commit, then patch public docs.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 119 -- Checkpoint 000002: implement or audit primary scoped work.
