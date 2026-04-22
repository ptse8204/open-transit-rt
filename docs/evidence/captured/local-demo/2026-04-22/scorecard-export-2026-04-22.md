# Scorecard Export Evidence

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run

## Export Details

- Format: JSON.
- Export trigger: manual API call.
- Export timestamp (UTC): 2026-04-22T21:37:22Z.
- Command/endpoint used: `POST http://localhost:8081/admin/compliance/scorecard` with Bearer admin token, followed by `GET http://localhost:8081/admin/compliance/scorecard`.
- Artifact location:
  - `artifacts/scorecard/scorecard-2026-04-22T2137Z.json`
  - `artifacts/scorecard/latest-scorecard-2026-04-22T2137Z.json`
- Artifact checksum/hash:
  - `scorecard-2026-04-22T2137Z.json`: `6c7cc7c11d099cde3e64140b2ecb28da12190d1c6b7b0773a2ccd3df368ab4a1`
  - `latest-scorecard-2026-04-22T2137Z.json`: `62b17d248d4df55c7b65479159686a5ff0b977618fe20a5b975fe46c827ec1c2`

## Timestamped History

- Prior export timestamps/artifacts:
  - 2026-04-22T21:33:57Z from the demo flow, stored transiently in the demo temp directory before the final export.
  - 2026-04-22T21:37:22Z committed in this packet.
- Evidence that generation is repeatable: `make demo-agency-flow` generated a scorecard, and a later manual POST generated a new stored snapshot after additional validator records were captured.

## Scorecard Result

The final local scorecard reported:

- `overall_status`: `red`
- `schedule_status`: `yellow`
- `vehicle_positions_status`: `yellow`
- `trip_updates_status`: `green`
- `alerts_status`: `green`
- `validation_status`: `red`
- `discoverability_status`: `green`
- `consumer_ingestion_status`: `red`

Important detail: this is local dev scorecard evidence, not proof of compliance. `validation_status` is red because canonical validator runs failed. `consumer_ingestion_status` is red because workflow records are `not_started`.

## Retention

- Retention policy: missing production policy.
- Storage boundary: committed redacted local evidence packet.
- Redaction notes: no secrets or tokens were stored in the committed scorecard artifacts.

## Scheduled Job Evidence

- Job/scheduler reference: Missing.
- Last successful run: Missing.
- Recent run history excerpt or link: Missing.

## Manual Fallback Evidence

- Reason manual export was used: no scheduler/job exists in this repo or local environment.
- Operator notes: manual export is sufficient to prove local scorecard extraction but not scheduled production history.
