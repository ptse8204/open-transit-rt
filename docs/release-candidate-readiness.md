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
