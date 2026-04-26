# Runbook: Small-Agency Pilot Operations

This runbook is the Phase 17 deployment profile for a small-agency pilot. It turns the OCI pilot path into repeatable operations without claiming universal production readiness, CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, or hosted SaaS availability.

The Phase 16 local app package remains a local evaluation path. It uses `http://localhost:8080`, development defaults, and local container networking. This pilot profile uses deployment-owned DNS, HTTPS/TLS, secrets, backup retention, monitoring, and evidence review.

## Deployment Shape

- Runtime: compiled Go binaries, PostgreSQL/PostGIS, Caddy or equivalent reverse proxy, and systemd services/timers.
- Public edge: expose only stable anonymous feed paths:
  - `/public/gtfs/schedule.zip`
  - `/public/feeds.json`
  - `/public/gtfsrt/vehicle_positions.pb`
  - `/public/gtfsrt/trip_updates.pb`
  - `/public/gtfsrt/alerts.pb`
- Private/internal surfaces: keep `/admin/*`, `/admin/debug/*`, `/public/gtfsrt/*.json`, `/v1/events`, and `/metrics` behind admin auth plus deployment network controls.
- Operator helper: `scripts/pilot-ops.sh` provides `validator-cycle`, `backup`, `restore-drill`, `feed-monitor`, and `scorecard-export` subcommands. Every subcommand supports `--dry-run` and prints its target environment before doing work.

## Environment Variable Matrix

Use a private environment file such as `/opt/open-transit-rt/env` for service runtime and `/opt/open-transit-rt/ops/pilot-ops.env` for scheduled operations. Keep both `0600` and operator-only. Do not commit real values.

| Variable | Required | Example value | Secret | Used by |
| --- | --- | --- | --- | --- |
| `APP_ENV` | Required | `production` | No | all services |
| `BIND_ADDR` | Required | `127.0.0.1` | No | all services |
| `DATABASE_URL` | Required | `postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt?sslmode=disable` | Yes | all services, `pilot-ops.sh backup` |
| `MIGRATIONS_DIR` | Required | `/opt/open-transit-rt/app/db/migrations` | No | `cmd/migrate` |
| `AGENCY_ID` | Required | `demo-agency` | No | feed services, admin workflows, scripts |
| `PUBLIC_BASE_URL` | Required | `https://feeds.example.org` | No | feed metadata, `pilot-ops.sh feed-monitor`, evidence collection |
| `FEED_BASE_URL` | Required | `https://feeds.example.org/public` | No | publication metadata |
| `VEHICLE_POSITIONS_FEED_URL` | Required | `https://feeds.example.org/public/gtfsrt/vehicle_positions.pb` | No | Trip Updates adapter |
| `SCHEDULE_FEED_URL` | Optional | `https://feeds.example.org/public/gtfs/schedule.zip` | No | operator checks |
| `REALTIME_VALIDATION_BASE_URL` | Required for hosted validation fallback | `https://feeds.example.org/public` | No | agency-config validation |
| `TRIP_UPDATES_FEED_URL` | Optional | `https://feeds.example.org/public/gtfsrt/trip_updates.pb` | No | validation fallback |
| `ALERTS_FEED_URL` | Optional | `https://feeds.example.org/public/gtfsrt/alerts.pb` | No | validation fallback |
| `TECHNICAL_CONTACT_EMAIL` | Required | `transit-data@example.org` | No | publication metadata, scorecard |
| `FEED_LICENSE_NAME` | Required | `CC BY 4.0` | No | publication metadata, scorecard |
| `FEED_LICENSE_URL` | Required | `https://creativecommons.org/licenses/by/4.0/` | No | publication metadata, scorecard |
| `PUBLICATION_ENVIRONMENT` | Required | `production` | No | scorecard severity |
| `ADMIN_JWT_SECRET` | Required | generated high-entropy value | Yes | admin auth |
| `ADMIN_JWT_OLD_SECRETS` | Optional | generated high-entropy value during rotation | Yes | admin auth |
| `ADMIN_JWT_ISSUER` | Required | `open-transit-rt-pilot` | No | admin auth |
| `ADMIN_JWT_AUDIENCE` | Required | `open-transit-rt-admin` | No | admin auth |
| `ADMIN_JWT_TTL` | Optional | `8h` | No | admin auth |
| `CSRF_SECRET` | Required | generated high-entropy value | Yes | admin browser flows |
| `DEVICE_TOKEN_PEPPER` | Required | generated high-entropy value | Yes | telemetry ingest, device rebind |
| `METRICS_ENABLED` | Optional | `false` | No | services |
| `LOG_LEVEL` | Optional | `info` | No | services |
| `GTFS_VALIDATOR_PATH` | Required for canonical static validation | `/opt/open-transit-rt/.cache/validators/gtfs-validator-7.1.0-cli.jar` | No | agency-config validation |
| `GTFS_VALIDATOR_VERSION` | Required for evidence | `v7.1.0` | No | validation evidence |
| `GTFS_RT_VALIDATOR_PATH` | Required for canonical realtime validation | `/opt/open-transit-rt/.cache/validators/gtfs-rt-validator-wrapper.sh` | No | agency-config validation |
| `GTFS_RT_VALIDATOR_VERSION` | Required for evidence | pinned MobilityData image digest | No | validation evidence |
| `GTFS_RT_VALIDATOR_READY_TIMEOUT_SECONDS` | Optional | `60` | No | realtime validator wrapper |
| `VALIDATOR_TOOLING_MODE` | Required | `pinned` | No | validation checks |
| `ENVIRONMENT_NAME` | Required for operational helpers | `pilot-agency-prod` | No | `pilot-ops.sh`, evidence collection |
| `EVIDENCE_OUTPUT_DIR` | Required for operational helpers | `/opt/open-transit-rt/evidence/2026-04-26` | No | `pilot-ops.sh` |
| `ADMIN_BASE_URL` | Required for admin helper operations | `http://127.0.0.1:8081` | No, if loopback; private if internal host | `pilot-ops.sh validator-cycle`, `scorecard-export` |
| `ADMIN_TOKEN` | Required for admin helper operations | `replace-with-redacted-admin-token` | Yes | `pilot-ops.sh validator-cycle`, `scorecard-export` |
| `BACKUP_DIR` | Required for backups | `/opt/open-transit-rt/backups` | Private path | `pilot-ops.sh backup` |
| `BACKUP_RETENTION_DAYS` | Optional | `7` | No | `pilot-ops.sh backup` |
| `RESTORE_DATABASE_URL` | Required for restore drill | `postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt_restore?sslmode=disable` | Yes | `pilot-ops.sh restore-drill` |
| `RESTORE_BACKUP_FILE` | Required for restore drill | `/opt/open-transit-rt/backups/open-transit-rt-20260426T000000Z.dump` | Private path | `pilot-ops.sh restore-drill` |
| `NOTIFY_WEBHOOK_URL` | Optional | `replace-with-redacted-webhook-url` | Yes | `pilot-ops.sh feed-monitor` |
| `NOTIFY_EMAIL_TO` | Optional | `ops@example.org` | Private | `pilot-ops.sh feed-monitor` |
| `CAPTURE_DATE_UTC` | Optional | `2026-04-26` | No | evidence helpers |

