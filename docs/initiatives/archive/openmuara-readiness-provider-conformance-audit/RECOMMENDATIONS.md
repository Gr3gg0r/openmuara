> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# OpenMuara Readiness — Provider Conformance Audit Recommendations

> **Created:** 2026-07-09
> **Last Updated:** 2026-07-09
> **Status:** ✅ Complete — recommendations implemented; framework extended and tests added.

---

These recommendations define a gold-standard conformance framework for OpenMuara provider emulation.

## Conformance maturity model

| Level | Name | Coverage | State after initiative |
|---|---|---|---|
| L0 | Static surface | Routes, methods, paths, versions | ✅ All providers |
| L1 | Request contract | Required fields, headers, content types, validation errors | ✅ All P0 providers |
| L2 | Response contract | JSON shapes, status codes, error payloads, idempotency | ✅ All P0 providers |
| L3 | Signature verification | HMAC/SHA256, key derivation, negative/tampering tests | ✅ All P0 providers |
| L4 | Webhook dispatch | Payload shape, signature headers, retries, idempotency | ✅ All P0 providers |
| L5 | State transitions | Charge → authorize → capture → refund → fail | ✅ Stripe + Fawry explicit; regional gateways partial |
| L6 | External validation | Provider team review sign-off | ⬜ Request prepared; sending deferred |

## Current state assessment

After implementation, all P0 providers reach L1–L4. L5 is explicit for Stripe + Fawry and partial for regional gateways via engine scenarios.

| Provider | L0 | L1 | L2 | L3 | L4 | L5 | Notes |
|---|---|---|---|---|---|---|---|
| fawry | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `conformance_test.go` + `scenario_test.go`; v1/v2 charge, status, webhook, escape |
| stripe | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | `conformance_test.go` + `scenario_test.go`; Checkout + PaymentIntent flows |
| senangpay | ✅ | ✅ | ✅ | ✅ | ✅ | partial | `conformance_test.go`; charge + webhook; refund scenario deferred |
| billplz | ✅ | ✅ | ✅ | ✅ | ✅ | partial | `conformance_test.go`; collection + bill + webhook; refund scenario deferred |
| toyyibpay | ✅ | ✅ | ✅ | ✅ | ✅ | partial | `conformance_test.go`; bill + pay page + webhook; refund scenario deferred |
| ipay88 | ✅ | ✅ | ✅ | ✅ | ✅ | partial | `conformance_test.go`; entry + backend + webhook; refund scenario deferred |
| default | ✅ | partial | partial | n/a | n/a | partial | Reference provider; lower priority |

### Key gaps closed

1. ✅ L2 response-contract golden files added for all P0 providers.
2. ✅ L1 negative request tests added for missing required fields and invalid signatures.
3. ✅ L3 invalid/tampered signature tests added for all P0 providers.
4. ✅ L4 webhook payload and signature-header tests added for all P0 providers.
5. ✅ L5 state-transition scenarios added for Stripe and Fawry.
6. ✅ `internal/provider/conform` extended with `AssertJSONEqual` for behavior snapshots.
7. ⬜ External validation request prepared but not yet sent (deferred follow-up).

## Provider documentation references

During contract mapping (P02), record the official provider doc URLs and version numbers. Examples:

| Provider | Reference docs | Notes |
|---|---|---|
| fawry | Fawry API docs (v1/v2) | Versioned charge and webhook formats |
| stripe | Stripe API reference — PaymentIntents, Checkout Sessions, Webhooks | Versioned via `Stripe-Version` header |
| billplz | Billplz API v3 | Bills, collections, payment methods |
| toyyibpay | ToyyibPay API | CreateBill, callback, category |
| senangpay | SenangPay API | Payment form, callback, hash verification |
| ipay88 | iPay88 Enterprise API | Entry page, requery, signature |

*Exact URLs to be added during P02 mapping.*

## Priority matrix

