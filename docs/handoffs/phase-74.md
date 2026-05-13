# Phase 74 Handoff -- GitHub Pages And Agency UI Product Polish

## Status

Checkpoint 000007 is complete for site and UI product-polish closeout.

Phase 74 is authorized by maintainer instruction as `GitHub Pages And Agency
UI Product Polish`. This supersedes the previous Phase 74 connector-maturity
slot, which is postponed to a later phase.

Phase 73 CP000006 is complete for bounded agency UI acceptance closeout.
Phase 72 remains `needs_review` and is not release-ready.

## Goal

Align the public GitHub Pages documentation site and the private Operations
Console around the same browser-first path:

1. Start in the browser.
2. Open `Agency Operations Cockpit / Start Here`.
3. Review setup.
4. Import or review GTFS.
5. Check the five public feed paths.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

## Worktrees

Main branch worktree:

```text
/Users/edwintse/Downloads/open-transit-rt
```

GitHub Pages worktree:

```text
/Users/edwintse/Downloads/open-transit-rt-gh-pages
```

Do not switch the main worktree to `gh-pages`. Use fixed working directories
for every command.

## Current Truth

- GitHub Pages is documentation only, not hosted SaaS.
- Public screenshots are local/demo documentation aids only, not retained
  evidence.
- Phase 72 is not release-ready.
- Consumer tracker targets remain exactly seven prepared-only records:
  Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database,
  and transit.land.
- No external party may be contacted.
- No retained evidence may be created.
- No consumer status may be moved.

## Checkpoint Sequence

1. `CP000001 -- add site and agency UI polish plan`
2. `CP000002 -- refresh GitHub Pages product story`
3. `CP000003 -- refresh GitHub Pages UI tour and quickstart`
4. `CP000004 -- improve Operations Console visual hierarchy`
5. `CP000005 -- improve first-run empty states and next actions`
6. `CP000006 -- align docs README wiki and site navigation`
7. `CP000007 -- close site and UI product polish review`

## CP000001 Plan

CP000001 adds:

- `docs/phase-74-github-pages-and-agency-ui-product-polish.md`
- `docs/handoffs/phase-74.md`

The plan records:

- why GitHub Pages is stale;
- why the private UI still needs product polish;
- which `gh-pages` files are expected to change;
- which private UI files are expected to change;
- what will not change;
- validation commands;
- screenshot policy;
- claim boundaries.

## CP000001 Result

Changed files:

- `docs/phase-74-github-pages-and-agency-ui-product-polish.md`
- `docs/handoffs/phase-74.md`

Sub-agents used or simulated:

- Context / Repo Truth Sub-Agent, intended GPT-5.5 x-high: real, complete.
- Planning Sub-Agent, intended GPT-5.5 x-high: real, complete.
- QA Sub-Agent, intended GPT-5.5 high: real, complete.
- UI/UX Sub-Agent, intended GPT-5.5 high: real, complete.
- Documentation / IA Sub-Agent, intended GPT-5.5 high: real, complete.
- Implementation Sub-Agent, intended GPT-5.5 high: simulated after agent
  error.
- Claim-Boundary Sub-Agent, intended GPT-5.5 high: simulated after agent
  error.

Validation run:

```bash
git status --short
git diff --check
make check
go test ./cmd/agency-config
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

Also run in the `gh-pages` worktree:

```bash
git status --short
git diff --check
```

Validation result: passed. The main worktree shows only the new Phase 74 plan
and handoff. The `gh-pages` worktree remains unchanged. Protected evidence,
consumer tracker, migration, and module paths have no status. The consumer
tracker remains exactly seven prepared-only targets.

Master review: CP000001 is approved. No required edits remain for the planning
checkpoint.

Next checkpoint: CP000002 -- refresh GitHub Pages product story.

## CP000002 Result

Changed `gh-pages` files:

- `assets/site.css`
- `connectors.html`
- `contribute.html`
- `how-it-works.html`
- `index.html`
- `quickstart.html`
- `readiness.html`
- `status.html`
- `ui-tour.html`

Summary:

- The public site navigation now leads with `Start Here`.
- The home page now makes `Agency Operations Cockpit / Start Here` the first
  browser concept.
- The home page separates no-developer browser review from technical-helper
  startup.
- The home page links to README, Small Agency Quick Start, Browser-First
  Setup, Operations Console Tour, No Command Line First Run, and Roadmap
  Status.
- Status/readiness copy now states GitHub Pages is documentation only and not
  a hosted app, SaaS offer, public feed root, public proof, or stronger claim.
- `assets/site.css` removed the decorative `border-left: 3px` list accent and
  added simple styling for the new documentation-only boundary copy.

Validation run:

```bash
# gh-pages worktree
git status --short
git diff --check
rg -n "Start Here|Agency Operations Cockpit / Start Here|documentation only|hosted SaaS|Small Agency Quick Start|Browser-First Setup|Operations Console Tour|No Command Line First Run|technical-helper" *.html

# main worktree
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

Validation result: passed. Protected evidence, migration, and module paths had
no status. The consumer tracker remains exactly seven prepared-only targets.
No `gh-pages` push was performed.

