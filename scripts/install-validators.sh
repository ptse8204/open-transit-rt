#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
LOCK_FILE="$ROOT_DIR/tools/validators/validators.lock.json"

json_value() {
  key="$1"
  sed -n "s/.*\"$key\"[[:space:]]*:[[:space:]]*\"\\([^\"]*\\)\".*/\\1/p" "$LOCK_FILE" | head -n 1
}

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "missing required tool for validator install: $1" >&2
    exit 1
  fi
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

need curl
need docker

static_url="$(json_value url)"
static_sha="$(json_value sha256)"
static_path="$ROOT_DIR/$(json_value local_path)"
rt_image="$(json_value image)"
rt_wrapper="$ROOT_DIR/$(json_value wrapper_path)"

mkdir -p "$(dirname "$static_path")"

if [ ! -f "$static_path" ] || [ "$(sha256_file "$static_path")" != "$static_sha" ]; then
  tmp="$static_path.tmp"
  echo "Installing pinned static GTFS validator..."
  curl -L --fail --silent --show-error -o "$tmp" "$static_url"
  actual="$(sha256_file "$tmp")"
  if [ "$actual" != "$static_sha" ]; then
    rm -f "$tmp"
    echo "misconfigured pinned tooling: static GTFS validator checksum mismatch: got $actual want $static_sha" >&2
    exit 12
  fi
  mv "$tmp" "$static_path"
fi

echo "Pulling pinned GTFS-RT validator image..."
docker pull "$rt_image" >/dev/null

cat >"$rt_wrapper" <<EOF
#!/usr/bin/env sh
set -eu
IMAGE="$rt_image"
workdir=""
prev=""
rewritten=""
for arg in "\$@"; do
  case "\$prev" in
    --schedule|--realtime|--output_dir|-i|-o)
      if [ -z "\$workdir" ]; then
        workdir="\$(dirname -- "\$arg")"
      fi
      rewritten="\$rewritten \$(printf '%s\n' "\$arg" | sed "s#^\$workdir#/work#")"
      prev=""
      continue
      ;;
  esac
  rewritten="\$rewritten \$arg"
  case "\$arg" in
    --schedule|--realtime|--output_dir|-i|-o) prev="\$arg" ;;
    *) prev="" ;;
  esac
done
if [ -z "\$workdir" ]; then
  echo "gtfs-rt validator wrapper requires at least one file/output path argument" >&2
  exit 2
fi
# shellcheck disable=SC2086
exec docker run --rm -v "\$workdir:/work" "\$IMAGE" \$rewritten
EOF
chmod 755 "$rt_wrapper"

cat <<EOF
Pinned validators installed.

Set these environment variables when running services or smoke checks:
  GTFS_VALIDATOR_PATH=$static_path
  GTFS_RT_VALIDATOR_PATH=$rt_wrapper
  GTFS_RT_VALIDATOR_VERSION=$rt_image

Use VALIDATOR_TOOLING_MODE=stub only for deterministic smoke/test stubs.
EOF
