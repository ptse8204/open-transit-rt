# Open Transit RT — Post-60 Codex Master Plan

## Executive decision

Do not start a generic Phase 61. Treat the repo as post-Phase-60 and move into a post-60 productization program whose default mission is:

> Make Open Transit RT the easiest, cheapest, safest open-source path for a small transit agency, civic technologist, or operator to publish GTFS and GTFS Realtime feeds, connect existing systems through adapters, and evaluate CAL-ITP/Caltrans-style readiness without making unsupported compliance/adoption claims.

External proof remains optional and evidence-gated. The default path is product quality, external connector maturity, clean installability, release-candidate readiness, and operator usability.

## Non-negotiable claim boundaries

Do not claim any of the following unless retained public-safe evidence exists:

- CAL-ITP/Caltrans compliance
- consumer submission, review, acceptance, ingestion, listing, or display
- agency adoption, endorsement, or approval
- agency-owned final-root proof
- hosted SaaS, paid support, or SLA-backed service
- universal production readiness or production multi-tenant hosting
- certified vendor compatibility or hardware certification
- production-grade ETA quality or real-world ETA accuracy
- public launch completion

All seven consumer/aggregator targets must remain `prepared` unless target-originated evidence supports an exact target-specific status change.

## Master/sub-agent governance

Codex should run as a master agent. If native sub-agents are available, use them. Otherwise simulate the pattern with explicit passes.

Roles:

1. Master agent
   - reads repo source-of-truth files;
   - selects the next approved stage;
   - asks a planning sub-agent for a plan;
   - reviews the plan against boundaries;
   - asks an execution sub-agent to implement only the approved scope;
   - inspects repo changes directly;
   - requires patches until validation and boundaries pass;
   - commits using checkpoint format.

2. Planning sub-agent
   - produces exact files to edit, files not to edit, tests, risks, and claim boundaries.

3. Execution sub-agent
   - implements only the approved plan;
   - does not create evidence unless explicitly authorized;
   - does not move consumer statuses;
   - runs required checks;
   - reports blockers exactly.

Commit format:

```text
Post-60 -- Checkpoint 000001: realign default roadmap to external connection quality
Post-60 -- Checkpoint 000002: add external connection readiness check
Post-60 -- Checkpoint 000003: add connector plugin contract
```

Rules:

- six-digit checkpoint numbers;
- monotonic within post-60;
- no mixed-stage commits;
- no unapproved scope expansion;
- no secrets, raw private telemetry, raw vendor payloads, credentials, private portal links, or unredacted artifacts.

## Source files to read before every stage

- `AGENTS.md`
- `README.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/integration-adapter-kit.md`
- `docs/dependencies.md`
- `docs/decisions.md`
- `docs/backlog.md`
- `docs/open-questions.md`
- `docs/phase-60-final-claim-review-and-public-closeout.md`
- `docs/handoffs/phase-60.md`
- `docs/evidence/redaction-policy.md`
- `SECURITY.md`
- `Makefile`

## Stage 0 — Repo truth and validation baseline

Purpose: establish the exact repo state before changing anything.

Tasks:

- Verify phases 0 through 60 are closed for documented scopes.
- Verify Phase 60 closeout and unsupported-claim language are intact.
- Verify all seven consumer targets remain `prepared`.
- Verify no new files have appeared under `docs/evidence` except explicitly approved historical evidence.
- Verify current README still describes the repo as a self-hosted backend, not SaaS/compliance/adoption proof.

Required checks:

```bash
make audit-final-claim-review
make validate
make test
make smoke
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path
expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config
```

Exit gate:

- If checks pass, proceed to Stage 1.
- If checks fail, patch only truth-state, validation, or stale-doc inconsistencies before product work.

## Stage 1 — Post-60 product direction realignment

Purpose: make the repo’s next direction unambiguous: product quality and external-connection maturity first; optional evidence tracks only when authorized.

Files likely to edit:

- `docs/handoffs/latest.md`
- `docs/current-status.md`
- `docs/roadmap-to-calitp-compliance-and-gap-closure.md`
- `docs/integration-adapter-kit.md`
- `README.md` only if bounded wording needs a light update
- add `docs/post-60-product-roadmap.md`

Implementation requirements:

- State that no Phase 61 is being started.
- Describe missing agency/final-root/consumer/vendor/ETA evidence as limiting claims, not blocking product development.
- Make the default roadmap:
  - external adapter conformance;
  - release-candidate hardening;
  - clean install testing;
  - integration contract tests;
  - observability/export surfaces;
  - security/redaction hardening;
  - self-hosted integrator docs.
- Preserve optional evidence tracks for final root, real pilot, consumer submission, vendor data, and observed-arrival/ETA validation.

Do not:

- write to `docs/evidence`;
- move consumer statuses;
- claim compliance, adoption, production readiness, vendor compatibility, public launch, or ETA quality.

Commit:

```text
Post-60 -- Checkpoint 000001: realign default roadmap to external connection quality
```

## Stage 2 — Release-candidate readiness and stabilization

Purpose: make the repo trustworthy to install, test, and evaluate from a fresh clone.

Deliverables:

- `docs/release-candidate-readiness.md`
- `scripts/release-candidate-check.sh`
- `make release-candidate-check`

Checks covered by the new target:

- Go tests and validation;
- static shell/script checks already used by repo validation;
- Docker Compose config;
- local app startup path or documented blocker;
- GTFS import path or fixture path;
- five public feed paths;
- telemetry simulator path;
- deployment doctor;
- validator health;
- operations reliability;
- operations notifications;
- final claim audit;
- release package audit if artifacts are generated.

Rules:

- Release artifacts must remain under `.cache` unless maintainer explicitly approves publishing.
- No GitHub release, tag, package, or public announcement without explicit approval.

Commit:

```text
Post-60 -- Checkpoint 000002: add release-candidate readiness check
```

## Stage 3 — External connector plugin contract

Purpose: turn external integrations into a safe, cheap, open-source plugin/sidecar system without making the core depend on vendors.

Architecture decision:

Use a sidecar + manifest + conformance-test model, not arbitrary dynamic code loading. This keeps deployment cheap, auditable, and safer for small agencies.

Connector classes:

1. Telemetry source connectors
   - Transform AVL/device/vendor/CSV/HTTP payloads into authenticated `POST /v1/telemetry` calls.
   - No vendor assumptions inside matching, Vehicle Positions, Trip Updates, or token lifecycle.

2. Prediction connectors
   - Implement optional prediction sidecars behind `internal/prediction.Adapter`.
   - Deterministic prediction remains default.
   - Vehicle Positions must remain independent of external predictor health.

3. Validator connectors
   - Use server-side allowlisted validator IDs only.
   - No browser-supplied commands, paths, URLs, argv, or timeouts.

4. Monitoring/export connectors
   - Export sanitized local/reference summaries to deployment-owned monitoring.
   - No notification sending by default.

5. Consumer/discovery connectors
   - Document stable public GTFS/GTFS-RT URLs and packet generation.
   - Do not automate submissions or change statuses without target-originated evidence.

Deliverables:

- `docs/connectors/plugin-contract.md`
- `docs/external-connection-readiness.md`
- `internal/connectors/manifest.go`
- `internal/connectors/manifest_test.go`
- `schemas/connector-manifest.schema.json` if schema validation is already consistent with repo style
- `testdata/connectors/valid/*.json`
- `testdata/connectors/invalid/*.json`
- `scripts/external-connection-check.sh`
- `make external-connection-check`

Manifest fields:

- `schema_version`
- `connector_id`
- `connector_type`
- `display_name`
- `description`
- `mode`: `sidecar`, `http`, `csv_replay`, `static_example`, `deployment_owned`
- `input_contract`
- `output_contract`
- `auth_requirements`
- `redaction_policy`
- `failure_behavior`
- `claim_boundary`
- `example_command`
- `docs_url`

Validation rules:

