# Runbook: Production Validator Evidence

This runbook records how to package validator results from a real deployment.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/validator-record-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/validator-record-2026-04-24.md`

The local packet contains records for schedule, Vehicle Positions, Trip Updates, and Alerts, but all four validator runs failed in the local environment. It is retained as failure evidence and must not be used to claim validator-clean feeds. The OCI pilot packet contains hosted validator records for that recorded pilot scope only.

## Purpose

Retain complete validator history for schedule and realtime feeds without selectively omitting warnings or errors.

## Required Feed Types

Collect records for:

- Static GTFS (`schedule.zip`).
- GTFS-RT Vehicle Positions.
- GTFS-RT Trip Updates.
- GTFS-RT Alerts.

Phase 17 scheduled validation may use:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=<redacted-admin-token> \
scripts/pilot-ops.sh validator-cycle --dry-run
```

The live run writes `validator-cycle-YYYY-MM-DD.json` and per-feed response files to `EVIDENCE_OUTPUT_DIR`.

## Required Record Fields

Each validator record must include:

- Environment name.
- Feed type.
- Validator ID/tool name.
- Validator version.
- Run timestamp (UTC).
- Feed revision/version context.
- Full result summary (pass/warn/fail).
- Link or attachment to full validator output.

Do not remove warnings/errors from archived output.

## Retention Rules

- Keep latest successful run plus recent failed/warn runs.
- Retain enough history to demonstrate operational trends.
- Keep redaction notes if sensitive infra details are removed.

## Output Artifact

Use:

- `docs/evidence/templates/validator-record-template.md`

Store under `docs/evidence/captured/<environment>/validator-record-YYYY-MM-DD.md`.

Evidence labels:

- `validator-cycle-YYYY-MM-DD.json`: `private/operator-only` until reviewed; `safe-to-commit-after-review` if redacted and complete.
- Full validator responses: usually `safe-to-commit-after-review`, but review for private paths/internal hosts.
- Admin token, private admin URL, temporary validator work dirs containing private paths, and unredacted logs: `never-commit`.

## Truthfulness Guardrail

Validator success is quality evidence, not consumer acceptance evidence.
