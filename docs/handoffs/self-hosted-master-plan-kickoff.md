# Self-Hosted Master Plan Kickoff Prompt For Codex

```text
Read first:
1. AGENTS.md
2. README.md
3. docs/current-status.md
4. docs/handoffs/latest.md
5. docs/handoffs/phase-34.md
6. docs/runbooks/small-agency-pilot-operations.md
7. docs/tutorials/agency-first-run.md
8. docs/tutorials/public-gtfs-local-pilot.md
9. docs/california-readiness-summary.md
10. docs/compliance-evidence-checklist.md
11. docs/evidence/consumer-submissions/status.json

Important direction:
Do not pursue external proof as the default roadmap. The next roadmap is about making the code easy for agencies to adapt and reuse, using the existing OCI/OCL-style pilot server as the reference deployment path.

Immediate task:
Start Phase 35 — README And Roadmap Realignment.

Required:
1. Restore root README.md as the Open Transit RT product front door.
2. Move roadmap-export wording out of root README.
3. Update latest handoff objective so the default next work is self-hosted agency reuse and OCI/OCL reference deployment productization, not external proof.
4. Patch docs/phase-plan.md if it still describes the Phase 34 static validator as only blocked by Java without the later Homebrew Java 17 retry.
5. Preserve all claim boundaries.

Do not change:
- runtime code
- schema
- migrations
- APIs
- consumer statuses
- final-root evidence
- external evidence packets

Required checks:
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null

Confirm all seven consumer targets remain prepared.
```
