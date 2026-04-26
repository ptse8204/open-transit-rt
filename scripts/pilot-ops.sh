#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'USAGE'
scripts/pilot-ops.sh <subcommand> [--dry-run] [--force]

Subcommands:
  validator-cycle   Run schedule, Vehicle Positions, Trip Updates, and Alerts validation.
  backup            Create a Postgres custom-format backup and cleanup old backups.
  restore-drill     Restore a backup into the configured target DB and verify public feeds.
  feed-monitor      Check public feed availability and record monitor evidence.
  scorecard-export  Export the compliance scorecard JSON through the admin API.

Required for every subcommand:
  ENVIRONMENT_NAME      explicit environment name, e.g. pilot-agency-prod
  EVIDENCE_OUTPUT_DIR   evidence output directory for this environment/date

Common optional:
  CAPTURE_DATE_UTC      YYYY-MM-DD; defaults to current UTC date

validator-cycle requires:
  ADMIN_BASE_URL        admin/origin URL, e.g. http://127.0.0.1:8081
  ADMIN_TOKEN           bearer token; never printed

backup requires:
  DATABASE_URL          source Postgres URL
  BACKUP_DIR            private backup destination
  BACKUP_RETENTION_DAYS optional retention cleanup window

restore-drill requires:
  RESTORE_DATABASE_URL  target Postgres URL; destructive
  RESTORE_BACKUP_FILE   backup dump to restore
  PUBLIC_BASE_URL       public HTTPS feed root for post-restore fetch checks

feed-monitor requires:
  PUBLIC_BASE_URL       public HTTPS feed root
  NOTIFY_WEBHOOK_URL    optional; do not commit real values
  NOTIFY_EMAIL_TO       optional; do not commit real values

scorecard-export requires:
  ADMIN_BASE_URL        admin/origin URL
  ADMIN_TOKEN           bearer token; never printed

Safety:
  All subcommands support --dry-run.
  State-changing operations require explicit target env vars; no deployment
  defaults are assumed. restore-drill requires typed confirmation unless
  --force is passed.
USAGE
}

die() {
  echo "ERROR: $*" >&2
  exit 2
}

need() {
  command -v "$1" >/dev/null 2>&1 || die "missing required tool: $1"
}

require_env() {
  name="$1"
  value="$(eval "printf '%s' \"\${$name:-}\"")"
  [ -n "$value" ] || die "missing required environment variable: $name"
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

capture_date() {
  if [ -n "${CAPTURE_DATE_UTC:-}" ]; then
    printf '%s' "$CAPTURE_DATE_UTC"
  else
    date -u '+%Y-%m-%d'
  fi
}

timestamp_utc() {
  date -u '+%Y-%m-%dT%H:%M:%SZ'
}

require_common() {
  require_env ENVIRONMENT_NAME
  require_env EVIDENCE_OUTPUT_DIR
  case "$ENVIRONMENT_NAME" in
    *[!A-Za-z0-9._-]*) die "ENVIRONMENT_NAME may contain only letters, digits, dot, underscore, and hyphen" ;;
  esac
}

print_target() {
  mode="live"
  [ "$DRY_RUN" = "dry-run" ] && mode="dry-run"
  echo "==> Target environment: ${ENVIRONMENT_NAME}"
  echo "==> Evidence output: ${EVIDENCE_OUTPUT_DIR}"
  echo "==> Mode: ${mode}"
}

ensure_output_dir() {
  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would create evidence directory: $EVIDENCE_OUTPUT_DIR"
  else
    mkdir -p "$EVIDENCE_OUTPUT_DIR"
  fi
}

public_paths() {
  cat <<'EOF'
/public/feeds.json
/public/gtfs/schedule.zip
/public/gtfsrt/vehicle_positions.pb
/public/gtfsrt/trip_updates.pb
/public/gtfsrt/alerts.pb
EOF
}

