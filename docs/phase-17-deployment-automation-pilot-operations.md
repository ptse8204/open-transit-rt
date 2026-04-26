# Phase 17 — Deployment Automation And Pilot Operations

## Status

Planned phase. Not implemented until `docs/handoffs/latest.md` marks it active.

## Purpose

Phase 17 converts the successful OCI pilot evidence into repeatable deployment operations. Phase 12 proved one hosted pilot; Phase 17 makes that deployment path easier to reproduce and maintain.

## Scope

1. Deployment profiles and runbooks.
2. Reverse proxy and TLS automation.
3. Validator schedules.
4. Backup/restore automation.
5. Monitoring/alerting deployment assets.
6. Operational maintenance playbooks.

## Required Work

### 1) Deployment Profile

Document and/or implement a reusable deployment profile for a small agency pilot:

- environment variables;
- service units or containers;
- reverse proxy config;
- validator tooling;
- backup location;
- monitoring location;
- evidence export path.

### 2) Reverse Proxy/TLS

Provide a reusable, redacted Caddy or Nginx example that:

- exposes only public feed paths;
- keeps admin/debug private or auth protected;
- handles TLS renewal;
- documents HTTP→HTTPS redirect behavior.

### 3) Validator Scheduling

Add a recommended scheduled validation workflow for:

- static GTFS;
- Vehicle Positions;
- Trip Updates;
- Alerts.

### 4) Backup/Restore Automation

Document and/or script:

- scheduled backups;
- retention cleanup;
- restore drill command sequence;
- operator verification steps.

### 5) Monitoring/Alerting Assets

Add example monitoring configuration where practical:

- feed availability checks;
- feed freshness checks;
- validator failure alerts;
- service readiness alerts;
- notification path placeholder.

### 6) Evidence Refresh

Add a documented process for refreshing deployment evidence periodically without overstating it as compliance.

## Acceptance Criteria

Phase 17 is complete only when:

- a deployment operator has a repeatable pilot operations guide;
- proxy/TLS setup is documented safely;
- validator scheduling is documented or scripted;
- backup/restore automation is documented or scripted;
- monitoring/alerting examples exist;
- evidence refresh process is clear;
- no secrets are committed.

## Required Checks

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Run smoke/demo checks if deployment scripts touch local flows.

## Explicit Non-Goals

Phase 17 does not:

- claim official production readiness for every agency;
- claim consumer acceptance;
- require a specific cloud vendor;
- introduce complex Kubernetes requirements unless explicitly approved.
