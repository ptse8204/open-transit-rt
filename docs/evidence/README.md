# Evidence Folder

This folder holds repo-owned evidence templates, captured evidence packets, and later-phase evidence scaffolds.

## Structure

- `templates/`: repo-owned templates and checklists committed to git.
- `captured/`: location for captured evidence packets. Existing local packets may be partial; hosted deployment artifacts should be added here when available.
- `real-agency-gtfs/`: templates for future real-agency GTFS import review packets. It contains templates only until real agency-approved, public-safe evidence exists.
- `device-avl/`: templates for future device, GPS emitter, or vendor AVL integration review packets. It contains templates only until public-safe reviewed evidence exists.
- `device-avl/templates/integration-review-template.md`: the reusable template for future device, GPS emitter, sidecar, or vendor AVL review records, including Phase 29B dry-run adapter review fields.

Current captured packet:

- `captured/local-demo/2026-04-22/`: real local demo evidence packet. It is not hosted HTTPS production proof.
- `captured/hosted-pending/2026-04-22/`: hosted evidence intake packet with command artifacts and pending fields. It is not proof until an operator replaces pending entries with real hosted outputs.
- `captured/oci-pilot/2026-04-24/`: hosted OCI pilot packet with public feed, TLS, auth-boundary, validation, monitoring, backup/restore, rollback, and scorecard evidence for the recorded pilot scope.

Supporting hygiene docs:

- `redaction-policy.md`: rules for public-safe evidence, required redactions, checksum refreshes, and secret response.
- `archive-inventory.md`: committed archive inventory for `docs/evidence/captured/**`.

## Important

Do not fabricate evidence.

If real deployment artifacts are not yet collected, leave placeholders and mark status as pending.

## Suggested Workflow

1. Read `docs/runbooks/deployment-evidence-overview.md`.
2. Use runbook-specific templates from `templates/`.
3. Save environment-specific outputs under `captured/<environment>/`.
4. Redact sensitive details as needed and note all redactions.
5. Keep claims aligned with available evidence.

Before committing a new captured packet, review it against `redaction-policy.md`
and add every committed archive to `archive-inventory.md`.

## Phase 17 Evidence Refresh

Pilot operations helpers write private/operator-owned outputs to `EVIDENCE_OUTPUT_DIR` before any public evidence packet is assembled:

- `validator-cycle-YYYY-MM-DD.json`
- `backup-run-YYYY-MM-DD.txt`
- `restore-drill-YYYY-MM-DD.txt`
- `feed-monitor-YYYY-MM-DD.txt`
- `scorecard-export-YYYY-MM-DD.json`

Raw backups, admin tokens, database URLs with passwords, TLS private material, webhook URLs, notification credentials, and private operator artifacts are never committed. Redacted summaries may be copied into `captured/<environment>/<UTC-date>/` only after review.

Every hosted evidence refresh must end with:

```sh
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

Do not call refreshed evidence complete unless that audit passes.

## Phase 28 Operations Evidence Notes

Use `docs/runbooks/production-operations-hardening.md` for day-to-day operations, alert delivery proof, capacity checks, incident response, secret rotation, restore events, and operator handover.

Template files under `docs/runbooks/templates/` are not evidence by themselves. Do not commit fake incidents, fake alert delivery proof, fake rotation records, fake restore events, or placeholder operational artifacts.

Current backup, restore, export, and evidence workflows are deployment/DB scoped. They are not tenant-safe multi-agency workflows. Phase 27 selected isolation tests do not prove production multi-tenant operations or tenant-safe backup/restore/export/evidence handling.
