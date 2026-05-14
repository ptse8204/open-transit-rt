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

## Checkpoint 000002 Evaluator And Operations Staff Task Trials

Browser plugin status:
`blocked`.

The in-app Browser plugin was connected successfully, but attempts to open both
`http://localhost:8080/healthz` and
`http://127.0.0.1:8080/admin/operations` failed before page load with
`net::ERR_BLOCKED_BY_CLIENT`. No authenticated Operations Console page loaded
in the browser.

Fallback used:
Terminal-authenticated local route checks against the running local app. A
short-lived local demo admin token was generated and kept in process memory
only; it was not written to repository files, docs, retained evidence, or
diagnostic summaries.

Diagnostic summary path:
`.cache/phase-93/cp000002-evaluator-ops.json`

| Flow | Route | Result | Notes |
| --- | --- | --- | --- |
| Unauthenticated admin boundary | `/admin/operations` | passed | Returned local `401` without an admin token. |
| New agency evaluator | `/admin/operations` | passed | Authenticated route returned `200`; expected Start Here, setup, help, and readiness terms were present. |
| New agency evaluator | `/admin/operations/setup-wizard` | passed | Authenticated route returned `200`; expected agency setup, first-step, technical-helper, and local-scope terms were present. |
| New agency evaluator | `/admin/operations/help` | passed | Authenticated route returned `200`; expected help, first-week, glossary, and boundary terms were present. |
| New agency evaluator | `/admin/operations/readiness` | passed | Authenticated route returned `200`; expected readiness, prepared-only, and does-not-prove terms were present. |
| Operations staff | `/admin/operations/realtime` | passed | Authenticated route returned `200`; expected Realtime, Vehicle Positions, Trip Updates, Alerts, and stale-state terms were present. |
| Operations staff | `/admin/operations/feed-health` | passed | Authenticated route returned `200`; expected Feed Health, Vehicle Positions, Trip Updates, Alerts, and freshness terms were present. |
| Operations staff | `/admin/operations/maintenance` | passed | Authenticated route returned `200`; expected Maintenance, backup, upgrade, and validator terms were present. |

### Usability Notes

- Browser automation itself is a blocker in this environment, so the CP000002
  result is a local authenticated route fallback rather than a successful
  browser task trial.
- The new agency evaluator path has enough route coverage and boundary copy to
  proceed, but the Start Here page remains dense for first-time users.
- The operations staff path exposes the required realtime, feed-health, and
  maintenance concepts, but daily-state review remains spread across multiple
  pages rather than summarized in one compact blocker view.

## Checkpoint 000002 Report

Checkpoint:
Phase 93 -- Checkpoint 000002: run evaluator and operations staff task trials.

