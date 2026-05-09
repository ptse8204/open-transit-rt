#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
CAPTURE_DATE_UTC="${CAPTURE_DATE_UTC:-$(date -u '+%Y-%m-%d')}"
OUTPUT_DIR="${OUTPUT_DIR:-}"
FORCE="${FORCE:-false}"
ALLOW_CAPTURED_EVIDENCE_WRITE="${ALLOW_CAPTURED_EVIDENCE_WRITE:-false}"
FINAL_ROOT_BASE_URL="${FINAL_ROOT_BASE_URL:-}"
FINAL_ROOT_ENVIRONMENT_NAME="${FINAL_ROOT_ENVIRONMENT_NAME:-final-root}"
FINAL_ROOT_APPROVAL_ARTIFACT="${FINAL_ROOT_APPROVAL_ARTIFACT:-}"
FINAL_ROOT_APPROVAL_SUMMARY="${FINAL_ROOT_APPROVAL_SUMMARY:-}"
ADMIN_BASE_URL="${ADMIN_BASE_URL:-}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
CONNECT_TIMEOUT_SECONDS="${CONNECT_TIMEOUT_SECONDS:-5}"
REQUEST_TIMEOUT_SECONDS="${REQUEST_TIMEOUT_SECONDS:-20}"
MAX_FEED_BYTES="${MAX_FEED_BYTES:-20971520}"
BLOCKER_ONLY="false"
DRY_RUN="false"
RETAIN_CAPTURED="false"

usage() {
  cat <<'EOF'
Usage:
  scripts/collect-final-root-evidence.sh [--help] [--blocker-only] [--dry-run] [--retain-captured]

Environment:
  FINAL_ROOT_BASE_URL              Approved final public feed root. HTTPS required except loopback HTTP test roots.
  FINAL_ROOT_ENVIRONMENT_NAME      Safe environment folder label, default final-root.
  FINAL_ROOT_APPROVAL_ARTIFACT     Readable redacted approval artifact for real evidence collection.
  FINAL_ROOT_APPROVAL_SUMMARY      Optional short public-safe approval summary.
  CAPTURE_DATE_UTC                 YYYY-MM-DD, defaults to current UTC date.
  OUTPUT_DIR                       Defaults to .cache/final-root-evidence/<UTC timestamp>.
  FORCE                            true|false; allow reuse of a non-empty output directory.
  ALLOW_CAPTURED_EVIDENCE_WRITE    Must be true with --retain-captured before writing docs/evidence/captured.
  ADMIN_BASE_URL                   Optional admin origin for validator runs; must not contain credentials.
  ADMIN_TOKEN                      Optional bearer token for validator runs; token value is never written.
  CONNECT_TIMEOUT_SECONDS          Positive integer, default 5.
  REQUEST_TIMEOUT_SECONDS          Positive integer, default 20.
  MAX_FEED_BYTES                   Positive integer cap per fetched feed artifact, default 20971520.

Safety:
  Default and blocker-only output stays under ignored .cache storage and writes
  exactly blocker.json, blocker.md, manifest.json, and manifest.md. Retained
  docs/evidence/captured writes require --retain-captured,
  ALLOW_CAPTURED_EVIDENCE_WRITE=true, a valid final root, and a readable
  redacted approval artifact. This helper does not contact consumers, refresh
  consumer packets, change consumer statuses, or create compliance, agency
  adoption, hosted SaaS, production-readiness, SLA/uptime, vendor-compatibility,
  consumer-acceptance, or production-grade ETA claims.
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
    --blocker-only)
      BLOCKER_ONLY="true"
      ;;
    --dry-run)
      DRY_RUN="true"
      ;;
    --retain-captured)
      RETAIN_CAPTURED="true"
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
  shift
done

