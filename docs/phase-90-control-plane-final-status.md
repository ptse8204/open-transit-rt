# Phase 90 Control Plane Final Status

This document is the final product-track review artifact for the authorized
Phase 75-90 Consumer-Grade Control Plane work.

It is not a release tag, package, image, hosted deployment, retained evidence
packet, consumer action record, final-root proof, compliance proof,
production-readiness proof, vendor/device proof, SLA/uptime proof, or
production-grade ETA proof.

## Status Summary

Phases 75 through 90 close a private browser-first control-plane product track
for Open Transit RT. The track made the self-hosted Operations Console clearer,
safer, and broader for small-agency operations across setup, GTFS import and
review, feed health, validation, realtime diagnostics, telemetry/devices,
connectors, prediction diagnostics, maintenance, roles/access/audit,
nontechnical training, and release-candidate review.

The repository remains evidence-bounded:

- Phase 72 remains `needs_review`, not release-ready.
- Phase 89 remains the current local `v0.1.0-rc.1` gate result and closes as
  `needs_review`.
- No release tag, package, checksum, SBOM/provenance file, image publication,
  GitHub release, retained evidence collection, external contact, consumer
  status movement, or public launch action exists in this track.
- Optional evidence gates remain future work and require separate written
  authorization before any retained proof work.

## Route Inventory

These are private Operations Console or related private admin surfaces. This
inventory is a product review aid only and does not prove public deployment,
consumer action, compliance, or release readiness.

| Area | Private routes |
| --- | --- |
| Start | `/admin/operations`, `/admin/operations.json`, `/admin/operations/launchpad`, `/admin/operations/launchpad.json`, `/admin/operations/setup-wizard`, `/admin/operations/setup-wizard.json`, `/admin/operations/setup` |
| Schedule | `/admin/operations/gtfs-workbench`, `/admin/operations/gtfs-workbench.json`, `/admin/operations/gtfs-import`, `/admin/gtfs-studio`, `/admin/operations/feeds`, `/admin/operations/feed-health`, `/admin/operations/feed-health.json`, `/admin/operations/gtfs-quality`, `/admin/operations/validation-health`, `/admin/operations/validation-health.json`, `/admin/operations/validation-health/refresh.json` |
| Realtime | `/admin/operations/realtime`, `/admin/operations/realtime.json`, `/admin/operations/prediction-lab`, `/admin/operations/prediction-lab.json`, `/admin/operations/telemetry`, `/admin/operations/devices`, `/admin/operations/telemetry-simulator`, `/admin/operations/telemetry-simulator.json`, `/admin/alerts/console` |
| Connectors | `/admin/operations/connectors`, `/admin/operations/connectors.json`, `/admin/operations/connectors/workbench`, `/admin/operations/connectors/workbench.json`, `/admin/operations/connectors/tests`, `/admin/operations/connectors/tests.json` |
| Health | `/admin/operations/validation-center`, `/admin/operations/validation-center.json`, `/admin/operations/readiness`, `/admin/operations/readiness.json`, `/admin/operations/checklist`, `/admin/operations/checklist.json`, `/admin/operations/reliability`, `/admin/operations/reliability.json` |
| Maintain | `/admin/operations/maintenance`, `/admin/operations/maintenance.json`, `/admin/operations/access`, `/admin/operations/access.json`, `/admin/operations/audit`, `/admin/operations/audit.json` |
| Learn | `/admin/operations/help`, `/admin/operations/help.json`, `/admin/operations/consumers`, `/admin/operations/evidence` |
| Progressive enhancement | `/admin/operations/assets/operations.js` |

Route boundary notes:

- Private route inventory does not add a public admin surface.
- Browser POST surfaces remain bounded by existing admin role, auth, CSRF,
  size, agency-scope, and server-owned path/command rules.
- The private JavaScript asset fetches only relative private Operations
  Console routes and stores only local UI preferences.
- Public feed paths remain the existing feed publication paths, not admin
  control routes.

## Feature Inventory

