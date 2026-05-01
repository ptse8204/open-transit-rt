# Phase 32 — Public Launch And Ecosystem Outreach

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

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

## Acceptance Criteria

Phase 32 is complete only when:

- public messaging is short and understandable;
- no compliance/acceptance/endorsement claim is introduced;
- agency one-pager exists;
- demo outline exists;
- GitHub star/support wording remains friendly and non-pushy;
- launch checklist includes truthfulness and security review.

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
- create marketing that outruns evidence.

## Likely Files

- `docs/public-launch-checklist.md`
- `docs/agency-one-pager.md`
- `docs/demo-video-outline.md`
- `README.md`
- `wiki/README.md`
- `docs/roadmap-status.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-32.md`
