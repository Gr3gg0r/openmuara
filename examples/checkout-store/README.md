# Checkout Store Example

A modern product landing page and checkout SPA built with **React**, **TypeScript**,
**Vite + SWC**, **Tailwind CSS**, and **DaisyUI**. It is backed by a small Go server
that accepts one-time payments through **ToyyibPay** and sends confirmation emails
through **Mailpit**.

The example is **locked to ToyyibPay** by default. The integration is written so the
same code runs against the real ToyyibPay sandbox and the OpenMuara emulator — only
`TOYYIBPAY_BASE_URL` (plus the matching key and category code) changes. This is the
seamless "develop against OpenMuara, ship against the real gateway" workflow.

## What it demonstrates

- React SPA with client-side routing.
- Beautiful, responsive landing and checkout pages using DaisyUI.
- Server-side creation of ToyyibPay bills (real sandbox or OpenMuara emulator).
- Redirect flows back to the SPA after payment.
- Webhook receiver that updates payment status and sends email.
- Local email testing with Mailpit.
- Switching between the real ToyyibPay sandbox and OpenMuara by changing one env var.

## Architecture

```
Browser → checkout-store (:8080) → ToyyibPay gateway
                          │        (real dev.toyyibpay.com OR OpenMuara :9000 emulator)
                          ↑
                          └────── webhook ← ToyyibPay / OpenMuara
                          ↓
                     Mailpit (:1025 SMTP / :8025 UI)
```

## Prerequisites

- Go 1.22+
- Node.js 20+ and npm
- Docker + Docker Compose
- OpenMuara checked out at the repo root

## Run

### 1. Start OpenMuara and Mailpit

```bash
cd examples/checkout-store
docker compose up --build -d
```

This starts OpenMuara on `http://127.0.0.1:9000` and Mailpit on
`http://127.0.0.1:8025`.

### 2. Configure OpenMuara (emulator path only)

This step is only needed when you point `TOYYIBPAY_BASE_URL` at the OpenMuara
emulator. If you are going straight to the real ToyyibPay sandbox, skip to step 3.

If this is the first time, initialize the workspace:

```bash
cd ../..
go run ./cmd/muara init
```

Then enable ToyyibPay and point webhooks at the checkout-store:

```yaml
# .muara/config.yml
providers:
  toyyibpay:
    enabled: true
    config:
      user_secret_key: muara-toyyibpay-secret
      category_code: cat_openmuara
webhook:
  url: "http://host.docker.internal:8080/webhook"
  max_retries: 3
```

Restart the OpenMuara container:

```bash
cd examples/checkout-store
docker compose restart muara
```

### 3. Build the frontend

The Go server embeds the built frontend from `web/dist`. Build it once:

```bash
cd examples/checkout-store/web
npm install
npm run build
```

For frontend development with hot reload, use Vite's dev server instead:

```bash
cd examples/checkout-store/web
npm run dev
```

The dev server proxies API calls to the Go backend on `:8080`.

### 4. Start the checkout store

```bash
cd examples/checkout-store
go run .
```

The store is now at `http://127.0.0.1:8080`.

### 5. Buy the product

1. Open `http://127.0.0.1:8080`.
2. Click **Buy now**.
3. Enter your name and email, then pay with **ToyyibPay**.
4. Complete payment on the ToyyibPay page (the real sandbox or the OpenMuara
   emulator, depending on `TOYYIBPAY_BASE_URL`).
5. You are redirected back to the success page.
6. Check Mailpit at `http://127.0.0.1:8025` for the confirmation email.

## ToyyibPay: real sandbox vs OpenMuara

ToyyibPay is disabled in OpenMuara by default. Enable it in your
`.muara/config.yml` before running the example against OpenMuara:

```yaml
providers:
  toyyibpay:
    enabled: true
    config:
      user_secret_key: muara-toyyibpay-secret
      category_code: cat_openmuara
```

The example reads three ToyyibPay variables. To switch between the real sandbox
and the emulator, change only `TOYYIBPAY_BASE_URL` and the matching key and
category code — no code changes.

| Target | `TOYYIBPAY_BASE_URL` | `TOYYIBPAY_USER_SECRET_KEY` | `TOYYIBPAY_CATEGORY_CODE` |
|---|---|---|---|
| Real ToyyibPay sandbox | `https://dev.toyyibpay.com` | your dev key | your dev category |
| OpenMuara emulator | `http://127.0.0.1:9000` | `muara-toyyibpay-secret` | `cat_openmuara` |

