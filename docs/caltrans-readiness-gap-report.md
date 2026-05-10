# Caltrans Readiness Gap Report

## Purpose

`make caltrans-readiness-check` writes a private `.cache` gap report for
operators and maintainers who want to review Caltrans-style transit data
readiness signals without making unsupported claims.

This is a gap-review workflow only. It does not create retained evidence,
refresh official requirements, fetch official pages, contact consumers, automate
portals, change consumer statuses, or prove CAL-ITP/Caltrans compliance.
It does not claim CAL-ITP/Caltrans compliance for the repository or any
deployment.

## Official-Source Context

The report uses the repository's Phase 54 official-source context by date:

- Phase 54 refresh date: May 9, 2026.
- Caltrans California Transit Data Guidelines: Version 4.0, dated December 11,
  2024, as recorded during Phase 54.
- Caltrans California Transit Data Guidelines FAQ: Version 4.0, as recorded
  during Phase 54.
- FTA 2025 NTD Reporting Policy Manual page: last updated April 15, 2026, as
  recorded during Phase 54.

The check does not verify whether those official pages have changed since Phase
54. Before any public compliance claim or final compliance closeout, a
maintainer must approve a separate official requirements refresh.

## What The Check Reviews

The generated report has rows for:

- stable feed URLs;
- static GTFS;
- GTFS Realtime Vehicle Positions;
- GTFS Realtime Trip Updates;
- GTFS Realtime Alerts;
- public fetchability;
- HTTPS;
- open license metadata;
- technical/feed contact metadata;
- validator signals;
- freshness signals;
- trip ID consistency signals;
- consumer packet preparedness;
- unsupported-claim boundaries.

Allowed row statuses are only:

- `present`
- `partial`
- `missing`
- `not_checked`
- `needs_review`
- `blocked`

The report intentionally does not use statuses such as `ok`, `passed`,
`compliant`, `certified`, `accepted`, `ingested`, `listed`, or `displayed`.
Those words can be mistaken for claim outcomes.

## Usage

Run the default cache-only review:

```sh
make caltrans-readiness-check
```

Run the script directly:

```sh
scripts/caltrans-readiness-check.sh
```

Fast dry-run:

```sh
scripts/caltrans-readiness-check.sh --dry-run
```

Optional inputs:

```sh
PUBLIC_BASE_URL=https://feeds.example.org \
FEEDS_JSON_PATH=.cache/operator/feeds.json \
VALIDATOR_HEALTH_SUMMARY=.cache/validator-health/<timestamp>/summary.json \
OPERATIONS_RELIABILITY_SUMMARY=.cache/operations-reliability/<timestamp>/summary.json \
TRIP_ID_CONSISTENCY_SUMMARY=.cache/trip-id-consistency/<timestamp>/summary.json \
make caltrans-readiness-check
```

Live public fetch checks are disabled by default. To perform bounded GET/HEAD
fetch checks against configured feed URLs, opt in explicitly:

```sh
PUBLIC_BASE_URL=https://feeds.example.org \
RUN_PUBLIC_FETCH=true \
FETCH_TIMEOUT_SECONDS=5 \
make caltrans-readiness-check
```

Use opt-in fetches only for a root the operator is allowed to check. Missing
or unavailable fetch inputs remain `missing`, `not_checked`, `needs_review`, or
`blocked`; they must not be converted silently to `present`.

## Outputs

The script writes exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `gap-review.txt`

The default output root is:

```text
.cache/caltrans-readiness-check/<timestamp>
```

These files are private local diagnostics. They are not retained evidence and
must not be copied under `docs/evidence` without a separately approved evidence
intake.

## Claim Boundary

The report keeps all claim flags false. It does not prove:

- CAL-ITP/Caltrans compliance;
- consumer submission, review, acceptance, ingestion, listing, or display;
- consumer acceptance;
- agency adoption, endorsement, or approval;
- agency-owned final-root readiness;
- public launch;
- hosted SaaS or paid support availability;
- SLA or uptime coverage;
- production readiness;
- production multi-tenant hosting;
- vendor compatibility, production AVL reliability, or certified hardware
  support;
- production-grade ETA quality.

Validator output, public fetchability, HTTPS, license/contact metadata, and
consumer packet preparedness are supporting signals only. They are not
compliance, consumer acceptance, or correctness proof.
