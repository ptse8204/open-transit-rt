#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

tmp_dir="$(mktemp -d "${TMPDIR:-/tmp}/otrt-install-confidence-test.XXXXXX")"
cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT

scripts/install-confidence.sh --help >/dev/null
sh -n scripts/install-confidence.sh

if INSTALL_CONFIDENCE_OUTPUT_DIR="$tmp_dir/outside-cache" scripts/install-confidence.sh >/dev/null 2>"$tmp_dir/outside.err"; then
  echo "expected output-dir guard to reject outside .cache" >&2
  exit 1
fi
if ! grep -F "INSTALL_CONFIDENCE_OUTPUT_DIR must be under .cache/install-confidence" "$tmp_dir/outside.err" >/dev/null 2>&1; then
  echo "expected output-dir guard message" >&2
  exit 1
fi

if INSTALL_CONFIDENCE_MODE=invalid scripts/install-confidence.sh >/dev/null 2>"$tmp_dir/mode.err"; then
  echo "expected invalid mode to fail" >&2
  exit 1
fi
if ! grep -F "INSTALL_CONFIDENCE_MODE must be clone or archive" "$tmp_dir/mode.err" >/dev/null 2>&1; then
  echo "expected invalid mode message" >&2
  exit 1
fi

if INSTALL_CONFIDENCE_MODE=archive INSTALL_CONFIDENCE_SOURCE="$tmp_dir/missing.tar.gz" INSTALL_CONFIDENCE_OUTPUT_DIR=.cache/install-confidence/test-missing INSTALL_CONFIDENCE_FORCE=true scripts/install-confidence.sh >/dev/null 2>"$tmp_dir/archive.err"; then
  echo "expected missing archive to fail" >&2
  exit 1
fi
if ! grep -F "archive source not found" "$tmp_dir/archive.err" >/dev/null 2>&1; then
  echo "expected missing archive message" >&2
  exit 1
fi

echo "install-confidence script tests passed"
