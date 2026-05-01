Read and obey AGENTS.md first.

Then read:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/handoffs/track-b-roadmap.md
4. docs/track-b-productization-roadmap.md
5. docs/phase-22-release-distribution-hardening.md
6. docs/release-process.md
7. CONTRIBUTING.md
8. SECURITY.md
9. docs/evidence/redaction-policy.md
10. README.md
11. docs/dependencies.md
12. docs/decisions.md

You are starting Phase 22 — Release And Distribution Hardening only.

Goal:
Make Open Transit RT easier to version, release, install, and upgrade without changing backend product behavior.

Scope:
- changelog
- release checklist
- versioning/tagging guidance
- upgrade/migration notes
- Docker image/release artifact plan
- install verification commands
- rollback/restore notes for releases
- release note template
- CI/release workflow documentation if appropriate

Do not:
- change backend API behavior
- change database schema unless explicitly needed and justified
- change public feed URLs
- change consumer statuses
- claim compliance or consumer acceptance
- add hosted SaaS claims
- add external integrations

Acceptance criteria:
- maintainers can cut a tagged release consistently
- users can tell what version they are running
- release notes include migrations, operations changes, security notes, evidence/claim changes, known limitations, and checks
- install/upgrade/rollback docs exist
- no unsupported claims are introduced

Run before editing:

```bash
make validate
make test
git diff --check
```

Run after editing:

```bash
make validate
make test
git diff --check
```

If release docs touch deployment assumptions, also run:

```bash
docker compose -f deploy/docker-compose.yml config
make smoke
```

Add:
- docs/handoffs/phase-22.md

Update:
- docs/current-status.md
- docs/handoffs/latest.md
- docs/phase-22-release-distribution-hardening.md
- docs/release-process.md
- README.md only if a short release/install link is needed

The Phase 22 handoff must include:
1. what was implemented
2. what was intentionally deferred
3. release/distribution docs added
4. versioning/tagging guidance
5. install/upgrade/rollback guidance
6. commands run
7. blocked commands
8. known remaining release/distribution gaps
9. exact recommendation for Phase 23
