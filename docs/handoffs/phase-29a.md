# Phase 29A Handoff

## Phase

Phase 29A — External Predictor Adapter Evaluation

## Status

- Complete for the adapter contract documentation, candidate-only TheTransitClock feasibility review, and test-only mock adapter contract checks.
- Active phase after this handoff: Phase 29B — AVL / Vendor Adapter Pilot Implementation.

## What Was Implemented

- Documented the external predictor adapter contract in `docs/phase-29a-external-predictor-adapter-evaluation.md`.
- Added strict Trip Updates adapter output validation before protobuf serialization.
- Added test-only mock external adapter coverage for valid output passing through existing normalization and diagnostics persistence.
- Added test-only rejection coverage for missing active-feed trips, impossible stop sequences, stale prediction timestamps, wrong-agency output, wrong-feed-version output, unsupported added-trip predictions, low confidence, and missing confidence.
- Added adapter failure coverage for timeout, unavailable service, and malformed-response style errors producing visible diagnostics and empty valid Trip Updates output.
- Documented Vehicle Positions independence from external predictor availability.
- Recorded TheTransitClock public source URLs reviewed and the review date.

## What Was Designed But Intentionally Not Implemented Yet

- No TheTransitClock or external predictor runtime integration.
- No production runtime adapter wiring, environment variables, service clients, network calls, subprocess calls, Java/Maven/Tomcat invocation, or external service requirements.
- No automatic runtime external-to-deterministic fallback chain; deterministic prediction remains the default configured adapter and safe fallback path for a later approved phase.
- No real-world observed-arrival/departure comparison or production ETA-quality evaluation.

## Schema And Interface Changes

- Added internal `prediction.TripUpdate.AgencyID`, `prediction.TripUpdate.FeedVersionID`, and `prediction.TripUpdate.Confidence` for adapter-output scope and quality checks. These are internal Go fields only and do not change public GTFS-RT protobuf contracts.
- Added internal diagnostics reason `adapter_output_rejected`.
- Public feed URLs, GTFS-RT protobuf output shape, auth boundaries, and API contracts are unchanged.

## Dependency Changes

- Updated `docs/dependencies.md` to mark TheTransitClock as candidate-only.
- TheTransitClock is not vendored, not installed, not required, and not a runtime dependency.
- TheTransitClock remains GPL-3.0 licensed; vendoring or linking code requires later explicit maintainer/license review.

## Migrations Added

- None.

## Tests Added And Results

- Added focused Trip Updates adapter output validation tests under `internal/feed/tripupdates`.
- Focused `go test ./internal/prediction ./internal/feed/tripupdates ./internal/realtimequality` passed.

## Checks Run And Blocked Checks

- `make validate` — passed.
- `make realtime-quality` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed.
- Blocked checks: none so far.

## Known Issues

- Public-source TheTransitClock review is not runtime compatibility proof.
- Phase 29A mock adapter tests do not prove better ETAs, production-grade ETA quality, real-world predictor compatibility, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or vendor equivalence.
- Runtime external predictor integration still needs a later approved phase with dependency review, health/readiness behavior, failure/fallback behavior, and real compatibility evidence.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-29a.md`
  - `docs/phase-29a-external-predictor-adapter-evaluation.md`
  - `docs/phase-29b-avl-vendor-adapter-pilot.md`
  - `docs/tutorials/device-avl-integration.md`
  - `docs/evidence/redaction-policy.md`
  - `SECURITY.md`
- First files likely to edit:
  - `docs/phase-29b-avl-vendor-adapter-pilot.md`
  - `docs/handoffs/phase-29b.md`
  - telemetry adapter test-only helpers selected for Phase 29B
  - docs/status files after implementation
- Commands to run before coding:
  - `make validate`
  - `make realtime-quality`
  - `make test`
  - `git diff --check`
- Known blockers:
  - No real vendor AVL data or credentials should be used.
  - No named vendor runtime dependency should be added without explicit approval.
- Recommended first implementation slice:
  - Build a synthetic vendor payload adapter pilot behind the existing `/v1/telemetry` boundary, using public-safe fixtures and tests only.
