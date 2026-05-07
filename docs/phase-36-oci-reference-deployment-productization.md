# Phase 36 — OCI/OCL Reference Deployment Productization

## Goal

Turn the existing OCI/OCL pilot server pattern into the reference deployment that agencies and operators can reuse.

## Required work

Create or update:

```text
docs/deployment/oci-reference-deployment.md
docs/deployment/oci-reference-env.example
docs/deployment/oci-reference-smoke-checklist.md
```

## Must cover

- server prerequisites;
- user/group layout;
- directory layout under `/opt/open-transit-rt`;
- env file placeholders;
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

## Goal state

A technically capable operator can reproduce the existing pilot-server style deployment from docs, without reading old evidence packets as instructions.

## Boundaries

This is reference deployment productization. It does not claim hosted SaaS, final-root proof, consumer acceptance, or production readiness for all agencies.
