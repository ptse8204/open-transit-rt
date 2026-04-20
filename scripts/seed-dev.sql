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
