# Open Questions

These questions do not block the next phase.

## Deployment

- What production hosting target should be documented first: single VM, managed container platform, or Kubernetes?
- Should production HTTPS termination be owned by this repo or by deployment infrastructure?

## Authentication

- Which auth provider should be used when admin/operator login is implemented?
- Should device credentials use opaque bearer tokens only, or support signed device JWTs later?
- Should Phase 1's local debug `/v1/events` endpoint be removed, protected behind admin auth, or moved under a separate admin route before production deployment?

## Validation Tooling

- Which exact GTFS static validator distribution should be pinned first?
- Which GTFS-Realtime validator should be used for CI versus scheduled production checks?
- Phase 8 added canonical validator command adapters and `validation_report` persistence, but did not pin validator binary download/install steps. Which distributions should local dev, CI, and production standardize on?

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
