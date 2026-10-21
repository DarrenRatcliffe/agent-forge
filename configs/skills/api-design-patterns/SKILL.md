---
name: api-design-patterns
description: Provides design patterns, trade-offs, and best practices for building scalable, maintainable REST and GraphQL APIs. Use when designing API contracts, choosing pagination strategy, versioning approach, or evaluating API design trade-offs.
---

# API Design Patterns

## Scope

This skill provides structured guidance for common API design decisions and patterns. It complements `api-dev` (which implements endpoints) and `api-reviewer` (which checks rule compliance). This skill covers design patterns that the rules file doesn't prescribe: pagination strategies, versioning approaches, error response envelopes, filtering and sorting, webhooks, and authentication patterns.

Use this skill when:
- Designing a new API from scratch
- Choosing pagination for a large dataset endpoint
- Deciding on versioning strategy
- Designing error response format
- Building webhook system
- Supporting filtering and sorting
- Choosing REST vs GraphQL
- Designing rate limiting and quotas

## Methodology

### 1. REST vs GraphQL decision

**Choose REST if:**
- Simple CRUD operations on resources
- Clear, hierarchical resource structure
- Cacheable responses (HTTP caches, CDNs work well)
- Client diversity low (1-2 client types)
- API complexity low to medium
- Team familiar with REST patterns
- Performance debugging simpler (each HTTP request is atomic)

