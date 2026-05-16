# Release Status -- `v0.1.0-rc.1`

Status date: 2026-05-16

This artifact records the current release-candidate publication status for
Open Transit RT `v0.1.0-rc.1`.

It does not publish a release, create a tag, create retained evidence, move
consumer statuses, or claim stable release readiness, production readiness,
compliance, adoption, consumer acceptance, final-root readiness, hosted
service availability, SLA/uptime, vendor compatibility, hardware
certification, or production-grade ETA quality.

## Current Conclusion

`blocked_public_distribution_review`

The local package helper and package audit passed for a clean local package
generated from commit `ee8e186bcc4c879fa0e64ddfbfe5464f550e66de`, but the
source archive public-distribution gate is blocked because the generated source
archive contains tracked protected evidence and consumer-submission paths.

Do not publish `v0.1.0-rc.1` until this is resolved by an explicit maintainer
decision, a revised release-archive policy, or a source package that excludes
public-distribution blockers without touching protected paths.

## Local Tooling And Remote State

- Local branch during Phase 112 audit: `main`.
- `gh` installed: `gh version 2.92.0`.
- `gh auth status`: active authenticated account `ptse8204` with `repo` scope.
- `gh repo view ptse8204/open-transit-rt --json nameWithOwner,visibility,viewerPermission`:
  repository is `PUBLIC` and viewer permission is `ADMIN`.
- `git ls-remote --tags origin v0.1.0-rc.1`: no matching remote tag returned.
- `gh release view v0.1.0-rc.1 --repo ptse8204/open-transit-rt`: `release not found`.

These facts show publication tooling is available locally, but they do not
override the failed source-archive public-distribution review.

## Local Package Result

- Package directory: `.cache/release-package/v0.1.0-rc.1`
- Generated at: `20260516T023241Z`
- Package commit: `ee8e186bcc4c879fa0e64ddfbfe5464f550e66de`
- Dirty checkout: `false`
- Package helper status: `release_ready`
- Package audit: passed locally
- Source archive:
  `.cache/release-package/v0.1.0-rc.1/artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- Source archive SHA-256:
  `5a1caa29e5bf00d5b9fd6303aee6b067a386a9834e3d05a86c263ef8766535bf`
- Checksum manifest SHA-256:
  `6f3531b69a3e84879e7d89f5e776781571f3a7f8ed45960800ef9186115086f4`
- SBOM SHA-256:
  `180e705ef622c3e559e4e53b374ce33c6c2c7449b56d3baeb2323e7a700d17c6`
- Provenance SHA-256:
  `e8ac054e813f70ac6fec31ea382876575d6f1e50645d3641d77a888d4ae95939`
- Manifest SHA-256:
  `0a7b6842a6a1ccd5a756198de735e6710f5ea953618a6c06d9e4a88acfc4d2e8`

The helper's `release_ready` field means the local package structure and
metadata passed the helper's checks. It does not mean public distribution is
approved.

## Source Archive Review

Archive scan command:

```bash
tar -tzf .cache/release-package/v0.1.0-rc.1/artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz |
  rg '(^|/)docs/evidence/(captured|consumer-submissions/(status\.json|current|artifacts|packets))'
```

Source archive counts:

| Item | Count |
| --- | ---: |
| Total archive entries | 1263 |
| Total archive files | 1263 |
| `docs/evidence/**` entries | 215 |
| `docs/evidence/**` files | 215 |
| Protected-path entries | 182 |
| Protected-path files | 182 |

First protected-path examples in the archive:

- `docs/evidence/captured/README.md`
- `docs/evidence/captured/hosted-pending/2026-04-22/README.md`
- `docs/evidence/captured/local-demo/2026-04-22/SHA256SUMS.txt`
- `docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/public_feeds.json`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/google-maps.md`
- `docs/evidence/consumer-submissions/artifacts/google-maps/README.md`
- `docs/evidence/consumer-submissions/packets/google-maps/README.md`

This is a public-distribution blocker for Phase 112. The audit did not modify
or delete any protected path.

## Claim And Secret Scan Notes

Local final-claim audit passed. A bounded source-archive text scan found
claim-like strings only in negated examples, audit fixtures, prompts, or tests,
such as "do not use wording that says the project is compliant" and test data
that intentionally verifies unsafe wording is rejected.

A conservative secret-pattern scan produced expected hits in examples, tests,
templates, placeholder environment files, and documented command snippets
because it matches generic strings such as `Authorization:`, `postgres://`,
`token=`, and placeholder passwords. Phase 112 did not classify those expected
examples as direct leaked production credentials. Protected-path archive
content remains the blocking public-distribution issue.

## Release-Candidate Diagnostic

- Output directory: `.cache/release-candidate-check/20260516T023259Z`
- Helper exit: `0`
- Overall status: `not_checked`
- Counts: 36 passed, 0 blocker, 0 `needs_review`, 3 `not_checked`
- Local app startup and five public feed fetches: passed
- Release package audit row: passed
- Not checked by helper: `make validate`, `make test`, `make smoke`
- Local app cleanup: `make agency-app-down` passed

This diagnostic is local/private only. It is not retained evidence or proof of
consumer acceptance, compliance, hosted service availability, production
readiness, vendor compatibility, or SLA/uptime.

## Publication Gate Matrix

| Gate | Current status | Evidence |
| --- | --- | --- |
| Clean local package generation | passed locally | Package generated from clean commit `ee8e186bcc4c879fa0e64ddfbfe5464f550e66de`. |
| Local package audit | passed locally | `RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package`. |
| Source archive public-distribution review | blocked | Archive contains 182 protected-path entries. |
| Release notes bounded/truthful | needs_review | Draft refreshed in Phase 112; must be reread before Phase 115. |
| Protected paths untouched | passed | Protected-path git status check returned no output. |
| Consumer tracker prepared-only | passed | Exact seven targets remain `prepared`. |
| GitHub tooling | available locally | `gh` installed/authenticated; repository viewer permission `ADMIN`. |
| Remote tag | absent | `git ls-remote --tags origin v0.1.0-rc.1` returned no tag. |
| GitHub Release | absent | `gh release view ...` returned `release not found`. |
| Publication decision | blocked | Do not publish while the source-archive public-distribution gate is blocked. |

## Unlock Conditions

Publication could be reconsidered in Phase 115 only after all of the following
are true:

- the source archive public-distribution blocker is resolved without touching
  protected paths or weakening claim boundaries;
- release notes are refreshed for the exact release commit and remain bounded;
- local package generation and audit pass again from the intended final commit;
- protected-path and prepared-only consumer tracker checks pass;
- final claim audit passes;
- the worktree is clean;
- `gh` authentication and repository permissions remain available.
