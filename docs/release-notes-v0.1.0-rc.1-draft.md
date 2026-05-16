# Open Transit RT `v0.1.0-rc.1` Draft Release Notes

Draft date: 2026-05-16

This is a local pre-tag draft refreshed during Phase 112. It
is not a GitHub Release, tag, hosted deployment, retained evidence packet,
public package publication, image publication, consumer action, or release
readiness claim.

## Source

- Candidate label: `v0.1.0-rc.1`
- Git tag: None; not tagged.
- Source package: local `.cache` diagnostic package at
  `.cache/release-package/v0.1.0-rc.1`.
- Source package commit:
  `c9fd75765837a5e4812e94b6e91250fd0f3d679b`
- Release notes link: local draft at
  `docs/release-notes-v0.1.0-rc.1-draft.md`
- Current gate references:
  - `docs/phase-89-rc1-gate-results.md`
  - `docs/phase-92-clean-checkout-release-candidate-gate.md`
  - `docs/phase-95-v0-1-0-rc-1-candidate-cut.md`
  - `docs/phase-108-post-rc-bug-bash-and-stabilization.md`
  - `docs/phase-111-goal-activation-and-public-release-roadmap-pack.md`
  - `docs/phase-112-public-release-artifact-and-claim-blocking-audit.md`
  - `docs/release-status-v0.1.0-rc.1.md`
- Artifact checksums: local checksum manifest at
  `.cache/release-package/v0.1.0-rc.1/checksums/SHA256SUMS.txt`.
- Source archive checksum:
  `0a3e5476983724b82eea65e4654771c88652bf7c6c25faf245c9898525d16069`
- SBOM/provenance: local metadata generated under
  `.cache/release-package/v0.1.0-rc.1`; SBOM status `present`.
- Published Docker image: None.
- GitHub Release: None; `gh release view v0.1.0-rc.1 --repo
  ptse8204/open-transit-rt` returned `release not found` during Phase 112.

## Summary

`v0.1.0-rc.1` remains a maintainer-evaluation milestone for the self-hosted
small-agency product path. Since the earlier Phase 89 draft, the autonomous
post-90 roadmap has added route/product stabilization, a clean-checkout local
release-candidate gate, local/private browser task trials, and an Operations
Console route-registry refactor. Later post-90 phases added operator-facing
GTFS versioning, quality planning, realtime usefulness review, prediction
backtesting, alerts operations, connector recipes, device onboarding,
monitoring exports, small-host guidance, multi-agency/role hardening, staff
training material, public docs/contributor alignment, and this post-RC bug
bash review.

The candidate remains blocked for public distribution review. Phase 115 is
authorized to attempt publication only if release gates pass and authenticated
tooling is available. Phase 112 found that the regenerated source archive
contains protected evidence and consumer-submission paths, so publication must
not proceed until that blocker is resolved or explicitly accepted by a
maintainer release decision.

## User-Facing Changes

- Private Operations Console route maps and product docs are aligned around
  the newer center-style routes: GTFS Workbench, Realtime Center, Validation
  Center, Connector Workbench, Prediction & ETA Lab, Access & Roles, and Audit
  Log.
- Start Here now includes role-based entry points for a new agency evaluator,
  daily operations staff, technical helper, release reviewer, and connector
  evaluator.
- Telemetry Simulator guidance now calls out the first local/synthetic
  `dry-run` safety check and states that it does not test a live vendor, live
  AVL API, real device, or public feed consumer.
- Operations Console route metadata is centralized in a private route registry
  that drives navigation, page titles, tests, and route inventory audit logic.
- The route inventory audit now covers 28 canonical private HTML routes, 20
  canonical private JSON routes, one command route, and two external admin
  surfaces.
- Post-RC bug bash reran the route inventory audit and strict docs mode; both
  passed locally before the full Phase 108 validation rerun.

## Install And Upgrade Notes

