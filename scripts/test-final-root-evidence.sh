#!/usr/bin/env sh
set -eu

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

python3 - "$ROOT_DIR" <<'PY'
import contextlib
import hashlib
import http.server
import json
import os
import pathlib
import shutil
import socketserver
import subprocess
import tempfile
import threading

ROOT = pathlib.Path(__import__("sys").argv[1]).resolve()
BASE = ROOT / ".cache" / "final-root-evidence-tests"
COLLECT = ROOT / "scripts" / "collect-final-root-evidence.sh"
AUDIT = ROOT / "scripts" / "audit-final-root-evidence.sh"


def run(args, env=None, ok=True):
    merged = os.environ.copy()
    if env:
        merged.update({k: str(v) for k, v in env.items()})
    proc = subprocess.run(args, cwd=ROOT, env=merged, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if ok and proc.returncode != 0:
        raise AssertionError(f"command failed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    if not ok and proc.returncode == 0:
        raise AssertionError(f"command unexpectedly passed: {args}\nstdout={proc.stdout}\nstderr={proc.stderr}")
    return proc


def assert_files_exact(path, expected):
    names = sorted(p.name for p in pathlib.Path(path).iterdir())
    if names != sorted(expected):
        raise AssertionError(f"{path} files mismatch: {names} != {sorted(expected)}")


def sha(path):
    h = hashlib.sha256()
    with pathlib.Path(path).open("rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        bodies = {
            "/public/feeds.json": (b'{"feeds":[]}\n', "application/json"),
            "/public/gtfs/schedule.zip": (b"PK\x05\x06" + b"\x00" * 18, "application/zip"),
            "/public/gtfsrt/vehicle_positions.pb": (b"\x0a\x03rtv", "application/octet-stream"),
            "/public/gtfsrt/trip_updates.pb": (b"\x0a\x03rtt", "application/octet-stream"),
            "/public/gtfsrt/alerts.pb": (b"\x0a\x03rta", "application/octet-stream"),
            "/public/gtfsrt/large.pb": (b"x" * (2 * 1024 * 1024), "application/octet-stream"),
        }
        body, content_type = bodies.get(self.path, (b"not found\n", "text/plain"))
        status = 200 if self.path in bodies else 404
        self.send_response(status)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *_):
        pass


@contextlib.contextmanager
def server():
    srv = socketserver.TCPServer(("127.0.0.1", 0), Handler)
    thread = threading.Thread(target=srv.serve_forever, daemon=True)
    thread.start()
    try:
        yield f"http://127.0.0.1:{srv.server_address[1]}"
    finally:
        srv.shutdown()
        srv.server_close()
        thread.join(timeout=5)


def write_approval(path):
    pathlib.Path(path).write_text(
        "# Redacted Approval\n\nApproved final public feed root for test fixture.\n",
        encoding="utf-8",
    )


def make_real_packet(base_url, out, approval):
    run([
        str(COLLECT),
    ], env={
        "FINAL_ROOT_BASE_URL": base_url,
        "FINAL_ROOT_APPROVAL_ARTIFACT": approval,
        "FINAL_ROOT_APPROVAL_SUMMARY": "redacted local test approval",
        "FINAL_ROOT_ENVIRONMENT_NAME": "local-test",
        "OUTPUT_DIR": out,
        "FORCE": "true",
    })


def write_checksums(packet):
    lines = []
    for path in sorted(packet.rglob("*")):
        if path.is_file() and path.name != "SHA256SUMS.txt":
            lines.append(f"{sha(path)}  {path.relative_to(packet).as_posix()}\n")
    (packet / "SHA256SUMS.txt").write_text("".join(lines), encoding="utf-8")


def make_auditable_real_fixture(packet):
    manifest = json.loads((packet / "manifest.json").read_text())
    for row in manifest["validator_records"]:
        row["status"] = "passed"
        (packet / row["artifact"]).write_text(json.dumps({"status": "passed"}) + "\n", encoding="utf-8")
    manifest["validator_records"] = manifest["validator_records"]
    (packet / "manifest.json").write_text(json.dumps(manifest, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    rows = [(row["feed_type"], row["status"], row["artifact"]) for row in manifest["validator_records"]]
    (packet / "validator-record.md").write_text(
        "# Validator Record\n\n"
        "| Feed | Status | Artifact |\n| --- | --- | --- |\n"
        + "".join(f"| `{feed}` | `{status}` | `{artifact}` |\n" for feed, status, artifact in rows)
        + "\nValidator records must be `passed` before a real audit can support the narrow final-root evidence claim.\n",
        encoding="utf-8",
    )
    write_checksums(packet)


if BASE.exists():
    shutil.rmtree(BASE)
BASE.mkdir(parents=True)

run([str(COLLECT), "--help"])
run([str(AUDIT), "--help"])

blocker = BASE / "blocker"
run([str(COLLECT), "--blocker-only"], env={"OUTPUT_DIR": blocker})
assert_files_exact(blocker, ["blocker.json", "blocker.md", "manifest.json", "manifest.md"])
run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": blocker, "AUDIT_MODE": "blocker"})

dry = BASE / "dry-run"
run([str(COLLECT), "--dry-run"], env={"OUTPUT_DIR": dry})
assert_files_exact(dry, ["blocker.json", "blocker.md", "manifest.json", "manifest.md"])

run([str(COLLECT)], env={"FINAL_ROOT_BASE_URL": "https://user:pass@example.com", "OUTPUT_DIR": BASE / "bad-url"}, ok=False)
run([str(COLLECT)], env={"FINAL_ROOT_BASE_URL": "https://example.com/path?token=x", "OUTPUT_DIR": BASE / "bad-query"}, ok=False)
run([str(COLLECT)], env={"FINAL_ROOT_BASE_URL": "http://example.com", "OUTPUT_DIR": BASE / "bad-http"}, ok=False)

symlink_target = BASE / "symlink-target"
symlink_target.mkdir()
symlink_path = BASE / "symlink-output"
symlink_path.symlink_to(symlink_target, target_is_directory=True)
run([str(COLLECT)], env={"OUTPUT_DIR": symlink_path}, ok=False)

approval = BASE / "approval.md"
write_approval(approval)
run([str(COLLECT), "--retain-captured"], env={
    "FINAL_ROOT_BASE_URL": "http://127.0.0.1:1",
    "FINAL_ROOT_APPROVAL_ARTIFACT": approval,
    "FINAL_ROOT_ENVIRONMENT_NAME": "local-test",
    "OUTPUT_DIR": BASE / "captured-denied",
}, ok=False)

with server() as base_url:
    real = BASE / "real"
    make_real_packet(base_url, real, approval)
    assert_files_exact(real, [
        "README.md",
        "SHA256SUMS.txt",
        "approval.md",
        "artifacts",
        "dns-tls-redirect.md",
        "manifest.json",
        "manifest.md",
        "proxy-config-summary.md",
        "public-fetches.md",
        "redaction-notes.md",
        "validator-record.md",
    ])
    if (ROOT / "docs" / "evidence" / "captured" / "local-test").exists():
        raise AssertionError("collector wrote docs/evidence/captured without retain opt-in")

    audit_fixture = BASE / "real-audit-pass"
    shutil.copytree(real, audit_fixture)
    make_auditable_real_fixture(audit_fixture)
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": audit_fixture, "AUDIT_MODE": "real"})

    non_200 = BASE / "non-200"
    shutil.copytree(audit_fixture, non_200)
    m = json.loads((non_200 / "manifest.json").read_text())
    m["public_fetches"][0]["status"] = 301
    (non_200 / "manifest.json").write_text(json.dumps(m, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    write_checksums(non_200)
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": non_200, "AUDIT_MODE": "real"}, ok=False)

    root_mismatch = BASE / "root-mismatch"
    shutil.copytree(audit_fixture, root_mismatch)
    m = json.loads((root_mismatch / "manifest.json").read_text())
    m["final_root_base_url"] = "http://example.com"
    (root_mismatch / "manifest.json").write_text(json.dumps(m) + "\n", encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": root_mismatch, "AUDIT_MODE": "real"}, ok=False)

    placeholder = BASE / "placeholder"
    shutil.copytree(audit_fixture, placeholder)
    (placeholder / "README.md").write_text("placeholder\n", encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": placeholder, "AUDIT_MODE": "real"}, ok=False)

    unsafe = BASE / "unsafe"
    shutil.copytree(audit_fixture, unsafe)
    (unsafe / "artifacts" / "operator-supplied" / "unsafe.txt").write_text("Authorization: Bearer abcdefghijklmnop\n", encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": unsafe, "AUDIT_MODE": "real"}, ok=False)

    checksum = BASE / "checksum"
    shutil.copytree(audit_fixture, checksum)
    (checksum / "artifacts" / "public" / "feeds.json").write_text('{"changed":true}\n', encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": checksum, "AUDIT_MODE": "real"}, ok=False)

    checksum_missing = BASE / "checksum-missing"
    shutil.copytree(audit_fixture, checksum_missing)
    checksum_lines = (checksum_missing / "SHA256SUMS.txt").read_text(encoding="utf-8").splitlines()
    if not checksum_lines:
        raise AssertionError("expected checksum entries")
    (checksum_missing / "SHA256SUMS.txt").write_text("\n".join(checksum_lines[1:]) + "\n", encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": checksum_missing, "AUDIT_MODE": "real"}, ok=False)

    validator_fail = BASE / "validator-fail"
    shutil.copytree(audit_fixture, validator_fail)
    m = json.loads((validator_fail / "manifest.json").read_text())
    m["validator_records"][0]["status"] = "failed"
    (validator_fail / "manifest.json").write_text(json.dumps(m) + "\n", encoding="utf-8")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": validator_fail, "AUDIT_MODE": "real"}, ok=False)

    max_bytes = BASE / "max-bytes"
    run([str(COLLECT)], env={
        "FINAL_ROOT_BASE_URL": base_url,
        "FINAL_ROOT_APPROVAL_ARTIFACT": approval,
        "FINAL_ROOT_APPROVAL_SUMMARY": "redacted local test approval",
        "FINAL_ROOT_ENVIRONMENT_NAME": "local-test",
        "OUTPUT_DIR": max_bytes,
        "FORCE": "true",
        "MAX_FEED_BYTES": "1",
    })
    if (ROOT / "docs" / "evidence" / "captured" / "local-test").exists():
        raise AssertionError("too-low MAX_FEED_BYTES run wrote docs/evidence/captured")
    run([str(AUDIT)], env={"FINAL_ROOT_PACKET_DIR": max_bytes, "AUDIT_MODE": "real"}, ok=False)

print("final-root evidence script tests passed")
PY
