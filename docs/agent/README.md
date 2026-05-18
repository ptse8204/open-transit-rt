# AI-Agent Documentation

This directory is the entry point for Codex and other AI-agent maintenance
context. Human readers should usually start with the [Docs Index](../index.md)
or the [README](../../README.md) instead.

Long project-history files are preserved on the `archive/agent-history` branch.
They are not the normal product docs for agency staff.

## Agent Read Order

For a new Codex session, read in this order:

1. [AGENTS.md](../../AGENTS.md)
2. [Codex Task Brief](../codex-task.md)
3. [Current Status](../current-status.md)
4. [Latest Handoff](../handoffs/latest.md)
5. [External Connector Runtime Roadmap](../roadmaps/external-connector-runtime-integration/README.md)
6. [Post-rc2 Browser-First Closeout](../roadmaps/post-rc2-browser-first-product/closeout.md)
7. [Decisions](../decisions.md)
8. [Dependencies](../dependencies.md)
9. [Final Claim Review Rules](../phase-60-final-claim-review-and-public-closeout.md)

## Agent-Only Areas

- [Handoffs](handoffs/README.md): latest continuation note and archive pointer.
- [Roadmaps](roadmaps/README.md): current roadmap pointers and archive pointer.
- [Prompts](prompts/README.md): current prompt context and archive pointer.
- [History](history/README.md): archive branch and stale-doc handling.

## Stable Branch Boundary

The `stable` branch is filtered for product/user-facing source and docs. It
intentionally omits `AGENTS.md`, `docs/agent/**`, `docs/handoffs/**`,
`docs/prompts/**`, `docs/roadmaps/**`, Codex task files, conversation
summaries, phase ledgers, and roadmap planning packs.

The filtered stable branch still omits agent-only paths. Run this local check
before changing the filter workflow or exclude list:

```bash
make check-stable-filter
```

That check simulates the filtered tree, verifies preserved product paths,
verifies AI-agent-only paths are absent, checks dry-run/no-force-push workflow
guards, and confirms the consumer tracker remains seven prepared-only targets.

## Human-Facing Boundary

Do not make agency users read this area to understand the product. When an
agent updates docs, keep human-facing changes in:

- [Docs Index](../index.md)
- [Tutorials](../tutorials/README.md)
- [Connector Catalog](../connectors/catalog.md)
- [Wiki Home](../../wiki/README.md)
- [Static Site](../../site/index.html)

## Avoiding Stale-Doc Drift

When implementation changes behavior:

- update the human guide that teaches the workflow;
- update the relevant roadmap or handoff only as history/status;
- update [Decisions](../decisions.md) for architecture-significant choices;
- preserve unsupported-claim wording;
- never move consumer tracker status without retained target-originated
  evidence.

This area indexes AI-agent records. It does not prove CAL-ITP/Caltrans
compliance, production readiness, consumer acceptance, agency adoption, hosted
service availability, vendor compatibility, SLA/uptime, production AVL
reliability, or ETA quality.
