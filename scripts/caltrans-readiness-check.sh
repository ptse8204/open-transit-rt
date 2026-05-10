#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
OUTPUT_DIR="${OUTPUT_DIR:-.cache/caltrans-readiness-check/$TIMESTAMP}"
FORCE="${FORCE:-false}"
RUN_PUBLIC_FETCH="${RUN_PUBLIC_FETCH:-false}"
FETCH_TIMEOUT_SECONDS="${FETCH_TIMEOUT_SECONDS:-5}"
PUBLIC_BASE_URL="${PUBLIC_BASE_URL:-}"
FEEDS_JSON_PATH="${FEEDS_JSON_PATH:-}"
VALIDATOR_HEALTH_SUMMARY="${VALIDATOR_HEALTH_SUMMARY:-}"
OPERATIONS_RELIABILITY_SUMMARY="${OPERATIONS_RELIABILITY_SUMMARY:-}"
TRIP_ID_CONSISTENCY_SUMMARY="${TRIP_ID_CONSISTENCY_SUMMARY:-}"
CONSUMER_STATUS_PATH="${CONSUMER_STATUS_PATH:-docs/evidence/consumer-submissions/status.json}"
DRY_RUN="false"

usage() {
  cat <<'EOF'
Usage:
  scripts/caltrans-readiness-check.sh [--help] [--dry-run]

Environment:
  OUTPUT_DIR                       Default .cache/caltrans-readiness-check/<UTC timestamp>
  FORCE                            true|false; allow non-empty output reuse
  RUN_PUBLIC_FETCH                 true|false; opt-in bounded GET/HEAD fetch checks
  FETCH_TIMEOUT_SECONDS            Positive integer timeout for opt-in fetches
  PUBLIC_BASE_URL                  Optional final/root URL to derive feed paths
  FEEDS_JSON_PATH                  Optional safe local feeds.json summary
  VALIDATOR_HEALTH_SUMMARY         Optional safe .cache validator-health summary.json
  OPERATIONS_RELIABILITY_SUMMARY   Optional safe .cache operations-reliability summary.json
  TRIP_ID_CONSISTENCY_SUMMARY      Optional safe .cache trip consistency summary.json
  CONSUMER_STATUS_PATH             Default docs/evidence/consumer-submissions/status.json

Safety:
  This helper writes private local gap diagnostics only under .cache. It writes
  exactly summary.json, summary.md, manifest.json, manifest.md, and
  gap-review.txt. It does not refresh official requirements, create retained
  evidence, write docs/evidence, contact consumers, automate portals, change
  consumer statuses, or claim CAL-ITP/Caltrans compliance, consumer acceptance,
  agency adoption, final-root proof, production readiness, hosted SaaS, SLA,
  vendor compatibility, production AVL reliability, or production-grade ETA
  quality.
EOF
}

fail() {
  printf 'ERROR: %s\n' "$1" >&2
  exit 1
}

bool_var() {
  case "$2" in
    true|false) ;;
    *) fail "$1 must be true or false" ;;
  esac
}

positive_int() {
  case "$2" in
    ''|*[!0-9]*) fail "$1 must be a positive integer" ;;
    0) fail "$1 must be greater than zero" ;;
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
bool_var RUN_PUBLIC_FETCH "$RUN_PUBLIC_FETCH"
positive_int FETCH_TIMEOUT_SECONDS "$FETCH_TIMEOUT_SECONDS"

python3 - "$ROOT_DIR" "$OUTPUT_DIR" "$TIMESTAMP" "$FORCE" "$DRY_RUN" "$RUN_PUBLIC_FETCH" "$FETCH_TIMEOUT_SECONDS" "$PUBLIC_BASE_URL" "$FEEDS_JSON_PATH" "$VALIDATOR_HEALTH_SUMMARY" "$OPERATIONS_RELIABILITY_SUMMARY" "$TRIP_ID_CONSISTENCY_SUMMARY" "$CONSUMER_STATUS_PATH" <<'PY'
import json
import os
import pathlib
import re
import shutil
import sys
import urllib.error
import urllib.parse
import urllib.request