- reject unknown connector types;
- reject missing redaction policy;
- reject claim words like “certified,” “accepted,” “compliant,” “production-ready,” “vendor-compatible” unless explicitly marked as forbidden wording examples;
- reject credentials, tokens, private keys, raw payloads, private URLs, and local absolute private paths in committed manifests;
- ensure sidecar connector docs say adapters are optional and fail closed.

Commit:

```text
Post-60 -- Checkpoint 000003: add connector plugin contract
```

## Stage 4 — Adapter conformance suite

Purpose: let agencies and integrators test connectors before touching real vendor data.

Deliverables:

- `cmd/adapter-conformance`
- `docs/tutorials/external-adapter-conformance.md`
- `testdata/adapter-conformance/telemetry/`
- `testdata/adapter-conformance/prediction/`
- `testdata/adapter-conformance/validator/`
- `testdata/adapter-conformance/monitoring/`

Test categories:

Telemetry:

- valid on-route event;
- malformed payload;
- missing device ID;
- unknown device;
- stale timestamp;
- future timestamp;
- wrong agency;
- duplicate event;
- out-of-order event;
- low-quality GPS;
- redaction of raw payload.

Prediction:

- deterministic fallback;
- external timeout;
- malformed predictor output;
- wrong agency/feed;
- low confidence;
- stale prediction;
- shadow mode comparison;
- no production ETA claim.

Validator:

- allowlisted validator IDs;
- static vs realtime mapping;
- command/path injection rejection;
- no raw stdout/stderr exposure.

Monitoring:

- sanitized summary export;
- no webhook/email sending by default;
- destination presence booleans only;
- redaction scan.

Commit:

```text
Post-60 -- Checkpoint 000004: add adapter conformance suite
```

## Stage 5 — Generic connector examples

Purpose: provide cheap reusable examples without vendor lock-in.

Deliverables:

- `examples/connectors/telemetry-http-poller/`
- `examples/connectors/telemetry-csv-replay/`
- `examples/connectors/predictor-sidecar-stub/`
- `examples/connectors/monitoring-export/`
- `docs/connectors/README.md`

Rules:

- synthetic data only;
- no real credentials;
- no real vendor payloads;
- no named vendor compatibility claim;
- examples must run locally or document exact blockers;
- examples must include redaction and claim-boundary notes.

Commit:

```text
Post-60 -- Checkpoint 000005: add generic connector examples
```

## Stage 6 — Agency launchpad UX hardening

Purpose: reduce adoption friction for non-expert small agencies.

Deliverables:

- `docs/tutorials/agency-launchpad.md`
- private authenticated admin page, if existing app structure supports it: `/admin/operations/launchpad`
- JSON endpoint: `/admin/operations/launchpad.json`

Workflow rows:

1. Start local/reference deployment.
2. Import or author GTFS.
3. Set agency metadata, license, and technical contact.
4. Verify five public feed paths.
5. Connect synthetic telemetry.
6. Run validator health.
7. Review CAL-ITP-style readiness gaps.
8. Try connector conformance.
9. Generate support bundle.
10. Decide whether this remains local/reference, internal RC, or authorized evidence path.

Rules:

- private authenticated route only;
- no public unauthenticated route;
- no evidence creation;
- no consumer status change;
- no approval/compliance flags.

Commit:

```text
Post-60 -- Checkpoint 000006: add agency launchpad workflow
```

## Stage 7 — CAL-ITP/Caltrans readiness hardening

Purpose: align the readiness workflow with current Caltrans guidance while staying claim-bounded.

Deliverables:

- `docs/caltrans-readiness-gap-report.md`
- `scripts/caltrans-readiness-check.sh`
- `make caltrans-readiness-check`

Rows to check and report:

- stable public GTFS URL;
- stable public GTFS-RT Trip Updates URL;
- stable public GTFS-RT Vehicle Positions URL;
- stable public GTFS-RT Service Alerts URL;
- public fetchability without unreasonable barriers;
- HTTPS;
- open license metadata;
- technical contact metadata;
- MobilityData static validator status;
- MobilityData realtime validator status;
- feed freshness and timestamp health;
- trip ID consistency signals;
- consumer packet preparedness;
- explicit unsupported-claim reminders.

