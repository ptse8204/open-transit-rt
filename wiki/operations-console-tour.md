# Operations Console Tour

The Operations Console is the private browser UI for local evaluation and
operator review. It is not public, not a portal automation surface, and not an
evidence collector.

Open it locally after startup:

```text
http://localhost:8080/admin/operations
```

## Home

The home page starts with **Start Here**. Use it as the first-run acceptance
path. It shows:

- no-developer and developer paths;
- ordered first-run tasks;
- five copyable public feed URLs;
- local demo / deployment / evidence boundary;
- all-false claim flags.

## Launchpad

Open:

```text
/admin/operations/launchpad
```

Launchpad keeps a broader private operator workflow with setup, GTFS, metadata,
five feeds, telemetry, validators, readiness, connector conformance, support
bundle, and decision-gate sections.

## Setup Wizard

Open:

```text
/admin/operations/setup-wizard
```

The wizard explains the staged setup path. It is read-only and points to the
existing pages that perform admin-only work.

## GTFS Import

Open:

```text
/admin/operations/gtfs-import
```

Admins can upload a GTFS ZIP or import from a safe URL. The import path uses the
existing importer and validation records. It does not create retained evidence
or agency approval.

## Feed Health

Open:

```text
/admin/operations/feed-health
```

Feed Health summarizes feeds.json, schedule, Vehicle Positions, Trip Updates,
and Alerts. Use this page to see whether URLs, validation context, freshness,
and next actions are visible.

## Readiness

Open:

```text
/admin/operations/readiness
```

Readiness is a private CAL-ITP-style checklist. It distinguishes UI signals from
missing deployment evidence and keeps missing data visible.

## Connector Hub

Open:

```text
/admin/operations/connectors
```

Connector Hub explains optional sidecars, manifests, command adapters, telemetry
sources, prediction boundaries, validator adapters, monitoring/export, and
consumer/discovery boundaries. It does not load arbitrary backend plugins or
prove vendor compatibility.

## Telemetry Simulator

Open:

```text
/admin/operations/telemetry-simulator
```

This page explains synthetic scenarios and copyable shell commands. It does not
collect tokens or execute commands in the browser.

## GTFS Quality

Open:

```text
/admin/operations/gtfs-quality
```

Use this page to review static validator and internal importer findings. Admins
can rerun the allowlisted static validator from the existing admin action.

## Validation Health

Open:

```text
/admin/operations/validation-health
```

Validation Health shows validator tooling state and feed artifact availability.
It runs only server-side allowlisted validator actions when an admin posts the
existing form.

## Reliability

Open:

```text
/admin/operations/reliability
```

Reliability shows private diagnostics from existing records. It does not prove
SLA coverage or uptime.

## Help

Open:

```text
/admin/operations/help
```

Help explains the concepts behind each major section and links back to the
right UI pages and docs.

## What Needs A Technical Helper

The UI is enough for first evaluation, but a technical helper is still needed
for local Docker issues, real deployment hosting, validator installation,
custom telemetry connectors, token security, and any future authorized evidence
intake.

## What This Does Not Prove

Using the Operations Console does not prove CAL-ITP/Caltrans compliance, agency
approval or adoption, consumer acceptance, final-root readiness, hosted SaaS,
production readiness, vendor compatibility, hardware certification, SLA/uptime
coverage, or production-grade ETA quality.
