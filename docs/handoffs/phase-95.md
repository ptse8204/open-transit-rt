# Phase 95 Handoff -- v0.1.0-rc.1 Candidate Cut

## Status

Phase 95 is complete for the local `v0.1.0-rc.1` candidate cut scope. A local
source package was generated and audited under `.cache`, and package-enabled
release-candidate diagnostics were run. No tag, GitHub Release, public package
distribution, image publication, retained evidence action, consumer action, or
release-ready claim was created.

## Completed Checkpoints

- Phase 95 -- Checkpoint 000001: add v0.1.0-rc.1 candidate cut plan.
- Phase 95 -- Checkpoint 000002: refresh candidate docs.
- Phase 95 -- Checkpoint 000003: generate and audit local candidate package.
- Phase 95 -- Checkpoint 000004: close v0.1.0-rc.1 candidate cut review.

## Local Package Result

- Package path: `.cache/release-package/v0.1.0-rc.1`
- Source archive:
  `.cache/release-package/v0.1.0-rc.1/artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- Source archive SHA-256:
  `ef7f667cf8e0a4238d78ebbb2812c40250e40857057a75d55d2640c781724214`
- Package commit:
  `9684403b9090c948477870636de59b485df42009`
- Package dirty flag: `false`
- SBOM status: `present`
- SBOM module count: 73
- Image metadata status: `not_configured`
- Package audit: passed locally

The package was generated from the clean CP000002 commit. Later CP000003 and
CP000004 documentation commits record package results and closeout; they are
not part of the generated archive.

## Release-Candidate Diagnostic

- Output directory: `.cache/release-candidate-check/20260514T235805Z`
- Helper overall status: `not_checked`
- `git_clean`: `passed`
- `claim_audit`: `passed`
- `local_app_five_feeds`: `passed`
- `release_package_audit`: `passed`

The helper still records `make validate`, `make test`, and `make smoke` as
follow-up rows outside its bounded output. Phase 95 closeout separately ran
`make validate` and `make test`; `make smoke` was not part of the required
closeout baseline and was not run.

## Validation

Passed:

- `git diff --check`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make test-release-package`
- `make release-package` with `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1`,
  `RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1`,
  `RELEASE_PACKAGE_ALLOW_DIRTY=false`, `RELEASE_PACKAGE_STRICT=true`, and
  `RELEASE_PACKAGE_FORCE=true`
- `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package`
- `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true RUN_LOCAL_APP=true make release-candidate-check`
- `make agency-app-down`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- protected-path status check

Blocked:

- `git tag`
- `git push --tags`
- GitHub Release creation
- Public package distribution or release asset upload
- Image publication
- Retained evidence creation
- Consumer action or status movement
- External contact
- Public archive-content review

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched by
tracked changes. Generated package and release-candidate diagnostics stayed
under ignored `.cache`.

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

Phase 95 generated local package diagnostics only. It makes no release
readiness, compliance, adoption, consumer acceptance, production readiness,
final-root readiness, hosted-service availability, vendor compatibility,
hardware certification, SLA/uptime, or ETA-quality claim.

## Security/Auth Status

No route, auth, credential, CSRF, token, or protected data behavior changed.
Package metadata was audited for unsafe strings and false claim flags. The
generated source archive contents still require separate public-distribution
review before any future release action.

## Data/Migration Status

No persistence, migration, GTFS data model, tenant model, or realtime data
model change is included.

## Checkpoint Report

Checkpoint:
Phase 95 -- Checkpoint 000004: close v0.1.0-rc.1 candidate cut review.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high; real Release / Supply-Chain Sub-Agent -- GPT-5.5
high; real Claim-Boundary / Security QA Sub-Agent -- GPT-5.5 high.
Implementation, QA, and Documentation closeout roles were simulated by the
Master Agent. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/release-notes-v0.1.0-rc.1-draft.md`;
`docs/phase-95-v0-1-0-rc-1-candidate-cut.md`;
`docs/handoffs/phase-95.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`.

Validation run:
`git status --short`; `git diff --check`; `make check`; `make validate`;
`make test`; `docker compose -f deploy/docker-compose.yml config`; `make
test-release-package`; `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1
RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1
RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true
RELEASE_PACKAGE_FORCE=true make release-package`;
`RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make
audit-release-package`; `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1
RUN_RELEASE_PACKAGE=true RUN_LOCAL_APP=true make release-candidate-check`;
`make agency-app-down`; `make audit-product-acceptance`; `make
audit-final-claim-review`; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null`; exact
prepared-only consumer tracker assertion; protected-path status check.

Blocked checks:
`git tag`, `git push --tags`, GitHub Release creation, public package
distribution, image publication, retained evidence creation, consumer action,
external contact, and public archive-content review are blocked by Phase 95
scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched by
tracked changes. Generated package and release-candidate diagnostics stayed
under ignored `.cache`.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven consumer targets remain present in order and all remain `prepared`.

Claim-boundary status:
Phase 95 generated local package diagnostics only. It makes no release
readiness, compliance, adoption, consumer acceptance, production readiness,
final-root readiness, hosted-service availability, vendor compatibility,
hardware certification, SLA/uptime, or ETA-quality claim.

Security/auth status:
No route, auth, credential, CSRF, token, or protected data behavior changed.
Package metadata was audited for unsafe strings and false claim flags. The
generated source archive contents still require separate public-distribution
review before any future release action.

Data/migration status:
No persistence, migration, GTFS data model, tenant model, or realtime data
model change is included.

Master review:
Approved. The phase satisfied the local package-generation authorization,
recorded package/audit diagnostics, and preserved all no-tag, no-release,
no-publication, no-evidence, no-consumer-status, and no-claim boundaries.

Required edits:
None.

Decision:
Close Phase 95 and continue immediately to Phase 96.

Next checkpoint:
Phase 96 -- Checkpoint 000001: add GTFS versioning diff and rollback workbench
plan.
