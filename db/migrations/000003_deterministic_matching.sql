-- +goose Up
ALTER TABLE vehicle_trip_assignment
  ALTER COLUMN service_date DROP NOT NULL,
  ADD COLUMN block_id TEXT,
  ADD COLUMN telemetry_event_id BIGINT REFERENCES telemetry_event(id) ON DELETE SET NULL,
  ADD COLUMN degraded_state TEXT NOT NULL DEFAULT 'none'
    CHECK (degraded_state IN ('none', 'unknown', 'stale', 'ambiguous', 'missing_schedule_data', 'missing_shape', 'low_confidence')),
  ADD COLUMN score_details_json JSONB NOT NULL DEFAULT '{}'::jsonb;

ALTER TABLE incident
  ADD COLUMN vehicle_trip_assignment_id BIGINT REFERENCES vehicle_trip_assignment(id) ON DELETE SET NULL;

CREATE INDEX vehicle_trip_assignment_current_vehicle_idx
  ON vehicle_trip_assignment (agency_id, vehicle_id, active_to, active_from DESC);

CREATE INDEX vehicle_trip_assignment_degraded_idx
  ON vehicle_trip_assignment (agency_id, degraded_state, created_at DESC)
  WHERE degraded_state <> 'none';

CREATE INDEX vehicle_trip_assignment_telemetry_event_idx
  ON vehicle_trip_assignment (telemetry_event_id)
  WHERE telemetry_event_id IS NOT NULL;

CREATE INDEX incident_assignment_idx
  ON incident (vehicle_trip_assignment_id)
  WHERE vehicle_trip_assignment_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS incident_assignment_idx;
DROP INDEX IF EXISTS vehicle_trip_assignment_telemetry_event_idx;
DROP INDEX IF EXISTS vehicle_trip_assignment_degraded_idx;
DROP INDEX IF EXISTS vehicle_trip_assignment_current_vehicle_idx;

ALTER TABLE incident
  DROP COLUMN IF EXISTS vehicle_trip_assignment_id;

ALTER TABLE vehicle_trip_assignment
  DROP COLUMN IF EXISTS score_details_json,
  DROP COLUMN IF EXISTS degraded_state,
  DROP COLUMN IF EXISTS telemetry_event_id,
  DROP COLUMN IF EXISTS block_id,
  ALTER COLUMN service_date SET NOT NULL;
