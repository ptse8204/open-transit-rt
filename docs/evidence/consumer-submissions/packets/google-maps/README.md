# Google Maps Prepared Submission Packet

## Status

- Target name: Google Maps
- Status: `prepared`
- Prepared at: `2026-04-27T04:42:17Z`
- Prepared by: Codex Phase 20 docs/evidence pass
- Evidence snapshot: OCI pilot final current-live recheck at `2026-04-24T16:38:46Z`
- OCI packet reference: `docs/evidence/captured/oci-pilot/2026-04-24/`
- `feeds.json` snapshot reference: `docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/public_feeds.json`
- Validator records reference: `docs/evidence/captured/oci-pilot/2026-04-24/validator-record-2026-04-24.md`
- Phase 19 replay/quality summary reference: `docs/handoffs/phase-19.md`

## Operator Warning

Do not submit from repo docs alone. Before any actual submission, the operator must review feed URLs, license/contact metadata, validation status, agency identity, consumer-specific requirements, and redactions.

## Submission Method

- Submission method: not verified
- Official submission URL/contact: not verified
- Verified as current: not verified
- Notes: This packet intentionally does not guess a Google Maps submission path. The operator must identify and verify the current official Google transit partner workflow outside this repo before submission.

## Prepared Feed URLs

- Feed root: `https://open-transit-pilot.duckdns.org`
- Feed discovery: `https://open-transit-pilot.duckdns.org/public/feeds.json`
- Static GTFS schedule: `https://open-transit-pilot.duckdns.org/public/gtfs/schedule.zip`
- GTFS-RT Vehicle Positions: `https://open-transit-pilot.duckdns.org/public/gtfsrt/vehicle_positions.pb`
- GTFS-RT Trip Updates: `https://open-transit-pilot.duckdns.org/public/gtfsrt/trip_updates.pb`
- GTFS-RT Alerts: `https://open-transit-pilot.duckdns.org/public/gtfsrt/alerts.pb`

## License, Contact, And Identity Metadata

- Agency ID from `feeds.json`: `demo-agency`
- Agency name from `feeds.json`: `Demo Transit`
- Technical contact email from `feeds.json`: `ops@example.org`
- License name from `feeds.json`: `CC-BY-4.0`
- License URL from `feeds.json`: `https://creativecommons.org/licenses/by/4.0/`
- Publication environment from `feeds.json`: `production`

## Evidence References

- Phase 12 hosted/operator evidence: `docs/evidence/captured/oci-pilot/2026-04-24/`
- Public feed proof: `docs/evidence/captured/oci-pilot/2026-04-24/public-feed-proof-2026-04-24.md`
- Final `feeds.json` snapshot: `docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/public_feeds.json`
- Validator records: `docs/evidence/captured/oci-pilot/2026-04-24/validator-record-2026-04-24.md`
- Scorecard export: `docs/evidence/captured/oci-pilot/2026-04-24/scorecard-export-2026-04-24.md`
- Phase 19 replay/quality summary: `docs/handoffs/phase-19.md`

## Phase 19 Realtime Quality Boundary

Phase 19 replay metrics measure current deterministic behavior and make unknown, ambiguous, stale, withheld, and degraded cases visible. They do not prove production-grade ETA quality or Google Maps acceptance.

## Redaction Notes

This packet contains public feed URLs, public metadata, and repo evidence paths only. It includes no portal credentials, private ticket links, private correspondence, tokens, DB URLs, private operator artifacts, or raw personal data.

## Next Action

Operator reviews the packet, verifies Google Maps' current official submission requirements and contact path, confirms the agency identity and metadata are approved for submission, then submits outside the repo only if authorized. Store redacted receipt or correspondence evidence before changing status to `submitted` or later.

## Allowed Public Wording

"A Google Maps submission packet has been prepared, but no submission or acceptance is confirmed."

## Claim Boundary

This packet does not prove Google Maps submission, review, ingestion, acceptance, consumer display, CAL-ITP/Caltrans compliance, marketplace/vendor equivalence, hosted SaaS availability, agency endorsement, or production-grade ETA quality.
