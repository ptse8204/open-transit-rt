# Open Transit RT `v0.1.0-rc.1` Draft Release Notes

Draft date: 2026-05-12

This is a local pre-tag draft created in Phase 72 Checkpoint 000006 and updated
during Checkpoint 000007 closeout. It is not a GitHub release, tag, package,
hosted deployment, retained evidence packet, or release-readiness claim.

## Source

- Git tag: None; not tagged.
- Commit SHA at draft time: `b0e354d29b14ba5c5ab0afab3269ec5eb3405fa7`
- Dirty/clean state at draft time: `b0e354d-dirty`
- Release notes link: local draft at
  `docs/release-notes-v0.1.0-rc.1-draft.md`
- Artifact checksums: None.
- Release package: None.
- SBOM/provenance: None.

## Summary

`v0.1.0-rc.1` is the next maintainer evaluation milestone for the self-hosted
small-agency product path. The draft is based on Phase 72 local hardening
results through CP000007: clean-checkout evaluator rehearsal,
release-candidate diagnostic hardening, local app and Operations Console
walkthrough checks, five local public feed fetches, local synthetic
connector/adaptor conformance checks, final local validation, and closeout. The
candidate still has unresolved blockers before it can be treated as a complete
release-candidate review.

## User-Facing Changes

- Browser-first agency operations are represented by the private
  `Agency Operations Cockpit / Start Here` path and required Operations Console
  routes.
- Local agency-app startup and five public feed path fetches have recorded
  local diagnostic results.
- Connector/adaptor examples and conformance checks have recorded local
  synthetic pass results.
- Release-candidate diagnostics now record pinned validator tooling blockers
  explicitly instead of allowing ambient stub validator mode to mask them.

## Install And Upgrade Notes

- Clean install from source tag: None; no tag exists for this draft.
- Local app verification: CP000004 recorded `make agency-app-up` and
  `make agency-app-down` passing in the primary checkout.
- Release-candidate diagnostic: CP000002 evaluator `make
  release-candidate-check` exited `2` because pinned static GTFS validator
  tooling was missing. CP000004 primary `RUN_LOCAL_APP=true` diagnostic exited
  `0` with `overall_status=needs_review` because release-package auditing was
  intentionally not enabled.
- Local release package: None.
- Local Docker image build: None for this draft.
- Published production Docker image: None.

## Migrations

- None.

Before any future upgrade, operators must back up the database and run:

```bash
make migrate-status
make migrate-up
make migrate-status
```

## Operations Changes

- None; CP000004 and CP000005 recorded local diagnostics only.

## Security Notes

- None.

## Dependency Changes

- None.

## Evidence Or Claim Changes

- None.

This draft does not create retained evidence, change consumer statuses, contact
external services, publish artifacts, or add claims of CAL-ITP/Caltrans
compliance, consumer submission/review/acceptance/ingestion/listing/display,
agency adoption/approval, final-root readiness, hosted SaaS availability, paid
support, SLA/uptime, production readiness, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality,
validator-clean feeds, or release readiness.

## Known Blockers Matrix

| Area | Status | Exact blocker or note | Next owner |
| --- | --- | --- | --- |
| Primary source state | needs_review | Primary checkout is dirty at `b0e354d-dirty`; do not describe it as clean. | Future release-cut cleanup |
| Clean evaluator release-candidate diagnostic | blocker | CP000002 evaluator `make release-candidate-check` exited `2`: `missing pinned tooling: static GTFS validator not installed at /private/var/folders/_g/bvzl9cms7cx1d0wdpc981n9w0000gn/T/tmp.BA6ULYuxdo/open-transit-rt-evaluator/.cache/validators/gtfs-validator-7.1.0-cli.jar; run make validators-install`. | Future release-cut cleanup or later authorized tooling step |
| Release package | not_checked | No local release package was generated or audited in CP000006 or CP000007. | Future release-cut cleanup if approved |
| Artifact checksums | not_checked | No package artifact exists for this draft. | Future release-cut cleanup if approved |
| SBOM/provenance | not_checked | No package artifact exists for this draft. | Future release-cut cleanup if approved |
| Browser automation | needs_review | CP000004 in-app Browser review with safe bearer-header support was unavailable; terminal authenticated GET checks were used as a safe substitute. | Future UI review if a fresh browser signal is required |
| Final validation | needs_review | CP000007 final focused checks passed where run; `RUN_LOCAL_APP=true make release-candidate-check` exited `0` with `overall_status=needs_review`, `git_clean=needs_review`, and `release_package_audit=not_checked`. | Future release-cut cleanup |
| External evidence and consumer state | unchanged | No retained evidence collection, external contact, consumer submission, or consumer status movement is authorized or performed. | Authorization-gated only |

