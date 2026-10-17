#!/usr/bin/env python3
"""Scan a directory tree for likely hardcoded secrets.

Usage: scan-secrets.py [path]

Matched values are masked in output (first 4 characters + ***).
Exit code is 1 if anything was found, 0 otherwise.
"""
import os
import re
import sys

SKIP_DIRS = {".git", "node_modules", ".venv", "venv", "__pycache__", "dist", "build", ".next"}
MAX_BYTES = 1_000_000

PATTERNS = [
    ("AWS access key", re.compile(r"\b(AKIA|ASIA)[0-9A-Z]{16}\b")),
    ("Private key header", re.compile(r"-----BEGIN (?:RSA |EC |DSA |OPENSSH |PGP )?PRIVATE KEY")),
    ("GitHub token", re.compile(r"\bgh[pousr]_[A-Za-z0-9]{30,}\b")),
    ("Slack token", re.compile(r"\bxox[abprs]-[A-Za-z0-9-]{10,}\b")),
    ("Generic API key/secret/token", re.compile(
        r"""(?i)\b(api[_-]?key|secret|token|passwd|password)\b\s*[:=]\s*['"]([A-Za-z0-9_\-/+=]{16,})['"]""")),
]


def mask(value: str) -> str:
    return value[:4] + "***"


def is_binary(chunk: bytes) -> bool:
    return b"\0" in chunk


def scan_file(path: str):
    try:
        if os.path.getsize(path) > MAX_BYTES:
            return
        with open(path, "rb") as f:
            raw = f.read()
    except OSError:
        return
    if is_binary(raw[:2048]):
        return
    text = raw.decode("utf-8", errors="replace")
    for lineno, line in enumerate(text.splitlines(), 1):
        for label, rx in PATTERNS:
            m = rx.search(line)
            if m:
                value = m.group(m.lastindex) if m.lastindex else m.group(0)
                yield lineno, label, mask(value)


def main() -> int:
    root = sys.argv[1] if len(sys.argv) > 1 else "."
    findings = 0
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d not in SKIP_DIRS]
        for name in filenames:
            path = os.path.join(dirpath, name)
            for lineno, label, masked in scan_file(path):
                findings += 1
                print(f"{path}:{lineno}: {label}: {masked}")
    print(f"{findings} potential secret(s) found")
    return 1 if findings else 0


if __name__ == "__main__":
    sys.exit(main())
