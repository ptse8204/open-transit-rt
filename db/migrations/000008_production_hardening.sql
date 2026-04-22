-- +goose Up
ALTER TABLE device_credential
  ADD COLUMN IF NOT EXISTS vehicle_id TEXT,
  ADD COLUMN IF NOT EXISTS last_used_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS device_credential_active_lookup_idx
  ON device_credential (agency_id, device_id, token_hash)
  WHERE status = 'active' AND revoked_at IS NULL AND rotated_at IS NULL;

WITH ranked AS (
  SELECT
    id,
    active_from,
    row_number() OVER (
      PARTITION BY agency_id, vehicle_id
      ORDER BY active_from DESC, id DESC
    ) AS row_rank
  FROM vehicle_trip_assignment
  WHERE active_to IS NULL
)
UPDATE vehicle_trip_assignment v
SET active_to = ranked.active_from
FROM ranked
WHERE v.id = ranked.id
  AND ranked.row_rank > 1;

CREATE UNIQUE INDEX IF NOT EXISTS vehicle_trip_assignment_current_uidx
  ON vehicle_trip_assignment (agency_id, vehicle_id)
  WHERE active_to IS NULL;

ALTER TABLE incident
  ADD COLUMN IF NOT EXISTS dedupe_key TEXT,
  ADD COLUMN IF NOT EXISTS last_seen_at TIMESTAMPTZ,
  ADD COLUMN IF NOT EXISTS occurrence_count INTEGER NOT NULL DEFAULT 1;

UPDATE incident
SET last_seen_at = COALESCE(last_seen_at, created_at),
    occurrence_count = GREATEST(occurrence_count, 1);

CREATE UNIQUE INDEX IF NOT EXISTS incident_prediction_review_active_dedupe_uidx
  ON incident (agency_id, incident_type, dedupe_key)
  WHERE incident_type = 'prediction_review'
    AND dedupe_key IS NOT NULL
    AND status IN ('open', 'acknowledged', 'deferred');

-- +goose Down
DROP INDEX IF EXISTS incident_prediction_review_active_dedupe_uidx;

ALTER TABLE incident
  DROP COLUMN IF EXISTS occurrence_count,
  DROP COLUMN IF EXISTS last_seen_at,
  DROP COLUMN IF EXISTS dedupe_key;

DROP INDEX IF EXISTS vehicle_trip_assignment_current_uidx;
DROP INDEX IF EXISTS device_credential_active_lookup_idx;

ALTER TABLE device_credential
  DROP COLUMN IF EXISTS last_used_at,
  DROP COLUMN IF EXISTS vehicle_id;
