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

## Checkpoint 000002 Route/User-Task Audit

### Current Private Route Inventory

The Phase 90 private route inventory remains broadly accurate, but the public
README/wiki route maps lag the private Operations Console navigation. The
runtime navigation currently exposes these primary user-facing private routes:

| Area | Navigation routes |
| --- | --- |
| Start Here | `/admin/operations`, `/admin/operations/launchpad`, `/admin/operations/setup-wizard`, `/admin/operations/setup` |
| Schedule | `/admin/operations/gtfs-workbench`, `/admin/operations/gtfs-import`, `/admin/gtfs-studio`, `/admin/operations/feeds`, `/admin/operations/feed-health`, `/admin/operations/gtfs-quality`, `/admin/operations/validation-health` |
| Realtime | `/admin/operations/realtime`, `/admin/operations/prediction-lab`, `/admin/operations/telemetry`, `/admin/operations/devices`, `/admin/operations/telemetry-simulator`, `/admin/alerts/console` |
| Connectors | `/admin/operations/connectors`, `/admin/operations/connectors/workbench`, `/admin/operations/connectors/tests` |
| Health | `/admin/operations/validation-center`, `/admin/operations/readiness`, `/admin/operations/checklist`, `/admin/operations/reliability` |
| Maintain | `/admin/operations/maintenance`, `/admin/operations/access`, `/admin/operations/audit` |
| Learn | `/admin/operations/help`, `/admin/operations/consumers`, `/admin/operations/evidence` |

Private JSON/read-only companion routes exist for the newer center-style pages:
`/admin/operations.json`, `/admin/operations/launchpad.json`,
`/admin/operations/setup-wizard.json`, `/admin/operations/feed-health.json`,
`/admin/operations/validation-center.json`, `/admin/operations/readiness.json`,
`/admin/operations/telemetry-simulator.json`,
`/admin/operations/realtime.json`, `/admin/operations/prediction-lab.json`,
`/admin/operations/connectors.json`,
`/admin/operations/connectors/workbench.json`,
`/admin/operations/connectors/tests.json`,
`/admin/operations/gtfs-workbench.json`,
`/admin/operations/validation-health.json`,
`/admin/operations/reliability.json`,
`/admin/operations/maintenance.json`, `/admin/operations/access.json`,
`/admin/operations/audit.json`, and `/admin/operations/help.json`.

The private progressive JavaScript asset remains
`/admin/operations/assets/operations.js`. The only read-only command JSON route
in this surface remains
`POST /admin/operations/validation-health/refresh.json`.

### Route/User-Task Matrix

| User task | Primary route | Supporting route(s) | Current audit result | Required edit |
| --- | --- | --- | --- | --- |
| New agency evaluator starts review | `/admin/operations` | `/admin/operations/setup-wizard`, `/admin/operations/help` | Coherent. README and wiki still start at the correct private URL and first-click label. | None required. |
| Operations staff checks daily state | `/admin/operations/realtime` | `/admin/operations/feed-health`, `/admin/operations/maintenance` | Runtime route exists and operator training points to it, but README/wiki route maps omit it. | Patch route maps and task docs to include Realtime Center. |
| Technical helper imports/reviews GTFS | `/admin/operations/gtfs-workbench` | `/admin/operations/gtfs-import`, `/admin/operations/gtfs-quality`, `/admin/operations/validation-center` | Runtime route exists and operator training points to it, but README/wiki route maps still lead mostly through import/quality/validation-health. | Patch route maps to make GTFS Workbench the schedule review front door. |
| Maintainer release reviewer audits RC state | `/admin/operations/validation-center` | `/admin/operations/maintenance`, `/admin/operations/consumers` | Runtime route exists and Phase 90 inventory includes it, but README/wiki route maps omit it. | Patch route maps to include Validation Center. |
| Connector evaluator chooses safe integration path | `/admin/operations/connectors/workbench` | `/admin/operations/connectors/tests`, `/admin/operations/telemetry-simulator`, `/admin/operations/prediction-lab` | Runtime route exists and integration docs mention it, but README/wiki route maps omit it. | Patch route maps and connector task guidance to include Connector Workbench. |

### Detailed Audit Findings

