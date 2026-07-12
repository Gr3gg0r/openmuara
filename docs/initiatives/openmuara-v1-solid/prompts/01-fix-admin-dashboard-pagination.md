> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 01 — Fix Admin Dashboard for Paginated Responses

### Context

`/_admin/transactions` and `/_admin/webhooks` now return paginated envelopes (`{ limit, offset, results }`) instead of arrays. The embedded dashboard in `internal/ui/index.html` still expects arrays, so it displays no data.

### Current State

- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Target Files:** `internal/ui/index.html`, `internal/ui/handler_test.go`

### Scope

- **In scope:**
  - Update `loadTransactions()` to read `response.results`.
  - Update `loadWebhooks()` to read `response.results`.
  - Add or extend a test that verifies the dashboard HTML references the paginated shape.
- **Out of scope:**
  - Adding pagination controls to the UI.
  - Changing server response shape.

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./internal/ui/...
```

### Execution

1. Read `internal/ui/index.html`.
2. Update `loadTransactions()`:
   - `const data = await res.json();`
   - `const txs = data.results || [];`
3. Update `loadWebhooks()`:
   - `const data = await res.json();`
   - `const attempts = data.results || [];`
4. In `internal/ui/handler_test.go`, add a test that serves the dashboard and asserts the script contains `data.results` or equivalent.

### Quality Gates

```bash
go build ./...
go test ./internal/ui/...
go test ./...
golangci-lint run
```

### Commit

```bash
git add internal/ui/index.html internal/ui/handler_test.go
git commit -m "fix(ui): handle paginated admin API responses"
```

### Post-completion

1. Update `TRACKING.md` Step 01 to ✅.
2. Update `HANDOFF.md`.
