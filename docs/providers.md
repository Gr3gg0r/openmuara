---
id: providers
title: OpenMuara Providers
---

# OpenMuara Providers

OpenMuara emulates the payment providers below. Each provider is local-only and uses sample credentials so you can test without real accounts.

## Built-in providers

| Provider | Category | Real provider | First route | Docs |
|----------|----------|---------------|-------------|------|
| [fawry](providers/fawry.md) | regional | Fawry | `POST /fawry/charge` | [docs](providers/fawry.md) |
| [stripe](providers/stripe.md) | card | Stripe, Stripe Checkout, Stripe PaymentIntents | `POST /v1/checkout/sessions` | [docs](providers/stripe.md) |
| [billplz](providers/billplz.md) | redirect | Billplz | `POST /api/v3/bills` | [docs](providers/billplz.md) |
| [toyyibpay](providers/toyyibpay.md) | redirect | ToyyibPay | `POST /index.php/api/createBill` | [docs](providers/toyyibpay.md) |
| [senangpay](providers/senangpay.md) | redirect | SenangPay | `POST /senangpay/payment` | [docs](providers/senangpay.md) |
| [ipay88](providers/ipay88.md) | redirect | iPay88 | `POST /ePayment/entry.asp` | [docs](providers/ipay88.md) |
| [default](providers/default.md) | diy | OpenMuara Default | `POST /default/charge` | [docs](providers/default.md) |

## Choosing a provider

- **Testing Stripe integrations:** Use `stripe`.
- **Testing Fawry (Egypt):** Use `fawry`.
- **Testing Malaysian gateways:** Use `billplz`, `toyyibpay`, `senangpay`, or `ipay88`.
- **Quick experiments:** Use `default`.

Enable a provider in `.muara/config.yml`:

```yaml
providers:
  fawry:
    enabled: true
    config:
      merchant_code: muara-merchant-code
      merchant_security_key: muara-fawry-secret
      webhook_secret: muara-webhook-secret
      version: v1
```

## Contributor checklist

To add a new provider that appears in the wizard, dashboard, and docs:

1. Create a package under `internal/<provider>/` implementing `provider.Provider`.
2. Register it in `provider.Default()` during `init()`.
3. Add a `WizardChoice` entry in `internal/config/wizard.go`.
4. Add `providerCategory` and `realProvidersFor` entries in `internal/server/admin_api.go`.
5. Create `docs/providers/<provider>.md` with first route and config example.
6. Update this file's table.
7. Add tests and run all quality gates.
