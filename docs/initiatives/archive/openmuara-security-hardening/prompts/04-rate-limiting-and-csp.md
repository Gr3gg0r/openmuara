> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 04 — Rate Limiting and Security Headers

## Goal

Add per-IP rate limiting and security headers for the admin surface.

## Context

Rate limiting protects against brute-force auth and accidental DoS from replay endpoints. It must be in-memory, bounded, and low-overhead — no Redis.

Security headers protect against clickjacking, MIME sniffing, and XSS on the admin UI.

## Required Output

1. Implement in-memory rate limiter:
   - Token bucket or sliding window per IP.
   - Bounded map size with TTL eviction.
   - Config: `rate_limit.enabled`, `rate_limit.requests_per_minute`.
   - Admin endpoints always rate-limited when enabled.
   - Provider endpoints rate-limited only when `hardened: true`.
2. Add security headers middleware for `/_admin/*`:
   - `Content-Security-Policy: default-src 'self'`
   - `X-Frame-Options: DENY`
   - `X-Content-Type-Options: nosniff`
   - `Referrer-Policy: strict-origin-when-cross-origin`
   - `Strict-Transport-Security` when TLS is enabled.
3. Preserve existing CORS and CSRF behavior.
4. Add tests for rate limiting and headers.
5. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Rate limiter uses O(1) or O(n) bounded memory.
- Stale entries are evicted automatically.
- Provider endpoints are unaffected unless hardened mode is on.
- Headers add negligible per-request overhead.

## Quality Gate

- `go test ./...`
- Benchmark rate limiter memory under many IPs.
- Verify headers on `/_admin` responses.
