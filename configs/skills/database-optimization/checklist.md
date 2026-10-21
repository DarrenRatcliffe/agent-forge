# Database Optimization Checklist

## Pre-analysis
- [ ] Database state is clear (test/dev database, not production, or read replica)
- [ ] Query in question is identified and reproducible
- [ ] Baseline performance metric available (query time, row count affected, frequency)

## Execution Plan Review
- [ ] EXPLAIN ANALYZE output captured
- [ ] Scan types identified (sequential, index, bitmap)
- [ ] Join strategies noted (nested loop, hash, merge)
- [ ] Filter selectivity assessed (% of rows filtered at each step)
- [ ] Actual vs estimated rows compared (planner accuracy)
- [ ] Node timing analyzed (which step is bottleneck)

## Index Analysis
- [ ] Existing indexes catalogued
- [ ] Planner confirmed to use the expected indexes for the query in question
- [ ] Unused indexes identified
- [ ] Missing indexes proposed (with selectivity justification)
- [ ] Index type justified (B-tree, GiST, GIN, hash)
- [ ] Index write-side cost considered

## Query Patterns
- [ ] N+1 query patterns identified
- [ ] Batch operation opportunities noted
- [ ] Join reordering or rewrite considered
- [ ] Subquery vs JOIN trade-off evaluated
- [ ] Pagination or LIMIT strategy for large result sets

## Transaction & Lock Analysis
- [ ] Transaction scope checked (as short as possible)
- [ ] Lock risk identified for schema changes
- [ ] Migration duration estimated (table size, how long locks are held)
- [ ] Maintenance window requirement noted (if long-running)
- [ ] CONCURRENTLY options used where available
- [ ] Batching strategy for bulk operations

## Caching & Denormalization
- [ ] Staleness tolerance assessed
- [ ] Cache layer opportunity identified (if applicable)
- [ ] Denormalization trade-off evaluated (if applicable)
- [ ] Invalidation strategy considered

## Connection Pooling (if applicable)
- [ ] Pool size and pooler (e.g. PgBouncer) reviewed against the connection limit

## Reporting
- [ ] Findings stated clearly with specific file/line or query
- [ ] Trade-offs explained (performance vs complexity vs maintenance)
- [ ] Recommendations are actionable (specific SQL, index definition, etc.)
- [ ] Recommendations are safe to implement (no destructive operations suggested)
- [ ] Severity levels assigned (critical bottleneck, good-to-have optimization, etc.)
- [ ] Next steps clear (implement query change, add index, profile in production, etc.)

