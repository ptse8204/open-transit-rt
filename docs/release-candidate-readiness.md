# Release-Candidate Readiness

Current reader note: for the published `v0.1.0-rc.2` release-candidate status,
start with [Release Status v0.1.0-rc.2](release-status-v0.1.0-rc.2.md) and
[Release Download Replay v0.1.0-rc.2](release-download-replay-v0.1.0-rc.2.md).
This file remains as maintainer preflight history and still references earlier
rc1 review work.

## Purpose

`make release-candidate-check` is a local preflight for maintainers who want to
evaluate whether a checkout is ready for a release-candidate review. It writes
private diagnostics under `.cache/release-candidate-check/<timestamp>` by
default.

`v0.1.0-rc.1` is now published as a public release candidate for
local/self-hosted evaluation:
`https://github.com/ptse8204/open-transit-rt/releases/tag/v0.1.0-rc.1`.
Do not treat this first RC as a full `v0.1.0` release gate or production
readiness proof.

Phase 72 completed the bounded hardening review for `v0.1.0-rc.1`; see
[`docs/phase-72-v0.1.0-rc.1-release-candidate-hardening.md`](phase-72-v0.1.0-rc.1-release-candidate-hardening.md).
Phase 115 published the public rc1 prerelease, Phase 116 verified release
downloads and recorded the source-archive `make check` limitation, and
Phase 117 verified public fresh-clone install confidence.
Phase 72 Checkpoint 000001 is planning/status only. It does not run this gate,
does not prove any check passed, and does not draft release notes. Checkpoint
000003 hardens this diagnostic so its validator tooling row explicitly runs in
pinned mode even when the ambient shell has `VALIDATOR_TOOLING_MODE=stub`, and
so validator-tooling failures surface the exact final
`scripts/check-validators.sh` blocker line in `summary.json` and `summary.md`.
It does not install validators and does not turn missing pinned tooling into a
validator-clean, release-ready, or compliance claim.

Checkpoint 000004 verified local app startup, the required private Operations
Console route checks, five anonymous public feed fetches, and the
unauthenticated admin-route boundary. Checkpoint 000005 verified local synthetic
connector and adapter conformance gates through `make external-connection-check`,
`make adapter-conformance`, and `make test-connector-examples`. Phase 72
Checkpoint 000006 added a local pre-tag release notes and known blockers draft
at `docs/release-notes-v0.1.0-rc.1-draft.md`. Phase 73 Checkpoint 000001 is
complete for documentation-only agency UI acceptance planning. Phase 73
Checkpoint 000002 is complete for local no-developer browser walkthrough
review. Phase 73 Checkpoint 000003 is complete for local technical-helper
walkthrough review. Phase 73 Checkpoint 000004 is complete for narrow UI copy,
route-label, Devices/Telemetry boundary-copy, and browser-first tutorial
patching. Phase 73 Checkpoint 000005 is complete for small-agency docs and
wiki navigation freeze. The current default next work is Phase 73 Checkpoint
000006 close agency UI acceptance review.
Phase 72 closeout completed with `needs_review` release-candidate diagnostics,
not a release-ready pass.

The check does not tag, publish, push images, create retained evidence, contact
external services, change consumer statuses, or make production-readiness,
consumer-acceptance, agency-approval, hosted-service, vendor-compatibility, or
CAL-ITP/Caltrans compliance claims.

## What It Covers

The readiness check summarizes:

- fresh-clone prerequisites and required repo files;
- clean checkout source metadata;
- pinned validator tooling status;
- Docker Compose configuration;
- final claim audit status;
- Go validation, test, and smoke command follow-ups;
- local app startup path;
- GTFS import fixture path and a public GTFS trial when explicitly run;
- five public feed paths;
- telemetry simulator dry-run path;
- connector/adaptor conformance follow-ups;
- deployment doctor;
- validator health;
- operations reliability;
- operations notification draft;
- optional audit of an existing local `.cache` release package when explicitly
  enabled.

Some checks are intentionally opt-in because they start local services or rely
on a package generated outside this bounded preflight. Missing opt-in checks are
recorded as `not_checked`, not as a passed production signal. The script keeps
its retained repository output to the five files listed below.

## Usage

Run the lightweight evaluator check first. It does not install validators, pull
Docker images, start services, contact external systems, or write retained
evidence:

