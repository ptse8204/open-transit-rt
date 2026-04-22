# Production Closure Plan After Phase 8

## Summary

This document formalizes the next repo-directed work after Phase 8. It is intentionally split into three short, evidence-driven phases instead of one large catch-all effort.

The goal is not just to add more code. The goal is to close the repo toward production-directed quality, make it understandable to agencies, and build truthful compliance/readiness evidence without overstating the system's maturity.

These phases should be treated as follow-on closure phases:

- **Phase 9 — Production Closure**
- **Phase 10 — Docs, Tutorials, Deployment, and Demo**
- **Phase 11 — Compliance Evidence and Optional External Integrations**

## Guiding Rules

- Do not claim production readiness without implementation, tests, docs, and smoke evidence.
- Do not claim CAL-ITP / Caltrans compliance without explicit evidence for the checklist items recorded in the requirements docs and deployment artifacts.
- Do not claim external integrations unless they are actually wired and tested.
- Do not allow docs to drift from code, commands, routes, or environment variables.
- Prefer closure with evidence over closure by assertion.

## Phase 9 — Production Closure

### Goal

Close the remaining production-readiness gaps in the codebase and make the current public/private surface operationally defensible.

### Scope

1. Finish any remaining hardening gaps from the hardening slice.
2. Verify anonymous protobuf publication remains stable.
3. Verify admin/debug/mutation surfaces are protected and correctly scoped.
4. Verify GTFS import, GTFS Studio, matcher, Trip Updates, Alerts, and publication behavior still work together after hardening.
5. Improve observability and deployment correctness.

### Required Work

- Complete realtime validator artifact derivation for `/admin/validation/run`.
- Pin validator tooling in local/dev/CI/prod setup docs and scripts.
- Finish admin auth/roles and agency scoping consistency.
- Add or complete structured request logging and request IDs.
- Add or complete metrics and deeper readiness checks.
- Add or complete strong smoke/e2e checks covering:
  - anonymous protobuf feed fetches
  - protected admin/debug routes
  - authenticated telemetry ingest
  - GTFS import and publish
  - validator execution flow
  - publication bootstrap and metadata flows

### Required Checks

- `gofmt -w ./cmd ./internal`
- `go mod tidy`
- `go test ./...`
- `make validate`
- `make test-integration`
- `make smoke`
- `docker compose -f deploy/docker-compose.yml config`
- `git diff --check`

### Phase 9 Exit Criteria

Phase 9 is complete only when:

- all public `.pb` endpoints remain anonymous and stable
- all admin/debug/mutation endpoints are protected
- realtime validator runs derive their own server-owned artifacts
- smoke/e2e tests exercise the full intended pilot path
- docs accurately describe the hardened runtime behavior
- no stale claims remain in README or status docs

## Phase 10 — Docs, Tutorials, Deployment, and Demo

### Goal

Make the repo understandable and demoable by agencies, contributors, and evaluators.

### Scope

1. Rewrite the README to current system reality.
2. Add runnable tutorials.
3. Add a simple agency demo flow.
4. Improve repo presentation and onboarding.
5. Generate diagrams/assets for docs.

### Required Work

Add or update:

- `README.md`
- `docs/tutorials/local-quickstart.md`
- `docs/tutorials/deploy-with-docker-compose.md`
- `docs/tutorials/agency-demo-flow.md`
- `docs/tutorials/production-checklist.md`
- `docs/tutorials/calitp-readiness-checklist.md`
- `docs/assets/README.md`

Package a simple agency demo flow that covers:

1. bootstrap DB
2. import sample GTFS
3. publish schedule
4. ingest authenticated telemetry
5. fetch public Vehicle Positions / Trip Updates / Alerts protobuf feeds
6. show protected admin/debug routes
7. run validator flow
8. view `feeds.json`, scorecard, and consumer-ingestion records

### Docs/Asset Rules

- Every doc must match the current codebase exactly.
- Every tutorial must be executable from the committed repo.
- Use generated diagrams and images for docs assets where available.
- If image generation is unavailable, fall back to Mermaid and document the blocker.

### Phase 10 Exit Criteria

Phase 10 is complete only when:

- README reflects current Phase 8 + hardening reality
- tutorials are runnable as written
- agency demo flow is documented and executable
- deployment guide is concrete enough for a small agency pilot
- docs assets exist and match the current architecture/surface

## Phase 11 — Compliance Evidence and Optional External Integrations

### Goal

Create truthful technical-readiness evidence for California-facing compliance expectations and optional external backend integrations.

### Scope

1. Produce a technical-readiness mapping to current California transit data expectations.
2. Distinguish code-complete vs deployment-required vs third-party-confirmation-required evidence.
3. Review and, where appropriate, wire optional external integrations mentioned in repo dependency docs.
4. Keep Open Transit RT as source of truth even when optional predictors are integrated.

### Required Work

- Add an evidence checklist that distinguishes:
  - code-complete
  - deployment-required
  - external-consumer-confirmation-required
- Review `docs/dependencies.md` and identify each external repo/tool originally mentioned.
- For each intended external integration:
  - either wire and smoke-test it through the documented adapter boundary
  - or explicitly document why it remains optional/deferred
- If integrating an optional predictor backend such as TheTransitClock:
  - keep it behind `PredictionAdapter` only
  - do not let it become the source of truth
  - add smoke/integration or fixture-based contract coverage
  - document setup, failure behavior, replacement strategy, inputs, and outputs

### Phase 11 Exit Criteria

Phase 11 is complete only when:

- the compliance/readiness checklist is truthful and evidence-based
- docs clearly separate technical capability from deployment proof and consumer acceptance
- optional external integrations are either tested or explicitly deferred with reasons
- no compliance or ingestion claims are made without evidence

## Expected Deliverables

- updated `README.md`
- new tutorial docs
- docs assets generated with the image-generation tool, or Mermaid fallback with blocker documented
- updated `docs/dependencies.md`
- updated `docs/current-status.md`
- updated `docs/handoffs/latest.md`
- `docs/handoffs/production-closure.md`
- `docs/handoffs/docs-demo.md`
- `docs/handoffs/compliance-evidence.md`

## Recommended Read Order For Future Codex Runs

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/latest.md`
4. this document
5. `docs/dependencies.md`
6. `docs/decisions.md`
7. relevant tutorial/checklist docs
