# Open Transit RT — Master Planner Export For Remaining Product Work

**Generated:** 2026-05-11  
**Prepared for:** maintainer / next Codex master agent  
**Repository:** `ptse8204/open-transit-rt`

## 1. Why this file exists

Earlier responses evaluated the plan and current implementation, but did not export a standalone Markdown handoff. That was a miss.

The first instruction in this workstream was to make the Codex instructions expansive enough to implement, ground them in the current repo and project history, keep the work cohesive toward the long-term product vision, clean the README/wiki/docs, add UI, and improve the user experience. This export turns that into a usable master-planner artifact.

## 2. Current truth to preserve

- Phases 0-60 remain closed.
- Phase 61+ roadmap naming is approved.
- Phases 61-67 are complete.
- Phase 68+ is closed blocker-only / authorization-gated.
- Phase 69 is complete for maintainer product acceptance and UI-first agency usability.
- Phase 70 is complete for the GitHub Pages product explainer site.
- Phase 71 is complete for adoption-first productization and no-CLI agency operations.
- Phase 72 is complete for bounded `v0.1.0-rc.1` agency evaluation release-candidate hardening review.
- Phase 72 ended with `needs_review` release-candidate diagnostics, not a release-ready pass.
- Phase 73 Checkpoint 000001 is complete for documentation-only agency UI acceptance planning.
- Phase 73 Checkpoint 000002 is complete for local no-developer browser walkthrough review.
- Phase 73 Checkpoint 000003 is complete for local technical-helper walkthrough review.
- Phase 73 Checkpoint 000004 is complete for narrow UI copy, route-label, Devices/Telemetry boundary-copy, and browser-first tutorial patching.
- Phase 73 Checkpoint 000005 is complete for small-agency docs and wiki navigation freeze.
- Phase 73 Checkpoint 000006 is complete for bounded agency UI acceptance closeout.
- Phase 74 Checkpoints 000001 through 000008 are complete for GitHub Pages and agency UI product polish.
- Phase 74 Checkpoint 000008 reconciled and published the actual `gh-pages` branch with the Phase 74 closeout.
- GitHub Pages product story is refreshed.
- Private Operations Console first-run hierarchy is improved.
- Docs/site/UI now point to the same browser-first product path.
- Phase 75 is complete for the authorized Consumer-Grade Control Plane roadmap pack.
- Phase 76 is complete for Design System And App Shell.
- Phase 77 is complete for the private Admin Control API And Command Model scope.
- Phase 78 is complete for Frontend Routing, State, And Data Loading.
- Phase 79 is complete for Agency Setup V3.
- Phase 80 is complete for GTFS Workbench.
- Phase 81 is complete for Realtime Operations Center.
- Phase 82 is complete for Feed Health And Validation Center.
- The previous Phase 74 connector-maturity slot is postponed to a later separately authorized phase.
- Evidence/adoption/compliance tracks remain optional and authorization-gated.

Do not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- final-root readiness;
- hosted SaaS availability;
- paid support or SLA;
- production readiness;
- vendor compatibility or hardware certification;
- production-grade ETA quality.

All seven consumer / aggregator targets must remain exactly `prepared` unless retained, redacted, target-originated evidence supports a specific change:

1. Google Maps
2. Apple Maps
3. Transit App
4. Bing Maps
5. Moovit
6. Mobility Database
7. transit.land

## 3. Original goal restated

Open Transit RT should become a self-hosted, open-source GTFS / GTFS-Realtime operations platform that a small agency, civic technologist, or developer integrator can start from a clean checkout, open in a browser, import or author GTFS, publish GTFS and all three GTFS-Realtime feed types, connect telemetry through documented connector paths, monitor feed health, review CAL-ITP-style readiness, and understand exactly what remains before public deployment or stronger claims.

The key strategy is:

> Build great software first. Treat real external proof as optional, authorization-gated work.

## 4. Master-agent operating model

The next Codex instance should act as the **Master Agent**.

If real spawned sub-agents are available, use them. If not, simulate the roles in clearly labeled sections.

### Model Assignment

Use these model levels for the master/sub-agent workflow:

| Role | Model level |
| --- | --- |
| Master Agent | GPT-5.5 x-high |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high |
| Planning Sub-Agent | GPT-5.5 x-high |
| Implementation Sub-Agent | GPT-5.5 high |
| QA Sub-Agent | GPT-5.5 high |
| UI/UX Sub-Agent | GPT-5.5 high |
| Documentation / Information Architecture Sub-Agent | GPT-5.5 high |
| Claim-Boundary Sub-Agent | GPT-5.5 high |

