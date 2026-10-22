---
name: testing-patterns
description: Provides testing patterns, strategies, and best practices by language and framework. Use when designing test suite, choosing testing approach, or improving test coverage and quality.
---

# Testing Patterns

## Scope

This skill provides structured guidance for testing strategies across different languages and frameworks. It complements `test-writer` (which writes tests) and domain reviewers (which check test compliance). This skill covers when to test what and how, not implementation details.

Use this skill when:
- Designing test strategy for new project
- Choosing testing framework
- Assessing test coverage adequacy
- Improving test quality (reducing flaky tests, better patterns)
- Deciding what to test and what not to test
- Cross-language testing pattern comparison

## Methodology

### 1. Test pyramid strategy

**Levels (from bottom to top):**
1. **Unit tests** (foundation, many)
   - Test single function/class in isolation
   - Mock external dependencies
   - Fast to run (milliseconds)
   - Good coverage possible
   - Catch logic errors early

2. **Integration tests** (middle, some)
   - Test components working together
   - Real database (not mocked), real HTTP client to mock service
   - Slower than unit tests (seconds)
   - Catch interface/data flow issues
   - Require test environment setup

3. **End-to-end tests** (top, few)
   - Test complete workflow from user perspective
   - Real or staging environment
   - Slowest (minutes)
   - Test against actual external services (or realistic mocks)
   - Catch environment-specific issues
   - Flaky (environmental factors, timing)

**Rule:** aim for 70% unit, 20% integration, 10% E2E

**In CI/CD:**
- Unit tests: always run (fast, reliable)
- Integration tests: always run (catch regressions)
- E2E tests: run nightly or staging only (slow, flaky, environment-specific)

### 2. Coverage targets by domain

**Backend (API, service):**
- Happy path: 80%+ coverage
- Error paths: 60%+ coverage
- Edge cases: prioritize critical (payment, auth, data loss)
- Don't test third-party libraries or framework code (assume it works)
- Test your usage of them

**Frontend (React, Vue, etc.):**
- Component behavior: 80%+ coverage (what user sees and does)
- Edge cases: focus on interactive (modal, form submission, error states)
- Styles: rarely test (implementation detail)
- Router/navigation: integration test (not unit)
- API integration: mock backend, test happy path + error case

**Database:**
- Schema: version control, test migrations (are they additive/safe?)
- Queries: test important ones (complex JOINs, business-critical)
- Don't test ORM internals (test your usage of ORM)

**DevOps/Infrastructure:**
- Terraform/IaC: test syntax (`terraform validate`), resource count
- CI/CD: test pipeline runs, artifacts created, secrets not exposed
- End-to-end: deploy to test environment, run smoke tests

### 3. What NOT to test

- **Framework internals:** don't test React hooks or Express middleware (they're tested by framework)
- **Third-party libraries:** don't test npm packages (trust they work)
- **Implementation details:** test behavior, not internal state (easier to refactor)
- **UI appearance:** colors, fonts (test with visual regression if critical)
- **Timing-dependent:** flaky tests (sleep, `setTimeout`) → use event/callback instead
- **Environment setup:** don't test database/server startup in unit tests (do in integration tests)

### 4. Framework selection by language

**JavaScript/Node.js:**
- **Jest:** default choice (mocking, snapshots, coverage)
- **Vitest:** fast alternative to Jest (better ESM support)
- **Mocha + Chai:** lightweight, flexible
- **React Testing Library:** component testing (user-centric)
- **Cypress/Playwright:** E2E testing (browser automation)

**Python:**
- **pytest:** de facto standard (fixtures, plugins, parametrization)
- **unittest:** built-in, verbose but works
- **hypothesis:** property-based testing (generate random inputs)
- **mock/unittest.mock:** built-in mocking
- **pytest-cov:** coverage reporting

**Go:**
- **testing (built-in):** standard, simple
- **testify:** assertions and mocks
- **table-driven tests:** Go idiom (`t.Run(tc.name, ...)` subtest per case)
- **httptest:** HTTP mocking
- **Integration tests in _test.go files alongside source

**Rust:**
- **cargo test (built-in):** unit tests in a `#[cfg(test)] mod tests` at the bottom of each source file; integration tests in the crate's `tests/` directory
- **`#[bench]`:** built-in benchmarks are nightly-only; use `criterion` on stable
- **rstest:** parametrized tests
- **mockall:** general-purpose mocking (traits and structs)
- **mockito:** HTTP server mocking (stub external HTTP APIs)
- **proptest:** property-based testing

### 5. Test data and fixtures

**Fixtures:**
- Reusable test data (user object, API response)
- Stored as JSON, YAML, or factory functions
- Keep realistic (use actual data shape, not simplified)
- Share via `beforeEach()` or fixture setup

**Factory functions (better than static fixtures):**
- Generate test data on-the-fly with defaults
- Easy to override specific fields
- Clearer intent in test (shows what matters, not irrelevant fields)
- Example:
```python
def make_user(name="Alice", email="alice@example.com", password_hash=None, role="user", verified=True):
    return {"id": 1, "name": name, "email": email, "password_hash": password_hash,
            "role": role, "verified": verified}

# In test:
user = make_user(role="admin")  # Clear which field matters
```

**Database fixtures (integration tests):**
- Real database (test database, not production)
- Seed with known data via migrations or fixtures
- Clean between tests (rollback or truncate)
- Use transactions to isolate tests

