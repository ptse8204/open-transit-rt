#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/release-candidate-check/$TIMESTAMP}"
FORCE="${FORCE:-false}"
ALLOW_UNIGNORED_OUTPUT_DIR="${ALLOW_UNIGNORED_OUTPUT_DIR:-false}"
RUN_LOCAL_APP="${RUN_LOCAL_APP:-false}"
RUN_RELEASE_PACKAGE="${RUN_RELEASE_PACKAGE:-false}"
RELEASE_PACKAGE_DIR="${RELEASE_PACKAGE_DIR:-}"
DRY_RUN="false"
TMP_DIR=""

usage() {
  cat <<'EOF'
Usage:
  scripts/release-candidate-check.sh [--help] [--dry-run]

Environment:
  OUTPUT_DIR                    Default .cache/release-candidate-check/<UTC timestamp>
  FORCE                         true|false; allow non-empty output reuse
  ALLOW_UNIGNORED_OUTPUT_DIR    true|false; allow output outside .cache except evidence-like paths
  RUN_LOCAL_APP                 true|false; when true, start local app and fetch five public paths
  RUN_RELEASE_PACKAGE           true|false; when true, audit RELEASE_PACKAGE_DIR
  RELEASE_PACKAGE_DIR           Optional existing .cache/release-package/... directory to audit

Safety:
  This helper creates private local release-candidate diagnostics only. It
  writes exactly summary.json, summary.md, manifest.json, manifest.md, and
  check-log.txt. It never tags, publishes, pushes images, contacts consumers,
  creates retained evidence, writes docs/evidence, changes consumer statuses,
  or claims compliance, agency approval, hosted service availability,
  production readiness, vendor compatibility, or production-grade ETA quality.
EOF
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

cleanup() {
  if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT INT TERM

bool_var() {
  case "$2" in
    true|false) ;;
    *) fail "$1 must be true or false" ;;
  esac
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --dry-run)
      DRY_RUN="true"
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

bool_var FORCE "$FORCE"
bool_var ALLOW_UNIGNORED_OUTPUT_DIR "$ALLOW_UNIGNORED_OUTPUT_DIR"
bool_var RUN_LOCAL_APP "$RUN_LOCAL_APP"
bool_var RUN_RELEASE_PACKAGE "$RUN_RELEASE_PACKAGE"

mkdir -p .cache
TMP_DIR="$(mktemp -d ".cache/release-candidate-check-tmp.XXXXXX")"
chmod 700 "$TMP_DIR"

python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$TIMESTAMP" "$FORCE" "$ALLOW_UNIGNORED_OUTPUT_DIR" <<'PY' >"$TMP_DIR/output-dir"
import pathlib
import shutil
import sys

root = pathlib.Path(sys.argv[1]).resolve()
raw = pathlib.Path(sys.argv[2])
force = sys.argv[4] == "true"
allow = sys.argv[5] == "true"
out = raw if raw.is_absolute() else root / raw
resolved = out.resolve(strict=False)
cache = (root / ".cache").resolve(strict=False)

def evidence_like(path):
    text = str(path).replace("\\", "/").lower()
    parts = [p.lower() for p in pathlib.Path(path).parts]
    return "docs/evidence" in text or "evidence" in parts or "submission" in parts or "proof" in parts

def has_symlink(path):
    probe = pathlib.Path(path)
    if not probe.is_absolute():
        probe = root / probe
    current = pathlib.Path(probe.anchor) if probe.anchor else pathlib.Path(".")
    parts = probe.parts[1 if probe.anchor else 0:]
    for part in parts:
        current = current / part
        if current.exists() and current.is_symlink():
            return True
    return False

if evidence_like(raw) or evidence_like(resolved):
    raise SystemExit("OUTPUT_DIR must not be evidence-like or under docs/evidence")
if has_symlink(raw):
    raise SystemExit("OUTPUT_DIR must not contain symlink directories")
if not allow:
    try:
        resolved.relative_to(cache)
    except ValueError:
        raise SystemExit("OUTPUT_DIR must resolve under repo .cache unless ALLOW_UNIGNORED_OUTPUT_DIR=true")
if resolved.exists() and not resolved.is_dir():
    raise SystemExit("OUTPUT_DIR must be a directory")
