# Post-rc2 Browser-First Product Phase Plan

This file persists the maintainer-provided post-`v0.1.0-rc.2` roadmap so a
future Codex context reset can resume from the same phase list and boundaries.

Current release baseline:

- `v0.1.0-rc.2` is a public release candidate for local/self-hosted evaluation.
- It is not a stable release.
- It does not prove production readiness, compliance, agency adoption,
  consumer acceptance, final-root readiness, hosted service availability,
  vendor compatibility, hardware certification, SLA/uptime, production AVL
  reliability, production-grade ETA quality, or real-world ETA accuracy.
- Protected evidence paths must remain untouched.
- `docs/evidence/consumer-submissions/status.json` must remain exactly seven
  prepared-only targets unless a later authorized evidence workflow provides
  retained target-originated evidence.

## Grounding

Read these before planning or editing:

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
- `.github/workflows/*`

Also inspect current implementation under:

- `cmd/agency-config/`
- `internal/admincontrol/`
- `internal/connectors/`
- `internal/gtfsrtconformance/`
- `site/`, if present
- `docs/`
- `wiki/`

Before planning or modifying UI, website, README visual structure, tutorial
flow, screenshots, video tutorial flow, or public docs pages, load:

`/Users/edwintse/.agents/skills/web-design-engineer/SKILL.md`

Record in the relevant phase report that the Web Design Engineer skill was
loaded and summarize how it influenced the work.

## Master-Agent Workflow

Operate as the Master Agent. If real sub-agents are available, spawn:

- Planning sub-agent, highest effort: audit UI, docs, website, CI, workflows,
  connector docs, and CAL-ITP-style readiness surfaces; define implementation
  plan, stale docs, repeated docs, confusing language, gaps, route/file/test
  impact, and validation commands.
- Implementation sub-agent, high effort: implement each phase, patch code,
  UI, docs, tests, workflows, and website; keep each phase independently
  reviewable; commit after each phase.
- Review/QA sub-agent, high effort: run validation, inspect UI routes and docs,
  verify claim boundaries, verify protected paths and consumer tracker, and
  require patching until the phase meets the goal.

If real sub-agents are unavailable, simulate these roles explicitly.

The Master Agent must not accept a phase until:

- the phase goal is met;
- validation ran or blockers are documented;
- claim boundaries are preserved;
- protected evidence paths are untouched;
- consumer tracker remains prepared-only;
- the phase commit exists.

## Hard Boundaries

Do not:

- write to `docs/evidence/captured/**`;
- write to `docs/evidence/consumer-submissions/current/**`;
- write to `docs/evidence/consumer-submissions/artifacts/**`;
- write to `docs/evidence/consumer-submissions/packets/**`;
- change `docs/evidence/consumer-submissions/status.json`;
- contact external agencies, consumers, vendors, or portals;
- create new evidence packets;
- claim CAL-ITP/Caltrans compliance;
- claim production readiness;
- claim agency adoption or approval;
- claim consumer submission, review, acceptance, ingestion, listing, or display;
- claim final-root readiness;
- claim hosted SaaS or hosted-service availability;
- claim paid support, SLA, or uptime;
- claim vendor compatibility or hardware certification;
- claim production AVL reliability;
- claim production-grade ETA quality or real-world ETA accuracy;
- claim stable release readiness;
- publish another release or tag unless separately authorized.

Allowed wording:

- Open Transit RT supports CAL-ITP-style readiness workflows.
- Open Transit RT provides local/self-hosted evaluation workflows.
- Open Transit RT provides synthetic/local connector examples and conformance
  checks.
- Open Transit RT does not prove compliance, production readiness, consumer
  acceptance, vendor compatibility, or ETA quality without future retained
  evidence.

## Commit Format

Use these exact phase commit names unless a later repo convention supersedes
them:

