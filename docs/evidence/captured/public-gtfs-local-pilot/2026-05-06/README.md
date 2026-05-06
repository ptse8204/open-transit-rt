# Public GTFS Local/Pilot Attempt — 2026-05-06

## Outcome

Outcome B — attempted public-GTFS run blocked.

This packet records a blocked Phase 33 attempt to import LA Metro Bus GTFS into
the local Open Transit RT app flow. It is a blocker summary and command/log
inventory only. It does not include public feed fetch proof, validator proof, or
evidence that the active published local schedule is LA Metro GTFS.

## Blocker Summary

The LA Metro Bus GTFS ZIP downloaded successfully into ignored local storage,
and its SHA-256 checksum was recorded. The first import attempt used the default
local app agency, `demo-agency`, and failed because LA Metro `agency.txt` uses
`agency_id=LACMTA`.

A second attempt seeded a local `LACMTA` agency row and imported with
`-agency-id LACMTA`. That attempt reached the repo-supported importer and parsed
the dataset counts, but the publish step blocked on the importer context
deadline while inserting the large `stop_times.txt` dataset:

```text
gtfs import publish failed and failure report could not be stored: publish error: insert stop time 70060001831656-DEC25/53: timeout: context deadline exceeded; report error: context deadline exceeded
```

Because import/publish did not complete, Outcome C was not reached. No
five-path local public fetch proof, validator proof, telemetry simulator proof,
or admin/private boundary proof is claimed in this packet.

## Source And Catalog References

Catalog facts were checked on `2026-05-06T20:35:48Z` and are time-sensitive.

Mobility Database reference:
`https://mobilitydatabase.org/feeds/gtfs/mdb-29`

- Feed: Los Angeles County Metropolitan Transportation Authority (LA Metro),
  Bus GTFS schedule feed.
- Official feed: yes, as shown by the catalog page at check time.
- Producer URL:
  `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip`
- Routes: 114 bus routes.
- Service date range shown by the catalog page: December 8, 2025 through
  April 1, 2027.
- Latest dataset history row at check time: downloaded April 29, 2026, zipped
  size 21.17 MB, unzipped size 310.16 MB.
- License/terms reference:
  `http://developer.metro.net/the-basics/policies/terms-and-conditions/`

Transitland secondary reference:
`https://www.transit.land/feeds/f-9q5-metro~losangeles`

- Current static GTFS URL at check time:
  `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip`
- Last fetch shown by the catalog page: May 6, 2026.
- License URL:
  `http://developer.metro.net/the-basics/policies/terms-and-conditions/`
- Attribution fields shown at check time included required attribution text:
  `provided by LA Metro`.

Do not treat catalog facts as permanent. Re-check both catalogs before a future
attempt.

## Source Download

| Field | Value |
| --- | --- |
| Source URL | `https://gitlab.com/LACMTA/gtfs_bus/raw/master/gtfs_bus.zip` |
| Access date (UTC) | `2026-05-06T20:37:00Z` |
| License/terms URL | `http://developer.metro.net/the-basics/policies/terms-and-conditions/` |
| Source ZIP SHA-256 | `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94` |
| Source ZIP local handling | downloaded to ignored `.cache/public-gtfs-local-pilot/2026-05-06/`; not committed |

## Attempted Import Summary

| Attempt | Result |
| --- | --- |
| `demo-agency` import | Failed validation because `agency.txt` did not contain requested `agency_id="demo-agency"`. |
| `LACMTA` import | Blocked during publish with `context deadline exceeded` while inserting `stop_times.txt`; import result reported `report_stored=false`. |

Parsed counts from the blocked `LACMTA` import attempt:

| GTFS file/entity | Count |
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

## Checks Not Claimed

These checks were not completed because the core import/publish step blocked:

- published `/public/gtfs/schedule.zip` fetch for LA Metro GTFS;
- proof that fetched schedule ZIP is the imported public GTFS rather than the
  repo sample feed;
- five-path public fetch;
- validators;
- telemetry simulator or dry-run proof;
- admin/private/debug boundary check for the LA Metro local run.

## Claim Boundary

This packet proves only that the Phase 33 LA Metro public-GTFS local run was
attempted and blocked at import/publish time.

It does not claim agency adoption, agency endorsement, agency approval,
official agency feed status, agency-owned final-root proof, consumer
submission, consumer review, consumer acceptance, consumer ingestion/listing/
display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production
readiness, production multi-tenant hosting, real vendor AVL compatibility,
real-world ETA accuracy, or production-grade ETA quality.

The post-Phase-32 final-root blocker remains unchanged. Consumer statuses
remain unchanged.

## Inventory

- `README.md`: blocker summary, source references, attempted import summary,
  checks not claimed, and claim boundary.
- `command-log-inventory-2026-05-06.md`: command inventory and public-safe
  output summaries.
