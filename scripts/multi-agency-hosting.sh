#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/multi-agency-hosting/$TIMESTAMP}"
FORCE="${FORCE:-false}"

usage() {
  cat <<'EOF'
Usage:
  scripts/multi-agency-hosting.sh [--help]

Environment:
  OUTPUT_DIR  Default .cache/multi-agency-hosting/<UTC timestamp>
  FORCE       true|false; allow reuse of a non-empty output directory

Safety:
  This helper creates private local multi-agency hosting diagnostics only. It
  writes exactly summary.json, summary.md, manifest.json, and manifest.md under
  ignored .cache output by default. It does not create retained evidence,
  contact consumers, change consumer tracker state, publish tenant backups,
  restore tenant data, or claim hosted SaaS, production multi-tenant hosting,
  production readiness, SLA/uptime, compliance, agency adoption, consumer
  acceptance, vendor compatibility, marketplace approval, or production-grade
  ETA quality.
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
  shift
done

case "$FORCE" in
  true|false) ;;
  *) fail "FORCE must be true or false" ;;
esac

python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$TIMESTAMP" "$FORCE" <<'PY'
import json
import os
import pathlib
import re
import shutil
import sys

root_arg, output_arg, timestamp, force_arg = sys.argv[1:5]

ROOT = pathlib.Path(root_arg).resolve()
OUTPUT_RAW = pathlib.Path(output_arg)
FORCE = force_arg == "true"

OUTPUT_FILES = ("summary.json", "summary.md", "manifest.json", "manifest.md")
REQUIRED_PUBLIC_ROUTES = (
    "/public/agencies/*/feeds.json",
    "/public/agencies/*/gtfs/schedule.zip",
    "/public/agencies/*/gtfsrt/vehicle_positions.pb",
    "/public/agencies/*/gtfsrt/trip_updates.pb",
    "/public/agencies/*/gtfsrt/alerts.pb",
)
FORBIDDEN_PUBLIC_EDGE_PATTERNS = (
    r"/admin/",
    r"/admin\*",
    r"/admin",
    r"/v1/telemetry",
    r"/v1/events",
    r"/metrics",
    r"/gtfs-studio",
    r"vehicle_positions\.json",
    r"trip_updates\.json",
    r"alerts\.json",
)
CLAIM_FLAGS = {
    "hosted_saas_claimed": False,
    "production_multi_tenant_hosting_claimed": False,
    "production_readiness_claimed": False,
    "sla_or_uptime_claimed": False,
    "compliance_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "vendor_compatibility_claimed": False,
    "marketplace_approval_claimed": False,
    "production_grade_eta_claimed": False,
    "retained_evidence_created": False,
    "consumer_statuses_changed": False,
}


def fail(message):
    raise SystemExit(f"ERROR: {message}")


def is_relative_to(path, base):
    try:
        pathlib.Path(path).resolve(strict=False).relative_to(pathlib.Path(base).resolve(strict=False))
        return True
    except ValueError:
        return False


def path_has_symlink(path):
    probe = pathlib.Path(path)
    if not probe.is_absolute():
        probe = ROOT / probe
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False


def is_evidence_like(path):
    raw = str(path).replace("\\", "/").lower()
    parts = [part.lower() for part in pathlib.Path(path).parts]
    return "docs/evidence" in raw or "evidence" in parts or "proof" in parts or "submission" in parts


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def resolve_output_dir():
    out = OUTPUT_RAW if OUTPUT_RAW.is_absolute() else ROOT / OUTPUT_RAW
    resolved = out.resolve(strict=False)
    cache = (ROOT / ".cache").resolve(strict=False)
    if not is_relative_to(resolved, cache):
        fail("OUTPUT_DIR must resolve under repo .cache")
    if is_evidence_like(out) or is_evidence_like(resolved):
        fail("OUTPUT_DIR must not be evidence-like or under docs/evidence")
    if path_has_symlink(out):
        fail("OUTPUT_DIR must not contain symlink directories")
    if resolved.exists() and not resolved.is_dir():
        fail("OUTPUT_DIR must be a directory")
    if resolved.exists() and any(resolved.iterdir()):
        if not FORCE:
            fail("OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it")
        for child in resolved.iterdir():
            if child.is_symlink() or child.is_file():
                child.unlink()
            else:
                shutil.rmtree(child)
    resolved.mkdir(parents=True, exist_ok=True)
    os.chmod(resolved, 0o700)
    return resolved


