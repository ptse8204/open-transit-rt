# Command And Log Inventory — 2026-05-06

This is a public-safe summary of commands and relevant outputs from the Phase
33 LA Metro public-GTFS local Outcome C run. Raw ZIPs, fetched protobufs,
fetched schedule ZIPs, generated service logs, local admin token summaries, and
intermediate inspection files were kept under ignored `.cache/` paths.

## Source Download

```bash
mkdir -p .cache/public-gtfs-local-pilot/2026-05-06
curl --fail --location --silent --show-error \
  --output .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip \
  https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip
shasum -a 256 .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip
```

Result:

```text
ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94  .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip
```

## Isolated Local Database

```bash
make db-up
docker compose -f deploy/docker-compose.yml exec -T postgres \
  psql -U postgres -d postgres \
  -c "DROP DATABASE IF EXISTS open_transit_rt_phase33_20260506" \
  -c "CREATE DATABASE open_transit_rt_phase33_20260506"
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt_phase33_20260506?sslmode=disable" \
  go run ./cmd/migrate up
```

Result summary:

- Postgres became ready.
- A dedicated local Phase 33 database was created.
- Migrations applied through version 8.

## Local LACMTA Setup

```sql
INSERT INTO agency (id, name, timezone, contact_email, public_url)
VALUES (
  'LACMTA',
  'Los Angeles County Metropolitan Transportation Authority',
  'America/Los_Angeles',
  'dev@example.com',
  'https://www.metro.net'
);
```

Result summary:

- Local `LACMTA` agency row created or updated.
- Local `phase33-admin@example.com` admin/editor/operator/read-only role
  bindings created for authenticated local admin calls.

This was local setup only and does not imply agency approval.

## Import

```bash
DATABASE_URL="postgres://postgres:postgres@localhost:55432/open_transit_rt_phase33_20260506?sslmode=disable" \
  go run ./cmd/gtfs-import \
    -agency-id LACMTA \
    -zip .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip \
    -actor-id phase-33-outcome-c-attempt \
    -notes "Phase 33 public GTFS local/pilot Outcome C retry" \
    -timeout 15m
```

Public-safe result summary:

```text
status=published
feed_version_id=gtfs-import-1
errors=0
warnings=0
report_stored=true
routes=114
stops=11884
trips=33642
stop_times=2105503
shapes=343530
calendar=130
calendar_dates=8432
frequencies=0
```

## Local Services And Public Root

Repo service binaries were started with `AGENCY_ID=LACMTA`,
`DATABASE_URL=postgres://postgres:postgres@localhost:55432/open_transit_rt_phase33_20260506?sslmode=disable`,
and `PUBLIC_BASE_URL=http://localhost:19080`.

Ports:

- `19081`: `cmd/agency-config`
- `19082`: `cmd/telemetry-ingest`
- `19083`: `cmd/feed-vehicle-positions`
- `19084`: `cmd/feed-trip-updates`
- `19085`: `cmd/feed-alerts`
- `19086`: `cmd/gtfs-studio`
- `19080`: local public proxy for the five public paths

Publication metadata bootstrap was performed through the direct local admin
service port with a runtime-generated local admin token. The token was not
retained.

Bootstrap result:

```text
{"stored":true}
readyz={"service":"agency-config","status":"ready"}
```

## Five-Path Fetch

```bash
curl --fail --location --dump-header <ignored-header-path> \
  --output <ignored-output-path> http://localhost:19080/public/feeds.json
curl --fail --location --dump-header <ignored-header-path> \
  --output <ignored-output-path> http://localhost:19080/public/gtfs/schedule.zip
curl --fail --location --dump-header <ignored-header-path> \
  --output <ignored-output-path> http://localhost:19080/public/gtfsrt/vehicle_positions.pb
curl --fail --location --dump-header <ignored-header-path> \
  --output <ignored-output-path> http://localhost:19080/public/gtfsrt/trip_updates.pb
curl --fail --location --dump-header <ignored-header-path> \
  --output <ignored-output-path> http://localhost:19080/public/gtfsrt/alerts.pb
```

