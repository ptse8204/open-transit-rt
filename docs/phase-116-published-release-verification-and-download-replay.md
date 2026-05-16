# Phase 116 -- Published Release Verification And Download Replay

## Goal

Verify the published `v0.1.0-rc.1` tag and GitHub Release can be fetched,
checksummed, unpacked, protected-path scanned, and used for bounded local
evaluation.

Phase 116 verifies publication mechanics and download replay only. It does not
claim stable release readiness, production readiness, compliance, adoption,
consumer acceptance, final-root readiness, hosted service availability,
SLA/uptime, vendor compatibility, hardware certification, or production-grade
ETA quality.

## Current Repo Context

- Phase 115 published the public GitHub prerelease:
  `https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`.
- The annotated tag `v0.1.0-rc.1` dereferences to
  `497f99a97baff630af147c83a7e1249bb08e32da`.
- The uploaded local package source archive SHA-256 is
  `dedf67537b1ed90c24921db32f0df7770aa42968c2d7cbe4927ec9a49a110e6f`.
- Phase 115 verified the local package source archive had zero protected-path
  entries after the root `.gitattributes` `export-ignore` policy.
- Phase 116 must independently verify downloaded release assets and
  GitHub-generated tag archives rather than relying only on local `.cache`
  artifacts.

## Scope

- Download published release assets through GitHub release tooling.
- Verify downloaded checksums against the published `SHA256SUMS.txt` asset.
- Verify downloaded release metadata remains draft `false`, prerelease `true`,
  and points at the expected tag.
- Scan the downloaded uploaded source archive for protected evidence and
  consumer-submission paths.
- Download GitHub-generated tag archives and scan them for the same protected
  paths.
- Extract downloaded archives in ignored `.cache` paths and run bounded local
  replay checks such as `make check`, bootstrap preflight, and package audit
  where the archive structure supports it.
- Record a release download replay report and closeout handoff.

## Protected Paths

Do not modify, reformat, delete, stage, or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

The consumer tracker must remain exactly seven targets in order and all
`prepared`.

## Verification Checklist

- `gh release view v0.1.0-rc.1 --repo ptse8204/open-transit-rt ...`
- `gh release download v0.1.0-rc.1 --repo ptse8204/open-transit-rt --dir .cache/release-download/v0.1.0-rc.1 --clobber`
- `sha256sum -c SHA256SUMS.txt` from the download directory for listed assets.
- Uploaded source archive protected-path scan returns `0`.
- GitHub-generated `tar.gz` tag archive protected-path scan returns `0`.
- GitHub-generated `zip` tag archive protected-path scan returns `0`.
- Extracted uploaded source archive runs `make check`.
- Extracted uploaded source archive runs bootstrap preflight where supported.
- Reconstructed or local package audit evidence is recorded without claiming
  production readiness or consumer acceptance.
- `make check`, claim audits, prepared-only consumer assertion, and
  protected-path status checks pass before each checkpoint commit.

## Deliverables

- `docs/release-download-replay-v0.1.0-rc.1.md`
- `docs/handoffs/phase-116.md`
- Source-of-truth status updates for Phase 116 closeout.
- Ignored local download/extract artifacts under `.cache`.

## Implementation Plan

1. Add this Phase 116 plan and commit checkpoint 000001.
2. Download and verify published release assets and GitHub-generated archives.
3. Run extraction/replay checks, patch repo-caused blockers if any, and record
   exact evidence.
4. Close Phase 116 with handoff/status docs and continue immediately to
   Phase 117.

## Checkpoint Plan