if resolved.exists() and any(resolved.iterdir()):
    if not force:
        raise SystemExit("OUTPUT_DIR exists and is non-empty; use FORCE=true to reuse it")
    for child in resolved.iterdir():
        if child.is_symlink() or child.is_file():
            child.unlink()
        else:
            shutil.rmtree(child)
resolved.mkdir(parents=True, exist_ok=True)
resolved.chmod(0o700)
print(resolved)
PY
OUT_REAL="$(cat "$TMP_DIR/output-dir")"
CHECKS_TSV="$TMP_DIR/checks.tsv"
LOG_FILE="$OUT_REAL/check-log.txt"
printf 'id\tlabel\tstatus\tdetail\n' >"$CHECKS_TSV"
: >"$LOG_FILE"

add_check() {
  id="$1"
  label="$2"
  status="$3"
  detail="$4"
  case "$status" in
    passed|needs_review|not_checked|blocker) ;;
    *) fail "invalid check status: $status" ;;
  esac
  printf '%s\t%s\t%s\t%s\n' "$id" "$label" "$status" "$detail" >>"$CHECKS_TSV"
}

log_cmd() {
  printf '\n## %s\n' "$1" >>"$LOG_FILE"
  printf '$ %s\n' "$2" >>"$LOG_FILE"
}

run_check() {
  id="$1"
  label="$2"
  command="$3"
  if [ "$DRY_RUN" = "true" ]; then
    add_check "$id" "$label" "not_checked" "dry-run: $command"
    return
  fi
  log_cmd "$label" "$command"
  if sh -c "$command" >>"$LOG_FILE" 2>&1; then
    add_check "$id" "$label" "passed" "completed"
  else
    add_check "$id" "$label" "blocker" "command failed for $id; see check-log.txt"
  fi
}

run_check_final_blocker_detail() {
  id="$1"
  label="$2"
  command="$3"
  if [ "$DRY_RUN" = "true" ]; then
    add_check "$id" "$label" "not_checked" "dry-run: $command"
    return
  fi
  log_cmd "$label" "$command"
  command_output="$TMP_DIR/$id.out"
  if sh -c "$command" >"$command_output" 2>&1; then
    cat "$command_output" >>"$LOG_FILE"
    add_check "$id" "$label" "passed" "completed"
  else
    cat "$command_output" >>"$LOG_FILE"
    detail="$(awk 'NF { line = $0 } END { print line }' "$command_output" | tr '\t' ' ')"
    if [ -z "$detail" ]; then
      detail="command failed for $id; see check-log.txt"
    fi
    add_check "$id" "$label" "blocker" "$detail"
  fi
}

file_check() {
  id="$1"
  label="$2"
  path="$3"
  if [ -e "$path" ]; then
    add_check "$id" "$label" "passed" "$path present"
  else
    add_check "$id" "$label" "blocker" "$path missing"
  fi
}

for path in \
  go.mod \
  Makefile \
  deploy/docker-compose.yml \
  scripts/check-validators.sh \
  scripts/audit-final-claim-review.sh \
  scripts/audit-product-acceptance.sh \
  scripts/audit-product-language.sh \
  scripts/audit-ui-layout.sh \
  scripts/audit-operations-route-inventory.sh \
  scripts/api-contract-check.sh \
  scripts/check-internal-links.sh \
  scripts/check-stable-filter.sh \
  scripts/external-connection-check.sh \
  scripts/agency-local-app.sh \
  scripts/agency-pilot-onboard.sh \
  scripts/product-ui-smoke.sh \
  scripts/telemetry-simulator.sh \
  scripts/deployment-doctor.sh \
  scripts/validator-health.sh \
  scripts/operations-reliability.sh \
  scripts/operations-notify.sh \
  testdata/gtfs/valid-small/agency.txt \
  docs/release-candidate-readiness.md \
  docs/release-process.md \
  docs/release-checklist.md \
  docs/release-notes-template.md \
  docs/api-contracts.md
do
  file_check "file_$path" "Required file $path" "$path"
done

if command -v go >/dev/null 2>&1; then
  add_check "tool_go" "Go toolchain" "passed" "$(go version)"
else
  add_check "tool_go" "Go toolchain" "blocker" "go is not installed"
fi

if command -v docker >/dev/null 2>&1; then
  add_check "tool_docker" "Docker CLI" "passed" "docker CLI present"
else
  add_check "tool_docker" "Docker CLI" "needs_review" "docker CLI missing; local app and realtime validator wrapper may be unavailable"
