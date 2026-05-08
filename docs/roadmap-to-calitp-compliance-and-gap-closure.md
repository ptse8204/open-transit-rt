# Roadmap To Cal-ITP Compliance And GTFS-RT Gap Closure

**Status:** proposed review draft committed on a review branch  
**Current repo baseline:** Phase 40 closed for the guided self-hosted operator trial scope  
**Purpose:** plan the path from the current self-hosted agency-reuse prototype to an evidence-backed open-source GTFS / GTFS-Realtime operations stack that can support California transit data readiness work and close the practical gaps in existing GTFS-RT open-source tooling.

This roadmap is evidence-bounded. It does not claim that Open Transit RT is CAL-ITP/Caltrans compliant today, accepted by any consumer, used by an agency, production-ready for all deployments, vendor-compatible, or marketplace/vendor-equivalent.

## Official Requirements Verification Gate

This roadmap is based on the repository's current requirements and readiness docs:

- `docs/requirements-calitp-compliance.md`
- `docs/california-readiness-summary.md`
- `docs/compliance-evidence-checklist.md`
- `docs/master-plan-self-hosted-agency-reuse.md`
- `docs/handoffs/latest.md`

Before any public compliance claim or final compliance closeout, a maintainer must re-check the current official Caltrans / Cal-ITP guidance and update this roadmap if requirements changed.

Allowed wording until retained evidence supports more:

```text
Open Transit RT supports CAL-ITP-style readiness workflows and implements technical foundations for stable GTFS and GTFS-Realtime publication.
```

Do not use yet:

```text
Open Transit RT is CAL-ITP compliant.
Open Transit RT is Caltrans compliant.
Open Transit RT is accepted by trip planners.
Open Transit RT is production-ready for all agencies.
```

## Current Baseline

As of Phase 40, Open Transit RT has GTFS ZIP import, GTFS Studio draft/publish, stable public feed paths, authenticated telemetry ingest, conservative deterministic trip matching, Vehicle Positions, Trip Updates behind `internal/prediction.Adapter`, Alerts, pinned MobilityData validators, local app packaging, OCI/OCL-style reference deployment docs, `make agency-pilot-up`, an Integration Adapter Kit, an authenticated CAL-ITP-style readiness page, a guided self-hosted operator trial, synthetic AVL adapter fixtures, deterministic realtime-quality replay fixtures, prepared consumer packets for seven targets, and hosted/operator OCI pilot evidence for the recorded pilot root.

The repo still lacks retained proof for agency-owned or agency-approved final public root, agency adoption or approval, final-root DNS/TLS/redirect/five-feed evidence, final-root validator-clean records, consumer submission/review/acceptance/listing/display/ingestion, official CAL-ITP/Caltrans compliance, real vendor AVL compatibility, real-world ETA quality, production-grade ETA quality, production multi-tenant hosting, hosted SaaS availability, and paid support/SLA commitments.

## End State

The end state has three separate levels that must not be collapsed into one claim.

### Level 1 — Self-Hosted Technical Maturity

A small agency, operator, or civic technologist can install Open Transit RT from a release, deploy the reference server pattern, import or author GTFS, publish stable schedule and GTFS-RT feeds, validate all feeds, connect telemetry through a documented adapter path, monitor feed freshness and validation status, operate backup/restore/rollback/incident workflows, and produce a redaction-safe support bundle.

This level supports self-hosted reuse. It is not a compliance or consumer acceptance claim.

### Level 2 — Evidence-Backed California Data Compliance Readiness

For a specific deployment and feed root, retained evidence exists for stable public HTTPS URLs, public static GTFS, all three GTFS-Realtime feeds, current no-error validation results, agency-approved license/contact metadata, public discoverability metadata, current operations evidence, consumer packet readiness, and official requirements mapping checked against current Caltrans / Cal-ITP guidance.

This level may support stronger readiness language, but still does not prove consumer acceptance unless target-originated evidence exists.

### Level 3 — Full External Closeout

For the exact feed scope and URL root, retained redacted artifacts prove agency/operator authorization, final-root approval, official consumer submission path verification, target-specific submissions, target-originated status artifacts, and final claim review against current official requirements.

Only this level can support public language about consumer acceptance or final compliance status for a named deployment.

## Plugin / Adapter Component Strategy

Open Transit RT should close GTFS-RT tooling gaps through adapter-bound components, not through a closed all-in-one system. In this roadmap, `plugin` means a replaceable component with a stable contract. It does not require Go dynamic plugins. Safer shapes include command adapters, sidecars, HTTP adapters, fixtures, manifests, and deployment-owned processes.

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

