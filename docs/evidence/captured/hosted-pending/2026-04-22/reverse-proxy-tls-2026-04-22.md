# Hosted Reverse Proxy and TLS Evidence

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Status: missing

## Required Evidence

- Reverse proxy or load balancer routing map.
- TLS certificate issuer, subject, SANs, and validity window.
- Renewal mechanism and last renewal/check timestamp.
- HTTP-to-HTTPS redirect behavior if HTTP is exposed.
- Admin/debug protection boundary.

## Collection Commands

```sh
mkdir -p "$ENVIRONMENT_NAME/tls"

curl -sS -I "$PUBLIC_BASE_URL/public/feeds.json" \
  | tee "$ENVIRONMENT_NAME/tls/https-feeds-headers.txt"

curl -sS -I "http://$(printf '%s' "$PUBLIC_BASE_URL" | sed 's#^https://##')/public/feeds.json" \
  | tee "$ENVIRONMENT_NAME/tls/http-redirect-headers.txt"

openssl s_client -connect "$(printf '%s' "$PUBLIC_BASE_URL" | sed 's#^https://##'):443" \
  -servername "$(printf '%s' "$PUBLIC_BASE_URL" | sed 's#^https://##')" \
  </dev/null 2>/dev/null \
  | openssl x509 -noout -issuer -subject -dates -ext subjectAltName \
  | tee "$ENVIRONMENT_NAME/tls/certificate.txt"
```

Attach a redacted proxy config snippet or load balancer route export as:

- `reverse-proxy-config-redacted.txt`
- `route-map.txt`
- `renewal-evidence.txt`

## Required Summary To Fill

- Public host:
- Proxy/load balancer product:
- Public path routing map:
- TLS issuer:
- TLS not before:
- TLS not after:
- Renewal mechanism:
- Last renewal/check evidence:
- HTTP-to-HTTPS result:
- Admin/debug anonymous result:

## Blocker

No hosted proxy, TLS certificate, redirect, or renewal artifacts were available when this intake packet was created.
