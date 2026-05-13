# Phase 74 -- GitHub Pages And Agency UI Product Polish

## Status

Checkpoint 000008 is complete for reconciling and publishing the actual
`gh-pages` browser-first refresh.
Checkpoint 000006 is complete for navigation alignment. Checkpoint 000005 is
complete for first-run empty states and next actions. Checkpoint 000004 is
complete for Operations Console visual hierarchy polish. Checkpoint 000003 is
complete for GitHub Pages quickstart and UI tour refresh. This phase
supersedes the previous Phase 74 connector-maturity slot; connector maturity
is postponed to a later separately authorized phase.

Phase 73 CP000006 is complete for bounded agency UI acceptance closeout.
Phase 72 remains `needs_review` and is not release-ready. This phase is not a
release-cut, evidence-collection, production-readiness, consumer-submission,
or compliance phase.

## Goal

Make the public GitHub Pages documentation site and the private Operations
Console tell the same product story:

1. Start in the browser.
2. Open `Agency Operations Cockpit / Start Here`.
3. Review setup.
4. Import or review GTFS.
5. Check the five public feed paths.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

The public site remains documentation only. The private Operations Console
remains authenticated server-rendered Go UI.

## Why GitHub Pages Is Stale

The GitHub Pages branch was created in Phase 70 as a static product explainer.
Since then, Phases 71 through 73 improved the browser-first agency path, the
private `Agency Operations Cockpit / Start Here`, Browser GTFS Import, Feed
Health, Device Credentials, Telemetry Freshness, Telemetry Simulator,
Connector Tests, Maintenance Center, and Operations Console Help.

The public site still orients readers around an older explainer/quickstart
shape. It does not make `Agency Operations Cockpit / Start Here` the first
browser concept, does not separate no-developer review from technical-helper
startup strongly enough, and does not consistently show the exact five-feed
path review as part of the first product story.

## Why The Private UI Still Needs Product Polish

Phase 73 closed the agency UI acceptance review with no remaining required
edits for that bounded scope. Phase 74 now raises the product-polish bar:
first-time evaluators should see the most important path before secondary
context, understand no-developer and technical-helper paths immediately, copy
or inspect the five feed URLs without hunting through dense diagnostics, and
recover from empty states without maintainer narration.

This is visual hierarchy and copy polish only. It must not change auth,
routes, database schema, telemetry ingest semantics, GTFS-RT protobuf
contracts, connector schemas, or evidence/consumer workflows.

## Checkpoint Sequence

1. `CP000001 -- add site and agency UI polish plan`
   - Add this plan.
   - Add `docs/handoffs/phase-74.md`.
   - Record the authorized scope, worktree split, validation commands,
     screenshot policy, protected paths, and claim boundaries.
2. `CP000002 -- refresh GitHub Pages product story`
   - Update the `gh-pages` branch so the public site starts with browser
     review and `Agency Operations Cockpit / Start Here`.
   - Keep the site documentation-only and claim-bounded.
3. `CP000003 -- refresh GitHub Pages UI tour and quickstart`
   - Update quickstart and tour pages around the current browser-first flow
     and required private UI surfaces.
   - Use text cards or diagrams when screenshots are stale or unavailable.
4. `CP000004 -- improve Operations Console visual hierarchy`
   - Improve first-screen hierarchy in the private Operations Console.
   - Keep server-rendered Go templates and existing private route boundaries.
5. `CP000005 -- improve first-run empty states and next actions`
   - Patch important empty/blocker states so evaluators understand what they
     are seeing, whether it is bad, what to do next, whether the browser can
     do it, when a technical helper is needed, and what it does not prove.
6. `CP000006 -- align docs README wiki and site navigation`
   - Ensure README, wiki, docs, and GitHub Pages share the same browser-first
     path.
7. `CP000007 -- close site and UI product polish review`
   - Update status, handoff, roadmap, planner, and closeout text.
   - Record validation, protected-path status, consumer tracker status, and
     claim-boundary review.
