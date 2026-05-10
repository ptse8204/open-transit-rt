# Roadmap To Cal-ITP Compliance And GTFS-RT Gap Closure

**Status:** Phase 60 closed for final claim review and public closeout
**Intended repo path:** `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
**Generated for review:** 2026-05-08
**Phase 54 refresh date:** 2026-05-09
**Current repo baseline:** Phases 0 through 60 are closed for their documented scopes. Phase 60 closed for final claim review and public closeout.

This roadmap plans the path from the current self-hosted agency-reuse prototype
to a fully evidence-backed, open-source GTFS / GTFS-Realtime operations stack
that can support California transit data compliance work and close the most
important practical gaps in existing open-source GTFS-RT tooling.

This document is intentionally evidence-bounded. It does not claim that Open
Transit RT is CAL-ITP/Caltrans compliant today, accepted by any consumer,
adopted or approved by an agency, available as hosted SaaS, production-ready
for all deployments, vendor-compatible, marketplace approved, agency-owned
final-root proven, or production-grade ETA quality proven.

## Official Requirements Verification Gate

This roadmap is based on the repository's current requirements and readiness
docs, especially:

- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/master-plan-self-hosted-agency-reuse.md`
- `docs/handoffs/latest.md`
- `docs/phase-54-official-requirements-refresh.md`

Phase 54 re-checked current official public sources on May 9, 2026. The
Caltrans California Transit Data Guidelines page was current as Version 4.0
with visible date December 11, 2024, and the FAQ was current as Version 4.0.
The refresh mapped current guidance to repo requirements only. It did not
create compliance evidence, consumer evidence, final-root evidence, or a
readiness claim.

Phase 55 added local `.cache` compliance/readiness packet generation and audit
guards. It did not create retained evidence, contact consumers, fetch live
feeds, change consumer statuses, prove final-root readiness, or claim
compliance.

Phase 59 closed blocker-only because no retained real pilot authorization,
kickoff, operations, feedback, and decision artifact set was available in the
repository. It did not create evidence, change consumer statuses, refresh
earlier pilot packets, or strengthen agency, consumer, compliance, operations,
hosted service, support, vendor, marketplace, or ETA-quality claims.

Phase 60 added a final claim-to-evidence review, unsupported-claim table, local
read-only audit helper, mutation-style script tests, Make targets, and handoff.
It did not create evidence, write `docs/evidence`, contact external parties,
change consumer statuses, or strengthen public launch, compliance, agency,
consumer, hosted service, SLA/uptime, production-readiness, vendor,
marketplace, or ETA-quality claims.

Before any public compliance claim or final compliance closeout, a maintainer
must re-check the official Caltrans / Cal-ITP and FTA sources again and update
this roadmap if the requirements have changed.

The repo should keep using this wording until retained evidence supports more:

```text
Open Transit RT supports CAL-ITP-style readiness workflows and implements
technical foundations for stable GTFS and GTFS-Realtime publication.
```

Do not use wording that says the project is compliant, accepted by trip
planners, or production-ready for all agencies.

## Current Baseline

As of Phase 60, Open Transit RT has:

- GTFS ZIP import and GTFS Studio draft/publish workflows;
- stable public feed paths for:
  - `/public/feeds.json`;
  - `/public/gtfs/schedule.zip`;
  - `/public/gtfsrt/vehicle_positions.pb`;
  - `/public/gtfsrt/trip_updates.pb`;
  - `/public/gtfsrt/alerts.pb`;
- authenticated telemetry ingest with device bearer tokens;
- conservative deterministic trip matching;
- GTFS-RT Vehicle Positions;
- GTFS-RT Trip Updates behind `internal/prediction.Adapter`;
- GTFS-RT Alerts;
- pinned MobilityData validator workflows;
- local app packaging;
- OCI/OCL-style reference deployment docs;
- reusable agency onboarding through `make agency-pilot-up`;
- an Integration Adapter Kit;
- a CAL-ITP-style Operations Console readiness page;
- a guided self-hosted operator trial;
- operator smoke checks and redaction-safe support bundles;
- a read-only OCI/OCL-style reference deployment doctor;
- a private authenticated setup/readiness checklist in HTML and JSON;
- a synthetic telemetry simulator that uses real device-token auth and
  `POST /v1/telemetry`;
- private validator-health diagnostics;
- private local operations notification drafts from existing diagnostics;
- synthetic AVL adapter fixtures;
- deterministic realtime-quality replay fixtures;
- prepared consumer packets for seven targets;
- hosted/operator OCI pilot evidence for the recorded pilot root;
- current official-source requirement mappings for stable public schedule and
  realtime URLs, all three realtime feed types, canonical no-error validation,
  open licensing, major trip-planner acceptance as a third-party requirement,
  provider website source-of-truth links, technical/feed contact, Transitland
  and Mobility Database availability, and API-key registration constraints.
- final claim review and local audit tooling that checks bounded public/status
  docs, prepared-only consumer tracker state, README-only consumer artifact
  directories, unsafe private strings, and required Phase 60 sections.

The repo still lacks retained proof for:

- agency-owned or agency-approved final public root;
- agency adoption or approval;
- final-root DNS/TLS/redirect/five-feed evidence;
- final-root validator-clean records;
- consumer submission, review, acceptance, listing, display, or ingestion;
- official CAL-ITP/Caltrans compliance;
- real authorized pilot closeout evidence and agency/operator feedback;
- real vendor AVL compatibility;
- real-world ETA quality;
- production-grade ETA quality;
- production multi-tenant hosting;
- hosted SaaS availability;
- paid support/SLA commitments.

Phases 54, 59, and 60 do not reduce any of these evidence gaps.

## End State

The long-term end state has three levels. They must not be collapsed into one
claim.

### Level 1 — Self-Hosted Technical Maturity

A small agency, operator, or civic technologist can:

1. install Open Transit RT from a release;
2. deploy the reference server pattern;
3. import or author GTFS;
4. publish stable schedule and GTFS-RT feeds;
5. validate all feeds;
6. connect telemetry through a documented adapter path;
7. monitor feed freshness and validation status;
8. operate backup/restore, rollback, and incident workflows;
9. produce a redaction-safe support bundle.

This level supports self-hosted reuse. It is not yet a compliance or consumer
acceptance claim.

### Level 2 — Evidence-Backed California Data Compliance Readiness

For a specific deployment and feed root, retained evidence exists for:

1. stable public HTTPS URLs;
2. public static GTFS;
3. all three GTFS-Realtime feeds;
4. current no-error validation results for schedule and realtime feeds;
5. agency-approved license/contact metadata;
6. public discoverability metadata;
7. current operations evidence;
8. consumer packet readiness;
9. official requirements mapping checked against current Caltrans / Cal-ITP
   guidance.

This level may support a stronger readiness claim, but still does not prove
consumer acceptance unless target-originated evidence exists.

### Level 3 — Full External Closeout

For the exact feed scope and URL root, retained redacted artifacts prove:

1. agency or operator authorization;
2. final root approval;
3. official consumer submission path verification;
4. target-specific submissions;
5. target-originated under-review, rejected, blocked, accepted, ingested,
   listed, or displayed statuses;
6. final claim review against current official requirements.

Only this level can support public language about consumer acceptance or final
compliance status for a named deployment.

## Plugin / Adapter Component Strategy

Open Transit RT should close open-source GTFS-RT gaps through adapter-bound
components, not through a closed all-in-one system.

The term `plugin` in this roadmap means a replaceable component with a stable
contract. It does not require Go dynamic plugins. Safer shapes include command
adapters, sidecars, HTTP adapters, fixtures, manifests, and deployment-owned
processes.

| Component area | Current repo status | Long-term target | Boundary rule |
| --- | --- | --- | --- |
| Static GTFS input | ZIP import and GTFS Studio exist | import, edit, validate, triage, publish, rollback | published GTFS stays separate from draft GTFS |
| Telemetry ingest | authenticated `/v1/telemetry` exists | device, simulator, and AVL adapter paths | vendor payloads transform before telemetry ingest |
| AVL/vendor adapters | synthetic dry-run adapter kit exists | private/authorized real adapters can be added outside core | no named vendor claim without retained evidence |
| Prediction engines | deterministic adapter exists; external predictor boundary evaluated | internal predictor, TheTransitClock-like external predictor, or later replacements | all predictors stay behind `internal/prediction.Adapter` |
| Validators | pinned static and realtime workflows exist | scheduled, strict, reportable validators with operator triage | validators verify outputs; they do not own core state |
| Monitoring | metrics and pilot ops helpers exist | Prometheus/Grafana/OpenTelemetry optional packs | observability must not break feed generation |
| Consumer workflows | prepared packets and tracker exist | target-specific workflows with retained artifacts | no target status change without target-originated evidence |
| Evidence/redaction | policies and packet patterns exist | support bundles and evidence packets with strict redaction | no secrets, raw telemetry, or private operator artifacts committed |
| Deployment | OCI/OCL-style reference docs exist | preflighted, release-based, rollback-safe deployment | public feed edge and private admin edge stay separate |

## Mapping To Repository Compliance Requirements

The repo's `docs/requirements-calitp-compliance.md` defines RQ-4A through
RQ-4G. The roadmap below closes those requirements in stages.

| Requirement | Current status | Roadmap closure path |
| --- | --- | --- |
| RQ-4A complete realtime feed set | all three GTFS-RT feed paths exist; Phase 54 confirms Caltrans completeness maps to Trip Updates, Vehicle Positions, and Service Alerts | Phases 44, 46, 50, and 51 add simulator, validation-health, backtesting, and reliability diagnostics that can support future real-deployment evidence when retained artifacts exist. |
| RQ-4B stable public production URLs | stable paths exist; final root missing; Phase 54 confirms stable public schedule and realtime URLs remain official-source requirements | Phases 42 and 52 create preflight and final-root evidence workflow |
| RQ-4C validator-clean feeds | tooling exists; final no-error production evidence missing; Phase 54 confirms regular canonical validation with no errors remains required for compliant wording | Phases 45, 46, and 55 add triage, scheduled validators, and compliance packet generation |
| RQ-4D open license and discoverability | metadata model exists; agency-approved metadata and source-of-truth website proof are missing; Phase 54 confirms provider/regional website links, technical contact, Transitland, and Mobility Database availability mappings | Phases 43, 52, and 55 harden metadata workflow and final-root public discovery proof |
| RQ-4E consumer ingestion workflow | prepared packets exist; no submissions or acceptance; Phase 54 confirms major trip-planner acceptance is a separate third-party requirement | Phases 53 and 55 handle authorized target submission and artifact-backed statuses |
| RQ-4F marketplace/vendor equivalence | explicitly not claimed | Phase 58 creates optional service/procurement pack without turning repo into SaaS |
| RQ-4G dashboard and scorecard | readiness page exists | Phases 43, 51, and 55 improve operator remediation and exportable compliance reporting |