| Finding | Severity | Evidence | Required follow-up |
| --- | --- | --- | --- |
| Roadmap pack checkpoint drift | medium | `docs/roadmaps/post-90-agency-grade-gtfs-rt-product/02-phases-and-checkpoints.md` and Phase 91 phase prompt still use generic four-checkpoint wording, while the autonomous kickoff requires five concrete Phase 91 checkpoints. | Patch roadmap pack checkpoint language in CP000004. |
| README and wiki route maps are stale | high | `README.md`, `wiki/README.md`, and `wiki/small-agency-quick-start.md` omit `/admin/operations/gtfs-workbench`, `/admin/operations/realtime`, `/admin/operations/validation-center`, `/admin/operations/connectors/workbench`, `/admin/operations/prediction-lab`, `/admin/operations/access`, and `/admin/operations/audit`. | Patch route maps in CP000004. |
| Browser-first setup docs underrepresent center-style pages | medium | `wiki/browser-first-setup.md`, `wiki/operations-console-tour.md`, and `docs/tutorials/no-cli-agency-first-run.md` still emphasize older route order and omit several Phase 80-86 center routes. | Patch only highest-value task-path references in CP000004. |
| Runtime route registration is coherent | none | Navigation routes have matching handlers through explicit `mux.Handle` registrations or the `/admin/operations/` catch-all and `operationsRoot` switch. Focused tests cover grouped navigation, active route state, titles, private/scoped/no-store behavior for representative newer routes. | Add a read-only route inventory helper in CP000003 to make this repeatable. |
| Route truth is duplicated | medium | Route metadata is repeated across `cmd/agency-config/main.go`, `operationsRoot`, `operationsNavGroups`, page-title switches, tests, templates, and docs. | CP000003 helper should detect drift now; Phase 94 should centralize route metadata. |
| Trailing-slash aliases are implicit | low | `operationsRoot` trims `/admin/operations/` suffixes, so some trailing-slash aliases can return the same page without being canonical route inventory entries. | CP000003 helper should audit canonical routes, not treat implicit aliases as primary route map entries. |
| Legacy generic pages need cache-header review | medium | `feeds`, `telemetry`, `devices`, `consumers`, `evidence`, and `setup` share `renderOperations` and are not as visibly wrapped with `Cache-Control: no-store` as newer routes. | CP000003 helper should flag routes requiring explicit no-store review; CP000004 can patch if needed. |
| Progressive JS allowlist is narrower than JSON inventory | low | `operations_admin.js` allows root JSON, single-segment `/admin/operations/*.json`, and one POST refresh route. Nested JSON routes such as `/admin/operations/connectors/workbench.json` are linked but not JS-fetchable. | No immediate bug unless future JS uses nested JSON. Phase 94 should decide whether route metadata owns the JS allowlist. |
| No public admin route expansion found | none | Public routes remain `/public/*`; Operations Console routes remain under authenticated `/admin/*`. | Preserve boundary. |
| Claim-boundary language remains present | none | Route templates repeatedly distinguish diagnostics from evidence, compliance, consumer acceptance, release readiness, production readiness, vendor compatibility, hardware certification, SLA/uptime, and ETA-quality proof. | Preserve boundary in CP000004 wording. |

### CP000002 Report

Checkpoint:
Phase 91 -- Checkpoint 000002: audit private routes user tasks and docs drift.

Sub-agents used or simulated, including intended model level:
Planning Sub-Agent -- GPT-5.5 x-high, real; Claim-Boundary/Security Sub-Agent
-- GPT-5.5 high, real; Context / Repo Truth Sub-Agent -- GPT-5.5 x-high,
real; Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`docs/phase-91-maintainer-route-product-audit.md`.

Validation run:
`git diff --check` passed. `make check` passed. `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod
go.sum` returned clean.

Blocked checks:
No code/script route helper exists yet; CP000003 will add it. No app startup or
browser trial was required for this audit-only checkpoint.

Protected path status:
No protected evidence path was modified.

Consumer tracker status:
No consumer tracker edit was made. All seven targets are required to remain
exactly `prepared`.

Claim-boundary status:
Audit wording remains limited to private route/product coherence and local docs
drift. It makes no compliance, adoption, consumer action, final-root,
hosted-service, production-readiness, vendor, hardware, SLA/uptime,
release-ready, public-launch, or ETA-quality claim.

