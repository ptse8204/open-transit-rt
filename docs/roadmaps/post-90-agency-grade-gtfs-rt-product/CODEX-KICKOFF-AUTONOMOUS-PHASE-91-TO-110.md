# CODEX KICKOFF PROMPT -- Autonomous Phase 91 To 110 Roadmap

You are the Master Agent for Open Transit RT. This kickoff authorizes the
post-90 agency-grade GTFS-RT product roadmap from Phase 91 through Phase 110.

## Current Truth

- Phases 0-60 are closed.
- Phase 61+ roadmap naming is approved.
- Phases 61-67 are complete.
- Phase 68+ is closed blocker-only / authorization-gated.
- Phases 69-74 are complete.
- Phases 75-90 are complete for the Consumer-Grade Control Plane track.
- Phase 72 remains `needs_review`, not release-ready.
- Phase 89 remains the current local `v0.1.0-rc.1` gate result and closes as
  `needs_review`.
- Phase 90 is complete for final control-plane closeout and future evidence
  gate stubs.

## Autonomous Sequence

Continue phase by phase:

```text
Phase 91 -> Phase 92 -> Phase 93 -> Phase 94 -> Phase 95 -> Phase 96 ->
Phase 97 -> Phase 98 -> Phase 99 -> Phase 100 -> Phase 101 -> Phase 102 ->
Phase 103 -> Phase 104 -> Phase 105 -> Phase 106 -> Phase 107 -> Phase 108 ->
Phase 109 -> Phase 110
```

Do not stop after roadmap-pack reconciliation, Phase 91, a documentation-only
phase, a successful validation, a `needs_review` diagnostic, or a blocker-only
optional-evidence phase.

## Phase List

- Phase 91 -- Maintainer Route/Product Audit And Stabilization.
- Phase 92 -- Clean Checkout Release-Candidate Gate.
- Phase 93 -- Browser End-To-End Agency Task Trials.
- Phase 94 -- Operations Console Architecture Refactor.
- Phase 95 -- v0.1.0-rc.1 Candidate Cut.
- Phase 96 -- GTFS Versioning, Diff, And Rollback Workbench.
- Phase 97 -- GTFS Quality Fix Planner And Safe Draft Suggestions.
- Phase 98 -- Realtime Operations QA And Feed Usefulness.
- Phase 99 -- Prediction / ETA Conformance And Backtesting V2.
- Phase 100 -- Alerts Operations And Disruption Workflow.
- Phase 101 -- Connector Maturity And Adapter Recipes V2.
- Phase 102 -- Device / AVL Fleet Onboarding V2.
- Phase 103 -- Monitoring, Notifications, And Export Surfaces.
- Phase 104 -- Small-Host Deployment And Upgrade Hardening.
- Phase 105 -- Multi-Agency Isolation And Operator Roles V2.
- Phase 106 -- Staff Training, Demo Datasets, And Adoption Kit.
- Phase 107 -- Public Docs/Site Freeze And Contributor Onboarding.
- Phase 108 -- Post-RC Bug Bash And Stabilization.
- Phase 109 -- Optional Evidence Intake Gate Pack.
- Phase 110 -- Long-Term Extensibility And Plugin Governance.

## Boundaries

This kickoff does not authorize protected evidence writes, external contact,
real credentials, real private payloads, consumer status movement beyond
`prepared`, tags, GitHub Releases, public package/image publication, or claims
of release readiness, compliance, adoption, consumer acceptance, production
readiness, hosted SaaS, SLA/uptime, vendor compatibility, hardware
certification, final-root readiness, or production-grade ETA quality.

Phase 95 authorizes only local `.cache` release-package generation and audit
when repo tooling supports it. It does not authorize publication.

## Protected Paths

Do not modify or generate files under:

```text
docs/evidence/captured/**
docs/evidence/consumer-submissions/status.json
docs/evidence/consumer-submissions/current/**
docs/evidence/consumer-submissions/artifacts/**
docs/evidence/consumer-submissions/packets/**
```

All seven consumer and aggregator targets must remain exactly `prepared`.

## Required Checkpoint Report

Each checkpoint and phase handoff must record:

```text
Checkpoint:
Sub-agents used or simulated, including intended model level:
Changed files:
Validation run:
Blocked checks:
Protected path status:
Consumer tracker status:
Claim-boundary status:
Security/auth status:
Data/migration status:
Master review:
Required edits:
Decision:
Next checkpoint:
```
