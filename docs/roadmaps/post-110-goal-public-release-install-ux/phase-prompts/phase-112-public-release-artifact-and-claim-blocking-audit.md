# Phase 112 Prompt — Public Release Artifact And Claim Blocking Audit

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Audit the local rc1 package, release notes, source archive contents, claim
boundaries, protected paths, and public-distribution readiness before release.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 112 -- Checkpoint 000001: add public release artifact and claim blocking audit plan
Phase 112 -- Checkpoint 000002: implement or audit primary scoped work
Phase 112 -- Checkpoint 000003: run validation and patch required gaps
Phase 112 -- Checkpoint 000004: close public release artifact and claim blocking audit review
```

## Required handoff

```text
docs/handoffs/phase-112.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
