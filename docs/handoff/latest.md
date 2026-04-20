# Latest Handoff

This file tells the next Codex instance where the project stands right now.

## Active phase
Phase 0 — Scaffolding and repo hardening

## Phase status
- Not yet completed
- This is the current working phase unless a newer handoff file says otherwise

## Read these files first
1. `AGENTS.md`
2. `docs/phase-plan.md`
3. `docs/current-status.md`
4. `docs/codex-task.md`
5. `docs/repo-gaps.md`
6. `docs/dependencies.md`

## Current objective
Add the missing repo scaffolding and make the repository ready for reliable phased development.

## Highest-priority tasks right now
- add `.env.example`
- add `Taskfile.yml` or expand build tasks
- add `cmd/migrate`
- add versioned migrations
- add `scripts/bootstrap-dev.sh`
- add `testdata/` fixtures
- add `docs/decisions.md`
- ensure `docs/dependencies.md` stays aligned with implementation

## Deliverables required before closing this phase
- scaffolding exists in the repo
- status/handoff docs are updated
- bootstrap path exists
- migration path exists
- baseline checks are run or blocked reasons are documented

## Constraints to preserve
- mostly Go
- Vehicle Positions first
- Trip Updates pluggable
- draft GTFS separate from published GTFS
- conservative matching
- no rider apps, payments, or dispatcher CAD
- external integrations must stay behind adapters

## On completion of this phase
When Phase 0 is complete:
- update `docs/current-status.md`
- create or update `docs/handoffs/phase-00.md`
- point this file to the next active phase
- summarize:
  - what was added
  - what checks ran
  - what is still blocked
  - exact next action for the next Codex instance

## If this file becomes stale
The active Codex instance must overwrite this file at the end of its phase so that a fresh instance can resume from here safely.