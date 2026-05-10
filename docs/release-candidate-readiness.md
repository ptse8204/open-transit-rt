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
