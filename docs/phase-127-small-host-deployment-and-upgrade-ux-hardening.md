# Phase 127 -- Small-Host Deployment And Upgrade UX Hardening

## Goal

Harden small-host deployment and upgrade UX: preflight, off-host validation,
backup/restore guidance, resource constraints, and recovery.

Phase 127 is not a production operations, hosted-service, SLA/uptime,
production-readiness, compliance, consumer acceptance, public launch, vendor
compatibility, hardware certification, adoption, stable release-readiness,
evidence, or deployment-success proof phase.

## Current Repo Context

- Phase 104 added deployment-doctor resource, service dependency, proxy,
  Postgres-capacity, backup/restore, and upgrade/rollback checks.
- The private Maintenance Center already summarizes backup/restore,
  upgrade/rollback, support bundle, cadence, monitoring/export, and
  infrastructure rows.
- `scripts/deployment-doctor.sh`, `scripts/oci-reference-check.sh`, and
  `scripts/validate-public-feeds.sh` provide local/reference diagnostics and
  dry-run flows without creating retained evidence.
- The Phase 115 public rc1 prerelease remains published, and Phase 117 public
  fresh-clone install confidence remains passed.

## Scope

- Improve small-host deployment and upgrade UX using existing private
  Maintenance Center and deployment-doctor boundaries.
- Prefer read-only summary rows, checklist clarity, and operator/technical
  helper separation over browser-executed operations.
- Keep destructive actions, live backup/restore, service control, migration
  execution, package publication, retained evidence, and consumer status
  changes outside the phase.
- Preserve private path, credential, token, backup dump, raw log, and raw env
  non-disclosure.

## Web Design Skill Use

The Web Design Skill at `~/.agents/skills/web-design-engineer` was loaded
because Phase 127 may touch private Maintenance Center UX.

Design decisions:

- Color palette: reuse the existing Operations Console neutral/status system.
- Typography: keep the existing server-rendered admin console type scale.
- Spacing system: keep compact operational spacing for scan-heavy pages.
- Border-radius strategy: reuse the existing low-radius panel/table style.
- Shadow hierarchy: no new decorative elevation.
- Motion style: no new animation; preserve existing focus/hover behavior.

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

- Small-host deployment and upgrade UX hardening or audit.
- Focused tests for any changed Maintenance Center or deployment-doctor
  behavior.
- `docs/handoffs/phase-127.md`
- Source-of-truth status updates for Phase 127 closeout.

## Implementation Plan

1. Add this Phase 127 plan and commit checkpoint 000001.
2. Inspect Maintenance Center, deployment-doctor, OCI/reference, off-host
   validation, backup/restore, and upgrade/rollback surfaces.
3. Implement a narrow UX hardening patch that improves preflight/recovery
   review without adding browser execution or stronger claims.
4. Run focused script/Go tests plus baseline validation; patch repo-caused
   failures.
5. Close Phase 127 with handoff/status docs and continue immediately to Phase
   128.

## Checkpoint Plan

- `Phase 127 -- Checkpoint 000001: add small host deployment and upgrade ux hardening plan`
- `Phase 127 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 127 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux hardening review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 127 -- Checkpoint 000001: add small host deployment and upgrade ux
hardening plan.

Goal status:
Active. Phase 126 is closed and Phase 127 has started.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, UI/UX, Web Design Skill, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Release, Install Confidence,
Connector, and GTFS-RT Domain roles are simulated by the Master Agent.

Changed files:
`docs/phase-127-small-host-deployment-and-upgrade-ux-hardening.md`.

Validation run:
Initial inspection reviewed the Phase 127 prompt, Phase 104 handoff,
Maintenance Center model/tests, deployment-doctor script, and validation and
claim-boundary instructions. The Web Design Skill was loaded for the planned
UX-facing work.

Blocked checks:
Implementation, focused tests, script checks, and closeout validation are
scheduled for later Phase 127 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release readiness, production readiness,
deployment success, compliance, adoption, agency approval, consumer
acceptance, consumer ingestion/listing/display, final-root readiness, hosted
service availability, paid support, SLA/uptime, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality, and
real-world ETA accuracy claims.

Security/auth status:
The plan does not add public routes, browser-executed service actions,
credential collection, token handling, raw env rendering, raw backup rendering,
live backup/restore, migration execution, service control, or package
publication.

Data/migration status:
No migration, durable state, dependency, or Go module change is planned.

Release/publication status:
The public rc1 prerelease remains published. Phase 127 does not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
The Web Design Skill was loaded before UX-facing Maintenance Center work. The
phase will preserve the existing dense private admin console visual system.

Master review:
Approved. The plan scopes Phase 127 to private small-host deployment and
upgrade UX hardening with no destructive browser operations and no stronger
public claims.

Required edits:
Commit checkpoint 000001, then implement the scoped UX hardening or audit.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 127 -- Checkpoint 000002: implement or audit primary scoped work.
