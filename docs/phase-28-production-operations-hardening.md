# Phase 28 — Production Operations Hardening

## Status

Planned Track B phase. Not implemented until selected in `docs/handoffs/latest.md`.

## Purpose

Move from pilot operations examples toward stronger operational maturity for real agency deployments.

Phase 17 added repeatable pilot operations runbooks and helper scripts. Phase 28 should harden long-running operations, secret rotation, upgrade practices, incident response, and evidence refresh cadence.

## Scope

1. Long-running monitoring runbooks.
2. Alert delivery proof pattern.
3. Backup/restore cadence.
4. Upgrade and migration operations.
5. Incident response templates.
6. Validator failure response.
7. Capacity guidance.
8. Secret rotation runbook.
9. Operator handover checklist.

## Required Work

### 1) Monitoring And Alerting

Document or improve:

- feed availability checks;
- freshness checks;
- service readiness checks;
- validation failure alerts;
- disk/database capacity alerts;
- notification destination handling;
- alert lifecycle evidence.

### 2) Backup And Restore Cadence

Define:

- backup frequency;
- retention;
- restore drill schedule;
- restore verification;
- evidence capture;
- private dump handling.

### 3) Upgrade Operations

Document:

- pre-upgrade backup;
- migration status check;
- post-upgrade validation;
- rollback/restore decision path;
- release note review.

### 4) Secret Rotation

Create a runbook for rotating:

- admin JWT secret;
- CSRF secret;
- device token pepper;
- device tokens;
- database password;
- TLS/ACME material;
- notification credentials.

### 5) Incident Response

Templates for:

- public feed outage;
- validation failure;
- telemetry outage;
- stale Trip Updates;
- secret exposure;
- consumer complaint or rejection;
- rollback/restore event.

## Acceptance Criteria

Phase 28 is complete only when:

- production operations runbooks are stronger and actionable;
- secret rotation is documented;
- incident templates exist;
- backup/restore cadence is clear;
- upgrade operations are tied to release process;
- no paid SLA or universal readiness is implied.

## Required Checks

```bash
make validate
make test
git diff --check
```

If operations scripts change:

```bash
sh -n <changed scripts>
make smoke
docker compose -f deploy/docker-compose.yml config
```

## Explicit Non-Goals

Phase 28 does not:

- create paid support;
- promise SLA coverage;
- implement a full observability stack unless explicitly approved;
- claim production readiness for all agencies;
- change backend product behavior without a separate phase.

## Likely Files

- `docs/runbooks/`
- `docs/support-boundaries.md`
- `docs/release-process.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-28.md`
- `scripts/pilot-ops.sh` only if safe improvements are needed