## Master Phase Roadmap

The following phases are closed for their documented scopes. Completed local,
diagnostic, template, or blocker-only phases do not imply retained external
evidence, compliance, consumer acceptance, agency adoption, hosted SaaS,
production readiness, vendor compatibility, marketplace approval, final-root
proof, or production-grade ETA quality.

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 41 | Operator Smoke And Support Bundle | Complete. Repeatable local/reference smoke helper and redaction-safe support bundle. |
| 42 | Reference Deployment Doctor | Complete. Read-only checks for env, services, reverse proxy posture, validators, DB, backups, restore readiness, and route boundaries. |
| 43 | Operator UX Setup V2 | Complete. Private grouped checklist for setup/readiness remediation and JSON export. |
| 44 | Telemetry Simulator And Device Trial | Complete. Safe simulator path that sends synthetic telemetry through real device-token ingest. |
| 45 | GTFS Quality Triage Loop | Complete. Private authenticated static GTFS quality triage and admin-only rerun diagnostics without evidence, auto-edit, or compliance/consumer claims. |
| 46 | Validator Automation And Health Gates | Complete. Private validator-health diagnostics and local script/export workflow without evidence, publish blocking, or compliance/production claims. |
| 47 | Self-Hosted Operations Notifications | Complete. Private local notification drafts from existing diagnostic summaries only; no sending or evidence creation. |
| 48 | AVL Adapter Runtime Path | Complete. Private authenticated AVL adapter send mode through `/v1/telemetry` with env-only credentials and redacted `.cache` diagnostics. |
| 49 | External Predictor Runtime Adapter | Complete. Optional generic external HTTP predictor adapter behind `internal/prediction.Adapter` without named predictor or ETA-quality claim. |
| 50 | Realtime Quality Backtesting | Complete. Private local observed-arrival quality backtesting CLI/library with aggregate diagnostics only. |
| 51 | Operations Reliability And SLO Readiness | Complete. Private reliability diagnostics and Vehicle Positions health persistence without SLA/uptime or production-readiness claims. |
| 52 | Final Public Root Evidence Workflow | Complete blocker-only. Guarded final-root collector/audit workflow exists; no real final-root or approval evidence was retained. |
| 53 | Authorized Consumer Submission Execution | Complete blocker-only. No authorization, official target path verification, or target-originated artifact existed; all targets remain `prepared`. |
| 54 | Official Requirements Refresh | Complete. Re-checked current official public sources and updated requirement mapping without creating evidence or compliance claims. |
| 55 | Compliance Evidence Packet Generator | Complete. Generates ignored `.cache` blocker/draft readiness packets with strict claim-gated audit. |
| 56 | Multi-Agency Hosting Hardening | Complete. Tenant-safe public feed routing, proxy checks, private diagnostics, and backup/restore/export/evidence boundaries without hosted SaaS or production multi-tenant claims. |
| 57 | Release Packaging And Supply Chain | Complete. Local `.cache` source packages with checksums, SBOM/provenance metadata, optional local image metadata, and audit checks; no hosted service or production image claim. |
| 58 | Optional Marketplace / Vendor-Equivalent Pack | Complete. Template-only BYOD/hardware, support, SLA/KPI, implementation, and procurement pack with local claim-boundary audit. |
| 59 | Real Pilot Closeout | Complete blocker-only. No retained real pilot authorization, feedback, operations, or decision artifact set was available; no evidence or claim status changed. |
| 60 | Final Claim Review And Public Closeout | Complete. Final claim-to-evidence review and local audit tooling; no evidence or claim status changed. |

## Phase 41 — Operator Smoke And Support Bundle

### Goal

Turn the Phase 40 guided operator trial into a repeatable smoke workflow and
redaction-safe support bundle.

### Why This Comes First

The repo already has a guided trial, but operators still need a dependable way
to run the same checks and share safe diagnostics without leaking secrets or
mistaking trial output for evidence.

### Artifacts

Add:

- `scripts/operator-smoke.sh`
- `scripts/support-bundle.sh`
- `docs/tutorials/operator-smoke-and-support-bundle.md`
- `docs/phase-41-operator-smoke-support-bundle.md`
- `docs/handoffs/phase-41.md` after implementation

Update:

- `Makefile`
- `README.md`
- `docs/README.md`
- `docs/tutorials/README.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/handoffs/latest.md`

### Required Behavior

The smoke helper should:

- run in local/reference mode;
- accept `PUBLIC_BASE_URL`, optional `ADMIN_BASE_URL`, and optional
  `ADMIN_TOKEN`;
- fetch the five public feed paths;
- check pinned validator tooling status;
- run validator API calls when configured and authenticated;
- check `/admin/operations/readiness` only through the private/admin boundary;
- run the synthetic AVL adapter dry-run;
- write local ignored output under `.cache/operator-smoke/<timestamp>/`;
- print a redaction-safe summary;
- state that no evidence packet was created.

