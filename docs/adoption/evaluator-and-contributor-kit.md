# Evaluator And Contributor Kit

Use this kit when an agency evaluator, civic technologist, or open-source
contributor wants to try Open Transit RT and send public-safe feedback or a
first contribution.

This kit is not adoption proof. It does not claim agency approval, hosted
service availability, SLA/uptime, production readiness, CAL-ITP/Caltrans
compliance, consumer submission/review/acceptance/ingestion/listing/display,
vendor compatibility, hardware certification, final-root readiness, or
production-grade ETA quality.

## Pick A Path

| Path | Best for | Time box | Start here | Output |
| --- | --- | --- | --- | --- |
| Local browser evaluation | Non-developer agency review with a technical helper | 30-60 minutes | `make agency-app-up`, then `/admin/operations` | Public-safe notes about setup, GTFS import, feed health, telemetry, connectors, readiness, and maintenance clarity |
| Release-candidate install trial | Technical helper checking rc1 installability | 60-120 minutes | `git checkout v0.1.0-rc.1`, `make check`, `make agency-app-up` | Exact commands, environment blockers, and public-safe install notes |
| Synthetic connector review | Integrator or vendor-adjacent contributor using fake data | 30-90 minutes | `docs/connectors/vehicle-avl-starter-kits.md` | Synthetic fixture, manifest, docs, or conformance feedback |
| First contribution | New contributor | 30-120 minutes | `docs/contributor-first-issues.md` | Small docs fix, synthetic fixture, focused test, connector manifest edit, or troubleshooting improvement |
| Feedback-only review | Agency or evaluator without code changes | 15-45 minutes | `docs/agency-feedback-template.md` | Public-safe feedback with private values removed |

## No-Claim Trial Rules

- Use public GTFS, committed demo data, or synthetic fixtures unless the agency
  has explicitly approved another data source.
- Keep tokens, database URLs, private endpoints, raw logs, raw telemetry, and
  private screenshots out of issues and docs.
- Treat local feed checks, screenshots, and demo outputs as evaluation notes,
  not retained evidence.
- Keep consumer packet statuses unchanged unless a separate retained-evidence
  phase authorizes a specific update.
- Record exact blockers instead of converting missing data into success.

## Local Evaluator Flow

1. Ask a technical helper to start the app:

   ```bash
   make agency-app-up
   ```

2. Open the private Operations Console:

   ```text
   http://localhost:8080/admin/operations
   ```

3. Review these paths in order:

   | Step | Page | What to note |
   | --- | --- | --- |
   | Start | Agency Operations Cockpit / Start Here | Whether first-run tasks and missing setup signals are understandable |
   | Schedule | GTFS Import or GTFS Workbench | Whether importing or reviewing GTFS feels clear |
   | Feeds | Feed Health and Validation Center | Whether all five feed rows explain current status and next action |
   | Realtime | Realtime Center and Prediction Lab | Whether Vehicle Positions, Trip Updates, withheld predictions, and confidence language are clear |
   | Connectors | Connector Hub and Workbench | Whether CSV, polling, webhook-sidecar, and synthetic paths are discoverable |
   | Maintain | Maintenance Center | Whether backup/restore, off-host validation, small-host readiness, and support guidance are understandable |
   | Learn | Help and Consumers | Whether claim and evidence boundaries are visible |

4. Fill in `docs/agency-feedback-template.md` using only public-safe details.

## Demo Paths

- Browser-first local demo: `docs/tutorials/no-cli-agency-first-run.md`
- Agency demo flow: `docs/tutorials/agency-demo-flow.md`
- Staff training demo kit: `docs/tutorials/staff-training-demo-kit.md`
- Telemetry simulator and device trial:
  `docs/tutorials/telemetry-simulator-and-device-trial.md`
- Small agency maintenance guide:
  `docs/tutorials/small-agency-maintenance-guide.md`
- Public rc1 install confidence:
  `docs/public-install-confidence-v0.1.0-rc.1.md`

## Feedback Templates

Use `docs/agency-feedback-template.md` for evaluator notes. For bug reports or
feature requests, use the GitHub issue templates and keep the report
reproducible with public-safe commands, synthetic fixtures, redacted output,
and exact expected/actual behavior.

For release-candidate feedback, support-bundle sharing, and maintainer triage
lanes, use
`docs/support/community-support-and-issue-triage-kit.md`.

Include:

- command or browser route reviewed;
- expected result;
- actual result;
- public-safe environment detail, such as operating system and Docker/Compose
  setup;
- redacted error text or fixture link;
- whether the issue affects setup, GTFS import, feeds, telemetry, connectors,
  validation, readiness, maintenance, or docs.

Do not include:

- secrets, tokens, DB URLs, private hostnames, raw logs, raw telemetry, raw
  GTFS, private screenshots, private ticket links, or unredacted operator
  artifacts;
- claims that the agency adopted, approved, certified, launched, submitted to
  a consumer, reached compliance, or achieved production readiness.

## First Contribution Map

| Contributor interest | Good first issue | Avoid at first |
| --- | --- | --- |
| Docs | Fix one command, missing prerequisite, stale link, or unclear boundary | Broad README rewrites that change public claims |
| Fixtures | Add one small synthetic fixture with deterministic IDs | Real agency, vendor, device, or portal data |
| Tests | Add one focused test for existing documented behavior | Public feed contract changes or schema migrations |
| Connectors | Improve one manifest, starter-kit row, or synthetic conformance case | Live vendor sends, hardware claims, or credential handling |
| Operations | Clarify one support-bundle, maintenance, or deployment-doctor step | Live backup/restore, service control, or evidence collection |
| Triage | Reproduce with public-safe commands and link related docs | Requesting private artifacts in public issues |

Before opening a PR, run:

```bash
git diff --check
make check
```

If code changes, also run:

```bash
make test
```

If connector examples change, also run:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

If public-claim, readiness, evidence, release, or support wording changes,
also run:

```bash
make audit-product-acceptance
make audit-final-claim-review
```

## What Maintainers Need

For feedback:

- exact public-safe reproduction steps;
- what was confusing or blocked;
- whether the issue is docs, UI, install, connector, validation, realtime, or
  maintenance;
- whether a private technical helper is needed.

For PRs:

- concise summary;
- checks run;
- docs updated when behavior or setup changed;
- statement that protected evidence paths and consumer statuses were not
  touched;
- statement that no secrets or private operator artifacts were added.

## Boundaries To Preserve

Evaluation, feedback, stars, discussions, issues, and PRs are useful project
signals. They are not proof of adoption, compliance, production readiness,
consumer acceptance, hosted availability, SLA coverage, or vendor/hardware
compatibility.
