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
  automatic matching.
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
