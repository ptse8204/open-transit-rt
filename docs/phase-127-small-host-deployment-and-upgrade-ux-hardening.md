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

## Checkpoint Report -- 000002

Checkpoint:
Phase 127 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Phase 127 implemented the scoped private Maintenance Center
small-host readiness UX hardening.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, UI/UX, Web Design Skill, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Release, Install Confidence,
Connector, and GTFS-RT Domain roles are simulated by the Master Agent.

Changed files:
`cmd/agency-config/operations_maintenance.go`,
`cmd/agency-config/operations.go`, `cmd/agency-config/main_test.go`,
`docs/tutorials/small-agency-maintenance-guide.md`, and this phase report.

Implementation summary:
Added a read-only `small_host_readiness` Maintenance Center panel and rendered
it in the private Maintenance Center HTML. The panel gives a compact
preflight/recovery checklist for deployment-doctor diagnostics, dry-run feed
checks, off-host validation choices, resource-budget review, backup/restore
recovery path, and upgrade stop points. It does not add browser-executed
deployment actions, backup/restore, migration execution, service control,
package publication, release actions, evidence writes, or consumer status
changes. The small-agency maintenance guide now points operators to the new
panel before upgrade or recovery work.

Validation run:
`gofmt` passed on touched Go files. `go test ./cmd/agency-config -run
'TestOperationsMaintenance|TestOperationsNavigation|TestOperationsRouteTitles|TestOperationsConsole'`
passed. `git diff --check` passed. `scripts/check-consumer-tracker.sh`
passed. Protected-path git status returned no output.

Blocked checks:
None for this checkpoint. Full repo validation is scheduled for checkpoint
000003.

Protected path status:
`git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output. No
protected evidence path, migration, or module file was modified.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
The new panel explicitly avoids deployment success, production readiness,
hosted availability, SLA/uptime, release readiness, compliance, final-root
readiness, consumer acceptance, vendor compatibility, hardware certification,
and ETA-quality claims.

Security/auth status:
The Maintenance Center remains private and read-only. No route auth, CSRF
behavior, credential handling, token handling, public exposure, raw env
rendering, raw backup rendering, private payload handling, or browser command
execution behavior was changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 127 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
The Web Design Skill was loaded before the UX-facing Maintenance Center
change. The patch preserves the dense private admin console style, low-radius
table/panel vocabulary, and claim-boundary copy pattern.

Master review:
Approved for full validation. The implementation adds explicit small-host UX
guidance without expanding browser execution or public claims.

Required edits:
Run checkpoint 000003 full validation and patch any repo-caused failures.

Decision:
Proceed to checkpoint 000002 commit.

Next checkpoint:
Phase 127 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 127 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. Full Phase 127 validation passed with no repo-caused failures.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, UI/UX, Web Design Skill, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Release, Install Confidence,
Connector, and GTFS-RT Domain roles are simulated by the Master Agent.

Changed files:
`docs/phase-127-small-host-deployment-and-upgrade-ux-hardening.md`.

Validation run:
`git status --short` returned clean at validation start. `git diff --check`
passed. `python3 -m json.tool
docs/evidence/consumer-submissions/status.json` passed.
`scripts/check-consumer-tracker.sh` passed. Protected-path git status returned
no output. `sh -n scripts/deployment-doctor.sh scripts/oci-reference-check.sh
scripts/validate-public-feeds.sh scripts/pilot-ops.sh` passed. `go test
./cmd/agency-config -run
'TestOperationsMaintenance|TestOperationsNavigation|TestOperationsRouteTitles|TestOperationsConsole'`
passed. `OUTPUT_DIR=.cache/phase-127/deployment-doctor FORCE=true
scripts/deployment-doctor.sh` passed and reported expected local configuration
blockers while all claim flags remained false. `PUBLIC_BASE_URL=https://feeds.example.org
OUTPUT_DIR=.cache/phase-127/validate-public-feeds FORCE=true
scripts/validate-public-feeds.sh --dry-run` passed. `PUBLIC_BASE_URL=https://feeds.example.org
OUTPUT_DIR=.cache/phase-127/oci-reference-check FORCE=true
scripts/oci-reference-check.sh --dry-run` passed. Generated Phase 127 helper
JSON parsed successfully. `make check` passed. `make validate` passed. `make
test` passed. `docker compose -f deploy/docker-compose.yml config` passed.
`make audit-product-acceptance` passed. `make audit-final-claim-review`
passed. `make external-connection-check` passed. `make adapter-conformance`
passed. `make gtfsrt-conformance` passed. Final `git status --short`, `git
diff --check`, and protected-path git status returned clean.

Blocked checks:
None.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
`scripts/check-consumer-tracker.sh` reported exactly seven prepared-only
targets.

Claim-boundary status:
Product acceptance and final claim audits passed. The deployment-doctor dry
run reported local backup/restore and environment blockers as private
configuration blockers, not as production, uptime, SLA, compliance, deployment
success, or release-readiness proof.

Security/auth status:
No route auth, CSRF behavior, credential handling, token handling, public
exposure, raw env rendering, raw backup rendering, private payload handling, or
browser command execution behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, or Go module change was made.

Release/publication status:
The public rc1 prerelease remains published. Phase 127 did not create or
modify a release.

Install confidence status:
Phase 117 public fresh-clone install confidence remains passed.

Web design skill status:
The Web Design Skill was used for the Phase 127 private Maintenance Center UX
change. Validation found no follow-up visual/code patch required.

Master review:
Approved for Phase 127 closeout.

Required edits:
None.

Decision:
Proceed to checkpoint 000003 commit.

Next checkpoint:
Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux
hardening review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux
hardening review.

Goal status:
Active. Phase 127 is closed and the goal continues to Phase 128.

Sub-agents used or simulated:
The agent thread limit prevents new real sub-agents. Context / Repo Truth,
Planning, Implementation, QA, UI/UX, Web Design Skill, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, Release, Install Confidence,
Connector, and GTFS-RT Domain roles are simulated by the Master Agent.

Closeout summary:
Phase 127 added a private read-only `small_host_readiness` Maintenance Center
panel and aligned the small-agency maintenance guide. It improves small-host
preflight and recovery review without adding browser-executed deployment,
backup, restore, migration, service-control, package, rollback, release,
evidence, or consumer-status actions.

Changed files:
`docs/handoffs/phase-127.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`; and this phase
report.

Validation run:
Full Phase 127 validation passed before closeout docs. Focused closeout
validation passed after closeout docs: `git diff --check`, `make check`,
`make audit-product-acceptance`, `make audit-final-claim-review`,
`scripts/check-consumer-tracker.sh`, and protected-path git status.

Blocked checks:
No Phase 127 check remains blocked.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 127 remains bounded to private small-host deployment and upgrade UX
guidance and makes no stable release readiness, deployment success, production
readiness, compliance, adoption, agency approval, consumer acceptance,
consumer ingestion/listing/display, final-root readiness, hosted service
availability, paid support, SLA/uptime, vendor compatibility, hardware
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
The Web Design Skill was used for the private Maintenance Center UX change.

Master review:
Approved. Phase 127 closes with tested small-host readiness UX.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 128.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 128.

Next checkpoint:
Phase 128 -- Checkpoint 000001: add contributor and agency evaluator adoption
kit plan.
