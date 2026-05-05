# Phase 31 — Agency Pilot Program Package

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

Phase 31 starts from the Phase 30 prepared-only consumer state. It must not
assume submission, review, acceptance, rejection, blocker, ingestion, listing,
display, or adoption evidence exists.

## Purpose

Package Open Transit RT for real-world agency pilot evaluation.

The repo has code, docs, evidence, operations, and governance. A real agency needs a clear pilot plan: prerequisites, responsibilities, timeline, success criteria, risks, and support boundaries.

## Scope

1. Agency pilot checklist.
2. Onboarding script/agenda.
3. Training materials.
4. Prerequisite checklist.
5. Success criteria.
6. Risk register.
7. Feedback template.
8. Support boundary summary.
9. Public launch readiness checklist.

## Required Work

### 1) Pilot Checklist

Document:

- agency data prerequisites;
- GTFS ownership;
- telemetry/device plan;
- domain plan;
- validator plan;
- consumer submission plan;
- staff/operator roles.

### 2) Training Materials

Create docs or deck outline for:

- what Open Transit RT does;
- what it does not do;
- running local demo;
- uploading/importing GTFS;
- device token safety;
- validation and evidence;
- submission/acceptance boundaries.

### 3) Success Criteria

Define pilot success/failure criteria:

- feed availability;
- validation status;
- telemetry freshness;
- operator ability to use console;
- evidence packet completeness;
- submission readiness;
- support load.

### 4) Feedback Loop

Add templates for:

- agency feedback;
- bug reports;
- feature requests;
- operations issues;
- training gaps.

## Acceptance Criteria

Phase 31 is complete only when:

- a small agency can understand what a pilot requires;
- operator and maintainer responsibilities are clear;
- success/failure criteria are defined;
- no paid support or SLA is implied unless explicitly offered;
- pilot docs stay truthfulness-bounded.

## Required Checks

```bash
make validate
make test
git diff --check
```

If demo docs change:

```bash
make smoke
make demo-agency-flow
```

## Explicit Non-Goals

Phase 31 does not:

- create paid support;
- promise agency endorsement;
- create legal/procurement commitments;
- claim compliance or consumer acceptance;
- assume consumer submission, review, acceptance, rejection, blocker, ingestion,
  listing, display, or adoption evidence exists;
- add backend product features.

## Likely Files

- `docs/agency-pilot-program.md`
- `docs/agency-pilot-checklist.md`
- `docs/agency-training-outline.md`
- `docs/agency-feedback-template.md`
- `docs/support-boundaries.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-31.md`
