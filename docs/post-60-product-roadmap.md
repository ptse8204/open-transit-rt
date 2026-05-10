# Post-60 Product Roadmap

## Status

Post-60 productization is the active default roadmap. It is not Phase 61, and
it does not reopen Phases 0 through 60.

The default next work is product quality and external connection maturity:
release-candidate checks, connector contracts, adapter conformance, generic
examples, operator launchpad workflow, and readiness gap reporting. Real-world
evidence tracks remain optional and require separate authorization, intake, and
public-safe retention rules before any evidence tool is run.

## Default Workstream

Post-60 work should make Open Transit RT easier and safer to self-host:

- repeatable release-candidate checks from a fresh clone;
- manifest-bound connector contracts for telemetry, prediction, validators,
  monitoring export, and consumer discovery;
- synthetic adapter conformance tests that run offline;
- generic connector examples with no real credentials or vendor payloads;
- private operator launchpad workflow for setup-to-readiness review;
- `.cache`-only Caltrans-style readiness gap summaries;
- explicit evidence-track routing for later authorized real-world work.

This work improves self-hosted evaluation and integration quality. It is not
consumer acceptance, agency approval, hosted service availability, or
CAL-ITP/Caltrans compliance evidence.

## External Connection Model

Use sidecars, manifests, command adapters, fixtures, and conformance tests.
Do not add arbitrary dynamic plugin loading to the core backend.

Connector classes:

- telemetry source connectors transform external observations before calling
  authenticated `POST /v1/telemetry`;
- prediction connectors stay behind `internal/prediction.Adapter`;
- validator connectors use server-side allowlisted validator IDs;
- monitoring/export connectors emit private redacted diagnostics and do not
  send by default;
- consumer/discovery connectors prepare public URL and metadata workflows
  without automating submissions or changing target statuses.

Vehicle Positions must remain independent of external predictor availability,
and deterministic prediction remains the default Trip Updates path.

## Optional Evidence Tracks

Optional evidence tracks are not default product work. They require retained,
public-safe, claim-specific artifacts before any stronger wording is allowed.

Future tracks include:

- final public root and agency/operator authorization;
- target-originated consumer or aggregator status;
- real pilot closeout and feedback;
- real device, AVL, or vendor integration review;
- deployment operations evidence;
- real-world realtime and ETA quality review.

Without an approved intake, do not run evidence tools, contact external
parties, automate portals, change consumer statuses, or write retained
evidence.

## Claim Boundary

Post-60 productization does not claim CAL-ITP/Caltrans compliance, consumer
submission or acceptance, agency adoption or approval, agency-owned final-root
proof, hosted SaaS, paid support, SLA-backed service, universal deployment
fitness, vendor compatibility, hardware certification, public launch, real
fleet reliability, or production-grade ETA quality.

All seven consumer and aggregator targets remain `prepared` unless future
target-originated evidence supports a specific status transition.