```sh
make check
```

Then run the release-candidate diagnostic:

```sh
make release-candidate-check
```

Fast dry-run used by repository validation:

```sh
scripts/release-candidate-check.sh --dry-run
```

Optional local app startup and five-feed fetch path:

```sh
RUN_LOCAL_APP=true make release-candidate-check
```

Optional audit of an existing local source package:

```sh
RELEASE_PACKAGE_DIR=.cache/release-package/<version> RUN_RELEASE_PACKAGE=true make release-candidate-check
```

## First `v0.1.0-rc.1` Review Sequence

Use this sequence for the first maintainer release-candidate review from a
clean checkout. The output is a local review packet, not proof of production
readiness.

1. Confirm the source state:

   ```sh
   git status --short
   git describe --tags --always --dirty
   git rev-parse HEAD
   ```

   Do not tag or publish from a dirty checkout. Dirty checkouts may still be
   used for local diagnostics while implementation work is active.

2. Run the bounded repo check:

   ```sh
   make check
   ```

   Fix JSON, shell syntax, protected-status, or claim-audit failures before
   proceeding.

3. Generate the private release-candidate diagnostic:

   ```sh
   make release-candidate-check
   ```

   Review `summary.md`, `summary.json`, `manifest.md`, `manifest.json`, and
   `check-log.txt` under `.cache/release-candidate-check/<timestamp>/`. The
   summary captures source metadata such as `git describe --tags --always
   --dirty`, commit SHA, branch, dirty state, and pre-tag review mode for
   release-note drafting.

4. Run validator, test, package, and Docker Compose checks as the local
   environment allows:

   ```sh
   scripts/bootstrap-dev.sh --check
   make validate
   make test
   make test-release-package
   docker compose -f deploy/docker-compose.yml config
   ```

   If Java, Docker, network access, pinned validator assets, or local ports are
   unavailable, record the exact command and blocker. Do not convert a skipped
   or blocked check into a readiness or compliance claim.

5. Generate and audit a local source package only when a maintainer needs a
   package review:

   ```sh
   RELEASE_PACKAGE_VERSION=v0.1.0-rc.1 make release-package
   RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 make audit-release-package
   RELEASE_PACKAGE_DIR=.cache/release-package/v0.1.0-rc.1 RUN_RELEASE_PACKAGE=true make release-candidate-check
   ```

   Release packages under `.cache` are local diagnostics until a maintainer
   cuts an actual release.

6. Run external-connection maturity checks:

   ```sh
   make external-connection-check
   make adapter-conformance
   make test-connector-examples
   ```

   Confirm telemetry connector paths send only to authenticated
   `POST /v1/telemetry`, external predictors stay behind
   `internal/prediction.Adapter`, deterministic Trip Updates remain the safe
   fallback, monitoring/export surfaces are redacted/no-send by default, and
   feed-consumer URL/metadata workflows do not automate submissions or status
   changes.

7. Draft release notes from `docs/release-notes-template.md`. State `None`
   for unchanged migration, security, dependency, operations, and evidence or
   claim sections. List blocked commands exactly.

## Post-RC Stabilization Review

After local package review, maintainers can run a post-RC bug bash before
considering any later release action. This review should stay local and
diagnostic:

```sh
make audit-operations-route-inventory
OPERATIONS_ROUTE_AUDIT_STRICT_DOCS=true scripts/audit-operations-route-inventory.sh
make check
make validate
make test
make external-connection-check
make adapter-conformance
make test-connector-examples
RUN_LOCAL_APP=true make release-candidate-check
docker compose -f deploy/docker-compose.yml config
make audit-product-acceptance
make audit-final-claim-review
```

Record exact pass, blocked, `needs_review`, and `not_checked` rows in the
phase or release-note blocker matrix. A post-RC bug bash does not tag a
release, publish a package, push an image, create retained evidence, move
consumer statuses, contact external parties, or approve stronger public
claims.

## Validation Matrix

