# Phase 29 Handoff

## Phase

Phase 29 — Real-World Realtime Quality Expansion

## Status

- Complete for the synthetic replay evidence expansion scope.
- Active phase after this handoff: Phase 30 — recommended evidence-refresh or adoption-readiness slice.

## What Was Implemented

- Added richer deterministic replay fixtures under `testdata/replay/`.
- Added replay trip `frequencies` fixture support.
- Added optional replay manual override `expires_at` support.
- Updated the replay telemetry repository to return latest telemetry per vehicle for feed snapshots.
- Strengthened replay comparison for already-recorded cancellation alert linkage and unsupported disruption-withheld metrics.
- Added focused realtime-quality tests for Phase 29 replay behavior.
- Updated replay docs, Phase 29 docs, current status, and latest handoff.

## What Was Intentionally Deferred

- Real-world observed-arrival/departure evidence and ETA accuracy comparison.
- Real route/time-period quality metrics, because no real deployment or observed-arrival evidence exists in the repo.
- External predictor integration, including TheTransitClock.
- Operations Console changes.
- Public feed URL, GTFS-RT protobuf contract, consumer status, auth boundary, external dependency, and evidence-claim changes.

## Replay Fixtures Added

- `after-midnight-service.json`
- `frequency-exact-window.json`
- `frequency-non-exact-window.json`
- `block-continuity-transition.json`
- `long-layover-withheld.json`
- `sparse-telemetry-near-stale-threshold.json`
- `noisy-off-shape-gps-degraded.json`
- `stale-ambiguous-hard-pattern.json`
- `cancellation-alert-linkage.json`
- `manual-override-before-expiry.json`
- `manual-override-after-expiry.json`

All fixtures use synthetic, public-safe identifiers and telemetry only.

## Metrics Added Or Confirmed

- Confirmed existing denominator-aware metrics remain the primary public diagnostics shape.
- Confirmed `not_applicable` zero-denominator behavior remains intact.
- Confirmed unknown, ambiguous, stale, degraded, withheld-by-reason, degraded-by-reason, manual override, cancellation alert linkage, and unsupported disruption-withheld metrics remain visible.
- Added stricter replay assertions for:
  - `canceled_trips_emitted`
  - `cancellation_alert_links_expected`
  - `cancellation_alert_links_missing`
  - `added_trips_withheld`
  - `short_turns_withheld`
  - `detours_withheld`
  - `degraded_by_reason`

Real route/time-period quality metrics were not added. Synthetic fixture routes and timestamps are scenario labels only, not deployment coverage metrics.

## Deterministic Improvements

- No production matcher or predictor algorithm was changed.
- Frequency additions are replay-fixture modeling support. Production frequency matching and prediction behavior were already present in existing matcher, Vehicle Positions, and deterministic predictor code.
- Manual override expiry additions are replay-harness modeling support. Production expiry semantics were already reflected by the state model and repository contract; Phase 29 made them representable in replay fixtures.
- Replay latest-telemetry behavior now matches the production telemetry repository contract for feed snapshots.

## Observed-Arrival Comparison Status

Unavailable. Phase 29 expands synthetic replay coverage only. No real observed stop-time, arrival, departure, or retained real-world ETA evidence exists in the repo for Phase 29, so no real-world observed-arrival ETA accuracy claim is supported.

## Operations Console Visibility

None added. Existing authenticated Operations Console quality summaries remain unchanged.

## Tests Added

- Replay fixture glob now covers all Phase 19 and Phase 29 fixtures.
- Added focused tests for:
  - after-midnight service and Trip Updates output;
  - exact frequency scheduled identity;
  - non-exact frequency conservative unscheduled identity;
  - manual override before expiry;
  - manual override after expiry returning to automatic matching;
  - cancellation alert linkage counts;
  - zero-denominator `not_applicable`;
  - stale/ambiguous hard-pattern visibility.

## Commands Run

- Pre-edit/planning `make validate` — passed.
- Pre-edit/planning `make realtime-quality` — passed.
- Pre-edit/planning `make test` — passed.
- Pre-edit/planning `make smoke` — passed.
- Pre-edit/planning `make test-integration` — passed.
- Pre-edit/planning `git diff --check` — passed.
- Pre-edit/planning `docker compose -f deploy/docker-compose.yml config` — passed.
- Implementation focused `go test ./internal/realtimequality` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make test` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.

## Blocked Commands

- None so far.

## Known Remaining Realtime Quality Gaps

- No real observed-arrival/departure comparison.
- No real route/time-period coverage metrics.
- No real agency AVL quality evidence.
- No production-grade ETA quality evidence.
- No external predictor contract or integration beyond the existing adapter boundary.
- Added trips, short turns, and detours remain explicitly withheld rather than predicted.
- Cancellation alert linkage still identifies missing alert linkage; it does not author Service Alerts automatically.
- Manual override authoring UI remains limited; Phase 29 only expanded replay modeling.

## Exact Recommendation For Phase 30

Start Phase 30 with evidence refresh and claim-boundary work, not ETA claim expansion. Recommended first slice: produce an evidence inventory that separates synthetic replay evidence, template-only evidence, hosted/operator pilot evidence, and missing real-world inputs. Explicitly list the real observed-arrival/departure data, real deployment route/time-period coverage, and agency-approved AVL evidence required before any stronger ETA quality or production readiness claim can be considered.