| Phase | Capability added or hardened | Boundary |
| --- | --- | --- |
| 75 | Consumer-Grade Control Plane roadmap pack and future phase prompts. | Planning only; no UI/API/evidence/release work. |
| 76 | Shared private Operations Console design tokens, grouped shell, active navigation, and safer status wording. | Kept Go server-rendered shell; no heavy SPA or public admin route. |
| 77 | Bounded `internal/admincontrol` command result model and private validator-health refresh command. | Private, role-checked, CSRF-aware, no arbitrary command execution. |
| 78 | Buildless progressive UI state for copy controls, bounded filters/sorting, and validation-health refresh. | Relative private routes only; no frontend secret storage. |
| 79 | Agency Setup V3 with profile, diagnostics, role visibility, and technical-helper escalation. | No final-root readiness or agency approval claim. |
| 80 | GTFS Workbench with import history, ZIP identity/checksum, required-file checklist, previews, validation triage, and draft/publish review. | No silent service edits; no compliance claim. |
| 81 | Realtime Operations Center for feed cards, telemetry freshness, device state, matching states, Vehicle Positions, Trip Updates, and Alerts context. | Unknown remains safer than false certainty; no production ETA claim. |
| 82 | Validation Center with schedule/realtime distinction, feed URL review, stale/missing/blocked tool states, and sanitized issue drilldowns. | Validator success remains supporting signal only. |
| 83 | Connector Workbench with recipe chooser, manifest review, dry-run guidance, synthetic normalization preview, and conformance results. | Synthetic/local only; no vendor compatibility or browser network sends. |
| 84 | Prediction & ETA Lab with deterministic fallback, external predictor shadow/fail-closed review, local backtest browsing, and withheld Trip Updates explanation. | No production-grade ETA or real-world accuracy claim. |
| 85 | Maintenance Center V2 with backup/restore guidance, upgrade/rollback review, deployment-doctor summaries, support-bundle guidance, and cadence rows. | No browser-executed destructive maintenance actions. |
| 86 | Multi-agency scope display, role/access guidance, metadata-only audit review, and accessibility shell hardening. | No production multi-tenant readiness claim. |
| 87 | Public feed readiness review, feed URL copy/share guidance, docs portal alignment, and prepared-only consumer packet explanation. | No final-root evidence or consumer status movement. |
| 88 | Nontechnical training: role tours, first-week checklist, glossary, recovery guidance, quick tasks, staff handoff, and operator training guide. | Private guidance only; no public claim expansion. |
| 89 | Local `v0.1.0-rc.1` gate with product, route, connector/backend diagnostics, draft release notes, and blockers. | `needs_review`; no tag/package/image/release action. |
| 90 | Final track status, inventories, validation matrix, blocker matrix, protected-path and claim-boundary review, and future evidence gate stubs. | Closeout only; no retained evidence collection. |

## Validation Matrix

Checkpoint 000003 ran the final validation matrix on 2026-05-14. Package
creation and package audit remain blocked because they require separate
explicit maintainer authorization.

| Check | Status | Notes |
| --- | --- | --- |
| `git status --short` | passed | Clean before final validation and clean after final validation before this file update. |
| `git diff --check` | passed | No whitespace errors. |
| `make check` | passed | Lightweight no-network/no-Docker/no-validator-install checks passed. |
| `make validate` | passed | Validation smoke passed with pinned validator tooling available. |
| `make test` | passed | `go test ./...` passed. |
| `RUN_LOCAL_APP=true make release-candidate-check` | passed | Local diagnostics wrote ignored output under `.cache/release-candidate-check/20260514T112136Z`; package audit remains blocked/not checked unless separately authorized. |
| `make external-connection-check` | passed | Connector manifests and examples remain sidecar/manifest/conformance bounded. |
| `make adapter-conformance` | passed | Adapter conformance suite passed against `testdata/adapter-conformance`. |
| `make test-connector-examples` | passed | Connector example packages passed. |
| `make audit-product-acceptance` | passed | Product acceptance guard passed. |
| `make audit-final-claim-review` | passed | Final claim-boundary guard passed. |
| Consumer tracker JSON parse | passed | `docs/evidence/consumer-submissions/status.json` parsed as JSON. |
| Exact seven-target prepared-only tracker check | passed | All seven targets remain exactly `prepared`: Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, transit.land. |
| Protected-path status check | passed | No status under `docs/evidence/consumer-submissions`, `docs/evidence/captured`, `db/migrations`, `go.mod`, or `go.sum`. |
| `make release-package` | blocked | Requires separate explicit maintainer authorization; not run in Phase 90. |
| `make audit-release-package` | blocked | Requires a package artifact and separate explicit maintainer authorization; not run in Phase 90. |

## Blocker Matrix

