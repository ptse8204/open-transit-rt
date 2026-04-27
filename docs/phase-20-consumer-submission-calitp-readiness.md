# Phase 20 — Consumer Submission Execution And CAL-ITP Readiness Program

## Status

Complete for the Phase 20 docs/evidence packet preparation and readiness summary scope.

## Purpose

Phase 20 moves beyond internal tracking by preparing real consumer/aggregator submission packets where an operator may later have permission to submit. It also prepares a truthful readiness packet for California-facing review.

No external portal was contacted, no submission was automated, and no acceptance or compliance claim was added.

## Scope

1. Submission packet preparation.
2. External consumer/aggregator submissions only when an operator has authorization and retained evidence.
3. Evidence updates for submitted/under_review/accepted/rejected states only when target-originated evidence exists.
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

Phase 20 prepared packets for all seven targets under `docs/evidence/consumer-submissions/packets/`.

Each packet includes evidence freshness fields, submission-method fields marked `not verified`, all five feed URLs, license/contact metadata, Phase 12 evidence links, validator evidence links, Phase 19 replay/quality boundary, redaction notes, next action, allowed wording, and an operator warning not to submit from repo docs alone.

The machine-readable snapshot is `docs/evidence/consumer-submissions/status.json`.

### 2) External Submission Evidence

Only update a target status if there is real redacted evidence:

- receipt;
- ticket;
- portal screenshot;
- email correspondence;
- rejection reason;
- acceptance confirmation.

Phase 20 moved targets to `prepared` only because complete packets exist. No target moved to `submitted`, `under_review`, `accepted`, `rejected`, or `blocked`.

### 3) California Readiness Summary

Prepare a document that says exactly:

- what is code-complete;
- what is deployment-proven;
- what is consumer-submitted;
- what is consumer-accepted;
- what remains missing.

Created `docs/california-readiness-summary.md`. It explicitly records that agency-owned stable URL/domain proof remains missing and that the OCI pilot DuckDNS domain is pilot evidence, not agency-domain production proof.

### 4) Marketplace/Vendor Gap Review

Record what would still be needed for vendor-like adoption:

- support model;
- SLAs;
- procurement docs;
- onboarding package;
- maintenance plan;
- hardware/device support strategy.

Created `docs/marketplace-vendor-gap-review.md`.

## Acceptance Criteria

Phase 20 is complete only when:

- at least one real submission packet exists or the phase clearly records why no submission was possible;
- every updated consumer status is evidence-backed;
- no acceptance is claimed without third-party proof;
- readiness docs separate code/deployment/submission/acceptance;
- next gaps are clear.

Phase 20 meets these criteria for packet preparation. It does not claim external submission execution because no retained third-party submission evidence exists.

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
- contact external portals;
- guess submission paths;
- add backend API behavior;
- claim consumer ingestion, hosted SaaS availability, agency endorsement, or production-grade ETA quality.
