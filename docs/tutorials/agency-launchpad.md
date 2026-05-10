# Agency Launchpad

This tutorial is a private, authenticated launchpad workflow for agency
operators and maintainers. It combines existing local/reference checks into one
review sequence for setup, GTFS import, publication metadata, the five feed
paths, telemetry, validators, readiness, connector conformance, and support
bundle collection.

The workflow is private diagnostics only. It performs no evidence creation, no
external contact, and no consumer/status mutation. It makes no approval,
compliance, production, public-launch, vendor, ETA, or SLA claim. All claim
flags remain `false`.

## 1. Setup

Start from a clean local shell in the repository root:

```bash
cd <repo-root>
make deps
make validators-check
```

If validator tooling is missing and the launchpad review should include
validator execution, install the pinned tooling first:

```bash
make validators-install
make validators-check
```

Bring up the local/reference app package when the review is local:

```bash
make agency-app-up
```

For an already-running private deployment boundary, use a loopback, VPN, SSH
tunnel, or otherwise admin-protected `ADMIN_BASE_URL`. Do not expose admin,
debug, GTFS Studio, or Operations Console routes on the anonymous feed edge.

## 2. GTFS Import

For the no-external-contact launchpad path, serve a reviewed GTFS ZIP from
loopback. This example uses the committed fixture:

```bash
tmpdir="$(mktemp -d)"
(cd testdata/gtfs/valid-small && zip -qr "$tmpdir/valid-small.zip" .)
python3 -m http.server 18080 --directory "$tmpdir"
```

Then import it in another terminal. The `AGENCY_ID` must match `agency.txt`.

```bash
make agency-pilot-up \
  AGENCY_ID=demo-agency \
  GTFS_URL=http://127.0.0.1:18080/valid-small.zip \
  SKIP_VALIDATORS=true
```

This verifies the command path and keeps output local/private by default.
For a non-loopback GTFS source, use
[Reusable Agency Onboarding](reusable-agency-onboarding.md) and record that
the launchpad no-external-contact path was not used.

## 3. Metadata

Supply feed metadata during onboarding when the operator has reviewed values
for the private launchpad environment:

```bash
scripts/agency-pilot-onboard.sh \
  --agency-id demo-agency \
  --gtfs-url http://127.0.0.1:18080/valid-small.zip \
  --technical-contact-email ops@example.org \
  --feed-license-name "CC BY 4.0" \
  --feed-license-url https://example.org/license
```

If these values are omitted, the script uses obvious placeholders. Treat
placeholder metadata as a blocker for any later externally shared packet or
agency-facing review.

## 4. Five Feeds

Verify that the five expected feed paths respond from the configured public
base URL:

```bash
BASE=http://localhost:8080
mkdir -p .cache/agency-launchpad-fetch
curl -fsS "$BASE/public/feeds.json" -o .cache/agency-launchpad-fetch/feeds.json
curl -fsS "$BASE/public/gtfs/schedule.zip" -o .cache/agency-launchpad-fetch/schedule.zip
curl -fsS "$BASE/public/gtfsrt/vehicle_positions.pb" -o .cache/agency-launchpad-fetch/vehicle_positions.pb
curl -fsS "$BASE/public/gtfsrt/trip_updates.pb" -o .cache/agency-launchpad-fetch/trip_updates.pb
curl -fsS "$BASE/public/gtfsrt/alerts.pb" -o .cache/agency-launchpad-fetch/alerts.pb
```

The five paths are:

- `/public/feeds.json`
- `/public/gtfs/schedule.zip`
- `/public/gtfsrt/vehicle_positions.pb`
- `/public/gtfsrt/trip_updates.pb`
- `/public/gtfsrt/alerts.pb`

The protobuf paths are anonymous feed paths. JSON debug, admin, and GTFS
Studio routes remain authenticated or private.

## 5. Telemetry

Use the simulator for synthetic telemetry through the authenticated ingest
path:

