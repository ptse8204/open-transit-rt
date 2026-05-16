# Phase 116 Prompt — Published Release Verification And Download Replay

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Verify the published tag/release can be fetched, checksummed, unpacked, and
used for local evaluation.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 116 -- Checkpoint 000001: add published release verification and download replay plan
Phase 116 -- Checkpoint 000002: implement or audit primary scoped work
Phase 116 -- Checkpoint 000003: run validation and patch required gaps
Phase 116 -- Checkpoint 000004: close published release verification and download replay review
```

## Required handoff

```text
docs/handoffs/phase-116.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
