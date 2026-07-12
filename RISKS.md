# Risks

This document tracks the main operational and security risks for OpenMuara v1.

## Data loss

OpenMuara stores all state in a local SQLite file (`.muara/data/ledger.db`) by
default. There is no replication, automated backup, or remote persistence. If the
workspace directory is deleted or the disk fails, transaction history and audit
logs are lost.

**Mitigation**

- Back up `.muara/` regularly when running long-lived test suites.
- Use the `memory` persistence mode only for ephemeral smoke tests.

## Not a real payment provider

OpenMuara emulates provider protocols but does not settle funds, validate real
credit cards, or interact with live App Store / Play Store / Stripe / Fawry
backends. A passing test against OpenMuara does not guarantee the same behavior
in production.

**Mitigation**

- Run a final integration test against each provider's sandbox before shipping.
- Treat OpenMuara as a development and CI tool, not a production payment gateway.

## No authentication by design

The HTTP API and `/_admin` dashboard are intentionally unauthenticated to keep
local development frictionless. Exposing OpenMuara to a public network would let
anyone create transactions, replay webhooks, and inspect the ledger.

**Mitigation**

- Bind to `127.0.0.1` (default) and never expose port `9000` publicly.
- Run OpenMuara inside a private network or behind an authenticated reverse proxy
  if remote access is required.

## Local secrets in configuration

Provider credentials and webhook secrets live in plain YAML files under
`.muara/config.yml`. These values are emulated defaults, but users may paste real
sandbox keys during testing.

**Mitigation**

- Add `.muara/` to `.gitignore` and never commit it.
- Prefer environment variables (`MUARA_*`) for real keys.
- Rotate any secret that was accidentally written to disk.

## Webhook misconfiguration

An invalid or unreachable `webhook.url` silently skips deliveries. A typo can
cause tests to pass while webhooks never leave OpenMuara.

**Mitigation**

- Check the startup log and `/_admin/webhooks` for failed attempts.
- Validate the URL with `muara config validate`.

## Provider emulation gaps

Some provider-specific behaviors are simplified:

- App Store / Play Store receipts are treated as lookup keys; no real crypto
  validation is performed.
- SenangPay, iPay88, Billplz, and Razer Merchant Services have partial route
  coverage compared to Stripe / Fawry.
- Provider webhooks are dispatched synchronously in-process, not from an
  external queue.

**Mitigation**

- Track known limitations in `KNOWN_ISSUES.md`.
- Add sandbox-level contract tests before relying on a provider emulation for
  release gating.
