# Phase 26 — Admin UX Setup Wizard

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Expand the minimal Operations Console into a more complete guided setup experience for agency operators.

Phase 18 added a practical console. Phase 26 should reduce remaining command-line dependence for common setup and operating tasks while preserving security and truthfulness boundaries.

## Scope

1. Guided setup wizard or checklist.
2. Agency metadata and license/contact setup.
3. GTFS upload/import path if feasible.
4. Publication bootstrap guidance/actions.
5. Validator run/status UI.
6. Device credential flow improvements.
7. Alert authoring improvements.
8. Manual assignment override UI, if feasible.
9. Consumer packet/status viewer improvements.

## Required Work

### 1) Setup Wizard

If feasible, add browser-guided steps for:

- agency metadata;
- license/contact metadata;
- GTFS import or Studio draft path;
- publication bootstrap;
- device token setup;
- first validation run;
- public feed verification.

If a wizard is too large, implement a stronger setup checklist with action links and status signals.

### 2) GTFS Import UX

Evaluate whether GTFS upload/import can be made browser-accessible safely. If not, document the limitation and keep command-line import guidance.

### 3) Validation UX

Add or improve safe authenticated UI for:

- running validators;
- viewing last result;
- explaining warnings/errors;
- linking to evidence docs.

### 4) Alerts And Overrides

Improve operator flows for:

- alert creation/editing;
- publish/archive;
- manual assignment overrides if supported by existing backend.

### 5) Safety

Preserve:

- admin auth;
- role boundaries;
- CSRF for unsafe forms;
- one-time token handling;
- no public exposure of admin/debug JSON.

## Acceptance Criteria

Phase 26 is complete only when:

- common setup tasks are easier from the browser;
- operators can see what step is next;
- auth/CSRF/secret handling remains intact;
- no public feed URLs or protobuf contracts change;
- no unsupported compliance or acceptance claims are introduced.

## Required Checks

```bash
make validate
make test
make smoke
make demo-agency-flow
git diff --check
```

If local app is touched:

```bash
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

Add handler tests for new routes/forms.

## Explicit Non-Goals

Phase 26 does not:

- build a large SPA;
- add rider-facing apps;
- add fares/payments;
- weaken auth;
- claim consumer acceptance;
- replace all command-line operations.

## Likely Files

- `cmd/agency-config/`
- `cmd/gtfs-studio/`
- `cmd/feed-alerts/`
- `internal/*` only as needed for safe existing data access
- `docs/tutorials/agency-first-run.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-26.md`
