# Repo Placement Instructions

Place this roadmap under:

```text
docs/roadmaps/agency-first-connector-platform/
```

Do not place it under:

```text
docs/evidence/
docs/evidence/captured/
docs/evidence/consumer-submissions/
```

This is product roadmap and Codex execution guidance, not evidence.

## Install From ZIP

From repo root:

```bash
unzip /path/to/open-transit-rt-phase61-roadmap-repo-pack.zip -d .
find docs/roadmaps/agency-first-connector-platform -maxdepth 3 -type f | sort
```

## Add Source-Of-Truth Links

Add this sentence to `docs/handoffs/latest.md`, `docs/current-status.md`, and `docs/README.md`:

```md
The Phase 61+ agency-first connector platform roadmap lives at `docs/roadmaps/agency-first-connector-platform/README.md`; Codex should start from `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md` before planning new product phases.
```

Also update wording that says future work is “not Phase 61” to say:

```md
Phases 0 through 60 remain closed. After the Post-60 productization audit, the maintainer authorized Phase 61+ as the forward product roadmap naming for agency-first UI and connector-platform work.
```

## Commit

```bash
git add docs/roadmaps/agency-first-connector-platform docs/handoffs/latest.md docs/current-status.md docs/README.md
git commit -m "Phase 61 -- Checkpoint 000001: add agency-first connector platform roadmap"
```

## Next Codex Prompt

```text
Read `docs/roadmaps/agency-first-connector-platform/00-CODEX-READ-ME-FIRST.md`.

Then plan Phase 61 -- Checkpoint 000002: implement agency-first UI and connector hub.

Use the master/sub-agent model described in the roadmap pack. Plan first. Do not execute until the master review passes. Preserve all protected evidence and consumer-status boundaries. Do not create retained evidence or unsupported claims.
```
