# Roadmap Status

This page gives a public-readable status summary without requiring readers to understand every phase handoff.

It does not claim CAL-ITP/Caltrans compliance, consumer acceptance, agency endorsement, hosted SaaS availability, paid support, SLA coverage, marketplace/vendor equivalence, production-grade ETA quality, or universal production readiness.

![Illustrative evidence maturity ladder from code exists to hosted evidence, prepared packet, submitted, under review, and accepted.](assets/evidence-maturity-ladder.png)

## What Works Today

Open Transit RT has technical foundations for:

- importing static GTFS ZIP files;
- editing typed GTFS Studio drafts;
- publishing public schedule and GTFS Realtime feed paths;
- ingesting authenticated vehicle telemetry;
- conservative trip matching;
- Vehicle Positions publication;
- Trip Updates behind a prediction adapter;
- basic Alerts authoring and publication;
- validation records and scorecard workflows;
- local app packaging and pilot operations helpers;
- reusable agency GTFS onboarding from an agency ID and GTFS URL;
- an integration adapter kit that maps telemetry, AVL, predictor, validator,
  monitoring, and consumer workflow boundaries;
- an authenticated CAL-ITP-style readiness workflow that shows readiness gaps,
  status sources, and next actions without compliance overclaims;
- a guided self-hosted operator trial that ties reference deployment,
  reusable onboarding, readiness review, validators, and the synthetic AVL
  dry-run together without creating evidence;
- operator smoke checks and redaction-safe support bundles for local/reference
  diagnostics without creating evidence;
- a synthetic telemetry simulator that uses real device-token auth and
  `POST /v1/telemetry` for local/reference diagnostics without creating
  evidence;
- a private synthetic AVL adapter send mode that posts transformed records to
  authenticated `POST /v1/telemetry` with env-referenced tokens and redacted
  `.cache/` diagnostics only;
- a local compliance/readiness packet generator and audit that writes ignored
  `.cache` summaries only;
- consumer packet preparation workflows.

## Evidence That Exists

Current evidence includes:

- local demo and validation workflows;
- hosted/operator evidence for the OCI pilot;
- Phase 33 Outcome C public GTFS local/pilot evidence for the May 6, 2026 LA
  Metro Bus GTFS run, with public-safe retained summaries;
- replay fixtures and metrics that measure current realtime behavior;
- prepared consumer and aggregator packet drafts for seven targets;
- an operator workflow for official-path verification, pre-submission checks,
  evidence intake, and artifact redaction;
- draft public launch materials for review;
- redaction and evidence policies.

OCI pilot evidence is useful pilot evidence. It is not agency-owned production proof.

## What Remains Missing

The repo does not currently have:

- third-party consumer submission, review, acceptance, rejection, or blocker evidence;
- agency-owned stable URL/domain proof for the OCI pilot feed set;
- validator-clean or no-warning static GTFS evidence for the Phase 33
  public-GTFS packet; a post-Phase-34 retry executed and reported 3 warning
  notices;
- production-grade ETA quality evidence;
- full hosted multi-tenant implementation;
- paid support or SLA commitments;
- marketplace/vendor-equivalent service packaging.

## How To Interpret Prepared Packets

Prepared means a packet draft exists for operator review. Prepared does not mean submitted, under review, accepted, rejected, listed, ingested, or compliant.

Target status can move beyond `prepared` only when retained, redacted, target-originated evidence exists for that specific target and feed scope.

Artifact directories for target-originated evidence must remain README-only
until real redacted evidence exists.

## How To Interpret Replay Metrics

Phase 19 replay metrics measure current behavior against committed fixtures. They help catch regressions and document conservative handling of stale, ambiguous, and withheld cases.

Replay metrics do not prove production-grade ETA quality or consumer acceptance.

## Future Roadmap Meaning

Future phases describe intended work. They are not commitments to hosted service, support coverage, agency endorsement, or consumer acceptance.

Track B is the planned productization path for release, agency onboarding,
device/AVL integration guidance, setup UX, multi-agency isolation proof,
operations hardening, realtime quality expansion, optional external predictor
adapter evaluation, AVL/vendor adapter pilot work, and truthful public
positioning.

