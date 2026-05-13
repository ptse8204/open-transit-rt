# Open Transit RT — Consumer-Grade Control Plane Roadmap Pack

Generated: 2026-05-13
Project: `ptse8204/open-transit-rt`
Scope: future product roadmap only; no implementation starts automatically.

## Purpose

This pack is a Codex-ready roadmap for turning Open Transit RT from a capable backend with a functional operations console into a consumer-grade, browser-first control plane for small transit agencies, civic technologists, and developer integrators.

The roadmap is not just visual polish. It expands the product toward the original mission:

> Help small or resource-constrained agencies publish trustworthy GTFS and GTFS Realtime feeds, connect devices and external systems safely, monitor readiness, operate without vendor lock-in, and understand what remains before public deployment or stronger claims.

## Use order

1. `00-CODEX-READ-ME-FIRST.md`
2. `01-roadmap-overview.md`
3. `02-phases-and-checkpoints.md`
4. `03-master-subagent-operating-manual.md`
5. `04-validation-and-claim-boundaries.md`
6. relevant `phase-prompts/phase-XX-*.md`
7. relevant `audit-prompts/*.md`

## Current status assumption

This pack assumes:

- Phases 0-74 are closed as of the Phase 74 CP000008 closeout.
- Phase 72 remains `needs_review`, not release-ready.
- No new phase starts automatically.
- Phase 68+ evidence tracks remain closed unless separately authorized.
- Protected evidence paths and consumer tracker records remain untouched.
- All seven consumer targets remain exactly `prepared`.

## Proposed next phase family

Phase 75 is authorized only for this roadmap pack. Phase 76+ requires separate
maintainer authorization before implementation starts.

This pack proposes a future phase family:

```text
Phase 75+ — Consumer-Grade Control Plane And Operator Workflows
```

The first phase should only add this roadmap pack and, where needed, link it
from source-of-truth status docs with explicit planning-only wording. It should
not implement UI, APIs, migrations, evidence collection, or consumer tracker
changes.

## Commit rule

Every phase must end with a git commit. A phase is not closed until there is a final closeout commit and a master-agent checkpoint report.

Use this pattern:

```text
Phase XX -- Checkpoint 000001: add <phase> plan
Phase XX -- Checkpoint 000002: implement <major scope>
Phase XX -- Checkpoint 000003: close <phase> review
```

For longer phases, use more checkpoints. Do not mix phase scopes in one commit.
