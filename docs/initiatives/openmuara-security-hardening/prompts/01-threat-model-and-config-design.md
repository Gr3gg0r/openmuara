> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# 01 — Threat Model and Config Design

## Goal

Define the security threat model and the configuration schema for OpenMuara's hardening features.

## Context

OpenMuara is often run locally, but it may also run in CI/CD or on developer machines with ports forwarded. Without controls, `/_admin` and provider endpoints can leak or be abused.

**Critical constraint:** OpenMuara is a drop-in payment emulator. A user must be able to integrate it by changing only the base URL in their production provider client. Therefore, security controls must apply only to the admin/internal surface (`/_admin`, admin JSON APIs, simulation endpoints). Provider emulation endpoints (`/fawry/charge`, `/v1/checkout/sessions`, etc.) must remain contract-faithful to the real providers and must not gain additional auth requirements.

## Required Output

1. Update `DECISIONS.md` with resolved decisions D001–D005.
2. Propose config schema additions in `internal/config/`:
   - `server.bind` (default `127.0.0.1`)
   - `server.tls_cert`, `server.tls_key`
   - `admin.enabled` (default `false` for local dev)
   - `admin.username`
   - `admin.password_hash` (bcrypt) or `admin.token`
   - `hardened` (preset)
   - `rate_limit.enabled`, `rate_limit.requests_per_minute`
3. Create `docs/security.md` with the threat model and hardening guide outline.
4. Update `TRACKING.md` prompt 01 to `✅`.
5. Update `HANDOFF.md`.

## Decision Criteria

- Local dev remains zero-config and unauthenticated by default.
- Explicit opt-in is required for auth, TLS, and network exposure.
- Secrets must be configurable via environment variables.
- Hardened mode must be a single toggle that turns on a secure preset.

## Quality Gate

- Human review of `DECISIONS.md` and proposed config schema.
- `go build ./...` must still pass after config schema additions (no breaking changes).
