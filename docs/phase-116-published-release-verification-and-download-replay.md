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
