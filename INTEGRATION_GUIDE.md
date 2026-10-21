# agent-forge Integration Guide

**Version:** 1.0  
**Date:** October 19, 2026  
**Status:** Production Ready

The installer has been extended with CLI and HTTP API capabilities for programmatic access to agents and skills. This enables seamless integration with IDEs, services like Codex, Replit, and custom tools.

---

## Quick Start

### Install agents/skills
```bash
agent-install -install
```

### Query an agent (CLI)
```bash
agent-install -agent api-dev
agent-install -agent go-reviewer
agent-install -agent fastapi-dev
agent-install -agent django-dev
```

### Query a skill (CLI)
```bash
agent-install -skill testing-patterns
agent-install -skill frontend-accessibility
```

### Auto-detect reviewer (CLI)
```bash
agent-install -review mycode.py    # → python-reviewer
agent-install -review script.go    # → go-reviewer
agent-install -review app.js       # → node-reviewer
```

### Start HTTP API server
```bash
agent-install -serve
# Server starts on http://127.0.0.1:8081 (this machine only)
```

---

## CLI Commands

### Installation & Updates

```bash
# List files to be installed
agent-install -contents

# Preview installation (dry-run)
agent-install -dryrun

# Install agents/skills to ~/.claude
agent-install -install

# Install to custom location
agent-install -install -target ~/my-agents

# Install specific version
agent-install -install -version v1.0.0

# Install from different repo
agent-install -install -repo myorg/agents -version main
```

### Agent Queries

```bash
# Get agent guidance
agent-install -agent api-dev
agent-install -agent go-reviewer
agent-install -agent python-reviewer
agent-install -agent test-writer
agent-install -agent fastapi-dev
agent-install -agent django-dev

# List available agents
ls ~/.claude/agents/
```

Use `fastapi-dev` for FastAPI routers, dependencies, Pydantic schemas, async services, and endpoint tests. Use `django-dev` for Django views, forms, templates, models, migrations, admin, and DRF APIs where DRF is already used. Both agents follow the shared Python rules and the API rules when applicable; neither requires a new framework-specific rule file.

Use `python-dev` for general Python or Flask work, and `api-dev` for framework-neutral API contracts. Framework implementation can be reviewed by `python-reviewer`, with `api-reviewer` for API changes.

### Skill Queries

```bash
# Get skill methodology + checklist
agent-install -skill testing-patterns
agent-install -skill frontend-accessibility
agent-install -skill performance-assessment
agent-install -skill database-optimization

# List available skills
ls ~/.claude/skills/
```

### Reviewer Suggestions

```bash
# Auto-detect reviewer by file extension
agent-install -review backend.py       # → python-reviewer
agent-install -review service.go       # → go-reviewer
agent-install -review component.tsx    # → node-reviewer
agent-install -review library.rs       # → rust-reviewer

# Supported extensions:
# .py       → python-reviewer
# .go       → go-reviewer
# .js/.ts   → node-reviewer
# .jsx/.tsx → node-reviewer
# .rs       → rust-reviewer
# (other)   → code-reviewer
```

### HTTP Server

```bash
# Start HTTP API server on 127.0.0.1:8081 (this machine only)
agent-install -serve

# Listen on a different address or port
agent-install -serve -addr 127.0.0.1:9090

# Server runs indefinitely, press Ctrl+C to stop
```

By default the server only accepts connections from the same machine. It has
no authentication, so only use `-addr 0.0.0.0:8081` (all network interfaces)
when you need it, for example inside a container, and don't expose it to an
untrusted network.

---

## HTTP API Reference

### Base URL
```
http://localhost:8081
```

### Endpoints

#### Health Check
```bash
GET /health

curl http://localhost:8081/health
# Response: {"status":"ok"}
```

#### List All Agents
```bash
GET /agents

curl http://localhost:8081/agents
# Response: {"agents":["api-dev","code-reviewer",...]}
```

#### Get Agent Details
```bash
GET /agent/{name}

curl http://localhost:8081/agent/api-dev
# Response: Full agent markdown content
```

#### List All Skills
```bash
GET /skills

curl http://localhost:8081/skills
# Response: {"skills":["api-design-patterns","testing-patterns",...]}
```

#### Get Skill Details
```bash
GET /skill/{name}
GET /skill/{name}/SKILL.md
GET /skill/{name}/checklist.md

curl http://localhost:8081/skill/testing-patterns
curl http://localhost:8081/skill/testing-patterns/checklist.md
```

#### Root Help
```bash
GET /

curl http://localhost:8081/
# Response: API documentation
```

---

## Integration Examples

### 1. Shell Script Wrapper

Create `agent-query.sh`:
```bash
#!/bin/bash

case "$1" in
  agent)
    agent-install -agent "$2"
    ;;
  skill)
    agent-install -skill "$2"
    ;;
  review)
    agent-install -review "$2"
    ;;
  serve)
    agent-install -serve
    ;;
  *)
    echo "Usage: agent-query.sh {agent|skill|review|serve} [name]"
    exit 1
    ;;
esac
```

