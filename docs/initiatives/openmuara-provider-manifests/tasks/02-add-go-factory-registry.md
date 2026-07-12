> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** Task spec for prompt 02. See `prompts/02-add-go-factory-registry.md` for the high-level prompt.

# Task 02 — Add a Go Factory Registry

## Objective

Create a registry that maps provider names to factory functions, allowing Go providers to be activated by a manifest declaring `runtime.type: go`.

## Background

We want Go providers to register a constructor without side effects, so the manifest controls whether they are instantiated.

## Detailed Steps

1. **Define the factory type**

   Suggested signature (adapt to actual provider interface):

   ```go
   type Factory func(cfg map[string]any) (provider.Provider, error)
   ```

   If shared dependencies are needed, extend to:

   ```go
   type Factory func(cfg map[string]any, deps Deps) (provider.Provider, error)
   ```

   Apply pinned decisions D006 and D008 from `DECISIONS.md`.

2. **Create the registry package**
   - Path: `internal/provider/factory/registry.go` (per recommendation RD006).
   - Thread-safe map.
   - `Register(name string, factory Factory)`.
   - `Get(name string) (Factory, bool)`.
   - `Names() []string` for discovery/doctor commands.
   - Panic or return error on duplicate registration; document behavior.

3. **Add registration files for each Go provider**
   - `internal/ipay88/register.go`
   - `internal/billplz/register.go`
   - `internal/toyyibpay/register.go`
   - `internal/stripe/register.go`

   Each file's `init()` should only call `hybrid.Register(name, factory)`.

4. **Wire into loader**
   - In task 01's loader, replace the placeholder with a real registry lookup.

5. **Update CLI doctor (optional)**
   - `muara doctor` can list registered Go factories for debugging.

## Inputs

- Existing Go provider constructors in `internal/<provider>/provider.go`.
- Provider interface definition (find in `internal/provider/`).

## Outputs

- `internal/provider/factory/registry.go`
- `internal/<provider>/register.go` for each Go provider
- Updated `internal/config/provider_loader.go`
- Updated `internal/cli/doctor.go` (optional)

## Test Plan

- Unit test: register and retrieve a factory.
- Unit test: factory is not called unless manifest is loaded.
- Unit test: duplicate registration behavior is defined.
- Integration test: load a Go provider from a manifest fixture.

## Definition of Done

- `go test ./internal/provider/... ./internal/config/...` passes.
- `go vet ./...` passes.
- `golangci-lint run ./...` passes.
- All Go providers have a `register.go` that only registers a factory.
- No global mutable state leaks after `init()`.
- `TRACKING.md` and `HANDOFF.md` updated.

## Risks to Watch

- R002: Registry design becoming too complex.
- R011: Mutable global state in registry.
