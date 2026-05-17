# Agent History

This repo has a long agent-maintained history layer. Phase 11 keeps that layer
discoverable while moving normal readers toward shorter browser-first docs.

## Inventory

Phase 11 inventory found:

- 125 root `docs/phase-*.md` implementation ledgers;
- 146 files in `docs/handoffs/`;
- roadmap packs under `docs/roadmaps/`;
- Codex and kickoff prompt files under `docs/prompts/` and roadmap folders.

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

- do not delete it just to simplify the tree;
- prefer adding a short redirect note from current docs;
- update current human docs first;
- update `docs/handoffs/latest.md` only for current continuation state;
- keep protected evidence paths and the prepared-only consumer tracker
  unchanged unless a maintainer explicitly opens an evidence track.
