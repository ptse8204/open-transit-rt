# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 16 — Agency Onboarding Product Packaging is complete for the approved local packaging, first-run documentation, device onboarding, and operator-friendly output scope.

Phase 12 remains closed for the OCI pilot hosted/operator evidence scope. Phase 13 remains closed for the initial consumer-submission evidence tracker structure. Phase 14 remains closed for public launch polish. Phase 15 remains closed for targeted public repo hygiene and evidence redaction review. Do not reopen earlier phases unless a blocking truthfulness or safety issue directly affects the next task.

## Phase Status

- Phases 0 through 16 are closed for their documented scopes.
- Phase 16 added the local Compose `app` profile, `deploy/Dockerfile.local`, `deploy/Caddyfile.local`, `scripts/agency-local-app.sh`, `scripts/device-onboarding.sh`, `make agency-app-up`, `make agency-app-down`, `make agency-app-logs`, and `make agency-app-reset`.
- `make agency-app-up` starts the full local stack behind `http://localhost:8080`, applies migrations, seeds demo data, imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for readiness, verifies public feed URLs, and prints next steps.
- The local reverse proxy is demo packaging only. Production still requires HTTPS/TLS and deployment-owned admin network boundaries.
- All seven current consumer/aggregator records are still `not_started`; no current repo evidence supports submitted, under-review, accepted, rejected, or blocked claims for any consumer target.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `README.md`
4. `wiki/README.md`
5. `docs/README.md`
6. `docs/handoffs/phase-16.md`
7. `docs/phase-17-deployment-automation-pilot-operations.md`
8. `SECURITY.md`
9. `docs/evidence/redaction-policy.md`
10. `docs/compliance-evidence-checklist.md`
11. `docs/prompts/calitp-truthfulness.md`
12. `docs/tutorials/agency-first-run.md`
13. `docs/tutorials/local-quickstart.md`
14. `docs/tutorials/agency-demo-flow.md`
15. `docs/tutorials/deploy-with-docker-compose.md`
16. `docs/tutorials/production-checklist.md`
17. `docs/dependencies.md`
18. `docs/decisions.md`

## Current Objective

The next planned phase is Phase 17 — Deployment Automation And Pilot Operations. Use the Phase 16 local app package as the operator-experience baseline, but do not treat it as production hosting, TLS, admin network policy, or compliance evidence.

Do not claim CAL-ITP/Caltrans compliance, production readiness, marketplace/vendor equivalence, agency endorsement, or consumer acceptance from repo evidence, validator success, public fetch proof, workflow records, stars, the local app package, or the Phase 12 hosted packet alone.

## Exact First Commands

```bash
make validate
make test
make smoke
make agency-app-up
docker compose -f deploy/docker-compose.yml --profile app config
git diff --check
```

If touching demo validation or public feed workflows, also run:

```bash
make demo-agency-flow
```

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` includes hosted/operator proof for public HTTPS feed fetches, TLS/redirect behavior, auth boundaries, clean hosted validation, monitoring/alert lifecycle, backup/restore, deployment rollback, and scorecard export job history.
- Phase 15 found real secrets only in ignored local `.cache` files, not in tracked files or history for those `.cache` paths. Rotation/revocation is still required before further real pilot use.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, or raw private operator artifacts.
- `docs/evidence/redaction-policy.md` is the evidence safety rule for future packets.
- Validator success and public fetch proof are not consumer acceptance.
- Consumer-ingestion workflow records are not third-party acceptance.

## First Files Likely To Edit

- `docs/phase-17-deployment-automation-pilot-operations.md`
- `scripts/oci-pilot.sh`
- `deploy/oci/`
- `deploy/systemd/`
- `docs/runbooks/`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-17.md`

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Do not change backend API contracts, database schema, public feed URLs, consumer-submission statuses, or evidence claims unless the active phase explicitly requires it.
- Do not add consumer submission APIs or automate fake submissions.
- Do not invent acceptance, rejection, receipt, blocker, compliance, or endorsement evidence.
- Keep `.cache` ignored and do not mount or commit it unless a specific safe need is documented and reviewed.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.
- Do not echo raw long-lived secrets in final command output.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/roadmap-post-phase-14.md` as the roadmap source of truth.

The next planned phase is:

- Phase 17 — Deployment Automation And Pilot Operations

Future roadmap docs:

- `docs/phase-17-deployment-automation-pilot-operations.md`
- `docs/phase-18-admin-ux-agency-operations-console.md`
- `docs/phase-19-realtime-quality-eta-improvement.md`
- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`
