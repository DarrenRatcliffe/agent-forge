---
paths:
  - "src/server/**/*.ts"
  - "src/server/**/*.js"
  - "src/server/**/*.py"
---

# API Development Rules

## Design first
- No endpoint gets built without agreeing the contract first: method, path, request shape, response shape, status codes. If that's not written down anywhere, write it down before coding.
- REST resource naming: plural nouns, no verbs in the path (`/orders`, not `/getOrders`). Actions that don't map to CRUD get a clear sub-resource or action-style endpoint, documented as the exception it is.
- Version the API from the start (`/v1/...` or a header-based scheme — pick one, be consistent). Don't bolt on versioning after the first breaking change forces the issue.

## Request and response contracts
- Every request body validated against a schema before it touches business logic. Reject bad input at the boundary with a clear 400, not three layers deep with a confusing 500.
- Consistent response envelope across the API (or consistently no envelope — pick one, don't mix per-endpoint).
- Consistent error response shape: an error code, a human-readable message, and enough detail for the caller to fix their request. Not a raw exception dump.
- Pagination on any list endpoint that could return more than a page's worth of results, from day one, not retrofitted once it's slow.

## Status codes
- Use them correctly: 400 for bad input, 401 for missing/invalid auth, 403 for authenticated-but-not-permitted, 404 for missing resource, 409 for conflict, 422 for semantically invalid input, 500 only for actual server failure.
- Don't return 200 with an error payload inside it. That's not what 200 means.

## Auth and security
- Authentication and authorization are separate concerns and both get checked. "Logged in" is not the same as "allowed to do this."
- No secrets, tokens, or keys in URLs or query strings — headers or body only.
- Rate limiting on public endpoints. Document the limits.
- CORS configured explicitly for known origins, not wildcarded in anything handling real user data.

## Documentation
- OpenAPI/Swagger spec kept in sync with the actual implementation, generated from code where the framework supports it (FastAPI does this natively — use it).
- Every endpoint documents required auth, expected request shape, possible response codes, and at least one example.

## Backward compatibility
- Additive changes (new optional fields, new endpoints) don't require a version bump. Anything that changes or removes existing behaviour does.
- Deprecate before removing. State the removal date, don't just delete an endpoint a client might still be calling.

## Observability
- Every request loggable with a correlation/request id that can be traced through logs across services.
- Track latency and error rate per endpoint. If you can't answer "which endpoint is slow or failing right now," that's a gap, not a nice-to-have.