fi

if command -v java >/dev/null 2>&1; then
  JAVA_VERSION="$(java -version 2>&1 | sed -n '1p')"
  add_check "tool_java" "Java runtime" "passed" "$JAVA_VERSION"
else
  add_check "tool_java" "Java runtime" "needs_review" "java missing; pinned MobilityData GTFS Validator execution may be unavailable"
fi

for tool in python3 git curl make; do
  if command -v "$tool" >/dev/null 2>&1; then
    add_check "tool_$tool" "$tool tool" "passed" "$tool present"
  else
    add_check "tool_$tool" "$tool tool" "needs_review" "$tool missing; review fresh-clone prerequisites"
  fi
done

if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  if [ -z "$(git status --short)" ]; then
    add_check "git_clean" "Git worktree state" "passed" "worktree clean"
  else
    add_check "git_clean" "Git worktree state" "needs_review" "worktree has uncommitted changes; acceptable for local diagnostics but review before tagging or publishing"
  fi
else
  add_check "git_clean" "Git worktree state" "needs_review" "not a git worktree"
fi

run_check "compose_config" "Docker Compose config" "docker compose -f deploy/docker-compose.yml config >/dev/null"
run_check_final_blocker_detail "validators_check" "Pinned validator install/check" "VALIDATOR_TOOLING_MODE=pinned scripts/check-validators.sh"
run_check "check_links" "Internal documentation link check" "scripts/check-internal-links.sh"
run_check "product_ui_smoke" "Product UI smoke" "scripts/product-ui-smoke.sh"
run_check "auth_production_boundary" "Production auth boundary, disabled local login, rotated JWT rejection, admin_session, CSRF, and public edge feed access" "go test ./cmd/agency-config -run 'TestProductionAuthBoundaryRegression$'"
run_check "auth_password_login" "Password login issues admin_session cookie without token leaks" "go test ./cmd/agency-config ./internal/auth -run 'Test(AdminPasswordLoginIssuesProductionSessionCookie|AdminPasswordLoginFailureIsGenericAndStateIsSingleUse|PasswordHashUsesArgon2IDAndDoesNotStorePlaintext|PasswordPolicyRejectsUnsafeValues)$'"
run_check "auth_bootstrap_single_use" "First-admin bootstrap token hash, TTL, one-time output, unsafe agency rejection, and browser setup" "go test ./cmd/agency-config ./internal/auth -run 'Test(FirstAdminSetupConsumesTokenSetsPasswordAndIssuesSession|BootstrapAdminLinkConfigRejectsUnsafeAgencyIDBeforeDB|BootstrapAdminLinkConfigNormalizesTTLAndBaseURL|BootstrapAdminLinkOutputShowsTokenOnceAndNotHash|BootstrapTokenHashIsHashOnlyAndStable|NormalizeBootstrapTTL)'"
run_check "auth_logout_expiry" "Logout requires cookie CSRF and expires admin_session" "go test ./cmd/agency-config -run 'TestAdminLogoutRequiresExistingCookieCSRFAndExpiresSession$'"
run_check "auth_cookie_post_csrf" "Cookie-authenticated unsafe POST without CSRF is rejected" "go test ./cmd/agency-config -run 'Test(ProductionAuthBoundaryRegression|OperationsCookiePostRequiresCSRF|OperationsSetupCookiePostRequiresCSRF|GTFSQualityCookiePostRequiresCSRF)$'"
run_check "dashboard_issue_priority" "Dashboard top-three issue priority and healthy fallback" "go test ./cmd/agency-config -run 'Test(TopDashboardIssuesSkipsHealthyAndCapsAtThree|DashboardHealthyFallbackFillsWhenFewerThanThreeIssues|OperationsDashboardFirstRunAcceptanceWorkflow)$'"
run_check "setup_wizard_skip_reminder" "Setup wizard skip path and session-scoped incomplete setup reminder" "go test ./cmd/agency-config -run 'Test(OperationsDashboardSetupReminderDismissalIsSessionScoped|SetupWizardRoutesPrivateScopedGETOnlyNoStore|SetupWizardJSONShapeFlagsAndStages|SetupWizardHTMLBoundariesNoFormsAndEscapes)$'"
run_check "connector_examples_vs_configured" "Connector examples remain separate from configured instances and explicit states" "go test ./cmd/agency-config ./internal/connectors -run 'Test(ConnectorHubSeparatesExamplesFromConfiguredInstances|ConnectorHubShowsConfiguredInstanceStateWithoutSecrets|ConnectorInstanceStatesAreExplicitAndStable|ConnectorInstanceConfigSummaryIsMetadataOnly)$'"
run_check "connector_dry_run_redaction" "Connector dry-run records redacted results and rejects raw payload signals" "go test ./cmd/agency-config ./internal/connectors -run 'Test(VehicleAVLDryRunPostStoresRedactedResult|VehicleAVLDryRunPostRejectsRawSummary|ConnectorDryRunJobValidationRejectsRawPayloadSignals)$'"
run_check "product_acceptance_audit" "Product acceptance audit" "scripts/audit-product-acceptance.sh"
run_check "product_language_audit" "Product language audit" "scripts/audit-product-language.sh"
run_check "ui_layout_audit" "UI layout audit" "scripts/audit-ui-layout.sh"
run_check "operations_route_inventory" "Operations route inventory audit" "scripts/audit-operations-route-inventory.sh"
run_check "api_contract_check" "API/feed/extension contract check" "scripts/api-contract-check.sh"
run_check "stable_filter_check" "Stable branch filter check" "scripts/check-stable-filter.sh --skip-ref-check"
run_check "external_connection_check" "External connection check" "scripts/external-connection-check.sh"
run_check "adapter_conformance" "Adapter conformance suite" "go run ./cmd/adapter-conformance run --suite testdata/adapter-conformance"
run_check "connector_examples" "Synthetic connector example tests" "go test ./examples/connectors/..."
run_check "gtfsrt_conformance" "GTFS-RT conformance harness" "go test ./internal/gtfsrtconformance ./cmd/gtfsrt-conformance"
run_check "claim_audit" "Final claim audit" "scripts/audit-final-claim-review.sh"
add_check "validate" "Repository validation command" "not_checked" "run make validate after reviewing this summary; this helper keeps repo output bounded to its five diagnostics files"
add_check "test" "Go unit tests command" "not_checked" "run make test after reviewing this summary; this helper keeps repo output bounded to its five diagnostics files"
add_check "smoke" "HTTP smoke command" "not_checked" "run make smoke after reviewing this summary; this helper keeps repo output bounded to its five diagnostics files"
run_check "gtfs_import_dry_run" "GTFS import/onboarding dry-run" "scripts/agency-pilot-onboard.sh --agency-id rc-dryrun --gtfs-url http://127.0.0.1/example.zip --dry-run >/dev/null"
run_check "telemetry_simulator" "Telemetry simulator dry-run" "OUTPUT_DIR='$TMP_DIR/telemetry-simulator' FORCE=true ALLOW_UNIGNORED_OUTPUT_DIR=true scripts/telemetry-simulator.sh --scenario on-route --dry-run --force --allow-unignored-output-dir >/dev/null"
run_check "deployment_doctor" "Deployment doctor" "OUTPUT_DIR='$TMP_DIR/deployment-doctor' FORCE=true ALLOW_UNIGNORED_OUTPUT_DIR=true scripts/deployment-doctor.sh >/dev/null"
run_check "validator_health" "Validator health dry-run" "OUTPUT_DIR='$TMP_DIR/validator-health' FORCE=true ALLOW_UNIGNORED_OUTPUT_DIR=true scripts/validator-health.sh --dry-run >/dev/null"
run_check "operations_notify" "Operations notification draft dry-run" "OUTPUT_DIR='$TMP_DIR/operations-notify' FORCE=true ALLOW_UNIGNORED_OUTPUT_DIR=true ALLOW_UNIGNORED_SOURCE_DIR=true VALIDATOR_HEALTH_SUMMARY='$TMP_DIR/missing-validator/summary.json' DEPLOYMENT_DOCTOR_SUMMARY='$TMP_DIR/missing-doctor/summary.json' scripts/operations-notify.sh --dry-run >/dev/null"
run_check "operations_reliability" "Operations reliability dry-run" "OUTPUT_DIR='$TMP_DIR/operations-reliability' FORCE=true ALLOW_UNIGNORED_OUTPUT_DIR=true ALLOW_UNIGNORED_SOURCE_DIR=true VALIDATOR_HEALTH_SUMMARY='$TMP_DIR/missing-validator/summary.json' DEPLOYMENT_DOCTOR_SUMMARY='$TMP_DIR/missing-doctor/summary.json' OPERATIONS_NOTIFY_SUMMARY='$TMP_DIR/missing-notify/summary.json' scripts/operations-reliability.sh --dry-run >/dev/null"

