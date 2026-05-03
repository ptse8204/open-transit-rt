# Runbook: Backup and Restore Evidence

This runbook captures evidence that deployment data can be backed up and restored with a repeatable operator process.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/backup-restore-drill-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/backup-restore-drill-2026-04-24.md`

The local packet records a one-time Postgres dump, isolated restore database, restored row counts, and public feed fetches against the restored database. The OCI pilot packet records deployment/operator proof for that recorded pilot scope.

## Phase 27 Operations Boundary

Current backup and restore helpers are deployment/DB scoped. They are appropriate for separate single-agency deployments and pilot environments, but they are not tenant-safe multi-agency workflows. Phase 27 selected isolation tests prove repository-level isolation for selected paths only; they do not prove production multi-tenant operations or tenant-safe backup/restore/export/evidence handling.

## Backup Evidence

Record:

- Backup schedule and cadence.
- Retention period.
- Backup storage location and access boundary.
- Verification that jobs actually completed.

Default small-pilot cadence:

- Run database backups daily.
- Keep at least 7 days of backups unless agency policy requires longer retention.
- Take and verify a backup before every source tag, binary, image, or migration upgrade.
- Review backup storage growth weekly and retention cleanup monthly.
- Keep raw dumps, checksum files, private backup paths, and DB URLs with passwords outside public evidence.

Phase 17 backup helper dry-run:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
DATABASE_URL=postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt?sslmode=disable \
BACKUP_DIR=/opt/open-transit-rt/backups \
scripts/pilot-ops.sh backup --dry-run
```

The live backup writes `backup-run-YYYY-MM-DD.txt` to `EVIDENCE_OUTPUT_DIR` and a private dump under `BACKUP_DIR`.

Backup evidence should record a redacted backup path, completion timestamp, checksum, retention policy, and operator-reviewed access boundary. If the backup path reveals sensitive infrastructure layout, commit only a redacted summary.

## Restore Procedure Evidence

Record step-by-step restore instructions with command placeholders for the deployment.

At minimum include:

1. How to isolate the target environment.
2. How to restore database snapshots/backups.
3. How to verify feed-serving integrity after restore.
4. How to run validator checks post-restore.

Restore drills should run monthly for long-running pilots and before stronger operations wording is considered. Restore over a live database only during an approved incident response. Use an isolated restore database for routine drills.

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
- Backup checksum verification.
- Restore target and isolation method.
- Restore duration and outcome.
- Public feed fetch checks after restore.
- Schedule and realtime validation checks after restore.
- Follow-up issues.

Use `docs/runbooks/templates/restore-event-template.md` for incident restores, rollback restores, and restore drills that need a retained operator record.

## Outage and Validator-Failure Notes

Maintain response notes for:

- public feed outage
- validator failure in production run

Include escalation path and rollback criteria.

If restore is part of outage mitigation, record the rollback decision, migration state before and after restore, post-restore public feed checks, post-restore validator checks, evidence retained, and redaction review.

## Output Artifact

Use:

- `docs/evidence/templates/backup-restore-drill-template.md`

Store under `docs/evidence/captured/<environment>/backup-restore-drill-YYYY-MM-DD.md`.

Evidence labels:

- Raw database dumps and checksum files in `BACKUP_DIR`: `never-commit`.
- `backup-run-YYYY-MM-DD.txt`: `private/operator-only`; `safe-to-commit-after-review` if paths and private infrastructure details are redacted.
- `restore-drill-YYYY-MM-DD.txt`: `private/operator-only`; redacted summaries may be committed after review.
- DB URLs with passwords, private backup paths, private restore target names, raw logs, and private operator artifacts are `never-commit` unless fully redacted into a public-safe summary.
