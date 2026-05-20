# Deployment Reference

This guide is documentation for reproducing a self-hosted deployment pattern. It is not evidence that a deployment was run, accepted, compliant, agency-approved, or production-ready.

This folder documents the self-hosted OCI/OCL-style reference path for Open
Transit RT. It is an operator reference for reusing the existing pilot-server
pattern on a deployment-owned host. It is not hosted SaaS, not an evidence
packet, and not final-root proof.

## Documents

- [OCI/OCL Reference Deployment](oci-reference-deployment.md) - end-to-end
  operator guide for the reference path.
- [Reference Environment Example](oci-reference-env.example) - placeholder-only
  environment file structure.
- [Reference Smoke Checklist](oci-reference-smoke-checklist.md) - repeatable
  verification checklist for an operator-run deployment.
- [Reference Deployment Doctor](reference-deployment-doctor.md) - read-only
  private diagnostics for env, small-host resources, services, proxy exposure,
  validators, DB/PostGIS/pool guidance, backups, restore-drill readiness, and
  upgrade/rollback plus install/recovery stop-point posture.
- [OCI Reference Check](oci-reference-check.md) - private reference diagnostic
  helper for public five-feed fetches, optional SSH loopback health, backup/
  restore configuration presence, and all-false claim flags.
- [Off-Host Public Feed Validation](off-host-validation.md) - fetch public
  artifacts and run validators from an operator machine when the server is too
  small for validator workloads.
- [OCI Reference Diagnostic Runs](oci-reference-diagnostic-runs.md) -
  public-safe summaries of completed reference diagnostics. These are not
  retained evidence packets.
- [Self-Hosted Operator Trial](../tutorials/self-hosted-operator-trial.md) -
  guided local/reference evaluation path tying deployment prep, GTFS
  onboarding, readiness review, validators, and the synthetic AVL dry-run
  together without creating evidence.
- [Operator Smoke And Support Bundle](../tutorials/operator-smoke-and-support-bundle.md) -
  repeatable smoke checks and redaction-safe diagnostics for local/reference
  operators.

Phase 56 adds optional path-routed public feed URLs under
`/public/agencies/{agency_id}/...` plus `make multi-agency-hosting` for private
route/proxy diagnostics. This is repository-boundary hardening only; it is not
a hosted SaaS offer, production multi-tenant certification, SLA, compliance
proof, agency adoption proof, or consumer acceptance proof.

Review `docs/evidence/redaction-policy.md` before turning any operator output
from this path into public evidence.

For adoption-first use, start routine review in the private browser UI:

```text
/admin/operations
/admin/operations/feed-health
/admin/operations/maintenance
```

Use `make oci-reference-check` and `make validate-public-feeds` as private
technical-helper diagnostics, not as evidence collection.
