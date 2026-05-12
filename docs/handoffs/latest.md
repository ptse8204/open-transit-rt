# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phases 0 through 60 remain closed for their documented scopes. After the
Post-60 agency-ready productization audit, the maintainer authorized Phase 61+
as the forward product roadmap naming for agency-first UI, guided setup,
connector-platform, release, and product-polish work.

The Phase 61+ agency-first connector platform roadmap lives at
`docs/roadmaps/agency-first-connector-platform/README.md`; Codex should start
from `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`
before planning new product phases.

Phase 61 -- Agency-First UI And Connector Hub is complete. Phase 62 -- Guided
Setup and Browser GTFS Import is complete. Phase 63 -- Feed Health and
Readiness UX is complete. Phase 64 -- Connector Platform and SDKs is complete.
Phase 65 -- Operator Workflow and Data Quality UX is complete. Phase 66 --
Release Candidate and Installability is complete. Phase 67 -- Product Polish,
Accessibility, and In-App Help is complete. Phase 68+ -- Optional Authorized
Evidence Tracks is closed blocker-only / authorization-gated for the current
no-authorization review after Checkpoint 000002. Phase 69 -- Maintainer
Product Acceptance And UI-First Agency Usability Trial is complete. Phase 70
-- GitHub Pages Product Explainer Site is complete. Phase 71 --
Adoption-First Productization And No-CLI Agency Operations is complete. The
current default next work is `v0.1.0-rc.1` agency evaluation release-candidate
hardening: run the clean-checkout release-candidate gate, review the browser
Agency Operations Cockpit / Start Here path, verify off-host validation and
reference diagnostics, and review connector/adaptor conformance. The canonical
review is `docs/roadmap-status.md#review-and-recommendations`. Any future
evidence intake is optional and requires explicit written authorization first.

Phase 61 added the roadmap pack, replaced stale "not Phase 61" source-of-truth
wording, improved `/admin/operations` with agency-first action cards, and
added private read-only Connector Hub routes at `/admin/operations/connectors`
and `/admin/operations/connectors.json`. Connector Hub explains telemetry,
prediction, validator, monitoring/export, and consumer/discovery connector
paths, and defines plugins as optional sidecars, manifests, command adapters,
or connector processes rather than arbitrary dynamic backend code loading.

Phase 61 created no retained evidence, wrote nothing under protected evidence
paths, contacted no external party, changed no consumer status, changed no
public feed URL, added no migration, and made no compliance, agency approval,
agency adoption, consumer acceptance, hosted-service, paid-support, SLA,
production-readiness, vendor-compatibility, hardware-certification, public
launch, production AVL reliability, or production-grade ETA claim. All seven
consumer and aggregator targets remain `prepared`.

Phase 62 added `/admin/operations/setup-wizard` and
`/admin/operations/setup-wizard.json` as private setup guidance, plus
`/admin/operations/gtfs-import` as an admin-only browser GTFS import path for
ZIP upload or safe URL import. The route derives
agency and actor from the authenticated principal, stores raw ZIP bytes only in
temporary runtime storage, reuses `internal/gtfs.ImportService.ImportZip`,
renders bounded import/validation results, blocks client-supplied paths and
validator/argv fields, and keeps CLI import and GTFS Studio available. It does
not add public routes, migrations, evidence writes, consumer status changes, or
stronger compliance/adoption/acceptance/hosted-service/vendor/production/ETA
claims.

Phase 63 made feed health and CAL-ITP-style readiness easier to understand
without requiring operators to read diagnostic tables. It added private
read-only feed health routes at `/admin/operations/feed-health` and
`/admin/operations/feed-health.json`, summarizing `feeds.json`, schedule,
Vehicle Positions, Trip Updates, and Alerts with plain-language status,
freshness, validator context, health context, next actions, and "does not
prove" boundaries. It also added private readiness checklist v2 at
`/admin/operations/readiness` and `/admin/operations/readiness.json`. The JSON
returns only the v2 model, keeps all claim flags false, and explains readiness
item, status, current signal, meaning, why it matters, next action, what it
does not prove, and admin/doc links for discovery, feed health, GTFS quality,
Vehicle Positions, Trip Updates, Alerts, validator health, reliability,
telemetry/devices, scorecard, and consumer prepared tracker signals. Phase 63
is closed. Do not claim SLA/uptime proof or CAL-ITP/Caltrans compliance.

Phase 63 closeout lives in
`docs/phase-63-feed-health-and-readiness-ux.md`. Phase 64 closeout lives in
`docs/phase-64-connector-platform-and-sdks.md` and
`docs/handoffs/phase-64.md`. Phase 65 closeout lives in
`docs/phase-65-operator-workflow-and-data-quality-ux.md` and
`docs/handoffs/phase-65.md`. Phase 66 closeout lives in
`docs/handoffs/phase-66.md`. Phase 67 closeout lives in
`docs/handoffs/phase-67.md`. Phase 68+ closeout lives in
`docs/handoffs/phase-68-plus.md`. Phase 69 closeout lives in
`docs/handoffs/phase-69.md`.

Phase 67 -- Product Polish, Accessibility, and In-App Help is complete.
Phase 67 Checkpoint 000001 added the scoped product polish, accessibility, and
in-app help plan in
`docs/phase-67-product-polish-accessibility-in-app-help.md`, corrected stale
current-roadmap wording, and kept protected evidence/status paths untouched.
Checkpoint 000002 grouped the private Operations Console navigation by
operator intent, added active-page state, preserved existing route paths, and
updated the Launchpad decision-gate docs link to the current Phase 61+
roadmap. It did not add public routes, evidence writes, consumer status
changes, runtime contract changes, or stronger claims.
Checkpoint 000003 improved accessibility-oriented shared layout markup,
visible keyboard focus styles, responsive/mobile constraints, table overflow
behavior, and explicit form labels and submit buttons while preserving
existing POST names, actions, role checks, CSRF behavior, and private route
boundaries.
Checkpoint 000004 added private GET-only help routes at
`/admin/operations/help` and `/admin/operations/help.json`, a static/derived
help model for GTFS, GTFS-RT, connectors, readiness, validators, telemetry,
and claim/evidence boundaries, contextual shared-layout help panels, and
all-false help claim flags.
Checkpoint 000005 closed Phase 67 with validation, protected-path,
consumer-tracker, and claim-boundary review.

Phase 68+ optional authorized evidence-track readiness is closed
blocker-only / authorization-gated for the current no-authorization review.
Checkpoint 000001 documented the optional tracks and gates. Checkpoint 000002
closed the phase with docs/status/handoff updates only. No retained evidence
was collected, no external party was contacted, no final public root was
fetched or verified, no consumer status was moved, and protected evidence paths
and consumer tracker status records remained unchanged.

The current default next work is `v0.1.0-rc.1` product readiness and
external-connection maturity, not real agency pilot work or public proof by
default. The gate should cover a clean checkout, `make check`,
`make validate`, `make test`, local app startup, a public GTFS trial when
allowed, five public feed fetches, validator health, telemetry simulator,
connector/adaptor conformance, and final claim audit. If no explicit written
authorization, exact claim target, allowed tools, public-safe retention plan,
redaction rules, and stop conditions exist, keep evidence tracks closed as
authorization-gated blockers without collecting retained evidence or contacting
anyone.

Phase 69 improved the UI-first product acceptance path without evidence
intake. It added the private Operations Console first-click label, `Agency Operations Cockpit / Start Here`,
five-feed URL copy section, no-developer/developer paths, browser-first
small-agency walkthroughs, README/wiki/docs task navigation,
capability-versus-evidence cleanup, and `make audit-product-acceptance` /
`make test-product-acceptance`. It created no retained evidence, contacted no
external party, changed no consumer status, and made no compliance, agency
adoption, consumer acceptance, final-root, hosted SaaS,
production-readiness, vendor-compatibility, SLA, or ETA-quality claim.

Phase 71 improved adoption-first product quality without evidence intake. It
added the Agency Operations Cockpit JSON/HTML surface, maintenance routes,
feed health realtime-usefulness rows, GTFS import source review, browser-first
validator/quality/telemetry guidance, `make oci-reference-check`,
`make validate-public-feeds`, and adoption docs for no-CLI first run,
maintenance, off-host validation, and OCI/reference diagnostics. It created no
retained evidence, contacted no external party, changed no consumer status,
added no public admin routes, and made no compliance, agency adoption,
consumer acceptance, final-root, hosted-service, production-readiness,
vendor-compatibility, SLA, uptime, or ETA-quality claim.

