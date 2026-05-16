# Phase 114 Web Design Skill UX Review

Status: `phase_114_polish_applied`

This review used
`/Users/edwintse/.agents/skills/web-design-engineer/SKILL.md` for the private
Operations Console UX pass. It records design review and local verification
only. It is not release publication, retained evidence, production readiness,
compliance proof, consumer acceptance, hosted-service availability,
SLA/uptime proof, vendor compatibility, hardware certification, production AVL
reliability, production-grade ETA quality, or real-world ETA accuracy proof.

## Design System Read

- Product surface: private server-rendered Operations Console for small-agency
  operators and technical helpers.
- Visual tone: dense, utilitarian, status-forward, and evidence-bounded.
- Color: existing neutral surfaces and semantic status colors; no new
  decorative palette.
- Type: system UI fonts with compact dashboard headings.
- Layout: responsive grids, bordered sections, tables with horizontal overflow
  on small screens, and no nested decorative card stacks.
- Interaction: no-JS fallback first, progressive JavaScript only for copy,
  review filtering, and bounded private refresh controls.

## Audit Findings

Resolved in Phase 114:

- Missing feed URLs no longer expose a copy action for the literal value
  `missing`.
- Missing first-run feed URL cells now display `Not configured yet`, while the
  underlying `copy_value` stays empty.
- The progressive copy enhancer now rejects empty, `missing`, `unknown`,
  `not available`, and `not configured` sentinel values before adding a copy
  button.
- The first-run realtime task label now uses
  `Realtime feeds: Vehicle Positions, Trip Updates, Alerts` instead of
  `VP/TU/Alerts`.

Reviewed and kept:

- The existing server-rendered design system is appropriate for the operator
  workflow and should not be replaced by a frontend framework for this release
  candidate.
- Wide-table and mobile behavior already uses responsive overflow and compact
  controls; no Phase 114 structural rewrite was needed.
- JSON/debug links remain visible for technical helpers, but release-critical
  first-run copy defects were prioritized for this phase.

Remaining UX candidates for later phases:

- Further separate primary operator actions from JSON/debug links on dense
  pages.
- Consider adding table captions or row-group summaries to very wide review
  tables where the first viewport carries too much diagnostic text.
- Re-check the post-release state in Phase 118 after Phase 115/116/117 release
  outcomes are known.

## Verification

Focused checks passed:

- `node --test cmd/agency-config/operations_admin_test.mjs`
- `go test ./cmd/agency-config -run 'FirstRun|Launchpad|Progressive|FeedHealth|RealtimeOperations|Operations'`
- authenticated local HTML check confirmed the plain-language realtime label
  and no `data-copy-value="missing"` marker.

Local visual review:

- `make agency-app-up` started the local app at `http://localhost:8080`.
- Playwright CLI captured desktop and mobile screenshots into ignored
  `.cache/phase-114-ux/`.
- Desktop viewport checked: `1366x900`.
- Mobile viewport checked: `390x844`.
- Screenshots showed the private Operations Console shell rendering with
  stable header, agency scope, navigation, and mobile stacking.

The Playwright auth storage state and screenshots are local ignored artifacts;
they are not retained evidence and are not committed.

## Claim Boundaries

Phase 114 did not publish a release, create a tag, upload assets, contact
external parties, write retained evidence, move consumer statuses, or modify
protected evidence paths. The review does not claim stable release readiness,
production readiness, compliance, adoption, consumer acceptance, final-root
readiness, hosted service availability, paid support, SLA/uptime, vendor
compatibility, hardware certification, production AVL reliability,
production-grade ETA quality, or real-world ETA accuracy.

