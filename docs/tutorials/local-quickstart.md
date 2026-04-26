# Local Quickstart

This tutorial starts the current Open Transit RT development environment on one machine. It is local-dev only and uses seeded demo credentials.

![Exact-behavior local quickstart flow showing database bootstrap, validator install, sample GTFS import, service startup, publication bootstrap, telemetry ingest, public feed fetches, validation run, and scorecard inspection.](../assets/quickstart-flow.png)

*Exact-behavior flow diagram for the committed local demo path, rendered from a reviewed SVG spec.*

## Prerequisites

- Docker with Compose support
- `curl`
- Go matching `go.mod` for development commands outside the local app package
- `zip` and `unzip` for the legacy demo script
- Java if you want the static GTFS validator JAR to execute successfully

The GTFS-RT validator workflow uses Docker because the repo-supported path is a pinned wrapper around a MobilityData container image.

## Full Local App Package

The simplest agency/evaluator path is:

```bash
make agency-app-up
```

This starts the full local stack behind:

```text
http://localhost:8080
```

The command imports `testdata/gtfs/valid-small`, publishes it as the active local feed, bootstraps publication metadata, verifies feed discovery, verifies public protobuf feed URLs, and prints admin/device/log/validation next steps.

Useful companion commands:

```bash
make agency-app-logs
make agency-app-down
make agency-app-reset
```

`make agency-app-reset` is destructive and prompts before removing local containers, the Compose Postgres volume, local demo database state, and container logs.

The local reverse proxy is only for demo packaging. Admin/debug routes may be reachable through `localhost:8080`, but they still require auth. Production deployments need HTTPS/TLS and deployment-owned admin network boundaries.

## Development Bootstrap

```bash
cp .env.example .env
make dev
```

`make dev` starts Postgres/PostGIS, applies migrations, seeds demo agencies, and prints a local admin token plus the seeded telemetry device token.

Seeded local credentials:

```text
agency_id=demo-agency
admin subject=admin@example.com
device_id=device-1
vehicle_id=bus-1
device token=dev-device-token
```

Install and verify pinned validators:

```bash
make validators-install
make validators-check
```

## Run The Demo

The fastest end-to-end check is:

```bash
make demo-agency-flow
```

This script imports sample GTFS, starts all current services, creates a temporary local public proxy, bootstraps publication metadata, ingests authenticated telemetry, fetches public feeds, verifies protected debug/admin access, runs validator flows, and reads the scorecard plus consumer-ingestion records.

The demo explicitly verifies:

- `/public/gtfs/schedule.zip`
- `/public/feeds.json`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`
- protected JSON debug routes
- protected `/admin/gtfs-studio` access
- protected GTFS Studio draft subroute access

## Manual Service Startup

To run services manually, use separate terminals:

```bash
make run-agency-config
make run-telemetry-ingest
make run-feed-vehicle-positions
make run-feed-trip-updates
make run-feed-alerts
make run-gtfs-studio
```

Default local service ports:

| Service | URL |
| --- | --- |
| agency-config | `http://localhost:8081` |
| telemetry-ingest | `http://localhost:8082` |
| feed-vehicle-positions | `http://localhost:8083` |
| feed-trip-updates | `http://localhost:8084` |
| feed-alerts | `http://localhost:8085` |
| gtfs-studio | `http://localhost:8086` |

## Import Sample GTFS

The runtime importer accepts a ZIP file. Create one from the committed fixture:

```bash
tmp_zip="$(mktemp "${TMPDIR:-/tmp}/open-transit-rt-gtfs.XXXXXX.zip")"
(cd testdata/gtfs/valid-small && zip -qr "$tmp_zip" .)
go run ./cmd/gtfs-import -agency-id demo-agency -zip "$tmp_zip" -actor-id quickstart -notes "local quickstart"
```

## Generate An Admin Token

