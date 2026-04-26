# Runbook: Reverse Proxy and TLS Evidence

This runbook is for deployment/operator evidence capture. Do not use it to claim hosted proof until artifacts are actually collected.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/public-feed-proof-2026-04-22.md`
- `docs/evidence/captured/local-demo/2026-04-22/reverse-proxy-tls-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/reverse-proxy-tls-2026-04-24.md`

The local packet is loopback HTTP only. The OCI pilot packet is hosted deployment/operator proof for that recorded pilot scope.

## Purpose

Document and retain proof that public feeds are served from a stable HTTPS boundary with protected admin/debug routes.

## Inputs

- Deployment hostname and DNS ownership details.
- Reverse proxy config (redacted as needed).
- TLS certificate details and renewal method.
- Public/admin route map.

## Evidence Checklist

### A) Stable HTTPS Public Feed Root

Capture:

- Canonical host (for example `https://transit.example.gov`).
- Stable public paths:
  - `/public/gtfs/schedule.zip`
  - `/public/feeds.json`
  - `/public/gtfsrt/vehicle_positions.pb`
  - `/public/gtfsrt/trip_updates.pb`
  - `/public/gtfsrt/alerts.pb`
- Anonymous fetch proof for each path.
- Before/after publish proof showing URLs did not change.
- Rollback proof showing URLs still did not change.

### B) Routing and Boundary

Capture:

- Proxy routing map from public URLs to backend services.
- Explicit admin/debug boundary (`/admin/*`, `/admin/debug/*`, and protected debug JSON routes).
- HTTP→HTTPS redirect behavior if port 80 is exposed.

Default public Caddy behavior should expose only:

```text
/public/gtfs/*
/public/feeds.json
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

Admin, debug JSON, Studio, telemetry debug listings, and metrics should be absent from the public edge, SSH-tunneled, or protected by a separate auth-aware private edge.

### C) TLS Ownership and Renewal

Capture:

- TLS termination owner (load balancer, ingress, or edge proxy).
- Certificate issuer, validity window, and renewal strategy.
- Renewal check process and cadence.

## Output Artifact

Create one evidence file using:

- `docs/evidence/templates/public-feed-proof-template.md`

Store it in `docs/evidence/captured/<environment>/` with UTC date.

Evidence labels:

- Redacted route maps and TLS certificate metadata: `safe-to-commit-after-review`.
- Public feed headers and public redirect headers: `safe-to-commit-after-review`.
- TLS private keys, ACME account material, internal origin URLs, private hostnames, and raw access logs: `never-commit`.

## Failure Notes

If any public path is not anonymous, not HTTPS, or not stable across publish/rollback, record the mismatch explicitly and mark deployment evidence incomplete.
