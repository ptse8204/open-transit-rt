# Public GTFS Local/Pilot Runbook

This guide explains how to repeat a Phase 33-style public GTFS local/pilot run.

It is for development and evaluation only. It does not imply agency approval, official feed status, consumer acceptance, Caltrans/CAL-ITP compliance, production readiness, real realtime data, or ETA-quality proof.

For a command-level helper that performs the reusable local/reference import
and five-path verification flow, see
[Reusable Agency Onboarding](reusable-agency-onboarding.md). Use this runbook
when you need retained public-safe summaries or a fuller manual repeatability
record.

## What this run can prove

A completed local/pilot run can prove:

- a public GTFS ZIP was downloaded from a documented public source;
- source URL, access date, license/terms reference, and checksum were recorded;
- Open Transit RT imported the public GTFS for a local-only agency setup;
- the active local schedule feed was fetched;
- the fetched schedule was verified as the imported public dataset, not the repo sample feed;
- all five public paths were fetched;
- validators ran or blockers were documented;
- telemetry dry-run was exercised if supported;
- admin/private boundaries were checked;
- public-safe summaries were retained.

## What this run cannot prove

It cannot prove:

- agency adoption;
- agency endorsement;
- agency approval;
- official agency feed status;
- agency-owned final-root proof;
- consumer submission/review/acceptance;
- consumer ingestion/listing/display;
- Caltrans/CAL-ITP compliance;
- hosted SaaS availability;
- production readiness;
- real device/vendor AVL integration;
- real agency realtime data;
- real-world ETA accuracy;
- production-grade Trip Updates.

## Inputs to collect before running

For the selected public GTFS feed, record:

```text
Feed name:
Producer/source URL:
Catalog/reference URL:
Checked-at timestamp:
License/terms URL:
Agency ID expected in agency.txt:
Why this public GTFS is safe to use locally:
```

Download raw GTFS only to ignored storage such as `.cache/`.

Do not commit the raw ZIP unless maintainers have reviewed the license/terms and intentionally accepted it.

## Suggested directory layout

```text
${PILOT_DIR}/
  source.zip
  source.sha256
  fetched-schedule.zip
  fetched-schedule.sha256
  unpacked-fetched-schedule/

${EVIDENCE_DIR}/
  README.md
  command-log-inventory-${RUN_DATE}.md
  retained-summaries-${RUN_DATE}.md
```

## High-level procedure

The examples below use shell variables so paths and URLs stay consistent:

```bash
RUN_DATE="$(date -u +%F)"
PILOT_DIR=".cache/public-gtfs-local-pilot/${RUN_DATE}"
EVIDENCE_DIR="docs/evidence/captured/public-gtfs-local-pilot/${RUN_DATE}"
GTFS_URL="https://example.com/path/to/gtfs.zip"
AGENCY_ID="example-agency"
ACTOR_ID="local-public-gtfs-pilot"
ROOT="http://localhost:19080"
```

Set `GTFS_URL` and `AGENCY_ID` to the selected public GTFS dataset before
running commands. For a retained evidence packet, use a stable `RUN_DATE` value
that matches the packet directory.

### 1. Download the public GTFS ZIP

```bash
mkdir -p "${PILOT_DIR}"
curl -L -o "${PILOT_DIR}/source.zip" "${GTFS_URL}"
sha256sum "${PILOT_DIR}/source.zip"
```

Record the checksum in the evidence packet.

### 2. Prepare a local-only agency setup

The local agency ID must match the GTFS `agency_id` that will be imported.

For example, the Phase 33 LA Metro run used `LACMTA` locally. That local setup was for evidence only and did not imply LA Metro approval.

Do not reuse production credentials, tokens, or private operator data.

### 3. Import with a suitable timeout

Use the repo-supported import command.

```bash
go run ./cmd/gtfs-import \
  -agency-id "${AGENCY_ID}" \
  -zip "${PILOT_DIR}/source.zip" \
  -actor-id "${ACTOR_ID}" \
  -notes "public GTFS local/pilot evaluation only" \
  -timeout 15m
```

Use `GTFS_IMPORT_TIMEOUT` or `-timeout` only as documented by the repo. Do not hide timeouts or failed imports.

### 4. Start or route the local public root

Use a local root dedicated to the public-GTFS pilot run, such as:

```text
http://localhost:19080
```

Avoid confusing this with the default `make agency-app-up` demo flow if that flow imports the repo sample feed for `demo-agency`.

### 5. Fetch the five public paths

```bash
curl -fsS -D "${PILOT_DIR}/feeds.headers.txt" \
  "${ROOT}/public/feeds.json" \
  -o "${PILOT_DIR}/feeds.json"

curl -fsS -D "${PILOT_DIR}/schedule.headers.txt" \
  "${ROOT}/public/gtfs/schedule.zip" \
  -o "${PILOT_DIR}/fetched-schedule.zip"

curl -fsS -D "${PILOT_DIR}/vehicle_positions.headers.txt" \
  "${ROOT}/public/gtfsrt/vehicle_positions.pb" \
  -o "${PILOT_DIR}/vehicle_positions.pb"

curl -fsS -D "${PILOT_DIR}/trip_updates.headers.txt" \
  "${ROOT}/public/gtfsrt/trip_updates.pb" \
  -o "${PILOT_DIR}/trip_updates.pb"

curl -fsS -D "${PILOT_DIR}/alerts.headers.txt" \
  "${ROOT}/public/gtfsrt/alerts.pb" \
  -o "${PILOT_DIR}/alerts.pb"
```

Record HTTP status, byte counts, and SHA-256 checksums.

### 6. Prove the fetched schedule is the imported public GTFS

```bash
mkdir -p "${PILOT_DIR}/unpacked-fetched-schedule"
unzip -q "${PILOT_DIR}/fetched-schedule.zip" \
  -d "${PILOT_DIR}/unpacked-fetched-schedule"
```

Summarize public-safe facts from:

```text
agency.txt
routes.txt
calendar.txt
calendar_dates.txt
```

At minimum record:

- agency ID;
- agency name;
- timezone;
- route count;
- route type counts;
- service date minimum and maximum.

This check prevents false positives where the server is still publishing the sample feed.

### 7. Run validators or record blockers

Static GTFS validation requires a runnable static validator environment. If Java or the pinned validator is unavailable, record the blocker clearly.

GTFS-RT validation should be run for:

```text
vehicle_positions.pb
trip_updates.pb
alerts.pb
```

If feeds are empty valid protobuf publications, say that. Do not imply real agency realtime data.

### 8. Run telemetry dry-run if supported

Use dry-run mode only unless authorized real telemetry is available.

```bash
TARGET="${ROOT}" \
AGENCY_ID="${AGENCY_ID}" \
DEVICE_ID="public-gtfs-dryrun-device" \
VEHICLE_ID="public-gtfs-dryrun-vehicle" \
  scripts/device-onboarding.sh simulate --dry-run
```

Record that no telemetry was sent if dry-run was used.

### 9. Check admin/private boundaries

Check that public root admin/debug paths are not exposed and that direct local admin/debug paths require authentication.

Do not retain runtime admin tokens in the evidence packet.

### 10. Write public-safe retained summaries

The evidence packet should include:

- source/catalog summary;
- command inventory;
- checksums;
- import summary;
- fetched schedule proof;
- five-path fetch summary;
- validator summary or blockers;
- telemetry dry-run summary;
- admin/private boundary summary;
- what the local review does not prove.

## Suggested Wording For What This Does Not Prove

Use wording like:

```text
This packet proves only local/pilot handling of a real public static GTFS dataset in the recorded development environment.
```

Do not overclaim.