cmd_validator_cycle() {
  require_common
  require_env ADMIN_BASE_URL
  require_env ADMIN_TOKEN
  print_target
  ensure_output_dir
  out="$EVIDENCE_OUTPUT_DIR/validator-cycle-$(capture_date).json"
  echo "==> Validator cycle evidence: $out"

  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would POST /admin/validation/run for schedule, vehicle_positions, trip_updates, alerts"
    echo "DRY RUN: ADMIN_TOKEN is required but not printed"
    return 0
  fi

  need curl
  need sed
  need date

  tmp="${out}.tmp"
  printf '{\n  "environment":"%s",\n  "captured_at":"%s",\n  "records":[\n' "$ENVIRONMENT_NAME" "$(timestamp_utc)" >"$tmp"
  first=1
  for feed_type in schedule vehicle_positions trip_updates alerts
  do
    if [ "$feed_type" = "schedule" ]; then
      validator_id="static-mobilitydata"
    else
      validator_id="realtime-mobilitydata"
    fi
    body="$EVIDENCE_OUTPUT_DIR/validator-cycle-$(capture_date)-$feed_type.response.json"
    status="$(
      curl -fsS -X POST "$ADMIN_BASE_URL/admin/validation/run" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -H "Content-Type: application/json" \
        --data "{\"validator_id\":\"$validator_id\",\"feed_type\":\"$feed_type\"}" \
        | tee "$body" \
        | sed -n 's/.*"status"[[:space:]]*:[[:space:]]*"\([^"]*\)".*/\1/p' \
        | head -n 1
    )"
    status="${status:-unknown}"
    hash="$(sha256_file "$body")"
    [ "$first" -eq 1 ] || printf ',\n' >>"$tmp"
    first=0
    printf '    {"feed_type":"%s","validator_id":"%s","status":"%s","response_file":"%s","sha256":"%s"}' \
      "$feed_type" "$validator_id" "$status" "$(basename "$body")" "$hash" >>"$tmp"
  done
  printf '\n  ]\n}\n' >>"$tmp"
  mv "$tmp" "$out"
  echo "==> Validator cycle complete."
}

cmd_backup() {
  require_common
  require_env DATABASE_URL
  require_env BACKUP_DIR
  print_target
  ensure_output_dir
  ts="$(date -u '+%Y%m%dT%H%M%SZ')"
  backup="$BACKUP_DIR/open-transit-rt-$ts.dump"
  evidence="$EVIDENCE_OUTPUT_DIR/backup-run-$(capture_date).txt"
  echo "==> Backup destination: $backup"
  echo "==> Backup evidence: $evidence"

  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would run pg_dump -Fc to BACKUP_DIR and write backup evidence"
    if [ -n "${BACKUP_RETENTION_DAYS:-}" ]; then
      echo "DRY RUN: would delete backup dumps older than $BACKUP_RETENTION_DAYS day(s)"
    fi
    return 0
  fi

  need pg_dump
  need date

  mkdir -p "$BACKUP_DIR"
  pg_dump -Fc "$DATABASE_URL" -f "$backup"
  chmod 640 "$backup"
  hash="$(sha256_file "$backup")"
  printf '%s  %s\n' "$hash" "$backup" >"$backup.sha256"
  {
    echo "environment=$ENVIRONMENT_NAME"
    echo "timestamp_utc=$(timestamp_utc)"
    echo "backup_file=$backup"
    echo "sha256=$hash"
    echo "retention_days=${BACKUP_RETENTION_DAYS:-not-configured}"
  } >"$evidence"
  if [ -n "${BACKUP_RETENTION_DAYS:-}" ]; then
    find "$BACKUP_DIR" -name 'open-transit-rt-*.dump' -mtime +"$BACKUP_RETENTION_DAYS" -print -delete >>"$evidence"
  fi
  echo "==> Backup complete."
}

cmd_restore_drill() {
  require_common
  require_env RESTORE_DATABASE_URL
  require_env RESTORE_BACKUP_FILE
  require_env PUBLIC_BASE_URL
  print_target
  ensure_output_dir
  evidence="$EVIDENCE_OUTPUT_DIR/restore-drill-$(capture_date).txt"
  echo "==> Restore backup file: $RESTORE_BACKUP_FILE"
  echo "==> Restore evidence: $evidence"
  echo "WARNING: restore-drill is destructive for RESTORE_DATABASE_URL."

  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would restore RESTORE_BACKUP_FILE into RESTORE_DATABASE_URL"
    echo "DRY RUN: would fetch public feed paths from PUBLIC_BASE_URL after restore"
    return 0
  fi

  need pg_restore
  need curl
  need date

  [ -f "$RESTORE_BACKUP_FILE" ] || die "RESTORE_BACKUP_FILE does not exist: $RESTORE_BACKUP_FILE"
  if [ "$FORCE" != "force" ]; then
    printf 'Type "restore %s" to continue: ' "$ENVIRONMENT_NAME" >&2
    read answer
    [ "$answer" = "restore $ENVIRONMENT_NAME" ] || die "restore confirmation did not match; aborting"
  fi

  {
    echo "environment=$ENVIRONMENT_NAME"
    echo "started_at_utc=$(timestamp_utc)"
    echo "backup_file=$RESTORE_BACKUP_FILE"
    echo "restore_status=running"
  } >"$evidence"
  pg_restore --clean --if-exists -d "$RESTORE_DATABASE_URL" "$RESTORE_BACKUP_FILE" >>"$evidence" 2>&1
  echo "restore_finished_at_utc=$(timestamp_utc)" >>"$evidence"
  echo "post_restore_public_fetches:" >>"$evidence"
  public_paths | while IFS= read -r path
  do
    status="$(curl -sS -o /dev/null -w '%{http_code}' "$PUBLIC_BASE_URL$path" || true)"
    printf '%s %s%s\n' "$status" "$PUBLIC_BASE_URL" "$path" >>"$evidence"
  done
  echo "==> Restore drill complete."
}

