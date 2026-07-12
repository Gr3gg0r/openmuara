# Contributing Providers

This guide explains how to add a new payment provider to OpenMuara. Every
non-default provider is discovered through a `gateway.yml` manifest in
`plugins/<name>/`. The manifest decides whether the provider uses the YAML-driven
`simple` runtime or the Go `go` runtime.

## Choose a runtime

| Runtime | Use when | Location |
|---------|----------|----------|
| `simple` | REST-ish JSON flows, standard signature algorithms, templated responses, optional escape page. | Only `plugins/<name>/gateway.yml`. |
| `go` | Complex state machines, form-encoded flows, custom admin pages, unsupported signature schemes, multiple API versions. | `plugins/<name>/gateway.yml` + `internal/<name>/`. |

Start with `simple` and graduate to `go` when the simple runtime cannot cover a
behavior faithfully.

## Quick start: simple provider

Create a manifest:

```bash
muara provider init my-gateway
```

This creates `plugins/my-gateway/gateway.yml` with a starter configuration. Edit
it to match your provider's request shape, signature algorithm, and response
template.

### Required manifest sections

- `metadata`: `name`, `version`, `description`, `author`.
- `runtime.type: simple` and `runtime.simple`:
  - `charge_route`: action name of the charge route.
  - `currency`: default transaction currency.
  - `reference_field`: JSON field used as the transaction reference.
  - `amount_field`: JSON field used as the amount, or `charge_items` to sum line items.
  - `response_template`: JSON body returned by the charge handler.
  - Optional `escape_page` for the payment simulation page.
- `routes`: HTTP methods, paths, and action names.
- `schemas.requests`: validation rules for required fields.
- `signature`: algorithm, fields, and dotted `secret_key` path into provider config.
- `webhooks`: outgoing notification template.
- `fixtures`: sample requests/responses used by `muara provider test <name>`.

### Signature algorithms

| Algorithm       | Description                                            |
|-----------------|--------------------------------------------------------|
| `fawry_sha256`  | Fawry-style SHA256 of concatenated fields + secret.    |
| `hmac_sha256`   | HMAC-SHA256 of sorted `key+value` pairs joined by `\|`.|
| `senangpay_md5` | MD5 of `secret + detail + amount(2dp) + order_id`.     |
| `md5_concat`    | MD5 of `secret + concatenated field values`.           |
| `stripe_v1`     | Stripe-style webhook signature (documentation only).   |
| `none`          | Skip signature verification.                           |

### Testing

```bash
muara provider test my-gateway
```

The test command loads the manifest, exercises the charge handler with the
configured fixture, and prints the HTTP status and response body.

### Enabling the provider

Add to `.muara/config.yml`:

```yaml
providers:
  my-gateway:
    enabled: true
    config:
      secret_key: your-test-secret
```

The config key under `providers.my-gateway.config` must match the dotted path
in `signature.secret_key`.

## Built-in Go provider

Use the `go` runtime when the gateway requires logic the simple runtime cannot
express.

1. Create `plugins/<name>/gateway.yml` with `runtime.type: go`.
2. Create the Go package at `internal/<name>/`.
3. Export `NewProvider() provider.Provider` and implement the provider
   interface.
4. Add `internal/<name>/register.go` that registers a factory:

   ```go
   func init() {
       factory.MustRegister("<name>", func(_ map[string]any) (provider.Provider, error) {
           return NewProvider(), nil
       })
   }
   ```

5. Do **not** register the provider instance in `init()`. The manifest controls
   activation.

6. Import the new provider package in `internal/cli/start.go` so the factory is
   registered at build time:

   ```go
   _ "github.com/Gr3gg0r/openmuara/internal/<name>"
   ```

## Conventions

- Keep provider names lowercase and hyphenated.
- Re-use `internal/engine` for transactions and `internal/webhook` for
  dispatch.
- Return `errcode` errors from `Init` and handlers.
- Add tests in `internal/provider/simple/` for simple-runtime behavior and in
  `internal/<provider>/` for built-ins.
- Do not break existing provider contract tests in
  `internal/provider/conform`.