Master review: CP000002 is approved. No required edits remain for the public
site product-story checkpoint.

Next checkpoint: CP000003 -- refresh GitHub Pages UI tour and quickstart.

## CP000003 Result

Changed `gh-pages` files:

- `quickstart.html`
- `ui-tour.html`

Summary:

- Quickstart now starts from the private Operations Console URL and
  `Agency Operations Cockpit / Start Here`.
- Quickstart separates no-developer browser review from technical-helper
  startup commands.
- Quickstart lists the five local public feed paths and maps the expected
  private UI surfaces.
- UI tour now marks screenshots as local/demo documentation aids only.
- UI tour states stale or unavailable screenshots are not faked and uses
  annotated text cards instead.
- UI tour covers Agency Operations Cockpit / Start Here, Browser GTFS Import,
  Feed Health, Readiness, GTFS Quality, Validation Health, Device
  Credentials, Telemetry Freshness, Telemetry Simulator, Connector Hub,
  Connector Tests, Maintenance Center, and Operations Console Help.

Validation run:

```bash
# gh-pages worktree
git status --short
git diff --check
rg -n "Agency Operations Cockpit / Start Here|Browser GTFS Import|Feed Health|Readiness|GTFS Quality|Validation Health|Device Credentials|Telemetry Freshness|Telemetry Simulator|Connector Hub|Connector Tests|Maintenance Center|Operations Console Help|local/demo documentation aids only|does not fake|documentation only|hosted SaaS|/public/feeds.json|/public/gtfs/schedule.zip|/public/gtfsrt/vehicle_positions.pb|/public/gtfsrt/trip_updates.pb|/public/gtfsrt/alerts.pb" quickstart.html ui-tour.html

# main worktree
git diff --check
make check
go test ./cmd/agency-config
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

Validation result: passed. Protected evidence, migration, and module paths had
no status. The consumer tracker remains exactly seven prepared-only targets.
No screenshots were generated or treated as evidence. No `gh-pages` push was
performed.

Master review: CP000003 is approved. No required edits remain for the
quickstart and UI-tour checkpoint.

Next checkpoint: CP000004 -- improve Operations Console visual hierarchy.

## CP000004 Result

Changed main-branch files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/main_test.go`

Summary:

- Dashboard contextual help now renders after `Agency Operations Cockpit /
  Start Here` and the first-run panel, not before them.
- Contextual help remains available on dashboard and existing sections.
- First-run paths are visually distinct cards.
- `Developer path` was renamed to `Technical-helper path`.
- The five feed URLs now render as copy-card style items before the dense
  first-run acceptance task table.
- Setup progress and primary action statuses now use plain-language status
  chips.
- First-run claim flags remain inside the existing details disclosure.

Sub-agents used or simulated:

- Context / Repo Truth Sub-Agent, intended GPT-5.5 x-high: real, complete.
- Planning Sub-Agent, intended GPT-5.5 x-high: real prior checkpoint report
  used.
- QA Sub-Agent, intended GPT-5.5 high: simulated for CP000004 after
  thread-limit error.
- UI/UX Sub-Agent, intended GPT-5.5 high: simulated for CP000004 after
  thread-limit error.
- Documentation / IA Sub-Agent, intended GPT-5.5 high: real prior checkpoint
  report used.
- Implementation Sub-Agent, intended GPT-5.5 high: simulated.
- Claim-Boundary Sub-Agent, intended GPT-5.5 high: simulated after
  thread-limit error.

Validation run:

```bash
gofmt -w cmd/agency-config/operations.go cmd/agency-config/operations_first_run.go cmd/agency-config/main_test.go
git diff --check
go test ./cmd/agency-config
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum

# gh-pages worktree
git status --short
git diff --check
```

Validation result: passed. Protected evidence, migration, and module paths had
no status. The consumer tracker remains exactly seven prepared-only targets.
No public unauthenticated operations route, retained evidence, external
contact, consumer status change, DB migration, telemetry ingest change,
protobuf contract change, connector schema change, or forbidden claim was
introduced.

Master review: CP000004 is approved. No required edits remain for the visual
hierarchy checkpoint.

Next checkpoint: CP000005 -- improve first-run empty states and next actions.

## CP000005 Result

Changed main-branch files:

- `cmd/agency-config/operations.go`
- `cmd/agency-config/main_test.go`

Summary:

- Added compact empty/blocker guidance blocks to Browser GTFS Import, Feed
  Health, GTFS Quality, Validation Health, Device Credentials, Telemetry
  Freshness, Telemetry Simulator, Connector Hub, Connector Tests, Maintenance
  Center, and Operations Console Help.
- Each block answers what the user is seeing, whether it is bad, what to do
  next, whether the browser can handle it, when a technical helper is needed,
  and what the state does not prove.
- Added a focused route test for the shared first-run empty-state question
  set.
- Reworded one feed-health boundary phrase from `SLA coverage` to a neutral
  service-level boundary to satisfy existing safe-string tests.

Validation run:

```bash
gofmt -w cmd/agency-config/operations.go cmd/agency-config/main_test.go
git diff --check
go test ./cmd/agency-config
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum

# gh-pages worktree
git status --short
git diff --check
```

Validation result: passed. Protected evidence, migration, and module paths had
no status. The consumer tracker remains exactly seven prepared-only targets.
No retained evidence, external contact, consumer status change, public
unauthenticated operations route, DB migration, telemetry ingest change,
protobuf contract change, connector schema change, or forbidden claim was
introduced.

Master review: CP000005 is approved. No required edits remain for the
empty-state checkpoint.

Next checkpoint: CP000006 -- align docs README wiki and site navigation.

## CP000006 Result

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

Summary:

- README, docs home, wiki home, Small Agency Quick Start, Browser-First Setup,
  Operations Console Tour, No Command Line First Run, and GitHub Pages now use
  the same seven-step browser-first product path.
- GitHub Pages index, quickstart, and status pages now show the same order:
  start in the browser; open Agency Operations Cockpit / Start Here; review
  setup; import/review GTFS; check five feed URLs; review feed health,
  readiness, validation, telemetry, connectors, and maintenance; understand
  what remains before deployment or stronger claims.
- Primary browser-first docs now say `technical-helper path` instead of
  `developer path` for the helper role.

Validation run:

```bash
git diff --check
make check
go test ./cmd/agency-config
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
# scoped rg checks for the seven-step path in README, docs, wiki, and gh-pages
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum

# gh-pages worktree
git status --short
git diff --check
```

Validation result: passed. Protected evidence, migration, and module paths had
no status. The consumer tracker remains exactly seven prepared-only targets.
No retained evidence, external contact, consumer status change, public route
change, DB migration, telemetry ingest change, protobuf contract change,
connector schema change, or forbidden claim was introduced.

Master review: CP000006 is approved. No required edits remain for navigation
alignment.

Next checkpoint: CP000007 -- close site and UI product polish review.

## CP000007 Result

Changed main-branch closeout files:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-74.md`
- `docs/phase-74-github-pages-and-agency-ui-product-polish.md`
- `docs/roadmap-status.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`

Summary:

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

Sub-agents used or simulated:

- Context / Repo Truth Sub-Agent, intended GPT-5.5 x-high: real prior
  checkpoint reports used.
- Planning Sub-Agent, intended GPT-5.5 x-high: real prior checkpoint report
  used.
- QA Sub-Agent, intended GPT-5.5 high: simulated for CP000007 after
  thread-limit error.
- UI/UX Sub-Agent, intended GPT-5.5 high: simulated for CP000007 after
  thread-limit error.
- Documentation / IA Sub-Agent, intended GPT-5.5 high: simulated for
  closeout status alignment using prior real report.
- Implementation Sub-Agent, intended GPT-5.5 high: simulated.
- Claim-Boundary Sub-Agent, intended GPT-5.5 high: simulated.

Validation run:

```bash
git status --short
git diff --check
make check
go test ./cmd/agency-config
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
# exact seven-target prepared-only consumer tracker check
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum

# gh-pages worktree
git status --short
git diff --check
```

Validation result: passed. Protected evidence, consumer-submission,
migration, and module paths had no status. The consumer tracker remains
exactly seven prepared-only targets. No retained evidence, external contact,
consumer status change, public route change, DB migration, telemetry ingest
change, protobuf contract change, connector schema change, or forbidden claim
was introduced. Optional local app startup was not rerun for this docs/status
closeout checkpoint.

Master review: CP000007 is approved. No required edits remain for the Phase
74 GitHub Pages and agency UI product-polish review.

Next checkpoint: none for Phase 74. The next recommendation is maintainer
review of the Phase 74 closeout, then separately authorize future release-cut
cleanup/release-candidate gating, postponed connector maturity, or another
product phase.

## Protected Paths

Do not edit or generate files under:

- `docs/evidence/captured/**`;
- `docs/evidence/consumer-submissions/status.json`;
- `docs/evidence/consumer-submissions/current/**`;
- `docs/evidence/consumer-submissions/artifacts/**`;
- `docs/evidence/consumer-submissions/packets/**`;
- `db/migrations/**`;
- `go.mod`;
- `go.sum`.

## Claim Boundaries

Phase 74 may claim only:

- GitHub Pages product story refreshed;
- private Operations Console first-run hierarchy improved;
- docs/site/UI aligned around the same browser-first product path;
- no retained evidence created;
- no external party contacted;
- no consumer status changed;
- no compliance, adoption, consumer, final-root, SaaS, production, vendor,
  SLA, or ETA-quality claim added.

Phase 74 must not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption, approval, endorsement, or final-root readiness;
- consumer submission, review, acceptance, ingestion, listing, or display;
- hosted SaaS;
- production readiness;
- vendor compatibility;
- hardware certification;
- SLA or uptime coverage;
- production-grade ETA quality;
- release readiness.

## Validation

Required validation before phase closeout:

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

Optional local UI validation if environment supports it:

```bash
make agency-app-up
# manually or script-check the private UI pages
make agency-app-down
```

For the `gh-pages` worktree:

```bash
git status --short
git diff --check
```
