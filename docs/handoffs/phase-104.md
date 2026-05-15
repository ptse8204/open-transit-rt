# Phase 104 Handoff -- Small-Host Deployment And Upgrade Hardening

## Status

Phase 104 is complete for private small-host deployment and upgrade hardening.
The deployment doctor now records bounded resource posture, service dependency,
proxy exposure, Postgres pool budget, off-host validator guidance, backup/
restore readiness aliases, and upgrade/rollback checklist posture. The private
Maintenance Center exposes those new categories as read-only summary rows.

The work remains private, local-diagnostic, and non-evidentiary. It adds no
live deployment action, backup/restore execution, migration execution, service
control, public admin route, publication, retained evidence, consumer action,
or stronger public claim.

## Completed Checkpoints

- Phase 104 -- Checkpoint 000001: add small-host deployment and upgrade hardening plan.
- Phase 104 -- Checkpoint 000002: implement primary scoped work.
- Phase 104 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 104 -- Checkpoint 000004: close small-host deployment and upgrade hardening review.

## Product Result

- `scripts/deployment-doctor.sh` now writes private small-host resource,
  service dependency/proxy, Postgres capacity, backup/restore readiness, and
  upgrade/rollback checklist summaries.
- Deployment-doctor top-level categories now include `small_host_resources`,
  `service_dependencies`, `proxy_exposure`, `postgres_capacity`, and
  `upgrade_rollback`, with expanded all-false claim flags.
- `scripts/oci-reference-check.sh` now probes `/healthz` for loopback service
  health and recognizes both restore-drill and pilot-ops restore env names.
- `scripts/pilot-ops.sh` now rejects `EVIDENCE_OUTPUT_DIR` values under
  protected repo evidence paths before running operator evidence helpers.
- The private Maintenance Center infrastructure panel now surfaces the new
  deployment-doctor categories without executing commands.
- Small-host deployment, off-host validation, reference environment, upgrade/
  rollback, pilot-ops, and maintenance docs now describe the new boundaries and
  `DB_MAX_CONNS=3` guidance.

## Changed Files

- `scripts/deployment-doctor.sh`
- `scripts/oci-reference-check.sh`
- `scripts/pilot-ops.sh`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations_maintenance_summaries.go`
- `cmd/agency-config/main_test.go`
- `deploy/oci/postgresql-tuning.conf`
- `deploy/oci/setup-instance.sh`
- `docs/deployment/README.md`
- `docs/deployment/oci-reference-check.md`
- `docs/deployment/oci-reference-deployment.md`
- `docs/deployment/oci-reference-env.example`
- `docs/deployment/off-host-validation.md`
- `docs/deployment/reference-deployment-doctor.md`
- `docs/runbooks/small-agency-pilot-operations.md`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `docs/upgrade-and-rollback.md`
- `docs/phase-104-small-host-deployment-and-upgrade-hardening.md`
- `docs/handoffs/phase-104.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- `git status --short`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `sh -n scripts/deployment-doctor.sh scripts/oci-reference-check.sh scripts/validate-public-feeds.sh scripts/pilot-ops.sh`
- `OUTPUT_DIR=.cache/phase-104/deployment-doctor FORCE=true scripts/deployment-doctor.sh`
- JSON validation for generated deployment-doctor summary, manifest, small-host resource, service-dependency, Postgres-capacity, and upgrade/rollback summaries
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-104/validate-public-feeds FORCE=true scripts/validate-public-feeds.sh --dry-run`
- JSON validation for generated validate-public-feeds summary and manifest
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-104/oci-reference-check FORCE=true scripts/oci-reference-check.sh --dry-run`
- JSON validation for generated oci-reference-check summary and manifest
- custom assertions for Phase 104 deployment-doctor categories, all-false claim flags, off-host validator guidance, `recommended_db_max_conns=3`, and OCI `/healthz` probing
- protected-path negative guard for `scripts/pilot-ops.sh`
- `go test ./cmd/agency-config -run 'DeploymentDoctor|ReferenceCheck|Maintenance|OperationsNavigation|RouteTitles'`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- final `git status --short`
- final protected-path status check
- final `git diff --check`

Blocked:

- Release-candidate diagnostics, package generation/audit, retained evidence,
  real public-root validation, live deployment actions, live backup/restore,
  consumer submission, public publication, and tag/release/package/image
  publication were not run because they are outside Phase 104 scope.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check for `docs/evidence/consumer-submissions`,
`docs/evidence/captured`, `db/migrations`, `go.mod`, and `go.sum` returned no
output.

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

## Claim-Boundary Status

Phase 104 makes no production-readiness, hosted-service, SLA, uptime,
compliance, consumer-acceptance, agency-adoption, vendor, hardware, final-root,
release-readiness, public-launch, or production-grade ETA claim. The new
diagnostics are private preflight/checklist signals only.

## Security/Auth Status

No public route, credential collection, raw env rendering, raw backup rendering,
live send, service-control action, migration execution, browser validator
execution, or external contact was added.

## Data/Migration Status

No migration, schema change, durable deployment state, backup metadata table,
restore table, service-control table, public feed contract change, or go module
dependency change was added.

## Commit List

- `660c768` -- Phase 104 -- Checkpoint 000001: add small-host deployment and upgrade hardening plan
- `023fc97` -- Phase 104 -- Checkpoint 000002: implement primary scoped work
- `a639622` -- Phase 104 -- Checkpoint 000003: run validation and patch required gaps
- Phase 104 -- Checkpoint 000004: close small-host deployment and upgrade hardening review

## Checkpoint Report

Checkpoint:
Phase 104 -- Checkpoint 000004: close small-host deployment and upgrade
hardening review.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain
closeout roles were simulated by the Master Agent. Master Agent -- GPT-5.5
x-high, current thread.

Changed files:
`docs/handoffs/phase-104.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-104-small-host-deployment-and-upgrade-hardening.md`.

Validation run:
Closeout relies on the checkpoint 000003 full validation pass: focused script
and Go tests, generated summary JSON checks, baseline checks, product
acceptance audit, final claim audit, `make validate`, `make test`, docker
compose config, and final status/protected-path/diff checks all passed.

Blocked checks:
Release-candidate diagnostics, package generation/audit, retained evidence,
real public-root validation, live deployment actions, live backup/restore,
consumer submission, public publication, and tag/release/package/image
publication remain blocked by scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check returned no output.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets remain present in order and all remain `prepared`.

Claim-boundary status:
No production-readiness, hosted-service, SLA, uptime, compliance,
consumer-acceptance, agency-adoption, vendor, hardware, final-root,
release-readiness, public-launch, or production-grade ETA claim was added.

Security/auth status:
No credential collection, raw env rendering, raw backup rendering, live send,
service-control action, migration execution, browser validator execution,
public admin route, or external contact was added.

Data/migration status:
No migration, schema change, durable deployment state, backup metadata table,
restore table, service-control table, public feed contract change, or module
dependency change was added.

Master review:
Approved. Phase 104 is complete and safe to close.

Required edits:
None for Phase 104.

Decision:
Close Phase 104 and continue immediately to Phase 105 -- Multi-Agency
Isolation And Operator Roles V2.

Next checkpoint:
Phase 105 -- Checkpoint 000001: add multi-agency isolation and operator roles
v2 plan.
