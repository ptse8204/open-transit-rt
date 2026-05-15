# Public Docs And Site Freeze Checklist

Use this checklist before changing README, wiki pages, docs home, public site
copy, screenshots, diagrams, or contributor-facing onboarding.

This checklist is a documentation control only. It does not publish a site,
create evidence, move consumer status, tag a release, or prove public launch.

## Freeze Rules

- Start from the browser-first path: private Operations Console, Start Here,
  setup, GTFS, five public feed URLs, validation, realtime, connectors,
  maintenance, Help.
- Keep public pages task-based; do not make phase history the first thing a
  small-agency reader sees.
- Link deeper maintainer history through `docs/current-status.md`,
  `docs/roadmap-status.md`, and `docs/handoffs/latest.md`.
- Keep the seven consumer targets described as `prepared` only unless retained
  target-specific evidence authorizes a narrower status update.
- Keep screenshots and diagrams local/demo documentation aids only.
- Do not add site-publish, social-post, release, or external-contact steps to
  ordinary docs edits.

## Required Public Boundaries

Public docs must not claim:

- CAL-ITP/Caltrans compliance;
- agency adoption, endorsement, deployment success, or approval;
- consumer submission, review, ingestion, listing, display, or acceptance;
- final-root readiness;
- hosted SaaS, paid support, SLA, or uptime;
- production readiness;
- vendor compatibility or hardware certification;
- production-grade ETA quality or real-world ETA accuracy;
- release readiness unless a release gate explicitly supports only the scoped
  wording used.

Prefer phrasing such as:

- "supports local evaluation";
- "implements technical foundations";
- "private diagnostic";
- "prepared packet";
- "synthetic/local fixture";
- "authorization-gated evidence track";
- "not evidence by itself."

## Screenshot And Diagram Rules

- Store product screenshots only under `docs/assets/product-screenshots/`.
- Store diagrams or diagram sources under `docs/assets/product-diagrams/` or the
  reviewed docs asset locations.
- Do not store screenshots, diagrams, or media under `docs/evidence`.
- Use synthetic, local-demo, or public-safe data only.
- Do not show tokens, private URLs, private hostnames, private device IDs,
  private emails, private logs, raw telemetry, database URLs, portal records,
  or IPs.
- Caption every product screenshot with "local/demo product screenshot".
- Leave missing screenshots as `not captured`; never fake a screenshot.

## Contributor-Onboarding Freeze

Contributor docs should direct first-time contributors toward:

- docs corrections;
- tutorial troubleshooting improvements;
- synthetic fixtures;
- small focused tests;
- connector manifests and examples using synthetic inputs;
- redaction and claim-boundary cleanup.

First PRs should avoid:

- migrations;
- public feed contract changes;
- auth or secret-handling changes;
- consumer tracker or protected evidence paths;
- real agency/vendor/device data;
- release actions;
- broad public-claim wording.

## Validation

Run at least:

```bash
git diff --check
make check
make audit-product-acceptance
make audit-final-claim-review
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Also run the exact prepared-only consumer tracker assertion and protected-path
status check before merging docs/site/contributor changes.
