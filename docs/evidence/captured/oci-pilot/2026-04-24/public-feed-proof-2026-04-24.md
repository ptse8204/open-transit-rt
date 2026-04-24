# Hosted Public Feed Root Proof

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials
- Canonical HTTPS host: `https://open-transit-pilot.duckdns.org`

## Anonymous Hosted Fetches

| Path | Status | Artifact | Header artifact |
| --- | ---: | --- | --- |
| `/public/gtfs/schedule.zip` | 200 | `artifacts/public/public_gtfs_schedule.zip` | `artifacts/public/public_gtfs_schedule.zip.headers.txt` |
| `/public/feeds.json` | 200 | `artifacts/public/public_feeds.json` | `artifacts/public/public_feeds.json.headers.txt` |
| `/public/gtfsrt/vehicle_positions.pb` | 200 | `artifacts/public/public_gtfsrt_vehicle_positions.pb` | `artifacts/public/public_gtfsrt_vehicle_positions.pb.headers.txt` |
| `/public/gtfsrt/trip_updates.pb` | 200 | `artifacts/public/public_gtfsrt_trip_updates.pb` | `artifacts/public/public_gtfsrt_trip_updates.pb.headers.txt` |
| `/public/gtfsrt/alerts.pb` | 200 | `artifacts/public/public_gtfsrt_alerts.pb` | `artifacts/public/public_gtfsrt_alerts.pb.headers.txt` |

## Public-Edge Auth Boundary

Artifact: `artifacts/public/auth-boundary/public-and-tunneled-auth-boundary.txt`

- Public schedule, `feeds.json`, and `.pb` routes returned 200 anonymously.
- Public `.json` debug routes returned 404 from the Caddy public edge.
- Public admin mutation routes returned 404 from the Caddy public edge.
- SSH-tunneled admin mutation returned 401 without Bearer auth and 200 with Bearer auth.

## Publish / Restore URL Stability

Artifacts:

- Before update: `artifacts/public/snapshots/feeds-before-update.json`
- After controlled update: `artifacts/public/snapshots/feeds-after-update.json`
- During transient restore drill update: `artifacts/public/snapshots/feeds-during-restore-drill-transient-update.json`
- After deployment data-restore rollback: `artifacts/public/snapshots/feeds-after-restore-rollback.json`
- Final after validator rerun: `artifacts/public/snapshots/feeds-final-after-validator-rerun.json`
- Final current-live recheck: `artifacts/public/snapshots/feeds-final-current-live-20260424T163846Z.json`
- Final current-live summary: `artifacts/public/snapshots/feeds-final-current-live-20260424T163846Z-summary.json`
- Current live feed-version history: `artifacts/public/snapshots/current-live-feed-version-history-20260424T163846Z.txt`

The controlled update moved the active feed from `gtfs-import-1` to `gtfs-import-3`. A transient drill update moved it to `gtfs-import-4`; restoring the clean backup returned the active feed to `gtfs-import-3`. The public base URL and all feed URLs stayed unchanged through the update and deployment data-restore rollback drill.

The final current-live recheck at `2026-04-24T16:38:46Z` showed active `gtfs-import-3` for schedule, Vehicle Positions, Trip Updates, and Alerts. The database history artifact did not show an intentional move to `gtfs-import-16`.
