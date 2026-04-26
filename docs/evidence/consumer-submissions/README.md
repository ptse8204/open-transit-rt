# Consumer Submission Evidence Tracker

This directory contains Phase 13 consumer and aggregator submission evidence records.

Tracker last reviewed timestamp: `2026-04-26T02:10:56Z`

Reviewed by: Codex documentation pass

Linked Phase 12 evidence packet: `docs/evidence/captured/oci-pilot/2026-04-24/`

## Claim Boundary

The OCI pilot hosted/operator evidence packet supports public URL, validation, TLS/auth-boundary, monitoring, backup/restore, and scorecard job-history claims for the OCI pilot. It does not prove consumer submission or third-party acceptance.

Validator success and public fetch proof are supporting evidence only. They are not consumer acceptance, consumer ingestion, CAL-ITP compliance, marketplace listing, or vendor equivalence.

## Current Records

| Target | Current status | Current record | Next action |
| --- | --- | --- | --- |
| Google Maps | `not_started` | `current/google-maps.md` | Operator prepares and submits the packet through the official Google transit partner workflow, then stores redacted receipt evidence. |
| Apple Maps | `not_started` | `current/apple-maps.md` | Operator prepares and submits the packet through the official Apple Maps transit data workflow, then stores redacted receipt evidence. |
| Transit App | `not_started` | `current/transit-app.md` | Operator prepares and submits the packet through the official Transit data partner workflow, then stores redacted receipt evidence. |
| Bing Maps | `not_started` | `current/bing-maps.md` | Operator prepares and submits the packet through the official Microsoft/Bing Maps transit data workflow, then stores redacted receipt evidence. |
| Moovit | `not_started` | `current/moovit.md` | Operator prepares and submits the packet through the official Moovit data partner workflow, then stores redacted receipt evidence. |
| Mobility Database | `not_started` | `current/mobility-database.md` | Operator prepares and submits or registers the feed through the Mobility Database workflow, then stores redacted receipt evidence. |
| transit.land | `not_started` | `current/transit-land.md` | Operator prepares and submits or registers the feed through the transit.land workflow, then stores redacted receipt evidence. |

## Templates

Reusable target templates live in `templates/`:

- `templates/google-maps.md`
- `templates/apple-maps.md`
- `templates/transit-app.md`
- `templates/bing-maps.md`
- `templates/moovit.md`
- `templates/mobility-database.md`
- `templates/transit-land.md`

## Status Definitions

| Status | Meaning |
| --- | --- |
| `not_started` | No submission has been made. |
| `prepared` | Packet prepared only, no submission. |
| `submitted` | Submission sent, no acceptance implied. |
| `under_review` | Consumer review in progress, no acceptance implied. |
| `accepted` | Acceptance may be claimed only for the named consumer, feed scope, URL root, and evidence date. |
| `rejected` | Rejection documented, no acceptance claim. |
| `blocked` | Submission blocked by named missing evidence/action. |

## Current OCI Pilot Feed URLs For Submission Packets

Feed root:

- `https://open-transit-pilot.duckdns.org`

Exact feed URLs available from the Phase 12 hosted/operator evidence packet:

- `https://open-transit-pilot.duckdns.org/public/feeds.json`
- `https://open-transit-pilot.duckdns.org/public/gtfs/schedule.zip`
- `https://open-transit-pilot.duckdns.org/public/gtfsrt/vehicle_positions.pb`
- `https://open-transit-pilot.duckdns.org/public/gtfsrt/trip_updates.pb`
- `https://open-transit-pilot.duckdns.org/public/gtfsrt/alerts.pb`

These URLs are not listed as submitted to any target until a target record's status changes to `submitted` or later with retained evidence.