Sub-agents used or simulated, including intended model level:
Real UI/UX Sub-Agent -- GPT-5.5 high reported Start Here density and daily
operations spread as likely blockers. Real Context / Repo Truth Sub-Agent --
GPT-5.5 x-high, Planning Sub-Agent -- GPT-5.5 x-high, and Claim-Boundary /
Security Sub-Agent -- GPT-5.5 high informed this checkpoint. Master Agent --
GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-93-browser-task-trials.md`.

Validation run:
`make agency-app-up` passed before the trial. Browser plugin setup succeeded,
but local URL navigation failed with `net::ERR_BLOCKED_BY_CLIENT` for both
localhost and loopback URLs. Terminal-authenticated fallback route checks
passed for seven evaluator/operations routes. Unauthenticated
`/admin/operations` returned `401`.

Blocked checks:
In-app Browser task execution is blocked in this environment. Fallback route
checks do not prove browser interaction success.

Protected path status:
No protected evidence path was edited or generated. Fallback diagnostics stayed
under ignored `.cache`.

Consumer tracker status:
The tracker file was not edited. All seven targets must remain exactly
`prepared`.

Claim-boundary status:
This checkpoint records local/private task-trial diagnostics only. It makes no
release readiness, compliance, adoption, consumer acceptance, production
readiness, final-root readiness, hosted-service availability, vendor
compatibility, hardware certification, SLA/uptime, or ETA-quality claim.

Security/auth status:
Admin routes remained non-public without a token. The local demo admin token
was kept in process memory only and was not written to committed files or
diagnostic summaries.

Data/migration status:
No persistence, migration, GTFS data model, or realtime data model change is
included.

Master review:
Approved. Browser automation is truthfully recorded as blocked, the fallback
kept admin credentials out of retained artifacts, and the evaluator/operations
route coverage is sufficient to proceed.

Required edits:
CP000004 should add a compact role-based entry panel or equivalent IA cue if
the existing Start Here density can be improved safely. No CP000002 edit is
required before commit.

Decision:
Proceed to CP000002 validation and commit, then CP000003 technical-helper,
maintainer, and connector-evaluator trials.

Next checkpoint:
Phase 93 -- Checkpoint 000003: run technical helper maintainer and connector trials.

## Checkpoint 000003 Technical Helper, Maintainer, And Connector Trials

Browser plugin status:
Still `blocked` from CP000002 with `net::ERR_BLOCKED_BY_CLIENT` on local app
URLs. This checkpoint used the same safe terminal-authenticated fallback.

Diagnostic summary path:
`.cache/phase-93/cp000003-technical-maintainer-connector.json`

| Flow | Route | Result | Notes |
| --- | --- | --- | --- |
| Technical helper | `/admin/operations/gtfs-workbench` | passed | Authenticated route returned `200`; expected GTFS Workbench, active, draft, and rollback terms were present. |
| Technical helper | `/admin/operations/gtfs-import` | passed | Authenticated route returned `200`; expected GTFS Import, ZIP, active, and review terms were present. |
| Technical helper | `/admin/operations/gtfs-quality` | passed | Authenticated route returned `200`; expected GTFS Quality, validation, fix, and owner terms were present. |
| Technical helper | `/admin/operations/validation-center` | passed | Authenticated route returned `200`; expected Validation Center, Vehicle Positions, prepared-only, and does-not-prove terms were present. |
| Maintainer release reviewer | `/admin/operations/validation-center` | passed | Authenticated route returned `200`; expected validation, prepared-only, does-not-prove, and blocker terms were present. |
| Maintainer release reviewer | `/admin/operations/maintenance` | passed | Authenticated route returned `200`; expected Maintenance, backup, upgrade, and validator terms were present. |
| Maintainer release reviewer | `/admin/operations/consumers` | passed | Authenticated route returned `200`; expected prepared-only consumer/evidence boundary terms were present. |
| Connector evaluator | `/admin/operations/connectors/workbench` | passed | Authenticated route returned `200`; expected Connector Workbench, synthetic, recipe, and vendor-boundary terms were present. |
| Connector evaluator | `/admin/operations/connectors/tests` | passed | Authenticated route returned `200`; expected connector test, synthetic, and conformance terms were present. |
| Connector evaluator | `/admin/operations/telemetry-simulator` | needs_review | Authenticated route returned `200`; expected Telemetry Simulator, device, and synthetic terms were present, but the exact `dry-run` cue was missing. |
| Connector evaluator | `/admin/operations/prediction-lab` | passed | Authenticated route returned `200`; expected Prediction, Trip Updates, withheld, and external-predictor boundary terms were present. |

### Usability Notes

- Technical-helper schedule routes expose active/draft and rollback guidance,
  but the workbench remains broad and should keep review-vs-mutate boundaries
  prominent.
- Maintainer release-review routes expose prepared-only and blocker language,
  but there is still no single release gate posture panel in the Operations
  Console; Phase 92 docs remain the current RC gate source.
- Connector evaluator routes are coherent overall. The Telemetry Simulator page
  should explicitly label copied commands as `dry-run` or local/synthetic
  safety checks to avoid confusion with real AVL/vendor testing.

## Checkpoint 000003 Report

Checkpoint:
Phase 93 -- Checkpoint 000003: run technical helper maintainer and connector trials.

Sub-agents used or simulated, including intended model level:
Real UI/UX Sub-Agent -- GPT-5.5 high reported connector-workbench density and
synthetic/local wording risks. Real Context / Repo Truth Sub-Agent -- GPT-5.5
x-high, Planning Sub-Agent -- GPT-5.5 x-high, and Claim-Boundary / Security
Sub-Agent -- GPT-5.5 high informed this checkpoint. Master Agent -- GPT-5.5
x-high, current thread.

Changed files:
`docs/phase-93-browser-task-trials.md`.

Validation run:
Terminal-authenticated fallback route checks passed for eleven technical
helper, maintainer release reviewer, and connector evaluator routes. One
expected wording cue, `dry-run`, was missing from `/admin/operations/telemetry-simulator`.

Blocked checks:
In-app Browser task execution remains blocked in this environment. Fallback
route checks do not prove browser interaction success.

Protected path status:
No protected evidence path was edited or generated. Fallback diagnostics stayed
under ignored `.cache`.

Consumer tracker status:
The tracker file was not edited. All seven targets must remain exactly
`prepared`.

Claim-boundary status:
This checkpoint records local/private task-trial diagnostics only. It makes no
release readiness, compliance, adoption, consumer acceptance, production
readiness, final-root readiness, hosted-service availability, vendor
compatibility, hardware certification, SLA/uptime, or ETA-quality claim.

Security/auth status:
Admin route checks used a local demo bearer token kept in process memory only.
No browser mutation, external submission, evidence collection, or risky
maintenance action was performed.

Data/migration status:
No persistence, migration, GTFS data model, or realtime data model change is
included.

Master review:
Approved. The fallback route checks found one small connector-evaluator copy
gap suitable for CP000004 and no unsafe continuation condition.

Required edits:
CP000004 should add explicit dry-run/local-synthetic wording to the Telemetry
Simulator guidance and, if safe, add a compact role-based entry cue to Start
Here.

Decision:
Proceed to CP000003 validation and commit, then CP000004 copy/IA patching.

Next checkpoint:
Phase 93 -- Checkpoint 000004: patch task trial copy and IA gaps.