### 6. Mocking strategy

**When to mock:**
- External services (API, payment processor)
- Database (unit tests only; use real DB in integration)
- Time/clock (testing time-based logic)
- File system (testing file operations)
- Random numbers (testing randomized behavior)

**When NOT to mock:**
- Internal functions (test the integration)
- Standard library (trust it works)
- Application code you're testing (test real behavior)

**Mocking levels:**
1. **Spy:** record calls to real function, still execute it
2. **Stub:** replace with fake implementation, don't execute original
3. **Mock:** stub + verify it was called (verify interaction)

**Best practice:** mock at boundary (API calls, DB queries), not internal

### 7. Flaky test prevention

**Common causes:**
- **Timing:** `sleep()`, `setTimeout()`, race conditions
- **Environment:** test assumes specific state not guaranteed
- **Randomness:** random data, random port selection
- **External services:** network timeouts, service unavailability
- **Shared state:** test affects another test, run order matters

**Prevention:**
- Use event/callback instead of sleep
- Random port selection and cleanup (don't hardcode)
- Use fixed random seed for reproducibility
- Mock external services (don't call real API)
- Isolate tests (each test cleans up after itself)
- Run tests in different orders (detect order dependency)

### 8. Test organization

**File structure:**
Follow the language's convention, and pick one layout per project:
- JavaScript/TypeScript: co-locate tests with source (`Button.tsx` and `Button.test.tsx`), or a `__tests__/` folder next to the code it tests (`src/utils/__tests__/math.test.ts`); not a single top-level folder far from the source
- Python (pytest): a top-level `tests/` directory mirroring the package structure (`src/app/utils.py` → `tests/test_utils.py`)
- Go: `_test.go` files in the same package directory as the code
- Rust: unit tests inline in `#[cfg(test)]` modules; integration tests in `tests/`
- Within any layout, organize by feature or domain, not by test type

**Test naming:**
- Describe user behavior or system behavior
- Bad: `test_function_works`
- Good: `displays_error_when_email_invalid` or `updates_user_on_successful_save`

**Test structure (AAA pattern):**
```python
def test_login_with_valid_credentials():
    # Arrange: set up data and state
    user = make_user(email="alice@example.com", password_hash=hash_password("secret"))
    db.insert(user)
    
    # Act: perform the action
    result = login("alice@example.com", "secret")
    
    # Assert: verify the result
    assert result.success is True
    assert result.user_id == user["id"]
```

### 9. Coverage metrics

**Lines of code (LOC):** what % of code is executed by tests
- Target: 80%+ for critical paths
- Don't aim for 100% (diminishing returns, incentivizes testing implementation details)
- Focus on branch coverage (all decision paths), not line coverage

**Branch coverage:** what % of conditions (if/else, loops) are tested
- More important than line coverage
- Catch edge cases better

**Tools:**
- pytest-cov (Python): `pytest --cov=module`
- Jest (JavaScript): `jest --coverage`
- cargo tarpaulin (Rust): `cargo tarpaulin`
- go test (Go): `go test -cover`

**Coverage report usage:**
- Identify uncovered code
- Focus on critical paths and error handling
- Don't chase 100% (test behavior, not lines)

### 10. Continuous testing practices

**Local development:**
- Run tests on save (watch mode): `jest --watch`, `pytest-watch`
- Run before commit (pre-commit hook)
- Run full suite before PR submission

**CI/CD:**
- All tests run on every commit
- Fail if coverage drops
- Separate unit/integration/E2E (different timing)
- Cache dependencies (faster runs)
- Run in parallel where possible

**Regression testing:**
- Keep tests for bugs that were fixed (prevent reoccurrence)
- Add test before fixing bug (verify fix works)
- Document in test why this matters

---

## Pattern Decision Matrix

| Scenario | Pattern | Rationale |
|----------|---------|-----------|
| New API endpoint | Happy path unit test + integration test | Fast iteration, test behavior |
| Bug fix | Add failing test first, then fix | Verify fix works, prevent regression |
| Complex query | Integration test with real DB | Can't unit test query correctness |
| React component | React Testing Library (behavior) | Test what user sees, not implementation |
| Authentication | Unit test logic, integration test flow | Auth is critical, test thoroughly |
| Payment processing | Mock external API + integration test | Real payment API too risky in test |
| Performance optimization | Benchmark before/after | Verify improvement actually happens |
| Flaky test | Identify root cause (timing, randomness) | Fix cause, not symptom (don't add sleep) |

---

## Checklist

Work through `${CLAUDE_SKILL_DIR}/checklist.md` (`checklist.md` in this skill's directory) before finalizing a test strategy. Mark any item that doesn't apply as N/A and say why; don't skip it silently.

## Report

List findings in a table, most severe first, with the columns Severity | Finding | Evidence | Fix. Evidence is a file and line, or the test or config file the finding is about. Use this scale:

- **Critical**: causes a bug, outage, data loss, or security problem now. Blocks shipping.
- **High**: a real problem with a plausible path to harm. Fix before release.
- **Medium**: adds risk or cost, but needs another condition before it causes harm.
- **Low**: an improvement or hardening opportunity, not a current problem.

End with what you didn't check and why. If an area is clean, say so in one line. An agent that defines its own report format uses that format instead.

