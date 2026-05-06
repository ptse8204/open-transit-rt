# Command And Log Inventory — 2026-05-06

This is a public-safe summary of commands and relevant outputs from the blocked
Phase 33 LA Metro public-GTFS local attempt. Raw ZIP and intermediate inspection
files were kept under ignored `.cache/` paths.

## Source Download

```bash
mkdir -p .cache/public-gtfs-local-pilot/2026-05-06
curl -L --fail --show-error \
  --output .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip \
  https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip
shasum -a 256 .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip
```

Result:

```text
ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94  .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip
```

## Local App Startup

```bash
make agency-app-up
```

Result summary:

- Local app started behind `http://localhost:8080`.
- The script imported the repo sample GTFS as the active `demo-agency` feed
  before the LA Metro attempt.
- Pinned validator tooling was reported ready.

## First Import Attempt

```bash
docker cp \
  .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip \
  deploy-agency-config-1:/tmp/la-metro-gtfs-bus.zip

docker compose -f deploy/docker-compose.yml --profile app exec -T agency-config \
  /app/bin/gtfs-import \
  -agency-id demo-agency \
  -zip /tmp/la-metro-gtfs-bus.zip \
  -actor-id phase-33-public-gtfs \
  -notes "Phase 33 public GTFS local pilot evidence: LA Metro Bus GTFS"
```

Result summary:

```text
gtfs import validation failed with 1 error(s)
```

Stored report summary:

```text
agency.txt must contain requested agency_id "demo-agency"
```

This was not treated as final because the public GTFS uses `agency_id=LACMTA`.

## Source File Spot Check

```bash
unzip -q \
  .cache/public-gtfs-local-pilot/2026-05-06/la-metro-gtfs-bus.zip \
  agency.txt routes.txt calendar.txt calendar_dates.txt \
  -d .cache/public-gtfs-local-pilot/2026-05-06/source-inspect
sed -n '1,8p' .cache/public-gtfs-local-pilot/2026-05-06/source-inspect/agency.txt
sed -n '1,5p' .cache/public-gtfs-local-pilot/2026-05-06/source-inspect/routes.txt
wc -l \
  .cache/public-gtfs-local-pilot/2026-05-06/source-inspect/routes.txt \
  .cache/public-gtfs-local-pilot/2026-05-06/source-inspect/calendar.txt \
  .cache/public-gtfs-local-pilot/2026-05-06/source-inspect/calendar_dates.txt
```

Public-safe summary:

- `agency.txt` contains one agency row with `agency_id=LACMTA`, agency name
  `Metro - Los Angeles`, agency URL `https://www.metro.net`, and timezone
  `America/Los_Angeles`.
- `routes.txt` contains 113 data rows plus header, matching the 114-line
  command output.
- `calendar.txt` contains 130 data rows plus header.
- `calendar_dates.txt` contains 8431 data rows plus header.

This was a source spot check only. It is not published schedule proof because
the import/publish step later blocked.

## Local LACMTA Agency Setup

```sql
INSERT INTO agency (id, name, timezone, contact_email, public_url)
VALUES ('LACMTA', 'Metro - Los Angeles', 'America/Los_Angeles', 'bakerro@metro.net', 'https://www.metro.net')
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    timezone = EXCLUDED.timezone,
    contact_email = EXCLUDED.contact_email,
    public_url = EXCLUDED.public_url;
```

Result summary:

```text
INSERT 0 1
INSERT 0 4
```

This was local setup only and does not imply agency approval.

## Second Import Attempt

```bash
docker compose -f deploy/docker-compose.yml --profile app exec -T agency-config \
  /app/bin/gtfs-import \
  -agency-id LACMTA \
  -zip /tmp/la-metro-gtfs-bus.zip \
  -actor-id phase-33-public-gtfs \
  -notes "Phase 33 public GTFS local pilot evidence: LA Metro Bus GTFS"
```

Result summary:

```text
gtfs import publish failed and failure report could not be stored: publish error: insert stop time 70060001831656-DEC25/53: timeout: context deadline exceeded; report error: context deadline exceeded
```

Public-safe import result summary:

```text
{"import_id":25,"agency_id":"LACMTA","status":"failed","error_count":1,"warning_count":0,"info_count":0,"counts":{"agency":1,"calendar":130,"calendar_dates":8432,"frequencies":0,"routes":114,"shapes":343530,"stop_times":2105503,"stops":11884,"trips":33642},"report_stored":false,"failure_message":"publish failed and failure report could not be stored"}
```

The corresponding `gtfs_import` row remained `started` because the context
expired before the failure update completed. This is part of the blocker and
should be reviewed before a future large-dataset retry.

## Not Run After Blocker

The following checks were intentionally not claimed after the import blocker:

- local public schedule fetch for imported LA Metro GTFS;
- active schedule ZIP proof;
- five-path fetch;
- validators;
- telemetry simulator/dry-run;
- admin/private/debug boundary check for the LA Metro run.

## Teardown

```bash
make agency-app-down
```

Result summary:

- Local app containers and network were stopped and removed.
- The local Postgres volume was not reset.
