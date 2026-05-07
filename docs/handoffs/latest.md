# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness is complete
for the docs-only status consistency and evidence-readiness scope.

Phases 0 through 34 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, pilot-readiness, submission-readiness, public-messaging, or public-GTFS evidence issue directly requires it.

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

Phase 33 evidence is completed only for local/pilot public static GTFS dataset
handling. It does not prove agency adoption, agency approval, official agency
feed status, agency-owned final-root readiness, consumer evidence, compliance,
hosted SaaS, production readiness, real LA Metro realtime data, real vendor AVL
compatibility, or ETA quality.

Phase 34 aligned post-Outcome-C status and roadmap docs, documented the
original Phase 33 static validator blocker, recorded the later Homebrew Java 17
static validator retry result, added/verified final-root request and public-GTFS
repeatability guidance, and added this phase handoff. Phase 34 created no
external evidence and made no runtime code, script, Makefile target, schema,
migration, API, consumer tracker, final-root evidence packet, target artifact,
or OCI pilot final-root wording changes.

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
5. `docs/future-roadmap-post-outcome-c.md`
6. `docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md`
7. `docs/phase-33-public-gtfs-local-pilot-evidence.md`
8. `docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md`
9. `docs/final-root-operator-request.md`
10. `docs/tutorials/public-gtfs-local-pilot.md`
11. `docs/roadmap-status.md`
12. `docs/california-readiness-summary.md`
13. `docs/compliance-evidence-checklist.md`
14. `docs/agency-owned-domain-readiness.md`
15. `docs/evidence/consumer-submissions/status.json`
16. `docs/evidence/consumer-submissions/submission-workflow.md`
17. `docs/evidence/redaction-policy.md`
18. `SECURITY.md`
19. `README.md`
20. `docs/dependencies.md`
21. `docs/decisions.md`
22. `docs/handoffs/final-root-evidence-follow-up.md`
23. `docs/handoffs/phase-33.md`
24. `docs/handoffs/phase-32.md`
25. `docs/phase-32-public-launch-ecosystem-outreach.md`
26. `docs/agency-one-pager.md`
27. `docs/demo-video-outline.md`
28. `docs/public-share-copy.md`
29. `docs/ecosystem-positioning.md`
30. `docs/public-launch-checklist.md`
31. `docs/agency-pilot-program.md`
32. `docs/agency-pilot-checklist.md`
33. `docs/agency-feedback-template.md`

## Current Objective

Do not make stronger public claims than the retained evidence supports. The next useful work should target another concrete evidence gap: agency-owned/final-root proof, authorized target-specific consumer submission evidence, real agency pilot evidence, real deployment operations refresh, real device/vendor AVL evidence, or real-world realtime quality evidence.

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
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit Next

Choose files based on the evidence target selected next:

- Agency-owned/final-root proof: `docs/agency-owned-domain-readiness.md`, `docs/california-readiness-summary.md`, `docs/compliance-evidence-checklist.md`, and a future redacted evidence packet.
- Authorized target-specific consumer submission: `docs/evidence/consumer-submissions/submission-workflow.md`, the selected packet under `docs/evidence/consumer-submissions/packets/`, and only real redacted target-originated artifacts.
- Real agency pilot evidence: `docs/agency-pilot-program.md`, `docs/agency-pilot-checklist.md`, `docs/agency-feedback-template.md`, and a future public-safe agency pilot evidence packet.
- Real deployment operations evidence: `docs/runbooks/`, `docs/compliance-evidence-checklist.md`, and a future public-safe operations evidence packet.
- Real device/vendor AVL evidence: `docs/evidence/device-avl/`, `docs/tutorials/device-avl-integration.md`, and a future public-safe device/AVL evidence packet.
- Real-world realtime quality evidence: `internal/realtimequality`, `testdata/replay/`, `docs/phase-19-realtime-quality-eta-improvement.md`, and a future retained real-world quality evidence packet.
- Public GTFS local/pilot repeatability guidance only: `docs/tutorials/public-gtfs-local-pilot.md` and `docs/phase-33-public-gtfs-local-pilot-evidence.md`, unless a new evidence run is intentionally selected and can retain public-safe artifacts.

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

Candidate evidence work remains:

- agency-owned or agency-approved final-root proof;
- authorized target-specific consumer submission evidence;
- real agency pilot evidence;
- real deployment operations refresh;
- real device/vendor AVL evidence;
- real-world realtime quality evidence.
