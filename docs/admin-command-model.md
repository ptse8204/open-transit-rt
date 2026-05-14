# Admin Command Model

Open Transit RT browser-controlled operations use a private command model for
bounded workflows. It is not a public API, shell runner, plugin loader,
evidence collector, release publisher, or consumer submission system.

## Result Shape

Command results use this bounded JSON shape:

```json
{
  "action": "validation_health.refresh",
  "status": "ok",
  "started_at": "2026-05-14T00:00:00Z",
  "completed_at": "2026-05-14T00:00:01Z",
  "summary": "Validation health summary refreshed from existing private records.",
  "next_actions": ["Review stale or blocked rows before stronger readiness language."],
  "claim_flags": {
    "external_evidence_created": false,
    "consumer_statuses_changed": false,
    "compliance_claimed": false,
    "production_readiness_claimed": false,
    "agency_approval_claimed": false,
    "consumer_acceptance_claimed": false,
    "public_launch_claimed": false,
    "hosted_saas_claimed": false,
    "vendor_compatibility_claimed": false,
    "hardware_certification_claimed": false,
    "sla_claimed": false,
    "uptime_guarantee_claimed": false,
    "production_grade_eta_claimed": false
  }
}
```

Allowed statuses are:

- `ok`: completed in the private console.
- `needs_review`: completed, but operator review is still needed.
- `blocked`: not run or not completed because a bounded precondition failed.
- `failed`: execution started and hit a bounded runtime/tooling failure.

These statuses are private workflow outcomes. They do not mean compliant,
accepted, approved, production-ready, vendor-compatible, SLA-backed, or
production-grade ETA.

## Safe Command Ladder

| Level | Meaning | Browser rule |
| --- | --- | --- |
| `read_only_refresh` | Re-read or recompute existing private state without writes. | Allowed for private read roles when CSRF and agency checks pass. |
| `dry_run` | Run a bounded check that writes no durable product state. | Role and CSRF requirements are action-specific. |
| `reversible_private_change` | Store private records or config with a clear review/supersede path. | Admin-only unless a later phase explicitly narrows it. |
| `publish_activate` | Change public feed bytes, active feed metadata, or externally visible state. | Admin-only with preview, confirmation, audit, and rollback copy. |
| `destructive_or_hard_to_reverse` | Delete, rotate, restore, overwrite, or irreversibly alter data/secrets. | Disabled by default unless separately authorized. |

## Safety Rules

Browser command requests must not supply:

- shell commands;
- argv/args;
- validator IDs or validator binary paths unless selected from a server-owned
  allowlist;
- artifact paths, output paths, report paths, schedule ZIP paths, realtime
  protobuf paths, or private filesystem paths;
- URLs for validator execution;
- timeouts;
- support-bundle destinations;
- evidence destinations.

Command responses must not include raw stdout, stderr, argv, raw validator
reports, raw external payloads, bearer tokens, cookies, database URLs, private
hostnames, or private filesystem paths.

## Current Phase 77 Actions

`validation_health.refresh` is the first migrated command. It refreshes the
derived validator-health summary from existing private records and server-owned
artifact checks. It writes nothing and changes no public feed output.

`validation_health.run_all` is documented as a future/admin command definition
for the existing validator run workflow. It is a private diagnostic write
because successful runs may store normal `validation_report` rows. It changes
no public feed output, creates no retained evidence, and moves no consumer
status.
