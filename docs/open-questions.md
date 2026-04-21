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

- Should the first real Trip Updates adapter be an internal deterministic ETA engine or TheTransitClock?
- What is the minimum diagnostic output required from each predictor backend?

## GTFS Studio

- Phase 5 resolved the first UI entity scope: minimal operational forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Should a later Studio UI add map editing for shape points and timetable-design interactions for stop_times, or keep those as row editors?
- Should draft publish require canonical MobilityData validator success in all environments once canonical validation is wired, or only production-like environments?
