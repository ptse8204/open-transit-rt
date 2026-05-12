# Off-Host Public Feed Validation

`scripts/validate-public-feeds.sh` fetches the five public Open Transit RT
feed paths from an operator machine and runs local/off-host validators where
the pinned tools are installed.

Use this path when a tiny server cannot run Java, Docker, or the validator
workload comfortably.

This helper is a diagnostic. Validator completion is a supporting signal only.
It is not CAL-ITP/Caltrans compliance, consumer acceptance, agency adoption,
final-root proof, hosted service availability, production readiness, vendor
compatibility, SLA coverage, or production-grade ETA quality.

## Run A Dry Run

```bash
OUTPUT_DIR=.cache/validate-public-feeds/dry-run \
FORCE=true \
scripts/validate-public-feeds.sh \
  --public-base-url https://feeds.example.org \
  --dry-run
```

## Run Against A Public Root

```bash
PUBLIC_BASE_URL=https://feeds.example.org make validate-public-feeds
```

The default output path is:

```text
.cache/validate-public-feeds/<timestamp>
```

The script fetches exactly:

```text
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
```

For each row it records the public path, configured URL, fetch state, HTTP
status, byte count, content type, SHA-256 checksum, validator state, next
action, and what the row does not prove.

## Validator Tooling

Install pinned validators on the operator machine when validation should run:

```bash
make validators-install
make validators-check
```

Static GTFS validation uses the pinned MobilityData static validator. GTFS-RT
validation uses the repo-supported pinned GTFS-RT validator wrapper.

If Java, Docker, the static validator JAR, or the GTFS-RT wrapper is missing,
the summary says `missing_tooling`. That is different from
`validation_failed`, which means the validator ran and returned a failure.

## Strict Mode

Use strict mode when missing tooling or failed validation should stop the
operator workflow:

```bash
scripts/validate-public-feeds.sh \
  --public-base-url https://feeds.example.org \
  --strict
```

Use `--skip-validators` only when the goal is public fetch/checksum review:

```bash
scripts/validate-public-feeds.sh \
  --public-base-url https://feeds.example.org \
  --skip-validators
```

## Output Files

The helper writes:

```text
summary.json
summary.md
manifest.json
manifest.md
artifacts/
headers/
validators/
```

Do not commit these files. They are ignored `.cache` diagnostics and stay
private unless a separate retained-evidence approval exists.
