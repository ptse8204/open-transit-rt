# Integration Adapter Kit

This guide is the central map for integrating Open Transit RT with telemetry
sources, prediction engines, validators, monitoring, and feed-consumer
workflows.

It points to the detailed operator guides instead of duplicating them. It does
not prove certified vendor compatibility, production AVL reliability,
production-grade ETA quality, consumer acceptance, CAL-ITP/Caltrans compliance,
agency adoption, hosted SaaS availability, or final-root proof.

Post-60 productization makes external connection quality the default next work.
The intended integration shape is sidecar plus manifest plus conformance
testing, not arbitrary dynamic plugin loading. These contracts help operators
evaluate adapters safely; they do not create evidence or prove third-party
compatibility.

## Adapter Decision Tree

Use the existing boundary that matches the system being integrated:

| Need | Boundary | Detailed docs |
| --- | --- | --- |
| Run one guided local/reference operator trial | Deployment + onboarding + readiness + synthetic AVL dry-run | [Self-Hosted Operator Trial](tutorials/self-hosted-operator-trial.md) |
| Run smoke checks or collect safe diagnostics | Public feed checks + admin boundary checks + support bundle redaction | [Operator Smoke And Support Bundle](tutorials/operator-smoke-and-support-bundle.md) |
| Import a specific agency GTFS ZIP first | `make agency-pilot-up` / `scripts/agency-pilot-onboard.sh` | [Reusable Agency Onboarding](tutorials/reusable-agency-onboarding.md) |
| Send synthetic telemetry through real ingest | `make telemetry-simulator` / authenticated `POST /v1/telemetry` | [Telemetry Simulator And Device Trial](tutorials/telemetry-simulator-and-device-trial.md) |
| Send device, GPS, or AVL observations | Transform to `POST /v1/telemetry` | [Device And AVL Integration](tutorials/device-avl-integration.md) |
| Demonstrate an AVL transform without sending data | `cmd/avl-vendor-adapter --dry-run` with synthetic fixtures | [AVL fixture manifest](../testdata/avl-vendor/README.md) |
| Validate an external connector manifest | Sidecar/manifest contract and local checks | [Connector Plugin Contract](connectors/plugin-contract.md), [External Connection Readiness](external-connection-readiness.md) |
| Run synthetic adapter conformance | Offline synthetic conformance suite | [External Adapter Conformance](tutorials/external-adapter-conformance.md) |
| Swap or evaluate Trip Updates prediction | `internal/prediction.Adapter` | [Dependencies](dependencies.md), [Trip Updates requirements](requirements-trip-updates.md) |
| Validate GTFS or GTFS-Realtime artifacts | Server-side allowlisted validator IDs | [GTFS Validation Triage](tutorials/gtfs-validation-triage.md), [Dependencies](dependencies.md) |
| Monitor feeds or operations | Deployment-owned monitoring around existing endpoints/helpers | [Small-Agency Pilot Operations](runbooks/small-agency-pilot-operations.md) |
| Review readiness gaps and next actions | `/admin/operations/readiness` | [CAL-ITP Readiness Checklist](tutorials/calitp-readiness-checklist.md) |
| Prepare consumer/aggregator workflow | Public feed URLs and prepared packet records | [Consumer submission workflow](evidence/consumer-submissions/submission-workflow.md) |

## Practical Usage Path

1. Use [Self-Hosted Operator Trial](tutorials/self-hosted-operator-trial.md)
   when you want the complete guided local/reference path.
2. Use [Operator Smoke And Support Bundle](tutorials/operator-smoke-and-support-bundle.md)
   when you need repeatable smoke output or a redaction-safe support bundle.
3. Import and publish GTFS with `make agency-pilot-up`.
4. Review the printed feed URLs and `/public/feeds.json` metadata.
5. Choose a telemetry adapter path: direct device POST, agency-owned script,
   sidecar service, vendor-owned middleware, or private operator process.
6. Run `make telemetry-simulator` when you need synthetic-only telemetry to
   enter through real device-token auth and `/v1/telemetry`.
7. Run the synthetic AVL dry-run adapter against committed fixtures to verify
   the transform pattern and diagnostics shape.
8. Map real private AVL payloads outside this public repo.
9. Send only validated telemetry to `/v1/telemetry` with deployment-owned
   device credentials.
10. Review the Operations Console and public Vehicle Positions output.
11. Use `/admin/operations/readiness` to review CAL-ITP-style readiness gaps
   without converting workflow status into a compliance or acceptance claim.

The Phase 37 onboarding flow establishes the active schedule/feed baseline.
Phase 38 integration work starts after that baseline exists; it does not
replace GTFS import, feed publication metadata review, or device credential
management.

## Connector Manifest Checks

Post-60 connector work uses `open-transit-rt.connector.v1` manifests. Run:

```bash
make external-connection-check
```

The check validates synthetic connector manifests for telemetry sources,
prediction, validators, monitoring/export, and consumer/discovery workflows.
It is local validation only. It does not load Go plugins, start sidecars,
contact consumers, send notifications, create evidence, mutate consumer
statuses, or prove compatibility, compliance, production readiness, or
production-grade ETA quality.

