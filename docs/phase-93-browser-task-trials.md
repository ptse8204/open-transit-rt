# Phase 93 -- Browser End-To-End Agency Task Trials

This phase tests local browser task flows for small-agency staff and technical
helpers. It is a private/local usability and product-coherence review only: it
is not retained evidence, release proof, compliance proof, consumer proof,
final-root proof, production-readiness proof, hosted-service proof, SLA/uptime
proof, vendor proof, hardware proof, or ETA-quality proof.

## Scope

- Run local browser task trials for five user roles:
  - new agency evaluator
  - operations staff
  - technical helper
  - maintainer release reviewer
  - connector evaluator
- Record a task-flow matrix with route path, expected decision, observed
  usability blockers, and safe follow-up.
- Patch small copy, information-architecture, or route-label issues only when
  obvious, safe, private, and covered by validation.
- Close with protected-path, consumer tracker, claim-boundary, security/auth,
  and data/migration review.

## Out Of Scope

- Protected evidence path writes.
- Consumer tracker status changes.
- External contact or proof gathering.
- Real agency/vendor/device credentials, real private payloads, or real
  agency/vendor/device data.
- Release tags, GitHub Releases, public package/image publication, or public
  announcements.
- Claims of release readiness, compliance, adoption, consumer acceptance,
  production readiness, final-root readiness, hosted SaaS, vendor
  compatibility, hardware certification, SLA/uptime, or production-grade ETA
  quality.
- Browser execution of risky maintenance, deployment, evidence, or external
  submission actions.

## Checkpoints

```text
Phase 93 -- Checkpoint 000001: add browser end-to-end agency task trials plan
Phase 93 -- Checkpoint 000002: run evaluator and operations staff task trials
Phase 93 -- Checkpoint 000003: run technical helper maintainer and connector trials
Phase 93 -- Checkpoint 000004: patch task trial copy and IA gaps
Phase 93 -- Checkpoint 000005: close browser task trials
```

## Task-Flow Matrix To Produce

| Flow | Starting route | Core question | Routes to trial | Success signal |
| --- | --- | --- | --- | --- |
| New agency evaluator | `/admin/operations` | Can a nontechnical evaluator find the safest first steps without knowing phase history? | `/admin/operations`, `/admin/operations/setup-wizard`, `/admin/operations/help`, `/admin/operations/readiness` | First-click path is clear, private status is obvious, and next actions avoid false readiness claims. |
| Operations staff | `/admin/operations/realtime` | Can daily staff understand whether realtime feeds are useful and what needs attention? | `/admin/operations/realtime`, `/admin/operations/feed-health`, `/admin/operations/maintenance` | Stale/unknown states, Vehicle Positions, Trip Updates, Alerts, and maintenance next actions are understandable. |
| Technical helper | `/admin/operations/gtfs-workbench` | Can a helper review schedule state, imports, validation, and safe rollback guidance? | `/admin/operations/gtfs-workbench`, `/admin/operations/gtfs-import`, `/admin/operations/gtfs-quality`, `/admin/operations/validation-center` | Active/draft boundaries, validation meaning, import review, and rollback limitations are clear. |
| Maintainer release reviewer | `/admin/operations/validation-center` | Can a maintainer understand RC blockers and no-claim boundaries? | `/admin/operations/validation-center`, `/admin/operations/maintenance`, `/admin/operations/consumers` | `needs_review`, prepared-only consumers, package/tag blockers, and claim boundaries remain visible. |
| Connector evaluator | `/admin/operations/connectors/workbench` | Can an evaluator choose a safe synthetic/local integration path? | `/admin/operations/connectors/workbench`, `/admin/operations/connectors/tests`, `/admin/operations/telemetry-simulator`, `/admin/operations/prediction-lab` | Connector recipes, test commands, redaction, and no vendor-compatibility claim are clear. |

## Browser Trial Method

Use the in-app Browser plugin against the local app when available. Keep the
browser in the background unless visual evidence needs to be shown to the user.
Start the local app with `make agency-app-up`, perform authenticated local
browser trials against `http://localhost:8080/admin/operations`, and stop the
stack with `make agency-app-down` after diagnostics.

If Browser automation is unavailable, record the exact blocker and run
terminal-authenticated route checks as a fallback. The fallback is a product
diagnostic, not a browser success.

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

If code, scripts, routes, UI behavior, examples, or tests change, also run:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```

For local app/browser diagnostics, run:

```bash
make agency-app-up
make agency-app-down
```

Add focused `go test ./cmd/agency-config` route/UI tests when copy, route
labels, handlers, or templates change.

## Checkpoint 000001 Report

Checkpoint:
Phase 93 -- Checkpoint 000001: add browser end-to-end agency task trials plan.

Sub-agents used or simulated, including intended model level:
Real Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, Planning Sub-Agent --
GPT-5.5 x-high, UI/UX Sub-Agent -- GPT-5.5 high, and Claim-Boundary /
Security Sub-Agent -- GPT-5.5 high are running for Phase 93. Browser skill
instructions were loaded for later local browser automation. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-93-browser-task-trials.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/02-phases-and-checkpoints.md`.

Validation run:
`git status --short` was clean before edits. `git diff --check` passed.

Blocked checks:
Browser task trials, local app checks, copy/IA patch validation, and closeout
validation are scheduled for later Phase 93 checkpoints.

Protected path status:
No protected evidence path is planned or required for this checkpoint.

Consumer tracker status:
All seven targets must remain exactly `prepared`; no tracker edit is planned.

Claim-boundary status:
The plan is bounded to local/private task trials and explicit blockers. It
does not claim release readiness, compliance, adoption, consumer acceptance,
production readiness, final-root readiness, hosted-service availability,
vendor compatibility, hardware certification, SLA/uptime, or ETA quality.

Security/auth status:
No route, auth behavior, token handling, credential path, public exposure, or
admin command behavior changed.

Data/migration status:
No persistence, migration, GTFS data model, or realtime data model change is
included.

Master review:
Approved. The Phase 93 plan covers all five required task flows, uses Browser
automation as the preferred local trial surface, and keeps protected-path,
consumer-status, no-release, no-evidence, and no-claim boundaries explicit.

Required edits:
None for CP000001.

Decision:
Proceed to CP000001 validation and commit, then CP000002 evaluator and
operations staff task trials.

Next checkpoint:
Phase 93 -- Checkpoint 000002: run evaluator and operations staff task trials.
