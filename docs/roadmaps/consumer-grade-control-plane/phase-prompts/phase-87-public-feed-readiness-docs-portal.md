# Phase 87 Prompt — Public Feed Readiness And Docs Portal

## Goal

Make public feed URLs, feed metadata, and readiness explanations clear for operators and technical consumers without changing consumer statuses or claiming acceptance.

## Scope

- Public feed preview pages or private preview of public metadata, depending on auth/public-surface review.
- Feed URL copy/check page.
- `feeds.json` explanation.
- Schedule and GTFS-RT endpoint explanation.
- Consumer packet prepared-only explanation.
- Off-host validation guidance.
- Source-of-truth website/domain guidance.

## Public route review

Any new public route must be separately reviewed by the Master Agent, QA Sub-Agent, and Claim-Boundary Sub-Agent. Prefer documentation/static public pages over dynamic public admin surfaces.

## Do not

- Move consumer statuses.
- Claim submission/acceptance/ingestion/listing/display.
- Claim final-root readiness.
- Claim compliance.
- Write evidence.

## Validation

Baseline validation plus:

```bash
make validate
make test
make audit-final-claim-review
```

## Commits

```text
Phase 87 -- Checkpoint 000001: add public feed readiness portal plan
Phase 87 -- Checkpoint 000002: add feed URL and metadata preview guidance
Phase 87 -- Checkpoint 000003: add prepared-only consumer packet explanation
Phase 87 -- Checkpoint 000004: add off-host validation and source-of-truth guidance
Phase 87 -- Checkpoint 000005: close public feed readiness portal review
```
