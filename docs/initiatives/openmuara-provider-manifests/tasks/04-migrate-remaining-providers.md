> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** Task spec for prompt 04. See `prompts/04-migrate-remaining-providers.md` for the high-level prompt.

# Task 04 — Migrate Remaining Providers to Manifests

## Objective

Ensure every non-default provider has a valid `gateway.yml` manifest and loads through the manifest-first system.

## Background

Fawry and SenangPay already have `runtime.type: simple` manifests. iPay88, Billplz, and ToyyibPay have `runtime.type: go` manifests but still auto-register. Stripe has no manifest.

## Detailed Steps

1. **Audit existing manifests**
   - `plugins/fawry/gateway.yml`
   - `plugins/senangpay/gateway.yml`
   - `plugins/ipay88/gateway.yml`
   - `plugins/billplz/gateway.yml`
   - `plugins/toyyibpay/gateway.yml`

   Verify each has:
   - Valid `name`
   - Valid `runtime.type`
   - Required fields for its runtime

   Use `appendices/a-provider-contract-checklist.md`.

2. **Create Stripe manifest**
   - New file: `plugins/stripe/gateway.yml`
   - `runtime.type: go`
   - Include Stripe-specific config schema under `config:`.

3. **Register Stripe factory**
   - `internal/stripe/register.go` (if not done in task 02).

4. **Validate all manifests**

   ```bash
   go run ./cmd/muara provider validate plugins/fawry/gateway.yml
   # repeat for each provider
   ```

5. **Run conformance tests**
   - Ensure each migrated provider has passing conformance tests.
   - Add missing tests where needed.

6. **End-to-end test**
   - Start Muara.
   - Run `examples/checkout-store` Fawry flow.
   - Run `examples/checkout-store` Stripe flow.

7. **Update documentation**
   - `docs/provider-contract.md`
   - `docs/contributing-providers.md`
   - `docs/migration/provider-manifests.md` (if user-facing breaking change)
   - `CHANGELOG.md` release-notes snippet

## Inputs

- Existing provider packages.
- Existing `plugins/*/gateway.yml` files.

## Outputs

- Normalized manifests.
- New `plugins/stripe/gateway.yml`.
- Updated docs.

## Test Plan

- Manifest validation for every provider.
- Conformance tests for every provider.
- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`
- Manual checkout-store smoke test for Fawry and Stripe.

## Definition of Done

- Every non-default provider has a manifest.
- All manifests validate.
- All gates pass.
- Documentation is accurate.
- `TRACKING.md`, `HANDOFF.md`, `DECISIONS.md`, `RISKS.md`, `KNOWN_ISSUES.md` updated.
- `appendices/c-post-initiative-checklist.md` complete.

## Risks to Watch

- R003: Inconsistent existing manifests.
- R004: Stripe has no manifest.
- R009: Provider protocol emulation regression.
- R010: Documentation becoming stale.