| Requirement | Current status | Roadmap closure path |
| --- | --- | --- |
| RQ-4A complete realtime feed set | all three GTFS-RT feed paths exist | Phases 44, 46, 50, and 51 prove freshness, validation, and quality for real deployments |
| RQ-4B stable public production URLs | stable paths exist; final root missing | Phases 42 and 52 create preflight and final-root evidence workflow |
| RQ-4C validator-clean feeds | tooling exists; final no-error production evidence missing | Phases 45, 46, and 55 add triage, scheduled validators, and compliance packet generation |
| RQ-4D open license and discoverability | metadata model exists; agency-approved metadata missing | Phases 43, 52, and 55 harden metadata workflow and final-root public discovery proof |
| RQ-4E consumer ingestion workflow | prepared packets exist; no submissions or acceptance | Phases 53 and 55 handle authorized target submission and artifact-backed statuses |
| RQ-4F marketplace/vendor equivalence | explicitly not claimed | Phase 58 creates optional service/procurement pack without turning repo into SaaS |
| RQ-4G dashboard and scorecard | readiness page exists | Phases 43, 51, and 55 improve operator remediation and exportable compliance reporting |

## Master Phase Roadmap

The following phases are proposed future phases. They are not complete until a maintainer accepts the scope, implementation happens, checks are run, and a phase handoff is created.

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 41 | Operator Smoke And Support Bundle | Repeatable local/reference smoke helper and redaction-safe support bundle. |
| 42 | Reference Deployment Doctor | Preflight checks for env, services, reverse proxy, validators, DB, backups, and route boundaries. |
| 43 | Operator UX Setup V2 | Stronger setup/readiness UI with grouped remediation and exportable operator checklist. |
| 44 | Telemetry Simulator And Device Trial | Safe simulator path that sends synthetic telemetry through real device-token ingest. |
| 45 | GTFS Quality Triage Loop | Operator-facing static GTFS validator warning/error triage linked to import/Studio actions. |
| 46 | Validator Automation And Health Gates | Scheduled validator runs, strict activation gates, history views, and deployment status exports. |
| 47 | Observability Plugin Pack | Optional Prometheus/Grafana/OpenTelemetry deployment-owned integration package. |
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

## Phase Detail

### Phase 41 — Operator Smoke And Support Bundle

Turn Phase 40 into a repeatable smoke workflow and redaction-safe support bundle. Add `scripts/operator-smoke.sh`, `scripts/support-bundle.sh`, `make operator-smoke`, `make support-bundle`, and docs explaining what is safe to share. Output must stay under ignored `.cache/` storage and must not create evidence.

### Phase 42 — Reference Deployment Doctor

Add a deployment doctor for the OCI/OCL-style path. Check environment values, missing secrets without printing them, DB connectivity, migrations, service health, reverse proxy public/private route boundaries, validator availability, backup path readiness, and rollback readiness. This is preflight only, not proof of production readiness.

### Phase 43 — Operator UX Setup V2

Improve `/admin/operations/setup` and `/admin/operations/readiness` so operators can see grouped readiness rows, plain-language remediation, source/evidence fields, exportable operator checklists, and claim boundaries. Do not expose admin surfaces publicly.

### Phase 44 — Telemetry Simulator And Device Trial

Add a safe simulator path that sends synthetic telemetry through the real `/v1/telemetry` contract using device bearer tokens. This closes the gap between dry-run transforms and actual ingest without using real vendor payloads or claiming hardware support.

### Phase 45 — GTFS Quality Triage Loop

Make static GTFS validation output actionable. Surface errors/warnings by category, link warnings to import/Studio remediation, retain warning history, and explain which warnings block publish in strict mode. This helps move from “validator ran” to “operator can fix the schedule.”

### Phase 46 — Validator Automation And Health Gates

Add scheduled validator runs, strict production activation gates where configured, validation history summaries, and feed health exports. Validator success remains supporting evidence only and does not imply consumer acceptance.

### Phase 47 — Observability Plugin Pack

Provide optional deployment-owned observability integrations for Prometheus, Grafana, and/or OpenTelemetry. The app must continue if observability backends are absent. Dashboards and alert rules are operational aids, not SLA evidence by themselves.

### Phase 48 — AVL Adapter Runtime Path

