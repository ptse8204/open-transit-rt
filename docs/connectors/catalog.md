# Connector Catalog

Use this catalog to choose the safest starting point for an external connector.
Every entry is local, synthetic, disabled by default, or deployment-owned until
a technical helper configures it outside the browser.

The private browser pages are:

```text
/admin/operations/connectors
/admin/operations/connectors/workbench
/admin/operations/connectors/tests
```

They review connector shapes and fixed local checks. They do not run sidecars,
load dynamic backend plugins, contact vendors or consumers, create evidence,
send telemetry, or change consumer tracker status.

In Open Transit RT, a plugin is an optional sidecar, command adapter, manifest,
or connector process. It is not arbitrary dynamic code loaded into the backend.

## Support Matrix At A Glance

### Works today / local-supported

- CSV replay telemetry adapter.
- HTTP polling telemetry adapter.
- Webhook sidecar telemetry adapter.
- Generic JSON transform adapter.
- Authenticated `POST /v1/telemetry`.
- Deterministic built-in predictor.
- External HTTP predictor adapter and shadow mode.
- MobilityData static GTFS validator wrapper.
- MobilityData GTFS Realtime validator wrapper.
- Monitoring/export helper.
- `/public/feeds.json` discovery metadata.
- Static GTFS URL.
- Vehicle Positions URL.
- Trip Updates URL.
- Alerts URL.

### Planned / candidate

These items are roadmap direction only. They are not compatibility,
production, vendor, consumer, or compliance claims.

- TheTransitClock behind `internal/prediction.Adapter`.
- Real vendor AVL payload adapters behind `/v1/telemetry`.
- SIRI / GTFS-RT bridge.
- GTFS-Flex / demand-response QA.
- GTFS-ride / ridership analytics QA.
- GTFS-Pathways accessibility metadata QA.
- GTFS-Fares v2 fare metadata QA.
- OpenTripPlanner compatibility checks.
- OneBusAway compatibility checks.
- MobilityData validator UX and report explanation.
- Transitland / Mobility Database discovery-readiness workflows.

## Ecosystem Grounding

These connector paths are aligned to current GTFS / GTFS-Realtime ecosystem
needs without claiming external acceptance:

- The GTFS Realtime reference defines Trip Updates, Vehicle Positions, and
  Alerts, and requires `start_time` plus `start_date` for frequency-based
  trip descriptors. See the [GTFS Realtime reference](https://gtfs.org/documentation/realtime/reference/).
- GTFS Realtime best practices emphasize public permanent URLs, directly
  accessible protobuf responses, frequent refresh, and bounded data age. See
  [GTFS Realtime best practices](https://gtfs.org/documentation/realtime/realtime-best-practices/).
- MobilityData maintains the canonical static GTFS validator with web, desktop,
  command-line, and Docker workflows. See
  [MobilityData/gtfs-validator](https://github.com/MobilityData/gtfs-validator).
- MobilityData publishes a GTFS Realtime validator container package, derived
  from the CUTR validator lineage. See
  [MobilityData GTFS Realtime validator package](https://github.com/orgs/MobilityData/packages/container/package/gtfs-realtime-validator).
- Transitland exposes static and GTFS Realtime feed discovery APIs, including
  latest realtime downloads by `alerts`, `trip_updates`, and
  `vehicle_positions`. See [Transitland feeds API](https://www.transit.land/documentation/rest-api/feeds).
- Mobility Database replaced TransitFeeds for current GTFS and GTFS-Realtime
  discovery data. See [Mobility Database](https://mobilitydatabase.org/about).
- OpenTripPlanner, OneBusAway, and TheTransitClock remain relevant open-source
  ecosystem projects to evaluate through adapters or compatibility checks. See
  [OpenTripPlanner docs](https://docs.opentripplanner.org/en/latest/),
  [OneBusAway](https://github.com/OneBusAway/onebusaway-application-modules),
  and [TheTransitClock](https://thetransitclock.github.io/). They are not
  built-in dependencies or supported integrations.

## Copy/Adapt Path

1. Pick the connector category below.
2. Copy the closest example under `examples/connectors/` or start from the
   named internal contract.
3. Keep real credentials, live endpoints, private payloads, private paths, and
   portal details out of the repo.
4. Run the first safe check.
5. Review the Connector Workbench before any deployment-owned send path is
   enabled.
6. Keep the claim boundary with the connector.

## Vehicle / GPS / AVL Connectors

| Connector | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| CSV replay adapter | `examples/connectors/telemetry-csv-replay` | `make test-connector-examples` | Production data quality, real device proof, vendor compatibility, hardware certification, compliance, consumer acceptance, or production readiness. |
| HTTP polling adapter | `examples/connectors/telemetry-http-poller` | `make external-connection-check` | Live endpoint behavior, named API support, vendor compatibility, production AVL reliability, compliance, or public launch. |
| Webhook sidecar adapter | `examples/connectors/telemetry-webhook-sidecar` | `go run ./cmd/adapter-conformance telemetry --suite testdata/adapter-conformance` | Vendor support, hardware certification, production AVL reliability, SLA, compliance, consumer acceptance, or agency approval. |
| Generic JSON transform adapter | `examples/connectors/generic-json-transform` | `go test ./examples/connectors/generic-json-transform` | Vendor compatibility, real AVL reliability, production readiness, consumer acceptance, compliance, or ETA quality. |
| Vendor payload transform | `cmd/avl-vendor-adapter --dry-run` with `testdata/avl-vendor` | `go test ./cmd/avl-vendor-adapter` | Vendor compatibility, real AVL reliability, production readiness, consumer acceptance, compliance, or ETA quality. |
| Vendor-shaped synthetic examples | `testdata/avl-vendor` fixture families | `go test ./cmd/avl-vendor-adapter` | Named vendor support, hardware certification, real source behavior, production AVL reliability, or agency approval. |
| Authenticated telemetry POST | `POST /v1/telemetry` with deployment-owned device token | `make telemetry-simulator` | Real device proof, vendor compatibility, production AVL reliability, compliance, consumer acceptance, or production readiness. |

## Prediction Connectors

| Connector | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| Deterministic built-in predictor | Existing `internal/prediction.Adapter` path | `go test ./internal/prediction` | Production-grade ETA quality, real-world ETA accuracy, consumer acceptance, or production readiness. |
| External HTTP predictor adapter | `examples/connectors/predictor-sidecar-stub` | `go run ./cmd/adapter-conformance prediction --suite testdata/adapter-conformance` | Named predictor compatibility, live service behavior, production readiness, consumer acceptance, or ETA quality. |
| Shadow-mode predictor | Compare sidecar output while public Trip Updates stay unchanged | `go run ./examples/connectors/predictor-sidecar-stub` | External predictor acceptance, production-grade ETA quality, real-world accuracy, or public-feed improvement. |
| Fail-closed predictor behavior | Timeout, malformed, stale, wrong-agency, low-confidence, missing Vehicle Positions, and public-mutation cases | `make adapter-conformance` | ETA quality, real-world accuracy, named predictor compatibility, production readiness, or consumer acceptance. |
| TheTransitClock candidate notes | Future candidate behind the external predictor adapter only | `make adapter-conformance` | TheTransitClock integration, compatibility, certification, production ETA quality, or real-world accuracy. |

## Validator Connectors

| Connector | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| MobilityData static GTFS validator | Server-owned static validator tooling when installed | `make validators-check` | Validator-clean proof, CAL-ITP/Caltrans compliance, consumer acceptance, public launch, or production readiness. |
| MobilityData GTFS Realtime validator | Server-owned realtime validator wrapper when installed | `make gtfsrt-conformance` | Validator-clean proof, consumer acceptance, production readiness, public launch, or compliance. |
| Allowlisted validator IDs | `examples/connectors/validator-allowlist` | `go run ./cmd/adapter-conformance validator --suite testdata/adapter-conformance` | Raw validator command safety beyond the allowlist, validator-clean output, compliance, or consumer acceptance. |
| Private validation health | `/admin/operations/validation-health` | `make validate` | External approval, consumer acceptance, compliance, production readiness, or public launch. |

## Monitoring / Export Connectors

| Connector | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| Local health summaries | Feed, validator, telemetry, maintenance, and reliability summaries | `make operations-reliability` | SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence. |
| Operations notify draft | Deployment-owned notification draft, no send by default | `make operations-notify` | Notification delivery, incident response maturity, SLA, uptime, hosted service availability, or evidence creation. |
| Monitoring/export helper | `examples/connectors/monitoring-export` | `go run ./cmd/adapter-conformance monitoring --suite testdata/adapter-conformance` | Notification delivery, SLA coverage, uptime guarantee, hosted service availability, production readiness, or retained evidence. |
| Deployment-owned monitoring boundary | Real destination, credential, retention, and alert routing outside this browser page | `make external-connection-check` | Hosted service availability, paid support, SLA coverage, uptime guarantee, production readiness, or agency approval. |

## Consumer / Discovery Connectors

| Connector | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| `/public/feeds.json` | Public feed metadata endpoint for a running instance | `make smoke` | Consumer submission, review, acceptance, ingestion, listing, display, compliance, or production readiness. |
| Static GTFS URL | `/public/gtfs/schedule.zip` when an active feed version exists | `make validate-public-feeds` | Validator-clean proof, consumer acceptance, public launch, compliance, or production readiness. |
| Vehicle Positions URL | `/public/gtfsrt/vehicle_positions.pb` | `make gtfsrt-conformance` | Consumer ingestion, public display, production AVL reliability, compliance, or production readiness. |
| Trip Updates URL | `/public/gtfsrt/trip_updates.pb` behind the prediction boundary | `make gtfsrt-conformance` | Production-grade ETA quality, real-world accuracy, consumer acceptance, compliance, or production readiness. |
| Alerts URL | `/public/gtfsrt/alerts.pb` plus the private Alerts Console | `make gtfsrt-conformance` | Consumer acceptance, complete disruption operations, public display, compliance, or production readiness. |
| Consumer packet preparedness | `examples/connectors/consumer-discovery-metadata` and prepared-only tracker review | `go run ./cmd/adapter-conformance consumer_discovery --suite testdata/adapter-conformance` | Consumer submission, review, acceptance, ingestion, listing, display, compliance, or public launch. |

## Roadmap-Only Connector Candidates

These are useful directions for CAL-ITP-style readiness and GTFS-RT ecosystem
fit. They are not supported integrations until implemented, tested, and
documented behind the same adapter boundaries.

| Candidate | Why it matters | Required boundary | Current status |
| --- | --- | --- | --- |
| TheTransitClock integration | External prediction software can inform future Trip Updates evaluation. | `internal/prediction.Adapter`, external HTTP shadow/fail-closed modes, no Vehicle Positions dependency. | Candidate only; no compatibility, ETA-quality, or deployment claim. |
| Real vendor AVL payload adapters | Agencies often receive GPS/AVL data in vendor-specific shapes. | Transform outside core into authenticated `POST /v1/telemetry`; fixtures must stay synthetic or redacted. | Roadmap only; no vendor compatibility or hardware certification claim. |
| SIRI / GTFS-RT bridge | Some agencies or regions have SIRI-like realtime systems rather than direct GTFS-RT producers. | Sidecar transform into GTFS-RT-oriented telemetry/prediction DTOs, with no raw private payload commits. | Investigation only. |
| GTFS-Flex / demand-response QA | Flexible service metadata may matter for demand-response or deviated-route evaluation. | Static GTFS QA/readiness helper, not a dispatch or booking system. | Future investigation only. |
| GTFS-ride / ridership analytics QA | Ridership standards can help future planning analytics without becoming an unrelated analytics product. | Optional import/QA helper separate from public GTFS-RT publication. | Future investigation only. |
| GTFS-Pathways accessibility QA | Pathways and accessibility metadata affect station wayfinding and readiness. | Static GTFS QA helper and validator/report explanation. | Future investigation only. |
| GTFS-Fares v2 metadata QA | Fare metadata is part of broader public-data readiness. | Static GTFS QA helper; no fare payment product. | Future investigation only. |
| OpenTripPlanner compatibility checks | OTP consumes GTFS/OpenStreetMap and applies realtime updates and alerts. | Offline feed-consumer compatibility check with synthetic/local data. | Investigation only; no OTP compatibility claim. |
| OneBusAway compatibility or adapter investigation | OneBusAway exports GTFS-RT and is a relevant open-source realtime information system. | Offline compatibility/adaptor review; no public consumer status changes. | Investigation only; no OBA compatibility claim. |
| MobilityData validator UX integration | Users need plain explanations of validator reports. | Server-owned allowlisted validator IDs and browser report explanation. | Partially supported through wrappers; richer UX is roadmap. |
| Transitland / Mobility Database discovery readiness | Discovery workflows need stable URLs, license/contact metadata, and feed metadata. | `/public/feeds.json`, prepared-only local review, no portal automation. | Prepared-readiness only; no submission or listing claim. |

## Future Connector Extension Model

| Rule | Start with | First safe check | Does not prove |
| --- | --- | --- | --- |
| Manifest-based sidecars | `docs/connectors/plugin-contract.md` | `make external-connection-check` | Runtime installation, automatic backend extension, vendor compatibility, compliance, or production readiness. |
| No arbitrary dynamic backend plugin loading | Safe plugin definition above | `make external-connection-check` | Third-party runtime distribution support, remote plugin installation, or unreviewed extension execution. |
| Conformance tests required | `testdata/adapter-conformance` and synthetic fixtures | `make adapter-conformance` | Real integration proof, production readiness, vendor compatibility, compliance, or consumer acceptance. |

## Local Checks

Run:

```sh
make external-connection-check
make adapter-conformance
make test-connector-examples
```

These checks are local quality signals only. They do not contact external
systems, write evidence, move consumer status, prove compliance, prove
production readiness, prove consumer acceptance, prove vendor compatibility,
or prove ETA quality.
