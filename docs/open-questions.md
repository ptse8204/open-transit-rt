# Open Questions

These questions do not block Phase 0.

## Deployment

- What production hosting target should be documented first: single VM, managed container platform, or Kubernetes?
- Should production HTTPS termination be owned by this repo or by deployment infrastructure?

## Authentication

- Which auth provider should be used when admin/operator login is implemented?
- Should device credentials use opaque bearer tokens only, or support signed device JWTs later?

## Validation Tooling

- Which exact GTFS static validator distribution should be pinned first?
- Which GTFS-Realtime validator should be used for CI versus scheduled production checks?

## Prediction Backends

- Should the first real Trip Updates adapter be an internal deterministic ETA engine or TheTransitClock?
- What is the minimum diagnostic output required from each predictor backend?

## GTFS Studio

- Which draft entities should the first UI expose before full table coverage?
- Should draft publish require validator success in all environments or only production-like environments?
