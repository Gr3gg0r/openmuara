> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# OpenMuara iPay88 — HOWTO

## Decomposition

This initiative has one prompt: **P01 — iPay88 Provider**.

### P01 breakdown

1. **Provider scaffold**
   - Create `internal/ipay88/provider.go`.
   - Implement `provider.Provider` interface.
   - Register as `ipay88`.

2. **Payment request handler**
   - Create `internal/ipay88/payment_types.go`.
   - Implement `POST /epayment/entry.asp`.

3. **Signature helper**
   - Create `internal/ipay88/signature.go`.
   - Implement SHA256 signature algorithm.

4. **Local payment page**
   - Create `internal/ui/ipay88-pay.html`.
   - Implement `GET/POST /_admin/ipay88/pay/{refNo}`.

5. **Response/backend callbacks**
   - Create `internal/ipay88/webhook.go`.
   - Implement `PayloadBuilder`, response redirect, backend post.

6. **Requery**
   - Implement `POST /epayment/enquiry.asp`.

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
go test ./internal/ipay88/...
```

After the full prompt:

```bash
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```
