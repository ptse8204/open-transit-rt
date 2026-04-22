# Runbook: Backup and Restore Evidence

This runbook captures evidence that deployment data can be backed up and restored with a repeatable operator process.

## Backup Evidence

Record:

- Backup schedule and cadence.
- Retention period.
- Backup storage location and access boundary.
- Verification that jobs actually completed.

## Restore Procedure Evidence

Record step-by-step restore instructions with command placeholders for the deployment.

At minimum include:

1. How to isolate the target environment.
2. How to restore database snapshots/backups.
3. How to verify feed-serving integrity after restore.
4. How to run validator checks post-restore.

## Restore Drill Evidence

For each drill capture:

- UTC timestamp.
- Operator identity/role.
- Backup source used.
- Restore duration and outcome.
- Validation/fetch checks after restore.
- Follow-up issues.

## Outage and Validator-Failure Notes

Maintain response notes for:

- public feed outage
- validator failure in production run

Include escalation path and rollback criteria.

## Output Artifact

Use:

- `docs/evidence/templates/backup-restore-drill-template.md`

Store under `docs/evidence/captured/<environment>/backup-restore-drill-YYYY-MM-DD.md`.
