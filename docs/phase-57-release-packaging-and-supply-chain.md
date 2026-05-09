# Phase 57 -- Release Packaging And Supply Chain

## Status

Planning accepted. Execution may add local release package generation, release
package audit tooling, checksums, SBOM/provenance metadata, tests, Make targets,
and docs/status/handoff updates. It must not publish artifacts, push images,
contact external services, or claim hosted service availability.

## Goal

Make versioned local install artifacts reviewable through source archives,
checksums, SBOM/provenance metadata, and release-package audit checks while
preserving all evidence and claim boundaries.

## Scope

- Add a local release package generator that writes ignored `.cache` output by
  default.
- Generate a source archive from the current git commit.
- Generate SHA-256 checksums for generated local artifacts.
- Generate local provenance metadata from git identity, build inputs, command
  versions, and dirty/clean state.
- Generate a practical Go-module SBOM summary from local `go list -m -json all`.
- Record optional local Docker image metadata only when an operator supplies an
  image tag; do not build, push, or publish images by default.
- Add an audit helper that validates release package file set, JSON, checksums,
  false claim flags, safe paths, and consumer tracker preservation.
- Add focused script tests and Make targets.
- Update release, install/upgrade, dependency, decision, roadmap/status, and
  handoff docs.

## Non-Goals

- No hosted SaaS claim.
- No hosted service or production image publication.
- No registry push.
- No package repository upload.
- No GitHub release creation.
- No signing key, key management, Sigstore, or external attestation service.
- No production-readiness claim.
- No CAL-ITP/Caltrans compliance claim.
- No consumer acceptance, consumer ingestion, or consumer status change.
- No agency adoption claim.
- No marketplace approval or vendor compatibility claim.
- No SLA/uptime or paid support claim.
- No production-grade ETA claim.
- No retained evidence creation.
- No `docs/evidence` writes.

## Files Likely To Change

- `scripts/release-package.sh`
- `scripts/audit-release-package.sh`
- `scripts/test-release-package.sh`
- `Makefile`
- `docs/phase-57-release-packaging-and-supply-chain.md`
- `docs/release-process.md`
- `docs/release-checklist.md`
- `docs/release-notes-template.md`
- `docs/upgrade-and-rollback.md`
- `docs/dependencies.md`
- `docs/decisions.md`
- `docs/current-status.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-57.md`

Execution must not change:

- `docs/evidence/consumer-submissions/status.json`
- `docs/evidence/consumer-submissions/current/`
- `docs/evidence/consumer-submissions/artifacts/`
- `docs/evidence/consumer-submissions/packets/`
- `docs/evidence/captured/`

## Safety Boundaries

Release package output must default under ignored
`.cache/release-package/<version-or-timestamp>/`, mode `0700`, and reject
symlink, traversal, `docs/evidence`, evidence-like, proof-like, submission-like,
and outside-repo output paths.

The generator must not read `.env`, `.env.*`, private keys, raw logs, raw
telemetry payloads, database dumps, generated evidence packets, consumer
artifact directories, or private operator outputs into release artifacts.
Source archives must come from `git archive HEAD`, not recursive filesystem
copying.

Generated metadata may include public-safe local tool versions, git commit
identity, dirty/clean state, Go module names/versions/checksums, local image tag
or ID when explicitly provided, and SHA-256 checksums of generated local
artifacts.

Dirty worktrees must be recorded as not release-ready. Strict mode may fail
when the checkout is dirty, but test/validation paths may use an explicit
dirty-allowed test mode.

## Evidence And Claim Boundaries

All claim flags remain false:

- compliance;
- consumer acceptance;
- consumer ingestion;
- agency adoption;
- hosted SaaS;
- hosted service;
- production readiness;
- SLA/uptime;
- vendor compatibility;
- marketplace approval;
- production-grade ETA.

Release packages are local install/review artifacts only. They are not retained
evidence, consumer submission packets, compliance packets, hosted-service proof,
or production readiness proof.

Consumer tracker state remains unchanged; all seven targets remain
`prepared`.

## Implementation Details

Add `scripts/release-package.sh`:

- `--help` only for CLI arguments; configuration through environment variables.
- `RELEASE_PACKAGE_VERSION` optional, default from `git describe --tags --always`.
- `RELEASE_PACKAGE_OUTPUT_DIR` optional, default
  `.cache/release-package/<safe-version-or-timestamp>`.
