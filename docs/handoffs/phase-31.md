# Phase 31 Handoff

## Phase

Phase 31 — Agency Pilot Program Package

## Status

- Complete for the docs-only agency pilot package scope.
- Active phase after this handoff: Phase 32 — Public Launch And Ecosystem Outreach.

## What Was Implemented

- Added an agency pilot overview with pilot goals, explicit non-goals, evidence boundaries, consumer submission boundaries, suggested non-SLA timeline, success criteria, failure/blocker criteria, risk register, and closeout summary.
- Added a kickoff agenda with attendee guidance, pre-kickoff preparation, 30-minute and 60-minute agendas, walkthrough topics, decisions, follow-up actions, and what not to collect.
- Added an agency pilot checklist with data prerequisites, GTFS ownership, metadata, domain/DNS/TLS, telemetry/device, validator, operations, security/redaction, consumer submission, staff-role, public-launch-readiness, and exit criteria sections.
- Added a responsibility matrix covering GTFS ownership, domain/DNS/TLS, device/AVL credentials, operations/backups, consumer submissions, incident response, and support expectations.
- Added a training outline for GTFS, GTFS Realtime, local demo flow, real GTFS onboarding, validation triage, GTFS Studio basics, device token safety, AVL/vendor adapter boundaries, Operations Console setup, validation/evidence, consumer submission boundaries, support, and security reporting.
- Added a public-safe agency feedback template.
- Updated Phase 31 status, current status, latest handoff guidance, Track B roadmap next phase, and documentation navigation.

## What Was Intentionally Deferred

- No backend features.
- No runtime integrations.
- No database schema changes.
- No public feed URL changes.
- No GTFS-RT protobuf contract changes.
- No consumer status changes.
- No evidence artifacts.
- No submissions, external contacts, portal automation, or official-path verification.
- No legal, procurement, paid support, SLA, hosted SaaS, agency endorsement, consumer acceptance, production-readiness, or CAL-ITP/Caltrans compliance commitments or claims.

## Pilot Docs Added

- `docs/agency-pilot-program.md`
- `docs/agency-pilot-kickoff-agenda.md`

## Checklist, Training, And Feedback Docs Added

- `docs/agency-pilot-checklist.md`
- `docs/agency-training-outline.md`
- `docs/agency-feedback-template.md`

## Success And Failure Criteria

- Pilot success criteria are evaluation criteria only: GTFS import or blocker review, public schedule review, metadata review, validation review, telemetry/device or simulator review, Vehicle Positions review, Trip Updates diagnostics review, Alerts path review, Operations Console understanding, operations runbook review, evidence or blocker documentation, operator usability, and support-load understanding.
- Failure/blocker criteria include missing agency permission, no public-safe GTFS, no operator owner, no metadata or domain plan, no device/AVL/simulator path, unresolved validation blockers, security/redaction concerns, no staff to operate, unsupported support expectations, or unauthorized consumer submission requests.

## Risk Register

- Added in `docs/agency-pilot-program.md`.
- Covers data ownership, private data leakage, secret leakage, unstable public URL, GTFS validation failure, device/AVL reliability, Trip Updates quality, operations capacity, consumer submission delay, support expectation, and multi-agency/tenant-boundary risk.
- Each risk includes description, likelihood, impact, mitigation, owner, and evidence needed.

## Support Boundary Summary

- Maintainer/community help remains best-effort.
- Agencies/operators own GTFS permission, metadata approval, DNS/TLS/hosting, runtime secrets, device/AVL credentials, backups, monitoring, incident response, consumer submissions, and continued operations.
- Phase 31 does not create paid support, response targets, SLA coverage, hosted SaaS availability, legal commitments, procurement commitments, or production operating guarantees.

## Public Launch Readiness Checklist

- Added in `docs/agency-pilot-checklist.md`.
- Requires agency permission, approved wording, no private data, pilot-vs-production wording, no compliance overclaim, no consumer overclaim, current blocker list, redaction review, and support boundary review.
- The checklist does not approve public launch, agency endorsement, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or production readiness. It only helps decide whether public messaging is safe and truthful.

## Schema And Interface Changes

- None.

## Dependency Changes

- None.

## Migrations Added

- None.

## Tests Added And Results

- No code tests were added because Phase 31 is docs-only.
- Required checks and targeted scans are recorded below.

## Commands Run

- Pre-implementation `make validate` — passed.
- Pre-implementation `make test` — passed.
- Pre-implementation `git diff --check` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make demo-agency-flow` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `python3 -m json.tool docs/evidence/consumer-submissions/status.json` — passed.
- Post-edit read-only consumer tracker status check — passed; all seven targets remain `prepared`.
- Post-edit target tracker/artifact diff check — passed; `status.json`, current target records, and artifact directories were not edited.
- Post-edit secret-like value scan — passed with no matches.
- Post-edit context-aware forbidden-claim scan — reviewed; matches are negative/boundary wording, previous phase history, or required claim-boundary language.
- Post-edit redaction-sensitive term scan — reviewed; matches are "do not collect", "do not commit", support boundary, security, and redaction rules.
- Post-edit `git diff --check` — initially found one extra blank line at EOF in `docs/handoffs/latest.md`; fixed, then final rerun passed.

## Blocked Commands

- None.

## Known Remaining Agency Pilot Gaps

- No real agency pilot evidence has been collected.
- No agency-owned or agency-approved final feed root exists in repo evidence.
- No real agency-approved GTFS import evidence exists in repo evidence.
- No real device, hardware, or vendor AVL integration evidence exists in repo evidence.
- Phase 29B remains synthetic dry-run adapter evidence only.
- Phase 29A remains adapter evaluation evidence only, not production ETA proof.
- Phase 28 runbooks/templates are not real operations evidence by themselves.
- All seven consumer and aggregator targets remain `prepared`; no target has submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption evidence.

## Exact Recommendation For Phase 32

Proceed to Phase 32 — Public Launch And Ecosystem Outreach only for truthful public messaging, agency evaluation framing, and contributor/community outreach materials.

Phase 32 must not assume agency adoption, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, paid support/SLA coverage, production readiness, marketplace/vendor equivalence, production multi-tenant hosting, or production-grade ETA evidence exists.
