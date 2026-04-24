# Hosted Validator Records

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials

## Final Current-Live Validator Runs

| Feed type | Validator ID | Artifact | Status |
| --- | --- | --- | --- |
| `schedule` | `static-mobilitydata` | `artifacts/validation/validate-schedule.json` | `passed` |
| `vehicle_positions` | `realtime-mobilitydata` | `artifacts/validation/validate-vehicle_positions.json` | `passed` |
| `trip_updates` | `realtime-mobilitydata` | `artifacts/validation/validate-trip_updates.json` | `passed` |
| `alerts` | `realtime-mobilitydata` | `artifacts/validation/validate-alerts.json` | `passed` |

The final validator artifacts were refreshed during the current-live recheck against active feed version `gtfs-import-3`. The final public `feeds.json` snapshot reports `canonical_validation_complete=true` with schedule, Vehicle Positions, Trip Updates, and Alerts all `passed`.
