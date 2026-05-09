#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
COMPLIANCE_PACKET_DEPLOYMENT_NAME="${COMPLIANCE_PACKET_DEPLOYMENT_NAME:-}"
COMPLIANCE_PACKET_ROOT_URL="${COMPLIANCE_PACKET_ROOT_URL:-}"
COMPLIANCE_PACKET_OUTPUT_DIR="${COMPLIANCE_PACKET_OUTPUT_DIR:-}"
COMPLIANCE_PACKET_FORCE="${COMPLIANCE_PACKET_FORCE:-false}"
COMPLIANCE_PACKET_MAX_SOURCE_BYTES="${COMPLIANCE_PACKET_MAX_SOURCE_BYTES:-262144}"
COMPLIANCE_PACKET_HUMAN_REVIEW="${COMPLIANCE_PACKET_HUMAN_REVIEW:-false}"
COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT="${COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT:-}"
COMPLIANCE_PACKET_FINAL_ROOT_EVIDENCE="${COMPLIANCE_PACKET_FINAL_ROOT_EVIDENCE:-}"
COMPLIANCE_PACKET_VALIDATION_EVIDENCE="${COMPLIANCE_PACKET_VALIDATION_EVIDENCE:-}"
COMPLIANCE_PACKET_OPERATIONS_EVIDENCE="${COMPLIANCE_PACKET_OPERATIONS_EVIDENCE:-}"
COMPLIANCE_PACKET_CONSUMER_EVIDENCE="${COMPLIANCE_PACKET_CONSUMER_EVIDENCE:-docs/evidence/consumer-submissions/status.json}"
COMPLIANCE_PACKET_OCI_PILOT_EVIDENCE="${COMPLIANCE_PACKET_OCI_PILOT_EVIDENCE:-docs/evidence/captured/oci-pilot/2026-04-24}"

