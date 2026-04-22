# Public Feed Root and TLS Proof Template

- Environment:
- Capture date (UTC):
- Operator:
- Canonical HTTPS host:

## Stable Public Paths

Record one anonymous HTTPS fetch for each path:

- [ ] `/public/gtfs/schedule.zip`
- [ ] `/public/feeds.json`
- [ ] `/public/gtfsrt/vehicle_positions.pb`
- [ ] `/public/gtfsrt/trip_updates.pb`
- [ ] `/public/gtfsrt/alerts.pb`

For each, include status code, key headers, and command transcript.

## Publish / Rollback URL Stability

- Before publish proof:
- After publish proof:
- After rollback proof:
- URL changed? (must be No):

## Reverse Proxy / TLS

- Routing map reference:
- TLS termination owner:
- Certificate issuer + validity:
- Renewal check process:
- HTTP→HTTPS behavior:

## Admin/Debug Protection Boundary

- `/admin/*` protection evidence:
- `/admin/debug/*` protection evidence:
- `/public/gtfsrt/*.json` protection evidence:

## Notes / Gaps

- Pending items:
- Known blockers:
