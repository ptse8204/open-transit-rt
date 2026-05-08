# Roadmap To Cal-ITP Compliance And GTFS-RT Gap Closure

**Status:** proposed review draft  
**Baseline:** Phase 40 closed for the guided self-hosted operator trial scope  
**Purpose:** plan Open Transit RT from the current self-hosted agency-reuse prototype to an evidence-backed open-source GTFS / GTFS-Realtime operations stack that can support California transit data compliance work and close practical gaps in existing open-source GTFS-RT tooling.

This roadmap is evidence-bounded. It does not claim that Open Transit RT is CAL-ITP/Caltrans compliant today, accepted by any consumer, adopted by any agency, production-ready for every deployment, vendor-compatible, or marketplace/vendor-equivalent.

## Official Requirements Verification Gate

This roadmap is based on the repo's current requirements and readiness docs:

- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/master-plan-self-hosted-agency-reuse.md`
- `docs/handoffs/latest.md`

Before any public compliance claim or final compliance closeout, a maintainer must re-check the current official Caltrans / Cal-ITP guidance and update this roadmap if requirements have changed.

Allowed wording until retained evidence supports more:

```text
Open Transit RT supports CAL-ITP-style readiness workflows and implements technical foundations for stable GTFS and GTFS-Realtime publication.
```

Unsupported wording today:

```text
Open Transit RT is CAL-ITP compliant.
Open Transit RT is Caltrans compliant.
Open Transit RT is accepted by trip planners.
Open Transit RT is production-ready for all agencies.
```

## Current Baseline

As of Phase 40, Open Transit RT has:

- GTFS ZIP import and GTFS Studio draft/publish workflows;
- stable public feed paths for `/public/feeds.json`, `/public/gtfs/schedule.zip`, Vehicle Positions, Trip Updates, and Alerts;
- authenticated telemetry ingest with device bearer tokens;
- conservative deterministic trip matching;
- GTFS-RT Vehicle Positions, Trip Updates, and Alerts publication paths;
- Trip Updates behind `internal/prediction.Adapter`;
- pinned MobilityData validator workflows;
- local app packaging and reusable agency onboarding through `make agency-pilot-up`;
- an OCI/OCL-style reference deployment path;
- an Integration Adapter Kit;
- an authenticated CAL-ITP-style Operations Console readiness page;
- a guided self-hosted operator trial;
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

The end state has three levels. They must not be collapsed into one claim.

### Level 1 — Self-Hosted Technical Maturity

A small agency, operator, or civic technologist can install from a release, deploy the reference server pattern, import or author GTFS, publish stable schedule and GTFS-RT feeds, validate all feeds, connect telemetry through documented adapters, monitor feed freshness and validation status, operate backup/restore and rollback, and produce a redaction-safe support bundle.

This level supports self-hosted reuse. It is not yet a compliance or consumer acceptance claim.

### Level 2 — Evidence-Backed California Data Readiness

For a specific deployment and feed root, retained evidence exists for stable public HTTPS URLs, public static GTFS, all three GTFS-Realtime feeds, current no-error validation results, agency-approved license/contact metadata, public discoverability metadata, operations evidence, consumer packet readiness, and official requirements mapping checked against current Caltrans / Cal-ITP guidance.

This level may support stronger readiness language, but still does not prove consumer acceptance unless target-originated evidence exists.

### Level 3 — Full External Closeout

For the exact feed scope and URL root, retained redacted artifacts prove agency/operator authorization, final-root approval, official consumer submission path verification, target-specific submissions, target-originated review/acceptance/ingestion/listing/display statuses, and final claim review against current official requirements.

Only this level can support public language about consumer acceptance or final compliance status for a named deployment.

## Plugin / Adapter Component Strategy

Open Transit RT should close open-source GTFS-RT gaps through adapter-bound components, not through a closed all-in-one CAD/AVL stack.

The word `plugin` in this roadmap means a replaceable component with a stable contract. It does not require Go dynamic plugins. Safer shapes include command adapters, sidecars, HTTP adapters, fixtures, manifests, and deployment-owned processes.

| Component area | Current repo status | Long-term target | Boundary rule |
| --- | --- | --- | --- |
| Static GTFS input | ZIP import and GTFS Studio exist | import, edit, validate, triage, publish, rollback | published GTFS stays separate from draft GTFS |
| Telemetry ingest | authenticated `/v1/telemetry` exists | device, simulator, and AVL adapter paths | vendor payloads transform before telemetry ingest |
| AVL/vendor adapters | synthetic dry-run adapter kit exists | private/authorized real adapters can be added outside core | no named vendor claim without retained evidence |
| Prediction engines | deterministic adapter exists; external boundary evaluated | internal predictor, TheTransitClock-like predictor, or replacement predictors | all predictors stay behind `internal/prediction.Adapter` |
| Validators | pinned static and realtime workflows exist | scheduled, strict, reportable validators with triage | validators verify outputs; they do not own core state |
| Monitoring | metrics and pilot ops helpers exist | Prometheus/Grafana/OpenTelemetry optional packs | observability must not break feed generation |
| Consumer workflows | prepared packets and tracker exist | target-specific workflows with retained artifacts | no target status change without target-originated evidence |
| Evidence/redaction | policies and packet patterns exist | support bundles and evidence packets with strict redaction | no secrets, raw telemetry, or private operator artifacts committed |
| Deployment | OCI/OCL-style reference docs exist | preflighted, release-based, rollback-safe deployment | public feed edge and private admin edge stay separate |

## Mapping To Repository Compliance Requirements

| Requirement | Current status | Roadmap closure path |
| --- | --- | --- |
| RQ-4A complete realtime feed set | all three GTFS-RT feed paths exist | Phases 44, 46, 50, 51 prove freshness, validation, and quality for real deployments |
| RQ-4B stable public production URLs | stable paths exist; final root missing | Phases 42 and 52 create preflight and final-root evidence workflow |
| RQ-4C validator-clean feeds | tooling exists; final no-error production evidence missing | Phases 45, 46, and 55 add triage, scheduled validators, and compliance packet generation |
| RQ-4D open license and discoverability | metadata model exists; agency-approved metadata missing | Phases 43, 52, and 55 harden metadata workflow and final-root public discovery proof |
| RQ-4E consumer ingestion workflow | prepared packets exist; no submissions or acceptance | Phases 53 and 55 handle authorized target submission and artifact-backed statuses |
| RQ-4F marketplace/vendor equivalence | explicitly not claimed | Phase 58 creates optional service/procurement pack without turning repo into SaaS |
| RQ-4G dashboard and scorecard | readiness page exists | Phases 43, 51, and 55 improve operator remediation and exportable compliance reporting |

## Proposed Phase Roadmap

These phases are proposed future work. They are not complete until a maintainer accepts scope, implementation happens, checks are run, and a handoff is created.

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 41 | Operator Smoke And Support Bundle | Repeatable local/reference smoke helper and redaction-safe support bundle. |
| 42 | Reference Deployment Doctor | Preflight checks for env, services, reverse proxy, validators, DB, backups, and route boundaries. |
| 43 | Operator UX Setup V2 | Stronger setup/readiness UI with grouped remediation and exportable operator checklist. |
| 44 | Telemetry Simulator And Device Trial | Safe simulator path that sends synthetic telemetry through real device-token ingest. |
| 45 | GTFS Quality Triage Loop | Operator-facing static GTFS validator warning/error triage linked to import/Studio actions. |
| 46 | Validator Automation And Health Gates | Scheduled validator runs, strict activation gates, history views, and deployment status exports. |
| 47 | Observability Plugin Pack | Optional Prometheus/Grafana/OpenTelemetry deployment-owned integration pattern. |
| 48 | AVL Adapter Runtime Path | Controlled path from dry-run transforms to authorized private telemetry send mode. |
| 49 | External Predictor Runtime Adapter | Optional external predictor runtime integration behind `internal/prediction.Adapter`. |
| 50 | Realtime Quality Backtesting | Observed-arrival and route/time-period quality metrics without ETA overclaims. |
| 51 | Operations Reliability And SLO Readiness | Long-running operations, incident, backup, restore, validation, and alerting evidence workflows. |
| 52 | Final Public Root Evidence Workflow | Agency-owned or agency-approved final-root proof packet. |
| 53 | Authorized Consumer Submission Execution | Target-specific submissions only with official-path and operator authorization artifacts. |
| 54 | Official Requirements Refresh | Dated official Caltrans / Cal-ITP requirements review before final claims. |
| 55 | Compliance Evidence Packet Generator | Deployment-specific packet that maps retained evidence to supported claims. |
| 56 | Multi-Agency Hosting Hardening | Tenant-safe routing, backup/restore/export/evidence tooling, and isolation tests. |
| 57 | Release Packaging And Supply Chain | Versioned release artifacts, checksums, optional images, and upgrade/rollback matrix. |
| 58 | Optional Marketplace / Vendor-Equivalent Pack | Non-code support/procurement/SLA/KPI templates without marketplace approval claims. |
| 59 | Real Pilot Closeout | Authorized real pilot closeout with public-safe retained outcome evidence. |
| 60 | Final Claim Review And Public Closeout | Every public claim linked to retained evidence, or removed. |

## Phase Details

### Phase 41 — Operator Smoke And Support Bundle

Create `scripts/operator-smoke.sh`, `scripts/support-bundle.sh`, `make operator-smoke`, and `make support-bundle`. The smoke helper checks the five public feeds, validator tooling status, optional admin readiness access through the private/admin path, optional allowlisted validation runs, and the synthetic AVL dry-run. The support bundle collects redaction-safe diagnostics only. It must exclude secrets, raw telemetry, private vendor payloads, tokens, cookies, keys, ACME material, raw `.env` files, and raw database dumps.

### Phase 42 — Reference Deployment Doctor

Add preflight checks for the OCI/OCL-style deployment path: environment variables, generated secrets presence without values, migrations, service health, reverse proxy public/private boundaries, validators, public feed reachability, backup path, restore drill readiness, and rollback readiness.

### Phase 43 — Operator UX Setup V2

Improve `/admin/operations/setup` and `/admin/operations/readiness` so nontechnical operators see grouped next actions, status source, claim boundary, links to docs, and exportable local checklist output. The page remains authenticated and does not create evidence or claim compliance.

### Phase 44 — Telemetry Simulator And Device Trial

Add a safe simulator path that sends synthetic telemetry through real device-token ingest. Include no-hardware tests for accepted, stale, duplicate, out-of-order, off-shape, and unmatched behavior. This is still not real AVL evidence.

### Phase 45 — GTFS Quality Triage Loop

Turn validator warnings and errors into operator-facing remediation steps linked to import reports and GTFS Studio where possible. Include triage guidance for common issues such as expired calendars, long route names, unused shapes, bad references, invalid times, missing accessibility metadata, and frequency/block issues.

### Phase 46 — Validator Automation And Health Gates

Add scheduled validator runs, strict production activation gates where configured, validation history views, feed health rollups, and deployment status exports. Validator success remains supporting evidence, not consumer acceptance.

### Phase 47 — Observability Plugin Pack

Define optional deployment-owned observability packs for Prometheus/Grafana and OpenTelemetry without making them required runtime dependencies. Provide metrics inventory, dashboard templates, alert examples, and failure behavior rules.

### Phase 48 — AVL Adapter Runtime Path

Move from dry-run-only synthetic adapter examples toward an authorized private runtime path. Network send mode may exist only with explicit credentials outside the repo, allowlisted target URLs, dry-run default, clear diagnostics, and no named vendor compatibility claim without evidence.

### Phase 49 — External Predictor Runtime Adapter

Add an optional external predictor runtime path behind `internal/prediction.Adapter`. Vehicle Positions, telemetry ingest, assignments, and GTFS import continue independently if the predictor fails. Any TheTransitClock-like integration requires dependency/license review and fail-closed output validation.

### Phase 50 — Realtime Quality Backtesting

Create real or approved observed-arrival backtesting workflows. Track eligible candidates, Trip Updates coverage, future stop coverage, stale/unknown/ambiguous rates, mean/median absolute error where observed arrivals exist, percentile error by route/stop/time of day, cancellation/alert linkage, and external predictor fallback rate.

### Phase 51 — Operations Reliability And SLO Readiness

Make long-running operations visible: feed freshness targets, incident records, alert delivery proof, backup/restore cadence, validator failure response, key rotation records, and operator handoff notes. This prepares for reliability evidence but does not create an SLA.

### Phase 52 — Final Public Root Evidence Workflow

Acquire and prove an agency-owned or agency-approved public feed root. Required packet: owner/approval artifact, DNS proof, TLS metadata, HTTP-to-HTTPS redirect proof, anonymous fetches for all five public paths, validator records, proxy/config summary, checksums, and redaction notes.

### Phase 53 — Authorized Consumer Submission Execution

Submit only to selected consumers/aggregators when authorization, official path verification, and retained target-originated artifacts exist. Only the selected target may move beyond `prepared`, and only to the status supported by evidence.

### Phase 54 — Official Requirements Refresh

Re-check current official Caltrans / Cal-ITP guidance. Record source version/date, requirement diffs, and doc updates before any final claim review.

### Phase 55 — Compliance Evidence Packet Generator

Generate a deployment-specific packet with deployment identity, final root, authorization, public URLs, license/contact metadata, validation records, freshness/health, operations evidence, consumer status, official requirements mapping, unsupported claims, and reviewer signoff. The packet does not automatically create a compliance claim.

### Phase 56 — Multi-Agency Hosting Hardening

If the project chooses production multi-agency hosting, add per-agency public feed routing, tenant-safe backup/restore/export/evidence tooling, global/tenant admin models, and isolation tests for all critical workflows.

### Phase 57 — Release Packaging And Supply Chain

Make releases practical: artifact matrix, optional Docker images, checksums, SBOM/provenance where practical, config migration notes, upgrade/rollback test matrix, and release validation checklist.

### Phase 58 — Optional Marketplace / Vendor-Equivalent Pack

Prepare non-code materials for procurement or marketplace-style discussions: BYOD/hardware strategy, implementation plan, support boundaries, SLA/KPI templates, training, responsibility matrix, and security/data-handling summary. Do not claim marketplace approval.

### Phase 59 — Real Pilot Closeout

Run a real authorized pilot and retain public-safe outcome evidence: authorization, kickoff notes, GTFS import/validation results, telemetry/adapter notes, operations records, agency/operator feedback, and blocker/continue/pause decision.

### Phase 60 — Final Claim Review And Public Closeout

Create a final claim-to-evidence table, unsupported-claims table, README/public wording update, evidence index, and maintainer signoff. Every public claim must have retained evidence or be removed.

## Recommended Execution Order

### Near-term product hardening

1. Phase 41 — Operator Smoke And Support Bundle
2. Phase 42 — Reference Deployment Doctor
3. Phase 43 — Operator UX Setup V2
4. Phase 44 — Telemetry Simulator And Device Trial
5. Phase 45 — GTFS Quality Triage Loop
6. Phase 46 — Validator Automation And Health Gates

These phases make the repo easier for small agencies and civic technologists to use. They do not require external authorization.

### Integration and operations maturity

7. Phase 47 — Observability Plugin Pack
8. Phase 48 — AVL Adapter Runtime Path
9. Phase 49 — External Predictor Runtime Adapter
10. Phase 50 — Realtime Quality Backtesting
11. Phase 51 — Operations Reliability And SLO Readiness

These phases close major open-source GTFS-RT gaps: monitoring, adapters, telemetry, predictor replacement, and quality measurement.

### Evidence and compliance path

12. Phase 52 — Final Public Root Evidence Workflow
13. Phase 53 — Authorized Consumer Submission Execution
14. Phase 54 — Official Requirements Refresh
15. Phase 55 — Compliance Evidence Packet Generator

These phases should run only when an operator has a real deployment context or when maintainers are ready to refresh official requirements.

### Scale, release, and final closeout

16. Phase 56 — Multi-Agency Hosting Hardening
17. Phase 57 — Release Packaging And Supply Chain
18. Phase 58 — Optional Marketplace / Vendor-Equivalent Pack
19. Phase 59 — Real Pilot Closeout
20. Phase 60 — Final Claim Review And Public Closeout

## Final Closeout Checklist

Before claiming CAL-ITP/Caltrans compliance for any deployment, all of these must be true for the exact feed scope and URL root:

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
- production-grade ETA claims, if any, have real route/time-period quality evidence;
- every public claim is linked to retained evidence.

If any item is missing, use readiness/support wording rather than compliance wording.