usage() {
  cat <<'EOF'
Usage:
  scripts/generate-compliance-evidence-packet.sh [--help]

Environment:
  COMPLIANCE_PACKET_DEPLOYMENT_NAME   Deployment/operator label for this local review packet.
  COMPLIANCE_PACKET_ROOT_URL          Deployment public root URL. HTTPS required except loopback HTTP test roots.
  COMPLIANCE_PACKET_OUTPUT_DIR        Defaults to .cache/compliance-evidence-packet/<UTC timestamp>.
  COMPLIANCE_PACKET_FORCE             true|false; allow reuse of a non-empty output directory.
  COMPLIANCE_PACKET_MAX_SOURCE_BYTES  Positive integer read cap for optional source summaries, default 262144.
  COMPLIANCE_PACKET_HUMAN_REVIEW      true|false; labels operator-supplied review text only.
  COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT Optional already-redacted human review note.
  COMPLIANCE_PACKET_FINAL_ROOT_EVIDENCE      Optional retained final-root evidence path summary.
  COMPLIANCE_PACKET_VALIDATION_EVIDENCE      Optional retained validation evidence path summary.
  COMPLIANCE_PACKET_OPERATIONS_EVIDENCE      Optional retained operations evidence path summary.
  COMPLIANCE_PACKET_CONSUMER_EVIDENCE        Optional tracker path, default docs/evidence/consumer-submissions/status.json.
  COMPLIANCE_PACKET_OCI_PILOT_EVIDENCE       Optional OCI pilot evidence path summary.

Safety:
  This generator is a local summarizer. It writes ignored .cache output only,
  creates no retained evidence, contacts no consumers, fetches no live feeds,
  changes no consumer tracker state, and makes no compliance, acceptance,
  ingestion, adoption, hosted SaaS, production readiness, SLA, vendor
  compatibility, marketplace approval, or production-grade ETA claim.
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

python3 - "$ROOT_DIR" "$TIMESTAMP" \
  "$COMPLIANCE_PACKET_DEPLOYMENT_NAME" "$COMPLIANCE_PACKET_ROOT_URL" \
  "$COMPLIANCE_PACKET_OUTPUT_DIR" "$COMPLIANCE_PACKET_FORCE" \
  "$COMPLIANCE_PACKET_MAX_SOURCE_BYTES" "$COMPLIANCE_PACKET_HUMAN_REVIEW" \
  "$COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT" \
  "$COMPLIANCE_PACKET_FINAL_ROOT_EVIDENCE" "$COMPLIANCE_PACKET_VALIDATION_EVIDENCE" \
  "$COMPLIANCE_PACKET_OPERATIONS_EVIDENCE" "$COMPLIANCE_PACKET_CONSUMER_EVIDENCE" \
  "$COMPLIANCE_PACKET_OCI_PILOT_EVIDENCE" <<'PY'
import datetime as dt
import hashlib
import json
import os
import pathlib
import re
import shutil
import sys
import urllib.parse

(
    root_arg,
    timestamp,
    deployment_name_arg,
    root_url_arg,
    output_arg,
    force_arg,
    max_source_bytes_arg,
    human_review_arg,
    human_review_text_arg,
    final_root_path_arg,
    validation_path_arg,
    operations_path_arg,
    consumer_path_arg,
    oci_pilot_path_arg,
) = sys.argv[1:15]

ROOT = pathlib.Path(root_arg).resolve()
FORCE = force_arg == "true"
HUMAN_REVIEW = human_review_arg == "true"
DEPLOYMENT_NAME = deployment_name_arg.strip()
ROOT_URL = root_url_arg.strip().rstrip("/")
MAX_SOURCE_BYTES = int(max_source_bytes_arg)

BLOCKER_FILES = ("blocker.json", "blocker.md", "manifest.json", "manifest.md")
PACKET_FILES = (
    "summary.json",
    "summary.md",
    "readiness-packet.md",
    "evidence-map.json",
    "evidence-map.md",
    "missing-evidence.md",
    "human-review.md",
    "manifest.json",
    "manifest.md",
)
EXPECTED_CONSUMERS = [
    "Google Maps",
    "Apple Maps",
    "Transit App",
    "Bing Maps",
    "Moovit",
    "Mobility Database",
    "transit.land",
]
CLAIM_FLAGS = {
    "compliance_claimed": False,
    "consumer_acceptance_claimed": False,
    "consumer_ingestion_claimed": False,
    "agency_adoption_claimed": False,
    "hosted_saas_claimed": False,
    "production_readiness_claimed": False,
    "sla_or_uptime_claimed": False,
    "vendor_compatibility_claimed": False,
    "production_grade_eta_claimed": False,
    "marketplace_approval_claimed": False,
}
OFFICIAL_REQUIREMENTS_REVIEW_DATE = "2026-05-09"
ALLOWED_STATUSES = {"present", "partial", "missing", "blocked", "pilot_only", "needs_review"}

UNSAFE_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"authorization\s*:",
        r"cookie\s*:",
        r"\bbearer\s+[A-Za-z0-9._~+/=-]{8,}",
        r"postgres(?:ql)?://",
        r"database_url\s*[:=]",
        r"begin [a-z ]*private key",
        r"acme[_ -]?account",
        r"private[_ -]?key",
        r"secret\s*[:=]",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|diagnostic|telemetry|correspondence)",
        r"consumer[_ -]?portal",
        r"db[_ -]?url\s*[:=]",
        r"token\s*[:=]\s*[A-Za-z0-9._~+/=-]{8,}",
    )
]


def fail(message):
    raise SystemExit(f"ERROR: {message}")


def bool_value(name, value):
    if value not in {"true", "false"}:
        fail(f"{name} must be true or false")


def positive_int(name, value):
    if not re.fullmatch(r"[1-9][0-9]*", value):
        fail(f"{name} must be a positive integer")


def is_relative_to(path, base):
    try:
        pathlib.Path(path).resolve(strict=False).relative_to(pathlib.Path(base).resolve(strict=False))
        return True
    except ValueError:
        return False


