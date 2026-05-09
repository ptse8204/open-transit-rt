# Phase 57 Handoff -- Release Packaging And Supply Chain

## Status

Phase 57 is complete for the approved local release packaging and supply-chain
scaffolding scope.

## What Changed

- Added `scripts/release-package.sh` for local `.cache` release packages.
- Added `scripts/audit-release-package.sh` for package file-set, JSON,
  checksum, claim-flag, wording, and consumer tracker audits.
- Added `scripts/test-release-package.sh` for local script coverage and audit
  mutation tests.
- Added `make release-package`, `make audit-release-package`, and
  `make test-release-package`.
- Updated release, upgrade, dependency, decision, roadmap, backlog,
  open-question, status, and latest handoff docs.

## Package Contract

The default package output is ignored `.cache/release-package/<version>/` and
contains:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `provenance.json`
- `provenance.md`
- `sbom.json`
- `image.json`
- `artifacts/open-transit-rt-<version>.source.tar.gz`
- `checksums/SHA256SUMS.txt`

The source archive is created from `git archive HEAD`, not recursive working
tree copying. Dirty local packages are marked not release-ready. Actual release
use should run from a clean checkout with strict settings.

## Boundaries Preserved

Phase 57 created local packaging scaffolding only. It did not publish artifacts,
push images, create a GitHub release, upload to registries, create retained
evidence, write under `docs/evidence`, contact consumers, or change consumer
statuses.

No hosted SaaS, hosted service, production image publication,
production-readiness, compliance, agency adoption, consumer acceptance, vendor
compatibility, marketplace approval, SLA/uptime, or production-grade ETA claim
was created.

`docs/evidence/consumer-submissions/status.json`, current target records,
consumer artifact/packet directories, and `docs/evidence/captured` remain
unchanged. All seven consumer and aggregator targets remain `prepared`.

## Verification

Final verification was run from `/Users/edwintse/Downloads/open-transit-rt`.

- `sh -n scripts/release-package.sh scripts/audit-release-package.sh scripts/test-release-package.sh`
- `./scripts/test-release-package.sh`
- `make release-package`
- `RELEASE_PACKAGE_DIR=<generated-dir> make audit-release-package`
- `make validate`
- `make test`
- `make smoke`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured`
- consumer artifact directory scan printed no files
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured`
- `docker compose -f deploy/docker-compose.yml config`
- `INTEGRATION_TESTS=1 make test-integration`

## Next Phase

Proceed to Phase 58 -- Optional Marketplace / Vendor-Equivalent Pack. Keep
BYOD/hardware, support, SLA/KPI template, and procurement documentation
strictly template/boundary-oriented unless retained claim-specific evidence
exists.
