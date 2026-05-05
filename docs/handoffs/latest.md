# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 31 — Agency Pilot Program Package is complete for the docs-only agency
pilot package scope.

Phases 0 through 31 are closed for their documented scopes. Track A is also
closed for its docs-only external-proof workflow scope. Do not reopen earlier
phases unless a blocking truthfulness, safety, security, realtime-quality,
evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL
onboarding, admin-UX, operations-hardening, pilot-readiness, or
submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 32 — Public Launch And
Ecosystem Outreach. Phase 32 must proceed from the prepared-only consumer state
and must not assume agency adoption, consumer submission, review, acceptance,
rejection, blocker, ingestion, listing, display, CAL-ITP/Caltrans compliance,
hosted SaaS availability, paid support/SLA coverage, or production readiness
evidence exists.

## Phase 31 Summary

- Added `docs/agency-pilot-program.md` with pilot overview, non-goals,
  suggested non-SLA timeline, responsibilities, evidence boundaries, consumer
  submission boundary, success criteria, failure/blocker criteria, risk
  register, and closeout summary.
- Added `docs/agency-pilot-kickoff-agenda.md` with attendees, pre-kickoff
  preparation, 30-minute and 60-minute agendas, walkthrough topics, decisions,
  follow-up actions, and what not to collect.
- Added `docs/agency-pilot-checklist.md` with data prerequisites, GTFS
  ownership, metadata, domain/DNS/TLS, telemetry/device, validators,
  operations, security/redaction, consumer submission, staff roles,
  responsibility matrix, launch/readiness review, and exit criteria.
- Added `docs/agency-training-outline.md` with training topics for GTFS, GTFS
  Realtime, local demo, real GTFS onboarding, validation triage, GTFS Studio,
  device tokens, AVL/vendor boundary, Operations Console, evidence, consumer
  submissions, support, and security reporting.
- Added `docs/agency-feedback-template.md` with public-safe feedback prompts and
  claim-boundary review.
- Updated Phase 31 status, current status, Track B roadmap next phase, docs
  navigation, README navigation, and this latest handoff.

## Phase 30 Outcome B Summary

- Phase 30 closed as Outcome B — blocker-documented closure only.
- No authorized submission, official-path verification evidence, or
  target-originated artifact was available.
- No Phase 30 target was selected.
- Target selection is deferred until an operator is authorized and either
  official-path verification or target-originated evidence can be retained.
- No individual target status changed to `blocked` because no target-specific
  blocker artifact exists.
- `docs/evidence/consumer-submissions/status.json` and all current target
  records under `docs/evidence/consumer-submissions/current/` were left
  unchanged.
- Artifact directories remain README-only; no receipts, screenshots, tickets,
  correspondence, blocker notes, or placeholder artifacts were added.

## Truthfulness And Evidence Boundary

- All seven consumer and aggregator targets remain `prepared` only.
- No target has submitted, under-review, accepted, rejected, blocked, ingestion,
  listing, display, or adoption evidence.
- No agency-owned or agency-approved final public feed root exists in repo
  evidence.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable
  URL/domain proof.
- Phase 29B is synthetic adapter pattern and dry-run transform evidence only,
  not real vendor compatibility proof, production integration evidence, or AVL
  reliability evidence.
- Phase 29A is adapter evaluation evidence only, not production ETA proof.
- Phase 28 runbooks/templates are not real operations evidence by themselves.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal
production readiness, production multi-tenant hosting, consumer acceptance,
CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor
equivalence, real-world ETA accuracy, production-grade ETA quality, certified
hardware support, vendor compatibility, production AVL reliability, or agency
adoption.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-31.md`
4. `docs/phase-31-agency-pilot-program-package.md`
5. `docs/agency-pilot-program.md`
6. `docs/agency-pilot-checklist.md`
7. `docs/agency-pilot-kickoff-agenda.md`
8. `docs/agency-training-outline.md`
9. `docs/agency-feedback-template.md`
10. `docs/evidence/consumer-submissions/submission-workflow.md`
11. `docs/evidence/consumer-submissions/status.json`
12. `docs/evidence/redaction-policy.md`
13. `SECURITY.md`
14. `docs/california-readiness-summary.md`
15. `docs/compliance-evidence-checklist.md`
16. `README.md`
17. `docs/dependencies.md`
18. `docs/decisions.md`

## Current Objective

Start Phase 32 from the Phase 31 pilot package and prepared-only consumer state.
Phase 32 should focus on truthful public messaging, agency evaluation framing,
and contributor/community outreach materials.

Consumer or aggregator submission work remains available only when a future
operator is authorized, a target is selected, official target paths are
verified, and target-originated evidence can be retained and redacted. Product
improvements, validator success, pilot packaging, or prepared packets alone must
not advance target statuses.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

Run these when Phase 32 work touches relevant surfaces:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

## Checks Run For Phase 31

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
- Post-edit `git diff --check` — initially found one extra blank line at EOF in this file; fixed, then final rerun passed.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/`
  remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or
  agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real
  agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real
  public-safe device or AVL integration evidence exists.
- Phase 29B synthetic fixtures are not real vendor AVL evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not
  submissions.
- Phase 30 did not select a target, verify an official path, submit a packet,
  add artifacts, or change consumer statuses.
- Phase 31 did not add real pilot evidence, consumer evidence, agency adoption
  evidence, operations evidence, final-root proof, or device/AVL proof.
- Consumer-ingestion workflow records and docs tracker records are not
  third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin
  tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs,
  notification credentials, raw telemetry payloads, unredacted correspondence,
  private portal credentials, private ticket links, raw logs with credentials,
  private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 32

- `docs/phase-32-public-launch-and-ecosystem-outreach.md`, if added.
- `README.md`
- `docs/README.md`
- `wiki/`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-32.md`

Do not edit target-specific consumer records,
`docs/evidence/consumer-submissions/status.json`, or artifact directories unless
retained, redacted, target-originated evidence supports a target-specific status
transition.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs unless explicitly approved and backed by
  current target documentation.
- Do not automate submissions, contact external portals, guess submission paths,
  or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed,
  agency-endorsed, marketplace/vendor equivalent, universally production ready,
  production multi-tenant hosted, production-grade ETA proven, real-world
  ETA-accuracy proven, certified hardware supported, vendor-compatible, or
  agency-adopted.

## Exact Next-Step Recommendation

Start Phase 32 — Public Launch And Ecosystem Outreach as truthful public
messaging and ecosystem-facing documentation only. Do not assume agency adoption,
consumer acceptance, compliance, hosted SaaS, paid support/SLA, or production
readiness evidence exists.