- `Phase 116 -- Checkpoint 000001: add published release verification and download replay plan`
- `Phase 116 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 116 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 116 -- Checkpoint 000004: close published release verification and download replay review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 116 -- Checkpoint 000001: add published release verification and
download replay plan.

Goal status:
Active. Phase 115 is closed and Phase 116 has started.

Sub-agents used or simulated:
Release/Install Confidence sub-agent review was delegated for independent
download replay checklist recommendations. Planning, Implementation, QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Web Design
Skill, and GTFS-RT Domain roles are simulated by the Master Agent for this
plan checkpoint.

Changed files:
`docs/phase-116-published-release-verification-and-download-replay.md`.

Validation run:
Initial Phase 116 inspection reviewed the Phase 115 handoff, release status
artifact, roadmap phase definition, and validation/claim-boundary rules.
Focused checkpoint validation is scheduled before commit.

Blocked checks:
Release asset download, checksum verification, uploaded source archive scan,
GitHub-generated archive scan, extraction replay, package audit replay, and
full closeout validation are scheduled for later Phase 116 checkpoints.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected
path writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and exactly `prepared`.

Claim-boundary status:
The plan explicitly forbids stable release, production readiness, compliance,
adoption, consumer acceptance, final-root readiness, hosted service, paid
support, SLA/uptime, vendor compatibility, hardware certification, production
AVL reliability, and ETA-quality claims.

Security/auth status:
The plan uses existing GitHub release download tooling and does not change
application route auth, CSRF, credential handling, token handling, public
exposure, private payload handling, or operator command behavior.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
is planned.

Release/publication status:
Phase 115 already published the prerelease. Phase 116 verifies public download
replay and does not create a new release or tag.

Install confidence status:
Phase 113 remains the current local install-confidence result until Phase 116
and Phase 117 add public download/archive replay evidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 116 does not touch UX.

Master review:
Approved. The plan targets the remaining public-release verification gap
without weakening protected-path, consumer-tracker, auth, or claim boundaries.

Required edits:
Run checkpoint 000001 validation and commit, then download and verify release
assets.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 116 -- Checkpoint 000002: implement or audit primary scoped work.

## Checkpoint Report -- 000002

Checkpoint:
Phase 116 -- Checkpoint 000002: implement or audit primary scoped work.

Goal status:
Active. Published release assets and GitHub-generated archives were
downloaded and scanned; a release-archive `make check` blocker was found and
patched in the current repository for future archives.

Sub-agents used or simulated:
Release/Install Confidence sub-agent independently recommended Phase 116
commands and risks. Implementation, QA, Documentation / IA, Claim-Boundary,
Security/Auth, Data/Migration, Web Design Skill, and GTFS-RT Domain roles
were simulated by the Master Agent.

Changed files:
`Makefile`; `scripts/check-consumer-tracker.sh`;
`scripts/audit-product-acceptance.sh`;
`scripts/audit-final-claim-review.sh`; `scripts/audit-release-package.sh`;
`docs/release-download-replay-v0.1.0-rc.1.md`;
`docs/phase-116-published-release-verification-and-download-replay.md`.

Validation run:
Passed:

- `gh release download v0.1.0-rc.1 --repo ptse8204/open-transit-rt --dir .cache/release-download/v0.1.0-rc.1 --clobber`
- published `SHA256SUMS.txt` verification for all downloaded uploaded assets
- uploaded source archive protected-path scan: `0`
- GitHub-generated tag `tar.gz` protected-path scan: `0`
- GitHub-generated tag `zip` protected-path scan: `0`
- current repository `scripts/check-consumer-tracker.sh`
- current repository `make check`
- current repository `make audit-product-acceptance`
- current repository `make audit-final-claim-review`
- current repository `make test-release-package`
- export-like copy without protected paths: `make check`
- export-like copy without protected paths: `scripts/bootstrap-dev.sh --check`
- export-like copy with downloaded package under `.cache/release-package/v0.1.0-rc.1`: `scripts/audit-release-package.sh`

Blocked checks:
Published rc1 uploaded source archive and GitHub-generated tag archive
extraction replay both failed `make check` before the Phase 116 patch because
the protected consumer tracker is correctly excluded from public archives but
the published rc1 `make check` still requires it. This is recorded in
`docs/release-download-replay-v0.1.0-rc.1.md` and cannot be fixed inside the
already published rc1 archive.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.
Published archive scans returned zero protected-path hits.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared` in the repository checkout. Exported archives skip
the tracker check only when the protected tracker is missing outside a `.git`
checkout and the committed `.gitattributes` export policy documents that
absence.

Claim-boundary status:
The replay report records publication and download mechanics only. It does not
claim stable release, production readiness, compliance, adoption, consumer
acceptance, final-root readiness, hosted service, paid support, SLA/uptime,
vendor compatibility, hardware certification, production AVL reliability, or
ETA quality.

Security/auth status:
No application auth, CSRF, credential handling, token handling, private
payload handling, route exposure, or operator command behavior changed. The
new helper reads only local files.

Data/migration status:
No migration, schema, durable state, public feed contract, dependency, or Go
module change was added.

Release/publication status:
The public rc1 prerelease remains published. No new tag, release, or asset
upload was created in checkpoint 000002.

Install confidence status:
Published asset download and checksum replay passed, but published rc1 source
archive `make check` replay is blocked by the documented protected-tracker
check bug. Phase 117 should use a public fresh clone for independent install
confidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 116 does not touch UX.

Master review:
Approved. The checkpoint records the public download facts truthfully and
patches the current repo without weakening protected-path or prepared-only
consumer checks in normal repository checkouts.

Required edits:
Commit checkpoint 000002, then run broader validation and decide whether any
additional archive replay patching is needed.

Decision:
Proceed to checkpoint 000002 commit and checkpoint 000003 validation.

Next checkpoint:
Phase 116 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 116 -- Checkpoint 000003: run validation and patch required gaps.

Goal status:
Active. The Phase 116 patch was validated in the normal repository checkout
and in a post-patch exported source archive.

