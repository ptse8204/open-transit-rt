# Agent History

This repo has a long agent-maintained history layer. Main now keeps only the
active maintenance context. The long archive is preserved on
`archive/agent-history`.

## Inventory

Before archive, the inventory included:

- 125 root `docs/phase-*.md` implementation ledgers;
- 146 files in `docs/handoffs/`;
- roadmap packs under `docs/roadmaps/`;
- Codex and kickoff prompt files under `docs/prompts/` and roadmap folders.

Use `git switch archive/agent-history` to inspect those files locally, then
return to `main` before normal product work.

## How To Use History

Use history files to answer questions like:

- why a boundary exists;
- when a route or workflow was introduced;
- which validation commands were run for a phase;
- which unsupported claims must stay unsupported.

Do not require agency staff, technical helpers, or connector developers to read
history files before they can run the product.

## Stale-Doc Handling

When a historical file is stale:

- do not restore it to main just to simplify archaeology;
- prefer a short pointer to `archive/agent-history`;
- update current human docs first;
- update `docs/handoffs/latest.md` only for current continuation state;
- keep protected evidence paths and the prepared-only consumer tracker
  unchanged unless a maintainer explicitly opens an evidence track.
