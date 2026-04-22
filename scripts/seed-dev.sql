INSERT INTO agency (id, name, timezone, contact_email, public_url)
VALUES
  ('demo-agency', 'Demo Agency', 'America/Vancouver', 'dev@example.com', 'http://localhost'),
  ('overnight-agency', 'Overnight Agency', 'America/Vancouver', 'dev@example.com', 'http://localhost'),
  ('freq-agency', 'Frequency Agency', 'America/Vancouver', 'dev@example.com', 'http://localhost')
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    timezone = EXCLUDED.timezone,
    contact_email = EXCLUDED.contact_email,
    public_url = EXCLUDED.public_url;

WITH upserted AS (
  INSERT INTO agency_user (agency_id, email, display_name, auth_subject)
  VALUES ('demo-agency', 'admin@example.com', 'Local Admin', 'admin@example.com')
  ON CONFLICT (agency_id, email) DO UPDATE
  SET display_name = EXCLUDED.display_name,
      auth_subject = EXCLUDED.auth_subject
  RETURNING id, agency_id
)
INSERT INTO role_binding (agency_id, agency_user_id, role)
SELECT agency_id, id, role
FROM upserted
CROSS JOIN (VALUES ('admin'), ('editor'), ('operator'), ('read_only')) AS roles(role)
ON CONFLICT (agency_id, agency_user_id, role) DO NOTHING;

INSERT INTO device_credential (agency_id, device_id, vehicle_id, token_hash, status)
VALUES (
  'demo-agency',
  'device-1',
  'bus-1',
  'hmac-sha256:c965e4c7dc64fbe7790da93b59525a896c37065e5be31e01e792378ddd16ff21',
  'active'
)
ON CONFLICT (agency_id, device_id) DO UPDATE
SET vehicle_id = EXCLUDED.vehicle_id,
    token_hash = EXCLUDED.token_hash,
    status = EXCLUDED.status,
    revoked_at = NULL,
    rotated_at = NULL;
