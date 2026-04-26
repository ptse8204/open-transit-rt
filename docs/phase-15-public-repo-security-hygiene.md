# Phase 15 — Targeted Public Repo Hygiene And Evidence Redaction Review

## Status

Complete for the targeted Phase 15 scope.

## Purpose

Phase 15 prepares the repository for broader public attention by reviewing the
delta since the earlier Phase 14 security cleanup checkpoint, then tightening
evidence redaction, archive inventory, local-artifact ignores, and security
reporting guidance.

The chosen review baseline is commit `839efd6` (`Phase 14 -- Checkpoint 4 --
Security Cleanup`). This phase is delta-focused and does not claim a complete
historical security audit.

## Scope

1. Files added or changed since `839efd6`, using `git diff --name-only
   839efd6..HEAD`.
2. Tracked high-risk filename patterns from `git ls-files`, including archives,
   `.env`-like files, keys/certs, logs, local OS files, and generated evidence
   artifacts.
3. Phase 12 hosted OCI evidence, Phase 13 consumer evidence docs, Phase 14 public
   docs/assets, `.gitignore`, local artifacts, and generated/archive files.
4. Secret scanner attempts plus manual high-risk searches.
5. Concrete inventory of committed archives under `docs/evidence/captured/**`.

## Public Evidence Safety Rule

Public evidence may include public URLs, validation status, TLS metadata,
checksums, public HTTP status and headers, and redacted operational summaries.

Public evidence must not include raw credentials, bearer tokens, admin URLs with
secrets, private SSH paths, unredacted IP logs, private keys, database passwords,
or internal hostnames unless explicitly justified as public-safe.

## Completed Work

- Added `SECURITY.md` with private reporting guidance and evidence/secret
  handling rules.
- Added `docs/evidence/redaction-policy.md`.
- Added `docs/evidence/archive-inventory.md` with the two committed evidence
  ZIP archives and their contents.
- Expanded `.gitignore` for local environment files, key material, local logs,
  OS files, and root-level archives.
- Removed ignored local `.DS_Store` files from the working tree.
- Removed ignored local secret files found under `.cache/` from the working
  tree; they were not tracked and did not appear in git history for those paths.
- Redacted unnecessary OCI operator detail from committed evidence:
  - raw public client IP and remote ports in the Caddy route-map artifact;
  - OCI instance hostname in three operator-supplied artifacts.
- Refreshed the OCI pilot `SHA256SUMS.txt` entries for the modified artifacts.

## Required Rotation Or Revocation

The Docker gitleaks scan found real secrets in ignored local `.cache/` files:

- `.cache/duckdns-pilot/admin-token`
- `.cache/duckdns-pilot/device-rebind.json`
- `.cache/duckdns-pilot/env`
- `.cache/duckdns-pilot/caddy-data/.../default.key`
- `.cache/duckdns-pilot/caddy-data/.../open-transit-pilot.duckdns.org.key`
- `.cache/oci-admin-token`

Before further pilot use, the operator should rotate or revoke the affected
admin tokens, device token, admin JWT secret, CSRF secret, device token pepper,
ACME account key, and TLS certificate private key. The Phase 15 targeted check
found no git-tracked copy or git history record for those `.cache/` paths, so
destructive git history rewriting is not indicated from this finding.

## Acceptance Notes

- Public evidence artifacts were reviewed in the targeted Phase 15 scope.
- Committed archives are inventoried and retained only because they contain
  expected public/demo GTFS schedule files.
- Accidental local `.DS_Store` files and ignored local secret files were removed
  from the working tree.
- No backend runtime behavior, APIs, database schema, public feed URLs, consumer
  statuses, or evidence claims were changed.
- No compliance, consumer acceptance, production readiness, or vendor-equivalence
  claim is added by this phase.

## Checks

Pre-edit checks:

```bash
make validate
make test
git diff --check
```

Post-edit checks:

```bash
make validate
make test
git diff --check
make smoke
make demo-agency-flow
EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence
```

Scanner and manual review:

```bash
command -v gitleaks
command -v trufflehog
docker run --rm -v "$PWD:/repo:ro" zricethezav/gitleaks:latest dir /repo --redact --verbose --no-banner
git diff --name-only 839efd6..HEAD
git ls-files
find docs/evidence/captured -type f \( -name '*.zip' -o -name '*.tar' -o -name '*.tgz' -o -name '*.tar.gz' -o -name '*.gz' -o -name '*.7z' -o -name '*.rar' \) -print
```

The PATH scanner binaries were unavailable. Docker was available. The first
Docker gitleaks directory scan found ignored local `.cache/` secrets; after
removing those local files, the Docker gitleaks directory scan reported no
leaks. Manual high-risk searches over tracked and non-cache working-tree files
did not find committed private keys, cloud tokens, GitHub tokens, Slack tokens,
OpenAI-style API keys, or literal Bearer credentials.
