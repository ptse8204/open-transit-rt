# Docs Assets

These assets are repo-owned documentation visuals for Phase 10.

The image-generation tool was used first to create draft bitmap concepts for the required diagrams. Because generated drafts introduced minor label and endpoint inaccuracies, the final checked-in PNGs were rendered from exact SVG specs in this directory so the repository-owned assets remain truthful to the Phase 9 codebase.

## `architecture-overview.png`

Source:
- `architecture-overview.svg`

Alt text:
- Architecture overview showing GTFS ZIP import and GTFS Studio feeding published GTFS, telemetry feeding matcher assignments, Vehicle Positions and Trip Updates publication through a prediction adapter, Service Alerts, validation/scorecard, and public feed outputs.

Generation prompt/spec:
- Clean technical architecture overview. Include GTFS import, GTFS Studio, published GTFS, telemetry ingest, matcher/assignments, Vehicle Positions, prediction adapter, Trip Updates, Service Alerts, validation/scorecard, and public outputs: `schedule.zip`, `vehicle_positions.pb`, `trip_updates.pb`, `alerts.pb`, and `feeds.json`.

## `agency-deployment.png`

Source:
- `agency-deployment.svg`

Alt text:
- Agency deployment diagram showing public internet through TLS reverse proxy, anonymous public feed paths, protected admin/debug paths, Go services, Postgres/PostGIS, pinned validators, and optional predictor adapter boundary.

Generation prompt/spec:
- Clean small-agency deployment diagram. Show TLS reverse proxy, anonymous public feeds, protected admin/debug, Go services, Postgres/PostGIS, pinned validators, and optional predictor adapter. Do not show unsupported SSO, consumer acceptance, or production-ready claims.

## `quickstart-flow.png`

Source:
- `quickstart-flow.svg`

Alt text:
- Local quickstart flow showing database bootstrap, validator install, sample GTFS import, service startup, publication bootstrap, telemetry ingest, public feed fetches, validation run, and scorecard inspection.

Generation prompt/spec:
- Clean numbered quickstart flow. Include bootstrap DB, install validators, import sample GTFS, start services, bootstrap publication, ingest telemetry, fetch public feeds, run validation, and inspect scorecard.

## `public-vs-admin-endpoints.png`

Source:
- `public-vs-admin-endpoints.svg`

Alt text:
- Public versus protected endpoint boundary showing anonymous public schedule, feeds.json, and realtime protobuf routes on one side, and protected GTFS Studio, JSON debug, telemetry events, validation, scorecard, Alerts admin, and device rebinding routes on the other.

Generation prompt/spec:
- Clean two-column endpoint boundary. Public side includes `schedule.zip`, `feeds.json`, `vehicle_positions.pb`, `trip_updates.pb`, and `alerts.pb`. Protected side includes GTFS Studio, GTFS Studio draft subroutes, JSON debug feeds, telemetry events, validation run, scorecard, Alerts admin, and device rebinding. Include the Bearer JWT/admin-cookie auth boundary.
