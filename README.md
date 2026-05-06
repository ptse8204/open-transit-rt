# Open Transit RT — Post-Outcome-C Roadmap Export

Generated: 2026-05-06

This export updates the prior roadmap package to match the current repo state where **Phase 33 — Public GTFS Local/Pilot Evidence is complete as Outcome C**.

Use these files as repo-ready drafts for Codex. They are written to be copied into the project at the paths shown below, but they should still be reviewed before committing.

## Current baseline assumed by this export

- Phases 0 through 33 are closed for their documented scopes.
- Phase 33 is complete as Outcome C for **local/pilot public static GTFS dataset handling**.
- The LA Metro Bus public GTFS run completed locally with retained public-safe summaries.
- The large-import blocker exposed during the first LA Metro attempt was fixed inside the Phase 33 work: configurable import timeout, `pgx.CopyFrom` for large stop-time and shape-point publish inserts, and safer failure-report persistence.
- Phase 33 does **not** prove agency adoption, agency approval, official feed status, agency-owned final-root proof, consumer submission/review/acceptance, Caltrans/CAL-ITP compliance, hosted SaaS, production readiness, real vendor AVL compatibility, real LA Metro realtime data, real-world ETA accuracy, or production-grade ETA quality.
- Consumer targets remain `prepared` only.
- The post-Phase-32 final-root blocker remains unresolved.

## Files in this export

Copy or review these files in the corresponding repo locations:

```text
docs/handoffs/repo-evaluation-forward-roadmap-2026-05-06-post-outcome-c.md
docs/phase-34-post-outcome-c-status-consistency-and-evidence-readiness.md
docs/handoffs/phase-34-kickoff.md
docs/future-roadmap-post-outcome-c.md
docs/status-maintenance-patch-post-outcome-c.md
docs/final-root-operator-request.md
docs/tutorials/public-gtfs-local-pilot.md
```

## Recommended use

1. Start Codex with `docs/handoffs/phase-34-kickoff.md`.
2. Treat `docs/handoffs/latest.md` as the source of truth before editing.
3. Preserve all claim boundaries.
4. Do not update consumer statuses, final-root evidence, or target-specific artifacts unless retained external evidence exists.
5. After Codex applies the patch, run:

```bash
make validate
make test
git diff --check
```

Run additional checks when touched surfaces require them:

```bash
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
```

## Why this export exists

The earlier export assumed Phase 33 remained blocked. The repo now documents Outcome C, so the next work should shift away from large-import hardening and toward:

1. post-Outcome-C status/docs consistency;
2. repeatability of the public-GTFS local/pilot run;
3. final-root operator request packaging;
4. the next real retained-evidence track.
