# Roadmap Status

Active status now lives in:

- [Current Status](current-status.md)
- [Latest Handoff](handoffs/latest.md)
- [External Connector Runtime Roadmap](roadmaps/external-connector-runtime-integration/README.md)
- [Post-rc2 Browser-First Closeout](roadmaps/post-rc2-browser-first-product/closeout.md)
- [Agent History Archive](agent-history-archive.md)

Long historical roadmap status tables were archived on the
`archive/agent-history` branch at commit `ad6cc9c`.

Current product direction:

- reduce the technical capability required to set up, evaluate, and operate
  GTFS / GTFS-Realtime workflows;
- keep normal local evaluation browser-first after startup;
- keep connector support split between works-today/local-supported paths and
  planned/candidate items;
- keep unsupported production, compliance, consumer, vendor, SLA, AVL, and ETA
  claims unsupported.

## Better Software Product-Quality Backlog

Phase 141 merges the browser-first UI reset, external connector runtime plan,
and product-quality work into one implementation track. The executable guard is
`make audit-product-roadmap-baseline`; it checks the private route registry,
public site source, connector examples, release/audit Make targets, roadmap
status categories, and protected consumer tracker boundaries.

Use this grouped backlog for Phases 142 through 160:

| Area | Product-quality direction |
| --- | --- |
| Operator workflow | Keep Start, Help, diagnostics, issue center, access, and maintenance pages focused on the next safe action. |
| GTFS data quality | Improve import preview, active-vs-previous review, validation explanations, and safe recovery guidance without mixing draft and published schedule state. |
| GTFS-RT usefulness | Explain Vehicle Positions, Trip Updates, and Alerts states in private operator language while preserving conservative matching and withheld-output behavior. |
| Connectors | Make CSV, HTTP polling, webhook sidecar, generic transform, predictor, validator, monitoring, and discovery examples dry-run-first, redacted, and no-contact by default. |
| Deployment and observability | Strengthen local/self-hosted doctors, health digests, backup/restore review, and monitoring exports without SLA or hosted-service claims. |
| Security and redaction | Keep tokens, private URLs, raw payloads, database URLs, and private paths out of HTML, logs, support bundles, diagnostics, and screenshots. |
| Release gates | Keep release-candidate diagnostics repeatable, keep `stable` filtering explicit, and keep evidence/consumer status checks separate from software capability. |
