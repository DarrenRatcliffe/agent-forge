---
paths:
  - "**/migrations/**"
  - "**/alembic/**"
  - "**/*.sql"
---

# Database Rules — PostgreSQL / SQL

## Schema design
- Every table gets a primary key. Prefer a surrogate key (UUID, or `bigint GENERATED ALWAYS AS IDENTITY` rather than `bigserial`) over natural keys unless there's a specific reason.
- Foreign keys are enforced at the database level, not just assumed in application code. If referential integrity matters, the constraint exists.
- `NOT NULL` by default. A nullable column needs a reason, not the other way round.
- Use the narrowest correct type. Don't default everything to `text` or `varchar(255)` out of habit — dates as `timestamptz`, money as `numeric`, not `float`.
- Naming: snake_case for tables and columns, plural table names, singular column names. Consistent across the whole schema, not per-migration taste.

## Migrations
- Every schema change goes through a migration file, checked into version control. No manual `ALTER TABLE` against a live database, ever.
- Migrations are forward-only and reviewed like code: fix a bad migration with a new forward migration, not by running a down-migration in production. Each non-trivial migration still has a written rollback plan (a down-migration or documented steps), but don't rely on it in place of testing the up-migration properly first.
- Additive changes (new nullable column, new table) can ship without downtime coordination. Anything that changes or removes an existing column needs a deploy plan: add new, backfill, migrate reads, remove old. Don't do it in one step against a live system.
- Large table migrations (adding a NOT NULL column, adding an index) get checked for lock behaviour before running against production. Some of these lock the table; know which ones before you run them. Build indexes on large live tables with `CREATE INDEX CONCURRENTLY`, which can't run inside a transaction (so put it in its own non-transactional migration).

## Queries
- Parameterised queries only. No string-built SQL with interpolated user input, no exceptions.
- No `SELECT *` in application code. Name the columns you need.
- Every query that filters or joins on a column needs to know whether that column is indexed. If it's a hot path and it's not indexed, that's a finding, not something to discover in production.
- Avoid N+1 query patterns. If you're looping and querying inside the loop, batch it or join instead.

## Indexes
- Index foreign key columns. Postgres doesn't do this automatically and unindexed FKs are a common slow-join cause.
- Don't index speculatively. Every index has a write cost — add it because a real query needs it, not because it seems safe.
- Composite index column order follows the query pattern: equality-filtered columns first, then range or sort columns. Not by selectivity, alphabetically, or gut feel.

## Transactions
- Anything that must succeed or fail as a unit goes in a transaction. Don't rely on application-level "if this fails, manually undo that" logic for multi-step writes.
- Keep transactions short. No external API calls or slow operations held inside an open transaction.
- Know your isolation level and why. Default is usually fine; if you've changed it, document why.

## Managed Postgres (e.g. Azure Database for PostgreSQL)
- Connection pooling configured explicitly (PgBouncer or the platform equivalent) — don't let application instances open unbounded raw connections against the managed server.
- Backups and point-in-time restore are the provider's job to run, but it's your job to know the retention window and to have actually tested a restore, not just assumed it works.

## Testing
- Migrations tested against a copy of production-shaped data where the change is non-trivial, not just against an empty local dev database.
- Integration tests that hit a real (test) Postgres instance for anything with non-trivial query logic. Don't validate SQL correctness with a mocked ORM.