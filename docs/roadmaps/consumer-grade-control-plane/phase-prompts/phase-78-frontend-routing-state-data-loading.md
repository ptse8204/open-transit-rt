# Phase 78 Prompt — Frontend Routing, State, And Data Loading

## Goal

Add a progressive frontend interaction layer that makes the control plane feel modern without a speculative SPA rewrite.

## Scope

- Choose a low-dependency progressive pattern.
- Define no-build or lightweight-build policy.
- Add a frontend asset organization that is easy to audit.
- Add polling/progress pattern for long-running private tasks.
- Add safe client-side filtering/sorting for derived public-safe summaries.
- Add tests/lints appropriate to the selected lightweight stack.

## Recommended approach

Prefer server-rendered HTML plus small scoped JavaScript modules. Use JSON endpoints already present or added through Phase 77. Do not introduce React/Vue/Svelte/etc. until an explicit later architecture decision proves the need.

## Candidate interactions

- import progress status;
- validation result filtering;
- connector test progress;
- issue table filtering;
- inline help drawers;
- safe confirmation modal;
- copy-to-clipboard for public feed URLs.

## Forbidden

- No unauthenticated admin API.
- No client-owned security decisions.
- No localStorage secrets.
- No raw logs or tokens in browser state.
- No evidence writes.

## Validation

Baseline validation plus:

```bash
make validate
make test
RUN_LOCAL_APP=true make release-candidate-check
```

## Commits

```text
Phase 78 -- Checkpoint 000001: add frontend interaction architecture plan
Phase 78 -- Checkpoint 000002: add progressive UI runtime and asset policy
Phase 78 -- Checkpoint 000003: add private task progress pattern
Phase 78 -- Checkpoint 000004: apply interaction pattern to selected routes
Phase 78 -- Checkpoint 000005: close frontend interaction review
```