def has_symlink(path):
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


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def write_json(path, data):
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def unsafe_hits(text):
    return [pattern.pattern for pattern in UNSAFE_PATTERNS if pattern.search(text)]


def validate_url(raw):
    parsed = urllib.parse.urlparse(raw)
    if parsed.username or parsed.password:
        fail("COMPLIANCE_PACKET_ROOT_URL must not contain credentials")
    if parsed.query or parsed.fragment:
        fail("COMPLIANCE_PACKET_ROOT_URL must not contain query strings or fragments")
    if not parsed.scheme or not parsed.netloc:
        fail("COMPLIANCE_PACKET_ROOT_URL must be an absolute URL")
    host = parsed.hostname or ""
    loopback = host in {"localhost", "127.0.0.1", "::1"} or host.startswith("127.")
    if parsed.scheme == "https" or (parsed.scheme == "http" and loopback):
        return raw.rstrip("/")
    fail("COMPLIANCE_PACKET_ROOT_URL must use HTTPS except loopback HTTP test roots")


def validate_output_dir():
    default = ROOT / ".cache" / "compliance-evidence-packet" / timestamp
    raw = pathlib.Path(output_arg) if output_arg else default
    out = raw if raw.is_absolute() else ROOT / raw
    resolved = out.resolve(strict=False)
    if any(part == ".." for part in raw.parts):
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must not contain traversal components")
    if has_symlink(out):
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must not contain symlink directories")
    if not is_relative_to(resolved, ROOT):
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must stay inside the repository")
    if is_relative_to(resolved, ROOT / "docs" / "evidence"):
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must not be under docs/evidence")
    if not is_relative_to(resolved, ROOT / ".cache"):
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must resolve under ignored .cache")
    parts = set(resolved.relative_to(ROOT / ".cache").parts)
    if {"docs", "evidence"}.issubset(parts) or "captured" in parts or "consumer-submissions" in parts:
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must not use evidence-like path components")
    if resolved.exists() and not resolved.is_dir():
        fail("COMPLIANCE_PACKET_OUTPUT_DIR must be a directory")
    if resolved.exists() and any(resolved.iterdir()):
        if not FORCE:
            fail("COMPLIANCE_PACKET_OUTPUT_DIR exists and is non-empty; use COMPLIANCE_PACKET_FORCE=true")
        for child in resolved.iterdir():
            if child.is_symlink() or child.is_file():
                child.unlink()
            else:
                shutil.rmtree(child)
    resolved.mkdir(parents=True, exist_ok=True)
    os.chmod(resolved, 0o700)
    return resolved


def load_consumer_tracker(path_text):
    path = pathlib.Path(path_text) if path_text else ROOT / "docs/evidence/consumer-submissions/status.json"
    full = path if path.is_absolute() else ROOT / path
    if has_symlink(full):
        fail("consumer tracker path must not contain symlink directories")
    data = json.loads(full.read_text(encoding="utf-8"))
    records = data.get("targets", [])
    targets = [{"target": row.get("target"), "status": row.get("status"), "packet_path": row.get("packet_path")} for row in records]
    return full, data, targets


