# Phase 34 Kickoff Prompt — Post-Outcome-C Status Consistency And Evidence Readiness

Use this as the copy-paste prompt for a fresh Codex agent.

```text
Implement Phase 34 — Post-Outcome-C Status Consistency And Evidence Readiness for Open Transit RT.

Current repo baseline:
Phases 0 through 33 are closed for their documented scopes. Phase 33 is complete as Outcome C — public-GTFS local/pilot run completed with public-safe retained summaries. The LA Metro Bus public GTFS run completed locally, the large import blocker was fixed, all five public paths were fetched, and the fetched schedule was verified as the imported LA Metro public GTFS rather than the repo sample feed.

Important boundary:
Phase 33 proves only local/pilot handling of a real public static GTFS dataset. It does not prove agency adoption, agency endorsement, agency approval, official agency feed status, agency-owned final-root proof, consumer submission/review/acceptance, consumer ingestion/listing/display, Caltrans/CAL-ITP compliance, hosted SaaS availability, production readiness, production multi-tenant hosting, real vendor AVL compatibility, real LA Metro realtime data, real-world ETA accuracy, or production-grade ETA quality.

Read first:
1. AGENTS.md
2. README.md
3. docs/current-status.md
4. docs/handoffs/latest.md
5. docs/handoffs/phase-33.md
6. docs/phase-33-public-gtfs-local-pilot-evidence.md
7. docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md
8. docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/command-log-inventory-2026-05-06.md
9. docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/retained-summaries-2026-05-06.md
10. docs/roadmap-status.md
11. docs/track-b-productization-roadmap.md
12. docs/repo-gaps.md
13. docs/phase-plan.md
14. docs/evidence/consumer-submissions/status.json
15. docs/agency-owned-domain-readiness.md
16. docs/compliance-evidence-checklist.md
17. docs/evidence/redaction-policy.md

Goal:
Patch post-Outcome-C status inconsistencies, add missing operator-facing final-root request guidance, add a repeatable public-GTFS local/pilot guide, and preserve all evidence boundaries.

Required implementation:
1. Update docs/current-status.md so it no longer opens by calling the project an early-stage starter without qualification. It should reflect that Phase 33 Outcome C is complete while preserving all missing external-evidence gaps.
2. Update docs/roadmap-status.md so it no longer says public-GTFS local/pilot evidence is only an attempted blocker or still missing. It should describe the Phase 33 Outcome C packet and its limits.
3. Update docs/track-b-productization-roadmap.md so it no longer recommends Phase 32 as next. It should show Phases 22 through 33 as closed and Phase 34 or the next retained-evidence fork as current.
4. Refresh, retire, or clearly mark docs/repo-gaps.md as historical. It must not list already-completed starter scaffolding as current missing work.
5. Update docs/phase-plan.md to point future agents to the post-Outcome-C roadmap and latest handoff.
6. Update docs/README.md so the public-GTFS evidence link is labeled as Outcome C evidence rather than merely an attempt.
7. Add docs/final-root-operator-request.md as a plain-language request package for obtaining agency-owned or agency-approved final public feed root evidence. This is not an evidence packet.
8. Add docs/tutorials/public-gtfs-local-pilot.md as a repeatable guide for rerunning a public-GTFS local/pilot evidence flow without implying agency approval.
9. Clarify the Phase 33 static GTFS validator blocker: Java was unavailable, so static validation did not execute. Do not claim a static validator pass unless real retained validator evidence exists.
10. Add docs/handoffs/phase-34.md using the handoff template. Record changed files, checks, blockers, evidence/claim boundaries, and exact next recommendation.
11. Update docs/handoffs/latest.md to make Phase 34 closed only for the status/evidence-readiness scope and to identify the next retained-evidence forks.

Do not change unless retained external evidence exists:
- docs/evidence/consumer-submissions/status.json
- consumer target records
- target-specific consumer artifact directories
- final-root evidence packets
- OCI pilot final-root wording

Do not claim:
- agency adoption
- agency endorsement
- agency approval
- official agency feed status
- agency-owned final-root proof
- consumer submission/review/acceptance
- consumer ingestion/listing/display
- Caltrans/CAL-ITP compliance
- hosted SaaS availability
- paid support/SLA coverage
- production readiness
- production multi-tenant hosting
- real vendor AVL compatibility
- real LA Metro realtime data
- real-world ETA accuracy
- production-grade ETA quality
- public launch completion

Required checks:
make validate
make test
git diff --check

Run if relevant:
make realtime-quality
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
make agency-app-up
make agency-app-down

Closure criteria:
- Status/roadmap docs are consistent with Phase 33 Outcome C.
- Final-root operator request package exists and is clearly not evidence.
- Public-GTFS local/pilot repeatability guide exists.
- Static validator blocker is clear and not overclaimed.
- Consumer statuses remain unchanged.
- No final-root, consumer, agency, compliance, production, vendor, realtime-data, or ETA-quality claims are added.
- Latest handoff identifies next retained-evidence forks.
```
