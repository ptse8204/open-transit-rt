# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 18 — Admin UX And Agency Operations Console is complete for the approved minimal operator-console scope.

Phases 0 through 18 are closed for their documented scopes. Do not reopen earlier phases unless a blocking truthfulness, safety, or admin-UX issue directly requires it.

## Phase Status

- Phase 18 added authenticated server-rendered operations pages under `/admin/operations`.
- Alerts now have a minimal authenticated browser console under `/admin/alerts/console`.
- The local app package prints `/admin/operations`, and local proxy routing includes `/admin/operations*`.
- Consumer/aggregator records remain evidence records only; no current repo evidence supports submitted, under-review, accepted, rejected, or blocked claims for any consumer target unless a deployment updates those records with retained evidence.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `README.md`
4. `wiki/README.md`
5. `docs/README.md`
6. `docs/handoffs/phase-18.md`
7. `docs/phase-19-realtime-quality-eta-improvement.md`
8. `docs/runbooks/small-agency-pilot-operations.md`
9. `SECURITY.md`
10. `docs/evidence/redaction-policy.md`
11. `docs/compliance-evidence-checklist.md`
12. `docs/prompts/calitp-truthfulness.md`
13. `docs/dependencies.md`
14. `docs/decisions.md`

## Current Objective

The next planned phase is Phase 19 — Realtime Quality And ETA Improvement.

Use the Phase 18 Operations Console as an operator visibility surface for existing state, but do not treat it as proof of consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, agency endorsement, or universal production readiness.

Do not change public feed URLs, GTFS-RT protobuf contracts, Trip Updates adapter boundaries, consumer-submission statuses, or public-edge production proxy routes unless the active phase explicitly requires it.

## Exact First Commands

```bash
make validate
make test
make smoke
git diff --check
docker compose -f deploy/docker-compose.yml config
```

If touching local app/demo docs or scripts, also run:

```bash
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

## Current Evidence And Security Boundary

- Phase 18 did not collect new hosted evidence.
- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Validator success and public fetch proof are supporting evidence only, not consumer acceptance.
- Consumer-ingestion workflow records and Phase 13 docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, or raw private operator artifacts.
- The Operations Console intentionally omits token hashes, raw device tokens except immediate one-time rebind responses, raw telemetry payload JSON, full assignment score details, DB URLs, private URLs, and private debug blobs.

## First Files Likely To Edit

- `internal/prediction/`
- `internal/state/`
- `internal/feed/tripupdates/`
- `cmd/feed-trip-updates/`
- `cmd/agency-config/` if Phase 19 adds quality/coverage visibility to the console
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-19.md`

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs, automate fake submissions, or invent acceptance/rejection/compliance evidence.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/roadmap-post-phase-14.md` as the roadmap source of truth.

The next planned phase is:

- Phase 19 — Realtime Quality And ETA Improvement

Future roadmap docs:

- `docs/phase-19-realtime-quality-eta-improvement.md`
- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`
