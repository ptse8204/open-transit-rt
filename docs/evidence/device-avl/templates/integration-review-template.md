# Device / AVL Integration Review Template

This template is for future public-safe evidence only. Do not fill it with fake device data, fake vendor approvals, fake hardware certifications, private AVL payloads, credentials, or raw private telemetry.

## Summary

| Field | Value |
| --- | --- |
| Review date |  |
| Reviewer |  |
| Agency or pilot environment |  |
| Integration type | Simulator / pilot device / real device / vendor adapter / agency sidecar |
| Source |  |
| Permission recorded | Yes / No |
| Public-safe to commit | Yes / No |
| Redaction review completed | Yes / No |

## Identifiers

| Identifier type | Value or redaction | Public-safe or synthetic? | Notes |
| --- | --- | --- | --- |
| Agency ID |  |  |  |
| Device ID |  |  |  |
| Vehicle ID |  |  |  |
| Trip hint, if used |  |  |  |
| Vendor identifier, if relevant |  |  |  |

Do not include private device serial numbers, private vehicle IDs, vendor account IDs, or tokens unless the agency has explicitly approved them as public-safe. Prefer synthetic identifiers in committed evidence.

## Integration Path

Describe how telemetry reached Open Transit RT:

- endpoint used:
- authentication handled by:
- adapter, sidecar, or vendor-owned middleware:
- timestamp source:
- GPS source:
- field validation before forwarding:

## Commands Or Checks Run

List redacted commands or checks. Do not include bearer tokens, admin tokens, DB URLs with passwords, private URLs, or raw private logs.

```text

```

## Results

| Check | Result | Evidence path or note |
| --- | --- | --- |
| Telemetry accepted |  |  |
| Operations Console freshness reviewed |  |  |
| Vehicle Positions reviewed |  |  |
| Trip Updates reviewed, if relevant |  |  |
| Validator reviewed, if relevant |  |  |

## Redactions

List every redaction and why it was needed:

-

## What This Proves

-

## What This Does Not Prove

- This does not prove certified hardware or vendor support unless separate retained evidence says so.
- This does not prove production AVL reliability without real operating evidence.
- This does not prove consumer acceptance.
- This does not prove CAL-ITP/Caltrans compliance.

## Next Action

-
