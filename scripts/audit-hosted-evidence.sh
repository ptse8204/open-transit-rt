#!/usr/bin/env sh
set -eu

usage() {
  cat <<'USAGE'
Audit a Phase 12 hosted evidence packet for completion blockers.

Required environment:
  EVIDENCE_PACKET_DIR   path to one evidence packet directory

Example:
  EVIDENCE_PACKET_DIR=docs/evidence/captured/pilot-agency-prod/2026-04-22 \
    ./scripts/audit-hosted-evidence.sh

The audit is intentionally conservative. It fails when placeholder/pending
markers remain, when validator outputs are missing or failed, or when required
artifact classes are absent.
USAGE
}

if [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
fi

if [ -z "${EVIDENCE_PACKET_DIR:-}" ]; then
  echo "missing required environment variable: EVIDENCE_PACKET_DIR" >&2
  usage >&2
  exit 2
fi

PACKET_DIR="${EVIDENCE_PACKET_DIR%/}"

if [ ! -d "$PACKET_DIR" ]; then
  echo "evidence packet directory does not exist: $PACKET_DIR" >&2
  exit 2
fi

failures=0

record_failure() {
  failures=$((failures + 1))
  printf 'FAIL: %s\n' "$1"
}

record_pass() {
  printf 'PASS: %s\n' "$1"
}

require_file() {
  file="$1"
  label="$2"
  if [ -f "$file" ]; then
    record_pass "$label exists"
  else
    record_failure "$label missing at $file"
  fi
}

require_glob() {
  pattern="$1"
  label="$2"
  # shellcheck disable=SC2086
  set -- $pattern
  if [ "$#" -gt 0 ] && [ -e "$1" ]; then
    record_pass "$label exists"
  else
    record_failure "$label missing for pattern $pattern"
  fi
}

require_glob "$PACKET_DIR/public-feed-proof-"*.md "public feed proof"
require_glob "$PACKET_DIR/reverse-proxy-tls-"*.md "reverse proxy/TLS proof"
require_glob "$PACKET_DIR/validator-record-"*.md "validator record"
require_glob "$PACKET_DIR/monitoring-alert-"*.md "monitoring/alert evidence"
require_glob "$PACKET_DIR/backup-restore-drill-"*.md "backup/restore evidence"
require_glob "$PACKET_DIR/scorecard-export-"*.md "scorecard export evidence"
require_file "$PACKET_DIR/SHA256SUMS.txt" "artifact checksum manifest"

for subdir in public tls validation scorecard operator-supplied
do
  if [ -d "$PACKET_DIR/artifacts/$subdir" ]; then
    record_pass "artifacts/$subdir directory exists"
  else
    record_failure "artifacts/$subdir directory missing"
  fi
done

if find "$PACKET_DIR" -type f -name '*.md' -exec grep -EIn '(^|[^A-Za-z])(pending|Pending|missing|Missing|not completed|not proof|Status: missing)' {} + >/tmp/open_transit_rt_evidence_pending.$$ 2>/dev/null; then
  record_failure "pending/missing markers remain in packet markdown"
  sed -n '1,80p' /tmp/open_transit_rt_evidence_pending.$$
else
  record_pass "no pending/missing markers found in packet markdown"
fi
rm -f /tmp/open_transit_rt_evidence_pending.$$

if [ -d "$PACKET_DIR/artifacts/validation" ]; then
  for feed_type in schedule vehicle_positions trip_updates alerts
  do
    file="$PACKET_DIR/artifacts/validation/validate-$feed_type.json"
    if [ ! -f "$file" ]; then
      record_failure "validator output missing for $feed_type"
      continue
    fi
    if grep -Eq '"status"[[:space:]]*:[[:space:]]*"passed"' "$file"; then
      record_pass "validator output passed for $feed_type"
    else
      record_failure "validator output is not passed for $feed_type"
    fi
  done
fi

for path in \
  public_gtfs_schedule.zip \
  public_feeds.json \
  public_gtfsrt_vehicle_positions.pb \
  public_gtfsrt_trip_updates.pb \
  public_gtfsrt_alerts.pb
do
  if [ -f "$PACKET_DIR/artifacts/public/$path" ]; then
    record_pass "public artifact exists: $path"
  else
    record_failure "public artifact missing: $path"
  fi
done

if [ -f "$PACKET_DIR/artifacts/tls/certificate.txt" ] && [ -s "$PACKET_DIR/artifacts/tls/certificate.txt" ]; then
  if grep -Eq 'notAfter=|issuer=|subject=' "$PACKET_DIR/artifacts/tls/certificate.txt"; then
    record_pass "TLS certificate details captured"
  else
    record_failure "TLS certificate artifact lacks expected certificate fields"
  fi
else
  record_failure "TLS certificate artifact missing or empty"
fi

if [ -f "$PACKET_DIR/artifacts/tls/http-redirect-headers.txt" ] && [ -s "$PACKET_DIR/artifacts/tls/http-redirect-headers.txt" ]; then
  if grep -Eq '^HTTP/.* 30[12378]' "$PACKET_DIR/artifacts/tls/http-redirect-headers.txt"; then
    record_pass "HTTP redirect evidence captured"
  else
    record_failure "HTTP redirect artifact does not show a 3xx redirect"
  fi
else
  record_failure "HTTP redirect artifact missing or empty"
fi

operator_required="reverse-proxy monitoring alert-lifecycle backup restore scorecard-job"
for label in $operator_required
do
  if [ -d "$PACKET_DIR/artifacts/operator-supplied" ] && find "$PACKET_DIR/artifacts/operator-supplied" -type f -iname "*$label*" | grep -q .; then
    record_pass "operator-supplied artifact found for $label"
  else
    record_failure "operator-supplied artifact missing for $label"
  fi
done

if [ "$failures" -gt 0 ]; then
  printf '\nHosted evidence audit failed with %d blocker(s).\n' "$failures" >&2
  exit 1
fi

printf '\nHosted evidence audit passed.\n'
