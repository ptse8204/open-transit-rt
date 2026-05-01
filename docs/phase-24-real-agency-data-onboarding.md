# Phase 24 — Real Agency Data Onboarding

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Make Open Transit RT practical for importing and publishing a real agency’s GTFS data, not only the committed demo fixture.

The repo can import and publish sample GTFS. A real agency needs a guided path for validation errors, service calendars, shapes, frequencies, metadata, and publish review.

## Scope

1. Real GTFS onboarding checklist.
2. GTFS import failure triage guide.
3. Common validator issue explanations.
4. GTFS Studio operator flow improvements, if needed.
5. Agency metadata approval flow.
6. Publish review checklist.
7. Redaction rules for real agency data.

## Required Work

### 1) Agency GTFS Intake Checklist

Create a checklist covering:

- source of GTFS ZIP;
- agency permission;
- license/contact metadata;
- routes/stops/trips/stop_times/calendar coverage;
- shapes and frequencies expectations;
- service date/timezone review;
- validation command path;
- publish approval.

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

### 3) Import/Publish Flow

Improve docs or small UI guidance for:

- uploading/importing GTFS;
- reviewing validation output;
- correcting errors;
- publishing the active feed;
- verifying public feed URLs.

### 4) Real Data Redaction

If real agency data is used in repo examples, require explicit review:

- no private contracts;
- no private contact info;
- no private operator notes;
- no real non-public vehicle/device identifiers;
- use synthetic fixtures when possible.

## Acceptance Criteria

Phase 24 is complete only when:

- a real-agency GTFS onboarding guide exists;
- common import/validation failures are explained;
- publish review checklist exists;
- docs are understandable to non-expert agency staff;
- real data privacy/redaction boundaries are clear;
- no real private agency data is committed without explicit review.

## Required Checks

```bash
make validate
make test
git diff --check
```

If GTFS import code or fixtures change:

```bash
make test-integration
make demo-agency-flow
```

## Explicit Non-Goals

Phase 24 does not:

- claim that a real agency accepted or endorsed the repo;
- publish private agency data;
- add fare/payment/rider features;
- replace GTFS Studio with a full design suite;
- claim validator success means consumer acceptance.

## Likely Files

- `docs/tutorials/real-agency-gtfs-onboarding.md`
- `docs/tutorials/gtfs-validation-triage.md`
- `docs/tutorials/agency-first-run.md`
- `docs/tutorials/production-checklist.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-24.md`