Phase 66 checkpoint history: Checkpoint 000001 added the Phase 66 plan in
`docs/phase-66-release-candidate-and-installability.md`.
Checkpoint 000002 prepared the first release-candidate workflow in
`docs/release-candidate-readiness.md`, the release process docs, and the
private `.cache` release-candidate diagnostic summary.
Checkpoint 000003 added a local bootstrap preflight and clearer first-run
blocker guidance in `scripts/bootstrap-dev.sh`, quickstart docs, and tests.
Phase 64 Checkpoint 000001 added the Phase 64 plan in
`docs/phase-64-connector-platform-and-sdks.md`. Checkpoint 000002 added the
private Connector Hub manifest registry UI and bounded JSON registry model
from committed synthetic example manifests. Checkpoint 000003 added private
GET-only connector test instructions at `/admin/operations/connectors/tests`
without backend command execution. Checkpoint 000004 improved telemetry
connector SDK-style examples with a shared synthetic dry-run normalization
helper. Checkpoint 000005 improved the prediction connector SDK-style example
with sanitized dry-run request/response helpers, Vehicle Positions
independence, no-public-mutation flags, and withheld-output diagnostics.
Checkpoint 000006 added a monitoring/export SDK-style helper and synthetic
validator allowlist example. Checkpoint 000007 closed Phase 64 with
validation, protected-path, consumer-tracker, and claim-boundary review.
Phase 65 improved device/vehicle onboarding, telemetry simulator UX, GTFS
quality guidance, and operator troubleshooting without changing protected
runtime contracts or creating production fleet reliability/compliance claims.
Checkpoint 000001 added the Phase 65 plan in
`docs/phase-65-operator-workflow-and-data-quality-ux.md`. The next checkpoint
after Checkpoint 000001 was device and vehicle onboarding UI. Checkpoint
000002 improved `/admin/operations/devices` with guided use cases, non-admin
mutation-form gating, one-time token boundaries, and per-device telemetry
freshness/assignment next actions. Checkpoint 000003 added
`/admin/operations/telemetry-simulator` and
`/admin/operations/telemetry-simulator.json` as private GET-only guide
surfaces for committed synthetic scenarios and copyable operator-shell
commands. Checkpoint 000004 improved private GTFS quality guidance with likely
owners, affected files, safe fix paths, verification steps, escalation
triggers, and all-false claim flags. Checkpoint 000005 closed Phase 65 with
validation, protected-path, consumer-tracker, and claim-boundary review.
Checkpoint 000001 added the Phase 66 release-candidate and installability plan.
Checkpoint 000002 prepared the first release-candidate workflow with an
ordered review sequence, validation matrix, release-note inputs, package audit
matrix, and private diagnostic summary fields. Checkpoint 000003 added a local
bootstrap preflight, clearer missing-tool and Docker/DB readiness messages, and
first-run blocker docs. Checkpoint 000004 documented the Docker image
publishing decision: no published app image, no registry push, and no hosted
service claim. Checkpoint 000005 added the docs/demo site plan as
repository-native planning only with no hosting, marketing launch, evidence
collection, or public-launch claim. Checkpoint 000006 closed Phase 66 with
validation, protected-path, consumer-tracker, and claim-boundary review.
Optional evidence tracks remain future authorization-gated work and are not the
default continuation path.

Post-60 Checkpoints 000009 through 000012 completed the agency-ready
productization pass:

- added the MIT license and rewired the README/wiki front door around agency,
  civic technologist, operator, and developer evaluator paths;
- added agency-facing wiki pages for fit, trial checklist, connector cookbook,
  CAL-ITP-style readiness, and how agencies can help;
- added `examples/README.md`, refreshed `testdata/README.md`, and improved
  contribution guidance for connector examples and first PRs;
- added `make help`, lightweight no-network `make check`, Taskfile convenience
  aliases, lean no-claim GitHub Actions checks, and release-candidate
  environment blocker documentation.

Formal agency approval, final feed-root evidence, and consumer acceptance are
not required to use or improve the software. They remain future evidence
milestones only for agencies that choose public launch or compliance claims.

Phase 60 — Final Claim Review And Public Closeout is complete for the approved
docs/scripts/status/handoff scope.

Track A is also closed for its docs-only external-proof workflow scope.

Phase 60 added the final claim-to-evidence review, unsupported-claim table,
retained-evidence boundary, `scripts/audit-final-claim-review.sh`,
`scripts/test-final-claim-review.sh`, `make audit-final-claim-review`,
`make test-final-claim-review`, validation scaffolding, and
`docs/handoffs/phase-60.md`. The audit is local-only and read-only: it scans
bounded public/status docs, required Phase 60 sections, unsafe private strings,
unsupported positive claims, the exact seven-target prepared-only consumer
tracker state, and README-only consumer artifact directories.

Phase 60 created no retained evidence, wrote nothing under `docs/evidence`,
contacted no external party, sent no webhook or email, automated no portal,
changed no consumer status, refreshed no consumer packet or artifact directory,
and added no public launch, compliance, agency adoption, consumer acceptance,
hosted service, paid support, SLA/uptime, production-readiness, vendor,
marketplace, or ETA-quality claim. All seven consumer and aggregator targets
remain `prepared`.

Before any future stronger public wording, rerun `make audit-final-claim-review`
and require retained public-safe claim-specific evidence for the exact claim.

Before productization stages, reread the repository source-of-truth docs listed
in `AGENTS.md` and this handoff, then keep changes inside the stage scope. Do
not write retained evidence, contact external parties, automate portals, or
change consumer tracker statuses without a separately approved evidence intake.

Previous phase context: Phase 59 — Real Pilot Closeout is closed blocker-only
for the approved docs/status/handoff scope.

Do not reopen earlier phases unless a blocking truthfulness, safety, security,
realtime-quality, evidence, agency-boundary, auth, data-isolation,
agency-domain, device/AVL onboarding, admin-UX, operations-hardening,
pilot-readiness, submission-readiness, public-messaging, public-GTFS evidence,
or self-hosted-reuse issue directly requires it.

Phase 59 found no retained Phase 59 pilot authorization record, kickoff note,
agency/operator feedback record, operations closeout, or continue/pause/close
decision artifact in the repository. It closed blocker-only, created no pilot
evidence packet, wrote nothing under `docs/evidence`, contacted no external
party, changed no consumer records or artifact directories, and left all seven
consumer and aggregator targets at `prepared`.

Existing OCI and local public-GTFS pilot packets remain earlier-scope evidence
only. They do not prove real pilot authorization, agency feedback, agency
adoption, final-root readiness, consumer acceptance, CAL-ITP/Caltrans
compliance, production readiness, hosted service availability, vendor
compatibility, SLA/uptime, or production-grade ETA quality.

Phase 58 added `docs/vendor-equivalent-pack/` templates,
`scripts/audit-vendor-equivalent-pack.sh`,
`scripts/test-vendor-equivalent-pack.sh`, Make targets, validation
scaffolding, and `docs/handoffs/phase-58.md`. The template pack covers BYOD
hardware intake, implementation planning, support boundaries, SLA/KPI planning,
and procurement responses.

The audit helper verifies required templates, boundary phrases, unsafe private
strings, unsupported positive claim wording, placeholders, and the
seven-target prepared-only consumer tracker.

Phase 58 created no retained evidence, submitted to no marketplace, contacted
no marketplace, vendor, consumer, or agency, published no paid support offer,
changed no consumer-submission current records or artifact directories, changed
no `docs/evidence/consumer-submissions/status.json`, wrote nothing under
`docs/evidence/captured`, and made no marketplace approval, paid support,
vendor compatibility, hardware certification, hosted service, hosted SaaS,
SLA/uptime, production-readiness, compliance, agency adoption, consumer
acceptance, or production-grade ETA claim.

Phase 32 produced draft public launch materials only. No announcement was posted, no social copy was published, no agency was contacted, no reporter was contacted, no consumer or aggregator was contacted, and no public launch occurred.

The post-Phase-32 final-root evidence follow-up is complete as
blocker-documented only. No agency-owned or agency-approved final public feed
root was available, no root was used, and no owner/approval evidence was
available. No DNS, TLS, redirect, public feed fetch, validator, proxy/config,
packet README, or checksum evidence was collected.

Phase 33 added a public GTFS local/pilot evidence phase and templates,
attempted the preferred LA Metro Bus GTFS local run, fixed the large-import
timeout exposed by that run, and retried Outcome C. The final retained packet
downloaded the public GTFS ZIP to ignored `.cache/` storage, recorded the source
checksum, imported the feed through `cmd/gtfs-import -agency-id LACMTA`, fetched
`/public/gtfs/schedule.zip`, verified the fetched schedule as the imported LA
Metro public GTFS rather than the repo sample feed, fetched all five public
paths, recorded validator results or blockers, recorded a telemetry dry-run
summary, and checked admin/private route boundaries.

