# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 22 — Release And Distribution Hardening is complete.

Phases 0 through 22 are closed for their documented scopes. Track A is also closed for its docs-only external-proof workflow scope. Do not reopen earlier phases unless a blocking truthfulness, safety, governance, multi-agency, or submission-readiness issue directly requires it.

## Phase Status

- Phase 21 added `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, GitHub issue templates, and a PR template.
- Phase 21 added `docs/governance.md`, `docs/release-process.md`, `docs/support-boundaries.md`, `docs/multi-agency-strategy.md`, and `docs/roadmap-status.md`.
- Governance docs explicitly state who can merge PRs, cut releases, approve docs/evidence wording, and resolve competing design decisions.
- Issue templates warn users not to paste tokens, DB URLs, private keys, admin URLs with secrets, private portal screenshots, private ticket links, raw logs with credentials, or unredacted operator artifacts.
- Teaching visuals were generated under `docs/assets/` and documented in `docs/assets/README.md`.
- Phase 21 did not change backend behavior, API contracts, database schema, public feed URLs, consumer-submission statuses, external integrations, or evidence claims.
- Track A added `docs/evidence/consumer-submissions/submission-workflow.md`.
- Track A added README-only target artifact directories under `docs/evidence/consumer-submissions/artifacts/`.
- Track A added `docs/agency-owned-domain-readiness.md`.
- Track A added no helper scripts, no backend behavior, no portal automation, no public feed URL changes, and no consumer status changes.
- Track B roadmap docs have been added for Phase 22 through Phase 32.
- `docs/track-b-productization-roadmap.md` is the Track B roadmap source.
- `docs/handoffs/track-b-roadmap.md` records the Track B roadmap handoff.
- Phase 22 added `CHANGELOG.md`, `docs/release-checklist.md`, `docs/upgrade-and-rollback.md`, and `docs/release-notes-template.md`.
- Phase 22 expanded `docs/release-process.md` with release-from-main, tag, version verification, artifact, install, upgrade, rollback, release note, and evidence packet version-linkage guidance.
- Phase 22 explicitly documents that current distribution guidance supports source tags and local Docker builds only; published/versioned production Docker images are deferred.
- Phase 23 — Agency-Owned Deployment Proof is the recommended next implementation phase.
- Track B must preserve truthfulness, redaction, and security boundaries.
- All seven consumer and aggregator targets remain `prepared` only. No target has submitted, under-review, accepted, rejected, or blocked evidence.
- The OCI pilot DuckDNS hostname remains pilot evidence, not agency-owned stable URL/domain proof.
- No unsupported compliance, consumer-ingestion, marketplace-equivalence, hosted SaaS, paid support/SLA, agency-endorsement, or production-grade ETA claim was added.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/track-a-external-proof.md`
4. `docs/handoffs/track-b-roadmap.md`
5. `docs/track-b-productization-roadmap.md`
6. `docs/phase-22-release-distribution-hardening.md`
7. `docs/handoffs/phase-22.md`
8. `docs/phase-23-agency-owned-deployment-proof.md`
9. `docs/release-process.md`
10. `docs/release-checklist.md`
11. `docs/upgrade-and-rollback.md`
12. `docs/release-notes-template.md`
13. `CHANGELOG.md`
14. `docs/evidence/consumer-submissions/README.md`
15. `docs/evidence/consumer-submissions/submission-workflow.md`
16. `docs/evidence/consumer-submissions/status.json`
17. `docs/evidence/consumer-submissions/artifacts/README.md`
18. `docs/agency-owned-domain-readiness.md`
19. `docs/california-readiness-summary.md`
20. `docs/marketplace-vendor-gap-review.md`
21. `docs/compliance-evidence-checklist.md`
22. `docs/prompts/calitp-truthfulness.md`
23. `docs/evidence/redaction-policy.md`
24. `SECURITY.md`
25. `docs/roadmap-status.md`
26. `README.md`
27. `docs/dependencies.md`
28. `docs/decisions.md`