See [Connector Plugin Contract](connectors/plugin-contract.md) for the manifest
shape and [External Connection Readiness](external-connection-readiness.md) for
operator review questions.

## Generic Connector Examples

Stage 5 adds generic connector examples under `examples/connectors/` for:

- telemetry HTTP poller shape;
- telemetry CSV replay shape;
- prediction sidecar stub shape;
- monitoring/export summary shape.

Telemetry examples share a small stdlib-only SDK-style helper under
`examples/connectors/sdk/telemetry` for synthetic dry-run normalization,
fail-closed drop reasons, quality checks, stale/future timestamp handling, and
no-send event output. It is example helper code for adapter authors, not a
production vendor SDK, hardware certification, production AVL reliability
claim, or vendor compatibility claim.

The predictor sidecar stub uses `examples/connectors/sdk/prediction` for a
sanitized dry-run request/response shape, withheld-output diagnostics,
Vehicle Positions independence, and explicit no public Trip Updates mutation.
It is not a named predictor integration, runtime compatibility claim, or
production-grade ETA quality claim.

These examples use synthetic fixtures only and are included in local manifest
and conformance checks. They are developer examples, not real vendor adapters,
not consumer/discovery automation, not evidence packets, and not proof of
vendor compatibility, production AVL reliability, production-grade ETA quality,
consumer acceptance, production readiness, or CAL-ITP/Caltrans compliance.

## Adapter Conformance Suite

Use the offline conformance suite to check required synthetic cases:

```bash
make adapter-conformance
```

The private Operations Console also renders generated connector test
instructions at `/admin/operations/connectors/tests`. That page is guidance
only: it does not execute commands from the browser, run validators, start
sidecars, contact external systems, write evidence, or change consumer
statuses.

The CLI covers telemetry malformed/stale/future/wrong-agency/unknown-device/
low-quality/duplicate/out-of-order cases, prediction timeout/malformed/stale/
wrong-agency/low-confidence cases, validator allowlisting, and monitoring
redaction/no-send defaults. It does not send network traffic, run validators,
contact consumers, write evidence, or change repo state. See
[External Adapter Conformance](tutorials/external-adapter-conformance.md).

## Telemetry Adapter Path

Telemetry integrations should transform external payloads into the existing
Open Transit RT telemetry event shape before calling `/v1/telemetry`.

Acceptable adapter ownership patterns are:

- agency-owned adapter scripts;
- deployment-owned sidecar services;
- vendor-owned middleware that calls Open Transit RT;
- private operator integration processes.

The detailed contract, payload fields, response behavior, token handling,
Operations Console checks, and troubleshooting matrix live in
[Device And AVL Integration](tutorials/device-avl-integration.md). That guide
is the source for operator-level telemetry instructions.

## Synthetic Telemetry Simulator

`cmd/telemetry-simulator` and `scripts/telemetry-simulator.sh` provide a
synthetic-only way to test the real authenticated ingest path:

```bash
make telemetry-simulator
```

The simulator posts to `/v1/telemetry` with a device bearer token. Optional
`RUN_MATCHER=true` mode reads accepted rows back from Postgres after HTTP
ingest and runs the existing matcher and Vehicle Positions debug builder for
private diagnostics.

Simulator output under `.cache/telemetry-simulator/` is private diagnostics. It
is not an evidence packet, not vendor compatibility proof, not production AVL
reliability proof, and not production-grade ETA proof.

## `/v1/telemetry` Contract Pointer

Adapters call:

```text
POST /v1/telemetry
Authorization: Bearer <device-token>
Content-Type: application/json
```

Required payload fields are `agency_id`, `device_id`, `vehicle_id`,
`timestamp`, `lat`, and `lon`. Optional fields include `driver_id`, `bearing`,
`speed_mps`, `accuracy_m`, and `trip_hint`.

Do not treat this summary as a second API definition. Use
[Device And AVL Integration](tutorials/device-avl-integration.md) for the
current confirmed request and response details.

## Synthetic AVL Adapter Examples

`internal/avladapter` and `cmd/avl-vendor-adapter` provide a dry-run adapter
kit example. The command reads a synthetic mapping file and a synthetic payload
fixture, then prints transformed Open Transit RT telemetry JSON to stdout and
diagnostics JSON to stderr.