- Clean install from source tag: None; no tag exists for this draft.
- Local app verification: Phase 95 package-enabled
  `RUN_LOCAL_APP=true make release-candidate-check` completed local app startup
  and five-feed diagnostics.
- Release-candidate diagnostic:
  `.cache/release-candidate-check/20260514T235805Z`; helper overall
  `not_checked` because follow-up `make validate`, `make test`, and `make
  smoke` rows are intentionally outside the bounded helper output.
- Latest Phase 112 diagnostic:
  `.cache/release-candidate-check/20260516T022918Z`; helper overall
  `not_checked` with 36 passed rows, 0 blockers, 0 `needs_review` rows, and 3
  intentionally `not_checked` rows.
- Local release package: generated and audited locally under
  `.cache/release-package/v0.1.0-rc.1`.
- Local Docker image build: None.
- Published production Docker image: None.

## Migrations

- No migration file changed during Phase 112.

Before any future upgrade, operators must back up the database and run:

```bash
make migrate-status
make migrate-up
make migrate-status
```

## Operations Changes

- Phase 91 added a route inventory audit helper and route-map alignment.
- Phase 92 recorded a clean-checkout local release-candidate gate with a
  `needs_review` conclusion.
- Phase 93 exercised local/private agency task flows and patched narrow copy
  and IA gaps.
- Phase 94 centralized Operations Console route metadata without changing
  route behavior.
- Phase 95 generated and audited a local `.cache` source package only.
- Phase 96 through Phase 107 added operator workbench, QA, training,
  contributor, and governance-facing documentation/UI improvements while
  preserving release, evidence, and claim boundaries.
- Phase 108 reran the route inventory audit and refreshed this blocker matrix
  before the full stabilization validation checkpoint.

## Security Notes

- Private Operations Console route boundaries, role checks, agency query
  guards, CSRF expectations for browser POSTs, request-body caps, and no-store
  handling remain in force.
- Phase 112 local package tooling must keep outputs under ignored `.cache`
  directories and must not write protected evidence paths.
- The local source archive is generated by `git archive HEAD`; it includes
  tracked repository files. Because tracked evidence-path material exists in
  the repository, any public distribution requires a separate review before
  publication. Phase 112 does not publish the archive and records the current
  source-archive review as blocked.
- Operators must not copy secrets, tokens, database URLs, raw private output,
  credentials, or private records into public docs, issue trackers, release
  text, or package metadata.

## Dependency Changes

- None in Phase 112 release artifact audit.
- Phase 95 package generation recorded local SBOM/provenance metadata. The
  local package helper reported clean source metadata, but this is not a
  release-readiness, publication, or public-distribution claim.

## Evidence Or Claim Changes

- None.

This draft does not create retained evidence, change consumer statuses, contact
external services, publish artifacts, tag a release, distribute a package, or
add claims of CAL-ITP/Caltrans compliance, consumer
submission/review/acceptance/ingestion/listing/display, agency
adoption/approval, final-root readiness, hosted service availability, paid
support, SLA/uptime, production readiness, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality,
validator-clean feeds, or release readiness.

## Draft Tag Command

Draft only. Do not run before the Phase 115 publication gate passes.

```bash
git tag -a v0.1.0-rc.1 <final-reviewed-clean-main-sha> -m "Open Transit RT v0.1.0-rc.1"
```

No `git tag`, `git push --tags`, or tag publication is authorized by this
draft. Phase 115 is the first phase that may attempt publication if all gates
pass.

## Draft GitHub Release Text

Title:

```text
Open Transit RT v0.1.0-rc.1
```

Body draft:

```md
This is draft text for a future release candidate. It has not been published.

Scope: local source package review, checksum manifest, SBOM/provenance
metadata, release notes refresh, and maintainer review.

Artifacts: local `.cache` package only unless a maintainer later creates
release assets.

Not included: no tag created, no GitHub Release created, no image pushed, no
package registry publication, no consumer status movement, no retained
evidence writes, and no release-ready, production-readiness, compliance,
consumer-acceptance, public-launch, final-root-readiness, hosted-service,
vendor, SLA, or ETA-quality claim.

Known blockers: record any failed or skipped validation commands exactly
before any later release action.
```

