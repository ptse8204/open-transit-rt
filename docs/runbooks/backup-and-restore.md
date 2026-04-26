# Runbook: Backup and Restore Evidence

This runbook captures evidence that deployment data can be backed up and restored with a repeatable operator process.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/backup-restore-drill-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/backup-restore-drill-2026-04-24.md`

The local packet records a one-time Postgres dump, isolated restore database, restored row counts, and public feed fetches against the restored database. The OCI pilot packet records deployment/operator proof for that recorded pilot scope.

## Backup Evidence

Record:

- Backup schedule and cadence.
- Retention period.
- Backup storage location and access boundary.
- Verification that jobs actually completed.

Phase 17 backup helper dry-run:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
DATABASE_URL=postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt?sslmode=disable \
BACKUP_DIR=/opt/open-transit-rt/backups \
scripts/pilot-ops.sh backup --dry-run
```

The live backup writes `backup-run-YYYY-MM-DD.txt` to `EVIDENCE_OUTPUT_DIR` and a private dump under `BACKUP_DIR`.

## Restore Procedure Evidence

Record step-by-step restore instructions with command placeholders for the deployment.

At minimum include:

1. How to isolate the target environment.
2. How to restore database snapshots/backups.
3. How to verify feed-serving integrity after restore.
4. How to run validator checks post-restore.

Restore helpers are destructive for `RESTORE_DATABASE_URL`. They must warn clearly and require typed confirmation unless `--force` is passed:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
RESTORE_DATABASE_URL=postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt_restore?sslmode=disable \
RESTORE_BACKUP_FILE=/opt/open-transit-rt/backups/open-transit-rt-YYYYMMDDTHHMMSSZ.dump \
PUBLIC_BASE_URL=https://feeds.example.org \
scripts/pilot-ops.sh restore-drill --dry-run
```

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

Evidence labels:

- Raw database dumps and checksum files in `BACKUP_DIR`: `never-commit`.
- `backup-run-YYYY-MM-DD.txt`: `private/operator-only`; `safe-to-commit-after-review` if paths and private infrastructure details are redacted.
- `restore-drill-YYYY-MM-DD.txt`: `private/operator-only`; redacted summaries may be committed after review.
