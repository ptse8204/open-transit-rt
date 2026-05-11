# Release Notes Template

Use this template for tagged Open Transit RT releases. Replace every placeholder
with release-specific facts. If a section has no changes, write `None`.

## Open Transit RT `<tag>`

Release date: `<YYYY-MM-DD>`

Source:

- Git tag: `<tag>`
- Commit SHA: `<sha>`
- Dirty/clean state: `<output of git describe --tags --always --dirty>`
- Release notes link: `<link>`
- Artifact checksums: `<links or None>`
- Release package: `<.cache path, release asset link, or None>`
- SBOM/provenance: `<status or None>`

## Summary

`<One short paragraph describing the release scope.>`

## User-Facing Changes

- `<change>`

## Install And Upgrade Notes

- Clean install from source tag: `<supported / notes>`
- Local app verification: `<commands or link>`
- Release-candidate diagnostic: `<.cache/release-candidate-check path and status, or None>`
- Local release package: `<package path and audit result or None>`
- Local Docker image build: `<image tag or None>`
- Published production Docker image: `None; deferred unless a future release says otherwise.`

## Migrations

- `<None, or list migration files and required order>`

Before upgrade, operators must back up the database and run:

```bash
make migrate-status
make migrate-up
make migrate-status
```

## Operations Changes

- `<None, or operational changes>`

## Security Notes

- `<None, or security-relevant changes>`

Do not include secrets, private URLs, raw logs with credentials, private portal
artifacts, or unredacted operator evidence.

## Dependency Changes

- `<None, or dependency/tooling changes>`

## Evidence Or Claim Changes

- `<None, or evidence-backed changes>`

Do not claim CAL-ITP/Caltrans compliance, consumer submission/review/acceptance,
agency endorsement, hosted SaaS availability, paid support, SLA coverage,
marketplace/vendor equivalence, universal production readiness, or
production-grade ETA quality unless retained evidence supports that exact claim.

## Known Limitations

- `<limitation>`

## Checks

Record release-candidate and final check results:

```bash
make check
make release-candidate-check
make validate
make test
make test-release-package
make realtime-quality
make smoke
make audit-final-claim-review
docker compose -f deploy/docker-compose.yml config
git diff --check
```

Blocked checks:

- `<None, or exact command and reason>`

Release-candidate summary:

- Output directory: `<.cache/release-candidate-check/... or None>`
- Overall status: `<passed / needs_review / not_checked / blocker>`
- Package audit status: `<passed / needs_review / not_checked / blocker / None>`
- Claim-boundary result: `<no retained evidence, consumer status change, image push, hosted-service claim, production-readiness claim, or compliance claim / notes>`
