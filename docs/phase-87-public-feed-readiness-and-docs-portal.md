# Phase 87 -- Public Feed Readiness And Docs Portal

## Purpose

Phase 87 improves how private operators review public feed URLs, source-of-truth metadata, prepared consumer packet context, and docs portal guidance. The work is a private product/UI/docs phase in the authorized Phase 75-90 Consumer-Grade Control Plane track.

This phase must not collect final-root proof, write retained evidence, modify consumer packet artifacts, change consumer tracker statuses, contact external parties, publish release artifacts, or claim outside acceptance.

## Current Truth

- Phases 0-86 are complete for their bounded scopes.
- Phase 72 remains `needs_review`, not release-ready.
- Phase 74 CP000008 remains the latest GitHub Pages publication at commit `a8b250e`.
- Phase 86 closed with private agency-scope, access, audit, and accessibility hardening.
- Evidence/adoption/compliance tracks remain optional and require separate written authorization.
- All seven consumer targets remain exactly `prepared`: Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

## Sub-Agent Plan

- Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, simulated: confirm feed-health, readiness, setup, consumer, evidence, docs, and public feed route state before edits.
- Planning Sub-Agent -- GPT-5.5 x-high, simulated: keep checkpoints private, browser-first, evidence-separated, and route-stable.
- UI/UX Sub-Agent -- GPT-5.5 high, simulated: review whether a small-agency operator can copy, review, and explain feed URLs without confusing readiness with proof.
- Documentation / IA Sub-Agent -- GPT-5.5 high, simulated: align operator docs, docs portal pointers, screenshot/diagram guidance, and prepared-packet explanations.
- Claim-Boundary Sub-Agent -- GPT-5.5 high, simulated: block final-root, consumer, compliance, hosted-service, production, vendor, SLA, and ETA-quality claims.
- Security/Auth Sub-Agent -- GPT-5.5 high, simulated: preserve private Operations Console auth/role/CSRF boundaries and avoid exposing private paths or raw diagnostics.
- QA Sub-Agent -- GPT-5.5 high, simulated: run baseline validation, exact prepared-only tracker checks, and protected-path checks.
- Data/Migration Sub-Agent -- GPT-5.5 high, simulated only if persisted model changes are proposed. Default expectation: no migration.

## Master Approval Before Implementation

Approved bounded scope:

- Prefer private Operations Console improvements and documentation/static guidance.
- Preserve existing public feed routes and avoid new public dynamic surfaces unless separately justified and reviewed.
- Keep feed URL review and copy workflows informational unless an existing safe private helper already supports a bounded check.
- Keep prepared packet explanation strictly prepared-only.
- Keep final-root and evidence checklist items as future authorization gates, not current evidence collection.

Required edits before implementation: none.

## Checkpoints

### Checkpoint 000001 -- Plan

Deliver this plan and link it from current source-of-truth docs as the active Phase 87 plan.

Expected files:

- `docs/phase-87-public-feed-readiness-and-docs-portal.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

### Checkpoint 000002 -- Feed URL And Metadata Preview Guidance

Improve the private browser path for reviewing exactly the expected public feed URLs:

- `feeds.json`
- static schedule ZIP
- Vehicle Positions
- Trip Updates
- Alerts

Expected product behavior:

- show configured URL, public path, metadata source, license/contact status, validation/fetch context where already available, and copy/share guidance;
- explain what the URL record means and what it does not prove;
- no external fetch from the browser unless it is an existing private/local check;
- no final-root or consumer status claim.

Likely files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/main_test.go`

### Checkpoint 000003 -- Prepared-Only Consumer Packet Explanation

Improve the private explanation for prepared consumer packet records without touching protected packet/status files.

Expected product behavior:

- exactly seven prepared targets remain visible;
- copy says prepared means local packet readiness only;
- no submission, review, acceptance, listing, display, ingestion, or consumer outcome wording;
- link to protected tracker docs as read-only references only.

Likely files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_checklist.go`
- `cmd/agency-config/main_test.go`
- docs/tutorial or help files, if needed.

### Checkpoint 000004 -- Docs Portal, Screenshot, Diagram, Off-Host Guidance

Align docs portal guidance with the private feed readiness path:

- source-of-truth website/domain checklist;
- off-host validation guidance;
- screenshot and diagram capture guidance for future docs only;
- future final-root/evidence checklist that says it requires separate written authorization.

Expected files:

- docs/tutorials or deployment docs as needed;
- `wiki/` docs only if source-of-truth navigation needs a private-product link.

### Checkpoint 000005 -- Closeout

Write Phase 87 handoff and update current source-of-truth docs. Confirm:

- no protected path writes;
- exact prepared-only consumer tracker;
- no public admin routes;
- no new release artifacts;
- no evidence collection;
- no forbidden claims.

## Protected Paths

Do not modify or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

## Claim Boundaries

Phase 87 may say:

- private feed URL readiness review;
- configured feed URL copy/share guidance;
- source-of-truth metadata checklist;
- prepared packet records remain prepared only;
- future evidence gate requires separate written authorization.

Phase 87 must not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- final-root readiness;
- hosted SaaS;
- paid support;
- SLA or uptime;
- production readiness;
- vendor compatibility;
- hardware certification;
- production-grade ETA quality;
- real-world ETA accuracy;
- public launch completion.

## Validation

Run at least:

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

When code changes:

```bash
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

## Closeout Report Requirements

The closeout must include:

```text
Phase:
Sub-agents used or simulated, including intended model level:
Goal:
Changed files:
Routes added/changed:
Commands added/changed:
Migrations:
Validation run:
Blocked checks:
Known blockers:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Accessibility status:
Docs/site/wiki alignment:
Commit list:
Master review:
Required edits:
Decision:
Next phase:
```
