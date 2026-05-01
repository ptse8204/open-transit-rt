# Real Agency GTFS Onboarding

This guide helps an agency operator prepare, validate, review, and publish a real GTFS schedule through Open Transit RT.

It does not prove agency endorsement, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or agency-owned production-domain readiness. Demo GTFS remains demo-only.

## Before You Start

Use a real agency GTFS ZIP only when the agency or feed owner has approved its use and the data is public-safe to handle in this repo or deployment.

Do not commit:

- private contracts
- private contact information
- private operator notes
- private ticket links or portal screenshots
- non-public vehicle or device identifiers
- raw private telemetry
- credentials, tokens, private keys, or database URLs with passwords
- private GTFS files unless they are explicitly reviewed and approved as public-safe

Prefer synthetic fixtures for examples. If using public GTFS, record the source, license or permission, and why it is safe to keep.

## Intake Checklist

Complete this checklist before importing a real agency GTFS ZIP.

| Area | What to confirm |
| --- | --- |
| GTFS ZIP source | Record where the ZIP came from, who provided it, and whether it is public GTFS or agency-provided private data. |
| Agency permission | Confirm the agency or feed owner approved using the ZIP for this deployment or review. |
| License/contact metadata | Confirm the open license, license URL, and monitored technical contact email. |
| Agency identity | Confirm `agency.txt` uses the approved public agency name and URL. |
| Timezone | Confirm `agency_timezone` is correct for the agency-local service day. |
| Required files | Confirm `agency.txt`, `routes.txt`, `stops.txt`, `trips.txt`, `stop_times.txt`, and either usable `calendar.txt` or usable `calendar_dates.txt` exist. |
| Routes | Confirm route IDs, short names, long names, and `route_type` values represent the agency's public service. |
| Stops | Confirm stop coordinates are public stops, not private facility-only or operator-only points unless approved for publication. |
| Trips and stop times | Confirm trips reference valid routes and service IDs, and stop times reference valid stops and trips. |
| Calendar coverage | Confirm service dates cover the intended launch or pilot period and do not leave the feed with no usable service. |
| Shapes | Confirm shape IDs, point order, and optional distance fields are consistent enough for matching and review. |
| Frequencies | Confirm `frequencies.txt` rows use valid time ranges and `exact_times` values when the agency uses headway service. |
| Block IDs | Confirm `block_id` values are present when the agency expects vehicle continuity across trips. |
| Service date review | Confirm after-midnight trips and service-day boundaries make sense for the agency timezone. |
| Validation command path | Confirm who will run import validation and canonical validation, and where the output will be stored. |
| Publish approval | Confirm who is allowed to activate the imported feed. |
| Redaction/privacy review | Confirm no private data, secrets, private notes, or private operator artifacts will be committed. |
| Final public-feed review | Confirm whether the feed root is agency-owned, agency-approved, local/demo-only, or pilot-only. |

## Metadata Approval

Do not treat demo metadata as agency-approved. Record approval for each field before publishing or using the values in evidence.

| Field | Approved value | Approved by | Approval date | Notes |
| --- | --- | --- | --- | --- |
| Agency name |  |  |  |  |
| Agency URL |  |  |  |  |
| Timezone |  |  |  |  |
| Technical contact email |  |  |  |  |
| License name |  |  |  |  |
| License URL |  |  |  |  |
| Public feed root |  |  |  |  |

The public feed root must be handled carefully. If no agency-owned or agency-approved root exists, final public-feed review is limited to local/demo or pilot evidence. That does not prove agency-domain production readiness; Phase 23 remains closed as blocker-documented until final-root evidence exists.

## Import And Publish Path

Use the existing path that matches how the schedule is being prepared.

| Step | Existing path | Operator review |
| --- | --- | --- |
| Import GTFS ZIP | Use `cmd/gtfs-import` for a prepared ZIP. GTFS Studio supports typed draft authoring and publishing, not ZIP upload replacement. | Confirm source, permission, metadata, and redaction review before import. |
| Review validation | Read the import validation result and any stored validation report. Use the Operations Console and [GTFS Validation Triage](gtfs-validation-triage.md) for plain-language troubleshooting. | Resolve blocking errors before publication. Review warnings before deciding whether to proceed. |
| Publish active feed | The GTFS ZIP import flow activates a valid imported feed. GTFS Studio publishes typed drafts through the existing Studio publish flow. | Confirm the correct feed version is active and old data was not partially activated after a failure. |
| Verify public feed | Fetch `/public/feeds.json` and `/public/gtfs/schedule.zip`; inspect headers and unzip the schedule ZIP. | Confirm URLs, license/contact fields, and active schedule contents match the approved feed. |
| Review readiness | Use `/admin/operations` to review feed URLs, validation state, telemetry freshness, setup checklist, and evidence links. | Treat readiness as an operator review surface, not consumer acceptance or compliance proof. |

Example CLI import for a reviewed ZIP:

```bash
go run ./cmd/gtfs-import \
  -agency-id <agency-id> \
  -zip /path/to/reviewed-agency-gtfs.zip \
  -actor-id <operator-or-reviewer> \
  -notes "real agency GTFS import after approval review"
```

Do not commit the ZIP path, raw output, or notes if they reveal private local paths, private contacts, private operator notes, or non-public agency data.

## Publish Review Checklist

Complete this before treating a real GTFS import as ready for public feed review.

| Check | Required result |
| --- | --- |
| GTFS source approved | The source and permission are recorded. |
| Validation reviewed | Import validation and canonical validation status are reviewed. |
| No private data included | The ZIP and evidence contain no private contracts, private notes, private contacts, credentials, or non-public identifiers. |
| License/contact approved | License, license URL, and technical contact are agency-approved. |
| Service dates reviewed | Calendars and service dates cover the intended operating period. |
| Public feed URLs verified | `feeds.json` and `schedule.zip` are reachable at the reviewed root. |
| Publication environment understood | The operator knows whether this is local, pilot, or production-directed. |
| Consumer acceptance not claimed | Validation and public fetches are not described as consumer ingestion or acceptance. |
| Final-root status understood | If no agency-owned or agency-approved root exists, no agency-domain production proof is claimed. |

## Final Public-Feed Review

Review the five public feed URLs together:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

For Phase 24, the key real-agency schedule checks are `feeds.json` and `schedule.zip`. Realtime feed quality still depends on telemetry, matching, Trip Updates behavior, and Alerts operations.

Tie this review back to Phase 23:

- If the root is local, it is local demo/evaluation evidence only.
- If the root is the existing DuckDNS pilot, it remains hosted/operator pilot evidence only.
- If the root is agency-owned or agency-approved, collect the Phase 23-style domain, TLS, public fetch, validator, and approval evidence before claiming agency-domain production proof.

## Evidence Packet Template

Use `docs/evidence/real-agency-gtfs/templates/import-review-template.md` for future real-agency import packets. The template is intentionally empty until real agency-approved, public-safe evidence exists.
