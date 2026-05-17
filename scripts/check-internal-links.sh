#!/usr/bin/env sh
set -eu

python3 - "$@" <<'PY'
from __future__ import annotations

import pathlib
import re
import sys
from urllib.parse import unquote, urlparse


ROOT = pathlib.Path.cwd().resolve()
SITE_PREFIX = "/open-transit-rt/"
GITHUB_PREFIX = "https://github.com/ptse8204/open-transit-rt/"
PAGES_PREFIX = "https://ptse8204.github.io/open-transit-rt/"

MARKDOWN_LINK = re.compile(r"!?\[[^\]\n]*\]\(([^)\n]+)\)")
REFERENCE_LINK = re.compile(r"^\s*\[[^\]]+\]:\s*(\S+)", re.MULTILINE)
HTML_LINK = re.compile(r"""(?:href|src)=["']([^"']+)["']""", re.IGNORECASE)
FENCED_CODE = re.compile(r"^```.*?^```", re.MULTILINE | re.DOTALL)


def candidate_files() -> list[pathlib.Path]:
    files: list[pathlib.Path] = []
    root_readme = ROOT / "README.md"
    if root_readme.exists():
        files.append(root_readme)
    for base in ("docs", "wiki", "site"):
        root = ROOT / base
        if not root.exists():
            continue
        for path in root.rglob("*"):
            if path.suffix.lower() in {".md", ".html"}:
                files.append(path)
    return sorted(set(files))


def strip_wrappers(link: str) -> str:
    link = link.strip()
    if len(link) >= 2 and link[0] == "<" and link[-1] == ">":
        link = link[1:-1].strip()
    return link


def skip_link(link: str) -> bool:
    lower = link.lower()
    return (
        not link
        or link.startswith("#")
        or lower.startswith(("mailto:", "tel:", "data:", "javascript:"))
    )


def without_fragment_or_query(link: str) -> str:
    parsed = urlparse(link)
    if parsed.scheme or parsed.netloc:
        return link
    return link.split("#", 1)[0].split("?", 1)[0]


def exists_local(path: pathlib.Path) -> bool:
    return path.exists()


def resolve_local(source: pathlib.Path, link: str) -> tuple[pathlib.Path | None, str | None]:
    link_path = without_fragment_or_query(link)
    if not link_path:
        return None, None
    if link_path.startswith("/"):
        return None, None
    target = (source.parent / unquote(link_path)).resolve()
    try:
        target.relative_to(ROOT)
    except ValueError:
        return target, "escapes repository root"
    if source.relative_to(ROOT).parts[0] == "site":
        try:
            target.relative_to(ROOT / "site")
        except ValueError:
            return target, "site link escapes site/; use a GitHub source URL for repo docs"
    if exists_local(target):
        return target, None
    return target, "missing local target"


def resolve_github(link: str) -> tuple[pathlib.Path | None, str | None]:
    if not link.startswith(GITHUB_PREFIX):
        return None, None
    parsed = urlparse(link)
    parts = pathlib.PurePosixPath(parsed.path).parts
    # /ptse8204/open-transit-rt/blob/main/docs/index.md
    if len(parts) < 6 or parts[3] not in {"blob", "tree"} or parts[4] != "main":
        return None, None
    rel = pathlib.PurePosixPath(*parts[5:])
    target = (ROOT / pathlib.Path(*rel.parts)).resolve()
    try:
        target.relative_to(ROOT)
    except ValueError:
        return target, "GitHub source link escapes repository root"
    if parts[3] == "tree":
        return (target if target.is_dir() else target), (None if target.is_dir() else "missing GitHub source directory")
    return (target if target.is_file() else target), (None if target.is_file() else "missing GitHub source file")


def resolve_pages(link: str) -> tuple[pathlib.Path | None, str | None]:
    if not link.startswith(PAGES_PREFIX):
        return None, None
    parsed = urlparse(link)
    path = parsed.path
    if not path.startswith(SITE_PREFIX):
        return None, None
    rel = path[len(SITE_PREFIX):]
    if rel == "":
        rel = "index.html"
    if rel.endswith("/"):
        rel += "index.html"
    target = (ROOT / "site" / unquote(rel)).resolve()
    try:
        target.relative_to(ROOT / "site")
    except ValueError:
        return target, "GitHub Pages link escapes site/"
    return (target if target.exists() else target), (None if target.exists() else "missing GitHub Pages target")


def iter_links(path: pathlib.Path) -> list[tuple[str, int]]:
    text = path.read_text(encoding="utf-8")
    parse_text = text
    if path.suffix.lower() == ".md":
        parse_text = FENCED_CODE.sub(lambda match: "\n" * match.group(0).count("\n"), text)
    patterns = [MARKDOWN_LINK, REFERENCE_LINK]
    if path.suffix.lower() == ".html":
        patterns.append(HTML_LINK)
    links: list[tuple[str, int]] = []
    for pattern in patterns:
        for match in pattern.finditer(parse_text):
            link = strip_wrappers(match.group(1))
            line = parse_text.count("\n", 0, match.start()) + 1
            links.append((link, line))
    return links


def check() -> int:
    failures: list[str] = []
    checked = 0
    for source in candidate_files():
        rel_source = source.relative_to(ROOT)
        for link, line in iter_links(source):
            if skip_link(link):
                continue
            checked += 1
            if link.startswith(("http://", "https://")):
                target, reason = resolve_github(link)
                if reason is None and target is not None:
                    continue
                if reason is not None:
                    failures.append(f"{rel_source}:{line}: {reason}: {link}")
                    continue
                target, reason = resolve_pages(link)
                if reason is None and target is not None:
                    continue
                if reason is not None:
                    failures.append(f"{rel_source}:{line}: {reason}: {link}")
                continue
            target, reason = resolve_local(source, link)
            if reason is not None:
                failures.append(f"{rel_source}:{line}: {reason}: {link}")
    if failures:
        print("internal link check failed:")
        for failure in failures:
            print(f"  {failure}")
        return 1
    print(f"internal link check passed ({checked} links checked)")
    return 0


sys.exit(check())
PY
