# Prompt — Start Phase 29A External Predictor Adapter Evaluation

```text
Read and obey AGENTS.md first.

Then read:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/handoffs/phase-29.md
4. docs/phase-29a-external-predictor-adapter-evaluation.md
5. docs/phase-29-realtime-quality-expansion.md
6. testdata/replay/README.md
7. docs/tutorials/device-avl-integration.md
8. docs/evidence/redaction-policy.md
9. SECURITY.md
10. docs/dependencies.md
11. docs/decisions.md

You are starting Phase 29A — External Predictor Adapter Evaluation only.

Goal:
Evaluate external prediction adapter feasibility, including TheTransitClock as a candidate, while keeping deterministic prediction as default and avoiding runtime external dependencies unless explicitly approved.

Do not implement full runtime integration unless the phase plan is updated and approved.

Focus on:
- internal/prediction.Adapter contract review
- external predictor candidate feasibility
- mock/test adapter contract if useful
- replay comparison against Phase 29 fixtures
- failure/fallback behavior
- licensing/dependency review
- docs/handoffs/phase-29a.md

Run:
- make validate
- make realtime-quality
- make test
- make smoke
- make test-integration
- docker compose -f deploy/docker-compose.yml config
- git diff --check
```
