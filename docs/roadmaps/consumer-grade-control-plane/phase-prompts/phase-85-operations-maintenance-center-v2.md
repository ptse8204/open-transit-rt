# Phase 85 Prompt — Operations And Maintenance Center V2

## Goal

Make small-host/reference deployment maintenance safer and less CLI-dependent.

## Scope

- Backup status guidance.
- Restore-drill guidance.
- Upgrade and rollback checklist.
- Secret/device-token rotation guidance.
- Support-bundle generation guidance.
- Incident/reliability review UX.
- Small-host validator/runtime blocker explanations.
- Preview-only/default-safe browser guidance for risky maintenance tasks.
- Clear `requires technical helper` markers where browser execution would be
  risky, destructive, or environment-dependent.

## Do not

- Run destructive restore actions from the browser without a later explicit plan.
- Turn restore, rollback, secret rotation, purge, or live support-bundle
  operations into browser-executed actions without separate authorization.
- Expose secrets.
- Claim SLA/uptime or production readiness.
- Write evidence.
- Create hosted service claims.

## Validation

Baseline validation plus:

```bash
make validate
make test
make deployment-doctor
make operations-reliability
make operations-notify
```

Record blockers if optional commands need env/local app.

## Commits

```text
Phase 85 -- Checkpoint 000001: add maintenance center v2 plan
Phase 85 -- Checkpoint 000002: add backup and restore drill guidance
Phase 85 -- Checkpoint 000003: add upgrade rollback guidance
Phase 85 -- Checkpoint 000004: add secret rotation and support bundle guidance
Phase 85 -- Checkpoint 000005: add incident and reliability review UX
Phase 85 -- Checkpoint 000006: close maintenance center v2 review
```
