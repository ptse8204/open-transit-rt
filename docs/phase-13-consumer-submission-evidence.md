# Phase 13 — Consumer Submission And Acceptance Evidence

## Status

Planned documentation track. Not implemented in this repository yet.

## Purpose

Phase 13 defines how Open Transit RT deployments should track downstream consumer submission and acceptance evidence without overclaiming ingestion status.

This phase is workflow and evidence hardening only. It does not add submission API integrations or backend feature changes.

## Scope

Track the real-world status of feed submission/ingestion workflows for:
- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

Maintain per-consumer records with strict evidence categories and explicit proof links.

## Evidence Categories

Each consumer must be in exactly one category at a given time:

- **not_started** — no submission packet prepared or sent.
- **submitted** — packet sent with verifiable timestamp and destination.
- **under_review** — consumer acknowledged review in progress.
- **accepted** — consumer explicitly confirmed acceptance/ingestion.
- **rejected** — consumer explicitly rejected or requested resubmission.

Category transitions must be audit-loggable and accompanied by supporting artifacts.

## Workflow Tracker Requirements

For each consumer, track at minimum:
- canonical consumer name,
- current category,
- status effective timestamp,
- who changed the status,
- submission packet/version reference,
- evidence link(s): email thread, portal screenshot, ticket ID, or API response,
- notes with next action and due date.

Recommended tracker format:
- structured table in docs plus linked JSON artifacts,
- or structured admin export documented in runbook,
- with immutable historical snapshots per status transition.

## Overclaiming Rules (Strict)

Do not claim acceptance unless there is explicit third-party confirmation.

Required truthfulness rules:
1. `submitted` is **not** `accepted`.
2. `under_review` is **not** `accepted`.
3. validator-clean feeds are **not** automatic consumer acceptance.
4. Mobility Database and transit.land listing/processing are not proof that Google/Apple/Transit/Bing/Moovit accepted feeds.
5. "No rejection received" is **not** acceptance.
6. verbal/informal statements are not enough without retained evidence.

Allowed wording examples:
- "Submission sent on <date>; currently under review."
- "Consumer acceptance not yet confirmed."

Forbidden wording without proof:
- "Accepted by major trip planners"
- "Already ingested everywhere"
- "Consumer ready" (without consumer-specific confirmation)

## Acceptance Criteria

Phase 13 is complete only when all are true:

- A consumer workflow tracker exists for all seven named targets.
- Every target has one valid category and timestamped evidence state.
- Overclaiming guardrails are documented in tracker and handoff docs.
- Claims in README/status/handoff materials remain evidence-bounded.
- A repeatable update process exists for status changes and artifact retention.

## Commands / Artifacts / Docs To Update

### Commands

Run repo checks while updating workflow docs:

```bash
make validators-check
make validate
make test
make smoke
git diff --check
```

### Artifacts

Expected evidence artifacts (deployment-owned):
- submission packet copies,
- sent timestamps,
- acknowledgement/review receipts,
- acceptance/rejection confirmations,
- status history snapshots.

### Docs To Update During Execution

At minimum:
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/compliance-evidence-checklist.md` (if categories or interpretation rules change)
- relevant runbook docs that define artifact storage and review cadence

## Explicit Non-Goals

Phase 13 does **not**:
- implement direct submission integrations/APIs,
- guarantee or imply acceptance by any consumer,
- change runtime feed generation behavior,
- introduce new backend product features,
- reopen Phases 9–11 implementation,
- claim full CAL-ITP/Caltrans compliance.