If Codex can spawn real sub-agents, assign those model levels to the corresponding agents.

If Codex cannot spawn real sub-agents, simulate the roles in clearly labeled sections and still label each simulated role with the intended model level.

The Master Agent must approve the plan before implementation starts and must approve every checkpoint after reviewing all sub-agent reports. The Master Agent may move forward only when no required edits remain.

### Master Agent

Owns scope, repo truth, checkpoint sequencing, evidence boundaries, claim boundaries, and final approval.

### Context / Repo Truth Sub-Agent

Reads:

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`
- `docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md`
- `docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md`
- relevant phase docs and handoffs.

Reports current phase state, protected paths, consumer tracker state, UI routes, and validation gates.

### Planning Sub-Agent

Drafts checkpoint-level plans with files to read, files to edit, files not to edit, tests, stop conditions, and forbidden claims.

### Implementation Sub-Agent

Executes only approved checkpoint scope. No scope expansion.

### UI/UX Sub-Agent

Checks whether a small-agency operator can understand and use the browser-first workflow without phase-history knowledge or heavy CLI dependence.

### Documentation / Information Architecture Sub-Agent

Reviews README, wiki, docs, tutorials, connector docs, deployment docs, release-candidate docs, and public explainer alignment.

### QA Sub-Agent

Runs validation, records blockers, checks protected paths, and verifies the prepared-only consumer tracker.

### Claim-Boundary Sub-Agent

Blocks unsupported public claims and verifies all claim flags remain false.

### Required checkpoint report

After every checkpoint, the master must report:

```text
Checkpoint:
Sub-agents used or simulated, including intended model level:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Master review:
Required edits:
Decision:
Next checkpoint:
```

## 5. Completed immediate patch before Phase 72

### Phase 71 -- Checkpoint 000005: tighten adoption path labels and phase traceability

Status: complete.

Purpose:

- Standardize the first-click label across README/wiki/docs/UI guidance.
- Use one label everywhere: `Agency Operations Cockpit / Start Here`.
- Add a traceability note to `docs/handoffs/phase-69.md` explaining that Phase 69 product work appears to have landed as a bundled implementation commit even though the handoff records a checkpoint ledger. This is a process note, not a product rejection.
- Do not weaken Phase 69 closeout.

Completion note:

- The first-click label is standardized as `Agency Operations Cockpit / Start Here` in the reviewed README/wiki/docs/UI guidance.
- The Phase 69 traceability note is explicitly a process traceability note, not a product rejection, and does not weaken Phase 69 closeout.
- Protected evidence paths, consumer statuses, evidence references, route names, filenames, internal identifiers, and claim boundaries remain unchanged.

Likely files:

- `README.md`
- `wiki/README.md`
- `docs/README.md`
- `wiki/small-agency-quick-start.md`
- `wiki/browser-first-setup.md`
- `wiki/operations-console-tour.md`
- `docs/tutorials/small-agency-acceptance-script.md`
- `docs/handoffs/phase-69.md`
- `docs/handoffs/latest.md`, only if needed.

Validation:

```bash
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
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
```

## 6. Remaining roadmap

### Completed Phase 72 — v0.1.0-rc.1 Release Candidate Hardening

Completed result: the project now has a bounded `v0.1.0-rc.1` release-candidate evaluation path, with final closeout still `needs_review`.

This is not a final public release, release-ready pass, production-readiness claim, compliance claim, or evidence phase.

Remaining `needs_review` items: dirty primary checkout, release package audit `not_checked`, no tag/package/published image, no retained evidence, and the CP000004 browser automation limitation. Consumer tracker statuses remain unchanged as a protected boundary, not a `needs_review` item.

Checkpoint sequence:

```text
Phase 72 -- Checkpoint 000001: add release-candidate hardening plan
Phase 72 -- Checkpoint 000002: run clean-checkout local evaluator gate
Phase 72 -- Checkpoint 000003: harden release-candidate diagnostics and blockers
Phase 72 -- Checkpoint 000004: verify browser-first agency operations walkthrough
Phase 72 -- Checkpoint 000005: verify connector and adapter conformance gates
Phase 72 -- Checkpoint 000006: prepare rc1 release notes and known blockers
Phase 72 -- Checkpoint 000007: close rc1 hardening review
```

Covered:

- `make check`
- `make validate`
- `make test`
- `RUN_LOCAL_APP=true make release-candidate-check`
- local app startup
- all five local public feed fetches
- browser walkthrough:
  - `/admin/operations`
  - `/admin/operations/gtfs-import`
  - `/admin/operations/feed-health`
  - `/admin/operations/gtfs-quality`
  - `/admin/operations/validation-health`
  - `/admin/operations/devices`
  - `/admin/operations/telemetry`
  - `/admin/operations/telemetry-simulator`
  - `/admin/operations/connectors`
  - `/admin/operations/maintenance`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- Docker Compose config.

Deliverables produced/updated:

- `docs/phase-72-v0.1.0-rc.1-release-candidate-hardening.md`
- `docs/handoffs/phase-72.md`
- updates to `docs/release-candidate-readiness.md`
- updates to `docs/roadmap-status.md`
- updates to `docs/handoffs/latest.md`
- `docs/release-notes-v0.1.0-rc.1-draft.md`
- known blockers matrix.

Success condition:

> A maintainer can say the repo has a repeatable, bounded `v0.1.0-rc.1` release-candidate evaluation path, with clear blockers, local UI walkthrough, connector checks, and no stronger external claims.

### Phase 73 — Agency UI Acceptance And Documentation Freeze

Status: complete after Checkpoint 000006.

Goal: freeze the browser-first agency path so a non-expert evaluator can follow it without maintainer narration.

Checkpoint sequence:

```text
Phase 73 -- Checkpoint 000001: add agency UI acceptance plan
Phase 73 -- Checkpoint 000002: run no-developer browser walkthrough
Phase 73 -- Checkpoint 000003: run technical-helper walkthrough
Phase 73 -- Checkpoint 000004: patch UI copy and empty states
Phase 73 -- Checkpoint 000005: freeze small-agency docs and wiki navigation
Phase 73 -- Checkpoint 000006: close agency UI acceptance review
```

Must cover:

- Agency Operations Cockpit / Start Here.
- No-developer path.
- Developer / technical-helper path.
- GTFS import in browser.
- Feed health.
- GTFS quality.
- Validation health.
- Device / telemetry readiness.
- Telemetry simulator guidance.
- Connector Hub.
- Maintenance Center.
- Help.

Deliverables:

- updated no-CLI first-run guide;
- updated small-agency maintenance guide;
- updated operations console tour;
- updated acceptance checklist for UI pages;
- screenshot guidance only if images are real local/demo captures and not evidence.

Completion note: CP000006 closed Phase 73 as a bounded documentation/status acceptance closeout after CP000001 through CP000005. The browser-first agency path has no remaining Phase 73 required edits after the local walkthroughs, UI/docs copy patch, and docs/wiki navigation freeze. CP000006 did not rerun the local app, create retained evidence, move consumer statuses, tag, package, publish, or claim release readiness.

### Phase 74 — GitHub Pages And Agency UI Product Polish

Completed result: Phase 74 refreshed the public GitHub Pages
documentation-only product story, refreshed the GitHub Pages quickstart and UI
tour, improved the private Operations Console first-run hierarchy, improved
empty states and next actions, aligned README/docs/wiki/site navigation around
the same browser-first product path, reconciled and published the actual
`gh-pages` branch, and closed the product-polish review.

Closeout boundary: no retained evidence was created, no external party was
contacted, no consumer status changed, and no
compliance/adoption/consumer/final-root/SaaS/production/vendor/SLA/ETA claim
was added. Connector maturity was explicitly postponed to a later separately
authorized phase.

Phase 74 -- Checkpoint 000001: add site and agency UI polish plan
Phase 74 -- Checkpoint 000002: refresh GitHub Pages product story
Phase 74 -- Checkpoint 000003: refresh GitHub Pages UI tour and quickstart
Phase 74 -- Checkpoint 000004: improve Operations Console visual hierarchy
Phase 74 -- Checkpoint 000005: improve first-run empty states and next actions
Phase 74 -- Checkpoint 000006: align docs README wiki and site navigation
Phase 74 -- Checkpoint 000007: close site and UI product polish review
Phase 74 -- Checkpoint 000008: reconcile and publish GitHub Pages browser-first refresh

### Future Phase — Connector Maturity And Adapter Recipes

Goal: make external connections easier to understand and test without claiming real vendor compatibility.

Checkpoint sequence:

```text
Future connector maturity -- Checkpoint 000001: add connector maturity plan
Future connector maturity -- Checkpoint 000002: improve telemetry connector recipes
Future connector maturity -- Checkpoint 000003: improve prediction sidecar shadow/fail-closed docs
Future connector maturity -- Checkpoint 000004: improve validator and monitoring/export connector recipes
Future connector maturity -- Checkpoint 000005: expand adapter conformance coverage
Future connector maturity -- Checkpoint 000006: close connector maturity review
```

Recipes should cover:

- “I have a CSV of vehicle locations.”
- “I have a GPS API.”
- “I have an AVL vendor that can POST data.”
- “I want synthetic telemetry only.”
- “I want an external prediction engine.”
- “I want monitoring summaries.”
- “I want to verify public feed URLs off-host.”

Each recipe should include:

```text
What you have
Which connector path to use
What maps to /v1/telemetry or prediction adapter
How to test with synthetic/local data
What success looks like
What this does not prove
```

Deliverables:

- expanded existing `wiki/connector-cookbook.md`;
- updated `docs/integration-adapter-kit.md`;
- updated examples README;
- new conformance fixtures if gaps are found.

Boundaries:

- no real vendor payloads;
- no credentials;
- no named vendor compatibility claims.

### Superseded numbering note for old future slots

The following small-host, public-explanation, release-cut, and evidence themes
were drafted before the maintainer authorized Phase 75 as the
Consumer-Grade Control Plane Roadmap Pack. They are not active phase
assignments. Treat them as postponed backlog themes that may be remapped by
`docs/roadmaps/consumer-grade-control-plane/README.md` or by a later
separately authorized phase.

### Postponed backlog theme — Small-Host Deployment Operations Hardening

Goal: make small-host/reference deployment review more realistic and clear.

Possible checkpoint sequence if later re-authorized:

```text
Future small-host operations -- Checkpoint 000001: add small-host operations hardening plan
Future small-host operations -- Checkpoint 000002: harden off-host validation docs and scripts
Future small-host operations -- Checkpoint 000003: harden backup and restore-drill operator guidance
Future small-host operations -- Checkpoint 000004: harden upgrade and rollback guidance
Future small-host operations -- Checkpoint 000005: harden maintenance center docs and UI guidance
Future small-host operations -- Checkpoint 000006: close small-host operations review
```

Focus:

- tiny-server validator blockers;
- off-host validation;
- backup configuration;
- restore-drill configuration;
- upgrade rollback;
- service health diagnostics;
- support-bundle redaction;
- maintenance cadence.

Deliverables:

- improved `docs/deployment/off-host-validation.md`;
- improved `docs/deployment/oci-reference-check.md`;
- improved `docs/tutorials/small-agency-maintenance-guide.md`;
- script patches only if actual gaps are found.

### Postponed backlog theme — Public Explanation And Docs/Site Freeze

Goal: align README, wiki, docs, and public explainer site so outside readers receive the same bounded product story.

Possible checkpoint sequence if later re-authorized:

```text
Future public explanation freeze -- Checkpoint 000001: add public explanation freeze plan
Future public explanation freeze -- Checkpoint 000002: audit README wiki docs and gh-pages alignment
Future public explanation freeze -- Checkpoint 000003: patch public explainer text and links
Future public explanation freeze -- Checkpoint 000004: patch visual asset guidance and screenshot captions
Future public explanation freeze -- Checkpoint 000005: close public explanation review
```

Must preserve:

- GitHub Pages is documentation only.
- Screenshots are local/demo docs only.
- No public launch claim.
- No compliance/adoption/consumer/production/vendor/SLA/ETA claim.

### Postponed backlog theme — v0.1.0-rc.1 Candidate Cut

Goal: prepare a source release candidate after Phase 72-76 hardening.

Possible checkpoint sequence if later re-authorized:

```text
Future rc1 candidate cut -- Checkpoint 000001: add rc1 candidate cut plan
Future rc1 candidate cut -- Checkpoint 000002: run final rc1 verification matrix
Future rc1 candidate cut -- Checkpoint 000003: prepare rc1 release notes and source package
Future rc1 candidate cut -- Checkpoint 000004: audit package, claims, docs, and protected paths
Future rc1 candidate cut -- Checkpoint 000005: close rc1 candidate cut
```

Deliverables:

- release notes draft from existing `docs/release-notes-template.md`; use
  `docs/release-notes-v0.1.0-rc.1-draft.md` until a future release checkpoint
  approves a different tagged-release path;
- local `.cache` source package through existing release package tooling;
- checksum manifest;
- SBOM/provenance metadata where existing tools support it;
- clear “not production-ready proof” boundary.

Validation:

```bash
git status --short
git diff --check
make check
make validate
make test
make release-candidate-check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-product-acceptance
make audit-final-claim-review
# only with explicit maintainer authorization; otherwise record blocked/not_checked
make release-package
make audit-release-package
docker compose -f deploy/docker-compose.yml config
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

