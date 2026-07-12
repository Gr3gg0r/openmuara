# Provider Contract

OpenMuara providers are discovered through `plugins/<name>/gateway.yml` manifests
and implement the `provider.Provider` interface defined in
`internal/provider/provider.go`. This document describes the contract, the
routing model, and the two runtimes: `simple` (YAML-driven) and `go`
(factory-driven).

## The `provider.Provider` interface

```go
type Provider interface {
    Name() string
    Init(cfg map[string]any) error
    Routes() []Route
    ChargeHandler() http.Handler
    WebhookHandler() http.Handler
    PayloadBuilder() func(ctx context.Context, tx Transaction) ([]byte, error)
    EscapeHandler() http.Handler
}
```

A provider must:

- Return a stable, lowercase `Name` that matches `metadata.name` in
  `gateway.yml` and the key in `providers.<name>` configuration.
- Validate its config in `Init` and return an `errcode` error for missing or
  invalid values.
- Declare all HTTP routes it wants to register via `Routes()`.
- Provide a charge handler (or a no-op handler if the provider uses a different
  primary endpoint).
- Provide a webhook handler and a payload builder for outgoing notifications.
- Return `nil` from `EscapeHandler()` when it has no escape/simulation page.

## Runtime wiring

1. `config.LoadEnabledProvidersWithFallback` walks `plugins/` and loads each
   `gateway.yml`.
2. `runtime.type: simple` providers are built from the manifest by
   `internal/provider/simple`.
3. `runtime.type: go` providers are built by looking up the provider name in
   `internal/provider/factory` and calling the registered factory.
4. The loader calls `Init` with the user's `providers.<name>.config`.
5. `cli/start.go` wires each loaded provider with the shared ledger, base URL,
   and dispatcher.
6. `server/router.go` registers each provider's routes from the loaded provider
   set.

The `default` provider is the only hard-coded exception: it has no manifest and
is registered by `internal/provider/defaultplugin` for bootstrapping and tests.

## Simple runtime

The simple runtime in `internal/provider/simple` is sufficient for gateways
that need:

- JSON request validation against a schema.
- Signature verification (`fawry_sha256`, `hmac_sha256`, `senangpay_md5`,
  `md5_concat`, `stripe_v1`, or `none`).
- Transaction creation in the shared ledger.
- A templated JSON response.
- An optional escape/simulation page.
- Optional webhook dispatch after a simulated payment outcome.

Action names in `routes` are matched heuristically:

- `*_charge` or `charge` → charge handler.
- `*_webhook` or `webhook` → webhook acknowledgement handler.
- `*_escape_page` or `escape_page` → escape page.
- `*_escape_action` or `escape_action` → escape outcome action.
- `*_status` or `status` → status query handler.

## Go runtime

Complex providers use `runtime.type: go`. They live in `internal/<provider>/`
and register a factory in `internal/provider/factory` from a `register.go` file:

```go
func init() {
    factory.MustRegister("my-provider", func(_ map[string]any) (provider.Provider, error) {
        return NewProvider(), nil
    })
}
```

The factory returns an unconfigured provider; the loader calls `Init`.
`runtime.type: go` providers still ship a `gateway.yml` manifest for discovery,
documentation, and validation.

Built-in Go providers must not call `provider.Register` in `init()`. The manifest
is the single source of truth for activation.

## Migration status

| Provider    | Runtime | Notes |
|-------------|---------|-------|
| fawry       | go      | Custom signature, multiple API versions, and escape page. |
| senangpay   | simple  | Fully expressible in `gateway.yml`. |
| stripe      | go      | Checkout Sessions, PaymentIntents, and signed webhooks. |
| ipay88      | go      | Form-encoded flow and local pay page. |
| billplz     | go      | Collection/bill storage and x_signature. |
| toyyibpay   | go      | Category/bill storage. |

To migrate a provider, add or update `plugins/<name>/gateway.yml`, remove any
`init()` registration, and add a `register.go` factory if the provider uses the
`go` runtime.
