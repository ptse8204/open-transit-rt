# Copy-Paste Codex Kickoff Prompt

```text
You are the Master Agent for Open Transit RT.

Task: Start Phase 75 -- Consumer-Grade Control Plane Roadmap Pack.

This is a planning/artifact phase only. Do not implement UI, APIs, migrations, evidence collection, consumer submission, release tags, packages, or published images.

Read first:
- AGENTS.md
- README.md
- docs/current-status.md
- docs/handoffs/latest.md
- docs/handoffs/phase-74.md
- docs/open-transit-rt-master-planner-remaining-work.md
- docs/roadmap-status.md#review-and-recommendations
- docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md
- docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md
- docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md
- docs/evidence/redaction-policy.md
- docs/evidence/consumer-submissions/status.json
- AGENTS.md binding docs: docs/codex-task.md, docs/architecture.md, docs/conversation-summary.md, docs/requirements-2a-2f.md, docs/requirements-trip-updates.md, docs/requirements-calitp-compliance.md, docs/repo-gaps.md, docs/dependencies.md

Current truth:
Phases 0-74 are closed. Phase 72 remains needs_review, not release-ready. Phase 74 CP000008 reconciled and published gh-pages at commit a8b250e. No active phase starts automatically. Evidence/adoption/compliance tracks remain optional and authorization-gated. All seven consumer targets remain prepared.

Goal:
Add a new roadmap pack under:

docs/roadmaps/consumer-grade-control-plane/

The roadmap should plan a future consumer-grade browser control plane that can safely control backend workflows: agency setup, GTFS import/review/edit/publish, feed health, validation, realtime operations, telemetry/devices, connector workbench, prediction/ETA lab, maintenance, training, release-cut cleanup, and optional evidence gates.

Required master/sub-agent workflow:
- Master Agent: GPT-5.5 x-high
- Context / Repo Truth Sub-Agent: GPT-5.5 x-high
- Planning Sub-Agent: GPT-5.5 x-high
- Implementation Sub-Agent: GPT-5.5 high
- QA Sub-Agent: GPT-5.5 high
- UI/UX Sub-Agent: GPT-5.5 high
- Documentation / IA Sub-Agent: GPT-5.5 high
- Claim-Boundary Sub-Agent: GPT-5.5 high

If real sub-agents are unavailable, simulate these roles in labeled sections. The Master Agent must approve the plan before implementation starts and may move forward only when all sub-agent reviews have no required edits.

Protected paths:
Do not modify or generate files under:
- docs/evidence/captured/**
- docs/evidence/consumer-submissions/status.json
- docs/evidence/consumer-submissions/current/**
- docs/evidence/consumer-submissions/artifacts/**
- docs/evidence/consumer-submissions/packets/**

Consumer tracker:
All seven targets must remain exactly prepared:
Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land.

Forbidden claims:
Do not claim CAL-ITP/Caltrans compliance, agency adoption/approval, consumer submission/review/acceptance/ingestion/listing/display, final-root readiness, hosted SaaS, production readiness, vendor compatibility, hardware certification, SLA/uptime, or production-grade ETA quality.

Phase 75 deliverables:
- docs/roadmaps/consumer-grade-control-plane/README.md
- docs/roadmaps/consumer-grade-control-plane/00-CODEX-READ-ME-FIRST.md
- docs/roadmaps/consumer-grade-control-plane/01-roadmap-overview.md
- docs/roadmaps/consumer-grade-control-plane/02-phases-and-checkpoints.md
- docs/roadmaps/consumer-grade-control-plane/03-master-subagent-operating-manual.md
- docs/roadmaps/consumer-grade-control-plane/04-validation-and-claim-boundaries.md
- docs/roadmaps/consumer-grade-control-plane/phase-prompts/phase-75-consumer-grade-control-plane-roadmap-pack.md
- phase prompts through Phase 90+ optional evidence gates
- audit prompts for UI and claim-boundary review

Link the new pack from current source-of-truth docs only if needed. The links must say this is a proposed/authorized roadmap pack, not implemented future work.

Validation:
- git status --short
- git diff --check
- make check
- make audit-product-acceptance
- make audit-final-claim-review
- python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
- exact seven-target prepared-only consumer tracker check
- git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum

Commit sequence:
Phase 75 -- Checkpoint 000001: add consumer-grade control plane roadmap
Phase 75 -- Checkpoint 000002: link consumer-grade roadmap from status docs
Phase 75 -- Checkpoint 000003: close consumer-grade roadmap planning review

A phase is not closed until the final closeout commit exists and the Master Agent reports:
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
Commit created:
Next checkpoint:
```
