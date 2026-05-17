# Release Download Replay -- `v0.1.0-rc.2`

Status date: 2026-05-16

This artifact records public download replay for Open Transit RT
`v0.1.0-rc.2`.

Open Transit RT v0.1.0-rc.2 is a public release candidate for local/self-hosted evaluation.

It is not a stable release and does not prove production readiness, compliance, agency adoption, consumer acceptance, final-root readiness, hosted service availability, vendor compatibility, hardware certification, SLA/uptime, production AVL reliability, production-grade ETA quality, or real-world ETA accuracy.

## Conclusion

`download_replay_passed`

The public GitHub Release exists, uploaded release assets downloaded,
downloaded asset checksums matched the published `SHA256SUMS.txt`, uploaded
and GitHub-generated archives had zero protected-path hits, and the downloaded
uploaded source archive passed extracted archive replay checks.

## Release Metadata

- Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.2`
- Tag: `v0.1.0-rc.2`
- GitHub Release draft: `false`
- GitHub Release prerelease: `true`
- Published at: `2026-05-16T22:32:45Z`
- Annotated tag object: `3b7e9b616d98dee908b48645108035117a68e5dc`
- Tag target commit: `15a0ec7cbdacf2301ac906ff3ecbe655371fccc6`
- Download directory: `.cache/release-download-replay/v0.1.0-rc.2`

## Downloaded Uploaded Assets

Downloaded with:

```bash
gh release download v0.1.0-rc.2 \
  --repo ptse8204/open-transit-rt \
  --dir .cache/release-download-replay/v0.1.0-rc.2
```

Downloaded assets:

- `SHA256SUMS.txt`
- `image.json`
- `manifest.json`
- `manifest.md`
- `open-transit-rt-v0.1.0-rc.2.source.tar.gz`
- `provenance.json`
- `provenance.md`
- `sbom.json`
- `summary.json`
- `summary.md`

The downloaded `SHA256SUMS.txt` matched all listed downloaded assets. The
downloaded `SHA256SUMS.txt` file checksum was:

`d76167e2e6434683afca5dc1d5809c53c5e5dea323a6070a6cc5ef0e2fee72df`

## Asset Checksums

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

GitHub-generated archive checksums recorded during replay:

| Archive | SHA-256 |
| --- | --- |
| `github-v0.1.0-rc.2.tar.gz` | `a1dfdc54921806a18366fa410cf170d6181fc40df2fad0aabfb4946f1f0546a1` |
| `github-v0.1.0-rc.2.zip` | `2d63979be72aff7eb862c7a89d166bff329a04ea54eefc2c988f2b7f9be0f1e7` |

## Protected-Path Archive Scans

Protected-path pattern set:

```text
docs/evidence/captured
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current
docs/evidence/consumer-submissions/artifacts
docs/evidence/consumer-submissions/packets
```

| File | Scan mode | Protected-path hits |
| --- | --- | ---: |
| Uploaded `open-transit-rt-v0.1.0-rc.2.source.tar.gz` | tar listing | 0 |
| GitHub-generated `github-v0.1.0-rc.2.tar.gz` | tar listing | 0 |
| GitHub-generated `github-v0.1.0-rc.2.zip` | zip listing | 0 |
| Uploaded metadata/checksum assets | text scan | 0 |

Total protected-path hits: `0`.

## Extracted Archive Replay

Downloaded uploaded source archive extraction:

- Extracted to:
  `.cache/release-download-replay/v0.1.0-rc.2/extract-uploaded-source/open-transit-rt-v0.1.0-rc.2`
- `make check`: passed.
- `scripts/bootstrap-dev.sh --check`: passed.
- `make validators-install`: passed.
- `make validate`: passed.
- `make test`: passed.

Pre-publication archive install-confidence from the generated source archive
also passed at `.cache/install-confidence/v0.1.0-rc.2-archive-rerun2`, including
local app startup and fetches for:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

## Protected And Consumer Status

- No protected evidence path was edited, generated, reformatted, or touched.
- `docs/evidence/consumer-submissions/status.json` was not edited.
- The consumer tracker remains exactly seven targets in order and all remain
  `prepared`.
- This replay did not contact consumers, submit feeds, or move any consumer
  status.
