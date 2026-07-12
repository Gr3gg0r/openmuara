# OpenMuara

[![CI](https://github.com/openmuara/openmuara/actions/workflows/ci.yml/badge.svg)](https://github.com/openmuara/openmuara/actions/workflows/ci.yml)
[![Release](https://github.com/openmuara/openmuara/actions/workflows/release.yml/badge.svg)](https://github.com/openmuara/openmuara/actions/workflows/release.yml)
[![Coverage](https://img.shields.io/badge/coverage-80.3%25-brightgreen)](runbooks/quality-gates.md)
[![Container](https://img.shields.io/badge/container-ghcr.io-blue)](https://github.com/openmuara/openmuara/pkgs/container/openmuara)
[![License](https://img.shields.io/github/license/openmuara/openmuara)](LICENSE)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/openmuara/openmuara/badge)](https://scorecard.dev/viewer/?uri=github.com/openmuara/openmuara)

Local-first billing and payment virtualization layer. Emulate payment providers
(Stripe Checkout, Fawry, SenangPay, iPay88, Billplz, ToyyibPay) offline, fast,
and headlessly.

**Documentation:** [https://openmuara.github.io/openmuara/](https://openmuara.github.io/openmuara/)

## Install

### Recommended: release binary

Use the install script (macOS / Linux):

```bash
curl -sSL https://raw.githubusercontent.com/openmuara/openmuara/main/scripts/install.sh | bash
```

Or download a pre-built binary for your platform from the
[GitHub Releases](https://github.com/openmuara/openmuara/releases) page and extract it.

### From source

Clone the repo and build locally:

```bash
git clone https://github.com/openmuara/openmuara.git
cd openmuara
go build -o bin/muara ./cmd/muara
./bin/muara init
./bin/muara start
```

Or install directly with Go:

```bash
go install github.com/openmuara/openmuara/cmd/muara@latest
```

### Docker

A container image is published to `ghcr.io/openmuara/openmuara` on every release:

```bash
docker run --rm -p 127.0.0.1:9000:9000 \
  -e MUARA_SERVER_HOST=0.0.0.0 \
  -v "$(pwd)/.muara:/app/.muara" \
  ghcr.io/openmuara/openmuara:latest
```

Or use Docker Compose, which initializes a default config on first start:

```bash
docker compose up
```

## Quick Start

```bash
# Initialize a local workspace
muara init

# Start the server
muara start
```

The server listens on `127.0.0.1:9000` by default. Open
`http://127.0.0.1:9000/_admin` to see the ledger.

For guided paths tailored to developers, AI agents, testers, and contributors,
see [`docs/quickstart.md`](docs/quickstart.md).

## Examples

A maintained checkout-store example lives in [`examples/checkout-store/`](examples/checkout-store/):

```bash
cd examples/checkout-store/web
npm install && npm run build    # build React + DaisyUI SPA

cd ..
docker compose up --build -d    # OpenMuara + Mailpit
go run .                        # checkout store on :8080
```

It demonstrates a React + TypeScript + Vite + DaisyUI product landing page and
checkout SPA, Fawry and Stripe one-time payments, webhooks, and Mailpit email.
See [`examples/checkout-store/README.md`](examples/checkout-store/) for details.

## Configuration

OpenMuara reads `.muara/config.yml` at startup. Run `muara init` to create it from the bundled
defaults, or copy `muara.yml.example` and edit:

- `server.host` / `server.port` — bind address
- `log.level` — `debug`, `info`, `warn`, `error`
- `persistence.type` — `sqlite` (default) or `memory`
- `providers.<name>.enabled` — activate provider plugins
- `webhook.url` — local URL to receive outgoing webhooks
- `cors` / `csrf` — optional security settings

Environment variables override YAML values with the `MUARA_` prefix, e.g. `MUARA_SERVER_PORT=8080`.

See [`docs/operations.md`](docs/operations.md) for deployment, observability, and runbooks.

- `GET /healthz` — liveness probe.
- `GET /readyz` — readiness probe, lists enabled and available providers.
- Admin API responses (`/_admin/transactions`, `/_admin/webhooks`) are paginated as `{ limit, offset, results }`.
- Request body size is limited to 1 MiB by default.

## Example: Fawry-Style Charge Request

```bash
MERCHANT_CODE="muara-merchant-code"
SECRET="muara-fawry-secret"
REF="ref-$(date +%s)"
RETURN_URL="http://127.0.0.1:9999/callback"
ITEM_ID="prod_test_123"
PRICE="99.99"
QUANTITY="1"

MSG="${MERCHANT_CODE}${REF}user-123${RETURN_URL}${ITEM_ID}${QUANTITY}${PRICE}${SECRET}"
SIGNATURE=$(printf '%s' "$MSG" | shasum -a 256 | awk '{print $1}')

curl -X POST http://127.0.0.1:9000/fawry/charge \
  -H "Content-Type: application/json" \
  -d "{\"merchantCode\":\"${MERCHANT_CODE}\",\"merchantRefNum\":\"${REF}\",\"customerEmail\":\"test@example.com\",\"customerName\":\"Test\",\"customerProfileId\":\"user-123\",\"paymentExpiry\":9999999999999,\"language\":\"ar-eg\",\"chargeItems\":[{\"itemId\":\"${ITEM_ID}\",\"price\":${PRICE},\"quantity\":${QUANTITY}}],\"returnUrl\":\"${RETURN_URL}\",\"signature\":\"${SIGNATURE}\"}"
```

Visit `/_admin/fawry-escape?ref=<REF>&returnUrl=<RETURN_URL>` to simulate payment success or cancel.

## Development Commands

```bash
task check      # run fmt, vet, lint, tests, coverage, and UI checks
task test       # run tests with race detector and coverage
task coverage   # enforce 80% minimum coverage
task lint       # run golangci-lint
task vuln       # run govulncheck if installed
task security   # run gosec if installed
task secrets    # run gitleaks if installed
task smoke      # run the E2E smoke test
task ui:build   # build the dashboard SPA into the Go embed directory
task ui:test    # run dashboard unit tests
task dev        # run Go server + Vite dev server with HMR
./bin/muara doctor
```

The dashboard is a Vite + Preact SPA in `web/dashboard/`. Built assets are
embedded into the Go binary at `internal/ui/dashboard-dist/`. A tracked
placeholder `internal/ui/dashboard-dist/index.html` lets `go build ./...` work
on a fresh clone; run `task ui:build` to overwrite it with the real SPA. Do not
commit generated files inside `internal/ui/dashboard-dist/`.

The dashboard and all provider simulation pages support light and dark modes.
The theme follows the OS preference by default and can be toggled from the
dashboard header or persists automatically across pages via `localStorage`.

## Local Quality Gates

All commits should pass `task check`. Optionally install pre-commit hooks:

```bash
pre-commit install
pre-commit run --all-files
```

See [`runbooks/quality-gates.md`](runbooks/quality-gates.md) for the full local quality workflow.

## Security

OpenMuara is local-first by default: it binds to `127.0.0.1` and does not require admin
authentication. For CI/CD or shared environments, enable the opt-in hardening controls:

```yaml
# .muara/config.yml
server:
  host: 0.0.0.0
  tls_cert: /path/to/cert.pem
  tls_key: /path/to/key.pem

admin:
  enabled: true
  username: admin
  password_hash: "$2a$10$..."  # muara security hash-password

rate_limit:
  enabled: true
  requests_per_minute: 200

hardened: true
```

- `/_admin/*` and admin JSON APIs support HTTP Basic Auth or a bearer token.
- Provider emulation endpoints remain public and contract-faithful.
- Security features are lazy-initialized and add no overhead when disabled.

Helpers:

```bash
muara security hash-password --password mypassword
muara security gen-cert --host localhost --cert-out cert.pem --key-out key.pem
muara security audit
```

See [`docs/security.md`](docs/security.md) for the full hardening guide.

## Webhooks

OpenMuara can dispatch outgoing Fawry V2 webhooks to a local URL so you can test your webhook
handler without ngrok. Configure `webhook.url` in `.muara/config.yml`, then use the escape page or
`muara webhook` commands to inspect and replay.

See [`docs/webhooks.md`](docs/webhooks.md) for details.

## OpenAPI

The API spec is available at [`docs/openapi.yaml`](docs/openapi.yaml) and served live at `GET /openapi.yaml`.

## About Fawry Emulation

The default gateway emulates the Fawry Express Checkout contract used by `mkp/v1`, including the
SHA256 concatenated signature scheme. It is intended for local integration testing, not production
payment processing.

## Governance and Security

- [`GOVERNANCE.md`](GOVERNANCE.md) — maintainer roles and decision-making.
- [`SECURITY.md`](SECURITY.md) — reporting security issues.
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) — community standards.

## Contributing

We welcome contributions. Please read [CONTRIBUTING.md](CONTRIBUTING.md) and
[CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before participating.

## AI-assisted development

OpenMuara's code, documentation, and runbooks are developed with the assistance
of AI coding agents and reviewed by human maintainers. We treat AI-generated
output as a draft: it is tested, linted, and validated against the same quality
gates as human-written code before it is merged.

## License

MIT
