# Validator Records

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run

## Evidence Category

- Repo-proven capability: validation endpoints execute allowlisted validator IDs and persist normalized records.
- Deployment/operator proof: partial local proof only. These are not production validator records.
- Third-party proof: none.

## Static GTFS Validation

- Validator ID/version: `static-mobilitydata`, MobilityData GTFS Validator `v7.1.0`
- Input artifact/version: generated `schedule.zip` for feed version `gtfs-import-7`
- Run timestamp (UTC): 2026-04-22T21:33:56.906207Z
- Result summary: failed, 1 error, 0 warnings, 0 info
- Full output location: `artifacts/validation/validate-schedule.json`
- Failure reason: Java runtime missing on this workstation.

## GTFS-RT Validation

### Vehicle Positions

- Validator ID/version: `realtime-mobilitydata`, `ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89`
- Feed URL/artifact: `artifacts/public/vehicle_positions.pb`
- Feed generated/header timestamp or revision: feed version `gtfs-import-7`; public fetch at 2026-04-22T21:36:10Z
- Run timestamp (UTC): 2026-04-22T21:33:57.674307Z
- Result summary: failed, 1 error, 0 warnings, 0 info
- Full output location: `artifacts/validation/validate-vehicle-positions.json`
- Failure reason: pinned wrapper invocation reached the Docker image, but the image rejected `--schedule` as an unrecognized option.

### Trip Updates

- Validator ID/version: `realtime-mobilitydata`, `ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89`
- Feed URL/artifact: `artifacts/public/trip_updates.pb`
- Feed generated/header timestamp or revision: feed version `gtfs-import-7`; public fetch at 2026-04-22T21:36:10Z
- Run timestamp (UTC): 2026-04-22T21:36:39.508987Z
- Result summary: failed, 1 error, 0 warnings, 0 info
- Full output location: `artifacts/validation/validate-trip-updates.json`
- Failure reason: pinned wrapper invocation reached the Docker image, but the image rejected `--schedule` as an unrecognized option.

### Alerts

- Validator ID/version: `realtime-mobilitydata`, `ghcr.io/mobilitydata/gtfs-realtime-validator@sha256:5d2a3c14fba49983e1968c4a715e8ca624d4062bf4afede74aeca26322436c89`
- Feed URL/artifact: `artifacts/public/alerts.pb`
- Feed generated/header timestamp or revision: feed version `gtfs-import-7`; public fetch at 2026-04-22T21:36:10Z
- Run timestamp (UTC): 2026-04-22T21:36:40.336033Z
- Result summary: failed, 1 error, 0 warnings, 0 info
- Full output location: `artifacts/validation/validate-alerts.json`
- Failure reason: pinned wrapper invocation reached the Docker image, but the image rejected `--schedule` as an unrecognized option.

## Retention and Completeness

- Warnings/errors preserved without omission? Yes. Raw JSON outputs are committed in `artifacts/validation/`.
- Retention location: `docs/evidence/captured/local-demo/2026-04-22/artifacts/validation/`.
- Redactions applied: none.

## Blockers

- Clean validator evidence is missing.
- Static validation requires Java or an equivalent pinned runtime in this environment.
- Realtime validation requires correcting or replacing the pinned GTFS-RT validator wrapper invocation before clean validator evidence can be collected.
