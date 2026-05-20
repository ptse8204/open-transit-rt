# Phase 141-160 Better Software Roadmap Closeout

Date: 2026-05-20

This closeout summarizes the better-software roadmap that moved Open Transit RT
from browser-first local evaluation toward a more durable open-source GTFS and
GTFS-Realtime operations platform. It records software capability only. It is
not retained evidence, not a release publication, not compliance proof, not
consumer acceptance, not agency adoption, not hosted service availability, not
vendor compatibility, and not production readiness.

## Completed Commits

- `f67c014` Phase 141: executable product-quality baseline and grouped backlog.
- `38f36a3` Phase 142: primary Operations Console UI vocabulary/layout cleanup.
- `d68fac2` Phase 143: unified private operator issue center.
- `a3328c4` Phase 144: GTFS import preview, diff, and recovery review.
- `ae3fe5c` Phase 145: validation explanations and safe fix guidance.
- `b297421` Phase 146: telemetry ingest and device onboarding hardening.
- `5c93722` Phase 147: Vehicle Positions usefulness diagnostics.
- `a1cb9e6` Phase 148: Trip Updates shadow diagnostics and withheld reasons.
- `199f1de` Phase 149: Alerts and disruption workflow improvements.
- `2a89853` Phase 150: connector runtime pack and synthetic examples.
- `e9854fd` Phase 151: private connector health and setup review.
- `341aab3` Phase 152: redacted monitoring and reliability exports.
- `532a4a8` Phase 153: self-hosted install, upgrade, backup, and restore UX.
- `1c30aee` Phase 154: agency isolation and role guidance hardening.
- `7aa0aa9` Phase 155: public feed discovery and sharing preparation.
- `9be732f` Phase 156: support bundle redaction and capture safety.
- `4697690` Phase 157: nontechnical staff help and role quick paths.
- `42294f5` Phase 158: API, feed, connector, and adapter contract checks.
- `bec3012` Phase 159: release-candidate gate refresh for product-quality checks.

Phase 160 is this closeout, status, and handoff update.

## Product Capability Added

- Operator workflow: the console now has a prioritized issue center, role-based
  help paths, clearer private setup and maintenance guidance, and less internal
  audit vocabulary in primary pages.
- GTFS data quality: import preview, active-vs-previous review, rollback
  guidance, validation explanations, owner categories, safe fix paths, and
  malformed report handling are stronger.
- GTFS-Realtime usefulness: Vehicle Positions, Trip Updates, and Alerts expose
  clearer withheld, stale, unmatched, suppressed, shadow, and lifecycle signals
  without implying false trip certainty or ETA quality.
- Connectors: examples, manifests, adapter conformance, dry-run behavior,
  redacted diagnostics, and private connector health review are more usable
  while remaining synthetic and no-contact by default.
- Deployment and monitoring: local deployment doctor, recovery guidance,
  redacted health summaries, operations notification drafts, reliability
  exports, and release-candidate gates are tied to executable checks.
- Security and redaction: support bundle output, capture state handling,
  tenant-safe route behavior, agency path rejection, no-store behavior, role
  guidance, and protected-status audits are stricter.
- Contracts and release readiness: `/v1/telemetry`, public feed paths,
  `/public/feeds.json`, private JSON companions, connector manifests, adapter
  conformance fixtures, prediction DTOs, stable-filter rules, and release gates
  now have explicit docs and tests.

## Validation Snapshot

Phase 159 ran the broad release-candidate gate on 2026-05-20. The refreshed
release-candidate diagnostic under local `.cache` produced 56 passed rows, zero
blockers, one dirty-worktree `needs_review` row, and four follow-up rows that
were run separately: `make validate`, `make test`, `make smoke`, and package
audit. The local app fetch row remains opt-in for the helper; default
install-confidence clone/bootstrap checks passed separately.

Phase 160 reran the closeout validation set after the status and handoff
updates. The final local release-candidate diagnostic again reported 56 passed
rows, zero blockers, the expected dirty-worktree `needs_review` row, and the
same four helper follow-up rows, with the follow-up commands run directly.

The following commands passed during the Phase 159 and Phase 160 gate review:

- `git diff --check`
- `go test ./...`
- `make check`
- `make test`
- `make smoke`
- `make validate`
- `make check-links`
- `make product-ui-smoke`
- `make audit-product-language`
- `make audit-ui-layout`
- `make audit-product-acceptance`
- `make audit-final-claim-review`
- `make audit-operations-route-inventory`
- `make external-connection-check`
- `make adapter-conformance`
- `make test-connector-examples`
- `make gtfsrt-conformance`
- `make api-contract-check`
- `make check-stable-filter`
- `make install-confidence`
- `make test-release-candidate-check`
- `make test-release-package`
- local release package generation and `make audit-release-package`
- `docker compose -f deploy/docker-compose.yml config`
- `scripts/check-consumer-tracker.sh`
- `python3 -m json.tool docs/evidence/consumer-submissions/status.json`
- `git diff --exit-code -- docs/evidence/consumer-submissions/status.json`

`make validate` also stopped checking obsolete phase-ledger files and now checks
current durable release, API, connector, deployment, roadmap, and tutorial
contracts.

## Protected Boundaries

No phase in this roadmap intentionally edited `docs/evidence/**`.
`docs/evidence/consumer-submissions/status.json` remained unchanged, and the
consumer tracker remained exactly seven `prepared` targets. The roadmap did not
contact agencies, vendors, consumers, portals, live validators, or live connector
services.

Unsupported claims remain unsupported: CAL-ITP/Caltrans compliance, production
readiness, agency adoption, consumer acceptance, final-root readiness, hosted
service availability, vendor compatibility, hardware certification, SLA/uptime,
production AVL reliability, production-grade ETA quality, and real-world ETA
accuracy.

## Recommended Next Track

The best next software track is deeper realtime correctness:

- observed-arrival and departure evaluation harnesses;
- delay propagation and cancellation pairing behavior;
- frequency-service and after-midnight edge cases;
- block continuity and repeated trip instances;
- conservative Trip Updates quality measures that keep withheld and unknown
  states visible.

This is the highest-leverage next step because the roadmap strengthened import,
connector, operator, deployment, and release surfaces, while realtime correctness
still determines whether agencies can safely publish richer GTFS-Realtime beyond
Vehicle Positions. GTFS extensions, real connector runtime hardening with
authorized data, or release publication can follow, but real connector evidence
and release publication both require explicit maintainer authorization.