def summarize_path(label, path_text):
    if not path_text:
        return {
            "label": label,
            "path": "",
            "status": "missing",
            "exists": False,
            "classification": "missing",
            "notes": "No path configured for this packet run.",
        }
    path = pathlib.Path(path_text)
    full = path if path.is_absolute() else ROOT / path
    if has_symlink(full):
        fail(f"{label} path must not contain symlink directories")
    summary = {
        "label": label,
        "path": rel(full),
        "exists": full.exists(),
        "classification": "needs_review",
        "status": "present" if full.exists() else "missing",
    }
    if not full.exists():
        summary["notes"] = "Configured path does not exist in this checkout."
        return summary
    if full.is_file():
        with full.open("rb") as f:
            data = f.read(MAX_SOURCE_BYTES + 1)
        if len(data) > MAX_SOURCE_BYTES:
            summary["truncated"] = True
            data = data[:MAX_SOURCE_BYTES]
        text = data.decode("utf-8", errors="ignore")
        hits = unsafe_hits(text)
        if hits:
            fail(f"{label} source contains unsafe private strings")
        h = hashlib.sha256()
        with full.open("rb") as f:
            for chunk in iter(lambda: f.read(1024 * 1024), b""):
                h.update(chunk)
        summary.update({
            "kind": "file",
            "bytes": full.stat().st_size,
            "sha256": h.hexdigest(),
        })
        return summary
    if full.is_dir():
        files = []
        total = 0
        for child in sorted(full.rglob("*")):
            if child.is_symlink():
                fail(f"{label} source contains symlink: {rel(child)}")
            if child.is_file():
                total += 1
                if len(files) < 25:
                    files.append(rel(child))
        summary.update({
            "kind": "directory",
            "file_count": total,
            "sample_files": files,
            "notes": "Directory contents were not copied or read into the packet; only bounded inventory metadata is summarized.",
        })
        return summary
    summary["notes"] = "Path exists but is not a regular file or directory."
    return summary


def consumer_summary(targets):
    statuses = {}
    for row in targets:
        statuses[row["status"]] = statuses.get(row["status"], 0) + 1
    exact_prepared = [row["target"] for row in targets] == EXPECTED_CONSUMERS and all(row["status"] == "prepared" for row in targets)
    return {
        "targets": targets,
        "status_counts": statuses,
        "all_expected_targets_prepared": exact_prepared,
        "claim_boundary": "Prepared packets are not submissions, review, acceptance, ingestion, listing, or display.",
    }


def readiness_rows(final_root_summary, validation_summary, operations_summary, oci_summary, consumers):
    final_root_present = final_root_summary["exists"]
    validation_present = validation_summary["exists"]
    operations_present = operations_summary["exists"]
    oci_present = oci_summary["exists"]
    prepared_only = consumers["all_expected_targets_prepared"]
    rows = [
        {
            "requirement": "RQ-4A",
            "title": "Complete realtime feed set",
            "status": "pilot_only" if oci_present else "partial",
            "evidence_signal": "Repository exposes Schedule, Vehicle Positions, Trip Updates, and Alerts paths; OCI pilot is hosted/operator pilot evidence only.",
            "missing_or_review_needed": "Final deployment must prove all three realtime feeds are live, fresh, synchronized, and validator-clean.",
        },
        {
            "requirement": "RQ-4B",
            "title": "Stable public production URLs",
            "status": "present" if final_root_present else "missing",
            "evidence_signal": final_root_summary["path"] if final_root_present else "No retained final-root evidence path configured.",
            "missing_or_review_needed": "Agency-owned or agency-approved final-root DNS/TLS/public-fetch proof is required.",
        },
        {
            "requirement": "RQ-4C",
            "title": "Validator-clean feeds",
            "status": "present" if validation_present else ("pilot_only" if oci_present else "missing"),
            "evidence_signal": validation_summary["path"] if validation_present else "No final deployment validation evidence path configured.",
            "missing_or_review_needed": "Current no-error validation records are required for the exact deployment root and feed scope.",
        },
        {
            "requirement": "RQ-4D",
            "title": "Open license and discoverability",
            "status": "needs_review",
            "evidence_signal": "Repo supports metadata and prepared packet fields; operator must review deployment-specific license/contact/source-of-truth pages.",
            "missing_or_review_needed": "Agency-approved license/contact and provider or agreed regional source-of-truth page proof remain required.",
        },
        {
            "requirement": "RQ-4E",
            "title": "Consumer ingestion workflow",
            "status": "needs_review" if prepared_only else "blocked",
            "evidence_signal": "All seven targets remain prepared only." if prepared_only else "Consumer tracker is not in the exact prepared-only state.",
            "missing_or_review_needed": "No submission, review, acceptance, ingestion, listing, or display evidence is present.",
        },
        {
            "requirement": "RQ-4F",
            "title": "Marketplace/vendor equivalence separation",
            "status": "missing",
            "evidence_signal": "No marketplace/vendor equivalence evidence is represented in this packet.",
            "missing_or_review_needed": "Support packaging, SLA/KPI commitments, procurement artifacts, and external approvals are outside this packet.",
        },
        {
            "requirement": "RQ-4G",
            "title": "Compliance dashboard and scorecard",
            "status": "partial" if operations_summary["exists"] else "needs_review",
            "evidence_signal": operations_summary["path"] if operations_summary["exists"] else "Operations/readiness tooling exists, but no deployment operations packet path was configured.",
            "missing_or_review_needed": "Deployment-specific export/review evidence is required before stronger readiness language.",
        },
    ]
    for row in rows:
        if row["status"] not in ALLOWED_STATUSES:
            fail(f"unsupported readiness status generated: {row['status']}")
    return rows


