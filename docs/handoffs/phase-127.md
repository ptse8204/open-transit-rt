# Phase 127 Handoff -- Small-Host Deployment And Upgrade UX Hardening

## Status

Phase 127 is complete for Small-Host Deployment And Upgrade UX Hardening.

The private Maintenance Center now includes a read-only
`small_host_readiness` panel for small-host preflight and recovery review. The
panel keeps deployment-doctor diagnostics, dry-run feed checks, off-host
validation choices, resource-budget review, backup/restore recovery path, and
upgrade stop points visible before operator or technical-helper action.

## Completed Checkpoints

- Phase 127 -- Checkpoint 000001: add small host deployment and upgrade ux
  hardening plan.
- Phase 127 -- Checkpoint 000002: implement or audit primary scoped work.
- Phase 127 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux
  hardening review.

## Product Result

Small-host maintainers now get a compact private checklist before upgrade or
recovery work:

- deployment-doctor, five-feed dry run, validator plan, backup review, and
  rollback owner sequence
- off-host validation choice for constrained hosts
- conservative resource budget and DB pool guidance
- backup/restore recovery path review
- stop points before migration, service restart, package switch, rollback, or
  public-root announcement

The work is private and read-only. It does not execute deployment preflight,
validators, backups, restore drills, migrations, service changes, package
publication, rollback, release actions, evidence writes, or consumer status
changes.

## Changed Files

- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/main_test.go`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `docs/phase-127-small-host-deployment-and-upgrade-ux-hardening.md`
- `docs/handoffs/phase-127.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- Web Design Skill loaded before private Maintenance Center UX work
- `gofmt -w cmd/agency-config/operations_maintenance.go cmd/agency-config/main_test.go cmd/agency-config/operations.go`
- `go test ./cmd/agency-config -run 'TestOperationsMaintenance|TestOperationsNavigation|TestOperationsRouteTitles|TestOperationsConsole'`
- `git status --short`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json`
- `scripts/check-consumer-tracker.sh`
- protected-path git status check
- `sh -n scripts/deployment-doctor.sh scripts/oci-reference-check.sh scripts/validate-public-feeds.sh scripts/pilot-ops.sh`
- `OUTPUT_DIR=.cache/phase-127/deployment-doctor FORCE=true scripts/deployment-doctor.sh`
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-127/validate-public-feeds FORCE=true scripts/validate-public-feeds.sh --dry-run`
- `PUBLIC_BASE_URL=https://feeds.example.org OUTPUT_DIR=.cache/phase-127/oci-reference-check FORCE=true scripts/oci-reference-check.sh --dry-run`
- generated Phase 127 helper JSON parse check
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make external-connection-check`
- `make adapter-conformance`
- `make gtfsrt-conformance`

Blocked:

- None for Phase 127.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched.

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

Phase 127 makes no stable release readiness, deployment success, production
readiness, compliance, adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, real-world ETA
accuracy, or consumer display claim.

## Security/Auth Status

The Maintenance Center remains private and read-only. No route auth, CSRF
behavior, credential handling, token handling, public exposure, raw env
rendering, raw backup rendering, private payload handling, or browser command
execution behavior was changed.

## Data/Migration Status

No migration, durable state, runtime dependency, or Go module change was
added.

## Release/Publication Status

The Phase 115 public `v0.1.0-rc.1` prerelease remains published. Phase 127 did
not publish, republish, retag, upload assets, or patch the public rc1 release.

## Install Confidence Status

Phase 117 public fresh-clone install confidence remains passed. Phase 127 is
current-source hardening after rc1 and is not part of the published rc1 tag.

## Web Design Skill Status

The Web Design Skill at `~/.agents/skills/web-design-engineer` was loaded
before the private Maintenance Center UX work. The patch preserves the
existing dense private admin table style and low-radius panel vocabulary.

## Commit List

- `471b579` -- Phase 127 -- Checkpoint 000001: add small host deployment and
  upgrade ux hardening plan
- `ca7b9f5` -- Phase 127 -- Checkpoint 000002: implement or audit primary
  scoped work
- `60866aa` -- Phase 127 -- Checkpoint 000003: run validation and patch
  required gaps
- Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux
  hardening review

## Checkpoint Report

Checkpoint:
Phase 127 -- Checkpoint 000004: close Small-Host Deployment And Upgrade UX
Hardening review.

Goal status:
Active. Phase 127 is closed and the goal continues to Phase 128.

Sub-agents used or simulated:
Context / Repo Truth, Planning, Implementation, QA, UI/UX, Web Design Skill,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Release,
Install Confidence, Connector, and GTFS-RT Domain roles were simulated by the
Master Agent because the agent thread limit prevented new real sub-agents.

Changed files:
`docs/handoffs/phase-127.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-127-small-host-deployment-and-upgrade-ux-hardening.md`.

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
guidance and makes no stronger public claim.

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
