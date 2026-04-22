# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 12 — Deployment Evidence Hardening is in progress.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is active.
- **Phase 12 Step 1 is complete** as repo-side docs/runbooks/evidence-template scaffolding.
- Real hosted evidence collection for Phase 12 is still pending.
- Phases 13 and 14 remain planning/docs tracks only.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/phase-12-deployment-evidence-hardening.md`
4. `docs/handoffs/phase-12-step-1.md`
5. `docs/compliance-evidence-checklist.md`
6. `docs/dependencies.md`
7. `docs/prompts/calitp-truthfulness.md`
8. `README.md`
9. `docs/runbooks/deployment-evidence-overview.md`
10. `docs/evidence/README.md`
11. `docs/tutorials/production-checklist.md`
12. `docs/tutorials/calitp-readiness-checklist.md`

## Current Objective

Execute the next Phase 12 slice after Step 1 by collecting **real deployment/operator evidence** using the committed runbooks/templates.

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

## Known Blockers

- Docker must be installed/running for `make demo-agency-flow`, DB-backed integration flow, and `docker compose ... config`.
- Pinned validator tooling must be installed (`make validators-install`) before `make validators-check`, `make validate`, and `make smoke` can pass.
- Consumer submission APIs remain out of scope; workflow records are not third-party acceptance proof.

## First Files Likely To Edit

- `docs/evidence/captured/<environment>/*` (real operator evidence artifacts when available)
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-12-step-2.md` (or equivalent next-step handoff)

## Constraints To Preserve

- Keep claims evidence-bounded and truthful.
- Keep Trip Updates pluggable and architecture boundaries unchanged.
- Do not add unrelated product scope (rider apps, fares, CAD/dispatch).
- Do not reopen implementation work from Phases 9–11 during this docs/evidence pass.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
