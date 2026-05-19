# Open Transit RT — Phase 133–140 Product UI, Language, Tutorial, and Operator Flow Reset

## Goal

Use Codex goal-driven work to complete a real product reset, not a docs-only or artifact-only pass.

**Goal statement:** Make Open Transit RT feel like a usable self-hosted transit operations product for small agencies, civic technologists, and operator staff. Replace AI-agent/maintainer/audit-style wording, reduce card clutter, redesign the Operations Console and GitHub Pages layouts, make the tutorial real and self-hosted-aware, and add regression checks so these problems do not return.

**Completion condition:** The public GitHub Pages site, README/docs/wiki entry points, and private Operations Console have been changed in actual source code and committed phase-by-phase. The live/publishable `site/` source and `gh-pages` branch no longer show the flagged wording or card-heavy layout patterns. The Operations Console primary path presents clear action lists and configuration pages without forcing users through a huge header/nav/card wall before settings. The UI tour/tutorial uses actual app screenshots or capture-backed assets and covers both local evaluation and self-hosted/reference operation. Tests and audits enforce the new language/layout rules.

**Non-goals:** Do not collect evidence, contact outside parties, move consumer statuses, claim compliance/adoption/production readiness, or create more phase-ledger documentation as the main output. This is product implementation work.

## Repository Context

Repository: `https://github.com/ptse8204/open-transit-rt`

Ground this work in the actual current code and deployed site source, not only handoff markdown. Read at minimum:

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-status.md`
- `docs/roadmaps/**/README.md` only as needed for source-of-truth context
- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_design_system.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_route_registry.go`
- `cmd/agency-config/operations_admin.js`
- `cmd/agency-config/*_test.go`
- `site/index.html`
- `site/try-locally.html`
- `site/ui-tour.html`
- `site/check-feeds.html`
- `site/connect-vehicles.html`
- `site/connector-support.html`
- `site/readiness.html`
- `site/video.html`
- `site/assets/site.css`
- current `gh-pages` branch files for the same public pages
- `docs/deployment/oci-reference-deployment.md`
- `docs/tutorials/no-cli-agency-first-run.md`
- `docs/tutorials/video-recording-guide.md`
- `wiki/README.md`
- `Makefile`
- relevant scripts under `scripts/`

Current confirmed issues to fix in actual code/site:

1. GitHub Pages still uses user-hostile or insider wording. Examples that must be removed or rewritten:
   - `Common Next Actions`
   - `Follow the same path staff use in the console.`
   - `What stays technical?`
   - repeated “technical helper” phrasing where “administrator”, “deployment owner”, or “operator” would be clearer
   - repeated “does not prove” blocks on every card instead of one concise “Limits” section
   - raw internal claim-flag vocabulary on primary HTML views
2. `site/ui-tour.html` is a screenshot gallery, not a real tutorial. It is too local-first and does not explain the self-hosted/reference deployment path.
3. The Operations Console still has too many cards, status tiles, diagnostics tables, details sections, and caveat panels. The header/nav plus context panels consume too much vertical space before actual settings/actions.
4. The public site and console overuse card grids. Replace them with tighter task lists, progressive disclosure, split-pane layouts, accordions, tables where useful, and action-first pages.
5. Previous UX reviews tolerated a “dense, utilitarian, diagnostic, claim-bounded” interface. That is not enough. This phase must use a product UX standard for prospective users and operator staff.
6. This must improve the actual product routes (`/admin/operations*`) and the actual publishable site (`site/` and `gh-pages`), not just local-demo text or markdown handoffs.

## Quality Priority

Quality is more important than speed. Prefer correctness, maintainability, accessible UI, clear operator flow, concrete validation, and durable product value over quick wording patches.

Avoid “artifact slop.” Do not satisfy this request by adding another roadmap document while leaving the app, site, and tutorial essentially unchanged.

## Master Agent Operating Mode

Use a high-effort master Codex agent. The master owns planning, implementation, validation, review, final decisions, and commits.

