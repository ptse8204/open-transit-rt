# Operator Smoke And Support Bundle

This tutorial explains the Phase 41 operator smoke and support bundle helpers.
They are local/reference diagnostics, not evidence packets. They do not create
external evidence, final-root evidence, consumer submission artifacts, agency
approval evidence, CAL-ITP/Caltrans compliance proof, vendor compatibility
proof, production-readiness proof, or production-grade ETA quality evidence.

## When To Use Each Command

Use operator smoke when a local/reference app is running and you want a strict
answer to whether the five public feed paths respond, the public edge denies
admin readiness access, optional authenticated readiness works through a safe
admin URL, validator tooling is installed or clearly blocked, and the synthetic
AVL dry-run fixture still works.

Use support bundle when you need a redaction-safe diagnostic packet for a
maintainer. It can run even when the app is not available; unavailable public
feeds, admin checks, database status, or validators are recorded as skipped,
unavailable, or blocker instead of failing the whole bundle.

## Operator Smoke

Run against the local app package or a reference deployment:

```bash
make operator-smoke SKIP_VALIDATORS=true
```

For authenticated readiness and validator API checks, provide a safe admin URL
and a bearer token:

```bash
PUBLIC_BASE_URL=http://localhost:8080 \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-private-token \
make operator-smoke
```

The script reads environment variables directly. The Make target does not echo
secret values.

Supported variables include:

- `PUBLIC_BASE_URL`
- `ADMIN_BASE_URL`
- `ADMIN_TOKEN`
- `SKIP_VALIDATORS`
- `STRICT_VALIDATORS`
- `OUTPUT_DIR`
- `CONNECT_TIMEOUT_SECONDS`
- `REQUEST_TIMEOUT_SECONDS`
- `MAX_FEED_BYTES`
- `FORCE`
- `ALLOW_UNIGNORED_OUTPUT_DIR`

Output defaults to:

```text
.cache/operator-smoke/<timestamp>/
```

The helper rejects custom output directories outside repo `.cache/` unless
`ALLOW_UNIGNORED_OUTPUT_DIR=true`. It rejects pre-existing symlink output
directories and non-empty output directories unless `FORCE=true`.

## Support Bundle

Run:

```bash
make support-bundle
```

Output defaults to:

```text
.cache/support-bundles/<timestamp>/
```

To include authenticated readiness status, opt in explicitly:

```bash
PUBLIC_BASE_URL=http://localhost:8080 \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=replace-with-private-token \
INCLUDE_ADMIN_READINESS=true \
make support-bundle
```

The support bundle stores summaries only. It does not retain raw public feed
bodies, full validation reports, raw `.env` files, raw database dumps, raw
telemetry, private vendor payloads, unredacted logs, or credential values. If
`DATABASE_URL` is supplied, migration status is collected by running the Go
migrator directly inside the script and sanitizing captured output.

## Admin URL And Token Rules

`PUBLIC_BASE_URL` must be `http://` or `https://` and must not contain userinfo,
query strings, fragments, or embedded credentials. The scripts never send auth
headers to `PUBLIC_BASE_URL`.

`ADMIN_BASE_URL` may default to `PUBLIC_BASE_URL` only when the public URL is
loopback. For a non-loopback public URL, authenticated admin checks require an
explicit admin URL. Non-loopback admin URLs must use HTTPS. Admin URLs with
userinfo, query strings, fragments, embedded credentials, or plain HTTP on a
non-loopback host are rejected.

Authenticated admin calls use exactly:

```text
Authorization: Bearer <ADMIN_TOKEN>
```

Cookie auth is not supported by these scripts. The scripts must not print or
write admin tokens, Authorization headers, cookies, JWTs, or CSRF values.

Admin requests do not follow redirects. Public feed requests may record
redirect metadata.

## Validation Behavior

Operator smoke always records validator tooling state through:

```bash
scripts/check-validators.sh
```

