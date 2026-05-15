# Phase 100 -- Alerts Operations And Disruption Workflow

## Scope

Phase 100 improves the private Alerts Console and supporting documentation for
small-agency disruption workflows. The implementation must stay inside the
existing persisted Service Alerts and canceled-trip reconciliation boundaries.

Allowed work:

- improve private alert authoring review guidance;
- add a private lifecycle dashboard derived from existing alert records;
- surface canceled-trip linkage and missing-alert hints using existing
  cancellation reconciliation concepts;
- improve disruption templates for common operator use cases;
- add GTFS-RT Alerts validation guidance and public-feed usefulness review;
- add focused tests and documentation.

Not allowed:

- public launch or consumer acceptance claims;
- external contact or third-party portal automation;
- retained evidence collection;
- protected evidence path writes;
- consumer tracker status changes;
- real agency/vendor/device data or credentials;
- release, tag, package publication, or GitHub Release creation;
- new public admin routes or public mutation routes;
- broad Alert/Trip Update data-model rewrites unless a safety issue requires
  stopping and re-planning.

## Existing Surfaces

Alert authoring and feed rendering already exist:

- `cmd/feed-alerts/main.go` serves `/admin/alerts/console`, `/admin/alerts`,
  `/admin/alerts/{id}/publish`, `/admin/alerts/{id}/archive`, and
  `/admin/alerts/reconcile-cancellations`.
- `internal/alerts` owns persisted Service Alerts, lifecycle state,
  audit logging, and canceled-trip reconciliation.
- `internal/feed/alerts` owns GTFS-RT Alerts protobuf/debug rendering from
  published active alerts.
- `testdata/replay/cancellation-alert-linkage.json` and
  `testdata/replay/disruption-diagnostics-baseline.json` keep cancellation
  and missing-alert metrics visible.

Phase 100 should not replace these surfaces. It should make the operator path
clearer and safer.

## Master-Approved Plan

1. Add this plan and record the Phase 100 claim/security/data boundaries.
2. Implement bounded private Alerts Console view-model/template/test updates
   for lifecycle review, disruption templates, cancellation linkage guidance,
   validation guidance, and feed usefulness review.
3. Run focused alerts tests plus required baseline/heavier validation; patch
   changed-code failures only.
4. Close Phase 100 with a handoff, status-doc updates, protected-path status,
   prepared-only consumer tracker confirmation, and exact blockers.

The Master Agent approves implementation only if it remains private,
role-gated, no-external-contact, no-evidence, no-claim, and uses existing
alerts persistence/authoring contracts.

## Sub-Agent Plan

| Role | Intended model | Use in Phase 100 |
| --- | --- | --- |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high | Read-only inspection of Alerts Console, alerts domain, Alerts feed builder, cancellation linkage, tests, and requirements. |
| Planning Sub-Agent | GPT-5.5 x-high | Read-only checkpoint plan, validation plan, and guardrail review. |
| Implementation Sub-Agent | GPT-5.5 high | Simulated by Master unless a bounded disjoint edit becomes useful. |
| QA Sub-Agent | GPT-5.5 high | Simulated by Master through focused alerts tests and full required validation. |
| UI/UX Sub-Agent | GPT-5.5 high | Simulated by Master for Alerts Console wording/table shape. |
| Documentation / IA Sub-Agent | GPT-5.5 high | Simulated by Master for phase docs, tutorials/status, and handoff. |
| Claim-Boundary Sub-Agent | GPT-5.5 high | Simulated by Master with claim audit and forbidden wording review. |
| Security/Auth Sub-Agent | GPT-5.5 high | Simulated by Master; preserve private role gates, CSRF behavior, and no public mutation. |
| Data/Migration Sub-Agent | GPT-5.5 high | Simulated by Master because no migration is planned. Stop before adding persistence. |

## Checkpoints

```text
Phase 100 -- Checkpoint 000001: add alerts operations and disruption workflow plan
Phase 100 -- Checkpoint 000002: implement primary scoped work
Phase 100 -- Checkpoint 000003: run validation and patch required gaps
Phase 100 -- Checkpoint 000004: close alerts operations and disruption workflow review
```

## Validation Plan

Focused checks:

```bash
go test ./cmd/feed-alerts ./internal/alerts ./internal/feed/alerts
go test ./cmd/agency-config -run 'Realtime|FeedHealth|Readiness|OperationsNavigation|RouteTitles'
go test ./internal/realtimequality -run 'Cancellation|Disruption|Replay'
```