```bash
make telemetry-simulator
RUN_MATCHER=true make telemetry-simulator
```

Use the synthetic AVL dry-run when reviewing adapter transform shape without
sending telemetry:

```bash
go run ./cmd/avl-vendor-adapter --dry-run \
  --reference-time 2026-05-04T12:00:00Z \
  --mapping testdata/avl-vendor/mapping.json \
  testdata/avl-vendor/minimal-gps.json
```

Keep real device tokens, bearer tokens, private payloads, and database URLs out
of committed docs and shared notes.

## 6. Validators

Run validator tooling checks:

```bash
make validators-check
```

For onboarding, choose one validator mode deliberately:

```bash
make agency-pilot-up \
  AGENCY_ID=demo-agency \
  GTFS_URL=http://127.0.0.1:18080/valid-small.zip \
  STRICT_VALIDATORS=true
```

or:

```bash
make agency-pilot-up \
  AGENCY_ID=demo-agency \
  GTFS_URL=http://127.0.0.1:18080/valid-small.zip \
  SKIP_VALIDATORS=true
```

Use strict mode when validator blockers should stop the launchpad review. Use
skip mode only when the goal is import and feed-path inspection.

## 7. Readiness

Open the authenticated readiness page through the private admin boundary:

```text
/admin/operations/readiness
```

Review schedule, Vehicle Positions, Trip Updates, Alerts, metadata,
validation, telemetry freshness, operations status, and consumer packet
preparedness rows. This page is read-only for launchpad review and should be
used to decide what remains blocked, missing, or ready for a later operator
decision.

For a bounded private preflight summary, run:

```bash
make release-candidate-check
```

To include local app startup and five-feed fetches:

```bash
RUN_LOCAL_APP=true make release-candidate-check
```

Inspect `.cache/release-candidate-check/<timestamp>/summary.json` and confirm
every value under `claim_flags` is `false`.

## 8. Connector Conformance

Run connector and adapter checks before using optional telemetry, prediction,
validator, monitoring, or discovery sidecars:

```bash
make external-connection-check
make adapter-conformance
make test-connector-examples
```

These checks use local manifests and synthetic fixtures. They do not load
dynamic plugins, start sidecars, send notifications, automate portals, or
change tracked consumer records.

## 9. Support Bundle

Collect a redaction-safe private support bundle:

```bash
make support-bundle
```

To include authenticated readiness status, opt in explicitly:

```bash
PUBLIC_BASE_URL=http://localhost:8080 \
ADMIN_BASE_URL=http://127.0.0.1:8081 \
ADMIN_TOKEN=<redacted-admin-token> \
INCLUDE_ADMIN_READINESS=true \
make support-bundle
```

Review the generated manifest before sharing anything. The support bundle must
not contain admin tokens, device tokens, JWTs, cookies, API keys, raw `.env`
files, database dumps, raw private telemetry, private payloads, or unredacted
logs.

## 10. Decision Gate

Use this private gate before proceeding:

| Area | Continue only when |
| --- | --- |
| Setup | Required tools and local/reference services are available or blockers are documented. |
| GTFS | The intended GTFS ZIP imports and the active feed is not the sample feed by accident. |
| Metadata | Contact and license values are reviewed or placeholders are explicitly blocked. |
| Five feeds | All five feed paths respond from the intended base URL. |
| Telemetry | Synthetic telemetry or private device telemetry reaches authenticated ingest as expected. |
| Validators | Validator state is passed, skipped by intent, or blocked with a clear next action. |
| Readiness | `/admin/operations/readiness` has been reviewed through the private admin boundary. |
| Connectors | Local connector and adapter conformance checks pass for the intended boundary. |
| Support bundle | A private redaction-safe bundle exists when maintainers need diagnostics. |
| Claims | All claim flags remain `false`. |

If any row is blocked, stop and fix the blocker or record the next action in
private operator notes. Do not convert this launchpad output into external
claims or tracked consumer record changes.
