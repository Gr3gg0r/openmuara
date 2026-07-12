# Migration: Provider Manifests

In OpenMuara v1, provider discovery moved from hard-coded Go auto-registration
to manifest-first discovery via `plugins/<name>/gateway.yml`.

## What changed

- Every non-default provider now requires `plugins/<name>/gateway.yml`.
- `runtime.type: simple` providers need no Go code.
- `runtime.type: go` providers register a factory in
  `internal/provider/factory` instead of calling `provider.Register` in `init()`.
- Built-in Go providers are no longer auto-activated just because their package
  is compiled; the manifest controls activation.

## Am I affected?

You are affected only if your `.muara/config.yml` lists a provider and the
matching `plugins/<name>/gateway.yml` is missing. OpenMuara currently prints a
deprecation warning but still loads the provider for backwards compatibility.

```
provider "fawry" is configured but has no gateway.yml manifest;
auto-loading built-in providers is deprecated.
See docs/migration/provider-manifests.md
```

A future release will fail hard instead.

## Migration steps

1. Confirm the provider manifest exists:

   ```bash
   ls plugins/<name>/gateway.yml
   ```

2. If the file is missing, copy the existing manifest from the repository or
   create one. All built-in providers ship manifests:

   - `plugins/fawry/gateway.yml`
   - `plugins/senangpay/gateway.yml`
   - `plugins/stripe/gateway.yml`
   - `plugins/ipay88/gateway.yml`
   - `plugins/billplz/gateway.yml`
   - `plugins/toyyibpay/gateway.yml`

3. Restart muara. The warning should disappear.

## Custom providers

If you wrote a custom Go provider, add a `gateway.yml` manifest and replace any
`init()` registration with a factory registration. See
`docs/contributing-providers.md` for the exact `register.go` pattern.