**Choose GraphQL if:**
- Complex queries with deep nesting (mobile app needs specific fields)
- Multiple client types with different data needs (web, mobile, embedded)
- Rapid iteration on client requirements (flexible query language)
- Team comfortable with GraphQL ecosystem
- Real-time requirements (subscriptions fit naturally)
- Internal API (trust clients don't abuse query complexity)
- Trade-off: GraphQL queries harder to cache, more complex to secure

**Hybrid approach:**
- REST for simple, high-volume endpoints (public data, search)
- GraphQL for complex, low-volume queries (mobile app data needs)

### 2. Pagination patterns

**Offset/Limit (most common)**
- Usage: `GET /items?offset=0&limit=20`
- Pros: Simple, stateless, supports direct jumps to page N
- Cons: Inefficient for large offsets (database still scans N rows before returning), inconsistent results if data changes between requests
- Use when: Small datasets, user-controlled jumps between pages, consistency not critical
- Example response:
```json
{
  "items": [...],
  "total": 5000,
  "offset": 0,
  "limit": 20,
  "hasMore": true
}
```

**Cursor-based (recommended for high-volume APIs)**
- Usage: `GET /items?cursor=abc123&limit=20`
- Pros: Efficient (index seeks to cursor), consistent results even if data changes, works well at scale
- Cons: Can't jump to arbitrary page, client must follow links forward/backward
- Use when: High-volume API, data changes frequently, want consistent results
- Example response:
```json
{
  "items": [...],
  "nextCursor": "xyz789",
  "prevCursor": "abc000",
  "hasMore": true
}
```

**Keyset pagination (for very large datasets)**
- Usage: `GET /items?after_id=1000&limit=20` (where after_id is the last item's ID)
- Pros: Most efficient, works at any scale
- Cons: Must order by unique key, can't scroll backward easily
- Use when: Massive datasets (millions of rows), forward-only scrolling acceptable

**Time-based (for activity feeds)**
- Usage: `GET /feed?since=2026-01-01T10:00:00Z&limit=20`
- Pros: Intuitive for time-series data
- Cons: Requires timestamp, must handle duplicates if multiple events at same time
- Use when: Event streams, activity feeds, time-series data

### 3. Versioning strategies

**URL path versioning (most explicit)**
- Usage: `/v1/users`, `/v2/users`
- Pros: Clear in URL, easy to route
- Cons: URL bloat, can't gradually migrate clients
- Use when: Major breaking changes expected, different code paths acceptable

**Header-based versioning (most RESTful)**
- Usage: `Accept: application/vnd.myapi.v2+json`
- Pros: URL stays clean, can use content negotiation
- Cons: Less visible, debugging harder (must check headers)
- Use when: Prefer URL stability, client sophistication high

**Version from the start, then evolve additively (best practice)**
- Pick URL or header versioning for the first release, even if only `v1` ever ships (the team rule in `api-dev` requires this; don't bolt versioning on after the first breaking change)
- Within a version, never break existing fields/endpoints; only add new ones (new optional fields, new endpoints don't need a version bump)
- Use optional fields with default values for new behavior
- Deprecate field through documentation before removing (6+ month notice)
- Pros: Breaking changes have a home (`v2`) when they come, while most changes need no version management
- Cons: Requires discipline, schema can accumulate cruft

**Deprecation strategy:**
1. Add new endpoint/field with new behavior
2. Keep old endpoint working, document deprecation
3. State removal date (e.g., "Deprecated 2026-01-01, remove 2026-07-01")
4. Track usage and notify active clients
5. Remove only after significant time window

### 4. Error response formats

**Consistent error structure (required):**
Consider the standard Problem Details format (RFC 9457, `Content-Type: application/problem+json`, with `type`, `title`, `status`, `detail`, `instance` plus extension members) instead of a custom envelope. If you use a custom envelope like the one below, use it everywhere:
```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Email is required",
    "details": {
      "fields": {
        "email": ["Email is required"]
      }
    }
  }
}
```

**Fields to include:**
- `code`: machine-readable error code (enum), used for i18n and client logic
- `message`: human-readable message, safe to show to user
- `details`: specific error details (field-level errors, which validation rule failed, etc.)
- `requestId`: correlation ID for logging/debugging
- HTTP status code: 400 (bad input), 401 (auth), 403 (authz), 404 (not found), 409 (conflict), 422 (semantic validation), 429 (rate limited)

**Per-field vs error-list:**
- Per-field: show validation error next to each input (better UX)
- Error-list: show all errors at once (sometimes preferred)
- Use both if possible

### 5. Filtering and sorting

**Query parameter filtering:**
- Usage: `GET /users?status=active&role=admin`
- Pros: Simple, RESTful, cacheable (sort by field allows caching)
- Cons: Can't express complex queries (AND, OR, nested)
- Use for: Simple equality filters

**Advanced filtering options:**
1. Operators: `GET /items?price__gt=10&price__lt=100` (gt, gte, lt, lte, eq, ne)
2. Ranges: `GET /items?createdAt__between=2026-01-01,2026-12-31`
3. Arrays: `GET /items?tags__in=red,blue` (contains any of)
4. Prefix/substring: `GET /items?name__startswith=John`

**Sorting:**
- Usage: `GET /items?sort=-createdAt,name` (- for descending)
- Pros: Simple, supports multi-field sorts
- Cons: Limited in what can be sorted (must be indexed)
- Limit sortable fields (only indexed ones), document which fields support sort

**Full-text search:**
- Separate endpoint: `GET /search?q=user%20profile`
- Different from filtering (uses text index, fuzzy matching, ranking)
- Return ranked results, not just filter

### 6. Rate limiting and quotas

**Rate limiting levels:**
1. **Unauthenticated:** per IP address (e.g., 60 requests/minute)
2. **Authenticated:** per user or API key (e.g., 10,000 requests/day)
3. **Endpoint-specific:** some endpoints more limited (e.g., uploads 10/day)

**Header-based rate limit info:**
```
X-RateLimit-Limit: 10000
X-RateLimit-Remaining: 9876
X-RateLimit-Reset: 1640995200
```

**When the limit is exceeded:** respond `429 Too Many Requests` with the standard `Retry-After` header (seconds, or an HTTP date):
```
HTTP/1.1 429 Too Many Requests
Retry-After: 5
```
The `X-RateLimit-*` headers are a common convention, not a standard. The IETF `RateLimit` and `RateLimit-Policy` headers (draft-ietf-httpapi-ratelimit-headers) are the emerging standard alternative.

**Quota vs rate limit:**
- Rate limit: requests per time window (minute, hour)
- Quota: total requests per billing period (day, month, year)
- Implement both for different scale tiers

**Strategies:**
- Token bucket: smooth traffic, allow bursts
- Sliding window: consistent enforcement
- Leaky bucket: strict rate limiting
- Use Redis or similar for distributed rate limiting

### 7. Authentication and authorization patterns

**Authentication methods:**
- API Keys: simple, stateless, limited security (store securely, rotate regularly)
- OAuth 2.0: industry standard, supports delegation, good for user-facing APIs
- JWT: stateless, can include claims, verify without database lookup
- mTLS: certificate-based, high security, complex setup

**Authorization patterns:**
1. **Role-based access control (RBAC):** user has role (admin, user), role has permissions
2. **Attribute-based access control (ABAC):** permissions based on attributes (user.dept == resource.dept)
3. **Resource ownership:** resource owner can access, others cannot
4. **Scope-based:** OAuth scopes (read, write, admin), fine-grained permissions

**Implement as:**
- Middleware (check auth before reaching handler)
- Separate: never mix auth code into business logic

### 8. Webhook design

**Webhook triggering:**
- Event-driven: when resource created/updated/deleted, send HTTP POST to registered URL
- Retry strategy: exponential backoff (1s, 2s, 4s, 8s...), max 5 retries
- Signature: sent in a header, never inside the body it signs. Compute HMAC-SHA256(secret, timestamp + "." + raw request body) and send the timestamp alongside it
- Replay protection: receiver recomputes the HMAC over the raw bytes received (before JSON parsing), compares in constant time, and rejects requests whose timestamp is more than ~5 minutes old

**Request design:**
```
POST /webhooks/orders HTTP/1.1
Content-Type: application/json
Webhook-Id: evt_01HV3K8Z9X
Webhook-Timestamp: 1767261600
Webhook-Signature: v1=5257a869e7ecebeda32affa62cdca3fa51cad7e77a0e56ff536d0ce8e108d8bd
```
```json
{
  "id": "evt_01HV3K8Z9X",
  "event": "order.created",
  "timestamp": "2026-01-01T10:00:00Z",
  "data": {
    "id": "order_123",
    "amount": 9999,
    "currency": "USD"
  }
}
```
- Money as an integer in minor units (`9999` = $99.99) plus a currency code, never a float
- Retry attempt number can go in a header (e.g. `Webhook-Attempt: 2`); keep the body identical across retries so the event ID stays stable

**Idempotency:**
- Include a unique event ID (`id`) in every payload
- Client should deduplicate (if seen same event ID, skip)
- Webhook system should deduplicate (don't send same event twice)

**Best practices:**
- Async: send webhook without blocking primary operation
- No secrets in payload (client authenticates via signature)
- Delivery status tracking (success, failed, too many retries)
- Dashboard to replay webhooks

### 9. Consistency and caching strategy

**Immediate consistency:**
- Read immediately after write sees new data
- Simple model, but harder to scale
- Use: financial transactions, user account changes

**Eventual consistency:**
- Read immediately after write may see old data
- Better for scale, acceptable delay (seconds to minutes)
- Use: activity feeds, social media likes, non-critical data
- Communicate to client: "your comment will appear shortly"

**Caching layers:**
1. HTTP caching (Cache-Control header, CDN)
2. Application cache (in-memory, Redis)
3. Database cache (query result caching)
4. Client cache (browser, mobile app)

**Cache invalidation:**
- Time-based: TTL (cache-control: max-age=300)
- Event-based: clear cache when data changes
- Versioning: include version in URL/ETag, change version when data changes

---

## Pattern Decision Matrix

| Scenario | Pattern | Rationale |
|----------|---------|-----------|
| First API, unknown requirements | REST + offset pagination | Simple, easy to learn, adjust later |
| Mobile app, complex data needs | GraphQL or REST with cursor pagination | Flexible queries, efficient paging |
| Public API, high volume | REST + cursor pagination + rate limiting | Efficient, scalable, predictable |
| Internal tool, small team | GraphQL | Rapid iteration, developer experience |
| Financial/payment API | REST + versioning + per-field errors | Stability, clarity, compliance |
| Activity feed | REST + time-based pagination | Natural for time-series |
| Search functionality | REST + full-text search endpoint | Separate from filtering |
| Third-party integrations | REST + webhooks + OAuth | Standard, secure, push notifications |

---

## Checklist

Before finalizing API design:

- [ ] REST vs GraphQL decision made and justified
- [ ] Pagination strategy chosen (offset/limit, cursor, keyset, time-based)
- [ ] Versioning approach decided from the start (URL or header), with additive-only changes within a version
- [ ] Error response format standardized (RFC 9457 problem+json, or code, message, details)
- [ ] Filtering and sorting strategy defined (simple, operators, full-text search)
- [ ] Rate limiting and quotas planned (levels, headers, enforcement)
- [ ] Authentication method chosen (API key, OAuth, JWT, mTLS)
- [ ] Authorization model selected (RBAC, ABAC, ownership, scope)
- [ ] Caching strategy documented (HTTP, app-level, invalidation)
- [ ] Consistency model stated (immediate vs eventual)
- [ ] If webhooks needed: payload, retry, signature strategy defined
- [ ] API contracts documented (OpenAPI/Swagger)
- [ ] Backward compatibility plan in place (deprecation path)

