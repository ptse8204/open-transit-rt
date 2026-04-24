# Phase 12 Hosted Evidence Packet: oci-pilot

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials
- Canonical HTTPS host: `https://open-transit-pilot.duckdns.org`
- Status: operator-reviewed hosted pilot packet

## Claim Boundary

This packet supports Phase 12 hosted/operator evidence for the OCI pilot. It does not claim Cal-ITP compliance or third-party consumer acceptance.

## Evidence Included

- Anonymous public HTTPS fetches for schedule, `feeds.json`, Vehicle Positions, Trip Updates, and Alerts.
- Public-edge auth boundary probes showing public `.pb` routes are anonymous 200 while debug/admin routes are absent from the public edge.
- SSH-tunneled admin auth probes showing Bearer auth is required for admin mutation.
- TLS certificate and HTTP to HTTPS redirect evidence.
- Hosted validator records for schedule, Vehicle Positions, Trip Updates, and Alerts, all passed after the data-restore rollback.
- Deployment data-restore rollback drill with pre-update, post-update, transient-update, and post-restore `feeds.json` snapshots.
- Operator-supplied reverse proxy, monitoring, alert lifecycle, backup, restore, and scorecard job artifacts.
