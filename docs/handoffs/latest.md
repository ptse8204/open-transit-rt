# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 32 — Public Launch And Ecosystem Outreach is complete for the docs-only draft public launch materials scope.

Phases 0 through 32 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, pilot-readiness, submission-readiness, or public-messaging issue directly requires it.

Phase 32 produced draft public launch materials only. No announcement was posted, no social copy was published, no agency was contacted, no reporter was contacted, no consumer or aggregator was contacted, and no public launch occurred.

The post-Phase-32 final-root evidence follow-up is complete as
blocker-documented only. No agency-owned or agency-approved final public feed
root was available, no root was used, and no owner/approval evidence was
available. No DNS, TLS, redirect, public feed fetch, validator, proxy/config,
packet README, or checksum evidence was collected.

The recommended next roadmap step is to pause stronger public claims and pursue real retained evidence. Candidate next work should be one of: agency-owned or agency-approved final-root proof, authorized target-specific consumer submission evidence, real agency pilot evidence, or real deployment operations evidence.

## Phase 32 Summary

- Added `docs/agency-one-pager.md` with problem, solution, audience, current capabilities, pilot path, requirements, readiness boundaries, evidence boundaries, and agency next steps.
- Added `docs/demo-video-outline.md` with a truthful local demo script covering startup, GTFS import/demo feed, public feed URLs, Operations Console setup, device telemetry or dry-run adapter path, validation/evidence view, consumer packet boundary, and pilot next step.
- Added `docs/public-share-copy.md` with draft-only short, medium, and longer copy for GitHub launch, agency/evaluator, contributor, and transit/open-data audiences.
- Added `docs/ecosystem-positioning.md` covering GTFS/GTFS Realtime, validators, Caltrans/CAL-ITP-style readiness, downstream consumers and aggregators, agency-owned domains, external predictor adapters, AVL/vendor adapters, and open-source transit tooling.
- Added `docs/public-launch-checklist.md` with public-message safety checks, no-logo/no-affiliation rule, and claim-to-evidence table.
- Updated README, docs navigation, Phase 32 status, roadmap status, current status, and this latest handoff.
- Added `docs/handoffs/phase-32.md`.

## Truthfulness And Evidence Boundary

- All seven consumer and aggregator targets remain `prepared` only.
- Phase 30 selected no target and made no submissions.
- No target has submitted, under-review, accepted, rejected, blocked, ingestion, listing, display, or adoption evidence.
- No agency-owned or agency-approved final public feed root exists in repo evidence.
- The post-Phase-32 final-root evidence follow-up confirmed the final-root blocker remains unresolved and created no evidence packet.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- Phase 29A is adapter evaluation evidence only, not production ETA proof.
- Phase 29B is synthetic dry-run transform evidence only, not real vendor compatibility proof, production integration evidence, or AVL reliability evidence.
- Phase 31 is a docs-only pilot package and does not prove agency adoption.
- Phase 32 is draft launch materials only and does not prove launch, adoption, acceptance, compliance, endorsement, or production readiness.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, real-world ETA accuracy, production-grade ETA quality, certified hardware support, vendor compatibility, production AVL reliability, agency adoption, consumer submission, or public launch completion.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-32.md`
4. `docs/phase-32-public-launch-ecosystem-outreach.md`
5. `docs/agency-one-pager.md`
6. `docs/demo-video-outline.md`
7. `docs/public-share-copy.md`
8. `docs/ecosystem-positioning.md`
9. `docs/public-launch-checklist.md`
10. `docs/agency-pilot-program.md`
11. `docs/agency-pilot-checklist.md`
12. `docs/agency-feedback-template.md`
13. `docs/roadmap-status.md`
14. `docs/california-readiness-summary.md`
15. `docs/compliance-evidence-checklist.md`
16. `docs/agency-owned-domain-readiness.md`
17. `docs/evidence/consumer-submissions/status.json`
18. `docs/evidence/consumer-submissions/submission-workflow.md`
19. `docs/evidence/redaction-policy.md`
20. `SECURITY.md`
21. `README.md`
22. `docs/dependencies.md`
23. `docs/decisions.md`
24. `docs/handoffs/final-root-evidence-follow-up.md`

## Current Objective

Do not make stronger public claims until real retained evidence exists. The next useful work should target one concrete evidence gap: agency-owned/final-root proof, authorized target-specific consumer submission evidence, real agency pilot evidence, or real deployment operations evidence.

Consumer or aggregator submission work remains available only when a future operator is authorized, a target is selected, official target paths are verified, and target-originated evidence can be retained and redacted. Product improvements, validator success, pilot packaging, prepared packets, and draft launch materials alone must not advance target statuses.

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
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit Next

Choose files based on the evidence target selected next:

- Agency-owned/final-root proof: `docs/agency-owned-domain-readiness.md`, `docs/california-readiness-summary.md`, `docs/compliance-evidence-checklist.md`, and a future redacted evidence packet.
- Authorized target-specific consumer submission: `docs/evidence/consumer-submissions/submission-workflow.md`, the selected packet under `docs/evidence/consumer-submissions/packets/`, and only real redacted target-originated artifacts.
- Real agency pilot evidence: `docs/agency-pilot-program.md`, `docs/agency-pilot-checklist.md`, `docs/agency-feedback-template.md`, and a future public-safe agency pilot evidence packet.
- Real deployment operations evidence: `docs/runbooks/`, `docs/compliance-evidence-checklist.md`, and a future public-safe operations evidence packet.

Do not edit target-specific consumer records, `docs/evidence/consumer-submissions/status.json`, or artifact directories unless retained, redacted, target-originated evidence supports a target-specific status transition.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs unless explicitly approved and backed by current target documentation.
- Do not automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, production-grade ETA proven, real-world ETA-accuracy proven, certified hardware supported, vendor-compatible, agency-adopted, consumer-accepted, or publicly launched.

## Exact Next-Step Recommendation

Pause stronger public claims and pursue real retained evidence before any stronger launch, adoption, readiness, or consumer-status wording.

Recommended candidate next work should be one of:

- agency-owned or agency-approved final-root proof;
- authorized target-specific consumer submission evidence;
- real agency pilot evidence;
- real deployment operations evidence.
