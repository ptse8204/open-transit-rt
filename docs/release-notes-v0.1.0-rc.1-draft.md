# Open Transit RT `v0.1.0-rc.1` Draft Release Notes

Draft date: 2026-05-14

This is a local pre-tag draft refreshed during Phase 108 Checkpoint 000002. It
is not a GitHub Release, tag, hosted deployment, retained evidence packet,
public package publication, image publication, consumer action, or release
readiness claim.

## Source

- Candidate label: `v0.1.0-rc.1`
- Git tag: None; not tagged.
- Source package: local `.cache` diagnostic package at
  `.cache/release-package/v0.1.0-rc.1`.
- Source package commit:
  `9684403b9090c948477870636de59b485df42009`
- Release notes link: local draft at
  `docs/release-notes-v0.1.0-rc.1-draft.md`
- Current gate references:
  - `docs/phase-89-rc1-gate-results.md`
  - `docs/phase-92-clean-checkout-release-candidate-gate.md`
  - `docs/phase-95-v0-1-0-rc-1-candidate-cut.md`
  - `docs/phase-108-post-rc-bug-bash-and-stabilization.md`
- Artifact checksums: local checksum manifest at
  `.cache/release-package/v0.1.0-rc.1/checksums/SHA256SUMS.txt`.
- Source archive checksum:
  `ef7f667cf8e0a4238d78ebbb2812c40250e40857057a75d55d2640c781724214`
- SBOM/provenance: local metadata generated under
  `.cache/release-package/v0.1.0-rc.1`; SBOM status `present` with 73 modules.
- Published Docker image: None.
- GitHub Release: None.

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

The candidate remains `needs_review`. Phase 95 authorizes only local `.cache`
package generation and audit; it does not authorize tag creation, GitHub
Release creation, public package distribution, image publication, retained
evidence creation, consumer status movement, or a release-readiness claim.

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
- Latest post-RC diagnostic:
  `.cache/release-candidate-check/20260515T032047Z`; helper overall
  `not_checked` with 35 passed rows, 0 blockers, 0 `needs_review` rows, and 4
  intentionally `not_checked` rows.
- Local release package: generated and audited locally under
  `.cache/release-package/v0.1.0-rc.1`.
- Local Docker image build: None.
- Published production Docker image: None.

## Migrations

- No migration file changed during Phase 108.

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
- Phase 95 local package tooling must keep outputs under ignored `.cache`
  directories and must not write protected evidence paths.
- The local source archive is generated by `git archive HEAD`; it includes
  tracked repository files. Because tracked evidence-path material exists in
  the repository, any public distribution requires a separate review before
  publication. Phase 95 does not publish the archive.
- Operators must not copy secrets, tokens, database URLs, raw private output,
  credentials, or private records into public docs, issue trackers, release
  text, or package metadata.

## Dependency Changes

- None in Phase 108 documentation refresh.
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

Draft only. Do not run in Phase 95.

```bash
git tag -a v0.1.0-rc.1 <final-reviewed-clean-main-sha> -m "Open Transit RT v0.1.0-rc.1"
```

No `git tag`, `git push --tags`, or tag publication is authorized by this
draft.

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

No GitHub Release creation is authorized by this draft.

## Known Blockers Matrix

| Area | Status | Exact blocker or note | Next owner |
| --- | --- | --- | --- |
| Release action authorization | blocked | No authorization exists to tag, publish, distribute packages, publish images, create a GitHub Release, or claim release readiness. | Maintainer |
| Local package generation | passed locally | `RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.1 RELEASE_PACKAGE_ALLOW_DIRTY=false RELEASE_PACKAGE_STRICT=true make release-package` wrote local `.cache` output only. | Master Agent |
| Local package audit | passed locally | `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package` passed metadata, checksum, claim-flag, and consumer-tracker checks. | Master Agent |
| Post-RC route audit | passed locally | `make audit-operations-route-inventory` and strict docs mode passed during Phase 108 Checkpoint 000002. | Master Agent |
| Full post-RC validation rerun | passed locally | Phase 108 Checkpoint 000003 reran baseline, route inventory, validation, tests, connector checks, release-candidate diagnostic, compose config, prepared-only tracker assertion, protected-path status, and claim audits. | Master Agent |
| Diagnostic helper Java detail | needs_review | Phase 108 helper summary still marked the Java tool row `passed` while its detail contained the macOS system-stub message about no Java runtime; independent pinned validator check and `make validate` passed. | Future release reviewer |
| Package publication | blocked | No public package distribution or GitHub Release asset upload is authorized. | Maintainer only |
| Source archive review | needs_review | `git archive HEAD` includes tracked repository files; public publication requires a separate review for tracked evidence-path material. | Future release reviewer |
| Git tag | blocked | No tag is authorized in Phase 108. | Maintainer only |
| Published image | blocked | No image build or publication is authorized in Phase 108. | Future authorized release work |
| Evidence tracks | blocked | Final-root, consumer, real agency pilot, real vendor/device, ETA-quality, and compliance evidence gates require separate written authorization. | Authorization-gated only |
| Final conclusion | needs_review | Local diagnostics may pass, but release action and publication remain blocked. | Maintainer review |

## Known Limitations

- This draft is not tied to a release tag.
- Phase 95 local package output, once generated, remains a local diagnostic
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

Blocked checks:

- `git tag`: blocked; not authorized.
- `git push --tags`: blocked; not authorized.
- GitHub Release creation: blocked; not authorized.
- Public image/package publication: blocked; not authorized.

Release-candidate summary:

- Output summary: `.cache/release-candidate-check/20260514T235805Z`.
- Overall status: `needs_review`.
- Package audit status: passed locally.
- Claim-boundary result: no retained evidence, consumer status change, image
  push, hosted-service claim, production-readiness claim, release-readiness
  claim, or compliance claim is authorized.
