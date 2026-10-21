---
name: performance-assessment
description: Analyzes performance characteristics, identifies bottlenecks, and recommends optimization strategies across CPU, memory, I/O, and latency domains. Use when asked to profile code, identify performance issues, benchmark changes, or assess performance readiness.
---

# Performance Assessment

## Scope

This skill provides structured guidance for analyzing application performance across multiple dimensions: CPU profiling, memory profiling, I/O efficiency, latency analysis, and bottleneck identification. It complements language-specific development agents and identifies optimization opportunities that are difficult to spot through code review alone.

Use this skill when:
- An endpoint or operation is slow (high latency)
- Memory usage is higher than expected
- CPU utilization is inefficient
- I/O (database, files, network) is a bottleneck
- Comparing performance of two implementations
- Assessing performance before production deployment
- Investigating user-reported slowness

## Methodology

### 1. Profiling baseline and context

**When to establish baseline:**
- Before any optimization attempt
- Before comparing implementations
- When performance regression is suspected

**How to establish baseline:**
- Measure the operation in its typical context (not synthetic)
- Record: time, CPU usage, memory usage, I/O operations (disk reads/writes, network calls)
- Use realistic data volume (not toy 10-row dataset if production has millions)
- Run enough iterations to see the variance: a few runs can show a large difference, but percentiles need many samples (30+ at minimum, ideally hundreds, or a load-testing tool such as k6, wrk, or Locust)
- Calculate mean, median, and 95th percentile if timing is critical; don't report a p95 from a handful of runs
- Document environment (CPU cores, memory available, OS, other processes)

**Tools by language:**
- **Python:** `cProfile`, `py-spy`, `memory_profiler`, `timeit`
- **Node.js:** `clinic.js`, `node --prof`, `0x`, built-in `perf_hooks`
- **Go:** `pprof`, `go tool trace`, `runtime/debug`
- **Rust:** `flamegraph`, `perf`, `criterion` (benchmarking)

### 2. CPU profiling and hotspot analysis

**When to profile CPU:**
- Operation is slow; unclear why
- High CPU usage observed in production
- Comparing two implementations for efficiency

**How to profile:**
1. Run operation with CPU profiler enabled
2. Capture call stack samples (typically 100+ samples/sec)
3. Generate flame graph or tree view showing % time per function
4. Identify hotspots: functions consuming >5% CPU time
5. Check if hotspots are expected (e.g., expected hash function in crypto, OK; parsing 100MB JSON repeatedly, not OK)

**What to look for:**
- Functions called more often than expected (N+1 pattern in-process)
- Unnecessary allocations or garbage collection pauses
- Recursive functions with deep stacks
- Inefficient string/list operations in tight loops
- Sorting/searching with suboptimal algorithms
- Lock contention if multithreaded

**Actionable fixes:**
- Cache results instead of recomputing (memoization)
- Batch operations (process 1000 items at once, not 1 at a time)
- Use appropriate data structures (hash set for lookups, not list)
- Reduce allocations (reuse buffers, avoid intermediate copies)
- Use SIMD or vectorized operations if available
- Profile after each optimization to confirm improvement

### 3. Memory profiling and leak detection

**When to profile memory:**
- Memory usage grows over time (leak suspected)
- Memory usage higher than expected
- Running in memory-constrained environment (container, embedded)

**How to profile:**
1. Run with memory profiler (tracks allocations)
2. Take snapshot of heap at key points (startup, after 1 hour, before OOM)
3. Compare snapshots to identify growth
4. Identify objects consuming most memory (by count, by total size)
5. Trace references back to find what's holding them

**What to look for:**
- Unbounded collections (lists/maps) that grow indefinitely
- Circular references or closures preventing garbage collection
- Caches without eviction policy
- Long-lived temporary objects (should be garbage collected earlier)
- Memory fragmentation (many small allocations)

**Actionable fixes:**
- Add bounds or eviction to caches (LRU, TTL)
- Break circular references (set to null, weak references)
- Free resources explicitly (close files/sockets, disconnect from DB)
- Reuse objects where possible (object pools for hot-path allocations)
- Reduce retained data size (store IDs instead of full objects)
- Profile again to confirm memory is released

### 4. I/O efficiency and latency