8. `CP000008 -- reconcile and publish GitHub Pages refresh`
   - Verify the actual `gh-pages` branch content.
   - Patch and publish any mismatch between the public branch and Phase 74
     closeout wording.
   - Update main-branch status docs with the published result.

## CP000001 Result

Checkpoint 000001 added this plan and `docs/handoffs/phase-74.md`.

Sub-agent reports used:

- Context / Repo Truth Sub-Agent, intended GPT-5.5 x-high: confirmed Phase 73
  closeout, Phase 72 `needs_review`, prepared-only consumer tracker state,
  public/private route boundaries, and the Phase 74 connector-maturity naming
  conflict superseded by maintainer instruction.
- Planning Sub-Agent, intended GPT-5.5 x-high: recommended fixed worktree use,
  checkpoint gates, protected-path checks, and master approval before each
  phase step.
- QA Sub-Agent, intended GPT-5.5 high: confirmed validation commands and the
  exact seven-target prepared-only tracker check.
- UI/UX Sub-Agent, intended GPT-5.5 high: recommended moving no-developer path
  before commands, improving dashboard hierarchy, compacting five feed URLs,
  and standardizing empty states.
- Documentation / IA Sub-Agent, intended GPT-5.5 high: confirmed README/wiki
  are mostly aligned and identified GitHub Pages as the largest current gap.
- Implementation and Claim-Boundary Sub-Agents were simulated after agent
  errors; the simulated review kept the scope to docs-only CP000001 changes
  and preserved all forbidden claim boundaries.

Validation run for CP000001:

- `git status --short`
- `git diff --check`
- `make check`
- `go test ./cmd/agency-config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `git status --short` and `git diff --check` in the `gh-pages` worktree

All CP000001 validation passed. Protected evidence, migration, and module
paths had no status. The `gh-pages` worktree remained unchanged. Consumer
tracker records remained exactly seven prepared-only targets.

## CP000002 Result

Checkpoint 000002 refreshed the GitHub Pages product story in the `gh-pages`
worktree. The site now puts `Start Here` first in the navigation, makes
`Agency Operations Cockpit / Start Here` the first browser concept on the
home page, separates no-developer review from technical-helper startup, links
to the canonical README, Small Agency Quick Start, Browser-First Setup,
Operations Console Tour, No Command Line First Run, and Roadmap Status docs,
and states that GitHub Pages is documentation only, not hosted SaaS or a
public feed root.

Changed `gh-pages` files:

- `index.html`
- `quickstart.html`
- `ui-tour.html`
- `how-it-works.html`
- `connectors.html`
- `readiness.html`
- `status.html`
- `contribute.html`
- `assets/site.css`

The detailed quickstart and UI-tour surface coverage remains for CP000003.

Validation run for CP000002:

- `git status --short` in the `gh-pages` worktree
- `git diff --check` in the `gh-pages` worktree
- scoped `rg` check for `Start Here`, `Agency Operations Cockpit / Start
  Here`, documentation-only wording, technical-helper wording, and required
  canonical documentation links
- `git diff --check` in the main worktree
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

All CP000002 validation passed. Protected evidence, migration, and module
paths had no status. Consumer tracker records remained exactly seven
prepared-only targets. No `gh-pages` push was performed.

## CP000003 Result

Checkpoint 000003 refreshed the GitHub Pages quickstart and UI tour around
the current browser-first Operations Console flow.

Changed `gh-pages` files:

- `quickstart.html`
- `ui-tour.html`

The quickstart now separates no-developer browser review from
technical-helper startup commands, points reviewers to `Agency Operations
Cockpit / Start Here`, lists the five local public feed paths, and maps the
private UI surfaces a first-time evaluator should inspect. The UI tour now
marks existing screenshots as local/demo documentation aids only, states that
stale or unavailable screenshots are not faked, uses annotated text cards for
missing surfaces, and covers:

- `Agency Operations Cockpit / Start Here`
- `Browser GTFS Import`
- `Feed Health`
- `Readiness`
- `GTFS Quality`
- `Validation Health`
- `Device Credentials`
- `Telemetry Freshness`
- `Telemetry Simulator`
- `Connector Hub`
- `Connector Tests`
- `Maintenance Center`
- `Operations Console Help`

Validation run for CP000003:

- `git status --short` in the `gh-pages` worktree
- `git diff --check` in the `gh-pages` worktree
- scoped `rg` check for all required UI-tour and quickstart surfaces, the
  five public feed paths, screenshot policy wording, and hosted SaaS boundary
- `git diff --check` in the main worktree
- `make check`
- `go test ./cmd/agency-config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

