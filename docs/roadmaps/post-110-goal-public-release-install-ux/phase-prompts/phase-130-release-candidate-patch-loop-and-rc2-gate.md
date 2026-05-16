# Phase 130 Prompt — Release Candidate Patch Loop And rc2 Gate

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Run a release candidate patch loop and decide whether an rc2 gate is needed;
prepare rc2 only if repo-caused blockers require it.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 130 -- Checkpoint 000001: add release candidate patch loop and rc2 gate plan
Phase 130 -- Checkpoint 000002: implement or audit primary scoped work
Phase 130 -- Checkpoint 000003: run validation and patch required gaps
Phase 130 -- Checkpoint 000004: close release candidate patch loop and rc2 gate review
```

## Required handoff

```text
docs/handoffs/phase-130.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
