#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

usage() {
  cat <<'EOF'
Usage:
  scripts/external-connection-check.sh [--help]

Checks connector manifests, fixtures, and local validation code.

Safety:
  This helper runs local checks only. It does not load dynamic plugins, start
  sidecars, contact consumers, automate portals, send notifications, mutate
  consumer statuses, write docs/evidence, or claim compliance, production
  readiness, agency approval, consumer acceptance, vendor compatibility, or
  production-grade ETA quality.
EOF
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
done

test -f docs/connectors/plugin-contract.md
test -f docs/external-connection-readiness.md
test -d internal/connectors
test -d testdata/connectors/valid
test -d testdata/connectors/invalid

for file in testdata/connectors/valid/*.json testdata/connectors/invalid/*.json; do
  python3 -m json.tool "$file" >/dev/null
done

go test ./internal/connectors

printf 'external connection check passed: connector manifests remain sidecar/manifest/conformance bounded\n'
