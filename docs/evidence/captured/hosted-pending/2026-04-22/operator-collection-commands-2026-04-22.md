# Hosted Evidence Collection Commands

- Environment: `hosted-pending`
- Capture date (UTC): 2026-04-22
- Operator: pending
- Status: command kit only

## Setup

```sh
export ENVIRONMENT_NAME="<hosted-environment-name>"
export PUBLIC_BASE_URL="https://<canonical-feed-host>"
export FEED_BASE_URL="$PUBLIC_BASE_URL/public"
export ADMIN_BASE_URL="https://<admin-or-origin-host-if-different>"
export ADMIN_TOKEN="<redacted-admin-token>"
mkdir -p "$ENVIRONMENT_NAME/public" \
  "$ENVIRONMENT_NAME/tls" \
  "$ENVIRONMENT_NAME/validation" \
  "$ENVIRONMENT_NAME/monitoring" \
  "$ENVIRONMENT_NAME/backup" \
  "$ENVIRONMENT_NAME/scorecard"
```

## Preferred Repo Collector

Use the script when possible:

```sh
ENVIRONMENT_NAME="$ENVIRONMENT_NAME" \
PUBLIC_BASE_URL="$PUBLIC_BASE_URL" \
ADMIN_BASE_URL="$ADMIN_BASE_URL" \
ADMIN_TOKEN="$ADMIN_TOKEN" \
./scripts/collect-hosted-evidence.sh
```

The command writes to `docs/evidence/captured/$ENVIRONMENT_NAME/<UTC-date>/`.

After attaching deployment-owned artifacts, run:

```sh
EVIDENCE_PACKET_DIR="docs/evidence/captured/$ENVIRONMENT_NAME/<UTC-date>" \
make audit-hosted-evidence
```

## Public Feed Fetches

```sh
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

## TLS / Redirect

```sh
host="$(printf '%s' "$PUBLIC_BASE_URL" | sed 's#^https://##')"
curl -sS -I "$PUBLIC_BASE_URL/public/feeds.json" | tee "$ENVIRONMENT_NAME/tls/https-feeds-headers.txt"
curl -sS -I "http://$host/public/feeds.json" | tee "$ENVIRONMENT_NAME/tls/http-redirect-headers.txt"
openssl s_client -connect "$host:443" -servername "$host" </dev/null 2>/dev/null \
  | openssl x509 -noout -issuer -subject -dates -ext subjectAltName \
  | tee "$ENVIRONMENT_NAME/tls/certificate.txt"
```

## Validators

```sh
for feed_type in schedule vehicle_positions trip_updates alerts
do
  if [ "$feed_type" = "schedule" ]; then
    validator_id="static-mobilitydata"
  else
    validator_id="realtime-mobilitydata"
  fi

  curl -sS -X POST "$ADMIN_BASE_URL/admin/validation/run" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    --data "{\"validator_id\":\"$validator_id\",\"feed_type\":\"$feed_type\"}" \
    | tee "$ENVIRONMENT_NAME/validation/validate-$feed_type.json"
done
```

## Scorecard

```sh
timestamp="$(date -u '+%Y-%m-%dT%H%M%SZ')"
curl -sS -X POST "$ADMIN_BASE_URL/admin/compliance/scorecard" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  --data '{}' \
  | tee "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.json"
shasum -a 256 "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.json" \
  | tee "$ENVIRONMENT_NAME/scorecard/scorecard-$timestamp.sha256.txt"
```

## Operator-Supplied Attachments

Add redacted deployment-owned exports that cannot be collected by generic repo commands:

- reverse proxy or load balancer route config;
- certificate renewal evidence;
- monitoring dashboard export;
- alert rules and incident lifecycle;
- backup policy, job history, restore transcript;
- scorecard scheduler/job definition and history.

## Redaction Rule

Do not commit private keys, tokens, database passwords, private endpoint inventories, or unredacted personal data.
