# Phase 41 Kickoff — Operator Smoke And Support Bundle

**Status:** proposed kickoff, not implementation handoff  
**Do not mark Phase 41 complete from this file alone.**

## Goal

Create a repeatable operator smoke workflow and redaction-safe support bundle
for local/reference Open Transit RT deployments.

## Read First

1. `AGENTS.md`
2. `docs/current-status.md`
3. `docs/handoffs/latest.md`
4. `docs/tutorials/self-hosted-operator-trial.md`
5. `docs/integration-adapter-kit.md`
6. `docs/evidence/redaction-policy.md`
7. `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
8. `docs/phase-41-operator-smoke-support-bundle.md`

## Required Outputs

- `scripts/operator-smoke.sh`
- `scripts/support-bundle.sh`
- `make operator-smoke`
- `make support-bundle`
- `docs/tutorials/operator-smoke-and-support-bundle.md`
- completed `docs/handoffs/phase-41.md`
- navigation/status updates

## Claim Boundaries

Phase 41 creates no external evidence, no final-root evidence, no consumer
submission evidence, no agency adoption evidence, no vendor evidence, and no
compliance evidence.

Consumer statuses must remain unchanged.

## Checks

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Run additional relevant checks if the implementation touches Docker, smoke,
validators, deployment, or admin routes.