Security/auth status:
The audit found no required public admin route expansion. Follow-up helper must
stay read-only, local-only, no-network, no-Docker, and must not dump private
JSON bodies or execute commands.

Data/migration status:
No persistence or migration change.

Master review:
Approved. The next implementation checkpoint should make route inventory drift
auditable through a local helper before patching route maps.

Required edits:
Run focused validation and commit CP000002. CP000003 must add the read-only
route inventory audit helper. CP000004 must patch stale route maps and generic
roadmap checkpoint wording.

Decision:
Proceed to CP000003 after CP000002 validation and commit.

Next checkpoint:
Phase 91 -- Checkpoint 000003: add route inventory audit helper.

## Checkpoint 000003 Route Inventory Audit Helper

Phase 91 added a local, read-only route inventory audit helper:

```bash
make audit-operations-route-inventory
make test-operations-route-inventory
```

The helper parses committed source and docs only. It does not start the app,
fetch routes, call external URLs, run validators, contact consumers, write
`.cache` outputs, write evidence, or mutate data.

The helper currently checks:

- canonical private Operations Console HTML routes have navigation and handler
  coverage;
- canonical private JSON routes have handler coverage;
- the read-only command JSON route is explicitly registered and allowlisted;
- GTFS Studio and Alerts Console are marked as external admin surfaces in nav;
- no public admin route registration is present;
- Phase 90 route inventory still contains the canonical private routes;
- README/wiki route maps contain newer center-style routes when
  `OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true`.

Current helper result:

```text
PASS: 28 canonical private HTML routes have nav and handler coverage
PASS: 19 canonical private JSON routes have handler coverage
PASS: 1 private command route is explicit and allowlisted
PASS: 2 external admin surfaces are marked in nav
PASS: no public admin route registration detected
WARN: route map docs omit newer Operations route: /admin/operations/access
WARN: route map docs omit newer Operations route: /admin/operations/audit
WARN: route map docs omit newer Operations route: /admin/operations/connectors/workbench
WARN: route map docs omit newer Operations route: /admin/operations/gtfs-workbench
WARN: route map docs omit newer Operations route: /admin/operations/prediction-lab
WARN: route map docs omit newer Operations route: /admin/operations/realtime
WARN: route map docs omit newer Operations route: /admin/operations/validation-center
```

The warnings are the CP000004 work queue. They remain warnings rather than
failures in default mode so the helper can be introduced before the route-map
copy patch.

### CP000003 Report

Checkpoint:
Phase 91 -- Checkpoint 000003: add route inventory audit helper.

Sub-agents used or simulated, including intended model level:
Planning Sub-Agent -- GPT-5.5 x-high, real; Claim-Boundary/Security Sub-Agent
-- GPT-5.5 high, real; Context / Repo Truth Sub-Agent -- GPT-5.5 x-high,
real; Implementation Sub-Agent -- GPT-5.5 high, simulated by Master; QA
Sub-Agent -- GPT-5.5 high, simulated by Master; Master Agent -- GPT-5.5
x-high, current thread.

Changed files:
`Makefile`; `scripts/audit-operations-route-inventory.sh`;
`scripts/test-operations-route-inventory.sh`;
`docs/phase-91-maintainer-route-product-audit.md`.

Validation run:
`scripts/audit-operations-route-inventory.sh` passed with expected route-map
warnings. `scripts/test-operations-route-inventory.sh` passed.
`make test-operations-route-inventory` passed. `git diff --check` passed.
`make check` passed. `git status --short --
docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod
go.sum` returned clean.

Blocked checks:
Full closeout validation, `make validate`, `make test`, and Docker Compose
config are deferred to Phase 91 closeout because this checkpoint added scripts
and Make targets but did not change Go route behavior.

Protected path status:
No protected evidence path was modified.

Consumer tracker status:
No consumer tracker edit was made. All seven targets are required to remain
exactly `prepared`.

Claim-boundary status:
The helper reports local private route inventory only. It does not claim
compliance, adoption, consumer action, final-root readiness, hosted-service
availability, production readiness, release readiness, vendor compatibility,
hardware certification, SLA/uptime, public launch, or ETA quality.

