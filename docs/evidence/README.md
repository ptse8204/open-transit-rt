# Deployment Evidence Folder

This folder is the Phase 12 Step 1 evidence scaffold.

## Structure

- `templates/`: repo-owned templates and checklists committed to git.
- `captured/`: location for captured evidence packets. Existing local packets may be partial; hosted deployment artifacts should be added here when available.

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