Do not split planning, execution, and review into separate full-context agents by default. Use bounded read-only scouts only when they materially improve quality.

Allowed scouts:

- `plan-risk-scout`: read-only critique of plan risks, missing files, route gaps, and product UX blind spots.
- `writing-ux-scout`: read-only critique of public/user-facing wording, tone, and overexplaining.
- `diff-review-scout`: read-only review of final diff for regressions, accessibility, layout, tests, and claim-boundary mistakes.
- `domain-scout`: read-only GTFS/GTFS-RT/operator workflow critique when product changes touch feed, validation, realtime, connector, or readiness semantics.

Scout rules:

- At most 1–2 scouts per phase unless there is a concrete reason.
- Scouts may inspect files, diffs, tests, and command output.
- Scouts must not edit files, format code, install dependencies, commit, change configuration, spawn subagents, or do broad repo exploration without a narrow reason.
- Every scout output must contain only: Scope inspected; Files/symbols reviewed; Findings; Evidence; Recommended action; Confidence level.
- The master synthesizes scout feedback and makes the final decision.

## Skills Requirement

Before implementation, search for and read any available local skills or guidance for:

- Web design / UX design / web-design-engineer
- Product copy / plain-language documentation / user documentation
- Accessibility / HTML forms / progressive disclosure
- Go server-rendered UI tests, if present

If a skill exists, use it. If no skill is available, do not fake it; apply the explicit product language and layout rules below.

The prior UX tolerance for “dense, utilitarian, diagnostic” pages is no longer acceptable for primary operator pages. Diagnostic depth can remain, but it must live behind progressive disclosure and support pages, not dominate the first view.

## Product Language Rules

Apply these rules to README, docs entry points, wiki entry points, GitHub Pages, and primary Operations Console HTML.

### Voice

Write for prospective users, small-agency staff, operator staff, civic technologists, and people seeking help. Do not write like a phase ledger, maintainer status report, AI-agent handoff, compliance audit log, or generated site readout.

Use direct plain language:

- “Import a schedule” instead of “GTFS import or GTFS Studio path”.
- “Check feed URLs” instead of “Configured public route URLs”.
- “Fix feed issues” instead of “Common Next Actions”.
- “When an administrator is needed” instead of “What stays technical?”.
- “Tour the operator workflow” instead of “Follow the same path staff use in the console.”
- “Limits” or “Before you share this” instead of repeating “Does not prove” on every card.

### Banned or strongly discouraged primary-page phrases

Remove these from public pages and primary Operations Console HTML, except in tests that verify absence or in an explicitly collapsed advanced/admin-only section:

- `Common Next Action`
- `Common Next Actions`
- `Follow the same path staff use in the console`
- `What stays technical?`
- `site readout`
- `phase closed`
- `checkpoint`
- `AI agent`
- `Codex`
- `claim flags`
- raw flag names such as `external_evidence_created`, `consumer_statuses_changed`, `production_grade_eta_claimed`, `hosted_saas_claimed`, `dynamic_backend_plugin_loading_enabled`

Avoid repeating:

- “does not prove” in every card
- “technical helper” in every section
- long lists of unsupported claims in the middle of workflows
- raw route lists on the first screen
- status-source/evidence-source language when a user needs the next step first

### Boundaries still required

Keep truthfulness. Do not add unsupported claims. But make boundaries concise and placed well:

- Public pages: one clear “Limits” section per page is enough unless a specific dangerous action needs an inline warning.
- Operations Console: primary pages may have one compact limits/permissions note; detailed claim boundaries and all-false flags should be JSON-only or under a collapsed “Advanced safety details” section, not visible by default.
- Evidence/legal/compliance docs may remain more formal, but still remove needless phase-ledger wording.

## Layout Rules

### Public GitHub Pages

