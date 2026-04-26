# Docs Assets

These are repo-owned documentation visuals. They are used to teach the current Open Transit RT architecture, local workflows, and public/admin boundaries without making unsupported product or compliance claims.

Public-facing copies of selected PNG assets also live under `wiki/assets/` so `/wiki` pages do not depend on internal relative paths.

## Visual Review Rule

Every generated or generated-assisted image must be manually reviewed before it is referenced in docs.

- Labels must match project behavior.
- Captions must say whether the image is illustrative or exact-behavior based.
- Alt text must describe the useful content of the image, not just repeat the filename.
- Generated drafts may be refined into checked-in SVG/PNG assets when text or labels need cleanup.
- Do not check in visuals that imply hosted SaaS, SSO, CAD/AVL replacement, consumer acceptance, marketplace/vendor equivalence, full CAL-ITP/Caltrans compliance, or universal production readiness.

## Current Assets

### `agency-journey-to-public-feeds.png`

Source:
- `agency-journey-to-public-feeds.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG because generated text and dense labels can be unreliable.

Type:
- Illustrative teaching graphic, not a screenshot.

Used in:
- `README.md`
- `docs/tutorials/agency-demo-flow.md`
- `wiki/agency-demo.md`

Alt text:
- Illustrative agency journey from GTFS import or GTFS Studio drafts through schedule publication, authenticated telemetry, validation, and public GTFS plus GTFS Realtime feeds.

Prompt/spec:
- Create a clean, modern, easy-to-understand agency path graphic. Show the simple path from preparing GTFS, publishing a schedule, adding vehicle data, validating, and publishing public feeds. Keep labels large and readable. Do not show fake UI or unsupported claims.

Truthfulness notes:
- The visual says the path is illustrative and does not claim hosted SaaS, CAD/AVL replacement, consumer acceptance, or full compliance.

### `docs-choose-your-path.png`

Source:
- `docs-choose-your-path.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG.

Type:
- Illustrative docs navigation graphic.

Used in:
- `docs/README.md`
- `wiki/README.md`

Alt text:
- Illustrative documentation guide showing paths for trying locally, running the agency demo, planning deployment, reviewing evidence, and contributing.

Prompt/spec:
- Create a clean docs navigation graphic with five large destination cards: Try locally, Agency demo, Deploy, Evidence, and Contribute. Keep labels large and readable. Do not imply SaaS hosting, consumer acceptance, full compliance, or agency endorsement.

Truthfulness notes:
- The visual points to documentation paths only. It is not a product UI and does not imply agency endorsement.

### `data-flow-through-system.png`

Source:
- `data-flow-through-system.svg`
- Regenerated image draft created with the image-generation tool, then manually reviewed and refined into a simpler SVG-derived PNG.

Type:
- Illustrative system explainer.

Used in:
- `docs/README.md`
- `wiki/how-it-works.md`

Alt text:
- Illustrative data-flow diagram showing GTFS import, GTFS Studio drafts, vehicle telemetry, Open Transit RT state, validation, and public feed outputs.

Prompt/spec:
- Create a clean four-column system explainer showing Inputs, Open Transit RT, Validation, and Public feeds. Keep labels large and readable. Do not imply SSO, SaaS hosting, CAD/AVL replacement, full compliance, consumer acceptance, or external predictor integration.

Truthfulness notes:
- The visual summarizes current boundaries. Trip Updates remain behind the prediction adapter; validation checks generated artifacts rather than proving consumer acceptance.

### `architecture-overview.png`

Source:
- `architecture-overview.svg`

Type:
- Exact-behavior architecture summary rendered from a reviewed SVG spec.

Used in:
- Available as a deeper architecture reference; not currently used on the README front door.

Alt text:
- Architecture overview showing schedule inputs, vehicle state, feed builders, public feed outputs, and validation/readiness evidence as four clear areas.

Prompt/spec:
- Clean architecture summary with four bounded areas: schedule inputs, vehicle state, feed builders, and public feeds. Use large readable labels, no crossing lines, and no phase-history wording.

### `agency-deployment.png`

Source:
- `agency-deployment.svg`

Type:
- Exact-behavior deployment-shape diagram rendered from a reviewed SVG spec.

Used in:
- Deployment/readiness documentation as needed.

Alt text:
- Small agency deployment diagram showing public internet, TLS reverse proxy, anonymous public feeds, protected admin tools, Open Transit RT services, pinned validators, prediction adapter, and Postgres/PostGIS.

Prompt/spec:
- Clean small-agency deployment diagram. Show TLS reverse proxy, anonymous public feeds, protected admin tools, Open Transit RT services, Postgres/PostGIS, pinned validators, and optional predictor adapter. Keep text bounded and avoid overlapping connectors. Do not show unsupported SSO, consumer acceptance, or production-ready claims.

### `quickstart-flow.png`

Source:
- `quickstart-flow.svg`

Type:
- Exact-behavior local workflow diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/README.md`
- `docs/tutorials/local-quickstart.md`

Alt text:
- Exact-behavior local quickstart flow grouped into setup, publish-and-ingest, and review-output steps.

Prompt/spec:
- Clean numbered quickstart flow grouped into setup, publish-and-ingest, and review-output steps. Use large bounded labels and avoid long connector lines.

### `public-vs-admin-endpoints.png`

Source:
- `public-vs-admin-endpoints.svg`

Type:
- Exact-behavior endpoint-boundary diagram rendered from a reviewed SVG spec.

Used in:
- `docs/tutorials/agency-demo-flow.md`

Alt text:
- Exact-behavior endpoint boundary diagram showing anonymous public feed routes separated from protected admin, debug, telemetry, validation, and scorecard examples.

Prompt/spec:
- Clean two-column endpoint boundary. Public side includes `schedule.zip`, `feeds.json`, `vehicle_positions.pb`, `trip_updates.pb`, and `alerts.pb`. Protected side includes examples such as GTFS Studio, JSON debug feeds, telemetry events, validation, and scorecard. Keep labels bounded and readable.
