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
- Phase 4 defined and stores an internal GTFS import validation report contract but did not wire canonical validator install/download steps; which static validator distribution should be pinned for Phase 8 compliance?

## Static GTFS Publication

- Should later public schedule-feed serving store original uploaded ZIP bytes, regenerate ZIPs from published tables, or support both with checksum comparison?

## Prediction Backends

- Phase 7 chose an internal deterministic predictor as the first real Trip Updates adapter. Should a later phase add TheTransitClock as an alternate adapter, and what deployment profile should own it?
- What quality threshold should be required before claiming production-grade ETA quality rather than conservative schedule-deviation predictions?
- What historical telemetry retention and backtesting workflow should support MAE by route, stop, and time of day?

## Alerts

- Phase 7 persists cancellation-to-alert linkage signals with `expected_alert_missing=true`, but public Alerts authoring/persistence remains deferred. Should Alerts be operator-authored, incident-derived, or both?
- What minimal Alerts authoring workflow should clear or satisfy a canceled-trip `expected_alert_missing` review signal?

## GTFS Studio

- Phase 5 resolved the first UI entity scope: minimal operational forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Should a later Studio UI add map editing for shape points and timetable-design interactions for stop_times, or keep those as row editors?
- Should draft publish require canonical MobilityData validator success in all environments once canonical validation is wired, or only production-like environments?
