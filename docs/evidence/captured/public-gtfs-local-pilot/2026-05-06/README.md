# Public GTFS Local/Pilot Evidence — 2026-05-06

## Outcome

Outcome C — public-GTFS local/pilot run completed with public-safe retained
summaries.

This packet records a local Phase 33 LA Metro Bus GTFS run. It proves that this
repository, in the recorded local environment, imported a current public GTFS
ZIP, published it as the active local schedule feed, fetched the five public
paths, and retained public-safe summaries.

This packet does not prove agency adoption, agency endorsement, agency
approval, official agency feed status, agency-owned final-root proof, consumer
submission, consumer review, consumer acceptance, consumer ingestion/listing/
display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production
readiness, production multi-tenant hosting, real vendor AVL compatibility, real
LA Metro realtime data, real-world ETA accuracy, or production-grade ETA
quality.

## Source And Catalog References

Catalog facts were checked on `2026-05-06T21:15Z` and are time-sensitive.

Mobility Database reference:
`https://mobilitydatabase.org/feeds/gtfs/mdb-29`

- Feed: Los Angeles County Metropolitan Transportation Authority (LA Metro),
  Bus GTFS schedule feed.
- Official feed: yes, as shown by the catalog page at check time.
- Producer URL:
  `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip`.
- Routes: 114 bus routes.
- Service date range shown by the catalog page: December 8, 2025 through
  April 1, 2027.
- Latest dataset history row at check time: downloaded April 29, 2026, zipped
  size 21.17 MB, unzipped size 310.16 MB.
- License/terms reference:
  `http://developer.metro.net/the-basics/policies/terms-and-conditions/`.

Transitland secondary reference:
`https://www.transit.land/feeds/f-9q5-metro~losangeles`

- Current static GTFS URL at check time:
  `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip`.
- Last fetch shown by the catalog page: May 6, 2026.
- License URL:
  `http://developer.metro.net/the-basics/policies/terms-and-conditions/`.
- Attribution fields shown at check time included required attribution text:
  `provided by LA Metro`.

Do not treat catalog facts as permanent. Re-check catalogs before a future
attempt.

## Core Evidence Summary

| Field | Value |
| --- | --- |
| Local public root | `http://localhost:19080` |
| Source URL | `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip` |
| Access date (UTC) | `2026-05-06T21:25Z` |
| License/terms URL | `http://developer.metro.net/the-basics/policies/terms-and-conditions/` |
| Source ZIP SHA-256 | `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94` |
| Source ZIP local handling | downloaded to ignored `.cache/public-gtfs-local-pilot/2026-05-06/`; not committed |
| Import result | published `gtfs-import-1` for local `LACMTA` agency setup |
| Fetched schedule ZIP SHA-256 | `1819fade012ca53a58d880285bb3ab85a0fce0a1b241d20cf15320e8542503ab` |
| Fetched schedule ZIP bytes | 16,419,986 |

The local `LACMTA` agency row and local admin identity were created only to run
the repository locally against the public GTFS `agency_id`. This setup does not
imply LA Metro approval or official publication by Open Transit RT.

## Imported Dataset Summary

The import completed through `cmd/gtfs-import` with `-agency-id LACMTA` and
`-timeout 15m`.

| Entity | Count |
| --- | ---: |
| `agency` | 1 |
| `routes` | 114 |
| `stops` | 11884 |
| `trips` | 33642 |
| `stop_times` | 2105503 |
| `shapes` | 343530 |
| `calendar` | 130 |
| `calendar_dates` | 8432 |
| `frequencies` | 0 |

## Published Schedule Proof

`/public/gtfs/schedule.zip` was fetched from the local public root, then
unzipped in ignored `.cache/` storage. Public-safe summaries from the fetched
ZIP show:

- `agency.txt`: one row with `agency_id=LACMTA`, agency name
  `Metro - Los Angeles`, timezone `America/Los_Angeles`, and agency URL
  `https://www.metro.net`.
- `routes.txt`: 114 routes, all `route_type=3` bus routes.
- Sample route short names include `10/48`, `102`, `105`, `108`, `110`, `111`,
  `115`, and `117`.
- `calendar.txt`: 130 rows.
- `calendar_dates.txt`: 8,432 rows.
- Service-date coverage from `calendar.txt` and `calendar_dates.txt`:
  `20251208` through `20270401`.

