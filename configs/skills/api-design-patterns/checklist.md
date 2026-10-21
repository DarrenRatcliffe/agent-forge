# API Design Patterns Checklist

## REST vs GraphQL Decision
- [ ] API requirements and complexity assessed
- [ ] Client diversity and needs evaluated
- [ ] Performance requirements considered
- [ ] Team expertise and preferences noted
- [ ] Decision documented with rationale
- [ ] Decision reversal cost understood (if needed later)

## Pagination Strategy
- [ ] Dataset size and growth rate estimated
- [ ] Access patterns understood (jumps to page N, or forward-only scrolling?)
- [ ] Pagination method chosen: offset/limit, cursor, keyset, or time-based
- [ ] Consistency requirements clear (can results change between requests?)
- [ ] Response format defined (total count, hasMore, nextCursor, etc.)
- [ ] Maximum page size set (prevent abuse, e.g., limit 20-100 max)

## Versioning Strategy
- [ ] Versioning approach decided from the start: URL path or header (additive-only changes within a version)
- [ ] Rationale documented (why this choice over others)
- [ ] Deprecation process defined (timeline, communication)
- [ ] Old version support duration stated (e.g., 12 months)
- [ ] Additive-only constraints understood (never remove/change existing fields)

## Error Response Format
- [ ] One format used across the API: RFC 9457 problem+json, or a code, message, details envelope
- [ ] Error code enum defined (machine-readable codes)
- [ ] Message field standardized (human-readable, user-safe)
- [ ] Details field structure designed (per-field errors, etc.)
- [ ] HTTP status codes mapped to error conditions
- [ ] Correlation/request ID included (for debugging)
- [ ] Error response examples created

## Filtering and Sorting
- [ ] Simple filters designed (equality matches)
- [ ] Advanced filter operators decided (if needed)
- [ ] Sortable fields identified and indexed
- [ ] Sort order documented (ascending/descending syntax)
- [ ] Full-text search strategy (if applicable)
- [ ] Query complexity limits set (prevent expensive queries)

## Rate Limiting and Quotas
- [ ] Rate limit levels defined (unauthenticated, per-user, per-endpoint)
- [ ] Limit values set (requests/minute, requests/day)
- [ ] Rate limit headers chosen (X-RateLimit-*, Retry-After)
- [ ] Enforcement mechanism planned (Redis, in-memory, etc.)
- [ ] Quota strategy designed (billing tiers, hard limits)
- [ ] Client notification strategy (headers, error responses)

## Authentication & Authorization
- [ ] Authentication method chosen (API key, OAuth, JWT, mTLS)
- [ ] Authorization model selected (RBAC, ABAC, ownership, scope)
- [ ] Token/credential rotation policy defined
- [ ] Scope/permission granularity determined
- [ ] Audit logging planned (who accessed what, when)
- [ ] Token revocation mechanism (logout, key rotation)

## Caching Strategy
- [ ] Cacheability assessed (can response be cached?)
- [ ] Cache-Control headers designed
- [ ] TTL values set (appropriate for each resource)
- [ ] Cache invalidation strategy (time vs event-based)
- [ ] CDN strategy (if public API)
- [ ] Cache layer architecture (client, app, DB)

## Consistency Model
- [ ] Immediate vs eventual consistency chosen
- [ ] Consistency guarantees documented
- [ ] Client expectations managed (potential delays communicated)
- [ ] Data freshness requirements understood
- [ ] Conflict resolution strategy (if distributed)

## Webhooks (if applicable)
- [ ] Events to be sent identified
- [ ] Payload structure designed
- [ ] Event ID for deduplication included
- [ ] Signature sent in a header, computed over timestamp + raw body; old timestamps rejected (replay protection)
- [ ] Retry strategy designed (backoff, max attempts)
- [ ] Delivery status tracking planned
- [ ] Webhook replay functionality considered

## Backward Compatibility
- [ ] Additive changes documented (safe without version bump)
- [ ] Breaking changes identified (require version bump)
- [ ] Deprecation timeline for old versions
- [ ] Migration guide for clients (if breaking change)
- [ ] Communication plan (when to notify clients)

## Documentation
- [ ] OpenAPI/Swagger spec created
- [ ] All patterns documented with examples
- [ ] Rate limit limits and quotas clear
- [ ] Pagination examples shown for chosen strategy
- [ ] Error codes and meanings documented
- [ ] Authentication/authorization flow explained
- [ ] Caching strategy and headers documented

