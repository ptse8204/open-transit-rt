# Self-Hosted Operations Notifications

Use this tutorial when a self-hosted operator wants one private local draft that
summarizes existing validation-health and reference deployment diagnostics.

This workflow does not send anything. It creates a local draft only.

## Run A Local Draft

```bash
make operations-notify
```

By default the script writes:

```text
.cache/operations-notify/<timestamp>/
```

with exactly:

```text
summary.json
summary.md
manifest.json
manifest.md
notification.txt
```

The draft begins with `DRAFT — NOT SENT`. It is meant for private operator
review before deciding what to do next.

`summary.json` includes a limited `health_digest` object and
`channel_guidance` object. The digest is a local template with severity, loaded
source counts, blocked/needs-review counts, capped next actions, and explicit
"does not prove" wording. Channel guidance records only booleans such as
`webhook.present` and `email.present`; it keeps `send_enabled=false` and
`destination_value_recorded=false`.

## Dry Run

```bash
scripts/operations-notify.sh --dry-run
```

Dry-run writes the same five local draft files. It does not require source
summaries, network, webhook, email, database, Docker, admin token, app, or a
running service.

## Inputs

Default mode selects the latest timestamp-named summaries under:

```text
.cache/validator-health/*/summary.json
.cache/deployment-doctor/*/summary.json
```

To provide explicit summaries:

```bash
VALIDATOR_HEALTH_SUMMARY=.cache/validator-health/20260509T120000Z/summary.json \
DEPLOYMENT_DOCTOR_SUMMARY=.cache/deployment-doctor/20260509T120000Z/summary.json \
scripts/operations-notify.sh
```

Explicit sources must stay under `.cache` unless
`ALLOW_UNIGNORED_SOURCE_DIR=true`. Evidence-like paths are always rejected.

## Strict Mode

Default mode exits `0` when summary generation succeeds, even when sources are
missing, malformed, oversized, or blocked. It records safe next actions instead.

Use strict mode for local automation that should fail on unhealthy inputs:

```bash
STRICT_OPERATIONS_NOTIFY=true scripts/operations-notify.sh
```

Strict mode exits nonzero for missing sources, malformed sources, oversized
sources, blocked generated severity, and unhealthy loaded source statuses.
Redaction failures exit nonzero in all modes.

## Destination Placeholders

If these variables are set:

```bash
NOTIFY_WEBHOOK_URL=https://example.invalid/placeholder
NOTIFY_EMAIL_TO=ops@example.invalid
```

the script records only booleans:

```json
{
  "destinations": {
    "webhook_present": true,
    "email_present": true
  }
}
```

It never writes destination values and never sends to webhook or email.

The reliability helper can consume this no-send summary:

```bash
make operations-reliability
```

Its `summary.json` includes `monitoring_export` and `private_ops_summary`
sections for local dashboards or technical-helper review. These sections are
still private summaries only; they do not send notifications, upload data,
create evidence, or prove hosted monitoring.

## What To Review

Open `summary.md` for a short operator view. Use `summary.json` for local
automation that needs the stable schema. Use `notification.txt` only as an
unsent draft for an operator to review.

The draft summarizes:

- validator-health status;
- public feed and deployment-doctor blocker counts;
- capped next actions;
- omitted item counts through `overflow_count`;
- a safety section showing that no external claim was created.

It does not copy raw reports, private paths, source JSON, evidence packets,
consumer submissions, tokens, webhook values, email recipients, cookies,
database URLs, stdout, stderr, or argv fields.

## Boundary

This workflow is private local diagnostics. It is not evidence, not a
compliance proof artifact, not production health proof, not CAL-ITP readiness
proof, and not consumer-readiness proof. It is not a compliance gate. It does
not contact consumers, change consumer statuses, edit GTFS, block publishing,
create final-root evidence, or claim production readiness.