Security/auth status:
The helper is local-only and read-only. It does not start the app, fetch
private pages, dump JSON bodies, execute validators, call external services, or
inspect credentials.

Data/migration status:
No persistence or migration change.

Master review:
Approved. The helper makes the CP000002 drift repeatable and is safe to run as
part of `make check`.

Required edits:
CP000004 must patch stale route maps and reconcile generic roadmap checkpoint
language. It should also decide whether to make the current route-map warnings
clean in default audit mode.

Decision:
Proceed to CP000004 after CP000003 validation and commit.

Next checkpoint:
Phase 91 -- Checkpoint 000004: patch highest priority IA copy and route gaps.

## Checkpoint 000004 IA And Route-Gap Patches

Phase 91 patched the highest-priority route/product drift found in CP000002 and
made CP000003 route-map warnings clean:

- README and wiki route maps now include GTFS Workbench, Realtime Center,
  Validation Center, Connector Workbench, Prediction & ETA Lab, Access & Roles,
  and Audit Log.
- Browser-first setup and no-CLI tutorial guidance now points operators to the
  center-style schedule, realtime, validation, and connector pages before
  lower-level diagnostics.
- The Operations Console tour now explains GTFS Workbench, Connector Workbench,
  Validation Center, Realtime Center, Prediction & ETA Lab, Access & Roles, and
  Audit Log in the same private/no-claim style as the existing pages.
- The post-90 roadmap pack now uses the required five Phase 91 checkpoint names
  and the actual Phase 91 plan path.
- Legacy generic private Operations pages now set `Cache-Control: no-store`
  before GET or POST handling.

Strict route inventory audit now reports:

```text
PASS: 28 canonical private HTML routes have nav and handler coverage
PASS: 19 canonical private JSON routes have handler coverage
PASS: 1 private command route is explicit and allowlisted
PASS: 2 external admin surfaces are marked in nav
PASS: no public admin route registration detected
PASS: README/wiki route maps include newer center-style routes
```

### CP000004 Report

Checkpoint:
Phase 91 -- Checkpoint 000004: patch highest priority IA copy and route gaps.

Sub-agents used or simulated, including intended model level:
Planning Sub-Agent -- GPT-5.5 x-high, real; Claim-Boundary/Security Sub-Agent
-- GPT-5.5 high, real; Context / Repo Truth Sub-Agent -- GPT-5.5 x-high,
real; Implementation Sub-Agent -- GPT-5.5 high, simulated by Master; QA
Sub-Agent -- GPT-5.5 high, simulated by Master; UI/UX Sub-Agent -- GPT-5.5
high, simulated by Master; Documentation / IA Sub-Agent -- GPT-5.5 high,
simulated by Master; Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`README.md`; `wiki/README.md`; `wiki/small-agency-quick-start.md`;
`wiki/browser-first-setup.md`; `wiki/operations-console-tour.md`;
`docs/tutorials/no-cli-agency-first-run.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/02-phases-and-checkpoints.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/phase-prompts/phase-91-maintainer-route-product-audit-and-stabilization.md`;
`cmd/agency-config/operations.go`; `cmd/agency-config/main_test.go`;
`docs/phase-91-maintainer-route-product-audit.md`.

Validation run:
`scripts/audit-operations-route-inventory.sh` passed.
`OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true scripts/audit-operations-route-inventory.sh`
passed. `go test ./cmd/agency-config -run
'OperationsLegacyPrivatePagesUseNoStore|OperationsConsoleNavigation|OperationsRouteTitles|OperationsConsoleNavigationActiveState'`
passed. `git diff --check` passed. `make check` passed.
`make audit-product-acceptance` passed. `make audit-final-claim-review`
passed. `git status --short -- docs/evidence/consumer-submissions
docs/evidence/captured db/migrations go.mod go.sum` returned clean.

Blocked checks:
Full closeout validation, `make validate`, `make test`, and Docker Compose
config are deferred to Phase 91 closeout.

Protected path status:
No protected evidence path was modified.

Consumer tracker status:
No consumer tracker edit was made. All seven targets remain required to stay
exactly `prepared`.

Claim-boundary status:
Copy changes preserve private/local diagnostic wording and explicitly avoid
claiming compliance, adoption, consumer action, final-root readiness,
hosted-service availability, production readiness, release readiness, vendor
compatibility, hardware certification, SLA/uptime, public launch, or ETA
quality.

