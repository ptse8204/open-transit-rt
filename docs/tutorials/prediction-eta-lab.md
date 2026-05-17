# Prediction And ETA Lab

Use the private Prediction and ETA Lab when an operator needs to understand why
Trip Updates or ETA-like output is missing, partial, deterministic, shadowed, or
withheld.

This tutorial is for local/self-hosted operator review. It does not collect
retained evidence, contact external predictors, move consumer statuses, prove
real-world ETA accuracy, prove production-grade ETA quality, or make a
compliance, consumer, vendor, hardware, SLA, hosted-service, public-launch, or
release-readiness claim.

## Open The Lab

Start the local app package, then open the authenticated Operations Console:

```bash
make agency-app-up
```

In the browser, go to:

```text
/admin/operations/prediction-lab
```

For a private JSON summary, use:

```text
/admin/operations/prediction-lab.json
```

Both routes are private, authenticated, read-only, and agency-scoped. They do
not execute predictors, run commands, upload data, create evidence, or contact
sidecars.

## What To Review

Start at **Trip Updates Decision**. It summarizes the latest private Trip
Updates diagnostics: adapter name, status, active feed version, eligible
candidates, emitted updates, and withheld counts.

Use **Safe Fallback** and **Deterministic Predictor Diagnostics** to confirm that
Vehicle Positions can remain independent while Trip Updates stay withheld,
empty, or partial.

Use **Why ETAs Are Missing** to identify the first likely owner:

- telemetry freshness;
- assignment confidence;
- active schedule/feed state;
- future-stop availability;
- alert lifecycle for canceled trips.

Use **External Predictor Shadow Review** only as limited shadow/fail-closed
diagnostics. The page does not start a sidecar, test a URL, reveal credentials,
or claim named predictor support.

Use **Backtest Summary** to review safe aggregate outputs under:

```text
.cache/realtime-quality-backtest/<timestamp>
```

The browser reads only exact aggregate result files:

```text
summary.json
summary.md
metrics.json
metrics.md
manifest.json
```

It rejects unsafe root shapes, symlinked files, unexpected files, schema
mismatches, raw-row persistence flags, evidence-like paths, and forbidden claim
flags.

The backtest summary also exposes a synthetic conformance signal when the local
fixture output includes after-midnight service, frequency/headway service,
service-calendar start instances, unknown/ambiguous withholding, and
shadow/fail-closed predictor handling. This is a checklist for local review,
not evidence of real-world ETA quality.

## Local Commands

Run checks from an operator shell, not from the browser:

```bash
go test ./internal/prediction -run Deterministic
make realtime-quality
make realtime-quality-backtest
```

These commands are synthetic/local diagnostics. Passing them does not prove
consumer display, compliance, real-world ETA accuracy, production-grade ETA
quality, vendor compatibility, hardware certification, SLA coverage, hosted
service readiness, or release readiness.

## Conservative Handling

If telemetry is stale, first check device power, clock sync, network delay, and
reporting cadence.

If assignment is unknown, review service day, route/trip hints, block
continuity, after-midnight service, and operator override evidence.

If assignment is ambiguous, prefer withheld Trip Updates over a fabricated trip
descriptor or future stop sequence.

If confidence is low or no future stops are safe, review GTFS Workbench, active
feed version, current stop sequence, and withheld reasons before changing
thresholds.

The safe default is: unknown is better than false certainty.

## Future Proof Gates

The lab lists future proof gates so operators know what is missing. Each gate
requires separate written authorization before collecting, retaining, or
publishing evidence:

- real observed arrival/departure comparison;
- operating-day and route coverage review;
- external predictor or device/vendor proof;
- consumer, release, and public claim gates.

Do not use this page as evidence for final-root readiness, consumer acceptance,
compliance, public launch, hosted SaaS, vendor compatibility, hardware
certification, SLA/uptime, production readiness, or production-grade ETA
quality.
