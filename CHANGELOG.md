# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.1] - 2026-07-12

### Documentation
- Completed a full documentation completeness audit (`DOCS01`):
  - Fixed inaccurate provider-status claims in `docs/mkp-billing-requirements.md`.
  - Reconciled `CONTRIBUTING.md` and `docs/contributing.md` into a single source
    of truth.
  - Hardened all provider docs with runnable `curl` examples, signature
    algorithms, simulation routes, webhook payloads, and error tables.
  - Added `docs/cli.md`, `docs/accessibility.md`, `docs/DOCS_STYLE.md`, and
    `docs/migration/openmuara-to-openmuara.md`.
  - Added `GOVERNANCE.md`, a docs issue template, and updated the Docusaurus
    sidebar to surface all reference and community docs.
  - Added `.markdownlint.yml` and verified the website build passes.

### Fixed
- Simple-runtime providers now correctly resolve `signature.secret_key` when the
  gateway manifest uses the full dotted path (e.g.
  `providers.senangpay.config.secret_key`). This makes the SenangPay charge
  example work out of the box.

### Added
- Dark mode support across the admin dashboard, provider simulation pages, and
  example mini-apps. Theme follows OS preference by default and persists via
  `localStorage`.
- Accessibility enhancements for the admin dashboard: skip-to-main-content link,
  `prefers-contrast: more` support, and improved WCAG AA contrast in dark mode.
- Playwright E2E accessibility smoke tests covering keyboard navigation, theme
  toggle, and axe-core critical violation checks.
- Automated WCAG AA color-contrast regression check (`npm run a11y:contrast`)
  and a corresponding CI job.

### Provider Manifests
- Manifest-first provider discovery: `config.LoadEnabledProvidersWithFallback`
  reads `plugins/<name>/gateway.yml` before falling back to built-ins.
- Go factory registry at `internal/provider/factory/` with package-level default
  registry and per-provider `internal/<provider>/register.go` files.
- Removed built-in `init()` auto-registration for non-default providers; the
  `default` provider remains hard-coded as the bootstrap fallback.
- Added `plugins/stripe/gateway.yml` so Stripe is discovered and activated like
  every other non-default provider.
- Normalized `plugins/*/gateway.yml` manifests for Fawry, SenangPay, iPay88,
  Billplz, ToyyibPay, and Stripe.
- Updated `docs/provider-contract.md` and `docs/contributing-providers.md` to
  describe the `simple` and `go` runtimes and the factory registration pattern.
- Added `docs/migration/provider-manifests.md` for users migrating from
  auto-registered built-in providers.

### Improved
- Unified provider payment and simulation pages (Fawry, Stripe Checkout,
  Stripe PaymentIntent, Billplz, ToyyibPay, iPay88) under a shared
  `internal/ui/payment-pages.css` stylesheet with a modern card layout,
  responsive design, consistent OpenMuara branding, and improved accessibility.
- Added `GET /__muara/payment-pages.css` route on both the combined and
  provider routers so payment pages load styles in single-port and dual-port
  modes.

### Bug Hunt / Quality
- Visual baseline capture script (`npm run test:visual-baseline`) and captured
  dashboard baselines under `web/dashboard/e2e/baselines/`.
- GitHub issue and PR templates aligned with the bug register format.
- CI hardening: `govulncheck`, `npm audit --production`, and `golangci-lint`
  required checks.
- Fuzz/property tests for signature verification, idempotency keys, and the
  transaction state machine.
- Provider contract conformance tests with golden files for every supported
  provider and version.
- Webhook dispatcher chaos tests for non-2xx responses, retry exhaustion, and
  timeouts.
- Root `KNOWN_ISSUES.md` register for intentionally deferred bugs.
- PR coverage-regression comment workflow.
- Error-code taxonomy in `internal/errcode` adopted across providers and
  webhook dispatch.

### Quality Automation Follow-Up
- Visual baseline diff job in CI for dashboard changes.
- Mutation testing gate with Gremlins for `internal/webhook`, `internal/engine`,
  and `internal/fawry` (70% threshold).
- Per-module coverage regression script and workflow, hardened against cached
  `go test` output and non-blocking during the phased rollout.
- Provider-wide adoption of `internal/errcode` for signature, config, validation,
  transaction, and webhook errors; public HTTP message text preserved.
- Updated `runbooks/quality-gates.md` with every new gate and its local command.

## [0.1.0] - 2026-06-29

### Added
- Initial stable release of OpenMuara.
- Fawry Express Checkout emulation with SHA256 signature verification.
- Stripe Checkout session emulation with webhook signature verification.
- SenangPay provider stub.
- SQLite and in-memory transaction ledger.
- Outgoing webhook dispatcher with replay and provider-specific payload builders.
- Admin dashboard at `/_admin` with transactions, webhooks, and provider status.
- Prometheus metrics endpoint at `/metrics`.
- Structured audit logging with `/_admin/audit` API and `muara audit list` CLI.
- CORS and CSRF protection for admin endpoints.
- OpenAPI spec served at `/openapi.yaml`.
- Docker and Docker Compose support.
- GitHub Actions CI and release workflows.
