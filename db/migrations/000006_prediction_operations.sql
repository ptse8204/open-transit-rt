-- +goose Up
ALTER TABLE incident
  DROP CONSTRAINT IF EXISTS incident_status_check,
  ADD CONSTRAINT incident_status_check CHECK (status IN ('open', 'acknowledged', 'resolved', 'deferred'));

-- +goose Down
UPDATE incident
SET status = 'open'
WHERE status = 'deferred';

ALTER TABLE incident
  DROP CONSTRAINT IF EXISTS incident_status_check,
  ADD CONSTRAINT incident_status_check CHECK (status IN ('open', 'acknowledged', 'resolved'));
