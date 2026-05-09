# Repo Gaps

This document is the current gap list after Phase 34.

Historical note: the original version of this file described starter-repo gaps
before Phases 0 through 33. Those starter scaffolding items are no longer the
current missing work. The repo now has `.env.example`, `cmd/migrate`,
one-command local bootstrap flows, integration fixtures, `docs/decisions.md`,
`docs/dependencies.md`, a Makefile, and a Taskfile.

Use these files for current status before starting work:

- `docs/current-status.md`
- `docs/roadmap-status.md`
- `docs/handoffs/latest.md`
- `docs/future-roadmap-post-outcome-c.md`

Phase 54 refreshed official-source requirement mappings only. Phase 55 added a
local compliance/readiness packet generator and audit guard only. Neither phase
closed any deployment, consumer, final-root, operations, vendor, SLA, or
ETA-quality evidence gap.

## Current Evidence Gaps

- Agency-owned or agency-approved final public feed root proof.
- Retained final-root approval, DNS, TLS, redirect, public fetch, validator,
  proxy/config, README, and checksum evidence.
- Authorized target-specific consumer submission evidence.
- Consumer review, acceptance, ingestion, listing, display, rejection, or
  blocker evidence from named targets.
- Real agency pilot evidence and agency feedback.
- Real deployment operations refresh evidence beyond the existing OCI pilot
  packet.
- Real device or vendor AVL evidence beyond local simulator/no-hardware and
  synthetic adapter fixtures.
- Real-world observed-arrival/departure ETA accuracy evidence.
- Real route/time-period realtime quality metrics.
- Production-grade ETA quality evidence.
- Validator-clean or no-warning static GTFS evidence for the Phase 33
  public-GTFS packet. The original Outcome C static validator attempt was
  blocked because Java was unavailable; a later post-Phase-34 retry executed
  with Homebrew Java 17 and reported process exit code `0`, system error count
  `0`, and 3 warning notices. That retry is not validator-clean or compliance
  evidence.
- Provider or agreed regional source-of-truth website proof listing the
  canonical GTFS Schedule and all three GTFS Realtime feed links for a final
  deployment.
- Public technical contact or feed-contact proof for a final deployment,
  including feed-contact fields where supported by the active GTFS feed.
- Transitland and Mobility Database availability proof for the exact final feed
  scope if those discoverability claims are made.
- Realtime API-key registration proof if a deployment chooses authenticated
  realtime access.

Phase 55 can summarize these gaps into ignored `.cache` review packets, but
those generated summaries are not retained evidence and do not close any gap by
themselves.

## Current Product And Operations Gaps

- Agency-owned final-root deployment flow backed by retained evidence.
- Hosted login/SSO and server-side admin JWT `jti` replay tracking.
- Full operator UI for manual override workflows.
- Production SLO dashboards and alerting beyond current lightweight
  feed-monitor examples, Operations Console pages, request logs, request IDs,
  readiness checks, and optional `/metrics`.
- OpenTelemetry tracing/exporter wiring and Prometheus/Grafana deployment
  assets.
- Runtime external predictor integration such as TheTransitClock behind the
  documented adapter boundary.
- External consumer submission API integrations, if ever explicitly approved
  and backed by current target documentation.
- Production multi-tenant hosting proof.
- Marketplace/vendor-equivalent service packaging and support commitments.

## Claim Boundary

These gaps mean the repo must not claim agency adoption, agency endorsement,
agency approval, official agency feed status, agency-owned final-root proof,
consumer submission/review/acceptance, consumer ingestion/listing/display,
Caltrans/CAL-ITP compliance, hosted SaaS availability, paid support/SLA
coverage, production readiness, production multi-tenant hosting, real vendor
AVL compatibility, real-world ETA accuracy, production-grade ETA quality, or
public launch completion without retained evidence for that specific claim.
