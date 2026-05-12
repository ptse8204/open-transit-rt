# Product Screenshots

This directory is for local/demo product screenshots used as documentation
aids. Screenshots here are not retained evidence and must not be described as
proof of production readiness, CAL-ITP/Caltrans compliance, agency adoption or
approval, consumer submission or acceptance, final-root readiness, hosted SaaS
availability, SLA coverage, vendor compatibility, hardware certification, or
production-grade ETA quality.

Do not place screenshots or media under `docs/evidence`.

## Capture Rules

- Use only synthetic, local-demo, or public-safe data.
- Do not show credentials, tokens, private URLs, private device IDs, private
  agency/operator names, private emails, IPs, secrets, unredacted logs, or
  private infrastructure details.
- Do not use agency logos or real transit photos unless permission and license
  are explicit.
- Caption every referenced screenshot with the phrase
  "local/demo product screenshot".
- Keep commands copy-pasteable in docs; screenshots must not be the only place
  where a command appears.
- If a screenshot cannot be captured safely, leave a manifest row with
  `not captured` and do not fake an image.

## Manifest

| File | Source route or command | Caption | Status | Redaction review |
| --- | --- | --- | --- | --- |
| `operations-console-overview.png` | `http://localhost:8080/admin/operations` | Operations Console overview, local/demo product screenshot. | not captured | Capture only from local demo seed; verify no token, private URL, private email, private operator name, IP, or secret is visible. |
| `agency-operations-cockpit.png` | `http://localhost:8080/admin/operations` | Agency Operations Cockpit, local/demo product screenshot. | not captured | Capture setup progress and primary action cards only from local demo seed; verify every visible claim boundary remains negative/supporting. |
| `operations-feed-health.png` | `http://localhost:8080/admin/operations/feed-health` | Feed health command center, local/demo product screenshot. | not captured | Capture only from local demo seed; verify the five public paths are visible and status is not framed as consumer acceptance or compliance proof. |
| `operations-maintenance-center.png` | `http://localhost:8080/admin/operations/maintenance` | Maintenance Center, local/demo product screenshot. | not captured | Capture only local/demo rows; verify backup/restore values are configured/not configured only and no secret values are visible. |
| `operations-readiness.png` | `http://localhost:8080/admin/operations/readiness` | Readiness review, local/demo product screenshot. | not captured | Capture only from local demo seed; verify no retained evidence wording or private identifiers are visible. |
| `operations-checklist.png` | `http://localhost:8080/admin/operations/checklist` | Operator checklist, local/demo product screenshot. | not captured | Capture only from local demo seed; verify no token, private URL, private device ID, IP, or secret is visible. |
| `operations-gtfs-quality.png` | `http://localhost:8080/admin/operations/gtfs-quality` | GTFS quality triage, local/demo product screenshot. | not captured | Capture only from local demo seed; verify validation text is presented as diagnostics, not compliance proof. |
| `operations-validation-health.png` | `http://localhost:8080/admin/operations/validation-health` | Validator health, local/demo product screenshot. | not captured | Capture only from local demo seed; verify validator status is not framed as consumer acceptance or compliance proof. |
| `public-feed-path-check.png` | local/demo terminal summary for five public feed fetches | Public feed path check, local/demo product screenshot. | not captured | Capture only a redacted command/output summary; avoid tokens, private URLs, private IPs, and private paths. |
| `telemetry-simulator.png` | `http://localhost:8080/admin/operations/telemetry-simulator` | Telemetry simulator guide, local/demo product screenshot. | not captured | Capture only committed synthetic scenario guidance; verify no real device ID, token, private URL, or secret is visible. |
| `connector-hub.png` | `http://localhost:8080/admin/operations/connectors` | Connector Hub maturity review, local/demo product screenshot. | not captured | Capture only committed synthetic manifests and safe connector wording; verify no vendor compatibility, consumer status, or production-readiness claim is visible. |

## Capture Checklist

1. Start the local demo package:

   ```bash
   make agency-app-up
   ```

2. Generate an admin token without writing it into docs:

   ```bash
   docker compose -f deploy/docker-compose.yml --profile app exec -T agency-config \
     /app/bin/admin-token -sub admin@example.com -agency-id demo-agency
   ```

3. Capture only local/demo routes with an `Authorization: Bearer ...` header.
4. Review each image before committing it.
5. Update the manifest status to `captured` only for files that exist and pass
   redaction review.
6. Stop the local app:

   ```bash
   make agency-app-down
   ```

## What This Proves

These screenshots can help readers recognize local/demo UI surfaces and follow
documentation faster.

## What This Does Not Prove

These screenshots do not prove production readiness, compliance, adoption,
approval, consumer status, final-root readiness, hosted service availability,
SLA coverage, vendor compatibility, hardware certification, or production-grade
ETA quality.
