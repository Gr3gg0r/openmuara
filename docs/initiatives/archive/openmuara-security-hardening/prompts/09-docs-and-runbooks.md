> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 09 — Docs and Runbooks

## Goal

Document the security hardening features for users, operators, and contributors.

## Context

Security features are opt-in and must be discoverable. Docs should cover local dev, CI/CD, Docker, and shared environments.

## Required Output

1. Create/update `docs/security.md`:
   - Threat model summary.
   - Config reference (`server.bind`, TLS, admin auth, rate limits, hardened mode).
   - CLI helper reference.
   - CI/CD hardening guide.
   - Docker / shared environment guide.
2. Update `runbooks/local-development.md` with auth/TLS helpers.
3. Update `runbooks/on-call.md` with security event triage.
4. Update `runbooks/quality-gates.md` with `gosec` and secret scanning.
5. Update root `README.md` with a security section and link to `docs/security.md`.
6. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Docs are accurate and match implemented config schema.
- Examples are copy-pasteable.
- Secrets are never shown in examples.

## Quality Gate

- Docs build passes (if docs website exists by then).
- Markdown lint clean.
- Link checker passes.
