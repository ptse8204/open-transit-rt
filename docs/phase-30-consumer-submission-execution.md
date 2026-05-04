# Phase 30 — Consumer Submission Execution

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md` and only when an operator is authorized to proceed.

Phase 29A — External Predictor Adapter Evaluation and Phase 29B — AVL / Vendor Adapter Pilot Implementation now precede this phase in the Track B roadmap. Consumer statuses must still remain unchanged until retained, redacted, target-originated evidence supports a specific target update.

## Purpose

Move prepared packets to real consumer or aggregator submissions when an operator has permission and retained evidence.

This phase is evidence-first. It does not automate private portals or infer acceptance from validation.

## Scope

1. Verify one or more official submission paths.
2. Submit only when authorized.
3. Store redacted receipts/tickets/correspondence.
4. Update the specific target from `prepared` to evidence-backed status.
5. Maintain tracker/status consistency.
6. Document next review/acceptance/rejection steps.

## Required Work

### 1) Official Path Verification

For a selected target, retain evidence of the official path:

- public target-owned page;
- target-originated email;
- authorized portal screenshot;
- official support documentation.

Do not guess.

### 2) Submission

If authorized, submit the packet outside the repo.

Record:

- submitted feed root;
- submitted feed URLs;
- date/time;
- operator;
- target;
- evidence artifact;
- limitations.

### 3) Status Update

Update only the named target:

- current record;
- `status.json`;
- human tracker;
- artifact directory;
- current status docs if needed.

### 4) Follow-Up

Track:

- under-review acknowledgement;
- rejection/change request;
- acceptance confirmation;
- blocker state.

## Acceptance Criteria

Phase 30 is complete only when:

- at least one target has real submitted evidence, or the phase records why no submission was possible;
- no target status changes without target-originated evidence;
- tracker and `status.json` stay aligned;
- no acceptance is claimed without proof;
- private portal and personal data are redacted.

## Required Checks

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json
```

Run redaction scan over new/edited evidence.

## Explicit Non-Goals

Phase 30 does not:

- automate submissions;
- scrape portals;
- fake receipts;
- claim acceptance from submission;
- claim compliance from submission;
- change backend behavior.

## Likely Files

- `docs/evidence/consumer-submissions/artifacts/<target>/`
- `docs/evidence/consumer-submissions/current/<target>.md`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-30.md`
