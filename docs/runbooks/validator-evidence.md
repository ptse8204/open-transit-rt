# Runbook: Production Validator Evidence

This runbook records how to package validator results from a real deployment.

## Purpose

Retain complete validator history for schedule and realtime feeds without selectively omitting warnings or errors.

## Required Feed Types

Collect records for:

- Static GTFS (`schedule.zip`).
- GTFS-RT Vehicle Positions.
- GTFS-RT Trip Updates.
- GTFS-RT Alerts.

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

## Truthfulness Guardrail

Validator success is quality evidence, not consumer acceptance evidence.
