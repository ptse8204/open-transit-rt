# Phase 117 -- Independent Public Install Confidence Trial

## Goal

Run independent public install confidence from a fresh public clone and/or
release archive, then patch documentation or UX blockers found during replay.

Phase 117 validates local install/evaluation mechanics only. It does not
claim stable release readiness, production readiness, compliance, adoption,
consumer acceptance, final-root readiness, hosted service availability,
SLA/uptime, vendor compatibility, hardware certification, or production-grade
ETA quality.

## Current Repo Context

- Phase 115 published public prerelease `v0.1.0-rc.1`.
- Phase 116 verified public release assets and checksums, and found that
  extracted published rc1 source archives fail `make check` because the
  protected consumer tracker is export-ignored while the rc1 tag still
  requires it.
- Phase 116 patched the current repo so future exported source archives can run
  lightweight checks without restoring protected evidence.
- A public git clone of the rc1 tag should include the protected tracker file,
  so it is the primary Phase 117 independent install confidence path.

## Scope

- Run `scripts/install-confidence.sh` against the public GitHub repository and
  tag `v0.1.0-rc.1`.
- Include local app startup and five local public feed fetches if the local
  Docker environment permits it.
- Run `make validate` and `make test` in the cloned tree if feasible.
- Optionally run a current-branch fresh clone replay after the Phase 116 patch
  to confirm future-source install confidence.
- Record exact output directories, commit SHAs, statuses, and blockers.
- Patch docs or scripts only for repo-caused blockers discovered in replay.

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

Primary public tag replay:

```bash
INSTALL_CONFIDENCE_MODE=clone \
INSTALL_CONFIDENCE_SOURCE=https://github.com/ptse8204/open-transit-rt.git \
INSTALL_CONFIDENCE_REF=v0.1.0-rc.1 \
INSTALL_CONFIDENCE_RUN_LOCAL_APP=true \
INSTALL_CONFIDENCE_RUN_VALIDATE=true \
INSTALL_CONFIDENCE_RUN_TEST=true \
scripts/install-confidence.sh
```

Current source replay:

```bash
INSTALL_CONFIDENCE_MODE=clone \
INSTALL_CONFIDENCE_SOURCE=/Users/edwintse/Downloads/open-transit-rt \
INSTALL_CONFIDENCE_REF=HEAD \
INSTALL_CONFIDENCE_RUN_LOCAL_APP=false \
INSTALL_CONFIDENCE_RUN_VALIDATE=true \
INSTALL_CONFIDENCE_RUN_TEST=true \
scripts/install-confidence.sh
```

Closeout validation:

- `make check`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- exact prepared-only consumer tracker assertion
- protected-path git status check

## Deliverables

- `docs/public-install-confidence-v0.1.0-rc.1.md`
- `docs/handoffs/phase-117.md`
- Source-of-truth status updates for Phase 117 closeout.
- Ignored local diagnostics under `.cache/install-confidence`.

## Implementation Plan

1. Add this Phase 117 plan and commit checkpoint 000001.
2. Run the public tag fresh-clone install-confidence harness and record exact
   diagnostics.
3. Run current-source replay if needed to verify Phase 116 patch behavior and
   patch any repo-caused install blockers.
4. Close Phase 117 with handoff/status docs and continue immediately to
   Phase 118.

## Checkpoint Plan

- `Phase 117 -- Checkpoint 000001: add independent public install confidence trial plan`
- `Phase 117 -- Checkpoint 000002: implement or audit primary scoped work`
- `Phase 117 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 117 -- Checkpoint 000004: close independent public install confidence trial review`

## Checkpoint Report -- 000001

Checkpoint:
Phase 117 -- Checkpoint 000001: add independent public install confidence
trial plan.

Goal status:
Active. Phase 116 is closed and Phase 117 has started.

Sub-agents used or simulated:
Install Confidence sub-agent use was attempted, but the environment reported
the sub-agent thread limit was reached. Install Confidence, Planning, QA,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, Release,
Web Design Skill, and GTFS-RT Domain roles are simulated by the Master Agent
for this plan checkpoint.

Changed files:
`docs/phase-117-independent-public-install-confidence-trial.md`.

Validation run:
Initial Phase 117 inspection reviewed the Phase 116 handoff, release download
replay report, install-confidence harness, and roadmap phase definition.
Focused checkpoint validation is scheduled before commit.

Blocked checks:
Fresh public clone replay, local app/five-feed fetches, validate/test replay,
current-source replay, and closeout validation are scheduled for later Phase
117 checkpoints.

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
The plan uses public clone/install diagnostics and does not change route auth,
CSRF, credential handling, token handling, public exposure, private payload
handling, or operator command behavior.

Data/migration status:
No migration, schema, durable state, public feed contract, or Go module change
is planned.

Release/publication status:
Phase 117 does not create or modify a release. It verifies public install
confidence from the already published rc1 tag and current source as needed.

Install confidence status:
Phase 113 local install confidence passed. Phase 116 public archive replay was
partial due a published rc1 source-archive `make check` blocker. Phase 117
will run public fresh-clone install confidence.

Web design skill status:
Phase 114 Web Design Skill artifact is complete. Phase 117 does not touch UX.

Master review:
Approved. The plan targets independent public fresh-clone install confidence
while preserving the documented rc1 release-archive blocker.

Required edits:
Run checkpoint 000001 validation and commit, then execute the public clone
trial.

Decision:
Proceed to checkpoint 000001 validation and commit.

Next checkpoint:
Phase 117 -- Checkpoint 000002: implement or audit primary scoped work.