The support bundle helper should:

- collect selected public feed fetch metadata, not raw private data by default;
- collect env key presence, not env values;
- collect validator status and command versions;
- collect service health summaries;
- collect docs/status pointers;
- exclude secrets, tokens, raw telemetry, private vendor payloads, ACME material,
  private keys, DB passwords, cookies, JWTs, CSRF secrets, and unredacted logs;
- write to ignored `.cache/support-bundles/<timestamp>/`;
- include a manifest explaining what was included and excluded.

### Definition Of Done

- `make operator-smoke` exists.
- `make support-bundle` exists.
- Both commands support dry-run or safe local output behavior.
- Output is ignored by git.
- Docs explain when output can and cannot become evidence.
- Existing consumer statuses remain unchanged.
- No compliance, production, consumer, agency, vendor, or ETA-quality claim is added.

### Checks

Run:

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Run when relevant:

```bash
docker compose -f deploy/docker-compose.yml config
make smoke
```

## Phase 42 — Reference Deployment Doctor

### Goal

Add a deployment preflight / doctor workflow for the OCI/OCL-style reference
server path.

### Artifacts

- `scripts/deployment-doctor.sh`
- `make deployment-doctor`
- `docs/deployment/reference-deployment-doctor.md`
- phase doc and handoff

### Required Checks

The doctor should check:

- required env var presence without printing secret values;
- generated-secret presence and minimum length where applicable;
- DB connectivity and migration status;
- PostGIS availability;
- service health endpoints;
- public feed edge anonymous fetches;
- private/admin route boundary behavior;
- reverse proxy HTTPS/redirect configuration when running on a deployed host;
- pinned validator tooling status;
- backup target presence;
- restore-drill readiness;
- release version and git SHA reporting.

### Definition Of Done

An operator can run one command and learn what prevents a safe reference
deployment from being used for a pilot.

### Non-Goals

No final-root evidence is created. No deployment is called production-ready from
preflight success alone.

### Phase 42 Closeout

Phase 42 is complete for the read-only diagnostic tooling scope. The doctor
does not source private env files, does not run migrations/backups/restores,
does not write to `EVIDENCE_OUTPUT_DIR`, and does not create evidence or
change consumer statuses. Default mode reports blockers while exiting `0`;
`STRICT_DOCTOR=true` is the mode that fails on blockers.

## Phase 43 — Operator UX Setup V2

### Goal

Make the Operations Console useful for non-expert operators completing setup
and readiness tasks.

### Artifacts

- enhanced `/admin/operations/setup`;
- enhanced `/admin/operations/readiness`;
- operator checklist export;
- docs explaining row-level remediation.

### Required Improvements

- group checklist rows into setup, feeds, validation, telemetry, operations,
  and consumer workflow;
- show neutral status values with plain next actions;
- link validation errors/warnings to triage docs;
- show whether metadata is missing, placeholder-like, or operator-entered with
  approval unknown;
- show whether public URLs are local, pilot/reference, final-root candidate
  unverified, or missing final-root evidence;
- export a private operator checklist.

### Definition Of Done

A non-expert operator can tell what is missing without reading phase handoffs.

### Non-Goals

Do not expose admin routes publicly. Do not claim compliance from UI state.

### Phase 43 Closeout

Phase 43 is complete for the private authenticated checklist scope. The
Operations Console now exposes `/admin/operations/checklist` and
`/admin/operations/checklist.json`, both backed by one deterministic model with
fixed group order, stable row IDs, neutral statuses, repo-relative docs links,
heuristic labels, and explicit false claim flags. Dashboard, setup, and
readiness pages link to both checklist routes. The phase also patched local
routing so exact `/` remains a `200` convenience page while unmatched local
paths return `404`, and the deployment doctor now checks `/admin/gtfs-studio`
instead of exact `/admin/gtfs`.

Phase 43 created no external evidence, did not create final-root proof, did
not change consumer statuses, and added no compliance, consumer acceptance,
agency approval/adoption, hosted SaaS, production-readiness,
vendor-compatibility, or production-grade ETA claim.

## Phase 44 — Telemetry Simulator And Device Trial

### Goal

Create a safe, repeatable way to test the real telemetry ingest path without
real hardware or vendor data.

### Artifacts

- `scripts/telemetry-simulator.sh` or Go equivalent;
- fixture scenarios for on-route, stale, out-of-order, unknown, block
  transition, and after-midnight cases;
- tutorial for simulator use with device tokens;
- admin UI pointers to simulator results.

### Required Behavior

- use real device bearer-token authentication;
- post to `/v1/telemetry` rather than bypassing ingest;
- avoid raw private telemetry;
- support local/reference deployments;
- label output as synthetic.

### Definition Of Done

An operator can run synthetic telemetry through the real ingest path and see
Vehicle Positions / matching behavior update.

### Non-Goals

No real AVL reliability, vendor compatibility, or ETA-quality claim.

### Phase 44 Closeout

Phase 44 is complete for the synthetic-only local/reference simulator scope.
The repo now has `cmd/telemetry-simulator`, `scripts/telemetry-simulator.sh`,
`make telemetry-simulator`, fixtures under `testdata/telemetry-simulator/`,
optional DB-backed matcher and Vehicle Positions debug diagnostics after
accepted HTTP ingest, and operator docs. It posts only to `/v1/telemetry` with
real device bearer-token auth.

