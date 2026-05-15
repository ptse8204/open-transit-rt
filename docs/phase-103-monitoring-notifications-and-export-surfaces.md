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

## Checkpoint Report -- 000002

Checkpoint:
Phase 103 -- Checkpoint 000002: implement primary scoped work.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain roles
were simulated by the Master Agent. Master Agent -- GPT-5.5 x-high, current
thread.

Changed files:
`scripts/operations-notify.sh`;
`scripts/operations-reliability.sh`;
`cmd/agency-config/operations_maintenance.go`;
`cmd/agency-config/operations.go`;
`cmd/agency-config/main_test.go`;
`cmd/agency-config/operations_notify_script_test.go`;
`cmd/agency-config/operations_reliability_script_test.go`;
`docs/tutorials/self-hosted-operations-notifications.md`;
`docs/phase-103-monitoring-notifications-and-export-surfaces.md`.

Validation run:
`go test ./cmd/agency-config -run
'OperationsNotify|OperationsReliability|Maintenance|Reliability|Monitoring|OperationsNavigation|RouteTitles|Help'`
passed. `sh -n scripts/operations-notify.sh scripts/operations-reliability.sh`
passed. `OUTPUT_DIR=.cache/phase-103/operations-notify FORCE=true
scripts/operations-notify.sh --dry-run` passed, and generated `summary.json`
and `manifest.json` parsed as JSON. `OUTPUT_DIR=.cache/phase-103/operations-
reliability FORCE=true OPERATIONS_NOTIFY_SUMMARY=.cache/phase-103/operations-
notify/summary.json scripts/operations-reliability.sh --dry-run` passed, and
generated `summary.json` and `manifest.json` parsed as JSON. A local assertion
verified `notification.not_sent=true`, channel send flags false, destination
value recording false, `monitoring_export.not_sent=true`, and
`private_ops_summary.notification_not_sent=true`. `git diff --check` passed.
`git status --short -- .cache docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned no output.
`python3 -m json.tool docs/evidence/consumer-submissions/status.json
>/dev/null` passed, and the exact prepared-only consumer tracker assertion
passed.

Blocked checks:
Full closeout baseline, connector/example checks, `make validate`, `make
test`, and docker compose configuration validation are reserved for checkpoint
000003. Release-candidate diagnostics, package generation/audit, retained
evidence, live notification sends, external monitoring services, and consumer
actions remain out of scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check returned no output.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The required
prepared-only assertion passed and will run again in checkpoint 000003.

Claim-boundary status:
The implementation is private and no-send by default. It makes no hosted
monitoring, SLA, uptime, production readiness, compliance, consumer acceptance,
public launch, agency adoption, vendor, hardware, release-readiness,
production-grade ETA, or real-world ETA claim.

Security/auth status:
No live webhook/email send, browser send control, notification destination
value rendering, secret storage, public route, scheduler, queue, external
network dependency, raw report display, or credential collection was added.

Data/migration status:
No migration, durable notification state, delivery-attempt table, monitoring
backend, queue, scheduler, telemetry contract change, public feed contract
change, or module dependency change was added.

Master review:
Approved. The implementation stays inside existing no-send `.cache` helpers
and private Maintenance review surfaces.

Required edits:
Run full validation, patch any failures caused by Phase 103, then record
checkpoint 000003.

Decision:
Proceed to validation checkpoint 000003.

Next checkpoint:
Phase 103 -- Checkpoint 000003: run validation and patch required gaps.

## Checkpoint Report -- 000003

Checkpoint:
Phase 103 -- Checkpoint 000003: run validation and patch required gaps.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. QA, Claim-Boundary, Security/Auth,
Documentation / IA, UI/UX, Data/Migration, Implementation, and
Release/Supply-Chain roles were simulated by the Master Agent for validation.
Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-103-monitoring-notifications-and-export-surfaces.md`.

