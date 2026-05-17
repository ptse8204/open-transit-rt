# Post-rc2 Browser-First Product Roadmap

This roadmap continues after the public release candidate `v0.1.0-rc.2`.
The goal is to make Open Transit RT feel like a browser-first product for
small agencies while preserving the current evidence and claim boundaries.

`v0.1.0-rc.2` remains a public release candidate for local/self-hosted
evaluation. It is not a stable release and does not prove production readiness,
CAL-ITP/Caltrans compliance, agency adoption, consumer acceptance, hosted
service availability, final-root readiness, vendor compatibility, hardware
certification, paid support, SLA/uptime, production AVL reliability,
production-grade ETA quality, real-world ETA accuracy, or consumer submission,
review, ingestion, listing, display, or acceptance.

## Start Here

Use these files together:

- [Phase Plan](phase-plan.md): the durable 15-phase roadmap and validation
  commands.
- [UI Audit](ui-audit.md): Phase 01 route inventory, workflow gaps, language
  issues, and browser-first information architecture.
- [Closeout](closeout.md): completed phase commits, validation, changed
  behavior, remaining limits, and next recommended path.

## Phase 01 Status

Phase 01 is a planning and audit phase. It does not implement product UI
changes yet. It records what exists, what remains command-line-only after app
startup, and how later phases should reshape the product.

Required Phase 01 artifacts:

- `docs/roadmaps/post-rc2-browser-first-product/README.md`
- `docs/roadmaps/post-rc2-browser-first-product/ui-audit.md`
- `docs/roadmaps/post-rc2-browser-first-product/phase-plan.md`

## Web Design Engineer Skill

The Web Design Engineer skill was loaded before this planning work because the
roadmap changes UI, website, README structure, tutorial flow, screenshots, and
public docs pages. It shaped the Phase 01 direction in these ways:

- Treat the Operations Console as a product surface for agency staff, not only
  an internal diagnostic panel.
- Preserve the repo's existing Go server-rendered UI and no-JS fallback instead
  of introducing a heavy frontend stack.
- Use clear information architecture, consistent page hierarchy, visible next
  actions, and human-readable status language before visual ornament.
- Keep public website and README work concise, role-based, interactive where
  useful, and free of unsupported marketing claims.
- Avoid decorative or generic design patterns; prioritize dense, scannable,
  operations-friendly layouts.

## Phase 07 Report

The Web Design Engineer skill was loaded before the Phase 07 readiness work.
It shaped the implementation toward a scan-friendly browser workflow map that
fits the existing server-rendered Operations Console, keeps route links and
status language compact, and repeats "what this helps with" / "what this does
not prove" language instead of marketing-style compliance claims.

Phase 07 adds ten readiness focus areas to `/admin/operations/readiness`:
public feed URLs, static GTFS, Vehicle Positions, Trip Updates, Alerts,
validation, license/contact metadata, uptime and operations signals,
telemetry/device state, and consumer preparedness. It keeps URL readiness
separate from license/contact readiness and keeps prepared consumer tracker
records separate from runtime workflow notes.

## Phase 08 Report

The Web Design Engineer skill was loaded before the Phase 08 documentation
cleanup. It shaped the work toward a plain, role-based reading path instead of
a large internal route ledger: the README now starts with audience, normal
browser flow, GTFS import, feed URLs, connectors, readiness, and unsupported
claims; `docs/index.md` separates new users, agency staff, technical helpers,
connector developers, maintainers, and AI agents; and the wiki home now points
readers to task guides before project history.

## Phase 09 Report

The Web Design Engineer skill was loaded before the Phase 09 website work. It
shaped the static site toward a concise public product explainer with a shared
local CSS system, role tabs, flow cards, generated browser-capture panels, and
no external scripts, tracking, analytics, or external fonts. The site pages
stay documentation-only and repeat unsupported-claim limits instead of
marketing proof language.

## Phase 10 Report

Phase 10 adds a repeatable video recording workflow for maintainers. The guide
uses the same browser-first route order as the product UI, includes six
storyboards, requires public-safe local/demo data, keeps tokens and private
records off screen, requires captions or transcripts before publication, and
keeps large video binaries outside the repository unless separately
authorized. The static `site/video.html` page now points to the detailed guide.

## Phase 11 Report

Phase 11 separates AI-agent and Codex continuation material from the normal
human reading path. The new `docs/agent/` hub indexes handoffs, roadmap packs,
prompt files, and historical phase ledgers while keeping canonical paths such
as `docs/handoffs/latest.md` in place for compatibility. `docs/index.md`,
`docs/README.md`, `docs/codex-task.md`, `docs/conversation-summary.md`, and
`docs/handoffs/latest.md` now tell human readers to start with browser-first
guides instead of agent handoffs.

## Phase 12 Report

Phase 12 defines `main`, `stable`, and `gh-pages` branch roles in
`docs/branching-and-release-policy.md`. It adds a non-force-push stable sync
workflow with a dry-run dispatch mode and an explicit exclude list for
AI-agent-only docs, handoffs, prompt files, roadmap packs, and phase ledgers.
The stable branch remains a filtered product branch, not a stable release or
production-readiness claim.

## Phase 13 Report

Phase 13 keeps Go tests in CI because `go test ./...` passed locally and
remains useful. The GitHub Actions repair changes the default workflow into a
fast PR/push path with Go tests, `make check`, the consumer tracker script, and
the final claim audit. Validator-heavy, connector, conformance, product
acceptance, and release-package checks move to a manual release-gates workflow.
`docs/ci.md` explains when to use each path.

## Phase 14 Report

Phase 14 started the local evaluator stack with `make agency-app-up`, walked 23
private Operations Console routes, verified GTFS Studio and Alerts Console,
checked all five public feed URLs, and confirmed unauthenticated
`/admin/operations` returns `401`. The acceptance record is
`docs/product-acceptance/post-rc2-browser-first-acceptance.md`.

The Web Design Engineer skill was loaded before writing the acceptance record.
It shaped the report toward a role and workflow summary rather than a raw
diagnostic dump, with local route details kept under `.cache/`.

## Phase 15 Report

Phase 15 closes the roadmap in [Closeout](closeout.md), updates current status
and latest handoff, records final validation, and recommends a release-candidate
gate followed by stable branch review and external connector runtime
integration work. Optional evidence tracks remain separately authorization
gated.

## Phase Acceptance Rule

Do not accept a phase until:

- the phase goal is met;
- validation ran or blockers are recorded;
- protected evidence paths were not touched;
- `docs/evidence/consumer-submissions/status.json` remains the same seven
  prepared-only targets;
- unsupported claims remain unsupported;
- the phase has its own commit.

## Protected Paths

Do not write to:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`
- `docs/evidence/consumer-submissions/status.json`

Do not contact agencies, consumers, vendors, portals, or other external targets
as part of this roadmap unless a later maintainer authorization explicitly
opens an evidence track.
