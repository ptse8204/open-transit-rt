# Deployment Evidence Overview (Phase 12 Step 1)

This runbook defines **where** deployment evidence should live and **how** to collect it truthfully.

Phase 12 Step 1 is repo-side scaffolding only. It does not provide hosted deployment evidence by itself.

## Claim Boundaries

Use three evidence buckets in all deployment notes:

1. **Repo-proven capability** (code and local checks in this repository).
2. **Deployment/operator proof** (hosted environment artifacts captured by operators).
3. **Third-party proof** (consumer/aggregator responses such as acceptance emails).

Do not merge these categories in the same claim.

## Artifact Locations

Use `docs/evidence/` as the root.

- `docs/evidence/templates/` contains repo-owned templates/checklists committed to git.
- `docs/evidence/captured/` is the placeholder location for real deployment artifacts collected later.
  - Real evidence may be redacted before committing.
  - If artifacts cannot be committed, store a redacted summary plus a reference pointer in this folder.

## Required Evidence Packs

Collect one pack per deployment environment (for example, `pilot-agency-prod`).

1. Stable HTTPS public feed root proof.
2. Reverse proxy and TLS proof.
3. Production validator records.
4. Monitoring and alerting evidence.
5. Backup/restore evidence.
6. Scorecard export evidence.

Use the runbooks in this folder plus templates in `docs/evidence/templates/`.

## Minimum Naming Convention

Suggested naming under `docs/evidence/captured/<environment>/`:

- `public-feed-proof-YYYY-MM-DD.md`
- `reverse-proxy-tls-YYYY-MM-DD.md`
- `validator-record-YYYY-MM-DD.md`
- `monitoring-alert-YYYY-MM-DD.md`
- `backup-restore-drill-YYYY-MM-DD.md`
- `scorecard-export-YYYY-MM-DD.md`

Keep filenames date-stamped in UTC.

## Required Links

- Reverse proxy/TLS runbook: `docs/runbooks/reverse-proxy-and-tls.md`
- Validator evidence runbook: `docs/runbooks/validator-evidence.md`
- Monitoring/alerting runbook: `docs/runbooks/monitoring-and-alerting.md`
- Backup/restore runbook: `docs/runbooks/backup-and-restore.md`
- Scorecard export runbook: `docs/runbooks/scorecard-export.md`