| Priority | Area | Recommendation | Effort | Impact |
|---|---|---|---|---|
| P0 | Extend `conform` framework | Add `RequestSnapshot` and `ResponseSnapshot` to `internal/provider/conform` so golden files cover behavior, not just routes | S | High |
| P0 | Fawry contract completion | Map Fawry v1/v2 docs; close request/response/signature/webhook gaps; invite Fawry team to review | L | High |
| P0 | Stripe contract completion | Map Stripe Checkout + PaymentIntents; ensure idempotency-key behavior matches | L | High |
| P1 | Regional gateway contract matrix | Billplz, ToyyibPay, SenangPay, iPay88 — document emulated version and known limitations | M | High |
| P1 | Signature negative tests | Every P0 provider must have invalid-signature and tampering tests | M | High |
| P1 | Webhook conformance | Standardize webhook payload + signature tests across providers | M | High |
| P2 | State transition scenarios | Add explicit charge/refund/capture/fail scenarios for Stripe and Fawry | M | Medium |
| P2 | OpenAPI/schema generation | Future: generate provider request/response schemas from `gateway.yml` | L | Medium |
| P2 | Conformance dashboard | Surface per-provider conformance level in `/_admin` dashboard | M | Low |
| P3 | External review program | Establish process for provider teams to submit conformance corrections | S | Low |

## Conformance test patterns

### L1 — Request contract

```go
func TestFawryChargeRequiresMerchantCode(t *testing.T) {
    req := newFawryChargeRequest(t)
    delete(req.Body, "merchantCode")
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)
    assert.Equal(t, http.StatusBadRequest, rec.Code)
    assert.Contains(t, rec.Body.String(), "merchantCode")
}
```

### L2 — Response contract

```go
func TestFawryChargeResponseShape(t *testing.T) {
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, newFawryChargeRequest(t))
    var resp map[string]any
    require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
    assert.NotEmpty(t, resp["reference"])
    assert.NotEmpty(t, resp["status"])
}
```

### L3 — Signature verification

```go
func TestFawryChargeRejectsTamperedSignature(t *testing.T) {
    req := newFawryChargeRequest(t)
    req.Body["signature"] = "tampered"
    rec := httptest.NewRecorder()
    handler.ServeHTTP(rec, req)
    assert.Equal(t, http.StatusBadRequest, rec.Code)
}
```

### L4 — Webhook dispatch

```go
func TestFawryWebhookPayloadShape(t *testing.T) {
    p := newFawryProvider(t)
    tx := provider.Transaction{Reference: "ref-1", Status: "PAID"}
    payload, err := p.PayloadBuilder()(ctx, tx)
    require.NoError(t, err)
    var got map[string]any
    require.NoError(t, json.Unmarshal(payload, &got))
    assert.Equal(t, "ref-1", got["reference"])
}
```

## Recommended tool stack

| Purpose | Tool | Where |
|---|---|---|
| Static route snapshots | `internal/provider/conform` | Go tests |
| Request/response snapshots | Extended `conform` package | Go tests |
| Signature tests | Provider-specific `signature_test.go` + fuzz | Go tests |
| Webhook tests | `httptest` + `internal/webhook` | Go tests |
| Contract docs | `docs/providers/<provider>.md` | Docs |
| Limitation registry | `KNOWN_ISSUES.md` | Initiative docs |
| CI gate | Existing `go test ./...` + golden-file diff check | `.github/workflows/ci.yml` |
| Golden-file update | `UPDATE_GOLDEN=1 go test ./internal/provider/conform` | Local/CI |

## Standards mapping

| Recommendation | OpenSSF Scorecard | SLSA | CNCF |
|---|---|---|---|
| Contract tests for every provider | — | — | Testing |
| Golden-file regression | — | Build L2 | Quality |
| External provider review | — | — | Security |
| Documented limitations | — | — | Best Practice |

## Copy-paste command reference

```bash
# Run all provider conformance tests
go test ./internal/provider/conform/...

# Update golden files after intentional contract changes
UPDATE_GOLDEN=1 go test ./internal/provider/conform/...

# Run a single provider's contract tests
go test ./internal/fawry/... ./internal/provider/conform/...

# Check for undocumented deviations
grep -R "TODO.*provider\|FIXME.*provider\|not implemented" internal/ plugins/ docs/providers/
```

## What not to do

- Do **not** chase pixel-perfect emulation of provider quirks that do not affect integration correctness.
- Do **not** add conformance tests without linking them to provider doc references.
- Do **not** leave deviations undocumented.
- Do **not** block release on L5/L6 completion; ship L1–L4 first.

## Related documents

- [`TRACKING.md`](TRACKING.md) — execution phases
- [`KNOWN_ISSUES.md`](KNOWN_ISSUES.md) — gaps and deviations
- [`RISKS.md`](RISKS.md) — risk register
- [`DECISIONS.md`](DECISIONS.md) — decision log
- [`EXECUTION_PLAN.md`](EXECUTION_PLAN.md) — milestones and dependencies