Example:

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/minimal-gps.json
```

Dry-run output is transform output only; it is not telemetry ingest status and
it is not vendor compatibility proof. Phase 48 also adds a private `--send`
mode that posts transformed synthetic/operator-reviewed records only to the
existing authenticated `/v1/telemetry` boundary with env-referenced credentials
and redacted `.cache` diagnostics. Send mode is still not named vendor support,
real vendor compatibility proof, production AVL reliability proof, or retained
evidence.

See [testdata/avl-vendor/README.md](../testdata/avl-vendor/README.md) for the
fixture manifest and exact diagnostic codes.

## Mapping And Diagnostics Contract

The mapping file is the authority for emitted `agency_id`, `device_id`, and
`vehicle_id`. Vendor-looking identifiers in payload fixtures are lookup keys
only and cannot override mapped Open Transit RT identifiers.

Dry-run diagnostics are stable JSON objects with:

- `code`
- `severity`
- `message`
- optional `index`

Diagnostic codes are defined in `internal/avladapter`. Duplicate and
out-of-order dry-run diagnostics are batch-level review observations only.
They are not database ingest statuses such as `accepted`, `duplicate`, or
`out_of_order`.

Do not put tokens, endpoint URLs, auth headers, passwords, private keys,
database URLs, real vendor account IDs, private device identifiers, private
vehicle identifiers, or raw private telemetry in mappings or fixtures.

## External Predictor Lifecycle

External predictors must stay behind `internal/prediction.Adapter`. Vehicle
Positions, telemetry ingest, GTFS import, assignments, and audit state remain
owned by Open Transit RT.

Phase 49 adds a generic disabled-by-default HTTP sidecar boundary through
`TRIP_UPDATES_ADAPTER=external_http` and
`TRIP_UPDATES_ADAPTER=external_http_shadow`. The endpoint path is fixed to
`/v1/predict/trip-updates`; configuration requires an exact host allowlist,
safe URL shape, HTTPS except loopback test stubs, no redirects, byte/time caps,
and env-name-only bearer token lookup. Requests use sanitized DTOs and never
send raw telemetry events, raw assignments, payload JSON, device IDs, driver
IDs, score details, manual override IDs, audit fields, raw override reason
text, credentials, headers, or cookies.

Use `external_http_shadow` when evaluating a sidecar without changing the
public Trip Updates output. Shadow diagnostics are bounded redacted counts and
status only; they are not ETA-quality evidence.

Named predictor work should still follow this lifecycle:

1. Candidate review: document the predictor role, inputs, outputs, deployment
   shape, and replacement path.
2. Dependency and license review: update `docs/dependencies.md` before adding
   runtime coupling, vendored code, packaged services, or deployment assets.
3. Adapter contract tests: prove request construction, output normalization,
   wrong-agency/wrong-feed rejection, stale output rejection, and failure
   diagnostics.
4. Shadow or dry-run evaluation: compare outputs without changing public feed
   behavior unless explicitly approved.
5. Output validation: validate Trip Updates against active GTFS and GTFS-RT
   protobuf requirements before serialization.
6. Failure fallback: Vehicle Positions, telemetry ingest, assignments, admin
   workflows, and static GTFS publication must continue when the predictor is
   unavailable.
7. Evidence review: retain claim-specific evidence before making any stronger
   ETA quality, compatibility, compliance, or consumer-readiness claim.

TheTransitClock and other named predictors remain deferred optional integrations
unless a later phase explicitly implements and documents that named runtime
integration.

## Validator Integration Boundary

Validators are external tooling invoked through server-side allowlisted
validator IDs. Open Transit RT derives artifacts and normalizes reports; request
callers do not supply arbitrary commands, paths, argv, or output directories.

Validator success helps assess feed quality, but it is not consumer acceptance,
CAL-ITP/Caltrans compliance, agency approval, or production readiness proof by
itself.

## Monitoring Integration Boundary

Monitoring remains deployment-owned. The repo exposes readiness checks,
lightweight metrics when enabled, Operations Console views, and pilot
operations helpers.

Phase 38 does not add Prometheus/Grafana deployment assets, OpenTelemetry
SDK/exporter wiring, production SLO dashboards, or alert delivery proof.

## Consumer And Feed Workflow Boundary

Consumers and aggregators use the public GTFS and GTFS-Realtime URLs, plus
discoverability metadata such as `/public/feeds.json`. The repo has prepared
packet drafts and workflow records, but all current consumer and aggregator
targets remain `prepared`.

Do not add consumer submission APIs, runtime calls to external consumers,
portal automation, guessed submission paths, or target status changes without
retained target-originated evidence.

For detailed workflow rules, see
[Consumer Submission Workflow](evidence/consumer-submissions/submission-workflow.md)
and [California Readiness Summary](california-readiness-summary.md).

## Redaction And Evidence Boundary

Follow [Evidence Redaction Policy](evidence/redaction-policy.md) before
capturing or committing any integration material.

Never commit real device tokens, API keys, vendor credentials, private
identifiers, raw private payloads, unredacted logs, private portal
correspondence, webhook URLs, private database URLs, or private infrastructure
details.

Synthetic fixtures and dry-run output are examples and conformance aids. They
are not external evidence packets.

## What This Does Not Prove

This kit does not prove:

- certified vendor compatibility;
- real vendor AVL integration;
- production AVL reliability;
- production-grade ETA quality;
- runtime external predictor compatibility;
- consumer submission, review, acceptance, listing, or ingestion;
- CAL-ITP/Caltrans compliance;
- agency adoption, approval, or endorsement;
- agency-owned final-root proof;
- hosted SaaS availability.
