> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# OpenMuara Billplz — HOWTO

## Decomposition

This initiative has one prompt: **P01 — Billplz Provider**.

### P01 breakdown

1. **Provider scaffold**
   - Create `internal/billplz/provider.go`.
   - Implement `provider.Provider` interface.
   - Register as `billplz`.

2. **Collection types and handlers**
   - Create `internal/billplz/collection_types.go`.
   - Implement create/retrieve handlers.

3. **Bill types and handlers**
   - Create `internal/billplz/bill_types.go`.
   - Implement create/retrieve/delete handlers.
   - Link bills to collections.

4. **Local payment page**
   - Create `internal/ui/billplz-pay.html`.
   - Implement `GET/POST /_admin/billplz/pay/{id}`.

5. **Signature helper**
   - Create `internal/billplz/signature.go`.
   - Implement HMAC-SHA256 `X-Signature` generation and verification.

6. **Webhook/callback**
   - Create `internal/billplz/webhook.go`.
   - Implement `PayloadBuilder` and `PayloadHeaders`.
   - Implement callback redirect handler.

7. **Provider wiring**
   - Register routes in `provider.go`.
   - Mount in `internal/server/router.go`.

8. **OpenAPI**
   - Update `docs/openapi.yaml` and `internal/server/openapi.yaml`.

9. **Test**
   - Unit tests for all handlers and signature.
   - Smoke test.

## Verification

After each sub-step, run:

```bash
go build ./...
go test ./internal/billplz/...
```

After the full prompt:

```bash
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```