Usage:
```bash
chmod +x agent-query.sh
./agent-query.sh agent api-dev
./agent-query.sh skill testing-patterns
./agent-query.sh review mycode.py
./agent-query.sh serve
```

### 2. Systemd Service (Linux/macOS with Homebrew)

Create `/etc/systemd/user/agent-forge.service` (Linux) or use LaunchAgent (macOS):

```ini
[Unit]
Description=agent-forge HTTP Server
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/agent-install -serve
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=default.target
```

Start:
```bash
systemctl --user enable agent-forge
systemctl --user start agent-forge
systemctl --user status agent-forge
```

### 3. VS Code Extension (Custom)

In `.vscode/settings.json`:
```json
{
  "agent-forge.enable": true,
  "agent-forge.installerPath": "~/.local/bin/agent-install",
  "agent-forge.serverPort": 8081
}
```

In VS Code task:
```json
{
  "label": "agent-forge: Suggest Reviewer",
  "type": "shell",
  "command": "agent-install -review ${file}",
  "presentation": { "reveal": "always" }
}
```

### 4. GitHub Copilot Integration

Add to `.instructions.md` in your repo:
```markdown
# Agent-Forge Integration

Reference these agents when generating code:

**API Development:**
- Use: `~/.claude/agents/api-dev.md`
- Review with: `~/.claude/agents/api-reviewer.md`
- Patterns in: `~/.claude/skills/api-design-patterns/`

**Testing:**
- Patterns in: `~/.claude/skills/testing-patterns/SKILL.md`
- Checklist in: `~/.claude/skills/testing-patterns/checklist.md`

**Frontend:**
- Accessibility: `~/.claude/skills/frontend-accessibility/`
- Standards: `~/.claude/rules/frontend-dev.md`

When writing code, ensure it follows the relevant rules file.
```

Then in Copilot Chat:
```
@copilot refer to ~/.claude/agents/api-dev for REST patterns
@copilot check against ~/.claude/skills/testing-patterns
```

### 5. Replit Integration

Add to `.replit`:
```bash
run = "agent-install -serve"

[env]
AGENT_FORGE_PORT = "8081"
```

In Replit.nix:
```nix
pkgs.with (import <nixpkgs> {}); [
  go
  curl
]
```

Create `replit-agents.md` in project:
```markdown
# Using agent-forge on Replit

1. Start the server:
   ```
   claude-install -serve
   ```

2. In another terminal, query agents:
   ```
   curl http://localhost:8081/agent/api-dev
   curl http://localhost:8081/skill/testing-patterns
   ```

3. Use with Replit AI by referencing the API
```

### 6. Docker Container

Create `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN ./build.sh

FROM alpine:latest
RUN apk add --no-cache curl
COPY --from=builder /app/bin/agent-install-linux-amd64 /usr/local/bin/agent-install
RUN chmod +x /usr/local/bin/agent-install
RUN agent-install -install

EXPOSE 8081
# Inside a container the server must listen on all interfaces so the
# published port can reach it.
CMD ["agent-install", "-serve", "-addr", "0.0.0.0:8081"]
```

Build and run:
```bash
docker build -t agent-forge .
# Publish the port on the host's loopback only
docker run -p 127.0.0.1:8081:8081 agent-forge
```

Query from host:
```bash
curl http://localhost:8081/agents
```

### 7. Python Client Library

Create `agent_forge.py`:
```python
import requests
import json
from pathlib import Path

class AgentForge:
    def __init__(self, base_url="http://localhost:8081"):
        self.base_url = base_url
    
    def health(self):
        return requests.get(f"{self.base_url}/health").json()
    
    def agents(self):
        return requests.get(f"{self.base_url}/agents").json()
    
    def agent(self, name):
        return requests.get(f"{self.base_url}/agent/{name}").text
    
    def skills(self):
        return requests.get(f"{self.base_url}/skills").json()
    
    def skill(self, name, file="SKILL.md"):
        return requests.get(f"{self.base_url}/skill/{name}/{file}").text
    
    def detect_reviewer(self, filepath):
        ext_map = {
            ".py": "python-reviewer",
            ".go": "go-reviewer",
            ".js": "node-reviewer",
            ".ts": "node-reviewer",
            ".rs": "rust-reviewer",
        }
        ext = Path(filepath).suffix
        reviewer = ext_map.get(ext, "code-reviewer")
        return self.agent(reviewer)

# Usage
if __name__ == "__main__":
    af = AgentForge()
    
    # Get agent
    print(af.agent("api-dev")[:200])
    
    # Get skill
    print(af.skill("testing-patterns")[:200])
    
    # Auto-detect reviewer
    print(af.detect_reviewer("mycode.py")[:200])
```

### 8. Node.js Client