## Known Limitations

- The draft is not tied to a clean tagged source state.
- The static GTFS pinned validator tooling blocker remains unresolved in the
  clean evaluator release-candidate diagnostic.
- The optional release package audit did not run in CP000006 or CP000007.
- No artifact checksum, SBOM, provenance, image, GitHub release, or package
  exists for this draft.
- CP000005 connector/adaptor results are local synthetic checks only and do not
  prove real vendor/device compatibility.

## Checks

Recorded Phase 72 checks through CP000007:

| Check | Result | Scope |
| --- | --- | --- |
| `make check` | passed in CP000002 evaluator and later primary validation | Lightweight local repo check |
| `git diff --check` | passed in CP000007 primary checkout | Whitespace guard |
| `scripts/release-candidate-check.sh --dry-run` | passed in CP000002 evaluator | Private diagnostic rehearsal |
| `make release-candidate-check` | blocker in CP000002 evaluator | Missing pinned static GTFS validator tooling |
| `scripts/bootstrap-dev.sh --check` | passed in CP000002 evaluator | Local setup preflight |
| `make validate` | passed in CP000007 primary checkout | Local validation smoke |
| `make test` | passed in CP000002 evaluator, CP000003 validation, and CP000007 primary checkout | Local tests |
| `docker compose -f deploy/docker-compose.yml config` | passed in CP000002 evaluator and CP000007 primary checkout | Compose syntax/config |
| `make agency-app-up` | passed in CP000004 primary checkout | Local app startup diagnostic |
| Private Operations Console route checks | passed in CP000004 primary checkout | Authenticated local `200` route checks |
| `/admin/operations` without admin token | returned local `401` in CP000004 | Private admin boundary check |
| Five local public feed fetches | passed in CP000004 primary checkout | Anonymous local `200` fetches |
| `RUN_LOCAL_APP=true scripts/release-candidate-check.sh` | exited `0` in CP000004 primary checkout with `overall_status=needs_review` | Private local diagnostic |
| `RUN_LOCAL_APP=true make release-candidate-check` | exited `0` in CP000007 primary checkout with `overall_status=needs_review` | Private local diagnostic; dirty checkout and release package audit still need review |
| `make external-connection-check` | passed in CP000005 primary checkout | Local synthetic connector boundary check |
| `make adapter-conformance` | passed in CP000005 primary checkout | Local synthetic adapter conformance |
| `make test-connector-examples` | passed in CP000005 primary checkout | Local synthetic connector examples |
| `make audit-product-acceptance` | passed in CP000007 primary checkout | Product acceptance guard |
| `make audit-final-claim-review` | passed in CP000007 primary checkout | Claim-boundary guard |
| Consumer tracker prepared-only review | passed | Seven targets remain exactly `prepared` |
| Protected-path status review | passed | Protected evidence paths, migrations, `go.mod`, and `go.sum` unchanged |

Blocked checks:

- `make release-candidate-check` in the CP000002 evaluator remains blocked by
  missing pinned static GTFS validator tooling.
- Release package, checksum, SBOM/provenance, GitHub release, tag, and published
  image were not created or audited in CP000007.

Release-candidate summary:

- Output directory: `.cache/release-candidate-check/20260512T062130Z` for the
  CP000002 dry-run diagnostic; `.cache/phase-72-cp000004/release-candidate-check`
  for the CP000004 local-app diagnostic; `.cache/release-candidate-check/20260512T105655Z`
  for the CP000007 local-app diagnostic.
- Overall status: `blocker` for the CP000002 evaluator full diagnostic;
  `needs_review` for the CP000004 and CP000007 local-app diagnostics.
- Package audit status: `not_checked`.
- Claim-boundary result: no retained evidence, consumer status change, image
  push, hosted-service claim, production-readiness claim, release-readiness
  claim, or compliance claim was added.
