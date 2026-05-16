#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
MODE="${INSTALL_CONFIDENCE_MODE:-clone}"
SOURCE="${INSTALL_CONFIDENCE_SOURCE:-$ROOT_DIR}"
REF="${INSTALL_CONFIDENCE_REF:-HEAD}"
OUTPUT_DIR="${INSTALL_CONFIDENCE_OUTPUT_DIR:-.cache/install-confidence/$TIMESTAMP}"
FORCE="${INSTALL_CONFIDENCE_FORCE:-false}"
RUN_LOCAL_APP="${INSTALL_CONFIDENCE_RUN_LOCAL_APP:-false}"
RUN_VALIDATE="${INSTALL_CONFIDENCE_RUN_VALIDATE:-false}"
RUN_TEST="${INSTALL_CONFIDENCE_RUN_TEST:-false}"

usage() {
  cat <<'EOF'
Usage:
  scripts/install-confidence.sh [--help]

Environment:
  INSTALL_CONFIDENCE_MODE          clone|archive; default clone.
  INSTALL_CONFIDENCE_SOURCE        Clone URL/path or archive path. Default current repo path for clone mode.
  INSTALL_CONFIDENCE_REF           Git ref to checkout in clone mode. Default HEAD.
  INSTALL_CONFIDENCE_OUTPUT_DIR    Output under .cache/install-confidence/<timestamp> by default.
  INSTALL_CONFIDENCE_FORCE         true|false; allow replacing an existing output dir.
  INSTALL_CONFIDENCE_RUN_LOCAL_APP true|false; run make agency-app-up and five local feed fetches.
  INSTALL_CONFIDENCE_RUN_VALIDATE  true|false; run make validate.
  INSTALL_CONFIDENCE_RUN_TEST      true|false; run make test.

Safety:
  Runs local install-confidence checks only. Raw logs and fetched local feed
  files stay under ignored .cache output. This helper does not publish a
  release, contact consumers, create retained evidence, change consumer
  statuses, or claim compliance, adoption, hosted service availability,
  production readiness, vendor compatibility, SLA/uptime, or ETA quality.
EOF
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

case "${1:-}" in
  "")
    ;;
  --help|-h)
    usage
    exit 0
    ;;
  *)
    usage >&2
    fail "unknown argument: $1"
    ;;
esac

bool_arg() {
  name="$1"
  value="$2"
  case "$value" in
    true|false) ;;
    *) fail "$name must be true or false" ;;
  esac
}

bool_arg INSTALL_CONFIDENCE_FORCE "$FORCE"
bool_arg INSTALL_CONFIDENCE_RUN_LOCAL_APP "$RUN_LOCAL_APP"
bool_arg INSTALL_CONFIDENCE_RUN_VALIDATE "$RUN_VALIDATE"
bool_arg INSTALL_CONFIDENCE_RUN_TEST "$RUN_TEST"

case "$MODE" in
  clone|archive) ;;
  *) fail "INSTALL_CONFIDENCE_MODE must be clone or archive" ;;
esac

