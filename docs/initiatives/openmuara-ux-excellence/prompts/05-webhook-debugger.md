> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# P05 — Webhook Debugger

> **Initiative:** OpenMuara UX Excellence
> **Target:** `<repo-root>/`
> **Branch:** `feat/ux-excellence`
> **Depends on:** —

---

## Goal

Let users inspect webhook payloads, signature verification status, and retry history from the dashboard without reading server logs.

## Why now

Webhook failures are the most common "why isn't my integration working" question. The dashboard shows status but hides the payload and signature context.

## Scope

### In scope

- Extend `webhook.Attempt` with a `signature_valid` field (or `signature_status`) that the dispatcher populates after building the payload and headers.
- Add `GET /_admin/webhooks/{ref}` endpoint that returns:
  - `reference`
  - `provider`
  - `url`
  - `payload` — JSON body that was sent.
  - `headers` — relevant headers (redact any real secrets if present).
  - `signature_valid` — whether OpenMuara considers the signature valid.
  - `attempts` — list of `{time, status, error}`.
- Update `internal/ui/index.html` so clicking a webhook row opens a detail panel with the above.
- Keep one-click replay in the detail panel.
- Add tests for the new endpoint.

### Out of scope

- Editing payloads before replay.
- Webhook comparison with real provider signatures.

## Acceptance criteria

- [ ] `GET /_admin/webhooks/{ref}` returns payload, headers, signature status, and attempts.
- [ ] Dashboard shows a detail panel when a webhook row is clicked.
- [ ] Replay works from the detail panel.
- [ ] Tests cover the endpoint and redaction.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- The webhook dispatcher already stores attempts; payload may need to be stored if not already retained.
- Redact values for headers whose names case-insensitively contain `signature`, `authorization`, `token`, or `secret`.

## Deliverables

- Code changes on `feat/ux-excellence`.
- Updated `internal/server/admin_api_test.go` and webhook tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit.
