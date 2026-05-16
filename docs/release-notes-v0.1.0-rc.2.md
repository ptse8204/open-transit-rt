# Open Transit RT `v0.1.0-rc.2`

Release date: 2026-05-16

Open Transit RT v0.1.0-rc.2 is a public release candidate for local/self-hosted evaluation.

It is not a stable release and does not prove production readiness, compliance, agency adoption, consumer acceptance, final-root readiness, hosted service availability, vendor compatibility, hardware certification, SLA/uptime, production AVL reliability, production-grade ETA quality, or real-world ETA accuracy.

## Summary

This release candidate packages the current post-rc1 source for local
evaluation, with archive-safe source packaging, improved install replay
behavior, stronger GTFS-RT and connector conformance checks, and browser-first
Operations Console guidance for small-agency evaluation.

## Changes Since `v0.1.0-rc.1`

- Public source archives now exclude protected retained-evidence and protected
  consumer-submission paths while keeping extracted archive checks usable.
- Release package generation and audit flows record checksums, provenance,
  SBOM metadata, claim flags, and protected-path archive scan expectations.
- The release-candidate gate documents the previously observed concurrent
  test-output collision risk and requires sequential validation for affected
  checks.
- GTFS-RT conformance coverage now includes offline protobuf checks for Vehicle
  Positions, Trip Updates, and Alerts.
- Connector and adapter maturity improved through synthetic-only manifests,
  offline conformance fixtures, and starter examples for telemetry,
  prediction-sidecar, validator, monitoring/export, and discovery workflows.
- Browser-first Operations Console and evaluator documentation improved for
  setup, GTFS review/import, feed health, validation health, realtime review,
  connector review, maintenance, help, access, audit, consumers, and evidence
  boundary review.
- Realtime review surfaces summarize Vehicle Positions publication decisions,
  Trip Updates withheld/fallback diagnostics, Alerts review context, and
  conservative assignment reasons without converting local signals into
  deployment evidence.
- Public docs, handoffs, and status artifacts continue to separate software
  capability from deployment, compliance, consumer, final-root, hosted-service,
  vendor, SLA, and ETA-quality claims.

## Install And Upgrade Notes

- Clean install from a Git tag remains the preferred evaluator path.
- Source archive installability is explicitly checked by extracting the
  generated source archive and running lightweight/bootstrap/validator/test
  commands from the extracted tree.
- Local app startup and five public feed fetches remain release-gate checks
  through the release-candidate diagnostic and install-confidence workflows.
- No production Docker image is published by this release candidate.

## Migrations

None in this release-candidate cut beyond the migrations already present in the
tagged source tree.

Before upgrading any existing local database, operators should back up the
database and run:

```bash
make migrate-status
make migrate-up
make migrate-status
```

## Operations Notes

- Start no-developer review from the private local Operations Console at
  `/admin/operations` after a technical helper starts the local app.
- Public feed paths remain the evaluator-facing outputs; admin, validation,
  evidence, debug, and Operations Console routes remain private.
- Prepared consumer packet records remain prepared-only. This release does not
  submit to, contact, or receive a status from any consumer or aggregator.

## Security And Privacy Notes

- Protected retained-evidence and protected consumer-submission paths are not
  included in public source archives.
- Release assets must pass protected-path scans before publication.
- The release uses synthetic fixtures and local diagnostics only; it does not
  include real vendor payloads, credentials, private agency data, private
  consumer correspondence, or external submission artifacts.

## Known Boundaries

- This is a release candidate, not a stable release.
- It does not prove production readiness or production AVL reliability.
- It does not prove CAL-ITP/Caltrans compliance.
- It does not prove agency adoption, approval, endorsement, or public launch.
- It does not prove consumer submission, review, acceptance, ingestion,
  listing, or display.
- It does not prove final-root readiness or hosted service availability.
- It does not prove vendor compatibility, hardware certification, paid support,
  SLA coverage, uptime, production-grade ETA quality, or real-world ETA
  accuracy.

## Release Gate

Publication is allowed only after the rc2 gate passes from a clean checkout,
including repository checks, validation, tests, smoke checks, package
generation, package audit, release-candidate diagnostics, extracted archive
replay, protected archive scans, claim-boundary audits, and consumer tracker
verification.

Expected gate commands include:

```bash
git diff --check
scripts/check-consumer-tracker.sh
python3 -m json.tool docs/evidence/consumer-submissions/status.json
make check
make validate
make test
make smoke
make test-release-package
make audit-product-acceptance
make audit-final-claim-review
make external-connection-check
make adapter-conformance
make test-connector-examples
make gtfsrt-conformance
docker compose -f deploy/docker-compose.yml config
```

Package generation and audit:

```bash
RELEASE_PACKAGE_VERSION=v0.1.0-rc.2 \
RELEASE_PACKAGE_OUTPUT_DIR=.cache/release-package/v0.1.0-rc.2 \
RELEASE_PACKAGE_FORCE=true \
RELEASE_PACKAGE_ALLOW_DIRTY=false \
scripts/release-package.sh

RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.2 \
scripts/audit-release-package.sh
```