Phase 45 added authenticated `/admin/operations/gtfs-quality` triage. It
separates canonical MobilityData static validator results from Open Transit RT
internal import validation, uses bounded derived groups/samples, supports
admin-only rerun against the active schedule through the server-side
`static-mobilitydata` mapping, and keeps raw reports/stdout/stderr/argv/private
paths out of rendered HTML. The page is diagnostics only: no evidence packet,
consumer status, public route, GTFS auto-edit, consumer acceptance claim, or
CAL-ITP/Caltrans compliance claim was added.

Phase 46 added authenticated `/admin/operations/validation-health` and
`/admin/operations/validation-health.json`, plus `scripts/validator-health.sh`
and `make validator-health`. The health summary is bounded, emits exactly four
feed rows in fixed order, separates canonical static MobilityData validation
from realtime MobilityData validation, and keeps Open Transit RT internal GTFS
import validation out of canonical validator health. Admin-only `run_all` is
CSRF-aware, form-size-capped, strict about accepted fields, and uses only
server-owned validator mappings, artifacts, and timeouts. The script defaults
to `.cache/validator-health/<timestamp>`, refuses evidence-like output paths,
validates summary/manifest JSON, and does not call private admin routes without
`ADMIN_TOKEN`. The deployment doctor may GET validation-health JSON with a safe
admin URL and token, but never POSTs it. Phase 46 created no evidence packet,
changed no consumer statuses, auto-edited no GTFS, blocked no publish, and
added no compliance, consumer acceptance, agency adoption, hosted SaaS,
production-readiness, vendor-compatibility, or production-grade ETA claim.

Phase 47 added `scripts/operations-notify.sh` and `make operations-notify` for
private local notification drafts from existing validator-health and
deployment-doctor summaries. The script defaults to
`.cache/operations-notify/<timestamp>`, writes exactly `summary.json`,
`summary.md`, `manifest.json`, `manifest.md`, and `notification.txt`, caps
source/output sizes, rejects symlink and evidence-like paths, records webhook
and email presence as booleans only, and supports strict mode for local
automation. It sends no notifications, calls no webhook/email/admin route, runs
no validators, requires no DB/Docker/app/admin token, creates no evidence,
writes nothing under `docs/evidence`, contacts no consumers, changes no
consumer statuses, blocks no publish, auto-edits no GTFS, and adds no
compliance, consumer acceptance, agency adoption, hosted SaaS,
production-readiness, vendor-compatibility, or production-grade ETA claim.

Phase 48 added mutually exclusive `--dry-run` and `--send` modes to
`cmd/avl-vendor-adapter`. Dry-run preserves the existing no-network
stdout/stderr JSON behavior. Send mode validates an `avl-adapter-send.v1`
manifest with env-only token references, reads target/config from the Phase 48
`AVL_ADAPTER_*` env contract, validates `/v1/telemetry` targets and safe
non-evidence output paths, blocks stale/future batches, optionally blocks other
warnings via `AVL_ADAPTER_FAIL_ON_WARNINGS=true`, posts one transformed
`telemetry.Event` per request with bearer auth, retries only retryable
failures, stops on first terminal failure, marks later records
`skipped_after_failure`, and writes exactly `summary.json`, `summary.md`,
`manifest.json`, `manifest.md`, and `diagnostics.json` under `.cache/` by
default. The send diagnostics are redaction-safe private operator artifacts:
they use `record_index` and deterministic `credential_ref`, omit raw response
bodies, raw mapped IDs, token values, raw private host/path details, and false
claim flags remain false. Phase 48 did not change `/v1/telemetry`
payload/auth semantics, add public/admin APIs, add queues/schedulers/daemons,
write under `docs/evidence`, change consumer statuses, add named vendor
support, or add compliance, consumer acceptance, agency adoption, hosted SaaS,
production-readiness, vendor-compatibility, production AVL reliability, or
production-grade ETA claims.

Phase 49 added the optional generic external HTTP Trip Updates predictor
runtime behind `internal/prediction.Adapter`. The shared factory/config path is
used by both `cmd/feed-trip-updates` and `cmd/agency-config`; deterministic
remains default and `noop` remains available. `external_http` posts sanitized
DTOs only to a fixed `/v1/predict/trip-updates` endpoint with exact host
allowlisting, safe URL validation, no redirects, HTTPS except loopback test
stubs, byte/time caps, and env-name-only bearer token lookup. Failures produce
valid empty Trip Updates with `adapter_error` diagnostics.
`external_http_shadow` keeps deterministic public output and stores only
bounded redacted shadow diagnostics. Phase 49 added no TheTransitClock-specific
runtime, named service call, external process manager, evidence writes,
consumer status changes, compliance claim, vendor compatibility claim, or
production-grade ETA claim.

Phase 50 added `internal/realtimequality` backtesting support,
`cmd/realtime-quality-backtest`, synthetic fixtures under
`testdata/realtime-quality-backtest`, and `make realtime-quality-backtest`.
The workflow writes only private `.cache/realtime-quality-backtest/<timestamp>`
diagnostics with exactly `summary.json`, `summary.md`, `metrics.json`,
`metrics.md`, and `manifest.json`. It rejects evidence-like paths, unsafe
traversal, and symlink ancestors, and does not copy raw inputs or persist raw
rows. Phase 50 is local engineering diagnostics only and does not add DB
persistence, migrations, Operations Console routes, public APIs, evidence
writes, consumer tracker changes, publish gates, or stronger ETA/readiness
claims.

Phase 51 added authenticated GET-only `/admin/operations/reliability` and
`/admin/operations/reliability.json`, a fixed private reliability summary
model, DB-backed feed and incident rollups from existing tables only, bounded
best-effort Vehicle Positions health persistence into
`feed_health_snapshot`, `scripts/operations-reliability.sh`, and
`make operations-reliability`. The summary keeps feed rows ordered as
`schedule`, `vehicle_positions`, `trip_updates`, `alerts`; uses only `ok`,
`needs_review`, `missing`, `unknown`, and `unhealthy`; and keeps missing data
from becoming `ok`. The script defaults to
`.cache/operations-reliability/<timestamp>` and writes exactly
`summary.json`, `summary.md`, `manifest.json`, `manifest.md`, and
`reliability-review.txt`. Phase 51 added no public route, migration,
monitoring-stack dependency, evidence write, consumer-status change, publish
gate, SLA claim, uptime guarantee, production-readiness claim, compliance
claim, hosted SaaS claim, agency adoption claim, vendor compatibility claim,
consumer acceptance claim, or production-grade ETA claim.

Phase 52 added the guarded final-root workflow from
`docs/phase-52-final-public-root-evidence-workflow.md`. The repo now has
`scripts/collect-final-root-evidence.sh`,
`scripts/audit-final-root-evidence.sh`, `make collect-final-root-evidence`,
`make audit-final-root-evidence`, final-root workflow templates under
`docs/evidence/templates/`, and local-only script coverage. The default
collector output stays under ignored `.cache/final-root-evidence/<timestamp>`.
Retention under `docs/evidence/captured` requires explicit opt-in, a real final
root, `ALLOW_CAPTURED_EVIDENCE_WRITE=true`, and a readable redacted approval
artifact. Phase 52 closed blocker-only because no real final root or approval
artifact was available in repo evidence; no real final-root evidence was
retained, `docs/evidence/captured` was not changed, prepared consumer packets
were not refreshed, and consumer statuses remain unchanged.

Phase 53 closed the authorized consumer submission execution scope
blocker-only. No local operator authorization artifact, official target path
verification artifact, or target-originated artifact exists, so no target was
selected, no consumer or aggregator was contacted, no portal was automated or
scraped, no submission path was guessed, no submission was made, and no
artifact was added. All seven consumer and aggregator targets remain
`prepared`. Phase 53 did not change
`docs/evidence/consumer-submissions/status.json`, current target records,
target artifact directories, or `docs/evidence/captured`.

Phase 54 closed the official requirements refresh scope as docs-only. It
re-checked current public Caltrans / Cal-ITP and FTA sources and refreshed
requirement mappings for stable public GTFS Schedule and GTFS Realtime URLs,
all three standard realtime feed types, canonical no-error validation, major
trip-planner acceptance as a separate third-party requirement, open license
visibility, provider or regional source-of-truth website links, technical/feed
contact, Transitland and Mobility Database availability, and realtime API-key
registration constraints. Phase 54 did not change code, migrations, runtime
behavior, public routes, auth, validators, evidence files, consumer current
records, consumer artifact directories, or
`docs/evidence/consumer-submissions/status.json`. It did not create
compliance evidence, change consumer statuses, or claim compliance, consumer
acceptance, final-root readiness, agency adoption, hosted SaaS availability,
production readiness, vendor compatibility, SLA/uptime, marketplace approval,
or production-grade ETA quality.