- Reduce card grids. Do not use six-card grids where a simple ordered task flow, comparison table, or accordion is clearer.
- Navigation must include local evaluation and self-hosted/reference deployment, not only “Try locally.” Consider `Try & deploy`, `Tour`, `Feeds`, `Connectors`, `Readiness`, `Docs`, `GitHub`.
- `check-feeds.html` must become a practical troubleshooting page: five feeds, what good looks like, what to fix, and where to go next.
- `ui-tour.html` must become a real guided tutorial with step sections, not just a screenshot gallery.
- `video.html` must not imply a PPT/slideshow is a product tutorial. If video remains, it must be paired with capture-backed tutorial steps and clear self-hosted context.
- Every public page should answer: What is this? Who is it for? What can I do next? When do I need an administrator/deployment owner? What are the limits?

### Operations Console

- The first visible screen at 1366×768 must show the main action queue and current product status without forcing a scroll past a giant header/nav wall.
- Replace grid-of-cards navigation with a compact persistent sidebar, segmented top rail, or collapsed navigation menu.
- Use collapsible groups for secondary pages and diagnostic detail.
- Keep forms and actionable controls near the top of relevant pages.
- Prefer a straight-line task model:
  1. Set up agency
  2. Import schedule
  3. Check feed URLs
  4. Connect vehicle data
  5. Review realtime output
  6. Fix issues
  7. Maintain deployment
  8. Prepare external sharing only if authorized
- Do not show more than three summary cards above the fold on any core page.
- Diagnostic tables are acceptable only after the main action and short explanation.
- Raw claim flags must not appear in primary HTML. Keep them in JSON or collapsed advanced sections only.
- “Does not prove” repetition must be collapsed or consolidated.
- Ensure keyboard navigation, landmarks, skip links, focus states, and mobile behavior remain accessible.

## Roadmap

Complete all phases below unless a real blocker is reached. Do not stop after one phase.

### Phase 133 — Product Language And UI Guardrails

**Purpose:** Create hard guardrails before patching copy/layout so the repo cannot drift back into AI-agent/maintainer wording.

**User value:** Users see product language, not internal audit language.

**Technical scope:**

- Add `docs/product-language-guide.md` or `docs/content/product-language-guide.md` with the rules above, concise and practical.
- Add `scripts/audit-product-language.sh` and `make audit-product-language`.
- Add `scripts/audit-ui-layout.sh` and `make audit-ui-layout` if feasible for static layout checks.
- Audit public pages, README, wiki/docs entry points, and Operations Console templates for banned primary-page wording.
- Add an allowlist for tests, JSON-only fields, evidence/compliance docs, and intentionally advanced/collapsed sections.
- Update `Makefile help` with the new audit targets.

**Files or areas likely to change:**

- `docs/product-language-guide.md`
- `scripts/audit-product-language.sh`
- `scripts/audit-ui-layout.sh`
- `Makefile`
- tests or fixtures as needed

**Tests and validation:**

- `make audit-product-language`
- `make audit-ui-layout` if added
- `make check`
- `git diff --check`

**Risk review:**

- Do not over-broaden banned phrase checks so evidence/compliance docs break for legitimate reasons.
- Do not allow broad exceptions that make the audit useless.

**Expected commit message:**

`Phase 133 -- Checkpoint 000001: define product language and UI guardrails`

### Phase 134 — Public Site Information Architecture And Copy Redesign

**Purpose:** Redesign the public GitHub Pages site as a prospective-user product site and help center.

**User value:** A visitor can understand the product, try it, deploy it, check feeds, connect vehicle data, and find help without reading a maintainer status report.

**Technical scope:**

- Rewrite and redesign:
  - `site/index.html`
  - `site/try-locally.html`
  - `site/ui-tour.html`
  - `site/check-feeds.html`
  - `site/connect-vehicles.html`
  - `site/connector-support.html`
  - `site/readiness.html`
  - `site/video.html`
  - `site/assets/site.css`
