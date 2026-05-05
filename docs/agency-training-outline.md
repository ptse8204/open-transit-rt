# Agency Training Outline

This outline supports a small-agency Open Transit RT pilot. It is not a
certification program, paid support package, SLA, hosted service, or compliance
approval.

## Session 1: GTFS And GTFS Realtime Basics

Topics:

- static GTFS schedule data: routes, stops, trips, stop times, calendars,
  shapes, frequencies, and blocks;
- agency-local service day and after-midnight service;
- GTFS Realtime feed types: Vehicle Positions, Trip Updates, and Alerts;
- why Vehicle Positions are the first production-directed realtime output;
- why Trip Updates quality needs conservative diagnostics and real-world review.

Outcome:

- participants can explain the difference between schedule data, vehicle
  telemetry, Vehicle Positions, Trip Updates, and Alerts.

## Session 2: Local Demo And Agency App Flow

Reference: `docs/tutorials/agency-first-run.md`.

Topics:

- prerequisites: Docker with Compose support and `curl`;
- `make agency-app-up`;
- local public feed paths;
- Operations Console path;
- local demo limitations;
- stopping and resetting the local app.

Outcome:

- participants can run or observe the local demo and identify what it proves and
  what it does not prove.

## Session 3: Real GTFS Onboarding

Reference: `docs/tutorials/real-agency-gtfs-onboarding.md`.

Topics:

- GTFS ownership and permission;
- public-safe handling;
- metadata approval;
- service date, timezone, route, stop, shape, frequency, and block review;
- import and publish path;
- final public-feed root boundary.

Outcome:

- participants can prepare an approved GTFS source or document why the GTFS path
  is blocked.

## Session 4: GTFS Validation Triage

Reference: `docs/tutorials/gtfs-validation-triage.md`.

Topics:

- validator tooling;
- import validation versus canonical validation;
- blocking errors, warnings, and informational findings;
- how validation findings affect public claims;
- why validator success is not consumer acceptance.

Outcome:

- participants can read validation results and identify the next operator-owned
  fix or blocker.

## Session 5: GTFS Studio Basics

Topics:

- GTFS Studio typed draft editing;
- draft versus published feed separation;
- when to use import versus Studio editing;
- publish review and audit expectations;
- limitations of the current minimal server-rendered UI.

Outcome:

- participants understand that drafts are not active published data until
  reviewed and published.

## Session 6: Device Token Safety

Reference: `docs/tutorials/device-token-lifecycle.md`.

Topics:

- Bearer token behavior;
- agency/device/vehicle binding;
- rotate and rebind flow;
- one-time token display;
- secure storage expectations;
- compromise response.

Outcome:

- participants can explain how to keep device tokens out of public docs, logs,
  screenshots, and evidence.

## Session 7: AVL And Vendor Adapter Pilot Boundary

Reference: `docs/tutorials/device-avl-integration.md`.

Topics:

- `/v1/telemetry` payload contract;
- simulator path;
- agency-owned, deployment-owned, vendor-owned, or private adapter patterns;
- Phase 29B synthetic dry-run adapter;
- why synthetic adapter output is not real vendor compatibility proof.

Outcome:

- participants can choose a telemetry path or document the blocker without
  making hardware, vendor, or AVL reliability claims.

## Session 8: Operations Console Setup Checklist

Topics:

- `/admin/operations`;
- `/admin/operations/setup`;
- feed URLs and validation state;
- telemetry freshness;
- device bindings;
- evidence and consumer status views;
- protected admin/debug boundary.

Outcome:

- participants can use the console for pilot review without treating it as
  consumer acceptance or compliance proof.

## Session 9: Validation And Evidence

References:

- `docs/compliance-evidence-checklist.md`;
- `docs/evidence/redaction-policy.md`;
- `docs/agency-owned-domain-readiness.md`.

Topics:

- local demo evidence;
- hosted/operator pilot evidence;
- final-root evidence;
- validator records;
- scorecards;
- redaction review;
- evidence gaps.

Outcome:

- participants can separate implementation capability, pilot evidence,
  deployment proof, and third-party confirmation.

## Session 10: Consumer Submission Boundaries

Reference: `docs/evidence/consumer-submissions/submission-workflow.md`.

Topics:

- prepared packets;
- official-path verification;
- target-originated evidence;
- status transitions;
- why all seven targets remain `prepared` without later evidence.

Outcome:

- participants can explain why prepared packets are not submissions and why
  consumer acceptance cannot be inferred from validation or public fetches.

## Session 11: Support And Security Reporting

References:

- `docs/support-boundaries.md`;
- `SECURITY.md`.

Topics:

- community support;
- useful public-safe bug reports;
- private vulnerability reporting;
- what not to post publicly;
- no paid support, SLA, hosted operations, or legal commitment.

Outcome:

- participants know how to report issues without exposing secrets or private
  operator artifacts.

