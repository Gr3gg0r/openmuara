> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This prompt is subordinate to it.**

# Prompt 02 — Response Delay Config

> **Initiative:** OpenMuara MKP Fawry Integration
> **Target:** `<repo-root>/`
> **Branch:** `feat/mkp-fawry`
> **Depends on:** —

---

## Goal

Add a configurable `response_delay_ms` to the Fawry provider and apply it to the
outgoing webhook dispatch so MKP can test slow-gateway behavior.

## Why now

MKP wants to simulate real-world gateway latency without slowing down the
escape-page redirect itself.

## Scope

### In scope

- Add `response_delay_ms` to the Fawry plugin config.
- Parse and validate the value in `internal/fawry/plugin.go`.
- Pass the delay to the Fawry escape action handler.
- Sleep for the configured duration before calling `dispatcher.Dispatch`.
- Default to `0` (no delay).
- Add tests that verify the delay is applied.

### Out of scope

- Delaying the synchronous charge response.
- Per-request delay override headers.

## Acceptance criteria

- [ ] `fawry.response_delay_ms` is accepted in `.muara/config.yml` and via
      `MUARA_FAWRY_RESPONSE_DELAY_MS`.
- [ ] When set to a positive value, the outgoing webhook is delayed by that
      amount but the escape redirect still returns immediately.
- [ ] When set to `0`, behavior is unchanged.
- [ ] All quality gates pass:
  - [ ] `go build ./...`
  - [ ] `go test ./...`
  - [ ] `go vet ./...`
  - [ ] `golangci-lint run`
  - [ ] `./scripts/smoke-test.sh`

## Hints

- Do the sleep inside the dispatch goroutine or before dispatch; do not block
  the HTTP redirect on it.
- Consider using `time.AfterFunc` or a background goroutine so tests can assert
  dispatch timing without real waits.

## Deliverables

- Code changes on `feat/mkp-fawry`.
- Updated plugin and escape tests.
- Updated `TRACKING.md` and `HANDOFF.md`.
- Git commit with a clear message.
