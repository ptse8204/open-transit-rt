# Post-rc2 Browser-First Product Roadmap Closeout

Date: 2026-05-17 UTC / 2026-05-16 Pacific

The post-rc2 browser-first roadmap is complete through Phase 15. It turns the
post-`v0.1.0-rc.2` repository from a release-candidate backend into a
browser-first local/self-hosted evaluation product for small agencies.

`v0.1.0-rc.2` remains a public release candidate, not a stable release. This
roadmap does not prove production readiness, CAL-ITP/Caltrans compliance,
agency adoption, consumer acceptance, final-root readiness, hosted service
availability, vendor compatibility, hardware certification, SLA/uptime,
production AVL reliability, production-grade ETA quality, or real-world ETA
accuracy.

## Phase Commits

| Phase | Commit | Result |
| --- | --- | --- |
| 01 | `139710f` | Browser-first UI audit and architecture plan |
| 02 | `26dbf61` | Coherent app shell and navigation |
| 03 | `e919b8e` | Browser-first normal agency workflows |
| 04 | `47fe4e2` | GTFS Workbench and quality guidance |
| 05 | `77db1dd` | Realtime Operations Center and feed usefulness |
| 06 | `cde6b13` | External connector support and catalog |
| 07 | `72b0f6e` | CAL-ITP-style readiness workflow |
| 08 | `393c88e` | Human-centered README/docs/wiki cleanup |
| 09 | `c56bae7` | Interactive website and product explainer |
| 10 | `1cd5dfb` | Video tutorial recording workflow |
| 11 | `dad5980` | AI-agent docs separated from human docs |
| 12 | `f3297e8` | Stable branch policy and filtering automation |
| 13 | `60bb1e2` | GitHub Actions Go test repair |
| 14 | `655f6b2` | Browser-first product acceptance walkthrough |
| 15 | this closeout commit | Roadmap closeout and next-release recommendation |

## User-Facing Behavior

After a technical helper starts the local app, agency staff can use the
browser to review setup, import/review GTFS, inspect feed URLs and feed health,
review validation, inspect realtime operations, review devices and telemetry,
use connector guidance, inspect readiness, review maintenance and support
guidance, and follow help/tutorial flows.

The final route walkthrough verified 23 private Operations Console routes, the
GTFS Studio route, the Alerts Console route, all five public feed URLs, and the
unauthenticated admin boundary. Details are recorded in
[Product Acceptance](../../product-acceptance/post-rc2-browser-first-acceptance.md).

## Docs, Wiki, And Website

The README, docs index, wiki home, public site, connector docs, readiness docs,
CI docs, product acceptance docs, and video tutorial guide now lead with normal
human workflows instead of phase history. AI-agent files remain available but
are indexed separately under [AI-Agent Documentation](../../agent/README.md).

The public website is a static, no-tracking explainer with role navigation,
flow cards, generated UI tour panels, connector catalog, readiness explainer,
and video tutorial page.

## Video Workflow

[Video Recording Guide](../../tutorials/video-recording-guide.md) defines
repeatable storyboards for a 3-minute overview, 10-minute local setup,
browser-first GTFS import, feed health/readiness review, connector/AVL
overview, and maintenance/support workflow. It keeps raw and finished video
binaries out of git unless a maintainer separately authorizes release assets or
another storage path.

## Stable Branch And CI

Phase 12 created the `stable` branch from the `v0.1.0-rc.2` baseline and
added `.github/workflows/update-stable.yml` to filter product/user-facing files
from `main` into `stable` without force-pushing. The filter excludes
AI-agent-only docs, handoffs, prompt files, roadmap packs, and phase ledgers.

Phase 13 keeps `go test ./...` in Fast CI and adds a manual release-gates
workflow for validator-heavy checks, connector/conformance checks, product
acceptance, and release-package audit. See [Continuous Integration](../../ci.md).

Post-closeout CI follow-up `d8dfc3b` replaced Go 1.24-era `t.Context()` test
calls with explicit `context.Background()` so GitHub Actions can run the repo's
declared Go `1.23.2` toolchain. The follow-up Fast CI run on `main` passed,
and the update-stable workflow updated remote `stable` to `cf51dd7` with the
same clean-docs filter.

## Connector Support

Connector support is documented in README, docs, website, and UI for:

- Vehicle / GPS / AVL connectors: CSV replay adapter, HTTP polling adapter,
  webhook sidecar adapter, generic JSON transform adapter, vendor-shaped
  synthetic examples, and authenticated `POST /v1/telemetry`.
- Prediction connectors: deterministic built-in predictor, external HTTP
  predictor adapter, shadow-mode predictor, fail-closed behavior, and
  TheTransitClock candidate notes only.
- Validator connectors: MobilityData static GTFS validator, MobilityData GTFS
  Realtime validator, allowlisted validator IDs, and private validation health.
- Monitoring/export connectors: local health summaries, operations notify
  draft, monitoring/export helper, and deployment-owned monitoring boundary.
- Consumer/discovery connectors: `/public/feeds.json`, static GTFS URL,
  Vehicle Positions URL, Trip Updates URL, Alerts URL, and consumer packet
  preparedness without submission or acceptance claims.
- Future extension model: manifest-based sidecars, no arbitrary dynamic backend
  plugin loading, and conformance tests required.

## Readiness

Readiness is organized around public feed URLs, static GTFS, Vehicle Positions,
Trip Updates, Alerts, validation, license/contact metadata, uptime and
operations signals, telemetry/device state, and consumer preparedness. These
workflows support CAL-ITP-style readiness preparation but do not prove
CAL-ITP/Caltrans compliance or external acceptance.

## Validation

Phase 15 closeout validation:

| Command | Result |
| --- | --- |
| `git diff --check` | passed |
| `make check` | passed |
| `make validate` | passed |
| `make test` | passed |
| `make smoke` | passed |
| `make audit-product-acceptance` | passed |
| `make audit-final-claim-review` | passed |
| `make external-connection-check` | passed |
| `make adapter-conformance` | passed |
| `make gtfsrt-conformance` | passed |
| `scripts/check-consumer-tracker.sh` | passed |
| `git status --short` | checked before commit |

Protected evidence paths had no tracked diff. The consumer tracker remained
exactly seven prepared-only targets.

## Remaining Limitations

The technical helper still owns initial startup/shutdown, validator
installation, Docker/TLS/DNS/reverse proxy setup, release packaging, real
device secret handling, external integrations, and any evidence or consumer
workflow.

No claim has moved. Unsupported claims remain unsupported.

## Recommended Next Step

Recommended sequence:

1. Run a release-candidate gate for the post-rc2 browser-first product changes.
2. Review and, if desired, push the filtered `stable` branch after the
   `update-stable` automation has a clean dry run.
3. Start an external connector runtime integration roadmap focused on real
   adapter boundaries, still without vendor compatibility or hardware
   certification claims.
4. Open optional evidence tracks only with separate written authorization,
   retained target-originated evidence rules, and redaction gates.
