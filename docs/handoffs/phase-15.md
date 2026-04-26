# Phase 15 Handoff — Targeted Public Repo Hygiene And Evidence Redaction Review

## Phase

Phase 15 — Targeted Public Repo Hygiene And Evidence Redaction Review.

## Status

Complete for the targeted Phase 15 scope.

The review baseline was `839efd6` (`Phase 14 -- Checkpoint 4 -- Security
Cleanup`). The commit existed and was reachable, so no alternate baseline was
needed. This handoff records a targeted delta-focused review, not a complete
historical repository security audit.

## What Was Implemented

- Reviewed files changed since `839efd6` with `git diff --name-only
  839efd6..HEAD`; the delta was Phase 14 README/docs/roadmap/prompt material.
- Reviewed tracked high-risk filename patterns with `git ls-files`, including
  archives, `.env`-like files, keys/certs, logs, local OS files, generated
  evidence artifacts, checksum files, curl transcripts, and public protobuf
  evidence artifacts.
- Added `SECURITY.md` with a private reporting path and evidence/secret handling
  rules.
- Added `docs/evidence/redaction-policy.md` with the public evidence safety rule,
  required redactions, checksum refresh guidance, archive-inventory rule, and
  secret response workflow.
- Added `docs/evidence/archive-inventory.md` with every committed archive under
  `docs/evidence/captured/**`, including contents and keep decisions.
- Updated `docs/evidence/README.md` to link the redaction policy and archive
  inventory.
- Expanded `.gitignore` for local env files, secret/key material, local runtime
  logs, OS files, and root-level archives.
- Removed ignored local `.DS_Store` files from the working tree:
  - `./.DS_Store`
  - `./docs/.DS_Store`
- Removed ignored local `.cache` files containing real secrets from the working
  tree:
  - `.cache/duckdns-pilot/admin-token`
  - `.cache/duckdns-pilot/device-rebind.json`
  - `.cache/duckdns-pilot/env`
  - `.cache/duckdns-pilot/caddy-data/caddy/acme/acme-v02.api.letsencrypt.org-directory/users/default/default.key`
  - `.cache/duckdns-pilot/caddy-data/caddy/certificates/acme-v02.api.letsencrypt.org-directory/open-transit-pilot.duckdns.org/open-transit-pilot.duckdns.org.key`
  - `.cache/oci-admin-token`
- Redacted unnecessary operational detail from OCI evidence artifacts:
  - raw public client IP, remote ports, request IDs, process ID, and instance
    hostname in `artifacts/operator-supplied/reverse-proxy-caddy-route-map-redacted.txt`;
  - instance hostname in
    `artifacts/operator-supplied/scorecard-job-definition-and-history.txt`;
  - instance hostname in
    `artifacts/operator-supplied/monitoring-feed-monitor-history.txt`.
- Refreshed the corresponding entries in
  `docs/evidence/captured/oci-pilot/2026-04-24/SHA256SUMS.txt`.
- Updated `docs/current-status.md`, `docs/handoffs/latest.md`, and this Phase 15
  handoff.

## What Was Designed But Intentionally Not Implemented Yet

- No destructive git history rewriting was performed.
- No full historical repo-wide security scrub was repeated, because the user
  directed a delta-focused review unless a new risk required broader history
  cleanup.
- No backend runtime behavior, APIs, database schema, public feed URLs, consumer
  statuses, or evidence claims were changed.
- No consumer-submission claims were advanced.

## Schema And Interface Changes

None.

## Dependency Changes

No repo dependency files changed.

Operationally, Docker was used to run the gitleaks container because `gitleaks`
and `trufflehog` were not available on PATH.

## Migrations Added

None.

## Tests Added And Results

No code tests were added because Phase 15 changed docs, evidence text, ignore
rules, and local artifacts only.

Pre-edit checks run before repository edits:

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.

Post-edit checks run after repository edits:

- `make validate` — passed.
- `make test` — passed.
- `git diff --check` — passed.
- `make smoke` — passed.
- `make demo-agency-flow` — passed.
- `EVIDENCE_PACKET_DIR=docs/evidence/captured/oci-pilot/2026-04-24 make audit-hosted-evidence` — passed.

## Checks Run And Blocked Checks

Baseline and delta:

- `git rev-parse --verify 839efd6^{commit}` — passed; resolved to
  `839efd6bf4d0396388e37e86931ca11f18cfd640`.
- `git diff --name-only 839efd6..HEAD` — reviewed Phase 14 docs/README/roadmap
  delta.
- `git ls-files` — reviewed tracked high-risk patterns for archives, env-like
  files, keys/certs, logs, local OS files, and generated evidence artifacts.

Scanner attempts:

- `command -v gitleaks` — unavailable on PATH.
- `command -v trufflehog` — unavailable on PATH.
- `docker info --format '{{.ServerVersion}}'` — passed; Docker server version
  `29.4.0`.
- `docker run --rm -v "$PWD:/repo:ro" zricethezav/gitleaks:latest dir /repo --redact --verbose --no-banner` — first run found eight leaks in ignored local `.cache` files.
- After removing the ignored local secret files, rerunning the same Docker
  gitleaks command reported no leaks.

Manual high-risk searches:

- Searched tracked and non-cache working-tree files for private-key blocks,
  AWS-style keys, GitHub tokens, Slack tokens, OpenAI-style API keys, Google API
  keys, literal Bearer auth headers, DuckDNS/admin/JWT/CSRF/device/database
  secret assignments, and password assignments.
- Reviewed Phase 12 OCI evidence, Phase 13 consumer evidence docs, Phase 14
  public docs/assets, `.gitignore`, local artifacts, generated evidence
  artifacts, and root-level archive patterns.
- Searched OCI evidence for raw public IPs, Caddy request IP fields, and internal
  instance hostnames after redaction; no matches remained in the OCI packet for
  the redacted values.

Archive inventory:

- `find docs/evidence/captured -type f \( -name '*.zip' -o -name '*.tar' -o
  -name '*.tgz' -o -name '*.tar.gz' -o -name '*.gz' -o -name '*.7z' -o -name
  '*.rar' \) -print` found two committed archives.
- `unzip -l` verified both archives contain expected GTFS text files only.
- `sha256sum` recorded both archive hashes in
  `docs/evidence/archive-inventory.md`.

Git history check for discovered local secrets:

- `git ls-files -- <secret paths>` — no output.
- `git log --all --oneline -- <secret paths>` — no output.

## Known Issues

- Real local secrets were found in ignored `.cache` files. They were removed from
  the working tree, but the operator must rotate/revoke the affected admin
  tokens, device token, admin JWT secret, CSRF secret, device token pepper, ACME
  account key, and TLS private key before further pilot use.
- The targeted check found no tracked or historical git record for those `.cache`
  paths, so destructive git history cleanup is not indicated from the evidence
  available in this review.
- This phase does not claim a complete historical security audit.
- Local-demo evidence still contains loopback-only `127.0.0.1` logs. These are
  treated as public-safe local demo evidence under
  `docs/evidence/redaction-policy.md` because they do not expose real client IPs
  or credentials.

## Exact Next-Step Recommendation

Start Phase 16 — Agency Onboarding Product Packaging using
`docs/phase-16-agency-onboarding-product-packaging.md` and
`docs/roadmap-post-phase-14.md`.

Before any future public promotion or pilot reuse, complete the Phase 15
rotation/revocation actions for the ignored local `.cache` secrets discovered by
the scanner.
