# Phase 36 — OCI/OCL Reference Deployment Productization

## Status

Complete for docs-only OCI/OCL-style reference deployment productization.

Phase 36 turned the existing pilot-server pattern into a reusable self-hosted
reference deployment path. It created documentation only. It did not run a new
deployment, create external evidence, change runtime behavior, change schemas,
change migrations, change APIs, change consumer statuses, refresh final-root
proof, update OCI pilot evidence, or update validator artifacts.

## Deployment Docs

- [Deployment Reference Index](deployment/README.md)
- [OCI/OCL Reference Deployment](deployment/oci-reference-deployment.md)
- [Reference Environment Example](deployment/oci-reference-env.example)
- [Reference Smoke Checklist](deployment/oci-reference-smoke-checklist.md)

## Implemented Scope

The deployment guide covers:

- server prerequisites;
- user/group layout;
- directory layout under `/opt/open-transit-rt`;
- placeholder-only env files;
- secret generation;
- database setup/migration;
- service supervision;
- Caddy/reverse proxy public/private routing;
- validator install/check;
- GTFS import;
- five feed URL verification;
- backup/restore;
- feed monitor;
- scorecard export;
- update/rollback;
- support bundle/redaction rules.

The environment example uses placeholder-only values and grouped sections for
service runtime, public URL/feed metadata, admin/auth secrets, device token
secret, validator tooling, database, pilot-ops helpers, backup/restore, and
monitoring/notification placeholders.

## Goal State

A technically capable operator can reproduce the existing pilot-server style deployment from docs, without reading old evidence packets as instructions.

## Boundaries

This is reference deployment productization. It does not claim hosted SaaS,
agency adoption, agency endorsement, agency-owned final-root proof, consumer
submission/review/acceptance, consumer ingestion/listing/display,
CAL-ITP/Caltrans compliance, production readiness, vendor compatibility, or
production-grade ETA quality.

Each new deployment doc states:

> This guide is documentation for reproducing a self-hosted deployment pattern. It is not evidence that a deployment was run, accepted, compliant, agency-approved, or production-ready.

## Next Step

The next recommended phase is Phase 37 — Reusable Agency Onboarding Flow.
