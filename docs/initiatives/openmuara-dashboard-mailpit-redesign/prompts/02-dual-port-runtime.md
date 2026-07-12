> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P02 — Dual-Port Runtime

> **Initiative:** OpenMuara Dashboard — Mailpit-Style Redesign
> **Depends on:** —
> **Target files:** `internal/config/config.go`, `internal/server/server.go`, `cmd/muara/start.go`, `internal/server/router.go`, `web/dashboard/index.html`, `web/dashboard/src/api.ts`
> **Status:** ⬜

## Goal

Add an optional second port so the admin web UI/API can be separated from provider emulation endpoints, making it easy to expose only the API to a network while keeping the dashboard localhost-only.

## Tasks

- [ ] Add `admin_port` field to the server config (`internal/config/config.go`) with validation.
- [ ] Update `internal/server/server.go` to start a second listener when `admin_port` is set.
- [ ] Ensure provider emulation endpoints listen only on `server.port`.
- [ ] Ensure `/_admin` and all `/_admin/*` JSON endpoints listen only on `server.admin_port` when it is set.
- [ ] When `admin_port` is unset, preserve current single-port behavior.
- [ ] Update `cmd/muara/start.go` to log both URLs clearly on startup.
- [ ] Inject the admin API base URL into the dashboard HTML via a `<meta>` tag or `window.__MUARA_ADMIN_API__` so the SPA calls the correct port.
- [ ] Update `web/dashboard/src/api.ts` to read the admin API base URL from the injected value and fall back to `window.location`.
- [ ] Add backend tests for dual-port startup and endpoint isolation.

## Acceptance Criteria

- [ ] `server.admin_port` is optional and validated (must be a valid port number; must differ from `server.port` when both are set).
- [ ] When `admin_port` is set, `/_admin` is reachable on `admin_port` and provider routes are reachable on `port`.
- [ ] When `admin_port` is unset, the current single-port behavior is unchanged.
- [ ] The dashboard calls `/_admin/*` endpoints on the correct port in both modes.
- [ ] Startup logs show both the provider API URL and the admin UI URL.

## Quality Gates

Run before committing:

```bash
go build ./...
go test ./...
go vet ./...
golangci-lint run
cd web/dashboard && npm run test
cd web/dashboard && npm run build
node web/dashboard/scripts/check-bundle-size.js
```

## Notes

- This prompt changes core server startup, so review it as a P0 integration change per `AGENTS.md` if it touches routing or middleware ordering.
- Do not break existing tests or CLI behavior; the second listener must be strictly opt-in.
