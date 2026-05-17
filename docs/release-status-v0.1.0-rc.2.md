# Release Status -- `v0.1.0-rc.2`

Status date: 2026-05-16

This artifact records the public release-candidate publication status for Open
Transit RT `v0.1.0-rc.2`.

Open Transit RT v0.1.0-rc.2 is a public release candidate for local/self-hosted evaluation.

It is not a stable release and does not prove production readiness, compliance, agency adoption, consumer acceptance, final-root readiness, hosted service availability, vendor compatibility, hardware certification, SLA/uptime, production AVL reliability, production-grade ETA quality, or real-world ETA accuracy.

## Current Conclusion

`published_public_release_candidate`

Open Transit RT `v0.1.0-rc.2` was published as a public GitHub prerelease for
local/self-hosted evaluation after the rc2 release gate passed.

- Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.2`
- Tag: `v0.1.0-rc.2`
- GitHub Release title: `Open Transit RT v0.1.0-rc.2`
- GitHub Release draft: `false`
- GitHub Release prerelease: `true`
- GitHub latest stable endpoint after publication: `404 Not Found`
- Published at: `2026-05-16T22:32:45Z`
- Annotated tag object: `3b7e9b616d98dee908b48645108035117a68e5dc`
- Tag target commit: `15a0ec7cbdacf2301ac906ff3ecbe655371fccc6`

## Published Assets

The release was created with the audited local package assets from
`.cache/release-package/v0.1.0-rc.2`.

| Asset | SHA-256 |
| --- | --- |
| `open-transit-rt-v0.1.0-rc.2.source.tar.gz` | `da9cf24757498379eeb27a6812f4508d4238033900c1f2906f46e1929b2a2d76` |
| `SHA256SUMS.txt` | `d76167e2e6434683afca5dc1d5809c53c5e5dea323a6070a6cc5ef0e2fee72df` |
| `image.json` | `7034501a0cfcb67b44cc7316ef0f342d43a53eb8395f436146686f5df1e1947c` |
| `manifest.json` | `618df0f32af62ab813e978fe8bd9d847932d2f3a9040a3dd8fcbfae5e80d7a68` |
| `manifest.md` | `f715c3b2aff6d54f08f38123c6da59434a4c5b4404647d90f7c48d5f1242abbb` |
| `provenance.json` | `33b6d14ce0c87d6a8e6563a191d0d51a3d8a977bc9df24da459e981957fa0ffe` |
| `provenance.md` | `45f0935cbdba0c616e9f7e19856882c699d970bd83282e7c74f3839db5d7e769` |
| `sbom.json` | `05804a50e7d4ac2a2d8fbd214eacb21340ca56e76ea08a970b50e4054d0e7139` |
| `summary.json` | `bcbe404c072b626ef954515593e6fdbb870f9d1ea22cef812974e31af3d483e5` |
| `summary.md` | `094379250ad4fe9480103a7a29b694590e7379ec70fa675e8d5f087b125ad094` |

The package summary reported:

- schema: `open-transit-rt-release-package-summary.v1`
- status: `release_ready`
- version: `v0.1.0-rc.2`
- git commit: `15a0ec7cbdacf2301ac906ff3ecbe655371fccc6`
- dirty checkout: `false`
- source archive:
  `artifacts/open-transit-rt-v0.1.0-rc.2.source.tar.gz`
- SBOM status: `present`
- image status: `not_configured`

## Gate Matrix

| Gate | Status | Evidence |
| --- | --- | --- |
| Clean worktree before tagging | passed | `git status --short` returned no output. |
| Remote rc2 absence before tagging | passed | remote tag absent; GitHub Release returned `release not found`. |
| Lightweight repo check | passed | `make check`. |
| Validation | passed | `make validate`. |
| Unit/integration tests | passed | `make test`. |
| HTTP smoke | passed | `make smoke`. |
| Release package script tests | passed | `make test-release-package`. |
| Product acceptance audit | passed | `make audit-product-acceptance`. |
| Final claim audit | passed | `make audit-final-claim-review`. |
| External connection check | passed | `make external-connection-check`. |
| Adapter conformance | passed | `make adapter-conformance`. |
| Connector examples | passed | `make test-connector-examples`. |
| GTFS-RT conformance | passed | `make gtfsrt-conformance`. |
| Compose config | passed | `docker compose -f deploy/docker-compose.yml config`. |
| Package generation | passed | strict clean package generated under `.cache/release-package/v0.1.0-rc.2`. |
| Package audit | passed | `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.2 scripts/audit-release-package.sh`. |
| Release-candidate diagnostics | passed with bounded helper status | `.cache/release-candidate-check/v0.1.0-rc.2`; 36 passed, 0 blocker, 0 needs_review, 3 not_checked. The not_checked rows are `validate`, `test`, and `smoke`, which were run separately and passed. |
| Extracted generated archive replay | passed | `make check`, `scripts/bootstrap-dev.sh --check`, `make validators-install`, `make validate`, and `make test`. |
| Archive install-confidence | passed | `.cache/install-confidence/v0.1.0-rc.2-archive-rerun2`; local app startup and five local public feed fetches passed. |
| Protected-path archive scan | passed | Generated upload candidates had `0` protected-path hits. |
| Tag push | passed | `git push origin v0.1.0-rc.2`. |
| GitHub prerelease creation | passed | `gh release create ... --verify-tag --prerelease --latest=false`. |
| Release metadata verification | passed | release is draft `false`, prerelease `true`, URL recorded above. |

## Protected And Consumer Status

- No protected evidence path was edited, generated, reformatted, or touched.
- `docs/evidence/consumer-submissions/status.json` was not edited.
- The consumer tracker remains exactly seven targets in order and all
  `prepared`.
- Release source archive and upload candidates had zero protected-path hits.

## Claim Boundary Notes

- This release is a public release candidate for local/self-hosted evaluation.
- This release is not a stable release.
- This release does not prove production readiness, compliance, adoption,
  consumer acceptance, final-root readiness, hosted service availability,
  vendor compatibility, hardware certification, SLA/uptime, production AVL
  reliability, production-grade ETA quality, or real-world ETA accuracy.
