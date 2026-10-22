---
paths:
  - "**/Dockerfile*"
  - "**/docker-compose*.y*ml"
  - "**/compose.y*ml"
  - "**/.dockerignore"
---

# Docker and Container Development Rules

## Dockerfile design principles

### Base images
- Choose minimal base images: `alpine`, `distroless`, or language-specific slim images (`python:3.13-slim`, `node:24-alpine`).
- Pin base image versions explicitly to a specific, currently supported tag, ideally also by digest (`FROM alpine:3.22@sha256:...`). No `latest` tags and no end-of-life releases; check the upstream support schedule when choosing. Example: `FROM alpine:3.22`, not `FROM alpine`.
- Document why a specific base image was chosen if it's larger or unusual (e.g., Ubuntu for specific tooling).
- Prefer distroless images when the application needs minimal dependencies (compiled binaries, Node.js apps with dependencies bundled).

### Multi-stage builds
- Use multi-stage builds to separate build environments from runtime environments.
- Build stage: include build tools, compilers, dev dependencies, source code.
- Runtime stage: copy only artifacts (binaries, compiled assets) from build stage. Never include build tools or source in final image.
- Naming: use descriptive stage names (`builder`, `dependencies`, `final`) not generic ones.

Example pattern:
```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/server .

FROM alpine:3.22
COPY --from=builder /app/server /usr/local/bin/server
CMD ["server"]
```

### Layers and caching
- Order Dockerfile instructions from least-frequently-changed to most-frequently-changed to maximize layer cache hits.
- Copy dependency files (go.mod, requirements.txt, package.json) early, run dependency installation, then copy application code.
- RUN instructions that install dependencies should be consolidated: `RUN apt-get update && apt-get install ... && rm -rf /var/lib/apt/lists/*` (clean apt cache to reduce layer size).
- Avoid unnecessary layers: combine related RUN commands with `&&` to reduce image size.

### Image size optimization
- Remove build artifacts, caches, and temporary files in the same RUN instruction that creates them.
- Use `.dockerignore` to exclude unnecessary files (see below).
- Don't install unnecessary packages. Use `--no-install-recommends` with apt-get.
- Consider multi-stage builds to avoid shipping build tools and intermediate dependencies.

### .dockerignore
- Include `.dockerignore` in the root of the project alongside the Dockerfile.
- Exclude: `.git`, `.gitignore`, `node_modules`, `*.log`, `.env`, `.env.local`, build output, test coverage, IDE files.
- Pattern: one per line, glob patterns supported (e.g., `**/*.pyc`).

## Security practices

### Running as non-root
- Never run containers as root (UID 0).
- Create an unprivileged user in the Dockerfile: `RUN useradd -m -u 1000 appuser` (or equivalent for Alpine: `addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser`).
- Switch to that user before CMD: `USER appuser`.
- Ensure the application's working directory and files are owned by that user.

### File permissions
- Set appropriate permissions on application files. Writable by the app, readable/executable where needed.
- Avoid `chmod 777` or overly permissive settings. Example: `RUN chmod 755 /app/script.sh`.

### Secret handling
- Never embed secrets, API keys, passwords, or tokens in the image or Dockerfile.
- Use build args for non-sensitive configuration only: `ARG BUILD_VERSION=1.0`.
- Build-time secrets (private registry or package feed tokens): use BuildKit secret mounts, `RUN --mount=type=secret,id=npm_token ...` with `docker build --secret id=npm_token,src=...`. Never pass secrets via `ARG` or `ENV`; they persist in the image history.
- Runtime secrets: prefer file-mounted secrets (Compose `secrets:`, Kubernetes secret volumes) over environment variables. Env vars are visible in `docker inspect` and leak easily into logs and child processes; use them only when the app can't read a file.
- Use Docker secrets (Docker Swarm) or external secret management (Kubernetes, CI/CD platform) in production.

### Image scanning
- Scan images for vulnerabilities during CI/CD: `docker scout cves` or third-party tools (Trivy, Snyk). `docker scan` has been removed.
- Fix high/critical vulnerabilities before deployment.
- Use a container registry that scans on push (Docker Hub, ECR, GCR, ACR).

## Docker Compose (local development)

### File structure
- Use `compose.yaml` as the filename (the Compose Spec's preferred name; an existing `docker-compose.yml` still works).
- Don't add a top-level `version:` key. The Compose Spec treats it as obsolete and Compose ignores it with a warning.
- Top-level keys: `services`, `volumes`, `networks`, `secrets` (if needed).

### Service definition
- Service names use snake_case: `api_server`, `database`, `cache`, `worker`.
- Each service maps to one container. Don't run multiple processes in one container.
- Ports exposed only if needed locally for manual testing. Avoid ports taken by OS services (5000 and 7000 are used by AirPlay Receiver on macOS); `8080`, `8000` or `3000` are typical.
- Mounts (`volumes`) for development: source code, config files, mounted read-write for fast iteration.
- Mount `/node_modules`, `/venv`, `/.cargo/target` as volumes (anonymous or named, not bind mounts) to avoid permission/performance issues.
- Environment variables: use `.env` file or `environment:` key. Never hardcode secrets; use `env_file:` and `.gitignore` the `.env`.

### Example:
```yaml
services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@postgres:5432/appdb
      - LOG_LEVEL=debug
    volumes:
      - .:/app
      - /app/node_modules  # anonymous volume, stops the bind mount hiding the image's node_modules
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - app

  postgres:
    image: postgres:17-alpine
    environment:
      - POSTGRES_PASSWORD=localpass
      - POSTGRES_USER=user
      - POSTGRES_DB=appdb
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks:
      - app

volumes:
  postgres_data:

networks:
  app:
    driver: bridge
```

### Networking
- Use named networks (not `links`, which is legacy).
- Service-to-service communication: use service name as hostname (`postgres` in the example above).
- Expose ports only for external access (developer machine, tests).

### Startup and health
- Use `depends_on` with `condition: service_healthy` for order-sensitive services (app depends on DB).
- Define `healthcheck` for long-running services (database, broker) so Compose knows when they're ready.
- Distinguish between "container started" and "service ready" (e.g., PostgreSQL starts but needs time to initialize).

## Build and deployment

### CI/CD integration
- Dockerfile should build cleanly in CI with an explicit tag, e.g. `docker build -t myapp:${GIT_SHA} .` (not `latest`).
- Use consistent image naming: `{registry}/{namespace}/{repo}:{tag}`. Example: `ghcr.io/company/myapp:v1.0.0`.
- Pin build-time dependencies in FROM statements, not with `latest`.
- Tag images with semantic versioning for releases (`v1.0.0`, not `main` or `latest`).
- If using build args, document them in README or CI config.

### Production readiness
- Containers should be stateless. Don't rely on container-local disk for persistence (data loss on restart).
- Logs to stdout/stderr (not files). Container orchestrators (Docker, Kubernetes) capture logs.
- Health checks for production deployments.
- Resource limits: set `--memory` and `--cpus` or equivalent in orchestration config.

## Common pitfalls to avoid

- ❌ Running as root
- ❌ Hardcoding secrets in Dockerfile or environment
- ❌ Using `latest` tag on base images
- ❌ Including build tools in final image (no multi-stage build)
- ❌ Bloated images from unnecessary dependencies
- ❌ No `.dockerignore`, shipping `.git`, `node_modules`, etc.
- ❌ Tightly coupling containers (assuming container names resolve without networks)
- ❌ Mixing logic and configuration (hardcoding in Dockerfile instead of env vars)
- ❌ No health checks for services that need time to initialize
- ❌ Not cleaning apt/pip/npm caches in same RUN instruction

