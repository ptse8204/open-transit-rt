# Phase 104 -- Small-Host Deployment And Upgrade Hardening

## Goal

Make small-server deployment, validator offload, backup/restore drills,
upgrade/rollback review, service dependency checks, Caddy/proxy exposure, and
Postgres/PostGIS diagnostics clearer and safer for private operators without
claiming production readiness.

## Current Surface

- `scripts/deployment-doctor.sh` is the main private deployment diagnostic. It
  already checks env presence, generated-secret status, public feed fetches,
  private route boundaries, optional authenticated admin readiness, HTTPS
  posture, service health, database/migration/PostGIS status, pinned validator
  tooling, backup/restore readiness, git identity, recent diagnostics, and the
  prepared-only consumer tracker guard.
- `scripts/validate-public-feeds.sh` and
  `docs/deployment/off-host-validation.md` already describe off-host public feed
  validation when tiny servers should not run validator workloads.
- `scripts/oci-reference-check.sh`, `deploy/oci/*`, `deploy/Caddyfile.local`,
  and the OCI/reference deployment docs cover reference deployment and proxy
  boundaries.
- `scripts/pilot-ops.sh`, `docs/runbooks/backup-and-restore.md`, and
  `docs/upgrade-and-rollback.md` cover backup, restore, and upgrade guidance.
- `/admin/operations/maintenance` already exposes safe deployment-doctor and
  reliability summary pointers, but it must remain a private GET-only review
  surface and must not execute deployment actions.

## Master-Approved Plan

1. Add the Phase 104 plan and checkpoint report.
2. Implement the smallest safe deployment hardening over existing seams:
   - extend `scripts/deployment-doctor.sh` with read-only small-host resource,
     service dependency, proxy exposure, Postgres/PostGIS, validator off-host,
     and backup/restore/upgrade guidance fields;
   - patch safe route/proxy drift and environment guidance where current repo
     truth shows mismatches;
   - update docs so small-host operators know which checks are local private
     diagnostics and when to move validators off host;
   - optionally surface bounded new deployment-doctor summary fields in the
     private Maintenance Center without adding execution controls.
3. Do not add migrations, live backup/restore execution, service timers, hosted
   monitoring, external contact, real credentials, evidence collection, public
   routes, or consumer automation.
4. Add focused tests around script syntax, generated summary shape, route/proxy
   guards, no secret leakage, and claim-boundary flags.
5. Run focused checks, then required baseline/code-change validation.
6. Close with `docs/handoffs/phase-104.md`, `docs/handoffs/latest.md`, and
   roadmap/status updates.

## Non-Goals

- No production-readiness, hosted-service, SLA, uptime, compliance,
  consumer-acceptance, agency-adoption, vendor, hardware, final-root,
  release-readiness, or production-grade ETA claim.
- No retained evidence or protected-path writes.
- No consumer tracker edits or status movement.
- No real credentials, real private payloads, real backup content, or real
  restore execution.
- No tag, GitHub Release, image push, package publication, or public launch.
- No automatic migration, rollback, restore, or validation execution from the
  private browser UI.

## Checkpoint Plan

- `Phase 104 -- Checkpoint 000001: add small-host deployment and upgrade hardening plan`
- `Phase 104 -- Checkpoint 000002: implement primary scoped work`
- `Phase 104 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 104 -- Checkpoint 000004: close small-host deployment and upgrade hardening review`

## Focused Validation Targets

- `sh -n scripts/deployment-doctor.sh scripts/oci-reference-check.sh scripts/validate-public-feeds.sh scripts/pilot-ops.sh`
- `OUTPUT_DIR=.cache/phase-104/deployment-doctor FORCE=true scripts/deployment-doctor.sh`
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-104/validate-public-feeds FORCE=true scripts/validate-public-feeds.sh --dry-run`
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-104/oci-reference-check FORCE=true scripts/oci-reference-check.sh --dry-run`
- JSON validation for generated `.cache/phase-104/**/summary.json` and
  `manifest.json` files.
- `go test ./cmd/agency-config -run 'DeploymentDoctor|ReferenceCheck|Maintenance|OperationsNavigation|RouteTitles'`

Because this phase is expected to change code/docs/tests/scripts, closeout also
requires:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 104 -- Checkpoint 000001: add small-host deployment and upgrade
hardening plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
are simulated by the Master Agent for this plan checkpoint. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-104-small-host-deployment-and-upgrade-hardening.md`.

Validation run:
Initial Phase 104 repository inspection found existing deployment-doctor,
off-host validation, backup/restore, upgrade/rollback, proxy, Postgres/PostGIS,
and Maintenance Center diagnostic seams. Sub-agents completed read-only
inspection and planning. After adding the plan, `git status --short` showed
only `docs/phase-104-small-host-deployment-and-upgrade-hardening.md`; `git
diff --check` passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused script/UI tests, and closeout baseline checks are not
yet run because this checkpoint only approves the Phase 104 plan.
Release-candidate checks, package generation/audit, evidence collection,
publication, and live deployment actions are out of scope for Phase 104.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids production-readiness, hosted-service, SLA, uptime,
compliance, consumer-acceptance, agency-adoption, vendor, hardware, final-root,
release-readiness, public-launch, and production-grade ETA claims.

Security/auth status:
The plan preserves private diagnostics, avoids raw env/secret rendering, avoids
credential collection, and does not add browser execution of backup, restore,
migration, validation, or service-control actions.

Data/migration status:
No migration, durable deployment state, backup metadata table, service-control
table, or schema change is planned.

Master review:
Approved. The smallest safe Phase 104 implementation is to harden private
read-only diagnostics and guidance over existing deployment seams.

Required edits:
Extend bounded diagnostics and documentation, patch confirmed route/proxy/env
drift, add focused tests, and record validation results.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 104 -- Checkpoint 000002: implement primary scoped work.
