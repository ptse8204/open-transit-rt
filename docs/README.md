# Documentation Home

Use this page when the wiki does not go deep enough. Public readers should
start with the task guides first; maintainers can use the lower sections for
release, evidence, claim-boundary, and project-history references.

Start with the public product explainer site when you want the shortest
overview: [https://ptse8204.github.io/open-transit-rt/](https://ptse8204.github.io/open-transit-rt/).

## Browser-First Product Path

Use the same order in the README, wiki, docs, GitHub Pages, and private UI:

1. Start in the browser.
2. Open **Agency Operations Cockpit / Start Here**.
3. Review setup.
4. Import or review GTFS.
5. Check the five configured feed URLs.
6. Review feed health, readiness, validation, telemetry, connectors, and
   maintenance.
7. Understand what remains before deployment or stronger claims.

![Illustrative contribution paths: report a bug, improve docs, suggest a feature, submit code, and help with evidence runbooks.](assets/how-to-contribute-paths.png)

## Public User Docs

- [Product Explainer Site](https://ptse8204.github.io/open-transit-rt/)
- [Small Agency Quick Start](../wiki/small-agency-quick-start.md)
- [Browser-First Setup](../wiki/browser-first-setup.md)
- [No Command Line First Run](tutorials/no-cli-agency-first-run.md)
- [Small Agency Maintenance Guide](tutorials/small-agency-maintenance-guide.md)
- [Operations Console Tour](../wiki/operations-console-tour.md)
- [Can My Agency Use This?](../wiki/can-my-agency-use-this.md)
- [Agency Evaluation Checklist](../wiki/agency-adoption-checklist.md)
- [CAL-ITP Readiness Plain English](../wiki/calitp-readiness-plain-english.md)

These pages answer what you can do in the UI, what still needs a technical
helper, and what the local evaluation does not prove.

## Operator Docs

- [Review And Recommendations](roadmap-status.md#review-and-recommendations)
- [Agency First Run](tutorials/agency-first-run.md)
- [No Command Line First Run](tutorials/no-cli-agency-first-run.md)
- [Small Agency Maintenance Guide](tutorials/small-agency-maintenance-guide.md)
- [Small-Agency Acceptance Script](tutorials/small-agency-acceptance-script.md)
- [Reusable Agency Onboarding](tutorials/reusable-agency-onboarding.md)
- [Self-Hosted Operator Trial](tutorials/self-hosted-operator-trial.md)
- [Agency Launchpad](tutorials/agency-launchpad.md)
- [Operator Smoke And Support Bundle](tutorials/operator-smoke-and-support-bundle.md)
- [GTFS Validation Triage](tutorials/gtfs-validation-triage.md)
- [Prediction And ETA Lab](tutorials/prediction-eta-lab.md)
- [CAL-ITP Readiness Checklist](tutorials/calitp-readiness-checklist.md)
- [Deploy With Docker Compose](tutorials/deploy-with-docker-compose.md)
- [Reference Deployment Doctor](deployment/reference-deployment-doctor.md)
- [OCI Reference Check](deployment/oci-reference-check.md)
- [Off-Host Public Feed Validation](deployment/off-host-validation.md)

## Integrator Docs

- [Integration Adapter Kit](integration-adapter-kit.md)
- [Connector Plugin Contract](connectors/plugin-contract.md)
- [Contributing Connectors](connectors/contributing-connectors.md)
- [Connector Cookbook](../wiki/connector-cookbook.md)
- [Device And AVL Integration](tutorials/device-avl-integration.md)
- [Device Token Lifecycle](tutorials/device-token-lifecycle.md)
- [External Adapter Conformance](tutorials/external-adapter-conformance.md)
- [External Connection Readiness](external-connection-readiness.md)
- [Examples Index](../examples/README.md)
- [Test Fixture Index](../testdata/README.md)
- [Dependencies](dependencies.md)

External-connection maturity means synthetic/local review of telemetry to
`POST /v1/telemetry`, predictor adapters behind `internal/prediction.Adapter`,
validator tooling, monitoring/export surfaces, feed-consumer URL/metadata
expectations, and redaction checks. Real vendor proof remains optional
evidence only when authorized.

## Release-Candidate Docs

- [Release-Candidate Readiness](release-candidate-readiness.md)
- [Release Process](release-process.md)
- [Release Checklist](release-checklist.md)
- [Release Notes Template](release-notes-template.md)
- [Upgrade And Rollback](upgrade-and-rollback.md)
- [Demo And Documentation Site Plan](demo-docs-site-plan.md)
- [Changelog](../CHANGELOG.md)

Release-candidate checks are local maintainer diagnostics. They do not tag,
publish, push images, create retained evidence, or prove production readiness.
The next recommended release milestone is `v0.1.0-rc.1` before any full
`v0.1.0` release.

## Contributor Docs

- [Contributing](../CONTRIBUTING.md)
- [Contributor First Issues](contributor-first-issues.md)
- [Contributing Connectors](connectors/contributing-connectors.md)
- [Public Docs And Site Freeze Checklist](public-docs-site-freeze-checklist.md)
- [Test Fixture Index](../testdata/README.md)
- [Examples Index](../examples/README.md)

Contributor onboarding should favor small docs fixes, synthetic fixtures,
focused tests, connector examples, redaction cleanup, and task-based
navigation. First issues should not touch migrations, public feed contracts,
protected evidence paths, consumer tracker status, release publication, or
unsupported public claims.

## Visual Documentation

- [Product Screenshots](assets/product-screenshots/README.md)
- [Product Diagrams](assets/product-diagrams/README.md)
- [Documentation Assets](assets/README.md)
- [Public Docs And Site Freeze Checklist](public-docs-site-freeze-checklist.md)

Product screenshots and diagrams are local/demo documentation aids only. They
must not be stored under `docs/evidence`, called evidence, or used as proof of
production, compliance, adoption, consumer acceptance, final-root readiness,
vendor compatibility, or ETA quality.

## Evidence And Claim-Boundary Docs

- [Readiness And Evidence](../wiki/readiness-and-evidence.md)
- [Compliance Evidence Checklist](compliance-evidence-checklist.md)
- [California Readiness Summary](california-readiness-summary.md)
- [CAL-ITP / Caltrans Requirements](requirements-calitp-compliance.md)
- [Phase 60 Final Claim Review](phase-60-final-claim-review-and-public-closeout.md)
- [Consumer Submission Evidence](consumer-submission-evidence.md)
- [Consumer Submission Tracker](evidence/consumer-submissions/README.md)
- [Consumer Submission Workflow](evidence/consumer-submissions/submission-workflow.md)
- [Captured Evidence Index](evidence/captured/README.md)

Formal agency approval, final feed-root evidence, and consumer acceptance are
not required for local evaluation or open-source contribution. They are future
authorization-gated evidence tracks only when a maintainer explicitly supplies
the required scope, retention, redaction, and stop conditions.

## Architecture And Requirements

- [Architecture](architecture.md)
- [Admin Command Model](admin-command-model.md)
- [Dependencies](dependencies.md)
- [Decisions](decisions.md)
- [Known Gaps](repo-gaps.md)
- [Requirements 2A-2F](requirements-2a-2f.md)
- [Trip Updates Requirements](requirements-trip-updates.md)
- [CAL-ITP / Caltrans Requirements](requirements-calitp-compliance.md)

## Maintainer Docs And History

Detailed phase history and notes are retained as project history. They should
not be the first path for a small-agency evaluator.

- [Current Status](current-status.md)
- [Latest Handoff](handoffs/latest.md)
- [Phase 61+ Product Roadmap](roadmaps/agency-first-connector-platform/README.md)
- [Adoption Productization Roadmap](roadmaps/agency-first-connector-platform/adoption-productization-roadmap.md)
- [Consumer-Grade Control Plane Proposed/Authorized Roadmap Pack](roadmaps/consumer-grade-control-plane/README.md)
- [Historical Post-60 Product Roadmap](post-60-product-roadmap.md)
- [Roadmap Status](roadmap-status.md)
- [Backlog](backlog.md)
- [Open Questions](open-questions.md)
- [Maintainer Handoffs](handoffs/)
- [Documentation Assets](assets/README.md)

Public-facing pages should stay task-based and reader-friendly. Detailed
evidence records, operational notes, and implementation history belong in the
deeper docs instead of the top-level README.
