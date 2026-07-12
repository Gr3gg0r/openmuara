> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** This is prompt 3 of 4 for the `openmuara-provider-manifests` initiative. Read `README.md`, `DECISIONS.md`, `RISKS.md`, and `TRACKING.md` first.

# Prompt 03 — Remove Built-in Auto-Registration

## Goal

Remove all `init()` functions and side effects that auto-instantiate built-in providers. Provider creation must flow through the manifest loader → simple runtime or Go factory registry.

## Why

Built-in providers currently register themselves globally. This makes the loader's manifest-first behavior unreliable and pollutes tests.

## Acceptance Criteria

- [ ] `internal/fawry/provider.go` no longer auto-registers a provider instance.
- [ ] `internal/senangpay/provider.go` no longer auto-registers a provider instance.
- [ ] `internal/ipay88/provider.go` no longer auto-registers the provider instance.
- [ ] `internal/billplz/provider.go` no longer auto-registers the provider instance.
- [ ] `internal/toyyibpay/provider.go` no longer auto-registers the provider instance.
- [ ] `internal/stripe/provider.go` no longer auto-registers the provider instance.
- [ ] Factory `register.go` files from P02 are preserved.
- [ ] `internal/server/router.go` and any other callers do not depend on providers being pre-registered globally.
- [ ] `internal/provider/conform/conform_test.go` and `internal/server/providers_test.go` are updated to load manifests or use factories explicitly.
- [ ] A startup warning is added for configured providers that lack a manifest (resolves D007).
- [ ] `go test ./...` passes.

## Files to Touch

- `internal/fawry/provider.go`
- `internal/senangpay/provider.go`
- `internal/ipay88/provider.go`
- `internal/billplz/provider.go`
- `internal/toyyibpay/provider.go`
- `internal/stripe/provider.go`
- `internal/server/router.go`
- `internal/provider/conform/conform_test.go`
- `internal/server/providers_test.go`
- Any other files that call `provider.Get(...)` for built-in providers

## Out of Scope

- Changing provider business logic.
- Creating new manifests (that is prompt 04).
- `bridge` or `wasm` runtimes.

## Sign-off

- This touches P0 provider integration logic. Get human sign-off before committing.

## Verification

```bash
go test ./...
go vet ./...
golangci-lint run ./...
```

## Rollback Criteria

If a provider cannot be loaded through its manifest after auto-registration is removed, do not proceed to P04 until fixed.

## Observability

- Startup logs should warn when a configured provider has no manifest.
- `muara doctor` should report providers that are registered but not manifested.

## Commit Message

```
refactor(provider): remove built-in provider auto-registration

Providers are now instantiated only through gateway.yml + the simple runtime or Go factory registry.
```

## After This

Update `TRACKING.md` and `HANDOFF.md`, then proceed to `prompts/04-migrate-remaining-providers.md`.
