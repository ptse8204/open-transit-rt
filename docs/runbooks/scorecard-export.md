# Runbook: Scorecard Export Evidence

This runbook defines how to capture evidence for compliance scorecard exports from a real deployment.

## Required Evidence

Capture:

- Export format used (JSON, CSV, or both).
- Export trigger mode (scheduled and/or manual).
- Export timestamp in UTC.
- Environment identifier.
- Storage location and retention policy.

## Scheduled Export Evidence

If scheduled exports exist, record:

- scheduler/job definition reference
- run history excerpt
- at least one successful recent export

## Manual Export Evidence

If manual fallback is used, record:

- operator identity
- command or endpoint used
- resulting artifact path

## Output Artifact

Use:

- `docs/evidence/templates/scorecard-export-template.md`

Store under `docs/evidence/captured/<environment>/scorecard-export-YYYY-MM-DD.md`.

## Truthfulness Guardrail

A recorded export proves scorecard extraction for that deployment window only; it does not by itself prove full compliance or consumer acceptance.
