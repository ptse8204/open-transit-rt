# Backlog

This backlog is organized by phase. `docs/phase-plan.md` remains the phase contract.

## Phase 1 — Durable Telemetry Foundation

- Complete. Durable telemetry ingest, DB readiness, duplicate/out-of-order handling, parsed payload JSONB storage, agency-scoped debug listing, and DB-backed tests are implemented.

## Phase 2 — Deterministic Trip Matching

- Complete. The matcher resolves agency-local service days, handles after-midnight and frequency cases, persists explicit unknown rows, respects manual overrides, records reasons/degraded state/incidents, and has unit plus DB-backed tests.

## Phase 3 — Vehicle Positions Production Feed

- Complete. Vehicle Positions protobuf and JSON debug output are DB-backed, generated from the same snapshot model, preserve unknown/stale behavior, and have unit plus DB-backed tests.

## Phase 4 — GTFS Import And Publish

- Complete. GTFS ZIP import, internal validation, transactional feed-version publish, failed-import report storage, rollback-safe activation, block preservation, shape-line construction, and active feed switching tests are implemented.

## Phase 5 — GTFS Studio

- Complete. Typed GTFS draft tables, draft CRUD, minimal server-rendered Studio UI, soft discard, cloned-draft provenance, draft publish traceability, and direct shared validation/activation publishing are implemented.

## Phase 6 — Trip Updates And Alerts Architecture

- Complete. Trip Updates adapter contracts, no-op adapter, diagnostics persistence, stable empty Trip Updates endpoints, stable empty Alerts endpoints, and non-coupling tests are implemented.

## Phase 7 — Prediction Quality And Operations

- Complete for the first conservative operations scope. The deterministic Trip Updates adapter, override lifecycle repository, audit logging, prediction review queue, cancellation linkage signal, deadhead/layover suppression, disruption withholds, and coverage metrics are implemented.
- Remaining later work: production-grade ETA quality, backtesting, full operator UI, vehicle swap UI/workflow, and richer detour/short-turn prediction behavior.

## Phase 8 — Compliance And Consumer Workflow

- Complete for the first publication/compliance layer. Persisted Alerts, public Alerts feeds, canceled-trip alert reconciliation, public schedule ZIP publication, public feed metadata, license/contact workflows, consumer ingestion records, marketplace-gap records, validator report execution/recording, and compliance scorecard snapshots are implemented.
- Remaining later work: richer operator UI, production observability/SLO rollups, deployment-specific validator evidence, and external consumer acceptance evidence.

## Phase 9 — Production Closure

- Complete for the current codebase surface. Validator execution hardening, pinned validator install/check workflow, admin JWT/cookie auth, DB-backed roles, device token binding/rebinding, assignment current-row race protection, request IDs/logging, metrics gating, readiness checks, and smoke coverage are implemented.
- Remaining later work: hosted login/SSO, server-side admin JWT `jti` replay tracking, production SLO dashboards, and deployment-specific monitoring/alerting assets.

## Phase 10 — Docs, Tutorials, Deployment, And Demo

- Complete for the current codebase surface. README, local/deployment/demo/checklist tutorials, executable agency demo flow, docs assets, bootstrap output polish, and truthful CAL-ITP/Caltrans-aligned wording are implemented.
- Remaining later work: deployment-specific proof for any stronger readiness or consumer-ingestion claims.

## Phase 11 — Compliance Evidence And Optional External Integrations

- Complete for the selected evidence-only path. The Phase 11 evidence checklist, dependency reality table, README/tutorial truthfulness links, current-status update, and Phase 11 handoff are implemented.
- Remaining later work: deployment evidence hardening, real HTTPS feed proof, production validator records, monitored operations evidence, scorecard export, and third-party submission or acceptance records where they exist.

## Phase 35 — README And Roadmap Realignment

- Complete for docs-only README and roadmap realignment. The root README is
  restored as the Open Transit RT product front door, roadmap/status docs now
  make self-hosted agency reuse the default continuation path.

## Phase 36 — OCI/OCL Reference Deployment Productization

- Complete for docs-only reference deployment productization. The repo now has
  a self-hosted OCI/OCL-style reference deployment guide, placeholder-only env
  example, smoke checklist, deployment index, and closed phase handoff.
- Remaining later work: expand the integration adapter kit
  and make CAL-ITP-style readiness visible in product workflows without
  overclaiming compliance.

## Phase 37 — Reusable Agency Onboarding Flow

- Complete for the opt-in local/reference onboarding scope. The repo now has
  `scripts/agency-pilot-onboard.sh`, `make agency-pilot-up`, explicit Compose
  env interpolation defaults, reusable GTFS download/import/five-path fetch
  checks, schedule identity summary verification, metadata placeholder
  warnings, best-effort validator/blocker reporting, and a closed phase
  handoff.
- Remaining later work: Phase 39 CAL-ITP-style readiness workflow, without
  overclaiming compliance or consumer acceptance.

## Phase 38 — Integration Adapter Kit

- Complete for the navigation and conformance scope. The repo now has the
  central integration adapter kit, synthetic AVL fixture manifest, neutral
  dry-run fixture examples, CLI boundary wording, focused adapter tests, and
  dependency docs noting synthetic dry-run/developer support only.

## Phase 39 — CAL-ITP-Style Readiness Workflow

- Complete for the product-facing readiness workflow scope. The authenticated
  Operations Console now has `/admin/operations/readiness` with ten
  evidence-bounded readiness rows and next actions.

## Phase 40 — Guided Self-Hosted Operator Trial

- Complete for the docs/navigation guided trial scope. The repo now has
  `docs/tutorials/self-hosted-operator-trial.md`, a Phase 40 status doc,
  handoff, navigation updates, and validation checks for those files.

## Phase 41 — Operator Smoke And Support Bundle

- Complete for the local/reference diagnostic tooling scope. The repo now has
  `scripts/operator-smoke.sh`, `scripts/support-bundle.sh`, `make
  operator-smoke`, `make support-bundle`, an operator tutorial, a Phase 41
  handoff, and validation checks for those files.

## Phase 42 — Reference Deployment Doctor

- Complete for the read-only reference deployment diagnostic scope. The repo
  now has `scripts/deployment-doctor.sh`, `make deployment-doctor`,
  `docs/deployment/reference-deployment-doctor.md`, a Phase 42 handoff, and
  validation checks for those files.

## Phase 43 — Operator UX Setup V2

- Complete for the private authenticated Operations Console checklist scope.
  The repo now has `/admin/operations/checklist` and
  `/admin/operations/checklist.json`, a shared deterministic checklist model,
  grouped setup/feeds/validation/telemetry/operations/consumer_workflow rows,
  neutral statuses, heuristic metadata and URL labels, repo-relative docs
  links, explicit false claim flags, local Caddy fallback hardening, and a
  Phase 43 handoff.
- Remaining later work: only future optional retained-evidence paths such as
  final-root proof, real agency pilot evidence, real device/vendor AVL
  evidence, consumer submissions, or real-world realtime quality evidence when
  claim-specific artifacts exist.
