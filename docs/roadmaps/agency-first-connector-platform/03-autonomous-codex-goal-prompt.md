# Autonomous Codex `/goal` Prompt — Phase 61+ Roadmap

Use this after the roadmap is committed.

```text
/goal Complete the Open Transit RT Phase 61+ agency-first connector platform roadmap without stopping until the product has a friendly agency Operations Console, guided setup and GTFS import, visible Connector Hub, connector SDK examples, operator workflow UI, release-candidate installability, product polish, accessibility, and in-app help, while preserving all evidence and claim boundaries.
```

Then paste:

```text
You are the MASTER AGENT for Open Transit RT Phase 61+.

Current repo truth:
- Phases 0 through 60 are closed.
- Earlier Post-60 checkpoints through 000012 completed agency-ready productization.
- The maintainer has now authorized Phase 61+ naming for future product work.
- This is a forward naming change only; it does not reopen Phase 60 or earlier phases.
- Evidence/adoption work remains optional and authorization-gated.

Start from:
- docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md
- docs/roadmaps/agency-first-connector-platform/README.md
- docs/roadmaps/agency-first-connector-platform/02-phases-and-checkpoints.md
- docs/roadmaps/agency-first-connector-platform/04-master-subagent-operating-manual.md
- docs/roadmaps/agency-first-connector-platform/05-validation-and-claim-boundaries.md

Commit naming:
- Each phase has its own checkpoint sequence starting at 000001.
- Use `Phase XX -- Checkpoint 00000N: <short outcome>`.
- If a closed phase needs a fix, continue that phase's checkpoint sequence.

Run phases in order:
1. Phase 61 — Agency-First UI and Connector Hub
2. Phase 62 — Guided Setup and Browser GTFS Import
3. Phase 63 — Feed Health and Readiness UX
4. Phase 64 — Connector Platform and SDKs
5. Phase 65 — Operator Workflow and Data Quality UX
6. Phase 66 — Release Candidate and Installability
7. Phase 67 — Product Polish, Accessibility, and In-App Help

Do not start Phase 68+ optional evidence work unless the maintainer explicitly authorizes it.

Master/sub-agent mode:
- Planning sub-agent drafts a narrow phase plan.
- Master reviews and tightens plan before implementation.
- Implementation sub-agent executes only approved scope.
- QA sub-agent runs validations.
- UI/UX sub-agent checks small-agency usability.
- Claim-boundary sub-agent audits wording and protected files.

Protected paths:
- docs/evidence/captured/**
- docs/evidence/consumer-submissions/status.json
- docs/evidence/consumer-submissions/current/**
- docs/evidence/consumer-submissions/artifacts/**
- docs/evidence/consumer-submissions/packets/**
- db/migrations/** unless a phase explicitly requires schema and maintainer approves
- go.mod and go.sum unless dependency addition is explicitly approved

Forbidden claims:
- CAL-ITP/Caltrans compliance
- agency adoption or approval
- consumer submission/review/acceptance/ingestion/listing/display
- hosted SaaS
- paid support/SLA
- universal production readiness
- vendor compatibility/certification
- production-grade ETA quality

Default validation after each checkpoint:
- git diff --check
- make check
- go test ./cmd/agency-config when UI touched
- make test when feasible
- make external-connection-check when connector docs/examples touched
- make adapter-conformance when connector/conformance touched
- make test-connector-examples when examples touched
- make audit-final-claim-review
- python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
- exact seven-target prepared-only consumer tracker check
- git diff --exit-code -- docs/evidence/consumer-submissions/status.json
- git diff --exit-code -- docs/evidence/captured

Stop only when Phase 67 closes and the final roadmap completion audit passes, or when the maintainer explicitly stops the run, or when a hard blocker prevents safe progress.
```
