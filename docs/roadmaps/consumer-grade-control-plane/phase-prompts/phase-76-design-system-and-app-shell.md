# Phase 76 Prompt — Design System And App Shell

## Goal

Make the private Operations Console feel like one coherent product instead of separate primitive admin pages.

This phase is visual/workflow polish, not a backend behavior change.

## Product outcome

A small-agency evaluator should be able to land on `Agency Operations Cockpit / Start Here`, understand the page hierarchy, recognize primary actions, and move between setup, GTFS import, feed health, readiness, validation, telemetry, connectors, maintenance, and help without reading phase docs.

## Read first

- Phase 75 roadmap pack
- `cmd/agency-config/operations.go`
- current templates/layout helpers
- current CSS/static assets
- tests covering operations routes
- README/wiki/docs pages describing the UI

## Scope

- Shared app shell.
- Design tokens: spacing, typography, colors, borders, focus, status badges.
- Components: cards, page headers, breadcrumbs, task panels, tables, forms, empty states, alerts, status badges, next-action blocks.
- Responsive layout.
- Keyboard-visible focus.
- Route-stable navigation.
- Page-level help panels.
- Snapshot-ish or string tests for key shell features where practical.

## Do not

- Replace the product with a heavy SPA.
- Change route URLs unless explicitly planned.
- Change auth/CSRF/role semantics.
- Add migrations.
- Add public admin routes.
- Add evidence writes or status changes.

## Acceptance criteria

- Core private routes use the same shell and component vocabulary.
- First-run empty states are attractive and helpful.
- Tables do not overflow badly on small screens.
- Forms have labels, descriptions, and clear submit states.
- Status chips use consistent labels and do not imply claims.
- UI has explicit `what this does not prove` copy where relevant.

## Validation

Run baseline validation from `04-validation-and-claim-boundaries.md`, plus:

```bash
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

If local app/browser automation is blocked, record exact blocker and run terminal-authenticated route checks as substitute.

## Commits

```text
Phase 76 -- Checkpoint 000001: add design system and app shell plan
Phase 76 -- Checkpoint 000002: implement shared layout tokens and components
Phase 76 -- Checkpoint 000003: apply shell to core Operations Console routes
Phase 76 -- Checkpoint 000004: add responsive and accessibility baseline checks
Phase 76 -- Checkpoint 000005: close design system and app shell review
```
