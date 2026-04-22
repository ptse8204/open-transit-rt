# Codex Production Closure Prompt

Read and obey `AGENTS.md` first.

You are not starting a brand-new feature phase. You are closing the repo to production-directed quality in three short phases:

- **Phase 9 — Production Closure**
- **Phase 10 — Docs, Tutorials, Deployment, and Demo**
- **Phase 11 — Compliance Evidence and Optional External Integrations**

## Mission

Close the remaining production-readiness gaps, make the repo presentable and understandable to agencies, and produce evidence that supports a truthful CAL-ITP / Caltrans-aligned technical-readiness claim.

## Non-Negotiable Truthfulness Rules

- Do not claim production readiness unless implementation, tests, docs, and smoke checks support it.
- Do not claim CAL-ITP / Caltrans compliance unless the repo plus deployment evidence supports the required checklist items.
- Do not claim external integration unless it is actually wired and tested.
- Do not let docs drift from code.

## Required Repo/Requirements Context

- Open Transit RT remains the source of truth for GTFS management, telemetry, assignments, audits, Vehicle Positions, publication, and compliance workflow records.
- External predictors such as TheTransitClock remain optional backends behind `PredictionAdapter`.
- Preserve stable public feed paths and GTFS lifecycle boundaries.
- Keep GTFS import, GTFS Studio, matcher, Trip Updates, Alerts, and publication behavior coherent.

## Phase 9 — Production Closure

### Goals

1. Finish any remaining hardening and evidence gaps.
2. Keep public protobuf endpoints anonymous and stable.
3. Keep admin/debug/mutation endpoints protected.
4. Verify all major subsystems still work together after hardening.

### Required Work

- complete realtime validator artifact derivation for `/admin/validation/run`
- pin validator tooling in CI/dev/prod setup
- tighten and verify admin auth/roles everywhere
- finish structured request logging and request IDs
- finish metrics and readiness depth
- add strong smoke/e2e coverage
- verify GTFS import, Studio, matcher, and feed publication interplay

### Required Checks

- `gofmt -w ./cmd ./internal`
- `go mod tidy`
- `go test ./...`
- `make validate`
- `make test-integration`
- `make smoke`
- `docker compose -f deploy/docker-compose.yml config`
- `git diff --check`

## Phase 10 — Docs, Tutorials, Deployment, and Demo

### Goals

1. Rewrite `README.md` to current reality.
2. Add quickstart/tutorial/deployment docs.
3. Package a simple agency demo flow.
4. Improve repo presentation and onboarding.
5. Use generated diagrams/assets in docs.

### Required Deliverables

- `README.md`
- `docs/tutorials/local-quickstart.md`
- `docs/tutorials/deploy-with-docker-compose.md`
- `docs/tutorials/agency-demo-flow.md`
- `docs/tutorials/production-checklist.md`
- `docs/tutorials/calitp-readiness-checklist.md`
- `docs/assets/README.md`

### Agency Demo Flow Must Cover

1. bootstrap DB
2. import sample GTFS
3. publish schedule
4. ingest authenticated telemetry
5. fetch public Vehicle Positions / Trip Updates / Alerts protobuf feeds
6. show protected admin/debug routes
7. run validator flow
8. view `feeds.json`, scorecard, and consumer-ingestion records

### Critical Docs Rule

Every doc must match the actual codebase and commands. Every tutorial must be executable from the repo as committed.

## Phase 11 — Compliance Evidence and Optional External Integrations

### Goals

1. Produce a truthful technical-readiness mapping to California-facing transit data expectations.
2. Separate code-complete from deployment-required and third-party-confirmation-required evidence.
3. Review and, where appropriate, wire external repos/tools mentioned in `docs/dependencies.md`.
4. Keep Open Transit RT as source of truth even when optional predictors are integrated.

### Required Work

- add an evidence checklist doc
- review every external repo/tool originally mentioned in repo dependency docs
- either wire and smoke-test intended integrations or explicitly document why they remain deferred
- if integrating optional predictors such as TheTransitClock:
  - keep them behind `PredictionAdapter`
  - do not let them become the source of truth
  - add smoke/integration or fixture-based contract coverage
  - document setup, inputs, outputs, failure behavior, and replacement strategy

## Deliverables

- updated `README.md`
- new tutorial docs
- docs assets generated via image-generation or Mermaid fallback
- updated `docs/dependencies.md`
- updated `docs/current-status.md`
- updated `docs/handoffs/latest.md`
- `docs/handoffs/production-closure.md`
- `docs/handoffs/docs-demo.md`
- `docs/handoffs/compliance-evidence.md`

## Acceptance Criteria

- repo can be demoed end-to-end with documented commands
- tutorials are runnable as written
- hardening gaps are closed with tests
- public/private surface is clear
- compliance/readiness checklist is truthful and evidence-based
- external integrations are either tested or explicitly deferred with reasons
- no stale docs remain

If a full closeout is unrealistic in one pass, stop and create the exact next phase handoff instead of pretending closure.
