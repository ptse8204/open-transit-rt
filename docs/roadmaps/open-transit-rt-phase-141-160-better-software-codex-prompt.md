# Open Transit RT — Phase 141–160 Better Software Roadmap Prompt

## Goal

Use a Codex Goal-driven workflow to turn Open Transit RT from a browser-first self-hosted evaluation product into a stronger open-source GTFS / GTFS-Realtime operations platform.

**Goal statement:** improve the actual product, not only the copy, site, or roadmap. Build durable software capabilities that help small agencies, civic technologists, operators, and deployment owners import and validate GTFS, publish useful GTFS-Realtime feeds, connect vehicle data safely, understand realtime failures, operate a self-hosted deployment, and prepare external sharing without unsupported claims.

**Completion condition:** phases 141 through 160 are complete, each phase is committed separately, actual code/tests/docs/scripts/UI are changed where needed, validation is run and recorded, protected evidence paths and consumer statuses remain unchanged, and the repo has a clearer path from local/self-hosted evaluation toward reliable agency operations and GTFS-RT adoption.

**Non-goals:** do not collect evidence, contact agencies/vendors/consumers, move consumer statuses, claim compliance/adoption/production readiness/vendor compatibility/SLA/ETA quality, or satisfy this roadmap with markdown-only phase ledgers. Documentation is allowed only when it unlocks or explains actual product behavior.

## Repository Context

Repository: <https://github.com/ptse8204/open-transit-rt>

Ground all decisions in current source code, tests, scripts, and route behavior. Do not rely only on the committed Phase 133–140 prompt or handoff markdown. Re-check the repo before each phase.

