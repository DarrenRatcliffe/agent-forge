---
paths:
  - "**/*.py"
  - "**/pyproject.toml"
---

# Backend Rules — Python

## Stack assumptions
- Match the existing framework: FastAPI, Flask, or Django. Don't introduce a different framework into an existing codebase.
- Python 3.11+ features are fine unless the project's runtime target says otherwise — check `pyproject.toml` / `requirements.txt` before assuming.
- Type hints on every function signature. Not optional, not "where convenient." Enforced by a type checker (`mypy` or `pyright`) in CI.

## Structure
- Layered separation: routers/views → services → data access (repository or ORM layer). No DB queries or business logic inline in route handlers.
- Pydantic models (FastAPI), explicit schema validation (Flask), or Django forms / DRF serializers where applicable for request validation. API response bodies use explicit schemas or serializers; native Django HTML responses use templates and context. No passing unvalidated raw dicts through the stack.
- Settings via a single validated config object (`BaseSettings` from the separate `pydantic-settings` package in Pydantic v2, or equivalent), or Django's settings and existing environment-loading mechanism, not scattered `os.environ.get()` calls.

## Style and tooling
- Formatted with `ruff format` or `black` (match what the project uses), linted with `ruff`. Both run in CI, not just locally on trust.
- No bare `except:`. Catch specific exceptions. If you genuinely need to catch everything, say why in a comment.
- f-strings for formatting. No `%`-formatting, no `.format()` unless working in an existing codebase that already standardises on it.

## Async
- FastAPI: use `async def` for I/O-bound route handlers (DB calls, HTTP calls to other services). Don't mix sync blocking calls inside an async handler — it blocks the event loop.
- If a route is CPU-bound, sync is fine, or offload to a background worker. Don't force async where it doesn't help.

## Error handling
- Custom exception classes for domain errors, caught by a centralised exception handler that maps them to HTTP responses. Don't hand-build error JSON in every handler.
- Never expose internal exception messages or stack traces to the client outside local dev.
- Log the exception with context at the point it's caught, not just re-raised silently.

## Dependencies and environments
- `pyproject.toml` with pinned or constrained versions, virtual environment per project (`venv`, `poetry`, or `uv` — pick one per project, don't mix).
- Justify new dependencies. Standard library first if it does the job.

## Testing
- `pytest`, with fixtures for setup, not copy-pasted setup code in every test. For Django, follow the existing Django test runner or pytest-django setup.
- Test the service/business logic layer directly, not just through the HTTP layer, so failures point at the actual broken code.
- Mock external calls (third-party APIs, non-test databases). Don't mock the thing you're actually trying to test.

## Data and DB
- Migrations via Alembic (or the project's existing tool) — no manual schema changes against a live database.
- Parameterised queries only. If using raw SQL, no string interpolation of user input, ever.