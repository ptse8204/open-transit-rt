# Release Download Replay -- `v0.1.0-rc.1`

Status date: 2026-05-16

This artifact records Phase 116 public download replay for Open Transit RT
`v0.1.0-rc.1`.

It verifies release publication mechanics and archive/download behavior only.
It does not claim stable release readiness, production readiness, compliance,
adoption, consumer acceptance, final-root readiness, hosted service
availability, SLA/uptime, vendor compatibility, hardware certification, or
production-grade ETA quality.

## Conclusion

`download_replay_needs_patch_for_release_archive_make_check`

The public GitHub Release exists, published release assets downloaded, uploaded
asset checksums matched the published `SHA256SUMS.txt`, and both the uploaded
source archive and GitHub-generated tag archives contained zero protected
evidence or protected consumer-submission paths.

However, extracting the published `v0.1.0-rc.1` source archive and running
`make check` fails because the protected consumer tracker
`docs/evidence/consumer-submissions/status.json` is intentionally excluded from
public source archives by `.gitattributes`, while the `make check` target in
the published rc1 tag still required that file.

Phase 116 patched the current repository so future source archives can run
lightweight checks without restoring protected evidence. This does not modify
the already published rc1 source archive, so the rc1 release-archive replay
must remain recorded with this blocker.

## Release Metadata

- Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`
- Tag: `v0.1.0-rc.1`
- GitHub Release draft: `false`
- GitHub Release prerelease: `true`
- Published at: `2026-05-16T03:09:40Z`
- Tag target commit: `497f99a97baff630af147c83a7e1249bb08e32da`
- Download directory: `.cache/release-download/v0.1.0-rc.1`

## Downloaded Uploaded Assets

Downloaded with:

```bash
gh release download v0.1.0-rc.1 \
  --repo ptse8204/open-transit-rt \
  --dir .cache/release-download/v0.1.0-rc.1 \
  --clobber
```

Downloaded assets:

- `SHA256SUMS.txt`
- `image.json`
- `manifest.json`
- `manifest.md`
- `open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- `provenance.json`
- `provenance.md`
- `sbom.json`
- `summary.json`
- `summary.md`

The published `SHA256SUMS.txt` matched all listed downloaded assets. The
source archive checksum was:

`dedf67537b1ed90c24921db32f0df7770aa42968c2d7cbe4927ec9a49a110e6f`

## Protected-Path Archive Scans

Protected-path pattern:

```text
(^|/)docs/evidence/(captured|consumer-submissions/(status\.json|current|artifacts|packets))
```

| Archive | SHA-256 | Protected-path hits |
| --- | --- | ---: |
| Uploaded `open-transit-rt-v0.1.0-rc.1.source.tar.gz` | `dedf67537b1ed90c24921db32f0df7770aa42968c2d7cbe4927ec9a49a110e6f` | 0 |
| GitHub-generated `v0.1.0-rc.1.tar.gz` | `b5ed4b6112b3f3e0d16cab3b38b5e04d68e2649eb96e4e7bcee8c3ff17092049` | 0 |
| GitHub-generated `v0.1.0-rc.1.zip` | `7a5dcb21d314fb5ae457d70db7e976a6fcb2df4e69759d4078335005d83d729c` | 0 |

The GitHub-generated archive SHA-256 values are recorded separately because
they are not listed in the uploaded `SHA256SUMS.txt`.

## Extraction Replay

Uploaded source archive extraction:

- Extracted to:
  `.cache/release-download/v0.1.0-rc.1/extract-uploaded-source/open-transit-rt-v0.1.0-rc.1`
- `make check`: failed.
- Failure:
  `FileNotFoundError: docs/evidence/consumer-submissions/status.json`

GitHub-generated tag tar extraction:

- Extracted to:
  `.cache/release-download/v0.1.0-rc.1/extract-github-tar/open-transit-rt-0.1.0-rc.1`
- `make check`: failed.
- Failure:
  `FileNotFoundError: docs/evidence/consumer-submissions/status.json`

The failure is a repository check bug for exported public source archives, not
a protected-path leak. The protected file is absent because the release archive
correctly excludes protected consumer-submission state.

## Phase 116 Patch

The current repository was patched after the rc1 replay result:

- Added `scripts/check-consumer-tracker.sh`.
- Updated `make check` to call the helper instead of inline Python.
- Updated claim/product/package audit scripts to skip protected consumer
  tracker and consumer artifact checks only when all of these are true:
  the source tree is outside a `.git` checkout, the protected file/directory
  is missing, and `.gitattributes` records the matching `export-ignore`
  policy.
- In normal repository checkouts, the exact seven prepared-only consumer
  tracker remains mandatory.

Patch validation:

- Current repository `scripts/check-consumer-tracker.sh`: passed with exact
  seven prepared-only targets.
- Current repository `make check`: passed.
- Current repository `make audit-product-acceptance`: passed.
- Current repository `make audit-final-claim-review`: passed.
- Current repository `make test-release-package`: passed.
- Export-like copy without protected paths: `make check` passed.
- Export-like copy without protected paths: `scripts/bootstrap-dev.sh --check`
  passed.
- Export-like copy with downloaded package copied to
  `.cache/release-package/v0.1.0-rc.1`: `scripts/audit-release-package.sh`
  passed.

## Claim And Boundary Notes

- No protected evidence path was edited, generated, reformatted, or touched.
- `docs/evidence/consumer-submissions/status.json` was not edited.
- The consumer tracker remains exactly seven targets in order and all
  `prepared`.
- The release download replay does not prove production readiness,
  compliance, adoption, consumer acceptance, hosted service availability,
  vendor compatibility, hardware certification, SLA/uptime, or ETA quality.
- Phase 117 should use a public fresh clone at the tag for independent install
  confidence, because the already published rc1 source archive has the
  `make check` replay blocker described above.
