# Phase 39 — CAL-ITP-Style Readiness Workflow

## Status

Complete for the productized readiness-workflow scope.

Phase 39 makes CAL-ITP-style readiness gaps visible and actionable for
self-hosted agency operators. It does not claim CAL-ITP/Caltrans compliance,
consumer acceptance, agency adoption or approval, agency-owned final-root proof,
hosted SaaS availability, production readiness, vendor compatibility, or
production-grade ETA quality.

## What Changed

- Added an authenticated Operations Console page at
  `/admin/operations/readiness`.
- The readiness page shows ten operator-facing rows:
  - stable public URLs;
  - static GTFS feed;
  - Vehicle Positions;
  - Trip Updates;
  - Alerts;
  - license/contact metadata;
  - validation status;
  - telemetry freshness;
  - operations status;
  - consumer packet preparedness.
- Each row includes status, status source, current evidence/signal, next
  action, and claim boundary.
- Added the readiness page to the Operations Console navigation and dashboard.
- Updated docs and validation checks so the workflow is discoverable from
  onboarding, deployment, readiness, and handoff paths.

## Source Boundaries

The page uses existing local/deployment state only:

- `FeedDiscovery` and `published_feed` metadata;
- validation records;
- feed health snapshots;
- latest telemetry and assignment summaries;
- Trip Updates diagnostics;
- scorecard snapshots;
- runtime consumer workflow records;
- docs/evidence tracker paths.

Viewing the page creates no external evidence, contacts no consumers or
agencies, runs no validators, and changes no consumer statuses.

## Consumer Status

All seven consumer and aggregator targets remain `prepared` only. Prepared
packet records are not submitted, under review, accepted, listed, displayed, or
ingested.

Consumer statuses may move only when retained, redacted, target-originated
evidence supports a target-specific transition.
