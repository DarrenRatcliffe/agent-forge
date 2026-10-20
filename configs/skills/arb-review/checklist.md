# ARB Review Checklists

Use these with the severity and condition rules in [SKILL.md](SKILL.md). Mark each item Met, Not met, Not addressed, or N/A.

## Business person

- [ ] The problem statement names who has the problem and how often
- [ ] Value is stated in measurable terms, and the basis for any number is given
- [ ] Build, run, and change costs are stated, or flagged as missing
- [ ] At least one alternative was considered, including do-nothing
- [ ] The proposal links to a stated strategic priority
- [ ] Customer impact is described, including changes for existing users
- [ ] Success measures and a review date are defined

## Lead architect

- [ ] Fits existing patterns in the repo, or the deviation is justified
- [ ] Integration points and contracts are identified
- [ ] Data ownership and the system of record are clear
- [ ] NFRs (latency, throughput, availability, scalability) have stated targets
- [ ] Failure modes and graceful degradation are described
- [ ] Observability (logs, metrics, traces) is addressed
- [ ] The decision is reversible, or the cost of reversal is stated
- [ ] Relevant ADRs exist or are proposed
- [ ] The standards in the matching installed rules are followed

## Chief engineer

- [ ] Scope is broken into deliverable increments
- [ ] Complexity is justified, and simpler alternatives were considered
- [ ] The test strategy covers the new behaviour (unit, integration, and end-to-end as relevant)
- [ ] Introduced tech debt is named, with a payback plan
- [ ] Operability is covered: deployment, rollback, runbooks, and on-call impact
- [ ] Run cost (infrastructure, licences, maintenance) is estimated
- [ ] The team has the skills, or a plan to acquire them
- [ ] New dependencies are maintained and version-pinned

## Lead CISO

- [ ] Data classification is stated for all data handled, and personal and regulated data is identified
- [ ] The threat model covers trust boundaries (STRIDE or equivalent)
- [ ] The authentication and authorisation model is defined, with least privilege applied
- [ ] Secrets are stored and rotated through an approved mechanism
- [ ] Encryption in transit and at rest is specified
- [ ] Sensitive actions are logged, with an audit trail
- [ ] Third-party and supply chain risk is assessed for new dependencies or vendors
- [ ] Compliance obligations are identified (for example GDPR, PCI DSS, or SOC 2, as relevant)
- [ ] Security testing is planned (SAST, dependency scanning, and penetration testing scaled to risk)

## Programme management lead

- [ ] Milestones and dates are given, or the proposal states that none exist
- [ ] Dependencies on other teams, systems, and vendors are listed with owners
- [ ] Resourcing (people, skills, and budget) is stated
- [ ] RAID items (risks, assumptions, issues, dependencies) are recorded with owners
- [ ] The rollout plan covers phases, feature flags, migration, and communication
- [ ] Change impact on users, support, and other teams is described
- [ ] Decision points and the governance cadence are defined
