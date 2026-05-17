# AI-Agent Documentation

This directory is the entry point for Codex and other AI-agent maintenance
context. Human readers should usually start with the [Docs Index](../index.md)
or the [README](../../README.md) instead.

The repository keeps long project-history files because they preserve release,
roadmap, and claim-boundary decisions. They are not the normal product docs for
agency staff.

## Agent Read Order

For a new Codex session, read in this order:

1. [AGENTS.md](../../AGENTS.md)
2. [Codex Task Brief](../codex-task.md)
3. [Current Status](../current-status.md)
4. [Latest Handoff](../handoffs/latest.md)
5. [Post-rc2 Browser-First Product Roadmap](../roadmaps/post-rc2-browser-first-product/README.md)
6. [Phase Plan](../roadmaps/post-rc2-browser-first-product/phase-plan.md)
7. [Decisions](../decisions.md)
8. [Dependencies](../dependencies.md)
9. [Final Claim Review Rules](../release-status-v0.1.0-rc.2.md)

## Agent-Only Areas

- [Handoffs](handoffs/README.md): latest and historical continuation notes.
- [Roadmaps](roadmaps/README.md): autonomous roadmap packs, phase plans, and
  internal status ledgers.
- [Prompts](prompts/README.md): Codex kickoff prompts and truthfulness guards.
- [History](history/README.md): phase-ledger inventory and stale-doc handling.

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
