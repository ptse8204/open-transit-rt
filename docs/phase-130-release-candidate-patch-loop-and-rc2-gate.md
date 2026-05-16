# Phase 130 -- Release Candidate Patch Loop And rc2 Gate

## Goal

Run a release candidate patch loop and decide whether an rc2 gate is needed;
prepare rc2 only if repo-caused blockers require it.

Phase 130 is not a stable release, rc2 publication authorization, production
readiness, compliance, adoption, consumer acceptance, hosted-service,
SLA/uptime, vendor compatibility, hardware certification, final-root
readiness, evidence, or ETA-quality proof phase.

## Current Repo Context

- Phase 115 published public prerelease `v0.1.0-rc.1`.
- Phase 116 verified release downloads and recorded that the already published
  rc1 source archives still fail `make check` because protected consumer
  tracker state is intentionally export-ignored while the rc1 tag still
  required it.
- Phase 117 showed a public fresh clone of the rc1 tag passes after the
  install-confidence harness installs validators before validate-enabled
  trials.
- Current source has post-rc1 patches and adoption/support hardening through
  Phase 129.

## Scope

- Audit current-source release-candidate health and known rc1 blockers.
- Generate and audit local rc2-style package artifacts only if useful for
  deciding the gate.
- Decide whether rc2 is required, recommended, or not needed, with exact
  evidence and without publishing rc2 unless separately authorized.
- Preserve the public rc1 release status and avoid changing tags, GitHub
  releases, consumer statuses, or protected evidence paths.

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

- Release-candidate patch-loop and rc2-gate artifact.
- Local package or release-candidate diagnostic evidence when safe and useful.
- `docs/handoffs/phase-130.md`
- Source-of-truth status updates for Phase 130 closeout.

## Implementation Plan

1. Add this Phase 130 plan and commit checkpoint 000001.
2. Inspect rc1 release, download replay, install-confidence, release process,
   package helpers, and current source status.
3. Run local release-candidate patch-loop diagnostics and record whether rc2
   is needed, recommended, or blocked.
4. Run release/package/claim/protected-path validation; patch repo-caused
   failures.
5. Close Phase 130 with handoff/status docs and continue immediately to Phase
   131.

## Checkpoint Plan

- `Phase 130 -- Checkpoint 000001: add release candidate patch loop and rc2 gate plan`
- `Phase 130 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 130 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 130 -- Checkpoint 000004: close release candidate patch loop and rc2 gate review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 130 -- Checkpoint 000001: add release candidate patch loop and rc2 gate
plan.

Goal status:
Active. Phase 129 is closed and Phase 130 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, Release/Supply-Chain, Install Confidence,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Connector,
GTFS-RT Domain, and UI/UX roles are simulated by the Master Agent.

Changed files:
`docs/phase-130-release-candidate-patch-loop-and-rc2-gate.md`.

Validation run:
Initial inspection reviewed the Phase 130 prompt, rc1 release status,
download replay, public install confidence, release package helpers, and
current working tree status.

Blocked checks:
Release-candidate diagnostics, package audit, rc2-gate decision, and closeout
validation are scheduled for later Phase 130 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, rc2 publication claims,
production readiness, compliance, adoption, agency approval, consumer
acceptance, consumer ingestion/listing/display, final-root readiness, hosted
service availability, paid support, SLA/uptime, vendor compatibility, hardware
certification, and production-grade ETA claims.

Security/auth status:
The plan does not add public routes, credential handling, token handling,
private payload handling, evidence collection, tag publication, GitHub release
publication, or service actions.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 130 does not itself
authorize rc2 publication.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
Not used for checkpoint 000001 because Phase 130 is release/process work and
does not touch a visual UI surface.

Master review:
Approved. The plan scopes Phase 130 to local rc2-gate decision work without
publishing or overclaiming.

Required edits:
Commit checkpoint 000001, then run the scoped release-candidate patch-loop
audit.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 130 -- Checkpoint 000002: implement or audit primary scoped work.