```text
Post-rc2 Product Roadmap -- Phase 01: browser-first UI audit and architecture plan
Post-rc2 Product Roadmap -- Phase 02: complete browser-first app shell and navigation
Post-rc2 Product Roadmap -- Phase 03: remove command-line dependency from normal agency workflows
Post-rc2 Product Roadmap -- Phase 04: improve GTFS Workbench and quality guidance
Post-rc2 Product Roadmap -- Phase 05: improve realtime operations and feed usefulness
Post-rc2 Product Roadmap -- Phase 06: clarify and expand external connector support
Post-rc2 Product Roadmap -- Phase 07: improve CAL-ITP-style readiness workflow
Post-rc2 Product Roadmap -- Phase 08: rewrite human-facing docs and reduce repetition
Post-rc2 Product Roadmap -- Phase 09: improve interactive website and product explainer
Post-rc2 Product Roadmap -- Phase 10: add video tutorial recording workflow
Post-rc2 Product Roadmap -- Phase 11: separate AI-agent docs from human docs
Post-rc2 Product Roadmap -- Phase 12: add stable branch policy and docs filtering automation
Post-rc2 Product Roadmap -- Phase 13: repair and rationalize GitHub Actions Go tests
Post-rc2 Product Roadmap -- Phase 14: complete browser-first product acceptance review
Post-rc2 Product Roadmap -- Phase 15: close browser-first product roadmap
```

If the repo requires checkpoint numbering instead, use:

```text
Post-rc2 -- Checkpoint 000001: browser-first UI audit and architecture plan
Post-rc2 -- Checkpoint 000002: complete browser-first app shell and task flow
```

Do not combine phases into one commit.

## Phase 01 - Browser-First UI Audit And Product Architecture

Goal: Audit the current Operations Console and define the browser-first product
architecture.

Required work:

- Load the Web Design Engineer skill.
- Inventory all private UI routes.
- Identify workflows that still require command-line use after the app starts.
- Identify confusing labels, stale internal language, overlong pages, missing
  next actions, and diagnostic-only screens.
- Define a browser-first information architecture for Start Here / Dashboard,
  Setup, GTFS Workbench, Feed Health, Validation, Realtime, Devices / AVL,
  Connectors, Alerts, Prediction / ETA Lab, Maintenance, Help / Tutorials,
  Readiness, and Support Bundle / Troubleshooting.
- Create:
  - `docs/roadmaps/post-rc2-browser-first-product/README.md`
  - `docs/roadmaps/post-rc2-browser-first-product/ui-audit.md`
  - `docs/roadmaps/post-rc2-browser-first-product/phase-plan.md`

Validation:

```bash
git diff --check
make check
scripts/check-consumer-tracker.sh
make audit-operations-route-inventory
make audit-final-claim-review
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 01: browser-first UI audit and architecture plan
```

## Phase 02 - Complete Browser-First App Shell And Navigation

Goal: Make the Operations Console feel like a coherent browser app.

Required work:

- Improve layout, navigation, breadcrumbs, page headings, status cards, and task
  grouping.
- Replace internal wording with user-facing wording.
- Add visible "what to do next" actions on every important page.
- Ensure every major workflow is reachable from the browser.
- Add route-level no-store and private/admin boundaries where missing.
- Keep no-JS fallback usable.
- Improve mobile/tablet usability.
- Add tests for route availability and navigation consistency.
- Avoid internal phrases such as "The product path is:", "claim flags",
  "bounded helper status", and "prepared-only signal" in user-facing UI.
- Prefer natural language such as "Start here to set up your feeds.",
  "Review what is working and what still needs attention.", and "These checks
  help you prepare, but they do not prove external approval."

Validation:

```bash
git diff --check
make check
make test
make smoke
scripts/check-consumer-tracker.sh
make audit-final-claim-review
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 02: complete browser-first app shell and navigation
```

## Phase 03 - Browser-First Setup And No-Command-Line Operations

Goal: After a technical helper starts the app, normal users should not need the
command line.

Required work:

- Build or improve browser UI for agency setup, GTFS upload/import by file and
  safe URL, feed metadata configuration, public feed URL review, validator
  install/status explanation, safe browser-triggered validation review, device
  and token setup guidance, telemetry simulator guidance and safe dry-run where
  appropriate, support bundle generation or browser-guided instructions, and
  maintenance checks/update guidance.
- Do not expose arbitrary command execution in the browser. Only server-owned,
  allowlisted, bounded actions are allowed.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make audit-product-acceptance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 03: remove command-line dependency from normal agency workflows
