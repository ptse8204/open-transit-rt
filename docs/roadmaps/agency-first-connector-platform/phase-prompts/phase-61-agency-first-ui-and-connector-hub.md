# Phase 61 — Agency-First UI and Connector Hub

## Goal

Make the private Operations Console feel like a real agency product UI and add a Connector Hub that advertises the external connector/plugin model.

## Commit Sequence

```text
Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap
Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub
Phase 61 -- Checkpoint 000003: close agency-first UI and connector hub
```

If audit gaps appear after close, use `Phase 61 -- Checkpoint 000004+`.

## Required Work

- Improve `/admin/operations` and `/admin/operations/launchpad` as agency-first pages.
- Add `/admin/operations/connectors`.
- Add `/admin/operations/connectors.json`.
- Add Connector Hub navigation.
- Show connector categories:
  - telemetry / GPS / AVL source;
  - prediction engine;
  - validator connector;
  - monitoring/export connector;
  - consumer/discovery workflow.
- Define plugin safely: optional sidecar, command adapter, manifest, or connector process; not arbitrary dynamic backend code loading.

## Tests

- auth required;
- agency scoping preserved;
- five connector categories rendered;
- JSON claim flags false;
- no unsupported claims;
- no evidence or consumer status changes.