- Add a self-hosted/reference deployment page if useful, such as `site/deploy.html` or `site/self-host.html`.
- Update navigation to include local trial and self-hosted/reference operation.
- Replace excessive card grids with task flows, accordions, concise comparison tables, and action strips.
- Remove the flagged wording exactly.
- Reduce unsupported-claim lists to concise “Limits” sections.
- Update README/docs/wiki entry links as needed so user-facing entry points align with the redesigned site.
- Publish/reconcile the actual `gh-pages` branch. Do not leave the live branch stale.

**Files or areas likely to change:**

- `site/*.html`
- `site/assets/site.css`
- `README.md`
- `docs/index.md` or docs entry points as needed
- `wiki/README.md` or wiki entry points as needed
- `gh-pages` branch files

**Tests and validation:**

- `make check-links`
- `make audit-product-language`
- `make audit-ui-layout` if added
- `git diff --check`
- verify `gh-pages` contains the same public-page fixes

**Risk review:**

- Do not turn the site into marketing overclaiming.
- Do not remove important local/self-hosted boundaries; compress and place them well.
- Do not publish stale generated or unrelated files to `gh-pages`.

**Expected commit message:**

`Phase 134 -- Checkpoint 000001: redesign public site for product users`

### Phase 135 — Operations Console Shell And Navigation Redesign

**Purpose:** Fix the oversized header/nav/menu problem and create a compact app shell for real operator work.

**User value:** Operators see their next action and configuration path immediately instead of scrolling through a wall of cards and headings.

**Technical scope:**

- Refactor shared shell templates in `cmd/agency-config/operations.go`.
- Refactor `cmd/agency-config/operations_design_system.go` to support:
  - compact header
  - persistent or collapsible sidebar/navigation
  - compact status bar
  - action list layout
  - primary/secondary content regions
  - accordions for diagnostics
  - fewer card styles
- Keep route paths stable unless a route is proven obsolete and tests are updated.
- Update `operations_navigation.go` and `operations_route_registry.go` labels to plain operator language.
- Update `operations_admin.js` only where needed for collapsible nav or action filters; keep buildless, safe, same-origin behavior.
- Ensure first view at 1366×768 shows useful status and action queue before secondary diagnostics.

**Files or areas likely to change:**

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_design_system.go`
- `cmd/agency-config/operations_navigation.go`
- `cmd/agency-config/operations_route_registry.go`
- `cmd/agency-config/operations_admin.js`
- route/shell tests

**Tests and validation:**

- `go test ./cmd/agency-config`
- `make audit-operations-route-inventory`
- `make audit-product-language`
- `make audit-ui-layout` if added
- `go test ./...`
- `git diff --check`

**Risk review:**

- Do not break authentication, CSRF, route guards, no-store headers, JSON endpoints, or local-login behavior.
- Do not add frontend dependencies or a build pipeline unless clearly justified and validated.
- Do not hide required admin controls behind inaccessible JS-only UI.

**Expected commit message:**

`Phase 135 -- Checkpoint 000001: redesign operations console shell and navigation`

### Phase 136 — Core Admin Page Action-First Redesign

**Purpose:** Convert the actual Operations Console pages from overlapping diagnostic/card surfaces into task-oriented operator workflows.

**User value:** Users know what to do next on each page and do not need to interpret repeated insights and caveats.

**Technical scope:**

Patch primary pages and their models/templates as needed:

- Start/dashboard
- Agency Setup
- Setup Details
- Import GTFS
- GTFS Workbench
- Schedule Quality
- Feed Health
- Validation Center
- Validator Health
- Realtime
- Prediction Lab
- Devices & Tokens
- Telemetry Freshness
- Telemetry Simulator
- Connector Hub
- Connector Workbench
- Connector Checks
- Readiness
- Maintenance
- Help

Required changes:

- Put one primary action and a short “current state” summary near the top.
- Reduce card grids. Use action lists, compact tables, and accordions.
- Move diagnostic detail below the first action path.
- Collapse or remove repetitive “does not prove” paragraphs from cards.
- Remove raw claim-flag tables from visible primary HTML; keep JSON outputs and advanced collapsed areas if needed.
- Make copy practical and non-cringe.
- Keep security boundaries and role gating unchanged.

**Files or areas likely to change:**

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_*go`
- `cmd/agency-config/*_test.go`
- `internal/admincontrol` only if needed for action-list data models

