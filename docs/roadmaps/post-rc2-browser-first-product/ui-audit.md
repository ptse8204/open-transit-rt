# Phase 01 UI Audit

This audit records the current Operations Console shape at the start of the
post-`v0.1.0-rc.2` browser-first roadmap.

It is planning material only. It does not create retained evidence, contact
external parties, move consumer statuses, publish a release, or prove
compliance, production readiness, consumer acceptance, hosted service
availability, vendor compatibility, SLA coverage, production AVL reliability,
production-grade ETA quality, or real-world ETA accuracy.

## Sources Reviewed

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/release-status-v0.1.0-rc.2.md`
- `docs/release-download-replay-v0.1.0-rc.2.md`
- `docs/release-notes-v0.1.0-rc.2.md`
- `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`
- `docs/roadmaps/agency-first-connector-platform/README.md`
- `docs/evidence/redaction-policy.md`
- `.github/workflows/test.yml`
- `cmd/agency-config/`
- `internal/admincontrol/`
- `internal/connectors/`
- `internal/gtfsrtconformance/`
- `docs/`
- `wiki/`

`site/` is not present in the main worktree. The public website currently lives
on the `gh-pages` branch according to existing docs and branch state.

## Web Design Engineer Skill

The Web Design Engineer skill was loaded before this audit. It influenced this
phase by keeping the plan focused on user task flow, route hierarchy, visual
density, readable labels, responsive behavior, and product polish while
preserving the existing Go server-rendered approach and no-JS fallback.

## Private UI Route Inventory

The current Operations Console route registry defines these private browser
surfaces.

| Area | Routes | Notes |
| --- | --- | --- |
| Start Here / Dashboard | `/admin/operations`, `/admin/operations.json`, `/admin/operations/launchpad`, `/admin/operations/launchpad.json` | Main cockpit, first-run tasks, role entry cards, launchpad diagnostics. |
| Setup | `/admin/operations/setup-wizard`, `/admin/operations/setup-wizard.json`, `/admin/operations/setup` | Guided setup plus advanced metadata forms. |
| GTFS Workbench | `/admin/operations/gtfs-workbench`, `/admin/operations/gtfs-workbench.json`, `/admin/operations/gtfs-import`, `/admin/operations/gtfs-quality` | Browser import exists; quality triage and workbench are separate. |
| Feed Links / Health | `/admin/operations/feeds`, `/admin/operations/feed-health`, `/admin/operations/feed-health.json` | Feed URLs, local health, validation context, and realtime usefulness. |
| Validation | `/admin/operations/validation-health`, `/admin/operations/validation-health.json`, `/admin/operations/validation-health/refresh.json`, `/admin/operations/validation-center`, `/admin/operations/validation-center.json` | Review plus allowlisted server-owned validation actions. |
| Realtime | `/admin/operations/realtime`, `/admin/operations/realtime.json`, `/admin/operations/telemetry`, `/admin/operations/prediction-lab`, `/admin/operations/prediction-lab.json` | Fleet freshness, Vehicle Positions, Trip Updates, withheld reasons, ETA lab. |
| Devices / AVL | `/admin/operations/devices`, `/admin/operations/telemetry-simulator`, `/admin/operations/telemetry-simulator.json` | Device credential actions and simulator guide; sending still points to shell commands. |
| Connectors | `/admin/operations/connectors`, `/admin/operations/connectors.json`, `/admin/operations/connectors/workbench`, `/admin/operations/connectors/workbench.json`, `/admin/operations/connectors/tests`, `/admin/operations/connectors/tests.json` | Synthetic/local connector catalog, workbench, and generated fixed command guidance. |
| Readiness | `/admin/operations/readiness`, `/admin/operations/readiness.json`, `/admin/operations/checklist`, `/admin/operations/checklist.json` | CAL-ITP-style readiness and private operator checklist. |
| Maintenance / Troubleshooting | `/admin/operations/reliability`, `/admin/operations/reliability.json`, `/admin/operations/maintenance`, `/admin/operations/maintenance.json`, `/admin/operations/access`, `/admin/operations/access.json`, `/admin/operations/audit`, `/admin/operations/audit.json` | Maintenance, reliability, role guidance, metadata-only audit. |
| Help / Tutorials | `/admin/operations/help`, `/admin/operations/help/`, `/admin/operations/help.json`, `/admin/operations/help.json/` | Role tours, glossary, first-week checklist, recovery guidance. |
| Consumer / Evidence Boundary | `/admin/operations/consumers`, `/admin/operations/evidence` | Prepared-only consumer tracker and evidence/runbook links. These are claim-boundary surfaces, not consumer action surfaces. |
| Assets | `/admin/operations/assets/operations.js` | Private buildless JS enhancement route. |
| External Admin Surfaces | `/admin/gtfs-studio`, `/admin/alerts/console` | Linked from the Operations Console but not fully integrated into the shell. |
| Legacy Private Admin | `/admin/publication/bootstrap`, `/admin/compliance/scorecard`, `/admin/consumer-ingestion`, `/admin/validation/run`, `/admin/devices/rebind` | Older private routes still exist outside the Operations Console shell. |

Most Operations Console routes are authenticated and set `Cache-Control:
no-store`. Phase 02 should keep route-level no-store/private boundaries covered
by tests as navigation and shell changes continue.

## Browser-First Progress Already Present

The current UI already supports many browser-first workflows:

- Start Here cockpit with setup progress, role cards, primary actions, and five
  feed URL review.
- Agency setup wizard and advanced setup details.
- Browser GTFS upload/import by file or safe URL for admin users.
- GTFS Workbench review and JSON companion.
- Feed Health, Validation Center, Validator Health, Readiness, Realtime,
  Prediction Lab, Connectors, Connector Workbench, Maintenance, Access, Audit,
  and Help pages.
- Private JSON companions for many major routes.
- A private buildless JavaScript layer for copy/filter/refresh enhancements,
  with no-JS HTML fallback.
- Strong claim-boundary language that prevents accidental compliance,
  production, vendor, consumer, hosted-service, SLA, or ETA-quality claims.

## Workflows Still Requiring Command-Line Help After Startup

These workflows still need a technical helper or operator shell after the app is
running. Later phases should decide which can be moved safely into bounded
server-owned browser actions and which must remain out of browser scope.

| Workflow | Current browser state | Why shell is still needed |
| --- | --- | --- |
| Local app start/stop/logs/reset | Browser starts after `make agency-app-up`. | Starting/stopping services and destructive reset are intentionally outside the app. |
| Validator installation | UI explains validator health. | `make validators-install` and environment setup remain shell tasks. |
| Heavy validation / off-host checks | Admins can request allowlisted validation where configured. | `make validate`, `make validate-public-feeds`, Docker-backed validators, and off-host fetch checks may need local tooling and network context. |
| Telemetry simulator send | Browser shows scenarios and commands. | Sending synthetic telemetry still uses `make telemetry-simulator`; device tokens must stay out of browser history/storage. |
| Support bundle | Maintenance explains when to generate one. | `make support-bundle` remains shell-only and writes private `.cache` diagnostics. |
| Operator smoke / deployment doctor | Maintenance and docs link the workflows. | `make operator-smoke`, `make deployment-doctor`, `make oci-reference-check`, and related diagnostics are shell-run. |
| Backup, restore, migrations, upgrade, rollback | Maintenance provides checklists. | Destructive or deployment-changing actions remain outside the browser. |
| Connector conformance and examples | Connector pages show fixed commands. | `make external-connection-check`, `make adapter-conformance`, and `make test-connector-examples` are shell-run. |
| Connector sidecar installation/configuration | Browser explains safe shapes. | Sidecar setup, credentials, manifests, and deployment-owned runtime wiring are not launched by the browser. |
| Large/scripted GTFS imports and rollback execution | Browser upload and safe URL import exist. | Very large/scripted imports and rollback/staged promotion workflows still point to technical-helper paths. |
| Release candidate and package checks | Validation/maintenance docs mention them. | Release gates, package audits, and publishing remain shell/GitHub workflows. |
| Public source-of-truth or consumer workflows | Readiness/consumer pages explain boundaries. | Final-root evidence, consumer submissions, and target communication are authorization-gated and must not be browser-automated in this roadmap. |

## Language And UX Issues

These issues should guide Phase 02 through Phase 10.

| Issue | Current examples | Recommended direction |
| --- | --- | --- |
| Internal product wording appears in human docs and UI. | `The product path is:`, `Private operations control plane`, `claim flags`, `prepared-only`, `bounded`, `source-of-truth`. | Use plain agency language: "Start here", "Review feed health", "These checks help you prepare, but do not prove outside approval." Keep internal terms in agent/internal docs. |
| Some labels feel diagnostic rather than product-oriented. | `Private Launchpad`, `Private Operator Checklist`, `Evidence`, `Consumers`, `Readiness Checklist V2`, `Diagnostic only`. | Rename or group around user tasks: Start Here, Setup, Schedule, Feed Health, Validation, Realtime, Devices, Connectors, Readiness, Maintenance, Help, Troubleshooting. |
| The dashboard is powerful but dense. | Start Here renders role cards, first-run panel, setup progress, primary actions, readiness, dashboard sections, feed URL table, and Trip Updates summary. | Preserve density for operators, but add clearer page zones, a shorter top task list, and stronger hierarchy for first-time users. |
| Claim-flag tables are too exposed for nontechnical users. | Help, launchpad, first-run, realtime, simulator, maintenance, and consumer pages render raw boolean flag names. | Keep flags in JSON/internal details, but replace human-visible tables with a short "What this does not prove" summary unless an advanced details disclosure is necessary. |
| Several pages mix daily operations with release/evidence concepts. | Dashboard and validation surfaces include scorecard, consumers, evidence, release-state wording, prepared tracker, and diagnostics. | Separate normal agency work from maintainer/release/evidence review. Use role-based entry points and hide evidence-only topics from normal first-run flows. |
| Some current signals expose implementation-shaped text. | `overall=...; tooling=...`, boolean fragments, validator IDs, cache/status terms. | Convert to sentence summaries and keep raw fields in JSON or advanced details. |
| Major external surfaces are not fully shell-integrated. | GTFS Studio and Alerts Console are marked as separate admin surfaces. | Phase 02 and Phase 05 should either integrate their navigation/headers or clearly explain why they open separate private tools. |
| Site source location is split. | Main branch has no `site/`; docs say GitHub Pages lives on `gh-pages`. | Phase 09 should use a clear website source strategy before editing and document whether site files live in main, `gh-pages`, or a worktree. |

## Diagnostic-Only Or Overlong Screens

Priority screens to simplify or restructure:

- `/admin/operations`: too much content for a first screen; needs a clearer
  product dashboard hierarchy.
- `/admin/operations/help`: valuable but very long; should become task-based
  help and tutorial entry points.
- `/admin/operations/connectors/workbench`: comprehensive but dense; Phase 06
  should convert it into a connector catalog/workbench flow.
- `/admin/operations/maintenance`: mixes support bundle, small-host, backup,
  upgrade, notification, restore, and cadence concerns; Phase 03 and Phase 14
  should clarify what normal users can do from the browser.
- `/admin/operations/validation-center`: good safety language, but should
  emphasize issue triage and next actions before raw diagnostics.
- `/admin/operations/consumers` and `/admin/operations/evidence`: useful for
  claim boundaries, but these should not dominate normal agency workflows.

## Browser-First Information Architecture

Use this target IA for later phases.

| Target section | Primary routes | Purpose |
| --- | --- | --- |
| Start Here / Dashboard | `/admin/operations` | One clear first screen: setup state, top alerts, next best actions, and links into major workflows. |
| Setup | `/admin/operations/setup-wizard`, `/admin/operations/setup`, `/admin/operations/access` | Agency metadata, public base/feed URL metadata, roles, environment state, and technical-helper handoff. |
| GTFS Workbench | `/admin/operations/gtfs-workbench`, `/admin/operations/gtfs-import`, `/admin/gtfs-studio`, `/admin/operations/gtfs-quality` | Import, review, draft/edit, quality triage, version comparison, publish review, rollback guidance. |
| Feed Health | `/admin/operations/feed-health`, `/admin/operations/feeds` | Five feed URLs, local fetch status, freshness, validation context, and what needs attention. |
| Validation | `/admin/operations/validation-center`, `/admin/operations/validation-health`, `/admin/operations/gtfs-quality` | Run/review allowlisted validation, triage issues, and explain what validation does not prove. |
| Realtime | `/admin/operations/realtime`, `/admin/operations/telemetry` | Vehicle Positions usefulness, stale/unmatched/suppressed vehicles, Trip Updates withheld/fallback state, and alerts links. |
| Devices / AVL | `/admin/operations/devices`, `/admin/operations/telemetry-simulator` | Device/token setup, vehicle binding, latest telemetry, simulator guidance, AVL connector handoff. |
| Connectors | `/admin/operations/connectors`, `/admin/operations/connectors/workbench`, `/admin/operations/connectors/tests` | Connector catalog, copy/adapt examples, conformance checks, redaction and fail-closed behavior. |
| Alerts | `/admin/alerts/console`, `/admin/operations/realtime`, `/admin/operations/feed-health` | Alert lifecycle, stale/missing links, disruption review, Alerts feed health. |
| Prediction / ETA Lab | `/admin/operations/prediction-lab`, `/admin/operations/realtime` | Prediction adapter status, withheld reasons, shadow/fail-closed behavior, no ETA-quality overclaim. |
| Maintenance | `/admin/operations/maintenance`, `/admin/operations/reliability` | Routine checks, backup/restore status, validator cadence, support boundaries, update readiness. |
| Help / Tutorials | `/admin/operations/help`, `docs/tutorials/*`, `wiki/*` | Role-based help, video/tutorial scripts, nontechnical glossary, troubleshooting. |
| Readiness | `/admin/operations/readiness`, `/admin/operations/checklist` | CAL-ITP-style readiness review without compliance claims. |
| Support Bundle / Troubleshooting | `/admin/operations/maintenance`, `/admin/operations/reliability`, `/admin/operations/audit`, `/admin/operations/access` | Safe support summary guidance, audit metadata, role/access troubleshooting, redaction steps. |

## Later Phase File And Test Impact

Likely code files:

- `cmd/agency-config/operations_route_registry.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_design_system.go`
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_admin.js`
- `cmd/agency-config/operations_cockpit.go`
- `cmd/agency-config/operations_first_run.go`
- `cmd/agency-config/operations_setup_wizard.go`
- `cmd/agency-config/operations_gtfs_workbench.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_validation_center.go`
- `cmd/agency-config/operations_realtime.go`
- `cmd/agency-config/operations_devices.go`
- `cmd/agency-config/operations_connectors.go`
- `cmd/agency-config/operations_connector_workbench.go`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations_help.go`
- `internal/admincontrol/`
- `internal/connectors/`
- `internal/gtfsrtconformance/`

Likely test files and checks:

- `cmd/agency-config/main_test.go`
- `cmd/agency-config/operations_admin_test.mjs`
- `cmd/agency-config/*_script_test.go`
- `internal/admincontrol/model_test.go`
- `internal/connectors/*_test.go`
- `internal/gtfsrtconformance/*_test.go`
- `scripts/audit-product-acceptance.sh`
- `scripts/audit-final-claim-review.sh`
- `scripts/check-consumer-tracker.sh`
- `.github/workflows/test.yml`

Likely docs/site/wiki files:

- `README.md`
- `docs/README.md`
- `wiki/README.md`
- `wiki/browser-first-setup.md`
- `wiki/operations-console-tour.md`
- `wiki/small-agency-quick-start.md`
- `docs/tutorials/no-cli-agency-first-run.md`
- `docs/tutorials/small-agency-maintenance-guide.md`
- `docs/tutorials/operator-smoke-and-support-bundle.md`
- `docs/tutorials/external-adapter-conformance.md`
- `docs/dependencies.md`
- `docs/decisions.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `gh-pages` site files or a future main-branch `site/` directory, depending on
  the Phase 09 source strategy.

## Phase 01 Validation Commands

Run before accepting Phase 01:

```bash
git diff --check
make check
scripts/check-consumer-tracker.sh
make audit-operations-route-inventory
make audit-final-claim-review
```

These commands are read-only or local repository checks for this phase. They do
not create retained evidence or move consumer status.
