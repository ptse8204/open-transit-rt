# Phase 26 — Admin UX Setup Wizard

## Status

Complete for the Phase 26 browser-guided setup checklist scope.

## Purpose

Expand the minimal Operations Console into a more complete guided setup experience for agency operators.

Phase 18 added a practical console. Phase 26 should reduce remaining command-line dependence for common setup and operating tasks while preserving security and truthfulness boundaries.

## Scope

1. Guided setup checklist under `/admin/operations/setup`.
2. Agency metadata and license/contact setup through the existing publication bootstrap/update repository behavior.
3. GTFS import path guidance with browser ZIP upload intentionally deferred.
4. Publication bootstrap/update form with server-derived agency ID.
5. Validator run/status UI that accepts only feed type from the browser and maps to server-side allowlisted validator IDs.
6. Device credential and telemetry setup guidance sourced from device bindings and telemetry repository summaries.
7. Alerts setup links to the existing Alerts Console.
8. Manual assignment override UI deferred because Phase 26 did not add a safe bounded summary view.
9. Consumer packet/status viewer improvements sourced from the Phase 20 docs/evidence tracker.

## Required Work

### 1) Setup Wizard

Implemented browser-guided steps for:

- agency metadata;
- license/contact metadata;
- GTFS import or Studio draft path;
- publication bootstrap;
- device token setup;
- first validation run;
- public feed verification;
- first telemetry event;
- Alerts setup;
- consumer packet/status review;
- evidence/readiness review.

Each setup step shows a named status source such as publication metadata, feed discovery, validation records, device bindings, telemetry repository, docs/evidence tracker, or evidence links. Missing evidence remains visible as missing, not run yet, or not observed yet.

### 2) GTFS Import UX

Browser GTFS ZIP upload is deferred. The setup page links operators to the real-agency GTFS onboarding guide, GTFS Studio, validation triage, and the existing CLI import path instead of adding a new upload surface.

### 3) Validation UX

Implemented safe authenticated UI for:

- running validators by feed type only;
- viewing latest validation status through feed discovery metadata;
- explaining that validation records are supporting evidence only;
- linking to evidence and validation triage docs.

### 4) Alerts And Overrides

Improved operator flows for:

- alert creation/editing/publish/archive by linking setup to `/admin/alerts/console`;
- manual assignment overrides are intentionally deferred until a safe summary-only UI is designed.

### 5) Safety

Preserved:

- admin auth;
- role boundaries;
- CSRF for unsafe forms;
- one-time token handling;
- no public exposure of admin/debug JSON;
- no raw long-lived tokens or token hashes in setup output;
- no browser-submitted agency ID trust for setup publication forms.

## Acceptance Criteria

Phase 26 is complete only when:

- common setup tasks are easier from the browser;
- operators can see what step is next and what source backs each status;
- auth/CSRF/secret handling remains intact;
- no public feed URLs or protobuf contracts change;
- no telemetry/device APIs, Trip Updates adapter boundaries, consumer statuses, external integrations, or evidence claims change;
- no unsupported compliance or acceptance claims are introduced.

## Required Checks

```bash
make validate
make test
make smoke
make demo-agency-flow
make realtime-quality
docker compose -f deploy/docker-compose.yml config
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
