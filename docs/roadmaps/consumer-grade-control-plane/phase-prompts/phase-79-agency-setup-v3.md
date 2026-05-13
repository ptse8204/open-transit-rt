# Phase 79 Prompt — Agency Setup V3

## Goal

Make agency setup a full browser-first workflow for nontechnical staff.

## User story

An agency staff member opens `Agency Operations Cockpit / Start Here`, follows setup steps, enters agency metadata, imports or reviews GTFS, sees feed URLs, sees readiness blockers, and knows what to do next.

## Scope

- Agency profile setup and review.
- Technical contact/license metadata review.
- Feed URL configuration explanation.
- GTFS import source review and safe preview.
- Setup progress model.
- Blocker-oriented next actions.
- Technical-helper escape hatches for shell commands.
- No-developer and technical-helper modes.

## Deliverables

- Updated setup wizard UI.
- Setup progress JSON model.
- Readiness integration.
- Docs/tutorial updates.
- Tests for route auth, safe rendering, and status model.

## Non-goals

- Do not claim final-root readiness.
- Do not contact external parties.
- Do not write retained evidence.
- Do not submit to consumers.

## Validation

Baseline validation plus:

```bash
make validate
make test
make audit-product-acceptance
```

## Commits

```text
Phase 79 -- Checkpoint 000001: add agency setup v3 plan
Phase 79 -- Checkpoint 000002: implement agency profile and metadata review
Phase 79 -- Checkpoint 000003: improve browser GTFS source import review
Phase 79 -- Checkpoint 000004: add setup progress and blocker next actions
Phase 79 -- Checkpoint 000005: close agency setup v3 review
```