No GitHub Release creation is authorized by this draft. Phase 115 is the first
phase that may attempt publication if all gates pass.

## Known Blockers Matrix

| Area | Status | Exact blocker or note | Next owner |
| --- | --- | --- | --- |
| Release action authorization | gated | Phase 115 authorizes public `v0.1.0-rc.1` publication only if all gates pass and authenticated tooling is available. | Maintainer / Phase 115 |
| Local package generation | passed locally | `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1 RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true make release-package` wrote local `.cache` output only. | Master Agent |
| Local package audit | passed locally | `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package` passed metadata, checksum, claim-flag, and consumer-tracker checks. | Master Agent |
| Source archive public-distribution review | blocked | Phase 112 archive scan found 182 protected-path entries under `docs/evidence/captured/**` and `docs/evidence/consumer-submissions/**` protected areas. | Maintainer / Phase 115 |
| Post-RC route audit | passed locally | `make audit-operations-route-inventory` and strict docs mode passed during Phase 108 Checkpoint 000002. | Master Agent |
| Full post-RC validation rerun | passed locally | Phase 108 Checkpoint 000003 reran baseline, route inventory, validation, tests, connector checks, release-candidate diagnostic, compose config, prepared-only tracker assertion, protected-path status, and claim audits. | Master Agent |
| Diagnostic helper Java detail | needs_review | Phase 108 helper summary still marked the Java tool row `passed` while its detail contained the macOS system-stub message about no Java runtime; independent pinned validator check and `make validate` passed. | Future release reviewer |
| Package publication | blocked | Public package distribution or GitHub Release asset upload must not happen while source archive public-distribution review is blocked. | Maintainer / Phase 115 |
| Git tag | blocked | Remote `v0.1.0-rc.1` tag is absent, and tag creation must not happen while source archive public-distribution review is blocked. | Maintainer / Phase 115 |
| GitHub Release | absent | `gh release view v0.1.0-rc.1 --repo ptse8204/open-transit-rt` returned `release not found`. | Maintainer / Phase 115 |
| Published image | blocked | No image build or publication is authorized while source archive public-distribution review is blocked. | Future authorized release work |
| Evidence tracks | blocked | Final-root, consumer, real agency pilot, real vendor/device, ETA-quality, and compliance evidence gates require separate written authorization. | Authorization-gated only |
| Final conclusion | blocked_public_distribution_review | Local diagnostics may pass, but publication remains blocked by source archive protected-path contents. | Maintainer review |

## Known Limitations

- This draft is not tied to a release tag.
- Phase 112 local package output, once generated, remains a local diagnostic
  artifact under `.cache`.
- Connector/adaptor results are synthetic/local checks only and do not prove
  real vendor/device compatibility.
- Validator and route diagnostics are local product signals only and do not
  prove compliance, consumer action, final-root readiness, production
  readiness, hosted-service availability, or SLA/uptime.

## Checks

Recorded Phase 95 package checks:

| Check | Result | Scope |
| --- | --- | --- |
| `make test-release-package` | passed | Local package helper tests. |
| `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1 RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true RELEASE_PACKAGE_FORCE=true make release-package` | passed | Local `.cache` source archive, checksum, SBOM, provenance, manifest, image metadata summary. |
| `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package` | passed | Local package metadata/checksum/claim-flag audit. Does not inspect source archive contents for public distribution. |
| `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true RUN_LOCAL_APP=true make release-candidate-check` | exited `0`; helper overall `not_checked` | Local app and five-feed diagnostics passed; package audit passed; bounded helper still records `make validate`, `make test`, and `make smoke` as follow-up rows. |
| `make agency-app-down` | passed | Stopped local app containers after diagnostics. |

