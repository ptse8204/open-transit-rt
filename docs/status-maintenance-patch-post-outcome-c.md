# Status Maintenance Patch — Post-Outcome-C

This file is the historical checklist used for the Phase 34 stale-docs patch
after Phase 33 Outcome C.

## Patch objective

Make the repo status docs agree on this truth state:

```text
Phase 33 is complete as Outcome C for local/pilot public static GTFS dataset handling.
The large LA Metro public-GTFS import blocker was fixed.
The Outcome C packet exists and is retained under docs/evidence/captured/public-gtfs-local-pilot/2026-05-06/.
External evidence gaps remain unchanged.
```

## Files to patch

Phase 34 completed this maintenance pass. Keep the notes below as historical
implementation guidance for why these files were touched.

### `docs/current-status.md`

Problems to fix:

- It opens by calling the repo an “early-stage starter,” which is misleading after Phases 0-33.
- It should summarize the repo as a technically broad, evidence-bounded open-source transit operations prototype.
- It should preserve all missing external-evidence gaps.

Suggested replacement opening:

```markdown
## Current Repository State

Open Transit RT is a technically broad, evidence-bounded open-source backend prototype for GTFS and GTFS Realtime publication. Phases 0 through 33 are closed for their documented scopes. Phase 33 is complete as Outcome C for local/pilot public static GTFS dataset handling using the LA Metro Bus public GTFS feed.

The repo has substantial local, hosted-pilot, validation, consumer-packet, operations, replay, adapter, and agency-pilot scaffolding. It still does not prove agency adoption, agency-owned final-root readiness, consumer acceptance, Caltrans/CAL-ITP compliance, hosted SaaS availability, production readiness, real vendor AVL compatibility, real realtime data, or production-grade ETA quality.
```

### `docs/roadmap-status.md`

Problems to fix:

- The “Evidence That Exists” section still refers to attempted public GTFS blocker evidence.
- The “What Remains Missing” section still says completed public-GTFS local/pilot evidence is missing.
- Later lines correctly say Phase 33 Outcome C is complete.

Suggested corrections:

- Replace stale attempted-blocker wording with “Phase 33 Outcome C public GTFS
  local/pilot evidence for the May 6, 2026 LA Metro Bus public GTFS run.”
- Remove stale missing-evidence wording for public-GTFS local/pilot evidence.
- Add that static GTFS validation for that packet is blocked by unavailable Java and not a pass.

### `docs/track-b-productization-roadmap.md`

Problems to fix:

- It still says the recommended next phase is Phase 32.
- It does not include Phase 33 in the phase sequence.
- It does not point to the next post-Outcome-C work.

Suggested corrections:

- Mark Phases 22-33 as closed for documented scopes.
- Add Phase 33 to the sequence.
- Add Phase 34 as the recommended next maintenance/evidence-readiness phase.
- Preserve the future evidence forks: final-root proof, real agency pilot, real deployment operations refresh, authorized consumer submission.

### `docs/repo-gaps.md`

Problems to fix:

- It still lists starter scaffolding that already exists.

Acceptable fixes:

1. replace it with a current gaps document; or
2. rename/mark the content as historical; or
3. add a top warning:

```markdown
> Historical note: this file described starter-repo gaps before Phases 0-33. It is no longer the current gap list. See docs/current-status.md, docs/roadmap-status.md, and docs/handoffs/latest.md.
```

A stronger fix is to replace it with current gaps:

- final-root proof;
- real agency pilot;
- consumer evidence;
- real deployment operations refresh;
- static GTFS validator execution for public-GTFS packet;
- repeatable public-GTFS local/pilot guide;
- real device/vendor AVL evidence;
- real-world realtime quality evidence;
- runtime external predictor integration;
- production multi-tenant proof.

### `docs/phase-plan.md`

Problems to fix:

- The phase overview does not show the full modern track or post-Outcome-C continuation.

Suggested corrections:

- Add a note near the top:

```markdown
For the current post-Phase-33 continuation state, use docs/handoffs/latest.md and docs/future-roadmap-post-outcome-c.md. Earlier phase definitions remain historical implementation guidance.
```

- Add Phase 34 to the phase overview.

### `docs/README.md`

Problem to fix:

- The link text “Public GTFS Local/Pilot Attempt” undersells/obscures Outcome C.

Suggested replacement:

```markdown
- [Public GTFS Local/Pilot Outcome C Evidence](evidence/captured/public-gtfs-local-pilot/2026-05-06/README.md)
```

## Files not to patch without external evidence

Do not edit:

```text
docs/evidence/consumer-submissions/status.json
consumer target records
target-specific consumer artifact directories
final-root evidence packet directories
OCI pilot final-root wording
```

## Checks

Run:

```bash
make validate
make test
git diff --check
```

Run additional checks if touched surfaces require them.

## Closure wording

Allowed:

```text
Post-Outcome-C status docs were made consistent. Phase 33 Outcome C remains local/pilot public static GTFS evidence only.
```

Not allowed:

```text
Open Transit RT is agency-adopted, final-root ready, consumer-accepted, compliant, production-ready, or ETA-quality proven.
```
