# Reusable Agency Onboarding

This guide is for operators who want to import a specific agency GTFS ZIP into
the local/reference Open Transit RT flow without manually editing the database.

It is an evaluation and reference-deployment helper. It does not prove agency
approval, agency adoption, consumer acceptance, CAL-ITP/Caltrans compliance,
production readiness, vendor compatibility, final-root ownership, real
realtime data, or production-grade ETA quality.

For the complete guided local/reference operator trial that includes this
onboarding command, five public feed checks, `/admin/operations/readiness`,
validator handling, synthetic AVL dry-run, and cleanup, see
[Self-Hosted Operator Trial](self-hosted-operator-trial.md).

## Command

For local Compose evaluation:

```bash
make agency-pilot-up \
  AGENCY_ID=example-agency \
  GTFS_URL=https://example.org/path/to/gtfs.zip
```

The Make target runs:

```bash
scripts/agency-pilot-onboard.sh
```

You can also call the script directly:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id example-agency \
  --gtfs-url https://example.org/path/to/gtfs.zip \
  --public-base-url http://localhost:8080
```

The script does not call `make agency-app-up` and does not import the demo
sample feed.

## Next Integration Step

After the requested GTFS imports and the five public feed paths verify, use the
[Integration Adapter Kit](../integration-adapter-kit.md) to choose the next
telemetry path. The usual progression is to connect a private telemetry source
or run the synthetic AVL dry-run adapter, then send only validated observations
through the existing `/v1/telemetry` contract.

After publication metadata, feed URLs, and telemetry setup begin to exist,
review `/admin/operations/readiness` in the authenticated Operations Console.
It shows CAL-ITP-style readiness gaps with status sources and next actions,
but it does not create external evidence or claim CAL-ITP/Caltrans compliance.

That next step still does not prove real vendor compatibility, production AVL
reliability, consumer acceptance, CAL-ITP/Caltrans compliance, agency adoption,
or production-grade ETA quality.

## Required Inputs

`AGENCY_ID` must match the `agency_id` in the GTFS `agency.txt` file. It may
contain only letters, numbers, dot, underscore, and hyphen.

`GTFS_URL` must be an `http://` or `https://` URL to a GTFS ZIP.

The raw GTFS ZIP is downloaded under ignored storage:

```text
.cache/agency-pilot/<agency-id>/source.zip
```

Do not commit raw GTFS unless maintainers have reviewed the license, privacy,
and redaction implications.

## Metadata Defaults

Publication metadata is used by `/public/feeds.json` and readiness workflows.
The onboarding script accepts:

```bash
--technical-contact-email ops@example.org
--feed-license-name "CC BY 4.0"
--feed-license-url https://example.org/license
```

Environment fallbacks are:

```bash
TECHNICAL_CONTACT_EMAIL
FEED_LICENSE_NAME
FEED_LICENSE_URL
```

If these are not supplied, the script uses obvious local/reference placeholders
and prints:

```text
Publication metadata is local/reference placeholder metadata unless the operator supplied agency-approved values.
```

Replace placeholders before treating feed discovery metadata as agency-approved.
Do not treat placeholder metadata as agency-approved, final-root evidence,
consumer evidence, CAL-ITP/Caltrans compliance, or production readiness.

## Modes

Local Compose mode is the default:

```bash
scripts/agency-pilot-onboard.sh \
  --mode local-compose \
  --agency-id example-agency \
  --gtfs-url https://example.org/gtfs.zip
```

It starts Postgres, builds the local app image, applies migrations, creates the
requested agency/admin rows, starts services, imports the requested GTFS, and
fetches public feeds.

Running mode is for an already-running local/reference deployment:

```bash
ADMIN_TOKEN=<redacted> \
DATABASE_URL=<redacted-db-url> \
scripts/agency-pilot-onboard.sh \
  --mode running \
  --agency-id example-agency \
  --gtfs-url https://example.org/gtfs.zip \
  --public-base-url https://feeds.example.org \
  --admin-base-url http://127.0.0.1:8081
```

Running mode requires `ADMIN_TOKEN`, `DATABASE_URL`, and an explicit
`--admin-base-url` or `ADMIN_BASE_URL`. The admin base URL should normally be a
loopback, VPN, SSH tunnel, or otherwise private/admin-protected URL. The script
uses `DATABASE_URL` to upsert the requested agency/admin rows before import; it
does not reuse `scripts/seed-dev.sql`.

## Local State

If a previous local Compose database volume exists, the script reuses it
non-destructively and prints a warning.

To reset local state:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id example-agency \
  --gtfs-url https://example.org/gtfs.zip \
  --reset-local-state
```

The reset requires typed confirmation:

```text
reset-agency-pilot
```

Use `--force` only for an intentional scripted reset.

## Validators

Validator calls are best-effort by default. Missing pinned tooling records a
blocked validator status but does not fail the whole onboarding flow.

Install and check validators with:

```bash
make validators-install
make validators-check
```

Use `--strict-validators` when validator blockers should fail the command.
Use `--skip-validators` when you only want import and public fetch verification.

Validation output is not consumer acceptance, compliance proof, or production
readiness proof.

## Dry Run

To validate inputs and print planned paths and URLs without downloading,
importing, starting services, or running validators:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id example-agency \
  --gtfs-url https://example.org/gtfs.zip \
  --dry-run
```

`--help` and `--dry-run` do not require Docker, network access, a database,
validators, or secrets.

## What The Script Verifies

The script fetches and checks that these paths are non-empty:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

It also unzips the fetched schedule and compares public-safe schedule facts
against the source GTFS summary so the flow does not silently leave the sample
feed active.

## Output Boundary

The script prints feed URLs, admin URL, checksum paths, validator status, and
next steps. It does not create final-root evidence, external evidence packets,
consumer submission artifacts, or consumer status changes.