All CP000003 validation passed. Protected evidence, migration, and module
paths had no status. Consumer tracker records remained exactly seven
prepared-only targets. No screenshots were generated or treated as evidence.
No `gh-pages` push was performed.

## CP000004 Result

Checkpoint 000004 improved the private Operations Console visual hierarchy
without adding routes, changing auth, changing backend semantics, or changing
public feed contracts.

Changed main-branch files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/main_test.go`

The dashboard now renders `Agency Operations Cockpit / Start Here` and the
first-run panel before dashboard contextual help. Contextual help remains
available on the dashboard, but no longer displaces the primary first-run
path. The first-run panel now uses clearer status chips, visually distinct
no-developer and technical-helper path cards, and a feed URL copy-card grid
before the dense acceptance task table. The second path is now labeled
`Technical-helper path` instead of `Developer path` to match the public site
story and small-agency evaluation language. Existing claim flags remain
inside a details disclosure.

Sub-agent reports used:

- Context / Repo Truth Sub-Agent, intended GPT-5.5 x-high: real, complete;
  confirmed shared-layout contextual help rendered before dashboard content,
  identified first-run path/feed URL locations, existing CSS helpers, and
  relevant tests.
- UI/UX Sub-Agent, intended GPT-5.5 high: simulated after thread-limit error;
  recommendations followed the earlier real UI/UX report and focused on
  hierarchy, copy cards, status chips, and no-developer/technical-helper
  clarity.
- QA + Claim-Boundary Sub-Agent, intended GPT-5.5 high: simulated after
  thread-limit error; review stayed within existing private UI tests and
  forbidden-claim boundaries.

Validation run for CP000004:

- `gofmt -w cmd/agency-config/operations.go cmd/agency-config/operations_first_run.go cmd/agency-config/main_test.go`
- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `git status --short` and `git diff --check` in the `gh-pages` worktree

All CP000004 validation passed. Protected evidence, migration, and module
paths had no status. Consumer tracker records remained exactly seven
prepared-only targets. No new public route, retained evidence, external
contact, consumer status change, DB migration, telemetry ingest change,
protobuf contract change, connector schema change, or forbidden claim was
introduced.

## CP000005 Result

Checkpoint 000005 improved first-run empty and blocker states across the
private Operations Console.

Changed main-branch files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/main_test.go`

The affected private pages now include compact empty/blocker guidance blocks
that answer:

- what the evaluator is seeing;
- whether the state is bad;
- what to do next;
- whether the browser can do it;
- when a technical helper is needed;
- what the state does not prove.

Pages covered:

- Browser GTFS Import
- Feed Health
- GTFS Quality
- Validation Health
- Device Credentials
- Telemetry Freshness
- Telemetry Simulator
- Connector Hub
- Connector Tests
- Maintenance Center
- Operations Console Help

A focused test now checks the first-run empty-state question set across those
routes. Existing no-leakage/no-claim tests remain in force; one wording
adjustment changed an exact disallowed `SLA coverage` phrase to a neutral
service-level boundary.

Validation run for CP000005:

- `gofmt -w cmd/agency-config/operations.go cmd/agency-config/main_test.go`
- `git diff --check`
- `go test ./cmd/agency-config`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `git status --short` and `git diff --check` in the `gh-pages` worktree

All CP000005 validation passed. Protected evidence, migration, and module
paths had no status. Consumer tracker records remained exactly seven
prepared-only targets. No retained evidence, external contact, consumer
status change, public unauthenticated operations route, DB migration,
telemetry ingest change, protobuf contract change, connector schema change,
or forbidden claim was introduced.

