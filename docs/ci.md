# Continuous Integration

CI is split into a fast path for pull requests and pushes, plus a manual
release-gate path for validator-heavy checks.

## Phase 13 Evaluation

Phase 13 did not reproduce a broken `go test ./...` result locally:

```bash
go test ./...
```

passed across the repository. The Go test workflow is useful and should stay.
The old workflow was still worth repairing because it did not match the current
repo stage: it ran connector checks in the default path, duplicated consumer
tracker logic inline, and did not run `make check` or
`make audit-final-claim-review`.

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
