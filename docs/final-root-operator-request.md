# Final Public Feed Root Operator Request

This document is an operator-facing request package. It helps an agency, operator, or deployment owner understand what Open Transit RT needs before the project can collect agency-owned or agency-approved final public feed root evidence.

This document is **not** evidence by itself.

## Why this matters

Open Transit RT can publish stable public feed paths, and the repo has local/pilot public-GTFS evidence. However, stronger public readiness claims require a final public feed root that is agency-owned or agency-approved.

A final public feed root is the stable base URL where public consumers fetch the feed set.

Example:

```text
https://gtfs.exampleagency.gov
```

With the five feed URLs:

```text
https://gtfs.exampleagency.gov/public/feeds.json
https://gtfs.exampleagency.gov/public/gtfs/schedule.zip
https://gtfs.exampleagency.gov/public/gtfsrt/vehicle_positions.pb
https://gtfs.exampleagency.gov/public/gtfsrt/trip_updates.pb
https://gtfs.exampleagency.gov/public/gtfsrt/alerts.pb
```

## Acceptable root types

Preferred:

```text
https://gtfs.exampleagency.gov
https://data.exampleagency.gov/transit
https://transit.exampleagency.gov/feeds
```

Potentially acceptable if approved:

```text
https://feeds.example-operator.org/exampleagency
https://transit-data.example-operator.org/exampleagency
```

Not enough by itself:

```text
https://open-transit-pilot.duckdns.org/
```

The DuckDNS OCI pilot is useful pilot/operator evidence. It is not agency-owned or agency-approved final-root proof unless an authorized agency/operator explicitly approves it for final public feed use and that approval is retained.

## Approval evidence requested

Please provide one retained artifact showing the root is approved for public feed use.

Acceptable examples:

- email from an authorized agency/operator contact;
- ticket from agency IT or operator IT;
- signed pilot or deployment agreement naming the feed root;
- DNS/domain ownership screenshot or export, redacted as needed;
- written approval from the deployment owner naming the root and feed paths.

Do not send private credentials, portal passwords, DNS write tokens, private keys, ACME account keys, or raw private correspondence that cannot be redacted.

## Required technical proof

For the final root, collect public-safe evidence for:

| Area | Required proof |
| --- | --- |
| DNS | `A`, `AAAA`, or `CNAME` resolution for the final hostname. |
| TLS | Valid HTTPS certificate metadata for the final hostname. |
| Redirect | HTTP-to-HTTPS redirect proof if HTTP is exposed. |
| Public reachability | Anonymous fetch proof for all five public feed URLs. |
| Static GTFS | Current `schedule.zip` fetch proof and static GTFS validator record. |
| GTFS Realtime | Vehicle Positions, Trip Updates, and Alerts fetch proof and validator records. |
| Metadata | Public `feeds.json` with agency-approved provider, license, contact, and feed URL metadata. |
| Proxy/config | Redacted summary showing public feed paths are exposed and admin/debug/private routes are not exposed on the public edge. |
| Inventory | README, command inventory, timestamps, and checksums for retained artifacts. |

## Five public URLs to verify

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

All should be fetchable anonymously over HTTPS at the approved root.

## What not to send or commit

Do not send or commit:

- secrets;
- API tokens;
- generated device tokens;
- JWT secrets;
- CSRF secrets;
- database passwords;
- private keys;
- ACME account material;
- DNS provider write credentials;
- private portal credentials;
- raw telemetry payloads;
- private backup paths;
- unredacted logs with credentials;
- private ticket links that cannot be redacted;
- unredacted correspondence with personal contact details.

Use redacted summaries where possible.

## What this proof would allow

If retained final-root evidence is collected, the repo can truthfully say something like:

```text
A final public feed root has retained approval and evidence for DNS, TLS, anonymous public feed fetches, and validator records for the recorded deployment scope.
```

It still would not prove consumer acceptance, consumer ingestion, Caltrans/CAL-ITP compliance, production-grade ETA quality, or universal production readiness by itself.

## What happens after approval

After an approved final root is available:

1. deploy all five public feed URLs at the root;
2. collect DNS, TLS, redirect, fetch, validator, proxy/config, README, and checksum evidence;
3. run the hosted evidence audit if applicable;
4. refresh prepared consumer packets with final-root URLs only after the evidence exists;
5. submit to consumers only when an operator authorizes the target and the official path is verified;
6. update consumer status only from target-originated evidence.

## Claim boundary

This request package is only a request package. It does not prove:

- final-root approval;
- final-root deployment;
- agency adoption;
- agency endorsement;
- consumer submission;
- consumer acceptance;
- consumer ingestion/listing/display;
- Caltrans/CAL-ITP compliance;
- hosted SaaS availability;
- production readiness;
- production-grade ETA quality.