Phase 44 created no evidence packet, changed no consumer statuses, added no
real vendor payloads or private telemetry, and added no vendor compatibility,
production AVL reliability, real realtime data, production-grade ETA,
CAL-ITP/Caltrans compliance, agency approval/adoption, hosted SaaS, or
production-readiness claim.

## Phase 45 — GTFS Quality Triage Loop

### Goal

Close the gap between validator output and operator action.

### Artifacts

- validator notice taxonomy docs;
- import/Studio triage guide;
- UI links from validation records to recommended fixes;
- examples for common warnings such as expired calendars, long route names,
  unused shapes, missing references, bad times, and duplicate IDs.

### Required Behavior

- distinguish blocking errors, warnings, and informational notices;
- keep canonical validator output separate from internal import validation;
- document when warnings can remain and when they block a claim;
- support re-run after fixes.

### Definition Of Done

An operator can move from validator report to fix path without needing a GTFS
expert for every notice.

### Phase 45 Closure

Implemented as private authenticated Operations Console triage at
`/admin/operations/gtfs-quality`. The route separates canonical MobilityData
static validator output from Open Transit RT internal import validation, caps
groups/samples, supports admin-only rerun against the authenticated agency
active schedule through the server-side static validator mapping, and keeps raw
reports/stdout/stderr/argv/private paths out of rendered HTML.

This is diagnostics only. It does not create evidence packets, write to
`docs/evidence`, auto-edit GTFS, or claim consumer acceptance or
CAL-ITP/Caltrans compliance.

### Non-Goals

Do not auto-edit agency GTFS without operator approval.

## Phase 46 — Validator Automation And Health Gates

### Goal

Make validation operationally meaningful for local/reference operators without
turning diagnostics into compliance or evidence claims.

### Artifacts

- private Operations Console validator-health HTML and JSON routes;
- `scripts/validator-health.sh` and `make validator-health`;
- deployment-doctor GET-only summary integration;
- fixed four-feed health status export;
- strict local/reference health gates through script exit behavior.

### Required Behavior

- static health uses only `schedule` plus `static-mobilitydata`;
- realtime health uses `vehicle_positions`, `trip_updates`, and `alerts` plus
  `realtime-mobilitydata`;
- internal GTFS import validation remains context only;
- warnings, failures, stale results, missing artifacts, and tooling blockers
  produce safe operator next actions;
- browser requests cannot supply commands, paths, URLs, artifacts, argv, args,
  or timeouts.

### Definition Of Done

A deployment can show current validation state for schedule, Vehicle Positions,
Trip Updates, and Alerts.

### Phase 46 Closure

Phase 46 is complete for private local/reference diagnostics. It does not add
automatic publish blocking, evidence creation, consumer submission, consumer
status mutation, or compliance/production-readiness claims.

### Non-Goals

Validator success still does not mean consumer acceptance, compliance, or
production readiness.

## Phase 47 — Self-Hosted Operations Notifications

### Goal

Provide a private local notification draft that summarizes existing
validator-health and deployment-doctor diagnostics for a self-hosted operator.
The output is a bounded local summary only.

### Artifacts

- `scripts/operations-notify.sh`;
- `make operations-notify`;
- `docs/tutorials/self-hosted-operations-notifications.md`;
- `docs/phase-47-self-hosted-operations-notifications.md`;
- `docs/handoffs/phase-47.md`.

### Required Behavior