(
    root_arg,
    output_arg,
    timestamp,
    force_arg,
    dry_run_arg,
    run_fetch_arg,
    fetch_timeout_arg,
    public_base_url,
    feeds_json_arg,
    validator_arg,
    reliability_arg,
    trip_consistency_arg,
    consumer_status_arg,
) = sys.argv[1:14]

ROOT = pathlib.Path(root_arg).resolve()
OUTPUT_RAW = pathlib.Path(output_arg)
FORCE = force_arg == "true"
DRY_RUN = dry_run_arg == "true"
RUN_FETCH = run_fetch_arg == "true"
FETCH_TIMEOUT = int(fetch_timeout_arg)

OUTPUT_FILES = ("summary.json", "summary.md", "manifest.json", "manifest.md", "gap-review.txt")
ALLOWED_STATUSES = {"present", "partial", "missing", "not_checked", "needs_review", "blocked"}
BAD_STATUS_WORDS = {"ok", "passed", "compliant", "certified", "accepted", "ingested", "listed", "displayed"}
EXPECTED_CONSUMERS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
FEED_PATHS = {
    "feeds_json": "/public/feeds.json",
    "schedule": "/public/gtfs/schedule.zip",
    "vehicle_positions": "/public/gtfsrt/vehicle_positions.pb",
    "trip_updates": "/public/gtfsrt/trip_updates.pb",
    "alerts": "/public/gtfsrt/alerts.pb",
}
CLAIM_FLAGS = {
    "caltrans_compliance_claimed": False,
    "consumer_acceptance_claimed": False,
    "consumer_ingestion_claimed": False,
    "agency_adoption_claimed": False,
    "final_root_proven": False,
    "production_readiness_claimed": False,
    "sla_or_uptime_claimed": False,
    "hosted_saas_claimed": False,
    "vendor_compatibility_claimed": False,
    "production_avl_reliability_claimed": False,
    "production_grade_eta_claimed": False,
}
OFFICIAL_CONTEXT = {
    "requirements_refresh_date": "2026-05-09",
    "refreshed_by_this_script": False,
    "caltrans_guidelines": "Version 4.0 dated December 11, 2024, recorded by Phase 54",
    "caltrans_faq": "Version 4.0, recorded by Phase 54",
    "fta_ntd_manual": "2025 NTD Reporting Policy Manual page last updated April 15, 2026, recorded by Phase 54",
    "boundary": "Requirements context only; this script does not verify current official sources and does not prove compliance.",
}


def fail(message):
    raise SystemExit(message)


def evidence_like(path):
    text = str(path).replace("\\", "/").lower()
    parts = [p.lower() for p in pathlib.Path(path).parts]
    return "docs/evidence" in text or "evidence" in parts or "proof" in parts or "submission" in parts


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


