# Phase 103 -- Monitoring, Notifications, And Export Surfaces

## Goal

Improve redacted monitoring/export surfaces and no-send notification previews
so small operators can understand unhealthy feeds without adding hosted SaaS
behavior, live delivery, evidence collection, consumer automation, or stronger
public claims.

## Current Surface

- `scripts/operations-notify.sh` creates a private local draft under
  `.cache/operations-notify/<timestamp>` with `summary.json`, `summary.md`,
  `manifest.json`, `manifest.md`, and `notification.txt`. It records
  destination presence as booleans and never sends webhook/email output.
- `scripts/operations-reliability.sh` reads validator, deployment-doctor, and
  notification summaries and writes private reliability diagnostics under
  `.cache/operations-reliability/<timestamp>`.
- `/admin/operations/maintenance` and `/admin/operations/maintenance.json`
  surface latest safe summary pointers for deployment doctor, reliability,
  operations notification, and support-bundle diagnostics.
- `/admin/operations/reliability` and `/admin/operations/reliability.json`
  already summarize private reliability diagnostics.
- Connector Workbench already documents monitoring/export recipes and no-send
  defaults, but Phase 103 should make the operator-facing monitoring export
  and health digest path more obvious.

## Master-Approved Plan

1. Add the Phase 103 plan and checkpoint report.
2. Inspect sub-agent findings before implementation.
3. Implement a bounded no-send monitoring/export improvement using existing
   surfaces first. Preferred scope:
   - improve local notification draft guidance and health digest templates;
   - expose a private operations-summary JSON or equivalent Maintenance
     Center section derived from existing view models;
   - add redacted webhook/email draft guidance that records destination
     presence only and never values;
   - clarify monitoring export summary inputs, safe outputs, and no-send
     default.
4. Avoid migrations, durable notification state, live webhooks, email sends,
   background schedulers, hosted monitoring stacks, external network contact,
   and new public routes unless a safety issue requires otherwise.
5. Add focused tests for private route access, JSON shape, no-send flags,
   redaction, claim flags, and docs wording.
6. Run focused checks, then required baseline/code-change validation.
7. Close with `docs/handoffs/phase-103.md`, `docs/handoffs/latest.md`, and
   roadmap/status updates.

## Non-Goals

- No live webhook/email send.
- No notification scheduler or hosted monitoring service.
- No public admin route.
- No retained evidence or protected-path writes.
- No consumer tracker edits.
- No real credentials, real private payloads, or real destination values.
- No SLA, uptime, hosted-service, production readiness, compliance, consumer
  acceptance, public launch, agency adoption, vendor, hardware, or ETA-quality
  claim.

## Checkpoint Plan

- `Phase 103 -- Checkpoint 000001: add monitoring, notifications, and export surfaces plan`
- `Phase 103 -- Checkpoint 000002: implement primary scoped work`
- `Phase 103 -- Checkpoint 000003: run validation and patch required gaps`
- `Phase 103 -- Checkpoint 000004: close monitoring, notifications, and export surfaces review`

## Focused Validation Targets

- `go test ./cmd/agency-config -run 'OperationsNotify|OperationsReliability|Maintenance|Reliability|Monitoring|OperationsNavigation|RouteTitles|Help'`
- `scripts/operations-notify.sh --help`
- `OUTPUT_DIR=.cache/phase-103/operations-notify FORCE=true scripts/operations-notify.sh --dry-run`
- `scripts/operations-reliability.sh --help`
- `OUTPUT_DIR=.cache/phase-103/operations-reliability FORCE=true OPERATIONS_NOTIFY_SUMMARY=.cache/phase-103/operations-notify/summary.json scripts/operations-reliability.sh --dry-run`
- JSON validation of any generated `.cache/phase-103/**/summary.json` and
  `manifest.json` files.

Because this phase is expected to change code/docs/tests/scripts, closeout also
requires:

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact prepared-only consumer tracker assertion
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`

## Checkpoint Report -- 000001

Checkpoint:
Phase 103 -- Checkpoint 000001: add monitoring, notifications, and export
surfaces plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, active for monitoring
and notification surface inspection. Real Planning Sub-Agent -- GPT-5.5
x-high, active for implementation planning. Implementation, QA, UI/UX,
Documentation / IA, Claim-Boundary, Security/Auth, Data/Migration, and
Release/Supply-Chain roles are simulated by the Master Agent for this plan
checkpoint. Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-103-monitoring-notifications-and-export-surfaces.md`.

Validation run:
Initial Phase 103 repository inspection found existing no-send
`operations-notify` and private `operations-reliability` helpers, Maintenance
Center diagnostic summary rows, and tests around notification no-send behavior.
After adding the plan, `git status --short` showed only
`docs/phase-103-monitoring-notifications-and-export-surfaces.md`; `git diff
--check` passed; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed; the exact
prepared-only consumer tracker assertion passed; and `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned no output.

Blocked checks:
Implementation, focused code tests, and closeout baseline checks are not yet
run because this checkpoint only approves the Phase 103 plan. Release-candidate
and package checks are out of scope for Phase 103.

Protected path status:
No protected evidence path is part of the plan. The plan forbids protected path
writes.

Consumer tracker status:
The consumer tracker is not part of the plan. The seven targets must remain in
order and `prepared`.

Claim-boundary status:
The plan explicitly forbids SLA, uptime, hosted-service, production readiness,
compliance, consumer acceptance, public launch, agency adoption, vendor,
hardware, release-readiness, production-grade ETA, and real-world ETA claims.

Security/auth status:
The plan preserves private/authenticated review surfaces, no-send defaults,
destination-value redaction, and no browser collection of notification
credentials.

Data/migration status:
No migration, durable notification state, scheduler, queue, or hosted
monitoring dependency is planned.

Master review:
Approved. The smallest safe Phase 103 implementation is to refine local
no-send monitoring/export guidance and private summaries over existing
diagnostic helpers.

Required edits:
Incorporate sub-agent findings, implement bounded private no-send monitoring
exports/guidance, update tests/docs, and record validation.

Decision:
Proceed to implementation checkpoint 000002 after plan validation and commit.

Next checkpoint:
Phase 103 -- Checkpoint 000002: implement primary scoped work.
