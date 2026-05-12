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
- The current default next work is `v0.1.0-rc.1` agency evaluation release-candidate hardening.
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
Sub-agents used or simulated:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Decision:
Next checkpoint:
```

## 5. Immediate patch before Phase 72

### Phase 71 -- Checkpoint 000005: tighten adoption path labels and phase traceability

Purpose:

- Standardize the first-click label across README/wiki/docs/UI guidance.
- Use one label everywhere: `Agency Operations Cockpit / Start Here`.
- Add a traceability note to `docs/handoffs/phase-69.md` explaining that Phase 69 product work appears to have landed as a bundled implementation commit even though the handoff records a checkpoint ledger. This is a process note, not a product rejection.
- Do not weaken Phase 69 closeout.

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

### Phase 72 — v0.1.0-rc.1 Release Candidate Hardening

Goal: prove from a clean checkout that the project can be evaluated as a bounded release candidate.

This is not a final public release, production-readiness claim, compliance claim, or evidence phase.

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

Must cover:

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
  - `/admin/operations/telemetry-simulator`
  - `/admin/operations/connectors`
  - `/admin/operations/maintenance`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- Docker Compose config.

Deliverables:

- `docs/phase-72-v0.1.0-rc.1-release-candidate-hardening.md`
- `docs/handoffs/phase-72.md`
- updates to `docs/release-candidate-readiness.md`
- updates to `docs/roadmap-status.md`
- updates to `docs/handoffs/latest.md`
- release notes draft for `v0.1.0-rc.1`
- known blockers matrix.

Success condition:

> A maintainer can say the repo has a repeatable, bounded `v0.1.0-rc.1` release-candidate evaluation path, with clear blockers, local UI walkthrough, connector checks, and no stronger external claims.

### Phase 73 — Agency UI Acceptance And Documentation Freeze

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

### Phase 74 — Connector Maturity And Adapter Recipes

Goal: make external connections easier to understand and test without claiming real vendor compatibility.

Checkpoint sequence:

```text
Phase 74 -- Checkpoint 000001: add connector maturity plan
Phase 74 -- Checkpoint 000002: improve telemetry connector recipes
Phase 74 -- Checkpoint 000003: improve prediction sidecar shadow/fail-closed docs
Phase 74 -- Checkpoint 000004: improve validator and monitoring/export connector recipes
Phase 74 -- Checkpoint 000005: expand adapter conformance coverage
Phase 74 -- Checkpoint 000006: close connector maturity review
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

- expanded `wiki/connector-cookbook.md` or `docs/tutorials/connector-recipes.md`;
- updated `docs/integration-adapter-kit.md`;
- updated examples README;
- new conformance fixtures if gaps are found.

Boundaries:

- no real vendor payloads;
- no credentials;
- no named vendor compatibility claims.

### Phase 75 — Deployment Operations Hardening For Small Hosts

Goal: make small-host/reference deployment review more realistic and clear.

Checkpoint sequence:

```text
Phase 75 -- Checkpoint 000001: add small-host operations hardening plan
Phase 75 -- Checkpoint 000002: harden off-host validation docs and scripts
Phase 75 -- Checkpoint 000003: harden backup and restore-drill operator guidance
Phase 75 -- Checkpoint 000004: harden upgrade and rollback guidance
Phase 75 -- Checkpoint 000005: harden maintenance center docs and UI guidance
Phase 75 -- Checkpoint 000006: close small-host operations review
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

### Phase 76 — Public Explanation And Docs/Site Freeze

Goal: align README, wiki, docs, and public explainer site so outside readers receive the same bounded product story.

Checkpoint sequence:

```text
Phase 76 -- Checkpoint 000001: add public explanation freeze plan
Phase 76 -- Checkpoint 000002: audit README wiki docs and gh-pages alignment
Phase 76 -- Checkpoint 000003: patch public explainer text and links
Phase 76 -- Checkpoint 000004: patch visual asset guidance and screenshot captions
Phase 76 -- Checkpoint 000005: close public explanation review
```

Must preserve:

- GitHub Pages is documentation only.
- Screenshots are local/demo docs only.
- No public launch claim.
- No compliance/adoption/consumer/production/vendor/SLA/ETA claim.

### Phase 77 — v0.1.0-rc.1 Candidate Cut

Goal: prepare a source release candidate after Phase 72-76 hardening.

Checkpoint sequence:

```text
Phase 77 -- Checkpoint 000001: add rc1 candidate cut plan
Phase 77 -- Checkpoint 000002: run final rc1 verification matrix
Phase 77 -- Checkpoint 000003: prepare rc1 release notes and source package
Phase 77 -- Checkpoint 000004: audit package, claims, docs, and protected paths
Phase 77 -- Checkpoint 000005: close rc1 candidate cut
```

Deliverables:

- `docs/releases/v0.1.0-rc.1.md` or release notes draft;
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
make release-package
make audit-release-package
docker compose -f deploy/docker-compose.yml config
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

### Phase 78+ — Optional Authorized Evidence Tracks

Only after RC hardening, and only with explicit written authorization.

Possible tracks:

- agency-owned final-root evidence;
- target-specific consumer submission evidence;
- real agency pilot closeout;
- real vendor/device AVL proof;
- real-world realtime quality study;
- production deployment evidence;
- compliance evidence packet for one specific deployment.

Do not start any Phase 78+ evidence track unless the maintainer supplies:

```text
exact claim target
allowed tools
public-safe retention rules
redaction rules
stop conditions
operator/agency authorization
```

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
Phases 0-60 are closed. Phase 61+ roadmap naming is approved. Phases 61-67 are complete. Phase 68+ is closed blocker-only / authorization-gated. Phase 69, Phase 70, and Phase 71 are complete. Default next work is v0.1.0-rc.1 release candidate hardening, not evidence intake.

Before Phase 72, perform:
Phase 71 -- Checkpoint 000005: tighten adoption path labels and phase traceability

Then start:
Phase 72 -- v0.1.0-rc.1 Release Candidate Hardening

Use the master/sub-agent workflow:
- Context / Repo Truth sub-agent
- Planning sub-agent
- Implementation sub-agent
- QA sub-agent
- UI/UX sub-agent
- Documentation / IA sub-agent
- Claim-Boundary sub-agent

If real sub-agents are unavailable, simulate those roles in labeled sections.

Do not implement until the master approves the checkpoint plan.

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

Primary Phase 72 goal:
Prove that a clean checkout can run a bounded release-candidate evaluation: make check, make validate, make test, local app startup, five local public feed fetches, browser-first Operations Console walkthrough, connector/adaptor conformance, off-host diagnostics dry-run, product acceptance audit, and final claim audit.

Required final validation:
git status --short
git diff --check
make check
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-product-acceptance
make audit-final-claim-review
docker compose -f deploy/docker-compose.yml config
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured

Report after each checkpoint:
Checkpoint:
Sub-agents used or simulated:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Decision:
Next checkpoint:
```

## 9. Master planner verdict

The project should now stop broad planning and move into disciplined release-candidate execution.

The product direction is correct:

- browser-first agency operations;
- simple README/wiki/docs;
- connector maturity;
- clean install;
- release-candidate diagnostics;
- claim discipline.

The remaining risk is not lack of roadmap. The remaining risk is incomplete release-candidate proof from a clean checkout.

The next master-agent action is:

```text
Patch Phase 71 label/traceability consistency.
Then start Phase 72 -- v0.1.0-rc.1 Release Candidate Hardening.
```