| Area | Status | Blocker or residual risk |
| --- | --- | --- |
| Release readiness | needs_review | Phase 89 local diagnostics passed where authorized, but release package creation/audit, tag, image publication, and release action remain blocked/not checked. |
| Phase 72 precedent | needs_review | Phase 72 remains a bounded diagnostic review, not a release-ready pass. |
| Package artifact | blocked | No package, checksum, SBOM/provenance file, GitHub release, or published image exists. |
| Evidence tracks | blocked | Final-root, consumer, agency pilot, vendor/device, ETA-quality, and compliance proof require separate written authorization. |
| Consumer tracker | prepared-only | All seven targets remain prepared; movement beyond prepared requires retained, redacted, target-originated evidence for the named target and feed scope. |
| External connectors | synthetic/local only | Connector/adaptor checks do not prove real vendor compatibility, hardware behavior, or device reliability. |
| ETA quality | diagnostic only | Prediction and backtest surfaces explain diagnostics and withholding reasons; they do not prove production-grade ETA quality or real-world accuracy. |
| Public deployment | not claimed | Local feed and route diagnostics do not prove public launch, managed hosting, SLA/uptime, or final-root readiness. |

## Protected Path Review

Phase 90 must not modify or generate files under:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`

Final CP000003 status: passed. No protected evidence path was modified, and
the protected-path status check for `docs/evidence/consumer-submissions`,
`docs/evidence/captured`, `db/migrations`, `go.mod`, and `go.sum` was clean.

## Consumer Tracker Review

The required consumer tracker state is exactly:

| Target | Required status |
| --- | --- |
| Google Maps | `prepared` |
| Apple Maps | `prepared` |
| Transit App | `prepared` |
| Bing Maps | `prepared` |
| Moovit | `prepared` |
| Mobility Database | `prepared` |
| transit.land | `prepared` |

Final CP000003 status: passed. The tracker JSON parsed successfully and all
seven targets remain exactly `prepared`. Phase 90 did not change the tracker
file or any consumer packet directories.

## Claim-Boundary Review

Allowed final wording:

- private browser-first control-plane product work;
- local/reference operations;
- synthetic/local connector conformance;
- prepared-only consumer packet records;
- release-candidate diagnostics with `needs_review`;
- future evidence gates requiring separate written authorization.

Forbidden final wording:

- CAL-ITP/Caltrans compliance;
- agency adoption or approval;
- consumer submission, review, acceptance, ingestion, listing, or display;
- final-root readiness;
- hosted SaaS or hosted-service availability;
- paid support;
- SLA or uptime guarantee;
- production readiness;
- vendor compatibility;
- hardware certification;
- production-grade ETA quality;
- real-world ETA accuracy;
- public launch completion;
- release readiness.

## Future Evidence Gate Stubs

Each gate below is optional future work. Each requires separate written
authorization before planning, collection, retention, external contact,
submission, status movement, or claim changes.

Required authorization for any gate:

- exact claim target;
- allowed tools;
- public-safe retention rules;
- redaction rules;
- stop conditions;
- operator/agency authorization;
- specific target, vendor, root, feed, or deployment scope when applicable.

| Future gate | Purpose | Required authorization boundary |
| --- | --- | --- |
| Final-root proof | Verify agency-controlled public feed roots and publication metadata. | Requires separate written authorization for the exact root, allowed fetches, retained public-safe artifacts, redaction rules, and stop conditions. |
| Consumer submission | Move one named prepared packet toward a consumer or aggregator workflow. | Requires separate written authorization plus official target path verification and target-originated or operator-retained evidence; prepared status must not move without retained target-specific evidence. |
| Real agency pilot | Record a named agency trial, feedback, decisions, and operational outcome. | Requires separate written agency/operator authorization, public-safe retention rules, redaction, and explicit claim target. |
| Real vendor/device AVL | Test real hardware, vendor payloads, or AVL credentials against the telemetry boundary. | Requires separate written authorization, credential handling rules, payload redaction, stop conditions, and explicit device/vendor scope. |
| Real-world ETA quality | Measure ETA quality against real operations or retained ground truth. | Requires separate written authorization, dataset scope, retention rules, accuracy methodology, and claim boundary approval. |
| Compliance packet | Prepare a claim-specific packet for CAL-ITP-style or other compliance review. | Requires separate written authorization, exact requirement target, retained source artifacts, redaction rules, and claim-boundary review. |

Default without complete authorization: close blocker-only, leave protected
paths unchanged, leave all consumer targets at `prepared`, and make no stronger
claim.
