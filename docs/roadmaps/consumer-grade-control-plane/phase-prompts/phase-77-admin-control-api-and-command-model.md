# Phase 77 Prompt — Admin Control API And Command Model

## Goal

Define safe private JSON/control APIs so the frontend can control backend workflows without shelling out, leaking internals, or weakening auth.

## Why

A polished frontend needs controlled backend actions: import GTFS, rerun validators, test connectors, run diagnostics, start simulator flows, prepare support bundles, and review results. These must be safer than exposing arbitrary commands or paths.

## Scope

- Define a shared private command-result schema.
- Define `action`, `status`, `started_at`, `completed_at`, `summary`, `next_actions`, `claim_flags`, and bounded `errors` fields.
- Define the safe command ladder: read-only refresh, dry run, reversible private
  change, publish/activate, and destructive or hard-to-reverse action.
- For each ladder level, define preview, confirmation, role, audit, rollback,
  disabled-state, and technical-helper handoff requirements.
- Add tests for auth, roles, CSRF, request caps, unsupported fields, and no raw command/path exposure.
- Migrate one low-risk workflow to the model, such as a read-only refresh/review action or existing validation-health trigger if safe.
- Document command rules for future phases.

## Safety rules

- Browser forms must not supply shell commands, argv, validator paths, output paths, artifact paths, or timeouts.
- Server-owned mappings only.
- Cookie-auth POSTs require CSRF.
- Bearer-token exceptions only where already intentionally supported.
- No raw stdout/stderr in UI.
- No raw private file paths.
- No evidence writes.
- Confirmation copy must state the exact object, agency scope, public/private
  impact, whether public feeds change, rollback path, and what the action does
  not prove.

## Deliverables

- Internal command model docs.
- Shared helper/model if appropriate.
- Tests.
- One migrated low-risk action.
- UI wording that explains command outcomes.

## Validation

Baseline validation plus:

```bash
make validate
make test
```

## Commits

```text
Phase 77 -- Checkpoint 000001: add admin control API and command model plan
Phase 77 -- Checkpoint 000002: define private command result contracts
Phase 77 -- Checkpoint 000003: add command safety tests
Phase 77 -- Checkpoint 000004: migrate one low-risk workflow to command model
Phase 77 -- Checkpoint 000005: close admin control API review
```
