# Phase 29A — External Predictor Adapter Evaluation

## Status

Complete for adapter contract documentation, candidate-only TheTransitClock feasibility review, and test-only mock adapter contract checks.

## Purpose

Evaluate whether Open Transit RT can support an external prediction engine, such as TheTransitClock or another standards-based predictor, through a clean adapter boundary without making the external predictor the default source of truth.

Phase 29 expanded the synthetic replay baseline. Phase 29A uses that baseline to define and test the integration contract for external prediction adapters before any production runtime dependency is introduced.

This phase is about **adapter evaluation and contract proof**, not broad runtime integration or production-grade ETA claims.

## Why This Phase Exists

Open Transit RT's original direction included the ability to connect with outside transit tools instead of locking agencies into one internal predictor. External predictors can be valuable, but they must not weaken the current safety posture:

- unknown is better than false certainty;
- withheld/degraded reasons must stay visible;
- deterministic fallback must remain available;
- no external predictor can be treated as accepted, certified, or production-grade without evidence.

## Scope

1. External predictor adapter contract.
2. Candidate predictor feasibility review.
3. TheTransitClock evaluation as a named candidate, if still appropriate.
4. Mock adapter and contract tests.
5. Replay comparison against Phase 29 fixtures.
6. Failure/fallback behavior.
7. Configuration and dependency boundary documentation.
8. Security, licensing, and operations review.
9. Handoff recommendation for whether to implement a runtime adapter later.

## Implemented Scope

Phase 29A added a contract-level validation layer for Trip Updates adapter output and test-only mock external adapter coverage. It did not add runtime external predictor wiring, service clients, environment variables, network calls, subprocess calls, Java/Maven/Tomcat invocation, or TheTransitClock integration.

Vehicle Positions generation remains independent of external predictor availability. External predictor evaluation affects only Trip Updates adapter evaluation paths; telemetry ingest, assignment persistence, and Vehicle Positions publication continue without consulting an external predictor.

## External Adapter Contract

External predictor adapters must implement `internal/prediction.Adapter` and are evaluated through the existing Trip Updates builder. They must accept:

- `agency_id` from the Trip Updates service configuration;
- active `feed_version_id` and active published GTFS feed metadata;
- latest vehicle telemetry and assignment context for the agency;
- telemetry freshness context through observation timestamps and current snapshot time;
- schedule/static GTFS context available through the active feed version and schedule repository;
- manual-override effects as represented in current assignments;
- canceled-trip overrides and unsupported disruption state only through documented prediction-operation inputs or diagnostics, not by mutating assignment state;
- canonical Vehicle Positions feed URL or equivalent feed data when a later runtime adapter is explicitly approved.

External predictor outputs must return:

- prediction status and reason;
- `agency_id` and `feed_version_id` scope when the adapter can echo or assert them;
- Trip Update trip descriptor fields: `trip_id`, `route_id`, `start_date`, `start_time`, and schedule relationship;
- ordered or orderable `stop_time_update` entries;
- confidence or an equivalent quality signal when the adapter is external;
- explicit withheld/degraded reason for unsafe or unsupported cases;
- diagnostics payload suitable for `feed_health_snapshot` persistence;
- timeout, error, malformed-response, or unavailable status when predictions cannot be trusted.

Required behavior:

- adapter output must never bypass existing Trip Updates normalization;
- adapter diagnostics must never bypass diagnostics persistence;
- stale, ambiguous, unknown, degraded, wrong-agency, wrong-feed, and unsupported cases must remain visible;
- active manual override authority remains upstream in assignment state and must not be weakened by external predictions;
- deterministic prediction remains the default and must remain available as fallback;
- wrong agency or wrong feed-version output must be rejected or withheld, never silently published.

## Adapter Output Validation

The Trip Updates builder now rejects unsafe adapter output before protobuf serialization. Rejected output produces a valid empty or partial Trip Updates feed with diagnostics reason `adapter_output_rejected` or `partial_predictions`, and `withheld_by_reason` records the rejection category.

Validated rejection categories include:

- trip not present in the active feed;
- adapter-declared or candidate trip scoped to the wrong agency;
- adapter-declared or candidate trip scoped to the wrong feed version;
- impossible or missing stop sequence;
- stale prediction timestamp;
- unsupported added-trip prediction;
- low confidence;
- missing confidence when the adapter declares or is identified as external.