Result summary:

| Path | HTTP | Bytes | SHA-256 |
| --- | ---: | ---: | --- |
| `/public/feeds.json` | 200 | 2,496 | `d2ecdc45245b7e336517fe485da997f9c3c5931be9d682f7b25f0e804f85f9a6` |
| `/public/gtfs/schedule.zip` | 200 | 16,419,986 | `1819fade012ca53a58d880285bb3ab85a0fce0a1b241d20cf15320e8542503ab` |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` |
| `/public/gtfsrt/trip_updates.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` |
| `/public/gtfsrt/alerts.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` |

## Published Schedule Inspection

```bash
unzip <ignored-fetched-schedule.zip> -d .cache/public-gtfs-local-pilot/2026-05-06/schedule-inspect
```

Public-safe summary:

- `agency.txt`: `agency_id=LACMTA`, `agency_name=Metro - Los Angeles`,
  `agency_timezone=America/Los_Angeles`, `agency_url=https://www.metro.net`.
- `routes.txt`: 114 routes, all `route_type=3`.
- `calendar.txt`: 130 rows.
- `calendar_dates.txt`: 8,432 rows.
- Service-date coverage: `20251208` through `20270401`.
- Source and fetched schedule summaries matched on agency summary, route count,
  route type counts, and service-date minimum/maximum.

## Validator Attempts

```bash
./scripts/check-validators.sh
curl -X POST http://localhost:19181/admin/validation/run \
  -H "Authorization: Bearer <local-admin-token>" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"static-mobilitydata","feed_type":"schedule"}'
curl -X POST http://localhost:19181/admin/validation/run \
  -H "Authorization: Bearer <local-admin-token>" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"vehicle_positions"}'
curl -X POST http://localhost:19181/admin/validation/run \
  -H "Authorization: Bearer <local-admin-token>" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"trip_updates"}'
curl -X POST http://localhost:19181/admin/validation/run \
  -H "Authorization: Bearer <local-admin-token>" \
  -H "Content-Type: application/json" \
  --data '{"validator_id":"realtime-mobilitydata","feed_type":"alerts"}'
```

Result summary:

- Validator tooling check: passed.
- Static GTFS validator: execution failed because Java runtime was unavailable
  in this local environment.
- Vehicle Positions GTFS-RT validator: passed, 0 errors, 0 warnings, 0 info.
- Trip Updates GTFS-RT validator: passed, 0 errors, 0 warnings, 0 info.
- Alerts GTFS-RT validator: passed, 0 errors, 0 warnings, 0 info.

## Telemetry Dry-Run

```bash
TARGET=http://localhost:19080 \
AGENCY_ID=LACMTA \
DEVICE_ID=phase33-dryrun-device \
VEHICLE_ID=phase33-dryrun-vehicle \
  scripts/device-onboarding.sh simulate --dry-run
```

Result summary:

- Three synthetic payloads were printed.
- No telemetry was sent.
- This is dry-run helper coverage only, not real LA Metro realtime data.

## Admin/Private Boundary Check

Result summary:

```text
public root /admin/operations without token: 404
public root /admin/debug/gtfsrt/vehicle_positions.json without token: 404
direct agency-config /admin/operations without token: 401
direct agency-config /admin/operations with token: 200
direct vehicle positions debug without token: 401
direct trip updates debug without token: 401
direct alerts admin without token: 401
```

## Teardown

```bash
make db-down
```

Result summary:

- Local services were stopped.
- Docker Compose Postgres was stopped.
- No raw GTFS ZIP, fetched schedule ZIP, protobuf artifacts, service logs, or
  tokens were committed.

## Required Checks

```bash
make validate
make test
git diff --check
```

Result summary:

- `make validate` passed.
- `make test` passed.
- `git diff --check` passed.