cmd_feed_monitor() {
  require_common
  require_env PUBLIC_BASE_URL
  print_target
  ensure_output_dir
  evidence="$EVIDENCE_OUTPUT_DIR/feed-monitor-$(capture_date).txt"
  echo "==> Feed monitor evidence: $evidence"
  if [ -n "${NOTIFY_WEBHOOK_URL:-}" ] || [ -n "${NOTIFY_EMAIL_TO:-}" ]; then
    echo "==> Notification destination: configured (value not printed)"
  else
    echo "==> Notification destination: not configured"
  fi

  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would check public feed URLs and record status"
    echo "DRY RUN: missing notification destination would be reported as notification not configured, not feed failure"
    return 0
  fi

  need curl
  need date

  failures=0
  {
    echo "environment=$ENVIRONMENT_NAME"
    echo "timestamp_utc=$(timestamp_utc)"
    if [ -n "${NOTIFY_WEBHOOK_URL:-}" ] || [ -n "${NOTIFY_EMAIL_TO:-}" ]; then
      echo "notification=configured"
    else
      echo "notification=not configured"
    fi
  } >"$evidence"
  public_paths | while IFS= read -r path
  do
    status="$(curl -sS -o /dev/null -w '%{http_code}' --connect-timeout 5 --max-time 15 "$PUBLIC_BASE_URL$path" || true)"
    case "$status" in
      200) result=ok ;;
      *) result=failed; failures=$((failures + 1)) ;;
    esac
    printf '%s path=%s status_code=%s status=%s\n' "$(timestamp_utc)" "$path" "$status" "$result" >>"$evidence"
  done
  if grep -q ' status=failed' "$evidence"; then
    if [ -n "${NOTIFY_WEBHOOK_URL:-}" ] || [ -n "${NOTIFY_EMAIL_TO:-}" ]; then
      echo "notification_action=operator notification destination configured; delivery is deployment-owned" >>"$evidence"
    else
      echo "notification_action=notification not configured" >>"$evidence"
    fi
    echo "ERROR: one or more feed checks failed; see $evidence" >&2
    exit 1
  fi
  echo "==> Feed monitor checks passed."
}

cmd_scorecard_export() {
  require_common
  require_env ADMIN_BASE_URL
  require_env ADMIN_TOKEN
  print_target
  ensure_output_dir
  out="$EVIDENCE_OUTPUT_DIR/scorecard-export-$(capture_date).json"
  echo "==> Scorecard evidence: $out"

  if [ "$DRY_RUN" = "dry-run" ]; then
    echo "DRY RUN: would POST /admin/compliance/scorecard and write JSON evidence"
    echo "DRY RUN: ADMIN_TOKEN is required but not printed"
    return 0
  fi

  need curl
  need date

  curl -fsS -X POST "$ADMIN_BASE_URL/admin/compliance/scorecard" \
    -H "Authorization: Bearer $ADMIN_TOKEN" \
    -H "Content-Type: application/json" \
    --data '{}' >"$out"
  sha256_file "$out" >"$out.sha256"
  echo "==> Scorecard export complete."
}

SUBCOMMAND="${1:-help}"
shift || true
DRY_RUN=""
FORCE=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --dry-run) DRY_RUN="dry-run" ;;
    --force) FORCE="force" ;;
    -h|--help|help) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
  shift
done

case "$SUBCOMMAND" in
  validator-cycle) cmd_validator_cycle ;;
  backup) cmd_backup ;;
  restore-drill) cmd_restore_drill ;;
  feed-monitor) cmd_feed_monitor ;;
  scorecard-export) cmd_scorecard_export ;;
  help|-h|--help) usage ;;
  *) die "unknown subcommand: $SUBCOMMAND" ;;
esac
