# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 13 — Consumer Submission and Acceptance Evidence is complete for the initial docs/evidence tracker structure.

Phase 12 — Deployment Evidence Hardening remains closed for the OCI pilot evidence scope. Do not reopen Phase 12 unless a blocking defect directly prevents evidence-bounded consumer submission work.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is closed for the OCI pilot hosted/operator evidence scope.
- Phase 13 created the consumer-submission evidence layer for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.
- All seven current consumer/aggregator records are `not_started`.
- No current repo evidence supports submitted, under-review, accepted, rejected, or blocked claims for any target.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-13.md`
4. `docs/consumer-submission-evidence.md`
5. `docs/evidence/consumer-submissions/README.md`
6. `docs/phase-13-consumer-submission-evidence.md`
7. `docs/phase-12-deployment-evidence-hardening.md`
8. `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
9. `docs/compliance-evidence-checklist.md`
10. `docs/dependencies.md`
11. `docs/prompts/calitp-truthfulness.md`
12. `README.md`
13. `docs/tutorials/calitp-readiness-checklist.md`
14. `docs/tutorials/production-checklist.md`

## Current Objective

Collect real external consumer or aggregator submission evidence only when an operator has redacted artifacts to add.

Do not claim CAL-ITP compliance, production readiness, marketplace/vendor equivalence, or consumer acceptance from repo evidence, validator success, public fetch proof, workflow records, or the Phase 12 hosted packet alone.

## Exact First Commands

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

## Current Evidence Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` includes hosted/operator proof for public HTTPS feed fetches, TLS/redirect behavior, auth boundaries, clean hosted validation, monitoring/alert lifecycle, backup/restore, deployment rollback, and scorecard export job history.
- The Phase 13 tracker at `docs/evidence/consumer-submissions/README.md` links to the Phase 12 packet as supporting evidence only.
- Validator success and public fetch proof are not consumer acceptance.
- Consumer-ingestion workflow records are not third-party acceptance.
- Acceptance may be claimed only for the named consumer, feed scope, URL root, and evidence date shown in a retained evidence artifact.

## First Files Likely To Edit

- `docs/evidence/consumer-submissions/current/<target>.md`
- `docs/evidence/consumer-submissions/README.md`
- a redacted artifact path under an evidence packet directory, if real external evidence exists
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-13.md`

## Constraints To Preserve

- Keep claims evidence-bounded and truthful.
- Keep Phase 12 closed unless a blocking defect directly prevents Phase 13 evidence work.
- Do not add consumer submission APIs unless explicitly required and supported by a public documented API.
- Do not automate fake submissions.
- Do not invent acceptance, rejection, receipt, or blocker evidence.
- Do not change public feed URLs.
- Do not add external predictor integrations or reopen GTFS/GTFS-RT product implementation.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
