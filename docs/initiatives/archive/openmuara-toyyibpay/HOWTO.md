> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# OpenMuara ToyyibPay — HOWTO

## Decomposition

This initiative has one prompt: **P01 — ToyyibPay Provider**.

### P01 breakdown

1. **Provider scaffold**
   - Create `internal/toyyibpay/provider.go`.
   - Implement `provider.Provider` interface.
   - Register as `toyyibpay`.

2. **Category types and handlers**
   - Create `internal/toyyibpay/category_types.go`.
   - Implement create/get handlers.

3. **Bill types and handlers**
   - Create `internal/toyyibpay/bill_types.go`.
   - Implement create and get-bill-transactions handlers.

4. **Local payment page**
   - Create `internal/ui/toyyibpay-pay.html`.
   - Implement `GET/POST /_admin/toyyibpay/pay/{billCode}`.

5. **Webhook/callback**
   - Create `internal/toyyibpay/webhook.go`.
   - Implement `PayloadBuilder`, callback redirect.

6. **Provider wiring**
   - Register routes in `provider.go`.
   - Mount in `internal/server/router.go`.

7. **OpenAPI**
   - Update `docs/openapi.yaml` and `internal/server/openapi.yaml`.

8. **Test**
   - Unit tests for all handlers.
   - Smoke test.

## Verification

After each sub-step, run:

```bash
go build ./...
go test ./internal/toyyibpay/...
```

After the full prompt:

```bash
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```
