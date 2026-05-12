# Open Transit RT Agency-First Connector Platform Roadmap

This directory is the canonical repo location for the forward product roadmap that starts at **Phase 61**.

It supersedes the earlier local artifact language that used `Post-60` or `Phase 60 continuation` checkpoint names. The content and product direction remain the same: make Open Transit RT easier for small agencies, civic technologists, and integrators to use through better UI, guided setup, connector/plugin visibility, and release-quality workflows.

The current canonical review and next-step sequence is
[`docs/roadmap-status.md#review-and-recommendations`](../../roadmap-status.md#review-and-recommendations).
The default next product gate is `v0.1.0-rc.1` readiness before any full
`v0.1.0` release, with external-connection maturity as a product-quality path
and real pilots/final-root/consumer/vendor proof left as optional
authorization-gated evidence tracks.

Phase 71 keeps that direction adoption-first: browser-first agency operations,
maintenance, off-host validation, reference diagnostics, and docs/wiki/site
guidance. A real agency pilot is not the default next step.

## Current State

- Phases 0 through 60 are closed.
- Earlier Post-60 productization checkpoints through 000012 made the repo more agency-facing.
- The maintainer has now chosen Phase 61+ as the naming convention for this roadmap.
- Phase 61+ does not reopen earlier phases and does not weaken evidence/claim boundaries.

## Phase Map

| Phase | Name | Main outcome | Commit pattern |
| --- | --- | --- | --- |
| 61 | Agency-First UI and Connector Hub | Friendly Operations Console landing, Connector Hub, plugin/sidecar explanation | `Phase 61 -- Checkpoint 00000N: ...` |
| 62 | Guided Setup and Browser GTFS Import | Browser-led setup, metadata, GTFS URL/upload import, validation status | `Phase 62 -- Checkpoint 00000N: ...` |
| 63 | Feed Health and Readiness UX | Plain-language health dashboard, readiness checklist v2, feed freshness UX | `Phase 63 -- Checkpoint 00000N: ...` |
| 64 | Connector Platform and SDKs | Manifest registry UI, connector test runner, telemetry/prediction/monitoring SDK examples | `Phase 64 -- Checkpoint 00000N: ...` |
| 65 | Operator Workflow and Data Quality UX | Device/vehicle onboarding UI, telemetry simulator UI, GTFS quality fix guidance | `Phase 65 -- Checkpoint 00000N: ...` |
| 66 | Release Candidate and Installability | RC packaging, installer/bootstrap UX, Docker image decision, demo/docs site | `Phase 66 -- Checkpoint 00000N: ...` |
| 67 | Product Polish, Accessibility, and In-App Help | IA cleanup, responsive/accessibility polish, in-app help system | `Phase 67 -- Checkpoint 00000N: ...` |
| 68+ | Optional Authorized Evidence Tracks | Real agency/final-root/consumer/vendor/ETA evidence only when authorized | `Phase 68+ -- Checkpoint 00000N: ...` |
| 69 | Maintainer Product Acceptance and UI-First Agency Usability Trial | UI-first product acceptance, README/wiki/docs path cleanup, small-agency walkthroughs, and local acceptance audits without evidence intake | `Phase 69 -- Checkpoint 00000N: ...` |
| 70 | GitHub Pages Product Explainer Site | Static public explainer site on `gh-pages`, local/demo UI screenshots, and claim-bounded reader paths | `Phase 70 -- Checkpoint 00000N: ...` |
| 71 | Adoption-First Productization And No-CLI Agency Operations | Agency Operations Cockpit, richer GTFS/feed/validator/telemetry/maintenance UI, off-host validation, OCI reference check, and adoption docs | `Phase 71 -- Checkpoint 00000N: ...` |

## How To Use This Folder

Start every Codex session here:

```text
docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md
```

Then choose the current phase prompt:

```text
phase-prompts/phase-61-agency-first-ui-and-connector-hub.md
phase-prompts/phase-62-guided-setup-and-browser-gtfs-import.md
phase-prompts/phase-63-feed-health-and-readiness-ux.md
phase-prompts/phase-64-connector-platform-and-sdks.md
phase-prompts/phase-65-operator-workflow-and-data-quality-ux.md
phase-prompts/phase-66-release-candidate-and-installability.md
phase-prompts/phase-67-product-polish-accessibility-in-app-help.md
phase-prompts/phase-68-plus-optional-evidence-tracks.md
```

Phase 69 is a product acceptance bridge after Phase 68+ closed as
authorization-gated. It is not evidence intake and does not reopen final-root,
consumer, agency, vendor, or ETA proof tracks.

Phase 71 is documented in
[`adoption-productization-roadmap.md`](adoption-productization-roadmap.md).
It keeps real agency pilot work out of the default next-step sequence.

## Commit Pattern

Each phase owns its own checkpoint sequence.

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub

Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan
Phase 62 -- Checkpoint 000002: implement guided setup wizard v1
Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow
Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import
```

If a completed phase needs a patch later, continue that same phase sequence:

```text
Phase 61 -- Checkpoint 000004: fix connector hub audit gaps
```

## Core Product Goal

Open Transit RT should become a self-hosted, open-source GTFS / GTFS-Realtime operations platform that a small agency can use through a friendly UI, connect to external GPS/AVL/prediction/monitoring systems through safe connector plugins, and publish reliable feed outputs without proprietary vendor lock-in.

Vehicle Positions remain the first high-quality realtime output. Optional
Vehicle Positions fields need reliability gates, and external predictors should
be tested in shadow or fail-closed modes while deterministic prediction remains
the safe Trip Updates fallback.

## Non-Goals

This roadmap does not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- hosted SaaS availability;
- paid support or SLA;
- production readiness for all deployments;
- vendor compatibility or hardware certification;
- production-grade ETA quality.

These may become optional future evidence phases only when authorization and retained public-safe evidence exist.