- `RELEASE_PACKAGE_FORCE=true|false` for reusing non-empty output.
- `RELEASE_PACKAGE_ALLOW_DIRTY=true|false` for test/local diagnostic runs.
- `RELEASE_PACKAGE_STRICT=true|false` to fail on dirty checkout or unavailable
  SBOM source.
- `RELEASE_PACKAGE_IMAGE_TAG` optional; inspect local image metadata when
  supplied, otherwise record `not_configured`.
- Create exactly documented default files and directories:
  - `summary.json`
  - `summary.md`
  - `manifest.json`
  - `manifest.md`
  - `provenance.json`
  - `provenance.md`
  - `sbom.json`
  - `image.json`
  - `artifacts/open-transit-rt-<version>.source.tar.gz`
  - `checksums/SHA256SUMS.txt`
- Use `git archive --format=tar --prefix=... HEAD | gzip -n` for the source
  archive.
- Compute SHA-256 checksums with Python or `shasum -a 256`.
- Generate SBOM metadata by parsing `go list -m -json all` output with bounded
  subprocess output and graceful unavailable status unless strict mode requires
  failure.
- Include false claim flags in summary, manifest, provenance, and image
  metadata.

Add `scripts/audit-release-package.sh`:

- Validate exact default file set.
- Validate JSON files.
- Recompute checksums and compare `checksums/SHA256SUMS.txt`.
- Fail on true claim flags.
- Fail on unsafe/private string patterns.
- Fail if the package claims hosted service, production readiness, compliance,
  consumer acceptance, agency adoption, marketplace approval, vendor
  compatibility, SLA/uptime, or production-grade ETA.
- Validate the seven-target prepared-only consumer tracker without writing it.

Add `scripts/test-release-package.sh`:

- Exercise help output.
- Generate a test package under `.cache` with dirty-allowed mode.
- Audit the generated package.
- Test output path rejection for outside-repo, symlink, and evidence-like
  paths.
- Test audit failures for checksum drift, true claim flags, extra files, and
  unsafe wording.

Update `Makefile`:

- `release-package`
- `audit-release-package`
- `test-release-package`
- validation scaffolding for scripts/docs.

## Tests

- Shell syntax tests for all new scripts.
- Local script integration tests.
- Audit mutation tests for checksum drift, extra files, true claim flags, and
  unsafe wording.
- Existing release docs must continue to pass `make validate`.
- Consumer tracker preservation checks must pass.

## Performance And Scale Tests

- Source archive must use `git archive HEAD` instead of walking the full working
  tree.
- SBOM generation must cap subprocess output and record `unavailable` rather
  than hanging indefinitely.
- Audit checksum verification is linear in generated package artifact size and
  runs only on `.cache` local output.

Benchmark results, if any, are local engineering diagnostics only and are not
SLA, capacity, production-readiness, compliance, or supply-chain proof.

## Docs, Status, And Handoff Updates

Close Phase 57 by updating:

- this phase document;
- `docs/handoffs/phase-57.md`;
- `docs/handoffs/latest.md`;
- `docs/current-status.md`;
- `docs/backlog.md`;
- `docs/open-questions.md`;
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`;
- release/upgrade docs;
- `docs/dependencies.md`;
- `docs/decisions.md`.

The handoff must explicitly state that Phase 57 creates local packaging
scaffolding only and no hosted service, production image publication,
production-readiness, compliance, consumer acceptance, agency adoption,
marketplace approval, vendor compatibility, SLA/uptime, retained evidence, or
production-grade ETA claim.

## Required Verification Commands

Run and report:

```bash
sh -n scripts/release-package.sh scripts/audit-release-package.sh scripts/test-release-package.sh
./scripts/test-release-package.sh
make release-package
RELEASE_PACKAGE_DIR=<generated-dir> make audit-release-package
make validate
make test
make smoke
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path

expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]

data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
git diff --exit-code -- docs/evidence/consumer-submissions/current docs/evidence/consumer-submissions/artifacts docs/evidence/consumer-submissions/packets docs/evidence/captured
find docs/evidence/consumer-submissions/artifacts -mindepth 2 -maxdepth 2 -type f ! -name README.md -print
git status --short -- docs/evidence/consumer-submissions docs/evidence/captured
docker compose -f deploy/docker-compose.yml config
```

Run `INTEGRATION_TESTS=1 make test-integration` if the local database is
available and record any environment blocker truthfully.

The `find` command must print no files for the current Phase 57 state.