Output:

- `.cache/caltrans-readiness/<timestamp>/summary.json`
- `.cache/caltrans-readiness/<timestamp>/summary.md`
- no writes to `docs/evidence` unless an explicitly approved evidence track promotes a redacted artifact.

Commit:

```text
Post-60 -- Checkpoint 000007: add Caltrans readiness gap check
```

## Stage 8 — Optional evidence-track router

Purpose: route real-world opportunities without confusing them with default product work.

Deliverable:

- `docs/evidence/evidence-track-router.md`

Decision table:

| Available real input | Track | First action |
| --- | --- | --- |
| agency/operator-approved final public root | Track C | final-root intake |
| agency/operator willing to evaluate | Track D | pilot authorization package |
| official consumer-submission authorization | Track E | target selection and official-path verification |
| real AVL/vendor data authorization | Track F | private adapter trial intake |
| observed arrival/departure data | Track G | ETA/backtesting intake |
| none | Track H | release/installability hardening |

No intake, no evidence phase.

Commit:

```text
Post-60 -- Checkpoint 000008: add evidence-track router
```

## Autonomous master prompt for Codex

Paste the following to Codex to run the post-60 program.

```text
You are the Codex master agent for ptse8204/open-transit-rt.

Mission:
Make Open Transit RT the easiest, cheapest, safest open-source self-hosted backend for small agencies and civic technologists to publish GTFS and GTFS Realtime, connect external systems through adapter/sidecar plugins, and evaluate CAL-ITP/Caltrans-style readiness without making unsupported claims.

Historical note: this prompt was superseded by the maintainer-authorized Phase 61+ agency-first connector platform roadmap in docs/roadmaps/agency-first-connector-platform/. Phases 0 through 60 remain closed for their documented scopes.

Governance:
- Act as a master agent.
- If sub-agents are available, spawn:
  1. planning sub-agent for each stage;
  2. execution sub-agent after master approval;
  3. review sub-agent or master review after execution.
- If sub-agents are not available, simulate with explicit planning, review, execution, and post-execution review passes.
- Do not execute a stage until its plan is complete and claim-safe.
- After execution, inspect repo files directly and patch until checks pass.

Read first before every stage:
- AGENTS.md
- README.md
- docs/current-status.md
- docs/handoffs/latest.md
- docs/roadmap-to-calitp-compliance-and-gap-closure.md
- docs/integration-adapter-kit.md
- docs/dependencies.md
- docs/decisions.md
- docs/backlog.md
- docs/open-questions.md
- docs/phase-60-final-claim-review-and-public-closeout.md
- docs/handoffs/phase-60.md
- docs/evidence/redaction-policy.md
- SECURITY.md
- Makefile

Global boundaries:
- Do not write to docs/evidence unless explicitly approved for an evidence track.
- Do not edit docs/evidence/consumer-submissions/status.json except for an approved target-originated evidence status transition.
- Do not automate consumer submissions or contact external portals.
- Do not claim CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption/approval, final-root proof, hosted SaaS, paid support/SLA, production readiness, production multi-tenant hosting, vendor compatibility, hardware certification, public launch, production AVL reliability, or production-grade ETA quality.
- Keep all external integrations adapter-bound, optional, documented, tested, redacted, and fail-closed.
- Keep Vehicle Positions independent of external predictor availability.
- Keep deterministic prediction as default.
- Treat validator success as supporting signal only, not compliance, consumer acceptance, or correctness proof.

Required baseline checks after every checkpoint:
make audit-final-claim-review
make validate
make test
make smoke
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
python3 - <<'PY'
import json
from pathlib import Path
expected = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
data = json.loads(Path("docs/evidence/consumer-submissions/status.json").read_text())
records = data.get("targets", [])
seen = {row["target"]: row.get("status") for row in records}
assert list(seen) == expected, seen
assert all(seen[name] == "prepared" for name in expected), seen
PY
git diff --exit-code -- docs/evidence/consumer-submissions/status.json
docker compose -f deploy/docker-compose.yml config

Proceed in order:

Stage 0 — Repo truth and validation baseline.
- Verify phases 0 through 60 closed, Phase 60 claim boundaries, consumer prepared-only tracker, and evidence preservation.
- Patch only stale truth-state/validation issues if found.

Stage 1 — Post-60 product direction realignment.
- Add/update docs so default next work is external-connection quality, not real agency/public pilot.
- Add docs/post-60-product-roadmap.md if missing.
- Commit: Post-60 -- Checkpoint 000001: realign default roadmap to external connection quality

Stage 2 — Release-candidate readiness and stabilization.
- Add docs/release-candidate-readiness.md.
- Add scripts/release-candidate-check.sh and make release-candidate-check.
- Do not tag or publish without explicit maintainer approval.
- Commit: Post-60 -- Checkpoint 000002: add release-candidate readiness check

Stage 3 — External connector plugin contract.
- Add docs/connectors/plugin-contract.md.
- Add docs/external-connection-readiness.md.
- Add internal/connectors manifest validation and tests.
- Add connector manifest testdata.
- Add scripts/external-connection-check.sh and make external-connection-check.
- Use a sidecar/manifest/conformance model, not arbitrary dynamic plugin loading.
- Commit: Post-60 -- Checkpoint 000003: add connector plugin contract

Stage 4 — Adapter conformance suite.
- Add cmd/adapter-conformance.
- Add testdata/adapter-conformance for telemetry, prediction, validator, and monitoring cases.
- Add docs/tutorials/external-adapter-conformance.md.
- Test malformed/stale/future/wrong-agency/unknown-device/low-quality/duplicate/out-of-order telemetry; predictor timeout/malformed/stale/wrong-agency; validator allowlist; monitoring redaction/no-send.
- Commit: Post-60 -- Checkpoint 000004: add adapter conformance suite

Stage 5 — Generic connector examples.
- Add examples/connectors/telemetry-http-poller.
- Add examples/connectors/telemetry-csv-replay.
- Add examples/connectors/predictor-sidecar-stub.
- Add examples/connectors/monitoring-export.
- Use synthetic data only. No real credentials, no vendor payloads, no vendor compatibility claims.
- Commit: Post-60 -- Checkpoint 000005: add generic connector examples

Stage 6 — Agency launchpad UX hardening.
- Add docs/tutorials/agency-launchpad.md.
- If feasible within existing app patterns, add private authenticated /admin/operations/launchpad and /admin/operations/launchpad.json.
- Show setup, GTFS, metadata, five feeds, telemetry, validators, readiness, connector conformance, support bundle, and next-decision steps.
- No evidence creation or approval/compliance flags.
- Commit: Post-60 -- Checkpoint 000006: add agency launchpad workflow

Stage 7 — CAL-ITP/Caltrans readiness hardening.
- Add docs/caltrans-readiness-gap-report.md.
- Add scripts/caltrans-readiness-check.sh and make caltrans-readiness-check.
- Report stable URLs, all three GTFS-RT feed types, public fetchability, HTTPS, open license, contact, validators, freshness, trip ID consistency signals, consumer packet preparedness, and unsupported-claim boundaries.
- Output to .cache only.
- Commit: Post-60 -- Checkpoint 000007: add Caltrans readiness gap check

Stage 8 — Optional evidence-track router.
- Add docs/evidence/evidence-track-router.md.
- Document Tracks C/D/E/F/G/H and the rule: no intake, no evidence phase.
- Do not run evidence tools unless real authorization and public-safe artifact retention rules exist.
- Commit: Post-60 -- Checkpoint 000008: add evidence-track router

Stop condition:
Stop after Stage 8 if all checks pass. Summarize implemented checkpoints, blocked commands, remaining optional evidence tracks, and the next maintainer decision: release candidate publish vs more installability testing vs authorized evidence path.
```