Phase 33 evidence is completed only for local/pilot public static GTFS dataset
handling. It does not prove agency adoption, agency approval, official agency
feed status, agency-owned final-root readiness, consumer evidence, compliance,
hosted SaaS, production readiness, real LA Metro realtime data, real vendor AVL
compatibility, or ETA quality.

Phase 34 aligned post-Outcome-C status and roadmap docs, documented the
original Phase 33 static validator blocker, recorded the later Homebrew Java 17
static validator retry result with 0 system errors and 3 warnings,
added/verified final-root request and public-GTFS repeatability guidance, and
added this phase handoff. Phase 34 created no external evidence and made no
runtime code, script, Makefile target, schema, migration, API, consumer tracker,
final-root evidence packet, target artifact, or OCI pilot final-root wording
changes.

Phase 35 restored the root README as the Open Transit RT product front door and
realigned roadmap/handoff docs around self-hosted agency reuse. The default next
work is now Phase 36 — OCI/OCL Reference Deployment Productization. External
proof tracks remain documented as future optional paths when retained,
claim-specific artifacts exist, but they are not the default roadmap. Phase 35
created no external evidence and made no runtime code, script, Makefile target,
schema, migration, API, consumer tracker, final-root evidence packet, target
artifact, validator artifact, or public feed contract changes.

Phase 36 created the self-hosted OCI/OCL-style reference deployment docs under
`docs/deployment/`, closed the Phase 36 phase reference, and updated roadmap,
status, backlog, open-question, tutorial, README, and handoff navigation. Phase
36 made no runtime behavior, script, Makefile target, schema, migration, API,
consumer tracker, final-root evidence packet, target artifact, OCI pilot
evidence packet, validator artifact, or public feed contract changes. No
external evidence was created and no consumer statuses changed. The next
recommended phase is Phase 37 — Reusable Agency Onboarding Flow.

Phase 37 added `scripts/agency-pilot-onboard.sh`, `make agency-pilot-up`, and
explicit Compose env interpolation defaults for agency/public feed settings.
The onboarding helper accepts an agency ID and GTFS URL, downloads the ZIP into
ignored `.cache/` storage, records a checksum, safely upserts the requested
agency/admin roles without using the demo seed, imports with configurable
timeout, bootstraps explicit publication metadata, verifies all five public
paths, checks the fetched schedule summary against the imported source, and
runs validators or reports blockers/skips. Phase 37 does not call
`make agency-app-up`, does not import the demo sample feed, does not create
external or final-root evidence, and does not change consumer statuses. The
post-completion hardening patch makes running mode upsert agency/admin rows via
`DATABASE_URL`, requires explicit running-mode `ADMIN_BASE_URL`, rejects
`.`/`..`/leading-dot agency IDs, adds a no-network dry-run validation check,
and passes advanced options through `make agency-pilot-up`. At Phase 37 close,
the recommended next phase was Phase 38 — Integration Adapter Kit.

Phase 38 added `docs/integration-adapter-kit.md` as the central adapter map,
linked it from onboarding/tutorial/navigation docs, added neutral synthetic AVL
fixtures and `testdata/avl-vendor/README.md`, refreshed the dry-run adapter CLI
help boundary wording, added focused fixture/help diagnostics tests, updated
`make validate` to check the adapter kit files, and recorded the synthetic
dry-run/developer support status in `docs/dependencies.md`. Phase 38 did not
add network send mode, named vendor support, real vendor payloads, credentials,
runtime external predictor integration, Prometheus/Grafana assets,
OpenTelemetry wiring, consumer APIs, final-root evidence, external evidence
packets, or consumer status changes. At Phase 38 close, the recommended next
phase was Phase 39 — CAL-ITP-Style Readiness Workflow.

Phase 39 added an authenticated Operations Console readiness page at
`/admin/operations/readiness`. The page shows CAL-ITP-style readiness rows for
stable public URLs, static GTFS, Vehicle Positions, Trip Updates, Alerts,
license/contact metadata, validation status, telemetry freshness, operations
status, and consumer packet preparedness. Each row includes a status, status
source, current evidence/signal, next action, and claim boundary. Phase 39
updated readiness, onboarding, deployment, adapter, roadmap, and handoff docs
and added validation checks for the Phase 39 docs/handoff. Phase 39 did not
create external evidence, did not change consumer statuses, did not add a
public unauthenticated route, and did not add any CAL-ITP/Caltrans compliance,
consumer acceptance, agency adoption, final-root proof, hosted SaaS,
production-readiness, vendor-compatibility, or production-grade ETA-quality
claim.

Phase 40 added a guided self-hosted operator trial tutorial that ties together
the Phase 36 reference deployment docs, Phase 37 reusable agency onboarding
flow, Phase 38 integration adapter kit, and Phase 39 readiness workflow. The
trial covers local/reference prep, `make agency-pilot-up`, a no-external-network
`demo-agency` fixture option, five public feed checks,
`/admin/operations/readiness`, validator run/skip/blocker handling, synthetic
AVL dry-run, next actions, and teardown. Phase 40 was docs/navigation only; the
only non-doc file changed was `Makefile` validation scaffolding. It created no
external evidence, changed no consumer statuses, and added no compliance,
consumer acceptance, agency approval/adoption, final-root, hosted SaaS,
production-readiness, vendor-compatibility, or production-grade ETA claim.

Phase 41 added strict operator smoke checks and redaction-safe support bundles
for local/reference diagnostics. It created no external evidence, changed no
consumer statuses, and added no compliance, consumer acceptance, agency
approval/adoption, final-root, hosted SaaS, production-readiness,
vendor-compatibility, or production-grade ETA claim.

Phase 42 added `scripts/deployment-doctor.sh`, `make deployment-doctor`, and
reference deployment doctor docs for read-only OCI/OCL-style deployment
diagnostics. The doctor inspects already-exported env vars only, records
env/secret presence without values, checks public feed metadata, public/private
route boundaries, loopback health endpoints, read-only DB/migration/PostGIS
status when supplied, validator tooling, backup/restore readiness, git/release
identity, and the prepared-only consumer tracker guard. Default mode exits `0`
while reporting blockers/skips/unavailable checks; `STRICT_DOCTOR=true` is the
mode that fails on blockers. Phase 42 created no external evidence, changed no
consumer statuses, and added no compliance, consumer acceptance, agency
approval/adoption, final-root, hosted SaaS, production-readiness,
vendor-compatibility, or production-grade ETA claim.

Phase 43 added a private authenticated Operations Console checklist at
`/admin/operations/checklist` and `/admin/operations/checklist.json`. Both
routes are derived from one deterministic model with fixed group order
(`setup`, `feeds`, `validation`, `telemetry`, `operations`,
`consumer_workflow`), stable row IDs, neutral statuses, source/current-signal
fields, next actions, claim boundaries, repo-relative docs links, heuristic
metadata/URL labels, and explicit false claim flags. Dashboard, setup, and
readiness pages link to both checklist routes. Phase 43 also patched local
routing so exact `/` returns the local app message with `200`, unmatched local
paths return `404`, and the deployment doctor checks `/admin/gtfs-studio`
instead of exact `/admin/gtfs`. Phase 43 created no external evidence, changed
no consumer statuses, and added no compliance, consumer acceptance, agency
approval/adoption, final-root proof, hosted SaaS, production-readiness,
vendor-compatibility, or production-grade ETA claim.

Phase 44 added `cmd/telemetry-simulator`, `scripts/telemetry-simulator.sh`,
`make telemetry-simulator`, and synthetic fixtures under
`testdata/telemetry-simulator/`. The simulator loads synthetic-only scenarios,
posts to the real authenticated `POST /v1/telemetry` path with device
bearer-token auth, records accepted/duplicate/out-of-order/rejected ingest
results, and can optionally run the existing DB-backed matcher plus Vehicle
Positions debug builder after accepted HTTP ingest. The Phase 44 safety closure
requires HTTPS for non-loopback credentialed sends, permits non-loopback HTTP
only for dry-run validation, confines outputs to repo `.cache` unless
explicitly overridden, always rejects `docs/evidence`, rejects symlink output
directories, redaction-scans generated diagnostics, and records private timing
fields. Phase 44 created no evidence packet, changed no consumer statuses,
added no real vendor payloads or private telemetry, and added no
vendor-compatibility, production AVL reliability, real realtime data,
production-grade ETA, CAL-ITP/Caltrans compliance, agency approval/adoption,
hosted SaaS, or production-readiness claim.