Define how authorized/private adapters can send transformed records to `/v1/telemetry` while keeping credentials and raw vendor payloads outside the public repo. Add a no-credential send-mode pattern, sidecar docs, and evidence boundaries. No named vendor support without retained authorization and proof.

### Phase 49 — External Predictor Runtime Adapter

Implement an optional runtime external predictor adapter behind `internal/prediction.Adapter`. Preserve Vehicle Positions independence, fail closed on malformed/stale/wrong-agency outputs, and document dependency/license review. No production-grade ETA claim from integration alone.

### Phase 50 — Realtime Quality Backtesting

Add a backtesting workflow for observed arrival/departure comparison where authorized data exists. Track route/time-period metrics, prediction coverage, unknown/stale/ambiguous rates, and fallback behavior. Define ETA maturity gates before any production-grade ETA language.

### Phase 51 — Operations Reliability And SLO Readiness

Make long-running operations visible: feed freshness, validation cadence, backup/restore drills, incident records, alert delivery proof, secret rotation records, and operator handoff notes. Do not claim an SLA without a separate support/service commitment.

### Phase 52 — Final Public Root Evidence Workflow

Acquire and prove an agency-owned or agency-approved public feed root. Required artifacts include owner/approval record, DNS proof, TLS metadata, HTTP-to-HTTPS redirect proof, anonymous fetches for all five public paths, current validators, proxy/config summary, and checksums.

### Phase 53 — Authorized Consumer Submission Execution

Submit to a selected consumer or aggregator only when operator authorization and official path evidence exist. Update only the selected target status and only from retained target-originated artifacts.

### Phase 54 — Official Requirements Refresh

Re-check official Caltrans / Cal-ITP guidance before final compliance work. Record the official-source review date, requirement diff, and docs updates required before final claim review.

### Phase 55 — Compliance Evidence Packet Generator

Create a deployment-specific packet generator that assembles root identity, authorization, public feed URLs, license/contact metadata, validation records, operations evidence, consumer status, official requirements mapping, unsupported claims, and reviewer signoff. A human maintainer must still review claims.

### Phase 56 — Multi-Agency Hosting Hardening

If the project chooses a multi-agency hosting path, prove per-agency public routing, tenant-safe backup/restore/export/evidence, global/tenant admin roles, and isolation tests for all critical workflows. Do not claim hosted SaaS without an actual service offering.

### Phase 57 — Release Packaging And Supply Chain

Make install/upgrade practical from known releases. Add release artifact matrix, optional published Docker images, checksums, SBOM/provenance where practical, config migration notes, and upgrade/rollback tests. Release packaging is not hosted SaaS.

### Phase 58 — Optional Marketplace / Vendor-Equivalent Pack

Prepare optional non-code packaging for vendor-equivalent conversations: BYOD/hardware strategy, implementation plan template, support boundaries, SLA/KPI templates, procurement one-pager, training material, responsibility matrix, and security/data-handling summary. Do not claim marketplace approval.

### Phase 59 — Real Pilot Closeout

Run a real authorized pilot and retain public-safe closeout evidence: authorization, kickoff notes, GTFS import/validation, telemetry/adapter notes, operations records, agency/operator feedback, and blocker/continue/pause decision. Pilot participation does not automatically mean endorsement or compliance.

### Phase 60 — Final Claim Review And Public Closeout

Create a final claim-to-evidence table, unsupported-claims table, public README wording update, archived evidence index, and maintainer signoff. Every public claim must have retained evidence or be removed.

## Recommended Execution Order

Near-term product hardening:

1. Phase 41 — Operator Smoke And Support Bundle
2. Phase 42 — Reference Deployment Doctor
3. Phase 43 — Operator UX Setup V2
4. Phase 44 — Telemetry Simulator And Device Trial
5. Phase 45 — GTFS Quality Triage Loop
6. Phase 46 — Validator Automation And Health Gates

Integration and operations maturity:

7. Phase 47 — Observability Plugin Pack
8. Phase 48 — AVL Adapter Runtime Path
9. Phase 49 — External Predictor Runtime Adapter
10. Phase 50 — Realtime Quality Backtesting
11. Phase 51 — Operations Reliability And SLO Readiness

Evidence and compliance path:

12. Phase 52 — Final Public Root Evidence Workflow
13. Phase 53 — Authorized Consumer Submission Execution
14. Phase 54 — Official Requirements Refresh
15. Phase 55 — Compliance Evidence Packet Generator

Scale, release, and final closeout:

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