```

## Phase 04 - GTFS Workbench Usability And Guided Data Quality

Goal: Make GTFS import, review, fix planning, and publish review understandable
to agency staff.

Required work:

- Improve GTFS Workbench UI.
- Add plain-language summaries for required files, row counts, service dates,
  routes, stops, trips, calendars, and import history.
- Add "what changed" views for GTFS versions.
- Improve validation issue triage with likely owner, plain-English meaning,
  suggested fix path, safe next action, and what this does not prove.
- Reduce raw diagnostic exposure.
- Add tests for GTFS Workbench pages and JSON companions.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 04: improve GTFS Workbench and quality guidance
```

## Phase 05 - Realtime Operations Center And Feed Usefulness

Goal: Make realtime feed operations easier to understand and act on.

Required work:

- Improve browser UI and docs for Vehicle Positions, including vehicle count,
  stale vehicles, unmatched vehicles, suppressed vehicles, trip descriptor
  coverage, and why vehicles are not published.
- Improve Trip Updates explanations for generated vs withheld output, fallback
  reasons, stale/ambiguous/low-confidence handling, and prediction source.
- Improve Alerts review for active alerts, stale alerts, missing cancellation
  or disruption links, and service disruption review.
- Add realtime health summary: what is healthy, what needs attention, and what
  is not proven.
- Improve synthetic/local realtime replay guidance.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make gtfsrt-conformance
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 05: improve realtime operations and feed usefulness
```

## Phase 06 - External Connector Support And Connector Catalog

Goal: Make external connector support concrete, visible, and useful.

Required work:

- Spell out exact connector categories in README, docs, website, and UI:
  vehicle / GPS / AVL connectors; CSV replay adapter; HTTP polling adapter;
  webhook sidecar adapter; generic JSON transform adapter; vendor-shaped
  synthetic examples; authenticated `POST /v1/telemetry`; prediction
  connectors; deterministic built-in predictor; external HTTP predictor
  adapter; shadow-mode predictor; fail-closed predictor behavior;
  TheTransitClock candidate notes only; validator connectors; MobilityData
  static GTFS validator; MobilityData GTFS Realtime validator; allowlisted
  validator IDs; private validation health; monitoring/export connectors;
  local health summaries; operations notify draft; monitoring/export helper;
  deployment-owned monitoring boundary; consumer/discovery connectors;
  `/public/feeds.json`; static GTFS URL; Vehicle Positions URL; Trip Updates
  URL; Alerts URL; consumer packet preparedness without submission/acceptance
  claim; manifest-based sidecars; no arbitrary dynamic backend plugin loading;
  conformance tests required.
- Improve `/admin/operations/connectors`.
- Improve `/admin/operations/connectors/workbench`.
- Add a connector catalog page in docs and website.
- Make examples easier to copy and adapt.
- Add connector conformance tests where missing.

Validation:

```bash
git diff --check
make check
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 06: clarify and expand external connector support
```

## Phase 07 - CAL-ITP-Style Readiness Workflow Without Overclaiming

Goal: Move the product closer to CAL-ITP-style readiness workflows while
preserving truthfulness.

Required work:

- Improve browser readiness flow.
- Organize readiness around public feed URLs, static GTFS, Vehicle Positions,
  Trip Updates, Alerts, validation, license/contact metadata, uptime/operations
  signals, telemetry/device state, and consumer preparedness.
- Add clearer "what this helps with" and "what this does not prove."
- Add concise docs explaining CAL-ITP-style readiness in plain language.
- Link to official/source docs only where already present or clearly allowed.
- Do not claim CAL-ITP/Caltrans compliance.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 07: improve CAL-ITP-style readiness workflow
```

## Phase 08 - Human-Centered README, Docs, And Wiki Cleanup

Goal: Make documentation shorter, clearer, less repetitive, and easier for
humans.

Required work:

- Rewrite README in normal, fluent language.
- Remove or rewrite phrases such as "The product path is:", "bounded", "claim
  flags", "prepared-only signal", and "source-of-truth matrix" unless they are
  in AI-agent/internal docs.
- Put the most important user flow first: What is this? Who is it for? Try it
  locally. Open the browser UI. Import GTFS. Check feeds. Connect vehicles.
  Review readiness. What is not proven yet?
- Reduce repeated docs.
- Add a docs index that clearly separates new users, agency staff, technical
  helpers, connector developers, maintainers, and AI agents.
