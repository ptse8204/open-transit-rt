# Retained Public-Safe Summaries — 2026-05-06

## Checksums

| Artifact | SHA-256 |
| --- | --- |
| Source LA Metro Bus GTFS ZIP | `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94` |
| Fetched local `/public/gtfs/schedule.zip` | `1819fade012ca53a58d880285bb3ab85a0fce0a1b241d20cf15320e8542503ab` |
| Fetched `/public/feeds.json` before validator rerun | `d2ecdc45245b7e336517fe485da997f9c3c5931be9d682f7b25f0e804f85f9a6` |
| Fetched empty GTFS-RT protobuf feeds | `89989fa94ed97951cd59ec0440a37ddd89521bbf6b38b20a999743da377c9d86` |

## Active Published Schedule Proof

The fetched local schedule ZIP was unzipped in ignored `.cache/` storage and
summarized.

| Field | Summary |
| --- | --- |
| Agency row | `agency_id=LACMTA`, `agency_name=Metro - Los Angeles`, `agency_timezone=America/Los_Angeles`, `agency_url=https://www.metro.net` |
| Route count | 114 |
| Route type counts | `route_type=3`: 114 |
| Sample route short names | `10/48`, `102`, `105`, `108`, `110`, `111`, `115`, `117` |
| `calendar.txt` rows | 130 |
| `calendar_dates.txt` rows | 8,432 |
| Service-date coverage | `20251208` through `20270401` |

Source and fetched schedule summaries matched on agency summary, route count,
route type counts, and service-date minimum/maximum. This verifies the active
local schedule as the imported public GTFS rather than the repo sample feed.

## Fetch Summary

| Path | HTTP | Bytes | Interpretation |
| --- | ---: | ---: | --- |
| `/public/feeds.json` | 200 | 2,496 | local feed discovery for `LACMTA` |
| `/public/gtfs/schedule.zip` | 200 | 16,419,986 | generated static GTFS from imported feed tables |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | 15 | empty valid protobuf publication |
| `/public/gtfsrt/trip_updates.pb` | 200 | 15 | empty valid protobuf publication |
| `/public/gtfsrt/alerts.pb` | 200 | 15 | empty valid protobuf publication |

## Validator Summary

| Feed type | Result |
| --- | --- |
| `schedule` | Static validator attempted; Java runtime unavailable, so execution failed. |
| `vehicle_positions` | GTFS-RT validator passed with 0 errors, 0 warnings, 0 info. |
| `trip_updates` | GTFS-RT validator passed with 0 errors, 0 warnings, 0 info. |
| `alerts` | GTFS-RT validator passed with 0 errors, 0 warnings, 0 info. |

## Dry-Run And Boundary Summary

- Telemetry simulator: `scripts/device-onboarding.sh simulate --dry-run`
  printed three synthetic payloads and submitted nothing.
- Public root admin paths checked returned `404`.
- Direct local admin/debug paths without token returned `401`.
- Direct local admin operations with a runtime-generated local admin token
  returned `200`.

## Claim Boundary

These summaries support Outcome C for local/pilot public GTFS dataset handling
only. They do not support agency adoption, official feed status, final-root
proof, consumer evidence, Caltrans/CAL-ITP compliance, hosted SaaS,
production-readiness, real vendor AVL compatibility, real LA Metro realtime
data, real-world ETA accuracy, or production-grade ETA quality.
