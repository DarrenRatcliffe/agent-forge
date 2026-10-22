#!/usr/bin/env python3
"""Scan a directory tree (or a single file) for likely hardcoded secrets.

Usage: scan-secrets.py [path]

Matched values are masked in output (first 4 characters + ***).
Exit code is 1 if anything was found, 0 if clean, 2 if the path is invalid.
Files larger than MAX_BYTES are skipped and listed on stderr.
"""
import os
import re
import sys

SKIP_DIRS = {".git", "node_modules", ".venv", "venv", "__pycache__", "dist", "build", ".next"}
MAX_BYTES = 1_000_000

# Specific patterns come first: when two patterns match overlapping text,
# only the first one is reported. A named group "v" marks the value to mask.
PATTERNS = [
    ("AWS access key", re.compile(r"\b(?:AKIA|ASIA)[0-9A-Z]{16}\b")),
    ("AWS secret access key", re.compile(
        r"""(?i)aws[\w.-]{0,30}?secret[\w.-]{0,30}['"]?\s*[:=]\s*['"]?(?P<v>[A-Za-z0-9/+=]{40})(?![A-Za-z0-9/+=])""")),
    ("Private key header", re.compile(r"-----BEGIN (?:[A-Z0-9]+ )*PRIVATE KEY(?: BLOCK)?-----")),
    ("GitHub token", re.compile(r"\bgh[pousr]_[A-Za-z0-9]{30,}\b")),
    ("GitHub fine-grained token", re.compile(r"\bgithub_pat_[A-Za-z0-9_]{22,}\b")),
    ("Slack token", re.compile(r"\bxox[abprs]-[A-Za-z0-9-]{10,}\b")),
    ("Stripe live key", re.compile(r"\b[sr]k_live_[A-Za-z0-9]{20,}\b")),
    ("Google API key", re.compile(r"\bAIza[0-9A-Za-z_\-]{35}(?![0-9A-Za-z_\-])")),
    ("Anthropic API key", re.compile(r"\bsk-ant-[A-Za-z0-9_\-]{20,}")),
    ("OpenAI API key", re.compile(r"\bsk-(?!ant-)(?:proj-|svcacct-|admin-)?[A-Za-z0-9_\-]{20,}")),
    ("JWT", re.compile(r"\beyJ[A-Za-z0-9_-]{10,}\.eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}")),
    # Keyword anywhere in an identifier (MY_API_KEY, DB_PASSWORD, client_secret),
    # optional closing quote after the key (JSON/YAML), quoted or unquoted value.
    ("Generic API key/secret/token", re.compile(
        r"""(?i)[\w.-]*(?:api[_-]?key|access[_-]?key|secret|token|passwd|password|credential)[\w.-]*['"]?"""
        r"""\s*(?<![=!<>])[:=](?![=>])\s*"""
        r"""(?:(?P<qt>['"])(?P<q>[^'"\s]{8,})(?P=qt)|(?P<u>[^\s'"`,;#()\[\]{}<>$]{8,})(?=[\s,;#]|$))""")),
]

PLACEHOLDER_WORDS = ("changeme", "change_me", "placeholder", "your_", "your-", "example",
                     "dummy", "redacted", "xxxx", "****", "....")
DOTTED_IDENTIFIER = re.compile(r"^[A-Za-z_$][\w$]*(?:\.[\w$]+)+$")


def mask(value: str) -> str:
    return value[:4] + "***"


def is_placeholder(value: str, quoted: bool) -> bool:
    low = value.lower()
    if low.startswith(("<", "${", "{{", "%(", "$")) or any(w in low for w in PLACEHOLDER_WORDS):
        return True
    if len(set(value)) <= 2:
        return True
    if not quoted:
        # Unquoted right-hand sides in code are usually expressions, not literals.
        if DOTTED_IDENTIFIER.match(value) or value.isdigit() or re.fullmatch(r"[A-Za-z_]+", value):
            return True
        if low in {"undefined", "required", "optional"}:
            return True
    return False


def extract(m):
    """Return (value, start, end) for a match, or None to ignore it."""
    groups = m.groupdict()
    for name, quoted in (("v", True), ("q", True), ("u", False)):
        if groups.get(name):
            value = groups[name]
            if name in ("q", "u") and is_placeholder(value, quoted):
                return None
            return value, m.start(name), m.end(name)
    return m.group(0), m.start(), m.end()


def decode(raw: bytes):
    """Decode file bytes to text, or return None for binary files."""
    if raw.startswith((b"\xff\xfe", b"\xfe\xff")):
        return raw.decode("utf-16", errors="replace")
    if raw.startswith(b"\xef\xbb\xbf"):
        return raw[3:].decode("utf-8", errors="replace")
    if b"\0" in raw[:2048]:
        return None
    return raw.decode("utf-8", errors="replace")


def scan_file(path: str):
    try:
        with open(path, "rb") as f:
            raw = f.read()
    except OSError:
        return
    text = decode(raw)
    if text is None:
        return
    for lineno, line in enumerate(text.splitlines(), 1):
        spans = []
        for label, rx in PATTERNS:
            for m in rx.finditer(line):
                hit = extract(m)
                if hit is None:
                    continue
                value, start, end = hit
                if any(start < e and s < end for s, e in spans):
                    continue
                spans.append((start, end))
                yield lineno, label, mask(value)


def iter_files(root: str):
    if os.path.isfile(root):
        yield root
        return
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d not in SKIP_DIRS]
        for name in filenames:
            yield os.path.join(dirpath, name)


def main() -> int:
    root = sys.argv[1] if len(sys.argv) > 1 else "."
    if not os.path.exists(root):
        print(f"error: path not found: {root}", file=sys.stderr)
        return 2
    if not (os.path.isdir(root) or os.path.isfile(root)):
        print(f"error: not a file or directory: {root}", file=sys.stderr)
        return 2
    findings = 0
    skipped_large = 0
    for path in iter_files(root):
        try:
            too_large = os.path.getsize(path) > MAX_BYTES
        except OSError:
            continue
        if too_large:
            skipped_large += 1
            print(f"skipped (larger than {MAX_BYTES} bytes): {path}", file=sys.stderr)
            continue
        for lineno, label, masked in scan_file(path):
            findings += 1
            print(f"{path}:{lineno}: {label}: {masked}")
    print(f"{findings} potential secret(s) found")
    if skipped_large:
        print(f"{skipped_large} file(s) skipped as too large to scan; check them manually", file=sys.stderr)
    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
