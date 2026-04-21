-- +goose Up
ALTER TABLE feed_config
  ADD COLUMN publication_environment TEXT NOT NULL DEFAULT 'dev'
  CHECK (publication_environment IN ('dev', 'production'));

CREATE TABLE service_alert (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  alert_key TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
  cause TEXT NOT NULL DEFAULT 'unknown_cause',
  effect TEXT NOT NULL DEFAULT 'unknown_effect',
  header_text TEXT NOT NULL,
  description_text TEXT,
  url TEXT,
  active_start TIMESTAMPTZ,
  active_end TIMESTAMPTZ,
  feed_version_id TEXT REFERENCES feed_version(id) ON DELETE SET NULL,
  source_type TEXT NOT NULL DEFAULT 'operator' CHECK (source_type IN ('operator', 'cancellation_reconciler')),
  source_id TEXT,
  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_by TEXT NOT NULL DEFAULT 'system',
  updated_by TEXT NOT NULL DEFAULT 'system',
  published_at TIMESTAMPTZ,
  archived_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, alert_key)
);

CREATE INDEX service_alert_public_idx
  ON service_alert (agency_id, status, active_start, active_end);

CREATE TABLE service_alert_informed_entity (
  id BIGSERIAL PRIMARY KEY,
  service_alert_id BIGINT NOT NULL REFERENCES service_alert(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  route_id TEXT,
  stop_id TEXT,
  trip_id TEXT,
  start_date TEXT,
  start_time TEXT,
  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX service_alert_informed_entity_alert_idx
  ON service_alert_informed_entity (service_alert_id);

CREATE INDEX service_alert_informed_entity_trip_idx
  ON service_alert_informed_entity (agency_id, trip_id, start_date, start_time)
  WHERE trip_id IS NOT NULL;

CREATE TABLE compliance_scorecard_snapshot (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  snapshot_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  publication_environment TEXT NOT NULL CHECK (publication_environment IN ('dev', 'production')),
  overall_status TEXT NOT NULL CHECK (overall_status IN ('red', 'yellow', 'green')),
  schedule_status TEXT NOT NULL CHECK (schedule_status IN ('red', 'yellow', 'green')),
  vehicle_positions_status TEXT NOT NULL CHECK (vehicle_positions_status IN ('red', 'yellow', 'green')),
  trip_updates_status TEXT NOT NULL CHECK (trip_updates_status IN ('red', 'yellow', 'green')),
  alerts_status TEXT NOT NULL CHECK (alerts_status IN ('red', 'yellow', 'green')),
  validation_status TEXT NOT NULL CHECK (validation_status IN ('red', 'yellow', 'green')),
  discoverability_status TEXT NOT NULL CHECK (discoverability_status IN ('red', 'yellow', 'green')),
  consumer_ingestion_status TEXT NOT NULL CHECK (consumer_ingestion_status IN ('red', 'yellow', 'green')),
  details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX compliance_scorecard_snapshot_agency_idx
  ON compliance_scorecard_snapshot (agency_id, snapshot_at DESC);

-- +goose Down
DROP INDEX IF EXISTS compliance_scorecard_snapshot_agency_idx;
DROP TABLE IF EXISTS compliance_scorecard_snapshot;

DROP INDEX IF EXISTS service_alert_informed_entity_trip_idx;
DROP INDEX IF EXISTS service_alert_informed_entity_alert_idx;
DROP TABLE IF EXISTS service_alert_informed_entity;

DROP INDEX IF EXISTS service_alert_public_idx;
DROP TABLE IF EXISTS service_alert;

ALTER TABLE feed_config
  DROP COLUMN IF EXISTS publication_environment;
