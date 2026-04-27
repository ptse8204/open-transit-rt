# Phase Handoff Template

## Phase

Phase 19 — Realtime Quality And ETA Improvement

## Status

- Complete for the approved measurement-first Phase 19 scope.
- Active phase after this handoff: Phase 20 — Consumer Submission / Cal-ITP Readiness, if the roadmap still applies.

## What Was Implemented

- Added deterministic replay evaluation in `internal/realtimequality`.
- Added replay fixtures in `testdata/replay/` for matched current behavior, stale telemetry, ambiguous assignment, low-confidence Trip Updates withholding, manual override precedence, and cancellation/added-trip/short-turn/detour diagnostics.
- Added `testdata/replay/README.md` documenting the replay fixture schema and required fields.
- Added explicit prediction quality metrics for assignment confidence buckets, unknown assignments, ambiguous assignments, degraded assignments, stale telemetry rows, manual override usage, and rate objects with numerators/denominators.
- Added explicit rate denominator definitions for unknown assignment rate, ambiguous assignment rate, stale telemetry rate, Trip Updates coverage, future-stop coverage, and withheld-by-reason counts.
- Added zero-denominator handling with `not_applicable` rate status rather than reporting misleading `0%` coverage.
- Added regression tests proving unknown, ambiguous, stale, withheld, and degraded cases remain visible in metrics and diagnostics.
- Added authenticated Operations Console Trip Updates quality summaries from existing `feed_health_snapshot` diagnostics.
- Added Operations Console fallback text: `no Trip Updates diagnostics recorded yet`.
- Added `make realtime-quality` as a focused replay command.

## What Was Designed But Intentionally Not Implemented Yet

- No TheTransitClock or other external predictor was integrated.
- No broad ETA or schedule-deviation algorithm change was made; current deterministic behavior is now measured first.
- No rider-facing app, fares/payments, CAD/dispatch, consumer submission API, or fake submission automation was added.
- No public feed URL, GTFS-RT protobuf contract, unauthenticated route, consumer-submission status, or evidence-claim change was made.

## Schema And Interface Changes

- No database migration was added.
- `prediction.Metrics` now includes explicit quality counts and `RateMetric` fields.
- `feed_health_snapshot.details_json.prediction_metrics` receives the expanded metrics shape through the existing persistence path.
- `internal/compliance.PostgresRepository` now exposes latest Trip Updates diagnostics summary for authenticated Operations Console rendering.
- `internal/realtimequality` defines replay scenario, expectation, report, and comparison helpers.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- Added `internal/realtimequality` replay tests for fixture execution, deterministic fixed-clock behavior, stable report comparison, uncertainty regression guards, and zero-denominator rate behavior.
- Added Operations Console tests for no-diagnostics fallback and safe Trip Updates quality summary rendering.
- Updated Vehicle Positions debug expectation for degraded unknown assignments so ambiguity/unknown degradation remains visible instead of being collapsed to generic not-in-service.
- Updated Trip Updates diagnostics persistence test to avoid a misleading fallback coverage percent when no eligible denominator is recorded.

## Checks Run And Blocked Checks

Pre-edit baseline:

- `make validate` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `git diff --check` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.

Focused implementation checks:

- `go test ./internal/realtimequality` — passed.
- `go test ./internal/prediction ./internal/feed ./cmd/agency-config ./internal/realtimequality` — passed.

Final post-edit checks:

- `make validate` — passed.
- `make realtime-quality` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `git diff --check` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make demo-agency-flow` — passed.
- `make agency-app-up` — passed.
- `make agency-app-down` — passed.
- `docker compose -f deploy/docker-compose.yml --profile app config` — passed.

Blocked checks:

- None.

## Known Issues

- Replay fixtures are small baseline scenarios, not proof of production-grade ETA quality.
- Trip Updates remain conservative and withhold many hard cases.
- Canceled trips can emit `CANCELED` TripUpdates, but missing Service Alerts remain a diagnostics/review signal when alert linkage is not present.
- Added-trip, short-turn, and detour Trip Updates remain intentionally withheld without a safe prediction model.
- The Operations Console quality summary depends on recorded `feed_health_snapshot` diagnostics and does not infer health when no diagnostics exist.

## Exact Next-Step Recommendation

- First files to read:
  - `AGENTS.md`
  - `docs/current-status.md`
  - `docs/handoffs/latest.md`
  - `docs/handoffs/phase-19.md`
  - `testdata/replay/README.md`
- First files likely to edit:
  - `internal/realtimequality/`
  - `testdata/replay/`
  - `internal/prediction/`
  - `internal/state/`
  - `internal/feed/tripupdates/`
- Commands to run before coding:
  - `make validate`
  - `make realtime-quality`
  - `make test`
  - `make smoke`
  - `make test-integration`
  - `git diff --check`
  - `docker compose -f deploy/docker-compose.yml config`
- Known blockers:
  - None.
- Recommended first implementation slice:
  - Phase 20 should focus on consumer-submission and Cal-ITP readiness evidence without changing Phase 19 quality claims. If more realtime work is chosen instead, add more replay fixtures before changing ETA logic, especially for after-midnight, frequency-window, block continuity, and real-world stale/ambiguous telemetry patterns.