This validation does not change the public GTFS-RT protobuf contract. It only controls which internal adapter results are allowed to reach the existing Trip Updates serializer.

## Failure And Fallback Semantics

Timeout, unavailable service, malformed response, and other adapter errors produce a valid empty Trip Updates feed with visible `adapter_error` diagnostics and persisted diagnostics records. Malformed or conflicting adapter output is rejected or withheld with visible diagnostics. Phase 29A does not add an automatic external-to-deterministic runtime fallback chain because no runtime external adapter exists in this phase; deterministic prediction remains the configured default and the safe fallback path for future approved runtime work.

Tests use only mock/test-only adapters. They do not start TheTransitClock, invoke Java, Maven, or Tomcat, make network calls, or require external services.

## Candidate Review — TheTransitClock

Review date: 2026-05-04.

Public sources reviewed:

- `https://thetransitclock.github.io/`
- `https://github.com/TheTransitClock/transitime`
- `https://raw.githubusercontent.com/TheTransitClock/transitime/develop/BUILD.md`
- `https://raw.githubusercontent.com/TheTransitClock/transitime/develop/LICENSE`

Public project information describes TheTransitClock as open-source arrival prediction software that takes a GTFS-Realtime Vehicle Positions feed as input and produces GTFS-Realtime Trip Updates. The public repository is Java-oriented, built with Maven, includes REST/API and webapp modules, and is licensed under GPL-3.0.

Feasibility conclusion: TheTransitClock remains a plausible candidate for a later process-level or network-level adapter because its stated input/output shape aligns with Open Transit RT's Vehicle Positions-first architecture. It is not suitable for Phase 29A runtime integration because that would introduce Java/service deployment, operational health checks, and GPL-3.0 license-review questions that require a later approved phase.

Public-source review is not runtime compatibility proof. The Phase 29A review and mock adapter tests do not prove better ETAs, production-grade ETA quality, real-world predictor compatibility, consumer acceptance, CAL-ITP/Caltrans compliance, hosted SaaS availability, or vendor equivalence.

## Required Work

### 1) Adapter Contract Review

Review the existing `internal/prediction.Adapter` boundary and document whether it is sufficient for an external predictor.

The review should answer:

- What inputs does the internal deterministic predictor currently consume?
- What outputs are required by Trip Updates generation?
- Which fields are required, optional, or unsafe to infer?
- How should unknown, stale, ambiguous, degraded, canceled, added-trip, short-turn, and detour cases be represented?
- What must remain adapter-independent?
- What cannot be passed to an external predictor because it is private, unstable, or unsupported?

If the existing adapter boundary is sufficient, document that. If not, propose the smallest contract refinement and test it without forcing runtime integration.

### 2) Candidate Predictor Feasibility Review

Evaluate candidate external predictors at the design level.

At minimum, include TheTransitClock as a candidate if current public project information is still relevant. Review:

- expected inputs;
- expected outputs;
- deployment shape;
- runtime dependency implications;
- configuration needs;
- startup/health behavior;
- failure modes;
- licensing implications;
- whether the predictor expects GTFS, GTFS-RT Vehicle Positions, AVL, or another input shape;
- whether it produces GTFS-RT Trip Updates or another output shape;
- how it handles disruptions, stale inputs, and low-confidence cases.

Do not vendor external source code into this repository in Phase 29A.

### 3) Mock External Predictor Adapter

Add a mock or test-only external predictor adapter if useful for contract tests.

The mock adapter should:

- be deterministic;
- use synthetic replay inputs only;
- support success and failure cases;
- support stale/ambiguous/unknown cases;
- expose how fallback works;
- not require a network service;
- not introduce a real external runtime dependency.

### 4) Replay Comparison Against Phase 29 Fixtures

Use Phase 29 replay fixtures to compare internal deterministic behavior with the external-adapter contract.

The goal is not to prove the external predictor is better. The goal is to prove:

- adapter inputs can be shaped safely;
- adapter outputs can be interpreted safely;
- fallback behavior is deterministic;
- unknown/withheld/degraded cases remain visible;
- unsupported cases do not become false ETAs.

Include cases from Phase 29 where possible:

- after-midnight service;
- exact and non-exact frequency windows;
- block continuity;
- long layover;
- sparse telemetry;
- noisy/off-shape GPS;
- stale/ambiguous hard pattern;
- cancellation/alert linkage;
- manual override before/after expiry.

