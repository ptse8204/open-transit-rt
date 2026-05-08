# Phase 37 — Reusable Agency Onboarding Flow

## Goal

Make it easy for an agency/operator to import and publish its own GTFS locally or on the reference server.

## Status

Closed for the reusable local/reference onboarding scope.

Phase 37 adds an executable onboarding helper and documentation for importing
an operator-provided GTFS ZIP by agency ID. It does not create final-root
evidence, external evidence packets, consumer artifacts, agency approval,
consumer acceptance, compliance, production-readiness, vendor-compatibility, or
ETA-quality claims.

## Command Shape

Local Compose:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

Direct script:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id agency \
  --gtfs-url https://example.org/gtfs.zip \
  --public-base-url http://localhost:8080
```

The script intentionally does not call `make agency-app-up` or
`scripts/agency-local-app.sh up`, because that path imports the demo sample
GTFS for `demo-agency`.

## Implemented Behavior

- validate required inputs;
- download GTFS into ignored storage;
- record checksum;
- create/seed local agency and admin roles safely without using
  `scripts/seed-dev.sql`;
- import with configurable timeout;
- start services or verify service readiness;
- fetch all five public paths;
- verify the fetched schedule is the imported GTFS summary, not an accidental
  leftover sample feed;
- run validators, skip validators, or document blockers;
- print next steps and support summary.

The helper accepts metadata flags and environment fallbacks:

```text
--technical-contact-email / TECHNICAL_CONTACT_EMAIL
--feed-license-name / FEED_LICENSE_NAME
--feed-license-url / FEED_LICENSE_URL
```

If not supplied, these use obvious local/reference placeholders. Script output
warns that publication metadata is placeholder metadata unless the operator
supplies agency-approved values.

## Must not do

- commit raw GTFS by default;
- create final-root evidence;
- claim agency approval;
- change consumer statuses.

## Follow-Up

If Phase 37 closes successfully, the next recommended productization phase is
Phase 38 — Integration Adapter Kit.