### Postponed backlog theme — Optional Authorized Evidence Tracks

Only after RC hardening, and only with explicit written authorization.

Possible tracks:

- agency-owned final-root evidence;
- target-specific consumer submission evidence;
- real agency pilot closeout;
- real vendor/device AVL proof;
- real-world realtime quality study;
- production deployment evidence;
- compliance evidence packet for one specific deployment.

Do not start any optional evidence track unless the maintainer supplies:

```text
exact claim target
allowed tools
public-safe retention rules
redaction rules
stop conditions
operator/agency authorization
```

Consumer status movement requires retained, redacted, target-originated
evidence for the named target and exact feed scope. Operator authorization
alone is not enough to move a target beyond `prepared`.

## 7. What not to do

Do not:

- chase real agency evidence by default;
- start consumer submissions;
- contact agencies, vendors, consumers, portals, or external systems;
- move consumer statuses beyond `prepared`;
- write under `docs/evidence/captured/**`;
- write protected consumer submission artifacts;
- claim compliance, adoption, consumer acceptance, final-root readiness, hosted SaaS, production readiness, vendor compatibility, SLA, uptime, or production-grade ETA quality;
- introduce a heavy frontend stack;
- replace the server-rendered admin UI with a speculative SPA;
- bury the quickstart under phase history;
- make docs impressive but unusable;
- let Codex mark a phase complete without repo-file inspection.

