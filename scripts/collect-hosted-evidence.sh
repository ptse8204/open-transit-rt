#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'USAGE'
Collect Phase 12 hosted deployment evidence.

Required environment:
  ENVIRONMENT_NAME   short evidence folder name, e.g. pilot-agency-prod
  PUBLIC_BASE_URL    canonical HTTPS feed root, e.g. https://feeds.example.org

Optional environment:
  ADMIN_BASE_URL     admin/origin base URL; defaults to PUBLIC_BASE_URL
  ADMIN_TOKEN        bearer token for validation and scorecard export
  CAPTURE_DATE_UTC   YYYY-MM-DD; defaults to current UTC date
  OUTPUT_ROOT        evidence output root; defaults to docs/evidence/captured

This script does not collect deployment-owned monitoring screenshots, backup
job history, reverse proxy config, or scorecard scheduler exports. Add those
artifacts manually to the generated packet.
USAGE
}

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required tool: $1" >&2
    exit 1
  fi
}

require_env() {
  name="$1"
  value="$(eval "printf '%s' \"\${$name:-}\"")"
  if [ -z "$value" ]; then
    echo "missing required environment variable: $name" >&2
    usage >&2
    exit 2
  fi
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

http_code() {
  curl -sS -o /dev/null -w '%{http_code}' "$1" || true
}

redacted_admin_base_url() {
  if [ -n "${ADMIN_BASE_URL:-}" ]; then
    printf '%s\n' "$ADMIN_BASE_URL"
  else
    printf '%s\n' "$PUBLIC_BASE_URL"
  fi
}

need curl
need openssl
need sed
need awk
need date

require_env ENVIRONMENT_NAME
require_env PUBLIC_BASE_URL

case "$ENVIRONMENT_NAME" in
  *[!A-Za-z0-9._-]*)
    echo "ENVIRONMENT_NAME may contain only letters, digits, dot, underscore, and hyphen" >&2
    exit 2
    ;;
esac

case "$PUBLIC_BASE_URL" in
  https://*) ;;
  *)
    echo "PUBLIC_BASE_URL must start with https:// for hosted evidence" >&2
    exit 2
    ;;
esac

ADMIN_BASE_URL="$(redacted_admin_base_url)"
CAPTURE_DATE_UTC="${CAPTURE_DATE_UTC:-$(date -u '+%Y-%m-%d')}"
OUTPUT_ROOT="${OUTPUT_ROOT:-docs/evidence/captured}"
PACKET_DIR="$OUTPUT_ROOT/$ENVIRONMENT_NAME/$CAPTURE_DATE_UTC"
STARTED_AT="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"

mkdir -p "$PACKET_DIR/artifacts/public" \
  "$PACKET_DIR/artifacts/tls" \
  "$PACKET_DIR/artifacts/validation" \
  "$PACKET_DIR/artifacts/scorecard" \
  "$PACKET_DIR/artifacts/operator-supplied" \
  "$PACKET_DIR/artifacts/logs"

public_host="$(printf '%s' "$PUBLIC_BASE_URL" | sed 's#^https://##' | sed 's#/.*$##')"

cat >"$PACKET_DIR/README.md" <<EOF
# Phase 12 Hosted Evidence Packet: $ENVIRONMENT_NAME

