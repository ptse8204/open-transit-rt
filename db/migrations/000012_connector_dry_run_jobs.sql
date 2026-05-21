-- +goose Up
CREATE TABLE connector_dry_run_job (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  connector_instance_id BIGINT NOT NULL REFERENCES connector_instance(id) ON DELETE CASCADE,
  status TEXT NOT NULL CHECK (status IN ('passed', 'failed', 'blocked')),
  started_at TIMESTAMPTZ NOT NULL,
  finished_at TIMESTAMPTZ NOT NULL,
  redacted_summary JSONB NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(redacted_summary) = 'object'),
  accepted_count INTEGER NOT NULL DEFAULT 0 CHECK (accepted_count >= 0),
  rejected_count INTEGER NOT NULL DEFAULT 0 CHECK (rejected_count >= 0),
  dropped_count INTEGER NOT NULL DEFAULT 0 CHECK (dropped_count >= 0),
  redaction_scan_status TEXT NOT NULL CHECK (redaction_scan_status IN ('passed', 'failed', 'blocked')),
  created_by TEXT NOT NULL DEFAULT 'system',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX connector_dry_run_job_instance_idx
  ON connector_dry_run_job (agency_id, connector_instance_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS connector_dry_run_job;
