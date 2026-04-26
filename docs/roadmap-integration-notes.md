# Roadmap Integration Notes

Copy the roadmap files from this bundle into the repository under `docs/`.

## Files To Add

- `docs/roadmap-post-phase-14.md`
- `docs/phase-15-public-repo-security-hygiene.md`
- `docs/phase-16-agency-onboarding-product-packaging.md`
- `docs/phase-17-deployment-automation-pilot-operations.md`
- `docs/phase-18-admin-ux-agency-operations-console.md`
- `docs/phase-19-realtime-quality-eta-improvement.md`
- `docs/phase-20-consumer-submission-calitp-readiness.md`
- `docs/phase-21-community-governance-multi-agency.md`

## Files To Update

Update `docs/handoffs/latest.md` after Phase 14 closes so future Codex instances know which phase is active.

Recommended addition to `docs/handoffs/latest.md` after Phase 14:

```markdown
## Future Roadmap

After Phase 14, use `docs/roadmap-post-phase-14.md` as the roadmap source of truth.
The next planned phase is Phase 15 — Public Repo Security Hygiene And Artifact Redaction.
```

Update `docs/current-status.md` similarly after Phase 14 closes.

## Codex Usage Pattern

Use one fresh Codex task per phase. Do not ask one instance to implement the whole roadmap.

Each phase should:

1. read `AGENTS.md`;
2. read `docs/current-status.md`;
3. read `docs/handoffs/latest.md`;
4. read `docs/roadmap-post-phase-14.md`;
5. read the phase-specific plan;
6. execute only that phase;
7. update status and handoff docs;
8. run required checks.
