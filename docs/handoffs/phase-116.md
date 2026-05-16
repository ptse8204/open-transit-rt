# Phase 116 Handoff -- Published Release Verification And Download Replay

## Status

Phase 116 is complete for published `v0.1.0-rc.1` release verification and
download replay.

The public release exists, uploaded assets downloaded, checksum verification
passed, and both uploaded and GitHub-generated source archives had zero
protected-path hits. The published rc1 source archive replay is not fully
install-clean: extracted rc1 archives fail `make check` because the protected
consumer tracker is intentionally excluded from public archives while the rc1
tag still requires that file in `make check`.

The current repository was patched so future source archives can run
lightweight checks without restoring protected evidence, while normal git
checkouts continue to require the exact seven prepared-only consumer tracker.

Phase 116 did not create a new tag, publish a new release, upload assets, push
images, create retained evidence, contact external parties, move consumer
statuses, modify protected evidence paths, or make stronger public claims.

## Completed Checkpoints

- Phase 116 -- Checkpoint 000001: add published release verification and
  download replay plan.
- Phase 116 -- Checkpoint 000002: implement or audit primary scoped work.
- Phase 116 -- Checkpoint 000003: run validation and patch required gaps.
- Phase 116 -- Checkpoint 000004: close published release verification and
  download replay review.

## Release Replay Result

Release replay report:
`docs/release-download-replay-v0.1.0-rc.1.md`.

- Release URL:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`
- Release draft: `false`
- Release prerelease: `true`
- Tag target commit: `497f99a97baff630af147c83a7e1249bb08e32da`
- Download directory: `.cache/release-download/v0.1.0-rc.1`
- Uploaded asset checksum verification: passed.
- Uploaded source archive protected hits: `0`.
- GitHub-generated tag tar protected hits: `0`.
- GitHub-generated tag zip protected hits: `0`.
- Published rc1 extracted source archive `make check`: blocked by missing
  protected consumer tracker file.
- Current post-patch source archive `make check`: passed.

GitHub-generated archive SHA-256 values:

- `v0.1.0-rc.1.tar.gz`:
  `b5ed4b6112b3f3e0d16cab3b38b5e04d68e2649eb96e4e7bcee8c3ff17092049`
- `v0.1.0-rc.1.zip`:
  `7a5dcb21d314fb5ae457d70db7e976a6fcb2df4e69759d4078335005d83d729c`

## Changed Files

- `Makefile`
- `scripts/check-consumer-tracker.sh`
- `scripts/audit-product-acceptance.sh`
- `scripts/audit-final-claim-review.sh`
- `scripts/audit-release-package.sh`
- `docs/phase-116-published-release-verification-and-download-replay.md`
- `docs/release-download-replay-v0.1.0-rc.1.md`
- `docs/handoffs/phase-116.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

## Validation

Passed:

- `gh release download v0.1.0-rc.1 --repo ptse8204/open-transit-rt --dir .cache/release-download/v0.1.0-rc.1 --clobber`
- published `SHA256SUMS.txt` verification for all downloaded uploaded assets
- uploaded source archive protected-path scan: `0`
- GitHub-generated tag `tar.gz` protected-path scan: `0`
- GitHub-generated tag `zip` protected-path scan: `0`
- current repository `scripts/check-consumer-tracker.sh`
- current repository `make check`
- current repository `make validate`
- current repository `make test`
- `docker compose -f deploy/docker-compose.yml config`
- current repository `make test-release-package`
- current repository `make audit-product-acceptance`
- current repository `make audit-final-claim-review`
- post-patch `git archive HEAD` protected-path scan: `0`
- post-patch extracted source archive `make check`
- post-patch extracted source archive `scripts/bootstrap-dev.sh --check`
- post-patch extracted source archive with downloaded package copied under
  `.cache/release-package/v0.1.0-rc.1`: `scripts/audit-release-package.sh`
- `git diff --check`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Blocked:

- Published rc1 extracted source archive `make check` remains blocked by the
  protected consumer tracker exclusion. This is patched in the current repo but
  cannot be retroactively fixed inside the already published rc1 archives.

## Protected Path Status

No protected evidence path was edited, generated, reformatted, or touched.
Public archive scans returned zero protected-path hits.

## Consumer Tracker Status

`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven targets remain present in order and all remain `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim Boundary Status

Phase 116 makes no stable release readiness, production readiness, compliance,
adoption, agency approval, consumer acceptance, consumer
ingestion/listing/display, final-root readiness, hosted service availability,
paid support, SLA/uptime, vendor compatibility, hardware certification,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy claim.

## Security/Auth Status

No application route auth, CSRF behavior, credential handling, token handling,
public exposure, private payload handling, external contact, or operator
command behavior changed. The new consumer tracker helper reads local files
only.

## Data/Migration Status

No migration, schema, durable state, dependency, runtime dependency, public
feed contract, or Go module change was added.

## Release/Publication Status

The Phase 115 public `v0.1.0-rc.1` prerelease remains published. Phase 116 did
not publish a new tag or release.

## Install Confidence Status

Public release asset download and checksum confidence passed. Published rc1
source archive install confidence is partial because `make check` is blocked
in extracted archives. Phase 117 should run independent public fresh-clone
install confidence from the tag.

## Web Design Skill Status

Phase 114 Web Design Skill artifact remains complete. Phase 118 remains the
post-release Web Design Skill UX validation pass.

## Commit List

- `183799e` -- Phase 116 -- Checkpoint 000001: add published release
  verification and download replay plan
- `4d6c36a` -- Phase 116 -- Checkpoint 000002: implement or audit primary
  scoped work
- `b9a7216` -- Phase 116 -- Checkpoint 000003: run validation and patch
  required gaps
- Phase 116 -- Checkpoint 000004: close published release verification and
  download replay review

## Checkpoint Report

Checkpoint:
Phase 116 -- Checkpoint 000004: close published release verification and
download replay review.

Goal status:
Active. Phase 116 is closed and the goal continues to Phase 117.

Sub-agents used or simulated:
Release/Install Confidence sub-agent recommendations were incorporated. QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Web Design
Skill, GTFS-RT Domain, Planning, and Implementation closeout roles were
simulated by the Master Agent.

Changed files:
`docs/handoffs/phase-116.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-116-published-release-verification-and-download-replay.md`.

Validation run:
Phase 116 full validation passed before closeout docs. After closeout docs are
updated, focused docs/protected-path validation is rerun before the checkpoint
000004 commit.

Blocked checks:
Published rc1 source archive `make check` replay remains blocked for the
already published archive. Current-source future archive replay passes.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 116 records download replay and a release-archive check blocker only. It
makes no stronger public claim.

Security/auth status:
No application security behavior changed.

Data/migration status:
No migration, schema, durable state, dependency, public feed contract, or Go
module change was added.

Release/publication status:
The public rc1 prerelease remains published. No new release action was taken.

Install confidence status:
Public asset download checks passed. Phase 117 starts next for independent
public fresh-clone install confidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 118 remains scheduled.

Master review:
Approved. Phase 116 closes with truthful public download verification, a
bounded blocker for published archive replay, and a current-source fix.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 117.

Decision:
Proceed to checkpoint 000004 commit and continue to Phase 117.

Next checkpoint:
Phase 117 -- Checkpoint 000001: add independent public install confidence
trial plan.
