# Phase 84 Handoff -- Prediction And ETA Lab

## Closeout Report

Phase: Phase 84 -- Prediction And ETA Lab

Sub-agents used or simulated, including intended model level:

- Master Agent, GPT-5.5 x-high: approved the plan, implementation checkpoints,
  validation results, and this closeout.
- Context / Repo Truth Sub-Agent, GPT-5.5 x-high: Locke returned repo-truth
  constraints for Phase 84, including the existing prediction adapter,
  Trip Updates diagnostics, `.cache/realtime-quality-backtest` output contract,
  protected paths, and hard-stop conditions.
- Planning Sub-Agent, GPT-5.5 x-high: Volta approved the route plan,
  checkpoint sequence, no-migration default, local aggregate backtest browsing,
  and stop conditions.
- Implementation Sub-Agent, GPT-5.5 high: simulated locally in the main
  rollout.
- QA Sub-Agent, GPT-5.5 high: Godel returned route, JSON, no-write,
  prepared-tracker, forbidden-claim, no-raw-output, and backtest safety
  requirements.
- UI/UX Sub-Agent, GPT-5.5 high: Hubble returned the Realtime nav placement,
  page IA, table structure, no-JS expectation, and operator copy requirements.
- Documentation / IA Sub-Agent, GPT-5.5 high: Avicenna returned tutorial,
  docs/status, and closeout alignment requirements.
- Claim-Boundary Sub-Agent, GPT-5.5 high: Fermat returned allowed wording and
  blocked claim language for deterministic diagnostics, shadow review,
  backtests, withheld-output explanations, and proof gates.
- Security/Auth Sub-Agent, GPT-5.5 high: simulated locally because the real
  agent slot was not available; review required private GET-only routes,
  no browser command execution, no raw/private output exposure, no external
  predictor calls, no uploads, no evidence writes, and no consumer tracker
  changes.
- Data/Migration Sub-Agent: not used; no migration or persisted model was
  added.

Goal:

Add a private browser-first Prediction and ETA Lab for explaining why Trip
Updates and ETA-like output is emitted, withheld, missing, shadowed, or failed
closed while keeping Vehicle Positions independent and avoiding production ETA,
consumer, compliance, vendor, hardware, hosted-service, SLA, release, or
evidence claims.

Changed files:

- `cmd/agency-config/main.go`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_prediction_lab.go`
- `docs/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-84.md`
- `docs/open-transit-rt-master-planner-remaining-work.md`
- `docs/phase-84-prediction-and-eta-lab.md`
- `docs/roadmap-status.md`
- `docs/tutorials/README.md`
- `docs/tutorials/prediction-eta-lab.md`
- `internal/compliance/model.go`
- `internal/compliance/postgres.go`
- `internal/realtimequality/browser.go`
- `internal/realtimequality/browser_test.go`

Routes added/changed:

- Added private `GET /admin/operations/prediction-lab`.
- Added private `GET /admin/operations/prediction-lab.json`.
- Added the Prediction & ETA Lab navigation item under the private Realtime
  group after Realtime Center.
- No public route was added.

Commands added/changed:

- No command was added or changed.
- The Lab displays fixed operator-shell guidance for existing local checks:
  `go test ./internal/prediction -run Deterministic`,
  `make realtime-quality`, and `make realtime-quality-backtest`.
- Browser requests execute none of those commands.

Migrations:

- None.

Validation run:

- `git status --short`
- `git diff --check`
- `go test ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'`
- `go test -count=1 ./cmd/agency-config -run 'PredictionLab|Realtime|OperationsNavigation|RouteTitles'`
- `go test ./cmd/realtime-quality-backtest ./internal/realtimequality ./internal/prediction ./internal/feed/tripupdates`
- `go test ./internal/realtimequality ./internal/prediction ./internal/feed/tripupdates`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make validate`
- `make test`
- `make realtime-quality`
- `make realtime-quality-backtest`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `RUN_LOCAL_APP=true make release-candidate-check`
  (`overall_status=needs_review`, diagnostic only, not release proof)
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`

Blocked checks:

- None.

Known blockers:

- Phase 72 remains `needs_review`, not release-ready.
- No release tag, release package, published image, final-root proof, consumer
  submission, compliance packet, real agency pilot, real vendor/device proof,
  SLA/uptime proof, real-world ETA accuracy proof, or production-grade ETA
  proof exists.
- Optional evidence tracks remain authorization-gated.

Protected path status:

- Protected evidence paths were not modified.
- `docs/evidence/consumer-submissions/status.json` was not modified.
- `docs/evidence/consumer-submissions/current/**`,
  `docs/evidence/consumer-submissions/artifacts/**`,
  `docs/evidence/consumer-submissions/packets/**`, and
  `docs/evidence/captured/**` were not modified.

Consumer tracker status:

- All seven targets remain exactly `prepared`: Google Maps, Apple Maps,
  Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.

Claim-boundary status:

- The Lab is described as private prediction diagnostics, deterministic
  fallback review, withheld-output explanation, external predictor
  shadow/fail-closed review, conservative handling guidance, and local
  aggregate backtest browsing only.
- It does not claim CAL-ITP/Caltrans compliance, agency adoption/approval,
  consumer submission/review/acceptance/ingestion/listing/display,
  final-root readiness, hosted SaaS, production readiness, vendor
  compatibility, hardware certification, SLA/uptime, public launch,
  real-world ETA accuracy, or production-grade ETA quality.
- Future proof gates are listed only as separate authorization-gated work.

Security/auth status:

- Lab routes are private, authenticated, role-checked for read-only operations
  roles, agency-query checked, no-store, and GET-only.
- No Lab form, POST route, browser predictor run, command execution, sidecar
  start, external URL test, credential capture, upload, raw observed-row
  persistence, evidence write, public API, release artifact, or consumer
  tracker mutation was added.
- Backtest browsing reads only exact aggregate `.cache/realtime-quality-backtest`
  output shapes, rejects unsafe roots/symlinks/unexpected files/schema/claim
  flags, and returns redacted aggregate metrics only.

Accessibility status:

- The route uses the existing server-rendered Operations Console shell,
  landmarks, headings, tables, status text, and responsive layout.
- The Lab remains usable with no JavaScript because all core content is
  server-rendered.

Docs/site/wiki alignment:

- Source-of-truth status docs now mark Phase 84 complete and Phase 85 next.
- `docs/tutorials/prediction-eta-lab.md`, `docs/tutorials/README.md`, and
  `docs/README.md` explain the private Lab workflow and claim boundaries.
- GitHub Pages content was not changed in this phase.

Commit list:

- `11a12f7 Phase 84 -- Checkpoint 000001: add prediction and ETA lab plan`
- `21b343d Phase 84 -- Checkpoint 000002: add deterministic predictor diagnostics view`
- `7c88779 Phase 84 -- Checkpoint 000003: add external predictor shadow review UI`
- `9cdaad1 Phase 84 -- Checkpoint 000004: add backtesting result browser`
- `1f13f17 Phase 84 -- Checkpoint 000005: add ETA quality caveats and withheld explanations`
- Phase 84 closeout checkpoint: `Phase 84 -- Checkpoint 000006: close prediction and ETA lab review`

Master review:

- Approved. The phase stayed inside the authorized private, read-only,
  browser-first diagnostics scope and preserved all claim boundaries.

Required edits:

- None at closeout.

Decision:

- Close Phase 84 and continue to Phase 85 -- Operations And Maintenance Center
  V2.

Next phase:

- Phase 85 -- Operations And Maintenance Center V2.