if [ "$RUN_LOCAL_APP" = "true" ]; then
  run_check "local_app_five_feeds" "Local app startup and five public feeds" "make agency-app-up >/dev/null && curl -fsS http://localhost:8080/public/feeds.json >/dev/null && curl -fsS http://localhost:8080/public/gtfs/schedule.zip >/dev/null && curl -fsS http://localhost:8080/public/gtfsrt/vehicle_positions.pb >/dev/null && curl -fsS http://localhost:8080/public/gtfsrt/trip_updates.pb >/dev/null && curl -fsS http://localhost:8080/public/gtfsrt/alerts.pb >/dev/null"
else
  add_check "local_app_five_feeds" "Local app startup and five public feeds" "not_checked" "set RUN_LOCAL_APP=true to run make agency-app-up and fetch /public/feeds.json, /public/gtfs/schedule.zip, /public/gtfsrt/vehicle_positions.pb, /public/gtfsrt/trip_updates.pb, and /public/gtfsrt/alerts.pb"
fi

if [ "$RUN_RELEASE_PACKAGE" = "true" ]; then
  if [ -n "$RELEASE_PACKAGE_DIR" ]; then
    run_check "release_package_audit" "Existing local release package audit" "RELEASE_PACKAGE_DIR='$RELEASE_PACKAGE_DIR' scripts/audit-release-package.sh >/dev/null"
  else
    add_check "release_package_audit" "Existing local release package audit" "needs_review" "RUN_RELEASE_PACKAGE=true requires RELEASE_PACKAGE_DIR pointing at an existing .cache release package; generation is outside this bounded check"
  fi
