---
paths:
  - "**/*.rs"
  - "Cargo.toml"
---

# Backend Rules — Rust

## Stack assumptions
- Use the edition and MSRV declared in `Cargo.toml`. Don't use features newer than the MSRV.
- `cargo fmt` and `cargo clippy --all-targets -- -D warnings` clean, enforced in CI, not on trust.
- One async runtime per project (usually `tokio`). Don't mix runtimes.

## Structure
- Small modules with one responsibility. Use a workspace for multi-crate projects.
- Keep `main.rs` thin: parse config, build dependencies, call into the library crate. Put logic in `lib.rs` and modules so it can be tested.
- Keep public APIs small. Default to private; use `pub(crate)` before `pub`.
- Config is parsed once at startup into a typed struct, not read ad hoc from the environment.

## Error handling
- Return `Result`. No `unwrap()` or `expect()` outside tests, or where the invariant is provable and documented in the `expect` message.
- Libraries define typed errors (`thiserror`). Applications may use `anyhow` at the top level with `.context(...)` on fallible calls.
- Use `?` to propagate. Don't swallow errors with `let _ =` without a comment explaining why that is safe.
- `panic!` is for bugs, not for expected failures like bad input or I/O errors.

## Ownership and API design
- Take `&str` and `&[T]` in parameters rather than `&String` and `&Vec<T>`. Take ownership only when the function needs it.
- Don't add `.clone()` just to satisfy the borrow checker. Restructure first, and justify the clone if it stays.
- Prefer iterators and combinators to index loops where they read clearer.
- Derive `Debug` on public types. Implement `Display` and `Error` for error types.
- Make invalid states unrepresentable: use enums and newtypes rather than booleans and bare strings.

## Unsafe
- No `unsafe` unless there is no safe alternative. Every `unsafe` block has a `// SAFETY:` comment stating the invariant that makes it sound.
- Keep `unsafe` blocks as small as possible and wrapped in a safe API. Call out any new `unsafe` in review.

## Concurrency and async
- Never block inside an `async fn` (no `std::thread::sleep`, blocking file or network I/O). Use `spawn_blocking` for blocking work.
- Don't hold a `std::sync::Mutex` guard across an `.await`.
- Spawned tasks have an owner that joins them or handles their errors. Support cancellation and graceful shutdown.
- Prefer message passing or well-scoped locks over shared mutable state.

## Logging
- Use `tracing` with structured fields. No `println!` debugging left in.
- No secrets, tokens, or PII in logs.

## Dependencies
- Justify each new crate: what it replaces and its maintenance status. Prefer the standard library when it does the job.
- Commit `Cargo.lock` for binaries. Keep the feature flags you enable minimal.
- Run `cargo audit` (or `cargo deny`) in CI.

## Testing
- Unit tests in `#[cfg(test)]` modules beside the code. Integration tests in `tests/`.
- Test behaviour through the public API. Use `Result`-returning tests rather than `unwrap` where useful.
- Doc examples compile and run (`cargo test` covers them), so keep them correct.
- Tests run in CI without manual setup.
