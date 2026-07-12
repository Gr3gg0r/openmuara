> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**

# Appendix E — Test Scenarios

Specific test cases that must pass before the initiative closes. Each scenario maps to one or more prompts.

---

## TS001 — Simple provider loads from manifest

**Prompt:** P01
**Type:** Unit / integration
**Setup:** Create a test fixture `plugins/test-simple/gateway.yml` with `runtime.type: simple`.
**Action:** Call the loader.
**Expected:** Provider is instantiated via `internal/provider/simple/`.

---

## TS002 — Missing manifest means provider is not loaded

**Prompt:** P01
**Type:** Unit / integration
**Setup:** Ensure no `plugins/ghost/gateway.yml` exists; the Go package `internal/ghost/` may or may not be present.
**Action:** Call the loader.
**Expected:** `ghost` provider is not in the loaded provider set.

---

## TS003 — Go provider loads via factory when manifest present

**Prompt:** P02
**Type:** Unit / integration
**Setup:** Register a test factory for `test-go`; create `plugins/test-go/gateway.yml` with `runtime.type: go`.
**Action:** Call the loader.
**Expected:** Factory is invoked and provider is returned.

---

## TS004 — Factory registered but no manifest = no provider

**Prompt:** P02
**Type:** Unit / integration
**Setup:** Register a test factory for `test-go-no-manifest`; do not create a manifest.
**Action:** Call the loader.
**Expected:** `test-go-no-manifest` provider is not loaded.

---

## TS005 — Default provider still loads without a manifest

**Prompt:** P01
**Type:** Unit / integration
**Setup:** Standard config with no `plugins/default/gateway.yml`.
**Action:** Call the loader.
**Expected:** `default` provider is available.

---

## TS006 — Invalid runtime type returns clear error

**Prompt:** P01
**Type:** Unit
**Setup:** Create `plugins/bad-runtime/gateway.yml` with `runtime.type: unknown`.
**Action:** Call the loader.
**Expected:** Error message includes file path and valid runtime types.

---

## TS007 — Provider loads after removing `init()` registration

**Prompt:** P03
**Type:** Integration
**Setup:** Remove `init()` registration from a built-in provider; keep its manifest and factory registration.
**Action:** Call the loader and run conformance tests.
**Expected:** Provider loads correctly; protocol emulation unchanged.

---

## TS008 — Configured provider without manifest triggers warning

**Prompt:** P03
**Type:** Integration
**Setup:** `.muara/config.yml` references `billplz`; temporarily remove `plugins/billplz/gateway.yml`.
**Action:** Start Muara.
**Expected:** Warning logged with migration guide link; provider still loads for backwards compatibility.

---

## TS009 — Stripe manifest validates and loads

**Prompt:** P04
**Type:** Integration
**Setup:** Create `plugins/stripe/gateway.yml` with `runtime.type: go`; register Stripe factory.
**Action:** Run `muara provider validate plugins/stripe/gateway.yml` and start Muara.
**Expected:** Validation passes; Stripe provider is available.

---

## TS010 — Conformance tests pass for all migrated providers

**Prompt:** P04
**Type:** Integration
**Setup:** All providers have manifests and factories registered.
**Action:** Run `go test ./internal/provider/conform/...`.
**Expected:** Fawry, SenangPay, iPay88, Billplz, ToyyibPay, Stripe all pass.

---

## TS011 — Checkout-store Fawry flow works end-to-end

**Prompt:** P04
**Type:** E2E / manual
**Setup:** Muara running with Fawry enabled.
**Action:** Run the Fawry checkout flow in `examples/checkout-store`.
**Expected:** Payment succeeds; webhook delivered.

---

## TS012 — Checkout-store Stripe flow works end-to-end

**Prompt:** P04
**Type:** E2E / manual
**Setup:** Muara running with Stripe enabled.
**Action:** Run the Stripe checkout flow in `examples/checkout-store`.
**Expected:** Payment succeeds; webhook delivered.

---

## TS013 — No phantom providers in tests after auto-registration removal

**Prompt:** P03
**Type:** Unit
**Setup:** Fresh test binary; no manifest fixtures loaded.
**Action:** Query the provider registry for built-in names.
**Expected:** Only `default` is present (if queried globally); built-ins require manifest.

---

## TS014 — Factory registry is thread-safe and read-only after init

**Prompt:** P02
**Type:** Unit / race
**Setup:** Multiple goroutines call `factory.Get` and `factory.Names`.
**Action:** Run with `go test -race`.
**Expected:** No races; no runtime registrations.

---

## Traceability Matrix

| Scenario | P01 | P02 | P03 | P04 | Gate |
|---|---|---|---|---|---|
| TS001 | ✅ | — | — | — | Unit |
| TS002 | ✅ | — | — | — | Unit |
| TS003 | — | ✅ | — | — | Unit |
| TS004 | — | ✅ | — | — | Unit |
| TS005 | ✅ | — | — | — | Unit |
| TS006 | ✅ | — | — | — | Unit |
| TS007 | — | — | ✅ | — | Integration |
| TS008 | — | — | ✅ | — | Integration |
| TS009 | — | — | — | ✅ | Integration |
| TS010 | — | — | — | ✅ | Integration |
| TS011 | — | — | — | ✅ | E2E |
| TS012 | — | — | — | ✅ | E2E |
| TS013 | — | — | ✅ | — | Unit |
| TS014 | — | ✅ | — | — | Race |