else
  add_check "release_package_audit" "Existing local release package audit" "not_checked" "set RUN_RELEASE_PACKAGE=true and RELEASE_PACKAGE_DIR=.cache/release-package/<version> to audit an existing local package"
fi

python3 - "$ROOT_DIR" "$OUT_REAL" "$CHECKS_TSV" "$TIMESTAMP" "$DRY_RUN" "$RUN_LOCAL_APP" "$RUN_RELEASE_PACKAGE" <<'PY'
import json
import pathlib
import subprocess
import sys

root = pathlib.Path(sys.argv[1]).resolve()
out = pathlib.Path(sys.argv[2])
checks_path = pathlib.Path(sys.argv[3])
timestamp = sys.argv[4]
dry_run = sys.argv[5] == "true"
run_local_app = sys.argv[6] == "true"
run_release_package = sys.argv[7] == "true"

def git_output(*args):
    try:
        result = subprocess.run(
            ["git", *args],
            cwd=root,
            check=True,
            capture_output=True,
            text=True,
        )
    except Exception:
        return "unavailable"
    value = result.stdout.strip()
    return value or "unavailable"

source = {
    "describe": git_output("describe", "--tags", "--always", "--dirty"),
    "commit_sha": git_output("rev-parse", "HEAD"),
    "branch": git_output("branch", "--show-current"),
}
source_status = git_output("status", "--short")
source["dirty"] = source_status not in ("", "unavailable")
source["pre_tag_review"] = True

rows = []
with checks_path.open() as fh:
    header = fh.readline()
    for line in fh:
        ident, label, status, detail = line.rstrip("\n").split("\t", 3)
        if ident != "validators_check":
            detail = detail[:260]
        rows.append({"id": ident, "label": label, "status": status, "detail": detail})

order = {"blocker": 4, "needs_review": 3, "not_checked": 2, "passed": 1}
overall = "passed"
for row in rows:
    if order[row["status"]] > order[overall]:
        overall = row["status"]

counts = {key: 0 for key in ("passed", "needs_review", "not_checked", "blocker")}
for row in rows:
    counts[row["status"]] += 1

public_feed_paths = [
    "/public/feeds.json",
    "/public/gtfs/schedule.zip",
    "/public/gtfsrt/vehicle_positions.pb",
    "/public/gtfsrt/trip_updates.pb",
    "/public/gtfsrt/alerts.pb",
]

