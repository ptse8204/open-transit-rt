# Codex Read Me First — Phase 61+ Agency-First Roadmap

This is the first file Codex should read before planning any new product work.

## Current Truth

- Phases 0 through 60 are closed for their documented scopes.
- The previous Post-60 productization work completed through Post-60 Checkpoint 000012.
- The maintainer has now authorized future roadmap work to use **Phase 61, Phase 62, Phase 63, ...** naming instead of more `Post-60` or `Phase 60 continuation` checkpoint names.
- This is a forward-looking naming change only. It does **not** reopen Phase 60 or any earlier phase.
- Optional evidence/adoption work remains authorization-gated and is not the default product path.

## Canonical Roadmap Location

```text
docs/roadmaps/agency-first-connector-platform/
```

Read in this order:

1. `00-CODEX-READ-ME-FIRST.md`
2. `README.md`
3. `01-roadmap-overview.md`
4. `02-phases-and-checkpoints.md`
5. `04-master-subagent-operating-manual.md`
6. `05-validation-and-claim-boundaries.md`
7. relevant `phase-prompts/phase-XX-*.md`
8. relevant `audit-prompts/*.md`

## Commit Naming Rule

Each phase has its own checkpoint sequence starting at `000001`.

Example:

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub
```

If a completed phase needs a later fix, continue that phase's checkpoint sequence:

```text
Phase 61 -- Checkpoint 000004: fix agency-first UI audit gaps
```

Do **not** use a global checkpoint counter across phases.

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

## First Commit To Add This Pack

After placing this directory in the repository, commit it as:

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
```

That commit should also add small links from current source-of-truth docs to this directory.
See `snippets-to-add-to-existing-docs.md`.

## First Implementation Phase

After the roadmap is committed, start:

```text
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
```

Use:

```text
phase-prompts/phase-61-agency-first-ui-and-connector-hub.md
```

## Hard Boundaries

Do not:

- write retained evidence;
- edit `docs/evidence/captured/**`;
- move consumer targets beyond `prepared`;
- contact agencies, vendors, consumers, portals, or external systems;
- claim CAL-ITP/Caltrans compliance;
- claim agency adoption or approval;
- claim consumer acceptance, ingestion, listing, or display;
- claim hosted SaaS, paid support, SLA, vendor compatibility, hardware certification, production readiness, or production-grade ETA quality.

## Master/Sub-Agent Mode

Use the operating model in `04-master-subagent-operating-manual.md`:

- master agent coordinates;
- planning sub-agent drafts phase plan;
- implementation sub-agent executes only after approval;
- QA sub-agent validates commands and tests;
- UI/UX sub-agent checks agency usability;
- claim-boundary sub-agent blocks unsupported wording.

If Codex cannot spawn real sub-agents, simulate these roles in sections.
