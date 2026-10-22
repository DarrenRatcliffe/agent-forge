---
paths:
  - ".github/workflows/**"
  - "**/azure-pipelines.yml"
  - "**/Jenkinsfile"
  - "**/*.tf"
  - "**/*.bicep"
  - "**/kustomization.yaml"
  - "**/k8s/**"
---

# DevOps Development Rules

## Infrastructure as Code (IaC)

### Naming and organization
- Resource names follow `{environment}-{service}-{component}`. Example: `prod-api-storage`, `dev-web-lb`.
- All resources tagged consistently: `Environment`, `Service`, `Owner`, `CostCenter`, `ManagedBy: IaC`.
- Terraform/Bicep organized by environment and component: `terraform/prod/main.tf`, `bicep/dev/app-service.bicep`.
- Never hardcode resource names in code; use variables and the naming convention.

### Secret and configuration management
- No secrets (passwords, keys, tokens, connection strings, certificates) in code or config files.
- Secrets stored in external secret management (Azure Key Vault, AWS Secrets Manager, HashiCorp Vault), injected at runtime through environment variables or secret references, and rotated independently of code.
- Environment-specific configuration (staging, prod) separated and never mixed in code.

### State management (Terraform)
- Terraform state stored remotely (Azure Storage, S3, or equivalent), never locally.
- State protected: encryption at rest and in transit, access control, locking enabled, backups automated.
- Never commit `terraform.tfstate` to version control.
- `terraform plan` reviewed before `terraform apply` (separate review step, no auto-apply).

### Idempotency and drift
- IaC must be idempotent: running twice produces the same result as running once.
- No manual changes to infrastructure outside IaC. If a manual change is unavoidable, update the IaC to match and re-apply.
- Existing resources are imported, not recreated.
- Destroy is safe for dev/test; production resources are protected from deletion.

## CI/CD pipelines

### Build and test phases
- Lint: syntax, code style, and a secret scan; IaC also gets a misconfiguration scan (e.g. checkov, tfsec, Azure Policy).
- Build: compile, package, create artifacts.
- Test: unit, integration, and contract tests (no e2e in CI).
- Quality gates: fail on critical lint violations, test coverage below the project's agreed threshold (or any drop from the main branch), and high/critical security findings.
- Artifacts versioned and named by commit or tag; old artifacts cleaned up.

### Deployment and promotion
- Development: auto-deploy on main branch.
- Staging: auto-deploy, or manual trigger with approval.
- Production: manual trigger with an approval gate, never auto-deploy.
- Promotion follows dev → staging → prod with the same artifact (build once). Environments differ by configuration, not code.
- Rollback plan: the previous version can be redeployed, or blue-green deployment.
- Production deployments recorded: who, when, what version, result. Failed deployments raise an alert.

### Secrets in pipelines
- Secrets stored in the pipeline secret store (GitHub Secrets, Azure Key Vault, pipeline-native), never in workflow files.
- Secrets and other sensitive data never logged or printed in pipeline output.
- Secret access audited (who accessed, when).

### GitHub Actions hardening
- Authenticate to cloud providers with OIDC federation (`id-token: write` plus a federated credential), not long-lived cloud keys stored as secrets
- Pin third-party actions to a full commit SHA (`uses: actions/checkout@<40-char-sha> # v4`), not a movable tag or branch
- Set minimal `permissions:` at workflow or job level (start from `contents: read`); never rely on the default token permissions

## Container practices (Docker/Kubernetes)

### Container image standards
- Non-root user and multi-stage builds with no build tools in the final image (see docker-dev.md).
- Security scanning before push (vulnerability scan, signed images if required).
- Image tagging: semantic versioning for releases. No `latest` tag in production deployments.
- Base images updated for security patches and deprecations.

### Kubernetes deployments (if applicable)
- Resource requests/limits defined (CPU, memory)
- Health checks: liveness and readiness probes
- Graceful shutdown: terminationGracePeriodSeconds, signal handling
- Pod Security Admission enforced with the `restricted` profile (no privileged containers, no root); PodSecurityPolicy was removed in Kubernetes 1.25
- Network policies restrict traffic between pods
- RBAC: service accounts with minimal permissions

## Security practices

### Access control
- Infrastructure access via IAM roles/RBAC, never shared credentials or shared SSH keys.
- Service accounts for automated tasks, with minimal permissions (least privilege).
- MFA for human access to production.
- Audit logging of who accessed or changed what resource, and when, retained per the compliance requirement.

### Network security
- Firewall rules: minimize exposed ports (443 for HTTPS, others internal only). Don't expose SSH/RDP (22/3389) to the internet; use a bastion (e.g. Azure Bastion) or just-in-time access
- VPC/VNet isolation: internal services not exposed to the internet
- DDoS protection and a WAF for public-facing web endpoints
- TLS everywhere: encrypt data in transit

## Backup and recovery
- Critical data (databases, file storage) backed up, with frequency set by the RPO and a retention policy defined per data set.
- Production backups geo-replicated to another region.
- Restores and failover tested regularly, not just documented. RTO defined, with a restore runbook.

## Monitoring and alerting
- Logs from all components sent to central logging in a structured format (JSON).
- Alerts for service down and for error rate or latency breaching the service's SLO (e.g. 5xx rate > 1% or p95 latency above target for 5 minutes), routed by severity, each with a runbook.
- Alert only on actionable issues, not information.

## Cost
- Dev/test use smaller instances than production and scale down when idle.
- Budget alerts set per subscription or account.