review_sequence = [
    {
        "step": "clean_checkout",
        "command": "git status --short",
        "expected": "no uncommitted files before tagging or publishing release notes",
        "blocker_handling": "do not tag from a dirty checkout; use diagnostics only until the tree is clean",
    },
    {
        "step": "lightweight_repo_check",
        "command": "make check",
        "expected": "no-network/no-Docker/no-validator-install checks pass",
        "blocker_handling": "fix repository, JSON, script syntax, or claim-audit failures before continuing",
    },
    {
        "step": "release_candidate_diagnostic",
        "command": "make release-candidate-check",
        "expected": "private .cache summary covering product UI, links, connectors, GTFS-RT conformance, API contracts, and claim boundaries",
        "blocker_handling": "record exact blockers; do not convert the summary into production or compliance proof",
    },
    {
        "step": "product_quality_followup",
        "command": "make product-ui-smoke && make external-connection-check && make adapter-conformance && make gtfsrt-conformance && make api-contract-check",
        "expected": "the same product-quality gates are reproducible outside the release-candidate summary",
        "blocker_handling": "fix product, connector, feed, or contract failures before drafting stronger release notes",
    },
    {
        "step": "package_audit",
        "command": "make release-package && RELEASE_PACKAGE_DIR=.cache/release-package/<version> make audit-release-package",
        "expected": "local source package manifest, checksums, SBOM/provenance metadata, and audit result",
        "blocker_handling": "dirty or failed packages remain local diagnostics and are not release-ready artifacts",
    },
    {
        "step": "release_notes_draft",
        "command": "edit docs/release-notes-template.md copy",
        "expected": "notes list scope, migrations, operations/security/dependency changes, limitations, checks, and blockers",
        "blocker_handling": "state None for unchanged sections and list blocked commands without stronger claims",
    },
]

release_note_inputs = {
    "source": [
        "git tag or planned tag",
        "commit SHA",
        "dirty/clean state",
        "release-candidate diagnostic output directory",
    ],
    "package": [
        "local release package path",
        "source archive checksum",
        "SBOM/provenance status",
        "local image metadata status when supplied",
    ],
    "validation": [
        "make check",
        "make validate or exact validator blocker",
        "make test",
        "make test-release-package",
        "make check-links",
        "make product-ui-smoke",
        "auth/setup/connector release-gate rows inside make release-candidate-check",
        "make external-connection-check",
        "make adapter-conformance",
        "make gtfsrt-conformance",
        "make api-contract-check",
        "docker compose -f deploy/docker-compose.yml config",
        "make audit-final-claim-review",
    ],
    "claim_boundaries": [
        "no retained evidence created",
        "no consumer statuses changed",
        "no release published by this helper",
        "no tag created by this helper",
        "no image pushed by this helper",
    ],
}

package_audit_matrix = [
    {
        "item": "source archive",
        "check": "make release-package records archive name and SHA-256 checksum",
        "claim_boundary": "local package diagnostic only until a maintainer cuts a release",
    },
    {
        "item": "manifest files",
        "check": "make audit-release-package verifies required package files and JSON shape",
        "claim_boundary": "manifest presence is not hosted service or production-readiness proof",
    },
    {
        "item": "SBOM and provenance metadata",
        "check": "audit confirms metadata files exist and match the local package summary",
        "claim_boundary": "metadata is local supply-chain context, not external acceptance evidence",
    },
    {
        "item": "dirty state",
        "check": "package summary records whether the checkout was dirty",
        "claim_boundary": "dirty packages are diagnostics and not release-ready artifacts",
    },
    {
        "item": "local image metadata",
        "check": "optional metadata is recorded only when RELEASE_PACKAGE_IMAGE_TAG is supplied",
        "claim_boundary": "local image metadata does not mean a published production image exists",
    },
]

claim_flags = {
    "retained_evidence_created": False,
    "consumer_statuses_changed": False,
    "compliance_claimed": False,
    "production_readiness_claimed": False,
    "hosted_saas_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "vendor_compatibility_claimed": False,
    "production_grade_eta_claimed": False,
    "release_published": False,
    "tag_created": False,
    "image_pushed": False,
}

