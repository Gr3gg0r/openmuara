> **⚠️ AI AGENT: Read `AGENTS.md` at the repo root first. This initiative is subordinate to it.**
>
> **Context:** This is prompt 1 of 4 for the `openmuara-provider-manifests` initiative. Read `README.md`, `DECISIONS.md`, `RISKS.md`, and `TRACKING.md` first.

# Prompt 01 — Make the Provider Loader Manifest-First

## Goal

Change `internal/config/provider_loader.go` so it discovers providers by reading `plugins/<name>/gateway.yml` first, and only uses built-in Go providers when the manifest explicitly declares `runtime.type: go`.

## Why

Today the loader knows about built-in providers and may load them regardless of whether a manifest exists. After this change, the manifest is the source of truth, which prevents phantom providers and makes the system predictable.

## Acceptance Criteria

- [ ] `LoadProviders` (or equivalent) walks `plugins/` and loads each `gateway.yml`.
- [ ] If a manifest declares `runtime.type: simple`, the provider is created via `internal/provider/simple/`.
- [ ] If a manifest declares `runtime.type: go`, the loader looks up a registered Go factory by provider name.
- [ ] If no manifest exists for a built-in provider, that built-in is **not** loaded.
- [ ] The `default` provider remains available as a fallback.
- [ ] Invalid manifests produce clear error messages that include the file path.
- [ ] Existing tests that assume built-ins are globally available are updated.
- [ ] Regression tests are added for manifest-first loading.

## Files to Touch

- `internal/config/provider_loader.go`
- `internal/config/validation.go` (if provider validation assumes built-ins)
- `internal/plugin/schema.go` (ensure `Runtime` block supports `type` lookup)
- Tests in `internal/config/`

## Out of Scope

- Implementing the Go factory registry itself (that is prompt 02).
- Removing `init()` registrations (that is prompt 03).
- Migrating providers that already have manifests.
- `bridge` or `wasm` runtimes.

## Sign-off

- No P0 integration logic changes expected.
- If the loader change affects config persistence or validation schemas, get human sign-off before committing.

## Verification

```bash
go test ./internal/config/...
go vet ./internal/config/...
golangci-lint run ./internal/config/...
```

## Rollback Criteria

If `go test ./...` fails after this change and cannot be fixed within the prompt, revert and update `RISKS.md`.

## Observability

- Add or update logs so startup shows which manifests were loaded and which runtime type was used.
- `muara doctor` should eventually list discovered providers (can be deferred to P02).

## Commit Message

```
feat(config): make provider loader manifest-first

Loader now discovers providers from plugins/<name>/gateway.yml.
Built-in Go providers are only loaded when runtime.type: go is declared.
```

## After This

Update `TRACKING.md` and `HANDOFF.md`, then proceed to `prompts/02-add-go-factory-registry.md`.
