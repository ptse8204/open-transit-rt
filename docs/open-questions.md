# Open Questions

These questions do not block the next phase.

## Deployment

- What production hosting target should be documented first: single VM, managed container platform, or Kubernetes?
- Should production HTTPS termination be owned by this repo or by deployment infrastructure?
- Phase 10 documents the current pilot path as Postgres/PostGIS through Compose plus deployment-owned Go service process management and a TLS reverse proxy. A fully packaged app-container or Kubernetes path remains open.

## Authentication

- Which auth provider should be used when admin/operator login is implemented?
- Should device credentials use opaque bearer tokens only, or support signed device JWTs later?
- Should Phase 1's local debug `/v1/events` endpoint be removed, protected behind admin auth, or moved under a separate admin route before production deployment?

## Validation Tooling

- Answered for the repo-supported path in Phase 9: MobilityData GTFS Validator `v7.1.0` and Docker-backed MobilityData GTFS-RT validator image digest in `tools/validators/validators.lock.json`.
- Open for deployments: whether a production environment should use the repo-supported Docker-backed GTFS-RT wrapper or document an equivalent checksum/digest contract for a native executable.

## Static GTFS Publication

- Phase 8 serves `/public/gtfs/schedule.zip` on demand from active published GTFS tables with deterministic bytes. Should a later phase add materialized ZIP caching or checksum comparison against original uploaded ZIP bytes?

## Prediction Backends

- Phase 7 chose an internal deterministic predictor as the first real Trip Updates adapter. Should a later phase add TheTransitClock as an alternate adapter, and what deployment profile should own it?
- What quality threshold should be required before claiming production-grade ETA quality rather than conservative schedule-deviation predictions?
- What historical telemetry retention and backtesting workflow should support MAE by route, stop, and time of day?

## Alerts

- Phase 8 chose both operator-authored and system-derived Alerts for canceled-trip reconciliation. Should later alert workflows include richer affected-route/stop/time selectors, multilingual text, and full operator UI review before publication?

## GTFS Studio

- Phase 5 resolved the first UI entity scope: minimal operational forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Should a later Studio UI add map editing for shape points and timetable-design interactions for stop_times, or keep those as row editors?
- Should draft publish require canonical MobilityData validator success in all environments once canonical validation is wired, or only production-like environments?
