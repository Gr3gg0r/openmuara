> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** Task spec for prompt 01. See `prompts/01-make-loader-manifest-first.md` for the high-level prompt.

# Task 01 — Make the Provider Loader Manifest-First

## Objective

Refactor `internal/config/provider_loader.go` so provider discovery is driven by `plugins/<name>/gateway.yml` rather than a hard-coded list of built-ins.

## Background

Currently the loader knows about built-in providers and may load them regardless of whether a manifest exists. After this change, the manifest is the source of truth.

## Detailed Steps

1. **Discover manifests**
   - Walk `plugins/` directory.
   - For each subdirectory, expect `gateway.yml`.
   - Parse with `internal/plugin/schema.go` types.
   - Validate using `internal/plugin/validator.go`.

2. **Route by runtime type**
   - `runtime.type == "simple"` → instantiate via `internal/provider/simple/`.
   - `runtime.type == "go"` → look up factory in registry (implemented in task 02; for now add a TODO/placeholder that compiles).
   - Unknown runtime type → validation error.

3. **Remove implicit built-in loading**
   - Do not instantiate providers that do not have a manifest.
   - Keep `default` provider as a hard-coded fallback.

4. **Update validation**
   - `internal/config/validation.go` should not assume any built-in provider names beyond `default`.

5. **Update tests**
   - Tests that assert built-ins are loaded must now load a manifest or use an explicit fixture.

6. **Observability**
   - Add startup logging: "loaded provider <name> with runtime <type>".
   - Log clear errors for missing/invalid manifests with file paths.

## Inputs

- `plugins/*/gateway.yml`
- `internal/plugin/schema.go`
- `internal/provider/simple/provider.go`

## Outputs

- Updated `internal/config/provider_loader.go`
- Updated `internal/config/validation.go` (if needed)
- Updated tests in `internal/config/`

## Test Plan

- Unit test: load a simple provider from a manifest fixture.
- Unit test: missing manifest means provider is not loaded.
- Unit test: invalid runtime type returns error.
- Regression: `default` provider still loads.

## Definition of Done

- `go test ./internal/config/...` passes.
- `go vet ./internal/config/...` passes.
- `golangci-lint run ./internal/config/...` passes.
- No references to hard-coded built-in provider names in loader logic.
- `TRACKING.md` and `HANDOFF.md` updated.

## Risks to Watch

- R001: Tests depending on global provider registration.
- R005: Loader ordering changes affecting existing users.
