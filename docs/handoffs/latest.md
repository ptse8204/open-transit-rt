# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 29 — Real-World Realtime Quality Expansion is complete for the synthetic replay evidence expansion scope.

Phases 0 through 29 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, realtime-quality, evidence, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, or submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 30 — Agency Pilot Evidence Refresh, with the evidence and claim boundaries below.

## Phase 29 Summary

- Added richer synthetic replay fixtures for after-midnight service, exact and non-exact frequency trips, block continuity, long layover withholding, sparse telemetry, noisy/off-shape GPS, stale/ambiguous hard patterns, cancellation alert linkage, and manual override before/after expiry.
- Added replay fixture support for `frequencies`.
- Added optional replay manual override `expires_at` support.
- Made the replay telemetry repository return the latest telemetry row per vehicle for feed snapshots, matching the production repository contract.
- Strengthened replay comparison for already-recorded cancellation alert linkage and unsupported disruption-withheld metrics.
- Added focused realtime-quality tests for Phase 29 scenarios.
- Updated replay fixture docs, Phase 29 docs, current status, and this handoff.

## Truthfulness And Evidence Boundary

- Phase 29 expands synthetic replay coverage only.
- No real-world observed-arrival/departure evidence exists in the repo for Phase 29.
- Do not claim real-world observed-arrival ETA accuracy or production-grade ETA quality.
- Real route/time-period quality metrics remain deferred because no real deployment or observed-arrival data exists.
- No public feed URLs changed.
- No GTFS-RT protobuf contracts changed.
- No consumer statuses changed.
- No auth boundaries changed.
- No external dependencies or predictor integrations changed.
- No TheTransitClock or other external predictor was integrated.
- No Operations Console surface changed.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, real-world ETA accuracy, or production-grade ETA quality.

All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.

The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-29.md`
4. `docs/phase-29-realtime-quality-expansion.md`
5. `docs/track-b-productization-roadmap.md`
6. `testdata/replay/README.md`
7. `docs/runbooks/production-operations-hardening.md`
8. `docs/evidence/redaction-policy.md`
9. `SECURITY.md`
10. `docs/prompts/calitp-truthfulness.md`
11. `docs/california-readiness-summary.md`
12. `docs/compliance-evidence-checklist.md`
13. `README.md`
14. `docs/dependencies.md`
15. `docs/decisions.md`

## Current Objective

Start Phase 30 only when maintainers are ready. Phase 30 should focus on evidence refresh or the next Track B objective without changing consumer statuses, public feed URLs, external integrations, or ETA claims unless retained, redacted evidence supports the change.

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

## Checks Run For Phase 29

- Pre-edit/planning `make validate` — passed.
- Pre-edit/planning `make realtime-quality` — passed.
- Pre-edit/planning `make test` — passed.
- Pre-edit/planning `make smoke` — passed.
- Pre-edit/planning `make test-integration` — passed.
- Pre-edit/planning `git diff --check` — passed.
- Pre-edit/planning `docker compose -f deploy/docker-compose.yml config` — passed.
- Implementation focused `go test ./internal/realtimequality` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make test` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 27 proves selected repository-level isolation paths with synthetic data only. It does not prove hosted multi-tenant production readiness, one-instance multi-agency public feed roots, or tenant-safe backup/restore/export/evidence operations.
- Phase 28 templates are templates only and are not evidence by themselves.
- Phase 29 replay fixtures are synthetic evidence only and are not real-world ETA accuracy evidence.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 30

- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/track-b-productization-roadmap.md`
- evidence docs or templates selected by maintainers for Phase 30

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
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

Start Phase 30 with an evidence-refresh or adoption-readiness slice that preserves Phase 29’s synthetic-only realtime quality boundary. Recommended first slice: refresh the local replay and validation evidence index, identify which evidence remains synthetic/template-only, and document the exact real-world inputs required before stronger production ETA or route/time-period quality claims can be made.
