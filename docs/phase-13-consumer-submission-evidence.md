# Phase 13 — Consumer Submission And Acceptance Evidence

## Status

Implemented documentation/evidence track for the initial tracker structure.

No current target has submitted, under-review, accepted, rejected, or blocked evidence in the repository.

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

- **not_started** — no submission has been made.
- **prepared** — packet prepared only, no submission.
- **submitted** — packet sent with verifiable timestamp and destination.
- **under_review** — consumer acknowledged review in progress.
- **accepted** — consumer explicitly confirmed acceptance/ingestion.
- **rejected** — consumer explicitly rejected or requested resubmission.
- **blocked** — submission blocked by named missing evidence/action.

Category transitions must be audit-loggable and accompanied by supporting artifacts.

## Allowed Claims By Status

| Status | Allowed claim |
| --- | --- |
| `not_started` | No submission has been made. |
| `prepared` | Packet prepared only, no submission. |
| `submitted` | Submission sent, no acceptance implied. |
| `under_review` | Consumer review in progress, no acceptance implied. |
| `accepted` | Acceptance may be claimed only for the named consumer, feed scope, URL root, and evidence date. |
| `rejected` | Rejection documented, no acceptance claim. |
| `blocked` | Submission blocked by named missing evidence/action. |

## Workflow Tracker Requirements

For each consumer, track at minimum:
- canonical consumer name,
- current category,
- status effective timestamp,
- who changed the status,
- tracker last reviewed timestamp,
- reviewed by,
- linked Phase 12 evidence packet,
- submission packet/version reference,
- feed root submitted,
- exact feed URLs submitted,
- validation evidence reference,
- Phase 12 evidence packet reference,
- acceptance-scope fields,
- evidence link(s): email thread, portal screenshot, ticket ID, or API response,
- redaction notes,
- next action,
- allowed public wording.

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
4. public fetch proof is **not** automatic consumer acceptance.
5. Mobility Database and transit.land listing/processing are not proof that Google/Apple/Transit/Bing/Moovit accepted feeds.
6. "No rejection received" is **not** acceptance.
7. verbal/informal statements are not enough without retained evidence.

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
- Every target has a current evidence record and reusable template.
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
make demo-agency-flow
make test-integration
docker compose -f deploy/docker-compose.yml config
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
- `docs/handoffs/phase-13.md`
- `docs/consumer-submission-evidence.md`
- `docs/evidence/consumer-submissions/README.md`
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