- Merge or archive duplicate docs where safe.
- Add redirect or "start here instead" notes for old long docs where deletion
  is risky.
- Keep all claim boundaries intact.

Validation:

```bash
git diff --check
make check
make audit-final-claim-review
make audit-product-acceptance
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 08: rewrite human-facing docs and reduce repetition
```

## Phase 09 - Interactive Website And Product Explainer

Goal: Make the website more interactive, clearer, and focused on essential
flows.

Required work:

- Load the Web Design Engineer skill.
- Improve the GitHub Pages/site experience.
- Make homepage concise and human.
- Add interactive flow cards: Try locally; Start in browser; Import GTFS;
  Review feed health; Connect vehicle data; Review readiness; Maintain the
  system.
- Add a UI tour with screenshots or generated browser captures.
- Add role navigation for agency manager, operations staff, technical helper,
  connector developer, and maintainer.
- Add connector catalog page.
- Add CAL-ITP-style readiness explainer.
- Add video tutorial page.
- Keep the website a documentation and product explainer surface, not a
  marketing website, public launch, hosted service, or production-readiness
  claim.
- Use only redacted local screenshots or generated browser captures that do not
  show admin tokens, private URLs, private evidence, consumer artifacts,
  credentials, personal data, or private operator details.
- Avoid marketing overclaims.
- Avoid external scripts, tracking, analytics, external fonts, or unnecessary
  dependencies unless already approved.

Validation:

