# Open Transit RT `v0.1.0-rc.1` Draft Release Notes

Draft date: 2026-05-14

This is a local pre-tag draft refreshed during Phase 89 Checkpoint 000005. It
is not a GitHub release, tag, package, hosted deployment, retained evidence
packet, release publication, or release-readiness claim.

## Source

- Git tag: None; not tagged.
- Commit SHA at draft time:
  `9ab1123911b0fc5867f4bdd5faa90b8ba32df765`
- Dirty/clean state at draft time: `9ab1123`
- Release notes link: local draft at
  `docs/release-notes-v0.1.0-rc.1-draft.md`
- Release-candidate gate results:
  `docs/phase-89-rc1-gate-results.md`
- Artifact checksums: None.
- Release package: None; package creation is blocked without separate
  maintainer authorization.
- SBOM/provenance: None.

## Summary

`v0.1.0-rc.1` remains the next maintainer evaluation milestone for the
self-hosted small-agency product path. Since the original Phase 72 draft, the
Phase 75-89 product track has added a much more complete private browser-first
control plane: app shell polish, bounded command model, progressive UI state,
agency setup, GTFS Workbench, Realtime Center, Validation Center, Connector
Workbench, Prediction & ETA Lab, Maintenance Center, agency scope/roles/audit,
feed readiness, consumer prepared-only guidance, and nontechnical training.

The Phase 89 local gate has passed clean local product checks, local route
diagnostics, and synthetic/local connector/backend diagnostics so far. The
candidate still remains `needs_review` because release package creation and
package audit are not authorized or run, no tag exists, and no release action
is authorized.

## User-Facing Changes

- Private Operations Console has a consumer-grade app shell with grouped
  navigation, status chips, route titles, contextual help, focus states, and
  mobile/high-contrast-friendly layout.
- Agency Setup, GTFS Import, GTFS Workbench, Feed Links & Health, Validation
  Center, Realtime Center, Prediction & ETA Lab, Connector Workbench,
  Maintenance Center, Access/Roles, Audit Log, Consumers, Evidence, and Help
  now present clearer private operator workflows.
- Help now includes role-based tours, first-week checklist, glossary, recovery
  guidance, quick tasks, staff handoff checklist, and a printable operator
  training guide.
- Public-facing docs and GitHub Pages remain documentation-only and point to
  the same browser-first product path.
- Prepared consumer packet records remain prepared-only and unchanged.

## Install And Upgrade Notes

- Clean install from source tag: None; no tag exists for this draft.
- Local app verification: Phase 89 CP000003 ran
  `RUN_LOCAL_APP=true make release-candidate-check`; local app startup and the
  five public feed path diagnostics completed.
- Release-candidate diagnostic: Phase 89 CP000003 helper exited `0` with
  helper overall `not_checked` because release package audit was intentionally
  not run.
- Local release package: None; blocked without separate maintainer
  authorization.
- Local Docker image build: None.
- Published production Docker image: None.

## Migrations

- None in Phase 75-89.

Before any future upgrade, operators must back up the database and run:

```bash
make migrate-status
make migrate-up
make migrate-status
```

## Operations Changes

- Operations Console workflows are substantially clearer and more complete,
  but remain private, self-hosted, and operator-controlled.
- Browser pages added or improved in the Phase 75-89 track do not execute
  destructive maintenance actions, release actions, package actions, evidence
  actions, portal actions, or consumer submission actions.

## Security Notes

- Private Operations Console route boundaries, role checks, agency query
  guards, CSRF expectations for browser POSTs, bearer-token telemetry/device
  boundaries, server-owned validator mappings, and no-raw-output HTML
  boundaries remain in force.
- Phase 89 did not add public admin routes or new browser mutation routes.
- Training guidance warns operators not to copy secrets, tokens, database URLs,
  raw private output, credentials, or private records into public docs or issue
  trackers.

## Dependency Changes

- None for Phase 89.
- Phase 89 local diagnostics found pinned validator tooling available in this
  checkout.

## Evidence Or Claim Changes

- None.

