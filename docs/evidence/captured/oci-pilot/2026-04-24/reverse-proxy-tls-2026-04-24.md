# Hosted Reverse Proxy and TLS Evidence

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials
- Public host: `open-transit-pilot.duckdns.org`

## TLS / Redirect Artifacts

- HTTPS headers: `artifacts/tls/https-feeds-headers.txt`
- HTTP redirect headers: `artifacts/tls/http-redirect-headers.txt`
- Certificate details: `artifacts/tls/certificate.txt`
- Redacted Caddy route map and service status: `artifacts/operator-supplied/reverse-proxy-caddy-route-map-redacted.txt`

## Boundary Summary

Caddy exposes only anonymous public feed paths on the public host. Admin, debug, and studio routes are not routed through the public edge; operators use SSH tunneling for admin access, and the application still requires Bearer admin auth on the tunneled admin route.