| Check | Expected Local Output | Blocker Handling |
| --- | --- | --- |
| `make check` | No-network repo validation passes | Fix before continuing |
| `make release-candidate-check` | `.cache/release-candidate-check/<timestamp>/` with five files | Record blocker rows; do not claim production readiness |
| `make validate` | Validator-backed repo validation passes | Record Java, Docker, network, or pinned-tool blocker exactly |
| `make test` | Go tests pass | Fix or record blocker before tagging |
| `make test-release-package` | Local package helper tests pass | Fix package helper before relying on package diagnostics |
| `docker compose -f deploy/docker-compose.yml config` | Compose config renders | Record Docker CLI/Compose blocker exactly |
| `make audit-final-claim-review` | Claim and consumer tracker audit passes | Fix unsupported wording or protected status drift before continuing |
| `make external-connection-check` | Connector manifests pass local checks | Fix manifest or boundary drift before stronger connector wording |
| `make adapter-conformance` | Synthetic adapter conformance passes | Record exact unsupported adapter cases |
| `make telemetry-simulator` | Synthetic/local telemetry reaches the authenticated ingest path when configured | Record app, token, or local service blockers |

## Release Note Inputs

When a release candidate becomes a tagged release, release notes should record:

- source tag or planned tag;
- commit SHA;
- dirty or clean state;
- release-candidate diagnostic output directory;
- local release package path and checksum manifest, if generated;
- SBOM and provenance metadata status, if generated;
- local Docker image tag, if built;
- validation commands and results;
- blocked commands with exact reasons;
- confirmation that no retained evidence, consumer status change, image push,
  hosted-service claim, production-readiness claim, or compliance claim was
  created by the release-candidate helper.

## Local Package Audit Matrix

| Package Item | Review Command | Boundary |
| --- | --- | --- |
| Source archive and checksum | `make release-package` | Local package diagnostic only until a release is cut |
| Manifest files | `make audit-release-package` | Manifest presence is not hosted-service proof |
| SBOM/provenance metadata | `make audit-release-package` | Metadata is local supply-chain context |
| Dirty state | `make audit-release-package` | Dirty packages are not release-ready artifacts |
| Optional local image metadata | `RELEASE_PACKAGE_IMAGE_TAG=<tag> make release-package` | Local image metadata is not a published production image |

## Environment Blockers

`make release-candidate-check` records validator and local tooling blockers
instead of treating them as product claims. If pinned validator tooling is not
installed, the check may exit non-zero with a blocker such as a missing
MobilityData GTFS Validator JAR, missing Java runtime, unavailable Docker CLI,
or missing GTFS Realtime validator image. Install pinned tooling only when the
review needs validator execution:

```sh
make validators-install
make validators-check
```

Installing validators may require network access and Docker image pulls. If
network, Docker, Java, or the pinned validator image is unavailable, record the
exact blocker, confirm non-network checks such as `make check`, and continue
productization review without converting the blocker into a compliance,
production-readiness, or consumer-acceptance claim.

The release-candidate diagnostic intentionally forces the validator tooling
check to `VALIDATOR_TOOLING_MODE=pinned`. `VALIDATOR_TOOLING_MODE=stub` remains
available for targeted tests and smoke runs that call validator-aware code
directly, but it cannot make the release-candidate `validators_check` row pass.
In dry-run mode, `validators_check` remains `not_checked`.

When `RUN_LOCAL_APP=true` is used, likely blockers are Docker daemon
availability, Docker Compose plugin/config errors, first-run image pull or Go
module network access, host ports `8080` and `55432`, existing local database
volume state, or service readiness timeouts. Record those as local setup
blockers. A blocked or skipped local app run is not compliance evidence and is
not production-readiness proof.

## Outputs

The script writes exactly:

- `summary.json`
- `summary.md`
- `manifest.json`
- `manifest.md`
- `check-log.txt`

These files are private local diagnostics only. They are not evidence packets,
not release artifacts, and not publication approval.

`summary.json` and `summary.md` also include source metadata, release-note
inputs, the ordered first release-candidate workflow, and the local package
audit matrix. These fields help maintainers draft release notes from the local
checkout; they do not tag, publish, push images, or approve a release.

## Claim Boundary

Passing this check means only that the local repository readiness workflow
completed. It does not prove:

- CAL-ITP/Caltrans compliance;
- consumer submission, review, acceptance, ingestion, listing, or display;
- consumer acceptance;
- agency adoption, endorsement, or approval;
- agency-owned final-root readiness;
- hosted SaaS or paid support availability;
- SLA or uptime coverage;
- production readiness;
- universal deployment fitness;
- vendor compatibility, vendor AVL compatibility, or certified hardware support;
- production-grade ETA quality.
