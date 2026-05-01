# Production Checklist

This checklist is for production-directed pilots. Completing it is not the same as proving full production readiness for every agency.

For the Phase 11 evidence separation between repo capability, deployment/operator proof, and third-party confirmation, see [Compliance Evidence Checklist](../compliance-evidence-checklist.md).

For the Phase 17 pilot operations profile, see [Small-Agency Pilot Operations](../runbooks/small-agency-pilot-operations.md).

## Runtime Configuration

- Set `APP_ENV=production`.
- Set explicit `DATABASE_URL`; do not rely on the local default.
- Use high-entropy values for `ADMIN_JWT_SECRET`, `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`.
- Set `ADMIN_JWT_ISSUER` and `ADMIN_JWT_AUDIENCE`.
- Set `ADMIN_JWT_OLD_SECRETS` only during intentional secret rotation.
- Keep `BIND_ADDR=127.0.0.1` unless the service is behind a TLS-terminating reverse proxy and network policy.
- Configure `PUBLIC_BASE_URL` and `FEED_BASE_URL` to the stable public HTTPS feed host.
- Configure `TECHNICAL_CONTACT_EMAIL`, `FEED_LICENSE_NAME`, and `FEED_LICENSE_URL`.
- Set `PUBLICATION_ENVIRONMENT=production` only when scorecards should treat missing validator evidence as red.
- Put deployment secrets in private environment files. Do not inline live DB passwords, admin tokens, device peppers, CSRF/JWT secrets, webhook URLs, private keys, or notification credentials in systemd units or committed docs.

## Database

- Provision PostgreSQL 16 with PostGIS support.
- Run migrations through `cmd/migrate`, not by editing `db/schema.sql`.
- Back up the database before imports, publishes, or operational credential rotation.
- Confirm `make migrate-status` reports the expected migration version.

## Public Feed Boundary

Expose these paths anonymously over stable HTTPS:

```text
/public/gtfs/schedule.zip
/public/feeds.json
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Do not require login for public protobuf feeds. Do not change URLs when the underlying feed version changes.

The local `http://localhost:8080` reverse proxy from `make agency-app-up` is demo packaging only. Do not treat it as the production public or admin network boundary.

## Admin And Debug Boundary

Protect these surfaces with admin auth and deployment network controls:

```text
/admin/*
/admin/debug/*
/public/gtfsrt/*.json
/v1/events
/metrics
```

Verify GTFS Studio is protected:

```bash
curl -s -o /dev/null -w '%{http_code}\n' "$PUBLIC_BASE_URL/admin/gtfs-studio"
```

Set `PUBLIC_BASE_URL` to the deployment host before running the command. The expected anonymous response is `401` or another deployment-level denial.

Production deployments must choose their own admin network boundary. The local app profile may route admin/debug paths through `localhost:8080` for convenience, but those routes still require auth and are not a production exposure recommendation.

## Validator Setup

- Run `make validators-install`.
- Run `make validators-check`.
- Configure `GTFS_VALIDATOR_PATH` to the pinned static validator JAR or an equivalently pinned artifact.
- Configure `GTFS_RT_VALIDATOR_PATH` to the pinned wrapper or an equivalently pinned executable.
- Run `/admin/validation/run` for schedule, Vehicle Positions, Trip Updates, and Alerts.
- Store and review validation results before claiming feeds are validator-clean.
- For scheduled pilot validation, dry-run `scripts/pilot-ops.sh validator-cycle --dry-run` before enabling `open-transit-validator-cycle.timer`.

## Publication Workflow

- Import or publish an active GTFS feed.
- For real agency GTFS, complete the [Real Agency GTFS Onboarding](real-agency-gtfs-onboarding.md) intake, metadata approval, redaction review, and publish review before treating the feed as production-directed.
- Confirm approved values for agency name, agency URL, timezone, technical contact email, license name, license URL, public feed root, approver, approval date, and notes.
- Bootstrap publication metadata through `/admin/publication/bootstrap`.
- Verify `/public/feeds.json` lists schedule, Vehicle Positions, Trip Updates, and Alerts.
- Confirm license and contact fields are complete.
- Confirm `schedule.zip` returns `ETag`, `Last-Modified`, and `X-Checksum-SHA256`.
- Confirm realtime feed timestamps and health records are fresh enough for the agency’s operating model.
- If no agency-owned or agency-approved feed root exists, keep final public-feed review limited to local/demo or pilot evidence. Do not claim agency-domain production proof.

## Telemetry And Devices

- Provision one opaque device token per device binding.
- Keep `DEVICE_TOKEN_PEPPER` secret and stable across restarts.
- Use `/admin/devices/rebind` to rotate a device token or change a vehicle binding.
- Confirm old tokens fail immediately after rebinding.
- Confirm telemetry payloads use RFC3339 timestamps with timezone or offset.

## Operations

- Capture service logs with request IDs.
- Keep `/metrics` internal if `METRICS_ENABLED=true`.
- Treat Prometheus/Grafana and OpenTelemetry as deployment-owned or future integrations; this repo does not currently ship dashboards, alert rules, collectors, exporters, or tracing configuration.
- Monitor readiness endpoints for database and active-feed dependencies.
- For lightweight pilot monitoring, dry-run `scripts/pilot-ops.sh feed-monitor --dry-run` before enabling `open-transit-feed-monitor.timer`.
- If webhook or email notification destinations are missing, record `notification not configured`; do not count that as a feed failure.
- Do not commit real webhook URLs, email credentials, or notification tokens.
- Dry-run `scripts/pilot-ops.sh backup --dry-run`, `restore-drill --dry-run`, and `scorecard-export --dry-run` before enabling scheduled operations. Restore is destructive and requires typed confirmation unless `--force` is passed.
- Define who can use admin roles: `admin`, `editor`, `operator`, and `read_only`.
- Keep audit logs for imports, publication metadata changes, alert edits, validation runs, overrides, and device rebinding.

## Evidence Before Stronger Claims

Before claiming a deployment is compliant or consumer-ready, collect evidence for:

- stable public HTTPS URLs
- successful static GTFS validation
- successful GTFS-RT validation for Vehicle Positions, Trip Updates, and Alerts
- complete license and technical contact metadata
- current `/public/feeds.json`
- scorecard status and validation history
- external consumer submission or acceptance evidence, if claimed

Without that evidence, use wording such as "supports deployment toward CAL-ITP/Caltrans-style readiness."

Hosted evidence refresh is complete only after:

```bash
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

The audit must pass before the packet is described as complete.