## CP000006 Result

Checkpoint 000006 aligned README, wiki, docs, and GitHub Pages around the
same browser-first product path.

Changed main-branch files:

- `README.md`
- `docs/README.md`
- `docs/tutorials/no-cli-agency-first-run.md`
- `wiki/README.md`
- `wiki/browser-first-setup.md`
- `wiki/operations-console-tour.md`
- `wiki/small-agency-quick-start.md`

Changed `gh-pages` files:

- `index.html`
- `quickstart.html`
- `status.html`

The aligned path is:

1. Start in the browser.
2. Open `Agency Operations Cockpit / Start Here`.
3. Review setup.
4. Import or review GTFS.
5. Check the five public feed URLs.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

The docs also now use `technical-helper path` consistently instead of
`developer path` for the helper role in the primary browser-first evaluation
flow.

Validation run for CP000006:

- `git diff --check`
- `make check`
- `go test ./cmd/agency-config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- scoped `rg` checks for the seven-step browser-first path in README, docs,
  wiki, and GitHub Pages files
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `git status --short` and `git diff --check` in the `gh-pages` worktree

All CP000006 validation passed. Protected evidence, migration, and module
paths had no status. Consumer tracker records remained exactly seven
prepared-only targets. No retained evidence, external contact, consumer
status change, public route change, DB migration, telemetry ingest change,
protobuf contract change, connector schema change, or forbidden claim was
introduced.

## CP000007 Result

Checkpoint 000007 closed the GitHub Pages and agency UI product-polish review.

Changed main-branch closeout files:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-74.md`
- `docs/phase-74-github-pages-and-agency-ui-product-polish.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

Closeout result:

- GitHub Pages product story is refreshed.
- Private Operations Console first-run hierarchy is improved.
- Docs/site/UI now point to the same browser-first product path.
- No retained evidence was created.
- No external party was contacted.
- No consumer status changed.
- No compliance/adoption/consumer/final-root/SaaS/production/vendor/SLA/ETA
  claim was added.
- Phase 74 connector maturity remains postponed to a later separately
  authorized phase.

Validation run for CP000007:

- `git status --short`
- `git diff --check`
- `make check`
- `go test ./cmd/agency-config`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- `git status --short` and `git diff --check` in the `gh-pages` worktree

All CP000007 required validation passed. The optional local app startup check
was not rerun for this docs/status closeout checkpoint.

## CP000008 Result

Checkpoint 000008 reconciled and published the actual `gh-pages` branch.

Changed `gh-pages` file:

- `status.html`

Changed main-branch status files:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-74.md`
- `docs/phase-74-github-pages-and-agency-ui-product-polish.md`
- `docs/roadmap-status.md`

The checkpoint verified that the published site now makes `Agency Operations
Cockpit / Start Here` the first browser concept, starts the quickstart from
`/admin/operations`, treats `make check` and `make agency-app-up` as
technical-helper startup, lists all five public feed paths, covers the
required Operations Console tour surfaces, keeps screenshots as local/demo
documentation aids only, and preserves the documentation-only / not hosted
SaaS boundary.

The `gh-pages` branch was committed and pushed at `a8b250e` with message
`Phase 74 -- Checkpoint 000008: publish GitHub Pages browser-first refresh`.

All CP000008 required validation passed. Protected evidence, consumer tracker,
migration, and module paths had no status. The consumer tracker remained
exactly seven prepared-only targets. No retained evidence was created, no
external party was contacted, no consumer status changed, no release was
tagged or packaged, no connector maturity work started, and no forbidden claim
was added.

## Expected GitHub Pages Changes

Use the existing `gh-pages` worktree at
`/Users/edwintse/Downloads/open-transit-rt-gh-pages`. Do not switch the main
worktree to `gh-pages`.

Expected `gh-pages` files:

- `index.html`
- `quickstart.html`
- `ui-tour.html`
- `status.html`
- `readiness.html`
- `connectors.html`
- `assets/site.css`

Optional only if navigation consistency requires it:

