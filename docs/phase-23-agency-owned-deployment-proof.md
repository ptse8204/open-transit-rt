# Phase 23 — Agency-Owned Deployment Proof

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Move from DuckDNS OCI pilot evidence toward agency-owned or agency-approved stable URL proof.

The existing OCI pilot is useful hosted/operator evidence, but it is not agency-owned production-domain proof. Stronger California-facing readiness requires a stable public URL root controlled or approved by the provider/operator.

## Scope

1. Agency-owned or agency-approved domain checklist.
2. Final public feed URL proof.
3. TLS and redirect proof for the final root.
4. Validator records for the final root.
5. Updated prepared packets for the final root.
6. Redirect/migration plan from pilot URL, if needed.
7. Evidence packet for the final root.

## Required Work

### 1) Domain Readiness

Document or collect:

- domain owner / approver;
- public feed root;
- DNS records;
- TLS certificate metadata;
- redirect behavior;
- operator approval for use in submissions.

### 2) Public Feed Proof

Capture final-root fetch evidence for:

- `/public/feeds.json`;
- `/public/gtfs/schedule.zip`;
- `/public/gtfsrt/vehicle_positions.pb`;
- `/public/gtfsrt/trip_updates.pb`;
- `/public/gtfsrt/alerts.pb`.

### 3) Validator Records

Run and retain validation evidence for the final URL root:

- static schedule;
- Vehicle Positions;
- Trip Updates;
- Alerts.

### 4) Packet Refresh

Refresh prepared consumer packets only when final-root fields are known:

- final feed root;
- final feed URLs;
- final validator evidence;
- final license/contact metadata;
- final agency identity.

### 5) Migration Plan

If a pilot URL was previously shared, document:

- whether old URLs redirect;
- whether consumers must be resubmitted;
- how long old URLs remain available;
- how packets and evidence records are updated.

## Acceptance Criteria

Phase 23 is complete only when:

- agency-owned or agency-approved URL root is documented or blockers are explicit;
- all five public feed URLs are proven at the final root, or blockers are explicit;
- validator records for final root exist, or blockers are explicit;
- prepared packets are refreshed only with evidence;
- DuckDNS pilot remains labeled as pilot evidence;
- no compliance or consumer-acceptance claim is introduced.

## Required Checks

```bash
make validate
make test
git diff --check
```

If hosted evidence is collected:

```bash
EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence
```

Run additional deployment checks as appropriate.

## Explicit Non-Goals

Phase 23 does not:

- submit to consumers by itself;
- claim acceptance;
- claim compliance from domain proof alone;
- implement hosted SaaS;
- change feed paths without a migration plan;
- commit private DNS credentials or TLS private keys.

## Likely Files

- `docs/agency-owned-domain-readiness.md`
- `docs/evidence/captured/<environment>/<date>/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/california-readiness-summary.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-23.md`
