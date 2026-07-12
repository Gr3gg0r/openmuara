> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** This is prompt 4 of 4 for the `openmuara-provider-manifests` initiative. Read `README.md`, `DECISIONS.md`, `RISKS.md`, and `TRACKING.md` first.

# Prompt 04 — Migrate Remaining Providers to Manifests

## Goal

Ensure every non-default provider has a valid `plugins/<name>/gateway.yml` manifest and loads correctly through the new manifest-first system.

## Why

Fawry and SenangPay already have `runtime.type: simple` manifests. iPay88, Billplz, and ToyyibPay have `runtime.type: go` manifests but still auto-register. Stripe has no manifest. This prompt closes those gaps.

## Acceptance Criteria

- [ ] `plugins/fawry/gateway.yml` declares `runtime.type: simple` and loads.
- [ ] `plugins/senangpay/gateway.yml` declares `runtime.type: simple` and loads.
- [ ] `plugins/ipay88/gateway.yml` declares `runtime.type: go` and loads via factory.
- [ ] `plugins/billplz/gateway.yml` declares `runtime.type: go` and loads via factory.
- [ ] `plugins/toyyibpay/gateway.yml` declares `runtime.type: go` and loads via factory.
- [ ] `plugins/stripe/gateway.yml` is created with `runtime.type: go` and loads via factory.
- [ ] All manifests pass validation (`muara provider validate` or equivalent).
- [ ] Each migrated provider has a conformance test or existing one passes.
- [ ] `examples/checkout-store` works for Fawry and Stripe.
- [ ] `docs/provider-contract.md` and `docs/contributing-providers.md` are updated.
- [ ] A migration guide exists at `docs/migration/provider-manifests.md` if the loader change is user-facing.

## Files to Touch

- `plugins/fawry/gateway.yml`
- `plugins/senangpay/gateway.yml`
- `plugins/ipay88/gateway.yml`
- `plugins/billplz/gateway.yml`
- `plugins/toyyibpay/gateway.yml`
- New: `plugins/stripe/gateway.yml`
- `internal/stripe/` factory registration
- `docs/provider-contract.md`
- `docs/contributing-providers.md`
- New: `docs/migration/provider-manifests.md` (if needed)

## Out of Scope

- Bridge providers (`providers.<name>.type: bridge`).
- WASM sandboxed plugins.
- Changing provider protocol behavior.

## Sign-off

- P0 provider migration; get human sign-off before final commit.
- Provider domain experts (e.g., Fawry team) should review manifests if invited.

## Verification

```bash
go test ./...
go vet ./...
golangci-lint run ./...
muara provider validate plugins/*/gateway.yml
go run ./cmd/muara start &
# Verify checkout-store Fawry + Stripe flows
```

## Rollback Criteria

If any provider fails conformance or checkout-store smoke test, fix before closing the initiative.

## Observability

- `muara provider list` (or equivalent) should show all manifested providers.
- Startup logs should confirm each provider loaded with its runtime type.

## Commit Message

```
feat(plugins): migrate all non-default providers to gateway.yml manifests

Fawry and SenangPay use runtime.type: simple.
iPay88, Billplz, ToyyibPay, and Stripe use runtime.type: go via factory registry.
```

## After This

Update `TRACKING.md`, `HANDOFF.md`, `DECISIONS.md`, `RISKS.md`, `KNOWN_ISSUES.md`, and `CHANGELOG.md`. Run final gates and complete `appendices/c-post-initiative-checklist.md`.
