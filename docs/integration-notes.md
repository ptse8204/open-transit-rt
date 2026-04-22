# Integration Notes For Repo Context Loading

These files are intended to be added to the repo using the same documentation-driven context pattern already used by the codebase.

## Recommended Paths

Add these files directly into the repo at:

- `docs/phase-plan-production-closure.md`
- `docs/prompts/codex-production-closure.md`
- `docs/prompts/calitp-truthfulness.md`
- `docs/prompts/docs-assets-image-generation.md`
- `docs/tutorials/README.md`

## Recommended Existing Docs To Update

After adding the new files, update these repo docs so future Codex runs can discover them naturally:

- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/phase-plan.md`
- `docs/dependencies.md`
- `README.md`

## Recommended Read Order For New Codex Runs

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/latest.md`
4. `docs/phase-plan.md`
5. `docs/phase-plan-production-closure.md`
6. `docs/prompts/codex-production-closure.md`
7. `docs/prompts/calitp-truthfulness.md`
8. `docs/prompts/docs-assets-image-generation.md`
9. `docs/dependencies.md`
10. `docs/decisions.md`

## Suggested Minimal Pointer Text

### For `docs/current-status.md`
Add a short note that post-Phase-8 work is driven by `docs/phase-plan-production-closure.md`.

### For `docs/handoffs/latest.md`
Point the next Codex instance to:
- the production closure plan
- the codex production closure prompt
- the CAL-ITP truthfulness guardrail
- the docs asset generation guidance

### For `README.md`
Once Phase 10 starts, link to the tutorial docs and docs assets.
