# OCI Reference Diagnostic Runs

This page records public-safe summaries of reference deployment diagnostics.
It is not a retained evidence packet and does not replace the guarded
`docs/evidence` workflows.

These records may include public URLs, HTTP status codes, byte counts,
checksums, service-state summaries, source URLs, and explicit blockers. They
must not include secrets, admin tokens, populated environment files, private
keys, database URLs, raw downloaded GTFS ZIPs, private service logs, raw
`.cache` output, or consumer portal artifacts.

## 2026-05-12 OCI Reference GTFS Diagnostic

### Scope

This was a product-quality and OCI reference deployment diagnostic. It
exercised the current self-hosted reference path, imported one real public
static GTFS feed, verified the five public feed paths, and recorded remaining
operator blockers.

The run does not prove CAL-ITP/Caltrans compliance, consumer acceptance,
agency adoption or approval, agency-owned final-root readiness, hosted SaaS
availability, production readiness, vendor compatibility, SLA coverage,
hardware certification, or production-grade ETA quality.

### Deployment Context

| Item | Value |
| --- | --- |
| Diagnostic environment | `oci-reference-diagnostic` |
| Deployed commit | `fe9be0eddadfb73d6cf7c89cc3d384e54c94128d` |
| Public base URL | `https://open-transit-pilot.duckdns.org` |
| Access pattern | Public feed paths over HTTPS; admin/token work kept on loopback/private boundary |
| Server record boundary | Public URL and service status only; SSH keys, tokens, DB URLs, env files, and private logs are not retained here |

### Repository Gate Results

The local repository was clean at the end of the run.

| Command | Result |
| --- | --- |
| `git status --short` | Passed; empty output |
| `make check` | Passed |
| `make test` | Passed |
| `make validators-install` | Passed |
| `make validate` | Passed |
| `make external-connection-check` | Passed |
| `make adapter-conformance` | Passed |
| `make test-connector-examples` | Passed |
| `make audit-final-claim-review` | Passed |
| `docker compose -f deploy/docker-compose.yml config` | Passed; rendered config captured outside the repo |
| `make release-candidate-check` | Passed; private `.cache` output only |

### OCI Deployment Results

The current commit was built, pushed, migrated, restarted, and checked with
the existing OCI helper path:

```sh
scripts/oci-pilot.sh build
scripts/oci-pilot.sh push
scripts/oci-pilot.sh migrate
scripts/oci-pilot.sh restart
scripts/oci-pilot.sh status
```

Migration `000009_public_schedule_indexes.sql` was applied on the reference
server. A follow-up cache path reduced the warmed loopback schedule ZIP fetch
to about `0.085s` after an initial warm request.

Final service and loopback health summary:

| Service / Endpoint | Result |
| --- | --- |
| `open-transit-agency-config` | `active` |
| `open-transit-telemetry-ingest` | `active` |
| `open-transit-feed-vehicle-positions` | `active` |
| `open-transit-feed-trip-updates` | `active` |
| `open-transit-feed-alerts` | `active` |
| `http://127.0.0.1:8081/healthz` | `200` |
| `http://127.0.0.1:8082/healthz` | `200` |
| `http://127.0.0.1:8083/healthz` | `200` |
| `http://127.0.0.1:8084/healthz` | `200` |
| `http://127.0.0.1:8085/healthz` | `200` |

### Public GTFS Source

The smaller public Marin Transit GTFS feed was used after a larger LA Metro
feed exceeded the useful diagnostic time budget on the small reference host.

| Item | Value |
| --- | --- |
| Source URL | `https://marintransit.gov/data/google_transit.zip` |
| Access time UTC | `2026-05-12T01:06:18Z` |
| HTTP status | `200` |
| Byte count | `1391013` |
| Content type | `application/zip` |
| SHA-256 | `2774f794f0e2fda3e781ca0a19909e2c106b88e66195e338c0edd1a7c3b3d44b` |
| Visible terms source | `https://marintransit.gov/developers` |
| Visible terms note | Page contained `Open Data Commons Attribution License` and `opendatacommons.org/licenses/by/1.0` |

### Import Result

Import succeeded through the running OCI deployment using the helper flow and
loopback admin URL. The admin token was generated on the host and was not
printed or retained.

| Item | Value |
| --- | --- |
| Imported agency ID | `33` |
| Active feed ID | `gtfs-import-39` |
| Active feed source | `gtfs_import` |
| Active feed state | `active` |
| Activated at | `2026-05-12 00:59:00.839134+00` |
| Routes | `20` |
| Stops | `545` |
| Trips | `1612` |
| Stop times | `39416` |
| Shape points | `94589` |
| Schedule identity check | Passed |

### Five Public Feed Path Results

The five public paths were fetched from outside the host.

| Public path | HTTP status | Bytes | Content type | SHA-256 |
| --- | ---: | ---: | --- | --- |
| `/public/feeds.json` | `200` | `2548` | `application/json` | `e5f99e6b78f930ea3d9701d4c6a7347e6377820f86e6da0d4d160d5cbca7f6aa` |
| `/public/gtfs/schedule.zip` | `200` | `1296729` | `application/zip` | `e0c7b85c1c46c7315a9f449317773a50cd6f3301b4af34fef0444ec61b05fb73` |
| `/public/gtfsrt/vehicle_positions.pb` | `200` | `15` | `application/x-protobuf` | `af886995f335d35b7039d6ad36e48519a3b0ac4fa603a01c227f6b3176276491` |
| `/public/gtfsrt/trip_updates.pb` | `200` | `15` | `application/x-protobuf` | `af886995f335d35b7039d6ad36e48519a3b0ac4fa603a01c227f6b3176276491` |
| `/public/gtfsrt/alerts.pb` | `200` | `15` | `application/x-protobuf` | `af886995f335d35b7039d6ad36e48519a3b0ac4fa603a01c227f6b3176276491` |

The fetched schedule ZIP contained the expected public agency row:

```text
agency_id,agency_name,agency_url,agency_timezone,agency_email
33,Marin Transit,https://marintransit.gov,America/Los_Angeles,info@marintransit.gov
```

### Remaining Blockers And Notes

- Static and realtime validator execution on the OCI micro host remains an
  operator tooling blocker. Java was not safely installed during this run, and
  the host is too memory-constrained for unbounded package installation during
  diagnostics.
- Backup and restore-drill values were not configured for this diagnostic.
- No device credential was available for agency `33`, so real telemetry
  simulator sends were not performed. A local dry run completed without
  sending telemetry.
- The deployed app artifact includes helper scripts but not the full Makefile
  or Go toolchain, so some diagnostics were run directly or locally against
  the public deployment.
- The public URL is a reference diagnostic root. Do not describe it as an
  agency-owned final public root or a production feed root.

### Follow-Up Work

Recommended next work:

1. Package or document remote diagnostic tooling so validator health,
   deployment doctor, operations reliability, and telemetry simulator have a
   consistent reference-server path.
2. Add a safe validator runtime path for small OCI hosts, or run validators
   off-host against fetched artifacts with recorded versions and checksums.
3. Configure backup and restore-drill values before using the reference server
   for longer-lived operator evaluation.
4. Add an operator-owned device credential only when real telemetry simulator
   sends are desired.
5. Prepare `v0.1.0-rc.1` only after the release-candidate gate and accepted
   deployment diagnostics are reviewed from a clean checkout.

