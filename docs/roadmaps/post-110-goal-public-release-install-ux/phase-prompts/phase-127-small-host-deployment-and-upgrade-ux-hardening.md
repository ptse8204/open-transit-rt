# Phase 127 Prompt — Small-Host Deployment And Upgrade UX Hardening

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Harden small-host deployment and upgrade UX: preflight, off-host validation,
backup/restore guidance, resource constraints, and recovery.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 127 -- Checkpoint 000001: add small host deployment and upgrade ux hardening plan
Phase 127 -- Checkpoint 000002: implement or audit primary scoped work
Phase 127 -- Checkpoint 000003: run validation and patch required gaps
Phase 127 -- Checkpoint 000004: close small host deployment and upgrade ux hardening review
```

## Required handoff

```text
docs/handoffs/phase-127.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
