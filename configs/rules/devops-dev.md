---
paths:
  - ".github/workflows/*.yml"
  - ".github/workflows/*.yaml"
  - "azure-pipelines.yml"
  - "Jenkinsfile"
  - "*.tf"
  - "*.bicep"
  - "kustomization.yaml"
  - "*.yaml"
---

# DevOps Development Rules

## Infrastructure as Code (IaC) principles

### Naming and organization
- Resource names follow a consistent pattern: `{environment}-{service}-{component}`. Example: `prod-api-storage`, `dev-web-lb`.
- All resources tagged consistently: `Environment`, `Service`, `Owner`, `CostCenter`, `ManagedBy: IaC`.
- Terraform/Bicep organized by environment and component: `terraform/prod/main.tf`, `bicep/dev/app-service.bicep`.
- Never hardcode resource names in code; use variables and naming conventions.

### Secret and configuration management
- No secrets (passwords, keys, tokens, connection strings) in code or config files.
- Secrets stored in external secret management: Azure Key Vault, AWS Secrets Manager, HashiCorp Vault.
- Environment-specific configuration (staging, prod) separated and never mixed in code.
- Database passwords, API keys, certificates: all external, never checked in, rotated regularly.
- Secrets referenced through environment variables or secret injection at runtime.

### State management (Terraform)
- Terraform state stored remotely (not locally): Azure Storage, S3, or equivalent.
- State file protected (encryption at rest, encryption in transit, access control).
- State locking enabled to prevent concurrent modifications.
- State backups automated and tested.
- Never commit `terraform.tfstate` to version control.
- `terraform plan` reviewed before `terraform apply` (separate review step, no auto-apply).

### Idempotency and repeatability
- IaC must be idempotent: running twice produces same result as running once.
- No manual changes to infrastructure outside IaC (drift). If manual change needed, update IaC and re-apply.
- IaC must handle existing resources gracefully (import existing resources if needed).
- Tear-down (destroy) must be safe and reversible (delete dev/test, keep prod).

## CI/CD Pipeline standards

### Build and test phases
- Lint: syntax, code style, security scan (no hardcoded secrets)
- Build: compile, package, create artifacts
- Test: unit tests, integration tests, contract tests (no e2e in CI)
- Quality gates: fail on critical lint violations, <X% test coverage, security findings
- Artifact storage: versioned, named by commit/tag, cleanup old artifacts

### Deployment strategies
- Development: auto-deploy on main branch
- Staging: auto-deploy, or manual trigger with approval
- Production: manual trigger with approval gate (not auto)
- Rollback plan: previous version easily restored, or blue-green deployment
- Deployment observability: logs, metrics, alerts on failure

### Environment promotion
- Promotions follow a clear path: dev → staging → prod
- Same artifacts promoted (build once, test multiple times, deploy same artifact)
- Configuration changes per environment (database URLs, API keys), not code changes
- Production deployments recorded: who, when, what version, result

### Secrets in pipelines
- Secrets stored in pipeline secret store (GitHub Secrets, Azure Key Vault, pipeline-native)
- Secrets never logged or printed in pipeline output
- Secrets rotated regularly and independently of code
- Pipeline logs don't contain sensitive data
- Secret access audited (who accessed, when)

## Container practices (Docker/Kubernetes)

### Container image standards
- Non-root user (see Docker rules in docker-dev.md)
- Multi-stage builds (no build tools in final image)
- Security scanning before push (vulnerability scan, signed images if required)
- Image tagging: semantic versioning for releases, `main` for latest development
- Base image updates tracked (security patches, deprecations)
- No `latest` tag in production deployments

### Kubernetes deployments (if applicable)
- Resource requests/limits defined (CPU, memory)
- Health checks: liveness and readiness probes
- Graceful shutdown: terminationGracePeriodSeconds, signal handling
- Pod security policies enforced (no privileged containers, no root)
- Network policies restrict traffic between pods
- RBAC: service accounts with minimal permissions

## Cost optimization

