# Phase 45 — GTFS Quality Triage Loop

Phase 45 adds a private Operations Console path from static GTFS validator
notices to bounded operator triage actions.

## Scope

- Authenticated route: `/admin/operations/gtfs-quality`.
- Canonical source: MobilityData static GTFS validator results stored in
  `validation_report`.
- Internal source: Open Transit RT GTFS import validation results stored in
  `validation_report`.
- Rerun action: admin-only POST that uses only the authenticated agency active
  published schedule ZIP and the server-side `static-mobilitydata` validator
  mapping.

## Boundaries

- The page is an authenticated admin/operator UI, not a public route.
- Validator output is diagnostics/supporting signal only.
- The route does not claim consumer acceptance, CAL-ITP/Caltrans compliance,
  agency adoption, hosted SaaS availability, production readiness, vendor
  compatibility, or production-grade ETA quality.
- The route does not auto-edit agency GTFS.
- The route does not create, export, download, retain, or write evidence
  packets, and it does not write to `docs/evidence`.
- Stored raw validator reports may contain `raw_report`, stdout, stderr, argv,
  or temp paths from the existing validation flow. The derived triage view model
  and rendered HTML omit those fields.

## Behavior

- GET permits read-only, operator, editor, and admin roles and sets
  `Cache-Control: no-store`.
- POST is admin-only, CSRF-protected for browser cookie auth, capped before form
  parsing, and sets `Cache-Control: no-store` for controlled responses.
- POST accepts only `action` and CSRF fields. Request-supplied validator IDs,
  commands, paths, URLs, argv/args, artifacts, report paths, schedule paths, and
  timeout fields are rejected before validator execution.
- Missing active schedule shows a blocker/action message instead of a 500.
- Unsupported methods are rejected safely.
- `?agency_id=<same>` passes; conflicting agency IDs are rejected.

## Triage Model

The derived triage model includes only bounded fields:

- source: `canonical_validator` or `internal_importer`;
- family and codes;
- severity: `blocking`, `needs_review`, `informational`, or `unknown`;
- count;
- operator summary;
- why it matters;
- recommended action;
- capped samples;
- overflow count.

Issue groups are capped at 100, samples are capped at 5 per group, sample text
is truncated, and ordering is deterministic by severity, family, code, and
count.

## Notice Families

The first taxonomy covers:

- `expired_calendar`;
- `route_short_name_too_long`;
- `unused_shape`;
- missing or foreign-key references;
- bad `stop_times`;
- duplicate IDs;
- calendar and service-date issues;
- shape ordering or `shape_dist_traveled` issues;
- frequency issues;
- `block_id` and block transition issues;
- unknown notices.

## Tests

Phase 45 adds coverage for auth/CSRF, cache headers, agency isolation, explicit
routing, unsupported methods, latest-result selection, stale active feed
behavior, strict POST forms, POST body-size capping, rerun success/failure,
source separation, malformed and hostile reports, HTML escaping, bounded output,
large report triage, large page render, and local benchmark diagnostics.

Benchmarks are not part of default `make test`; they are local engineering
diagnostics only and are not production capacity, SLA, evidence, compliance, or
consumer-readiness proof.
