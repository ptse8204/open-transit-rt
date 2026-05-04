# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 29A — External Predictor Adapter Evaluation is complete for the adapter contract documentation, candidate-only TheTransitClock feasibility review, and test-only mock adapter contract checks.

Phases 0 through 29A are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, or submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 29B — AVL / Vendor Adapter Pilot Implementation. Phase 30 — Consumer Submission Execution remains later and must not advance statuses without retained, redacted, target-originated evidence.

## Phase 29A Summary

- Documented the external predictor adapter contract, including inputs, outputs, diagnostics, failure modes, and strict wrong-agency/wrong-feed handling.
- Added Trip Updates adapter output validation before normalization/protobuf serialization.
- Added test-only mock external adapter coverage for happy-path normalization and diagnostics persistence.
- Added test-only rejection coverage for missing active-feed trips, impossible stop sequences, stale prediction timestamps, wrong agency/feed candidates, unsupported added-trip predictions, and low/missing confidence.
- Documented Vehicle Positions independence from external predictor availability.
- Reviewed TheTransitClock as candidate-only from public sources on 2026-05-04.
- Added `docs/handoffs/phase-29a.md`.

## Truthfulness And Evidence Boundary

- Phase 29A is adapter contract and candidate-feasibility evidence only.
- No TheTransitClock or external predictor runtime integration exists.
- No runtime external predictor config, service client, network call, subprocess call, Java/Maven/Tomcat invocation, or external service requirement was added.
- Vehicle Positions generation remains independent of external predictor availability.
- Public-source TheTransitClock review is not runtime compatibility proof.
- Do not claim better ETAs, real-world observed-arrival ETA accuracy, real-world predictor compatibility, or production-grade ETA quality.
- Real route/time-period quality metrics remain deferred because no real deployment or observed-arrival data exists.
- No public feed URLs changed.
- No GTFS-RT protobuf contracts changed.
- No consumer statuses changed.
- No auth boundaries changed.
- No database schema changed.
- No runtime external dependencies or predictor integrations changed.
- No Operations Console surface changed.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, real-world ETA accuracy, or production-grade ETA quality.

All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.

The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-29a.md`
4. `docs/phase-29a-external-predictor-adapter-evaluation.md`
5. `docs/track-b-productization-roadmap.md`
6. `docs/phase-29b-avl-vendor-adapter-pilot.md`
7. `docs/tutorials/device-avl-integration.md`
8. `docs/evidence/redaction-policy.md`
9. `SECURITY.md`
10. `docs/prompts/calitp-truthfulness.md`
11. `docs/california-readiness-summary.md`
12. `docs/compliance-evidence-checklist.md`
13. `README.md`
14. `docs/dependencies.md`
15. `docs/decisions.md`

## Current Objective

Start Phase 29B only when maintainers are ready. Phase 29B should implement a synthetic AVL/vendor adapter pilot pattern behind the existing telemetry boundary without adding real vendor data, credentials, named vendor runtime dependencies, public feed URL changes, consumer-status changes, or unsupported vendor-support claims. Phase 30 consumer submission execution remains later and must not advance consumer statuses without retained, redacted, target-originated evidence.

## Exact First Commands

```bash
make validate
make realtime-quality
make test
make smoke
make test-integration
git diff --check
docker compose -f deploy/docker-compose.yml config
```

## Checks Run For Phase 29A

- Focused `go test ./internal/prediction ./internal/feed/tripupdates ./internal/realtimequality` — passed.
- `make validate` — passed.
- `make realtime-quality` — passed.
- `make test` — passed.
- `make smoke` — passed.
- `make test-integration` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 27 proves selected repository-level isolation paths with synthetic data only. It does not prove hosted multi-tenant production readiness, one-instance multi-agency public feed roots, or tenant-safe backup/restore/export/evidence operations.
- Phase 28 templates are templates only and are not evidence by themselves.
- Phase 29 replay fixtures are synthetic evidence only and are not real-world ETA accuracy evidence.
- Phase 29A mock adapter tests are contract tests only and are not real-world predictor compatibility or ETA-quality evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 29B

- `docs/phase-29b-avl-vendor-adapter-pilot.md`
- `docs/handoffs/phase-29b.md`
- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/track-b-productization-roadmap.md`
- optional telemetry adapter test-only helpers selected by maintainers for Phase 29B

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Keep Vehicle Positions independent of external predictor availability.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve unknown/withheld/degraded/stale/ambiguous visibility in realtime diagnostics.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Keep `/v1/events` documented and operated as an authenticated admin/debug review path, not as a public or consumer-facing feed.
- Do not add consumer submission APIs, automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, production-grade ETA proven, or real-world ETA-accuracy proven.

## Exact Next-Step Recommendation

Start Phase 29B — AVL / Vendor Adapter Pilot Implementation. Keep vendor integrations behind the telemetry boundary, synthetic/test-only unless approved evidence exists, and truthfully described; do not add a named vendor runtime dependency, real private AVL data, credentials, public feed URL changes, consumer-status changes, or stronger vendor-support claims. Phase 30 consumer submission execution remains later and must not advance any target beyond `prepared` without target-originated evidence.