Read first:

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/product-language-guide.md`
- `docs/open-transit-rt-phase-133-140-codex-prompt.md`
- `docs/roadmaps/external-connector-runtime-integration/README.md`
- `docs/roadmaps/external-connector-runtime-integration/phase-plan.md`
- `docs/roadmaps/post-rc2-browser-first-product/closeout.md`
- `docs/integration-adapter-kit.md`
- `docs/external-connection-readiness.md`
- `docs/connectors/catalog.md`
- `docs/connectors/plugin-contract.md`
- `docs/tutorials/device-avl-integration.md`
- `docs/tutorials/external-adapter-conformance.md`
- `docs/deployment/oci-reference-deployment.md`
- `docs/evidence/redaction-policy.md`
- `Makefile`
- `.github/workflows/test.yml`
- `.github/workflows/release-gates.yml`

Inspect implementation areas as needed:

- `cmd/agency-config/`
- `cmd/adapter-conformance/`
- `cmd/avl-vendor-adapter/`
- `cmd/gtfs-import/`
- `cmd/telemetry-ingest/`
- `cmd/feed-vehicle-positions/`
- `cmd/feed-trip-updates/`
- `cmd/feed-alerts/`
- `internal/compliance/`
- `internal/connectors/`
- `internal/gtfs/`
- `internal/feed/`
- `internal/prediction/`
- `internal/realtimequality/`
- `internal/telemetry/`
- `internal/state/`
- `examples/connectors/`
- `testdata/adapter-conformance/`
- `testdata/gtfsrt-conformance/`
- `testdata/telemetry-simulator/`
- `site/`
- `wiki/`
- `scripts/`

## Quality Priority

Quality is more important than speed. Prefer correctness, maintainability, domain correctness, accessibility, durable tests, and useful product outcomes over quick patches.

Do not create artifact slop. Every phase must produce actual product value: code, tests, scripts, UI changes, validation harnesses, or operator workflows that can be exercised.

## Master Agent Operating Mode

Use one high-effort master Codex agent for the actual software engineering work. The master owns planning, implementation, validation, review, final decisions, and commits.

Do not split planning, execution, and review into separate full-context agents by default. Use bounded read-only scouts only when independent critique materially improves quality.

Allowed scouts:

- `plan-risk-scout`: critiques architecture risks, user-flow gaps, test gaps, and likely scope drift before implementation.
- `domain-scout`: checks GTFS, GTFS-RT, realtime transit, connector, or interoperability concerns when a phase touches domain behavior.
- `diff-review-scout`: reviews final diff before commit for regressions, security, accessibility, maintainability, and claim-boundary mistakes.
- `test-gap-scout`: checks whether tests cover edge cases and failure modes.

Scout rules:

- At most 1–2 scouts per phase unless a concrete reason exists.
- Scouts are read-only. They may inspect files, diffs, tests, and command output. They must not edit files, format code, install dependencies, commit, change configuration, spawn more agents, or perform broad exploration without a narrow reason.
- Scout output must contain only: Scope inspected; Files/symbols reviewed; Findings; Evidence; Recommended action; Confidence level.
- The master synthesizes scout findings and makes final decisions.

## Skills Requirement

Before UI, public-site, tutorial, copy, accessibility, or layout changes, load any available web design / UX / product copy / accessibility skill. In this project, the repo roadmap references a local web-design-engineer skill path; if it exists in the Codex environment, use it. If unavailable, do not fake it; apply `docs/product-language-guide.md` and the concrete requirements in this prompt.

Before GTFS / GTFS-RT / realtime / connector changes, use a domain-scout or equivalent domain review if the change affects feed correctness, schedule interpretation, protobuf output, connector input validation, prediction behavior, or validation semantics.

## Hard Boundaries

Do not:

- write to `docs/evidence/**`;
- change `docs/evidence/consumer-submissions/status.json`;
- move consumer targets beyond `prepared`;
- contact agencies, vendors, consumers, portals, live validators, or live connector services;
- use real private agency data, private vendor payloads, credentials, secrets, webhook URLs, private endpoint URLs, database URLs, or private paths;
- expose arbitrary command execution in the browser;
- make private admin routes public;
- weaken auth, authorization, CSRF, no-store, same-origin, path-safety, or redaction controls;
- claim CAL-ITP/Caltrans compliance, production readiness, agency adoption, consumer acceptance, final-root readiness, hosted service availability, vendor compatibility, hardware certification, SLA/uptime, production AVL reliability, production-grade ETA quality, or real-world ETA accuracy.

All examples, fixtures, connector tests, and screenshots must be local, synthetic, public-safe, redacted, and no-contact by default.

## Roadmap

Continue through every phase unless a real blocker is reached. Each phase must have a separate commit. Re-check whether a phase has already been partially implemented before duplicating work.

### Phase 141 — Current Product Re-Audit And Roadmap Merge

**Purpose:** Establish the actual post-140 baseline and merge the UI reset, external connector runtime plan, and product-quality goals into one implementation track.

**User value:** Future work stops bouncing between UI cleanup and disconnected connector plans.

**Technical scope:**

- Inspect current `main`, `stable`, `site/`, `gh-pages`, route registry, Make targets, audits, and public docs.
- Confirm which Phase 133–140 items are done, partially done, or still weak.
- Confirm which external connector runtime phases are still open.
- Add a concise implementation status page only if necessary, but do not make this phase docs-only: add at least one lightweight audit or test that catches a current roadmap drift risk.
- Identify a single product-quality backlog grouped by: operator workflow, GTFS data quality, GTFS-RT usefulness, connectors, deployment/observability, security/redaction, and release gates.

**Likely files:** `docs/current-status.md`, `docs/handoffs/latest.md`, `docs/roadmap-status.md`, `scripts/*audit*`, `Makefile`, targeted tests.

**Validation:** `git diff --check`, `make check`, `make audit-product-language`, `make audit-ui-layout`, `make audit-final-claim-review`, `scripts/check-consumer-tracker.sh`.

**Risk review:** Avoid creating a new phase ledger without executable checks. Avoid reopening closed evidence tracks.

**Expected commit:** `Phase 141 -- Checkpoint 000001: merge product quality roadmap into executable baseline`

### Phase 142 — Remove Remaining Primary-HTML Audit Vocabulary And Card Debt

**Purpose:** Finish the UI cleanup so primary HTML never feels like an internal audit surface.

**User value:** Operators see task flow, not raw flags, repeated caveats, or card piles.

**Technical scope:**

- Audit rendered primary Operations Console HTML, not just source strings.
- Remove raw claim-flag tables from primary HTML entirely; keep raw flags in JSON or tests. If a collapsed advanced safety section remains, make it plain-language and short.
- Further reduce `card-grid` and `status-grid` use on Help, first-run panels, connector pages, and diagnostics pages where tables/accordions/action lists are clearer.
- Convert repeated `Limits` cells into one page-level limits section plus targeted warnings only where needed.
- Strengthen `audit-product-language` and `audit-ui-layout` to render/check primary pages where feasible.

**Likely files:** `cmd/agency-config/operations.go`, `cmd/agency-config/operations_*go`, `cmd/agency-config/*_test.go`, `scripts/audit-product-language.sh`, `scripts/audit-ui-layout.sh`, `docs/product-language-guide.md`.

**Validation:** `go test ./cmd/agency-config`, `make audit-product-language`, `make audit-ui-layout`, `make product-ui-smoke`, `make audit-operations-route-inventory`, `git diff --check`.

**Risk review:** Do not remove safety data from JSON contracts. Do not hide essential warnings behind inaccessible JS.

**Expected commit:** `Phase 142 -- Checkpoint 000001: remove remaining audit vocabulary from primary UI`

### Phase 143 — Operator Issue Center V2

**Purpose:** Turn overlapping diagnostics into one prioritized work queue.

**User value:** Operators can answer: what is broken, who owns it, what should happen next, and where should I click?

**Technical scope:**

- Build or refine a unified private operator issue model with fields: label, severity, owner, current signal, why it matters, next action, route link, source, freshness, and deduplication key.
- Deduplicate feed health, validation, GTFS quality, telemetry, devices, connectors, readiness, reliability, and maintenance signals.
- Surface the issue center on Start and link from relevant pages.
- Add JSON output only if it helps scripts or tests.
- Use plain owner categories: operator, administrator, deployment owner, developer/integrator.

**Likely files:** `cmd/agency-config/operations.go`, `cmd/agency-config/operations_*go`, `internal/compliance/*` if a shared model belongs there, tests.

**Validation:** `go test ./cmd/agency-config`, `go test ./...`, `make product-ui-smoke`, `make audit-product-language`, `make audit-product-acceptance`.

**Risk review:** Read-only only. No feed mutation, evidence writes, external contact, or consumer status changes.

**Expected commit:** `Phase 143 -- Checkpoint 000001: add unified operator issue center`

### Phase 144 — GTFS Import, Diff, And Rollback Usability V3

**Purpose:** Make schedule import/review safer and more useful for small agencies.

**User value:** Operators can see what changed in a schedule, whether it is safe to publish, and how to recover.

**Technical scope:**

- Improve GTFS import preview with route/stop/trip/service-calendar summaries, service-date warnings, required-file status, row-count deltas, checksum comparison, and human-readable blockers.
- Improve active-vs-previous diff and rollback review guidance without adding unsafe rollback execution unless an existing safe mutation path supports it.
- Add bounded samples for affected routes/stops/trips and service calendar changes.
- Keep raw ZIP handling temporary and private.

**Likely files:** `internal/gtfs/`, `cmd/agency-config/operations_gtfs_import.go`, `cmd/agency-config/operations_gtfs_workbench.go`, GTFS tests, `testdata/gtfs/*`.

**Validation:** `go test ./internal/gtfs ./cmd/agency-config`, `go test ./...`, `make check`, `make audit-product-language`.

**Risk review:** Do not corrupt active schedule state. Do not expose raw private paths or ZIP contents beyond bounded review.

**Expected commit:** `Phase 144 -- Checkpoint 000001: improve GTFS import diff and recovery review`

### Phase 145 — GTFS Validation Explanation And Safe Fix Planner V2

**Purpose:** Convert validator output into practical remediation guidance.

**User value:** Less technical staff can understand validation issues and know who should fix them.

**Technical scope:**

- Improve mapping from MobilityData/internal validation notices to plain-language categories, likely owner, affected files, risk level, and safe fix path.
- Add issue grouping for common GTFS problems: calendar gaps, missing required files, bad references, shapes, stop times, frequencies, routes, agency metadata, license/contact metadata.
- Add tests for malformed reports, huge reports, hostile strings, and stale validation result handling.
- Keep validator commands allowlisted and server-owned.

**Likely files:** `internal/compliance/validation*`, `cmd/agency-config/operations_gtfs_quality.go`, `cmd/agency-config/operations_validation_center.go`, tests, fixtures.

**Validation:** `go test ./internal/compliance ./cmd/agency-config`, `make check`, `make audit-product-language`, `make audit-final-claim-review`.

**Risk review:** Do not imply validator-clean means compliance or acceptance.

**Expected commit:** `Phase 145 -- Checkpoint 000001: explain GTFS validation issues with safe fix guidance`

### Phase 146 — Telemetry Ingest And Device Onboarding Hardening V3

**Purpose:** Make real vehicle-data setup safer before named vendor integrations.

**User value:** Administrators can onboard devices/connectors and diagnose bad telemetry without exposing secrets.

**Technical scope:**

- Recheck `/v1/telemetry` request validation, token failures, agency-scope rejection, body size, timestamp validation, duplicate/out-of-order classification, unknown devices, low-quality GPS, and stale handling.
- Improve device/token onboarding UI with one-time token lifecycle, binding state, token rotation guidance, and clear rejection reasons.
- Add redacted ingest diagnostics visible in private browser where safe.
- Keep raw tokens, headers, payloads, and private endpoints out of HTML, logs, screenshots, and diagnostics.

**Likely files:** `cmd/telemetry-ingest/`, `internal/telemetry/`, `internal/devices/`, `cmd/agency-config/operations_devices.go`, tests.

**Validation:** `go test ./cmd/telemetry-ingest ./internal/telemetry ./internal/devices ./cmd/agency-config`, `make adapter-conformance`, `make external-connection-check`, `make audit-final-claim-review`.

**Risk review:** Do not weaken auth. Do not add browser token recovery.

**Expected commit:** `Phase 146 -- Checkpoint 000001: harden telemetry ingest and device onboarding`

### Phase 147 — Vehicle Positions Usefulness And Static-GTFS Interop V2

**Purpose:** Make Vehicle Positions output more understandable and standards-aware.

**User value:** Operators know when Vehicle Positions are publishable, empty, stale, suppressed, or missing schedule context.

**Technical scope:**

- Strengthen Vehicle Positions diagnostics around trip descriptor inclusion/omission, stale telemetry, missing active GTFS, low assignment confidence, unknown vehicles, and agency scope.
- Add or improve GTFS-RT conformance fixtures for Vehicle Positions edge cases.
- Ensure Vehicle Positions remains independent of predictor health.
- Surface plain-language public-feed usefulness on Realtime and Feed Health pages.

**Likely files:** `internal/feed/vehiclepositions/`, `cmd/feed-vehicle-positions/`, `cmd/agency-config/operations_realtime.go`, `testdata/gtfsrt-conformance/`, tests.

**Validation:** `go test ./internal/feed ./cmd/feed-vehicle-positions ./cmd/agency-config`, `make gtfsrt-conformance`, `make smoke`, `make audit-product-language`.

**Risk review:** Do not publish false trip certainty. Unknown is better than false certainty.

**Expected commit:** `Phase 147 -- Checkpoint 000001: improve Vehicle Positions usefulness diagnostics`

### Phase 148 — Trip Updates Prediction, Shadow, And Withheld-Reason V4

**Purpose:** Improve Trip Updates safety, explainability, and external predictor evaluation without overclaiming ETA quality.

**User value:** Operators understand why Trip Updates are missing or withheld and can compare predictor behavior safely.

**Technical scope:**

- Improve deterministic predictor diagnostics and withheld reasons.
- Improve external HTTP predictor shadow-mode summaries: timeout, malformed output, stale output, wrong agency/feed, low confidence, divergence, and fallback.
- Add local stub/fixture coverage where missing.
- Surface concise Trip Updates status in operator issue center, Realtime, Feed Health, and Prediction Lab.

**Likely files:** `internal/prediction/`, `internal/feed/tripupdates/`, `cmd/agency-config/operations_prediction_lab.go`, `examples/connectors/predictor-sidecar-stub/`, tests.

**Validation:** `go test ./internal/prediction ./internal/feed/tripupdates ./cmd/agency-config ./examples/connectors/predictor-sidecar-stub`, `make gtfsrt-conformance`, `make adapter-conformance`, `make audit-final-claim-review`.

**Risk review:** No production-grade ETA quality or real-world ETA accuracy claims.

**Expected commit:** `Phase 148 -- Checkpoint 000001: improve Trip Updates shadow diagnostics and withheld reasons`

### Phase 149 — Alerts And Disruption Workflow V3

**Purpose:** Make service alerts and disruption handling practical for operators.

**User value:** Staff can create/review alerts, understand affected services, and avoid unpaired cancellations.

**Technical scope:**

- Improve alert lifecycle review for drafts, active alerts, indefinite/stale alerts, entity scoping, and cancellation pairing.
- Add disruption templates or guided prompts if existing models support them safely.
- Improve GTFS-RT Alerts conformance fixtures and explanation.
- Link alert issues into the operator issue center.

**Likely files:** `cmd/agency-config/`, `cmd/feed-alerts/`, `internal/feed/alerts/`, alert console code, `testdata/gtfsrt-conformance/`, tests.

**Validation:** `go test ./cmd/agency-config ./cmd/feed-alerts ./internal/feed`, `make gtfsrt-conformance`, `make audit-product-language`.

**Risk review:** No public-feed mutation outside existing alert workflow. No consumer-display claim.

**Expected commit:** `Phase 149 -- Checkpoint 000001: improve alerts and disruption workflow`

### Phase 150 — Connector Runtime Pack V2

**Purpose:** Move connector examples closer to real deployment usefulness while staying synthetic/no-contact by default.

**User value:** Integrators can start from safe CSV, HTTP polling, webhook sidecar, and JSON transform shapes.

**Technical scope:**

- Improve connector manifests and examples for CSV replay, HTTP polling, webhook sidecar, generic JSON transform, prediction sidecar, validator wrapper, and monitoring export.
- Add bounded config validation, dry-run behavior, send-mode gates, timeout behavior, redacted diagnostics, and first-safe-check commands.
- Ensure no named vendor compatibility claims.
- Improve `/admin/operations/connectors*` only where it supports runtime setup clearly.

**Likely files:** `examples/connectors/`, `internal/connectors/`, `cmd/adapter-conformance/`, `cmd/avl-vendor-adapter/`, `docs/connectors/*`, tests.

**Validation:** `make external-connection-check`, `make adapter-conformance`, `make test-connector-examples`, `go test ./cmd/adapter-conformance ./internal/connectors ./examples/connectors/...`.

**Risk review:** No dynamic backend plugin loading. No browser command execution. No live sends by default.

**Expected commit:** `Phase 150 -- Checkpoint 000001: improve connector runtime pack and examples`

### Phase 151 — Connector Health And Setup UI V2

**Purpose:** Make connector setup review understandable in the browser.

**User value:** Operators and integrators can see connector readiness without decoding manifests or terminal output first.

**Technical scope:**

- Add private connector health summaries: configured, dry-run ready, send disabled/enabled, last synthetic check, redaction status, known blockers.
- Add copyable config checklists without secrets or endpoints.
- Link connector blockers to telemetry, prediction, validation, and monitoring issue categories.
- Keep commands as shell guidance only.

**Likely files:** `cmd/agency-config/operations_connectors.go`, `cmd/agency-config/operations_connector_workbench.go`, `internal/connectors/`, `operations_admin.js` if needed, tests.

**Validation:** `go test ./cmd/agency-config ./internal/connectors`, `make product-ui-smoke`, `make external-connection-check`, `make audit-product-language`.

**Risk review:** Do not render secrets, private destinations, private URLs, or raw payloads.

**Expected commit:** `Phase 151 -- Checkpoint 000001: add connector health review to private console`

### Phase 152 — Monitoring, Reliability, And Export V2

**Purpose:** Make self-hosted operations observable without promising SLA.

**User value:** Deployment owners can review health, freshness, validators, support summaries, and maintenance posture.

**Technical scope:**

- Improve redacted health digest and monitoring/export summaries.
- Add or refine local export formats for feed health, connector health, validator posture, telemetry freshness, and maintenance tasks.
- Keep destination sends disabled by default.
- Improve Operations Reliability and Maintenance pages around actionable owner tasks.

**Likely files:** `scripts/operations-reliability.sh`, `scripts/operations-notify.sh`, `examples/connectors/monitoring-export/`, `cmd/agency-config/operations_maintenance.go`, tests.

**Validation:** `make operations-reliability`, `make operations-notify`, `go test ./cmd/agency-config ./examples/connectors/...`, `make audit-final-claim-review`.

**Risk review:** No SLA/uptime, hosted-service, or paid-support claim.

**Expected commit:** `Phase 152 -- Checkpoint 000001: improve redacted monitoring and reliability exports`

### Phase 153 — Self-Hosted Install, Upgrade, Backup, And Restore UX V2

**Purpose:** Make small-host operation safer and easier.

**User value:** Deployment owners can install, upgrade, back up, restore, and roll back with fewer surprises.

**Technical scope:**

- Improve reference deployment docs/scripts around preflight, env validation, systemd/proxy checks, database migration status, backup/restore readiness, upgrade stop-points, and rollback guidance.
- Add or refine deployment doctor categories and product UI smoke checks for self-hosted settings.
- Keep all checks local/reference and no-contact.

**Likely files:** `docs/deployment/*`, `scripts/deployment-doctor.sh`, `scripts/product-ui-smoke.sh`, `scripts/oci-reference-check.sh`, `cmd/agency-config/operations_maintenance.go`, tests.

**Validation:** `make deployment-doctor`, `make oci-reference-check`, `make product-ui-smoke`, `docker compose -f deploy/docker-compose.yml config`, `make audit-final-claim-review`.

**Risk review:** Do not run destructive backup/restore/migration actions in tests.

**Expected commit:** `Phase 153 -- Checkpoint 000001: harden self-hosted install and recovery UX`

### Phase 154 — Multi-Agency, Roles, And Tenant-Safe Operations V3

**Purpose:** Strengthen isolation and role clarity without claiming production multi-tenant hosting.

**User value:** Operators avoid cross-agency leaks and understand who can do what.

**Technical scope:**

- Recheck agency query matching, role gates, audit metadata visibility, public path routing, debug JSON exposure, encoded slash/backslash agency rejection, and no-store behavior.
- Improve role-specific UI hints and access-denied guidance.
- Add tests for any identified gaps.

**Likely files:** `cmd/agency-config/`, `internal/auth/`, `internal/compliance/`, `scripts/audit-operations-route-inventory.sh`, tests.

**Validation:** `go test ./cmd/agency-config ./internal/auth`, `make audit-operations-route-inventory`, `make product-ui-smoke`, `make audit-final-claim-review`.

**Risk review:** No production multi-tenant hosting claim.

**Expected commit:** `Phase 154 -- Checkpoint 000001: strengthen agency isolation and role guidance`

### Phase 155 — Public Feed Discovery And External Sharing Prep V2

**Purpose:** Make `/public/feeds.json` and sharing preparation better without moving consumer statuses.

**User value:** Deployment owners can prepare feed metadata and know what is missing before authorized sharing.

**Technical scope:**

- Improve discovery metadata checks: license, contact, feed URLs, active schedule, realtime availability, stable base URL, HTTPS posture if configured, and publication environment.
- Add private sharing-prep guidance for Transitland/Mobility Database style metadata without submitting anything.
- Keep all consumer targets `prepared`; no portal automation or external contact.

**Likely files:** `internal/compliance/`, `cmd/agency-config/operations_readiness_v2.go`, `cmd/agency-config/operations_feed_health.go`, `docs/external-connection-readiness.md`, tests.

**Validation:** `go test ./internal/compliance ./cmd/agency-config`, `make validate-public-feeds` only if safe/configured, `make audit-final-claim-review`, `scripts/check-consumer-tracker.sh`.

**Risk review:** No consumer submission, acceptance, listing, display, ingestion, or final-root claim.

**Expected commit:** `Phase 155 -- Checkpoint 000001: improve public feed discovery and sharing preparation`

### Phase 156 — Security, Redaction, And Support Bundle V2

**Purpose:** Reduce accidental disclosure risk as product capabilities grow.

**User value:** Operators can produce useful support output without leaking secrets or private data.

**Technical scope:**

- Improve support bundle redaction for tokens, URLs, headers, payloads, DB URLs, private paths, endpoint names, and logs.
- Add tests with hostile/private-like strings.
- Recheck screenshots/tutorial capture path for secret exclusion.
- Ensure new connector/monitoring/deployment outputs remain private and redacted.

**Likely files:** `scripts/support-bundle.sh`, `scripts/capture-ui-tour.sh`, redaction docs, tests, fixtures.

**Validation:** `make support-bundle`, `make check`, `go test ./...` where relevant, `make audit-final-claim-review`.

**Risk review:** Do not commit generated bundles or private data.

**Expected commit:** `Phase 156 -- Checkpoint 000001: harden redaction and support bundle safety`

### Phase 157 — Nontechnical Staff Training And In-App Guidance V2

**Purpose:** Make the product teach itself without overwhelming users.

**User value:** Less technical staff can complete common workflows and know when to ask for help.

**Technical scope:**

- Improve Help, tutorial, docs, and public site based on actual product flows after phases 141–156.
- Add role-based quick paths for agency staff, administrator, deployment owner, and integrator.
- Replace remaining giant training tables with concise flows and expandable detail.
- Add tests/audits for primary help language.

**Likely files:** `cmd/agency-config/operations_help.go`, `site/*`, `docs/tutorials/*`, `wiki/*`, tests.

**Validation:** `make check-links`, `make audit-product-language`, `make product-ui-smoke`, `go test ./cmd/agency-config`.

**Risk review:** Do not turn help into a marketing page or phase ledger.

**Expected commit:** `Phase 157 -- Checkpoint 000001: improve nontechnical staff guidance`

### Phase 158 — API, Feed Contract, And Extension Governance V2

**Purpose:** Make extension points stable enough for open-source adoption.

**User value:** Civic technologists and integrators can build against known contracts.

**Technical scope:**

- Document and test stable contracts for `/v1/telemetry`, public feed paths, `/public/feeds.json`, admin JSON companion routes, connector manifests, adapter conformance fixtures, and prediction adapter DTOs.
- Add versioning/deprecation guidance where missing.
- Add contract tests for breaking-change detection where practical.

**Likely files:** `docs/api*`, `docs/connectors/*`, `internal/connectors/`, `cmd/adapter-conformance/`, `cmd/agency-config/*_test.go`, tests.

**Validation:** `go test ./...`, `make adapter-conformance`, `make gtfsrt-conformance`, `make external-connection-check`, `make check-links`.

**Risk review:** Do not promise indefinite API stability beyond documented release-candidate boundaries.

**Expected commit:** `Phase 158 -- Checkpoint 000001: clarify and test integration contracts`

### Phase 159 — Release Candidate Gate V3 And Stable Branch Readiness

**Purpose:** Turn the improved product into a release-candidate-ready state if validation supports it.

**User value:** Users can trust the install path and know what remains a release-candidate limitation.

**Technical scope:**

- Run/refresh release-candidate diagnostics, install-confidence checks, package audit, link checks, product UI smoke, connector checks, GTFS-RT conformance, and final claim audit.
- Update release status/draft notes only if source state supports it.
- Recheck stable branch filter if product files should be propagated.
- Do not publish a release/tag unless explicitly authorized; prepare the gate and report status.

**Likely files:** `docs/release-candidate-*.md`, `scripts/release-candidate-check.sh`, `scripts/install-confidence.sh`, `.github/workflows/*`, stable filter docs/scripts.

**Validation:** `git diff --check`, `go test ./...`, `make check`, `make test`, `make smoke`, `make check-links`, `make product-ui-smoke`, `make external-connection-check`, `make adapter-conformance`, `make gtfsrt-conformance`, `make audit-product-acceptance`, `make audit-final-claim-review`, `scripts/check-consumer-tracker.sh`.

**Risk review:** No release publication or release-ready claim unless the gate actually supports it and user authorizes publication.

**Expected commit:** `Phase 159 -- Checkpoint 000001: run release candidate gate for product-quality roadmap`

### Phase 160 — Better Software Roadmap Closeout And Next Track Decision

**Purpose:** Close this roadmap with evidence-backed status and choose the next best software track.

**User value:** The project exits with clearer software quality, real validation, and a defensible next direction.

**Technical scope:**

- Summarize phases 141–159 with what changed, what passed, what is still missing, and what claims remain unsupported.
- Confirm protected evidence paths and consumer tracker remain unchanged.
- Decide the next track based on repo evidence, not guesses. Candidate next tracks:
  - GTFS extensions and data-quality support: Flex, Pathways, Fares v2, GTFS-ride.
  - Deeper realtime correctness: observed-arrival evaluation harness, delay propagation, frequency service, cancellations, blocks.
  - Real connector runtime hardening with authorized deployment data.
  - Release publication, if the gate supports it and the user authorizes it.
  - Optional evidence path, only with explicit authorization and retained artifacts.
- Update current-status/latest handoff succinctly.

**Likely files:** `docs/current-status.md`, `docs/handoffs/latest.md`, optional closeout doc, no evidence paths.

**Validation:** full available validation set; always include `scripts/check-consumer-tracker.sh` and final claim audit.

**Risk review:** Do not convert roadmap closeout into a false claim. Separate software capability from external proof.

**Expected commit:** `Phase 160 -- Checkpoint 000001: close better software roadmap and recommend next track`

## Phase Execution Rules

For each phase:

1. Re-check current repo state and active branch.
2. Run `git status --short` and confirm protected evidence paths/consumer status are not already dirty.
3. Produce a concise phase plan.
4. Use a read-only scout when architecture, GTFS-RT correctness, security, UX, or release risk is meaningful.
5. Implement the smallest complete product change for the phase.
6. Run relevant tests, typecheck, lint/build, route audits, link checks, domain validation, and product-language/layout checks.
7. Review the diff before commit; optionally use a read-only diff-review scout.
8. Fix issues and rerun relevant validation.
9. Commit the phase with the expected commit message or a clear checkpoint variant.
10. Produce a concise phase summary with changed files, validations, and known risks.

Continue until Phase 160 is complete unless blocked by:

- missing credentials or secrets;
- destructive operations requiring user approval;
- ambiguous product requirements that would cause major rework;
- validation failures that cannot be resolved safely;
- repo state that makes continuation unsafe;
- unavailable browser/capture environment for screenshot-specific changes, in which case do not fake screenshots.

If blocked, stop, explain the blocker clearly, and provide the smallest next action required from the user.

## Git And Commit Policy

At the end of each completed phase:

- run `git status`;
- inspect the diff;
- ensure no secrets, credentials, generated junk, private data, or unrelated files are included;
- ensure protected evidence paths are unchanged;
- ensure consumer statuses remain unchanged;
- run relevant validation;
- create one commit for that phase.

Do not commit failing validation unless explicitly making a known-blocker checkpoint. If committing a blocker checkpoint, mark the blocker in the commit message and phase summary.

## Validation Policy

Discover the repo’s actual validation commands and use the relevant subset. Before final closeout, run the widest available set:

```bash
git diff --check
go test ./...
make check
make test
make smoke
make check-links
make product-ui-smoke
make audit-product-language
make audit-ui-layout
make audit-product-acceptance
make audit-final-claim-review
make audit-operations-route-inventory
make external-connection-check
make adapter-conformance
make test-connector-examples
make gtfsrt-conformance
scripts/check-consumer-tracker.sh
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config
```

If a command is unavailable, record the exact reason and run the closest safe substitute. Never report a skipped command as passed.

## Final Completion Criteria

The final Codex response must include:

- completed phases;
- commits created;
- files changed;
- validation performed;
- product capabilities added or improved;
- how the work improved GTFS, GTFS-RT, connectors, operator workflow, deployment, monitoring, security, or release readiness;
- confirmation that protected evidence paths and consumer statuses remained unchanged;
- confirmation that unsupported claims remain unsupported;
- known risks;
- follow-up recommendations.

Do not stop after only one phase. Continue until all planned phases are complete or until a real blocker is reached.
