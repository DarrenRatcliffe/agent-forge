---
name: database-optimization
description: Analyzes database performance, query execution plans, missing indexes, and migration safety concerns. Use when asked to optimize queries, review schema changes for performance impact, or assess database performance issues.
allowed-tools: Bash(psql *) Bash(grep *) Bash(cat *) Bash(find *)
---

# Database Optimization

## Scope

This skill provides structured guidance for analyzing database performance and optimizing PostgreSQL queries, indexes, and schema changes. It complements `database-dev` (which writes migrations) and `db-reviewer` (which checks rule compliance). This skill covers performance characteristics that the rules file doesn't detail: execution plans, index strategies, lock timing, and N+1 queries.

Use this skill when:
- A query is slow or missing an index
- A migration may cause long locks or performance issues
- Bulk operations need optimization strategy
- Caching or denormalization is worth considering
- Connection pooling or batch processing is needed

Read the checklist before analysis to confirm all relevant categories are covered.

## Methodology

### 1. Query execution analysis

**When to analyze:**
- Before optimizing a slow query
- After adding WHERE/JOIN/GROUP BY clauses
- When a query is part of a hot path (runs frequently)

**How to analyze:**
- Get the query and a representative database state
- Run `EXPLAIN ANALYZE` on the query to see actual execution plan
- Look for:
  - Sequential scans on large tables (should use index if possible)
  - Join order and strategy (nested loop vs hash vs merge)
  - Filter selectivity (% of rows filtered at each step)
  - Actual rows vs estimated rows (planner accuracy)
  - Node timing (which step takes most time)

**What to report:**
- Current plan and actual execution time
- Cost estimate vs actual time
- Bottleneck: which step is slow (scan, join, aggregate)?
- Potential fix (missing index, rewritten query, schema change)

**Limit queries to safe operations:**
- Read-only SELECT queries only
- Against a test/dev database or a read replica, never production
- Add `LIMIT 10` to EXPLAIN queries to cap execution
- Never run `VACUUM`, `ANALYZE`, or write operations

### 2. Index strategy

**When to assess:**
- After schema design review
- For any query with sequential scans on large tables
- When adding a new lookup field
- For hot-path queries that run frequently

**How to assess:**
- List existing indexes: `SELECT * FROM pg_indexes WHERE tablename = '...'`
- Identify missing indexes by examining EXPLAIN plans
- Check index selectivity: are indexes used for common queries?
- Evaluate cost-benefit: index speeds up SELECT but slows INSERT/UPDATE/DELETE

**What to recommend:**
- Which fields or combinations need indexes (single-column, compound)
- Index type: B-tree (default), hash, GiST, GIN (for arrays/JSON)
- Partial indexes for conditional lookups (e.g., only active records)
- Covering indexes to avoid table lookups (INCLUDE clause in newer PG versions)
- Which existing indexes might be unused or redundant

**Don't recommend:**
- Indexing everything (cost of maintenance outweighs benefit)
- Indexes without understanding the query patterns
- Indexes that duplicate existing ones

### 3. N+1 queries and batch operations

**When to identify:**
- When a query loop is suspected (one query plus N related queries)
- During code review of data-fetching patterns
- If slow page loads correlate with many DB connections

**How to identify:**
- Trace query logs or application profiling data
- Look for: one SELECT + loop { SELECT per item }
- Check if this could be one JOIN or batch SELECT

**Solutions to propose:**
- JOIN instead of loop
- Batch query: `WHERE id IN (list)` instead of loop
- Subqueries or CTEs
- Eager loading in ORMs

### 4. Transaction scope

**When to assess:**
- For long-running operations (bulk updates, migrations)
- For multi-step operations that must all succeed or all fail
- When reviewing lock contention

**What to check:**
- Transaction should be as short as possible (lock held only while needed)
- Avoid holding locks across I/O (network, file system)
- Use `SERIALIZABLE`, `REPEATABLE READ`, or `READ COMMITTED` appropriately
- If reading data then writing based on result, use explicit locking or application-level concurrency control

**What to avoid:**
- Transactions that span minutes (block other users)
- Interactive transactions (waiting for user input inside a transaction)
- Missing error handling (rollback on failure)

