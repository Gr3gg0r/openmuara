> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Dependency & License Incident Response Plan

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ⬜ Draft

---

This plan describes how to respond when a dependency or license issue is discovered in OpenMuara, such as an incompatible license, a reachable CVE, or a compromised package.

## Incident types

| Type | Example | Severity |
|---|---|---|
| License violation | A dependency is found to be GPL/AGPL or has no identifiable license | High |
| Reachable CVE | `govulncheck` or `npm audit` reports a reachable high/critical vulnerability | High |
| Compromised package | A dependency is hijacked or contains malicious code | Critical |
| Lockfile tampering | `go.sum` or `package-lock.json` is modified without a matching dependency change | Medium |
| False positive | A scanner flags a legitimate dependency incorrectly | Low |

## Response steps

### 1. Detect

- CI fails on `go-licenses check`, `govulncheck`, or `npm audit`.
- A user or security researcher reports an issue via `SECURITY.md`.
- Dependabot or a manual scan flags a problem.

### 2. Triage

- Identify the affected ecosystem(s): Go, npm, GitHub Actions, or container image.
- Determine whether the issue is reachable in production code.
- Check whether the issue is already tracked in `KNOWN_ISSUES.md` or `RISKS.md`.

### 3. Contain

- For a newly introduced bad dependency: revert the PR or commit that introduced it.
- For an existing dependency: pin to a safe version, apply a patch, or remove the dependency.
- For a compromised package: immediately remove it from lockfiles and block the version in CI if possible.

### 4. Remediate

- Replace the dependency with a compatible/safe alternative where possible.
- If no alternative exists, document the exception in `DECISIONS.md` with maintainer approval.
- Update `LICENSE_MATRIX.md` and SBOMs after the change.
- Run the full quality gate matrix.

### 5. Communicate

- For critical issues, open a security advisory per `SECURITY.md`.
- For high/medium issues, update `KNOWN_ISSUES.md` and notify maintainers.
- For false positives, add the package to the scanner allowlist and record the decision.

### 6. Learn

- After resolution, update `RISKS.md` and `DECISIONS.md` with lessons learned.
- Consider tightening CI checks if the incident revealed a gap.

## Rollback commands

```bash
# Revert a dependency change
git revert <commit-hash>

# Pin a Go dependency to a safe version
go get example.com/module@v1.2.3
go mod tidy

# Pin an npm dependency to a safe version
cd web/dashboard && npm install package@1.2.3
# or
npm update package

# Remove a dependency
cd web/dashboard && npm uninstall package
# Go: remove import and run go mod tidy
```

## Contacts

- Security issues: see `.github/SECURITY.md`.
- Dependency/license questions: open a discussion or issue in the repository.