Phase 95 closeout checks are recorded in `docs/handoffs/phase-95.md`.

Recorded Phase 108 stabilization checks:

| Check | Result | Scope |
| --- | --- | --- |
| `git status --short` | passed | Clean at start and after the validation rerun. |
| `git diff --check` | passed | No whitespace errors. |
| `make audit-operations-route-inventory` | passed | 28 HTML routes, 20 JSON routes, one command route, two external admin surfaces, no public admin route, README/wiki route maps aligned. |
| `OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true scripts/audit-operations-route-inventory.sh` | passed | Strict route/docs inventory mode passed. |
| `make check` | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| `make audit-product-acceptance` | passed | Product acceptance audit passed with prepared-only tracker and clean protected-path checks. |
| `make audit-final-claim-review` | passed | Final claim review audit passed. |
| `make validate` | passed | Validator tooling check and validation smoke passed. |
| `make test` | passed | `go test ./...` passed. |
| `make external-connection-check` | passed | Connector manifests remain sidecar/manifest/conformance bounded. |
| `make adapter-conformance` | passed | Adapter conformance suite passed. |
| `make test-connector-examples` | passed | Connector examples passed. |
| `RUN_LOCAL_APP=true make release-candidate-check` | exited `0`; helper overall `not_checked` | Local app and five-feed diagnostics passed; helper kept validation, tests, smoke, and package audit as bounded follow-up/not-checked rows. |
| `make agency-app-down` | passed | Stopped local app containers after diagnostics. |
| `docker compose -f deploy/docker-compose.yml config` | passed | Compose config rendered locally. |
| Consumer tracker JSON parse and prepared-only assertion | passed | Seven targets remain exactly `prepared`. |
| Protected-path status check | passed | No status under protected evidence paths, migrations, `go.mod`, or `go.sum`. |

Recorded Phase 112 release artifact audit checks:

| Check | Result | Scope |
| --- | --- | --- |
| `make test-release-package` | passed | Local package helper tests. |
| `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1 RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true RELEASE_PACKAGE_FORCE=true make release-package` | passed | Local `.cache` source archive, checksum, SBOM, provenance, manifest, and image metadata summary from clean commit `c9fd75765837a5e4812e94b6e91250fd0f3d679b`. |
| `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package` | passed | Local package metadata/checksum/claim-flag audit. Does not clear public distribution of source archive contents. |
| Source archive protected-path scan | blocked | Archive contains 182 protected-path entries under `docs/evidence/captured/**` and protected `docs/evidence/consumer-submissions/**` areas. |
| `RUN_LOCAL_APP=true RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true make release-candidate-check` | exited `0`; helper overall `not_checked` | Local app and five-feed diagnostics passed; package audit passed; bounded helper keeps `make validate`, `make test`, and `make smoke` as follow-up rows. |
| `make agency-app-down` | passed | Stopped local app containers after diagnostics. |
| `gh repo view ptse8204/open-transit-rt --json nameWithOwner,visibility,viewerPermission` | passed | Repository is public and local viewer permission is `ADMIN`; this does not override the source-archive blocker. |
| `gh release view v0.1.0-rc.1 --repo ptse8204/open-transit-rt` | release not found | No GitHub Release exists for the candidate. |

Blocked checks and actions:

- `git tag`: blocked by source archive public-distribution review.
- `git push --tags`: blocked by source archive public-distribution review.
- GitHub Release creation: blocked by source archive public-distribution review.
- Public image/package publication: blocked by source archive public-distribution review.

Release-candidate summary:

- Output summary: `.cache/release-candidate-check/20260516T022918Z`.
- Overall status: `not_checked`.
- Package audit status: passed locally.
- Claim-boundary result: no retained evidence, consumer status change, image
  push, hosted-service claim, production-readiness claim, release-readiness
  claim, or compliance claim is made.