Create `agent-forge.js`:
```javascript
const http = require('http');

class AgentForge {
  constructor(baseUrl = 'http://localhost:8081') {
    this.baseUrl = baseUrl;
  }

  async fetch(path) {
    return new Promise((resolve, reject) => {
      http.get(`${this.baseUrl}${path}`, (res) => {
        let data = '';
        res.on('data', chunk => data += chunk);
        res.on('end', () => {
          try {
            resolve(path.includes('.md') ? data : JSON.parse(data));
          } catch {
            resolve(data);
          }
        });
      }).on('error', reject);
    });
  }

  async health() {
    return this.fetch('/health');
  }

  async agents() {
    return this.fetch('/agents');
  }

  async agent(name) {
    return this.fetch(`/agent/${name}`);
  }

  async skill(name, file = 'SKILL.md') {
    return this.fetch(`/skill/${name}/${file}`);
  }

  detectReviewer(filepath) {
    const extMap = {
      '.py': 'python-reviewer',
      '.go': 'go-reviewer',
      '.js': 'node-reviewer',
      '.ts': 'node-reviewer',
      '.rs': 'rust-reviewer'
    };
    const ext = filepath.split('.').pop();
    const reviewer = extMap[`.${ext}`] || 'code-reviewer';
    return this.agent(reviewer);
  }
}

module.exports = AgentForge;
```

---

## Workflow Examples

### API Development Workflow

```bash
# 1. Get API development guidance
agent-install -agent api-dev

# 2. Build the endpoint
# ... write code ...

# 3. Get API patterns for this endpoint
agent-install -skill api-design-patterns

# 4. Get review guidance
agent-install -agent api-reviewer

# 5. Query specific testing patterns
agent-install -skill testing-patterns
```

### Performance Investigation Workflow

```bash
# 1. Run performance assessment
agent-install -skill performance-assessment

# 2. Implement optimization
agent-install -agent performance-optimizer

# 3. Verify with testing patterns
agent-install -skill testing-patterns
```

### Accessibility Audit Workflow

```bash
# 1. Get accessibility guidance
agent-install -skill frontend-accessibility

# 2. Review frontend code
agent-install -agent frontend-reviewer

# 3. Use checklist for validation
agent-install -skill frontend-accessibility | grep -A 50 "CHECKLIST"
```

---

## File Locations

After `agent-install -install`, files are installed to:

```
~/.claude/
├── agents/
│   ├── api-dev.md
│   ├── api-reviewer.md
│   ├── code-reviewer.md
│   ├── database-dev.md
│   ├── debugger.md
│   ├── dependency-manager.md
│   ├── devops-dev.md
│   ├── devops-reviewer.md
│   ├── django-dev.md
│   ├── docker-dev.md
│   ├── documentation-writer.md
│   ├── fastapi-dev.md
│   ├── frontend-dev.md
│   ├── frontend-reviewer.md
│   ├── git-manager.md
│   ├── go-dev.md
│   ├── go-reviewer.md
│   ├── node-dev.md
│   ├── node-reviewer.md
│   ├── performance-optimizer.md
│   ├── python-dev.md
│   ├── python-reviewer.md
│   ├── rust-dev.md
│   ├── rust-reviewer.md
│   ├── security-reviewer.md
│   ├── task-planner.md
│   ├── test-writer.md
│   └── ui-designer.md
│
├── rules/
│   ├── api-dev.md
│   ├── devops-dev.md
│   ├── docker-dev.md
│   ├── frontend-dev.md
│   ├── git-workflow.md
│   ├── go-dev.md
│   ├── nodejs-dev.md
│   ├── python-dev.md
│   ├── rust-dev.md
│   └── sql-dev.md
│
└── skills/
    ├── api-design-patterns/
    │   ├── SKILL.md
    │   └── checklist.md
    ├── database-optimization/
    │   ├── SKILL.md
    │   └── checklist.md
    ├── frontend-accessibility/
    │   ├── SKILL.md
    │   └── checklist.md
    ├── performance-assessment/
    │   ├── SKILL.md
    │   └── checklist.md
    ├── security-assessment/
    │   ├── SKILL.md
    │   ├── checklist.md
    │   └── scripts/scan-secrets.py
    └── testing-patterns/
        ├── SKILL.md
        └── checklist.md
```

---

## Troubleshooting

### Agent not found

```bash
# Ensure agents are installed
agent-install -install

# Verify installation
ls ~/.claude/agents/api-dev.md

# Query should work
agent-install -agent api-dev
```

### HTTP server won't start

```bash
# Check if port 8081 is in use
lsof -i :8081

# Use a different port
agent-install -serve -addr 127.0.0.1:9090

# OR kill existing process
kill -9 $(lsof -t -i:8081)

# Then try again
agent-install -serve
```

### Binary not found

```bash
# Add to PATH
export PATH="$PATH:~/.local/bin"

# OR use full path
/path/to/agent-install -agent api-dev
```

### Server connection refused

```bash
# Verify server is running
curl http://localhost:8081/health

# If it fails, start the server
agent-install -serve &

# Then try again
curl http://localhost:8081/agents
```

---

## Contributing

To add new agents or skills:

1. Create agent in `configs/agents/{name}.md` or skill in `configs/skills/{name}/SKILL.md`
2. Follow YAML frontmatter format
3. Add checklist to `configs/skills/{name}/checklist.md`
4. Commit to GitHub
5. Tag release (e.g., `v1.1.0`)
6. Users update with: `agent-install -install -version v1.1.0`

---

## Version History

- **v1.0 (Oct 19, 2026):** Initial release with CLI and HTTP API support

---

## License

Same as agent-forge project.
