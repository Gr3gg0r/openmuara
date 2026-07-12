---
id: quickstart
title: OpenMuara Quick Start
---

# OpenMuara Quick Start

This guide gets you from zero to your first emulated charge, no matter how you
plan to use OpenMuara. Pick the path that matches you:

- [Developer](#developer-path) — integrate OpenMuara into your app.
- [AI Agent](#ai-agent-path) — explore and drive OpenMuara from the CLI.
- [Tester](#tester-path) — exercise flows and inspect results in the dashboard.
- [Contributor](#contributor-path) — build or extend OpenMuara itself.

---

## Developer path

### 1. Install

```bash
go install github.com/Gr3gg0r/openmuara/cmd/muara@latest
# or build from source
go build -o bin/muara ./cmd/muara
```

### 2. Initialize a local workspace

```bash
muara init
```

This creates `.muara/config.yml` with sensible defaults. Use `--defaults` to
skip the interactive wizard.

### 3. Start the server

```bash
muara start
```

The server listens on `127.0.0.1:9000` by default. Verify it is ready:

```bash
curl http://127.0.0.1:9000/readyz
```

### 4. Point your app at OpenMuara

Update your app's base URL to `http://127.0.0.1:9000` and configure the
provider you want to emulate. At least one provider must be enabled in
`.muara/config.yml` under `providers.<name>.enabled`.

### 5. Send your first charge

For the Fawry gateway:

```bash
REF="ref-$(date +%s)"
curl -X POST http://127.0.0.1:9000/fawry/charge \
  -H "Content-Type: application/json" \
  -d "{\"merchantCode\":\"muara-merchant-code\",\"merchantRefNum\":\"${REF}\",\"customerEmail\":\"test@example.com\",\"customerName\":\"Test\",\"customerProfileId\":\"user-123\",\"paymentExpiry\":9999999999999,\"language\":\"ar-eg\",\"chargeItems\":[{\"itemId\":\"prod_test_123\",\"price\":99.99,\"quantity\":1}],\"returnUrl\":\"http://127.0.0.1:9999/callback\"}"
```

For Stripe:

```bash
curl -X POST http://127.0.0.1:9000/v1/checkout/sessions \
  -u sk_test_anything: \
  -d "success_url=http://127.0.0.1:9999/success" \
  -d "line_items[0][price_data][currency]=usd" \
  -d "line_items[0][price_data][unit_amount]=2000" \
  -d "line_items[0][price_data][product_data][name]=T-shirt" \
  -d "line_items[0][quantity]=1" \
  -d "mode=payment"
```

See [`docs/providers.md`](providers.md) and [`docs/providers/`](providers/) for
provider-specific routes.

### 6. Inspect the ledger

Open `http://127.0.0.1:9000/_admin` and switch to the **Ledger** tab. It shows
transactions and webhook attempts in a single time-ordered feed. Click any row
to inspect details or replay a webhook.

---

## AI Agent path

Start with the CLI help surface, then use structured output to drive the tool
programmatically.

### Discover commands

```bash
muara --help
muara start --help
muara doctor --help
muara scenario --help
```

Every command includes runnable examples in its help text.

### Check the environment

```bash
muara doctor --json
```

The output follows the schema in [`docs/cli-schemas/doctor.json`](cli-schemas/doctor.json).

### Simulate a payment outcome

```bash
muara scenario success tx-123
```

Use `--json` to get machine-readable results:

```bash
muara scenario success tx-123 --json
```

### Explore the API

The OpenAPI spec is available at [`docs/openapi.yaml`](openapi.yaml) and served
live at `GET /openapi.yaml`.

Useful admin endpoints:

- `GET /_admin/ledger` — unified transaction and webhook feed.
- `GET /_admin/transactions` — searchable transaction list.
- `GET /_admin/webhooks` — webhook attempt history.
- `GET /_admin/providers` — enabled providers and metadata.

---

## Tester path

### Start the server

```bash
go build -o bin/muara ./cmd/muara
./bin/muara init
./bin/muara start
```

### Use the dashboard

Open `http://127.0.0.1:9000/_admin`. The default **Ledger** tab auto-refreshes
every 2 seconds while the page is visible and pauses when you switch browser
tabs.

Keyboard shortcuts:

- `?` — show help.
- `/` — focus the ledger search box.
- `1` / `2` / `3` — switch to Ledger / Transactions / Webhooks.
- `Esc` — close detail panels and help.

### Simulate a payment

For Fawry, send a charge request, then open the escape page:

```
http://127.0.0.1:9000/_admin/fawry-escape?ref=<REF>&returnUrl=<URL>&amount=<AMOUNT>
```

Click **Simulate Paid** or **Simulate Cancel**.

For Stripe, create a Checkout Session and follow the `_admin/stripe` links in
the response, or call the simulate endpoints directly:

```bash
curl -X POST http://127.0.0.1:9000/_admin/stripe/success \
  -d "payment_intent=pi_test_xxx"
```

### Debug webhooks

Go to the **Webhooks** tab (or open a webhook row from the ledger) to see:

- Request payload and headers.
- Signature verification status.
- Retry timeline and last error.
- A **Replay** button to re-send the webhook.

---

## Contributor path

### Clone and build

```bash
git clone https://github.com/Gr3gg0r/openmuara.git
cd openmuara
go build ./...
```

### Run quality gates

```bash
./scripts/check-scripts.sh
./scripts/check-forbidden.sh
./scripts/check-sizes.sh
./scripts/check-coverage.sh
go vet ./...
golangci-lint run
./scripts/smoke-test.sh
```

The project uses `task` as a convenience wrapper:

```bash
task check   # fmt, vet, lint, tests, coverage gate
task smoke   # E2E smoke test
task test    # tests with race detector and coverage
```

### Project layout

- `cmd/muara/` — CLI entry point.
- `internal/server/` — HTTP router and admin API.
- `internal/engine/` — transaction ledger.
- `internal/webhook/` — webhook dispatch and replay.
- `internal/provider/` and `plugins/` — provider plugins.
- `internal/ui/` — embedded dashboard.
- `docs/` — product and provider documentation.

### Submitting changes

Work on `feat/<description>` branches. Do not commit directly to `main`. All
changes should pass `go build ./...`, `go test ./...`, `go vet ./...`, and
`golangci-lint run` with zero warnings.

See [`CONTRIBUTING.md`](https://github.com/Gr3gg0r/openmuara/blob/main/CONTRIBUTING.md) for the full contribution guide.
