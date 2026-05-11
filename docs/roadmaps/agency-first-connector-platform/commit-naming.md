# Commit Naming Convention For Phase 61+

The maintainer has chosen to continue the original phase/checkpoint style using **Phase 61, Phase 62, Phase 63, ...**.

## Rule

Each phase has its own local checkpoint sequence.

```text
Phase XX -- Checkpoint 000001: <short outcome>
Phase XX -- Checkpoint 000002: <short outcome>
Phase XX -- Checkpoint 000003: <short outcome>
```

Do not continue the old global Post-60 checkpoint count.

## Normal Phase Lifecycle

A typical phase should have:

```text
Phase XX -- Checkpoint 000001: add <phase name> plan
Phase XX -- Checkpoint 000002: implement <main scope>
Phase XX -- Checkpoint 000003: close <phase name>
```

If the phase is small and the roadmap already contains a plan, Codex may start with implementation:

```text
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
```

## Patch After Close

If a completed phase needs a fix, use the next checkpoint in that same phase.

```text
Phase 61 -- Checkpoint 000004: fix connector hub audit gaps
Phase 61 -- Checkpoint 000005: clarify connector hub claim boundary
```

Do not create Phase 62 for a Phase 61 bug fix unless the maintainer approves a new phase scope.

## First Commit For This Roadmap Pack

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
```

This commit should add this directory and small navigation links from source-of-truth docs.

## Examples

```text
Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan
Phase 62 -- Checkpoint 000002: implement guided setup wizard v1
Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow
Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import

Phase 64 -- Checkpoint 000001: add connector platform and SDK plan
Phase 64 -- Checkpoint 000002: implement connector manifest registry UI
Phase 64 -- Checkpoint 000003: implement connector test runner UI
Phase 64 -- Checkpoint 000004: improve telemetry connector SDK examples
Phase 64 -- Checkpoint 000005: improve prediction connector SDK examples
Phase 64 -- Checkpoint 000006: close connector platform and SDKs
```