summary = {
    "schema_version": "open-transit-rt.release_candidate_check.v1",
    "generated_at": timestamp,
    "overall_status": overall,
    "dry_run": dry_run,
    "run_local_app": run_local_app,
    "run_release_package": run_release_package,
    "source": source,
    "gtfs_import_fixture": "testdata/gtfs/valid-small",
    "public_feed_paths": public_feed_paths,
    "counts": counts,
    "checks": rows,
    "review_sequence": review_sequence,
    "release_note_inputs": release_note_inputs,
    "package_audit_matrix": package_audit_matrix,
    "claim_flags": claim_flags,
    "boundary": "Private local diagnostics only; not evidence, not a release publication, not compliance proof, not consumer acceptance, not agency approval, not hosted service availability, not vendor compatibility, and not production readiness.",
}
manifest = {
    "schema_version": "open-transit-rt.release_candidate_manifest.v1",
    "generated_at": timestamp,
    "output_files": ["summary.json", "summary.md", "manifest.json", "manifest.md", "check-log.txt"],
    "default_output_root": ".cache/release-candidate-check",
    "workflow_document": "docs/release-candidate-readiness.md",
    "retained_evidence_created": False,
    "consumer_statuses_changed": False,
    "external_parties_contacted": False,
    "release_published": False,
    "tag_created": False,
    "image_pushed": False,
}

def write_json(name, value):
    (out / name).write_text(json.dumps(value, indent=2, sort_keys=True) + "\n")

write_json("summary.json", summary)
write_json("manifest.json", manifest)

lines = [
    "# Release-Candidate Readiness Summary",
    "",
    f"- Generated at: `{timestamp}`",
    f"- Overall status: `{overall}`",
    f"- Dry run: `{str(dry_run).lower()}`",
    f"- Local app check: `{'requested' if run_local_app else 'not_checked'}`",
    f"- Release package check: `{'requested' if run_release_package else 'not_checked'}`",
    f"- Source describe: `{source['describe']}`",
    f"- Commit SHA: `{source['commit_sha']}`",
    f"- Branch: `{source['branch']}`",
    f"- Dirty checkout: `{str(source['dirty']).lower()}`",
    "- Review mode: `pre-tag local diagnostics`",
    "- GTFS import fixture: `testdata/gtfs/valid-small`",
    "- Public feed paths: `/public/feeds.json`, `/public/gtfs/schedule.zip`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, `/public/gtfsrt/alerts.pb`",
    "",
    "This is private local diagnostics only. It is not retained evidence, not a release publication, not compliance proof, not consumer acceptance, not agency approval, not hosted service availability, not vendor compatibility, and not production readiness.",
    "",
    "## Checks",
]
for row in rows:
    lines.append(f"- `{row['status']}` {row['label']}: {row['detail']}")
lines.extend([
    "",
    "## First Release-Candidate Workflow",
])
for item in review_sequence:
    lines.append(
        f"- `{item['step']}`: run `{item['command']}`; expected {item['expected']}; blocker handling: {item['blocker_handling']}."
    )
lines.extend([
    "",
    "## Release Note Inputs",
])
for group, values in release_note_inputs.items():
    lines.append(f"- `{group}`: " + "; ".join(values))
lines.extend([
    "",
    "## Local Package Audit Matrix",
])
for item in package_audit_matrix:
    lines.append(f"- `{item['item']}`: {item['check']}; boundary: {item['claim_boundary']}.")
(out / "summary.md").write_text("\n".join(lines) + "\n")

manifest_md = [
    "# Release-Candidate Readiness Manifest",
    "",
    f"- Generated at: `{timestamp}`",
    "- Output root: private `.cache` diagnostics by default",
    "- Output files: `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, `check-log.txt`",
    "- Retained evidence created: `false`",
    "- Consumer statuses changed: `false`",
    "- Release published: `false`",
    "- Tag created: `false`",
    "- Image pushed: `false`",
]
(out / "manifest.md").write_text("\n".join(manifest_md) + "\n")

actual = sorted(p.name for p in out.iterdir() if p.is_file())
expected = sorted(manifest["output_files"])
if actual != expected:
    raise SystemExit(f"unexpected output files: {actual}")

print(f"release-candidate diagnostics written to {out.relative_to(root)}")
PY

python3 - "$OUT_REAL/summary.json" <<'PY'
import json
import pathlib
import sys

summary = json.loads(pathlib.Path(sys.argv[1]).read_text())
raise SystemExit(1 if summary.get("overall_status") == "blocker" else 0)
PY