## 8. Copy-paste prompt for Codex Master Agent

```text
You are the master agent for Open Transit RT.

Start by reading:
- AGENTS.md
- README.md
- docs/current-status.md
- docs/handoffs/latest.md
- docs/roadmap-status.md
- docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md
- docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md
- docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md
- docs/handoffs/phase-69.md
- docs/handoffs/phase-71.md
- docs/phase-71-adoption-first-productization-no-cli-agency-operations.md
- docs/evidence/redaction-policy.md
- docs/evidence/consumer-submissions/status.json

Current truth:
Phases 0-60 are closed. Phase 61+ roadmap naming is approved. Phases 61-67 are complete. Phase 68+ is closed blocker-only / authorization-gated. Phase 69, Phase 70, Phase 71, and Phase 72 are complete. Phase 72 closed bounded v0.1.0-rc.1 hardening review with `needs_review` diagnostics, not release readiness. Phase 73 Checkpoints 000001 through 000006 are complete for agency UI acceptance. Phase 74 Checkpoints 000001 through 000008 are complete for GitHub Pages and agency UI product polish; CP000008 reconciled and published the actual `gh-pages` branch with the Phase 74 closeout. Phase 75 is complete for the Consumer-Grade Control Plane roadmap pack. Phase 76 is complete for Design System And App Shell. Phase 77 is complete for the private Admin Control API And Command Model scope. Phase 78 is complete for Frontend Routing, State, And Data Loading. Phase 79 is complete for Agency Setup V3. Phase 80 is complete for GTFS Workbench. Phase 81 is complete for Realtime Operations Center. Phase 82 is complete for Feed Health And Validation Center. Release-cut cleanup, postponed connector maturity, and optional evidence tracks remain separated by their phase gates. The default next work is not evidence intake.

Model assignment:
Use these model levels for the master/sub-agent workflow:

| Role | Model level |
| --- | --- |
| Master Agent | GPT-5.5 x-high |
| Context / Repo Truth Sub-Agent | GPT-5.5 x-high |
| Planning Sub-Agent | GPT-5.5 x-high |
| Implementation Sub-Agent | GPT-5.5 high |
| QA Sub-Agent | GPT-5.5 high |
| UI/UX Sub-Agent | GPT-5.5 high |
| Documentation / Information Architecture Sub-Agent | GPT-5.5 high |
| Claim-Boundary Sub-Agent | GPT-5.5 high |

If Codex can spawn real sub-agents, assign those model levels to the corresponding agents.

If Codex cannot spawn real sub-agents, simulate the roles in clearly labeled sections and still label each simulated role with the intended model level.

The Master Agent must approve the plan before implementation starts and must approve every checkpoint after reviewing all sub-agent reports. The Master Agent may move forward only when no required edits remain.

Phase 71 -- Checkpoint 000005 is complete. Phase 72 -- v0.1.0-rc.1 Release Candidate Hardening is complete for bounded review with `needs_review` diagnostics. Phase 73 -- Checkpoint 000001 is complete for documentation-only agency UI acceptance planning. Phase 73 -- Checkpoint 000002 is complete for local no-developer browser walkthrough review. Phase 73 -- Checkpoint 000003 is complete for local technical-helper walkthrough review. Phase 73 -- Checkpoint 000004 is complete for narrow UI copy, route-label, Devices/Telemetry boundary-copy, and browser-first tutorial patching. Phase 73 -- Checkpoint 000005 is complete for small-agency docs and wiki navigation freeze. Phase 73 -- Checkpoint 000006 is complete for bounded agency UI acceptance closeout. Phase 74 -- Checkpoint 000001 through Checkpoint 000008 are complete for GitHub Pages and agency UI product polish. Phase 75 through Phase 82 are complete for the authorized Consumer-Grade Control Plane product track through the private Feed Health And Validation Center. Continue with Phase 83 -- Connector Workbench, while keeping release-cut cleanup/release-candidate gating, connector maturity claims, and optional evidence tracks separated by their phase gates.

Use the master/sub-agent workflow gate:
- Context / Repo Truth sub-agent
- Planning sub-agent
- Implementation sub-agent
- QA sub-agent
- UI/UX sub-agent
- Documentation / IA sub-agent
- Claim-Boundary sub-agent

If real sub-agents are unavailable, simulate those roles in labeled sections.

Do not implement until the Master Agent approves the checkpoint plan. The Implementation Sub-Agent executes only approved checkpoint scope.

Protected paths:
- docs/evidence/captured/**
- docs/evidence/consumer-submissions/status.json
- docs/evidence/consumer-submissions/current/**
- docs/evidence/consumer-submissions/artifacts/**
- docs/evidence/consumer-submissions/packets/**

Consumer tracker:
All seven targets must remain exactly prepared:
Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land.

Forbidden claims:
No CAL-ITP/Caltrans compliance, agency adoption/approval, consumer submission, consumer acceptance, consumer ingestion/listing/display, final-root readiness, hosted SaaS, production readiness, vendor compatibility, hardware certification, SLA/uptime, or production-grade ETA quality.

Primary Phase 73 CP000006 result:
Closed the agency UI acceptance review by recording the final acceptance result, remaining blockers, validation status, protected path review, consumer tracker boundary, and exact next recommendation after CP000005 docs/wiki navigation freeze. Do not treat Phase 72 as a release-ready pass; clean-checkout release-cut proof remains future release-cut cleanup unless separately authorized.

Validation recorded for Phase 73 CP000006:
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured

Report after each checkpoint:
Checkpoint:
Sub-agents used or simulated, including intended model level:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Master review:
Required edits:
Decision:
Next checkpoint:
```

## 9. Master planner verdict

Phase 74 GitHub Pages and agency UI product-polish closeout is complete
through CP000008. The project should now stop broad planning and move into
maintainer review of the Phase 74 CP000008 closeout before any separately
authorized release-cut cleanup, postponed connector maturity, or future product
phase.

The product direction is correct:

- browser-first agency operations;
- simple README/wiki/docs;
- postponed connector maturity;
- clean install;
- completed release-candidate diagnostics with `needs_review` blockers;
- claim discipline.

The remaining risk is not lack of roadmap. Phase 82 records that the private
Validation Center path now combines five feed rows, validator health, GTFS
quality summary, sanitized issue drilldowns, readiness timeline, current
blockers, and prepared-only consumer tracker state. No retained evidence was
created, no external party was contacted, no consumer status changed, and no
compliance/adoption/consumer/final-root/SaaS/production/vendor/SLA/ETA claim
was added. Clean-checkout release-cut proof and connector maturity remain
separate future authorization.

The next master-agent action is:

```text
Continue the authorized Phase 75-90 product track with Phase 83 -- Connector
Workbench. Keep release-cut cleanup, connector maturity claims, and optional
evidence tracks separated by their phase gates.
```
