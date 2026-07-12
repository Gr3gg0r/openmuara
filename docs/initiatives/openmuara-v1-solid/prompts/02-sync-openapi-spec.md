> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

## 02 — Sync OpenAPI Spec with Current API

### Context

The OpenAPI spec (`docs/openapi.yaml` and embedded `internal/server/openapi.yaml`) is out of sync with recent changes: `/readyz` is missing, paginated responses are still arrays, and 409 Conflict is not documented for refund/scenario.

### Current State

- **Repo:** `<repo-root>/`
- **Branch:** `dev`
- **Target Files:** `docs/openapi.yaml`, `internal/server/openapi.yaml`, `internal/server/openapi_test.go`

### Scope

- **In scope:**
  - Add `GET /readyz` path and schema.
  - Update `/_admin/transactions` response to paginated envelope.
  - Update `/_admin/webhooks` response to paginated envelope.
  - Add 409 response to `POST /v1/refund/{ref}` and `POST /_admin/scenario/{outcome}`.
  - Update `openapi_test.go` if response-shape assertions need adjustment.
- **Out of scope:**
  - Adding new endpoints.
  - Changing API behavior.

### Pre-flight

```bash
cd <repo-root>/
git status
git branch --show-current  # must be dev
go test ./internal/server/... -run TestOpenAPI
```

### Execution

1. Read current `docs/openapi.yaml` and `internal/server/openapi.yaml`.
2. Add `/readyz` path under `paths:`.
3. Define `ReadyResponse` schema in `components/schemas`.
4. Update `/_admin/transactions` and `/_admin/webhooks` 200 schemas to object wrappers with `results` array.
5. Add 409 `$ref` to `/v1/refund/{ref}` and `/_admin/scenario/{outcome}`.
6. Mirror all changes into `internal/server/openapi.yaml`.
7. Run `go test ./internal/server/...` to ensure sync test passes.

### Quality Gates

```bash
go build ./...
go test ./internal/server/...
golangci-lint run
```

### Commit

```bash
git add docs/openapi.yaml internal/server/openapi.yaml internal/server/openapi_test.go
git commit -m "docs(openapi): sync spec with readyz, pagination, and 409 responses"
```

### Post-completion

1. Update `TRACKING.md` Step 02 to ✅.
2. Update `HANDOFF.md`.