- read the latest timestamp-named `.cache/validator-health/*/summary.json`;
- read the latest timestamp-named `.cache/deployment-doctor/*/summary.json`;
- support explicit source summary paths under `.cache`;
- reject symlink and evidence-like output or source paths;
- cap source sizes, output sizes, next actions, and copied fields;
- record webhook/email destination presence as booleans only;
- write `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, and
  `notification.txt`;
- keep `notification.txt` marked `DRAFT — NOT SENT`;
- support a strict mode for local automation failure semantics.

### Definition Of Done

A self-hosted operator can generate a private local draft that summarizes
current validation-health, feed availability, deployment-doctor blockers, and
safe next actions without sending anything or copying private source data.

### Non-Goals

No notification sending, webhook delivery, email delivery, public API, admin
route call, validator run, database requirement, Docker requirement, app
requirement, evidence creation, consumer contact, consumer status change, GTFS
auto-editing, publish blocking, CAL-ITP/Caltrans compliance proof, consumer
acceptance proof, agency adoption proof, hosted SaaS claim, production
readiness proof, vendor compatibility proof, or production-grade ETA proof.

Phase 47 is not a compliance gate, not production health proof, not evidence,
and not consumer-readiness proof.

## Phase 48 — AVL Adapter Runtime Path

### Goal

Move from synthetic dry-run adapter examples to a safe pattern for private or
authorized real AVL adapters.

### Artifacts

- send-mode reference design;
- private credential handling guide;
- adapter manifest schema;
- deployment-owned adapter process examples;
- redaction/evidence rules for real AVL payloads;
- conformance tests that do not include private vendor data.

### Required Behavior

- preserve `/v1/telemetry` as the core ingest boundary;
- keep vendor credentials out of the public repo;
- support mapping authority outside vendor payloads;
- fail closed on unknown devices, wrong agency, stale/future timestamps, and
  malformed coordinates;
- document retry/backoff and idempotency expectations.

### Definition Of Done

A maintainer can add or privately run a real adapter without contaminating core
matching, feed, or device-token code with vendor assumptions.

### Phase 48 Closure

Phase 48 is complete for the private authenticated send-mode scope. It keeps
`/v1/telemetry` as the runtime ingest boundary, uses env-only credential
references, writes redacted private diagnostics under `.cache`, and creates no
evidence packet, consumer-status change, named vendor support claim,
production AVL reliability claim, or production-readiness/compliance claim.

### Non-Goals

No named vendor support or certification unless retained authorization and test
evidence exist.

## Phase 49 — External Predictor Runtime Adapter

### Goal

Add an optional runtime path for external predictors behind
`internal/prediction.Adapter`.

### Artifacts

- predictor adapter config contract;
- process/HTTP adapter implementation or reference sidecar pattern;
- health and timeout behavior;
- output validation;
- shadow-mode comparison docs;
- dependency/license review checklist.

### Required Behavior

- Vehicle Positions continue if predictor is unavailable;
- Trip Updates fail closed or degrade clearly;
- wrong-agency/wrong-feed/stale/low-confidence outputs are rejected;
- diagnostics show adapter health and fallback behavior;
- no external predictor internals leak into telemetry, matching, GTFS import, or
  Vehicle Positions.

### Definition Of Done

A deployment can test an external predictor in shadow or controlled runtime mode
without weakening core source-of-truth boundaries.

### Phase 49 Closure

Phase 49 is complete for the optional disabled-by-default generic HTTP
predictor adapter boundary. Vehicle Positions remain independent, Trip Updates
fall back safely on adapter failure, and no named predictor, vendor
compatibility, consumer, evidence, compliance, or production-grade ETA claim
was added.

### Non-Goals

No production-grade ETA claim from adapter wiring alone.

## Phase 50 — Realtime Quality Backtesting

### Goal

Create the evidence path for ETA and realtime-quality maturity.

### Artifacts

- observed arrival/departure comparison schema;
- backtesting CLI or job;
- route/time-period metrics;
- coverage and withheld-case reports;
- shadow predictor comparison reports;
- Operations Console quality view updates.

### Required Metrics

- eligible prediction candidates;
- Trip Updates coverage;
- future stop coverage;
- unknown/ambiguous/stale rates;
- mean/median absolute error where observed arrivals exist;
- percentile error by route, stop, time of day, and service pattern;
- cancellation/alert linkage;
- fallback rate for external predictors.

### Definition Of Done

The repo can measure realtime quality on real or approved observed-operation
data without overclaiming.

### Phase 50 Closure

Phase 50 is complete for private local diagnostics. It added a backtesting
library/CLI and `.cache` aggregate outputs for overall, route, time-period, and
route/time-period metrics. It did not add DB persistence, public APIs,
Operations Console views, evidence writes, publish gates, consumer status
changes, production ETA proof, or readiness claims.

### Non-Goals

No production-grade ETA wording until route/time-period evidence meets a
maintainer-approved threshold.

## Phase 51 — Operations Reliability And SLO Readiness

### Goal

Make long-running operations visible and auditable.

### Artifacts

- feed freshness SLO template;
- incident record workflow;
- alert delivery proof pattern;
- backup/restore drill cadence;
- validator failure response workflow;
- key rotation records;
- operations status export.

### Required Evidence For Stronger Claims

- repeated public fetch success;
- validation cadence;
- restore drills;
- incident/response records;
- secret rotation records;
- alert delivery proof;
- operator handoff notes.

### Definition Of Done

Operators can run Open Transit RT over time and show what happened when feeds,
validators, telemetry, or deployments degraded.

### Phase 51 Closure

Phase 51 is complete for private operations reliability diagnostics. It added
authenticated reliability summaries, bounded Vehicle Positions health
persistence, and `.cache` reliability output only. It did not add public
routes, monitoring-stack dependencies, evidence writes, consumer status
changes, SLA/uptime claims, hosted SaaS claims, production-readiness claims, or
compliance claims.

### Non-Goals

No SLA commitment unless a separate support/service agreement exists.

## Phase 52 — Final Public Root Evidence Workflow

### Goal

Acquire and prove an agency-owned or agency-approved public feed root.

### Artifacts

- final-root approval intake template;
- DNS proof checklist;
- TLS and redirect proof checklist;
- five-feed public fetch packet;
- final-root validator record packet;
- redacted proxy/config summary;
- checksum manifest;
- prepared packet refresh instructions.

### Definition Of Done

For a real final-root run, a dated packet exists for a specific root with:

- owner/approval artifact;
- DNS proof;
- TLS certificate metadata;
- HTTP-to-HTTPS redirect proof;
- anonymous public fetches for all five public paths;
- current validation records;
- checksums;
- redaction notes.

### Phase 52 Closure

Phase 52 is complete blocker-only for the guarded final-root workflow. The
collector and audit tooling exist, but no real final public root or redacted
approval artifact was available, so no retained final-root evidence packet was
created and `docs/evidence/captured` remained unchanged.

### Non-Goals

Final-root proof does not by itself prove consumer acceptance or compliance.

## Phase 53 — Authorized Consumer Submission Execution

### Goal

Submit to selected consumers or aggregators only when authorization and official
path evidence exist.

### Artifacts

- target selection record;
- official path verification artifact;
- operator authorization artifact;
- submitted packet artifact;
- receipt/ticket/screenshot/email artifact;
- target-specific status transition update;
- redacted artifact directory.

### Definition Of Done

For a future authorized execution, only the selected target moves beyond
`prepared`, and only to the status supported by retained target-originated
evidence.

### Phase 53 Closure

Phase 53 is complete blocker-only. No operator authorization artifact,
official target path verification artifact, or target-originated artifact was
available. No target was selected, no consumer or aggregator was contacted, no
submission was made, no artifact was added, and all seven targets remain
`prepared`.

### Non-Goals

No guessed submission paths. No mass status changes.

## Phase 54 — Official Requirements Refresh

### Goal

Verify the current Caltrans / Cal-ITP requirements before final compliance work.

### Artifacts

- dated official-source review record;
- requirement diff against repo docs;
- updated `docs/requirements-calitp-compliance.md` if needed;
- updated readiness rows if needed;
- updated compliance packet checklist if needed.

### Definition Of Done

A maintainer can point to the exact official guidance version used for final
claim review.

### Phase 54 Closure

Phase 54 is complete for docs-only official-source mapping. It re-checked
current Caltrans / Cal-ITP and FTA sources and updated requirement mappings
only. It did not create evidence, change consumer statuses, change runtime
behavior, or claim compliance.

### Non-Goals

This phase does not itself create compliance evidence.

## Phase 55 — Compliance Evidence Packet Generator

### Goal

Generate a deployment-specific compliance/readiness packet from retained
artifacts without overstating claims.

### Artifacts

- `scripts/compliance-packet.sh` or equivalent;
- packet schema;
- redaction rules;
- claim-gating checklist;
- Operations Console export;
- template packet README.

### Required Packet Sections

- deployment identity and root;
- agency/operator authorization;
- public feed URLs;
- license/contact metadata;
- validation records;
- feed freshness and health;
- operations evidence;
- consumer packet/submission status;
- official requirements mapping;
- unsupported claims section;
- reviewer signoff.

### Definition Of Done

A packet can say exactly which readiness/compliance claims are supported and
which remain unsupported.

### Phase 55 Closure

Phase 55 is complete for local `.cache` blocker/draft compliance/readiness
packet generation and conservative audit tooling. It did not create retained
evidence, fetch live feeds, contact consumers, change consumer statuses, write
`docs/evidence/captured`, or claim compliance.

### Non-Goals

No automatic compliance claim. A human maintainer must review the packet.

## Phase 56 — Multi-Agency Hosting Hardening

### Goal

Harden repository-level multi-agency boundaries without claiming production
multi-tenant hosting.

### Artifacts

- complete for the current scope: validated per-agency public feed routing;
- complete for the current scope: private route/proxy diagnostics under
  ignored `.cache`;
- documented: backup/restore/export/evidence boundaries and blocked shared
  live-database tenant restore;
- documented: runtime admin remains agency-scoped; no global admin runtime
  model was added;
- focused route/proxy/tooling tests.

### Definition Of Done

Repository-level multi-agency route/proxy/tooling assumptions are proven
through tests and deployment docs. This does not certify production
multi-tenant hosting.

### Non-Goals

No hosted SaaS, production multi-tenant hosting, SLA/uptime, agency adoption,
consumer acceptance, production-readiness, vendor compatibility, marketplace
approval, or compliance claim.

## Phase 57 — Release Packaging And Supply Chain

### Goal

Make installation and upgrade practical from versioned local release packages.

### Artifacts

- complete for the current scope: local source package generation from
  `git archive HEAD`;
- complete for the current scope: SHA-256 checksum manifest;
- complete for the current scope: Go-module SBOM metadata;
- complete for the current scope: local provenance metadata;
- complete for the current scope: optional local image metadata when an
  operator supplies an image tag;
- complete for the current scope: audit and mutation tests.

### Definition Of Done

Maintainers can create and audit a local release package without publishing
artifacts or making hosted-service claims.

### Non-Goals

No registry push, production image publication, hosted service claim,
production-readiness claim, compliance claim, consumer acceptance claim, agency
adoption claim, marketplace approval, vendor compatibility claim, SLA/uptime
claim, or retained evidence creation.

## Phase 58 — Optional Marketplace / Vendor-Equivalent Pack

### Goal

Prepare optional non-code packaging needed for vendor-equivalent conversations,
without claiming marketplace approval. Complete for the template/audit scope.

### Artifacts

- complete for the current scope: BYOD/hardware intake template;
- complete for the current scope: implementation plan template;
- complete for the current scope: support boundaries template;
- complete for the current scope: SLA/KPI planning template;
- complete for the current scope: procurement response template;
- complete for the current scope: local audit script and mutation tests.

### Definition Of Done

The repo can support a procurement or marketplace-style discussion honestly
through templates while distinguishing code capability from service
commitments.

### Non-Goals

No marketplace listing, endorsement, paid support claim, vendor compatibility
claim, hardware certification claim, SLA/uptime claim, hosted service claim, or
production-readiness claim.

## Phase 59 — Real Pilot Closeout

### Goal

Document the closeout path for a real authorized pilot and retain public-safe
outcome evidence only when a real authorization and closeout artifact set
exists. Phase 59 closed blocker-only for the current repository state because
those artifacts were absent.

### Artifacts

- authorization record;
- pilot kickoff notes;
- GTFS import and validation results;
- telemetry or adapter notes;
- operations records;
- agency/operator feedback;
- blocker/continue/pause decision;
- closeout summary.

### Definition Of Done

A real pilot has a retained public-safe closeout packet that supports exactly
what happened and nothing more.

### Phase 59 Closeout

Phase 59 is complete blocker-only. No retained real pilot authorization,
kickoff, operations, feedback, and continue/pause/close decision artifact set
was available in the repository, so no pilot evidence packet was created.
Existing OCI/local pilot packets remain earlier-scope only, consumer tracker
state remains unchanged, and all seven consumer and aggregator targets remain
`prepared`.

### Non-Goals

Pilot participation does not automatically mean endorsement, compliance,
consumer acceptance, or production readiness.

## Phase 60 — Final Claim Review And Public Closeout

### Goal

Determine exactly what Open Transit RT can truthfully claim after all retained
evidence and official requirements review. Phase 60 is complete for the
current repository state and did not strengthen unsupported claims.

### Artifacts

- complete for the current scope: final claim-to-evidence table;
- complete for the current scope: unsupported-claims table;
- complete for the current scope: local read-only claim audit helper;
- complete for the current scope: mutation-style audit tests;
- complete for the current scope: bounded public/status docs updates;
- complete for the current scope: Phase 60 handoff.

### Definition Of Done

Every public claim has a retained evidence reference, is classified as
implementation capability or requirements context, or is removed/bounded.

### Phase 60 Closeout

Phase 60 is complete. It added `make audit-final-claim-review` and the final
claim review document. It created no retained evidence, wrote nothing under
`docs/evidence`, changed no consumer statuses, refreshed no packets or
artifact directories, and added no public launch, compliance, agency adoption,
consumer acceptance, hosted service, SLA/uptime, production-readiness, vendor,
marketplace, final-root, or production-grade ETA claim.

### Non-Goals

No aspirational language disguised as evidence. No public launch or compliance
statement was added because retained evidence does not support one.

## Recommended Execution Order

### Near-term product hardening

1. Phase 41 — Operator Smoke And Support Bundle — complete
2. Phase 42 — Reference Deployment Doctor — complete
3. Phase 43 — Operator UX Setup V2 — complete
4. Phase 44 — Telemetry Simulator And Device Trial — complete
5. Phase 45 — GTFS Quality Triage Loop — complete
6. Phase 46 — Validator Automation And Health Gates — complete
7. Phase 47 — Self-Hosted Operations Notifications — complete

These phases make the repo usable by small agencies and civic technologists.
They do not require external authorization.

### Integration and operations maturity

8. Phase 48 — AVL Adapter Runtime Path — complete
9. Phase 49 — External Predictor Runtime Adapter — complete
10. Phase 50 — Realtime Quality Backtesting — complete
11. Phase 51 — Operations Reliability And SLO Readiness — complete

These phases close later open-source GTFS-RT gaps: adapters, telemetry,
predictor replacement, quality measurement, and deployment-owned monitoring.

### Evidence and compliance path

12. Phase 52 — Final Public Root Evidence Workflow — complete blocker-only
13. Phase 53 — Authorized Consumer Submission Execution — complete blocker-only
14. Phase 54 — Official Requirements Refresh — complete
15. Phase 55 — Compliance Evidence Packet Generator — complete

These phases are complete for their current scopes. Phase 52 and Phase 53 are
blocker-only because the required real final-root and consumer-submission
artifacts were absent; Phase 54 and Phase 55 are requirements and local packet
workflow support only.

### Scale, release, and final closeout

16. Phase 56 — Multi-Agency Hosting Hardening
17. Phase 57 — Release Packaging And Supply Chain
18. Phase 58 — Optional Marketplace / Vendor-Equivalent Pack
19. Phase 59 — Real Pilot Closeout — complete blocker-only
20. Phase 60 — Final Claim Review And Public Closeout — complete

These phases make the project more reusable and honest at broader adoption
scale. Completion of the sequence does not reduce the listed missing evidence
gaps or authorize stronger public claims.

## Final Closeout Checklist

Before claiming CAL-ITP/Caltrans compliance for any deployment, all of these
must be true for the exact feed scope and URL root:

- official requirements were refreshed and recorded;
- agency/operator authorization exists;
- final public root approval exists;
- DNS/TLS/redirect evidence exists;
- all five public feed paths are anonymously fetchable;
- schedule validator has no blocking errors;
- Vehicle Positions validator has no blocking errors;
- Trip Updates validator has no blocking errors;
- Alerts validator has no blocking errors;
- license/contact metadata is agency-approved and publicly visible;
- feed freshness and operations evidence exists;
- backup/restore and incident workflows are documented for the deployment;
- consumer submission statuses match retained target-originated artifacts;
- production-grade ETA claims, if any, have real route/time-period quality
  evidence;
- every public claim is linked to retained evidence.

If any item is missing, use readiness/support wording rather than compliance
wording.
