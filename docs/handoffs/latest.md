# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 19 — Realtime Quality And ETA Improvement is complete for the approved measurement-first scope.

Phases 0 through 19 are closed for their documented scopes. Do not reopen earlier phases unless a blocking truthfulness, safety, or realtime-quality issue directly requires it.

## Phase Status

- Phase 19 added deterministic replay evaluation under `internal/realtimequality`.
- Replay fixtures live under `testdata/replay/` and document current behavior for matched, stale, ambiguous, low-confidence, manual override, canceled-trip, added-trip, short-turn, and detour cases.
- `prediction.Metrics` now records explicit quality counts and rate objects with numerator, denominator, denominator definition, status, and `not_applicable` zero-denominator handling.
- Unknown, ambiguous, stale, withheld, and degraded cases remain visible in metrics and diagnostics.
- Operations Console feed/dashboard views show safe Trip Updates quality summaries only when `feed_health_snapshot` diagnostics exist; otherwise they say `no Trip Updates diagnostics recorded yet`.
- No external predictor was integrated.
- No production-grade ETA quality, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, agency endorsement, or marketplace equivalence claim was added.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-19.md`
4. `testdata/replay/README.md`
5. `docs/phase-19-realtime-quality-eta-improvement.md`
6. `README.md`
7. `wiki/README.md`
8. `docs/README.md`
9. `docs/runbooks/small-agency-pilot-operations.md`
10. `SECURITY.md`
11. `docs/evidence/redaction-policy.md`
12. `docs/compliance-evidence-checklist.md`
13. `docs/prompts/calitp-truthfulness.md`
14. `docs/dependencies.md`
15. `docs/decisions.md`

## Current Objective

The next planned phase is Phase 20 — Consumer Submission / Cal-ITP Readiness, if the roadmap still applies.

Use Phase 19 replay metrics as measurement evidence for current realtime behavior, not as proof of production-grade ETA quality. If more realtime work is selected before Phase 20, add replay fixtures first and improve only cases justified by those fixtures or tests.

Do not change public feed URLs, GTFS-RT protobuf contracts, Trip Updates adapter boundaries, consumer-submission statuses, unauthenticated surfaces, or evidence claims unless the active phase explicitly requires it.

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

If touching local app/demo docs or scripts, also run:

```bash
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

## Current Evidence And Security Boundary

- Phase 19 did not collect new hosted evidence.
- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Replay fixtures measure current behavior only; they are not consumer acceptance, production-grade ETA proof, or CAL-ITP/Caltrans compliance.
- Consumer-ingestion workflow records and Phase 13 docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, or raw private operator artifacts.
- The Operations Console intentionally omits token hashes, raw device tokens except immediate one-time rebind responses, raw telemetry payload JSON, full assignment score details, DB URLs, private URLs, and private debug blobs.

## First Files Likely To Edit

- `internal/realtimequality/`
- `testdata/replay/`
- `internal/prediction/`
- `internal/state/`
- `internal/feed/tripupdates/`
- `cmd/agency-config/operations.go` if adding only safe authenticated quality summaries
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-20.md` when Phase 20 starts

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Keep replay deterministic: fixed timestamps, stable fixture data, stable report ordering, and no wall-clock dependency except injected clocks.
- Keep zero-denominator metrics honest with `not_applicable` or an explicit omission reason.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs, automate fake submissions, or invent acceptance/rejection/compliance evidence.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/roadmap-post-phase-14.md` as the roadmap source of truth.

The next planned phase is:

- Phase 20 — Consumer Submission / Cal-ITP Readiness

Future roadmap docs:

- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`