Phase closeout baseline:

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

Connector and release-candidate checks are not required for Phase 100 unless
the implementation unexpectedly changes connector or release surfaces.

## Implementation Summary

Checkpoint 000002 improved the private Alerts Console without adding routes or
migrations:

- lifecycle dashboard rows for draft, published, archived, active, upcoming,
  expired, operator-authored, cancellation-reconciled, unscoped, and
  indefinite alert review;
- a console form for the existing canceled-trip reconciliation path;
- disruption templates for canceled trips, detours, significant delays, stop
  closures, and modified/added service;
- GTFS-RT Alerts validation guidance, feed-health review guidance,
  missing-alert hints, and public-feed usefulness rows;
- all-false claim flags on the page;
- maintenance-guide instructions for weekly alert workflow review.

The implementation reuses existing `service_alert`,
`service_alert_informed_entity`, `internal/alerts`, and `internal/feed/alerts`
contracts. It does not change public Alerts feed semantics, Trip Updates
prediction ownership, or persisted schema.

## Checkpoint 000001 Report

Checkpoint: Phase 100 -- Checkpoint 000001: add alerts operations and disruption workflow plan

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high were launched read-only. Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration are simulated by the Master Agent for this planning checkpoint.

Changed files: `docs/phase-100-alerts-operations-and-disruption-workflow.md`

Validation run: `git status --short` before edits; source/doc read of Phase 100 prompt, Alerts Console, alerts domain repository/reconciliation, Alerts feed renderer, alerts tests, realtime replay cancellation fixtures, requirements, current status, latest handoff, roadmap status, and post-90 validation/operating manuals.

Blocked checks: Full validation deferred until implementation and closeout checkpoints. A mistaken read of `internal/feed/alerts/builder.go` failed because the feed builder lives in `internal/feed/alerts/alerts.go`; this was an exploration miss, not a product blocker.

Protected path status: No protected evidence path edits planned or made.

Consumer tracker status: Must remain exactly seven prepared-only targets; no status edits planned or made.

Claim-boundary status: Plan explicitly forbids public launch, consumer acceptance, compliance, production readiness, release readiness, hosted-service, vendor, hardware, SLA, and evidence claims.

Security/auth status: Plan preserves private role-gated Alerts Console and existing CSRF/mutation behavior; no public mutation or external contact planned.

Data/migration status: No migration, new persisted model, or raw private data storage planned.

Master review: Approved to proceed under private, role-gated, no-external-contact, no-evidence, no-claim constraints.

Required edits: Implement bounded private Alerts Console workflow guidance with tests and docs.

Decision: Continue to Checkpoint 000002.

Next checkpoint: Phase 100 -- Checkpoint 000002: implement primary scoped work

## Checkpoint 000002 Report

Checkpoint: Phase 100 -- Checkpoint 000002: implement primary scoped work

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high completed read-only review. Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration were simulated by the Master Agent for this bounded implementation checkpoint.

Changed files: `cmd/feed-alerts/main.go`, `cmd/feed-alerts/main_test.go`, `docs/tutorials/small-agency-maintenance-guide.md`, `docs/phase-100-alerts-operations-and-disruption-workflow.md`

Validation run: `gofmt -w cmd/feed-alerts/main.go cmd/feed-alerts/main_test.go`; `git diff --check`; `go test ./cmd/feed-alerts ./internal/feed/alerts ./internal/alerts`; `go test ./internal/realtimequality -run 'Cancellation|Disruption|Replay'`; `go test ./cmd/agency-config -run 'Realtime|FeedHealth|Readiness|OperationsNavigation|RouteTitles'`

Blocked checks: Full phase closeout validation deferred to Checkpoint 000003.

Protected path status: No protected evidence path edits made.

Consumer tracker status: No consumer tracker edits made; prepared-only status must be rechecked in Checkpoint 000003.

Claim-boundary status: Added console and docs wording keep Alerts workflow private/operator-facing and non-evidentiary; no public launch, consumer acceptance, compliance, production readiness, release readiness, vendor, hardware, hosted-service, or SLA claim added.

Security/auth status: Existing private role gates remain in place. The new console reconciliation form posts to the existing private `/admin/alerts/console/reconcile-cancellations` action and derives agency/actor from the authenticated principal.

