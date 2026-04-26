# Phase 20 — Consumer Submission Execution And CAL-ITP Readiness Program

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 20 moves beyond internal tracking by executing real consumer/aggregator submission workflows where an operator has permission to do so. It also prepares a truthful readiness packet for California-facing review.

## Scope

1. Submission packet preparation.
2. External consumer/aggregator submissions.
3. Evidence updates for submitted/under_review/accepted/rejected states.
4. California readiness summary.
5. Gap list for marketplace/vendor equivalence.

## Required Work

### 1) Submission Packets

Prepare target-specific packets using:

- OCI pilot feed URLs;
- Phase 12 hosted evidence;
- validator records;
- license/contact metadata;
- schedule/realtime feed list;
- operator contact notes.

### 2) External Submission Evidence

Only update a target status if there is real redacted evidence:

- receipt;
- ticket;
- portal screenshot;
- email correspondence;
- rejection reason;
- acceptance confirmation.

### 3) California Readiness Summary

Prepare a document that says exactly:

- what is code-complete;
- what is deployment-proven;
- what is consumer-submitted;
- what is consumer-accepted;
- what remains missing.

### 4) Marketplace/Vendor Gap Review

Record what would still be needed for vendor-like adoption:

- support model;
- SLAs;
- procurement docs;
- onboarding package;
- maintenance plan;
- hardware/device support strategy.

## Acceptance Criteria

Phase 20 is complete only when:

- at least one real submission packet exists or the phase clearly records why no submission was possible;
- every updated consumer status is evidence-backed;
- no acceptance is claimed without third-party proof;
- readiness docs separate code/deployment/submission/acceptance;
- next gaps are clear.

## Required Checks

```bash
make validate
make test
git diff --check
```

Run smoke/demo checks if docs or evidence references change materially.

## Explicit Non-Goals

Phase 20 does not:

- fake submissions;
- automate private portal workflows without permission;
- claim compliance from submission alone;
- claim marketplace equivalence without program approval.
