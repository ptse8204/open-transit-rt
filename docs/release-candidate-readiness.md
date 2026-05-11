# Release-Candidate Readiness

## Purpose

`make release-candidate-check` is a local preflight for maintainers who want to
evaluate whether a checkout is ready for a release-candidate review. It writes
private diagnostics under `.cache/release-candidate-check/<timestamp>` by
default.

The check does not tag, publish, push images, create retained evidence, contact
external services, change consumer statuses, or make production-readiness,
consumer-acceptance, agency-approval, hosted-service, vendor-compatibility, or
CAL-ITP/Caltrans compliance claims.

## What It Covers

The readiness check summarizes:

- fresh-clone prerequisites and required repo files;
- pinned validator tooling status;
- Docker Compose configuration;
- final claim audit status;
- Go validation, test, and smoke command follow-ups;
- local app startup path;
- GTFS import fixture path;
- five public feed paths;
- telemetry simulator dry-run path;
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

## First Release-Candidate Review Sequence

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
   RELEASE_PACKAGE_VERSION=<tag-or-rc> make release-package
   RELEASE_PACKAGE_DIR=.cache/release-package/<tag-or-rc> make audit-release-package
   RELEASE_PACKAGE_DIR=.cache/release-package/<tag-or-rc> RUN_RELEASE_PACKAGE=true make release-candidate-check
   ```

   Release packages under `.cache` are local diagnostics until a maintainer
   cuts an actual release.

6. Draft release notes from `docs/release-notes-template.md`. State `None`
   for unchanged migration, security, dependency, operations, and evidence or
   claim sections. List blocked commands exactly.

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
