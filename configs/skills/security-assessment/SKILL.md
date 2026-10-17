---
name: security-assessment
description: Performs a structured security assessment of a codebase covering secrets exposure, dependency vulnerabilities, injection risks, authentication and authorization gaps, insecure configuration, and cryptography misuse. Use when asked to run a security review, security audit, vulnerability assessment, or check code for security issues.
allowed-tools: Bash(npm audit *) Bash(pip-audit *) Bash(pip list *) Bash(git log *) Bash(git grep *) Bash(python3 ${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py *)
---

# Security Assessment

## Scope

This is a cross-cutting security sweep, not a code-quality review. It complements the domain reviewers (frontend-reviewer, node-reviewer, python-reviewer, api-reviewer, db-reviewer), which check style-level compliance against the rules files. This skill covers what those don't: secrets exposure, dependency vulnerabilities, injection classes, authn/authz weaknesses, insecure configuration, crypto misuse, and sensitive data exposure.

A condensed checklist of the same categories is in `${CLAUDE_SKILL_DIR}/checklist.md`. Read it before reporting, to confirm nothing was skipped.

## Methodology

Work through every category below. Skip a category only if it plainly doesn't apply (no dependency manifest, no auth system present), and say so explicitly rather than silently omitting it.

### 1. Secrets and credential exposure
- Run the bundled scanner against the working tree: `python3 ${CLAUDE_SKILL_DIR}/scripts/scan-secrets.py .`
- It flags likely API keys, AWS credentials, private key headers, tokens, and hardcoded password/secret assignments. It masks matched values in its output; never print a full secret value yourself either.
- Check git history separately, since a secret removed from the working tree is still compromised: `git log -p | grep -iE "(api[_-]?key|secret|password|token)\s*=\s*['\"][a-zA-Z0-9]{16,}"`
- Anything found in history is Critical regardless of current working-tree status, and needs rotation, not just deletion.

### 2. Dependency vulnerabilities
- Node: `npm audit --json` if `package.json` exists.
- Python: `pip-audit` if `requirements.txt` or `pyproject.toml` exists. If `pip-audit` isn't installed, fall back to `pip list --outdated` and note the gap rather than skipping silently.
- Separate vulnerabilities with an available fix from ones without. The remediation path differs.

### 3. Injection risks
- Grep for string-built SQL, string-built shell commands, and unsanitised input reaching `eval`, `exec`, `subprocess` with `shell=True`, or template rendering with `| safe`, `dangerouslySetInnerHTML`, or direct `innerHTML` assignment.
- Check every user-input entry point specifically, not just files that happen to be in a recent diff.

### 4. Authentication and authorization
- Confirm every route that should require auth actually checks it in code, not just that an auth mechanism exists somewhere in the project.
- Check authorization is verified separately from authentication. "Logged in" is not "allowed to touch this specific resource." A route checking only the former is a common, serious gap (IDOR).
- Check session and token handling: expiry set, tokens invalidated on logout, no sensitive data sitting in a JWT payload.

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