**Tests and validation:**

- Add or update HTML tests for:
  - banned phrase absence on primary pages
  - raw claim flags not visible on primary pages
  - route labels are operator-friendly
  - core pages expose one clear primary action
  - no-store/auth/CSRF behavior remains intact
- `go test ./cmd/agency-config`
- `go test ./...`
- `make audit-operations-route-inventory`
- `make audit-product-language`
- `make audit-ui-layout` if added
- `make audit-product-acceptance`

**Risk review:**

- Do not remove diagnostic capability entirely; move it into better disclosure.
- Do not break JSON companion routes or existing scripts.
- Do not create public admin routes.

**Expected commit message:**

`Phase 136 -- Checkpoint 000001: convert admin pages to action-first workflows`

### Phase 137 — Capture-Backed Tutorial, Screenshots, And Self-Hosted Story

**Purpose:** Replace the static screenshot/PPT-like tour with a real tutorial based on current UI captures and a self-hosted/reference deployment story.

**User value:** Users can learn the product from realistic screens and understand both local evaluation and remote operation.

**Technical scope:**

- Add or update a screenshot capture workflow/script if feasible.
- Generate or update actual screenshots from the current app using synthetic/public-safe data:
  - local sign-in
  - Start/dashboard
  - Agency Setup
  - Import GTFS or GTFS Workbench
  - Feed Health
  - Devices / telemetry setup
  - Realtime
  - Connectors
  - Readiness or Maintenance
  - Self-hosted/reference deployment review path or deployment settings/help page
- Update `site/ui-tour.html` into a tutorial with real step sections.
- Update `site/video.html` so it is a real walkthrough page, not only a local/demo slideshow.
- If video capture is not available, keep the page honest and provide a high-quality screenshot/tutorial flow. Do not create fake PPT-like screenshots or overstate capture.
- Keep private data, tokens, secrets, raw logs, and evidence paths out of screenshots/video.

**Files or areas likely to change:**

- `site/ui-tour.html`
- `site/video.html`
- `site/assets/screenshots/*`
- `site/assets/*.vtt` / video assets only if real capture is produced
- `scripts/capture-ui-tour.sh` or similar if added
- docs/tutorials/video or screenshot guidance

**Tests and validation:**

- Verify screenshots exist, are current enough to match redesigned UI, and are referenced correctly.
- `make check-links`
- `make audit-product-language`
- `git diff --check`
- If screenshot script exists, add help/test mode and run it where feasible.

**Risk review:**

- Do not create fake evidence.
- Do not use private agency data.
- Do not leave broken image/video paths.
- If the environment blocks browser automation, stop and report that blocker rather than fabricating assets.

**Expected commit message:**

`Phase 137 -- Checkpoint 000001: rebuild tutorial with current product captures`

### Phase 138 — Actual Product / Self-Hosted Verification, Not Local-Only Polish

**Purpose:** Prove the redesign applies to the actual product routes and self-hosted/reference operation path, not just the local demo.

**User value:** The product guidance works for the main use case: local evaluation that can progress to self-hosted operation.

**Technical scope:**

- Add or update a product smoke script/test that renders/fetches the Operations Console with deployment-like environment values, such as non-local public feed base URLs and reference deployment labels.
- Verify public copy and admin copy distinguish:
  - local evaluation
  - self-hosted/reference deployment
  - future authorized external sharing
- Update deployment docs to link to the redesigned operator pages and public tutorial.
- Update README/wiki/docs so they do not imply the product is only for local use.
- Confirm in tests or docs that `/admin/operations*` routes are the real product private routes for both local and deployed instances.

**Files or areas likely to change:**

- `scripts/release-candidate-check.sh` or new `scripts/product-ui-smoke.sh`
- `Makefile`
- `docs/deployment/oci-reference-deployment.md`
- `README.md`
- `wiki/*.md`
- `cmd/agency-config/*_test.go`

