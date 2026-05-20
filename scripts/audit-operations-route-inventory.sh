#!/usr/bin/env sh
set -eu

ROOT_DIR="${OPERATIONS_ROUTE_AUDIT_ROOT:-$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)}"
STRICT_DOCS="${OPERATIONS_ROUTE_AUDIT_STRICT_DOCS:-false}"

usage() {
  cat <<'EOF'
Usage:
  scripts/audit-operations-route-inventory.sh [--help]

Environment:
  OPERATIONS_ROUTE_AUDIT_ROOT          Optional repository root override.
  OPERATIONS_ROUTE_AUDIT_STRICT_DOCS   Set true to fail on stale README/wiki route maps.

Safety:
  Local read-only audit. It parses committed source and docs only. It does not
  start the app, fetch routes, call external URLs, run validators, contact
  consumers, write .cache outputs, write evidence, or mutate data.
EOF
}

case "${1:-}" in
  --help|-h)
    usage
    exit 0
    ;;
  "")
    ;;
  *)
    printf 'ERROR: unknown argument: %s\n' "$1" >&2
    exit 1
    ;;
esac

cd "$ROOT_DIR"

python3 - "$ROOT_DIR" "$STRICT_DOCS" <<'PY'
import pathlib
import re
import sys

ROOT = pathlib.Path(sys.argv[1]).resolve()
STRICT_DOCS = sys.argv[2].lower() in {"1", "true", "yes"}

DOC_ROUTE_MAP_FILES = [
    "README.md",
    "wiki/README.md",
    "wiki/small-agency-quick-start.md",
]

DOCS_SHOULD_INCLUDE = {
    "/admin/operations/gtfs-workbench",
    "/admin/operations/realtime",
    "/admin/operations/validation-center",
    "/admin/operations/connectors/workbench",
    "/admin/operations/prediction-lab",
    "/admin/operations/access",
    "/admin/operations/audit",
}

failures = []
warnings = []


def read(path):
    full = ROOT / path
    try:
        return full.read_text(encoding="utf-8")
    except FileNotFoundError:
        failures.append(f"missing required file: {path}")
        return ""


main_go = read("cmd/agency-config/main.go")
operations_go = read("cmd/agency-config/operations.go")
registry_go = read("cmd/agency-config/operations_route_registry.go")
operations_js = read("cmd/agency-config/operations_admin.js")
phase90 = read("docs/phase-90-control-plane-final-status.md")

mux_routes = set(re.findall(r'mux\.Handle(?:Func)?\("([^"]+)"', main_go))


def registry_block(name):
    match = re.search(rf"var {name} = \[\][A-Za-z0-9_]+\{{\n(.*?)\n\}}", registry_go, re.S)
    if not match:
        failures.append(f"missing route registry block: {name}")
        return ""
    return match.group(1)


def registry_rows(block, prefix):
    return [
        line.strip().rstrip(",")
        for line in block.splitlines()
        if line.strip().startswith(prefix)
    ]


def field(row, name):
    match = re.search(rf'{name}: "([^"]*)"', row)
    return match.group(1) if match else ""


group_labels = {}
for row in registry_rows(registry_block("operationsRouteGroupRegistry"), "{ID:"):
    group_id = field(row, "ID")
    label = field(row, "Label")
    if group_id and label:
        group_labels[group_id] = label

nav_items = []
HTML_ROUTES = {}
JSON_ROUTES = set()
EXTERNAL_ADMIN_SURFACES = {}
for row in registry_rows(registry_block("operationsRouteRegistry"), "{Section:"):
    section = field(row, "Section")
    path = field(row, "Path")
    json_path = field(row, "JSONPath")
    nav_label = field(row, "NavLabel")
    group_id = field(row, "GroupID")
    external = "ExternalAdminSurface: true" in row
    if not path:
        failures.append(f"registry row missing path for section: {section}")
        continue
    if nav_label:
        nav_items.append({
            "label": nav_label,
            "href": path,
            "section": section,
            "external": external,
        })
    if json_path:
        JSON_ROUTES.add(json_path)
    group_label = group_labels.get(group_id, group_id)
    if external:
        EXTERNAL_ADMIN_SURFACES[path] = group_label
    else:
        HTML_ROUTES[path] = group_label

COMMAND_ROUTES = {}
for row in registry_rows(registry_block("operationsCommandRouteRegistry"), "{Section:"):
    path = field(row, "Path")
    method = field(row, "Method")
    if path and method:
        COMMAND_ROUTES[path] = method

EXPECTED_COUNTS = {
    "html": 30,
    "json": 22,
    "command": 1,
    "external": 2,
}
if len(HTML_ROUTES) != EXPECTED_COUNTS["html"]:
    failures.append(f"canonical HTML route count = {len(HTML_ROUTES)}, want {EXPECTED_COUNTS['html']}")