- `how-it-works.html`
- `contribute.html`

The public site must:

- make `Agency Operations Cockpit / Start Here` the first browser concept;
- make "Start in the browser" prominent;
- separate no-developer review from technical-helper startup commands;
- link to README, Small Agency Quick Start, Browser-First Setup, Operations
  Console Tour, and No Command Line First Run;
- state that GitHub Pages is documentation only, not hosted SaaS;
- keep claim boundaries visible.

## Expected Private UI Changes

Expected main-branch private UI files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_cockpit.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_help.go`
- `cmd/agency-config/main_test.go`

Additional private UI files may be touched only for the requested empty-state
surfaces:

- `cmd/agency-config/operations_gtfs_import.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_gtfs_quality_guidance.go`
- `cmd/agency-config/operations_devices.go`
- `cmd/agency-config/operations_telemetry_simulator.go`
- `cmd/agency-config/operations_connectors.go`
- `cmd/agency-config/operations_connector_tests.go`
- `cmd/agency-config/operations_maintenance.go`

The UI polish should:

- put `Agency Operations Cockpit / Start Here` before contextual help on the
  dashboard;
- make no-developer and technical-helper paths visually obvious;
- make the five public feed paths easy to copy and understand;
- make "what this does not prove" visible but not overwhelming;
- improve empty states for no GTFS, no telemetry, missing validators,
  missing backup/restore configuration, and no connector setup;
- keep routes private/authenticated and server-rendered.

## What Will Not Change

Phase 74 must not:

- collect retained evidence;
- write under `docs/evidence/captured/**`;
- contact agencies, vendors, consumers, portals, or external systems;
- move consumer tracker statuses;
- add public unauthenticated admin, debug, operations, validation, scorecard,
  evidence, or authoring routes;
- change DB migrations;
- change telemetry ingest semantics;
- change GTFS-RT protobuf contracts;
- change connector schemas;
- change public feed path contracts;
- add a heavy frontend framework;
- tag, package, publish, push a branch, create a release, or claim release
  readiness;
- claim compliance, adoption, consumer acceptance, final-root readiness,
  hosted SaaS, production readiness, vendor compatibility, SLA/uptime,
  hardware certification, or production-grade ETA quality.

## Screenshot Policy

Screenshots and diagrams are local/demo documentation aids only. They are not
retained evidence and do not prove deployment, compliance, agency adoption,
consumer submission or acceptance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, SLA/uptime, hardware
certification, or production-grade ETA quality.

If screenshots on GitHub Pages are stale or unavailable, do not fake new
screenshots. Use annotated text cards or diagrams and label any existing
screenshots as local/demo documentation aids only.

## Claim Boundaries

Safe Phase 74 wording may say:

- GitHub Pages product story is refreshed.
- Private Operations Console first-run hierarchy is improved.
- Docs/site/UI point to the same browser-first product path.
- Local/demo screenshots are documentation aids only.
- Readiness and validation are supporting private diagnostics only.
- Consumer targets remain prepared-only.

Forbidden Phase 74 wording includes:

- CAL-ITP/Caltrans compliant;
- agency adopted, approved, endorsed, or live;
- submitted to, reviewed by, accepted by, listed by, displayed by, or ingested
  by any consumer;
- agency-owned final-root ready;
- hosted SaaS;
- production ready;
- vendor compatible;
- hardware certified;
- SLA-backed or uptime guaranteed;
- production-grade ETA quality;
- release-ready or `v0.1.0-rc.1` approved.

## Validation Commands

Run relevant checks after each checkpoint and the full set before closeout:

```bash
git status --short
git diff --check
make check
go test ./cmd/agency-config
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

Optional when environment supports it:

```bash
make agency-app-up
# manually or script-check the private UI pages
make agency-app-down
```

For the `gh-pages` worktree, also run:

```bash
git status --short
git diff --check
```

Do not push the `gh-pages` branch unless separately requested by the
maintainer. Editing the branch prepares the site update; pushing publishes it.

## Checkpoint Report Format

After every checkpoint, report:

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
