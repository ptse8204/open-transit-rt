-- +goose Up
ALTER TABLE gtfs_draft
  ADD COLUMN last_published_feed_version_id TEXT REFERENCES feed_version(id),
  ADD COLUMN last_publish_attempt_id BIGINT,
  ADD COLUMN discarded_at TIMESTAMPTZ,
  ADD COLUMN discarded_by TEXT,
  ADD COLUMN discard_reason TEXT;

ALTER TABLE gtfs_draft
  ADD CONSTRAINT gtfs_draft_id_agency_unique UNIQUE (id, agency_id);

CREATE TABLE gtfs_draft_agency (
  draft_id TEXT PRIMARY KEY,
  agency_id TEXT NOT NULL,
  name TEXT NOT NULL,
  timezone TEXT NOT NULL,
  contact_email TEXT,
  public_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft_route (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  id TEXT NOT NULL,
  short_name TEXT,
  long_name TEXT,
  route_type INTEGER NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, id),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft_stop (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  id TEXT NOT NULL,
  name TEXT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, id),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE,
  CHECK (lat >= -90 AND lat <= 90),
  CHECK (lon >= -180 AND lon <= 180)
);

CREATE TABLE gtfs_draft_calendar (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  service_id TEXT NOT NULL,
  monday BOOLEAN NOT NULL,
  tuesday BOOLEAN NOT NULL,
  wednesday BOOLEAN NOT NULL,
  thursday BOOLEAN NOT NULL,
  friday BOOLEAN NOT NULL,
  saturday BOOLEAN NOT NULL,
  sunday BOOLEAN NOT NULL,
  start_date TEXT NOT NULL,
  end_date TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, service_id),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft_calendar_date (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  service_id TEXT NOT NULL,
  date TEXT NOT NULL,
  exception_type INTEGER NOT NULL CHECK (exception_type IN (1, 2)),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, service_id, date),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft_trip (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  id TEXT NOT NULL,
  route_id TEXT NOT NULL,
  service_id TEXT NOT NULL,
  block_id TEXT,
  shape_id TEXT,
  direction_id INTEGER CHECK (direction_id IN (0, 1)),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, id),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE,
  FOREIGN KEY (draft_id, route_id) REFERENCES gtfs_draft_route(draft_id, id)
);

CREATE INDEX gtfs_draft_trip_route_idx ON gtfs_draft_trip (draft_id, route_id);

CREATE TABLE gtfs_draft_stop_time (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  trip_id TEXT NOT NULL,
  arrival_time TEXT,
  departure_time TEXT,
  stop_id TEXT NOT NULL,
  stop_sequence INTEGER NOT NULL,
  pickup_type INTEGER,
  drop_off_type INTEGER,
  shape_dist_traveled DOUBLE PRECISION,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, trip_id, stop_sequence),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE,
  FOREIGN KEY (draft_id, trip_id) REFERENCES gtfs_draft_trip(draft_id, id) ON DELETE CASCADE,
  FOREIGN KEY (draft_id, stop_id) REFERENCES gtfs_draft_stop(draft_id, id)
);

CREATE TABLE gtfs_draft_shape_point (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  shape_id TEXT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  sequence INTEGER NOT NULL,
  dist_traveled DOUBLE PRECISION,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, shape_id, sequence),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE,
  CHECK (lat >= -90 AND lat <= 90),
  CHECK (lon >= -180 AND lon <= 180)
);

CREATE TABLE gtfs_draft_frequency (
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  trip_id TEXT NOT NULL,
  start_time TEXT NOT NULL,
  end_time TEXT NOT NULL,
  headway_secs INTEGER NOT NULL CHECK (headway_secs > 0),
  exact_times INTEGER NOT NULL DEFAULT 0 CHECK (exact_times IN (0, 1)),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (draft_id, trip_id, start_time),
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE,
  FOREIGN KEY (draft_id, trip_id) REFERENCES gtfs_draft_trip(draft_id, id) ON DELETE CASCADE
);

CREATE TABLE gtfs_draft_publish (
  id BIGSERIAL PRIMARY KEY,
  draft_id TEXT NOT NULL,
  agency_id TEXT NOT NULL,
  feed_version_id TEXT REFERENCES feed_version(id) ON DELETE SET NULL,
  status TEXT NOT NULL CHECK (status IN ('started', 'failed', 'published')),
  error_count INTEGER NOT NULL DEFAULT 0,
  warning_count INTEGER NOT NULL DEFAULT 0,
  info_count INTEGER NOT NULL DEFAULT 0,
  report_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  actor_id TEXT,
  notes TEXT,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ,
  FOREIGN KEY (draft_id, agency_id) REFERENCES gtfs_draft(id, agency_id) ON DELETE CASCADE
);

CREATE INDEX gtfs_draft_publish_draft_started_idx
  ON gtfs_draft_publish (draft_id, started_at DESC);

ALTER TABLE gtfs_draft
  ADD CONSTRAINT gtfs_draft_last_publish_attempt_fk
  FOREIGN KEY (last_publish_attempt_id) REFERENCES gtfs_draft_publish(id) ON DELETE SET NULL;

ALTER TABLE validation_report
  ADD COLUMN gtfs_draft_publish_id BIGINT REFERENCES gtfs_draft_publish(id) ON DELETE SET NULL;

CREATE INDEX validation_report_gtfs_draft_publish_idx
  ON validation_report (gtfs_draft_publish_id)
  WHERE gtfs_draft_publish_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS validation_report_gtfs_draft_publish_idx;

ALTER TABLE validation_report
  DROP COLUMN IF EXISTS gtfs_draft_publish_id;

ALTER TABLE gtfs_draft
  DROP CONSTRAINT IF EXISTS gtfs_draft_last_publish_attempt_fk;

DROP INDEX IF EXISTS gtfs_draft_publish_draft_started_idx;
DROP TABLE IF EXISTS gtfs_draft_publish;
DROP TABLE IF EXISTS gtfs_draft_frequency;
DROP TABLE IF EXISTS gtfs_draft_shape_point;
DROP TABLE IF EXISTS gtfs_draft_stop_time;
DROP TABLE IF EXISTS gtfs_draft_trip;
DROP INDEX IF EXISTS gtfs_draft_trip_route_idx;
DROP TABLE IF EXISTS gtfs_draft_calendar_date;
DROP TABLE IF EXISTS gtfs_draft_calendar;
DROP TABLE IF EXISTS gtfs_draft_stop;
DROP TABLE IF EXISTS gtfs_draft_route;
DROP TABLE IF EXISTS gtfs_draft_agency;

ALTER TABLE gtfs_draft
  DROP CONSTRAINT IF EXISTS gtfs_draft_id_agency_unique;

ALTER TABLE gtfs_draft
  DROP COLUMN IF EXISTS discard_reason,
  DROP COLUMN IF EXISTS discarded_by,
  DROP COLUMN IF EXISTS discarded_at,
  DROP COLUMN IF EXISTS last_publish_attempt_id,
  DROP COLUMN IF EXISTS last_published_feed_version_id;