if len(JSON_ROUTES) != EXPECTED_COUNTS["json"]:
    failures.append(f"canonical JSON route count = {len(JSON_ROUTES)}, want {EXPECTED_COUNTS['json']}")
if len(COMMAND_ROUTES) != EXPECTED_COUNTS["command"]:
    failures.append(f"command route count = {len(COMMAND_ROUTES)}, want {EXPECTED_COUNTS['command']}")
if len(EXTERNAL_ADMIN_SURFACES) != EXPECTED_COUNTS["external"]:
    failures.append(f"external admin surface count = {len(EXTERNAL_ADMIN_SURFACES)}, want {EXPECTED_COUNTS['external']}")

nav_routes = {item["href"] for item in nav_items}

switch_cases = set()
for line in operations_go.splitlines():
    stripped = line.strip()
    if stripped.startswith("case "):
        switch_cases.update(re.findall(r'"([^"]+)"', stripped))


def child(path):
    if path == "/admin/operations":
        return ""
    if path == "/admin/operations.json":
        return "operations.json"
    return path.removeprefix("/admin/operations/")


def handler_present(path):
    if path in mux_routes:
        if path in {"/admin/operations", "/admin/operations.json"}:
            return True
        if path.startswith("/admin/operations/"):
            return child(path) in switch_cases or path.endswith("/assets/operations.js")
        return True
    return "/admin/operations/" in mux_routes and child(path) in switch_cases


for route in HTML_ROUTES:
    if route not in nav_routes:
        failures.append(f"canonical HTML route missing from Operations nav: {route}")
    if not handler_present(route):
        failures.append(f"canonical HTML route lacks handler/switch coverage: {route}")
    if route not in phase90:
        failures.append(f"canonical HTML route missing from Phase 90 inventory: {route}")

for route in JSON_ROUTES:
    if not handler_present(route):
        failures.append(f"canonical JSON route lacks handler/switch coverage: {route}")
    if route not in phase90:
        failures.append(f"canonical JSON route missing from Phase 90 inventory: {route}")

for route, method in COMMAND_ROUTES.items():
    if route not in mux_routes:
        failures.append(f"command route missing explicit mux registration: {route}")
    if route not in phase90:
        failures.append(f"command route missing from Phase 90 inventory: {route}")
    if method == "POST" and route not in operations_js:
        failures.append(f"command route missing from progressive JS allowlist/use: {route}")

for route in EXTERNAL_ADMIN_SURFACES:
    matches = [item for item in nav_items if item["href"] == route]
    if not matches:
        failures.append(f"external admin surface missing from nav: {route}")
    elif not matches[0]["external"]:
        failures.append(f"external admin surface is not marked external: {route}")
    if route not in phase90:
        failures.append(f"external admin surface missing from Phase 90 inventory: {route}")

for route in nav_routes:
    if route.startswith("/public/"):
        failures.append(f"public route was added to private Operations nav: {route}")
    if route.startswith("/admin/operations"):
        known = route in HTML_ROUTES or route in JSON_ROUTES or route == "/admin/operations/assets/operations.js"
        if not known and not route.endswith(".json"):
            failures.append(f"unexpected Operations nav route outside canonical inventory: {route}")

if re.search(r'mux\.Handle(?:Func)?\("/public/admin', main_go):
    failures.append("public admin route registration detected")

for route in sorted(DOCS_SHOULD_INCLUDE):
    present_files = [path for path in DOC_ROUTE_MAP_FILES if route in read(path)]
    if not present_files:
        message = f"route map docs omit newer Operations route: {route}"
        if STRICT_DOCS:
            failures.append(message)
        else:
            warnings.append(message)

nested_json = sorted(
    route for route in JSON_ROUTES
    if route.startswith("/admin/operations/") and route.removeprefix("/admin/operations/").count("/") > 0
)
for route in nested_json:
    if route in operations_js and not re.search(r'\[a-z0-9-/\]\+\\\.json|connectors/', operations_js):
        warnings.append(f"nested JSON route is present but JS allowlist may reject it: {route}")

if failures:
    for item in failures:
        print(f"FAIL: {item}")
    for item in warnings:
        print(f"WARN: {item}")
    raise SystemExit(1)

print(f"PASS: {len(HTML_ROUTES)} canonical private HTML routes have nav and handler coverage")
print(f"PASS: {len(JSON_ROUTES)} canonical private JSON routes have handler coverage")
print(f"PASS: {len(COMMAND_ROUTES)} private command route is explicit and allowlisted")
print(f"PASS: {len(EXTERNAL_ADMIN_SURFACES)} external admin surfaces are marked in nav")
print("PASS: no public admin route registration detected")
if warnings:
    for item in warnings:
        print(f"WARN: {item}")
else:
    print("PASS: README/wiki route maps include newer center-style routes")
PY
