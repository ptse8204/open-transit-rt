# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 12 — Deployment Evidence Hardening is in progress.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is active.
- **Phase 12 Step 1 is complete** as repo-side docs/runbooks/evidence-template scaffolding.
- **Phase 12 Step 2 is partially complete** with a real dated local evidence packet at `docs/evidence/captured/local-demo/2026-04-22/`.
- Real hosted HTTPS evidence collection for Phase 12 is still pending.
- Phases 13 and 14 remain planning/docs tracks only.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/phase-12-deployment-evidence-hardening.md`
4. `docs/handoffs/phase-12-step-2.md`
5. `docs/compliance-evidence-checklist.md`
6. `docs/dependencies.md`
7. `docs/prompts/calitp-truthfulness.md`
8. `README.md`
9. `docs/runbooks/deployment-evidence-overview.md`
10. `docs/evidence/captured/local-demo/2026-04-22/README.md`
11. `docs/evidence/README.md`
12. `docs/tutorials/production-checklist.md`
13. `docs/tutorials/calitp-readiness-checklist.md`

## Current Objective

Execute the next Phase 12 slice by collecting **hosted deployment/operator evidence** using the committed runbooks/templates. The local demo packet is useful but does not close the hosted evidence requirements.

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
- Java was not installed during Step 2, so local static validator execution failed even though pinned tooling was present.
- Local GTFS-RT validator execution failed because the current Docker-backed wrapper invocation passed unsupported `--schedule` arguments to the pinned image.
- No hosted HTTPS hostname, TLS certificate, production reverse proxy config, monitoring alert lifecycle, production backup policy, or consumer acceptance evidence has been captured.
- `docs/evidence/captured/hosted-pending/2026-04-22/` now contains an operator intake packet for the missing hosted artifacts; it is not completed evidence.
- `make collect-hosted-evidence` is available for hosted feed fetch, TLS, validation, and manual scorecard collection once `ENVIRONMENT_NAME` and `PUBLIC_BASE_URL` are set. `ADMIN_TOKEN` is needed for hosted validation and scorecard export.
- `make audit-hosted-evidence` is available for completed hosted packets once `EVIDENCE_PACKET_DIR` is set; it should fail until pending markers, failed validators, missing TLS evidence, and missing operator-supplied artifacts are resolved.
- Consumer submission APIs remain out of scope; workflow records are not third-party acceptance proof.

## First Files Likely To Edit

- `docs/evidence/captured/<hosted-environment>/*` (real hosted operator evidence artifacts when available)
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-12-step-3.md` (or equivalent next-step handoff)

## Constraints To Preserve

- Keep claims evidence-bounded and truthful.
- Keep Trip Updates pluggable and architecture boundaries unchanged.
- Do not add unrelated product scope (rider apps, fares, CAD/dispatch).
- Do not reopen implementation work from Phases 9–11 during this docs/evidence pass.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
