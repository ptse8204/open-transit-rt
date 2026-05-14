# Phase 88 -- Nontechnical Training And In-App Guidance

## Purpose

Phase 88 helps small-agency staff learn and operate Open Transit RT from the
private browser control plane without requiring a developer beside them. The
work is a private product/UI/docs phase in the authorized Phase 75-90
Consumer-Grade Control Plane track.

This phase must not collect evidence, contact external parties, publish
release artifacts, move consumer statuses, add public admin routes, automate
portal workflows, or claim outside acceptance.

## Current Truth

- Phases 0-87 are complete for their bounded scopes.
- Phase 72 remains `needs_review`, not release-ready.
- Phase 74 CP000008 remains the latest GitHub Pages publication at commit
  `a8b250e`.
- Phase 87 closed with private feed URL readiness, source-of-truth/off-host
  validation guidance, docs portal alignment guidance, prepared-only consumer
  packet explanation, and future authorization gates.
- Evidence/adoption/compliance tracks remain optional and require separate
  written authorization.
- All seven consumer targets remain exactly `prepared`: Google Maps, Apple
  Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

## Sub-Agent Plan

- Context / Repo Truth Sub-Agent -- GPT-5.5 x-high, simulated: confirm existing
  help, launchpad, setup, checklist, feed, telemetry, connector, maintenance,
  and docs surfaces before edits.
- Planning Sub-Agent -- GPT-5.5 x-high, simulated: keep checkpoints private,
  browser-first, nontechnical, route-stable, and evidence-separated.
- UI/UX Sub-Agent -- GPT-5.5 high, simulated: review whether director,
  operator, technical-helper, integrator, and no-developer evaluator paths are
  understandable without phase history.
- Documentation / IA Sub-Agent -- GPT-5.5 high, simulated: align glossary,
  quick tasks, troubleshooting decisions, printable training, and handoff
  material.
- Claim-Boundary Sub-Agent -- GPT-5.5 high, simulated: block compliance,
  adoption, consumer, final-root, hosted-service, production, vendor, hardware,
  SLA/uptime, ETA-quality, and public-launch claims.
- Security/Auth Sub-Agent -- GPT-5.5 high, simulated: preserve private
  Operations Console auth, role, CSRF, and no-secret/no-raw-output boundaries.
- QA Sub-Agent -- GPT-5.5 high, simulated: run focused help/training route
  tests, baseline validation, exact prepared-only tracker checks, and
  protected-path checks.
- Data/Migration Sub-Agent -- GPT-5.5 high, simulated only if persisted model
  changes are proposed. Default expectation: no migration.

## Master Approval Before Implementation

Approved bounded scope:

- Prefer private Operations Console help/training improvements and docs/static
  training material.
- Preserve Go server-rendered templates and buildless progressive enhancement.
- Keep training read-only unless it links to existing safe private workflows.
- Teach what to do next without adding evidence, portal, consumer, release, or
  public-launch actions.
- Avoid jargon-heavy copy in primary paths while keeping technical detail
  available through docs links.

Required edits before implementation: none.

## Checkpoints

### Checkpoint 000001 -- Plan

Deliver this plan and link it from current source-of-truth docs as the active
Phase 88 plan.

Expected files:

- `docs/phase-88-nontechnical-training-in-app-guidance.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

### Checkpoint 000002 -- Role-Based Tours And First-Week Checklist

Improve private in-app guidance for:

- no-developer evaluator;
- director or manager;
- daily operator;
- technical helper;
- integrator.

Expected product behavior:

- role paths explain what to review first, what to do next, when to ask for
  help, and what the path does not prove;
- a first-week checklist covers setup, GTFS import, feed health, validation,
  telemetry, connectors, prediction diagnostics, maintenance, and support
  bundle preparation;
- no evidence, consumer, release, portal, public-launch, or external-contact
  action is added.

Likely files:

- `cmd/agency-config/operations_help.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/main_test.go`

### Checkpoint 000003 -- Glossary And Recovery Guidance

Add operator-facing glossary and common mistake recovery guidance for:

- GTFS;
- GTFS-RT;
- telemetry and devices;
- Trip Updates and prediction;
- Alerts;
- validators and validation health;
- readiness, evidence, and claim boundaries;
- connectors;
- support bundles and maintenance.

Expected product behavior:

- glossary terms use plain language first, then technical detail;
- recovery rows explain what the operator sees, likely cause, safe next step,
  escalation trigger, and what the recovery action does not prove;
- guidance remains GET-only/read-only.

### Checkpoint 000004 -- Printable Operator Runbook And Handoff Checklist

Add a docs-based printable training guide and private UI links/checklists for
handoff between agency staff.

Expected product behavior:

- role-specific training paths for director, operator, technical helper, and
  integrator;
- quick tasks for importing GTFS, checking feed health, running validation,
  adding/simulating telemetry, reviewing connectors, preparing for a technical
  helper, and producing a support bundle;
- decision trees for "what should I do?" scenarios;
- no retained evidence and no stronger public claims.

Likely files:

- `docs/operator-training-guide.md`
- `cmd/agency-config/operations_help.go`
- `cmd/agency-config/main_test.go`

### Checkpoint 000005 -- Closeout

Write Phase 88 handoff and update current source-of-truth docs. Confirm:

- no protected path writes;
- exact prepared-only consumer tracker;
- no public admin routes;
- no release artifacts;
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

Phase 88 may say:

- private in-app guidance;
- training guide;
- role-based operator path;
- first-week checklist;
- glossary;
- troubleshooting decision tree;
- prepared packet records remain prepared only;
- optional evidence gates require separate written authorization.

Phase 88 must not claim:

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
