# Reference Deployment Doctor

The deployment doctor is a read-only diagnostic helper for the OCI/OCL-style
reference deployment path. It is private operator diagnostics, not evidence
that a deployment was run, accepted, compliant, agency-approved, or
production-ready.

Doctor output can become evidence only in a separate evidence phase that
reviews, redacts, dates, inventories, and intentionally retains specific
artifacts. Running the doctor does not create external evidence, final-root
evidence, consumer submission artifacts, CAL-ITP/Caltrans compliance proof,
agency adoption proof, hosted SaaS proof, production-readiness proof,
vendor-compatibility proof, or production-grade ETA quality proof.

## Run

From the repo checkout:

```sh
make deployment-doctor
```

By default this exits `0` even when the local checkout has no deployment env or
running app. Missing deployment values, stopped services, unavailable feeds, or
missing backup/restore settings are recorded as blockers, skipped checks, or
unavailable checks in the generated report.

Use strict mode only when a deployment operator wants blockers to fail the
command:

```sh
STRICT_DOCTOR=true make deployment-doctor
```

Output defaults to:

```text
.cache/deployment-doctor/<timestamp>/
```

The helper rejects custom output directories outside repo `.cache/` unless
`ALLOW_UNIGNORED_OUTPUT_DIR=true`, rejects symlink output directories, creates
the output directory mode `700`, and fails on a non-empty output directory
unless `FORCE=true`.

## Supported Variables

The doctor inspects already-exported environment variables only. It does not
source private env files. It derives the reference key list from
`docs/deployment/oci-reference-env.example` and records only whether each key
is present, missing, placeholder-like, too short, optional empty, or skipped.
It never writes the values.

Common variables:

- `PUBLIC_BASE_URL`
- `ADMIN_BASE_URL`
- `ADMIN_TOKEN`
- `DATABASE_URL`
- `MIGRATIONS_DIR`
- `BACKUP_DIR`
- `RESTORE_DATABASE_URL`
- `RESTORE_BACKUP_FILE`
- `STRICT_DOCTOR`
- `OUTPUT_DIR`
- `ALLOW_UNIGNORED_OUTPUT_DIR`
- `FORCE`

Timeout and size knobs:

- `CONNECT_TIMEOUT_SECONDS=5`
- `REQUEST_TIMEOUT_SECONDS=30`
- `MAX_FEED_BYTES=104857600`

Notification variables such as `NOTIFY_WEBHOOK_URL` and `NOTIFY_EMAIL_TO` are
optional. Empty values are allowed and are not blockers. If supplied, values
are never printed or written.

## What It Checks

The doctor checks:

- reference env key presence and placeholder status without values;
- generated-secret presence and minimum length for `ADMIN_JWT_SECRET`,
  `CSRF_SECRET`, and `DEVICE_TOKEN_PEPPER`;
- optional `ADMIN_TOKEN` status without making absence a blocker unless an
  authenticated admin check is requested;
- `PUBLIC_BASE_URL` and `ADMIN_BASE_URL` syntax and safe boundary rules;
- anonymous public fetch metadata for the five public feed paths;
- unauthenticated public-edge denial for private/admin routes;
- optional authenticated readiness through `ADMIN_BASE_URL` using only
  `Authorization: Bearer "$ADMIN_TOKEN"`;
- optional authenticated validator-health JSON through
  `/admin/operations/validation-health.json` using only `GET`, only when
  `ADMIN_TOKEN` and a safe `ADMIN_BASE_URL` are present;
- HTTPS and HTTP-to-HTTPS redirect posture for non-loopback HTTPS public roots;
- loopback service `/healthz` and `/readyz` status for ports `8081` through
  `8086`;
- read-only database reachability and migration status when `DATABASE_URL` is
  supplied;
- PostGIS availability when it can be probed safely;
- pinned validator tooling status;
- backup directory readiness without creating a backup;
- restore-drill readiness without running restore;
- git/release identity;
- consumer tracker shape, requiring the seven expected targets to remain
  `prepared`.

The private route boundary checks use `HEAD` first and fall back to `GET` only
when `HEAD` returns `405`. They never `POST` to admin routes. Public feed
checks fetch bodies only to temporary files, record status, size, checksum,
redirect count, effective URL, and content type, then delete the bodies.

The validator-health integration stores summary fields only: generated time,
agency ID, overall/tooling status, false claim flags, and bounded per-feed
status fields. It never runs validators, never POSTs the route, never stores
raw reports, and is private diagnostics only.

## Outputs

The output directory includes:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- supporting status summaries

`summary.json` and `manifest.json` are validated before the command exits.
The summary includes counts for `passed`, `blocker`, `warning`, `skipped`, and
`unavailable`, plus a deterministic `overall_status`:

- `blocker` if any blocker exists;
- `warning` if no blockers but warnings exist;
- `unavailable` if no blockers or warnings exist but required local checks are
  unavailable;
- `passed` only if required checks pass or are explicitly skipped by documented
  local/reference rules.

The report also records:

```text
external_evidence_created=false
final_root_evidence_created=false
consumer_statuses_changed=false
compliance_claimed=false
production_readiness_claimed=false
```

## Secret Safety

The doctor never prints or writes:

- admin tokens;
- Authorization headers;
- cookies;
- JWTs;
- CSRF secrets;
- device token peppers;
- database URLs;
- restore database URLs;
- webhook URLs;
- raw `.env` files;
- private keys or ACME material;
- raw database dumps;
- backup file contents.

It runs a final redaction scan over generated output and fails if it detects
secret-shaped content such as bearer credentials, JWT-looking values,
passworded Postgres URLs, cookie headers, private keys, raw secret assignment
lines, or webhook URLs with embedded token material.

## Boundaries

The deployment doctor is read-only. It does not run migrations, create backups,
restore databases, create evidence packets, create final-root proof, contact
consumers, change consumer statuses, or claim CAL-ITP/Caltrans compliance,
agency adoption, hosted SaaS availability, production readiness, vendor
compatibility, or production-grade ETA quality.
