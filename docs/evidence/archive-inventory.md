# Evidence Archive Inventory

This inventory lists committed archives under `docs/evidence/captured/**`.
Opaque archives are not allowed; every retained archive must list its contents.

## Reviewed Archives

### `docs/evidence/captured/local-demo/2026-04-22/artifacts/public/schedule.zip`

- SHA-256: `0956ed037a40ca9d2cca94a501bea1547d27dbd25a195c7ebefe1a34ffc78194`
- Purpose: local demo static GTFS schedule artifact used by the local evidence
  packet.
- Decision: keep. The archive contains expected demo GTFS text files only.
- Contents:
  - `agency.txt`
  - `routes.txt`
  - `stops.txt`
  - `trips.txt`
  - `stop_times.txt`
  - `calendar.txt`
  - `shapes.txt`

### `docs/evidence/captured/oci-pilot/2026-04-24/artifacts/public/public_gtfs_schedule.zip`

- SHA-256: `7b0b389c41faa0fa0c01a25611f019822453f7a03bc9f1b547f7b16b66d99d83`
- Purpose: hosted OCI pilot public static GTFS schedule artifact.
- Decision: keep. The archive contains expected public GTFS text files only.
- Contents:
  - `agency.txt`
  - `feed_info.txt`
  - `routes.txt`
  - `stops.txt`
  - `trips.txt`
  - `stop_times.txt`
  - `calendar.txt`
  - `shapes.txt`

## Review Command

The Phase 15 inventory was produced with:

```bash
find docs/evidence/captured -type f \( -name '*.zip' -o -name '*.tar' -o -name '*.tgz' -o -name '*.tar.gz' -o -name '*.gz' -o -name '*.7z' -o -name '*.rar' \) -print
unzip -l <archive>
sha256sum <archive>
```
