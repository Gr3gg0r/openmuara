> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** This is prompt 2 of 4 for the `openmuara-provider-manifests` initiative. Read `README.md`, `DECISIONS.md`, `RISKS.md`, and `TRACKING.md` first.

# Prompt 02 — Add a Go Factory Registry

## Goal

Create a registry where Go provider packages can register a factory function keyed by provider name. The loader uses this registry to instantiate providers whose manifest declares `runtime.type: go`.

## Why

We want Go providers to register a constructor without side effects, so the manifest controls whether they are instantiated. This removes phantom providers and makes provider presence explicit.

## Acceptance Criteria

- [ ] A new package exists at `internal/provider/factory/` (per recommendation RD006).
- [ ] It exposes `Register(name string, factory Factory)` and `Get(name string) (Factory, bool)`.
- [ ] `Factory` type signature is documented and minimal.
- [ ] The registry is safe for concurrent reads and is read-only after `init()`.
- [ ] Each existing Go provider (`ipay88`, `billplz`, `toyyibpay`, `stripe`) registers its factory in a dedicated `register.go` or `factory.go` file using an explicit `init()` that only calls `Register`.
- [ ] The loader from prompt 01 uses the registry for `runtime.type: go`.
- [ ] Removing a provider's manifest prevents it from loading, even if its factory is registered.
- [ ] Regression tests cover registry lookup, duplicate registration behavior, and missing factories.

## Files to Touch

- New: `internal/provider/factory/registry.go`
- `internal/ipay88/` — add factory registration
- `internal/billplz/` — add factory registration
- `internal/toyyibpay/` — add factory registration
- `internal/stripe/` — add factory registration
- `internal/config/provider_loader.go` — wire registry lookup

## Out of Scope

- Removing the old `init()` registrations that actually create the provider instance (that is prompt 03).
- Adding WASM support.
- Adding bridge support.

## Decisions Applied

- D006: Factory registry package is `internal/provider/factory/`.
- D008: Factory signature is `func(cfg map[string]any) (provider.Provider, error)`.

These are pinned in `DECISIONS.md`. Do not override without human sign-off.

## Sign-off

- If the factory signature changes the provider interface contract, get human sign-off.

## Verification

```bash
go test ./internal/provider/... ./internal/config/...
go vet ./internal/provider/... ./internal/config/...
golangci-lint run ./internal/provider/... ./internal/config/...
```

## Rollback Criteria

If the registry introduces global mutable state or data races, stop and refactor before proceeding.

## Observability

- `muara doctor` should list registered Go factories (optional but recommended).
- Startup logs should indicate when a factory is used for a provider.

## Commit Message

```
feat(provider): add Go factory registry for runtime.type: go

Go providers now register a factory instead of auto-instantiating.
The loader instantiates them only when their gateway.yml declares runtime.type: go.
```

## After This

Update `TRACKING.md` and `HANDOFF.md`, then proceed to `prompts/03-remove-builtin-auto-registration.md`.
