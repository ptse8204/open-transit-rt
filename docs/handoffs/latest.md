# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 21 — Community, Governance, And Multi-Agency Scale is complete for the approved docs/process/governance/teaching-visual scope.

Phases 0 through 21 are closed for their documented scopes. Do not reopen earlier phases unless a blocking truthfulness, safety, governance, multi-agency, or submission-readiness issue directly requires it.

## Phase Status

- Phase 21 added `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, GitHub issue templates, and a PR template.
- Phase 21 added `docs/governance.md`, `docs/release-process.md`, `docs/support-boundaries.md`, `docs/multi-agency-strategy.md`, and `docs/roadmap-status.md`.
- Governance docs explicitly state who can merge PRs, cut releases, approve docs/evidence wording, and resolve competing design decisions.
- Issue templates warn users not to paste tokens, DB URLs, private keys, admin URLs with secrets, private portal screenshots, private ticket links, raw logs with credentials, or unredacted operator artifacts.
- Teaching visuals were generated under `docs/assets/` and documented in `docs/assets/README.md`.
- Phase 21 did not change backend behavior, API contracts, database schema, public feed URLs, consumer-submission statuses, external integrations, or evidence claims.
- All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- No unsupported compliance, consumer-ingestion, marketplace-equivalence, hosted SaaS, paid support/SLA, agency-endorsement, or production-grade ETA claim was added.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/phase-21.md`
4. `CONTRIBUTING.md`
5. `docs/governance.md`
6. `docs/support-boundaries.md`
7. `docs/multi-agency-strategy.md`
8. `docs/roadmap-status.md`
9. `docs/evidence/consumer-submissions/README.md`
10. `docs/evidence/consumer-submissions/status.json`
11. `docs/california-readiness-summary.md`
12. `docs/marketplace-vendor-gap-review.md`
13. `docs/compliance-evidence-checklist.md`
14. `docs/prompts/calitp-truthfulness.md`
15. `docs/evidence/redaction-policy.md`
16. `README.md`
17. `SECURITY.md`
18. `docs/dependencies.md`
19. `docs/decisions.md`

## Current Objective

The next planned work should be selected explicitly by maintainers. A practical next step is a post-Phase-21 roadmap update that decides whether to prioritize multi-agency isolation tests, production operations hardening, or real consumer-submission evidence intake after target-originated artifacts exist.

If real consumer or aggregator artifacts arrive, update only the named target record, `docs/evidence/consumer-submissions/status.json`, and the human-readable tracker to the evidence-backed status. Do not infer submission, review, acceptance, rejection, blocker status, ingestion, or compliance from validator success, public fetch proof, OCI pilot evidence, replay metrics, or prepared packets.

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
- Issue templates and support docs require redacted logs and public-safe evidence.

## First Files Likely To Edit

- `docs/roadmap-post-phase-14.md` if maintainers define the post-Phase-21 roadmap.
- `docs/handoffs/latest.md` and `docs/current-status.md` when the next phase starts or closes.
- `docs/evidence/consumer-submissions/current/<target>.md` only after real target-originated evidence exists.
- `.github/ISSUE_TEMPLATE/*.yml` or `.github/pull_request_template.md` only if maintainers adjust triage process.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs, automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Keep `status.json` and the human-readable tracker aligned for target name, status, packet path, prepared timestamp, and evidence references.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, or universally production ready.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/roadmap-post-phase-14.md` as the roadmap source of truth until maintainers update it.

No post-Phase-21 phase is selected in this handoff.
