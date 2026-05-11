# Roadmap Overview — Phase 61+ Agency-First Connector Platform

## Vision

Open Transit RT should become a self-hosted, open-source GTFS / GTFS-Realtime operations platform that a small agency can use through a friendly UI, connect to external GPS/AVL/prediction/monitoring systems through safe connector plugins, and publish reliable feed outputs without proprietary vendor lock-in.

The product should feel less like:

```text
clone repo -> read many docs -> run commands -> inspect JSON -> understand phase history
```

and more like:

```text
open admin UI -> follow setup wizard -> import GTFS -> check feeds -> connect telemetry -> review readiness -> add connectors
```

## Why Phase 61+

Earlier source-of-truth docs correctly said Post-60 productization was not Phase 61. The maintainer has now changed the forward naming rule: future product work should use **Phase 61, Phase 62, ...** while preserving all existing evidence and claim boundaries.

This is a naming update for future product work. It does not reopen old phases.

## Roadmap Tracks

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

## Core Principles

1. Agency UI first: reduce command-line dependency.
2. Connector-ready: advertise safe plugin/sidecar boundaries.
3. No runtime overreach: avoid arbitrary dynamic plugin loading.
4. Evidence stays optional: final-root, adoption, consumer acceptance, and compliance are future evidence tracks only.
5. One phase, one outcome: finish and audit each phase before moving on.
6. Every completed phase gets a closeout checkpoint.
7. Patches to completed phases use the next checkpoint number inside that same phase.

## First Milestone

Phase 61 is the first roadmap phase.

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub
```