python3 - "$ROOT_DIR" "$TIMESTAMP" "$CAPTURE_DATE_UTC" "$OUTPUT_DIR" "$FORCE" \
  "$ALLOW_CAPTURED_EVIDENCE_WRITE" "$FINAL_ROOT_BASE_URL" \
  "$FINAL_ROOT_ENVIRONMENT_NAME" "$FINAL_ROOT_APPROVAL_ARTIFACT" \
  "$FINAL_ROOT_APPROVAL_SUMMARY" "$ADMIN_BASE_URL" "$ADMIN_TOKEN" \
  "$CONNECT_TIMEOUT_SECONDS" "$REQUEST_TIMEOUT_SECONDS" "$MAX_FEED_BYTES" \
  "$BLOCKER_ONLY" "$DRY_RUN" "$RETAIN_CAPTURED" <<'PY'
import base64
import datetime as dt
import hashlib
import json
import os
import pathlib
import re
import shutil
import socket
import ssl
import subprocess
import sys
import tempfile
import urllib.error
import urllib.parse
import urllib.request

(
    root_arg,
    timestamp,
    capture_date,
    output_arg,
    force_arg,
    allow_captured_arg,
    final_root_arg,
    environment_arg,
    approval_arg,
    approval_summary_arg,
    admin_base_arg,
    admin_token_arg,
    connect_timeout_arg,
    request_timeout_arg,
    max_feed_bytes_arg,
    blocker_only_arg,
    dry_run_arg,
    retain_captured_arg,
) = sys.argv[1:19]

ROOT = pathlib.Path(root_arg).resolve()
FORCE = force_arg == "true"
ALLOW_CAPTURED = allow_captured_arg == "true"
BLOCKER_ONLY = blocker_only_arg == "true"
DRY_RUN = dry_run_arg == "true"
RETAIN_CAPTURED = retain_captured_arg == "true"
FINAL_ROOT = final_root_arg.strip().rstrip("/")
ENVIRONMENT = environment_arg.strip()
APPROVAL_ARTIFACT = approval_arg.strip()
APPROVAL_SUMMARY = approval_summary_arg.strip()
ADMIN_BASE = admin_base_arg.strip().rstrip("/")
ADMIN_TOKEN_PRESENT = bool(admin_token_arg)
CONNECT_TIMEOUT = int(connect_timeout_arg)
REQUEST_TIMEOUT = int(request_timeout_arg)
MAX_FEED_BYTES = int(max_feed_bytes_arg)

BLOCKER_FILES = ("blocker.json", "blocker.md", "manifest.json", "manifest.md")
REAL_TOP_FILES = (
    "README.md",
    "approval.md",
    "dns-tls-redirect.md",
    "public-fetches.md",
    "validator-record.md",
    "proxy-config-summary.md",
    "redaction-notes.md",
    "manifest.json",
    "manifest.md",
    "SHA256SUMS.txt",
)
FEEDS = (
    ("feeds_json", "/public/feeds.json", "feeds.json"),
    ("schedule", "/public/gtfs/schedule.zip", "schedule.zip"),
    ("vehicle_positions", "/public/gtfsrt/vehicle_positions.pb", "vehicle_positions.pb"),
    ("trip_updates", "/public/gtfsrt/trip_updates.pb", "trip_updates.pb"),
    ("alerts", "/public/gtfsrt/alerts.pb", "alerts.pb"),
)
CLAIM_FLAGS = {
    "final_root_evidence_retained": False,
    "compliance_claimed": False,
    "consumer_statuses_changed": False,
    "consumer_contacted": False,
    "consumer_acceptance_claimed": False,
    "consumer_ingestion_claimed": False,
    "agency_adoption_claimed": False,
    "hosted_saas_claimed": False,
    "production_readiness_claimed": False,
    "sla_or_uptime_claimed": False,
    "vendor_compatibility_claimed": False,
    "production_grade_eta_claimed": False,
}