- Environment: \`$ENVIRONMENT_NAME\`
- Capture date (UTC): $CAPTURE_DATE_UTC
- Capture started (UTC): $STARTED_AT
- Operator: pending operator attribution
- Canonical HTTPS host: \`$PUBLIC_BASE_URL\`
- Status: collected by \`scripts/collect-hosted-evidence.sh\`; review required

## Claim Boundary

This packet contains command outputs from a hosted evidence collection run. It is not proof of compliance until an operator reviews every artifact, confirms validator status, attaches deployment-owned monitoring/backup/proxy evidence, and updates the summaries.

## Required Operator Attachments

Add redacted deployment-owned files under \`artifacts/operator-supplied/\` for:

- reverse proxy or load balancer route config;
- certificate renewal evidence;
- monitoring dashboard export;
- alert rules and one alert lifecycle;
- production backup policy, job history, and restore transcript;
- scorecard scheduler/job definition and history.
EOF

cat >"$PACKET_DIR/public-feed-proof-$CAPTURE_DATE_UTC.md" <<EOF
# Hosted Public Feed Root Proof

- Environment: \`$ENVIRONMENT_NAME\`
- Capture date (UTC): $CAPTURE_DATE_UTC
- Operator: pending operator attribution
- Canonical HTTPS host: \`$PUBLIC_BASE_URL\`

## Anonymous Hosted Fetches

| Path | Fetch timestamp UTC | Status | Bytes | SHA-256 | Header artifact |
| --- | --- | ---: | ---: | --- | --- |
EOF

for path in \
  /public/gtfs/schedule.zip \
  /public/feeds.json \
  /public/gtfsrt/vehicle_positions.pb \
  /public/gtfsrt/trip_updates.pb \
  /public/gtfsrt/alerts.pb
do
  safe_name="$(printf '%s' "$path" | sed 's#^/##; s#[/ ]#_#g')"
  out="$PACKET_DIR/artifacts/public/$safe_name"
  headers="$PACKET_DIR/artifacts/public/$safe_name.headers.txt"
  curl_meta="$PACKET_DIR/artifacts/public/$safe_name.curl.txt"
  fetch_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
  status="$(
    curl -sS -D "$headers" -o "$out" \
      -w 'url=%{url_effective}\nstatus=%{http_code}\ncontent_type=%{content_type}\nsize_download=%{size_download}\ntime_total=%{time_total}\n' \
      "$PUBLIC_BASE_URL$path" | tee "$curl_meta" | awk -F= '/^status=/{print $2}'
  )"
  bytes="$(wc -c <"$out" | awk '{print $1}')"
  hash="$(sha256_file "$out")"
  printf '%s  %s\n' "$hash" "artifacts/public/$safe_name" >"$PACKET_DIR/artifacts/public/$safe_name.sha256.txt"
  printf '| `%s` | %s | %s | %s | `%s` | `%s` |\n' "$path" "$fetch_at" "$status" "$bytes" "$hash" "artifacts/public/$safe_name.headers.txt" >>"$PACKET_DIR/public-feed-proof-$CAPTURE_DATE_UTC.md"
done

cat >>"$PACKET_DIR/public-feed-proof-$CAPTURE_DATE_UTC.md" <<'EOF'

## Publish / Rollback URL Stability

- Before publish proof: pending operator attachment.
- After publish proof: pending operator attachment.
- After rollback proof: pending operator attachment.
- URL changed? pending operator review.
EOF

cat >"$PACKET_DIR/reverse-proxy-tls-$CAPTURE_DATE_UTC.md" <<EOF
# Hosted Reverse Proxy and TLS Evidence

- Environment: \`$ENVIRONMENT_NAME\`
- Capture date (UTC): $CAPTURE_DATE_UTC
- Operator: pending operator attribution
- Public host: \`$public_host\`

## TLS / Redirect Artifacts

- HTTPS headers: \`artifacts/tls/https-feeds-headers.txt\`
- HTTP redirect headers: \`artifacts/tls/http-redirect-headers.txt\`
- Certificate details: \`artifacts/tls/certificate.txt\`

## Operator-Supplied Evidence Still Required

- Redacted reverse proxy or load balancer routing map.
- Renewal mechanism and last renewal/check timestamp.
- Admin/debug network boundary evidence.
EOF

curl -sS -I "$PUBLIC_BASE_URL/public/feeds.json" >"$PACKET_DIR/artifacts/tls/https-feeds-headers.txt" || true
curl -sS -I "http://$public_host/public/feeds.json" >"$PACKET_DIR/artifacts/tls/http-redirect-headers.txt" || true
openssl s_client -connect "$public_host:443" -servername "$public_host" </dev/null 2>/dev/null \
  | openssl x509 -noout -issuer -subject -dates -ext subjectAltName \
  >"$PACKET_DIR/artifacts/tls/certificate.txt" || true

cat >"$PACKET_DIR/validator-record-$CAPTURE_DATE_UTC.md" <<EOF
# Hosted Validator Records

- Environment: \`$ENVIRONMENT_NAME\`
- Capture date (UTC): $CAPTURE_DATE_UTC
- Operator: pending operator attribution

## Validator Runs

| Feed type | Validator ID | Artifact | Status |
| --- | --- | --- | --- |
EOF

if [ -n "${ADMIN_TOKEN:-}" ]; then
  for feed_type in schedule vehicle_positions trip_updates alerts
  do
    if [ "$feed_type" = "schedule" ]; then
      validator_id="static-mobilitydata"
    else
      validator_id="realtime-mobilitydata"
    fi
    out="$PACKET_DIR/artifacts/validation/validate-$feed_type.json"
    status="$(
      curl -sS -X POST "$ADMIN_BASE_URL/admin/validation/run" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        --data "{\"validator_id\":\"$validator_id\",\"feed_type\":\"$feed_type\"}" \
        | tee "$out" \
        | sed -n 's/.*"status"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
        | head -n 1
    )"
    status="${status:-unknown}"
    printf '| `%s` | `%s` | `%s` | `%s` |\n' "$feed_type" "$validator_id" "artifacts/validation/validate-$feed_type.json" "$status" >>"$PACKET_DIR/validator-record-$CAPTURE_DATE_UTC.md"
  done
else
  cat >"$PACKET_DIR/artifacts/validation/MISSING_ADMIN_TOKEN.txt" <<'EOF'
ADMIN_TOKEN was not set, so hosted validator runs were not collected.
EOF
  for feed_type in schedule vehicle_positions trip_updates alerts
  do
    printf '| `%s` | pending | `artifacts/validation/MISSING_ADMIN_TOKEN.txt` | missing |\n' "$feed_type" >>"$PACKET_DIR/validator-record-$CAPTURE_DATE_UTC.md"
  done
fi

cat >"$PACKET_DIR/scorecard-export-$CAPTURE_DATE_UTC.md" <<EOF
# Hosted Scorecard Export Evidence

- Environment: \`$ENVIRONMENT_NAME\`
- Capture date (UTC): $CAPTURE_DATE_UTC
- Operator: pending operator attribution

## Manual Export

EOF

if [ -n "${ADMIN_TOKEN:-}" ]; then
  score_ts="$(date -u '+%Y-%m-%dT%H%M%SZ')"
  score_out="$PACKET_DIR/artifacts/scorecard/scorecard-$score_ts.json"
  curl -sS -X POST "$ADMIN_BASE_URL/admin/compliance/scorecard" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    --data '{}' >"$score_out"
  score_hash="$(sha256_file "$score_out")"
  printf '%s  %s\n' "$score_hash" "artifacts/scorecard/scorecard-$score_ts.json" >"$PACKET_DIR/artifacts/scorecard/scorecard-$score_ts.sha256.txt"
  cat >>"$PACKET_DIR/scorecard-export-$CAPTURE_DATE_UTC.md" <<EOF
- Export timestamp UTC: $score_ts
- Artifact: \`artifacts/scorecard/scorecard-$score_ts.json\`
- SHA-256: \`$score_hash\`
EOF
else
  cat >"$PACKET_DIR/artifacts/scorecard/MISSING_ADMIN_TOKEN.txt" <<'EOF'
ADMIN_TOKEN was not set, so hosted scorecard export was not collected.
EOF
  cat >>"$PACKET_DIR/scorecard-export-$CAPTURE_DATE_UTC.md" <<'EOF'
- Export status: missing because `ADMIN_TOKEN` was not set.
EOF
fi

cat >>"$PACKET_DIR/scorecard-export-$CAPTURE_DATE_UTC.md" <<'EOF'

## Scheduled Job Evidence

- Job/scheduler reference: pending operator attachment.
- Recent run history: pending operator attachment.
- Retention policy: pending operator attachment.
EOF

cat >"$PACKET_DIR/monitoring-alert-$CAPTURE_DATE_UTC.md" <<'EOF'
# Hosted Monitoring and Alerting Evidence

- Status: pending operator attachment.

## Required Attachments

- Monitoring dashboard export or screenshot.
- Alert rule definitions.
- Notification destination evidence.
- One real alert lifecycle with detected, acknowledged, mitigated, and resolved timestamps.

Generic repository commands cannot collect these deployment-owned artifacts.
EOF

cat >"$PACKET_DIR/backup-restore-drill-$CAPTURE_DATE_UTC.md" <<'EOF'
# Hosted Backup and Restore Evidence

- Status: pending operator attachment.

## Required Attachments

- Backup schedule and retention policy.
- Backup storage location and access boundary.
- Last successful backup job history.
- Restore drill transcript.
- Post-restore feed fetch and validator checks.
- Outage and validator-failure runbook links.

Generic repository commands cannot collect these deployment-owned artifacts.
EOF

cat >"$PACKET_DIR/operator-collection-commands-$CAPTURE_DATE_UTC.md" <<EOF
# Hosted Evidence Collection Commands

This packet was generated by:

\`\`\`sh
ENVIRONMENT_NAME="$ENVIRONMENT_NAME" PUBLIC_BASE_URL="$PUBLIC_BASE_URL" ADMIN_BASE_URL="<redacted-if-internal>" ./scripts/collect-hosted-evidence.sh
\`\`\`

ADMIN_TOKEN and exact ADMIN_BASE_URL are intentionally omitted from this command record.
EOF

find "$PACKET_DIR/artifacts" -type f | sort | while IFS= read -r file
do
  rel="${file#"$PACKET_DIR/"}"
  hash="$(sha256_file "$file")"
  printf '%s  %s\n' "$hash" "$rel"
done >"$PACKET_DIR/SHA256SUMS.txt"

cat <<EOF
Hosted evidence packet written to:
  $PACKET_DIR

Review the generated markdown summaries, attach deployment-owned monitoring,
backup, proxy renewal, and scheduler artifacts, then replace pending fields with
operator-reviewed facts.
EOF
