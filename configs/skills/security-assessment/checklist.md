# Security Review Checklist

Use this to confirm every category in `SKILL.md` was covered. Mark an item N/A with a reason; don't leave it blank.

## 1. Secrets and credentials
- [ ] `scan-secrets.py` run against the working tree
- [ ] Git history checked for committed secrets
- [ ] API keys and tokens come from environment variables or a secret store

## 2. Dependencies
- [ ] `npm audit` run if `package.json` exists
- [ ] `pip-audit` run if a Python manifest exists (or the gap is noted)
- [ ] Findings split into fix-available and no-fix-available

## 3. Injection
- [ ] No string-built SQL; parameterised queries only
- [ ] No user input reaching `eval`, `exec`, or `subprocess` with `shell=True`
- [ ] No unescaped HTML sinks (`dangerouslySetInnerHTML`, `innerHTML`, `| safe`)
- [ ] File uploads validated by type and size
- [ ] Path traversal prevented on file operations

## 4. Authentication and authorization
- [ ] Every protected route checks authentication in code
- [ ] Authorization checked separately, per resource (no IDOR)
- [ ] Tokens and sessions expire and are invalidated on logout
- [ ] Passwords hashed with bcrypt or argon2
- [ ] No sensitive data in JWT payloads

## 5. Configuration
- [ ] CORS limited to known origins for user-data endpoints
- [ ] Cookies set `Secure`, `HttpOnly`, `SameSite`
- [ ] Security headers present (CSP, X-Content-Type-Options, HSTS)
- [ ] Debug mode and verbose errors off in production config

## 6. Cryptography
- [ ] No MD5, SHA1, or DES for security purposes
- [ ] No hardcoded keys or IVs; no home-rolled crypto

## 7. Sensitive data exposure
- [ ] No secrets, tokens, or PII in logs
- [ ] No stack traces or internal errors in API responses
