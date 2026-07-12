> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 08 — Automated Security Tests

## Goal

Add focused security tests that verify the hardened surface cannot be bypassed.

## Context

Security features are only valuable if they are tested. Add automated tests for auth, rate limiting, CSRF, and provider endpoint isolation.

## Required Output

1. Auth bypass tests:
   - `/_admin` returns 401 without credentials.
   - `/_admin` returns 200 with valid basic auth.
   - `/_admin` returns 200 with valid bearer token.
   - Wrong password/token is rejected.
2. Provider endpoint isolation tests:
   - `/fawry/charge`, `/v1/checkout/sessions`, `/senangpay/charge` remain unauthenticated.
   - Response shapes unchanged.
3. Rate limiting tests:
   - Admin endpoint blocks excessive requests from same IP.
   - Provider endpoints not rate-limited unless hardened.
4. CSRF tests:
   - State-changing admin actions still require CSRF token.
5. Header tests:
   - `/_admin` responses include CSP, X-Frame-Options, etc.
6. Config validation tests:
   - `0.0.0.0` + no auth/TLS triggers warning/audit event.
7. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Tests run as part of `go test ./...`.
- Tests do not require external services.
- Tests cover both hardened and non-hardened modes.

## Quality Gate

- `go test -race ./...`
- Coverage for new security packages ≥ 80%.