def write_json(path, data):
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def caddy_summary(path):
    text = (ROOT / path).read_text(encoding="utf-8")
    missing = [route for route in REQUIRED_PUBLIC_ROUTES if route not in text]
    if missing:
        status = "missing_required_routes"
    else:
        status = "required_routes_present"
    forbidden_matches = []
    if path == "deploy/oci/Caddyfile":
        for pattern in FORBIDDEN_PUBLIC_EDGE_PATTERNS:
            if re.search(pattern, text):
                forbidden_matches.append(pattern)
        if forbidden_matches:
            status = "forbidden_public_edge_route_present"
    return {
        "path": path,
        "status": status,
        "required_routes_present": not missing,
        "missing_routes": missing,
        "forbidden_public_edge_matches": forbidden_matches,
    }


output = resolve_output_dir()
caddy_files = [caddy_summary("deploy/Caddyfile.local"), caddy_summary("deploy/oci/Caddyfile")]
status = "passed" if all(item["status"] == "required_routes_present" for item in caddy_files) else "needs_review"
summary = {
    "schema": "multi-agency-hosting-diagnostic.v1",
    "generated_at": timestamp,
    "status": status,
    "output_directory": rel(output),
    "claim_flags": CLAIM_FLAGS,
    "public_route_contract": list(REQUIRED_PUBLIC_ROUTES),
    "single_agency_routes_preserved": [
        "/public/feeds.json",
        "/public/gtfs/schedule.zip",
        "/public/gtfsrt/vehicle_positions.pb",
        "/public/gtfsrt/trip_updates.pb",
        "/public/gtfsrt/alerts.pb",
    ],
    "caddy_files": caddy_files,
    "operations_model": {
        "preferred_isolation_order": [
            "separate deployment per agency",
            "separate database per agency",
            "shared database only with explicit tested route/export boundaries",
        ],
        "tenant_restore_into_shared_live_database": "blocked",
        "tenant_export_classification": "private redacted diagnostic only, not retained evidence and not a restorable backup",
    },
    "evidence_boundaries": {
        "docs_evidence_captured_written": False,
        "consumer_tracker_changed": False,
        "retained_evidence_created": False,
    },
}
manifest = {
    "schema": "multi-agency-hosting-manifest.v1",
    "generated_at": timestamp,
    "status": status,
    "files": list(OUTPUT_FILES),
    "claim_flags": CLAIM_FLAGS,
}

write_json(output / "summary.json", summary)
write_json(output / "manifest.json", manifest)
(output / "summary.md").write_text(
    "# Multi-Agency Hosting Diagnostic\n\n"
    f"- Status: {status}\n"
    "- Scope: repository route/proxy/operations boundary review only.\n"
    "- Tenant restore into a shared live database: blocked.\n"
    "- Claims: no hosted SaaS, production multi-tenant hosting, SLA/uptime, "
    "production readiness, compliance, agency adoption, consumer acceptance, "
    "vendor compatibility, marketplace approval, or production-grade ETA claim.\n",
    encoding="utf-8",
)
(output / "manifest.md").write_text(
    "# Multi-Agency Hosting Manifest\n\n"
    f"- Generated at: {timestamp}\n"
    f"- Output directory: {rel(output)}\n"
    "- Files: summary.json, summary.md, manifest.json, manifest.md\n"
    "- Retained evidence created: false\n",
    encoding="utf-8",
)

actual = sorted(child.name for child in output.iterdir())
if actual != sorted(OUTPUT_FILES):
    fail(f"unexpected output files: {actual}")
if status != "passed":
    fail("multi-agency hosting diagnostic needs review; inspect summary.json")
print(f"multi-agency hosting diagnostic: {rel(output)}")
PY
