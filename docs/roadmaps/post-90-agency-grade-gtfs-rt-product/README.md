# Open Transit RT — Post-90 Agency-Grade GTFS-RT Product Roadmap

Generated: 2026-05-14
Project: `ptse8204/open-transit-rt`
Proposed repo path: `docs/roadmaps/post-90-agency-grade-gtfs-rt-product/`

## Purpose

This is a Codex-ready roadmap for the next product arc after Phase 90. It is
intended to make Open Transit RT more releasable, more maintainable, easier for
less technical staff, and stronger as an open-source GTFS / GTFS-Realtime
platform.

It is not limited to frontend quality-of-life work. It covers release hardening,
browser task trials, safe browser command expansion, GTFS versioning and fix
planning, GTFS-RT usefulness, connector maturity, device onboarding, monitoring,
small-host operations, multi-agency hardening, training/adoption, and optional
future evidence gates.

## Baseline

- Phases 75-90 are complete for the Consumer-Grade Control Plane product track.
- Phase 72 remains `needs_review`.
- Phase 89 remains the current local `v0.1.0-rc.1` gate and closes as `needs_review`.
- No release tag, package, published image, retained evidence, consumer movement,
  final-root proof, compliance proof, production readiness proof, vendor proof,
  or production-grade ETA proof exists.
- This roadmap is active only when a maintainer explicitly authorizes it. The
  autonomous Phase 91-110 kickoff authorization is recorded in
  `CODEX-KICKOFF-AUTONOMOUS-PHASE-91-TO-110.md`.

## Use order

1. `00-CODEX-READ-ME-FIRST.md`
2. `01-roadmap-overview.md`
3. `02-phases-and-checkpoints.md`
4. `03-master-subagent-operating-manual.md`
5. `04-validation-and-claim-boundaries.md`
6. `05-autonomous-run-policy.md`
7. relevant `phase-prompts/phase-XX-*.md`
8. relevant `audit-prompts/*.md`
9. `CODEX-KICKOFF-AUTONOMOUS-PHASE-91-TO-110.md`
10. `05-open-source-gap-map.md` as supporting gap context only

## Commit rule

Every checkpoint and every phase closeout must end in a git commit.

```text
Phase XX -- Checkpoint 000001: add <phase> plan
Phase XX -- Checkpoint 000002: <implementation checkpoint>
Phase XX -- Checkpoint 000003: <validation or audit checkpoint>
Phase XX -- Checkpoint 000004: close <phase> review
```

Checkpoint numbers reset per phase. Do not mix phase scopes in one commit.
