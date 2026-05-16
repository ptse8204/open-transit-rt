# CODEX KICKOFF PROMPT — Goal-Based Public Release, Install Confidence, UX, and GTFS-RT Adoption Roadmap

You are the **Master Agent** for Open Transit RT.

This prompt must be used with the Goal in `goal/CODEX-GOAL.md`.

## Goal command

Paste this first:

```text
/goal Complete the Open Transit RT Phase 111-132 public release, independent install confidence, UX validation, and GTFS-RT adoption roadmap, verified by committed phase handoffs through docs/handoffs/phase-132.md, a release status artifact that either proves a published v0.1.0-rc.1 GitHub Release or records exact blocked-publication credentials/tooling evidence, independent fresh-clone or release-archive install reports, web-design-skill UX review artifacts, and passing validation/claim-boundary/protected-path checks. Preserve protected evidence paths, prepared-only consumer statuses, auth/security boundaries, and forbidden claim boundaries. Use the repository, local release tooling, GitHub release tooling only when authenticated and authorized, synthetic/local fixtures, and the web-design-engineer skill at ~/.agents/skills/web-design-engineer for UX phases. Between phases, commit every checkpoint, update handoffs/status docs, rerun the relevant validation, inspect blockers, fix repo-caused failures, and continue to the next phase until Phase 132 is closed. If a hard blocker prevents publication or a tool/credential is unavailable, record the exact evidence, close that phase truthfully as blocked or needs_review, continue safe downstream phases, and report what would unlock completion without faking release, evidence, acceptance, compliance, production readiness, SLA, vendor compatibility, or ETA-quality claims.
```

## Non-stop instruction

Do not stop after planning, Phase 111, a release recommendation, a UX pass, a
publication attempt, or any intermediate phase. Continue until Phase 132 is
closed and committed unless a hard safety blocker prevents safe continuation.

## Primary outcomes

1. Publish an actual public `v0.1.0-rc.1` release candidate if gates pass and
   release credentials/tooling are available.
2. Establish independent install confidence from fresh clone and/or release
   archive.
3. Use the Codex Web Design Skill:

```text
web-design-engineer ~/.agents/skills/web-design-engineer
```

4. Continue closing open-source GTFS-RT gaps with feed usefulness,
   interoperability, fixtures, connectors, realtime QA, alerts operations,
   safe operator commands, small-host deployability, and adoption support.

## Read first

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-95.md
docs/handoffs/phase-108.md
docs/handoffs/phase-110.md
docs/roadmap-status.md
docs/extension-governance.md
docs/release-notes-v0.1.0-rc.1-draft.md
docs/evidence/redaction-policy.md
docs/evidence/consumer-submissions/status.json
```

Also inspect relevant source:

```text
cmd/agency-config/
internal/admincontrol/
internal/compliance/
internal/gtfs/
internal/feed/
internal/prediction/
internal/realtimequality/
internal/connectors/
scripts/
Makefile
```

## Place roadmap pack

Add or reconcile this pack under:

```text
docs/roadmaps/post-110-goal-public-release-install-ux/
```

Then continue to Phase 112. Do not stop after adding the pack.

## Required workflow

Use the Master/sub-agent workflow in `03-master-subagent-operating-manual.md`.

## Hard boundaries

Protected paths must not be touched:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

All seven consumer targets must remain exactly `prepared`.

## Public release authorization

Phase 115 authorizes actual public `v0.1.0-rc.1` release candidate publication
if gates pass.

It does not authorize stable release, production readiness, compliance,
adoption, consumer acceptance, hosted SaaS, SLA/uptime, vendor compatibility,
hardware certification, final-root readiness, or production-grade ETA quality.

## Phase sequence

Run every phase:

- Phase 111 — Goal Activation And Public Release Roadmap Pack
- Phase 112 — Public Release Artifact And Claim Blocking Audit
- Phase 113 — Fresh Clone Install Harness And Release Dry Run
- Phase 114 — Web Design Skill UX Audit And Control Plane Polish
- Phase 115 — v0.1.0-rc.1 Public Release Cut
- Phase 116 — Published Release Verification And Download Replay
- Phase 117 — Independent Public Install Confidence Trial
- Phase 118 — Post-Release Web Design Skill UX Validation
- Phase 119 — Public Docs Site README And Quickstart Release Alignment
- Phase 120 — GTFS-RT Feed Usefulness And Reliability V2
- Phase 121 — GTFS-RT Interoperability Conformance Harness
- Phase 122 — GTFS-RT Fixture Library And Edge-Case Coverage
- Phase 123 — Vehicle AVL Connector Starter Kits
- Phase 124 — Realtime QA ETA Backtesting And Prediction Confidence V3
- Phase 125 — Alerts And Service Disruption Operations V2
- Phase 126 — Operator Assistant Safe Command Expansion
- Phase 127 — Small-Host Deployment And Upgrade UX Hardening
- Phase 128 — Contributor And Agency Evaluator Adoption Kit
- Phase 129 — Community Support Feedback And Issue Triage Kit
- Phase 130 — Release Candidate Patch Loop And rc2 Gate
- Phase 131 — Optional Evidence Gate Refresh Blocker-Only
- Phase 132 — Final Public Release Install UX Roadmap Closeout

## Validation

Use `04-validation-and-claim-boundaries.md`.

## Completion

Only after Phase 132 closeout commit, output:

```text
Goal status:
Completed phases:
Commit list by phase:
Public release status:
Independent install confidence status:
Web Design Skill UX status:
GTFS-RT gap improvements:
Validation summary:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Remaining recommended next steps:
```
