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