def rel_to_root(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def resolve_output_dir():
    out = OUTPUT_RAW if OUTPUT_RAW.is_absolute() else ROOT / OUTPUT_RAW
    resolved = out.resolve(strict=False)
    cache = (ROOT / ".cache").resolve(strict=False)
    if evidence_like(OUTPUT_RAW) or evidence_like(resolved):
        fail("OUTPUT_DIR must not be evidence-like or under docs/evidence")
    if path_has_symlink(OUTPUT_RAW):
        fail("OUTPUT_DIR must not contain symlink directories")
    try:
        resolved.relative_to(cache)
    except ValueError:
        fail("OUTPUT_DIR must resolve under repo .cache")
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


def unsafe_text(text):
    lower = text.lower()
    patterns = [
        "authorization:",
        "bearer ",
        "cookie:",
        "admin_session",
        "postgres://",
        "database_url",
        "password",
        "secret=",
        "secret:",
        "\"secret\"",
        "token_hash",
        "private_key",
        "begin private key",
        "payload_json",
        "raw_log",
        "webhook_url",
        "https://hooks.",
        "/users/",
        "/etc/",
        "/var/lib/",
    ]
    return any(pattern in lower for pattern in patterns)


def sanitize(value, default="unknown", limit=260):
    if value is None or isinstance(value, (dict, list)):
        return default
    text = str(value).replace("\r", " ").replace("\n", " ").strip()
    text = re.sub(r"\s+", " ", text)
    if not text:
        return default
    if unsafe_text(text):
        return "<redacted>"
    if len(text) > limit:
        return text[: limit - 15].rstrip() + " [truncated]"
    return text


def source_path_allowed(path, allow_consumer_status=False):
    resolved = path.resolve(strict=False)
    if path_has_symlink(path):
        return False
    if allow_consumer_status:
        expected = (ROOT / "docs/evidence/consumer-submissions/status.json").resolve(strict=False)
        return resolved == expected
    if evidence_like(path) or evidence_like(resolved):
        return False
    cache = (ROOT / ".cache").resolve(strict=False)
    try:
        resolved.relative_to(cache)
        return True
    except ValueError:
        return False


def read_json_source(kind, raw_arg, required=False, allow_consumer_status=False):
    if not raw_arg:
        return None, {
            "kind": kind,
            "status": "missing" if required else "not_checked",
            "source": "not_configured",
            "detail": "No source path configured.",
        }
    raw = pathlib.Path(raw_arg)
    path = raw if raw.is_absolute() else ROOT / raw
    if not source_path_allowed(path, allow_consumer_status=allow_consumer_status):
        fail(f"{kind} source path is not allowed")
    if not path.exists() or not path.is_file() or path.is_symlink():
        return None, {
            "kind": kind,
            "status": "missing",
            "source": rel_to_root(path),
            "detail": "Configured source is absent.",
        }
    if path.stat().st_size > 1024 * 1024:
        fail(f"{kind} source exceeds 1 MiB")
    text = path.read_text(encoding="utf-8", errors="replace")
    if unsafe_text(text):
        fail(f"{kind} source contains private or unsafe values")
    try:
        data = json.loads(text)
    except json.JSONDecodeError:
        return None, {
            "kind": kind,
            "status": "blocked",
            "source": rel_to_root(path),
            "detail": "Configured source is not valid JSON.",
        }
    return data, {
        "kind": kind,
        "status": "present",
        "source": rel_to_root(path),
        "detail": "Safe JSON source loaded.",
    }


def normalize_base_url(value):
    text = sanitize(value, "", 2048).rstrip("/")
    if not text:
        return ""
    parsed = urllib.parse.urlparse(text)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        return ""
    return text


def join_url(base, path):
    return base.rstrip("/") + path


def feed_urls_from_inputs(feeds_data):
    urls = {}
    source = "not_configured"
    base = normalize_base_url(public_base_url)
    if base:
        source = "PUBLIC_BASE_URL"
        for key, path in FEED_PATHS.items():
            urls[key] = join_url(base, path)
    if isinstance(feeds_data, dict):
        source = "FEEDS_JSON_PATH" if feeds_json_arg else source
        public_base = normalize_base_url(feeds_data.get("public_base_url"))
        if public_base:
            urls["feeds_json"] = join_url(public_base, FEED_PATHS["feeds_json"])
        for feed in feeds_data.get("feeds", []) if isinstance(feeds_data.get("feeds"), list) else []:
            if not isinstance(feed, dict):
                continue
            feed_type = sanitize(feed.get("feed_type"), "")
            url = normalize_base_url(feed.get("canonical_public_url"))
            if feed_type in {"schedule", "vehicle_positions", "trip_updates", "alerts"} and url:
                urls[feed_type] = url
    return urls, source


def row(area, status, signal, source, missing_or_review_needed, claim_boundary):
    if status not in ALLOWED_STATUSES:
        fail(f"invalid readiness row status {status!r} for {area}")
    if status in BAD_STATUS_WORDS:
        fail(f"claim-upgrading row status {status!r} for {area}")
    return {
        "area": area,
        "status": status,
        "signal": sanitize(signal),
        "source": sanitize(source),
        "missing_or_review_needed": sanitize(missing_or_review_needed),
        "claim_boundary": sanitize(claim_boundary),
    }


def url_status(urls, keys):
    present = [key for key in keys if urls.get(key)]
    if len(present) == len(keys):
        return "present"
    if present:
        return "partial"
    return "missing"


def feed_row(area, key, label, urls, source):
    status = "present" if urls.get(key) else "missing"
    signal = urls.get(key, f"{label} URL is missing")
    return row(
        area,
        status,
        signal,
        source,
        "Retain a stable public URL before using this row for any external packet." if status != "present" else "Review URL ownership and final-root evidence separately.",
        f"{label} availability is a readiness signal only, not compliance, consumer acceptance, or correctness proof.",
    )


def map_summary_status(value):
    raw = sanitize(value, "unknown").lower()
    if raw in {"recorded", "runnable", "configured", "installed", "stub", "info", "green", "yellow"}:
        return "present"
    if raw in {"passed"}:
        return "present"
    if raw in {"missing", "not_found", "not_run", "missing_tooling", "artifact_unavailable"}:
        return "missing"
    if raw in {"blocked", "failed", "unhealthy", "error", "misconfigured_tooling", "red"}:
        return "blocked"
    if raw in {"warning", "warnings", "needs_review", "stale", "degraded"}:
        return "needs_review"
    return "needs_review"


def license_status(feeds_data):
    if not isinstance(feeds_data, dict):
        return "missing", "No feeds.json summary was provided."
    license_obj = feeds_data.get("license") if isinstance(feeds_data.get("license"), dict) else {}
    name = sanitize(license_obj.get("name") or feeds_data.get("license_name"), "")
    url = sanitize(license_obj.get("url") or feeds_data.get("license_url"), "")
    if name and url:
        return "present", f"license name and URL are present: {name}"
    if name or url:
        return "partial", "license metadata is incomplete"
    for feed in feeds_data.get("feeds", []) if isinstance(feeds_data.get("feeds"), list) else []:
        if isinstance(feed, dict) and sanitize(feed.get("license_name"), "") and sanitize(feed.get("license_url"), ""):
            return "partial", "per-feed license metadata is present but root license metadata needs review"
    return "missing", "open license metadata is missing"


def contact_status(feeds_data):
    if not isinstance(feeds_data, dict):
        return "missing", "No feeds.json summary was provided."
    contact = sanitize(feeds_data.get("technical_contact_email"), "")
    if "@" in contact:
        return "present", "technical contact email is present"
    for feed in feeds_data.get("feeds", []) if isinstance(feeds_data.get("feeds"), list) else []:
        if isinstance(feed, dict) and "@" in sanitize(feed.get("contact_email"), ""):
            return "partial", "per-feed contact metadata is present but root technical contact needs review"
    return "missing", "technical contact metadata is missing"


def fetch_url(url):
    parsed = urllib.parse.urlparse(url)
    if parsed.scheme not in {"http", "https"} or not parsed.netloc:
        return "needs_review", "unsafe or invalid URL"
    for method in ("HEAD", "GET"):
        req = urllib.request.Request(url, method=method, headers={"User-Agent": "open-transit-rt-caltrans-readiness-check/1"})
        try:
            with urllib.request.urlopen(req, timeout=FETCH_TIMEOUT) as response:
                code = getattr(response, "status", 0)
                if 200 <= code < 400:
                    return "present", f"{method} returned HTTP {code}"
                return "needs_review", f"{method} returned HTTP {code}"
        except urllib.error.HTTPError as exc:
            if method == "HEAD" and exc.code in {405, 501}:
                continue
            return "needs_review", f"{method} returned HTTP {exc.code}"
        except Exception as exc:  # noqa: BLE001 - shell helper records bounded diagnostics only.
            if method == "HEAD":
                continue
            return "needs_review", sanitize(type(exc).__name__)
    return "needs_review", "fetch failed"


def public_fetch_row(urls):
    if DRY_RUN or not RUN_FETCH:
        return row(
            "public_fetchability",
            "not_checked",
            "public fetches are disabled by default",
            "RUN_PUBLIC_FETCH=false",
            "Set RUN_PUBLIC_FETCH=true only when the target root is approved for bounded GET/HEAD checks.",
            "Fetchability checks are point-in-time signals only, not compliance or consumer acceptance proof.",
        )
    if not urls:
        return row("public_fetchability", "missing", "no URLs available to fetch", "configured feed URLs", "Provide stable feed URLs before checking fetchability.", "Fetchability cannot prove compliance or consumer acceptance.")
    results = []
    statuses = []
    for key in ("feeds_json", "schedule", "vehicle_positions", "trip_updates", "alerts"):
        url = urls.get(key)
        if not url:
            statuses.append("missing")
            results.append(f"{key}: missing URL")
            continue
        status, detail = fetch_url(url)
        statuses.append(status)
        results.append(f"{key}: {status} ({detail})")
    if all(status == "present" for status in statuses):
        status = "present"
    elif all(status in {"missing", "needs_review"} for status in statuses):
        status = "needs_review"
    else:
        status = "partial"
    return row("public_fetchability", status, "; ".join(results), "bounded opt-in fetch", "Review non-present fetch results; do not convert them to success.", "Public fetch results are supporting signals only.")


def https_row(urls):
    if not urls:
        return row("https", "missing", "no URLs available for HTTPS review", "configured feed URLs", "Provide stable feed URLs.", "HTTPS review is not final-root proof.")
    schemes = [urllib.parse.urlparse(url).scheme for url in urls.values()]
    if all(scheme == "https" for scheme in schemes):
        return row("https", "present", "all configured URLs use HTTPS", "configured feed URLs", "Verify final-root ownership separately.", "HTTPS URLs are required signals but not compliance proof.")
    if any(scheme == "https" for scheme in schemes):
        return row("https", "partial", "some configured URLs use HTTPS", "configured feed URLs", "Move all public feed URLs to HTTPS before external review.", "Mixed URL schemes do not prove readiness.")
    return row("https", "needs_review", "configured URLs are not HTTPS", "configured feed URLs", "Use deployment-owned HTTPS before any external packet.", "HTTP or loopback URLs are not public compliance proof.")


def consumer_row(consumer_data, source_info):
    if not isinstance(consumer_data, dict):
        return row("consumer_packet_preparedness", "missing", source_info["detail"], source_info["source"], "Consumer tracker could not be read.", "Consumer packet state is separate from submission, acceptance, ingestion, listing, or display.")
    records = consumer_data.get("targets", [])
    seen = {record.get("target"): record.get("status") for record in records if isinstance(record, dict)}
    ordered = list(seen)
    if ordered == EXPECTED_CONSUMERS and all(seen.get(name) == "prepared" for name in EXPECTED_CONSUMERS):
        return row("consumer_packet_preparedness", "present", "seven tracked consumer and aggregator targets remain prepared only", source_info["source"], "No target-originated status exists in this gap report.", "Prepared packets do not prove submission, review, acceptance, ingestion, listing, or display.")
    return row("consumer_packet_preparedness", "needs_review", f"consumer tracker differs from prepared-only expectation: {seen}", source_info["source"], "Restore target-originated evidence discipline before any status change.", "Consumer status changes require target-originated retained evidence.")


def summary_row(kind, data, source_info, missing_status, boundary):
    if not isinstance(data, dict):
        return row(kind, missing_status, source_info["detail"], source_info["source"], "Provide a safe .cache summary to review this signal.", boundary)
    status = map_summary_status(data.get("overall_status") or data.get("status") or data.get("trip_id_consistency_status"))
    return row(kind, status, data.get("summary") or data.get("overall_status") or data.get("status") or "safe summary loaded", source_info["source"], "Review source detail; missing or stale inputs stay visible.", boundary)


feeds_data, feeds_source = read_json_source("feeds_json", feeds_json_arg)
validator_data, validator_source = read_json_source("validator_health", validator_arg)
reliability_data, reliability_source = read_json_source("operations_reliability", reliability_arg)
trip_data, trip_source = read_json_source("trip_id_consistency", trip_consistency_arg)
consumer_data, consumer_source = read_json_source("consumer_status", consumer_status_arg, required=True, allow_consumer_status=True)

urls, urls_source = feed_urls_from_inputs(feeds_data)
rows = [
    row(
        "stable_urls",
        url_status(urls, ("feeds_json", "schedule", "vehicle_positions", "trip_updates", "alerts")),
        "configured URLs: " + (", ".join(f"{key}={urls[key]}" for key in sorted(urls)) if urls else "none"),
        urls_source,
        "Stable public feed URLs must be supplied and owned or approved outside this cache-only report.",
        "Stable URLs are readiness inputs only; this report does not prove final-root ownership or compliance.",
    ),
    feed_row("static_gtfs", "schedule", "Static GTFS", urls, urls_source),
    feed_row("vehicle_positions", "vehicle_positions", "Vehicle Positions", urls, urls_source),
    feed_row("trip_updates", "trip_updates", "Trip Updates", urls, urls_source),
    feed_row("alerts", "alerts", "Alerts", urls, urls_source),
    public_fetch_row(urls),
    https_row(urls),
]
license_state, license_signal = license_status(feeds_data)
contact_state, contact_signal = contact_status(feeds_data)
rows.append(row("open_license", license_state, license_signal, feeds_source["source"], "Provide operator-reviewed open license metadata before external review.", "Open license metadata is not agency approval or compliance proof."))
rows.append(row("contact", contact_state, contact_signal, feeds_source["source"], "Provide monitored technical/feed contact metadata.", "Contact metadata is a readiness signal only."))
rows.append(summary_row("validators", validator_data, validator_source, "missing", "Validator output is a supporting signal only, not compliance, consumer acceptance, or correctness proof."))
rows.append(summary_row("freshness", reliability_data, reliability_source, "missing", "Freshness and reliability summaries are private diagnostics only, not production readiness, SLA, or uptime proof."))
rows.append(summary_row("trip_id_consistency", trip_data, trip_source, "not_checked", "Trip ID consistency signals do not prove production-grade ETA quality or consumer acceptance."))
rows.append(consumer_row(consumer_data, consumer_source))
rows.append(row("unsupported_claim_boundaries", "present", "all claim flags are false and Phase 54 context was not refreshed", "script policy", "No intake, retained evidence, or target-originated status is present.", "This gap report is not evidence and makes no unsupported claim."))

input_inventory = [feeds_source, validator_source, reliability_source, trip_source, consumer_source]
missing_inputs = [item for item in input_inventory if item["status"] in {"missing", "not_checked", "blocked"}]
counts = {status: 0 for status in sorted(ALLOWED_STATUSES)}
for item in rows:
    counts[item["status"]] += 1

rank = {"present": 0, "not_checked": 1, "partial": 2, "missing": 3, "needs_review": 4, "blocked": 5}
overall = "present"
for item in rows:
    if rank[item["status"]] > rank[overall]:
        overall = item["status"]

summary = {
    "schema_version": "open-transit-rt.caltrans_readiness_gap.v1",
    "generated_at": timestamp,
    "overall_status": overall,
    "dry_run": DRY_RUN,
    "official_source_context": OFFICIAL_CONTEXT,
    "inputs": input_inventory,
    "missing_or_not_checked_inputs": missing_inputs,
    "configured_feed_urls": urls,
    "rows": rows,
    "counts": counts,
    "consumer_tracker": {
        "source": consumer_source["source"],
        "expected_targets": EXPECTED_CONSUMERS,
        "prepared_only": isinstance(consumer_data, dict) and consumer_row(consumer_data, consumer_source)["status"] == "present",
    },
    "unsupported_claim_boundaries": [
        "No CAL-ITP/Caltrans compliance claim.",
        "No consumer submission, review, acceptance, ingestion, listing, or display claim.",
        "No agency adoption, endorsement, approval, final-root proof, or public launch claim.",
        "No hosted SaaS, paid support, SLA, uptime, production readiness, vendor compatibility, hardware certification, production AVL reliability, or production-grade ETA claim.",
        "Validator and fetch results are supporting signals only.",
    ],
    "claim_flags": CLAIM_FLAGS,
    "boundary": "Cache-only readiness gap report. Not retained evidence, not official-source refresh, not compliance proof, not consumer acceptance, not final-root proof, and not production readiness.",
}

manifest = {
    "schema_version": "open-transit-rt.caltrans_readiness_gap_manifest.v1",
    "generated_at": timestamp,
    "output_dir": rel_to_root(resolve_output_dir()),
    "output_files": list(OUTPUT_FILES),
    "default_output_root": ".cache/caltrans-readiness-check",
    "official_sources_refreshed": False,
    "retained_evidence_created": False,
    "docs_evidence_written": False,
    "consumer_statuses_changed": False,
    "external_parties_contacted": False,
    "public_fetch_opt_in": RUN_FETCH and not DRY_RUN,
    "claim_flags": CLAIM_FLAGS,
}

out = pathlib.Path(manifest["output_dir"])
out = ROOT / out

(out / "summary.json").write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8")
(out / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")

summary_md = [
    "# Caltrans Readiness Gap Summary",
    "",
    f"- Generated at: `{timestamp}`",
    f"- Overall status: `{overall}`",
    "- Phase 54 official-source context: `2026-05-09`, not refreshed by this script",
    "- Scope: private `.cache` gap diagnostics only",
    "- Claim flags: all `false`",
    "",
    "This is a readiness gap report only. It is not retained evidence, not an official requirements refresh, not CAL-ITP/Caltrans compliance proof, not consumer acceptance, not final-root proof, and not production readiness.",
    "",
    "## Rows",
]
for item in rows:
    summary_md.append(f"- `{item['status']}` `{item['area']}`: {item['signal']}")
(out / "summary.md").write_text("\n".join(summary_md) + "\n", encoding="utf-8")

manifest_md = [
    "# Caltrans Readiness Gap Manifest",
    "",
    f"- Generated at: `{timestamp}`",
    "- Output files: `summary.json`, `summary.md`, `manifest.json`, `manifest.md`, `gap-review.txt`",
    "- Output root: `.cache/caltrans-readiness-check` by default",
    "- Official sources refreshed: `false`",
    "- Retained evidence created: `false`",
    "- Consumer statuses changed: `false`",
    "- External parties contacted: `false`",
]
(out / "manifest.md").write_text("\n".join(manifest_md) + "\n", encoding="utf-8")

review = [
    "Caltrans-style readiness gap review",
    "",
    "This is a cache-only gap report, not evidence and not compliance proof.",
    "Phase 54 official-source context is used by date only: May 9, 2026.",
    "This script did not refresh official requirements.",
    f"Overall status: {overall}",
    "",
    "Rows:",
]
for item in rows:
    review.append(f"- {item['area']}: {item['status']} - {item['missing_or_review_needed']}")
review.extend([
    "",
    "Unsupported-claim boundary:",
    "No intake, no retained evidence, no target-originated status, and no official-source refresh means no compliance, consumer acceptance, final-root, agency approval, production readiness, SLA, hosted SaaS, vendor compatibility, production AVL reliability, or production-grade ETA claim.",
])
(out / "gap-review.txt").write_text("\n".join(review) + "\n", encoding="utf-8")

actual = sorted(p.name for p in out.iterdir() if p.is_file())
if actual != sorted(OUTPUT_FILES):
    fail(f"unexpected output files: {actual}")
for name in OUTPUT_FILES:
    text = (out / name).read_text(encoding="utf-8")
    if len(text.encode("utf-8")) > 131072:
        fail(f"{name} exceeds bounded output size")
for item in summary["rows"]:
    if item["status"] in BAD_STATUS_WORDS:
        fail(f"claim-upgrading status emitted: {item}")
if any(CLAIM_FLAGS.values()):
    fail("claim flags must remain false")

print(f"caltrans readiness gap diagnostics written to {rel_to_root(out)}")
PY