This draft does not create retained evidence, change consumer statuses, contact
external services, publish artifacts, tag a release, distribute a package, or
add claims of CAL-ITP/Caltrans compliance, consumer
submission/review/acceptance/ingestion/listing/display, agency
adoption/approval, final-root readiness, hosted service availability, paid
support, SLA/uptime, production readiness, vendor compatibility, hardware
certification, production AVL reliability, production-grade ETA quality,
validator-clean feeds, or release readiness.

## Known Blockers Matrix

| Area | Status | Exact blocker or note | Next owner |
| --- | --- | --- | --- |
| Release action authorization | blocked | No maintainer authorization exists to tag, publish, distribute packages, publish images, or claim release readiness. | Maintainer |
| Release package | blocked | `make release-package` was not run because package creation requires separate explicit authorization. | Future authorized release-cut work |
| Release package audit | blocked | `make audit-release-package` was not run because no package artifact exists and package auditing requires separate explicit authorization. | Future authorized release-cut work |
| Artifact checksums | not_checked | No package artifact exists for this draft. | Future authorized release-cut work |
| SBOM/provenance | not_checked | No package artifact exists for this draft. | Future authorized release-cut work |
| Git tag | blocked | No tag was created. | Maintainer only |
| Published image | blocked | No image was built or published. | Future authorized release-cut work |
| Evidence tracks | blocked | Final-root, consumer, real agency pilot, real vendor/device, ETA-quality, and compliance evidence gates require separate written authorization. | Authorization-gated only |
| Final conclusion | needs_review | Product diagnostics passed so far, but release actions and package checks remain unauthorized/not run. | Maintainer review |

## Known Limitations

- The draft is not tied to a release tag.
- No package artifact, checksum, SBOM/provenance metadata, published image, or
  GitHub release exists.
- Release package checks remain blocked without separate authorization.
- Connector/adaptor results are synthetic/local checks only and do not prove
  real vendor/device compatibility.
- Validator and route diagnostics are local product signals only and do not
  prove compliance, consumer action, final-root readiness, production
  readiness, hosted-service availability, or SLA/uptime.

## Checks

Recorded Phase 89 checks through CP000005:

| Check | Result | Scope |
| --- | --- | --- |
| `git status --short` | passed in CP000002 before diagnostics | Clean local source state |
| `git diff --check` | passed in CP000002 and follow-up doc checkpoints | Whitespace guard |
| `make check` | passed | Lightweight local repo check |
| `make validate` | passed | Local validation smoke |
| `make test` | passed | Local tests |
| `RUN_LOCAL_APP=true make release-candidate-check` | exited `0`; helper overall `not_checked` | Local app startup and five public feed diagnostics; package audit not run |
| Focused private Operations Console route tests | passed | Navigation, route titles, contextual help, Help, feed readiness, consumers, access, and audit |
| `make external-connection-check` | passed | Local synthetic connector boundary check |
| `make adapter-conformance` | passed | Local synthetic adapter conformance |
| `make test-connector-examples` | passed | Local synthetic connector examples |
| `docker compose -f deploy/docker-compose.yml config` | passed | Compose config rendered locally |
| `make audit-product-acceptance` | passed | Product acceptance guard |
| `make audit-final-claim-review` | passed | Claim-boundary guard |
| Consumer tracker prepared-only review | passed | Seven targets remain exactly `prepared` |
| Protected-path status review | passed | Protected evidence paths, migrations, `go.mod`, and `go.sum` unchanged |

Blocked checks:

- `make release-package` was not run; package creation requires separate
  explicit maintainer authorization.
- `make audit-release-package` was not run; no package artifact exists and
  package auditing requires separate explicit maintainer authorization.

Release-candidate summary:

- Output summary: see `docs/phase-89-rc1-gate-results.md`.
- Overall status: `needs_review`.
- Package audit status: `blocked/not_checked`.
- Claim-boundary result: no retained evidence, consumer status change, image
  push, hosted-service claim, production-readiness claim, release-readiness
  claim, or compliance claim was added.
