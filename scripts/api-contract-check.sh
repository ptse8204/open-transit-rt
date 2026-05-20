#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - <<'PY'
import json
import pathlib
import sys

root = pathlib.Path(".").resolve()
failures = []

def read(path):
    p = root / path
    if not p.is_file():
        failures.append(f"missing file: {path}")
        return ""
    return p.read_text(encoding="utf-8", errors="replace")

def require_contains(path, needle):
    text = read(path)
    if needle not in text:
        failures.append(f"{path} missing {needle!r}")

docs = read("docs/api-contracts.md")
for needle in [
    "POST /v1/telemetry",
    "/public/feeds.json",
    "/public/gtfs/schedule.zip",
    "/public/gtfsrt/vehicle_positions.pb",
    "/public/gtfsrt/trip_updates.pb",
    "/public/gtfsrt/alerts.pb",
    "open-transit-rt.connector.v1",
    "open-transit-rt.adapter_conformance.v1",
    "internal/prediction.Adapter",
    "release-candidate contracts",
]:
    if needle not in docs:
        failures.append(f"docs/api-contracts.md missing {needle!r}")

for path in [
    "/public/feeds.json",
    "/public/gtfs/schedule.zip",
    "/public/agencies/",
]:
    require_contains("cmd/agency-config/main.go", path)

for path, source in {
    "/v1/telemetry": "cmd/telemetry-ingest/main.go",
    "/public/gtfsrt/vehicle_positions.pb": "cmd/feed-vehicle-positions/main.go",
    "/public/gtfsrt/trip_updates.pb": "cmd/feed-trip-updates/main.go",
    "/public/gtfsrt/alerts.pb": "cmd/feed-alerts/main.go",
}.items():
    require_contains(source, path)

for path in [
    "/admin/operations.json",
    "/admin/operations/launchpad.json",
    "/admin/operations/checklist.json",
    "/admin/operations/feed-health.json",
    "/admin/operations/validation-center.json",
    "/admin/operations/readiness.json",
    "/admin/operations/telemetry-simulator.json",
    "/admin/operations/prediction-lab.json",
    "/admin/operations/connectors/tests.json",
    "/admin/operations/connectors/workbench.json",
    "/admin/operations/help.json",
    "/admin/operations/gtfs-workbench.json",
    "/admin/operations/validation-health.json",
    "/admin/operations/validation-health/refresh.json",
    "/admin/operations/reliability.json",
    "/admin/operations/maintenance.json",
]:
    require_contains("cmd/agency-config/main.go", path)
    if path not in docs:
        failures.append(f"docs/api-contracts.md missing admin companion route {path}")

prediction_model = read("internal/prediction/model.go")
for needle in [
    "type Adapter interface",
    "PredictTripUpdates(ctx context.Context, request Request) (Result, error)",
    "type Request struct",
    "ActiveFeedVersion   gtfs.FeedVersion",
    "Telemetry           []telemetry.StoredEvent",
    "Assignments         map[string]state.Assignment",
    "VehiclePositionsURL string",
    "type Result struct",
    "type Diagnostics struct",
]:
    if needle not in prediction_model:
        failures.append(f"internal/prediction/model.go missing {needle!r}")

try:
    suite = json.loads(read("testdata/adapter-conformance/suite.json"))
except json.JSONDecodeError as exc:
    failures.append(f"testdata/adapter-conformance/suite.json invalid JSON: {exc}")
    suite = {}
if suite.get("schema_version") != "open-transit-rt.adapter_conformance.v1":
    failures.append("adapter conformance suite schema_version drift")
if suite.get("synthetic_only") is not True:
    failures.append("adapter conformance suite must remain synthetic_only=true")
case_types = {case.get("type") for case in suite.get("cases", [])}
for required in {"telemetry", "prediction", "validator", "monitoring", "consumer_discovery"}:
    if required not in case_types:
        failures.append(f"adapter conformance suite missing {required} case type")

manifest_paths = sorted(root.glob("examples/connectors/*/connector.json"))
if not manifest_paths:
    failures.append("no example connector manifests found")
for manifest_path in manifest_paths:
    rel = manifest_path.relative_to(root).as_posix()
    try:
        manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    except json.JSONDecodeError as exc:
        failures.append(f"{rel} invalid JSON: {exc}")
        continue
    if manifest.get("schema_version") != "open-transit-rt.connector.v1":
        failures.append(f"{rel} schema_version drift")
    mode = manifest.get("mode", {})
    if mode.get("disabled_by_default") is not True:
        failures.append(f"{rel} must remain disabled_by_default=true")
    boundary = manifest.get("claim_boundary", {})
    not_claimed = " ".join(boundary.get("not_claimed", []))
    for phrase in ["vendor compatibility", "consumer acceptance", "compliance", "production readiness"]:
        if phrase not in not_claimed:
            failures.append(f"{rel} claim boundary missing {phrase!r}")

contract_refs = [
    "docs/api-contracts.md",
    "docs/connectors/plugin-contract.md",
    "docs/integration-adapter-kit.md",
    "docs/extension-governance.md",
]
for ref in contract_refs:
    if not (root / ref).is_file():
        failures.append(f"missing contract reference: {ref}")

if failures:
    print("api contract check failed:")
    for failure in failures:
        print(f"  - {failure}")
    sys.exit(1)
print("api contract check passed")
PY