UNSAFE_PATTERNS = [
    re.compile(pattern, re.I)
    for pattern in (
        r"authorization\s*:",
        r"cookie\s*:",
        r"\bbearer\s+[A-Za-z0-9._~+/=-]{8,}",
        r"postgres(?:ql)?://",
        r"database_url\s*=",
        r"begin [a-z ]*private key",
        r"acme[_ -]?account",
        r"private[_ -]?key",
        r"secret\s*[:=]",
        r"password\s*[:=]",
        r"set-cookie\s*:",
        r"raw[_-](log|payload|diagnostic)",
        r"dns[_ -]?provider[_ -]?(token|secret|export)",
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


def safe_environment(value):
    if not re.fullmatch(r"[A-Za-z0-9._-]+", value):
        fail("FINAL_ROOT_ENVIRONMENT_NAME may contain only letters, digits, dot, underscore, and hyphen")
    if value in {".", ".."} or value.startswith("."):
        fail("FINAL_ROOT_ENVIRONMENT_NAME must not be dot, dot-dot, or a leading-dot name")


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


def is_relative_to(path, base):
    try:
        pathlib.Path(path).resolve(strict=False).relative_to(pathlib.Path(base).resolve(strict=False))
        return True
    except ValueError:
        return False


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def parse_url(raw, name):
    if not raw:
        return None
    parsed = urllib.parse.urlparse(raw)
    if parsed.username or parsed.password:
        fail(f"{name} must not contain credentials")
    if parsed.query or parsed.fragment:
        fail(f"{name} must not contain query strings or fragments")
    if not parsed.scheme or not parsed.netloc:
        fail(f"{name} must be an absolute URL")
    host = parsed.hostname or ""
    loopback = host in {"localhost", "127.0.0.1", "::1"} or host.startswith("127.")
    if parsed.scheme == "https":
        return parsed
    if parsed.scheme == "http" and loopback:
        return parsed
    fail(f"{name} must use HTTPS except loopback HTTP test roots")


def validate_output_dir():
    if RETAIN_CAPTURED:
        if not ALLOW_CAPTURED:
            fail("docs/evidence/captured writes require ALLOW_CAPTURED_EVIDENCE_WRITE=true")
        if not FINAL_ROOT:
            fail("--retain-captured requires FINAL_ROOT_BASE_URL")
        if not APPROVAL_ARTIFACT:
            fail("--retain-captured requires FINAL_ROOT_APPROVAL_ARTIFACT")
        default = ROOT / "docs" / "evidence" / "captured" / ENVIRONMENT / capture_date
    else:
        default = ROOT / ".cache" / "final-root-evidence" / timestamp
    raw = pathlib.Path(output_arg) if output_arg else default
    out = raw if raw.is_absolute() else ROOT / raw
    resolved = out.resolve(strict=False)
    if path_has_symlink(out):
        fail("OUTPUT_DIR must not contain symlink directories")
    if not is_relative_to(resolved, ROOT):
        fail("OUTPUT_DIR must stay inside the repository")
    captured_root = ROOT / "docs" / "evidence" / "captured"
    if is_relative_to(resolved, captured_root):
        if not (RETAIN_CAPTURED and ALLOW_CAPTURED and FINAL_ROOT and APPROVAL_ARTIFACT):
            fail("writes under docs/evidence/captured require --retain-captured, ALLOW_CAPTURED_EVIDENCE_WRITE=true, FINAL_ROOT_BASE_URL, and FINAL_ROOT_APPROVAL_ARTIFACT")
    elif not is_relative_to(resolved, ROOT / ".cache"):
        fail("OUTPUT_DIR must resolve under .cache unless --retain-captured writes to docs/evidence/captured")
    if any(part == ".." for part in raw.parts):
        fail("OUTPUT_DIR must not contain traversal components")
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
    return [pat.pattern for pat in UNSAFE_PATTERNS if pat.search(text)]


def read_public_safe_approval(path_text):
    if not path_text:
        return None, "missing approval artifact"
    path = pathlib.Path(path_text)
    full = path if path.is_absolute() else ROOT / path
    if path_has_symlink(full):
        fail("FINAL_ROOT_APPROVAL_ARTIFACT must not contain symlink directories")
    if not full.exists() or not full.is_file():
        return None, "approval artifact is not readable"
    if full.suffix.lower() in {".log", ".sql", ".dump", ".bak", ".backup", ".pem", ".key", ".p12", ".pfx"}:
        fail("FINAL_ROOT_APPROVAL_ARTIFACT must be a redacted public-safe text artifact, not raw logs, dumps, backups, or key material")
    data = full.read_bytes()
    if len(data) > 1024 * 1024:
        fail("FINAL_ROOT_APPROVAL_ARTIFACT exceeds 1 MiB")
    text = data.decode("utf-8", errors="replace")
    hits = unsafe_text(text)
    if hits:
        fail("FINAL_ROOT_APPROVAL_ARTIFACT contains unsafe private strings")
    return (full, text), None


def write_json(path, data):
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def write_blocker(out, reason, root_valid):
    now = dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    flags = CLAIM_FLAGS.copy()
    data = {
        "schema": "final-root-blocker.v1",
        "generated_at": now,
        "capture_date_utc": capture_date,
        "mode": "blocker",
        "dry_run": DRY_RUN,
        "blocker_only": BLOCKER_ONLY,
        "reason": reason,
        "final_root_available": bool(FINAL_ROOT and root_valid),
        "approval_artifact_available": False,
        "real_final_root_evidence_retained": False,
        "captured_evidence_directory_created": False,
        "claim_flags": flags,
        "consumer_tracker_changed": False,
    }
    manifest = {
        "schema": "final-root-manifest.v1",
        "packet_type": "blocker",
        "generated_at": now,
        "packet_dir": rel(out),
        "expected_files": list(BLOCKER_FILES),
        "retained_captured": False,
        "docs_evidence_captured_changed": False,
        "claim_flags": flags,
    }
    write_json(out / "blocker.json", data)
    (out / "blocker.md").write_text(
        "# Final Public Root Evidence Blocker\n\n"
        f"- Generated UTC: {now}\n"
        f"- Capture date UTC: {capture_date}\n"
        f"- Reason: {reason}\n"
        "- Real final root evidence retained: false\n"
        "- Captured evidence directory created: false\n"
        "- Consumer tracker changed: false\n\n"
        "No real final public root and redacted approval artifact were available for this run. "
        "This blocker packet is not evidence of final-root approval, compliance, consumer acceptance, "
        "consumer ingestion, agency adoption, hosted SaaS availability, production readiness, SLA/uptime, "
        "vendor compatibility, or production-grade ETA quality.\n",
        encoding="utf-8",
    )
    write_json(out / "manifest.json", manifest)
    (out / "manifest.md").write_text(
        "# Final Root Blocker Manifest\n\n"
        f"- Packet type: blocker\n- Packet directory: `{rel(out)}`\n"
        f"- Expected files: `{', '.join(BLOCKER_FILES)}`\n"
        "- Captured evidence retained: false\n- Claim flags: all false\n",
        encoding="utf-8",
    )


def fetch_limited(url):
    with tempfile.TemporaryDirectory(prefix="open-transit-rt-final-root-fetch.") as tmp:
        tmp_path = pathlib.Path(tmp)
        headers_path = tmp_path / "headers.txt"
        body_path = tmp_path / "body.bin"
        cmd = [
            "curl",
            "-sS",
            "--connect-timeout",
            str(CONNECT_TIMEOUT),
            "--max-time",
            str(REQUEST_TIMEOUT),
            "--max-filesize",
            str(MAX_FEED_BYTES),
            "-A",
            "OpenTransitRT-FinalRootEvidence/1.0",
            "-D",
            str(headers_path),
            "-o",
            str(body_path),
            "-w",
            "status=%{http_code}\nurl_effective=%{url_effective}\n",
            url,
        ]
        proc = subprocess.run(cmd, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if proc.returncode != 0:
            raise RuntimeError(f"curl failed with exit {proc.returncode}")
        data = body_path.read_bytes() if body_path.exists() else b""
        if len(data) > MAX_FEED_BYTES:
            raise RuntimeError("artifact exceeds MAX_FEED_BYTES")
        status = 0
        effective = url
        for line in proc.stdout.splitlines():
            if line.startswith("status="):
                try:
                    status = int(line.split("=", 1)[1])
                except ValueError:
                    status = 0
            elif line.startswith("url_effective="):
                effective = line.split("=", 1)[1]
        headers = {}
        if headers_path.exists():
            for line in headers_path.read_text(encoding="utf-8", errors="replace").splitlines():
                if ":" in line:
                    key, value = line.split(":", 1)
                    headers[key.strip()] = value.strip()
        return status, headers, data, effective


def sha256_file(path):
    h = hashlib.sha256()
    with path.open("rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def write_dns_tls(out, parsed):
    host = parsed.hostname or ""
    dns_dir = out / "artifacts" / "dns"
    tls_dir = out / "artifacts" / "tls"
    rows = []
    try:
        infos = socket.getaddrinfo(host, parsed.port or (443 if parsed.scheme == "https" else 80), proto=socket.IPPROTO_TCP)
        addrs = sorted({info[4][0] for info in infos})
        (dns_dir / "host-resolution.txt").write_text("\n".join(addrs) + "\n", encoding="utf-8")
        rows.append(("DNS", "captured", "artifacts/dns/host-resolution.txt"))
    except Exception as exc:
        (dns_dir / "host-resolution.txt").write_text(f"unavailable: {type(exc).__name__}\n", encoding="utf-8")
        rows.append(("DNS", "unavailable", "artifacts/dns/host-resolution.txt"))

    if parsed.scheme == "https":
        try:
            ctx = ssl.create_default_context()
            with socket.create_connection((host, parsed.port or 443), timeout=CONNECT_TIMEOUT) as sock:
                with ctx.wrap_socket(sock, server_hostname=host) as ssock:
                    cert = ssock.getpeercert()
            lines = [
                f"subject={cert.get('subject')}",
                f"issuer={cert.get('issuer')}",
                f"notBefore={cert.get('notBefore')}",
                f"notAfter={cert.get('notAfter')}",
                f"subjectAltName={cert.get('subjectAltName')}",
            ]
            (tls_dir / "certificate.txt").write_text("\n".join(lines) + "\n", encoding="utf-8")
            rows.append(("TLS certificate", "captured", "artifacts/tls/certificate.txt"))
        except Exception as exc:
            (tls_dir / "certificate.txt").write_text(f"unavailable: {type(exc).__name__}\n", encoding="utf-8")
            rows.append(("TLS certificate", "unavailable", "artifacts/tls/certificate.txt"))
        http_url = urllib.parse.urlunparse(("http", parsed.netloc, parsed.path + "/public/feeds.json", "", "", ""))
        try:
            status, headers, body, effective = fetch_limited(http_url)
            lines = [f"status={status}", f"url_effective={effective}"] + [f"{k}: {v}" for k, v in headers.items()]
            (tls_dir / "http-redirect-headers.txt").write_text("\n".join(lines) + "\n", encoding="utf-8")
            rows.append(("HTTP redirect", "captured", "artifacts/tls/http-redirect-headers.txt"))
        except Exception as exc:
            (tls_dir / "http-redirect-headers.txt").write_text(f"unavailable: {type(exc).__name__}\n", encoding="utf-8")
            rows.append(("HTTP redirect", "unavailable", "artifacts/tls/http-redirect-headers.txt"))
    else:
        (tls_dir / "certificate.txt").write_text("not_applicable: loopback HTTP test root\n", encoding="utf-8")
        (tls_dir / "http-redirect-headers.txt").write_text("not_applicable: loopback HTTP test root\n", encoding="utf-8")
        rows.append(("TLS certificate", "not_applicable_loopback", "artifacts/tls/certificate.txt"))
        rows.append(("HTTP redirect", "not_applicable_loopback", "artifacts/tls/http-redirect-headers.txt"))
    return rows


def validate_json_if_needed(path):
    if path.name.endswith(".json"):
        json.loads(path.read_text(encoding="utf-8"))


def write_validator_records(out):
    validation_dir = out / "artifacts" / "validation"
    records = []
    if not (ADMIN_BASE and ADMIN_TOKEN_PRESENT):
        for feed_id, _, _ in FEEDS[1:]:
            artifact = validation_dir / f"validate-{feed_id}.json"
            write_json(artifact, {"feed_type": feed_id, "status": "unavailable", "reason": "ADMIN_BASE_URL and ADMIN_TOKEN not configured"})
            records.append((feed_id, "unavailable", f"artifacts/validation/{artifact.name}"))
        return records
    admin_parsed = parse_url(ADMIN_BASE, "ADMIN_BASE_URL")
    for feed_id, _, _ in FEEDS[1:]:
        validator_id = "static-mobilitydata" if feed_id == "schedule" else "realtime-mobilitydata"
        artifact = validation_dir / f"validate-{feed_id}.json"
        payload = json.dumps({"validator_id": validator_id, "feed_type": feed_id}).encode("utf-8")
        req = urllib.request.Request(
            ADMIN_BASE + "/admin/validation/run",
            data=payload,
            method="POST",
            headers={
                "Authorization": "Bearer " + admin_token_arg,
                "Content-Type": "application/json",
                "User-Agent": "OpenTransitRT-FinalRootEvidence/1.0",
            },
        )
        try:
            with urllib.request.urlopen(req, timeout=REQUEST_TIMEOUT) as resp:
                data = resp.read(min(MAX_FEED_BYTES, 1024 * 1024))
            artifact.write_bytes(data)
            validate_json_if_needed(artifact)
            try:
                status = json.loads(data.decode("utf-8", errors="replace")).get("status", "unknown")
            except Exception:
                status = "unknown"
        except Exception as exc:
            write_json(artifact, {"feed_type": feed_id, "status": "unavailable", "reason": type(exc).__name__, "admin_url_scheme": admin_parsed.scheme})
            status = "unavailable"
        records.append((feed_id, status, f"artifacts/validation/{artifact.name}"))
    return records


def scan_output_for_unsafe(out):
    for path in out.rglob("*"):
        if not path.is_file() or path.name == "SHA256SUMS.txt":
            continue
        data = path.read_bytes()
        if b"\x00" in data:
            continue
        text = data.decode("utf-8", errors="ignore")
        hits = unsafe_text(text)
        if hits:
            fail(f"unsafe private string detected in generated packet: {rel(path)}")


def write_sha256s(out):
    entries = []
    for path in sorted(out.rglob("*")):
        if not path.is_file() or path.name == "SHA256SUMS.txt":
            continue
        entries.append((sha256_file(path), path.relative_to(out).as_posix()))
    (out / "SHA256SUMS.txt").write_text("".join(f"{h}  {p}\n" for h, p in entries), encoding="utf-8")


def write_real_packet(out, parsed, approval):
    now = dt.datetime.now(dt.timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")
    for sub in ("artifacts/public", "artifacts/tls", "artifacts/dns", "artifacts/validation", "artifacts/operator-supplied"):
        (out / sub).mkdir(parents=True, exist_ok=True)
    approval_path, approval_text = approval
    copied_approval = out / "artifacts" / "operator-supplied" / "final-root-approval-artifact.md"
    copied_approval.write_text(approval_text, encoding="utf-8")

    dns_tls_rows = write_dns_tls(out, parsed)
    fetch_rows = []
    for feed_id, path, artifact_name in FEEDS:
        url = FINAL_ROOT + path
        artifact = out / "artifacts" / "public" / artifact_name
        headers_artifact = out / "artifacts" / "public" / f"{artifact_name}.headers.txt"
        summary_artifact = out / "artifacts" / "public" / f"{artifact_name}.summary.json"
        try:
            status, headers, data, effective = fetch_limited(url)
            artifact.write_bytes(data)
            header_lines = [f"status={status}", f"url_effective={effective}"] + [f"{k}: {v}" for k, v in headers.items()]
            headers_artifact.write_text("\n".join(header_lines) + "\n", encoding="utf-8")
            hash_value = sha256_file(artifact)
            write_json(summary_artifact, {
                "feed_type": feed_id,
                "url": url,
                "url_effective": effective,
                "status": status,
                "bytes": len(data),
                "sha256": hash_value,
            })
            fetch_rows.append((feed_id, path, status, len(data), hash_value, f"artifacts/public/{artifact_name}"))
        except Exception as exc:
            artifact.write_bytes(b"")
            headers_artifact.write_text(f"unavailable: {type(exc).__name__}\n", encoding="utf-8")
            write_json(summary_artifact, {
                "feed_type": feed_id,
                "url": url,
                "status": "unavailable",
                "reason": type(exc).__name__,
            })
            fetch_rows.append((feed_id, path, "unavailable", 0, "", f"artifacts/public/{artifact_name}"))
    validator_rows = write_validator_records(out)

    claim_flags = CLAIM_FLAGS.copy()
    claim_flags["final_root_evidence_retained"] = RETAIN_CAPTURED
    manifest = {
        "schema": "final-root-manifest.v1",
        "packet_type": "real",
        "generated_at": now,
        "capture_date_utc": capture_date,
        "packet_dir": rel(out),
        "final_root_base_url": FINAL_ROOT,
        "environment_name": ENVIRONMENT,
        "approval_artifact": "artifacts/operator-supplied/final-root-approval-artifact.md",
        "approval_artifact_source": rel(approval_path),
        "approval_summary": APPROVAL_SUMMARY or "redacted approval artifact retained in packet",
        "retained_captured": RETAIN_CAPTURED,
        "public_fetches": [
            {"feed_type": row[0], "path": row[1], "status": row[2], "bytes": row[3], "sha256": row[4], "artifact": row[5]}
            for row in fetch_rows
        ],
        "validator_records": [
            {"feed_type": row[0], "status": row[1], "artifact": row[2]}
            for row in validator_rows
        ],
        "claim_flags": claim_flags,
        "consumer_tracker_changed": False,
        "limitations": [
            "Final-root evidence does not prove compliance.",
            "Final-root evidence does not prove consumer acceptance or ingestion.",
            "Final-root evidence does not prove agency adoption, hosted SaaS availability, production readiness, SLA/uptime, vendor compatibility, or production-grade ETA quality.",
        ],
    }

    (out / "README.md").write_text(
        f"# Final Public Root Evidence Packet\n\n"
        f"- Environment: `{ENVIRONMENT}`\n- Final root: `{FINAL_ROOT}`\n- Capture date UTC: {capture_date}\n"
        f"- Generated UTC: {now}\n- Retained under docs/evidence/captured: {str(RETAIN_CAPTURED).lower()}\n\n"
        "## Narrow Supported Claim\n\n"
        "Final-root approval and technical fetch/validator evidence exists for the exact recorded root and date.\n\n"
        "## Claim Boundary\n\n"
        "This packet does not prove compliance, consumer acceptance, consumer ingestion, agency adoption, hosted SaaS availability, production readiness, SLA/uptime, vendor compatibility, or ETA quality.\n",
        encoding="utf-8",
    )
    (out / "approval.md").write_text(
        f"# Final Root Approval\n\n- Summary: {APPROVAL_SUMMARY or 'redacted approval artifact retained'}\n"
        "- Artifact: `artifacts/operator-supplied/final-root-approval-artifact.md`\n"
        "- Approval artifact was scanned for disallowed private strings before retention.\n",
        encoding="utf-8",
    )
    (out / "dns-tls-redirect.md").write_text(
        "# DNS, TLS, And Redirect Evidence\n\n"
        "| Area | Status | Artifact |\n| --- | --- | --- |\n"
        + "".join(f"| {area} | `{status}` | `{artifact}` |\n" for area, status, artifact in dns_tls_rows),
        encoding="utf-8",
    )
    (out / "public-fetches.md").write_text(
        "# Public Fetches\n\n"
        f"- Final root: `{FINAL_ROOT}`\n\n"
        "| Feed | Path | Status | Bytes | SHA-256 | Artifact |\n| --- | --- | ---: | ---: | --- | --- |\n"
        + "".join(f"| `{feed}` | `{path}` | `{status}` | {size} | `{sha}` | `{artifact}` |\n" for feed, path, status, size, sha, artifact in fetch_rows),
        encoding="utf-8",
    )
    (out / "validator-record.md").write_text(
        "# Validator Record\n\n"
        "| Feed | Status | Artifact |\n| --- | --- | --- |\n"
        + "".join(f"| `{feed}` | `{status}` | `{artifact}` |\n" for feed, status, artifact in validator_rows)
        + "\nValidator records must be `passed` before a real audit can support the narrow final-root evidence claim.\n",
        encoding="utf-8",
    )
    (out / "proxy-config-summary.md").write_text(
        "# Proxy And Public Edge Summary\n\n"
        "- Public feed paths expected: `/public/feeds.json`, `/public/gtfs/schedule.zip`, `/public/gtfsrt/vehicle_positions.pb`, `/public/gtfsrt/trip_updates.pb`, `/public/gtfsrt/alerts.pb`.\n"
        "- Admin, debug, telemetry ingest, and private diagnostics are not evidenced as public routes by this packet.\n"
        "- Add only redacted public-safe route summaries under `artifacts/operator-supplied/` when available.\n",
        encoding="utf-8",
    )
    (out / "redaction-notes.md").write_text(
        "# Redaction Notes\n\n"
        "- Approval artifact was supplied as a redacted public-safe file.\n"
        "- Authorization and Cookie headers are not written.\n"
        "- Database URLs, private TLS keys, ACME material, raw logs, private payloads, and unredacted diagnostics are disallowed.\n"
        "- Generated artifacts were scanned for configured unsafe strings before checksum generation.\n",
        encoding="utf-8",
    )
    write_json(out / "manifest.json", manifest)
    (out / "manifest.md").write_text(
        "# Final Root Manifest\n\n"
        f"- Packet type: real\n- Packet directory: `{rel(out)}`\n- Final root: `{FINAL_ROOT}`\n"
        f"- Retained captured evidence: {str(RETAIN_CAPTURED).lower()}\n- Consumer tracker changed: false\n",
        encoding="utf-8",
    )
    scan_output_for_unsafe(out)
    write_sha256s(out)


bool_value("FORCE", force_arg)
bool_value("ALLOW_CAPTURED_EVIDENCE_WRITE", allow_captured_arg)
positive_int("CONNECT_TIMEOUT_SECONDS", connect_timeout_arg)
positive_int("REQUEST_TIMEOUT_SECONDS", request_timeout_arg)
positive_int("MAX_FEED_BYTES", max_feed_bytes_arg)
safe_environment(ENVIRONMENT)
out = validate_output_dir()

root_parsed = None
root_valid = False
if FINAL_ROOT:
    root_parsed = parse_url(FINAL_ROOT, "FINAL_ROOT_BASE_URL")
    root_valid = True
if ADMIN_BASE:
    parse_url(ADMIN_BASE, "ADMIN_BASE_URL")

approval, approval_problem = read_public_safe_approval(APPROVAL_ARTIFACT)

if BLOCKER_ONLY or DRY_RUN or not (FINAL_ROOT and root_valid and approval):
    reason_parts = []
    if BLOCKER_ONLY:
        reason_parts.append("blocker-only requested")
    if DRY_RUN:
        reason_parts.append("dry-run requested")
    if not FINAL_ROOT:
        reason_parts.append("no final root provided")
    if not approval:
        reason_parts.append(approval_problem or "no approval artifact available")
    reason = "; ".join(reason_parts) or "no real final-root evidence retained"
    write_blocker(out, reason, root_valid)
    print(f"final-root blocker packet: {rel(out)}")
    sys.exit(0)

write_real_packet(out, root_parsed, approval)
print(f"final-root evidence packet: {rel(out)}")
PY
