> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 03 — Debuggability

> **Initiative:** OpenMuara v1 Solid Gold
> **Target:** `<repo-root>/`
> **Branch:** `feat/v1-solid-gold`
> **Depends on:** Prompt 01
> **⚠️ User approval required:** This prompt touches P0 webhook/provider integration logic.

---

## Goal

Make failed webhooks and provider calls trivial to trace end-to-end.

## Why now

When MKP or another consumer reports a missing or wrong webhook, the current
workflow requires reading server logs and correlating by reference and time.
A single trace ID across request → ledger → webhook dispatch would cut debug
 time dramatically.

## Scope

### In scope

- Propagate `X-Muara-Trace-ID` (or reuse the existing `trace_id`) on outgoing
  webhooks as an HTTP header only.
- Display the trace ID in the dashboard ledger and webhook detail panels.
- Add CLI commands:
  - `muara transaction inspect <ref>`
  - `muara webhook inspect <ref>`
- Add optional `/_admin/debug/pprof` endpoints gated by config.
- Update `runbooks/debugging.md`.

### Security / compatibility notes

- The trace-ID header must not change any webhook payload body, so provider
  signature verification stays intact.
- `pprof` must be disabled by default and enabled only in local/debug configs;
  OpenMuara is local-first, so no additional authentication is required.

### Out of scope

- Distributed tracing (OpenTelemetry).
- Changing webhook payload bodies.

## Acceptance criteria

- [ ] Outgoing webhooks include a trace-ID header.
- [ ] Dashboard shows the trace ID for transactions and webhook attempts.
- [ ] CLI inspect commands print payload, headers, attempts, and errors.
- [ ] pprof endpoints are off by default and enabled only via config.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Do not change the webhook payload body; only add an HTTP header.
- Reuse the request-scoped trace ID already logged by the router.

## Deliverables

- Code changes on `feat/v1-solid-gold`.
- Updated tests.
- Updated `runbooks/debugging.md`.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
