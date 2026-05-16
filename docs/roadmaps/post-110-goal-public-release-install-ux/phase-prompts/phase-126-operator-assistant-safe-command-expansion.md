# Phase 126 Prompt — Operator Assistant Safe Command Expansion

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Expand safe command-backed operator assistance through `internal/admincontrol`
without arbitrary shell execution or destructive browser actions.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 126 -- Checkpoint 000001: add operator assistant safe command expansion plan
Phase 126 -- Checkpoint 000002: implement or audit primary scoped work
Phase 126 -- Checkpoint 000003: run validation and patch required gaps
Phase 126 -- Checkpoint 000004: close operator assistant safe command expansion review
```

## Required handoff

```text
docs/handoffs/phase-126.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
