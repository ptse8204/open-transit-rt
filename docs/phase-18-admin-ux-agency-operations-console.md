# Phase 18 — Admin UX And Agency Operations Console

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 18 reduces command-line dependence by giving agency operators a simple web surface for common tasks and readiness checks. The goal is not a flashy product UI. The goal is clarity and operational confidence.

## Scope

1. Setup/readiness dashboard.
2. Feed URL and validation status view.
3. Telemetry freshness view.
4. Device credential management surface.
5. Alerts operations surface improvements.
6. Consumer submission evidence view.
7. Minimal setup wizard if feasible.

## Required Work

### 1) Dashboard

Add a minimal dashboard showing:

- active GTFS feed version;
- public feed URLs;
- latest validation status;
- telemetry freshness;
- service readiness;
- last scorecard;
- current consumer-submission statuses.

### 2) Setup Wizard

If feasible, create a small wizard for:

- agency metadata;
- license/contact metadata;
- GTFS upload/import;
- device token creation;
- publication bootstrap;
- first validation run.

### 3) Device Management

Expose safe admin flows for:

- creating device credentials;
- rotating/rebinding device tokens;
- listing active device/vehicle bindings;
- warning about secret handling.

### 4) Evidence Links

Make evidence state easier to find:

- Phase 12 hosted evidence;
- Phase 13 consumer tracker;
- scorecard history.

## Acceptance Criteria

Phase 18 is complete only when:

- an agency operator can inspect feed health and URLs from a browser;
- common setup tasks are easier than command-line-only workflows;
- device management is clear and safe;
- no admin route becomes public accidentally;
- authentication/CSRF boundaries remain intact;
- no overclaims are introduced.

## Required Checks

```bash
make validate
make test
make smoke
make demo-agency-flow
git diff --check
```

Add handler tests for any new admin/UI routes.

## Explicit Non-Goals

Phase 18 does not:

- add rider-facing apps;
- add payments/fares;
- replace all command-line operator workflows;
- build a large frontend stack unless explicitly approved;
- weaken auth.
