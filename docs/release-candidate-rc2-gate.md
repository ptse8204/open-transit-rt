# Release Candidate rc2 Gate

Status date: 2026-05-16

This artifact records the Phase 130 patch-loop and rc2 gate review for Open
Transit RT after the public `v0.1.0-rc.1` prerelease.

It does not publish, authorize, or claim a public `v0.1.0-rc.2` release. It
does not claim stable release readiness, production readiness, compliance,
adoption, consumer acceptance, final-root readiness, hosted service
availability, SLA/uptime, vendor compatibility, hardware certification, or
production-grade ETA quality.

## Current Conclusion

`local_rc2_gate_prepared_publication_not_authorized`

A local `v0.1.0-rc.2`-style package from current source was generated and
audited successfully. The local package fixes the already-published rc1 source
archive `make check` limitation because the extracted rc2-style source archive
passes `make check` without protected consumer tracker state.

No rc2 tag was pushed and no GitHub Release was created. Phase 130 does not
authorize public rc2 publication.

## rc2 Decision

| Question | Decision |
| --- | --- |
| Is rc2 publication required to keep rc1 usable for local evaluation? | No. Phase 117 public fresh-clone install confidence for rc1 passed. |
| Is rc2 publication useful if maintainers want an archive-first path with the current-source `make check` fix and post-rc1 docs/UX/GTFS-RT hardening? | Yes, recommended if a maintainer separately authorizes a public rc2 cut. |
| Was rc2 published in Phase 130? | No. Publication was not authorized by this phase. |
| Were protected paths or consumer statuses changed? | No. |

## Local rc2-Style Package

Command:

```bash
RELEASE_PACKAGE_VERSION=v0.1.0-rc.2 \
RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.2 \
RELEASE_PACKAGE_FORCE=true \
RELEASE_PACKAGE_ALLOW_DIRTY=false \
scripts/release-package.sh
```

Package:

- Directory: `.cache/release-package/v0.1.0-rc.2`
- Status: `release_ready`
- Version: `v0.1.0-rc.2`
- Git commit: `3805009cc77c6534725ce69a67b6fae95c23c24c`
- Git describe: `v0.1.0-rc.1-58-g3805009`
- Source archive:
  `.cache/release-package/v0.1.0-rc.2/artifacts/open-transit-rt-v0.1.0-rc.2.source.tar.gz`
- SBOM status: `present`
- Image status: `not_configured`

Audit:

```bash
RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.2 \
scripts/audit-release-package.sh
```

Result: passed.

Protected-path archive scan result: `0`.

## Extracted Archive Replay

The local rc2-style source archive was extracted under
`.cache/phase-130/extract-rc2-package`.

Checks:

- `make check`: passed.
- `scripts/bootstrap-dev.sh --check`: passed.

This demonstrates that the current-source archive path has the Phase 116
`make check` fix. It does not change the already-published rc1 archive.

## Release-Candidate Diagnostics

Command:

```bash
RUN_LOCAL_APP=true \
RUN_RELEASE_PACKAGE=true \
RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.2 \
OUTPUT_DIR=.cache/phase-130/release-candidate-check \
FORCE=true \
scripts/release-candidate-check.sh
```

Result:

- Output directory: `.cache/phase-130/release-candidate-check`
- Overall status: `not_checked`
- Passed: `36`
- Needs review: `0`
- Blockers: `0`
- Not checked: `3`

The three `not_checked` rows are `make validate`, `make test`, and
`make smoke`, which are intentionally recorded outside the bounded helper and
must be run separately before any public rc2 cut.

## Remote Publication Status

- `gh release view v0.1.0-rc.2 --repo ptse8204/open-transit-rt`: `release not found`
- `git ls-remote --tags origin refs/tags/v0.1.0-rc.2`: no output

## Claim And Boundary Notes

- rc1 remains the only published public release candidate.
- rc2 is a local gate artifact only until a maintainer separately authorizes
  tag push and GitHub Release creation.
- No protected evidence path was edited, generated, reformatted, or touched.
- `docs/evidence/consumer-submissions/status.json` was not edited.
- The consumer tracker remains exactly seven targets in order and all
  `prepared`.