```bash
# Real sandbox
TOYYIBPAY_BASE_URL=https://dev.toyyibpay.com \
TOYYIBPAY_USER_SECRET_KEY=your-dev-key \
TOYYIBPAY_CATEGORY_CODE=your-dev-category \
go run .

# OpenMuara emulator (defaults)
go run .
```

> Note: real ToyyibPay sends the bill amount in **sen** (1/100 MYR). The example
> sends the product price multiplied by 100, so a `49.99` product becomes a
> `4999` bill amount. Adjust `product().Price` in `main.go` if you want a round
> ringgit value.

## Provider selection & demo-mode banner

The store is locked to ToyyibPay by default (`PAYMENT_METHODS=toyyibpay`). The
checkout page reads `GET /api/config` and only renders the methods listed in
`PAYMENT_METHODS`, and the backend rejects checkout requests for any method not
in the list. To also exercise the OpenMuara Fawry/Stripe emulators, opt back in:

```bash
PAYMENT_METHODS=toyyibpay go run .          # default — ToyyibPay only
PAYMENT_METHODS=toyyibpay,fawry,stripe go run .
```

A yellow announcement bar appears at the top of every page whenever an enabled
provider still uses its default placeholder credentials (for example
`muara-toyyibpay-secret`). The bar disappears once you set real keys in `.env`,
so it doubles as a reminder that the store is still in demo mode. The backend
reports this per provider via `configured` in `/api/config`.

## End-to-end tests

The frontend includes a Playwright e2e suite in `web/e2e/`.

### Prerequisites

- [Mailpit](https://mailpit.axllent.org/docs/install/) binary in your `PATH`.
- Playwright browsers installed (`npx playwright install chromium`).

### Run

```bash
cd examples/checkout-store/web
npm install
npm run test:e2e
```

Playwright starts OpenMuara on `:9001`, Mailpit on `:9035` (UI) / `:9025`
(SMTP), and the checkout store on `:8080`, then runs the tests.

### What is covered

- Fawry charge → escape page → success page.
- Stripe session → pay page → redirect back to store.
- Webhook receiver → Mailpit confirmation email.

> Note: OpenMuara's synchronous webhook dispatcher currently deadlocks with
> SQLite when `webhook.url` is configured, so the suite disables outgoing
> webhooks and tests the checkout-store webhook endpoint with a manual POST.

## Environment variables

| Variable | Default | Purpose |
|---|---|---|
| `ADDR` | `:8080` | Address the checkout store listens on |
| `APP_URL` | `http://127.0.0.1:8080` | Public URL used for redirects |
| `PAYMENT_METHODS` | `toyyibpay` | Comma-separated allow-list of methods shown on the checkout page |
| `OPENMUARA_URL` | `http://127.0.0.1:9000` | OpenMuara base URL |
| `FAWRY_MERCHANT_CODE` | `muara-merchant-code` | Fawry merchant code |
| `FAWRY_SECURITY_KEY` | `muara-fawry-secret` | Fawry security key |
| `STRIPE_SECRET_KEY` | `sk_test_muara` | Stripe secret key |
| `TOYYIBPAY_BASE_URL` | `http://127.0.0.1:9000` | ToyyibPay API base URL (real or OpenMuara) |
| `TOYYIBPAY_USER_SECRET_KEY` | `muara-toyyibpay-secret` | ToyyibPay user secret key |
| `TOYYIBPAY_CATEGORY_CODE` | `cat_openmuara` | ToyyibPay category code |
| `MAILPIT_HOST` | `127.0.0.1` | Mailpit SMTP host |
| `MAILPIT_PORT` | `1025` | Mailpit SMTP port |
| `MAIL_FROM` | `store@example.com` | Sender address for emails |

## Project files

| File / Directory | Purpose |
|---|---|
| `main.go` | Go backend (API, webhook, email, static file serving) |
| `web/` | Vite + React + TypeScript + DaisyUI SPA |
| `web/src/pages/` | Landing, checkout, and status pages |
| `web/src/api/client.ts` | API client |
| `web/dist/` | Built frontend (embedded by Go) |
| `docker-compose.yml` | OpenMuara + Mailpit services |
| `README.md` | This file |