## Phase 32 Summary

- Added `docs/agency-one-pager.md` with problem, solution, audience, current capabilities, pilot path, requirements, readiness boundaries, evidence boundaries, and agency next steps.
- Added `docs/demo-video-outline.md` with a truthful local demo script covering startup, GTFS import/demo feed, public feed URLs, Operations Console setup, device telemetry or dry-run adapter path, validation/evidence view, consumer packet boundary, and pilot next step.
- Added `docs/public-share-copy.md` with draft-only short, medium, and longer copy for GitHub launch, agency/evaluator, contributor, and transit/open-data audiences.
- Added `docs/ecosystem-positioning.md` covering GTFS/GTFS Realtime, validators, Caltrans/CAL-ITP-style readiness, downstream consumers and aggregators, agency-owned domains, external predictor adapters, AVL/vendor adapters, and open-source transit tooling.
- Added `docs/public-launch-checklist.md` with public-message safety checks, no-logo/no-affiliation rule, and claim-to-evidence table.
- Updated README, docs navigation, Phase 32 status, roadmap status, current status, and this latest handoff.
- Added `docs/handoffs/phase-32.md`.

## Phase 33 Summary

- Added `docs/phase-33-public-gtfs-local-pilot-evidence.md` defining Outcome A,
  Outcome B, and Outcome C.
- Added template-only evidence packet scaffolding under
  `docs/evidence/captured/public-gtfs-local-pilot/templates/`.
- Added Outcome C local/pilot evidence packet at
  `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/`.
- Updated docs navigation, roadmap status, phase plan, current status, and this
  latest handoff.
- Kept the top-level README unchanged; docs navigation and evidence docs point
  to the retained packet.

## Phase 34 Summary

- Updated post-Outcome-C status and roadmap docs so Phase 33 Outcome C is no
  longer described as attempted or missing.
- Refreshed `docs/repo-gaps.md` from historical starter scaffolding gaps to
  current evidence and product gaps.
- Updated `docs/phase-plan.md`, `docs/future-roadmap-post-outcome-c.md`, and
  `docs/track-b-productization-roadmap.md` so the next path is an evidence
  fork, not another Phase 33 retry or Phase 32.
- Verified docs navigation links for the final-root request, public-GTFS
  local/pilot runbook, and Phase 33 Outcome C packet.
- Fixed final-root example URLs to use clear placeholder domains.
- Added `docs/handoffs/phase-34.md`.

## Phase 35 Summary

- Replaced the roadmap-export root README with an Open Transit RT product front
  door.
- Added clear README paths for local evaluation, public-GTFS local/pilot runs,
  and the OCI/OCL-style reference deployment path.
- Updated roadmap/status docs so self-hosted agency reuse and OCI/OCL reference
  deployment productization are the default continuation path.
- Reframed external-proof tracks as future optional paths without deleting the
  existing external-proof docs.
- Patched `docs/phase-plan.md` so Phase 34 records both the original missing
  Java static-validator blocker and the later Homebrew Java 17 retry result
  without calling it validator-clean or compliance evidence.
- Updated `docs/backlog.md` and `docs/open-questions.md`.
- Added `docs/handoffs/phase-35.md`.

## Phase 36 Summary

- Added `docs/deployment/README.md` as the deployment index for the
  self-hosted OCI/OCL-style reference path.
- Added `docs/deployment/oci-reference-deployment.md` with the operator
  copy/paste path, server prerequisites, user/group layout,
  `/opt/open-transit-rt` directory layout, placeholder-only env guidance,
  secret generation, database setup/migrations, systemd supervision, Caddy or
  equivalent reverse proxy, public/private route boundaries, validator
  install/check, GTFS import, five public feed URL verification,
  backup/restore, feed monitor, scorecard export, update/rollback, and
  redacted support bundle guidance.
- Added `docs/deployment/oci-reference-env.example` with grouped placeholder
  sections and generated-secret placeholders only.
- Added `docs/deployment/oci-reference-smoke-checklist.md` for repeatable
  operator verification without converting smoke output into evidence.
- Updated `docs/phase-36-oci-reference-deployment-productization.md`,
  `README.md`, `docs/current-status.md`, `docs/README.md`,
  `docs/tutorials/README.md`, `docs/roadmap-status.md`,
  `docs/track-b-productization-roadmap.md`, `docs/backlog.md`,
  `docs/open-questions.md`, and this latest handoff.
- Added `docs/handoffs/phase-36.md`.

## Phase 37 Summary

- Added `scripts/agency-pilot-onboard.sh` for reusable local/reference agency
  onboarding from `AGENCY_ID` and `GTFS_URL`.
- Added `make agency-pilot-up AGENCY_ID=... GTFS_URL=...`.
- Updated `deploy/docker-compose.yml` with explicit interpolation defaults for
  `AGENCY_ID`, `PUBLIC_BASE_URL`, `FEED_BASE_URL`,
  `SCHEDULE_FEED_URL`, `VEHICLE_POSITIONS_FEED_URL`,
  `TRIP_UPDATES_FEED_URL`, `ALERTS_FEED_URL`, and
  `REALTIME_VALIDATION_BASE_URL`.
- Added `docs/tutorials/reusable-agency-onboarding.md`.
- Updated README, tutorial, deployment, roadmap, backlog, status, and handoff
  navigation.
- Added `docs/handoffs/phase-37.md`.
- Post-completion hardening patched running-mode agency/admin upsert,
  running-mode admin URL requirements, agency ID validation, no-network dry-run
  validation, and Makefile option pass-through.

## Phase 38 Summary

- Added `docs/integration-adapter-kit.md` as the central map for telemetry,
  synthetic AVL, predictor, validator, monitoring, consumer/feed workflow, and
  redaction/evidence boundaries.
- Linked the kit from the root README, docs index, tutorials index, reusable
  agency onboarding guide, and device/AVL tutorial.
- Added neutral synthetic AVL fixtures and `testdata/avl-vendor/README.md`
  with fixture purposes, expected diagnostic codes, and evidence boundaries.
- Refreshed `cmd/avl-vendor-adapter` help text to dry-run adapter kit wording
  and added focused help/diagnostic/no-send-mode tests.
- Updated `docs/dependencies.md`, `docs/phase-38-integration-adapter-kit.md`,
  roadmap/status/backlog/open-question docs, `make validate`, and this handoff.
- Reviewed `docs/decisions.md`; no architecture-significant decision changed,
  so no edits were needed.

## Phase 39 Summary

- Added `/admin/operations/readiness` as an authenticated Operations Console
  page.
- Rendered ten readiness rows with status source, evidence/signal, next action,
  and claim boundary.
- Added explicit wording that the page supports CAL-ITP-style readiness
  workflows and does not claim CAL-ITP/Caltrans compliance.
- Added focused Operations Console tests for routing, agency scoping, row
  rendering, claim boundaries, and prepared-only consumer packet semantics.
- Updated readiness, deployment, onboarding, adapter, roadmap/status, and
  handoff docs.
- Updated `make validate` to check Phase 39 docs/handoff files.

## Phase 40 Summary

- Added `docs/tutorials/self-hosted-operator-trial.md` as a guided
  local/reference operator checklist.
- Added `docs/phase-40-guided-self-hosted-operator-trial.md` and
  `docs/handoffs/phase-40.md`.
- Linked the trial from README, docs index, tutorial index, deployment,
  onboarding, adapter, readiness, roadmap/status, backlog, open-question, phase
  plan, and handoff navigation.
- Updated `make validate` to check Phase 40 docs/handoff files.
- Created no external evidence and changed no consumer statuses.

## Phase 41 Summary

- Added `scripts/operator-smoke.sh` and `make operator-smoke` for strict
  local/reference smoke checks.
- Added `scripts/support-bundle.sh` and `make support-bundle` for
  redaction-safe diagnostics that can run even when the app is unavailable.
- Added `docs/tutorials/operator-smoke-and-support-bundle.md` and
  `docs/handoffs/phase-41.md`.
- Operator smoke checks the five public feed paths, unauthenticated admin
  boundary behavior, optional authenticated readiness through a safe admin URL,
  validator tooling state, optional allowlisted validation API summaries, and
  the deterministic synthetic AVL dry-run fixture.
- Support bundles store summaries, not raw feed bodies or full validation
  reports, and run a final redaction scan for secret-shaped values.
- Updated README, docs index, tutorial index, deployment navigation,
  integration adapter navigation, roadmap/status, backlog, open-question,
  current-status, and handoff docs.
- Updated `make validate` to check Phase 41 scripts and docs.
- Created no external evidence and changed no consumer statuses.

