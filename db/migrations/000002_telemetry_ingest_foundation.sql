-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
  payload_type TEXT;
BEGIN
  SELECT data_type
  INTO payload_type
  FROM information_schema.columns
  WHERE table_name = 'telemetry_event'
    AND column_name = 'payload_json';

  IF payload_type IS DISTINCT FROM 'jsonb' THEN
    IF payload_type = 'json' THEN
      ALTER TABLE telemetry_event
        ALTER COLUMN payload_json TYPE JSONB USING COALESCE(payload_json::jsonb, '{}'::jsonb);
    ELSE
      ALTER TABLE telemetry_event
        ALTER COLUMN payload_json TYPE JSONB USING COALESCE(to_jsonb(payload_json), '{}'::jsonb);
    END IF;
  END IF;
END
$$;
-- +goose StatementEnd

UPDATE telemetry_event
  SET payload_json = '{}'::jsonb
  WHERE payload_json IS NULL;

ALTER TABLE telemetry_event
  ALTER COLUMN payload_json SET DEFAULT '{}'::jsonb,
  ALTER COLUMN payload_json SET NOT NULL;

ALTER TABLE telemetry_event
  DROP CONSTRAINT IF EXISTS telemetry_event_agency_id_device_id_vehicle_id_observed_at_key;

CREATE UNIQUE INDEX IF NOT EXISTS telemetry_event_accepted_vehicle_observed_uidx
  ON telemetry_event (agency_id, vehicle_id, observed_at)
  WHERE ingest_status = 'accepted';

CREATE INDEX IF NOT EXISTS telemetry_event_latest_accepted_idx
  ON telemetry_event (agency_id, vehicle_id, observed_at DESC)
  WHERE ingest_status = 'accepted';

CREATE INDEX IF NOT EXISTS telemetry_event_agency_received_idx
  ON telemetry_event (agency_id, received_at DESC, id DESC);

-- +goose Down
DROP INDEX IF EXISTS telemetry_event_agency_received_idx;
DROP INDEX IF EXISTS telemetry_event_latest_accepted_idx;
DROP INDEX IF EXISTS telemetry_event_accepted_vehicle_observed_uidx;

ALTER TABLE telemetry_event
  ALTER COLUMN payload_json DROP NOT NULL,
  ALTER COLUMN payload_json DROP DEFAULT;

-- +goose StatementBegin
DO $$
BEGIN
  ALTER TABLE telemetry_event
    ADD CONSTRAINT telemetry_event_agency_id_device_id_vehicle_id_observed_at_key
    UNIQUE (agency_id, device_id, vehicle_id, observed_at);
EXCEPTION
  WHEN duplicate_object THEN
    NULL;
  WHEN unique_violation THEN
    RAISE NOTICE 'skipping telemetry_event uniqueness restoration because existing rows violate the old key';
END
$$;
-- +goose StatementEnd
