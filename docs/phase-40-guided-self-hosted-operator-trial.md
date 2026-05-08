# Phase 40 — Guided Self-Hosted Operator Trial

## Status

Complete for the docs/navigation guided operator trial scope.

Phase 40 ties together the Phase 36 reference deployment docs, Phase 37
reusable agency onboarding flow, Phase 38 integration adapter kit, and Phase 39
readiness workflow into one local/reference evaluation path.

This phase was docs/navigation only. It does not create external evidence,
final-root evidence, consumer submissions, consumer acceptance, agency approval
or adoption, CAL-ITP/Caltrans compliance, hosted SaaS availability, production
readiness, vendor compatibility, or production-grade ETA quality evidence.

## Related Context

Read these before running or extending the trial:

- [Phase 39 Handoff](handoffs/phase-39.md)
- [Phase 40 Handoff](handoffs/phase-40.md)
- [Self-Hosted Operator Trial](tutorials/self-hosted-operator-trial.md)
- [OCI/OCL Reference Deployment](deployment/oci-reference-deployment.md)
- [Reusable Agency Onboarding](tutorials/reusable-agency-onboarding.md)
- [Integration Adapter Kit](integration-adapter-kit.md)
- [CAL-ITP Readiness Checklist](tutorials/calitp-readiness-checklist.md)

## Goal

Make it easy for an operator to run one bounded self-hosted trial:

1. prepare or start a local/reference deployment;
2. import GTFS through `make agency-pilot-up`;
3. verify the five public feed paths;
4. review `/admin/operations/readiness`;
5. run, skip, or blocker-document validators;
6. run the synthetic AVL dry-run adapter;
7. identify next actions without creating evidence or stronger claims.

## What Changed

- Added a guided operator tutorial at
  `docs/tutorials/self-hosted-operator-trial.md`.
- Linked the trial from onboarding, deployment, adapter, readiness, roadmap,
  status, docs index, tutorial index, and handoff navigation.
- Updated `make validate` to assert the Phase 40 docs and handoff files exist.

## Boundaries

Operator output stays local/private by default. Logs, screenshots, validator
output, copied summaries, and support notes must not be committed unless a
later evidence phase explicitly reviews, redacts, and retains them.

The no-external-network fixture path uses `demo-agency` and
`testdata/gtfs/valid-small` for local demo evaluation only. It is not real
agency proof.

All seven consumer and aggregator targets remain `prepared` only. Consumer
statuses may move only when retained, redacted, target-originated evidence
supports a target-specific transition.

## Checks

Required Phase 40 close checks:

```bash
make validate
make test
git diff --check
python3 -m json.tool docs/evidence/consumer-submissions/status.json >/dev/null
```

Also verify all seven consumer targets remain `prepared`.
