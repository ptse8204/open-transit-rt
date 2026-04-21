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
- Should Phase 4 wire validator install/download steps, or only define the adapter/report contract first?

## Prediction Backends

- Should the first real Trip Updates adapter be an internal deterministic ETA engine or TheTransitClock?
- What is the minimum diagnostic output required from each predictor backend?

## GTFS Studio

- Which draft entities should the first UI expose before full table coverage?
- Should draft publish require validator success in all environments or only production-like environments?
