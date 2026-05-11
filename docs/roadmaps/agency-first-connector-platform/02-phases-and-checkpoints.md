# Phase 61+ Product Phases And Checkpoints

Each phase has its own checkpoint sequence starting at `000001`.

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

## Phase 61 — Agency-First UI and Connector Hub

**Goal:** Make the private Operations Console feel like a real agency product UI and add a Connector Hub.

**Checkpoint plan:**

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap` | Add this roadmap pack and source-of-truth links. |
| 000002 | `Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub` | Improve dashboard/launchpad and add `/admin/operations/connectors`. |
| 000003 | `Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub` | Update status/handoff/roadmap and run audit. |

Patches after close use `000004+`.

## Phase 62 — Guided Setup and Browser GTFS Import

**Goal:** Replace command-heavy first setup with a browser-led setup and GTFS import flow.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 62 -- Checkpoint 000001: add guided setup and browser GTFS import plan` | Scope forms, safety, upload/URL handling, and validation. |
| 000002 | `Phase 62 -- Checkpoint 000002: implement guided setup wizard v1` | Browser setup steps for agency profile, metadata, feeds, telemetry, validators, connectors. |
| 000003 | `Phase 62 -- Checkpoint 000003: implement browser GTFS import and validation flow` | Authenticated GTFS URL/upload import flow with status and validation display. |
| 000004 | `Phase 62 -- Checkpoint 000004: close guided setup and browser GTFS import` | Closeout docs/status/tests. |

## Phase 63 — Feed Health and Readiness UX

**Goal:** Give operators a simple health and readiness view for all feed outputs.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 63 -- Checkpoint 000001: add feed health and readiness UX plan` | Scope feed health model and readiness UX. |
| 000002 | `Phase 63 -- Checkpoint 000002: implement feed health dashboard` | Plain-language health for feeds.json, schedule, Vehicle Positions, Trip Updates, Alerts. |
| 000003 | `Phase 63 -- Checkpoint 000003: implement readiness checklist v2` | Agency-friendly CAL-ITP-style checklist with next actions. |
| 000004 | `Phase 63 -- Checkpoint 000004: close feed health and readiness UX` | Closeout docs/status/tests. |

## Phase 64 — Connector Platform and SDKs

**Goal:** Turn connector/plugin architecture into a visible, testable developer platform.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 64 -- Checkpoint 000001: add connector platform and SDK plan` | Scope connector registry, test runner, SDK examples. |
| 000002 | `Phase 64 -- Checkpoint 000002: implement connector manifest registry UI` | Show safe connector manifests in Connector Hub. |
| 000003 | `Phase 64 -- Checkpoint 000003: implement connector test runner UI` | UI for allowlisted connector/conformance checks or generated instructions. |
| 000004 | `Phase 64 -- Checkpoint 000004: improve telemetry connector SDK examples` | GPS/AVL adapter examples and helper docs. |
| 000005 | `Phase 64 -- Checkpoint 000005: improve prediction connector SDK examples` | Predictor sidecar/shadow examples and conformance cases. |
| 000006 | `Phase 64 -- Checkpoint 000006: improve monitoring export connector examples` | Redacted monitoring/export summaries. |
| 000007 | `Phase 64 -- Checkpoint 000007: close connector platform and SDKs` | Closeout docs/status/tests. |

## Phase 65 — Operator Workflow and Data Quality UX

**Goal:** Make day-to-day operator work easier: devices, simulator, and GTFS quality fixes.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 65 -- Checkpoint 000001: add operator workflow and data quality UX plan` | Scope device/vehicle onboarding, simulator UI, quality guidance. |
| 000002 | `Phase 65 -- Checkpoint 000002: implement device and vehicle onboarding UI` | Bind/rotate devices, show token once, latest telemetry. |
| 000003 | `Phase 65 -- Checkpoint 000003: implement telemetry simulator UI` | Choose synthetic scenarios and view results. |
| 000004 | `Phase 65 -- Checkpoint 000004: implement GTFS quality fix guidance UI` | Common validator/import issues translated into actionable guidance. |
| 000005 | `Phase 65 -- Checkpoint 000005: close operator workflow and data quality UX` | Closeout docs/status/tests. |

## Phase 66 — Release Candidate and Installability

**Goal:** Make the repo installable and releasable by real evaluators.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 66 -- Checkpoint 000001: add release candidate and installability plan` | Scope RC, install checks, Docker decision, docs site. |
| 000002 | `Phase 66 -- Checkpoint 000002: prepare first release candidate workflow` | Release notes, package audit, validation matrix. |
| 000003 | `Phase 66 -- Checkpoint 000003: improve installer and bootstrap UX` | Better preflight and missing-tool messages. |
| 000004 | `Phase 66 -- Checkpoint 000004: document Docker image publishing decision` | Source-only vs image publishing decision. |
| 000005 | `Phase 66 -- Checkpoint 000005: add demo site or documentation website plan` | Public-friendly website/docs structure. |
| 000006 | `Phase 66 -- Checkpoint 000006: close release candidate and installability` | Closeout docs/status/tests. |

## Phase 67 — Product Polish, Accessibility, and In-App Help

**Goal:** Make the product easier to navigate, read, and understand.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 67 -- Checkpoint 000001: add product polish and accessibility plan` | Scope IA, accessibility, in-app help. |
| 000002 | `Phase 67 -- Checkpoint 000002: improve operations console information architecture` | Clear nav groups and terminology. |
| 000003 | `Phase 67 -- Checkpoint 000003: improve accessibility and mobile layout` | Semantics, contrast, keyboard nav, responsive layout. |
| 000004 | `Phase 67 -- Checkpoint 000004: implement in-app help system` | Contextual explanations for GTFS, GTFS-RT, connectors, readiness, claims. |
| 000005 | `Phase 67 -- Checkpoint 000005: close product polish accessibility and help` | Closeout docs/status/tests. |

## Phase 68+ — Optional Authorized Evidence Tracks

Use only when explicitly authorized and evidence exists. These are not product blockers.

| Checkpoint | Commit message | Outcome |
| --- | --- | --- |
| 000001 | `Phase 68+ -- Checkpoint 000001: add optional evidence track blocker documentation` | Document Phase 68+ as authorization-gated and blocker-only when explicit written authorization is absent. |
| 000002 | `Phase 68+ -- Checkpoint 000002: close optional evidence tracks as authorization-gated` | Close Phase 68+ without evidence collection or status movement unless authorized intake and artifacts are present. |

Possible future phases:

- Phase 68 — Authorized Agency Trial Intake
- Phase 69 — Authorized Final Public Root Evidence
- Phase 70 — Authorized Consumer Submission Evidence
- Phase 71 — Authorized Real AVL/Vendor Integration Evidence
- Phase 72 — Authorized Real-World ETA Quality Study