`SKIP_VALIDATORS=true` makes missing validator tooling non-fatal and skips
validation API calls. `STRICT_VALIDATORS=true` makes missing or misconfigured
validator tooling fatal.

When validation API calls run, the scripts use only allowlisted bodies:

| Feed type | Validator ID |
| --- | --- |
| `schedule` | `static-mobilitydata` |
| `vehicle_positions` | `realtime-mobilitydata` |
| `trip_updates` | `realtime-mobilitydata` |
| `alerts` | `realtime-mobilitydata` |

Only summary fields are stored by default: feed type, validator ID, status,
error count, warning count, info count, validator name, and validator version.

For the Phase 46 validator-health flow, use:

```bash
make validator-health
```

`scripts/validator-health.sh --dry-run` requires no network, database, Docker,
`ADMIN_TOKEN`, or running app. Authenticated checks require `ADMIN_TOKEN` and a
safe `ADMIN_BASE_URL` unless `PUBLIC_BASE_URL` is loopback and safely defaults.
Without `ADMIN_TOKEN`, the script does not call private admin routes and
records local validator tooling status only. `RUN_VALIDATORS=true` performs the
admin-only `run_all` action against `/admin/operations/validation-health`;
`STRICT_VALIDATOR_HEALTH=true` exits non-zero on blocked, failed, missing,
misconfigured, artifact-unavailable, or stale health states.

The script writes `summary.json`, `summary.md`, `manifest.json`, and
`manifest.md` under `.cache/validator-health/<timestamp>` by default. It
rejects `docs/evidence` and evidence-like output paths even when custom output
directories are allowed. These files are private diagnostics only, not evidence
packets or consumer submission artifacts.

## Synthetic AVL Dry-Run

Both helpers record the synthetic AVL dry-run status using the deterministic
fixture command:

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/minimal-gps.json
```

This is transform-only synthetic output. It is not telemetry ingest proof, real
vendor compatibility proof, production AVL reliability proof, or ETA quality
proof.

## Synthetic Telemetry Simulator

Use the Phase 44 simulator when you want to send synthetic events through the
real authenticated ingest path instead of only testing AVL transform shape:

```bash
make telemetry-simulator
RUN_MATCHER=true make telemetry-simulator
```

The simulator writes private diagnostics under `.cache/telemetry-simulator/`.
It is still synthetic-only local/reference diagnostics, not a support bundle,
not an evidence packet, not real vendor compatibility proof, and not
production AVL reliability proof. See
[Telemetry Simulator And Device Trial](telemetry-simulator-and-device-trial.md).

## What Is Safe To Share

The copy/paste summaries printed by the scripts are intended to be shareable
after the operator reviews them. They include output directory, public feed
summary, admin boundary result, authenticated readiness status, validator
tooling status, validation API status, AVL dry-run status, and:

```text
external_evidence_created=false
consumer_statuses_changed=false
```

Before sharing bundle files, review the generated manifest and confirm no local
private material was added outside the script. The support bundle runs a final
redaction scan for secret-shaped values and fails closed if it detects obvious
credential material.

## What Must Never Be Shared

Do not share or commit:

- admin tokens, device tokens, JWTs, CSRF values, cookies, API keys, or bearer
  credentials;
- database URLs with passwords;
- private keys, SSH keys, certificates, ACME material, or private key paths;
- raw `.env` files;
- raw database dumps;
- raw private telemetry;
- private vendor payloads;
- webhook URLs or notification credentials;
- unredacted logs;
- private correspondence, portal screenshots, or target-originated consumer
  artifacts that have not been reviewed under the evidence workflow.

## Relationship To Evidence

Smoke output and support bundles are private operator diagnostics by default.
They are not evidence packets. A future evidence phase would need to review,
redact, retain, inventory, and claim-map specific artifacts before any output
could support stronger public claims.

Do not change consumer statuses from these commands. All seven consumer and
aggregator targets remain `prepared` unless retained target-originated evidence
supports a target-specific status change.
