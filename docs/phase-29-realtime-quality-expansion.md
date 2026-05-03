# Phase 29 — Real-World Realtime Quality Expansion

## Status

Complete for the synthetic replay evidence expansion scope.

Phase 29 expands committed replay coverage only. It does not add real-world observed-arrival/departure evidence, does not measure real-world ETA accuracy, and does not claim production-grade ETA quality.

## Purpose

Expand beyond Phase 19 baseline replay fixtures toward richer realtime quality evidence while preserving conservative matcher and predictor behavior.

Phase 19 made quality measurable. Phase 29 adds harder synthetic scenarios and tighter replay assertions so unknown, withheld, ambiguous, stale, and degraded outcomes remain visible.

## Implemented Scope

1. Added richer deterministic replay fixtures for after-midnight service, exact and non-exact frequency trips, block continuity, long layovers, sparse telemetry, noisy/off-shape GPS, stale/ambiguous patterns, cancellation alert linkage, and manual override expiry.
2. Added replay fixture modeling for `frequencies` and optional manual override `expires_at`.
3. Updated the replay telemetry repository to return the latest telemetry row per vehicle for feed snapshots, matching the production repository contract.
4. Strengthened replay metric comparison for already-recorded cancellation alert linkage and unsupported disruption-withheld counts.
5. Added focused realtime-quality tests for the Phase 29 scenarios.

## Evidence Limits

- All new Phase 29 fixtures use synthetic, public-safe identifiers and telemetry.
- No real private telemetry, private agency GTFS, private device IDs, private logs, or operator artifacts were added.
- Real-world observed-arrival/departure comparison remains unavailable because no retained real observed stop-time evidence exists in the repo.
- Real route/time-period quality metrics are explicitly deferred. The new fixtures cover synthetic routes and times only; they are not deployment coverage metrics and do not imply real-world route/time-period performance.
- TheTransitClock and other external predictors remain deferred behind the existing `internal/prediction.Adapter` boundary.

## Metrics

Phase 29 keeps the existing denominator-aware prediction diagnostics as the public quality shape:

- unknown assignment rate;
- ambiguous assignment rate;
- stale telemetry rate;
- Trip Updates coverage rate;
- future-stop coverage rate;
- withheld-by-reason counts;
- degraded-by-reason counts;
- manual override assignment counts;
- cancellation alert linkage counts;
- added-trip, short-turn, and detour withheld counts.

Zero-denominator rates continue to use `not_applicable` with no percent value.

## Deterministic Behavior

No production matcher, predictor, public feed URL, GTFS-RT protobuf contract, consumer status, auth boundary, or external dependency was changed.

Frequency support added in Phase 29 is replay-fixture modeling support for behavior already present in the matcher, Vehicle Positions, and deterministic predictor. Manual override expiry support added in Phase 29 is replay-harness modeling support for expiry semantics already reflected by the production state model and repository contract.

## Acceptance Criteria

Phase 29 is complete because:

- richer replay fixtures exist and pass;
- fixture documentation explains what each scenario proves and does not prove;
- metrics remain repeatable and denominator-aware;
- unknown/withheld/degraded cases remain visible;
- deterministic behavior was not broadened beyond test-backed existing behavior;
- observed-arrival comparison is documented as unavailable;
- real route/time-period quality metrics are deferred honestly;
- no production-grade ETA or real-world accuracy claim was introduced.

## Required Checks

```bash
go test ./internal/realtimequality
make validate
make realtime-quality
make test
make smoke
make test-integration
git diff --check
docker compose -f deploy/docker-compose.yml config
```

## Explicit Non-Goals

Phase 29 does not:

- claim production-grade ETA quality;
- claim real-world ETA accuracy;
- add observed-arrival evidence;
- add real route/time-period deployment coverage metrics;
- force external predictors into core;
- hide unknown/withheld/degraded cases;
- change consumer statuses;
- claim consumer acceptance, agency endorsement, hosted SaaS availability, marketplace equivalence, or CAL-ITP/Caltrans compliance.