```bash
export ADMIN_JWT_SECRET=dev-admin-jwt-secret-change-me
export ADMIN_JWT_ISSUER=open-transit-rt-local
export ADMIN_JWT_AUDIENCE=open-transit-rt-admin
export CSRF_SECRET=dev-csrf-secret-change-me
export DEVICE_TOKEN_PEPPER=dev-device-token-pepper-change-me

ADMIN_TOKEN="$(go run ./cmd/admin-token -sub admin@example.com -agency-id demo-agency | sed -n 's/^token=//p')"
```

## Bootstrap Publication Metadata

Run `agency-config` first, then:

```bash
curl -fsS -X POST http://localhost:8081/admin/publication/bootstrap \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{}'
```

The endpoint derives defaults from environment variables such as `PUBLIC_BASE_URL`, `FEED_BASE_URL`, `TECHNICAL_CONTACT_EMAIL`, `FEED_LICENSE_NAME`, `FEED_LICENSE_URL`, and `PUBLICATION_ENVIRONMENT`.

## Ingest Telemetry

The seeded device token works only for `demo-agency`, `device-1`, and `bus-1`.

```bash
observed_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
device_token="dev-device-token"
curl -fsS -X POST http://localhost:8082/v1/telemetry \
  -H "Authorization: Bearer ${device_token}" \
  -H "Content-Type: application/json" \
  --data "{
    \"agency_id\": \"demo-agency\",
    \"device_id\": \"device-1\",
    \"vehicle_id\": \"bus-1\",
    \"timestamp\": \"$observed_at\",
    \"lat\": 49.2827,
    \"lon\": -123.1207,
    \"bearing\": 120.0,
    \"speed_mps\": 8.4,
    \"accuracy_m\": 7.5,
    \"trip_hint\": \"trip-10-0800\"
  }"
```

Telemetry ingest persists the event if the token and agency/device/vehicle binding match. Rebinding a device through `/admin/devices/rebind` rotates the token and immediately invalidates the old binding.

## Fetch Feeds

Without a reverse proxy, local services expose their own ports:

```bash
curl -fsS -o /tmp/schedule.zip http://localhost:8081/public/gtfs/schedule.zip
unzip -t /tmp/schedule.zip

curl -fsS http://localhost:8081/public/feeds.json
curl -fsS -o /tmp/vehicle_positions.pb http://localhost:8083/public/gtfsrt/vehicle_positions.pb
curl -fsS -o /tmp/trip_updates.pb http://localhost:8084/public/gtfsrt/trip_updates.pb
curl -fsS -o /tmp/alerts.pb http://localhost:8085/public/gtfsrt/alerts.pb
```

The protobuf feeds may be valid and empty if no current assignment, prediction, or published alert is available.

## Verify Protected Access

Anonymous requests to debug/admin routes should fail:

```bash
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8086/admin/gtfs-studio
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8086/admin/gtfs-studio/drafts/demo-draft
curl -s -o /dev/null -w '%{http_code}\n' http://localhost:8083/public/gtfsrt/vehicle_positions.json
```

With the admin token, GTFS Studio and JSON debug routes are available:

```bash
curl -fsS -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8086/admin/gtfs-studio
curl -fsS -H "Authorization: Bearer $ADMIN_TOKEN" http://localhost:8083/public/gtfsrt/vehicle_positions.json
```

## Run Checks

```bash
make validate
make test
make smoke
```

For DB-backed integration coverage:

```bash
make test-integration
```

Docker must be running for DB-backed integration tests and the pinned GTFS-RT validator wrapper.

## Device Helper

Use the helper for local device onboarding and simulator-style telemetry:

```bash
scripts/device-onboarding.sh help
scripts/device-onboarding.sh sample --dry-run
scripts/device-onboarding.sh sample
scripts/device-onboarding.sh simulate --dry-run
scripts/device-onboarding.sh simulate
```

To rotate or bind a demo device token through the existing API:

```bash
scripts/device-onboarding.sh rebind --device-id device-1 --vehicle-id bus-1
```

The rebind command prints the one-time token returned by the API. Store it securely if you need it, and do not commit it.
