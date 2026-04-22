# Hosted Public Feed Root Proof

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Canonical HTTPS host: pending
- Status: missing

## Required Hosted Fetches

Record one anonymous HTTPS fetch for each path:

- [ ] `/public/gtfs/schedule.zip`
- [ ] `/public/feeds.json`
- [ ] `/public/gtfsrt/vehicle_positions.pb`
- [ ] `/public/gtfsrt/trip_updates.pb`
- [ ] `/public/gtfsrt/alerts.pb`

## Collection Commands

```sh
mkdir -p "$ENVIRONMENT_NAME/public"

for path in \
  /public/gtfs/schedule.zip \
  /public/feeds.json \
  /public/gtfsrt/vehicle_positions.pb \
  /public/gtfsrt/trip_updates.pb \
  /public/gtfsrt/alerts.pb
do
  name="$(basename "$path")"
  curl -sS -D "$ENVIRONMENT_NAME/public/$name.headers.txt" \
    -o "$ENVIRONMENT_NAME/public/$name" \
    -w "url=%{url_effective}\nstatus=%{http_code}\ncontent_type=%{content_type}\nsize_download=%{size_download}\ntime_total=%{time_total}\n" \
    "$PUBLIC_BASE_URL$path" | tee "$ENVIRONMENT_NAME/public/$name.curl.txt"
  shasum -a 256 "$ENVIRONMENT_NAME/public/$name" | tee "$ENVIRONMENT_NAME/public/$name.sha256.txt"
done
```

## Publish / Rollback URL Stability

- Before publish proof: pending.
- After publish proof: pending.
- After rollback proof: pending.
- URL changed? pending.

## Required Summary To Fill

| Path | Fetch timestamp UTC | Status | Bytes | SHA-256 | Header artifact |
| --- | --- | ---: | ---: | --- | --- |
| `/public/gtfs/schedule.zip` | pending | pending | pending | pending | pending |
| `/public/feeds.json` | pending | pending | pending | pending | pending |
| `/public/gtfsrt/vehicle_positions.pb` | pending | pending | pending | pending | pending |
| `/public/gtfsrt/trip_updates.pb` | pending | pending | pending | pending | pending |
| `/public/gtfsrt/alerts.pb` | pending | pending | pending | pending | pending |

## Blocker

No hosted HTTPS environment was available when this intake packet was created.