Phase 29A — External Predictor Adapter Evaluation is complete for contract and
candidate-only feasibility review. Phase 29B — AVL / Vendor Adapter Pilot
Implementation is complete for the synthetic dry-run adapter pilot scope. Phase
30 — Consumer Submission Execution closed as Outcome B — blocker-documented
closure only because no authorized submission, official-path verification
evidence, or target-originated artifact was available. Phase 31 — Agency Pilot
Program Package is complete for the docs-only pilot package scope. Phase 32 —
Public Launch And Ecosystem Outreach is complete for draft public launch
materials only. Phase 32 did not post announcements, contact agencies, contact
consumers, or complete a public launch. Phase 33 — Public GTFS Local/Pilot
Evidence is complete as Outcome C — public-GTFS local/pilot run completed with
public-safe retained summaries for the May 6, 2026 LA Metro Bus GTFS local
attempt. The packet proves local/pilot handling of a real public static GTFS
dataset only; it does not prove agency adoption, final-root readiness, consumer
evidence, compliance, production readiness, real realtime data, or ETA quality.
Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness is complete
for docs-only status alignment, final-root request guidance, and public-GTFS
repeatability guidance; it created no new external evidence. Phase 35 — README
And Roadmap Realignment is complete for docs-only README and roadmap
realignment; it makes self-hosted agency reuse and OCI/OCL reference deployment
productization the default continuation path. Phase 36 — OCI/OCL Reference
Deployment Productization is complete for docs-only self-hosted reference
deployment guidance and created no external evidence. Phase 37 — Reusable
Agency Onboarding Flow is complete for the opt-in local/reference onboarding
command and docs. Phase 38 — Integration Adapter Kit is complete for the
central adapter map, synthetic fixture manifest, dry-run AVL examples, and
focused adapter conformance tests. Phase 39 — CAL-ITP-Style Readiness
Workflow is complete for the authenticated Operations Console readiness page
and documentation/navigation updates. Phase 40 — Guided Self-Hosted Operator
Trial is complete for the docs/navigation trial checklist that ties Phase 36
through Phase 39 workflows into one local/reference operator path without
creating evidence or stronger claims. Phase 41 — Operator Smoke And Support
Bundle is complete for repeatable local/reference smoke checks and
redaction-safe support bundles without creating evidence or stronger claims.
Phase 42 — Reference Deployment Doctor is complete for read-only OCI/OCL-style
reference deployment diagnostics without creating evidence or stronger claims.
Phase 43 — Operator UX Setup V2 is complete for the private authenticated
checklist and local routing patch scope. Phase 44 — Telemetry Simulator And
Device Trial is complete for the synthetic-only simulator that posts through
real authenticated telemetry ingest and optionally runs private DB-backed
matcher/Vehicle Positions diagnostics without creating evidence or stronger
claims. Phase 45 — GTFS Quality Triage Loop, Phase 46 — Validator Automation
And Health Gates, and Phase 47 — Self-Hosted Operations Notifications are
complete for private diagnostics and local/reference workflows without
creating evidence or stronger claims. Phase 48 — AVL Adapter Runtime Path is
complete for private send-mode execution through authenticated
`POST /v1/telemetry` with env-referenced tokens and redacted private
diagnostics only. Phase 49 — External Predictor Runtime Adapter is complete
for the optional disabled-by-default generic HTTP sidecar boundary behind
`internal/prediction.Adapter`, with sanitized DTOs, strict URL/token-env
validation, redacted diagnostics, and no named predictor, compatibility, or
ETA-quality claim. Phase 50 — Realtime Quality Backtesting is complete for
private aggregate diagnostics. Phase 51 — Operations Reliability And SLO
Readiness is complete for private authenticated reliability summaries,
bounded Vehicle Positions health persistence, and `.cache` reliability output
without evidence writes or stronger operational claims. Phase 52 — Final
Public Root Evidence Workflow closed blocker-only. Phase 53 — Authorized
Consumer Submission Execution closed blocker-only. Phase 54 — Official
Requirements Refresh updated official-source mappings only. Phase 55 —
Compliance Evidence Packet Generator is complete for ignored `.cache`
blocker/draft packet generation and audit without retained evidence, live feed
fetching, consumer contact, consumer status changes, or compliance claims.
Phase 56 — Multi-Agency Hosting Hardening, Phase 57 — Release Packaging And
Supply Chain, Phase 58 — Optional Marketplace / Vendor-Equivalent Pack, Phase
59 — Real Pilot Closeout, and Phase 60 — Final Claim Review And Public
Closeout are complete for their bounded scopes. Phase 60 added local claim
audit tooling and a final claim-to-evidence review only; it created no retained
evidence, changed no consumer statuses, and added no stronger claim.
Phase 61 — Agency-First UI And Connector Hub is complete for its bounded
product scope. It adopted the Phase 61+ roadmap naming, added the private
Connector Hub and agency-first Operations Console dashboard cards, and kept
connector/plugin wording limited to optional sidecars, manifests, command
adapters, or connector processes. Phase 61 created no evidence, changed no
consumer statuses, added no dynamic backend plugin loading, and made no
stronger claim.
Phase 62 — Guided Setup And Browser GTFS Import is complete for its bounded
private Operations Console scope. It added the setup wizard and admin-only
browser GTFS import flow without public routes, migrations, retained evidence,
consumer status changes, or stronger claims.
Phase 63 — Feed Health And Readiness UX is complete for its bounded private
Operations Console scope. It added the feed-health dashboard and readiness
checklist v2 for existing private signals without public routes, schema
changes, retained evidence, consumer status changes, or stronger claims.
Phase 64 — Connector Platform And SDKs is complete. Checkpoint 000001 added the
scoped connector platform and SDK plan. Checkpoint 000002 added the private
Connector Hub manifest registry UI from safe committed synthetic examples.
Checkpoint 000003 added private generated connector test instructions without
backend command execution. Checkpoint 000004 improved telemetry connector
SDK-style examples with a shared synthetic dry-run helper. Checkpoint 000005
improved the prediction connector SDK-style example with sanitized dry-run
request/response helpers and no-public-mutation diagnostics. Checkpoint
000006 improved monitoring/export examples and added a validator allowlist
example without dynamic backend plugin loading, evidence writes, consumer
status changes, raw validator commands, SLA/uptime claims, named predictor
compatibility claims, or production-grade ETA claims. Checkpoint 000007 closed
Phase 64 with validation, protected-path, consumer-tracker, and claim-boundary
review.
Phase 65 — Operator Workflow And Data Quality UX is complete. Checkpoint 000001
added the scoped plan for device/vehicle onboarding UI, telemetry simulator UI,
GTFS quality fix guidance, and closeout. Checkpoint 000002 improved the
private device and vehicle onboarding UI without changing device credential
semantics, telemetry ingest contracts, public routes, evidence paths, consumer
statuses, or stronger claims. Checkpoint 000003 added a private GET-only
telemetry simulator guide UI and JSON export without command execution, token
collection, `.cache` reads, telemetry sends, evidence writes, consumer status
changes, or stronger claims. Checkpoint 000004 improved private GTFS quality
guidance without automatic GTFS edits, draft mutation, schedule publish,
validator semantic changes, evidence writes, consumer status changes, or
stronger claims. Checkpoint 000005 closed Phase 65 with validation,
protected-path, consumer-tracker, and claim-boundary review.

Phase 66 — Release Candidate And Installability is active. Checkpoint 000001
added the release-candidate and installability plan. Checkpoint 000002
prepared the first release-candidate workflow with an ordered review sequence,
validation matrix, release-note inputs, local package audit matrix, and private
diagnostic summary fields. The next checkpoint is `Phase 66 -- Checkpoint
000003: improve installer and bootstrap UX`.

External-proof tracks such as agency-owned/final-root proof,
authorized target-specific consumer submission evidence, real agency pilot
evidence, real device/vendor AVL evidence, or real-world realtime quality
evidence remain future optional paths when retained claim-specific artifacts are
available. Do not make stronger public claims from Phase 33 beyond local/pilot
public-GTFS dataset handling.

Use `docs/roadmaps/agency-first-connector-platform/README.md` for the Phase
61+ forward roadmap, `docs/track-b-productization-roadmap.md` for historical
Track B context, `docs/roadmap-post-phase-14.md` for older post-Phase-14
context, and `docs/handoffs/latest.md` for the current handoff state.