## Phase 42 Summary

- Added `scripts/deployment-doctor.sh` and `make deployment-doctor` for
  read-only OCI/OCL-style reference deployment diagnostics.
- Added `docs/deployment/reference-deployment-doctor.md`,
  `docs/phase-42-reference-deployment-doctor.md`, and
  `docs/handoffs/phase-42.md`.
- The doctor inspects already-exported env vars only and does not source
  private env files.
- It writes `summary.json`, `summary.md`, `manifest.json`, and `manifest.md`
  under ignored `.cache/deployment-doctor/<timestamp>/`.
- It checks reference env key presence, generated-secret status, URL safety,
  public feed fetch metadata, public/private route boundaries, optional
  authenticated readiness, HTTPS/redirect posture, service health,
  read-only DB/migration/PostGIS status, validator tooling, backup/restore
  readiness, git identity, and prepared-only consumer tracker status.
- Default mode exits `0` while reporting blockers/skips/unavailable checks;
  `STRICT_DOCTOR=true` fails on blockers.
- Updated README, docs index, tutorial/deployment navigation, roadmap/status,
  backlog, open-question, current-status, and handoff docs.
- Updated `make validate` to check the Phase 42 script, help path, and docs.
- Created no external evidence and changed no consumer statuses.

## Phase 43 Summary

- Patched the local Caddy fallback so exact `/` returns `200` and unmatched
  local paths return `404`.
- Patched the deployment doctor private route list from exact `/admin/gtfs` to
  `/admin/gtfs-studio`.
- Added static validation guards for the deployment-doctor route and local
  Caddy fallback shape.
- Added `/admin/operations/checklist` and `/admin/operations/checklist.json`
  as authenticated private routes derived from one shared checklist model.
- Added deterministic setup, feeds, validation, telemetry, operations, and
  consumer_workflow groups with stable row IDs, neutral statuses, next actions,
  claim boundaries, docs links, heuristic labels, and exact false claim flags.
- Added setup/readiness/dashboard navigation links to both checklist routes.
- Added focused tests for route registration, auth roles, agency scoping,
  method rejection, JSON shape, headers, deterministic ordering, classifiers,
  HTML escaping, docs-link safety, forbidden wording, consumer prepared-only
  semantics, deployment-doctor route regression, and local Caddy fallback
  shape.
- Added `docs/phase-43-operator-ux-setup-v2.md` and
  `docs/handoffs/phase-43.md`.
- Created no external evidence and changed no consumer statuses.

## Phase 44 Summary

- Patched the stale Current Objective wording in this handoff before starting
  Phase 44.
- Added `cmd/telemetry-simulator` for synthetic scenario loading, authenticated
  `/v1/telemetry` posting, status expectation checks, private diagnostics, and
  optional post-ingest matcher/Vehicle Positions debug diagnostics.
- Added `scripts/telemetry-simulator.sh` and `make telemetry-simulator`.
- Added synthetic fixtures for on-route, stale, out-of-order, unknown-device,
  low-quality GPS, after-midnight, and block-transition scenarios.
- Added command tests for dry-run redaction, authenticated `/v1/telemetry`
  posting, and evidence-directory output rejection.
- Patched URL safety, output-directory safety, final redaction scanning, timing
  diagnostics, and multi-event count consistency tests before starting Phase
  45.
- Added `docs/tutorials/telemetry-simulator-and-device-trial.md`,
  `docs/phase-44-telemetry-simulator-and-device-trial.md`, and
  `docs/handoffs/phase-44.md`.
- Updated README, docs navigation, current status, backlog, open questions,
  roadmap docs, adapter docs, and operator tutorials.
- Updated `make validate` to check simulator script/docs/fixtures and dry-run
  command behavior.
- Created no external evidence and changed no consumer statuses.

## Checks Run For Phase 43

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- Consumer tracker JSON syntax, exact target set, prepared-only statuses, and
  byte-for-byte unchanged check — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make agency-app-up` — passed.
- `PUBLIC_BASE_URL=http://localhost:8080 ADMIN_BASE_URL=http://localhost:8080 make deployment-doctor` — passed.
- Direct local route checks passed: `/` returned `200`, `/metrics` returned
  `404`, `/not-a-real-route` returned `404`, `/admin/gtfs` returned `404`,
  `/admin/gtfs-studio` returned `401`, and
  `/admin/debug/gtfsrt/vehicle_positions.json` returned `401`.
- `make agency-app-down` — passed.

## Checks Run For Phase 44

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` — passed; the tracker file is unchanged.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make agency-app-up` — passed.
- Local simulator check with
  `TARGET=http://localhost:8080 DEVICE_TOKEN=dev-device-token SCENARIO=testdata/telemetry-simulator/on-route.json RUN_MATCHER=true` — passed with a repeat-run-safe `out_of_order` ingest result because existing local telemetry was newer than the fixture timestamp.
- Fresh local simulator check with `REFERENCE_TIME=2026-05-11T15:05:00Z`
  and `RUN_MATCHER=true` — passed with HTTP `201`, ingest status
  `accepted`, matcher output, and private Vehicle Positions debug
  `trip_descriptor_published=true`.
- Phase 44 safety patch before Phase 45 — `make validate`, `make test`,
  `git diff --check`,
  `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null`,
  `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`,
  and `docker compose -f deploy/docker-compose.yml config` passed.
- Phase 44 safety patch local simulator check with
  `TARGET=http://localhost:8080 DEVICE_TOKEN=dev-device-token SCENARIO=testdata/telemetry-simulator/on-route.json RUN_MATCHER=true REFERENCE_TIME=2026-05-11T15:05:00Z make telemetry-simulator` — passed with HTTP `202`, ingest status `duplicate`, private timing fields, and no evidence/status changes because the reused local demo database already contained the same accepted synthetic timestamp from earlier Phase 44 verification.
- `make agency-app-down` after local verification — passed.

## Truthfulness And Evidence Boundary

- All seven consumer and aggregator targets remain `prepared` only.
- Phase 30 selected no target and made no submissions.
- No target has submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption evidence.
- No agency-owned or agency-approved final public feed root exists in repo evidence.
- The post-Phase-32 final-root evidence follow-up confirmed the final-root blocker remains unresolved and created no evidence packet.
- Phase 33 is Outcome C for local/pilot public-GTFS dataset handling only. It
  does not prove agency adoption, official agency feed status, final-root proof,
  consumer evidence, compliance, production readiness, real realtime data, or
  ETA quality.
- Phase 34 is status consistency and evidence-readiness only. It created no new
  external evidence and does not strengthen Phase 33 claims.
- Phase 35 is docs-only README and roadmap realignment. It created no new
  external evidence and does not strengthen any public claim.
- Phase 36 is docs-only reference deployment productization. It created no new
  external evidence, did not run a deployment, and does not strengthen any
  public claim.
- Phase 37 is reusable local/reference onboarding productization. It created no
  external evidence, did not create final-root proof, did not change consumer
  statuses, and does not strengthen agency approval, consumer acceptance,
  compliance, production-readiness, vendor-compatibility, or ETA-quality
  claims.
- Phase 38 is integration adapter kit productization. It created no external
  evidence, did not create final-root proof, did not change consumer statuses,
  and does not strengthen vendor-compatibility, production AVL reliability,
  production ETA quality, consumer acceptance, compliance, agency adoption, or
  final-root claims.
- Phase 39 is CAL-ITP-style readiness workflow productization. It created no
  external evidence, did not create final-root proof, did not change consumer
  statuses, and does not strengthen compliance, consumer acceptance, agency
  adoption, hosted SaaS, production-readiness, vendor-compatibility, or
  production ETA-quality claims.
- Phase 42 is reference deployment diagnostic tooling. It created no external
  evidence, did not create final-root proof, did not change consumer statuses,
  and does not strengthen compliance, consumer acceptance, agency adoption,
  hosted SaaS, production-readiness, vendor-compatibility, or production
  ETA-quality claims.
- Phase 40 is docs/navigation guided operator trial productization. It created
  no external evidence, did not create final-root proof, did not change
  consumer statuses, and does not strengthen compliance, consumer acceptance,
  agency approval/adoption, hosted SaaS, production-readiness,
  vendor-compatibility, or production ETA-quality claims.
- Phase 41 is local/reference diagnostic tooling. It created no external
  evidence, did not create final-root proof, did not change consumer statuses,
  and does not strengthen compliance, consumer acceptance, agency
  approval/adoption, hosted SaaS, production-readiness, vendor-compatibility,
  or production ETA-quality claims.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- Phase 29A is adapter evaluation evidence only, not production ETA proof.
