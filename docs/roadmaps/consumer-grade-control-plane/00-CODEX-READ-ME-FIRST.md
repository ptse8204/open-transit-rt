# Codex Read Me First — Consumer-Grade Control Plane Roadmap

This is the first file Codex should read before planning or implementing the consumer-grade frontend/control-plane roadmap.

## Current repo truth to preserve

- Phases 0-60 are closed.
- Phase 61+ roadmap naming is approved.
- Phases 61-67 are complete.
- Phase 68+ is closed blocker-only / authorization-gated.
- Phases 69-74 are complete for product acceptance, browser-first UX, GitHub Pages/product-site refresh, agency UI acceptance, and UI/site polish.
- Phase 74 CP000008 reconciled and published the actual `gh-pages` branch at commit `a8b250e`.
- GitHub Pages reflects the browser-first path and exact `Agency Operations Cockpit / Start Here` language.
- The private Operations Console first-run hierarchy is improved.
- Docs/site/UI point to the same browser-first product path.
- Phase 72 remains `needs_review`, not release-ready.
- No release tag, package, published image, or release-cut proof exists.
- Connector maturity is postponed.
- Evidence/adoption/compliance tracks remain optional and authorization-gated.

## Why this roadmap exists

The current UI is useful but still feels primitive because the product has historically been backend-first and docs/audit-first. The next product leap should make the browser UI the primary control plane for the backend:

- not just prettier pages;
- not a speculative SPA rewrite;
- not marketing polish;
- not evidence collection;
- not a production or compliance claim.

The control plane should let a nontechnical or semi-technical operator:

1. configure an agency;
2. import, inspect, and publish GTFS;
3. understand GTFS quality issues;
4. connect telemetry and devices through safe adapter paths;
5. monitor live/recent realtime status;
6. test Trip Updates/prediction behavior safely;
7. run validation and release-readiness checks;
8. manage maintenance, backups, support bundles, and upgrades;
9. understand what is ready, what is blocked, and what is not proven.

## Required reading before any new phase

Read these first:

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-74.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`
- `docs/roadmap-status.md#review-and-recommendations`
- `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`
- `docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md`
- `docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md`
- `docs/evidence/redaction-policy.md`
- `docs/evidence/consumer-submissions/status.json`

Also read and obey the binding product-contract docs named by `AGENTS.md`:

- `docs/codex-task.md`
- `docs/architecture.md`
- `docs/conversation-summary.md`
- `docs/requirements-2a-2f.md`
- `docs/requirements-trip-updates.md`
- `docs/requirements-calitp-compliance.md`
- `docs/repo-gaps.md`
- `docs/dependencies.md`

Then read this pack:

- `docs/roadmaps/consumer-grade-control-plane/README.md`
- `docs/roadmaps/consumer-grade-control-plane/01-roadmap-overview.md`
- `docs/roadmaps/consumer-grade-control-plane/02-phases-and-checkpoints.md`
- `docs/roadmaps/consumer-grade-control-plane/03-master-subagent-operating-manual.md`
- `docs/roadmaps/consumer-grade-control-plane/04-validation-and-claim-boundaries.md`

## Hard boundaries

Do not modify or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

Do not move consumer statuses. These seven targets must remain exactly `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

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

## Product architecture recommendation

Do not start with a heavy SPA rewrite. Start with a progressive control-plane architecture:

1. preserve existing Go server-rendered private admin routes;
2. introduce a shared design system and component/page shell;
3. add typed private JSON models where needed;
4. add small progressive JavaScript only where it materially improves usability;
5. add backend command APIs only after explicit scope and auth review;
6. consider a split SPA only after the control API contract is stable and tested.

This prevents frontend churn from breaking the backend, auth, CSRF, role boundaries, or claim discipline.

## Phase start rule

This pack is planning output only. A future Codex instance may add the pack in a docs-only phase. It may not implement Phase 76+ until the maintainer separately authorizes the phase.

First proposed phase:

```text
Phase 75 -- Consumer-Grade Control Plane Roadmap Pack
```

First proposed commit:

```text
Phase 75 -- Checkpoint 000001: add consumer-grade control plane roadmap
```

Phase 75 must be docs-only unless the maintainer explicitly authorizes otherwise.
