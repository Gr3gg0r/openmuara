> **⚠️ AI AGENT: Read `AGENTS.md` and the initiative `README.md` first.**

# OpenMuara Stripe FPX — HOWTO

> **Status:** ❄️ Archived / Superseded. Active development continues in [`docs/initiatives/openmuara-stripe-checkout-sessions/`](../openmuara-stripe-checkout-sessions/).

## Decomposition

This initiative has one prompt: **P01 — Stripe FPX and Card Charge and Escape**.

### P01 breakdown

1. **Define FPX and card types**
   - Create `internal/stripe/fpx_types.go` and `internal/stripe/card_types.go`.
   - Define `FPXChargeRequest`, `FPXChargeResponse`, `FPXEscapeData`, and the card equivalents.

2. **Implement FPX handlers**
   - Create `internal/stripe/fpx.go`.
   - `NewFPXChargeHandler` validates request, creates `engine.Transaction`, stores it in ledger, returns response.
   - `NewFPXEscapeHandler` renders HTML with bank dropdown.
   - `NewFPXEscapeActionHandler` confirms/cancels, dispatches webhook, redirects.
   - Default currency is `myr` if not provided.
   - Generate reference as `fpx_test_` + UUID.

3. **Implement card handlers**
   - Create `internal/stripe/card.go`.
   - Mirror the FPX pattern with `card_test_` references and a card confirmation page (no bank dropdown).
   - Default currency is `usd` if not provided.

4. **Wire routes**
   - In `Provider.Routes()`, add:
     - `POST /v1/stripe/fpx/charge`, `GET /v1/stripe/fpx/escape`, `POST /v1/stripe/fpx/escape`
     - `POST /v1/stripe/card/charge`, `GET /v1/stripe/card/escape`, `POST /v1/stripe/card/escape`

5. **Test**
   - Unit tests for FPX and card creation, confirm, cancel, and error paths.
   - Regression test proving Checkout sessions are unaffected.

## Verification

After each sub-step, run:

```bash
go build ./...
go test ./internal/stripe/...
```

After the full prompt:

```bash
go test ./...
golangci-lint run
./scripts/smoke-test.sh
```
