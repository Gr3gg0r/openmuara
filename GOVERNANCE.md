# Governance

This document describes how the OpenMuara project is governed.

## Roles

### Project lead (BDFL)

- Sets high-level direction and final tie-breaker on disputes.
- Owns the `main` branch and release tags.
- Delegates day-to-day decisions to maintainers.

### Maintainers

- Review and merge pull requests.
- Triage issues and prioritize the backlog.
- Enforce quality gates and security practices.
- Can approve routine changes and most provider hardening work.

### Committers

- Have merge rights for their areas of expertise.
- Work under maintainer oversight for cross-cutting changes.

### Contributors

- Anyone who opens issues, pull requests, or docs improvements.
- Follow the contribution guidelines in [`CONTRIBUTING.md`](CONTRIBUTING.md).

## Decision-making

| Type | Process |
|---|---|
| Routine changes | Lazy consensus among maintainers; silence for 72 hours is approval. |
| Breaking changes | Explicit approval from at least two maintainers, including the project lead. |
| Provider contract changes | Explicit approval from at least two maintainers; prefer review from a maintainer familiar with the real gateway. |
| Security issues | Private disclosure to maintainers; public discussion only after a fix is released. |
| New maintainers | Nominated by an existing maintainer, approved by the project lead. |

## Conflict resolution

1. Discuss the issue in the relevant PR or issue.
2. If unresolved, escalate to the maintainers.
3. If still unresolved, the project lead makes a final decision.

## Becoming a maintainer

- Sustained, high-quality contributions over time.
- Constructive code and docs review.
- Understanding of OpenMuara's architecture and quality gates.
- Nominated by an existing maintainer and approved by the project lead.

## Code of conduct

All participants must follow [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
