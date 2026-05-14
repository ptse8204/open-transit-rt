# Phase 91 -- Maintainer Route/Product Audit And Stabilization

This phase verifies that the post-90 private Operations Console is coherent,
route-complete, task-aligned, and source-of-truth aligned before adding more
feature breadth.

It does not collect retained evidence, contact external parties, move consumer
statuses, publish a package/tag/image/release, or claim release readiness,
compliance, adoption, consumer acceptance, production readiness, final-root
readiness, hosted-service availability, vendor compatibility, hardware
certification, SLA/uptime, or production-grade ETA quality.

## Scope

- Reconcile the post-90 roadmap pack with the autonomous Phase 91-110
  authorization.
- Audit private Operations Console routes against navigation, route handlers,
  JSON/HTML pairing, README/docs route maps, and small-agency user tasks.
- Add a local route inventory audit helper or test.
- Patch highest-priority information architecture, copy, or route-map gaps that
  are safe and source-of-truth aligned.
- Close with validation, protected-path review, consumer tracker review, and
  claim-boundary review.

## Out Of Scope

- Protected evidence path writes.
- Consumer tracker status changes.
- Evidence collection or external contact.
- Release tags, GitHub Releases, public package/image publication, or public
  announcements.
- Heavy frontend stack introduction.
- Production readiness, compliance, public launch, vendor, hardware, or ETA
  quality claims.

## Checkpoints

```text
Phase 91 -- Checkpoint 000001: add route product audit plan
Phase 91 -- Checkpoint 000002: audit private routes user tasks and docs drift
Phase 91 -- Checkpoint 000003: add route inventory audit helper
Phase 91 -- Checkpoint 000004: patch highest priority IA copy and route gaps
Phase 91 -- Checkpoint 000005: close route product audit
```

## Route/User-Task Matrix To Produce

| User task | Primary private route | Supporting route(s) | Audit focus |
| --- | --- | --- | --- |
| New agency evaluator starts review | `/admin/operations` | `/admin/operations/setup-wizard`, `/admin/operations/help` | First-click clarity, no phase-history dependency, claim boundaries. |
| Operations staff checks daily state | `/admin/operations/realtime` | `/admin/operations/feed-health`, `/admin/operations/maintenance` | Freshness, stale/unknown states, next action clarity. |
| Technical helper imports/reviews GTFS | `/admin/operations/gtfs-workbench` | `/admin/operations/gtfs-import`, `/admin/operations/gtfs-quality`, `/admin/operations/validation-center` | Active vs draft separation, validation source clarity, rollback guidance. |
| Maintainer release reviewer audits RC state | `/admin/operations/validation-center` | `/admin/operations/maintenance`, `/admin/operations/consumers` | `needs_review` truthfulness, no release-ready claim. |
| Connector evaluator chooses safe integration path | `/admin/operations/connectors/workbench` | `/admin/operations/connectors/tests`, `/admin/operations/telemetry-simulator`, `/admin/operations/prediction-lab` | Synthetic/local-only recipes, redaction, no vendor compatibility claim. |

## Validation Plan

Baseline closeout validation:

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
```

If code, scripts, build behavior, route behavior, examples, or tests change,
also run:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

## Checkpoint 000001 Report

Checkpoint:
Phase 91 -- Checkpoint 000001: add route product audit plan.

Sub-agents used or simulated, including intended model level:
Real sub-agents spawned for Phase 91 planning context: Context / Repo Truth
Sub-Agent -- GPT-5.5 x-high; Planning Sub-Agent -- GPT-5.5 x-high;
Claim-Boundary/Security Sub-Agent -- GPT-5.5 high. Master Agent -- GPT-5.5
x-high, current thread.

Changed files:
`docs/phase-91-maintainer-route-product-audit.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/README.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/00-CODEX-READ-ME-FIRST.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/01-roadmap-overview.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/04-validation-and-claim-boundaries.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/05-autonomous-run-policy.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/CODEX-KICKOFF-AUTONOMOUS-PHASE-91-TO-110.md`.

Validation run:
`git diff --check` passed. `make check` passed. `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod
go.sum` returned clean.

Blocked checks:
Full route audit, route helper, app validation, and closeout validation are
scheduled for later Phase 91 checkpoints.

Protected path status:
No protected evidence path is planned or required for this checkpoint.

Consumer tracker status:
All seven targets must remain exactly `prepared`; no tracker edit is planned.

Claim-boundary status:
Plan and roadmap reconciliation use only local/product-roadmap wording and do
not claim release readiness, compliance, adoption, consumer acceptance,
production readiness, final-root readiness, hosted-service availability,
vendor compatibility, hardware certification, SLA/uptime, or ETA quality.

Security/auth status:
No route, auth behavior, command route, credential path, or browser mutation is
changed in this checkpoint.

Data/migration status:
No persistence or migration change is included.

Master review:
Approved. The Phase 91 plan matches the autonomous authorization, preserves the
Vehicle Positions-first and pluggable Trip Updates boundaries, and keeps the
protected evidence and prepared-only consumer tracker constraints explicit.

Required edits:
None for CP000001.

Decision:
Proceed to CP000001 validation and commit, then CP000002 route/task/docs drift
audit.

Next checkpoint:
Phase 91 -- Checkpoint 000002: audit private routes user tasks and docs drift.