- Phase 29B is synthetic dry-run transform evidence only, not real vendor compatibility proof, production integration evidence, or AVL reliability evidence.
- Phase 31 is a docs-only pilot package and does not prove agency adoption.
- Phase 32 is draft launch materials only and does not prove launch, adoption, acceptance, compliance, endorsement, or production readiness.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, real-world ETA accuracy, production-grade ETA quality, certified hardware support, vendor compatibility, production AVL reliability, agency adoption, consumer submission, or public launch completion.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/latest.md`
4. `docs/handoffs/phase-34.md`
5. `docs/handoffs/phase-35.md`
6. `docs/handoffs/phase-36.md`
7. `docs/deployment/README.md`
8. `docs/deployment/oci-reference-deployment.md`
9. `docs/deployment/oci-reference-env.example`
10. `docs/deployment/oci-reference-smoke-checklist.md`
11. `docs/future-roadmap-post-outcome-c.md`
12. `docs/master-plan-self-hosted-agency-reuse.md`
13. `docs/phase-36-oci-reference-deployment-productization.md`
14. `docs/phase-37-agency-reusable-onboarding-flow.md`
15. `docs/handoffs/phase-37.md`
16. `docs/phase-38-integration-adapter-kit.md`
17. `docs/handoffs/phase-38.md`
18. `docs/phase-39-calitp-readiness-workflow.md`
19. `docs/handoffs/phase-39.md`
20. `docs/phase-40-guided-self-hosted-operator-trial.md`
21. `docs/tutorials/self-hosted-operator-trial.md`
22. `docs/handoffs/phase-40.md`
23. `docs/phase-41-operator-smoke-support-bundle.md`
24. `docs/tutorials/operator-smoke-and-support-bundle.md`
25. `docs/handoffs/phase-41.md`
26. `docs/deployment/reference-deployment-doctor.md`
27. `docs/phase-42-reference-deployment-doctor.md`
28. `docs/handoffs/phase-42.md`
29. `docs/integration-adapter-kit.md`
30. `docs/tutorials/reusable-agency-onboarding.md`
31. `docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md`
32. `docs/phase-33-public-gtfs-local-pilot-evidence.md`
33. `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md`
34. `docs/final-root-operator-request.md`
35. `docs/tutorials/public-gtfs-local-pilot.md`
36. `docs/runbooks/small-agency-pilot-operations.md`
37. `docs/roadmap-status.md`
38. `docs/california-readiness-summary.md`
39. `docs/compliance-evidence-checklist.md`
40. `docs/agency-owned-domain-readiness.md`
41. `docs/evidence/consumer-submissions/status.json`
42. `docs/evidence/consumer-submissions/submission-workflow.md`
43. `docs/evidence/redaction-policy.md`
44. `SECURITY.md`
45. `README.md`
46. `docs/dependencies.md`
47. `docs/decisions.md`
48. `docs/backlog.md`
49. `docs/open-questions.md`

## Current Objective

Make Open Transit RT easier to self-host, adapt, and integrate for small
agencies and civic technologists. Phase 60 is complete for final claim review
and public closeout without weakening the evidence boundaries.

External-proof tracks remain available later when a future operator is
authorized and retained claim-specific artifacts exist. Consumer or aggregator
submission work remains available only when a future operator is authorized, a
target is selected, official target paths are verified, and target-originated
evidence can be retained and redacted. Product improvements, validator success,
pilot packaging, prepared packets, and draft launch materials alone must not
advance target statuses.

For future public/status wording changes, run `make audit-final-claim-review`
and preserve the Phase 60 claim-to-evidence boundary unless retained
claim-specific evidence supports a narrower update.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

Run these when future work touches relevant surfaces:

```bash
make realtime-quality
make smoke
make demo-agency-flow
make test-integration
docker compose -f deploy/docker-compose.yml config
```

## Checks Run For Phase 32

- Pre-implementation `make validate` — passed.
- Pre-implementation `make test` — passed.
- Pre-implementation `git diff --check` — passed.
- Post-edit lightweight internal Markdown link/path check — passed.
- Post-edit consumer tracker status check — passed; all seven targets remain `prepared`.
- Post-edit targeted public-messaging scan — reviewed; matches are negative/boundary wording, current truth-state language, or required claim-to-evidence/checklist wording.
- Post-edit targeted secret/private-data scan — reviewed; no committed private artifacts found.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit final `git diff --check` — passed.
- Post-edit `make demo-agency-flow` — blocked during Docker image pull for the pinned GTFS-RT validator and was interrupted after no progress for several minutes. See `docs/handoffs/phase-32.md`.

## Checks Run For Phase 33

- Planning-pass baseline reportedly passed before implementation:
  `make validate`, `make test`, and `git diff --check`.
- Phase 33 attempted-run `make agency-app-up` — passed; local app started and
  imported the repo sample feed before the LA Metro attempt.
- Phase 33 attempted-run LA Metro source download — passed; raw ZIP kept in
  ignored `.cache/` storage.
- Phase 33 attempted-run LA Metro import with `demo-agency` — blocked by
  expected agency ID mismatch.
- Phase 33 attempted-run LA Metro import with local `LACMTA` setup — blocked by
  importer context timeout while inserting `stop_times.txt`.
- Phase 33 attempted-run `make agency-app-down` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `git diff --check` — passed.

## Post-Phase-33 Import Fix Results

- Added configurable `cmd/gtfs-import` timeout through `-timeout` and
  `GTFS_IMPORT_TIMEOUT`; default is now 15 minutes, and `0` disables the import
  timeout.
- Replaced per-row publish inserts for `gtfs_stop_time` and
  `gtfs_shape_point` with `pgx.CopyFrom`; shape point geometry is populated
  after bulk load.
- Changed publish-failure report recording to use a fresh short context, so a
  canceled import context does not leave the `gtfs_import` row stuck at
  `started`.
- Added focused CLI timeout tests and DB-backed timeout failure-report tests.
- Post-fix focused `go test ./cmd/gtfs-import ./internal/gtfs` — passed.
- Post-fix `make validate` — passed.
- Post-fix `make test` — passed.
- Post-fix `git diff --check` — passed.
- Post-fix first `make test-integration` attempt — blocked because Postgres was
  not ready immediately after `make db-up`; rerun after `pg_isready` passed.
- Post-fix `make test-integration` rerun — passed.
- Post-fix LA Metro import verification — passed locally; `gtfs-import-26`
  published with 114 routes, 11,884 stops, 33,642 trips, 2,105,503 stop_times,
  and 343,530 shape points.
- Post-fix `make db-down` — passed.

## Phase 33 Outcome C Retry Results

- Outcome C retry LA Metro source download — passed; raw ZIP kept in ignored
  `.cache/` storage.
- Outcome C retry source ZIP SHA-256 —
  `ce984bb5cc179d814fb0348878a6f7bd9ab6c940aaaec9fd4e97420583a0aa94`.
- Outcome C retry isolated local database migration — passed.
- Outcome C retry local `LACMTA` setup — passed for local evidence only.
- Outcome C retry import — passed; local feed version `gtfs-import-1`
  published with 114 routes, 11,884 stops, 33,642 trips, 2,105,503 stop_times,
  and 343,530 shape points.
- Outcome C retry local services and public proxy — passed at
  `http://localhost:19080`.
- Outcome C retry five-path public fetch — passed for `/public/feeds.json`,
  `/public/gtfs/schedule.zip`, Vehicle Positions, Trip Updates, and Alerts.
- Outcome C retry fetched schedule proof — passed; fetched schedule was
  verified as `LACMTA` public GTFS, not the repo sample feed.
- Outcome C retry static GTFS validator — attempted but failed to execute
  because Java runtime was unavailable in this local environment.
- Post-Phase-34 static GTFS validator retry — executed with Homebrew Java 17
  against the already-fetched local schedule ZIP; process exit code `0`, system
  error count `0`, and 3 warning notices:
  `expired_calendar`, `route_short_name_too_long`, and `unused_shape`. This is
  not a validator-clean or no-warning compliance claim.
- Outcome C retry Vehicle Positions, Trip Updates, and Alerts GTFS-RT
  validators — passed with 0 errors, 0 warnings, and 0 info notices against
  empty valid protobuf feeds.
- Outcome C retry telemetry dry-run — passed as dry-run-only synthetic payload
  display; no telemetry was sent.
- Outcome C retry admin/private boundary check — passed for recorded local
  checks.
- Outcome C retry packet update — added retained public-safe summaries under
  `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/`.
- Final post-Outcome-C-docs `make validate` — passed.
- Final post-Outcome-C-docs `make test` — passed.
- Final post-Outcome-C-docs `git diff --check` — passed.