## Evidence Output Locations

Default operator-owned locations:

- Private backups: `/opt/open-transit-rt/backups/` (`private/operator-only`, never commit raw dumps).
- Runtime operation evidence: `/opt/open-transit-rt/evidence/<UTC-date>/` (`private/operator-only` until reviewed).
- Redacted public evidence packets: `docs/evidence/captured/<environment>/<UTC-date>/` (`safe-to-commit-after-review`).
- Environment files, admin tokens, webhook URLs, private DB URLs, TLS private material, and raw access logs are `never-commit`.

Naming conventions:

- `validator-cycle-YYYY-MM-DD.json`
- `backup-run-YYYY-MM-DD.txt`
- `restore-drill-YYYY-MM-DD.txt`
- `feed-monitor-YYYY-MM-DD.txt`
- `scorecard-export-YYYY-MM-DD.json`

Review every artifact against `docs/evidence/redaction-policy.md` before committing. If an artifact contains credentials, private keys, raw public client IP logs, private hostnames, or notification credentials, keep it private and commit only a redacted summary.

## Operational Helper Dry Runs

Run dry-runs from the target operator environment before enabling timers:

```sh
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-redacted-admin-token \
./scripts/pilot-ops.sh validator-cycle --dry-run
```

```sh
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
DATABASE_URL=postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt?sslmode=disable \
BACKUP_DIR=/opt/open-transit-rt/backups \
./scripts/pilot-ops.sh backup --dry-run
```

```sh
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
RESTORE_DATABASE_URL=postgres://open_transit:REDACTED@127.0.0.1:5432/open_transit_rt_restore?sslmode=disable \
RESTORE_BACKUP_FILE=/opt/open-transit-rt/backups/open-transit-rt-YYYYMMDDTHHMMSSZ.dump \
PUBLIC_BASE_URL=https://feeds.example.org \
./scripts/pilot-ops.sh restore-drill --dry-run
```

```sh
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
PUBLIC_BASE_URL=https://feeds.example.org \
./scripts/pilot-ops.sh feed-monitor --dry-run
```

```sh
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-redacted-admin-token \
./scripts/pilot-ops.sh scorecard-export --dry-run
```

The restore drill is destructive for `RESTORE_DATABASE_URL`. Live restore requires typed confirmation `restore <ENVIRONMENT_NAME>` unless `--force` is passed by automation.

## Systemd Timer Examples

Example unit/timer files live under `deploy/systemd/`:

- `open-transit-validator-cycle.service` / `.timer`
- `open-transit-backup.service` / `.timer`
- `open-transit-feed-monitor.service` / `.timer`
- `open-transit-scorecard-export.service` / `.timer`

These examples load secrets and target values through `EnvironmentFile={{OCI_REMOTE_DIR}}/ops/pilot-ops.env`. Do not inline live tokens, DB passwords, webhook URLs, or notification credentials in unit files.

Enable timers only after dry-run checks pass and the private environment file has real deployment-owned values:

```sh
sudo systemctl enable --now open-transit-feed-monitor.timer
sudo systemctl enable --now open-transit-backup.timer
sudo systemctl enable --now open-transit-validator-cycle.timer
sudo systemctl enable --now open-transit-scorecard-export.timer
```

## Evidence Refresh Closure

After collecting and redacting a hosted evidence packet, run:

```sh
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

Do not call refreshed evidence complete unless the audit passes. A passing audit proves packet completeness for the documented evidence checklist only; it is not consumer acceptance or full compliance.
