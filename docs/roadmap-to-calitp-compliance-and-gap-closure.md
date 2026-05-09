# Roadmap To Cal-ITP Compliance And GTFS-RT Gap Closure

**Status:** proposed review draft, not yet committed  
**Intended repo path:** `docs/roadmap-to-calitp-compliance-and-gap-closure.md`  
**Generated for review:** 2026-05-08  
**Current repo baseline:** Phase 47 closed for private local operations notification summaries

This roadmap plans the path from the current self-hosted agency-reuse prototype
to a fully evidence-backed, open-source GTFS / GTFS-Realtime operations stack
that can support California transit data compliance work and close the most
important practical gaps in existing open-source GTFS-RT tooling.

This document is intentionally evidence-bounded. It does not claim that Open
Transit RT is CAL-ITP/Caltrans compliant today, accepted by any consumer, used
by an agency, production-ready for all deployments, vendor-compatible, or
marketplace/vendor-equivalent.

## Official Requirements Verification Gate

This roadmap is based on the repository's current requirements and readiness
docs, especially:

- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/master-plan-self-hosted-agency-reuse.md`
- `docs/handoffs/latest.md`

Before any public compliance claim or final compliance closeout, a maintainer
must re-check the current official Caltrans / Cal-ITP guidance and update this
roadmap if the requirements have changed.

The repo should keep using this wording until retained evidence supports more:

```text
Open Transit RT supports CAL-ITP-style readiness workflows and implements
technical foundations for stable GTFS and GTFS-Realtime publication.
```

Do not use this wording yet:

```text
Open Transit RT is CAL-ITP compliant.
Open Transit RT is Caltrans compliant.
Open Transit RT is accepted by trip planners.
Open Transit RT is production-ready for all agencies.
```

## Current Baseline

As of Phase 47, Open Transit RT has:

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
- hosted/operator OCI pilot evidence for the recorded pilot root.

The repo still lacks retained proof for:

- agency-owned or agency-approved final public root;
- agency adoption or approval;
- final-root DNS/TLS/redirect/five-feed evidence;
- final-root validator-clean records;
- consumer submission, review, acceptance, listing, display, or ingestion;
- official CAL-ITP/Caltrans compliance;
- real vendor AVL compatibility;
- real-world ETA quality;
- production-grade ETA quality;
- production multi-tenant hosting;
- hosted SaaS availability;
- paid support/SLA commitments.

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
| RQ-4A complete realtime feed set | all three GTFS-RT feed paths exist | Phases 44, 46, 50, 51 prove freshness, validation, and quality for real deployments |
| RQ-4B stable public production URLs | stable paths exist; final root missing | Phases 42 and 52 create preflight and final-root evidence workflow |
| RQ-4C validator-clean feeds | tooling exists; final no-error production evidence missing | Phases 45, 46, and 55 add triage, scheduled validators, and compliance packet generation |
| RQ-4D open license and discoverability | metadata model exists; agency-approved metadata missing | Phases 43, 52, and 55 harden metadata workflow and final-root public discovery proof |
| RQ-4E consumer ingestion workflow | prepared packets exist; no submissions or acceptance | Phases 53 and 55 handle authorized target submission and artifact-backed statuses |
| RQ-4F marketplace/vendor equivalence | explicitly not claimed | Phase 58 creates optional service/procurement pack without turning repo into SaaS |
| RQ-4G dashboard and scorecard | readiness page exists | Phases 43, 51, and 55 improve operator remediation and exportable compliance reporting |

## Master Phase Roadmap

The following phases are proposed future phases. They are not complete until a
maintainer accepts the scope, implementation happens, checks are run, and a
phase handoff is created.

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 41 | Operator Smoke And Support Bundle | Complete. Repeatable local/reference smoke helper and redaction-safe support bundle. |
| 42 | Reference Deployment Doctor | Complete. Read-only checks for env, services, reverse proxy posture, validators, DB, backups, restore readiness, and route boundaries. |
| 43 | Operator UX Setup V2 | Complete. Private grouped checklist for setup/readiness remediation and JSON export. |
| 44 | Telemetry Simulator And Device Trial | Complete. Safe simulator path that sends synthetic telemetry through real device-token ingest. |
| 45 | GTFS Quality Triage Loop | Operator-facing static GTFS validator warning/error triage linked to import/Studio actions. |
| 46 | Validator Automation And Health Gates | Scheduled validator runs, strict activation gates, history views, and deployment status exports. |
| 47 | Self-Hosted Operations Notifications | Complete. Private local notification drafts from existing diagnostic summaries only; no sending or evidence creation. |
| 48 | AVL Adapter Runtime Path | Private/authorized adapter send-mode pattern preserving `/v1/telemetry` and redaction boundaries. |
| 49 | External Predictor Runtime Adapter | Optional external predictor runtime path behind `internal/prediction.Adapter`. |
| 50 | Realtime Quality Backtesting | Real/simulated observed-arrival quality workflow with route/time-period metrics and ETA maturity gates. |
| 51 | Operations Reliability And SLO Readiness | Feed freshness, incident, backup/restore, alerting, and uptime evidence workflow. |
| 52 | Final Public Root Evidence Workflow | Agency-owned or agency-approved root acquisition, DNS/TLS/public fetch/validator evidence. |
| 53 | Authorized Consumer Submission Execution | Target-specific submission only with authorization and retained target-originated artifacts. |
| 54 | Official Requirements Refresh | Re-check current Caltrans / Cal-ITP guidance and update requirement mapping. |
| 55 | Compliance Evidence Packet Generator | Generate deployment-specific compliance/readiness packets with strict claim gating. |
| 56 | Multi-Agency Hosting Hardening | Tenant-safe routing, backup/restore/export/evidence, and operations model. |
| 57 | Release Packaging And Supply Chain | Versioned install artifacts, optional images, checksums, SBOM/provenance where practical. |
| 58 | Optional Marketplace / Vendor-Equivalent Pack | BYOD/hardware path, support boundaries, implementation templates, SLA/KPI templates. |
| 59 | Real Pilot Closeout | Run a real authorized pilot and retain public-safe feedback, operations, and blocker/continue evidence. |
| 60 | Final Claim Review And Public Closeout | Decide exactly what may be claimed after evidence and official requirements review. |

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

A dated packet exists for a specific root with:

- owner/approval artifact;
- DNS proof;
- TLS certificate metadata;
- HTTP-to-HTTPS redirect proof;
- anonymous public fetches for all five public paths;
- current validation records;
- checksums;
- redaction notes.

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

Only the selected target moves beyond `prepared`, and only to the status
supported by retained target-originated evidence.

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

### Non-Goals

No automatic compliance claim. A human maintainer must review the packet.

## Phase 56 — Multi-Agency Hosting Hardening

### Goal

Make production multi-agency hosting safe if the project chooses that path.

### Artifacts

- per-agency public feed routing;
- tenant-safe backup/restore/export/evidence tooling;
- global admin and tenant admin model;
- isolation tests for all critical workflows;
- operational docs.

### Definition Of Done

Multi-agency hosting assumptions are proven through tests and deployment docs,
not inferred from selected repository isolation tests.

### Non-Goals

No hosted SaaS claim unless an actual service offering exists.

## Phase 57 — Release Packaging And Supply Chain

### Goal

Make installation and upgrade practical from versioned releases.

### Artifacts

- release artifact matrix;
- optional published Docker images;
- checksums;
- SBOM/provenance where practical;
- versioned config migration notes;
- upgrade/rollback test matrix;
- release validation checklist.

### Definition Of Done

A small agency implementer can install a known version and upgrade or roll back
without building from a random local checkout.

### Non-Goals

Release packaging is not hosted SaaS or paid support.

## Phase 58 — Optional Marketplace / Vendor-Equivalent Pack

### Goal

Prepare optional non-code packaging needed for vendor-equivalent conversations,
without claiming marketplace approval.

### Artifacts

- BYOD/hardware strategy;
- implementation plan template;
- support boundaries;
- SLA/KPI templates;
- procurement-oriented one-pager;
- training material;
- operator responsibility matrix;
- security and data-handling summary.

### Definition Of Done

The repo can support a procurement or marketplace-style discussion honestly,
while distinguishing code capability from service commitments.

### Non-Goals

No marketplace listing, endorsement, or paid support claim.

## Phase 59 — Real Pilot Closeout

### Goal

Run a real authorized pilot and retain public-safe outcome evidence.

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

### Non-Goals

Pilot participation does not automatically mean endorsement, compliance,
consumer acceptance, or production readiness.

## Phase 60 — Final Claim Review And Public Closeout

### Goal

Determine exactly what Open Transit RT can truthfully claim after all retained
evidence and official requirements review.

### Artifacts

- final claim-to-evidence table;
- unsupported-claims table;
- public README wording update;
- public launch or compliance statement only if supported;
- archived evidence index;
- maintainer signoff.

### Definition Of Done

Every public claim has a retained evidence reference or is removed.

### Non-Goals

No aspirational language disguised as evidence.

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

8. Phase 48 — AVL Adapter Runtime Path
9. Phase 49 — External Predictor Runtime Adapter
10. Phase 50 — Realtime Quality Backtesting
11. Phase 51 — Operations Reliability And SLO Readiness

These phases close later open-source GTFS-RT gaps: adapters, telemetry,
predictor replacement, quality measurement, and deployment-owned monitoring.

### Evidence and compliance path

12. Phase 52 — Final Public Root Evidence Workflow
13. Phase 53 — Authorized Consumer Submission Execution
14. Phase 54 — Official Requirements Refresh
15. Phase 55 — Compliance Evidence Packet Generator

These phases should run only when an operator has a real deployment context or
when maintainers are ready to refresh official requirements.

### Scale, release, and final closeout

16. Phase 56 — Multi-Agency Hosting Hardening
17. Phase 57 — Release Packaging And Supply Chain
18. Phase 58 — Optional Marketplace / Vendor-Equivalent Pack
19. Phase 59 — Real Pilot Closeout
20. Phase 60 — Final Claim Review And Public Closeout

These phases make the project sustainable and honest at broader adoption scale.

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
