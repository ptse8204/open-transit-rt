# Consumer-Grade Control Plane — Roadmap Overview

## Product thesis

Open Transit RT should become the browser-first open-source control plane for small transit agencies that need to publish and operate GTFS/GTFS-Realtime without vendor lock-in.

The frontend should not be a decorative layer over scripts. It should become the primary way to control the backend safely.

## What “consumer-grade” means here

Consumer-grade does not mean a public rider app. It means the agency/operator UI should feel approachable, trustworthy, and coherent enough that a small agency staff member can use it without reading phase history or asking a developer for every step.

A consumer-grade agency control plane should have:

- polished visual hierarchy;
- responsive layout;
- accessible keyboard/focus behavior;
- clear page names and tasks;
- guided next actions;
- safe forms with previews and confirmations;
- empty states that explain what to do;
- a dashboard that shows health, readiness, and blockers;
- no raw logs, argv, private paths, secrets, or noisy internals;
- safe JSON/control APIs behind auth;
- strong evidence and claim boundaries.

## Strategic gaps this roadmap closes

| Gap | Current risk | Roadmap response |
| --- | --- | --- |
| UI feels primitive | Backend/product functionality is ahead of visual and workflow design. | Build a shared design system, app shell, page templates, and task-centered cockpit. |
| Backend is script-heavy | Useful capabilities are hard for nontechnical staff to discover or run. | Move common operator tasks into authenticated browser workflows. |
| GTFS editing is limited | Agencies need easier import review, diffing, triage, and safe publish/rollback. | Build a GTFS Workbench with import review, issue triage, map/table previews, draft publish, and rollback. |
| Realtime is hard to understand | Operators need to know whether vehicles, telemetry, matching, Trip Updates, and Alerts are useful now. | Build a Realtime Ops Center with live/recent health, fleet map, matching confidence, simulator controls, and issue queues. |
| Connectors are conceptually strong but adoption-heavy | Integrators need recipes, tests, and UI-backed conformance. | Build a Connector Workbench for CSV/API/webhook/prediction/monitoring recipes and synthetic conformance. |
| ETA quality is sensitive | The repo must improve prediction visibility without claiming production-grade ETA. | Build a Prediction/ETA Lab for deterministic vs external shadow comparison, backtesting, and diagnostics. |
| Release readiness is not proven | Phase 72 remains `needs_review`. | Add release-cut cleanup as a later, separate phase after UI/control-plane hardening. |
| Staff training is missing | Good software still fails if operators cannot learn it. | Add in-app training, role-based tours, glossary, and printable runbooks. |

## Non-goals

This roadmap does not authorize:

- public rider mobile app;
- evidence collection;
- consumer submission;
- external agency/vendor/consumer contact;
- hosted SaaS offering;
- production readiness claim;
- compliance claim;
- vendor compatibility claim;
- production-grade ETA claim;
- rewrite to a heavy SPA without prior API/control-plane proof.

## Recommended frontend architecture path

### Stage A — Design system and server-rendered control-plane shell

Use the existing authenticated Go admin routes. Add a consistent layout, typography, spacing, components, tables, cards, forms, alerts, empty states, and responsive behavior.

### Stage B — Private JSON models and command boundaries

Where the UI needs dynamic behavior, add typed private JSON routes and command endpoints with strict auth, CSRF, request caps, server-owned mappings, and all-false claim flags where relevant.

### Stage C — Progressive interactions

Add small scoped JavaScript for task flows:

- import progress polling;
- validation result filtering;
- map previews;
- connector test progress;
- timeline/history views;
- safe confirmation dialogs.

### Stage D — Optional future split frontend

Only after private API contracts are stable and tested, evaluate whether a separate frontend package is worthwhile. Do not start with this.

## Cross-phase control-plane rules

These rules apply to every future implementation phase in this roadmap.

### Information architecture

Use stable primary navigation groups so operators do not have to understand
phase history:

- Start Here: setup progress, primary actions, and current blockers.
- Schedule: GTFS import, review, edit, publish, history, and rollback.
- Realtime: telemetry, vehicles, matching, Vehicle Positions, Trip Updates,
  Alerts, and simulator workflows.
- Health: feed health, validation, readiness, incidents, and blocker history.
- Connectors: telemetry, AVL, prediction, validation, monitoring/export, and
  feed-consumer metadata workflows.
- Maintain: backups, restore drills, upgrades, secrets, support bundles, and
  release-cut checks.
