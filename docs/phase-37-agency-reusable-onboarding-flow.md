# Phase 37 — Reusable Agency Onboarding Flow

## Goal

Make it easy for an agency/operator to import and publish its own GTFS locally or on the reference server.

## Desired command shape

One of these patterns:

```bash
make agency-pilot-up AGENCY_ID=agency GTFS_URL=https://example.org/gtfs.zip
```

or:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id agency \
  --gtfs-url https://example.org/gtfs.zip \
  --public-base-url http://localhost:8080
```

## Required behavior

- validate required inputs;
- download GTFS into ignored storage;
- record checksum;
- create/seed local agency safely;
- import with configurable timeout;
- start services or verify service readiness;
- fetch all five public paths;
- run validators or document blockers;
- print next steps and support summary.

## Must not do

- commit raw GTFS by default;
- create final-root evidence;
- claim agency approval;
- change consumer statuses.
