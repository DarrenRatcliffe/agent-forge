# Testing Patterns Checklist

## Test Strategy
- [ ] Test pyramid ratio defined (70% unit, 20% integration, 10% E2E)
- [ ] Coverage targets set by domain (backend, frontend, database, infra)
- [ ] Framework selected for each language/domain
- [ ] What to test identified (behavior, critical paths, error cases)
- [ ] What not to test identified (framework internals, third-party libs)

## Test Organization
- [ ] File structure planned (co-located, mirrored source structure)
- [ ] Naming convention defined (behavior-focused, not function-focused)
- [ ] AAA pattern (Arrange, Act, Assert) understood
- [ ] Test data strategy chosen (fixtures, factories, real DB for integration)
- [ ] Mocking strategy defined (boundaries, levels, when to mock)

## Coverage
- [ ] Line coverage target set (80%+ for critical)
- [ ] Branch coverage target set (all decision paths)
- [ ] Coverage monitoring tool configured
- [ ] Coverage reporting in CI
- [ ] Alert if coverage drops

## Flaky Test Prevention
- [ ] Timing dependencies identified and eliminated
- [ ] Environment isolation verified
- [ ] Shared state issues addressed
- [ ] Random data handling managed
- [ ] External service mocking implemented
- [ ] Test order independence confirmed

## Local Development
- [ ] Watch mode setup (`jest --watch`, `pytest-watch`, etc.)
- [ ] Pre-commit hook for tests (optional but recommended)
- [ ] Test running instructions documented
- [ ] Common test commands documented (run all, run one, coverage)

## CI/CD Integration
- [ ] Unit tests run on every commit
- [ ] Integration tests run on every commit
- [ ] E2E tests run nightly or on staging
- [ ] Fail if coverage drops
- [ ] Dependencies cached (faster runs)
- [ ] Tests run in parallel where possible
- [ ] Test results reported (pass/fail clearly)

## Test Maintenance
- [ ] Regression tests added for fixed bugs
- [ ] Tests updated when feature changes
- [ ] Flaky tests investigated and fixed (not ignored)
- [ ] Test coverage reviewed regularly
- [ ] Deprecated tests removed

## Documentation
- [ ] Test strategy documented (pyramid, coverage targets)
- [ ] Framework choice documented with rationale
- [ ] Fixture/mock usage patterns documented
- [ ] Common issues and solutions documented
- [ ] How to run tests locally documented
- [ ] How to debug failing tests documented

