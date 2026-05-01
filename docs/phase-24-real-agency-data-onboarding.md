# Phase 24 — Real Agency Data Onboarding

## Status

Complete for the docs/process and evidence-template scope.

## Purpose

Make Open Transit RT practical for importing and publishing a real agency’s GTFS data, not only the committed demo fixture.

The repo can import and publish sample GTFS. A real agency needs a guided path for validation errors, service calendars, shapes, frequencies, metadata, and publish review.

## Scope

1. Real GTFS onboarding checklist.
2. GTFS import failure triage guide.
3. Common validator issue explanations.
4. Existing import, Studio publish, Operations Console, and public-feed review path documentation.
5. Agency metadata approval flow.
6. Publish review checklist.
7. Redaction rules for real agency data.
8. Template-only future real-agency evidence scaffold.

Phase 24 did not add real agency data, change backend behavior, change public feed URLs, change consumer statuses, or change Phase 23 final-root status.

## Required Work

### 1) Agency GTFS Intake Checklist

Create a checklist covering:

- source of GTFS ZIP;
- agency permission;
- license/contact metadata;
- agency identity;
- timezone;
- routes/stops/trips/stop_times/calendar coverage;
- shapes, frequencies, and block expectations;
- service date/timezone review;
- validation command path;
- publish approval;
- privacy/redaction review;
- final public-feed review.

Implemented in `docs/tutorials/real-agency-gtfs-onboarding.md`.

### 2) Validation Triage

Document common issues:

- missing required files;
- invalid route types;
- bad stop_times references;
- calendar/calendar_dates problems;
- times beyond 24:00:00;
- shapes ordering;
- block_id expectations;
- malformed frequencies.
- duplicate IDs;
- timezone/date mistakes;
- empty or all-suppressed service;
- importer errors versus validator errors;
- when to ask for technical help.

Implemented in `docs/tutorials/gtfs-validation-triage.md`.

### 3) Import/Publish Flow

Improve docs or small UI guidance for:

- uploading/importing GTFS;
- reviewing validation output;
- correcting errors;
- publishing the active feed;
- verifying public feed URLs.

Implemented as documentation of the existing `cmd/gtfs-import` CLI path, GTFS Studio typed draft publish path, `/public/feeds.json`, `/public/gtfs/schedule.zip`, and `/admin/operations` review path.

### 4) Real Data Redaction

If real agency data is used in repo examples, require explicit review:

- no private contracts;
- no private contact info;
- no private operator notes;
- no real non-public vehicle/device identifiers;
- use synthetic fixtures when possible.

Implemented in onboarding docs and `docs/evidence/real-agency-gtfs/`.

### 5) Evidence Template

Added `docs/evidence/real-agency-gtfs/README.md` and `docs/evidence/real-agency-gtfs/templates/import-review-template.md`.

The evidence directory contains templates only until real agency-approved, public-safe evidence exists. Placeholder GTFS ZIPs, fake validation outputs, fake approvals, and fake import evidence are not allowed.

## Acceptance Criteria

Phase 24 is complete only when:

- a real-agency GTFS onboarding guide exists;
- common import/validation failures are explained;
- publish review checklist exists;
- docs are understandable to non-expert agency staff;
- real data privacy/redaction boundaries are clear;
- no real private agency data is committed without explicit review.
- future real-agency import evidence templates exist.
- Phase 23 final-root limits remain clear.

## Required Checks

```bash
make validate
make test
git diff --check
```

If GTFS import code, fixtures, local app behavior, or demo flow change:

```bash
make test-integration
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

## Explicit Non-Goals

Phase 24 does not:

- claim that a real agency accepted or endorsed the repo;
- publish private agency data;
- add fare/payment/rider features;
- replace GTFS Studio with a full design suite;
- claim validator success means consumer acceptance.
- claim local or DuckDNS pilot review proves agency-owned production-domain readiness.

## Likely Files

- `docs/tutorials/real-agency-gtfs-onboarding.md`
- `docs/tutorials/gtfs-validation-triage.md`
- `docs/tutorials/agency-first-run.md`
- `docs/tutorials/production-checklist.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-24.md`