def missing_evidence(final_root_summary, validation_summary, operations_summary, consumers):
    items = []
    if not final_root_summary["exists"]:
        items.append({"area": "final_root", "status": "missing", "detail": "No retained agency-owned or agency-approved final public root evidence path was configured."})
    if not validation_summary["exists"]:
        items.append({"area": "validation", "status": "missing", "detail": "No current final-deployment no-error validation evidence path was configured."})
    if not operations_summary["exists"]:
        items.append({"area": "operations", "status": "needs_review", "detail": "No final deployment operations evidence path was configured."})
    items.extend([
        {"area": "consumer_submissions", "status": "missing", "detail": "Consumer tracker remains prepared-only; no submission, review, acceptance, ingestion, listing, or display evidence is present."},
        {"area": "agency_adoption", "status": "missing", "detail": "No agency adoption, endorsement, or approval evidence is represented."},
        {"area": "production_eta_quality", "status": "missing", "detail": "No real-world production-grade ETA quality evidence is represented."},
        {"area": "sla_uptime", "status": "missing", "detail": "No SLA or uptime commitment/evidence is represented."},
    ])
    if not consumers["all_expected_targets_prepared"]:
        items.append({"area": "consumer_tracker", "status": "blocked", "detail": "Consumer tracker is not in the exact seven-target prepared-only state."})
    return items


def scan_output(out):
    for path in out.rglob("*"):
        if not path.is_file():
            continue
        data = path.read_bytes()
        if b"\x00" in data:
            continue
        text = data.decode("utf-8", errors="ignore")
        hits = unsafe_hits(text)
        if hits:
            fail(f"unsafe private string detected in generated packet: {rel(path)}")
        if re.search(r'"status"\s*:\s*"compliant"', text, re.I):
            fail(f"disallowed compliant status generated: {rel(path)}")
        if re.search(r'\|\s*`?compliant`?\s*\|', text, re.I):
            fail(f"disallowed compliant status generated: {rel(path)}")


