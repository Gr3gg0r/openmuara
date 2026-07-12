# Checkout Store Example

A modern product landing page and checkout SPA built with **React**, **TypeScript**,
**Vite + SWC**, **Tailwind CSS**, and **DaisyUI**. It is backed by a small Go server
that accepts one-time payments through OpenMuara's **Fawry** and **Stripe** emulators
and sends confirmation emails through **Mailpit**.

## What it demonstrates

- React SPA with client-side routing.
- Beautiful, responsive landing and checkout pages using DaisyUI.
- Server-side creation of Fawry charges and Stripe Checkout sessions.
- Redirect flows back to the SPA after payment.
- Webhook receiver that updates payment status and sends email.
- Local email testing with Mailpit.

## Architecture

```
Browser → checkout-store (:8080) → OpenMuara (:9000) → Fawry/Stripe emulator
                          ↑
                          └────── webhook ← OpenMuara
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

### 2. Configure OpenMuara

If this is the first time, initialize the workspace:

```bash
cd ../..
go run ./cmd/muara init
```

Then enable Fawry and Stripe (they are enabled by default) and point webhooks
at the checkout-store:

```yaml
# .muara/config.yml
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
3. Enter your name and email, choose **Fawry** or **Stripe**, and pay.
4. For Fawry, complete payment on the OpenMuara escape page.
   For Stripe, complete payment on the OpenMuara checkout page.
5. You are redirected back to the success page.
6. Check Mailpit at `http://127.0.0.1:8025` for the confirmation email.

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
| `OPENMUARA_URL` | `http://127.0.0.1:9000` | OpenMuara base URL |
| `FAWRY_MERCHANT_CODE` | `muara-merchant-code` | Fawry merchant code |
| `FAWRY_SECURITY_KEY` | `muara-fawry-secret` | Fawry security key |
| `STRIPE_SECRET_KEY` | `sk_test_muara` | Stripe secret key |
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
