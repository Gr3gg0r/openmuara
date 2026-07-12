> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# 02 — Admin Authentication

## Goal

Add configurable authentication for `/_admin/*` and admin JSON APIs only. Provider emulation endpoints must remain unauthenticated.

## Context

Local dev remains zero-config and unauthenticated by default. Auth is opt-in via config or `--hardened`. We support two mechanisms:

- HTTP Basic Auth (`admin.username` + `admin.password_hash`).
- Bearer token (`admin.token`).

A request satisfies auth if either mechanism is valid. This makes CI and headless testing easy.

## Required Output

1. Implement `internal/server/auth.go` with:
   - Basic auth middleware.
   - Bearer token middleware.
   - Bcrypt password verification.
   - Constant-time token comparison.
2. Apply auth middleware only to `/_admin/*` routes and admin JSON APIs.
3. Add config fields:
   - `admin.enabled` (default `false`)
   - `admin.username`
   - `admin.password_hash`
   - `admin.token`
4. Support env vars: `MUARA_ADMIN_PASSWORD_HASH`, `MUARA_ADMIN_TOKEN`.
5. Use bcrypt cost 10 by default (configurable via constant).
6. Add unit and integration tests.
7. Update `TRACKING.md`, `DECISIONS.md`, `RISKS.md`, and `HANDOFF.md`.

## Decision Criteria

- Provider endpoints return unchanged response shapes and status codes.
- Bcrypt cost is sensible (default 10) to balance security and speed.
- Token comparison is constant-time.
- Auth middleware is skipped when `admin.enabled == false`.

## Quality Gate

- `go test ./...`
- Provider endpoint tests verify no auth requirement.
- Admin endpoint tests verify auth is enforced.
- `golangci-lint run`
