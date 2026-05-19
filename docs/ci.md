# Continuous Integration

CI is split into a fast path for pull requests and pushes, plus a manual
release-gate path for validator-heavy checks.

## Current Fast CI Baseline

The current repository baseline keeps `go test ./...` in the fast path:

```bash
go test ./...
```

It has passed across the repository in local verification. The Go test
workflow is useful and should stay. Fast CI also runs the repository's
lightweight checks and claim/consumer tracker audits so contributors see the
same baseline locally and in GitHub Actions.

GitHub Actions uses `actions/setup-go@v5` with `go-version-file: go.mod`.
`go.mod` currently declares `go 1.23.2` and does not set a separate
`toolchain` directive.

## Fast CI

`.github/workflows/test.yml` runs on pull requests and pushes to `main` or
`stable`.

It runs:

```bash
go test ./...
make check
scripts/check-consumer-tracker.sh
make audit-final-claim-review
```

This path avoids validator installation, Docker, external services, and
release packaging so normal PR feedback stays focused and repeatable.
`make check` also runs no-network repository guardrails, internal link checks,
final-claim/product-acceptance audits, product language and static layout
audits, the reference-settings product UI smoke, and a stable-branch filter
simulation with branch-ref checks skipped.

## Manual Release Gates

`.github/workflows/release-gates.yml` is a manual `workflow_dispatch` workflow.

It starts with the fast baseline, then can install pinned validators and run:

```bash
make validate
make smoke
make external-connection-check
make adapter-conformance
make test-connector-examples
make gtfsrt-conformance
make audit-product-acceptance
RELEASE_PACKAGE_ALLOW_DIRTY=true make release-package
make audit-release-package
```

Validator-heavy checks stay out of Fast CI because they depend on pinned
validator tooling and can be slower or more environment-sensitive. They remain
available for release-candidate review.

`make smoke` belongs in this manual workflow, after validator installation,
because it starts with `scripts/check-validators.sh`. It should not become a
required pull-request check unless the GitHub Actions runner has stable pinned
validator setup and the extra runtime cost is intentional.

The release-gates workflow is not required for every pull request. Run it
manually when preparing or reviewing a release candidate, or when a change
touches validator, connector, GTFS-RT conformance, product-acceptance, or
release-package behavior.

## Workflow Syntax Checks

When `actionlint` is available locally, use:

```bash
actionlint .github/workflows/*.yml
```

Without `actionlint`, the minimum local syntax pass is to parse the workflow
YAML files and run the fast baseline. GitHub Actions remains the source of
truth for remote workflow execution.

## Local Full Checks

Before a release-candidate gate, run the local full set from a clean checkout
when the environment supports validators:

```bash
git diff --check
go test ./...
make check
make test
make validate
make smoke
make audit-product-acceptance
make audit-final-claim-review
make external-connection-check
make adapter-conformance
make gtfsrt-conformance
scripts/check-consumer-tracker.sh
```

These checks do not prove production readiness, CAL-ITP/Caltrans compliance,
consumer acceptance, hosted service availability, vendor compatibility,
SLA/uptime, production AVL reliability, or ETA quality.
