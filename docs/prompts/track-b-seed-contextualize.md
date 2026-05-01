Read and obey AGENTS.md first.

Then read:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/handoffs/track-a-external-proof.md
4. docs/roadmap-post-phase-14.md
5. docs/roadmap-status.md
6. docs/california-readiness-summary.md
7. docs/marketplace-vendor-gap-review.md
8. docs/evidence/consumer-submissions/submission-workflow.md
9. docs/prompts/calitp-truthfulness.md
10. SECURITY.md
11. docs/evidence/redaction-policy.md

This is a docs-only Track B roadmap contextualization task.

The user has provided the Track B roadmap Markdown files. Copy them into the repository under the exact paths listed below, then update the repo’s status and roadmap pointers so future Codex instances can continue from repo-native context.

Do not invent new phase content unless needed to align file paths or repo conventions.
Do not implement Track B yet.
Do not change backend behavior.
Do not change API contracts.
Do not change database schema.
Do not change public feed URLs.
Do not change consumer statuses.
Do not add external integrations.
Do not claim compliance, consumer acceptance, vendor equivalence, hosted SaaS, agency endorsement, paid support, SLA coverage, or production-grade ETA quality.

Files to add:

- docs/track-b-productization-roadmap.md
- docs/phase-22-release-distribution-hardening.md
- docs/phase-23-agency-owned-deployment-proof.md
- docs/phase-24-real-agency-data-onboarding.md
- docs/phase-25-device-avl-integration-kit.md
- docs/phase-26-admin-ux-setup-wizard.md
- docs/phase-27-multi-agency-isolation-prototype.md
- docs/phase-28-production-operations-hardening.md
- docs/phase-29-realtime-quality-expansion.md
- docs/phase-30-consumer-submission-execution.md
- docs/phase-31-agency-pilot-program-package.md
- docs/phase-32-public-launch-ecosystem-outreach.md
- docs/handoffs/track-b-roadmap.md

Update these files to contextualize the new roadmap:

- docs/current-status.md
- docs/handoffs/latest.md
- docs/roadmap-post-phase-14.md
- docs/roadmap-status.md
- docs/README.md if useful for navigation

Required latest-handoff wording:

- Track A is closed for the docs-only external-proof workflow.
- Track B roadmap docs have been added.
- Phase 22 — Release And Distribution Hardening is the recommended next implementation phase.
- Track B must not advance consumer statuses unless target-originated evidence exists.
- Track B must preserve truthfulness, redaction, and security boundaries.

Run and record:

```bash
make validate
make test
git diff --check
```

If readiness/evidence docs are materially updated, also run:

```bash
make realtime-quality
make smoke
docker compose -f deploy/docker-compose.yml config
```

Acceptance criteria:

- All Track B roadmap files exist under docs/.
- Status/latest/roadmap docs point to Track B.
- Phase 22 is clearly named as the first recommended implementation phase.
- No consumer status changes are made.
- No unsupported claims are introduced.
- No backend/runtime behavior changes are made.
- A handoff for Track B roadmap seeding exists at docs/handoffs/track-b-roadmap.md.
