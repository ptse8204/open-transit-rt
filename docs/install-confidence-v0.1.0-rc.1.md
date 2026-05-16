# v0.1.0-rc.1 Install Confidence Report

Status: `local_install_confidence_passed_public_release_blocked`

This report summarizes Phase 113 local install-confidence diagnostics for the
`v0.1.0-rc.1` release-candidate track. It records bounded local replay
signals only. It is not retained evidence, release publication, production
readiness, compliance proof, consumer acceptance, agency approval, hosted
service availability, vendor compatibility, SLA/uptime, hardware
certification, production AVL reliability, production-grade ETA quality, or
real-world ETA accuracy proof.

## Inputs

- Active repository: `/Users/edwintse/Downloads/open-transit-rt`
- Local fresh-clone ref: `HEAD`
- Local fresh-clone commit:
  `f7fe95de9503093233d5393418e17d10af23bd04`
- Local fresh-clone describe: `f7fe95d`
- Archive source:
  `.cache/release-package/v0.1.0-rc.1/artifacts/open-transit-rt-v0.1.0-rc.1.source.tar.gz`
- Archive SHA-256:
  `a43610b61eb2a405408b6a9aabfefeb5b61c2b7cfa97c2813f69847ce8ea3413`
- Local tool metadata:
  - Go: `go1.26.2 darwin/amd64`
  - Docker: `29.4.2`
  - Docker Compose: `v5.1.3`

Raw logs and downloaded feed artifacts are intentionally kept under ignored
`.cache/install-confidence/**` and are not committed.

## Fresh Clone Replay

- Output directory: `.cache/install-confidence/phase-113-local-clone`
- Generated at: `20260516T024013Z`
- Mode: `clone`
- Source: local repository path
- Overall status: `passed`
- Local app startup: enabled
- `make validate`: not checked in this harness run
- `make test`: not checked in this harness run

Passed steps:

- `git clone`
- `git checkout`
- `make check`
- `scripts/bootstrap-dev.sh --check`
- `make agency-app-up`
- anonymous local fetch of `/public/feeds.json`
- anonymous local fetch of `/public/gtfs/schedule.zip`
- anonymous local fetch of `/public/gtfsrt/vehicle_positions.pb`
- anonymous local fetch of `/public/gtfsrt/trip_updates.pb`
- anonymous local fetch of `/public/gtfsrt/alerts.pb`

Fetched local artifact checksums:

- `feeds_json`:
  `1c75275995f038a96b3bf3f5639cbea8a3796542d092c81cb21b47d5cfd45fe1`
- `schedule_zip`:
  `efba37dc27bea1d261c2c5ec72a52349b8c88f3bef1c81f7278a2a4d0983456e`
- `vehicle_positions_pb`:
  `97d9569e12163920662545518397b444b3412c6c550ff9cf2ce53d150d21fb25`
- `trip_updates_pb`:
  `97d9569e12163920662545518397b444b3412c6c550ff9cf2ce53d150d21fb25`
- `alerts_pb`:
  `ad68f55a75381808980b78c4bab5ced96cd47bb969f3563118107ccaec9596b2`

## Release Archive Replay

- Output directory: `.cache/install-confidence/phase-113-archive`
- Generated at: `20260516T024049Z`
- Mode: `archive`
- Source: local generated source archive
- Overall status: `passed`
- Local app startup: not checked in this archive replay
- `make validate`: not checked in this harness run
- `make test`: not checked in this harness run

Passed steps:

- archive listing
- archive extraction
- `make check`
- `scripts/bootstrap-dev.sh --check`

The archive replay verifies the Phase 113 archive-aware `make check` behavior:
outside a Git worktree, the whitespace check skips `git diff --check` instead
of failing solely because `.git` is absent.

## Validation Context

Phase 113 also ran and passed the broader validation needed for the
scripts/Makefile change:

- `make test-install-confidence`
- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- consumer tracker JSON parse
- exact prepared-only consumer tracker assertion
- protected-path status check

## Remaining Publication Blocker

Install confidence does not override the Phase 112 public-distribution
blocker. Publication remains blocked until the maintainer resolves the source
archive public-distribution review recorded in
`docs/release-status-v0.1.0-rc.1.md`.

