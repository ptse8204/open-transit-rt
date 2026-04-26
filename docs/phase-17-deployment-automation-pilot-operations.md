# Phase 17 — Deployment Automation And Pilot Operations

## Status

Implemented for the Phase 17 deployment automation and pilot operations scope.

## Purpose

Phase 17 converts the successful OCI pilot evidence into repeatable deployment operations. Phase 12 proved one hosted pilot; Phase 17 makes that deployment path easier to reproduce and maintain without claiming universal production readiness, CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, or hosted SaaS availability.

## Scope

1. Deployment profiles and runbooks.
2. Reverse proxy and TLS automation.
3. Validator schedules.
4. Backup/restore automation.
5. Monitoring/alerting deployment assets.
6. Operational maintenance playbooks.
7. Evidence refresh and audit workflow.

## Required Work

### 1) Deployment Profile

Implemented in `docs/runbooks/small-agency-pilot-operations.md`.

The profile documents:

- an explicit deployment environment variable matrix;
- the systemd/Caddy pilot service layout;
- validator tooling and schedule;
- backup, restore, monitoring, and scorecard evidence paths;
- evidence naming conventions and commit-safety labels;
- how this differs from the Phase 16 local app package.

### 2) Reverse Proxy/TLS

Provided and documented through `deploy/oci/Caddyfile`, `docs/runbooks/reverse-proxy-and-tls.md`, and `docs/runbooks/small-agency-pilot-operations.md`.

The public edge exposes only stable public feed paths by default. Admin/debug/JSON/metrics surfaces remain internal, SSH-tunneled, or separately auth-protected. TLS renewal and HTTP to HTTPS redirect evidence remain deployment-owned and must be captured in evidence packets.

### 3) Validator Scheduling

Implemented through `scripts/pilot-ops.sh validator-cycle`, dry-run guidance, and `deploy/systemd/open-transit-validator-cycle.service` / `.timer` examples.

The workflow covers static GTFS, Vehicle Positions, Trip Updates, and Alerts. Outputs use `validator-cycle-YYYY-MM-DD.json` plus per-feed response files.

### 4) Backup/Restore Automation

Implemented through `scripts/pilot-ops.sh backup`, `scripts/pilot-ops.sh restore-drill`, dry-run guidance, and `deploy/systemd/open-transit-backup.service` / `.timer` examples.

Restore is explicitly destructive for `RESTORE_DATABASE_URL` and requires typed confirmation unless `--force` is passed. Backup outputs use `backup-run-YYYY-MM-DD.txt`; restore outputs use `restore-drill-YYYY-MM-DD.txt`.

### 5) Monitoring/Alerting Assets

Implemented through `scripts/pilot-ops.sh feed-monitor`, `deploy/systemd/open-transit-feed-monitor.service` / `.timer`, and monitoring runbook updates.

The helper checks public feed availability and records `feed-monitor-YYYY-MM-DD.txt`. Missing webhook/email destinations are reported as `notification not configured`, not as feed failures. Real webhook URLs and notification credentials must not be committed.

### 6) Evidence Refresh

Implemented in `docs/runbooks/deployment-evidence-overview.md`, `docs/evidence/README.md`, and `docs/runbooks/small-agency-pilot-operations.md`.

Every hosted evidence refresh must end with:

```sh
EVIDENCE_PACKET_DIR=<packet> make audit-hosted-evidence
```

Refreshed evidence is not complete unless this audit passes.

## Acceptance Criteria

Phase 17 is complete only when:

- a deployment operator has a repeatable pilot operations guide;
- proxy/TLS setup is documented safely;
- validator scheduling is documented or scripted;
- backup/restore automation is documented or scripted;
- monitoring/alerting examples exist;
- evidence refresh process is clear;
- no secrets are committed.

## Implemented Artifacts

- `docs/runbooks/small-agency-pilot-operations.md`
- `scripts/pilot-ops.sh`
- `deploy/systemd/open-transit-validator-cycle.service`
- `deploy/systemd/open-transit-validator-cycle.timer`
- `deploy/systemd/open-transit-backup.service`
- `deploy/systemd/open-transit-backup.timer`
- `deploy/systemd/open-transit-feed-monitor.service`
- `deploy/systemd/open-transit-feed-monitor.timer`
- `deploy/systemd/open-transit-scorecard-export.service`
- `deploy/systemd/open-transit-scorecard-export.timer`
- updated Phase 17 evidence, deployment, proxy/TLS, validator, backup/restore, monitoring, and scorecard runbooks.

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
