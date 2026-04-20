-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE agency (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  timezone TEXT NOT NULL,
  contact_email TEXT,
  public_url TEXT,
  branding_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE agency_user (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  email TEXT NOT NULL,
  display_name TEXT,
  auth_subject TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, email)
);

CREATE TABLE role_binding (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  agency_user_id BIGINT NOT NULL REFERENCES agency_user(id) ON DELETE CASCADE,
  role TEXT NOT NULL CHECK (role IN ('admin', 'editor', 'operator', 'read_only')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, agency_user_id, role)
);

CREATE TABLE device_credential (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  device_id TEXT NOT NULL,
  token_hash TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'revoked', 'rotated')),
  valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
  rotated_at TIMESTAMPTZ,
  revoked_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, device_id)
);

CREATE TABLE feed_config (
  agency_id TEXT PRIMARY KEY REFERENCES agency(id) ON DELETE CASCADE,
  public_base_url TEXT NOT NULL,
  feed_base_url TEXT NOT NULL,
  technical_contact_email TEXT NOT NULL,
  license_name TEXT NOT NULL,
  license_url TEXT,
  metrics_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  validator_strictness TEXT NOT NULL DEFAULT 'warn' CHECK (validator_strictness IN ('off', 'warn', 'block')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE feed_version (
  id TEXT PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  source_type TEXT NOT NULL CHECK (source_type IN ('gtfs_import', 'gtfs_studio', 'seed')),
  lifecycle_state TEXT NOT NULL DEFAULT 'staged' CHECK (lifecycle_state IN ('draft', 'staged', 'active', 'retired', 'failed')),
  published_at TIMESTAMPTZ,
  activated_at TIMESTAMPTZ,
  retired_at TIMESTAMPTZ,
  is_active BOOLEAN NOT NULL DEFAULT FALSE,
  validation_status TEXT NOT NULL DEFAULT 'not_run' CHECK (validation_status IN ('not_run', 'passed', 'warning', 'failed')),
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX feed_version_one_active_per_agency
  ON feed_version (agency_id)
  WHERE is_active;

CREATE TABLE published_feed (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  feed_type TEXT NOT NULL CHECK (feed_type IN ('schedule', 'vehicle_positions', 'trip_updates', 'alerts')),
  canonical_public_url TEXT NOT NULL,
  license_name TEXT NOT NULL,
  license_url TEXT,
  contact_email TEXT NOT NULL,
  revision_timestamp TIMESTAMPTZ,
  activation_status TEXT NOT NULL DEFAULT 'inactive' CHECK (activation_status IN ('inactive', 'active', 'rolled_back', 'unhealthy')),
  active_feed_version_id TEXT REFERENCES feed_version(id),
  metadata_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, feed_type)
);

CREATE TABLE gtfs_route (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  short_name TEXT,
  long_name TEXT,
  route_type INTEGER,
  PRIMARY KEY (id, feed_version_id)
);

CREATE TABLE gtfs_stop (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  geom geometry(Point, 4326),
  PRIMARY KEY (id, feed_version_id)
);

CREATE INDEX gtfs_stop_geom_idx ON gtfs_stop USING GIST (geom);

CREATE TABLE gtfs_calendar (
  service_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  monday BOOLEAN NOT NULL,
  tuesday BOOLEAN NOT NULL,
  wednesday BOOLEAN NOT NULL,
  thursday BOOLEAN NOT NULL,
  friday BOOLEAN NOT NULL,
  saturday BOOLEAN NOT NULL,
  sunday BOOLEAN NOT NULL,
  start_date TEXT NOT NULL,
  end_date TEXT NOT NULL,
  PRIMARY KEY (service_id, feed_version_id)
);

CREATE TABLE gtfs_calendar_date (
  service_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  date TEXT NOT NULL,
  exception_type INTEGER NOT NULL CHECK (exception_type IN (1, 2)),
  PRIMARY KEY (service_id, feed_version_id, date)
);

CREATE TABLE gtfs_trip (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  route_id TEXT NOT NULL,
  service_id TEXT NOT NULL,
  block_id TEXT,
  shape_id TEXT,
  direction_id INTEGER,
  PRIMARY KEY (id, feed_version_id),
  FOREIGN KEY (route_id, feed_version_id) REFERENCES gtfs_route(id, feed_version_id)
);

CREATE INDEX gtfs_trip_route_idx ON gtfs_trip (agency_id, feed_version_id, route_id);
CREATE INDEX gtfs_trip_block_idx ON gtfs_trip (agency_id, feed_version_id, block_id);

CREATE TABLE gtfs_stop_time (
  trip_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  arrival_time TEXT,
  departure_time TEXT,
  stop_id TEXT NOT NULL,
  stop_sequence INTEGER NOT NULL,
  pickup_type INTEGER,
  drop_off_type INTEGER,
  shape_dist_traveled DOUBLE PRECISION,
  PRIMARY KEY (trip_id, feed_version_id, stop_sequence),
  FOREIGN KEY (trip_id, feed_version_id) REFERENCES gtfs_trip(id, feed_version_id) ON DELETE CASCADE,
  FOREIGN KEY (stop_id, feed_version_id) REFERENCES gtfs_stop(id, feed_version_id)
);

CREATE TABLE gtfs_shape_point (
  shape_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  sequence INTEGER NOT NULL,
  dist_traveled DOUBLE PRECISION,
  geom geometry(Point, 4326),
  PRIMARY KEY (shape_id, feed_version_id, sequence)
);

CREATE INDEX gtfs_shape_point_geom_idx ON gtfs_shape_point USING GIST (geom);

CREATE TABLE gtfs_shape_line (
  shape_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  geom geometry(LineString, 4326),
  PRIMARY KEY (shape_id, feed_version_id)
);

CREATE INDEX gtfs_shape_line_geom_idx ON gtfs_shape_line USING GIST (geom);

CREATE TABLE gtfs_frequency (
  trip_id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  start_time TEXT NOT NULL,
  end_time TEXT NOT NULL,
  headway_secs INTEGER NOT NULL,
  exact_times INTEGER NOT NULL DEFAULT 0 CHECK (exact_times IN (0, 1)),
  PRIMARY KEY (trip_id, feed_version_id, start_time),
  FOREIGN KEY (trip_id, feed_version_id) REFERENCES gtfs_trip(id, feed_version_id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft (
  id TEXT PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'staged', 'published', 'discarded')),
  base_feed_version_id TEXT REFERENCES feed_version(id),
  created_by TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE gtfs_draft_record (
  id BIGSERIAL PRIMARY KEY,
  draft_id TEXT NOT NULL REFERENCES gtfs_draft(id) ON DELETE CASCADE,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  record_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (draft_id, entity_type, entity_id)
);

CREATE TABLE telemetry_event (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  device_id TEXT NOT NULL,
  vehicle_id TEXT NOT NULL,
  observed_at TIMESTAMPTZ NOT NULL,
  received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  lat DOUBLE PRECISION NOT NULL CHECK (lat >= -90 AND lat <= 90),
  lon DOUBLE PRECISION NOT NULL CHECK (lon >= -180 AND lon <= 180),
  geom geometry(Point, 4326),
  bearing DOUBLE PRECISION,
  speed_mps DOUBLE PRECISION,
  accuracy_m DOUBLE PRECISION,
  trip_hint TEXT,
  payload_json JSONB,
  ingest_status TEXT NOT NULL DEFAULT 'accepted' CHECK (ingest_status IN ('accepted', 'duplicate', 'out_of_order', 'rejected')),
  UNIQUE (agency_id, device_id, vehicle_id, observed_at)
);

CREATE INDEX telemetry_event_vehicle_observed_idx
  ON telemetry_event (agency_id, vehicle_id, observed_at DESC);
CREATE INDEX telemetry_event_geom_idx ON telemetry_event USING GIST (geom);

CREATE TABLE manual_override (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  vehicle_id TEXT NOT NULL,
  override_type TEXT NOT NULL CHECK (override_type IN ('trip_assignment', 'service_state', 'canceled_trip', 'added_trip', 'vehicle_swap', 'detour', 'short_turn')),
  route_id TEXT,
  trip_id TEXT,
  start_date TEXT,
  start_time TEXT,
  state TEXT NOT NULL CHECK (state IN ('unknown', 'in_service', 'layover', 'deadhead', 'out_of_service', 'canceled', 'added', 'detour', 'short_turn')),
  expires_at TIMESTAMPTZ,
  cleared_at TIMESTAMPTZ,
  reason TEXT,
  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX manual_override_active_idx
  ON manual_override (agency_id, vehicle_id, expires_at, cleared_at);

CREATE TABLE vehicle_trip_assignment (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  vehicle_id TEXT NOT NULL,
  feed_version_id TEXT REFERENCES feed_version(id),
  service_date TEXT NOT NULL,
  route_id TEXT,
  trip_id TEXT,
  start_date TEXT,
  start_time TEXT,
  current_stop_sequence INTEGER,
  shape_dist_traveled DOUBLE PRECISION,
  state TEXT NOT NULL CHECK (state IN ('unknown', 'in_service', 'layover', 'deadhead', 'out_of_service')),
  confidence DOUBLE PRECISION NOT NULL CHECK (confidence >= 0 AND confidence <= 1),
  assignment_source TEXT NOT NULL CHECK (assignment_source IN ('automatic', 'manual_override', 'imported', 'unknown')),
  reason_codes TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
  active_from TIMESTAMPTZ NOT NULL,
  active_to TIMESTAMPTZ,
  manual_override_id BIGINT REFERENCES manual_override(id),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX vehicle_trip_assignment_current_idx
  ON vehicle_trip_assignment (agency_id, vehicle_id, active_to, active_from DESC);

CREATE TABLE incident (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  incident_type TEXT NOT NULL,
  severity TEXT NOT NULL CHECK (severity IN ('info', 'warning', 'critical')),
  route_id TEXT,
  vehicle_id TEXT,
  trip_id TEXT,
  status TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'acknowledged', 'resolved')),
  details_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at TIMESTAMPTZ
);

CREATE INDEX incident_queue_idx ON incident (agency_id, status, severity, created_at DESC);

CREATE TABLE validation_report (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  feed_version_id TEXT REFERENCES feed_version(id) ON DELETE CASCADE,
  feed_type TEXT NOT NULL CHECK (feed_type IN ('schedule', 'vehicle_positions', 'trip_updates', 'alerts')),
  validator_name TEXT NOT NULL,
  validator_version TEXT,
  status TEXT NOT NULL CHECK (status IN ('not_run', 'passed', 'warning', 'failed')),
  error_count INTEGER NOT NULL DEFAULT 0,
  warning_count INTEGER NOT NULL DEFAULT 0,
  info_count INTEGER NOT NULL DEFAULT 0,
  report_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE feed_health_snapshot (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  feed_type TEXT NOT NULL CHECK (feed_type IN ('schedule', 'vehicle_positions', 'trip_updates', 'alerts')),
  snapshot_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  endpoint_available BOOLEAN,
  freshness_seconds DOUBLE PRECISION,
  generation_latency_ms DOUBLE PRECISION,
  invalid_response_percent DOUBLE PRECISION,
  matched_vehicle_percent DOUBLE PRECISION,
  coverage_percent DOUBLE PRECISION,
  details_json JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX feed_health_snapshot_idx ON feed_health_snapshot (agency_id, feed_type, snapshot_at DESC);

CREATE TABLE consumer_ingestion (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  consumer_name TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'submitted', 'accepted', 'rejected', 'pending_fix', 'resubmitted')),
  submitted_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  notes TEXT,
  packet_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  UNIQUE (agency_id, consumer_name)
);

CREATE TABLE marketplace_gap (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT REFERENCES agency(id) ON DELETE CASCADE,
  gap_key TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'not_started' CHECK (status IN ('not_started', 'in_progress', 'complete', 'not_applicable')),
  notes TEXT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, gap_key)
);

CREATE TABLE audit_log (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  actor_id TEXT NOT NULL,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT,
  old_value_json JSONB,
  new_value_json JSONB,
  reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX audit_log_agency_entity_idx
  ON audit_log (agency_id, entity_type, entity_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS marketplace_gap;
DROP TABLE IF EXISTS consumer_ingestion;
DROP TABLE IF EXISTS feed_health_snapshot;
DROP TABLE IF EXISTS validation_report;
DROP TABLE IF EXISTS incident;
DROP TABLE IF EXISTS vehicle_trip_assignment;
DROP TABLE IF EXISTS manual_override;
DROP TABLE IF EXISTS telemetry_event;
DROP TABLE IF EXISTS gtfs_draft_record;
DROP TABLE IF EXISTS gtfs_draft;
DROP TABLE IF EXISTS gtfs_frequency;
DROP TABLE IF EXISTS gtfs_shape_line;
DROP TABLE IF EXISTS gtfs_shape_point;
DROP TABLE IF EXISTS gtfs_stop_time;
DROP TABLE IF EXISTS gtfs_trip;
DROP TABLE IF EXISTS gtfs_calendar_date;
DROP TABLE IF EXISTS gtfs_calendar;
DROP TABLE IF EXISTS gtfs_stop;
DROP TABLE IF EXISTS gtfs_route;
DROP TABLE IF EXISTS published_feed;
DROP TABLE IF EXISTS feed_version;
DROP TABLE IF EXISTS feed_config;
DROP TABLE IF EXISTS device_credential;
DROP TABLE IF EXISTS role_binding;
DROP TABLE IF EXISTS agency_user;
DROP TABLE IF EXISTS agency;
DROP EXTENSION IF EXISTS postgis;
