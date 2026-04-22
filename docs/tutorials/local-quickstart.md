# Local Quickstart

This tutorial starts the current Open Transit RT development environment on one machine. It is local-dev only and uses seeded demo credentials.

## Prerequisites

- Go matching `go.mod`
- Docker with Compose support
- `curl`, `zip`, and `unzip`
- Java if you want the static GTFS validator JAR to execute successfully

The GTFS-RT validator workflow uses Docker because the repo-supported path is a pinned wrapper around a MobilityData container image.

## Bootstrap

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
curl -fsS -X POST http://localhost:8082/v1/telemetry \
  -H "Authorization: Bearer dev-device-token" \
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
