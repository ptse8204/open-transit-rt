# Hosted Backup And Restore Evidence

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials

## Backup Posture

- Backup system: `pg_dump -Fc` through `open-transit-backup.service`.
- Schedule: daily `open-transit-backup.timer`.
- Retention: 7 days for timer-created `open-transit-rt-*.dump` files on the pilot host.
- Storage: `/opt/open-transit-rt/backups` on the OCI pilot host.
- Access boundary: root/postgres on the pilot host.

## Restore Drill

- Backup artifact used: `/opt/open-transit-rt/backups/phase12-clean-20260424.dump`.
- Restore transcript: `artifacts/operator-supplied/restore-live-transcript.txt`.
- Backup policy and job history: `artifacts/operator-supplied/backup-policy-and-job-history.txt`.

## Outcome

The drill restored the live pilot database from the clean backup after a transient update. Services restarted successfully, public `feeds.json` returned HTTP 200, and the active feed returned from `gtfs-import-4` to `gtfs-import-3` with stable public URLs.
