# Docs Assets

These are repo-owned documentation visuals. They are used to teach the current Open Transit RT architecture, local workflows, and public/admin boundaries without making unsupported product or compliance claims.

## Visual Review Rule

Every generated or generated-assisted image must be manually reviewed before it is referenced in docs.

- Labels must match current repo behavior.
- Captions must say whether the image is illustrative or exact-behavior based.
- Alt text must describe the useful content of the image, not just repeat the filename.
- Generated drafts may be refined into checked-in SVG/PNG assets when text or labels need cleanup.
- Do not check in visuals that imply hosted SaaS, SSO, CAD/AVL replacement, consumer acceptance, marketplace/vendor equivalence, full CAL-ITP/Caltrans compliance, or universal production readiness.

## Current Assets

### `agency-journey-to-public-feeds.png`

Source:
- `agency-journey-to-public-feeds.svg`
- Generated image draft created with the image-generation tool, then manually reviewed and refined into an exact SVG-derived PNG because generated text can be unreliable.

Type:
- Illustrative teaching graphic, not a screenshot.

Used in:
- `README.md`
- `docs/tutorials/agency-demo-flow.md`

Alt text:
- Illustrative agency journey from GTFS import or GTFS Studio drafts through schedule publication, authenticated telemetry, validation, and public GTFS plus GTFS Realtime feeds.

Prompt/spec:
- Create an illustrative poster-style agency journey diagram for Open Transit RT. Show a small transit agency moving from GTFS import or GTFS Studio draft, to published static GTFS, to authenticated vehicle telemetry, to validation and public feeds. Include only these labels: "Agency journey", "Import or author GTFS", "Publish schedule.zip", "Ingest vehicle telemetry", "Validate feeds", and "Public GTFS + GTFS Realtime". Do not show a fake application UI or unsupported claims.

Truthfulness notes:
- The visual says the path is illustrative and does not claim hosted SaaS, CAD/AVL replacement, consumer acceptance, or full compliance.

### `docs-choose-your-path.png`

Source:
- `docs-choose-your-path.svg`
- Generated image draft created with the image-generation tool, then manually reviewed and refined into an exact SVG-derived PNG.

Type:
- Illustrative docs navigation graphic.

Used in:
- `docs/README.md`

Alt text:
- Illustrative documentation guide showing paths for trying locally, running the agency demo, planning deployment, reviewing evidence, and contributing.

Prompt/spec:
- Create a choose-your-path graphic that helps readers decide where to start in the docs. Use a central "Start here" signpost with paths for "Try locally", "Run agency demo", "Plan deployment", "Review evidence", and "Contribute". Do not imply SaaS hosting, consumer acceptance, full compliance, or agency endorsement.

Truthfulness notes:
- The visual points to documentation paths only. It is not a product UI and does not imply agency endorsement.

### `data-flow-through-system.png`

Source:
- `data-flow-through-system.svg`
- Generated image draft created with the image-generation tool, then manually reviewed and refined into an exact SVG-derived PNG.

Type:
- Illustrative system explainer.

Used in:
- `docs/README.md`

Alt text:
- Illustrative data-flow diagram showing GTFS import, GTFS Studio drafts, vehicle telemetry, Open Transit RT state, validation, and public feed outputs.

Prompt/spec:
- Create a simple system data-flow explainer showing inputs, Open Transit RT services and state, validation, and public feed outputs. Include GTFS import, GTFS Studio drafts, vehicle telemetry, published GTFS, assignments, prediction adapter, Alerts, validation, `schedule.zip`, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts. Do not imply SSO, SaaS hosting, CAD/AVL replacement, full compliance, consumer acceptance, or external predictor integration.

Truthfulness notes:
- The visual summarizes current boundaries. Trip Updates remain behind the prediction adapter; validation checks generated artifacts rather than proving consumer acceptance.

### `architecture-overview.png`

Source:
- `architecture-overview.svg`

Type:
- Exact-behavior architecture diagram rendered from a reviewed SVG spec.

Used in:
- Available as a deeper architecture reference; not currently used on the README front door.

Alt text:
- Architecture overview showing GTFS ZIP import and GTFS Studio feeding published GTFS, telemetry feeding matcher assignments, Vehicle Positions and Trip Updates publication through a prediction adapter, Service Alerts, validation/scorecard, and public feed outputs.

Prompt/spec:
- Clean technical architecture overview. Include GTFS import, GTFS Studio, published GTFS, telemetry ingest, matcher/assignments, Vehicle Positions, prediction adapter, Trip Updates, Service Alerts, validation/scorecard, and public outputs: `schedule.zip`, `vehicle_positions.pb`, `trip_updates.pb`, `alerts.pb`, and `feeds.json`.

### `agency-deployment.png`

Source:
- `agency-deployment.svg`

Type:
- Exact-behavior deployment-boundary diagram rendered from a reviewed SVG spec.

Used in:
- Deployment/readiness documentation as needed.

Alt text:
- Agency deployment diagram showing public internet through TLS reverse proxy, anonymous public feed paths, protected admin/debug paths, Go services, Postgres/PostGIS, pinned validators, and optional predictor adapter boundary.

Prompt/spec:
- Clean small-agency deployment diagram. Show TLS reverse proxy, anonymous public feeds, protected admin/debug, Go services, Postgres/PostGIS, pinned validators, and optional predictor adapter. Do not show unsupported SSO, consumer acceptance, or production-ready claims.

### `quickstart-flow.png`

Source:
- `quickstart-flow.svg`

Type:
- Exact-behavior local workflow diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/README.md`
- `docs/tutorials/local-quickstart.md`

Alt text:
- Exact-behavior local quickstart flow showing database bootstrap, validator install, sample GTFS import, service startup, publication bootstrap, telemetry ingest, public feed fetches, validation run, and scorecard inspection.

Prompt/spec:
- Clean numbered quickstart flow. Include bootstrap DB, install validators, import sample GTFS, start services, bootstrap publication, ingest telemetry, fetch public feeds, run validation, and inspect scorecard.

### `public-vs-admin-endpoints.png`

Source:
- `public-vs-admin-endpoints.svg`

Type:
- Exact-behavior endpoint-boundary diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/agency-demo-flow.md`

Alt text:
- Exact-behavior boundary diagram showing anonymous public schedule, feeds.json, and realtime protobuf routes separated from protected GTFS Studio, JSON debug, telemetry events, validation, scorecard, Alerts admin, and device rebinding routes.

Prompt/spec:
- Clean two-column endpoint boundary. Public side includes `schedule.zip`, `feeds.json`, `vehicle_positions.pb`, `trip_updates.pb`, and `alerts.pb`. Protected side includes GTFS Studio, GTFS Studio draft subroutes, JSON debug feeds, telemetry events, validation run, scorecard, Alerts admin, and device rebinding. Include the Bearer JWT/admin-cookie auth boundary.
