-- +goose Up
CREATE TABLE gtfs_import (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  feed_version_id TEXT REFERENCES feed_version(id) ON DELETE SET NULL,
  source_filename TEXT NOT NULL,
  source_sha256 TEXT NOT NULL,
  source_byte_size BIGINT NOT NULL CHECK (source_byte_size >= 0),
  status TEXT NOT NULL CHECK (status IN ('started', 'failed', 'published')),
  error_count INTEGER NOT NULL DEFAULT 0,
  warning_count INTEGER NOT NULL DEFAULT 0,
  info_count INTEGER NOT NULL DEFAULT 0,
  report_json JSONB NOT NULL DEFAULT '{}'::jsonb,
  actor_id TEXT,
  notes TEXT,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE INDEX gtfs_import_agency_started_idx
  ON gtfs_import (agency_id, started_at DESC);

ALTER TABLE validation_report
  ADD COLUMN gtfs_import_id BIGINT REFERENCES gtfs_import(id) ON DELETE SET NULL;

CREATE INDEX validation_report_gtfs_import_idx
  ON validation_report (gtfs_import_id)
  WHERE gtfs_import_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS validation_report_gtfs_import_idx;

ALTER TABLE validation_report
  DROP COLUMN IF EXISTS gtfs_import_id;

DROP INDEX IF EXISTS gtfs_import_agency_started_idx;
DROP TABLE IF EXISTS gtfs_import;
