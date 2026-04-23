# Deployment Evidence Overview (Phase 12 Step 1)

This runbook defines **where** deployment evidence should live and **how** to collect it truthfully.

Phase 12 Step 1 is repo-side scaffolding only. It does not provide hosted deployment evidence by itself.

## Captured Evidence Packets

- `docs/evidence/captured/local-demo/2026-04-22/`: Phase 12 Step 2 local demo packet. This packet contains real local artifacts for loopback HTTP feed retrieval, local validator run records, a local restore drill, and manual scorecard export. It does **not** prove hosted HTTPS deployment readiness, clean validator status, production monitoring/alerting, production backup retention, or consumer acceptance.

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

Generic hosted collection can be started with:

```sh
ENVIRONMENT_NAME="<hosted-environment>" \
PUBLIC_BASE_URL="https://<canonical-feed-host>" \
ADMIN_BASE_URL="https://<admin-or-origin-host>" \
ADMIN_TOKEN="<redacted-admin-token>" \
make collect-hosted-evidence
```

This collects feed fetches, TLS headers/certificate details, admin validation runs, and a manual scorecard export. Operators must still attach deployment-owned monitoring, alert lifecycle, backup/restore, reverse proxy renewal, and scheduler/job-history artifacts.

Before using the hosted validator collection path, run `make validators-install` and `make validators-check` on the collection host. The pinned static validator requires Java 17+; the pinned Docker-backed GTFS-RT wrapper requires Docker, `curl`, and `python3`.

After completing a hosted packet, audit it before making closure claims:

```sh
EVIDENCE_PACKET_DIR="docs/evidence/captured/<hosted-environment>/<UTC-date>" \
make audit-hosted-evidence
```

The audit fails while pending markers, failed validators, missing public artifacts, missing TLS redirect/certificate evidence, or missing operator-supplied monitoring/backup/scheduler artifacts remain.

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
