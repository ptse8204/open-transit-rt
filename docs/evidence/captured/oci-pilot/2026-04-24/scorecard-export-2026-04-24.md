# Hosted Scorecard Export Evidence

- Environment: `oci-pilot`
- Capture date (UTC): 2026-04-24
- Operator: Codex operator session using OCI pilot admin credentials

## Export Details

- Manual export artifacts: `artifacts/scorecard/latest-scorecard.json` and timestamped `scorecard-*.json` files.
- Scheduled job evidence: `artifacts/operator-supplied/scorecard-job-definition-and-history.txt`.
- Latest scorecard status after the current-live recheck: validation green, discoverability green, consumer ingestion red.

## Claim Boundary

The scorecard proves repeatable hosted export and retained job history. The consumer-ingestion score remains red because external consumer acceptance is outside Phase 12.
