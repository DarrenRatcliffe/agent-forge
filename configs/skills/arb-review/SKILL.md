---
name: arb-review
description: Per-persona checklists, severity scale, and condition standards for an Architecture Review Board review of a proposal, design, or ADR. Use when running an ARB review or structuring architecture review findings.
---

# ARB Review

## Scope

This skill is used by the `arb-panel` agent. It defines what each persona checks, how findings are rated, and what a valid condition looks like. The per-persona checklists are in [checklist.md](checklist.md).

## Severity

- **Critical**: an incorrect premise or an unacceptable risk. Blocks approval.
- **High**: must be resolved, or made a testable condition, before implementation starts.
- **Medium**: should be addressed. It may become a condition with a later due date.
- **Low**: recorded for the file. Does not affect the decision.

## Conditions

Each condition must be specific, testable, and have an owner role.

- Weak: "Improve security."
- Strong: "Encrypt personal data at rest with keys from the approved KMS. Evidence: the storage configuration and a passing check in CI. Owner: lead CISO."

## Evidence

Cite the file, ADR, or statement in the proposal behind each finding. If the proposal makes a claim without support, mark it as unverified rather than accepting it.

## Checklist Status

For each checklist item, record one of the following:

- **Met**: the proposal addresses the item with evidence.
- **Not met**: the proposal addresses the item, but the answer is inadequate.
- **Not addressed**: the proposal does not mention the item. This counts as a finding.
- **N/A**: the item does not apply to this proposal. Say why.