## Current Objective

The recommended next implementation phase is Phase 23 — Agency-Owned Deployment Proof. Do not implement the next Track B phase until maintainers explicitly start it.

For external proof work, use the Track A workflow before verifying official paths, submitting packets, or recording target-originated evidence. Track B productization work must not advance consumer statuses unless retained, redacted, target-originated evidence exists for the named target.

If real consumer or aggregator artifacts arrive, update only the named target record, `docs/evidence/consumer-submissions/status.json`, and the human-readable tracker to the evidence-backed status. Do not infer submission, review, acceptance, rejection, blocker status, ingestion, or compliance from validator success, public fetch proof, OCI pilot evidence, replay metrics, or prepared packets.

## Exact First Commands

```bash
make validate
make test
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
git diff --check
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
- Track A artifact directories are README-only until real redacted target-originated evidence exists.
- Track B roadmap docs and Phase 22 release docs do not prove compliance, acceptance, hosted service, paid support, SLA coverage, agency endorsement, vendor equivalence, or production-grade ETA quality.
- Replay fixtures measure current realtime behavior only; they are not consumer acceptance, production-grade ETA proof, or CAL-ITP/Caltrans compliance.
- Consumer-ingestion workflow records and docs tracker records are not third-party acceptance unless retained evidence from the named target exists.
- Do not rely on old local `.cache` credentials.
- Do not commit secrets, generated tokens, private keys, ACME material, admin tokens, device tokens, JWT secrets, CSRF secrets, DB passwords, webhook URLs, notification credentials, raw telemetry payloads, unredacted correspondence, private portal credentials, private ticket links, or raw private operator artifacts.
- Issue templates and support docs require redacted logs and public-safe evidence.

## First Files Likely To Edit

- `docs/phase-23-agency-owned-deployment-proof.md` when maintainers begin Phase 23 implementation.
- `docs/handoffs/latest.md` and `docs/current-status.md` when the next phase starts or closes.
- `docs/roadmap-post-phase-14.md` only if maintainers revise the roadmap structure.
- `docs/evidence/consumer-submissions/current/<target>.md` only after real target-originated evidence exists.
- `docs/evidence/consumer-submissions/artifacts/<target>/` only after real redacted target-originated evidence exists.
- `.github/ISSUE_TEMPLATE/*.yml` or `.github/pull_request_template.md` only if maintainers adjust triage process.

## Constraints To Preserve

- Keep Trip Updates pluggable and Vehicle Positions first.
- Preserve conservative matching: unknown is better than false certainty.
- Preserve admin auth, role checks, CSRF behavior, and token/secret handling.
- Do not expose admin/debug/JSON surfaces on the production public edge.
- Do not add consumer submission APIs, automate submissions, contact external portals, guess submission paths, or invent acceptance/rejection/compliance evidence.
- Keep `prepared` conditional on packet completeness.
- Keep artifact directories empty except README files unless real redacted target-originated evidence is provided.
- Keep `status.json` and the human-readable tracker aligned for target name, status, packet path, prepared timestamp, and evidence references.
- Keep local `http://localhost:8080` wording scoped to local-demo packaging only.
- Do not describe Open Transit RT as hosted SaaS, paid support, SLA-backed, agency-endorsed, marketplace/vendor equivalent, or universally production ready.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.

## Future Roadmap

Use `docs/track-b-productization-roadmap.md` as the forward roadmap source of truth for Track B. Use `docs/roadmap-post-phase-14.md` for historical post-Phase-14 context.

Phase 23 — Agency-Owned Deployment Proof is the recommended next implementation phase.

## Checks Run For Phase 22

- pre-edit `make validate` — passed.
- pre-edit `make test` — passed.
- pre-edit `git diff --check` — passed.
- `make validate` — passed.
- `make test` — passed.
- `make realtime-quality` — passed.
- `make smoke` — passed.
- `docker compose -f deploy/docker-compose.yml config` — passed.
- `git diff --check` — passed.