### Resource sizing
- Right-size instances (don't over-provision, don't under-provision)
- Use smaller instances for dev/test (staging can be minimal)
- Auto-scaling policies: scale down when load low, scale up when needed
- Spot/preemptible instances for non-critical workloads (dev, batch jobs)
- Storage tiering: hot, warm, cold storage based on access patterns

### Cost monitoring
- Budget alerts: notify if spending exceeds forecast
- Cost breakdown by service/environment: know what costs what
- Reserved instances or savings plans for predictable workloads
- Cleanup: remove unused resources (old DBs, orphaned disks, unused IPs)

## Disaster recovery and backup

### Backup strategy
- Critical data backed up (databases, file storage)
- Backup frequency based on RPO (recovery point objective): hourly, daily, weekly
- Backup tested regularly (restore from backup to verify it works)
- Backup retention policy: keep backups for X days/weeks/months
- Backups geo-replicated (backup in different region if primary region fails)

### Disaster recovery plan
- RTO (recovery time objective) defined: how fast must service recover?
- Runbook documented: steps to restore service from backup
- Failover procedure tested: actually test failover, not just documented
- Recovery from natural disaster: infrastructure in multiple regions/availability zones

## Monitoring and alerting

### Observability
- Logs aggregated: logs from all components sent to central logging (DataDog, Splunk, etc.)
- Log format consistent and parseable (JSON, structured logging)
- Metrics collected: CPU, memory, disk, network, application-specific metrics
- Traces captured for distributed tracing (OpenTelemetry, Jaeger)

### Alerting
- Alerts for critical issues: service down, error rate > X%, latency > Y ms
- Alert routing: ops on-call via PagerDuty, email, Slack based on severity
- Alert silence window for planned maintenance
- Alert runbooks: how to respond to each alert
- No alert fatigue: only alert on actionable issues, not information

### Dashboards
- Dashboards for each service: health, performance, errors
- Incident dashboard: correlate errors, latency, resource usage
- Business metrics dashboard: user signups, transactions, revenue impact
- On-call dashboard: quick view of system health

## Security practices

### Access control
- Infrastructure access via IAM roles/RBAC (not shared credentials)
- Service accounts for automated tasks, minimal permissions (principle of least privilege)
- MFA for human access to production
- SSH key rotation, no shared keys
- Audit logging: who accessed what resource, when

### Network security
- Firewall rules: minimize exposed ports (443 for HTTPS, 22 for SSH, others internal only)
- VPC/network isolation: internal services not exposed to internet
- DDoS protection for public-facing endpoints
- WAF (Web Application Firewall) for web apps
- TLS/SSL everywhere: encrypt data in transit

### Compliance and auditing
- Audit logs: all infrastructure changes logged and retained
- Compliance framework: understand requirements (SOC2, PCI-DSS, HIPAA, etc.)
- Configuration compliance: automated checks that config follows policy
- Data retention policy: logs kept for X days/years per compliance

## Documentation and runbooks

### Infrastructure documentation
- Diagram of infrastructure: services, databases, load balancers, networks
- Component descriptions: what each service does, dependencies
- Deployment procedure: how to deploy new version
- Rollback procedure: how to revert to previous version
- Known issues and limitations: documented clearly

### Operational runbooks
- Service dependencies: what else must be running for this service to work
- Troubleshooting guide: common issues and resolution
- On-call playbook: escalation path, who to contact
- Incident response: steps to take when service is down
- Recovery procedures: how to restore from backup, failover

## Common anti-patterns to avoid

- ❌ Hardcoding secrets or connection strings in IaC/code
- ❌ Manual infrastructure changes (drift, inconsistency)
- ❌ No version control for IaC (can't audit changes, can't rollback)
- ❌ Auto-deploying to production without approval
- ❌ No backups or untested recovery procedures
- ❌ Over-provisioned resources in dev/test (cost waste)
- ❌ No monitoring or alerting (find out about problems from users)
- ❌ Shared credentials or API keys (audit trail lost, can't rotate individually)
- ❌ No runbooks for incident response (chaos during outage)
- ❌ Never updating infrastructure (security patches missed, dependencies stale)