### 5) Failure And Fallback Behavior

Document and test:

- predictor unavailable;
- timeout;
- invalid response;
- stale response;
- missing trip update;
- conflicting prediction;
- low-confidence prediction;
- adapter returns output for unsupported disruption;
- deterministic fallback remains available.

The default behavior should remain conservative:

- keep deterministic predictor as default unless explicitly configured otherwise;
- fail closed to unknown/withheld/degraded where appropriate;
- do not make the public feed look better by hiding uncertainty;
- do not emit unsupported ETAs merely because an external predictor supplied them.

### 6) Configuration Boundary

Document a future configuration model without enabling a production external predictor by default.

Include:

- how an operator would choose predictor adapter;
- what environment variables would be needed;
- what should remain private;
- what health/readiness checks are required;
- how to disable the external predictor quickly;
- how to return to deterministic fallback.

Do not add secrets, credentials, private predictor configuration, `TRIP_UPDATES_ADAPTER=external`, production runtime external predictor toggles, service clients, or network calls in Phase 29A.

### 7) Licensing And Dependency Review

Document licensing and dependency implications before any external predictor becomes a runtime dependency.

This is especially important for predictors that are GPL-licensed or require Java/service deployment.

The review should state:

- whether any code is vendored;
- whether the integration is process-level, network-level, or library-level;
- whether dependencies are optional or required;
- whether license compatibility needs maintainer/legal review;
- whether `docs/dependencies.md` needs an update.

### 8) Documentation

Add or update:

- `docs/phase-29a-external-predictor-adapter-evaluation.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/handoffs/phase-29a.md`
- `docs/dependencies.md` only if a dependency status is introduced
- `docs/decisions.md` only if the adapter boundary changes
- `docs/track-b-productization-roadmap.md`
- `docs/roadmap-status.md` if needed

## Acceptance Criteria

Phase 29A is complete only when:

- the external predictor adapter contract is documented;
- TheTransitClock or another external predictor candidate has a bounded feasibility review;
- mock/contract tests exist if code-level adapter evaluation is implemented;
- replay comparison against Phase 29 fixtures exists or is explicitly deferred with reason;
- fallback/failure behavior is documented and tested where applicable;
- deterministic predictor remains default;
- no real external predictor is required at runtime;
- no public feed URL, GTFS-RT contract, consumer status, or evidence claim changes are introduced;
- no production-grade ETA quality claim is introduced;
- `docs/handoffs/phase-29a.md` exists and uses the repo handoff template.

## Required Checks

```bash
make validate
make realtime-quality
make test
make smoke
make test-integration
docker compose -f deploy/docker-compose.yml config
git diff --check
```

If adding focused adapter tests, also run:

```bash
go test ./internal/prediction ./internal/realtimequality ./internal/feed/tripupdates
```

If no runtime code changes are made, document why focused Go tests were not needed.

## Explicit Non-Goals

Phase 29A does not:

- make TheTransitClock or another predictor a production runtime dependency;
- vendor external predictor code;
- require Java, Docker, or external services for normal repo tests;
- replace the deterministic predictor;
- hide unknown/withheld/degraded cases;
- claim production-grade ETA quality;
- claim real-world ETA accuracy;
- claim consumer acceptance or CAL-ITP/Caltrans compliance;
- change public feed URLs;
- change GTFS-RT protobuf contracts;
- change consumer statuses;
- add hosted SaaS or paid support claims.

## Security And Privacy Boundaries

Do not commit:

- predictor credentials;
- private AVL payloads;
- real private telemetry;
- private agency GTFS;
- DB URLs with passwords;
- tokens;
- private keys;
- `.cache` files;
- raw private operator artifacts.

Use synthetic fixtures and public-safe examples only.

## Likely Files

- `internal/prediction/`
- `internal/realtimequality/`
- `internal/feed/tripupdates/`
- `testdata/replay/`
- `docs/phase-29a-external-predictor-adapter-evaluation.md`
- `docs/handoffs/phase-29a.md`
- `docs/track-b-productization-roadmap.md`
- `docs/current-status.md`
- `docs/handoffs/latest.md`
- `docs/dependencies.md` only if dependency status changes
- `docs/decisions.md` only if adapter policy changes
