# Docs Asset Generation Guidance

Use the image-generation tool when creating or updating documentation assets.

## Required Assets

Create or update these repo-owned assets:

- `docs/assets/architecture-overview.png`
- `docs/assets/agency-deployment.png`
- `docs/assets/quickstart-flow.png`
- `docs/assets/public-vs-admin-endpoints.png`

## Asset Rules

- diagrams must match the real repo architecture and endpoint surface
- prefer clean, technical, minimal diagrams over decorative art
- add alt text for each image where used in docs
- record the generation prompt/spec in `docs/assets/README.md`
- if screenshots are used, ensure they reflect the current app
- if image generation is unavailable in the environment, fall back to Mermaid and record the blocker in the handoff

## Content Expectations

### architecture-overview
Show:
- GTFS import and GTFS Studio feeding published feed versions
- telemetry ingest
- deterministic matcher / assignments
- public GTFS and GTFS-RT feeds
- alerts and compliance flows
- optional predictor adapters behind boundaries

### agency-deployment
Show:
- reverse proxy / TLS
- public protobuf endpoints
- private admin surfaces
- services
- Postgres/PostGIS
- optional validator tooling
- optional external predictor backend through adapter boundary

### quickstart-flow
Show:
- bootstrap
- migrate
- seed
- import GTFS
- ingest telemetry
- fetch feeds
- run validation
- inspect feeds metadata and scorecard

### public-vs-admin-endpoints
Show a clear split between:
- public anonymous endpoints
- admin authenticated endpoints
- debug/admin-only JSON routes
