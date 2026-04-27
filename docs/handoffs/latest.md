# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 20 — Consumer Submission Execution And CAL-ITP Readiness Program is complete for the approved docs/evidence packet-preparation scope.

Phases 0 through 20 are closed for their documented scopes. Do not reopen earlier phases unless a blocking truthfulness, safety, or submission-readiness issue directly requires it.

## Phase Status

- Phase 20 added complete prepared packet drafts for Google Maps, Apple Maps, Transit App, Bing Maps, Moovit, Mobility Database, and transit.land.
- The packet index and completeness table live at `docs/evidence/consumer-submissions/packets/README.md`.
- The machine-readable snapshot lives at `docs/evidence/consumer-submissions/status.json`.
- The human-readable tracker and `status.json` agree for target name, status, packet path, prepared timestamp, and evidence reference values.
- All seven targets are `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.
- Official submission methods/contact paths are marked `not verified`; no submission path was guessed.
- `docs/california-readiness-summary.md` separates code-complete capability, deployment-proven OCI pilot evidence, prepared packet evidence, submitted evidence, under-review evidence, accepted evidence, and missing evidence.
- `docs/marketplace-vendor-gap-review.md` records the remaining vendor-like adoption gaps.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- No external portal was contacted, no submission was automated, no backend API behavior was added, and no unsupported acceptance, compliance, consumer-ingestion, marketplace-equivalence, hosted SaaS, agency-endorsement, or production-grade ETA claim was added.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-20.md`
4. `docs/evidence/consumer-submissions/README.md`
5. `docs/evidence/consumer-submissions/status.json`
6. `docs/evidence/consumer-submissions/packets/README.md`
7. `docs/california-readiness-summary.md`
8. `docs/marketplace-vendor-gap-review.md`
9. `docs/compliance-evidence-checklist.md`
10. `docs/prompts/calitp-truthfulness.md`
11. `docs/evidence/redaction-policy.md`
12. `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
13. `docs/handoffs/phase-19.md`
14. `testdata/replay/README.md`
15. `README.md`
16. `SECURITY.md`
17. `docs/dependencies.md`
18. `docs/decisions.md`

## Current Objective

The next planned phase is Phase 21 — Community, Governance, And Multi-Agency Scale, if the roadmap still applies.

If real consumer or aggregator artifacts arrive before Phase 21 work, update only the named target record, `docs/evidence/consumer-submissions/status.json`, and the human-readable tracker to the evidence-backed status. Do not infer submission, review, acceptance, rejection, blocker status, ingestion, or compliance from validator success, public fetch proof, or prepared packets.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

If evidence/readiness docs change materially, also run:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

If local app/demo docs or scripts are touched, also run:

```bash
make demo-agency-flow
make agency-app-up
make agency-app-down
docker compose -f deploy/docker-compose.yml --profile app config
```

## Current Evidence And Security Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` remains the current hosted/operator evidence packet.
- Phase 20 prepared packets are operator review artifacts only; they are not submissions.
- Replay fixtures measure current realtime behavior only; they are not consumer acceptance, production-grade ETA proof, or CAL-ITP/Caltrans compliance.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, or raw private operator artifacts.
- The Operations Console intentionally omits token hashes, raw device tokens except immediate one-time rebind responses, raw telemetry payload JSON, full assignment score details, DB URLs, private URLs, and private debug blobs.

## First Files Likely To Edit

- `docs/evidence/consumer-submissions/current/<target>.md` only after real target-originated evidence exists
- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-21.md` when Phase 21 starts

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs, automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Keep `status.json` and the human-readable tracker aligned for target name, status, packet path, prepared timestamp, and evidence references.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/roadmap-post-phase-14.md` as the roadmap source of truth.

The next planned phase is:

- Phase 21 — Community, Governance, And Multi-Agency Scale

Future roadmap docs:

- `docs/phase-21-community-governance-multi-agency.md`
