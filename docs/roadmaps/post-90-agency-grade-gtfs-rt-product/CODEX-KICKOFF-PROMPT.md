# CODEX KICKOFF PROMPT — Post-90 Agency-Grade GTFS-RT Roadmap

You are the Master Agent for Open Transit RT.

Add this roadmap pack under:

```text
docs/roadmaps/post-90-agency-grade-gtfs-rt-product/
```

Then start **Phase 91 only**.

## Read first

```text
AGENTS.md
README.md
docs/current-status.md
docs/handoffs/latest.md
docs/handoffs/phase-90.md
docs/phase-90-control-plane-final-status.md
docs/phase-89-rc1-gate-results.md
docs/roadmap-status.md
docs/evidence/redaction-policy.md
docs/evidence/consumer-submissions/status.json
```

## Start

```text
Phase 91 -- Maintainer Route/Product Audit And Stabilization
```

## Phase 91 checkpoints

```text
Phase 91 -- Checkpoint 000001: add route product audit plan
Phase 91 -- Checkpoint 000002: audit private routes user tasks and docs drift
Phase 91 -- Checkpoint 000003: add route inventory audit helper
Phase 91 -- Checkpoint 000004: patch highest priority IA copy and route gaps
Phase 91 -- Checkpoint 000005: close route product audit
```

## Hard boundaries

Do not edit protected evidence paths, move consumer statuses, collect evidence,
contact external parties, tag/package/publish, or make stronger claims.

## Validation

```bash
git status --short
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured db/migrations go.mod go.sum
```

If code changes:

```bash
make validate
make test
docker compose -f deploy/docker-compose.yml config
```