### 5. Migration safety and lock risk

**When to assess:**
- Before running any schema-altering migration
- For operations on large tables (millions of rows)
- For operations that add constraints or remove columns

**How to assess:**
1. Identify the table size: `SELECT pg_size_pretty(pg_total_relation_size('table_name'))`
2. Identify the operation: `ALTER TABLE`, `CREATE INDEX`, `ADD COLUMN`, `DROP COLUMN`, `ADD CONSTRAINT`
3. Check PostgreSQL documentation for lock types:
   - Some operations (e.g., `ADD COLUMN` with no default, `ADD INDEX CONCURRENTLY`) don't lock writes
   - Others (`DROP COLUMN`, `ALTER COLUMN TYPE`, `ADD CONSTRAINT CHECK`) lock the table during the operation
4. Estimate duration: smaller tables (< 1GB) usually complete in seconds; larger tables may need longer maintenance windows

**Lock risk levels:**
- **Green (safe):** ADD COLUMN (no default), ADD INDEX CONCURRENTLY, DROP INDEX CONCURRENTLY
- **Yellow (caution):** ADD COLUMN (with default on new rows), UPDATE column in background, ADD CONSTRAINT UNIQUE
- **Red (risky):** DROP COLUMN, ALTER COLUMN TYPE, ADD CONSTRAINT CHECK on unverified column, REINDEX TABLE
- **Black (plan carefully):** Foreign key ADD CONSTRAINT between large tables, CLUSTER, full table REBUILD

**Mitigation strategies:**
- Run long operations during maintenance window (low traffic)
- Use `CONCURRENTLY` modifier where available (CREATE INDEX CONCURRENTLY, DROP INDEX CONCURRENTLY)
- For large tables, break into batches (update 10k rows, commit, repeat)
- Use `pg_stat_activity` to monitor locks and waiting transactions

### 6. Caching and denormalization trade-offs

**When to consider:**
- If a query is slow despite optimization (index, execution plan)
- If the same query runs frequently but data changes infrequently
- If complex calculations (JOINs, aggregates) are repeated

**Caching levels (outside DB):**
- Application-level cache (in-memory, short TTL): fast, but stale, limited to one instance
- Distributed cache (Redis, Memcached): fast, shared across instances, still stale
- Database query cache: built-in (query results cached), but limited control
- Pre-computed materialized view: results pre-calculated, refreshed on schedule or trigger

**Denormalization (inside DB):**
- Store pre-calculated values in the table (e.g., `order_total`, `user_post_count`)
- Trade-off: read speed vs write complexity (must keep denormalized fields in sync)
- Use triggers or application logic to keep denormalized fields updated
- Document denormalized fields clearly

**When to use:**
- Read-heavy, write-light workloads
- When the query is complex (many JOINs, aggregates)
- When staleness is acceptable (cache TTL, or eventual consistency)
- When the cost of computation outweighs the cost of storage

### 7. Connection pooling and batch processing

**When to assess:**
- If application creates many new connections per request
- If the workload is bursty (many requests at once)
- If prepared statements could reduce parsing overhead

**Recommendations:**
- Use a connection pool (PgBouncer, pgpool, or built-in ORM pool)
- Set pool size: typically 5-20 connections, depending on query complexity and concurrency
- Use prepared statements to avoid repeated parsing
- Batch multiple operations in a transaction (reduce round-trips)
- For bulk operations, use `COPY` (for inserts) or `UPDATE ... FROM` (for bulk updates)

---

## Checklist

Before finishing a database optimization analysis, confirm:

- [ ] Query execution plan analyzed (EXPLAIN ANALYZE output shown)
- [ ] Indexes reviewed: missing indexes identified, existing indexes used by planner
- [ ] N+1 queries or batch operations checked
- [ ] Transaction scope evaluated (locks, duration)
- [ ] Migration safety assessed (lock risk, duration estimate)
- [ ] Caching or denormalization considered (if applicable)
- [ ] Connection pooling strategy reviewed (if applicable)
- [ ] All recommendations include justification and trade-offs
- [ ] Recommendations are safe to implement (no destructive operations suggested)

