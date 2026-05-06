# Public GTFS Local/Pilot Evidence Packet Template

Status: template only. Do not fill this with fake evidence.

## Packet Metadata

| Field | Value |
| --- | --- |
| Outcome | Outcome A / Outcome B / Outcome C |
| Environment | local / pilot |
| Capture date (UTC) |  |
| Operator or reviewer |  |
| Public/local feed root |  |
| Raw GTFS retained in git? | No by default |

## Source Metadata

| Field | Value |
| --- | --- |
| Source agency/feed |  |
| Source URL |  |
| Access date (UTC) |  |
| License or terms reference |  |
| Catalog facts checked at (UTC) |  |
| Mobility Database reference |  |
| Transitland reference |  |
| Source ZIP SHA-256 |  |
| Raw ZIP local path handling | `.cache/` or other ignored path |

## Catalog Facts Summary

Record public catalog facts with a checked-at timestamp. Treat them as
time-sensitive, not permanent.

| Catalog | Facts recorded | Checked at (UTC) |
| --- | --- | --- |
| Mobility Database | Official-feed status, route count, service range, producer URL, dataset size |  |
| Transitland | Current URL, last fetch, license URL, attribution fields |  |

## Import Summary

| Field | Value |
| --- | --- |
| Import command, redacted |  |
| Import timestamp (UTC) |  |
| Import result | passed / failed / blocked |
| Active feed version after import |  |
| Import warnings/errors |  |
| Notes |  |

## Published Schedule Proof

| Field | Value |
| --- | --- |
| Fetched schedule URL |  |
| Fetched schedule SHA-256 |  |
| Fetched schedule unzip path | ignored temp path |
| `agency.txt` public-safe summary |  |
| `routes.txt` public-safe summary |  |
| Service-date coverage summary |  |
| Verified as imported public GTFS, not repo sample feed? | Yes / No / Blocked |

## Public Feed Fetch Proof

| Path | Result | Status/size/checksum summary | Notes |
| --- | --- | --- | --- |
| `/public/feeds.json` |  |  |  |
| `/public/gtfs/schedule.zip` |  |  |  |
| `/public/gtfsrt/vehicle_positions.pb` |  |  |  |
| `/public/gtfsrt/trip_updates.pb` |  |  |  |
| `/public/gtfsrt/alerts.pb` |  |  |  |

## Validator Summaries

| Feed | Result | Summary or blocker | Artifact/reference |
| --- | --- | --- | --- |
| Static GTFS schedule |  |  |  |
| Vehicle Positions |  |  |  |
| Trip Updates |  |  |  |
| Alerts |  |  |  |

Do not claim blocked validator checks passed.

## Known Warnings Or Errors

| Area | Warning/error/blocker | Disposition |
| --- | --- | --- |
| Import |  |  |
| Schedule proof |  |  |
| Public fetch |  |  |
| Validators |  |  |
| Realtime |  |  |

## Telemetry Simulator Or Dry-Run Summary

| Field | Value |
| --- | --- |
| Simulator/dry-run path attempted? | Yes / No |
| Result or blocker |  |
| Realtime labeling | empty / synthetic / withheld / not applicable |

GTFS-RT fetches prove endpoint availability and protobuf publication only. They
do not prove real agency realtime telemetry, real alerts operation, ETA quality,
or production-grade Trip Updates.

## Admin/Private Route Boundary Check

| Route/check | Result | Notes |
| --- | --- | --- |
| Anonymous admin/private route check |  |  |
| Anonymous debug route check |  |  |
| Authenticated admin route check, if applicable |  |  |

## Claim-Boundary Statement

This packet supports only the outcome recorded above.

It does not claim agency adoption, agency endorsement, agency approval,
official agency feed status, agency-owned final-root proof, consumer
submission, consumer review, consumer acceptance, consumer ingestion/listing/
display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production
readiness, production multi-tenant hosting, real vendor AVL compatibility,
real-world ETA accuracy, or production-grade ETA quality.

## Checksums Or Inventory Notes

| File/artifact | SHA-256 or inventory note |
| --- | --- |
| Source GTFS ZIP |  |
| Fetched schedule ZIP |  |
| Fetch summaries |  |
| Validator summaries |  |
| Command/log inventory |  |
