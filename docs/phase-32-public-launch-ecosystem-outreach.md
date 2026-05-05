# Phase 32 — Public Launch And Ecosystem Outreach

## Status

Complete for the docs-only draft public launch materials scope.

Phase 32 created review materials only. No announcement was posted, no social
copy was published, no agency was contacted, no reporter was contacted, no
consumer or aggregator was contacted, and no public launch occurred.

## Purpose

Prepare public messaging for a broader open-source launch without overclaiming.

The project should be easier for agencies, contributors, and transit technologists to understand. Public launch should emphasize what exists, what evidence exists, what remains missing, and how people can help.

## Scope

1. Public launch README review.
2. Agency one-pager.
3. Demo video outline or script.
4. Social/share copy.
5. Contributor call-to-action.
6. GitHub star/support wording.
7. Ecosystem positioning.
8. Public checklist for truthfulness.

## Required Work

### 1) Public Messaging

Prepare concise wording for:

- what Open Transit RT is;
- who it helps;
- what it can do today;
- what evidence exists;
- what it does not claim;
- how to try it;
- how to support the project.

### 2) Agency One-Pager

Create a one-page summary for agencies:

- problem;
- solution;
- demo path;
- requirements;
- readiness boundaries;
- pilot next steps.

### 3) Demo Video Outline

Create a script/outline showing:

- local app startup;
- GTFS import;
- public feed URLs;
- Operations Console;
- device telemetry sample;
- evidence/validation view;
- consumer packet boundary.

### 4) Ecosystem Positioning

Document how the project relates to:

- GTFS/GTFS Realtime standards;
- Caltrans/CAL-ITP-style readiness;
- downstream consumers;
- existing validators;
- other open-source transit tooling.

Avoid implying official affiliation or endorsement.

### 5) Draft-Only Share Copy

Social/share copy must be described as draft copy for review. It is not
evidence that anything was posted, submitted, accepted, endorsed, or launched.

### 6) Claim-To-Evidence Review

The public launch checklist must include a claim-to-evidence table covering
local demo behavior, OCI pilot evidence, prepared consumer packets, synthetic
AVL/vendor dry-run evidence, external predictor adapter evaluation, agency
pilot package existence, missing agency-owned final-root proof, and missing
consumer acceptance evidence.

### 7) No Logo Or Affiliation Rule

Do not use agency, Caltrans/CAL-ITP, consumer, vendor, validator, or
standards-body logos unless retained permission exists. Do not use wording that
implies affiliation, sponsorship, certification, acceptance, deployment
approval, or endorsement unless retained evidence supports that exact claim.

## Acceptance Criteria

Phase 32 is complete only when:

- public messaging is short and understandable;
- no compliance/acceptance/endorsement claim is introduced;
- agency one-pager exists;
- demo outline exists;
- public share copy exists and is marked draft-only;
- ecosystem positioning exists;
- public launch checklist exists with a claim-to-evidence table;
- contributor call-to-action points to contribution docs, issue templates, docs fixes, replay fixtures, agency pilot feedback, AVL/vendor adapter examples, operations runbooks, and public-safe evidence review;
- GitHub star/support wording remains friendly and non-pushy;
- launch checklist includes truthfulness, no-affiliation, and security review;
- status and handoff docs state that no announcement was posted, no agency or consumer was contacted, and no public launch occurred.

## Required Checks

```bash
make validate
make test
git diff --check
```

If README/demo docs change materially:

```bash
make smoke
make demo-agency-flow
```

## Explicit Non-Goals

Phase 32 does not:

- claim official endorsement;
- claim compliance or acceptance;
- create paid support;
- change product behavior;
- submit to consumers;
- create marketing that outruns evidence;
- post on social media;
- email agencies;
- contact reporters;
- contact consumers or aggregators;
- publish announcements;
- represent that an actual public launch occurred.

## Likely Files

- `docs/public-launch-checklist.md`
- `docs/agency-one-pager.md`
- `docs/demo-video-outline.md`
- `docs/public-share-copy.md`
- `docs/ecosystem-positioning.md`
- `README.md`
- `docs/README.md`
- `docs/roadmap-status.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-32.md`

## Completion Summary

Phase 32 added:

- `docs/agency-one-pager.md`
- `docs/demo-video-outline.md`
- `docs/public-share-copy.md`
- `docs/ecosystem-positioning.md`
- `docs/public-launch-checklist.md`
- `docs/handoffs/phase-32.md`

README and docs navigation were updated only to point readers to the new
evaluation and draft launch materials. The README remains a public front door,
not a phase ledger.

## Next-Step Recommendation

Pause stronger public claims and pursue real retained evidence before any
stronger launch, adoption, readiness, or consumer-status wording. The next
roadmap step should be one of:

- agency-owned or agency-approved final-root proof;
- authorized target-specific consumer submission evidence;
- real agency pilot evidence;
- real deployment operations evidence.
