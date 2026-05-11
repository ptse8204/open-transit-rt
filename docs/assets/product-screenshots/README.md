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
| `operations-console-overview.png` | `http://localhost:8080/admin/operations` | Operations Console overview, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed for no token, private URL, private email, private operator name, private IP, or secret. |
| `operations-readiness.png` | `http://localhost:8080/admin/operations/readiness` | Readiness review, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed for no retained evidence wording, private identifiers, or stronger claim. |
| `operations-checklist.png` | `http://localhost:8080/admin/operations/checklist` | Operator checklist, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed for no token, private URL, private device ID, private IP, or secret. |
| `operations-gtfs-quality.png` | `http://localhost:8080/admin/operations/gtfs-quality` | GTFS quality triage, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed that validation text is presented as diagnostics, not compliance proof. |
| `operations-validation-health.png` | `http://localhost:8080/admin/operations/validation-health` | Validator health, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed that validator status is not framed as consumer acceptance or compliance proof. |
| `public-feed-path-check.png` | local/demo terminal summary for five public feed fetches | Public feed path check, local/demo product screenshot. | not captured | Capture only a redacted command/output summary; avoid tokens, private URLs, private IPs, and private paths. |
| `telemetry-simulator.png` | `http://localhost:8080/admin/operations/telemetry-simulator` | Telemetry simulator guide, local/demo product screenshot. | captured | Captured from local demo seed on 2026-05-11; reviewed no real device ID, token, private URL, or secret is visible. |

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
