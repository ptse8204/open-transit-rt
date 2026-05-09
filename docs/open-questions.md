# Open Questions

These questions do not block the next phase.

## Deployment

- Answered in Phase 36 for the current reference path: document the
  OCI/OCL-style single-server pilot pattern first as the self-hosted reference
  deployment under `docs/deployment/`.
  Managed container platforms and Kubernetes remain later options.
- Answered in Phase 40 for the current guided evaluation path: use
  `docs/tutorials/self-hosted-operator-trial.md` to run one local/reference
  trial across reference deployment prep, reusable onboarding, readiness
  review, validators, and the synthetic AVL dry-run without creating evidence.
- Answered in Phase 41 for current diagnostics: use `make operator-smoke` when
  a local/reference app is running and `make support-bundle` when a maintainer
  needs redaction-safe diagnostics. These outputs are private diagnostics, not
  evidence packets.
- Answered in Phase 42 for current reference deployment diagnostics: use
  `make deployment-doctor` for a read-only OCI/OCL-style reference deployment
  doctor. Default mode exits `0` in a local checkout without deployment env
  while reporting blockers/skips/unavailable checks; `STRICT_DOCTOR=true` is
  the mode that fails on blockers. Output is private diagnostics, not evidence.
- Answered in Phase 43 for current operator setup/readiness diagnostics: use
  `/admin/operations/checklist` or `/admin/operations/checklist.json` from the
  authenticated Operations Console to see grouped private diagnostics,
  heuristic metadata/URL labels, and next actions. The checklist is not
  evidence, not compliance proof, not agency approval, not consumer
  acceptance, and not production readiness.
- Answered in Phase 44 for current telemetry simulator diagnostics: use
  `make telemetry-simulator` to send synthetic events through real
  device-token auth and `POST /v1/telemetry`. Use `RUN_MATCHER=true` for
  optional private DB-backed matcher and Vehicle Positions debug diagnostics.
  Output is private diagnostics, not evidence, not real vendor compatibility,
  not production AVL reliability, not real realtime data, not production-grade
  ETA proof, and not CAL-ITP/Caltrans compliance proof.
- Answered in Phase 47 for current notification drafts: use
  `make operations-notify` to write a private local draft from existing
  validator-health and deployment-doctor summaries. The draft is not sent, not
  evidence, not a compliance gate, not production health proof, and not
  consumer acceptance.
- Answered in Phase 51 for current reliability diagnostics: use authenticated
  `/admin/operations/reliability` or `/admin/operations/reliability.json` for
  DB-backed private feed/incident rollups, and `make operations-reliability`
  for a local private summary from safe existing diagnostics. These outputs
  are not evidence, not SLA or uptime proof, not production readiness, not
  compliance, and not consumer acceptance.
- Answered in Phase 52 for final-root workflow mechanics: use
  `make collect-final-root-evidence` and `make audit-final-root-evidence` for
  guarded final public root evidence. Without a real final root and public-safe
  redacted approval artifact, the workflow must close blocker-only under
  ignored `.cache` storage and leave `docs/evidence/captured` unchanged.
- Answered in Phase 53 for authorized consumer submission execution: without a
  retained operator authorization artifact, official target path verification
  artifact, and target-originated or operator-retained submission artifact, no
  target is selected, no submission is made, no artifact is added, and all seven
  consumer and aggregator targets remain `prepared`.
- Answered in Phase 54 for official-source mapping: current Caltrans /
  Cal-ITP and FTA source mappings were refreshed as docs-only requirements
  context. The refresh does not create evidence, change consumer statuses,
  prove compliance, prove consumer acceptance, prove final-root readiness,
  prove hosted SaaS availability, prove production readiness, prove vendor
  compatibility, create SLA/uptime evidence, or prove production-grade ETA
  quality.
- Answered in Phase 55 for local compliance/readiness packets: use
  `make generate-compliance-evidence-packet` to write ignored `.cache`
  blocker/draft packets and `COMPLIANCE_PACKET_DIR=<packet> make
  audit-compliance-evidence-packet` to enforce claim, redaction, file-set,
  tracker, artifact-directory, and misleading-wording guards. These packets are
  summaries only; they do not create retained evidence, contact consumers,
  fetch live feeds, change consumer statuses, or prove compliance.
- Answered in Phase 56 for current multi-agency hosting boundaries: use
  validated `/public/agencies/{agency_id}/...` public feed routes and
  `make multi-agency-hosting` for private local route/proxy diagnostics. These
  diagnostics do not create evidence, do not prove hosted SaaS or production
  multi-tenant readiness, and tenant restore into a shared live database
  remains blocked.
- Answered in Phase 57 for current release packaging: use
  `make release-package` and `RELEASE_PACKAGE_DIR=<dir> make
  audit-release-package` to create and audit local `.cache` source packages
  with checksums, SBOM, and provenance metadata. These packages do not publish
  artifacts, push images, create retained evidence, prove production readiness,
  or prove hosted service availability.
