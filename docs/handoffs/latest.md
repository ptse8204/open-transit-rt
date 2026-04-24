# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 13 planning/docs can proceed next. Phase 12 — Deployment Evidence Hardening is closed for the OCI pilot evidence scope.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is closed for the OCI pilot evidence scope.
- **Phase 12 Step 1 is complete** as repo-side docs/runbooks/evidence-template scaffolding.
- **Phase 12 Step 2 is partially complete** with a real dated local evidence packet at `docs/evidence/captured/local-demo/2026-04-22/`.
- **Phase 12 Step 3 is complete** as repo-side closure tooling hardening. It did not itself include hosted deployment evidence.
- **Phase 12 hosted closure evidence is complete** at `docs/evidence/captured/oci-pilot/2026-04-24/`.
- Phases 13 and 14 remain planning/docs tracks only.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/phase-12-deployment-evidence-hardening.md`
4. `docs/handoffs/phase-12-step-3.md`
5. `docs/compliance-evidence-checklist.md`
6. `docs/dependencies.md`
7. `docs/prompts/calitp-truthfulness.md`
8. `README.md`
9. `docs/runbooks/deployment-evidence-overview.md`
10. `docs/evidence/captured/local-demo/2026-04-22/README.md`
11. `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
12. `docs/evidence/README.md`
13. `docs/tutorials/production-checklist.md`
14. `docs/tutorials/calitp-readiness-checklist.md`

## Current Objective

Use `docs/evidence/captured/oci-pilot/2026-04-24/` as the latest hosted/operator evidence packet. It closes Phase 12 for the OCI pilot evidence scope because the hosted audit passed.

Do not claim production readiness, CAL-ITP compliance, or consumer acceptance without hosted and third-party evidence.

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

- The OCI pilot packet includes public HTTPS feed fetches, TLS/redirect evidence, public-edge auth-boundary proof, SSH-tunneled admin auth proof, clean hosted validator records, monitoring/alert lifecycle evidence, backup/restore evidence, deployment data-restore rollback proof, and scorecard export job-history proof.
- The closure command passed: `EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence`.
- The final current-live recheck on April 24, 2026 refreshed the packet with active `gtfs-import-3`, passed hosted validator artifacts for schedule, Vehicle Positions, Trip Updates, and Alerts, and `canonical_validation_complete=true`.
- The scorecard still has consumer ingestion red because external consumer acceptance is outside Phase 12.
- Consumer submission APIs remain out of scope; workflow records are not third-party acceptance proof.

## First Files Likely To Edit

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- later phase docs or evidence packets, depending on the next requested track

## Constraints To Preserve

- Keep claims evidence-bounded and truthful.
- Keep Trip Updates pluggable and architecture boundaries unchanged.
- Do not add unrelated product scope (rider apps, fares, CAD/dispatch).
- Do not reopen implementation work from Phases 9–11 during this docs/evidence pass.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
