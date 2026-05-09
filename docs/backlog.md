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

## Phase 44 — Telemetry Simulator And Device Trial

- Complete for the synthetic-only local/reference simulator scope. The repo now
  has `cmd/telemetry-simulator`, `scripts/telemetry-simulator.sh`,
  `make telemetry-simulator`, synthetic fixtures for on-route, stale,
  out-of-order, unknown-device, low-quality GPS, after-midnight, and
  block-transition scenarios, optional post-ingest DB-backed matcher/Vehicle
  Positions diagnostics, a tutorial, a Phase 44 status doc, a handoff, focused
  tests, and validation checks.
- Remaining later work: real device trials, real agency-approved telemetry,
  authorized vendor AVL integration evidence, and real-world realtime quality
  metrics only when retained claim-specific artifacts exist.

## Phase 45 — GTFS Quality Triage Loop

- Complete for the private authenticated static GTFS quality triage scope. The
  repo now has `/admin/operations/gtfs-quality`, bounded canonical
  MobilityData static validator triage, Open Transit RT internal import
  validation context, admin-only static rerun, strict browser field rejection,
  and a Phase 45 handoff.

## Phase 46 — Validator Automation And Health Gates

- Complete for the private local/reference validator-health diagnostic scope.
  The repo now has `/admin/operations/validation-health`,
  `/admin/operations/validation-health.json`, `scripts/validator-health.sh`,
  `make validator-health`, deployment-doctor summary integration, fixed
  four-feed health rows, strict admin-only `run_all`, bounded output, and a
  Phase 46 handoff.
- Remaining later work: production-owned scheduling, alerting, final-root
  validator evidence, and consumer-facing packet workflows only when retained
  claim-specific evidence exists. Phase 46 itself is not evidence, does not
  block publishing, and does not change consumer statuses.

## Phase 47 — Self-Hosted Operations Notifications

- Complete for the private local/reference notification-summary scope. The
  repo now has `scripts/operations-notify.sh`, `make operations-notify`, a
  tutorial, a Phase 47 reference, a handoff, validation checks, and focused
  tests for safe source parsing, no-send behavior, redaction, strict mode,
  source/output size caps, destination privacy, and docs wording.
- Remaining later work: deployment-owned alerting or monitoring integrations
  only as private diagnostics unless retained claim-specific evidence supports
  a narrower external proof path. Phase 47 itself sends nothing, creates no
  evidence, does not block publishing, and does not change consumer statuses.

## Phase 48 — AVL Adapter Runtime Path

- Complete for the private `/v1/telemetry` send-mode scope. The runtime path
  remains authenticated telemetry ingest only, with redacted `.cache`
  diagnostics and no named vendor, consumer, evidence, compliance, or
  production-readiness claim.

## Phase 49 — External Predictor Runtime Adapter

- Complete for the optional disabled-by-default generic HTTP predictor adapter
  boundary. The repo now has shared Trip Updates adapter factory/config
  validation, `external_http`, `external_http_shadow`, strict endpoint URL and
  token-env validation, sanitized external request/response DTOs, redacted
  diagnostics, valid-empty-feed failure behavior, and non-coupling coverage.
- Remaining later work: any named predictor such as TheTransitClock still
  requires a separately approved runtime phase with dependency/license review
  and claim-specific evidence before compatibility or ETA-quality claims.

## Phase 50 — Realtime Quality Backtesting

- Complete for the approved private local diagnostics scope. The repo now has
  a versioned observed-stop-event and prediction-sample backtesting library,
  `cmd/realtime-quality-backtest`, synthetic public-safe fixtures under
  `testdata/realtime-quality-backtest`, and `make realtime-quality-backtest`.
- Remaining later work: real-world observed arrival/departure quality review
  only when approved inputs and claim-specific evidence boundaries exist.

## Phase 51 — Operations Reliability And SLO Readiness

- Complete for the approved private operations reliability diagnostic scope.
  The repo now has authenticated GET-only `/admin/operations/reliability` and
  `/admin/operations/reliability.json`, bounded Vehicle Positions health
  persistence into the existing `feed_health_snapshot` table,
  `scripts/operations-reliability.sh`, `make operations-reliability`, focused
  runtime/script/model tests, and a Phase 51 handoff.
- Remaining later work: deployment-owned scheduling, production monitoring,
  alert delivery, backup/restore proof, availability evidence, and SLO/SLA
  workflows only when a later approved phase defines retained claim-specific
  evidence. Phase 51 itself creates no evidence, sends no notifications,
  changes no consumer statuses, adds no public route, adds no migrations, and
  makes no compliance, production-readiness, SLA, uptime, hosted SaaS, agency
  adoption, consumer acceptance, vendor compatibility, or production-grade ETA
  claim.

## Phase 52 — Final Public Root Evidence Workflow

- Complete blocker-only for the approved guarded final-root workflow scope.
  The repo now has final-root evidence templates,
  `scripts/collect-final-root-evidence.sh`,
  `scripts/audit-final-root-evidence.sh`, Make targets, validation scaffolding,
  and local-only script tests.
- Closure outcome: no real final root and no real redacted approval artifact
  were available in repo evidence, so no real final-root evidence was retained,
  `docs/evidence/captured` remained unchanged, prepared consumer packets were
  not refreshed, and consumer statuses remained `prepared`.
- Remaining later work: collect and audit real retained final-root evidence
  only when a real final root and public-safe redacted approval artifact exist.

## Phase 53 — Authorized Consumer Submission Execution

- Complete blocker-only for the approved authorized consumer submission
  execution scope.
- Closure outcome: no local operator authorization artifact, official target
  path verification artifact, or target-originated artifact exists, so no target
  was selected, no consumer or aggregator was contacted, no portal was automated
  or scraped, no submission was made, no artifact was added, and all seven
  consumer and aggregator targets remain `prepared`.
- Phase 53 did not change `docs/evidence/consumer-submissions/status.json`,
  current target records, target artifact directories, or
  `docs/evidence/captured`.
- Remaining later work: execute a target-specific submission only when retained
  authorization, official target path verification, and target-originated or
  operator-retained submission evidence exist for one named target.
