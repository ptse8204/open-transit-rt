# Phase 83 Prompt — Connector Workbench

## Goal

Make external connections easier to configure, test, and understand without claiming real vendor compatibility.

## User stories

- I have a CSV of vehicle locations.
- I have a GPS API.
- I have an AVL vendor that can POST data.
- I want synthetic telemetry only.
- I want an external prediction engine.
- I want monitoring summaries.
- I want to verify public feed URLs off-host.

## Scope

- Connector recipe chooser.
- Safe manifest review UI.
- CSV telemetry sandbox.
- API polling recipe.
- Webhook transform boundary guidance.
- Predictor sidecar recipe and shadow/fail-closed explanation.
- Monitoring/export recipe.
- UI-backed instructions for synthetic conformance checks.

## Do not

- Include real vendor payloads.
- Include credentials.
- Claim named vendor support.
- Claim vendor compatibility.
- Add dynamic backend plugin loading.
- Contact external systems.
- Write evidence.

## Validation

Baseline validation plus:

```bash
make validate
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
```

## Commits

```text
Phase 83 -- Checkpoint 000001: add connector workbench plan
Phase 83 -- Checkpoint 000002: add connector recipe chooser and manifest review
Phase 83 -- Checkpoint 000003: add CSV and API telemetry connector sandbox
Phase 83 -- Checkpoint 000004: add webhook and vendor transform boundary guidance
Phase 83 -- Checkpoint 000005: add predictor and monitoring connector recipe UI
Phase 83 -- Checkpoint 000006: add synthetic conformance runner guidance
Phase 83 -- Checkpoint 000007: close connector workbench review
```
