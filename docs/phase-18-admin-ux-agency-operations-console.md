# Phase 18 — Admin UX And Agency Operations Console

## Status

Complete for the approved minimal admin operations console scope.

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

Implemented as a guided setup checklist rather than a full wizard.

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

### 2) Setup Checklist

Phase 18 provides a small checklist for:

- agency metadata;
- license/contact metadata;
- GTFS import or GTFS Studio;
- supported device token rotate/rebind;
- publication bootstrap;
- first validation run.

It does not add a new first-time device credential API. The existing supported browser action is rotate/rebind, which may bootstrap a device credential through the existing service behavior. One-time tokens are displayed only in the immediate POST response returned by that flow.

### 3) Device Management

Expose safe admin flows for:

- rotating/rebinding device tokens;
- listing active device/vehicle bindings;
- warning about secret handling.

If a separate first-time creation API is needed later, it should be designed explicitly in a future phase.

### 4) Evidence Links

Make evidence state easier to find:

- Phase 12 hosted evidence;
- Phase 13 consumer tracker;
- scorecard history.

The console prefers database `consumer_ingestion` records where they exist. Targets not present in the running database link back to the file-backed Phase 13 docs tracker rather than inventing statuses.

## Routes Added

- `/admin/operations`
- `/admin/operations/feeds`
- `/admin/operations/telemetry`
- `/admin/operations/devices`
- `/admin/operations/consumers`
- `/admin/operations/evidence`
- `/admin/operations/setup`
- `/admin/alerts/console`

GTFS Studio links back to the Operations Console where practical. The local app proxy routes `/admin/operations*` only for local/demo packaging; the OCI public edge still exposes public feed paths only.

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