The source ZIP and fetched schedule ZIP have different SHA-256 values because
Open Transit RT regenerates the public schedule ZIP from imported database
tables. Public-safe comparison summaries show the source and fetched schedule
match on agency summary, route count, route type counts, and service-date
minimum/maximum. This is the retained proof that the active published local
schedule was the imported LA Metro public GTFS, not the repo sample feed.

## Five-Path Fetch Summary

| Path | HTTP | Bytes | SHA-256 | Notes |
| --- | ---: | ---: | --- | --- |
| `/public/feeds.json` | 200 | 2,496 | `d2ecdc45245b7e336517fe485da997f9c3c5931be9d682f7b25f0e804f85f9a6` | fetched before validator rerun |
| `/public/gtfs/schedule.zip` | 200 | 16,419,986 | `1819fade012ca53a58d880285bb3ab85a0fce0a1b241d20cf15320e8542503ab` | generated local schedule ZIP |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` | valid empty protobuf publication |
| `/public/gtfsrt/trip_updates.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` | valid empty protobuf publication |
| `/public/gtfsrt/alerts.pb` | 200 | 15 | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` | valid empty protobuf publication |

The GTFS-RT fetches prove endpoint availability and protobuf publication only.
They do not prove real LA Metro vehicle telemetry, real LA Metro alerts
operations, real-world ETA accuracy, or production-grade Trip Updates.

## Validator Summary

Pinned validator tooling was present, but the static GTFS validator execution
failed because this local macOS environment could not locate a Java runtime.
That failure is recorded as a validator execution blocker, not as a claim about
LA Metro GTFS data quality.

| Feed type | Validator | Result |
| --- | --- | --- |
| `schedule` | MobilityData GTFS Validator `v7.1.0` | failed to execute; Java runtime unavailable |
| `vehicle_positions` | MobilityData GTFS Realtime validator wrapper | passed; 0 errors, 0 warnings, 0 info |
| `trip_updates` | MobilityData GTFS Realtime validator wrapper | passed; 0 errors, 0 warnings, 0 info |
| `alerts` | MobilityData GTFS Realtime validator wrapper | passed; 0 errors, 0 warnings, 0 info |

After validator attempts, `/public/feeds.json` reported schedule validation
`failed`, the three realtime feed validations `passed`, and
`canonical_validation_complete=false`.

## Telemetry Dry-Run Summary

`scripts/device-onboarding.sh simulate --dry-run` was run against the local
root with synthetic `LACMTA` identifiers. It printed three payloads and did not
submit telemetry. This is dry-run helper coverage only. It is not real LA Metro
AVL, real vehicle telemetry, real device integration, vendor compatibility, or
ETA-quality evidence.

## Admin/Private Boundary Summary

Public root checks:

- `GET http://localhost:19080/admin/operations` without token: `404`.
- `GET http://localhost:19080/admin/debug/gtfsrt/vehicle_positions.json`
  without token: `404`.

Direct local service checks:

- `GET http://localhost:19081/admin/operations` without token: `401`.
- `GET http://localhost:19081/admin/operations` with local admin token: `200`.
- `GET http://localhost:19083/admin/debug/gtfsrt/vehicle_positions.json`
  without token: `401`.
- `GET http://localhost:19084/admin/debug/gtfsrt/trip_updates.json` without
  token: `401`.
- `GET http://localhost:19085/admin/alerts` without token: `401`.

The token used for local admin checks was generated at runtime and was not
retained in this packet.

## Post-Edit Checks

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.

## Claim Boundary

This packet proves only local/pilot handling of a real public static GTFS
dataset in the recorded development environment.

It does not claim agency adoption, agency endorsement, agency approval,
official agency feed status, agency-owned final-root proof, consumer
submission, consumer review, consumer acceptance, consumer ingestion/listing/
display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production
readiness, production multi-tenant hosting, real vendor AVL compatibility,
real-world ETA accuracy, production-grade ETA quality, or real LA Metro
realtime data.

The post-Phase-32 final-root blocker remains unchanged. Consumer statuses
remain unchanged.

## Inventory

- `README.md`: Outcome C summary, source/catalog references, imported dataset
  summary, fetch proof, validator summary, telemetry dry-run summary,
  admin/private boundary summary, and claim boundary.
- `command-log-inventory-2026-05-06.md`: command inventory and public-safe
  output summaries.
- `retained-summaries-2026-05-06.md`: compact retained public-safe summaries
  for checksums, fetches, schedule proof, validators, and boundaries.
