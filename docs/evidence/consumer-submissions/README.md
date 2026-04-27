# Consumer Submission Evidence Tracker

This directory contains Phase 13 and Phase 20 consumer and aggregator submission evidence records.

Tracker last reviewed timestamp: `2026-04-27T04:42:17Z`

Reviewed by: Codex Phase 20 docs/evidence pass

Linked Phase 12 evidence packet: `docs/evidence/captured/oci-pilot/2026-04-24/`

Machine-readable tracker snapshot: `docs/evidence/consumer-submissions/status.json`

## Claim Boundary

The OCI pilot hosted/operator evidence packet supports public URL, validation, TLS/auth-boundary, monitoring, backup/restore, and scorecard job-history claims for the OCI pilot. It does not prove consumer submission or third-party acceptance.

Validator success and public fetch proof are supporting evidence only. They are not consumer acceptance, consumer ingestion, CAL-ITP compliance, marketplace listing, or vendor equivalence.

## Current Records

The table below must agree exactly with `docs/evidence/consumer-submissions/status.json` for target name, status, packet path, prepared timestamp, and evidence reference values.

Evidence references for every prepared target:

- `oci_packet`: `docs/evidence/captured/oci-pilot/2026-04-24/`
- `feeds_json_snapshot`: `docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/public_feeds.json`
- `validator_records`: `docs/evidence/captured/oci-pilot/2026-04-24/validator-record-2026-04-24.md`
- `phase_19_replay_quality_summary`: `docs/handoffs/phase-19.md`

| Target name | Status | Current record | Packet path | Prepared timestamp | Next action |
| --- | --- | --- | --- | --- | --- |
| Google Maps | `prepared` | `current/google-maps.md` | `docs/evidence/consumer-submissions/packets/google-maps/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| Apple Maps | `prepared` | `current/apple-maps.md` | `docs/evidence/consumer-submissions/packets/apple-maps/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| Transit App | `prepared` | `current/transit-app.md` | `docs/evidence/consumer-submissions/packets/transit-app/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| Bing Maps | `prepared` | `current/bing-maps.md` | `docs/evidence/consumer-submissions/packets/bing-maps/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| Moovit | `prepared` | `current/moovit.md` | `docs/evidence/consumer-submissions/packets/moovit/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| Mobility Database | `prepared` | `current/mobility-database.md` | `docs/evidence/consumer-submissions/packets/mobility-database/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission/registration path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |
| transit.land | `prepared` | `current/transit-land.md` | `docs/evidence/consumer-submissions/packets/transit-land/README.md` | `2026-04-27T04:42:17Z` | Operator reviews packet, verifies official submission/registration path outside the repo, submits only if authorized, then stores redacted receipt evidence before changing status. |

## Prepared Packets

Packet completeness is tracked in `docs/evidence/consumer-submissions/packets/README.md`. A target may use `prepared` only when the packet includes all five public feed URLs, license/contact metadata, Phase 12 evidence link, validator evidence link, redaction note, next action, and allowed wording.

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
