-- +goose Up
ALTER TABLE agency_user
  ADD COLUMN IF NOT EXISTS disabled_at TIMESTAMPTZ;

CREATE TABLE admin_password_credential (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  agency_user_id BIGINT NOT NULL REFERENCES agency_user(id) ON DELETE CASCADE,
  password_hash TEXT NOT NULL,
  password_hash_params JSONB NOT NULL DEFAULT '{}'::jsonb,
  status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'disabled', 'locked', 'reset_required')),
  failed_attempts INTEGER NOT NULL DEFAULT 0 CHECK (failed_attempts >= 0),
  locked_until TIMESTAMPTZ,
  password_changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (agency_id, agency_user_id)
);

CREATE INDEX admin_password_credential_user_idx
  ON admin_password_credential (agency_id, agency_user_id);

CREATE TABLE admin_bootstrap_token (
  id BIGSERIAL PRIMARY KEY,
  agency_id TEXT NOT NULL REFERENCES agency(id) ON DELETE CASCADE,
  agency_user_id BIGINT NOT NULL REFERENCES agency_user(id) ON DELETE CASCADE,
  purpose TEXT NOT NULL CHECK (purpose IN ('first_admin', 'password_reset')),
  token_hash TEXT NOT NULL UNIQUE,
  expires_at TIMESTAMPTZ NOT NULL,
  used_at TIMESTAMPTZ,
  created_by TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX admin_bootstrap_token_lookup_idx
  ON admin_bootstrap_token (token_hash, expires_at, used_at);

-- +goose Down
DROP TABLE IF EXISTS admin_bootstrap_token;
DROP TABLE IF EXISTS admin_password_credential;
ALTER TABLE agency_user
  DROP COLUMN IF EXISTS disabled_at;
