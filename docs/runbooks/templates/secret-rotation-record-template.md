# Secret Rotation Record Template

Template only. Do not fill with fake rotation records or placeholder operational artifacts.

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

## Rotation Coverage

- Admin JWT secret:
- CSRF secret:
- Device token pepper:
- Device tokens:
- DB password:
- TLS/ACME material:
- Optional webhook/notification credentials:
- Phase 15 `.cache` secret findings reviewed:

## Verification

- New credential stored in private secret location:
- Old credential revoked or removed from active use:
- Services restarted:
- Readiness checks:
- Public feed checks:
- Device ingest checks, if applicable:
- Alert delivery proof, if applicable:

## Claim Boundary

Secret rotation records are operational evidence only. They do not prove hosted SaaS availability, paid support/SLA coverage, compliance, consumer acceptance, or production multi-tenant readiness.