Validation run:
Passed `git status --short`; `git diff --check`; `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null`; exact
prepared-only consumer tracker assertion; `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum`; `go test ./cmd/agency-config -run
'OperationsNotify|OperationsReliability|Maintenance|Reliability|Monitoring|OperationsNavigation|RouteTitles|Help'`;
`go test ./examples/connectors/sdk/monitoring ./examples/connectors/monitoring-export
./internal/connectors ./cmd/adapter-conformance`; `sh -n
scripts/operations-notify.sh scripts/operations-reliability.sh`;
`OUTPUT_DIR=.cache/phase-103/operations-notify FORCE=true
scripts/operations-notify.sh --dry-run`; JSON validation for generated
operations-notify `summary.json` and `manifest.json`;
`OUTPUT_DIR=.cache/phase-103/operations-reliability FORCE=true
OPERATIONS_NOTIFY_SUMMARY=.cache/phase-103/operations-notify/summary.json
scripts/operations-reliability.sh --dry-run`; JSON validation for generated
operations-reliability `summary.json` and `manifest.json`; generated no-send
assertions for notification, channel, monitoring export, and private ops
summary fields; `make check`; `make external-connection-check`; `make
adapter-conformance`; `make test-connector-examples`; `make
audit-product-acceptance`; `make audit-final-claim-review`; `make validate`;
`make test`; `docker compose -f deploy/docker-compose.yml config`; final `git
status --short`; final `git diff --check`; and final `.cache` plus protected
paths status check.

Blocked checks:
No Phase 103-required checks are blocked. Release-candidate diagnostics,
release package generation, retained evidence intake, live webhook/email
sending, hosted monitoring services, consumer submission, public publication,
and tag/release/package/image publication remain out of scope and blocked by
authorization boundaries.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched. The
protected-path status check returned no output before and after validation.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven targets remain in order and all remain `prepared`.

Claim-boundary status:
Product acceptance and final claim audits passed. Phase 103 still makes no
hosted monitoring, SLA, uptime, production readiness, compliance, consumer
acceptance, public launch, agency adoption, vendor, hardware,
release-readiness, production-grade ETA, or real-world ETA claim.

Security/auth status:
Validation preserved no-send defaults, destination-value redaction, private
Maintenance review, exact `.cache` output contracts, evidence-path rejection,
symlink/unsafe source rejection coverage, and no browser notification send.

Data/migration status:
The protected migration/module status check returned no output. No migration,
durable notification state, delivery-attempt table, queue, scheduler,
monitoring backend, telemetry contract change, public feed contract change, or
module dependency change was added.

Master review:
Approved. Validation passed after the scoped implementation, with no required
code patches at this checkpoint.

Required edits:
Close Phase 103 with handoff/status docs and final checkpoint commit.

Decision:
Proceed to closeout checkpoint 000004.

Next checkpoint:
Phase 103 -- Checkpoint 000004: close monitoring, notifications, and export
surfaces review.

## Checkpoint Report -- 000004

Checkpoint:
Phase 103 -- Checkpoint 000004: close monitoring, notifications, and export
surfaces review.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high; real Planning
Sub-Agent -- GPT-5.5 x-high. Implementation, QA, UI/UX, Documentation / IA,
Claim-Boundary, Security/Auth, Data/Migration, and Release/Supply-Chain
closeout roles were simulated by the Master Agent. Master Agent -- GPT-5.5
x-high, current thread.

Changed files:
`docs/phase-103-monitoring-notifications-and-export-surfaces.md`;
`docs/handoffs/phase-103.md`;
`docs/handoffs/latest.md`;
`docs/current-status.md`;
`docs/roadmap-status.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`.

Validation run:
Closeout relies on the checkpoint 000003 validation pass: focused script/UI
tests, connector/example tests, script syntax checks, generated no-send summary
checks, baseline checks, connector-conformance checks, product acceptance
audit, final claim audit, `make validate`, `make test`, docker compose config,
final protected-path check, and final `git diff --check`.

Blocked checks:
Release-candidate diagnostics, package generation/audit, retained evidence,
live webhook/email sends, hosted monitoring services, consumer submission,
public publication, and tag/release/package/image publication remain blocked
by scope.

Protected path status:
No protected evidence path was edited, generated, reformatted, or touched.

Consumer tracker status:
`docs/evidence/consumer-submissions/status.json` was not edited. The exact
seven targets remain in order and all remain `prepared`.

Claim-boundary status:
Phase 103 closes without hosted monitoring, SLA, uptime, production readiness,
compliance, consumer acceptance, public launch, agency adoption, vendor,
hardware, release-readiness, production-grade ETA, or real-world ETA claim.

Security/auth status:
No live send path, destination-value rendering, credential collection, public
route, scheduler, queue, external network dependency, or raw payload rendering
was added.

Data/migration status:
No migration, durable notification state, delivery-attempt table, monitoring
backend, queue, scheduler, telemetry contract change, public feed mutation, or
module dependency change was added.

Master review:
Approved. Phase 103 is complete and closed.

Required edits:
None.

Decision:
Continue immediately to Phase 104.

Next checkpoint:
Phase 104 -- Checkpoint 000001: add small-host deployment and upgrade
hardening plan.
