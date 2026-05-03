# Secret Exposure Incident Template

Template only. Do not fill with fake exposure records or placeholder operational artifacts.

## Required Fields

- Start time:
- Affected environment:
- Affected agency:
- Affected public URLs or services:
- Detection source:
- Operator:
- Severity:
- Timeline:
- Action taken:
- Evidence retained:
- Redaction review:
- Follow-up:
- Claim boundary:

## Exposure Details

- Secret category: admin JWT / CSRF / device token pepper / device token / DB password / TLS or ACME material / webhook or notification credential / other
- Exposure location:
- Was the value committed to git history:
- Rotation or revocation completed:
- Old credential verified invalid:
- Services restarted:
- Post-rotation verification:

## Required Response

Deleting a file is not enough when a real secret was exposed. Rotate or revoke the credential, assess history and backups for exposure, preserve redacted response notes, and do not rewrite git history without maintainer approval.

## Claim Boundary

Secret response records do not prove hosted SaaS availability, paid support/SLA coverage, compliance, consumer acceptance, or production multi-tenant readiness.
