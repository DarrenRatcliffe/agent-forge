---
name: security-assessment
description: Performs a structured security assessment of a codebase covering secrets exposure, dependency vulnerabilities, injection risks, authentication and authorization gaps, insecure configuration, and cryptography misuse. Use when asked to run a security review, security audit, vulnerability assessment, or check code for security issues.
allowed-tools: Bash(npm audit --json) Bash(npm audit --json --omit=dev) Bash(pnpm audit --json) Bash(yarn audit --json) Bash(yarn npm audit --json) Bash(pip-audit) Bash(pip-audit -r requirements.txt) Bash(pip list *) Bash(govulncheck ./...) Bash(cargo audit) Bash(git log *) Bash(git grep *) Bash(python3 "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" *) Bash(py -3 "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" *) Bash(python "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" *)
---

# Security Assessment

## Scope

This is a cross-cutting security sweep, not a code-quality review. It complements the domain reviewers (frontend-reviewer, node-reviewer, python-reviewer, api-reviewer, db-reviewer), which check style-level compliance against the rules files. This skill covers what those don't: secrets exposure, dependency vulnerabilities, injection classes, authn/authz weaknesses, insecure configuration, crypto misuse, and sensitive data exposure.

A condensed checklist of the same categories is in `${CLAUDE_SKILL_DIR}/checklist.md`. Read it before reporting, to confirm nothing was skipped.

## Methodology

Work through every category below. Skip a category only if it plainly doesn't apply (no dependency manifest, no auth system present), and say so explicitly rather than silently omitting it.

### 1. Secrets and credential exposure
- Run the bundled scanner against the working tree: `python3 "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" .` If `python3` isn't available (common on Windows, where it may be a Microsoft Store stub), use `py -3 "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" .` or `python "${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py" .` instead. Exit code 1 means findings, 0 clean, 2 a bad path; files too large to scan are listed on stderr, so check those by hand.
- It flags likely API keys, AWS credentials, private key headers, tokens, and hardcoded password/secret assignments. It masks matched values in its output; never print a full secret value yourself either.
- Check git history separately, since a secret removed from the working tree is still compromised. List only the commits and files that touched a likely secret, so no secret value enters the output: `git log --all -i -E -G"(api[_-]?key|secret|passw(or)?d|token)['\"]?[[:space:]]*[:=][[:space:]]*['\"]?[^'\"[:space:]]{8,}" --format=%h --name-only`. Report the commit and file; if you need to confirm a hit, look only at the matching line and mask any value you quote (first 4 characters + `***`).
- Anything found in history is Critical regardless of current working-tree status, and needs rotation, not just deletion.

### 2. Dependency vulnerabilities
- Node: `npm audit --json` if `package.json` exists (`npm audit --json --omit=dev` for production dependencies only). With a `pnpm-lock.yaml`, use `pnpm audit --json`; with a `yarn.lock`, `yarn audit --json` (Yarn 1) or `yarn npm audit --json` (Yarn 2+). Never run `audit fix`; this skill only reports.
- Python: `pip-audit` (or `pip-audit -r requirements.txt`) if `requirements.txt` or `pyproject.toml` exists. Never pass `--fix`. If `pip-audit` isn't installed, fall back to `pip list --outdated` and note the gap rather than skipping silently.
- Go: `govulncheck ./...` if `go.mod` exists.
- Rust: `cargo audit` if `Cargo.lock` exists and `cargo-audit` is installed; otherwise note the gap.
- Separate vulnerabilities with an available fix from ones without. The remediation path differs.

### 3. Injection risks
- Grep for string-built SQL, string-built shell commands, and unsanitised input reaching `eval`, `exec`, `subprocess` with `shell=True`, or template rendering with `| safe`, `dangerouslySetInnerHTML`, or direct `innerHTML` assignment.
- Check every user-input entry point specifically, not just files that happen to be in a recent diff.
- Path traversal: user input reaching file paths (`open`, `fs.readFile`, `send_file`, archive extraction) without normalising and confirming the result stays inside the intended directory.
- File uploads: type checked by content, not just extension or client-sent MIME type; size limited; stored outside the web root under a generated name.
- SSRF: user-supplied URLs fetched server-side (webhooks, URL previews, image fetchers) without an allowlist, or without blocking internal and metadata addresses (`127.0.0.1`, `169.254.169.254`, private ranges) after DNS resolution and redirects.
- Unsafe deserialization: untrusted data passed to `pickle`, `yaml.load` without `SafeLoader`, `marshal`, Java/.NET native deserialization, or `node-serialize`.

### 4. Authentication and authorization
- Confirm every route that should require auth actually checks it in code, not just that an auth mechanism exists somewhere in the project.
- Check authorization is verified separately from authentication. "Logged in" is not "allowed to touch this specific resource." A route checking only the former is a common, serious gap (IDOR).
- Check session and token handling: expiry set, tokens invalidated on logout, no sensitive data sitting in a JWT payload.
- Passwords are hashed with a slow, salted algorithm (argon2id, bcrypt, or scrypt), never stored plain, encrypted reversibly, or hashed with a fast hash like SHA-256.

### 5. Insecure configuration
- CORS: wildcarded origins on anything handling real user data.
- Cookies: `Secure`, `HttpOnly`, and `SameSite` flags on session cookies.
- Security headers: baseline set present on public-facing services (CSP, X-Content-Type-Options, HSTS).
- Debug or verbose error output left enabled in what looks like a production configuration.

### 6. Cryptography misuse
- Deprecated or weak algorithms (MD5, SHA1, DES) used for anything security-relevant, not general checksums.
- Hardcoded encryption keys or IVs, or home-rolled crypto instead of a maintained library.

### 7. Sensitive data exposure
- Secrets, tokens, or PII written to logs.
- Stack traces or internal error detail leaking into API responses.

## Reporting

Group findings by severity, since that's what determines what gets fixed first:
- **Critical** — actively exploitable, or a live/historical secret exposure. Blocks shipping.
- **High** — a real vulnerability with a plausible exploit path, not yet confirmed exploitable in this specific deployment.
- **Medium** — increases risk but needs another condition to be exploitable.
- **Low / Informational** — hardening opportunity, not active risk.

For every finding: what it is, where (file and line where applicable), why it matters concretely, not "best practice" but the actual consequence, and what fixing it looks like. Don't pad the report with informational noise to look thorough. A short accurate report beats a long padded one.

If a category is clean, say so in one line. Silence is not the same as "checked and clean."

This skill reports. It does not fix anything. Security findings, especially auth and crypto ones, need deliberate human review before a fix is applied, not an autonomous edit.
