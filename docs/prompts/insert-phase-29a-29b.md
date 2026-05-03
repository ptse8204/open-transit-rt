# Prompt — Insert Phase 29A And Phase 29B Into Track B

Use this prompt for Codex after copying the drafted files into the repo.

```text
Read and obey AGENTS.md first.

Then read:
1. docs/current-status.md
2. docs/handoffs/latest.md
3. docs/handoffs/phase-29.md
4. docs/track-b-productization-roadmap.md
5. docs/phase-29-realtime-quality-expansion.md
6. docs/tutorials/device-avl-integration.md
7. docs/tutorials/device-token-lifecycle.md
8. docs/evidence/device-avl/README.md
9. docs/prompts/calitp-truthfulness.md
10. docs/evidence/redaction-policy.md
11. SECURITY.md
12. docs/dependencies.md
13. docs/decisions.md

This is a docs-only Track B roadmap insertion task.

The user has provided two drafted phase files:
- docs/phase-29a-external-predictor-adapter-evaluation.md
- docs/phase-29b-avl-vendor-adapter-pilot.md

Add those files to the repo and update the roadmap/status docs so Phase 29A and Phase 29B are explicitly placed after Phase 29 and before Phase 30.

Do not implement Phase 29A or Phase 29B yet.
Do not change backend behavior.
Do not change API contracts.
Do not change database schema.
Do not change public feed URLs.
Do not change GTFS-RT protobuf contracts.
Do not change consumer statuses.
Do not add external dependencies.
Do not add external predictor runtime integration.
Do not add vendor AVL runtime integration.
Do not claim production-grade ETA quality, real-world ETA accuracy, certified vendor support, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, agency endorsement, or marketplace/vendor equivalence.

Files to add:
- docs/phase-29a-external-predictor-adapter-evaluation.md
- docs/phase-29b-avl-vendor-adapter-pilot.md

Files to update:
- docs/track-b-productization-roadmap.md
- docs/roadmap-status.md
- docs/current-status.md
- docs/handoffs/latest.md
- docs/handoffs/phase-29.md only if its next-step recommendation still points directly to Phase 30
- docs/phase-30-consumer-submission-execution.md only if it needs a short note that Phase 29A and 29B now precede it

Latest handoff should say:
- Phase 29 is complete for synthetic replay evidence expansion.
- Phase 29A — External Predictor Adapter Evaluation is the next recommended implementation phase.
- Phase 29B — AVL / Vendor Adapter Pilot Implementation follows Phase 29A.
- Phase 30 consumer submission execution remains later and must not advance statuses without target-originated evidence.
- External integrations must remain adapter-bound, tested, optional, and truthfully described.

Run and record:
- make validate
- make test
- git diff --check

Because roadmap/status docs are changing, also run:
- make realtime-quality
- make smoke
- docker compose -f deploy/docker-compose.yml config

Acceptance criteria:
- Phase 29A and Phase 29B files exist.
- Track B roadmap lists Phase 29A and 29B after Phase 29.
- latest handoff recommends Phase 29A next.
- No runtime behavior or dependency status changes are made.
- No unsupported claims are introduced.
```
```
