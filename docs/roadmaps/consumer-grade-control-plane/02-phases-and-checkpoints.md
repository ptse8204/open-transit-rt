# Consumer-Grade Control Plane — Phases And Checkpoints

## Phase 75 — Consumer-Grade Control Plane Roadmap Pack

Goal: add the roadmap pack and, where needed, bounded source-of-truth links. No
product implementation.

Checkpoints:

```text
Phase 75 -- Checkpoint 000001: add consumer-grade control plane roadmap
Phase 75 -- Checkpoint 000002: link consumer-grade roadmap from current status docs
Phase 75 -- Checkpoint 000003: close consumer-grade roadmap planning review
```

End-of-phase commit required: Checkpoint 000003.

## Cross-phase UI and claim gates

Every UI/control-plane implementation checkpoint after Phase 75 must include:

- keyboard, focus, and responsive/mobile review for touched pages;
- no-developer and technical-helper path review;
- command-safety review using the safe command ladder in
  `01-roadmap-overview.md`;
- status vocabulary review for queued/running/succeeded/failed/blocked/stale/
  unknown states where long-running tasks appear;
- copy audit for unsupported claims in labels, badges, help text, screenshots,
  readiness panels, and release notes;
- protected-path and seven-target prepared-only consumer tracker checks.

These gates are required before the Master Agent may close a future phase.

## Phase 76 — Design System And App Shell

Goal: make the private Operations Console feel like a coherent product.

Checkpoints:

```text
Phase 76 -- Checkpoint 000001: add design system and app shell plan
Phase 76 -- Checkpoint 000002: implement shared layout tokens and components
Phase 76 -- Checkpoint 000003: apply shell to core Operations Console routes
Phase 76 -- Checkpoint 000004: add responsive and accessibility baseline checks
Phase 76 -- Checkpoint 000005: close design system and app shell review
```

Core routes:

- `/admin/operations`
- `/admin/operations/setup-wizard`
- `/admin/operations/gtfs-import`
- `/admin/operations/feed-health`
- `/admin/operations/readiness`
- `/admin/operations/gtfs-quality`
- `/admin/operations/validation-health`
- `/admin/operations/devices`
- `/admin/operations/telemetry`
- `/admin/operations/telemetry-simulator`
- `/admin/operations/connectors`
- `/admin/operations/connectors/tests`
- `/admin/operations/maintenance`
- `/admin/operations/help`

End-of-phase commit required: Checkpoint 000005.

## Phase 77 — Admin Control API And Command Model

Goal: define how frontend pages safely control backend operations.

Checkpoints:

```text
Phase 77 -- Checkpoint 000001: add admin control API and command model plan
Phase 77 -- Checkpoint 000002: define private command/result contracts and shared response schema
Phase 77 -- Checkpoint 000003: add safe command audit tests and auth/CSRF fixtures
Phase 77 -- Checkpoint 000004: migrate one low-risk workflow to command model
Phase 77 -- Checkpoint 000005: close admin control API review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 78 — Frontend Routing, State, And Data Loading

Goal: add progressive frontend behavior without a risky SPA rewrite.

Checkpoints:

```text
Phase 78 -- Checkpoint 000001: add frontend interaction architecture plan
Phase 78 -- Checkpoint 000002: add small progressive UI runtime and no-build fallback policy
Phase 78 -- Checkpoint 000003: add polling/progress pattern for long-running private tasks
Phase 78 -- Checkpoint 000004: apply state/data-loading pattern to selected routes
Phase 78 -- Checkpoint 000005: close frontend interaction review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 79 — Agency Setup V3

Goal: turn setup into a guided browser workflow for nontechnical staff.

Checkpoints:

```text
Phase 79 -- Checkpoint 000001: add agency setup v3 plan
Phase 79 -- Checkpoint 000002: implement guided agency profile and feed metadata review
Phase 79 -- Checkpoint 000003: implement browser-first GTFS source/import review flow
Phase 79 -- Checkpoint 000004: implement setup progress, blockers, and readiness next actions
Phase 79 -- Checkpoint 000005: close agency setup v3 review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 80 — GTFS Workbench

Goal: help operators inspect, triage, edit, publish, and roll back static GTFS safely.

Checkpoints:

```text
Phase 80 -- Checkpoint 000001: add GTFS Workbench plan
Phase 80 -- Checkpoint 000002: add import diff and active/draft feed summary views
Phase 80 -- Checkpoint 000003: add stop route trip calendar table previews and filters
Phase 80 -- Checkpoint 000004: add safe draft edit and publish review improvements
Phase 80 -- Checkpoint 000005: add rollback and schedule history UX
Phase 80 -- Checkpoint 000006: close GTFS Workbench review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 81 — Realtime Operations Center

Goal: make live/recent realtime operations understandable from the browser.

Checkpoints:

```text
Phase 81 -- Checkpoint 000001: add realtime operations center plan
Phase 81 -- Checkpoint 000002: add fleet and telemetry freshness overview
Phase 81 -- Checkpoint 000003: add vehicle assignment and trip-matching explanation views
Phase 81 -- Checkpoint 000004: add Trip Updates and Alerts operational status views
Phase 81 -- Checkpoint 000005: add issue queue and simulator controls guidance
Phase 81 -- Checkpoint 000006: close realtime operations center review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 82 — Feed Health And Validation Center

Goal: unify feed status, validator results, quality triage, and readiness signals.

Checkpoints:

```text
Phase 82 -- Checkpoint 000001: add feed health and validation center plan
Phase 82 -- Checkpoint 000002: unify feed status and validator-history views
Phase 82 -- Checkpoint 000003: add validation issue drilldowns and fix-owner guidance
Phase 82 -- Checkpoint 000004: add readiness timeline and blocker history
Phase 82 -- Checkpoint 000005: close feed health and validation center review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 83 — Connector Workbench

