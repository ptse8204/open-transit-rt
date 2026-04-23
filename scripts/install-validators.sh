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
schedule=""
realtime=""
feed_type=""
output_dir=""

while [ "\$#" -gt 0 ]; do
  case "\$1" in
    --schedule)
      shift
      schedule="\${1:-}"
      ;;
    --realtime)
      shift
      realtime="\${1:-}"
      ;;
    --feed_type)
      shift
      feed_type="\${1:-}"
      ;;
    --output_dir|-o)
      shift
      output_dir="\${1:-}"
      ;;
    *)
      ;;
  esac
  shift || true
done

if [ -z "\$schedule" ] || [ -z "\$realtime" ] || [ -z "\$output_dir" ]; then
  echo "gtfs-rt validator wrapper requires --schedule, --realtime, and --output_dir" >&2
  exit 2
fi
if [ ! -s "\$schedule" ]; then
  echo "gtfs-rt validator wrapper schedule artifact is missing or empty: \$schedule" >&2
  exit 2
fi
if [ ! -s "\$realtime" ]; then
  echo "gtfs-rt validator wrapper realtime artifact is missing or empty: \$realtime" >&2
  exit 2
fi
if ! command -v python3 >/dev/null 2>&1; then
  echo "gtfs-rt validator wrapper requires python3 to serve local artifacts and normalize API output" >&2
  exit 2
fi
if ! command -v curl >/dev/null 2>&1; then
  echo "gtfs-rt validator wrapper requires curl" >&2
  exit 2
fi

mkdir -p "\$output_dir"
workdir="\$(dirname -- "\$schedule")"
schedule_name="\$(basename -- "\$schedule")"
realtime_name="\$(basename -- "\$realtime")"

free_port() {
  python3 - <<'PY'
import socket
with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
    sock.bind(("", 0))
    print(sock.getsockname()[1])
PY
}

file_port="\$(free_port)"
validator_port="\$(free_port)"
container_name="open-transit-rt-gtfsrt-validator-\$\$-\$validator_port"
client_id="open-transit-rt-wrapper-\$\$"

