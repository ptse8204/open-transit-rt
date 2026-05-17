# External Connector Runtime Integration Roadmap

This roadmap is the focused product-quality path after the post-`v0.1.0-rc.2`
browser-first work. It moves Open Transit RT from connector documentation,
synthetic examples, and local conformance checks toward stronger runtime
support for deployment-owned external connectors.

The goal is better self-hosted GTFS / GTFS-Realtime software with clear
external-connection boundaries. It is not an evidence track and it does not
prove named vendor compatibility, hardware certification, CAL-ITP/Caltrans
compliance, production readiness, consumer acceptance, hosted service
availability, SLA/uptime, production AVL reliability, production-grade ETA
quality, or real-world ETA accuracy.

## Start Here

Use these files together:

- [Phase Plan](phase-plan.md): implementation phases, acceptance checks, and
  validation commands.
- [Connector Catalog](../../connectors/catalog.md): current safe connector
  shapes and first local checks.
- [Integration Adapter Kit](../../integration-adapter-kit.md): current adapter
  contracts, redaction rules, and boundary notes.
- [External Connection Readiness](../../external-connection-readiness.md):
  claim-bounded readiness questions for connector work.

## Product Direction

Open Transit RT should make these connector paths practical for small agencies
and technical helpers while keeping external systems outside the core domain:

- Vehicle / GPS / AVL ingest through CSV replay, HTTP polling, webhook
  sidecars, generic JSON transforms, and authenticated `POST /v1/telemetry`.
- External HTTP predictor shadow mode that compares results without mutating
  public Trip Updates until explicitly enabled by deployment-owned config.
- Monitoring/export helpers that are redacted and no-send by default.
- Offline connector conformance that exercises failure modes with committed
  synthetic fixtures.
- Browser support for connector setup, health review, safe commands to run
  outside the browser, and warnings about secrets and private payloads.

## Phase Map

| Phase | Name | Main outcome |
| --- | --- | --- |
| 01 | Runtime Boundary Baseline | Freeze connector contracts, config ownership, redaction rules, and local validation inventory. |
| 02 | Vehicle / GPS / AVL Runtime Adapters | Improve CSV replay, HTTP polling, webhook sidecar, and generic JSON transform runtime paths toward authenticated telemetry ingest. |
| 03 | Authenticated Telemetry And Vehicle Positions Hardening | Tighten `/v1/telemetry` diagnostics and keep Vehicle Positions conservative and independent of external predictors. |
| 04 | External HTTP Predictor Shadow Mode | Make predictor config, health checks, timeout/malformed-output handling, fail-closed behavior, and comparison diagnostics repeatable. |
| 05 | Monitoring / Export Helper Runtime | Add deployment-owned export configuration, redacted summary formats, no-send defaults, and local delivery dry-runs. |
| 06 | Connector Conformance Expansion | Broaden adapter-conformance fixtures and fixture coverage across telemetry, prediction, monitoring, validator, and consumer/discovery boundaries. |
| 07 | Browser Connector Setup And Health | Improve private Connector Workbench, connector health review, and copyable deployment-owned config checklists. |
| 08 | Runtime Roadmap Closeout And Release Gate | Run connector-specific validation, update docs/status, and decide the next safe product-quality slice. |

## Success Definition

This roadmap succeeds when a technical helper can configure and test local or
deployment-owned connector runtimes with:

- clear runtime contracts and failure modes;
- no arbitrary command execution from the browser;
- no default network contact to external systems from tests;
- no committed secrets, private payloads, private endpoint URLs, or private
  infrastructure paths;
- local fixture and conformance coverage for safe first checks;
- browser-visible setup and health guidance for agency staff;
- public feed behavior that fails closed when connector inputs are missing,
  malformed, stale, low-confidence, or unavailable.

## Claim Boundary

Allowed wording:

- Open Transit RT supports local/self-hosted evaluation.
- Open Transit RT supports CAL-ITP-style readiness workflows.
- Open Transit RT provides synthetic/local connector examples and conformance
  checks.
- Open Transit RT provides deployment-owned connector runtime boundaries.

Unsupported without future retained evidence:

- vendor compatibility or hardware certification;
- production AVL reliability;
- production-grade ETA quality or real-world ETA accuracy;
- CAL-ITP/Caltrans compliance;
- consumer submission, review, acceptance, ingestion, listing, or display;
- agency adoption or approval;
- final-root readiness;
- hosted service availability, paid support, SLA, or uptime;
- production readiness.

## Protected Paths

This roadmap must not write to:

- `docs/evidence/captured/**`
- `docs/evidence/consumer-submissions/current/**`
- `docs/evidence/consumer-submissions/artifacts/**`
- `docs/evidence/consumer-submissions/packets/**`
- `docs/evidence/consumer-submissions/status.json`

The consumer tracker must remain exactly seven prepared-only targets unless a
separately authorized evidence workflow provides retained target-originated
evidence.