Goal: make telemetry, AVL, prediction, validator, monitoring, and feed-consumer connectors easier to configure and test with synthetic/local data.

Checkpoints:

```text
Phase 83 -- Checkpoint 000001: add connector workbench plan
Phase 83 -- Checkpoint 000002: add connector recipe chooser and safe manifest review
Phase 83 -- Checkpoint 000003: add CSV and API telemetry connector sandbox
Phase 83 -- Checkpoint 000004: add webhook and vendor-transform boundary guidance
Phase 83 -- Checkpoint 000005: add predictor sidecar and monitoring/export recipe UI
Phase 83 -- Checkpoint 000006: add UI-backed synthetic conformance runner instructions
Phase 83 -- Checkpoint 000007: close connector workbench review
```

End-of-phase commit required: Checkpoint 000007.

## Phase 84 — Prediction And ETA Lab

Goal: improve ETA transparency and experimentation without production-grade ETA claims.

Checkpoints:

```text
Phase 84 -- Checkpoint 000001: add prediction and ETA lab plan
Phase 84 -- Checkpoint 000002: add deterministic predictor diagnostics view
Phase 84 -- Checkpoint 000003: add external predictor shadow/fail-closed review UI
Phase 84 -- Checkpoint 000004: add backtesting result browser and fixture guidance
Phase 84 -- Checkpoint 000005: add ETA quality disclaimers and withheld-output explanations
Phase 84 -- Checkpoint 000006: close prediction and ETA lab review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 85 — Operations And Maintenance Center V2

Goal: make maintenance tasks discoverable and safe for small-host/reference deployments.

Checkpoints:

```text
Phase 85 -- Checkpoint 000001: add maintenance center v2 plan
Phase 85 -- Checkpoint 000002: add backup and restore-drill browser guidance/status
Phase 85 -- Checkpoint 000003: add upgrade rollback and release readiness guidance
Phase 85 -- Checkpoint 000004: add secret rotation and support bundle guidance
Phase 85 -- Checkpoint 000005: add incident and reliability review UX
Phase 85 -- Checkpoint 000006: close maintenance center v2 review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 86 — Multi-Agency, Roles, Audit, And Accessibility

Goal: make future agency/role growth safer and make the UI accessibility baseline stronger.

Checkpoints:

```text
Phase 86 -- Checkpoint 000001: add multi-agency roles audit accessibility plan
Phase 86 -- Checkpoint 000002: add agency switcher and scope visibility improvements
Phase 86 -- Checkpoint 000003: add role/permission explanation and access-denied UX
Phase 86 -- Checkpoint 000004: add audit log browser for admin operations
Phase 86 -- Checkpoint 000005: run accessibility and keyboard navigation review
Phase 86 -- Checkpoint 000006: close multi-agency roles audit accessibility review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 87 — Public Feed Readiness And Docs Portal

Goal: make public feed URLs, metadata, and readiness understandable without claiming consumer acceptance or compliance.

Checkpoints:

```text
Phase 87 -- Checkpoint 000001: add public feed readiness portal plan
Phase 87 -- Checkpoint 000002: add public feed preview and metadata documentation pages
Phase 87 -- Checkpoint 000003: add consumer packet prepared-only explanation UX
Phase 87 -- Checkpoint 000004: add off-host validation and source-of-truth link guidance
Phase 87 -- Checkpoint 000005: close public feed readiness portal review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 88 — Nontechnical Training And In-App Guidance

Goal: help small-agency staff learn the product.

Checkpoints:

```text
Phase 88 -- Checkpoint 000001: add training and in-app guidance plan
Phase 88 -- Checkpoint 000002: add role-based tours and first-week checklist
Phase 88 -- Checkpoint 000003: add glossary, examples, and mistake recovery guidance
Phase 88 -- Checkpoint 000004: add printable operator runbook and handoff checklist
Phase 88 -- Checkpoint 000005: close training and in-app guidance review
```

End-of-phase commit required: Checkpoint 000005.

## Phase 89 — Release-Cut Cleanup / v0.1.0-rc.1 Gate

Goal: rerun release-candidate cleanup after UI/control-plane work.

Checkpoints:

```text
Phase 89 -- Checkpoint 000001: add post-control-plane rc1 gate plan
Phase 89 -- Checkpoint 000002: run clean-checkout local product gate
Phase 89 -- Checkpoint 000003: run frontend and accessibility gate
Phase 89 -- Checkpoint 000004: run connector and backend diagnostics gate
Phase 89 -- Checkpoint 000005: prepare rc1 notes package and blockers matrix
Phase 89 -- Checkpoint 000006: close rc1 gate review
```

End-of-phase commit required: Checkpoint 000006.

## Phase 90+ — Optional Authorized Evidence Gates

Goal: only if explicitly authorized, collect claim-specific evidence.

Possible tracks:

```text
Phase 90A -- Final-root evidence gate
Phase 90B -- Authorized consumer submission gate
Phase 90C -- Real agency pilot gate
Phase 90D -- Real vendor/device AVL gate
Phase 90E -- Real-world ETA quality gate
Phase 90F -- Deployment/compliance packet gate
```

Do not start any Phase 90+ track without explicit written authorization, exact claim target, allowed tools, public-safe retention rules, redaction rules, stop conditions, and operator/agency authorization.
