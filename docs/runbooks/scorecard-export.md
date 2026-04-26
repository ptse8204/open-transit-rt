# Runbook: Scorecard Export Evidence

This runbook defines how to capture evidence for compliance scorecard exports from a real deployment.

Latest captured packets:

- `docs/evidence/captured/local-demo/2026-04-22/scorecard-export-2026-04-22.md`
- `docs/evidence/captured/oci-pilot/2026-04-24/scorecard-export-2026-04-24.md`

The local packet records a manual JSON scorecard export. The OCI pilot packet records deployment/operator scorecard job evidence for that recorded pilot scope.

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

Phase 17 scorecard export dry-run:

```sh
ENVIRONMENT_NAME=<environment> \
EVIDENCE_OUTPUT_DIR=/opt/open-transit-rt/evidence/<UTC-date> \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=<redacted-admin-token> \
scripts/pilot-ops.sh scorecard-export --dry-run
```

The live export writes `scorecard-export-YYYY-MM-DD.json` to `EVIDENCE_OUTPUT_DIR`.

## Output Artifact

Use:

- `docs/evidence/templates/scorecard-export-template.md`

Store under `docs/evidence/captured/<environment>/scorecard-export-YYYY-MM-DD.md`.

Evidence labels:

- `scorecard-export-YYYY-MM-DD.json`: `safe-to-commit-after-review`.
- Admin tokens, private admin origins, private job logs, and secret files used by scheduled exports: `never-commit`.

## Truthfulness Guardrail

A recorded export proves scorecard extraction for that deployment window only; it does not by itself prove full compliance or consumer acceptance.
