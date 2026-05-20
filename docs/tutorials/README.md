# Tutorials

Start here instead if you are choosing a guide by role:
[Docs Index](../index.md). These command-level tutorials remain detailed
references for technical helpers and maintainers. Public-facing versions live
in [`/wiki`](../../wiki/README.md).

They document what the current repo can run today. They do not claim hosted production readiness, consumer acceptance, or completed CAL-ITP/Caltrans compliance.

![Exact-behavior local quickstart flow showing database bootstrap, validator install, sample GTFS import, service startup, publication bootstrap, telemetry ingest, public feed fetches, validation run, and scorecard inspection.](../assets/quickstart-flow.png)

*Exact-behavior flow diagram for `make demo-agency-flow`, rendered from a reviewed SVG spec.*

## Choose By Role

| Role | Start with | Use when |
| --- | --- | --- |
| Agency staff | [Small-Agency Acceptance Script](small-agency-acceptance-script.md) and `/admin/operations/help` | Staff need a browser-first walkthrough, next actions, and clear limits without running commands. |
| Administrator | [Agency First Run](agency-first-run.md), [GTFS Validation Triage](gtfs-validation-triage.md), and [Operator Smoke And Support Bundle](operator-smoke-and-support-bundle.md) | Setup, import, validators, device tokens, and redacted diagnostics need an owner. |
| Deployment owner | [Self-Hosted Operator Trial](self-hosted-operator-trial.md), [Reference Deployment Doctor](../deployment/reference-deployment-doctor.md), and [Small-Agency Maintenance Guide](small-agency-maintenance-guide.md) | The host, public base URL, backup/restore, upgrade, monitoring, or recovery posture is the blocker. |
| Integrator | [Integration Adapter Kit](../integration-adapter-kit.md), [External Adapter Conformance](external-adapter-conformance.md), and [Device And AVL Integration](device-avl-integration.md) | Telemetry, connector, predictor, validator, monitoring, or discovery contracts need local/synthetic review. |

## Agency Operations Cockpit / Start Here

- [Agency First Run](agency-first-run.md): start the full local app package and understand the outputs.
- [Reusable Agency Onboarding](reusable-agency-onboarding.md): provide an agency ID and GTFS URL for local/reference onboarding without manual database edits.
- [Self-Hosted Operator Trial](self-hosted-operator-trial.md): run one guided local/reference evaluation across onboarding, public feed checks, readiness review, validators, and the synthetic AVL dry-run.
- [Agency Launchpad](agency-launchpad.md): run a private authenticated launchpad workflow across setup, GTFS, metadata, feeds, telemetry, validators, readiness, connector conformance, support bundle, and decision gate.
- [Operator Smoke And Support Bundle](operator-smoke-and-support-bundle.md): run strict smoke checks and collect redaction-safe diagnostics without creating evidence.
- [Reference Deployment Doctor](../deployment/reference-deployment-doctor.md): run read-only reference deployment diagnostics without creating evidence.
- [Self-Hosted Operations Notifications](self-hosted-operations-notifications.md): draft a private local notification summary from existing diagnostics without sending anything.
- [Telemetry Simulator And Device Trial](telemetry-simulator-and-device-trial.md): send synthetic telemetry through real authenticated ingest for local/reference diagnostics.
- [Prediction And ETA Lab](prediction-eta-lab.md): review private deterministic, shadow, withheld-output, and aggregate backtest diagnostics without making ETA-quality claims.
- [Staff Training Demo Kit](staff-training-demo-kit.md): run local/synthetic role paths, demo scenarios, recovery drills, trainer scripts, and technical-helper handoffs without creating evidence or adoption claims.
- [Video Recording Guide](video-recording-guide.md): record short local/demo tutorials without secrets, raw private data, or unsupported claims.
- [Real Agency GTFS Onboarding](real-agency-gtfs-onboarding.md): prepare, validate, review, and publish a real agency GTFS ZIP safely.
- [Public GTFS Local/Pilot Runbook](public-gtfs-local-pilot.md): repeat a real public GTFS local/pilot run without implying agency approval.
- [GTFS Validation Triage](gtfs-validation-triage.md): understand common import and validation failures and use the authenticated admin GTFS quality triage UI safely.
- [Integration Adapter Kit](../integration-adapter-kit.md): choose the correct telemetry, predictor, validator, monitoring, or consumer workflow boundary.
- [Device And AVL Integration](device-avl-integration.md): send authenticated telemetry from devices, vendors, adapters, or no-hardware simulator flows.
- [Device Token Lifecycle](device-token-lifecycle.md): rotate, rebind, store, and troubleshoot device bearer credentials safely.
- [Local Quickstart](local-quickstart.md): bring up the local development environment.
- [Agency Demo Flow](agency-demo-flow.md): run the executable agency/evaluator demo.
- [Deploy With Docker Compose](deploy-with-docker-compose.md): understand the current deployment path.
- [Production Checklist](production-checklist.md): review operational work still needed for a real deployment.
- [Small-Agency Pilot Operations](../runbooks/small-agency-pilot-operations.md): run the pilot operations profile.
- [OCI/OCL Reference Deployment](../deployment/oci-reference-deployment.md): reuse the self-hosted reference server pattern.
- [CAL-ITP Readiness Checklist](calitp-readiness-checklist.md): track readiness without overclaiming compliance.

For self-hosted deployment planning, see
[OCI/OCL Reference Deployment](../deployment/oci-reference-deployment.md),
[Self-Hosted Operator Trial](self-hosted-operator-trial.md),
[Operator Smoke And Support Bundle](operator-smoke-and-support-bundle.md),
[Reference Deployment Doctor](../deployment/reference-deployment-doctor.md),
[Self-Hosted Operations Notifications](self-hosted-operations-notifications.md),
[Small-Agency Pilot Operations](../runbooks/small-agency-pilot-operations.md),
and the [Self-Hosted Agency Reuse Master Plan](../master-plan-self-hosted-agency-reuse.md).

For broader navigation, see [Docs Home](../README.md). For detailed evidence boundaries, see [Compliance Evidence Checklist](../compliance-evidence-checklist.md).

## Truthfulness Rules For Tutorial Edits

- Every command must be runnable from the committed repo or clearly marked as deployment-specific.
- Every endpoint and environment variable must match the actual codebase.
- Public protobuf endpoints are anonymous; JSON debug, admin, and GTFS Studio routes are protected.
- `http://localhost:8080` is local-demo packaging only; production deployments need HTTPS/TLS and deployment-owned admin network boundaries.
- Use "supports" and "technical foundations for" when describing compliance readiness unless deployment and external-consumer evidence supports stronger wording.