- Answered in Phase 58 for optional vendor-equivalent materials: use
  `docs/vendor-equivalent-pack/` and `make audit-vendor-equivalent-pack` for
  template-only BYOD, support, KPI, and procurement review. These templates do
  not prove marketplace approval, paid support, vendor compatibility, hardware
  certification, SLA/uptime, production readiness, compliance, agency adoption,
  consumer acceptance, or hosted service availability.
- Answered in Phase 59 for real pilot closeout: no retained real pilot
  authorization plus kickoff, operations, feedback, and decision artifact set
  was available in the repository, so Phase 59 closed blocker-only. Existing
  OCI/local pilot packets remain earlier-scope only, no evidence was created,
  and all seven consumer and aggregator targets remain `prepared`.
- Answered in Phase 60 for final claim review: use
  `make audit-final-claim-review` for a local read-only scan of bounded
  public/status docs, required Phase 60 sections, unsafe private strings, the
  exact seven-target prepared-only consumer tracker state, and README-only
  consumer artifact directories. This audit does not create evidence, fetch
  live feeds, contact external parties, or support any stronger public launch,
  compliance, agency, consumer, hosted service, SLA/uptime,
  production-readiness, vendor, marketplace, or ETA-quality claim.
- Should production HTTPS termination be owned by this repo or by deployment infrastructure?
- Phase 10 documents the current pilot path as Postgres/PostGIS through Compose plus deployment-owned Go service process management and a TLS reverse proxy. A fully packaged app-container or Kubernetes path remains open.

## Authentication

- Which auth provider should be used when admin/operator login is implemented?
- Should device credentials use opaque bearer tokens only, or support signed device JWTs later?
- Should Phase 1's local debug `/v1/events` endpoint be removed, protected behind admin auth, or moved under a separate admin route before production deployment?

## Validation Tooling

- Answered for the repo-supported path in Phase 9: MobilityData GTFS Validator `v7.1.0` and Docker-backed MobilityData GTFS-RT validator image digest in `tools/validators/validators.lock.json`.
- Answered for Phase 54 source mapping: current official guidance still maps
  compliant GTFS Schedule and GTFS Realtime wording to regular canonical
  validation with no errors, but validation success remains separate from
  consumer acceptance and compliance evidence for a specific deployment.
- Open for deployments: whether a production environment should use the repo-supported Docker-backed GTFS-RT wrapper or document an equivalent checksum/digest contract for a native executable.

## Static GTFS Publication

- Phase 8 serves `/public/gtfs/schedule.zip` on demand from active published GTFS tables with deterministic bytes. Should a later phase add materialized ZIP caching or checksum comparison against original uploaded ZIP bytes?
- Answered for Phase 37 onboarding: reusable local/reference onboarding stores
  the downloaded source ZIP and checksum under ignored `.cache/` storage and
  compares public-safe source/fetched schedule summaries. This is an
  onboarding verification, not materialized production ZIP caching.

## Prediction Backends

- Phase 7 chose an internal deterministic predictor as the first real Trip Updates adapter. Should a later phase add TheTransitClock as an alternate adapter, and what deployment profile should own it?
- Phase 38 documented the required external predictor lifecycle: candidate
  review, dependency/license review, adapter contract tests, shadow/dry-run
  evaluation, output validation, failure fallback, and retained evidence before
  stronger ETA claims.
- Phase 49 answered the generic runtime-boundary question for the current
  scope: deployments may opt into `external_http` or `external_http_shadow`
  with a sanitized fixed-path HTTP sidecar behind `internal/prediction.Adapter`.
  Named predictor integration, external process ownership, deployment
  packaging, and ETA-quality evidence remain open for later approved phases.
- What quality threshold should be required before claiming production-grade ETA quality rather than conservative schedule-deviation predictions?
- Answered for Phase 50 private diagnostics: use
  `cmd/realtime-quality-backtest` with versioned local observed-stop-event and
  prediction-sample files to compute aggregate MAE, median, p90, coverage,
  stale/missing/withheld counts, and diagnostic maturity gates by route and
  agency-local time period. This is local engineering diagnostics only, not a
  database retention model or production ETA proof.

## Alerts

- Phase 8 chose both operator-authored and system-derived Alerts for canceled-trip reconciliation. Should later alert workflows include richer affected-route/stop/time selectors, multilingual text, and full operator UI review before publication?

## GTFS Studio

- Phase 5 resolved the first UI entity scope: minimal operational forms for agency metadata, routes, stops, trips, stop_times, calendars, calendar_dates, shape points, and frequencies.
- Should a later Studio UI add map editing for shape points and timetable-design interactions for stop_times, or keep those as row editors?
- Should draft publish require canonical MobilityData validator success in all environments once canonical validation is wired, or only production-like environments?
- Phase 46 answer for local/reference diagnostics: validator health does not
  block publishing automatically. It reports private health states and next
  actions only. Any future production publish gate needs a separate
  deployment-owned policy and evidence plan.