**Tests and validation:**

- New or updated product UI smoke script
- `make release-candidate-check` if feasible
- `make audit-product-acceptance`
- `make audit-product-language`
- `go test ./cmd/agency-config`
- `git diff --check`

**Risk review:**

- Do not add production-readiness claims.
- Do not require real DNS, TLS, secrets, or external contacts for this verification.
- Keep remote/self-hosted validation synthetic or config-based unless explicit authorization exists.

**Expected commit message:**

`Phase 138 -- Checkpoint 000001: verify admin console product paths beyond local demo`

### Phase 139 — Product Language, Layout, And Site Regression Gates

**Purpose:** Add durable checks so future Codex runs cannot reintroduce overexplained, card-heavy, AI-agent wording.

**User value:** The product keeps a clear user-facing tone across releases.

**Technical scope:**

- Integrate `audit-product-language` into `make check` if lightweight.
- Integrate `audit-ui-layout` into `make check` if lightweight and stable.
- Add focused tests for the public site and Operations Console visible HTML.
- Add a small review checklist for maintainers that is not a phase ledger.
- Make CI changes only if they are lightweight and do not require validators, Docker, network, or browser automation.

**Files or areas likely to change:**

- `Makefile`
- `.github/workflows/*` if appropriate
- `scripts/audit-product-language.sh`
- `scripts/audit-ui-layout.sh`
- `cmd/agency-config/*_test.go`
- `docs/product-language-guide.md`

**Tests and validation:**

- `make check`
- `make audit-product-language`
- `make audit-ui-layout` if added
- `make check-links`
- `go test ./...`
- `git diff --check`

**Risk review:**

- Avoid brittle checks that fail on legitimate evidence/compliance docs.
- Ensure stable branch filtering still behaves as expected if docs/site changes are filtered.

**Expected commit message:**

`Phase 139 -- Checkpoint 000001: add product language and layout regression gates`

### Phase 140 — Operator Issue Center And GTFS-RT Usefulness Improvements

**Purpose:** Bring one own set of improvements beyond copy/layout: reduce overlapping insights by adding a single prioritized operator issue model for feed, GTFS, realtime, connector, and maintenance problems.

**User value:** Operators get a straight line: what is broken, who owns it, what to do, where to click, and when to escalate.

**Technical scope:**

- Add a unified private operator issue/recommendation model, or refactor existing cockpit/readiness/feed-health models into one prioritized list.
- Surface it on Start/dashboard and relevant pages as “Fix these first” or similar.
- Deduplicate overlapping insights from feed health, validation, GTFS quality, telemetry, devices, connectors, readiness, reliability, and maintenance.
- Include issue fields such as:
  - label
  - status/severity
  - owner (`operator`, `administrator`, `deployment owner`, `developer/integrator`)
  - why it matters
  - next action
  - route link
  - source signal
- Improve GTFS-RT usefulness presentation for Vehicle Positions, Trip Updates, and Alerts:
  - publishable / missing / stale / withheld / blocked
  - reason in plain language
  - next fix
  - validator/feed-health connection
- Keep existing feed generation and protobuf contracts stable unless a real bug is found and fixed.

**Files or areas likely to change:**

- `cmd/agency-config/operations.go`
- `cmd/agency-config/operations_feed_health.go`
- `cmd/agency-config/operations_realtime.go`
- `cmd/agency-config/operations_validation_center.go`
- `cmd/agency-config/operations_maintenance.go`
- `cmd/agency-config/operations_connectors.go`
- `internal/compliance/*` only if a shared model belongs there
- tests

**Tests and validation:**

- Unit tests for issue prioritization and deduplication.
- HTML tests for Start/dashboard issue list.
- JSON tests for any new companion output.
- `go test ./cmd/agency-config`
- `go test ./...`
- `make gtfsrt-conformance`
- `make audit-product-language`
- `make audit-product-acceptance`

