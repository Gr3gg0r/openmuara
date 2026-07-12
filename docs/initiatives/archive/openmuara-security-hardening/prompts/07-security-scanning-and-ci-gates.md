> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 07 — Security Scanning and CI Gates

## Goal

Integrate static security scanners into CI and pre-commit workflows.

## Context

Automated scanning catches common Go security issues and leaked secrets before they reach `main`. We already run `govulncheck`; this prompt adds `gosec` and secret scanning.

## Required Output

1. Add `gosec` to CI:
   - New job in `.github/workflows/ci.yml` or extend `vuln` job.
   - Install `gosec` and run `gosec ./...`.
   - Fail on high/critical issues.
2. Add secret scanning to CI:
   - Use `gitleaks-action` or `trufflehog` in CI.
   - Optionally add to `.pre-commit-config.yaml`.
3. Add `scripts/check-secrets.sh` wrapper for local use (optional if using pre-commit).
4. Add `muara security audit` to smoke tests.
5. Update `runbooks/quality-gates.md` and `docs/security.md`.
6. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- CI remains fast; cache tool installs.
- Scans do not require paid services.
- `gosec` findings are actionable and not noisy.

## Quality Gate

- CI passes with new jobs.
- No high/critical `gosec` findings.
- No leaked secrets in history.
