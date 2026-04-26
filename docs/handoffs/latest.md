# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 17 — Deployment Automation And Pilot Operations is complete for the approved deployment automation, runbook, helper-script, systemd-example, and evidence-refresh scope.

Phases 0 through 17 are closed for their documented scopes. Do not reopen earlier phases unless a blocking truthfulness, safety, or deployment-automation issue directly requires it.

## Phase Status

- Phase 16 remains closed for local agency onboarding packaging. The local Compose app profile behind `http://localhost:8080` is still local-demo/evaluation only.
- Phase 17 added a production-directed small-agency pilot operations profile, dry-run-capable operational helpers, systemd timer examples, evidence output labels, and hosted evidence audit closure requirements.
- Consumer/aggregator records remain `not_started`; no current repo evidence supports submitted, under-review, accepted, rejected, or blocked claims for any consumer target.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `README.md`
4. `wiki/README.md`
5. `docs/README.md`
6. `docs/handoffs/phase-17.md`
7. `docs/runbooks/small-agency-pilot-operations.md`
8. `docs/phase-18-admin-ux-agency-operations-console.md`
9. `SECURITY.md`
10. `docs/evidence/redaction-policy.md`
11. `docs/compliance-evidence-checklist.md`
12. `docs/prompts/calitp-truthfulness.md`
13. `docs/tutorials/deploy-with-docker-compose.md`
14. `docs/tutorials/production-checklist.md`
15. `docs/dependencies.md`
16. `docs/decisions.md`

## Current Objective

The next planned phase is Phase 18 — Admin UX And Agency Operations Console. Use the Phase 17 runbooks and helper evidence outputs as operator workflow inputs, but do not turn them into compliance, consumer-acceptance, hosted-SaaS, or agency-endorsement claims.

Do not change backend API contracts, database schema, public feed URLs, GTFS-RT contracts, consumer-submission statuses, or evidence claims unless the active phase explicitly requires it.

## Exact First Commands

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
git diff --check
```

If touching deployment helper scripts, also run:

```bash
sh -n scripts/pilot-ops.sh scripts/oci-pilot.sh scripts/collect-hosted-evidence.sh
scripts/pilot-ops.sh help
ENVIRONMENT_NAME=dry-run-demo EVIDENCE_OUTPUT_DIR=/tmp/open-transit-rt-evidence ADMIN_BASE_URL=http://127.0.0.1:8081 ADMIN_TOKEN=redacted scripts/pilot-ops.sh validator-cycle --dry-run
ENVIRONMENT_NAME=dry-run-demo EVIDENCE_OUTPUT_DIR=/tmp/open-transit-rt-evidence DATABASE_URL=postgres://redacted BACKUP_DIR=/tmp/open-transit-rt-backups scripts/pilot-ops.sh backup --dry-run
ENVIRONMENT_NAME=dry-run-demo EVIDENCE_OUTPUT_DIR=/tmp/open-transit-rt-evidence RESTORE_DATABASE_URL=postgres://redacted RESTORE_BACKUP_FILE=/tmp/open-transit-rt-backups/example.dump PUBLIC_BASE_URL=https://feeds.example.org scripts/pilot-ops.sh restore-drill --dry-run
ENVIRONMENT_NAME=dry-run-demo EVIDENCE_OUTPUT_DIR=/tmp/open-transit-rt-evidence PUBLIC_BASE_URL=https://feeds.example.org scripts/pilot-ops.sh feed-monitor --dry-run
ENVIRONMENT_NAME=dry-run-demo EVIDENCE_OUTPUT_DIR=/tmp/open-transit-rt-evidence ADMIN_BASE_URL=http://127.0.0.1:8081 ADMIN_TOKEN=redacted scripts/pilot-ops.sh scorecard-export --dry-run
```

If touching local app/demo docs or scripts, also run:

```bash
make smoke
make agency-app-up
make agency-app-down
make demo-agency-flow
```

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` includes hosted/operator proof for public HTTPS feed fetches, TLS/redirect behavior, auth boundaries, clean hosted validation, monitoring/alert lifecycle, backup/restore, deployment rollback, and scorecard export job history.
- Phase 17 made this repeatable through `scripts/pilot-ops.sh`, systemd timer examples, and updated runbooks. It did not collect new live evidence.
- Phase 15 found real secrets only in ignored local `.cache` files, not in tracked files or history for those `.cache` paths. Rotation/revocation is still required before further real pilot use.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, or raw private operator artifacts.
- `docs/evidence/redaction-policy.md` is the evidence safety rule for future packets.
- Validator success and public fetch proof are not consumer acceptance.
- Consumer-ingestion workflow records are not third-party acceptance.
- Hosted evidence refresh is complete only after `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` passes.

## First Files Likely To Edit

- `docs/phase-18-admin-ux-agency-operations-console.md`
- `cmd/agency-config/`
- `cmd/gtfs-studio/`
- `internal/compliance/`
- `internal/devices/`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-18.md`

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

- Phase 18 — Admin UX And Agency Operations Console

Future roadmap docs:

- `docs/phase-18-admin-ux-agency-operations-console.md`
- `docs/phase-19-realtime-quality-eta-improvement.md`
- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`
