# Phase 35 — README And Roadmap Realignment

## Goal

Fix the repo’s first impression and align the roadmap with self-hosted agency reuse instead of external-proof chasing.

## Required work

1. Restore root `README.md` as the Open Transit RT product front door.
2. Move roadmap-export wording out of root README.
3. Update `docs/handoffs/latest.md` current objective:
   - focus on self-hosted reusable code;
   - use OCI/OCL pilot server as reference deployment;
   - stop listing external proof as the default next step.
4. Patch `docs/phase-plan.md` Phase 34 validator wording if stale.
5. Preserve claim boundaries.

## README must include

- what this is;
- who it is for;
- what it does;
- quick local run;
- public GTFS onboarding pointer;
- OCI/OCL reference deployment pointer;
- CAL-ITP-style readiness pointer;
- integration points;
- what this is not.

## Do not change

- runtime code;
- schema;
- migrations;
- consumer statuses;
- final-root evidence;
- external evidence packets.

## Checks

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```
