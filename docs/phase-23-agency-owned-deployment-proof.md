# Phase 23 — Agency-Owned Deployment Proof

## Status

Complete for blocker-documented closure only.

Phase 23 did not collect agency-owned or agency-approved final-root evidence
because no final public feed root was available or approved for this pass. The
DuckDNS OCI pilot remains useful hosted/operator evidence, but it is not
agency-owned production-domain proof.

## Purpose

Move from DuckDNS OCI pilot evidence toward agency-owned or agency-approved
stable URL proof.

The existing OCI pilot is useful hosted/operator evidence, but it is not
agency-owned production-domain proof. Stronger California-facing readiness
requires a stable public URL root controlled or approved by the
provider/operator.

## Phase 23 Result

Outcome B — blocker-documented closure.

No final-root evidence packet was created. No final-root validator records were
collected. Prepared consumer packets were not refreshed because final-root
fields are not evidence-backed.

| Area | Phase 23 result |
| --- | --- |
| Final public feed root | Not collected — final root unavailable. |
| Domain owner / approver | Not collected — final root unavailable. |
| DNS record proof | Not collected — final root unavailable. |
| TLS certificate metadata | Not collected — final root unavailable. |
| HTTP to HTTPS redirect behavior | Not collected — final root unavailable. |
| Operator approval for submissions | Not collected — final root unavailable. |
| All five final public feed URLs | Not collected — final root unavailable. |
| Final-root validator records | Not collected — final root unavailable. |
| Prepared packet refresh | Blocked — final-root evidence unavailable. |
| Migration / redirect proof | Blocked — no final root exists to migrate to. |

## Required Future Proof

Before claiming agency-owned or agency-approved production-domain readiness,
operators must collect or document:

- domain owner / approver;
- public feed root;
- DNS records;
- TLS certificate metadata;
- redirect behavior;
- operator approval for use in submissions;
- anonymous public fetch proof for all five feed URLs;
- canonical validator records for schedule, Vehicle Positions, Trip Updates,
  and Alerts at the final root;
- final-root packet refreshes backed by retained evidence.

## Final Public Feed URL Proof

The final root was unavailable, so no final-root fetch artifacts were collected
for:

- `/public/feeds.json`;
- `/public/gtfs/schedule.zip`;
- `/public/gtfsrt/vehicle_positions.pb`;
- `/public/gtfsrt/trip_updates.pb`;
- `/public/gtfsrt/alerts.pb`.

The existing OCI pilot proof remains under
`docs/evidence/captured/oci-pilot/2026-04-24/` and applies only to
`https://open-transit-pilot.duckdns.org`.

## Validator Records

Final-root validation did not run because no final root was available. Do not
reuse the DuckDNS OCI pilot validator records as agency-owned final-root
validator proof.

Future final-root evidence must retain validator records for:

- static schedule;
- Vehicle Positions;
- Trip Updates;
- Alerts.

## Packet Refresh

Prepared consumer packets were not refreshed in Phase 23. They remain pointed
at the OCI pilot and remain `prepared` only.

Refresh prepared packets only after final-root evidence exists for:

- final feed root;
- final feed URLs;
- final validator evidence;
- final license/contact metadata;
- final agency identity;
- final-root evidence packet reference.

Do not move any target beyond `prepared` without retained target-originated
evidence.

## Future Migration Plan Template

When a final agency-owned or agency-approved root exists, document:

- final root and exact five feed URLs;
- whether old DuckDNS URLs redirect to the final root;
- whether old DuckDNS URLs remain available without redirects;
- expected overlap period for old and new URLs;
- whether prepared packets need refresh;
- whether consumers or aggregators must be resubmitted;
- evidence packet paths for redirect, final-root fetch, TLS, and validation
  proof;
- how `docs/evidence/consumer-submissions/status.json` and the human tracker
  stay aligned if packet evidence references change.

No migration has occurred in Phase 23.

## Required Checks

For Phase 23 blocker-documented closure, run:

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Do not run `EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence` unless a
real final-root evidence packet is created. No such packet was created in Phase
23.

## Explicit Non-Goals

Phase 23 does not:

- submit to consumers by itself;
- claim acceptance;
- claim compliance from domain proof alone;
- implement hosted SaaS;
- change feed paths without a migration plan;
- commit private DNS credentials or TLS private keys;
- create fake final-root evidence;
- create a final-root evidence packet without real artifacts;
- refresh prepared consumer packets without final-root evidence.
