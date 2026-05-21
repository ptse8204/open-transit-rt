-- +goose Up
CREATE TABLE connector_instance (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  connector_type TEXT NOT NULL CHECK (connector_type IN ('telemetry_source', 'prediction', 'validator', 'monitoring_export', 'consumer_discovery')),
  connector_kind TEXT NOT NULL CHECK (length(trim(connector_kind)) > 0),
  display_name TEXT NOT NULL CHECK (length(trim(display_name)) > 0),
  state TEXT NOT NULL DEFAULT 'configured_not_tested' CHECK (state IN ('example_available', 'not_configured', 'configured_not_tested', 'dry_run_passed', 'ready_for_activation', 'active', 'blocked')),
  owner TEXT NOT NULL DEFAULT '',
  config_json JSONB NOT NULL DEFAULT '{}'::jsonb CHECK (jsonb_typeof(config_json) = 'object'),
  secret_refs JSONB NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(secret_refs) = 'array'),
  dry_run_status TEXT NOT NULL DEFAULT 'not_run',
  last_checked_at TIMESTAMPTZ,
  activated_at TIMESTAMPTZ,
  disabled_at TIMESTAMPTZ,
  created_by TEXT NOT NULL DEFAULT 'system',
  updated_by TEXT NOT NULL DEFAULT 'system',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, connector_type, connector_kind, display_name)
);

CREATE INDEX connector_instance_agency_type_state_idx
  ON connector_instance (agency_id, connector_type, state);

CREATE INDEX connector_instance_last_checked_idx
  ON connector_instance (agency_id, last_checked_at DESC);

-- +goose Down
DROP TABLE IF EXISTS connector_instance;
