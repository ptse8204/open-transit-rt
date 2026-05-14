# Phase 79 Handoff -- Agency Setup V3

## Phase

Phase 79 -- Agency Setup V3.

## Sub-Agents Used Or Simulated

Real sub-agents were used where available, with the intended model levels from
the Phase 75-90 track:

- Master Agent: GPT-5.5 x-high, main rollout.
- Context / Repo Truth Sub-Agent: GPT-5.5 x-high, Epicurus.
- Planning Sub-Agent: GPT-5.5 x-high, Russell.
- Security/Auth Sub-Agent: GPT-5.5 high, Kierkegaard.
- UI/UX Sub-Agent: GPT-5.5 high, Averroes.
- QA Sub-Agent: GPT-5.5 high, Singer.
- Claim-Boundary Sub-Agent: GPT-5.5 high, Galileo.
- Implementation Sub-Agent: GPT-5.5 high, simulated in the main rollout.
- Documentation / IA Sub-Agent: GPT-5.5 high, simulated in the main rollout.
- Data/Migration Sub-Agent: not used because Phase 79 added no migration or
  persisted model.

The claim-boundary review found one required wording edit around
`agency-approved open license`; Phase 79 replaced that with
operator-provided metadata language before closeout.

## Goal

Make first-run setup more usable for less technical staff while staying inside
the private Operations Console and preserving existing authorization,
CSRF, consumer-tracker, evidence, and claim boundaries.

## Changed Files

- `README.md`
- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_setup_wizard.go`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-79.md`
- `docs/phase-79-agency-setup-v3.md`
- `docs/roadmap-status.md`
- `docs/tutorials/agency-first-run.md`
- `wiki/operations-console-tour.md`
- `wiki/small-agency-quick-start.md`

## Routes Added/Changed

No new routes were added.

Changed private authenticated Operations Console surfaces:

- `GET /admin/operations/setup-wizard`
- `GET /admin/operations/setup-wizard.json`
- `GET /admin/operations/setup`
- `GET /admin/operations/gtfs-import`

`/admin/operations/setup` now hides mutation forms from non-admin roles while
the existing POST handlers remain admin-only.

## Commands Added/Changed

No commands were added or changed.

## Migrations

None.

## Validation Run

- `git status --short`
- `git diff --check`
- `make check`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`
- exact seven-target prepared-only consumer tracker check
- `git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum`
- focused setup/GTFS route tests:
  `go test ./cmd/agency-config -run 'Test(OperationsSetupRendersTruthfulMissingStates|SetupWizardRoutesPrivateScopedGETOnlyNoStore|SetupWizardJSONShapeFlagsAndStages|SetupWizardHTMLBoundariesNoFormsAndEscapes|OperationsSetupPublicationFormRequiresAdminAndDerivesAgencyID|OperationsSetupPublicationFormRejectsConflictingAgencyID|OperationsSetupValidationFormMapsFeedTypeServerSide|OperationsSetupValidationFormRejectsUnsafeBrowserFields|OperationsSetupCookiePostRequiresCSRF|OperationsConsoleFormsUseLabelsAndSubmitButtonsWithoutChangingContracts)'`
- `go test ./cmd/agency-config ./internal/auth ./internal/tenant`
- `node --check cmd/agency-config/operations_admin.js`
- `node --test cmd/agency-config/operations_admin_test.mjs`
- `make validate`
- `make test`
- `docker compose -f deploy/docker-compose.yml config`
- `RUN_LOCAL_APP=true make release-candidate-check`

The local release-candidate diagnostic wrote
`.cache/release-candidate-check/20260514T010924Z` and exited `0` with
`overall_status=not_checked`; the helper reported local app startup and five
public feeds as passed, while repository validation, Go tests, HTTP smoke, and
release package audit remain helper `not_checked` rows by design. `make
validate` and `make test` passed separately in this phase.

## Blocked Checks

- Release package creation/audit was not run because Phase 79 is a product UI
  phase and release packaging/publishing remains unauthorized.
- `make release-candidate-check` intentionally leaves release package audit and
  some bounded helper rows as `not_checked`; this was recorded as diagnostic
  status, not a release-readiness claim.

## Known Blockers

- Phase 72 remains `needs_review`, not release-ready.
- No release tag, package, published image, hosted SaaS endpoint, final-root
  proof, consumer submission, real agency pilot, real vendor/device proof, or
  production-grade ETA proof exists.
- Phase 80 still needs the broader GTFS Workbench buildout.

## Protected Path Status

Clean. No files under protected retained-evidence or consumer packet/status
paths were modified or generated.

## Consumer Tracker Status

All seven targets remain exactly `prepared`:

- Google Maps
- Apple Maps
- Transit App
- Bing Maps
- Moovit
- Mobility Database
- transit.land

## Claim-Boundary Status

Clean after the required wording edit. Phase 79 did not claim CAL-ITP/Caltrans
compliance, agency adoption or approval, consumer submission/review/acceptance,
final-root readiness, hosted SaaS availability, paid support, SLA/uptime,
production readiness, vendor compatibility, hardware certification, public
launch, production-grade ETA quality, or real-world ETA accuracy.

## Security/Auth Status

Preserved. Setup and setup-wizard routes remain private. Setup JSON is
GET-only and no-store. Setup mutations remain existing admin-only POST actions
with CSRF and agency derivation from the authenticated principal. The
read-only/operator/editor setup page no longer renders mutation forms.

## Accessibility Status

Preserved through the existing server-rendered Operations Console shell,
landmarks, focus styles, responsive tables, card grids, status chips, and
no-JS fallback. Phase 79 did not add a frontend dependency or SPA route.

## Docs/Site/Wiki Alignment

Updated source-of-truth status docs plus README/tutorial/wiki references so the
private setup path is described as Agency Setup with progress, diagnostics,
role visibility, technical-helper escalation, and admin-only mutations.

## Commit List

- `e26a5aa` -- Phase 79 -- Checkpoint 000001: add agency setup v3 plan
- `fb8af0c` -- Phase 79 -- Checkpoint 000002: implement agency profile and metadata review
- `9f452d6` -- Phase 79 -- Checkpoint 000003: improve browser GTFS source import review
- `25f7f27` -- Phase 79 -- Checkpoint 000004: add setup progress and blocker next actions
- Phase 79 -- Checkpoint 000005: close agency setup v3 review

## Master Review

Approved after implementation, QA, documentation/IA, security/auth, UI/UX, and
claim-boundary review. The required claim-boundary wording edit was applied and
rechecked.

## Required Edits

None remaining.

## Decision

Phase 79 is closed for the authorized Agency Setup V3 scope.

## Next Phase

Phase 80 -- GTFS Workbench.