**Risk review:**

- Do not mutate feed state from a read-only issue center.
- Do not create evidence, consumer status changes, or external contact.
- Do not claim production-grade realtime quality from local diagnostics.

**Expected commit message:**

`Phase 140 -- Checkpoint 000001: add prioritized operator issue center`

## Phase Execution Rules

For each phase:

1. Re-check the current repo state and active branch.
2. Confirm no protected evidence path or consumer status drift before starting.
3. Produce or update a concise phase plan.
4. Optionally run a read-only plan-risk scout when the phase has architecture, UX, content, or route risks.
5. Implement the smallest complete product change for the phase.
6. Run relevant tests, lint, build, route audits, link checks, and product-language/layout checks.
7. Optionally run a read-only diff-review scout before committing.
8. Fix issues found by tests or scouts.
9. Rerun validation.
10. Run `git status` and review the diff.
11. Commit the phase with the expected phase checkpoint message.
12. Produce a concise phase summary with changed files and validation results.

Continue through all phases until Phase 140 is complete unless blocked by:

- missing credentials or secrets;
- destructive operations requiring user approval;
- browser/screenshot capture not available for Phase 137;
- validation failures that cannot be resolved safely;
- repository state that makes continuation unsafe.

If blocked, stop, explain the blocker precisely, and provide the smallest next action required. Do not fake captures, evidence, or publication.

## Git And Commit Policy

At the end of each completed phase:

- run `git status`
- inspect the diff
- ensure no secrets, credentials, generated junk, private data, or unrelated files are included
- ensure protected evidence paths are unchanged unless explicitly authorized, which they are not in this roadmap
- ensure consumer statuses remain unchanged
- run relevant validation
- create a git commit for that phase

Use one commit per phase unless a phase needs a narrow fix commit. If a narrow fix commit is needed, use the same phase number and increment checkpoint number.

Do not commit if validation is failing unless explicitly making a known-blocker checkpoint. If committing with known blockers, the commit message and phase summary must say so clearly.

## Protected Boundaries

Do not:

- write to `docs/evidence/**`
- change `docs/evidence/consumer-submissions/status.json`
- move consumer targets beyond `prepared`
- contact external parties
- collect retained evidence
- fetch or validate a real final public root
- use real vendor/device payloads or credentials
- add hosted SaaS, production-readiness, compliance, agency adoption, consumer acceptance, vendor compatibility, hardware certification, SLA/uptime, public-launch, production AVL reliability, production-grade ETA, or real-world ETA accuracy claims
- remove authentication/authorization/CSRF/no-store boundaries
- make private admin routes public

Required consumer tracker preservation check:

```bash
scripts/check-consumer-tracker.sh
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
```

## Validation Policy

Discover and use the repo’s actual validation commands. At minimum, run the relevant subset per phase and the full set before final closeout where feasible:

```bash
git diff --check
go test ./...
make check
make test
make check-links
make audit-product-acceptance
make audit-final-claim-review
make audit-operations-route-inventory
make audit-product-language
make audit-ui-layout   # if added
make external-connection-check
make adapter-conformance
make gtfsrt-conformance
scripts/check-consumer-tracker.sh
docker compose -f deploy/docker-compose.yml config
```

If a command is unavailable in the environment, record the exact reason and run the closest safe substitute. Do not summarize a skipped command as passed.

## Final Completion Criteria

The final Codex response must include:

- phases completed
- commits created
- files changed
- validation performed
- public site changes, including `gh-pages` publication/reconciliation status
- Operations Console layout and wording changes
- tutorial screenshot/video changes and whether capture was real or blocked
- confirmation that admin console work changed actual product routes, not only the local demo
- confirmation that protected evidence paths and consumer statuses remained unchanged
- known risks and follow-up recommendations

The final response must not claim external evidence, compliance, production readiness, agency adoption, consumer acceptance, hosted service availability, vendor compatibility, SLA/uptime, production AVL reliability, production-grade ETA quality, or real-world ETA accuracy.
