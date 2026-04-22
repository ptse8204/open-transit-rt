# Phase 12 Step 2 Evidence Packet: local-demo

- Environment: `local-demo`
- Capture date (UTC): 2026-04-22
- Operator: Codex local run on developer workstation
- Evidence scope: local Docker/Postgres plus local Go services and temporary loopback public proxy
- Claim status: partial Phase 12 evidence only

## Claim Boundary

This packet is real evidence for one actual local environment. It is not hosted production evidence.

Proven here:

- repo-supported local demo can import GTFS, expose public feed paths through a loopback proxy, fetch all five public feed artifacts anonymously over HTTP, run validation endpoints, export a scorecard, and restore a Postgres dump into an isolated restore database;
- protected admin/debug routes reject anonymous requests in this local environment;
- backup/restore mechanics work for the local Postgres database copy used in this drill.

Not proven here:

- public HTTPS hosting;
- DNS ownership or stable public hostname;
- TLS certificate validity or renewal;
- HTTP-to-HTTPS redirect behavior;
- production monitoring dashboards, alert rules, notification delivery, or alert lifecycle;
- clean canonical validator results;
- production backup schedule or retention;
- URL permanence across production publish and rollback;
- third-party consumer or aggregator acceptance.

## Files

- `public-feed-proof-2026-04-22.md`
- `reverse-proxy-tls-2026-04-22.md`
- `validator-record-2026-04-22.md`
- `monitoring-alert-2026-04-22.md`
- `backup-restore-drill-2026-04-22.md`
- `scorecard-export-2026-04-22.md`
- `commands-run-2026-04-22.md`
- `SHA256SUMS.txt`
- `artifacts/`

## Phase 12 Closure

Phase 12 is not fully closed by this packet. It advances Step 2 by recording real local evidence and exact blockers, but the phase still needs a hosted HTTPS environment with production validator records, monitoring/alert lifecycle evidence, deployment backup policy, and URL permanence proof.
