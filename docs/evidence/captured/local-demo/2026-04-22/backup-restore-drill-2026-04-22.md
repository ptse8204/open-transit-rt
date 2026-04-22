# Backup and Restore Drill

- Environment: `local-demo`
- Drill date (UTC): 2026-04-22
- Operator: Codex local run

## Backup Posture

- Backup schedule: Missing for production. This packet used a one-time manual `pg_dump`.
- Retention policy: Missing for production. The committed packet stores restore transcript and checksums, not a production backup archive.
- Storage location: local temporary dump path `/var/folders/_g/bvzl9cms7cx1d0wdpc981n9w0000gn/T//open_transit_rt_phase12_20260422.sql`.
- Access boundary: developer workstation plus local Docker Postgres container.
- Last successful backup timestamp (UTC): 2026-04-22T21:37:02Z.

## Restore Procedure

- Isolation steps/commands: restored into separate database `open_transit_rt_restore_drill_20260422`.
- Database restore steps/commands:

```sh
docker compose -f deploy/docker-compose.yml exec -T postgres pg_dump -U postgres -d open_transit_rt --no-owner --no-privileges > "$DUMP"
docker compose -f deploy/docker-compose.yml exec -T postgres dropdb -U postgres --if-exists open_transit_rt_restore_drill_20260422
docker compose -f deploy/docker-compose.yml exec -T postgres createdb -U postgres open_transit_rt_restore_drill_20260422
docker compose -f deploy/docker-compose.yml exec -T postgres psql -U postgres -d open_transit_rt_restore_drill_20260422 < "$DUMP"
```

- Feed-serving verification steps/commands: started local services with `DATABASE_URL=postgres://postgres:postgres@localhost:55432/open_transit_rt_restore_drill_20260422?sslmode=disable` and fetched all five public paths through the local proxy.
- Post-restore validator steps/commands: not rerun against the restored database because validator execution was already failing for environment/tooling reasons recorded in `validator-record-2026-04-22.md`.
- Command transcript location: this file plus `artifacts/logs/restore-psql.log`.

## Restore Drill

- Backup artifact used: local plain SQL dump.
- Backup artifact checksum: `1291bea492176ce891b52fa83d7ed650b897d824037f5cc817e8084723300279`.
- Restore start (UTC): 2026-04-22T21:37:02Z.
- Restore finish (UTC): 2026-04-22T21:37:04Z.
- Duration: approximately 2 seconds.
- Outcome: local restore succeeded into the isolated restore database.

## Post-Restore Verification

Restored row counts:

- `feed_version_count=7`
- `validation_report_count=21`
- `service_alert_count=1`
- `telemetry_event_count=6`

Public feed fetch checks against restored database at 2026-04-22T21:38:44Z:

| Path | Status | Bytes | SHA-256 |
| --- | ---: | ---: | --- |
| `/public/gtfs/schedule.zip` | 200 | 1960 | `0956ed037a40ca9d2cca94a501bea1547d27dbd25a195c7ebefe1a34ffc78194` |
| `/public/feeds.json` | 200 | 2442 | `1a3708f96ab6c6ad1d52c7dd768ef8f9f3015424ca40f243577ddb41d9120ef8` |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | 63 | `7cd7a28ec7bf3fdc50b1c1f1faf7aa20d111a4c3aa8146259d8167b9286a1ebf` |
| `/public/gtfsrt/trip_updates.pb` | 200 | 15 | `56376073febc3856b1e33f4bbd779f14b0c67211387d0fc5344e21a399b2c93b` |
| `/public/gtfsrt/alerts.pb` | 200 | 135 | `79140b8bb93d82db9e41b12e6051631f645631f4add90e0dbd6ccbb86d1b6373` |

- Validator check summary: not rerun post-restore; known environment validator failures remain blockers.
- Known gaps after restore: no production backup schedule, no retention evidence, no encrypted storage/access evidence, and no production outage runbook exercise.

## Outage / Validator-Failure Response Notes

- Feed outage runbook link: missing production incident runbook.
- Validator failure runbook link: missing production incident runbook.
- Escalation path: missing.
- Rollback criteria: local rollback not exercised; production rollback criteria remain missing deployment evidence.
