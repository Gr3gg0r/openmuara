> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This project is subordinate to it.**

# OpenMuara v1 — Known Issues & Out-of-Scope List

> **Purpose:** Prevent the AI from wasting time on pre-existing bugs or out-of-scope problems.

---

## Pre-Existing Bugs (Do NOT Fix Unless Directly Caused by This Project)

| ID | Issue | Location | Impact | Why Out of Scope |
|----|-------|----------|--------|------------------|
| K01 | `task` CLI not in PATH after container restart | Environment | Blocks quality gates | Fixed: installed via `go install github.com/go-task/task/v3/cmd/task@latest` |
| K02 | `go` not in PATH after container restart | Environment | Blocks build/test | Fixed: add `/usr/local/go/bin` to PATH |
| K03 | `golangci-lint` not in PATH after container restart | Environment | Blocks lint | Fixed: add `<go-bin>` to PATH |
| K04 | `/tmp` mounted with `noexec` causes `fork/exec ... permission denied` in tests | Environment | Blocks `go test` | Workaround: set `TMPDIR=<tmp-dir>` before running tests |

---

## Out-of-Scope Areas

| Area | Reason | Boundary |
|------|--------|----------|
| RevenueCat adapter | Hard freeze — v2 only | Do not implement in v1 |
| App Store / Play Store receipt validation | Hard freeze — v2 only | Do not implement in v1 |
| Multi-port runtime | Deferred to v1.2+ | Do not implement unless explicitly added to scope |
| MCP server | Deferred to v1.2+ | Do not implement unless explicitly added to scope |
| SaaS / hosted service | Out of project vision | Do not implement |

---

## Pre-Existing Test Failures

If the AI runs tests and sees failures BEFORE making changes, it MUST log them here and NOT attempt to fix them unless they are directly caused by the current project's changes.

| Test Suite | Failing Test | Error | Logged Date |
|------------|-------------|-------|-------------|
| | | | |

---

## Operational Limitations

These are expected boundaries of OpenMuara v1, not bugs to fix.

| ID | Limitation | Impact | Mitigation |
|----|------------|--------|------------|
| L01 | No built-in authentication on `/_admin`, `/metrics`, or provider routes | Anyone with network access can view/modify state | Run behind a reverse proxy or on `127.0.0.1` only |
| L02 | SQLite is a single-writer store | High write concurrency can return `database is locked` | Use one instance per environment; avoid parallel load tests |
| L03 | Audit log grows unbounded | Disk usage increases over time | Periodically archive or prune old `audit_logs` rows |
| L04 | Webhook retries are immediate, no backoff or dead-letter queue | Bursty failures may retry quickly | Fix the consumer promptly; replay manually after fixing |
| L05 | Metrics endpoint is unauthenticated | Metric counts (not payloads) are public within network scope | Bind to localhost or protect with a reverse proxy |
| L06 | CORS and CSRF settings are global | Cannot configure per-provider or per-route rules | Set origins for the whole server; disable CSRF only in isolated environments |

---

## How to Use This File

1. **Before starting a step:** Scan this file to know what landmines to avoid.
2. **During execution:** If you encounter a pre-existing bug unrelated to your task, STOP trying to fix it. Log it here with a `K##` ID.
3. **In the prompt:** Reference this file if the step touches code near a known issue.
