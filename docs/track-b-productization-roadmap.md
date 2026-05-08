# Track B — Agency Productization, Release, And Real-World Adoption

## Status

Track B is active as the productization and self-hosted agency reuse roadmap.
Phases 22 through 39 are closed for their documented scopes. Use
`docs/handoffs/latest.md` and `docs/future-roadmap-post-outcome-c.md` for the
current continuation state.

## Purpose

Track B turns Open Transit RT from a technically capable open-source backend into an agency-adoptable pilot and release package.

The project now has core GTFS/GTFS Realtime publication, local app packaging,
hosted OCI pilot evidence, consumer packet preparation, governance docs,
external-proof workflow, and a self-hosted agency reuse master plan. Track B
now focuses on making the working code easier for agencies and implementers to
reuse:

- versioned releases;
- OCI/OCL-style reference deployment productization;
- reusable agency GTFS onboarding;
- device / AVL integration boundaries;
- setup and operator UX;
- multi-agency isolation proof;
- stronger operations hardening;
- real-world realtime quality evidence;
- optional consumer submissions when authorized and evidenced;
- agency pilot packaging;
- public ecosystem launch messaging.

## Track B Mission

Make Open Transit RT practical for real small-agency pilots and eventual production use without weakening the project’s evidence and truthfulness boundaries.

Track B should help operators answer:

- What version should I install?
- How do I upgrade or roll back?
- How do I deploy the OCI/OCL-style reference pattern?
- How do I import real GTFS and fix validation failures?
- How do I connect real telemetry devices or vendors?
- How do I complete setup without deep developer knowledge?
- Can the system safely isolate multiple agencies?
- What operations proof do I need over time?
- What realtime quality evidence exists?
- What does it take to submit to downstream consumers?
- What does a real agency pilot require?

## Truthfulness Boundary

Track B must not claim:

- CAL-ITP/Caltrans compliance without evidence;
- consumer submission, review, acceptance, ingestion, or display without retained target-originated evidence;
- agency endorsement;
- hosted SaaS availability;
- paid support or SLA coverage;
- marketplace/vendor equivalence;
- production-grade ETA quality without real-world evaluation;
- multi-tenant production readiness without isolation proof.

Allowed language includes:

- “technical foundations”;
- “local evaluation package”;
- “pilot operations path”;
- “agency-owned domain readiness checklist”;
- “prepared packets”;
- “operator workflow”;
- “replay measurement evidence.”

## Relationship To Track A

Track A is the external-proof workflow. It explains how an operator verifies official submission paths, submits only when authorized, records target-originated evidence, and updates a specific consumer target safely.

Track B is the productization path. It makes the product easier to self-host,
adapt, and operate. Track A remains available later when an operator has
claim-specific external artifacts.

Rules:

- If real consumer artifacts arrive, update only the relevant Track A target record and `status.json`.
- Do not advance consumer status from Track B work alone.
- Product improvements do not imply compliance or acceptance.
- Agency-owned domain proof is separate from the existing DuckDNS OCI pilot evidence.

## Phase Sequence

| Phase | Name | Primary outcome |
| --- | --- | --- |
| 22 | Release And Distribution Hardening | Versioned releases, changelog, release checklist, install/upgrade/rollback docs. |
| 23 | Agency-Owned Deployment Proof | Stable agency-owned or agency-approved URL proof and evidence packet. |
| 24 | Real Agency Data Onboarding | Real GTFS import, validation triage, metadata, and publish workflow. |
| 25 | Device And AVL Integration Kit | Real device/vendor telemetry onboarding and simulator/supporting docs. |
| 26 | Admin UX Setup Wizard | Browser-guided setup and richer operator workflows. |
| 27 | Multi-Agency Isolation Prototype | Testable agency isolation and multi-agency deployment assumptions. |
| 28 | Production Operations Hardening | Longer-running operational proof, rotation, upgrade, incident, and backup practices. |
| 29 | Real-World Realtime Quality Expansion | Synthetic replay quality expansion and conservative ETA/matching evidence boundaries. |
| 29A | External Predictor Adapter Evaluation | Adapter contract and candidate-only feasibility review for optional external predictors. |
| 29B | AVL / Vendor Adapter Pilot Implementation | Synthetic vendor payload adapter pilot pattern behind the telemetry boundary. |
| 30 | Consumer Submission Execution | Authorized target submissions and retained evidence-based status changes. |
| 31 | Agency Pilot Program Package | Pilot onboarding kit, training, success criteria, support boundaries. |
| 32 | Public Launch And Ecosystem Outreach | Public messaging, launch copy, agency one-pager, contributor call-to-action. |
| 33 | Public GTFS Local/Pilot Evidence | Outcome C public-safe local/pilot evidence for a real public static GTFS dataset. |
| 34 | Post-Outcome-C Status Consistency And Evidence Readiness | Consistent post-Outcome-C docs, final-root request package, and public-GTFS repeatability guide. |
| 35 | README And Roadmap Realignment | Product-front-door README and self-hosted agency reuse continuation path. |
| 36 | OCI/OCL Reference Deployment Productization | Complete docs-only self-hosted reference deployment path based on the existing pilot server pattern. |
| 37 | Agency Reusable Onboarding Flow | Complete. Guided agency GTFS onboarding for local and reference-server evaluation. |
| 38 | Integration Adapter Kit | Complete. Central adapter map, synthetic fixture manifest, and dry-run adapter conformance tests. |
| 39 | CAL-ITP-Style Readiness Workflow | Complete. Product-facing readiness workflow without compliance overclaims. |

## Recommended Next Phase

Phase 39 is complete. The next recommended phase should continue self-hosted
agency reuse productization from the latest handoff.

Reason: the repo now has local app packaging, self-hosted reference deployment
docs, public-GTFS local/pilot proof, pilot operations helpers, a reusable
agency onboarding command, a central adapter kit for integration boundaries,
and an authenticated readiness workflow that makes CAL-ITP-style gaps visible
without claiming compliance. This does not claim agency adoption, consumer
acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, production
readiness, real vendor AVL compatibility, real realtime data, or
production-grade ETA-quality.

External-proof tracks remain available later when retained, redacted,
claim-specific artifacts exist. They are not the default Track B roadmap.

## Cross-Phase Constraints

Every Track B phase must preserve:

- existing public feed paths unless the phase explicitly handles a migration plan;
- admin auth and CSRF boundaries;
- evidence redaction rules;
- consumer status truthfulness;
- conservative realtime quality language;
- no committed secrets or private operator artifacts;
- no fake consumer or compliance evidence.

## Common Checks

Most Track B phases should start with:

```bash
make validate
make test
git diff --check
```

Run these when relevant:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
make demo-agency-flow
make agency-app-up
make agency-app-down
```

## Handoff Requirement

Each implemented Track B phase must add a handoff under `docs/handoffs/` using the repo handoff template.

The handoff should state:

- what was implemented;
- what was deferred;
- what changed in docs/scripts/runtime, if anything;
- what evidence or claims changed, if any;
- commands run;
- blocked commands;
- known remaining gaps;
- exact recommendation for the next phase.
