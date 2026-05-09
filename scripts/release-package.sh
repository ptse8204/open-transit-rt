#!/usr/bin/env sh
set -eu
umask 077

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

TIMESTAMP="$(date -u '+%Y%m%dT%H%M%SZ')"
RELEASE_PACKAGE_VERSION="${RELEASE_PACKAGE_VERSION:-}"
RELEASE_PACKAGE_OUTPUT_DIR="${RELEASE_PACKAGE_OUTPUT_DIR:-}"
RELEASE_PACKAGE_FORCE="${RELEASE_PACKAGE_FORCE:-false}"
RELEASE_PACKAGE_ALLOW_DIRTY="${RELEASE_PACKAGE_ALLOW_DIRTY:-false}"
RELEASE_PACKAGE_STRICT="${RELEASE_PACKAGE_STRICT:-false}"
RELEASE_PACKAGE_IMAGE_TAG="${RELEASE_PACKAGE_IMAGE_TAG:-}"

usage() {
  cat <<'EOF'
Usage:
  scripts/release-package.sh [--help]

Environment:
  RELEASE_PACKAGE_VERSION      Optional release version/tag, default git describe --tags --always.
  RELEASE_PACKAGE_OUTPUT_DIR   Default .cache/release-package/<safe-version>.
  RELEASE_PACKAGE_FORCE        true|false; allow reuse of non-empty output.
  RELEASE_PACKAGE_ALLOW_DIRTY  true|false; allow dirty checkout for local diagnostics.
  RELEASE_PACKAGE_STRICT       true|false; fail on unavailable optional metadata.
  RELEASE_PACKAGE_IMAGE_TAG    Optional local Docker image tag to inspect; never builds or pushes.

Safety:
  This helper creates local release review artifacts under ignored .cache by
  default. It creates a source archive from git HEAD, checksums, SBOM summary,
  and provenance metadata. It does not publish artifacts, push images, contact
  registries or consumers, create retained evidence, change consumer statuses,
  or claim hosted service availability, production readiness, SLA/uptime,
  compliance, agency adoption, consumer acceptance, marketplace approval,
  vendor compatibility, or production-grade ETA quality.
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

python3 - "$ROOT_DIR" "$TIMESTAMP" "$RELEASE_PACKAGE_VERSION" \
  "$RELEASE_PACKAGE_OUTPUT_DIR" "$RELEASE_PACKAGE_FORCE" \
  "$RELEASE_PACKAGE_ALLOW_DIRTY" "$RELEASE_PACKAGE_STRICT" \
  "$RELEASE_PACKAGE_IMAGE_TAG" <<'PY'
import gzip
import hashlib
import json
import os
import pathlib
import re
import shutil
import subprocess
import sys

(
    root_arg,
    timestamp,
    version_arg,
    output_arg,
    force_arg,
    allow_dirty_arg,
    strict_arg,
    image_tag_arg,
) = sys.argv[1:9]

ROOT = pathlib.Path(root_arg).resolve()
FORCE = force_arg == "true"
ALLOW_DIRTY = allow_dirty_arg == "true"
STRICT = strict_arg == "true"
IMAGE_TAG = image_tag_arg.strip()
MAX_GO_LIST_BYTES = 4 * 1024 * 1024

ROOT_FILES = (
    "summary.json",
    "summary.md",
    "manifest.json",
    "manifest.md",
    "provenance.json",
    "provenance.md",
    "sbom.json",
    "image.json",
)
CLAIM_FLAGS = {
    "hosted_saas_claimed": False,
    "hosted_service_claimed": False,
    "production_readiness_claimed": False,
    "sla_or_uptime_claimed": False,
    "compliance_claimed": False,
    "agency_adoption_claimed": False,
    "consumer_acceptance_claimed": False,
    "consumer_ingestion_claimed": False,
    "vendor_compatibility_claimed": False,
    "marketplace_approval_claimed": False,
    "production_grade_eta_claimed": False,
    "retained_evidence_created": False,
    "consumer_statuses_changed": False,
}


def fail(message):
    raise SystemExit(f"ERROR: {message}")


def bool_arg(name, value):
    if value not in {"true", "false"}:
        fail(f"{name} must be true or false")


for name, value in (
    ("RELEASE_PACKAGE_FORCE", force_arg),
    ("RELEASE_PACKAGE_ALLOW_DIRTY", allow_dirty_arg),
    ("RELEASE_PACKAGE_STRICT", strict_arg),
):
    bool_arg(name, value)


def run(args, *, check=True, timeout=20, max_bytes=None):
    proc = subprocess.run(args, cwd=ROOT, text=False, stdout=subprocess.PIPE, stderr=subprocess.PIPE, timeout=timeout)
    if max_bytes is not None and len(proc.stdout) > max_bytes:
        if check:
            fail(f"{args[0]} output exceeded byte cap")
        proc.stdout = proc.stdout[:max_bytes]
    if check and proc.returncode != 0:
        stderr = proc.stderr.decode("utf-8", "replace").strip()
        fail(f"command failed: {' '.join(args)}: {stderr}")
    return proc


def text(args, *, check=True, timeout=20):
    proc = run(args, check=check, timeout=timeout)
    return proc.stdout.decode("utf-8", "replace").strip(), proc


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


def is_evidence_like(path):
    raw = str(path).replace("\\", "/").lower()
    parts = [part.lower() for part in pathlib.Path(path).parts]
    return "docs/evidence" in raw or "evidence" in parts or "proof" in parts or "submission" in parts


def rel(path):
    try:
        return pathlib.Path(path).resolve(strict=False).relative_to(ROOT).as_posix()
    except ValueError:
        return "<outside-repo>"


def safe_version(raw):
    value = raw.strip()
    if not value:
        fail("release package version is empty")
    if not re.fullmatch(r"[A-Za-z0-9][A-Za-z0-9._-]{0,80}", value):
        fail("release package version must be a conservative tag-like value")
    return value


def resolve_output_dir(version):
    if output_arg:
        out = pathlib.Path(output_arg)
    else:
        out = pathlib.Path(".cache") / "release-package" / version
    target = out if out.is_absolute() else ROOT / out
    resolved = target.resolve(strict=False)
    cache = (ROOT / ".cache").resolve(strict=False)
    if not is_relative_to(resolved, cache):
        fail("RELEASE_PACKAGE_OUTPUT_DIR must resolve under repo .cache")
    if is_evidence_like(target) or is_evidence_like(resolved):
        fail("RELEASE_PACKAGE_OUTPUT_DIR must not be evidence-like or under docs/evidence")
    if has_symlink(target):
        fail("RELEASE_PACKAGE_OUTPUT_DIR must not contain symlink directories")
    if resolved.exists() and not resolved.is_dir():
        fail("RELEASE_PACKAGE_OUTPUT_DIR must be a directory")
    if resolved.exists() and any(resolved.iterdir()):
        if not FORCE:
            fail("RELEASE_PACKAGE_OUTPUT_DIR exists and is non-empty; set RELEASE_PACKAGE_FORCE=true to reuse")
        for child in resolved.iterdir():
            if child.is_symlink() or child.is_file():
                child.unlink()
            else:
                shutil.rmtree(child)
    resolved.mkdir(parents=True, exist_ok=True)
    os.chmod(resolved, 0o700)
    (resolved / "artifacts").mkdir(mode=0o700)
    (resolved / "checksums").mkdir(mode=0o700)
    return resolved


def write_json(path, data):
    path.write_text(json.dumps(data, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def sha256_file(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1024 * 1024), b""):
            h.update(chunk)
    return h.hexdigest()


def parse_go_modules(raw):
    decoder = json.JSONDecoder()
    text_data = raw.decode("utf-8", "replace")
    idx = 0
    modules = []
    while idx < len(text_data):
        while idx < len(text_data) and text_data[idx].isspace():
            idx += 1
        if idx >= len(text_data):
            break
        obj, idx = decoder.raw_decode(text_data, idx)
        modules.append({
            "path": obj.get("Path", ""),
            "version": obj.get("Version", ""),
            "sum": obj.get("Sum", ""),
            "replace": obj.get("Replace", {}).get("Path", "") if isinstance(obj.get("Replace"), dict) else "",
            "main": bool(obj.get("Main")),
        })
    return modules


version_raw = version_arg.strip()
if not version_raw:
    version_raw, _ = text(["git", "describe", "--tags", "--always"], timeout=10)
version = safe_version(version_raw)
output = resolve_output_dir(version)

commit, _ = text(["git", "rev-parse", "HEAD"], timeout=10)
describe, _ = text(["git", "describe", "--tags", "--always", "--dirty"], timeout=10)
status_text, _ = text(["git", "status", "--porcelain"], timeout=10)
dirty = bool(status_text)
if dirty and not ALLOW_DIRTY:
    fail("working tree is dirty; set RELEASE_PACKAGE_ALLOW_DIRTY=true only for local diagnostics")

artifact_name = f"open-transit-rt-{version}.source.tar.gz"
artifact = output / "artifacts" / artifact_name
archive = subprocess.Popen(
    ["git", "archive", "--format=tar", f"--prefix=open-transit-rt-{version}/", "HEAD"],
    cwd=ROOT,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
)
with artifact.open("wb") as raw_out:
    with gzip.GzipFile(filename="", mode="wb", fileobj=raw_out, mtime=0) as gz:
        assert archive.stdout is not None
        shutil.copyfileobj(archive.stdout, gz)
stderr = archive.stderr.read().decode("utf-8", "replace") if archive.stderr else ""
code = archive.wait(timeout=30)
if code != 0:
    fail(f"git archive failed: {stderr.strip()}")

go_version, _ = text(["go", "version"], check=False, timeout=10)
git_version, _ = text(["git", "--version"], check=False, timeout=10)
docker_version, _ = text(["docker", "--version"], check=False, timeout=10)

sbom_status = "present"
sbom_error = ""
modules = []
go_proc = run(["go", "list", "-m", "-json", "all"], check=False, timeout=30, max_bytes=MAX_GO_LIST_BYTES)
if go_proc.returncode == 0:
    try:
        modules = parse_go_modules(go_proc.stdout)
    except Exception as exc:
        sbom_status = "unavailable"
        sbom_error = type(exc).__name__
else:
    sbom_status = "unavailable"
    sbom_error = go_proc.stderr.decode("utf-8", "replace").strip()[:240]
if STRICT and sbom_status != "present":
    fail(f"SBOM generation unavailable: {sbom_error}")

image_status = "not_configured"
image = {
    "schema": "open-transit-rt-release-image.v1",
    "status": image_status,
    "image_tag": IMAGE_TAG,
    "claim_flags": CLAIM_FLAGS,
}
if IMAGE_TAG:
    proc = run(["docker", "image", "inspect", IMAGE_TAG], check=False, timeout=20, max_bytes=MAX_GO_LIST_BYTES)
    if proc.returncode == 0:
        payload = json.loads(proc.stdout.decode("utf-8", "replace"))
        first = payload[0] if payload else {}
        image.update({
            "status": "local_image_found",
            "image_id": first.get("Id", ""),
            "repo_tags": first.get("RepoTags", []) or [],
            "repo_digests": first.get("RepoDigests", []) or [],
            "created": first.get("Created", ""),
        })
    else:
        image.update({
            "status": "local_image_unavailable",
            "error": proc.stderr.decode("utf-8", "replace").strip()[:240],
        })
        if STRICT:
            fail("local Docker image metadata unavailable")

release_ready = not dirty and sbom_status == "present"
status = "release_ready" if release_ready else "not_release_ready"
summary = {
    "schema": "open-transit-rt-release-package-summary.v1",
    "generated_at": timestamp,
    "version": version,
    "status": status,
    "release_ready": release_ready,
    "git_commit": commit,
    "git_describe": describe,
    "git_dirty": dirty,
    "output_directory": rel(output),
    "source_archive": f"artifacts/{artifact_name}",
    "sbom_status": sbom_status,
    "image_status": image["status"],
    "claim_flags": CLAIM_FLAGS,
    "boundaries": {
        "published_artifacts": False,
        "registry_push": False,
        "retained_evidence_created": False,
        "consumer_statuses_changed": False,
    },
}
manifest = {
    "schema": "open-transit-rt-release-package-manifest.v1",
    "generated_at": timestamp,
    "version": version,
    "files": list(ROOT_FILES) + [f"artifacts/{artifact_name}", "checksums/SHA256SUMS.txt"],
    "claim_flags": CLAIM_FLAGS,
}
provenance = {
    "schema": "open-transit-rt-release-provenance.v1",
    "generated_at": timestamp,
    "version": version,
    "source": {
        "git_commit": commit,
        "git_describe": describe,
        "git_dirty": dirty,
        "source_archive_from": "git archive HEAD",
    },
    "builder": {
        "type": "local",
        "go_version": go_version,
        "git_version": git_version,
        "docker_version": docker_version,
    },
    "claim_flags": CLAIM_FLAGS,
}
sbom = {
    "schema": "open-transit-rt-release-sbom.v1",
    "generated_at": timestamp,
    "status": sbom_status,
    "error": sbom_error,
    "module_count": len(modules),
    "modules": modules,
    "claim_flags": CLAIM_FLAGS,
}

write_json(output / "summary.json", summary)
write_json(output / "manifest.json", manifest)
write_json(output / "provenance.json", provenance)
write_json(output / "sbom.json", sbom)
write_json(output / "image.json", image)
(output / "summary.md").write_text(
    "# Release Package Summary\n\n"
    f"- Version: `{version}`\n"
    f"- Status: `{status}`\n"
    f"- Source archive: `artifacts/{artifact_name}`\n"
    f"- SBOM status: `{sbom_status}`\n"
    f"- Local image status: `{image['status']}`\n"
    "- Published artifacts: false\n"
    "- Retained evidence created: false\n"
    "- Consumer statuses changed: false\n",
    encoding="utf-8",
)
(output / "manifest.md").write_text(
    "# Release Package Manifest\n\n"
    f"- Generated at: {timestamp}\n"
    f"- Output directory: `{rel(output)}`\n"
    "- Includes source archive, checksums, local SBOM summary, image metadata, and provenance metadata.\n",
    encoding="utf-8",
)
(output / "provenance.md").write_text(
    "# Release Provenance\n\n"
    f"- Git commit: `{commit}`\n"
    f"- Git describe: `{describe}`\n"
    f"- Dirty checkout: `{str(dirty).lower()}`\n"
    "- Source archive command: `git archive HEAD`\n",
    encoding="utf-8",
)

checksums = []
for path in sorted(output.rglob("*")):
    if not path.is_file() or path.name == "SHA256SUMS.txt":
        continue
    checksums.append((sha256_file(path), path.relative_to(output).as_posix()))
(output / "checksums" / "SHA256SUMS.txt").write_text(
    "".join(f"{digest}  {name}\n" for digest, name in checksums),
    encoding="utf-8",
)

actual_root = sorted(child.name for child in output.iterdir())
expected_root = sorted(list(ROOT_FILES) + ["artifacts", "checksums"])
if actual_root != expected_root:
    fail(f"unexpected release package root files: {actual_root}")

print(f"release package: {rel(output)}")
PY
