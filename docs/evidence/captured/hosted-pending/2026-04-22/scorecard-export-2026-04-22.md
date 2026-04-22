# Hosted Scorecard Export and Job History Evidence

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Status: missing

## Required Evidence

- Timestamped compliance scorecard exports from the hosted environment.
- Artifact hashes/checksums.
- Export schedule or job definition.
- Recent successful job history.
- Retention policy and storage boundary.

## Manual Export Command

```sh
mkdir -p "$ENVIRONMENT_NAME/scorecard"

timestamp="$(date -u '+%Y-%m-%dT%H%M%SZ')"
curl -sS -X POST "$ADMIN_BASE_URL/admin/compliance/scorecard" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{}' \
  | tee "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.json"

shasum -a 256 "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.json" \
  | tee "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.sha256.txt"
```

## Required Attachments

- `scorecard-<timestamp>.json`
- `scorecard-<timestamp>.sha256.txt`
- `scorecard-export-job-definition.txt`
- `scorecard-export-job-history.txt`
- `scorecard-retention-policy.txt`

## Required Summary To Fill

- Export format:
- Export trigger:
- Latest export timestamp UTC:
- Latest export artifact:
- Latest export SHA-256:
- Scheduled job reference:
- Recent successful job history:
- Retention policy:
- Storage boundary:

## Blocker

No hosted scorecard export job definition, job history, retention policy, or hosted scorecard artifact was available when this intake packet was created.
