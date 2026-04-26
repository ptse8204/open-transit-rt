# Deploy With Docker Compose

This guide describes the current repository-supported Compose paths.

- Default Compose starts Postgres/PostGIS for local development.
- The `app` profile starts the full local demo package behind `http://localhost:8080`.

This is suitable as a small-agency pilot shape, not a complete hosted production package.

For the Phase 17 production-directed pilot operations profile, use [Small-Agency Pilot Operations](../runbooks/small-agency-pilot-operations.md). That profile covers systemd/Caddy deployment, scheduled validators, backup/restore, monitoring, scorecard export, and evidence refresh.

![Agency deployment](../assets/agency-deployment.png)

## Default Compose Database

```bash
docker compose -f deploy/docker-compose.yml config
docker compose -f deploy/docker-compose.yml up -d postgres
```

The database listens on host port `55432` by default and uses `postgis/postgis:16-3.4`.

Apply migrations:

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt?sslmode=disable" \
MIGRATIONS_DIR=db/migrations \
go run ./cmd/migrate up
```

## Full Local App Profile

For a non-expert local trial, use the orchestration helper:

```bash
make agency-app-up
```

This wraps:

```bash
docker compose -f deploy/docker-compose.yml --profile app ...
```

It starts Postgres/PostGIS, builds the local Go service image, applies migrations, seeds demo data, imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, waits for readiness, verifies the public feed URLs, and prints next steps.

Companion commands:

```bash
make agency-app-logs
make agency-app-down
make agency-app-reset
```

`make agency-app-reset` is destructive. It removes local containers, the Compose Postgres volume, local demo database state, and container logs after confirmation. Use `scripts/agency-local-app.sh reset --force` only for automation.

The local app profile uses `deploy/Dockerfile.local` and `deploy/Caddyfile.local`. It does not bake generated credentials or `.cache` material into the image.

## Production-Like Environment

Set real secrets and avoid the dev placeholders from `.env.example`:

```bash
export APP_ENV=production
export BIND_ADDR=127.0.0.1
export DATABASE_URL="postgres://user:password@db.example.internal:5432/open_transit_rt?sslmode=require"
export AGENCY_ID=demo-agency
export ADMIN_JWT_SECRET="replace-with-a-high-entropy-secret"
export ADMIN_JWT_ISSUER="open-transit-rt"
export ADMIN_JWT_AUDIENCE="open-transit-rt-admin"
export CSRF_SECRET="replace-with-a-high-entropy-csrf-secret"
export DEVICE_TOKEN_PEPPER="replace-with-a-high-entropy-device-token-pepper"
export PUBLIC_BASE_URL="https://feeds.example.org"
export FEED_BASE_URL="https://feeds.example.org/public"
export REALTIME_VALIDATION_BASE_URL="https://feeds.example.org/public"
export TECHNICAL_CONTACT_EMAIL="transit-data@example.org"
export FEED_LICENSE_NAME="CC BY 4.0"
export FEED_LICENSE_URL="https://creativecommons.org/licenses/by/4.0/"
export PUBLICATION_ENVIRONMENT=production
```

In `APP_ENV=production`, services fail fast when required secrets are missing or still look like development placeholders.

Do not inline live secrets in systemd units or committed docs. Store deployment values in private environment files such as `/opt/open-transit-rt/env` and `/opt/open-transit-rt/ops/pilot-ops.env`, with real admin tokens, DB passwords, device peppers, and notification credentials kept operator-only.

## Validator Tooling

Install and verify pinned validators on the host or bake equivalent pinned artifacts into your service image:

```bash
make validators-install
make validators-check
```

The repo-supported static validator path is:

```text
.cache/validators/gtfs-validator-7.1.0-cli.jar
```

The repo-supported GTFS-RT validator path is:

```text
.cache/validators/gtfs-rt-validator-wrapper.sh
```

The GTFS-RT wrapper requires Docker access because it runs the pinned MobilityData validator image by digest. A deployment that uses a native GTFS-RT validator executable should add an equivalent checksum/version contract before making stronger validation-readiness claims.

For scheduled pilot validation, run the dry-run first:

```bash
ENVIRONMENT_NAME=pilot-agency-prod \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/$(date -u +%Y-%m-%d) \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-redacted-admin-token \
scripts/pilot-ops.sh validator-cycle --dry-run
```

The live helper writes `validator-cycle-YYYY-MM-DD.json` to `EVIDENCE_OUTPUT_DIR`.

## Run Services

Run each service under your process manager with the shared environment above and a service-specific `PORT`.

```bash
PORT=8081 go run ./cmd/agency-config
PORT=8082 go run ./cmd/telemetry-ingest
PORT=8083 go run ./cmd/feed-vehicle-positions
PORT=8084 go run ./cmd/feed-trip-updates
PORT=8085 go run ./cmd/feed-alerts
PORT=8086 go run ./cmd/gtfs-studio
```

For packaged deployments, build binaries instead of using `go run`:

```bash
go build ./...
```

## Reverse Proxy Boundary

The local app profile includes a Caddy reverse proxy at:

```text
http://localhost:8080
```

That proxy is local-demo convenience only. Admin/debug routes may be routed through it, but they still require admin auth. Production deployments must choose their own admin network boundary, TLS policy, and reverse proxy controls.

Terminate TLS at a reverse proxy and expose only stable public feed paths anonymously:

```text
https://feeds.example.org/public/gtfs/schedule.zip
https://feeds.example.org/public/feeds.json
https://feeds.example.org/public/gtfsrt/vehicle_positions.pb
https://feeds.example.org/public/gtfsrt/trip_updates.pb
https://feeds.example.org/public/gtfsrt/alerts.pb
```

Route those public paths to the appropriate internal service:

| Public path | Internal service |
| --- | --- |
| `/public/gtfs/*` | `agency-config` |
| `/public/feeds.json` | `agency-config` |
| `/public/gtfsrt/vehicle_positions.pb` | `feed-vehicle-positions` |
| `/public/gtfsrt/trip_updates.pb` | `feed-trip-updates` |
| `/public/gtfsrt/alerts.pb` | `feed-alerts` |

Keep admin/debug/JSON paths protected by admin auth and, for production deployments, network or proxy policy:

```text
/admin/*
/admin/debug/*
/public/gtfsrt/*.json
/v1/events
/metrics
```

`/metrics` exists only when `METRICS_ENABLED=true`; treat it as an internal operations surface.

The redacted Caddy example in `deploy/oci/Caddyfile` follows this default-public-only boundary. It relies on Caddy-managed TLS renewal and keeps admin/debug paths absent from the public edge.

## Publish Feed Metadata

After importing or publishing an active GTFS feed, bootstrap publication metadata:

```bash
ADMIN_TOKEN="$(go run ./cmd/admin-token -sub admin@example.com -agency-id demo-agency | sed -n 's/^token=//p')"

curl -fsS -X POST http://localhost:8081/admin/publication/bootstrap \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{}'
```

The server fills unset bootstrap fields from environment variables. `PUBLICATION_ENVIRONMENT=production` makes missing canonical validation evidence red in scorecards.

## Operational Checks

Before treating a deployment as pilot-ready, run:

```bash
make validators-check
make validate
make test
make smoke
```

Then verify deployed URLs externally:

```bash
curl -fsS -o /tmp/schedule.zip "$PUBLIC_BASE_URL/public/gtfs/schedule.zip"
curl -fsS "$PUBLIC_BASE_URL/public/feeds.json"
curl -fsS -o /tmp/vehicle_positions.pb "$PUBLIC_BASE_URL/public/gtfsrt/vehicle_positions.pb"
curl -fsS -o /tmp/trip_updates.pb "$PUBLIC_BASE_URL/public/gtfsrt/trip_updates.pb"
curl -fsS -o /tmp/alerts.pb "$PUBLIC_BASE_URL/public/gtfsrt/alerts.pb"
```

Set `PUBLIC_BASE_URL` to the agency’s actual feed host before running these checks. Do not claim consumer acceptance until the agency has evidence from the specific consumer.

## Phase 17 Evidence Refresh

Deployment helper outputs should use these names in the private operator evidence directory:

- `validator-cycle-YYYY-MM-DD.json`
- `backup-run-YYYY-MM-DD.txt`
- `restore-drill-YYYY-MM-DD.txt`
- `feed-monitor-YYYY-MM-DD.txt`
- `scorecard-export-YYYY-MM-DD.json`

After copying reviewed, redacted artifacts into a hosted packet, finish with:

```bash
EVIDENCE_PACKET_DIR=docs/evidence/captured/<environment>/<UTC-date> make audit-hosted-evidence
```

Do not call refreshed evidence complete unless that audit passes.