def write_blocker(out, reason):
    now = dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    data = {
        "schema": "compliance-evidence-blocker.v1",
        "packet_version": "phase-55.v1",
        "generated_at": now,
        "mode": "blocker",
        "reason": reason,
        "deployment_identity_present": bool(DEPLOYMENT_NAME),
        "root_url_present": bool(ROOT_URL),
        "retained_evidence_created": False,
        "consumer_tracker_changed": False,
        "claim_flags": CLAIM_FLAGS.copy(),
    }
    manifest = {
        "schema": "compliance-evidence-manifest.v1",
        "packet_type": "blocker",
        "generated_at": now,
        "packet_dir": rel(out),
        "expected_files": list(BLOCKER_FILES),
        "retained_evidence_created": False,
        "claim_flags": CLAIM_FLAGS.copy(),
    }
    write_json(out / "blocker.json", data)
    (out / "blocker.md").write_text(
        "# Compliance Evidence Packet Blocker\n\n"
        f"- Generated UTC: {now}\n"
        f"- Reason: {reason}\n"
        "- Retained evidence created: false\n"
        "- Consumer tracker changed: false\n"
        "- Claim flags: all false\n\n"
        "This blocker packet exists because the deployment identity or public root URL was missing. "
        "It is not compliance evidence and does not prove consumer acceptance, consumer ingestion, "
        "agency adoption, hosted SaaS availability, production readiness, SLA/uptime, vendor compatibility, "
        "marketplace approval, or production-grade ETA quality.\n",
        encoding="utf-8",
    )
    write_json(out / "manifest.json", manifest)
    (out / "manifest.md").write_text(
        "# Compliance Evidence Blocker Manifest\n\n"
        f"- Packet type: blocker\n- Packet directory: `{rel(out)}`\n"
        f"- Expected files: `{', '.join(BLOCKER_FILES)}`\n"
        "- Retained evidence created: false\n- Claim flags: all false\n",
        encoding="utf-8",
    )