cleanup() {
  if [ -n "\${file_server_pid:-}" ]; then
    kill "\$file_server_pid" >/dev/null 2>&1 || true
    wait "\$file_server_pid" >/dev/null 2>&1 || true
  fi
  docker stop "\$container_name" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

python3 -m http.server "\$file_port" --bind 0.0.0.0 --directory "\$workdir" >"\$output_dir/artifact-http-server.log" 2>&1 &
file_server_pid="\$!"

docker run -d --rm --name "\$container_name" \
  --platform "\${GTFS_RT_VALIDATOR_DOCKER_PLATFORM:-linux/amd64}" \
  -p "127.0.0.1:\$validator_port:8080" \
  --add-host=host.docker.internal:host-gateway \
  "\$IMAGE" >"\$output_dir/validator-container-id.txt"

validator_base="http://127.0.0.1:\$validator_port"
for _ in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do
  if curl -fsS "\$validator_base/" >/dev/null 2>&1; then
    break
  fi
  sleep 1
done
if ! curl -fsS "\$validator_base/" >/dev/null 2>&1; then
  docker logs "\$container_name" >"\$output_dir/validator-container.log" 2>&1 || true
  echo "gtfs-rt validator webapp did not become ready" >&2
  exit 1
fi

schedule_url="http://host.docker.internal:\$file_port/\$schedule_name"
realtime_url="http://host.docker.internal:\$file_port/\$realtime_name"

gtfs_response="\$(curl -sS -H 'Accept: application/json' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode "gtfsurl=\$schedule_url" \
  --data 'enablevalidation=unchecked' \
  "\$validator_base/api/gtfs-feed")"
printf '%s\n' "\$gtfs_response" >"\$output_dir/gtfs-feed-response.json"
feed_id="\$(printf '%s' "\$gtfs_response" | sed -n 's/.*"feedId":[[:space:]]*\\([0-9][0-9]*\\).*/\\1/p')"
if [ -z "\$feed_id" ]; then
  docker logs "\$container_name" >"\$output_dir/validator-container.log" 2>&1 || true
  echo "gtfs-rt validator failed to load schedule feed" >&2
  exit 1
fi

rt_body="\$(printf '{"gtfsRtUrl":"%s","gtfsFeedModel":{"feedId":%s}}' "\$realtime_url" "\$feed_id")"
rt_response="\$(curl -sS -H 'Accept: */*' -H 'Content-Type: application/json' \
  --data "\$rt_body" \
  "\$validator_base/api/gtfs-rt-feed")"
printf '%s\n' "\$rt_response" >"\$output_dir/gtfs-rt-feed-response.json"
rt_id="\$(printf '%s' "\$rt_response" | sed -n 's/.*"gtfsRtId":[[:space:]]*\\([0-9][0-9]*\\).*/\\1/p')"
if [ -z "\$rt_id" ]; then
  docker logs "\$container_name" >"\$output_dir/validator-container.log" 2>&1 || true
  echo "gtfs-rt validator failed to load realtime feed" >&2
  exit 1
fi

monitor_response="\$(curl -sS -X PUT "\$validator_base/api/gtfs-rt-feed/monitor/\$rt_id?clientId=\$client_id&updateInterval=1&enableShapes=true")"
printf '%s\n' "\$monitor_response" >"\$output_dir/monitor-start-response.json"
sleep "\${GTFS_RT_VALIDATOR_MONITOR_SECONDS:-3}"

monitor_data="\$(curl -sS "\$validator_base/api/gtfs-rt-feed/monitor-data/\$rt_id?startTime=&summaryCurPage=1&summaryRowsPerPage=1000&toggledData=&logCurPage=1&logRowsPerPage=1000")"
printf '%s\n' "\$monitor_data" >"\$output_dir/monitor-data.json"
docker logs "\$container_name" >"\$output_dir/validator-container.log" 2>&1 || true

python3 - "\$feed_type" "\$monitor_data" <<'PY'
import json
import sys

feed_type = sys.argv[1]
raw = json.loads(sys.argv[2])
summary_items = raw.get("viewErrorSummaryModelList", [])
log_items = raw.get("viewErrorLogModelList", [])
items = [item for item in summary_items if isinstance(item, dict)]
if not items:
    items = [item for item in log_items if isinstance(item, dict)]

def item_count(item):
    try:
        return int(item.get("count", 1))
    except (TypeError, ValueError):
        return 1

errors = sum(item_count(item) for item in items if str(item.get("severity", "")).upper() in {"ERROR", "FATAL", "CRITICAL", "FAILURE", "FAILED"})
warnings = sum(item_count(item) for item in items if str(item.get("severity", "")).upper() in {"WARNING", "WARN"})
infos = sum(item_count(item) for item in items if str(item.get("severity", "")).upper() in {"INFO", "INFORMATIONAL", "NOTICE"})
status = "failed" if errors else "warning" if warnings else "passed"

print(json.dumps({
    "status": status,
    "feed_type": feed_type,
    "error_count": errors,
    "warning_count": warnings,
    "info_count": infos,
    "iteration_count": raw.get("iterationCount", 0),
    "unique_feed_count": raw.get("uniqueFeedCount", 0),
    "notices": items,
}, sort_keys=True))
PY
EOF
chmod 755 "$rt_wrapper"

cat <<EOF
Pinned validators installed.

Set these environment variables when running services or smoke checks:
  GTFS_VALIDATOR_PATH=$static_path
  GTFS_RT_VALIDATOR_PATH=$rt_wrapper
  GTFS_RT_VALIDATOR_VERSION=$rt_image

Use VALIDATOR_TOOLING_MODE=stub only for deterministic smoke/test stubs.
The GTFS-RT wrapper drives the pinned validator webapp through its local API;
it requires docker, curl, and python3 at runtime.
EOF