abs_path() {
  case "$1" in
    /*) printf '%s\n' "$1" ;;
    *) printf '%s\n' "$ROOT_DIR/$1" ;;
  esac
}

is_evidence_like() {
  case "$1" in
    *docs/evidence*|*evidence*|*proof*|*submission*) return 0 ;;
    *) return 1 ;;
  esac
}

OUTPUT_ABS="$(abs_path "$OUTPUT_DIR")"
case "$OUTPUT_ABS" in
  "$ROOT_DIR"/.cache/install-confidence/*) ;;
  *) fail "INSTALL_CONFIDENCE_OUTPUT_DIR must be under .cache/install-confidence" ;;
esac
if is_evidence_like "$OUTPUT_ABS"; then
  fail "INSTALL_CONFIDENCE_OUTPUT_DIR must not be evidence-like"
fi
if [ -e "$OUTPUT_ABS" ]; then
  if [ "$FORCE" != "true" ]; then
    fail "output directory exists; set INSTALL_CONFIDENCE_FORCE=true to replace it"
  fi
  rm -rf "$OUTPUT_ABS"
fi
mkdir -p "$OUTPUT_ABS/logs" "$OUTPUT_ABS/artifacts" "$OUTPUT_ABS/work"

SUMMARY="$OUTPUT_ABS/summary.md"
STATUS_TSV="$OUTPUT_ABS/status.tsv"
: >"$STATUS_TSV"

status_line() {
  printf '%s\t%s\t%s\n' "$1" "$2" "$3" >>"$STATUS_TSV"
}

run_step() {
  step="$1"
  label="$2"
  command="$3"
  log="$OUTPUT_ABS/logs/$step.log"
  if (cd "$WORKTREE" && sh -c "$command") >"$log" 2>&1; then
    status_line "$step" "passed" "$label"
    return 0
  fi
  status_line "$step" "failed" "$label"
  return 1
}

cleanup_local_app() {
  if [ "${WORKTREE:-}" ] && [ -d "$WORKTREE" ]; then
    (cd "$WORKTREE" && make agency-app-down >/dev/null 2>&1) || true
  fi
}

trap cleanup_local_app EXIT

if [ "$MODE" = "clone" ]; then
  WORKTREE="$OUTPUT_ABS/work/open-transit-rt"
  if ! git clone "$SOURCE" "$WORKTREE" >"$OUTPUT_ABS/logs/git-clone.log" 2>&1; then
    status_line "git_clone" "failed" "clone source"
    fail "git clone failed; see $OUTPUT_ABS/logs/git-clone.log"
  fi
  status_line "git_clone" "passed" "clone source"
  if ! (cd "$WORKTREE" && git checkout "$REF") >"$OUTPUT_ABS/logs/git-checkout.log" 2>&1; then
    status_line "git_checkout" "failed" "checkout ref"
    fail "git checkout failed; see $OUTPUT_ABS/logs/git-checkout.log"
  fi
  status_line "git_checkout" "passed" "checkout ref"
elif [ "$MODE" = "archive" ]; then
  source_abs="$(abs_path "$SOURCE")"
  [ -f "$source_abs" ] || fail "archive source not found: $source_abs"
  if is_evidence_like "$source_abs"; then
    fail "archive source must not be evidence-like"
  fi
  if ! tar -tzf "$source_abs" >/dev/null 2>"$OUTPUT_ABS/logs/archive-list.log"; then
    status_line "archive_list" "failed" "list archive"
    fail "archive listing failed; see $OUTPUT_ABS/logs/archive-list.log"
  fi
  status_line "archive_list" "passed" "list archive"
  if ! tar -xzf "$source_abs" -C "$OUTPUT_ABS/work" >"$OUTPUT_ABS/logs/archive-extract.log" 2>&1; then
    status_line "archive_extract" "failed" "extract archive"
    fail "archive extraction failed; see $OUTPUT_ABS/logs/archive-extract.log"
  fi
  WORKTREE="$(find "$OUTPUT_ABS/work" -mindepth 1 -maxdepth 1 -type d | sort | sed -n '1p')"
  [ -n "$WORKTREE" ] || fail "archive did not produce a top-level directory"
  status_line "archive_extract" "passed" "extract archive"
fi

{
  printf 'mode=%s\n' "$MODE"
  printf 'source=%s\n' "$SOURCE"
  printf 'ref=%s\n' "$REF"
  printf 'worktree=%s\n' "$WORKTREE"
  printf 'generated_at=%s\n' "$TIMESTAMP"
  printf 'go_version=%s\n' "$(go version 2>/dev/null || printf 'unavailable')"
  printf 'docker_version=%s\n' "$(docker --version 2>/dev/null || printf 'unavailable')"
  printf 'compose_version=%s\n' "$(docker compose version 2>/dev/null || printf 'unavailable')"
} >"$OUTPUT_ABS/environment.txt"

if [ -d "$WORKTREE/.git" ]; then
  (cd "$WORKTREE" && git rev-parse HEAD >"$OUTPUT_ABS/commit.txt") || true
  (cd "$WORKTREE" && git describe --tags --always --dirty >"$OUTPUT_ABS/describe.txt") || true
fi

overall="passed"
run_step "make-check" "make check" "make check" || overall="failed"
run_step "bootstrap-check" "scripts/bootstrap-dev.sh --check" "scripts/bootstrap-dev.sh --check" || overall="failed"

if [ "$RUN_VALIDATE" = "true" ]; then
  run_step "make-validate" "make validate" "make validate" || overall="failed"
else
  status_line "make_validate" "not_checked" "INSTALL_CONFIDENCE_RUN_VALIDATE=false"
fi

if [ "$RUN_TEST" = "true" ]; then
  run_step "make-test" "make test" "make test" || overall="failed"
else
  status_line "make_test" "not_checked" "INSTALL_CONFIDENCE_RUN_TEST=false"
fi

if [ "$RUN_LOCAL_APP" = "true" ]; then
  run_step "agency-app-up" "make agency-app-up" "make agency-app-up" || overall="failed"
  for pair in \
    "feeds_json /public/feeds.json" \
    "schedule_zip /public/gtfs/schedule.zip" \
    "vehicle_positions_pb /public/gtfsrt/vehicle_positions.pb" \
    "trip_updates_pb /public/gtfsrt/trip_updates.pb" \
    "alerts_pb /public/gtfsrt/alerts.pb"
  do
    name="$(printf '%s' "$pair" | awk '{print $1}')"
    path="$(printf '%s' "$pair" | awk '{print $2}')"
    if curl -fsS "http://localhost:8080$path" -o "$OUTPUT_ABS/artifacts/$name" >"$OUTPUT_ABS/logs/fetch-$name.log" 2>&1; then
      status_line "fetch_$name" "passed" "$path"
    else
      status_line "fetch_$name" "failed" "$path"
      overall="failed"
    fi
  done
else
  status_line "local_app" "not_checked" "INSTALL_CONFIDENCE_RUN_LOCAL_APP=false"
fi

if command -v shasum >/dev/null 2>&1; then
  find "$OUTPUT_ABS/artifacts" -type f ! -name SHA256SUMS.txt -print | sort | while IFS= read -r file; do
    shasum -a 256 "$file"
  done >"$OUTPUT_ABS/artifacts/SHA256SUMS.txt"
fi

{
  printf '# Install Confidence Summary\n\n'
  printf '- Generated at: `%s`\n' "$TIMESTAMP"
  printf '- Mode: `%s`\n' "$MODE"
  printf '- Source: `%s`\n' "$SOURCE"
  printf '- Ref: `%s`\n' "$REF"
  printf '- Overall status: `%s`\n' "$overall"
  printf '- Run local app: `%s`\n' "$RUN_LOCAL_APP"
  printf '- Run validate: `%s`\n' "$RUN_VALIDATE"
  printf '- Run test: `%s`\n\n' "$RUN_TEST"
  printf 'This is a local install-confidence diagnostic only. It is not retained evidence, release publication, production readiness, compliance proof, consumer acceptance, agency approval, hosted service availability, vendor compatibility, SLA/uptime, or ETA-quality proof.\n\n'
  printf '## Steps\n\n'
  while IFS='	' read -r step status label; do
    printf -- '- `%s` %s: %s\n' "$status" "$step" "$label"
  done <"$STATUS_TSV"
} >"$SUMMARY"

printf 'install-confidence diagnostics written to %s\n' "$OUTPUT_ABS"
[ "$overall" = "passed" ]