**When to assess I/O:**
- Operation is slow; CPU/memory usage is low (I/O bottleneck)
- Database queries are frequent
- File or network I/O is blocking

**How to assess:**
1. Instrument I/O calls (database, files, network) to measure latency
2. Count I/O operations: how many calls for a single user operation?
3. Measure total I/O time vs total operation time (% blocked)
4. Identify latency outliers (fast calls, slow calls)

**What to look for:**
- N+1 queries (loop with one query per item)
- Sequential I/O instead of parallel (download 1 file, then 1 file, then... vs parallel)
- Unnecessary round-trips (fetch data, inspect, fetch more related data, repeat)
- Synchronous waits that could be async
- Inefficient query (scanning table instead of using index)
- Connection overhead (new connection per request vs pooling)

**Actionable fixes:**
- Batch I/O (fetch 100 items in one query, not 100 queries)
- Use indexes or caches to reduce I/O operations
- Parallel I/O (download 10 files concurrently)
- Async I/O (don't block while waiting)
- Connection pooling (reuse connections, not new ones)
- Pagination (don't fetch 1M rows at once)

### 5. End-to-end latency analysis

**When to analyze latency:**
- User-facing operation is slow
- SLA target not being met (e.g., 95th percentile < 200ms)
- Comparing performance across deployment environments

**How to analyze:**
1. Instrument entire operation path (entry point to response)
2. Break down latency by component:
   - Request parsing
   - Auth/validation
   - Business logic (further broken down to sub-operations)
   - Database
   - Serialization/response sending
3. Measure each component multiple times (get distribution, not just average)
4. Identify which component consumes most time (often 80/20 rule)

**Percentile targets to measure:**
- p50 (median): typical user experience
- p95: slow users (if SLA, target this)
- p99: very slow users
- p999: outliers, but worth investigating if frequently hitting timeout limits

**Actionable fixes:**
- Optimize slowest component (focus on high-impact changes)
- Cache results (trade memory for speed)
- Parallelize independent operations
- Reduce data transfer (compression, fewer fields)
- Move slow operations off critical path (background jobs)

### 6. Throughput and scalability

**When to assess throughput:**
- Need to support more concurrent users
- Comparing designs for scalability
- Planning infrastructure

**How to assess:**
1. Load test with increasing concurrency (10, 50, 100, 500 users)
2. Measure: throughput (requests/sec), latency distribution (p50/p95/p99), error rate
3. Identify breaking point: where does latency degrade or errors increase?
4. Identify bottleneck: CPU exhausted? Memory? Lock contention? I/O?

**Scaling approaches:**
- Vertical: use faster hardware (limited, expensive)
- Horizontal: add more instances, load balance (limited by stateless requirement)
- Optimize: reduce resource per request (increase throughput per instance)
- Cache: reduce backend load (shared cache reduces I/O)
- Async: move work off critical path

### 7. Benchmarking and regression detection

**When to benchmark:**
- Before and after optimization
- Comparing two implementations
- Setting performance baselines

**How to benchmark:**
1. Use language-specific benchmarking tool:
   - Python: `timeit`, `pytest-benchmark`
   - Node.js: `tinybench` (`benchmark.js` is unmaintained), built-in `perf_hooks`
   - Go: `testing.B`, with `benchstat` to compare runs
   - Rust: `criterion`
2. Run same operation 100-1000x to get distribution
3. Measure: mean, median, std deviation, min, max
4. Compare: (new - old) / old = % change
5. Detect regressions: if p95 increased by >5%, investigate (repeat the run first to rule out noise)

**Regression prevention:**
- Add performance benchmarks to CI
- Alert when benchmark degrades >5%
- Store baseline results, trend over time

---

## Checklist

Before finishing a performance assessment:

- [ ] Baseline established (current performance measured)
- [ ] Profiling completed if CPU/memory suspected (with flame graph or call tree)
- [ ] I/O operations counted and latency measured
- [ ] End-to-end latency broken down by component
- [ ] Bottleneck identified and quantified (this component takes X% of time)
- [ ] Actionable recommendations provided (specific optimization, expected impact)
- [ ] Tradeoffs explained (speed vs memory, complexity, maintainability)
- [ ] Before/after comparison possible (can re-measure after optimization)
- [ ] Measurement approach documented (so others can reproduce)
- [ ] All measurements use realistic data volume and concurrency level

