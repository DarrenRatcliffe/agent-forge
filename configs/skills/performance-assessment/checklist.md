# Performance Assessment Checklist

## Pre-Assessment Setup
- [ ] Clear performance problem statement or optimization goal
- [ ] Realistic test environment (not toy data, appropriate concurrency)
- [ ] Baseline measurement in place (current performance known)
- [ ] Multiple runs planned (at least 3-5 iterations for variance)
- [ ] Documentation of environment (CPU, RAM, OS, other load)

## CPU Analysis
- [ ] CPU profiler run and output captured
- [ ] Hotspots identified (functions > 5% CPU time)
- [ ] Call stacks analyzed for unexpected patterns
- [ ] Flame graph or tree view reviewed
- [ ] Algorithm complexity verified (O(n²) vs O(n log n), etc.)
- [ ] Garbage collection pauses noted (if applicable)

## Memory Analysis
- [ ] Memory baseline established
- [ ] Heap snapshots taken at key points
- [ ] Growth rate identified (linear, exponential, stable)
- [ ] Top memory consumers identified (by count, by size)
- [ ] Leak suspects listed (unbounded caches, circular refs, etc.)
- [ ] GC efficiency checked (pause time, frequency)

## I/O Analysis
- [ ] I/O operations instrumented and counted
- [ ] Database queries measured (count, latency, indexes used)
- [ ] File/network I/O quantified
- [ ] N+1 patterns identified (if present)
- [ ] Batch operation opportunities noted
- [ ] Connection pooling effectiveness assessed

## Latency Breakdown
- [ ] Full operation instrumented (entry to exit)
- [ ] Latency attributed to components:
  - [ ] Request parsing/validation
  - [ ] Authentication/authorization
  - [ ] Business logic
  - [ ] Database
  - [ ] External APIs
  - [ ] Serialization/response
- [ ] Percentiles captured (p50, p95, p99)
- [ ] Slowest component identified (focus area)

## Throughput & Scalability
- [ ] Load test completed (10, 50, 100+ concurrent users)
- [ ] Throughput measured (requests/sec at each level)
- [ ] Latency distribution tracked across concurrency levels
- [ ] Breaking point identified (where does degradation occur)
- [ ] Resource bottleneck identified (CPU, memory, I/O, lock contention)

## Recommendations
- [ ] Specific optimization identified (not vague)
- [ ] Expected impact quantified (X% improvement, Y ms latency reduction)
- [ ] Tradeoffs listed (speed vs memory, complexity, maintainability)
- [ ] Implementation complexity assessed
- [ ] Side effects considered (does optimization affect other features?)

## Validation & Benchmarking
- [ ] Benchmark suite available (can re-run after optimization)
- [ ] Regression test planned (prevent degradation)
- [ ] Performance metric tracked (SLA target, p95, throughput)
- [ ] Before/after measurement will be comparable

## Reporting
- [ ] Findings clearly stated with numbers (not vague)
- [ ] Baseline and current performance compared
- [ ] Root cause identified (not just symptom)
- [ ] Recommendations prioritized (high impact first)
- [ ] Measurement approach documented (reproducible)

