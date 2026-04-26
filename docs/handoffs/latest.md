# Latest Handoff

This file is the source of truth for the next Codex instance.

## Active Phase

Phase 14 — Public Launch Polish and Repo Simplification is complete for the docs/presentation/navigation scope.

Phase 12 remains closed for the OCI pilot hosted/operator evidence scope. Phase 13 remains closed for the initial consumer-submission evidence tracker structure. Do not reopen either phase unless a blocking documentation truthfulness issue directly affects the next task.

## Phase Status

- Phases 0 through 11 are closed for their documented scope.
- Phase 12 is closed for the OCI pilot hosted/operator evidence scope.
- Phase 13 is closed for the consumer-submission evidence tracker structure.
- Phase 14 simplified README/docs navigation, split public docs into `wiki/`, marked `docs/` as internal reference material, and added reviewed teaching visuals.
- All seven current consumer/aggregator records are still `not_started`.
- No current repo evidence supports submitted, under-review, accepted, rejected, or blocked claims for any consumer target.

## Read These Files First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `README.md`
4. `wiki/README.md`
5. `docs/README.md`
6. `docs/handoffs/phase-14.md`
7. `docs/phase-14-public-launch-polish.md`
8. `docs/compliance-evidence-checklist.md`
9. `docs/consumer-submission-evidence.md`
10. `docs/evidence/consumer-submissions/README.md`
11. `docs/evidence/captured/oci-pilot/2026-04-24/README.md`
12. `docs/prompts/calitp-truthfulness.md`
13. `docs/tutorials/README.md`
14. `docs/tutorials/local-quickstart.md`
15. `docs/tutorials/agency-demo-flow.md`
16. `docs/tutorials/deploy-with-docker-compose.md`
17. `docs/tutorials/production-checklist.md`
18. `docs/tutorials/calitp-readiness-checklist.md`
19. `docs/assets/README.md`
20. `docs/dependencies.md`
21. `docs/decisions.md`

## Current Objective

Use the simplified README and `wiki/` as the public front door. Future docs work should preserve the concise README, keep public reader docs in `wiki/`, keep detailed evidence/history in `docs/`, and maintain the visual review rule for generated or generated-assisted assets.

Do not claim CAL-ITP/Caltrans compliance, production readiness, marketplace/vendor equivalence, agency endorsement, or consumer acceptance from repo evidence, validator success, public fetch proof, workflow records, stars, or the Phase 12 hosted packet alone.

## Exact First Commands

```bash
make validate
make test
git diff --check
```

If Docker is available, also run:

```bash
make smoke
make demo-agency-flow
```

## Current Evidence Boundary

- The OCI pilot packet at `docs/evidence/captured/oci-pilot/2026-04-24/` includes hosted/operator proof for public HTTPS feed fetches, TLS/redirect behavior, auth boundaries, clean hosted validation, monitoring/alert lifecycle, backup/restore, deployment rollback, and scorecard export job history.
- The Phase 13 tracker at `docs/evidence/consumer-submissions/README.md` links to the Phase 12 packet as supporting evidence only.
- Validator success and public fetch proof are not consumer acceptance.
- Consumer-ingestion workflow records are not third-party acceptance.
- Acceptance may be claimed only for the named consumer, feed scope, URL root, and evidence date shown in a retained evidence artifact.

## First Files Likely To Edit

- `README.md`
- `wiki/README.md`
- `docs/README.md`
- `docs/tutorials/README.md`
- `docs/assets/README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`

## Constraints To Preserve

- Keep README understandable by a non-technical agency reader in under 3 minutes.
- Keep README concise and easy to scan, ideally under 150 to 200 lines unless examples genuinely require more.
- Keep claims evidence-bounded and truthful.
- Keep captions explicit about illustrative versus exact-behavior visuals.
- Keep alt text descriptive and useful.
- Do not change backend runtime behavior, API contracts, database schema, public feed URLs, external integrations, or consumer-submission status.
- Do not add consumer submission APIs unless explicitly required and supported by a public documented API.
- Do not automate fake submissions.
- Do not invent acceptance, rejection, receipt, or blocker evidence.

## Handoff Template Requirement

All future phase handoff files must use `docs/handoffs/template.md` unless the phase explicitly documents a reason to diverge.