## Checks Run For Phase 34

- Initial `make validate` — blocked because the pinned GTFS-RT validator image
  was not installed locally.
- Initial `docker info` — blocked because the Docker client could not connect
  to the Docker daemon at `unix:///Users/edwintse/.docker/run/docker.sock`.
- Retry `docker info` — passed after Docker became reachable.
- Retry `make validators-install` — passed.
- Retry `make validators-check` — passed.
- Retry `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Consumer tracker status check — passed; all seven targets remain `prepared`.
- Targeted wording scan — reviewed; matches are negated/boundary wording,
  historical phase names, or allowed evidence-boundary language. See
  `docs/handoffs/phase-34.md` for terms.
- Direct `/usr/bin/java` probe — still blocked by the macOS shim, but
  Homebrew Java 17 was available at `/usr/local/opt/openjdk@17/bin/java`.
- Static GTFS validator retry — executed against the Phase 33 fetched schedule
  ZIP in ignored `.cache` storage; process exit code `0`, system error count
  `0`, and 3 warning notices. No validator-clean or no-warning claim was added.

## Checks Run For Phase 35

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only consumer tracker status check — passed; 7 targets found and all
  remain `prepared`.

## Checks Run For Phase 36

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only consumer tracker status check — passed; 7 targets found and all
  remain `prepared`.

## Checks Run For Phase 37

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- Read-only consumer tracker status check — passed; 7 targets found and all
  remain `prepared`.
- Optional local fixture smoke — passed using a generated ZIP from
  `testdata/gtfs/valid-small` served from ignored `.cache/` storage, with
  validators skipped.
- Post-completion hardening re-ran the required checks and preserved all seven
  consumer statuses as `prepared`.

## Checks Run For Phase 40

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Read-only consumer tracker status check — passed; 7 targets found and all
  remain `prepared`.

## Checks Run For Phase 41

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed after fixing one trailing-space issue.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Exact seven-target consumer status check — passed; Google Maps, Apple Maps,
  Transit App, Bing Maps, Moovit, Mobility Database, and transit.land all
  remain `prepared`.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` —
  passed; the tracker file is unchanged.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make smoke` — passed.
- `make support-bundle` — passed with no local/reference app running; the
  bundle recorded public checks as unavailable and passed redaction scanning.
- `make operator-smoke SKIP_VALIDATORS=true` — blocked because no
  local/reference app was running at `http://localhost:8080`; direct
  `/public/feeds.json` probe returned curl connection failure / HTTP `000`.

## Checks Run For Phase 42

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null` — passed.
- Exact seven-target consumer tracker check — passed; Google Maps, Apple Maps,
  Transit App, Bing Maps, Moovit, Mobility Database, and transit.land all
  remain `prepared`.
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json` —
  passed; the tracker file is unchanged.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `make deployment-doctor` — passed. With no local/reference deployment env or
  running app exported in this checkout, the doctor exited `0` as designed and
  reported 44 blockers, 0 warnings, and 12 unavailable checks. Public feed,
  admin boundary, database, migration, and PostGIS checks were skipped because
  no deployment values were supplied; validator tooling passed; backup and
  restore readiness reported blockers because deployment-owned paths/targets
  were not supplied.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 29B synthetic fixtures are not real vendor AVL evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Phase 30 did not select a target, verify an official path, submit a packet, add artifacts, or change consumer statuses.
- Phase 31 did not add real pilot evidence, consumer evidence, agency adoption evidence, operations evidence, final-root proof, or device/AVL proof.
- Phase 32 did not post announcements, contact agencies, contact consumers, launch publicly, add evidence artifacts, or change consumer statuses.
- The post-Phase-32 final-root evidence follow-up did not create a final-root packet, run hosted packet audit, or refresh prepared packet references.
- Phase 33 created an Outcome C public-GTFS local/pilot evidence packet only.
  It does not support final-root, consumer, compliance, production, or real
  realtime/ETA-quality claims.
- Phase 36 created deployment reference documentation only. It did not create
  external evidence, run a deployment, change consumer statuses, or strengthen
  final-root, compliance, production-readiness, or consumer claims.
- Phase 37 created reusable local/reference onboarding tooling only. It did not
  create external evidence, run a final-root proof, change consumer statuses,
  or strengthen agency, consumer, compliance, production-readiness, vendor, or
  ETA-quality claims.
- Phase 39 created an authenticated readiness workflow only. It did not create
  external evidence, run final-root proof, change consumer statuses, contact
  external parties, or strengthen agency, consumer, compliance,
  production-readiness, vendor, or ETA-quality claims.
- Phase 40 created guided trial docs/navigation only. It did not create
  external evidence, run final-root proof, change consumer statuses, contact
  external parties, or strengthen agency, consumer, compliance,
  production-readiness, vendor, or ETA-quality claims.
- Phase 53 did not select a consumer or aggregator target, verify an official
  path, contact a target, automate a portal, make a submission, add an
  artifact, or change consumer statuses because no operator authorization,
  official path verification, or target-originated artifact exists locally.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit Next

Use the next approved phase document and `docs/handoffs/latest.md` before
choosing files. Phase 60 is closed; do not
reopen Phase 39, Phase 40, Phase 41, Phase 42, Phase 43, Phase 44, Phase 45,
Phase 46, Phase 47, Phase 48, Phase 49, Phase 50, Phase 51, Phase 52, or Phase
53 unless a concrete readiness-workflow, guided-trial, diagnostics,
support-bundle, deployment-doctor, local-routing, checklist,
telemetry-simulator, GTFS quality triage, validator-health,
operations-notification, AVL send, predictor adapter, backtesting,
operations-reliability, final-root workflow, or authorized consumer-submission
boundary regression is found. Do not reopen Phase 56 unless a concrete
tenant-route, proxy exposure, backup/restore/export, evidence-boundary, or
hosted-claim regression is found. Do not reopen Phase 57 unless a concrete
release package, checksum, SBOM/provenance, image metadata, audit, or
hosted-service claim regression is found. Do not reopen Phase 58 unless a
concrete BYOD template, support-boundary, procurement-template, audit, or
marketplace/vendor/support claim regression is found.
Do not reopen Phase 59 unless a maintainer supplies a real retained
authorization plus public-safe kickoff, operations, feedback, and decision
artifact set for a named pilot.
Do not reopen Phase 60 unless a concrete public/status claim-boundary
regression is found, official requirements have materially changed, or retained
claim-specific evidence supports a narrower claim update.

External-proof docs remain available for later optional tracks. Do not edit
target-specific consumer records, `docs/evidence/consumer-submissions/status.json`,
final-root evidence packets, or artifact directories unless retained, redacted,
target-originated evidence supports a target-specific status transition.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs unless explicitly approved and backed by current target documentation.
- Do not automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, production-grade ETA proven, real-world ETA-accuracy proven, certified hardware supported, vendor-compatible, agency-adopted, consumer-accepted, or publicly launched.

## Exact Next-Step Recommendation

Proceed with maintainer review or a separately scoped product phase. Any
future evidence intake is optional and requires explicit written authorization
first. Phases 0 through 60 remain closed, Phases 61 through 67 are complete,
Phase 68+ is closed blocker-only / authorization-gated for the current
no-authorization review, and Phase 69 is complete for UI-first product
acceptance. The default next work is not evidence collection. For
public/status wording work, start with `make audit-final-claim-review` and
`make audit-product-acceptance`, and keep unsupported claims removed or
bounded.

Future final-root work should use `make collect-final-root-evidence` and
`make audit-final-root-evidence` only when a real final root and redacted
approval artifact exist. Without those artifacts, keep using blocker-only
closure and leave `docs/evidence/captured` unchanged.

Future consumer submission execution should proceed only when retained operator
authorization, official target path verification, and target-originated or
operator-retained submission evidence exist for one named target. Without those
artifacts, keep all seven targets at `prepared`.

Use the Phase 36 reference deployment docs, Phase 37 reusable onboarding flow,
Phase 38 adapter kit, Phase 39 readiness workflow, Phase 40 guided operator
trial, Phase 41 diagnostics, Phase 42 deployment doctor, Phase 43 private
operator checklist, Phase 44 telemetry simulator, Phase 45 GTFS quality triage,
Phase 46 validator health, Phase 47 operations notification summaries, Phase
50 realtime quality backtesting, Phase 51 operations reliability diagnostics,
Phase 52 final-root workflow tooling, Phase 56 multi-agency route/proxy
diagnostics, Phase 57 release package tooling, Phase 58 vendor-equivalent
templates, and Phase 59 blocker-only real-pilot closeout as the
self-hosted/integration baseline.
External-proof work remains a future optional path, not the default roadmap.
