#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
LOCK_FILE="$ROOT_DIR/tools/validators/validators.lock.json"
MODE="${VALIDATOR_TOOLING_MODE:-pinned}"

if [ "$MODE" = "stub" ]; then
  echo "validator tooling check: VALIDATOR_TOOLING_MODE=stub; pinned validators intentionally bypassed for deterministic test stubs"
  exit 0
fi
if [ "$MODE" != "pinned" ]; then
  echo "misconfigured pinned tooling: VALIDATOR_TOOLING_MODE must be pinned or stub" >&2
  exit 12
fi

json_value() {
  key="$1"
  sed -n "s/.*\"$key\"[[:space:]]*:[[:space:]]*\"\\([^\"]*\\)\".*/\\1/p" "$LOCK_FILE" | head -n 1
}

abs_path() {
  case "$1" in
    /*) printf '%s\n' "$1" ;;
    *) printf '%s\n' "$ROOT_DIR/$1" ;;
  esac
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

static_path="$(abs_path "$(json_value local_path)")"
static_sha="$(json_value sha256)"
rt_image="$(json_value image)"
rt_wrapper="$(abs_path "$(json_value wrapper_path)")"

if [ ! -f "$static_path" ]; then
  echo "missing pinned tooling: static GTFS validator not installed at $static_path; run make validators-install" >&2
  exit 11
fi
if ! command -v java >/dev/null 2>&1 || ! java -version >/dev/null 2>&1; then
  echo "missing pinned tooling: Java runtime is required for the static GTFS validator JAR" >&2
  exit 11
fi
actual_static_sha="$(sha256_file "$static_path")"
if [ "$actual_static_sha" != "$static_sha" ]; then
  echo "misconfigured pinned tooling: static GTFS validator checksum mismatch: got $actual_static_sha want $static_sha" >&2
  exit 12
fi
if [ -n "${GTFS_VALIDATOR_PATH:-}" ] && [ "$(abs_path "$GTFS_VALIDATOR_PATH")" != "$static_path" ]; then
  echo "misconfigured pinned tooling: GTFS_VALIDATOR_PATH must point to $static_path" >&2
  exit 12
fi

if ! command -v docker >/dev/null 2>&1; then
  echo "missing pinned tooling: docker is required for the repo-supported GTFS-RT validator wrapper" >&2
  exit 11
fi
if ! command -v curl >/dev/null 2>&1; then
  echo "missing pinned tooling: curl is required for the repo-supported GTFS-RT validator wrapper" >&2
  exit 11
fi
if ! command -v python3 >/dev/null 2>&1; then
  echo "missing pinned tooling: python3 is required for the repo-supported GTFS-RT validator wrapper" >&2
  exit 11
fi
if ! docker image inspect "$rt_image" >/dev/null 2>&1; then
  echo "missing pinned tooling: GTFS-RT validator image $rt_image is not installed; run make validators-install" >&2
  exit 11
fi
if [ ! -x "$rt_wrapper" ]; then
  echo "missing pinned tooling: GTFS-RT validator wrapper not installed at $rt_wrapper; run make validators-install" >&2
  exit 11
fi
if ! grep -F "$rt_image" "$rt_wrapper" >/dev/null 2>&1; then
  echo "misconfigured pinned tooling: GTFS-RT validator wrapper does not reference pinned image $rt_image" >&2
  exit 12
fi
if ! grep -F "/api/gtfs-rt-feed" "$rt_wrapper" >/dev/null 2>&1; then
  echo "misconfigured pinned tooling: GTFS-RT validator wrapper must drive the pinned webapp API" >&2
  exit 12
fi
if [ -n "${GTFS_RT_VALIDATOR_PATH:-}" ] && [ "$(abs_path "$GTFS_RT_VALIDATOR_PATH")" != "$rt_wrapper" ]; then
  echo "misconfigured pinned tooling: repo-supported GTFS_RT_VALIDATOR_PATH must point to $rt_wrapper; direct non-Docker executables are runtime-capable but not accepted by this pinned check" >&2
  exit 12
fi
if [ -n "${GTFS_RT_VALIDATOR_VERSION:-}" ] && [ "$GTFS_RT_VALIDATOR_VERSION" != "$rt_image" ]; then
  echo "misconfigured pinned tooling: GTFS_RT_VALIDATOR_VERSION must be $rt_image" >&2
  exit 12
fi

echo "validator tooling check passed: pinned static GTFS JAR and pinned Docker-backed GTFS-RT wrapper are installed"