Data/migration status: No migration, new durable model, public feed mutation, prediction adapter coupling, or connector runtime change added.

Master review: Approved. The implementation uses the existing Alerts-owned reconciliation and feed-rendering boundaries.

Required edits: Run full required validation and patch any changed-code failures.

Decision: Continue to Checkpoint 000003.

Next checkpoint: Phase 100 -- Checkpoint 000003: run validation and patch required gaps

## Checkpoint 000003 Report

Checkpoint: Phase 100 -- Checkpoint 000003: run validation and patch required gaps

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high completed read-only review. QA, Claim-Boundary, Security/Auth, Data/Migration, UI/UX, Documentation / IA, and Implementation were simulated by the Master Agent for validation and audit.

Changed files: `docs/phase-100-alerts-operations-and-disruption-workflow.md`

Validation run: `git status --short`; `git diff --check`; `go test ./cmd/feed-alerts ./internal/alerts ./internal/feed/alerts ./cmd/agency-config ./internal/architecture`; `go test ./internal/realtimequality -run 'Cancellation|Disruption|Replay'`; `make audit-operations-route-inventory`; `make check`; `make audit-product-acceptance`; `make audit-final-claim-review`; `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`; exact prepared-only consumer tracker assertion; `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`; `make validate`; `make test`; `docker compose -f deploy/docker-compose.yml config`; final `git status --short`; final `git diff --check`

Blocked checks: None for Phase 100.

Protected path status: Protected evidence and consumer-submission paths remained untouched; protected-path status check returned no output.

Consumer tracker status: `docs/evidence/consumer-submissions/status.json` parsed successfully and remained exactly seven targets, all `prepared`.

Claim-boundary status: Product acceptance and final claim audits passed; no public launch, consumer acceptance, compliance, production readiness, release readiness, vendor, hardware, hosted-service, SLA, evidence, or ETA-quality claim detected.

Security/auth status: Existing private role gates, agency scoping, CSRF token fields, no public mutation routes, and Alerts-owned reconciliation boundaries remained intact.

Data/migration status: No migration or persistence-model changes; `db/migrations`, `go.mod`, and `go.sum` status checks returned no output.

Master review: Approved. Validation passed with no required code patches after the earlier missing import was corrected before Checkpoint 000002 commit.

Required edits: Add Phase 100 closeout handoff/status updates.

Decision: Continue to Checkpoint 000004.

Next checkpoint: Phase 100 -- Checkpoint 000004: close alerts operations and disruption workflow review

## Checkpoint 000004 Report

Checkpoint: Phase 100 -- Checkpoint 000004: close alerts operations and disruption workflow review

Sub-agents used or simulated, including intended model level: Context / Repo Truth Sub-Agent GPT-5.5 x-high and Planning Sub-Agent GPT-5.5 x-high completed read-only review. Implementation, QA, UI/UX, Documentation / IA, Claim-Boundary, Security/Auth, and Data/Migration were simulated by the Master Agent for closeout.

Changed files: `docs/phase-100-alerts-operations-and-disruption-workflow.md`, `docs/handoffs/phase-100.md`, `docs/handoffs/latest.md`, `docs/current-status.md`, `docs/roadmap-status.md`, `docs/open-transit-rt-master-planner-remaining-work.md`

Validation run: Reused Checkpoint 000003 full validation for closeout plus `git diff --check`, exact prepared-only consumer tracker assertion, and protected-path status check after status-doc edits.

Blocked checks: Release-candidate diagnostics/package checks and connector-specific checks were not run because Phase 100 is not a release or connector phase. Retained evidence, external contact, consumer action, tag/release/package publication, and public claims remain blocked by scope.

Protected path status: No protected evidence path edits made.

Consumer tracker status: `docs/evidence/consumer-submissions/status.json` was not edited and remains prepared-only from the Checkpoint 000003 assertion.

Claim-boundary status: Closeout documents keep Alerts workflow private/operator-facing and non-evidentiary; no public launch, consumer acceptance, compliance, production readiness, release readiness, vendor, hardware, hosted-service, SLA, or adoption claim added.

Security/auth status: No security/auth behavior changed during closeout.

Data/migration status: No data or migration change during closeout.

Master review: Approved. Phase 100 is complete and safe to close.

Required edits: None.

Decision: Phase 100 complete. Continue to Phase 101.

Next checkpoint: Phase 101 -- Checkpoint 000001: add connector maturity and adapter recipes v2 plan
