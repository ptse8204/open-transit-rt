# Phase 31 — Agency Pilot Program Package

## Status

Complete for the docs-only agency pilot package scope.

Phase 31 starts from the Phase 30 prepared-only consumer state. It must not
assume submission, review, acceptance, rejection, blocker, ingestion, listing,
display, or adoption evidence exists.

## Purpose

Package Open Transit RT for real-world agency pilot evaluation.

The repo has code, docs, evidence, operations, and governance. A real agency
needs a clear pilot plan: prerequisites, responsibilities, timeline, success
criteria, risks, and support boundaries.

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
10. Pilot kickoff agenda.
11. Pilot closeout summary.

## Implemented Work

### 1) Pilot Checklist

Added `docs/agency-pilot-checklist.md`, covering:

- agency data prerequisites;
- GTFS ownership;
- telemetry/device plan;
- domain plan;
- validator plan;
- consumer submission plan;
- staff/operator roles;
- responsibility matrix;
- launch/readiness review;
- exit criteria.

### 2) Training Materials

Added `docs/agency-training-outline.md`, covering:

- what Open Transit RT does;
- what it does not do;
- running local demo;
- uploading/importing GTFS;
- device token safety;
- validation and evidence;
- submission/acceptance boundaries.

### 3) Success Criteria

Added success and failure criteria in `docs/agency-pilot-program.md`, covering:

- feed availability;
- validation status;
- telemetry freshness;
- operator ability to use console;
- evidence packet completeness;
- submission readiness;
- support load.

### 4) Risk Register

Added a risk register in `docs/agency-pilot-program.md`, covering data
ownership, private data leakage, secret leakage, unstable public URLs, GTFS
validation failures, device/AVL reliability, Trip Updates quality, operations
capacity, consumer submission delay, support expectations, and multi-agency
boundary risk.

### 5) Kickoff Agenda

Added `docs/agency-pilot-kickoff-agenda.md`, covering attendees, pre-kickoff
preparation, 30-minute and 60-minute agendas, walkthrough topics, decisions,
follow-up actions, and what not to collect.

### 6) Feedback Loop

Added `docs/agency-feedback-template.md`, covering:

- agency feedback;
- bug reports;
- feature requests;
- operations issues;
- training gaps.

### 7) Closeout Summary

Added a closeout mini-template in `docs/agency-pilot-program.md`, covering
continue/pause/close, what worked, blockers, evidence collected, evidence still
missing, next operator action, and claim boundaries.

## Acceptance Criteria

Phase 31 is complete only when:

- a small agency can understand what a pilot requires;
- operator and maintainer responsibilities are clear;
- success/failure criteria are defined;
- pilot risk register exists;
- training outline exists;
- feedback template exists;
- public launch readiness checklist exists;
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
- approve public launch;
- claim hosted SaaS availability;
- claim production readiness;
- assume consumer submission, review, acceptance, rejection, blocker, ingestion,
  listing, display, or adoption evidence exists;
- add backend product features.

## Files Added Or Updated

- `docs/agency-pilot-program.md`
- `docs/agency-pilot-kickoff-agenda.md`
- `docs/agency-pilot-checklist.md`
- `docs/agency-training-outline.md`
- `docs/agency-feedback-template.md`
- `docs/phase-31-agency-pilot-program-package.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-31.md`

## Phase 32 Boundary

The recommended next phase is Phase 32 — Public Launch And Ecosystem Outreach.
Phase 32 must not assume agency adoption, consumer acceptance, CAL-ITP/Caltrans
compliance, hosted SaaS availability, paid support/SLA coverage, or production
readiness evidence exists.
