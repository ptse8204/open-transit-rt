-- +goose Up
CREATE INDEX IF NOT EXISTS feed_version_active_schedule_idx
  ON feed_version (agency_id, is_active, activated_at DESC NULLS LAST, created_at DESC);

CREATE INDEX IF NOT EXISTS gtfs_route_schedule_export_idx
  ON gtfs_route (agency_id, feed_version_id, id);

CREATE INDEX IF NOT EXISTS gtfs_stop_schedule_export_idx
  ON gtfs_stop (agency_id, feed_version_id, id);

CREATE INDEX IF NOT EXISTS gtfs_calendar_schedule_export_idx
  ON gtfs_calendar (agency_id, feed_version_id, service_id);

CREATE INDEX IF NOT EXISTS gtfs_calendar_date_schedule_export_idx
  ON gtfs_calendar_date (agency_id, feed_version_id, service_id, date);

CREATE INDEX IF NOT EXISTS gtfs_trip_schedule_export_idx
  ON gtfs_trip (agency_id, feed_version_id, route_id, id);

CREATE INDEX IF NOT EXISTS gtfs_stop_time_schedule_export_idx
  ON gtfs_stop_time (agency_id, feed_version_id, trip_id, stop_sequence);

CREATE INDEX IF NOT EXISTS gtfs_shape_point_schedule_export_idx
  ON gtfs_shape_point (agency_id, feed_version_id, shape_id, sequence);

CREATE INDEX IF NOT EXISTS gtfs_frequency_schedule_export_idx
  ON gtfs_frequency (agency_id, feed_version_id, trip_id, start_time);

-- +goose Down
DROP INDEX IF EXISTS gtfs_frequency_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_shape_point_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_stop_time_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_trip_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_calendar_date_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_calendar_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_stop_schedule_export_idx;
DROP INDEX IF EXISTS gtfs_route_schedule_export_idx;
DROP INDEX IF EXISTS feed_version_active_schedule_idx;
