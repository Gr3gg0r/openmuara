> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This document is subordinate to it.**

# Appendix A — Security Checklist

> **Updated:** 2026-07-06

Run this checklist during P01 and verify it after every fix batch. If a fix touches any of these areas, the corresponding item must be re-checked.

## Secrets & Sensitive Data

- [ ] No secrets logged or returned by admin/health endpoints.
- [ ] Environment variables are surfaced as names only, never values.
- [ ] Config backup files (e.g., `.muara/config.yml.bak`) do not contain plaintext secrets that differ from the original file.
- [ ] CLI logs and audit logs do not include webhook secrets or signature keys.

## Webhook Security

- [ ] Webhook signature verification is not bypassed in dev mode.
- [ ] SSRF protections on webhook test endpoints remain effective.
- [ ] Outgoing webhooks include the correct `X-Trace-Id` header but no internal stack traces.
- [ ] Webhook replay endpoints require the same CSRF/admin protections as the original dispatch.

## Admin & Config Endpoints

- [ ] Admin endpoints are protected by existing admin middleware when `admin.enabled: true`.
- [ ] Config write endpoints require CSRF token/header when `server.csrf.enabled: true`.
- [ ] Config writes create a `.muara/config.yml.bak` and detect external changes (`409 Conflict`).
- [ ] Rate-limit middleware still covers write endpoints.

## Input Validation & Injection

- [ ] No new `panic` paths introduced without recovery.
- [ ] No SQL injection or command injection through provider callbacks or CLI args.
- [ ] Provider callback parameters are validated before use.
- [ ] File-path and identifier inputs are sanitized.

## Infrastructure

- [ ] Dual-port runtime does not accidentally expose admin routes on the provider port.
- [ ] Optional pprof endpoints remain off by default.
- [ ] Health/readiness endpoints do not leak internal state.