def write_packet(out):
    now = dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    root_url = validate_url(ROOT_URL)
    tracker_path, tracker_data, targets = load_consumer_tracker(consumer_path_arg)
    consumers = consumer_summary(targets)
    final_root_summary = summarize_path("final_root", final_root_path_arg)
    validation_summary = summarize_path("validation", validation_path_arg)
    operations_summary = summarize_path("operations", operations_path_arg)
    consumer_source_summary = summarize_path("consumer_tracker", str(tracker_path))
    oci_summary = summarize_path("oci_pilot", oci_pilot_path_arg)
    rows = readiness_rows(final_root_summary, validation_summary, operations_summary, oci_summary, consumers)
    missing = missing_evidence(final_root_summary, validation_summary, operations_summary, consumers)
    human_review_state = {
        "reviewed": HUMAN_REVIEW,
        "text_present": bool(human_review_text_arg.strip()),
        "text": human_review_text_arg.strip() if HUMAN_REVIEW else "",
        "claim_boundary": "Human review text is operator-supplied, must already be redacted, and does not create a compliance claim.",
    }
    if human_review_state["text"] and unsafe_hits(human_review_state["text"]):
        fail("COMPLIANCE_PACKET_HUMAN_REVIEW_TEXT contains unsafe private strings")
    evidence_sources = [
        final_root_summary,
        validation_summary,
        operations_summary,
        consumer_source_summary,
        oci_summary,
    ]
    summary = {
        "schema": "compliance-evidence-summary.v1",
        "packet_version": "phase-55.v1",
        "generated_at": now,
        "deployment": {
            "name": DEPLOYMENT_NAME,
            "root_url": root_url,
        },
        "source_evidence_path_summaries": evidence_sources,
        "official_requirements_review_date": OFFICIAL_REQUIREMENTS_REVIEW_DATE,
        "readiness_rows": rows,
        "consumer_tracker_summary": consumers,
        "missing_evidence": missing,
        "human_review": human_review_state,
        "claim_flags": CLAIM_FLAGS.copy(),
    }
    manifest = {
        "schema": "compliance-evidence-manifest.v1",
        "packet_type": "deployment",
        "generated_at": now,
        "packet_dir": rel(out),
        "expected_files": list(PACKET_FILES),
        "retained_evidence_created": False,
        "consumer_tracker_changed": False,
        "source_paths": [source["path"] for source in evidence_sources if source.get("path")],
        "claim_flags": CLAIM_FLAGS.copy(),
    }
    evidence_map = {
        "schema": "compliance-evidence-map.v1",
        "official_requirements_review_date": OFFICIAL_REQUIREMENTS_REVIEW_DATE,
        "deployment": summary["deployment"],
        "rows": rows,
        "sources": evidence_sources,
    }
    write_json(out / "summary.json", summary)
    (out / "summary.md").write_text(
        "# Compliance Evidence Packet Summary\n\n"
        f"- Deployment: `{DEPLOYMENT_NAME}`\n"
        f"- Root URL: `{root_url}`\n"
        f"- Generated UTC: {now}\n"
        f"- Official requirements review date: {OFFICIAL_REQUIREMENTS_REVIEW_DATE}\n"
        "- Retained evidence created: false\n"
        "- Consumer tracker changed: false\n"
        "- Claim flags: all false\n\n"
        "This packet summarizes configured local evidence paths for human review. It does not create evidence, "
        "contact consumers, fetch live feeds, or claim compliance.\n",
        encoding="utf-8",
    )
    (out / "readiness-packet.md").write_text(
        "# Readiness Packet\n\n"
        "| Requirement | Status | Evidence signal | Missing or review needed |\n"
        "| --- | --- | --- | --- |\n"
        + "".join(
            f"| {row['requirement']} {row['title']} | `{row['status']}` | {row['evidence_signal']} | {row['missing_or_review_needed']} |\n"
            for row in rows
        )
        + "\nAllowed status values are present, partial, missing, blocked, pilot_only, and needs_review. "
        "No status in this packet is a compliance claim.\n",
        encoding="utf-8",
    )
    write_json(out / "evidence-map.json", evidence_map)
    (out / "evidence-map.md").write_text(
        "# Evidence Map\n\n"
        "| Source | Status | Path | Notes |\n"
        "| --- | --- | --- | --- |\n"
        + "".join(
            f"| `{source['label']}` | `{source['status']}` | `{source.get('path', '')}` | {source.get('notes', source.get('classification', 'needs_review'))} |\n"
            for source in evidence_sources
        )
        + "\nThe OCI pilot, when present, is classified only as hosted/operator pilot evidence and not agency-owned final-root proof.\n",
        encoding="utf-8",
    )
    (out / "missing-evidence.md").write_text(
        "# Missing Evidence\n\n"
        + "".join(f"- `{item['area']}` (`{item['status']}`): {item['detail']}\n" for item in missing)
        + "\nThese gaps must remain visible until retained, public-safe, claim-specific evidence exists.\n",
        encoding="utf-8",
    )
    (out / "human-review.md").write_text(
        "# Human Review\n\n"
        f"- Human review marked complete: {str(HUMAN_REVIEW).lower()}\n"
        f"- Human review text present: {str(bool(human_review_text_arg.strip())).lower()}\n\n"
        + (human_review_state["text"] + "\n\n" if human_review_state["text"] else "")
        + "Operator-supplied text must be redacted and separately reviewed. Human review does not turn this packet into compliance evidence.\n",
        encoding="utf-8",
    )
    write_json(out / "manifest.json", manifest)
    (out / "manifest.md").write_text(
        "# Compliance Evidence Packet Manifest\n\n"
        f"- Packet type: deployment\n- Packet directory: `{rel(out)}`\n"
        f"- Expected files: `{', '.join(PACKET_FILES)}`\n"
        "- Retained evidence created: false\n"
        "- Consumer tracker changed: false\n"
        "- Claim flags: all false\n",
        encoding="utf-8",
    )
    scan_output(out)


bool_value("COMPLIANCE_PACKET_FORCE", force_arg)
bool_value("COMPLIANCE_PACKET_HUMAN_REVIEW", human_review_arg)
positive_int("COMPLIANCE_PACKET_MAX_SOURCE_BYTES", max_source_bytes_arg)
out = validate_output_dir()

if ROOT_URL:
    validate_url(ROOT_URL)

if not DEPLOYMENT_NAME or not ROOT_URL:
    reasons = []
    if not DEPLOYMENT_NAME:
        reasons.append("missing deployment identity")
    if not ROOT_URL:
        reasons.append("missing root URL")
    write_blocker(out, "; ".join(reasons))
    scan_output(out)
    print(f"compliance evidence blocker packet: {rel(out)}")
    sys.exit(0)

write_packet(out)
print(f"compliance evidence packet: {rel(out)}")
PY