```bash
git diff --check
make check
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

If a site-specific test or preview command exists, run it too.

Commit:

```text
Post-rc2 Product Roadmap -- Phase 09: improve interactive website and product explainer
```

## Phase 10 - Video Tutorial Recording Workflow

Goal: Add a repeatable video tutorial workflow so maintainers can record short
demos.

Required work:

- Add `docs/tutorials/video-recording-guide.md`.
- Add `site/video.html` or equivalent.
- Include scripts/storyboards for a 3-minute overview, 10-minute local setup,
  browser-first GTFS import, feed health/readiness review, connector/AVL
  overview, and maintenance/support workflow.
- Add a recording checklist for local app startup, demo data, avoiding secrets,
  avoiding private data, reset state, captions/transcript, and storing generated
  video files outside the repo or in release assets.
- Make clear that captures must not include admin tokens, private URLs, private
  evidence, consumer artifacts, credentials, personal data, or private operator
  details.
- If feasible and safe, add a browser route or static tutorial page that walks
  through the recording script.
- Do not commit large binary video files unless explicitly authorized.

Validation:

```bash
git diff --check
make check
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 10: add video tutorial recording workflow
```

## Phase 11 - AI-Agent Documentation Separation

Goal: Separate AI-agent/Codex docs from human-facing docs.

Required work:

- Inventory docs meant mainly for AI agents: handoffs, phase ledgers, master
  planner exports, Codex prompts, long phase status files, and internal roadmap
  packs.
- Move or index them under `docs/agent/`, `docs/agent/handoffs/`,
  `docs/agent/roadmaps/`, and `docs/agent/prompts/`.
- Keep human docs short and central.
- Add an AI-agent docs README explaining which files agents should read, which
  files humans should not need to read, and how to avoid stale-doc drift.
- Update links.
- Do not delete important history unless safe and validated.
- Keep current status and latest handoff discoverable for Codex.

Validation:

```bash
git diff --check
make check
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 11: separate AI-agent docs from human docs
```

## Phase 12 - Stable Branch And GitHub Automation

Goal: Create a stable branch and automate keeping it clean of AI-agent-only
docs.

Required work:

- Inspect current branch/release policy.
- Create or document creation of a stable branch from the appropriate release
  baseline.
- Add GitHub workflow automation so that when main is updated, stable can be
  updated with product/user-facing files while excluding AI-agent-only docs.
- Avoid destructive force pushes unless explicitly safe.
- Add a dry-run mode or workflow dispatch input.
- Exclude or filter `docs/agent/**`, long phase handoffs intended only for
  agents, Codex prompt artifacts, and internal planning ledgers.
- Exclude protected evidence and consumer-submission paths from automation:
  `docs/evidence/captured/**`, `docs/evidence/consumer-submissions/current/**`,
  `docs/evidence/consumer-submissions/artifacts/**`,
  `docs/evidence/consumer-submissions/packets/**`, and
  `docs/evidence/consumer-submissions/status.json`.
- Run a protected-path scan before any stable-branch update or archive-like
  publication step.
- Preserve README, user docs, website, source code, examples, tests, and
  release docs.
- Document branch policy in `docs/branching-and-release-policy.md`.

Validation:

```bash
git diff --check
make check
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

If workflow syntax validation is available, run it.

Commit:

```text
Post-rc2 Product Roadmap -- Phase 12: add stable branch policy and docs filtering automation
```

## Phase 13 - GitHub Actions And Go Test Repair

Goal: Evaluate current GitHub Actions Go tests, then fix, remove, or replace
them if they are not useful.

Required work:

- Inspect `.github/workflows/*`, `Makefile`, and Go test targets.
- Determine why the current Go test workflow is broken.
- Decide whether CI should run `go test ./...`, `make test`, lighter checks,
  install validators, avoid validators, or split fast CI from release gates.
- Fix the workflow if useful.
- Remove or replace it if completely useless at this stage.
- Add CI docs explaining fast PR checks, local full checks, release-candidate
  checks, and why validator-heavy checks are or are not in CI.
- Make GitHub Actions align with the current repo stage and avoid blocking on
  unavailable local-only tools unless intentionally configured.
- Prefer fast CI with checkout, setup Go, `go test ./...`, `make check`,
  `scripts/check-consumer-tracker.sh`, and `make audit-final-claim-review`.
- Prefer optional/manual release CI for validator install, `make validate`,
  `make smoke`, and release package audit.

Validation:

```bash
git diff --check
go test ./...
make check
make test
make audit-final-claim-review
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 13: repair and rationalize GitHub Actions Go tests
```

## Phase 14 - Final Product Acceptance And Route Walkthrough

Goal: Run a full product acceptance review from a small-agency user
perspective.

Required work:

- Start the app locally.
- Walk through browser routes for `/admin/operations`, setup, GTFS workbench,
  GTFS import, feeds, feed health, validation, realtime, devices, telemetry,
  connectors, readiness, maintenance, and help.
- Verify normal workflows are possible from browser UI after startup.
- Identify remaining command-line-only workflows.
- Patch remaining UI gaps if feasible.
- Record final product acceptance in
  `docs/product-acceptance/post-rc2-browser-first-acceptance.md`.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make audit-product-acceptance
make audit-final-claim-review
make external-connection-check
make adapter-conformance
make gtfsrt-conformance
scripts/check-consumer-tracker.sh
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 14: complete browser-first product acceptance review
```

## Phase 15 - Roadmap Closeout And Next-Release Recommendation

Goal: Close the roadmap and recommend the next release/evidence/product path.

Required work:

- Update `docs/current-status.md`, `docs/handoffs/latest.md`, relevant roadmap
  docs, and README if needed.
- Summarize all phases completed.
- Record validation.
- Record changed user-facing behavior.
- Record remaining limitations.
- Recommend whether the next step should be another product hardening roadmap,
  a release-candidate gate, optional evidence track, stable branch promotion,
  external connector runtime integration, or another scoped path.
- Preserve all claim boundaries.

Validation:

```bash
git diff --check
make check
make validate
make test
make smoke
make audit-product-acceptance
make audit-final-claim-review
make external-connection-check
make adapter-conformance
make gtfsrt-conformance
scripts/check-consumer-tracker.sh
git status --short
```

Commit:

```text
Post-rc2 Product Roadmap -- Phase 15: close browser-first product roadmap
```

## Final Report Required

When complete, report:

- Final branch and commit.
- Every phase commit.
- What changed in the UI.
- What normal users can now do without command line after app startup.
- What changed in README/docs/wiki/site.
- What website changes were made.
- What video tutorial workflow was added.
- How AI-agent docs were separated.
- Stable branch and automation status.
- GitHub Actions / Go test status.
- External connector support now documented, including exact connector
  categories.
- CAL-ITP-style readiness improvements.
- Validation commands and results.
- Protected evidence path status.
- Consumer tracker status.
- Unsupported claims that remain unsupported.
- Next recommended roadmap.
