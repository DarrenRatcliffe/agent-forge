---
paths:
  - "**/{server,api,routes,controllers,services,middleware}/**/*.{js,mjs,cjs,ts}"
  - "**/package.json"
---

# Backend Rules — Node.js

## Stack assumptions
- TypeScript, not plain JS, for anything beyond a throwaway script.
- Async/await throughout. No raw `.then()` chains, no callback-style APIs unless wrapping a legacy library.

## Structure
- Layered separation: routes/controllers → services → data access. No business logic sitting directly in route handlers.
- One responsibility per module. If a file is doing routing, validation, and DB queries, split it.
- Config and secrets come from environment variables via a single validated config module, not scattered `process.env` reads across the codebase.

## Error handling
- No silently swallowed errors. Every catch block either handles the error meaningfully or rethrows.
- Centralised error-handling middleware for HTTP errors. Don't hand-roll error responses in every route.
- Distinguish operational errors (bad input, not found, unauthorized) from programmer errors (bugs). They get handled differently — operational errors are expected and return clean responses; programmer errors should fail loudly in logs.
- Never return a raw stack trace or internal error message to the client in production.

## Async correctness
- Every promise is awaited or explicitly caught. No fire-and-forget async calls that can fail silently.
- Use `Promise.all` for independent concurrent operations. Don't serialise awaits that don't depend on each other.

## Logging
- Structured logging (JSON), not `console.log` scattered through the codebase.
- Log at the point of failure with enough context to debug without reproducing (request id, relevant ids, not full payloads with sensitive data).
- No logging of secrets, tokens, or PII, full stop.

## Dependencies
- Justify any new dependency before adding it — what does it replace, what's the maintenance cost. Don't pull in a library for something that's twenty lines of code.
- Lockfile committed. No floating versions in production dependencies.

## Testing
- Unit tests for business logic in the service layer.
- Integration tests for anything touching the database or external APIs, run against a real or realistic test instance, not mocked into meaninglessness.
- Tests must be able to run in CI without manual setup steps.

## Security baseline
- Input validation at the boundary (request schema validation), not scattered checks deep in the call stack.
- No string-concatenated queries. Parameterised queries or an ORM/query builder only.
- Rate limiting and basic security headers (helmet or equivalent) on any public-facing service.