Security/auth status:
No public route, admin auth expansion, credential path, external call, or
browser command was added. Legacy private generic pages now set no-store cache
headers consistently with newer private routes.

Data/migration status:
No persistence or migration change.

Master review:
Approved. The highest-priority IA drift and private cache-header gap are
patched without expanding product claims or protected paths.

Required edits:
Run Phase 91 closeout validation and create the Phase 91 handoff.

Decision:
Proceed to CP000005 closeout.

Next checkpoint:
Phase 91 -- Checkpoint 000005: close route product audit.

## Checkpoint 000005 Closeout

Checkpoint:
Phase 91 -- Checkpoint 000005: close route product audit.

Sub-agents used or simulated, including intended model level:
Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, real; Planning Sub-Agent --
GPT-5.5 x-high, real; Claim-Boundary/Security Sub-Agent -- GPT-5.5 high,
real; Implementation Sub-Agent -- GPT-5.5 high, simulated by Master; QA
Sub-Agent -- GPT-5.5 high, simulated by Master; UI/UX Sub-Agent -- GPT-5.5
high, simulated by Master; Documentation / IA Sub-Agent -- GPT-5.5 high,
simulated by Master; Data/Migration Sub-Agent -- GPT-5.5 high, simulated by
Master; Master Agent -- GPT-5.5 x-high, current thread.

Changed files:
`Makefile`; `README.md`; `cmd/agency-config/main_test.go`;
`cmd/agency-config/operations.go`; `docs/current-status.md`;
`docs/handoffs/latest.md`; `docs/handoffs/phase-91.md`;
`docs/open-transit-rt-master-planner-remaining-work.md`;
`docs/phase-91-maintainer-route-product-audit.md`;
`docs/roadmap-status.md`;
`docs/roadmaps/post-90-agency-grade-gtfs-rt-product/**`;
`docs/tutorials/no-cli-agency-first-run.md`;
`scripts/audit-operations-route-inventory.sh`;
`scripts/test-operations-route-inventory.sh`; `wiki/README.md`;
`wiki/browser-first-setup.md`; `wiki/operations-console-tour.md`;
`wiki/small-agency-quick-start.md`.

Validation run:
`git status --short` passed before closeout docs updates. `git diff --check`
passed. `make check` passed. `scripts/audit-operations-route-inventory.sh`
passed. `OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true
scripts/audit-operations-route-inventory.sh` passed.
`scripts/test-operations-route-inventory.sh` passed.
`make test-operations-route-inventory` passed.
`go test ./cmd/agency-config -run
'OperationsLegacyPrivatePagesUseNoStore|OperationsConsoleNavigation|OperationsRouteTitles|OperationsConsoleNavigationActiveState'`
passed. `make audit-product-acceptance` passed.
`make audit-final-claim-review` passed. `python3 -m json.tool
docs/evidence/consumer-submissions/status.json >/dev/null` passed. The exact
seven-target prepared-only consumer tracker check passed. `git status --short
-- docs/evidence/consumer-submissions docs/evidence/captured db/migrations
go.mod go.sum` returned clean. `make validate` passed. `make test` passed.
`docker compose -f deploy/docker-compose.yml config` passed.

Blocked checks:
None for Phase 91 closeout.

Protected path status:
No protected evidence path was modified.

Consumer tracker status:
All seven targets remain exactly `prepared`: Google Maps, Apple Maps, Transit
App, Bing Maps, Moovit, Mobility Database, and transit.land.

Claim-boundary status:
Phase 91 made no compliance, adoption, consumer action, final-root readiness,
hosted-service availability, production-readiness, release-ready, vendor,
hardware, SLA/uptime, public-launch, or ETA-quality claim.

Security/auth status:
No public admin route, auth role expansion, browser command route, credential
path, external call, raw private output exposure, or evidence path was added.
Legacy private generic pages now set no-store cache headers.

Data/migration status:
No persistence or migration change.

Master review:
Approved. Phase 91 is complete and safe to close.

Required edits:
None.

Decision:
Phase 91 is complete.

Next checkpoint:
Phase 92 -- Checkpoint 000001: add clean checkout rc gate plan.
