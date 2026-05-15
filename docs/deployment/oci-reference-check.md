# OCI Reference Check

`scripts/oci-reference-check.sh` is a private diagnostic helper for a
self-hosted OCI/OCL-style reference deployment.

It is product support tooling, not evidence. It does not write
`docs/evidence`, submit feeds, contact consumers, change consumer statuses,
require Go on the remote host, print secret values, or claim compliance,
agency adoption, consumer acceptance, final-root readiness, hosted service
availability, production readiness, vendor compatibility, SLA coverage,
uptime, or production-grade ETA quality.

## Run A Dry Run

```bash
PUBLIC_BASE_URL=https://feeds.example.org make oci-reference-check
```

For a no-network shape check:

```bash
OUTPUT_DIR=.cache/oci-reference-check/dry-run \
FORCE=true \
scripts/oci-reference-check.sh \
  --public-base-url https://feeds.example.org \
  --dry-run
```

The default output path is:

```text
.cache/oci-reference-check/<timestamp>
```

## What It Checks

- Public five-feed fetches through `scripts/validate-public-feeds.sh`.
- Public headers, byte counts, checksums, and validator state where available.
- Deployment helper environment presence, with values withheld.
- Optional SSH loopback health checks for the five services when `OCI_HOST`
  access is configured, using `/healthz` on each loopback service port.
- Backup and restore-drill configuration presence, with values withheld. The
  helper recognizes `RESTORE_DATABASE_URL` / `RESTORE_BACKUP_FILE` and the
  optional `RESTORE_DRILL_DATABASE_URL` / `RESTORE_DRILL_BACKUP_FILE` aliases.
- Telemetry simulator dry-run guidance.
- All-false claim flags.

The loopback health checks use `curl` on the remote host through SSH. They do
not require Go on the remote host.

## Optional SSH Loopback Checks

Set these only in the private operator shell:

```bash
OCI_HOST=replace-with-host
OCI_USER=opc
OCI_KEY=/path/to/private/key
```

The summary records whether these inputs are configured, but it does not print
the host, key path, tokens, database URLs, private keys, or populated env
values.

## Strict Mode

Use strict mode when a blocker should fail the command:

```bash
scripts/oci-reference-check.sh \
  --public-base-url https://feeds.example.org \
  --strict
```

Strict mode can fail for public fetch blockers, validator blockers, or SSH
loopback blockers. Missing optional environment values remain visible as
configuration gaps, not proof that the deployment is unhealthy.

## Output Files

The helper writes:

```text
summary.json
summary.md
manifest.json
manifest.md
public-feeds/
```

The nested `public-feeds/` directory contains the off-host validation summary,
fetched public artifacts when not in dry-run mode, headers, checksums, and
validator output when validators are installed.

Keep the output private. Review and redact before sharing with anyone.
