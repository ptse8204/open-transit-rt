# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 28 — Production Operations Hardening is complete for the docs-first operations hardening scope.

Phases 0 through 28 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, security, agency-boundary, auth, data-isolation, agency-domain, device/AVL onboarding, admin-UX, operations-hardening, or submission-readiness issue directly requires it.

The recommended next implementation phase is Phase 29 — Real-World Realtime Quality Expansion.

## Phase 28 Summary

- Added `docs/runbooks/production-operations-hardening.md`.
- Added template-only runbook records under `docs/runbooks/templates/`.
- Added incident/response templates for feed outage, validator failure, telemetry staleness, Trip Updates quality, secret exposure, consumer complaint or rejection, and restore events.
- Added secret rotation and operator handover templates.
- Added alert delivery proof guidance without requiring hosted monitoring SaaS or a full Prometheus/Grafana stack.
- Added capacity guidance for disk, database growth, backup storage, logs, and evidence artifacts.
- Hardened backup/restore, monitoring/alerting, validator evidence, deployment evidence, release checklist, and upgrade/rollback docs.
- Preserved Phase 27 language that current backup/restore/export/evidence workflows are deployment/DB scoped and are not tenant-safe multi-agency workflows.
- Updated `docs/phase-28-production-operations-hardening.md`, `docs/current-status.md`, docs navigation, and `CHANGELOG.md`.

## Truthfulness And Evidence Boundary

- No runtime APIs changed.
- No database schema changed.
- No public feed URLs changed.
- No GTFS-RT protobuf contracts changed.
- No consumer statuses changed.
- No external integrations changed.
- No systemd or Docker behavior changed.
- No evidence claims changed.
- No fake incidents, fake outage evidence, fake alert delivery proof, fake rotation records, fake restore events, or placeholder operational artifacts were added.

Do not claim hosted SaaS availability, paid support/SLA coverage, universal production readiness, production multi-tenant hosting, consumer acceptance, CAL-ITP/Caltrans compliance, agency endorsement, marketplace/vendor equivalence, or production-grade ETA quality.

All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.

The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-28.md`
4. `docs/phase-29-realtime-quality-expansion.md`
5. `docs/track-b-productization-roadmap.md`
6. `docs/runbooks/production-operations-hardening.md`
7. `docs/runbooks/backup-and-restore.md`
8. `docs/runbooks/monitoring-and-alerting.md`
9. `docs/runbooks/validator-evidence.md`
10. `docs/runbooks/deployment-evidence-overview.md`
11. `docs/evidence/redaction-policy.md`
12. `SECURITY.md`
13. `docs/prompts/calitp-truthfulness.md`
14. `docs/multi-agency-strategy.md`
15. `docs/agency-owned-domain-readiness.md`
16. `docs/california-readiness-summary.md`
17. `docs/compliance-evidence-checklist.md`
18. `README.md`
19. `docs/dependencies.md`
20. `docs/decisions.md`

## Current Objective

Start Phase 29 — Real-World Realtime Quality Expansion when maintainers are ready to continue Track B implementation.

Phase 29 should add richer quality fixtures and repeatable realtime metrics where evidence supports them. Preserve conservative uncertainty, keep Trip Updates pluggable, and do not claim production-grade ETA quality from replay-only evidence.

## Exact First Commands

```bash
make validate
make realtime-quality
make test
make smoke
make test-integration
git diff --check
```

## Checks Run For Phase 28

- Pre-edit/planning `make validate` — passed.
- Pre-edit/planning `make test` — passed.
- Pre-edit/planning `make test-integration` — passed.
- Pre-edit/planning `make realtime-quality` — passed.
- Pre-edit/planning `make smoke` — passed.
- Pre-edit/planning `docker compose -f deploy/docker-compose.yml config` — passed.
- Pre-edit/planning `git diff --check` — passed.
- Post-edit `make validate` — passed.
- Post-edit `make test` — passed.
- Post-edit `make test-integration` — passed.
- Post-edit `make realtime-quality` — passed.
- Post-edit `make smoke` — passed.
- Post-edit `docker compose -f deploy/docker-compose.yml config` — passed.
- Post-edit `git diff --check` — passed.
- Post-edit targeted redaction and forbidden-claim scan — passed.
- Blocked commands — none.

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 23 did not create final-root evidence. No agency-owned or agency-approved final public feed root is available in repo evidence.
- Phase 24 real-agency GTFS evidence scaffolding is template-only until real agency-approved, public-safe evidence exists.
- Phase 25 device/AVL evidence scaffolding is template-only until real public-safe device or AVL integration evidence exists.
- Phase 27 proves selected repository-level isolation paths with synthetic data only. It does not prove hosted multi-tenant production readiness, one-instance multi-agency public feed roots, or tenant-safe backup/restore/export/evidence operations.
- Phase 28 templates are templates only and are not evidence by themselves.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, raw logs with credentials, private backup paths, or raw private operator artifacts.

## First Files Likely To Edit For Phase 29

- `internal/realtimequality/`
- `testdata/replay/`
- `internal/prediction/`
- `internal/state/`
- `internal/feed/tripupdates/`
- `cmd/agency-config/operations.go` only for safe summary display
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-29.md`

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Keep `/v1/events` documented and operated as an authenticated admin/debug review path, not as a public or consumer-facing feed.
- Do not add consumer submission APIs, automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, universally production ready, production multi-tenant hosted, or production-grade ETA proven.

## Exact Next-Step Recommendation

Start Phase 29 with richer replay fixtures for after-midnight service, frequency-window trips, block continuity, long layovers, sparse telemetry, noisy GPS, stale/ambiguous real-world patterns, cancellation/alert linkage, and manual overrides over time. Add metrics only where denominators and evidence are explicit, and keep unknown/withheld outcomes visible.
