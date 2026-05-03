# Phase 28 — Production Operations Hardening

## Status

Complete for the docs-first operations hardening scope.

Phase 28 added runbooks and templates for repeatable small-agency pilot operations. It did not change runtime APIs, database schema, public feed URLs, GTFS-RT contracts, consumer statuses, external integrations, systemd/Docker behavior, or evidence claims.

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
- alert lifecycle evidence;
- alert delivery proof pattern without requiring hosted monitoring SaaS or a full Prometheus/Grafana stack.

### 2) Backup And Restore Cadence

Define:

- backup frequency;
- retention;
- restore drill schedule;
- restore verification;
- evidence capture;
- private dump handling;
- deployment/DB-scoped backup and restore boundary from Phase 27.

### 3) Upgrade Operations

Document:

- pre-upgrade backup;
- migration status check;
- post-upgrade validation;
- rollback/restore decision path;
- release note review.
- evidence packet version linkage;
- irreversible or untested migration handling.

### 4) Secret Rotation

Create a runbook for rotating:

- admin JWT secret;
- CSRF secret;
- device token pepper;
- device tokens;
- database password;
- TLS/ACME material;
- optional webhook/notification credentials;
- Phase 15 `.cache` secret findings.

Deleting a file is not enough when a real secret was exposed. Operators must rotate or revoke the credential, verify the old value no longer works, and assess history/backups.

### 5) Incident Response

Templates for:

- public feed outage;
- validation failure;
- telemetry outage;
- stale Trip Updates;
- secret exposure;
- consumer complaint or rejection;
- rollback/restore event.

Each template must include start time, affected environment, affected agency, affected public URLs or services, detection source, operator, severity, timeline, action taken, evidence retained, redaction review, follow-up, and claim boundary.

### 6) Operator Handover

The operator handover template includes:

- current release version;
- deployment environment;
- public feed URLs;
- admin access process without secrets;
- secret storage location without secrets;
- backup location;
- restore process;
- validator cadence;
- monitoring cadence;
- evidence packet location;
- known blockers;
- consumer status boundaries;
- agency-owned-domain status;
- multi-agency limitations.

### 7) Capacity Guidance

The operations runbook documents disk, database, backup storage, log, and evidence artifact growth thresholds plus next actions when thresholds are crossed.

## Acceptance Criteria

Phase 28 is complete only when:

- production operations runbooks are stronger and actionable;
- secret rotation is documented;
- incident templates exist;
- backup/restore cadence is clear;
- upgrade operations are tied to release process;
- alert delivery proof pattern is documented;
- capacity guidance is documented;
- operator handover checklist exists;
- Phase 27 deployment/DB-scoped operations boundary is preserved;
- no paid SLA or universal readiness is implied.

## Required Checks

```bash
make validate
make test
make test-integration
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Also run a targeted context-aware scan of changed docs for secrets, private operator artifacts, and unsupported claims. Negated boundary language such as "no SLA coverage" is allowed.

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
- `docs/phase-28-production-operations-hardening.md`
- `docs/support-boundaries.md`
- `docs/release-process.md`
- `docs/release-checklist.md`
- `docs/upgrade-and-rollback.md`
- `docs/evidence/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-28.md`
- `scripts/pilot-ops.sh` only if safe improvements are needed
