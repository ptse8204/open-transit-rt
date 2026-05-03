# Track B Integration Phase Insert — Phase 29A And Phase 29B

## Purpose

This insert adds the external integration work that belongs after Phase 29 and before consumer submission/adoption execution.

Phase 29 expanded synthetic realtime replay coverage and intentionally deferred external predictor integration. Track B should now explicitly add two integration phases so future work does not rely on chat history:

- Phase 29A — External Predictor Adapter Evaluation
- Phase 29B — AVL / Vendor Adapter Pilot Implementation

## Placement

Recommended Track B order after Phase 29:

1. Phase 29 — Real-World Realtime Quality Expansion
2. Phase 29A — External Predictor Adapter Evaluation
3. Phase 29B — AVL / Vendor Adapter Pilot Implementation
4. Phase 30 — Consumer Submission Execution
5. Phase 31 — Agency Pilot Program Package
6. Phase 32 — Public Launch And Ecosystem Outreach

This preserves the current Phase 30/31/32 numbering while inserting the missing integration phases as 29A and 29B.

## Why These Phases Belong Here

Phase 25 documented the device/AVL integration boundary. Phase 29 expanded replay fixtures and quality measurement. Together, they create the foundation needed to evaluate external prediction and AVL/vendor integrations safely.

External integrations should be added only after the replay baseline exists because integrations need comparison points, fallback tests, and clear evidence boundaries.

## Phase 29A Summary

Phase 29A evaluates external prediction adapters such as TheTransitClock behind the existing prediction adapter boundary.

It should:

- review `internal/prediction.Adapter`;
- define external predictor contract expectations;
- evaluate TheTransitClock as a candidate;
- add mock adapter/contract tests where useful;
- compare adapter behavior against Phase 29 replay fixtures;
- define fallback behavior;
- keep deterministic prediction as default;
- avoid required runtime dependencies unless explicitly approved;
- avoid production-grade ETA claims.

## Phase 29B Summary

Phase 29B implements or documents a minimal AVL/vendor adapter pilot pattern.

It should:

- define a generic vendor payload to Open Transit RT telemetry contract;
- use synthetic vendor payload fixtures;
- keep vendor credentials outside the repo;
- preserve `POST /v1/telemetry` as the boundary;
- add dry-run/transform behavior if safe;
- test synthetic payload transformations;
- avoid certified vendor or hardware reliability claims.

## Truthfulness Boundary

Do not claim:

- production-grade ETA quality;
- real-world ETA accuracy;
- certified vendor support;
- hardware reliability;
- consumer acceptance;
- CAL-ITP/Caltrans compliance;
- hosted SaaS availability;
- agency endorsement;
- marketplace/vendor equivalence.

## Files To Add

- `docs/phase-29a-external-predictor-adapter-evaluation.md`
- `docs/phase-29b-avl-vendor-adapter-pilot.md`

## Files To Update

- `docs/track-b-productization-roadmap.md`
- `docs/roadmap-status.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-29.md` only if its next-phase recommendation still skips the integration phases
- `docs/phase-30-consumer-submission-execution.md` only if it needs a note that Phase 29A/29B now come first

## Acceptance Criteria For This Insert

This insert is complete when:

- Phase 29A and Phase 29B docs exist;
- roadmap docs show Phase 29A and Phase 29B after Phase 29;
- latest handoff recommends Phase 29A next;
- no runtime code changes are made;
- no dependency status changes are made;
- no public feed URLs, consumer statuses, or evidence claims change.
