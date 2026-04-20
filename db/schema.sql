CREATE TABLE agency (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  timezone TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feed_version (
  id TEXT PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id),
  source_type TEXT NOT NULL,
  published_at TIMESTAMP,
  is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE route (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id),
  short_name TEXT,
  long_name TEXT,
  route_type INTEGER,
  PRIMARY KEY (id, feed_version_id)
);

CREATE TABLE stop (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id),
  name TEXT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  PRIMARY KEY (id, feed_version_id)
);

CREATE TABLE trip (
  id TEXT NOT NULL,
  feed_version_id TEXT NOT NULL REFERENCES feed_version(id),
  route_id TEXT NOT NULL,
  service_id TEXT NOT NULL,
  block_id TEXT,
  shape_id TEXT,
  direction_id INTEGER,
  PRIMARY KEY (id, feed_version_id)
);

CREATE TABLE telemetry_event (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id),
  device_id TEXT NOT NULL,
  vehicle_id TEXT NOT NULL,
  observed_at TIMESTAMP NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  bearing DOUBLE PRECISION,
  speed_mps DOUBLE PRECISION,
  accuracy_m DOUBLE PRECISION,
  payload_json TEXT
);

CREATE TABLE vehicle_trip_assignment (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id),
  vehicle_id TEXT NOT NULL,
  service_date TEXT NOT NULL,
  trip_id TEXT,
  start_date TEXT,
  start_time TEXT,
  state TEXT NOT NULL,
  confidence DOUBLE PRECISION NOT NULL,
  assignment_source TEXT NOT NULL,
  active_from TIMESTAMP NOT NULL,
  active_to TIMESTAMP
);
