# Implementation Plan Template

This template is for planning an operator-run Open Transit RT deployment. It is
not production-readiness proof, compliance proof, agency adoption proof,
consumer acceptance proof, hosted service proof, or vendor certification.

## Scope

- Operator: `<operator>`
- Deployment owner: `<owner>`
- Target environment: `<local/reference/pilot/other>`
- Public feed root: `<root or unresolved>`
- GTFS source: `<source>`
- AVL source: `<source>`
- Realtime feeds in scope: `<Vehicle Positions / Trip Updates / Alerts>`

## Work Plan

| Area | Owner | Inputs | Output | Review Gate |
| --- | --- | --- | --- | --- |
| Static GTFS | `<owner>` | `<input>` | `<output>` | `<gate>` |
| Vehicle telemetry | `<owner>` | `<input>` | `<output>` | `<gate>` |
| Vehicle Positions | `<owner>` | `<input>` | `<output>` | `<gate>` |
| Trip Updates | `<owner>` | `<input>` | `<output>` | `<gate>` |
| Alerts | `<owner>` | `<input>` | `<output>` | `<gate>` |
| Operations | `<owner>` | `<input>` | `<output>` | `<gate>` |

## Verification

- `make validate`: `<result>`
- `make test`: `<result>`
- `make smoke`: `<result>`
- Public feed fetch checks: `<result>`
- Validator checks: `<result or blocker>`
- Consumer submission status: `<prepared/submitted only with evidence>`

## Rollback

- Backup taken before change: `<yes/no/path redacted>`
- Migration status before: `<status>`
- Migration status after: `<status>`
- Restore or rollback owner: `<owner>`
- Restore drill result: `<result or blocker>`

## Claim Review

List each proposed external claim and the retained evidence that supports it.
If no retained evidence exists, mark the claim unsupported and remove it from
public-facing material.
