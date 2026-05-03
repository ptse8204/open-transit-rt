# Realtime Quality Replay Fixtures

Replay fixtures are deterministic JSON scenarios used by `internal/realtimequality`.
They measure current matcher, Vehicle Positions, and Trip Updates behavior without
live services.

Required top-level fields:

- `name`: stable scenario name.
- `input_fixture`: short label or source fixture path for the scenario inputs.
- `agency_id`, `timezone`, `feed_version_id`: fixed agency/feed context.
- `generated_at`: fixed replay clock timestamp. Replay tests must not use wall-clock
  time except through this injected value.
- `trips`: schedule candidates used by the matcher and deterministic predictor.
  Trips may include optional `frequencies` rows with `start_time`, `end_time`,
  `headway_secs`, and `exact_times` when a scenario needs frequency-window
  identity.
- `telemetry_events`: input telemetry rows with fixed IDs and timestamps.
- `expected.assignments`: expected current assignment reports.
- `expected.vehicle_positions`: expected Vehicle Positions publication decisions.
- `expected.trip_updates`: expected Trip Updates output summaries.
- `expected.withheld_reasons`: expected raw Trip Updates withheld-by-reason counts.
- `expected.metrics`: expected quality metrics and rate denominators.

Optional top-level fields:

- `input_assignments`: fixed current assignments used when the scenario is focused on
  predictor behavior instead of matcher output.
- `manual_overrides`: matcher manual overrides. These verify operator authority over
  automatic matching. Overrides may include `expires_at` for replay-only expiry
  modeling.
- `prediction_overrides`: Trip Updates disruption overrides such as canceled trips,
  added trips, short turns, and detours.

Rate denominator definitions:

- `unknown_assignment_rate`: current unknown assignments / current assignments considered.
- `ambiguous_assignment_rate`: current ambiguous assignments / current assignments considered.
- `stale_telemetry_rate`: stale latest telemetry rows / telemetry rows considered.
- `trip_updates_coverage_rate`: emitted non-canceled Trip Updates / eligible in-service ETA candidates.
- `future_stop_coverage_rate`: non-canceled Trip Updates with at least one future stop update / eligible in-service ETA candidates.
- `withheld_by_reason`: raw counts grouped by reason, not a rate.

Zero-denominator metrics must use `status: "not_applicable"` with no percent value.
Do not encode `0%` when there were no eligible candidates.

Reports must remain stable: fixture order is fixed, timestamps are fixed, and replay
output slices are sorted before comparison.

## Scenario Catalog

All scenarios are synthetic replay evidence. They prove deterministic behavior for
the committed fixture inputs only. They do not prove production-grade ETA quality,
real-world ETA accuracy, consumer acceptance, agency endorsement, hosted SaaS
availability, or CAL-ITP/Caltrans compliance.

| Fixture | Represents | Proves | Does Not Prove |
| --- | --- | --- | --- |
| `matched-current-behavior.json` | Baseline matched in-service trip. | Vehicle Positions can publish a trip descriptor and Trip Updates can emit one future stop for a defensible match. | Real-world accuracy or broad route coverage. |
| `stale-telemetry-withheld.json` | Stale latest telemetry. | Assignment becomes unknown/stale, Vehicle Positions withholds the trip descriptor, Trip Updates withholds with `stale_telemetry`, and stale rate denominator is explicit. | Recovery behavior after fresh telemetry arrives. |
| `ambiguous-assignment-visible.json` | Two equally plausible trips. | Ambiguity remains visible as unknown/degraded, Trip Updates are withheld, and zero eligible prediction denominators are `not_applicable`. | Automatic disambiguation on real telemetry. |
| `low-confidence-trip-updates-withheld.json` | Assignment below prediction confidence threshold. | Trip Updates withhold with `below_confidence_threshold`. | A tuned production confidence threshold. |
| `manual-override-precedence.json` | Manual override precedence. | Manual override is authoritative over automatic matching and its withheld reason remains visible. | Timed expiry; use the Phase 29 expiry fixtures for that. |
| `disruption-diagnostics-baseline.json` | Canceled trip plus unsupported added-trip, short-turn, and detour overrides. | Cancellation output and disruption withheld counts remain repeatable. | Full Alerts authoring linkage or support for added trips, short turns, or detours. |
| `after-midnight-service.json` | Trip with GTFS times after `24:00:00`. | Previous-service-day matching, Vehicle Positions trip descriptor publication, and scheduled Trip Updates remain deterministic after midnight. | Multi-day service lookback beyond the existing service-day behavior. |
| `frequency-exact-window.json` | `frequencies.txt` exact-times instance. | Exact frequency identity keeps a scheduled Trip Updates relationship and a distinct `start_time`. | Any production behavior change; this exercises existing matcher/predictor frequency support through replay fixtures. |
| `frequency-non-exact-window.json` | Non-exact frequency window. | Non-exact frequency matching remains conservative, visible through `frequency_non_exact_conservative`, and Trip Updates use `unscheduled`. | Exact vehicle-to-departure certainty inside non-exact headways. |
| `block-continuity-transition.json` | Same vehicle transitioning from one block trip to the next. | Latest telemetry per vehicle is used for feed snapshots, while matcher assignment reasons can expose `block_transition_match`. | Complex dispatch changes or unsupported block edits. |
| `long-layover-withheld.json` | Operator-marked long layover. | Vehicle Positions omits trip descriptor as not in service and Trip Updates withholds with `layover_no_prediction`. | Automated layover inference from real operations data. |
| `sparse-telemetry-near-stale-threshold.json` | Fresh but sparse telemetry near the stale threshold. | Stale denominator handling remains explicit and the case stays measured, not `not_applicable`. | Reliability under prolonged sparse AVL feeds. |
| `noisy-off-shape-gps-degraded.json` | Fresh off-shape/noisy GPS point. | Low-confidence/off-shape evidence remains unknown/degraded and withheld instead of becoming a false trip prediction. | Real GPS noise modeling or correction. |
| `stale-ambiguous-hard-pattern.json` | Combined stale and ambiguous vehicles. | Unknown, stale, ambiguous, degraded, withheld, and degraded-by-reason metrics stay visible together. | Real-world stale/ambiguous rates. |
| `cancellation-alert-linkage.json` | Canceled trip requiring alert linkage review. | Cancellation Trip Updates and missing-alert linkage counters remain repeatable. | Actual Service Alerts authoring or third-party alert acceptance. |
| `manual-override-before-expiry.json` | Timed manual override before expiry. | Manual override is authoritative before expiry and diagnostics show manual influence without raw private details. | Production UI for timed override authoring. |
| `manual-override-after-expiry.json` | Timed manual override after expiry. | Expired replay override returns to automatic matching and prediction eligibility. | A new production expiry algorithm; production expiry is already part of the state repository contract. |

Expected Vehicle Positions behavior is encoded in `expected.vehicle_positions` for
each fixture. Expected Trip Updates behavior is encoded in `expected.trip_updates`
and `expected.metrics.withheld_by_reason`. Expected metrics include explicit
numerators and denominators for rates; use `not_applicable` when a denominator is
zero.

Real route/time-period quality metrics are not represented here. Fixture names,
routes, and timestamps are synthetic scenario coverage only, not deployment
coverage or observed-arrival evidence.
