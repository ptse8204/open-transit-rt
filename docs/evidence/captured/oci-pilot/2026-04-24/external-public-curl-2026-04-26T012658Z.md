# External Public Curl Check

- Environment: `oci-pilot`
- Capture time (UTC): `2026-04-26T01:26:58Z`
- Network path: direct public HTTPS request from the operator workstation, with no SSH tunnel, no admin token, and no private origin URL.
- Endpoint: `https://open-transit-pilot.duckdns.org/public/feeds.json`

## Exact Command

```bash
curl --http1.1 --location --silent --show-error --dump-header docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/external-curl/feeds-json-headers-20260426T012658Z.txt --output docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/external-curl/feeds-json-body-20260426T012658Z.json https://open-transit-pilot.duckdns.org/public/feeds.json
```

## Header Evidence

Artifact: `artifacts/public/external-curl/feeds-json-headers-20260426T012658Z.txt`

- Status: `HTTP/1.1 200 OK`
- Date: `Sun, 26 Apr 2026 01:26:58 GMT`
- Cache-Control: not present
- ETag: not present
- Via: `1.1 Caddy`

## Body Evidence

Artifacts:

- Full body: `artifacts/public/external-curl/feeds-json-body-20260426T012658Z.json`
- Summary: `artifacts/public/external-curl/feeds-json-summary-20260426T012658Z.json`

The fetched body shows `canonical_validation_complete=true`. Schedule, Vehicle Positions, Trip Updates, and Alerts all report active feed version `gtfs-import-3` and `last_validation_status=passed`.

## Staleness Note

No stale response was observed in this external public curl. No mitigation was applied.
