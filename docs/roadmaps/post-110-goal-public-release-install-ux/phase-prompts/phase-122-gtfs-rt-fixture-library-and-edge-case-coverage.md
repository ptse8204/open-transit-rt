# Phase 122 Prompt — GTFS-RT Fixture Library And Edge-Case Coverage

Use this with:

- `goal/CODEX-GOAL.md`
- `CODEX-KICKOFF-AUTONOMOUS-PHASE-111-TO-132.md`

## Phase goal

Goal:
Expand GTFS/GTFS-RT fixtures and edge-case coverage for midnight rollover,
frequency service, canceled trips, stale telemetry, unknown vehicles, and
malformed realtime messages.

## Required workflow

Use the Master/sub-agent workflow. Include Web Design Skill review when this
phase touches UX, especially Phases 114 and 118.

## Minimum checkpoint commits

```text
Phase 122 -- Checkpoint 000001: add gtfs rt fixture library and edge case coverage plan
Phase 122 -- Checkpoint 000002: implement or audit primary scoped work
Phase 122 -- Checkpoint 000003: run validation and patch required gaps
Phase 122 -- Checkpoint 000004: close gtfs rt fixture library and edge case coverage review
```

## Required handoff

```text
docs/handoffs/phase-122.md
```

Also update source-of-truth status docs.

## Boundaries

Do not touch protected evidence paths. Do not move consumer statuses. Do not
make unsupported claims. If publication, browser automation, skill usage, or
release tooling is blocked, record exact evidence and continue safe downstream
phases.
