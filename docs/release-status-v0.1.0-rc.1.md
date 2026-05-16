# Release Status -- `v0.1.0-rc.1`

Status date: 2026-05-16

This artifact records the public release-candidate publication status for Open
Transit RT `v0.1.0-rc.1`.

It does not claim stable release readiness, production readiness, compliance,
adoption, consumer submission/review/acceptance/ingestion/listing/display,
final-root readiness, hosted service availability, SLA/uptime, vendor
compatibility, hardware certification, or production-grade ETA quality.

## Current Conclusion

`published_public_release_candidate`

Open Transit RT `v0.1.0-rc.1` was published as a public GitHub prerelease for
local/self-hosted evaluation after Phase 115 release gates passed.

- Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`
- Tag: `v0.1.0-rc.1`
- GitHub Release title: `Open Transit RT v0.1.0-rc.1`
- GitHub Release draft: `false`
- GitHub Release prerelease: `true`
- Published at: `2026-05-16T03:09:40Z`
- Annotated tag object: `52e91c7966e0fe1a5a4202277ab32173f8e78e67`
- Tag target commit: `497f99a97baff630af147c83a7e1249bb08e32da`

The GitHub Release exists and the annotated tag dereferences to the intended
release commit. GitHub-generated source archive download replay is intentionally
left to Phase 116.

## Published Assets

The release was created with the audited local package assets from
`.cache/release-package/v0.1.0-rc.1`.

| Asset | SHA-256 |
| --- | --- |
| `open-transit-rt-v0.1.0-rc.1.source.tar.gz` | `dedf67537b1ed90c24921db32f0df7770aa42968c2d7cbe4927ec9a49a110e6f` |
| `SHA256SUMS.txt` | GitHub asset digest `sha256:82d413da6764a244d19d48ec2f32015fa370be2544ae07ac0c5d20547cede00a` |
| `image.json` | `7034501a0cfcb67b44cc7316ef0f342d43a53eb8395f436146686f5df1e1947c` |
| `manifest.json` | `f3298919ff4e0b6b0ccd1b140aedf19c3bd6c7ea5062c234823f386e5c6bc954` |
| `manifest.md` | `35fa8e38961df8e3fbc7cded913a3729e68d725a43780fefc3de24d44bbcab80` |
| `provenance.json` | `aa3412cde677fc889e66e54e05426dda602bae219e742b9b16b4ccf4507f60f8` |
| `provenance.md` | `754226d4ae12cb9a9f4d200c9796fed25956187c3cf93f3c66d2a6c36df52ba1` |
| `sbom.json` | `34e127d0454c7e64eedb9f1fb9c439f0bfedd61fa0616617aac291a8c8a17a7d` |
| `summary.json` | `c3e0428c3f3dcd3fce698c91a571cdd14fc32f6e2aa48913bdc8f62cc6ef31d5` |
| `summary.md` | `deab40543ff396fb6595a1186f774caa40873099e7b3c0462d8eeb124402edbf` |

The local package summary reported:

- schema: `open-transit-rt-release-package-summary.v1`
- status: `release_ready`
- version: `v0.1.0-rc.1`
- git commit: `497f99a97baff630af147c83a7e1249bb08e32da`
- dirty checkout: `false`
- source archive:
  `artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- SBOM status: `present`
- image status: `not_configured`

## Source Archive Review

The local release package source archive was regenerated after the
`.gitattributes` `export-ignore` policy was committed.

Protected-path scan:

```bash
tar -tzf .cache/release-package/v0.1.0-rc.1/artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz |
  rg '(^|/)docs/evidence/(captured|consumer-submissions/(status\.json|current|artifacts|packets))' |
  wc -l
```

Result: `0`

The public-distribution blocker recorded in Phase 112 was resolved for the
audited local source archive without modifying protected evidence paths.

## Phase 115 Gate Matrix

| Gate | Status | Evidence |
| --- | --- | --- |
| Clean worktree | passed | `git status --short` returned no output before tagging. |
| Release package script tests | passed | `make test-release-package`. |
| Strict package generation | passed | Clean package generated at `.cache/release-package/v0.1.0-rc.1` from `497f99a97baff630af147c83a7e1249bb08e32da`. |
| Package audit | passed | `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package`. |
| Source archive protected-path scan | passed | protected hits: `0`. |
| Lightweight repo checks | passed | `make check`. |
| Validation | passed | `make validate`. |
| Unit/integration tests | passed | `make test`. |
| HTTP smoke | passed | `make smoke`. |
| Compose config | passed | `docker compose -f deploy/docker-compose.yml config`. |
| Release-candidate diagnostics | passed with bounded helper status | `.cache/release-candidate-check/20260516T030728Z`; helper exit `0`, 36 passed, 0 blocker, 0 needs_review, 3 not_checked. The three not_checked rows are `make validate`, `make test`, and `make smoke`, which were run separately and passed. |
| Product acceptance audit | passed | `make audit-product-acceptance`. |
| Final claim audit | passed | `make audit-final-claim-review`. |
| Consumer tracker | passed | exact seven targets remain in order and all `prepared`. |
| Protected-path status | passed | no tracked or untracked status under protected evidence paths. |
| GitHub tooling | passed | `gh auth status` active account `ptse8204`; repo `ptse8204/open-transit-rt` visibility `PUBLIC`, viewer permission `ADMIN`. |
| Remote tag/release absence before publication | passed | remote tag absent and `gh release view` returned `release not found` before tag push. |
| Tag push | passed | `git push origin v0.1.0-rc.1`. |
| GitHub prerelease creation | passed | `gh release create v0.1.0-rc.1 ... --prerelease`. |
| Release metadata verification | passed | release is draft `false`, prerelease `true`, URL recorded above. |

## Claim And Boundary Notes

- The release is a public `v0.1.0-rc.1` release candidate for
  local/self-hosted evaluation.
- The release is not a stable release.
- The release does not prove production readiness, compliance, adoption,
  consumer acceptance, final-root readiness, hosted service availability,
  vendor compatibility, hardware certification, SLA/uptime, or production-grade
  ETA quality.
- The release did not move consumer tracker statuses. All seven consumer
  targets remain exactly `prepared`.
- The release did not create retained evidence under protected evidence paths.
- Published package download replay and GitHub-generated archive replay are
  Phase 116 scope.
