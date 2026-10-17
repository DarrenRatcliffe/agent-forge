---
paths:
  - "**/*.go"
  - "go.mod"
---

# Backend Rules — Go

## Stack assumptions
- Use the Go version declared in `go.mod`. Don't use language features newer than it.
- Standard library first. A new dependency needs a reason: what it replaces and its maintenance cost.
- `gofmt` (or `goimports`) clean and `go vet` clean, enforced in CI, not on trust.

## Structure
- Small, focused packages named for what they provide (`orders`, not `utils` or `common`).
- `cmd/<app>/main.go` stays thin: parse config, wire dependencies, call into packages. `internal/` for code that shouldn't be imported elsewhere.
- Accept interfaces, return concrete types. Define an interface where it is consumed, and keep it small.
- Config comes from environment or flags, parsed once at startup into a struct. No scattered `os.Getenv` calls.

## Error handling
- Handle every error. Never assign one to `_` without a comment saying why that is safe.
- Wrap with context: `fmt.Errorf("load order %s: %w", id, err)`. Check with `errors.Is` / `errors.As`, not string comparison.
- Return errors; don't `panic` outside `main` startup or truly unrecoverable programmer errors.
- Error strings are lowercase and have no trailing punctuation.
- Don't log and return the same error. Do one or the other, at the layer that owns it.

## Concurrency
- Pass `context.Context` as the first parameter to anything doing I/O or that can block. Don't store contexts in structs.
- Every goroutine has a defined way to stop and an owner who waits for it (`sync.WaitGroup`, `errgroup`). No fire-and-forget goroutines.
- Share memory by communicating where it's clear; otherwise guard shared state with a mutex. Document what each mutex protects.
- The sender closes channels, never the receiver. Run `go test -race` on concurrent code.

## Style
- Follow Effective Go and the Go Code Review Comments. Short names in small scopes, descriptive names in wide ones.
- Exported identifiers have doc comments that start with the identifier name.
- Table-driven tests where cases share logic. Avoid package-level mutable state.
- Use `defer` for cleanup directly after acquiring the resource. Remember `defer` in a loop runs at function exit.

## Logging
- Structured logging with `log/slog`. No `fmt.Println` debugging left in.
- No secrets, tokens, or PII in logs.

## HTTP and I/O
- Set timeouts on `http.Client` and `http.Server`. The default client has none.
- Always close response bodies. Validate and bound input size at the boundary.

## Testing
- `go test ./...` runs in CI with no manual setup.
- Test behaviour through the public API of the package. Use `t.Helper()`, `t.Cleanup()`, and `t.TempDir()`.
- Use fakes for external services rather than mocking the thing under test. Integration tests hit a real or realistic instance.

## Security baseline
- Parameterised queries only. No string-built SQL.
- No `exec.Command` with shell strings built from user input. Pass arguments separately.
- Use `crypto/rand` for anything security-relevant, never `math/rand`.
- Run `govulncheck ./...` in CI.
