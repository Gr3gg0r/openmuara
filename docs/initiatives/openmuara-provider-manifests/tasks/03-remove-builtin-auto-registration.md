> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** Task spec for prompt 03. See `prompts/03-remove-builtin-auto-registration.md` for the high-level prompt.

# Task 03 — Remove Built-in Auto-Registration

## Objective

Delete all `init()` and global-registration side effects that instantiate built-in providers automatically. Provider instantiation must flow through the loader.

## Background

Built-in providers currently register themselves globally. This makes the loader's manifest-first behavior unreliable and pollutes tests.

## Detailed Steps

1. **Identify auto-registration points**
   - Search for `provider.Register`, `init()`, and global variable assignments in:
     - `internal/fawry/`
     - `internal/senangpay/`
     - `internal/ipay88/`
     - `internal/billplz/`
     - `internal/toyyibpay/`
     - `internal/stripe/`

2. **Remove provider-instance registration**
   - Delete `init()` functions that call a global registry to create instances.
   - Keep `register.go` files added in task 02; those only register factories.

3. **Update callers**
   - `internal/server/router.go` must get providers from the config/loader, not from a global registry.
   - Any code path that calls `provider.Get("<name>")` for built-ins must be updated.

4. **Add migration warning**
   - If a configured provider has no manifest, print a startup warning (resolves D007).
   - Point users to `docs/migration/provider-manifests.md`.

5. **Fix tests**
   - `internal/provider/conform/conform_test.go`
   - `internal/server/providers_test.go`
   - Other tests that assume global registration.

   Options:
   - Load a manifest fixture in test setup.
   - Use the factory registry directly.

6. **Run full gates**

## Inputs

- Code from tasks 01 and 02.

## Outputs

- Cleaned provider packages.
- Updated router and test files.
- Migration warning and guide.

## Test Plan

- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...`
- Verify no `provider.Register` calls remain in provider packages except factory registrations.

## Definition of Done

- No built-in provider auto-registers its instance.
- All gates pass.
- `examples/checkout-store/` still starts without errors.
- `TRACKING.md` and `HANDOFF.md` updated.

## Risks to Watch

- R001: Tests depending on global provider registration.
- R005: Loader ordering changes affecting existing users.
- R009: Provider protocol emulation regression.