Sub-agents used or simulated:
Release/Install Confidence sub-agent recommendations were incorporated.
QA, Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration,
Implementation, Web Design Skill, and GTFS-RT Domain roles were simulated by
the Master Agent.

Changed files:
`docs/release-download-replay-v0.1.0-rc.1.md`;
`docs/phase-116-published-release-verification-and-download-replay.md`.

Validation run:
Passed:

- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make test-release-package`
- post-patch `git archive HEAD` protected-path scan: `0`
- post-patch `git archive HEAD` extracted tree `scripts/check-consumer-tracker.sh`: skipped only because the protected tracker is export-ignored from source archive output
- post-patch `git archive HEAD` extracted tree `make check`
- post-patch `git archive HEAD` extracted tree `scripts/bootstrap-dev.sh --check`
- post-patch `git archive HEAD` extracted tree with downloaded rc1 package copied to `.cache/release-package/v0.1.0-rc.1`: `scripts/audit-release-package.sh`

Blocked checks:
The already published rc1 source archives still fail `make check` because the
fix is not present in the published tag. This remains a truthful Phase 116
replay blocker and is recorded in
`docs/release-download-replay-v0.1.0-rc.1.md`.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.
The current post-patch source archive scan returned zero protected-path hits.

Consumer tracker status:
The tracker was not edited. The repository checkout still requires the exact
seven prepared-only targets. Exported archive trees skip only the missing
protected tracker that `.gitattributes` explicitly excludes.

Claim-boundary status:
Validation records package/download/archive mechanics and a repo-caused
archive replay bug only. It does not claim stable release, production
readiness, compliance, adoption, consumer acceptance, final-root readiness,
hosted service, paid support, SLA/uptime, vendor compatibility, hardware
certification, production AVL reliability, or ETA quality.

Security/auth status:
No application auth, CSRF, credential handling, token handling, private
payload handling, route exposure, or operator command behavior changed.

Data/migration status:
No migration, schema, durable state, public feed contract, dependency, or Go
module change was added.

Release/publication status:
The public rc1 prerelease remains published. No new tag, release, or asset
upload was created in checkpoint 000003.

Install confidence status:
Published rc1 release archive replay remains blocked for `make check`, while
post-patch current source archive replay passes the lightweight check and
bootstrap preflight. Phase 117 should run public fresh-clone install
confidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 116 does not touch UX.

Master review:
Approved. The validation confirms the patch keeps prepared-only enforcement
in repository checkouts while making exported future archives usable without
protected evidence files.

Required edits:
Commit checkpoint 000003, then close Phase 116 with handoff/status docs.

Decision:
Proceed to checkpoint 000003 commit and Phase 116 closeout.

Next checkpoint:
Phase 116 -- Checkpoint 000004: close published release verification and
download replay review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 116 -- Checkpoint 000004: close published release verification and
download replay review.

Goal status:
Active. Phase 116 is closed and Phase 117 starts next.

Sub-agents used or simulated:
Release/Install Confidence sub-agent recommendations were incorporated. QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration,
Implementation, Web Design Skill, GTFS-RT Domain, and Planning closeout roles
were simulated by the Master Agent.

Changed files:
`docs/handoffs/phase-116.md`; `docs/handoffs/latest.md`;
`docs/current-status.md`; `docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-116-published-release-verification-and-download-replay.md`.

Validation run:
Closeout relies on checkpoint 000003 full validation. After closeout docs were
updated, focused docs/protected-path validation is rerun before the checkpoint
000004 commit.

Blocked checks:
Published rc1 source archive `make check` replay remains blocked because the
already published tag predates the archive-aware protected tracker check. This
is recorded truthfully in the replay report.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
The tracker was not edited. The exact seven consumer targets remain in order
and all remain `prepared`.

Claim-boundary status:
Phase 116 records download and archive verification only. It makes no stable,
production, compliance, adoption, consumer, final-root, hosted-service,
vendor, SLA, hardware, or ETA-quality claim.

Security/auth status:
No application route auth, CSRF behavior, credential handling, token handling,
private payload handling, public exposure, or operator command behavior
changed.

Data/migration status:
No migration, schema, durable state, public feed contract, dependency, or Go
module change was added.

Release/publication status:
The public rc1 prerelease remains published. No new tag, release, or asset
upload was created.

Install confidence status:
Public asset download and checksum verification passed. Published rc1 source
archive install confidence remains partial because extracted archives fail
`make check`; Phase 117 is required for independent public fresh-clone install
confidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 118 remains scheduled.

Master review:
Approved. Phase 116 closes with bounded release download verification, exact
blocker evidence for rc1 archive replay, and a current-source patch for future
archives.

Required edits:
Commit checkpoint 000004, then continue directly to Phase 117.

Decision:
Phase 116 is complete.

Next checkpoint:
Phase 117 -- Checkpoint 000001: add independent public install confidence
trial plan.