- Learn: help, glossary, training paths, and technical-helper handoffs.

Each page needs one clear owner, one primary action, a shared empty-state
pattern, and consistent status vocabulary. Avoid duplicating the same concept
across health, readiness, validation, and maintenance without a clear link
between them.

### Safe command ladder

Classify browser-controlled workflows before adding controls:

| Level | Examples | Required UI guard |
| --- | --- | --- |
| Read-only refresh | reload status, fetch derived summaries | show timestamp and source |
| Dry run | synthetic connector check, validation preview | label as no-mutation and synthetic/local where applicable |
| Reversible private change | save metadata, update draft-only state | preview changed object, agency scope, audit result, and rollback/undo path |
| Publish or activate | activate GTFS, update public feed metadata | explicit confirmation, public-feed impact, rollback path, and claim boundary |
| Destructive or hard-to-reverse | restore drill, secret rotation, purge, live rollback | default to guidance or technical-helper handoff unless separately authorized |

Confirmation copy must state the exact object, agency scope, public/private
impact, whether public feeds change, the rollback path, and what the action
does not prove. Disabled states must explain the missing role, missing data,
or blocker.

### Status visibility

Long-running or backend-controlled workflows must expose bounded statuses:
`queued`, `running`, `succeeded`, `failed`, `blocked`, `stale`, and `unknown`.
Each status view should include last updated time, owner, retry availability,
next action, and whether a technical helper is required.

### Accessibility and copy audits

Every UI phase must include keyboard/focus/mobile checks, not only the later
accessibility phase. Every UI/content checkpoint must run a copy audit for
unsupported claims in labels, badges, empty states, screenshots, readiness
panels, release notes, and help text.

### No-developer acceptance

Every workflow should answer: can a no-developer operator complete this task,
or can they safely hand it to a technical helper with the right context? Shell
commands remain technical-helper paths and should be labeled that way.

## Proposed phase sequence

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 75 | Consumer-Grade Control Plane Roadmap Pack | Add this roadmap, prompts, and source-of-truth links. |
| 76 | Design System And App Shell | Replace primitive UI feel with consistent shell, components, responsive layout, and accessibility baseline. |
| 77 | Admin Control API And Command Model | Define safe private JSON/control API patterns for frontend-backed operations. |
| 78 | Frontend Routing, State, And Data Loading | Add progressive UI interaction model without SPA churn. |
| 79 | Agency Setup V3 | Full browser-first agency onboarding, metadata, feed URLs, readiness gates, and setup progress. |
| 80 | GTFS Workbench | Import review, GTFS quality triage, table/map previews, safe draft publish, rollback. |
| 81 | Realtime Operations Center | Fleet/telemetry/trip-matching/Trip Updates/Alerts health, issue queues, simulator controls. |
| 82 | Feed Health And Validation Center | Unified feed status, validation history, issue owners, fix workflow, no evidence claims. |
| 83 | Connector Workbench | CSV/API/webhook AVL recipes, predictor sidecars, monitoring/export recipes, UI conformance. |
| 84 | Prediction And ETA Lab | Deterministic/external shadow comparison, backtesting UI, quality diagnostics, no ETA-quality claim. |
| 85 | Operations And Maintenance Center V2 | Backup, restore drill, upgrades, secrets, support bundle, incidents, small-host diagnostics. |
| 86 | Multi-Agency, Roles, Audit, And Accessibility | Safer tenant switching, role UX, audit log views, WCAG-oriented review. |
| 87 | Public Feed Readiness And Docs Portal | Public feed preview/readiness pages and documentation portal without consumer/evidence claims. |
| 88 | Nontechnical Training And In-App Guidance | Tours, role-based learning paths, glossary, checklists, printable training. |
| 89 | Release-Cut Cleanup / v0.1.0-rc.1 Gate | Clean checkout, packaging, release notes, blocked/not-ready truth state. |
| 90+ | Optional Authorized Evidence Gates | Final-root, consumer, agency, vendor, ETA evidence only with explicit authorization. |

## Success definition

The roadmap succeeds when a maintainer can say:

> Open Transit RT has a polished, browser-first, self-hosted agency control plane that lets operators configure, import, publish, monitor, validate, connect, and maintain GTFS/GTFS-Realtime workflows safely, while preserving evidence and claim boundaries.

It still must not claim external adoption, compliance, final-root readiness, consumer acceptance, vendor compatibility, hosted SaaS, SLA, production readiness, or production-grade ETA quality without retained evidence.
