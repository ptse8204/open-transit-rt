# Phase 12 Hosted Evidence Intake Packet

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Status: pending evidence collection

## Purpose

This packet creates the required Phase 12 hosted evidence artifact slots without fabricating proof.

Every artifact in this folder is intentionally marked pending until a deployment operator runs the listed commands against a real hosted environment and attaches or records the resulting outputs.

## Required Inputs

Set these before collecting evidence:

```sh
export ENVIRONMENT_NAME="<hosted-environment-name>"
export PUBLIC_BASE_URL="https://<canonical-feed-host>"
export FEED_BASE_URL="$PUBLIC_BASE_URL/public"
export ADMIN_BASE_URL="https://<admin-or-origin-host-if-different>"
export ADMIN_TOKEN="<redacted-admin-token>"
```

Do not commit secrets, bearer tokens, private keys, or unredacted internal hostnames unless the deployment owner approves that disclosure.

## Collection Script

The repo now includes an executable collector:

```sh
ENVIRONMENT_NAME="$ENVIRONMENT_NAME" \
PUBLIC_BASE_URL="$PUBLIC_BASE_URL" \
ADMIN_BASE_URL="$ADMIN_BASE_URL" \
ADMIN_TOKEN="$ADMIN_TOKEN" \
./scripts/collect-hosted-evidence.sh
```

Equivalent Make/Task targets:

```sh
make collect-hosted-evidence
task evidence:hosted
```

The script collects generic hosted fetch, TLS, validation, and manual scorecard artifacts. Deployment-owned monitoring exports, alert lifecycle records, backup job history, restore transcripts, reverse proxy configs, renewal proof, and scorecard scheduler history must still be attached by an operator.

After filling the packet, run:

```sh
EVIDENCE_PACKET_DIR="docs/evidence/captured/$ENVIRONMENT_NAME/<UTC-date>" \
make audit-hosted-evidence
```

The audit must pass before this packet can support Phase 12 closure.

## Artifact Files

- `public-feed-proof-2026-04-22.md`
- `reverse-proxy-tls-2026-04-22.md`
- `validator-record-2026-04-22.md`
- `monitoring-alert-2026-04-22.md`
- `backup-restore-drill-2026-04-22.md`
- `scorecard-export-2026-04-22.md`
- `operator-collection-commands-2026-04-22.md`

## Claim Boundary

This packet is not evidence that the hosted deployment is ready. It is the intake location for the missing hosted artifacts.

Phase 12 remains open until these pending artifacts are replaced or supplemented with real outputs from one hosted deployment environment.
