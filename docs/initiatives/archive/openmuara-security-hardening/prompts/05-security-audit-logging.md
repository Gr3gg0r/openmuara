> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 05 — Security Audit Logging

## Goal

Log security-relevant events to the existing audit store and expose them via CLI/API.

## Context

Security events must be auditable without building a separate logging system. Reuse `internal/audit` SQLite store.

## Required Output

1. Define security event types:
   - `auth.failure` — failed admin login.
   - `auth.success` — successful admin login.
   - `config.change` — change to auth, TLS, or bind settings.
   - `replay.action` — webhook replay or transaction simulation.
   - `tls.enabled` / `tls.disabled`.
   - `rate_limit.triggered`.
2. Emit events from auth middleware, config reload, replay handlers, and TLS setup.
3. Add CLI/API to list security events:
   - `muara audit list --type=auth.failure`
   - Admin API `GET /_admin/api/audit?type=auth.failure`
4. Ensure raw webhook payloads or secrets are never logged.
5. Add tests.
6. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Events reuse existing audit schema or an additive migration.
- Secrets and payloads are redacted.
- No performance regression for non-security events.

## Quality Gate

- `go test ./...`
- Audit log tests verify redaction and event types.
