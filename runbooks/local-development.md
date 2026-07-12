---
id: local-development
title: Local Development Runbook
---

# Local Development Runbook

## Boot Order

1. Ensure Go 1.25+ is installed.
2. Clone the repo.
3. Build the CLI: `go build -o bin/muara ./cmd/muara`
4. Initialize workspace: `./bin/muara init`
5. Start server: `./bin/muara start`
6. Server is available at `http://127.0.0.1:9000`.
7. Check readiness: `curl http://127.0.0.1:9000/readyz`.

## Configuration

Default config is written to `.muara/config.yml` by `muara init`. Override with env vars prefixed by `MUARA_`, e.g.:

```bash
MUARA_SERVER_PORT=8080 ./bin/muara start
```

If you expose OpenMuara through a reverse proxy or tunnel for remote testing, set `server.public_base_url` to the external URL so generated payment links work:

```yaml
server:
  public_base_url: "https://muara.example.com"
```

## Common Commands

Every CLI command includes runnable examples. Run `./bin/muara <command> --help` to see them, or use `--json` / `--quiet` for scripting.

| Command | Purpose |
|---------|---------|
| `task check` | Run all quality gates |
| `task security` | Run gosec if installed |
| `task secrets` | Run gitleaks if installed |
| `go test -race ./...` | Run tests with race detector |
| `./scripts/smoke-test.sh` | Run E2E smoke test |
| `./bin/muara doctor` | Check environment |
| `./bin/muara doctor --json` | Structured health report |
| `./bin/muara version --json` | Version metadata for scripts |
| `./bin/muara scenario success tx-123` | Simulate a payment outcome |
| `./bin/muara security hash-password --password <pw>` | Generate bcrypt hash for config |
| `./bin/muara security gen-cert` | Generate self-signed TLS cert/key |
| `./bin/muara security audit` | Print security posture |

JSON output schemas are documented in `docs/cli-schemas/`.

## Running with Docker (local only)

Build and start with Docker Compose:

```bash
docker compose up --build
```

Or build and run the image manually:

```bash
docker build -t openmuara:latest .
mkdir -p .muara
docker run --rm -v "$(pwd)/.muara:/app/.muara" openmuara:latest init
docker run --rm -p 127.0.0.1:9000:9000 -v "$(pwd)/.muara:/app/.muara" openmuara:latest start
```

The dashboard is available at `http://127.0.0.1:9000/`.
The Compose healthcheck uses `GET /healthz`; use `GET /readyz` to confirm provider readiness.

## Testing the Fawry Flow

1. Send a signed `POST /fawry/charge` request.
2. Open `/_admin/fawry-escape?ref=<ref>&returnUrl=<url>&amount=<amount>`.
3. Click "Simulate Paid" or "Simulate Cancel".
4. The browser redirects to `returnUrl` with `orderStatus` and `statusCode`.

## Dashboard

The dashboard is a Vite + Preact SPA embedded into the Go binary at build time.
Source lives in `web/dashboard/` and built assets are output to
`internal/ui/dashboard-dist/`.

### Building the dashboard

```bash
# One-shot build (required before go build if dist is not present)
cd web/dashboard && npm install && npm run build

# Or use the task wrapper from the repo root
task ui:build
```

`internal/ui/dashboard-dist/index.html` is a tracked placeholder so that
`go build ./...` works on a fresh clone. Running `npm run build` or
`task ui:build` overwrites it with the real SPA. Do not commit the overwritten
generated file; the root `.gitignore` ignores `internal/ui/dashboard-dist/assets/`.

### Running the dashboard dev server

```bash
# Concurrent Go server + Vite dev server with HMR
task dev
```

The Vite dev server proxies API calls to the Go server on `127.0.0.1:9000`.
Open `http://127.0.0.1:5173/_admin` (or the Go server URL) to view the SPA.

### Production dashboard

Open `http://127.0.0.1:9000/_admin` for the embedded dashboard served by the
Go binary.

- The **Ledger** tab is the primary view. It merges transactions and webhook
  attempts into a single time-ordered feed, auto-refreshes every 2 seconds,
  and pauses when the browser tab is hidden.
- Use `?` to open keyboard-shortcut help, `/` to focus the ledger search, and
  `1`/`2`/`3` to switch tabs.
- Click any ledger row to inspect details or replay a webhook.
- The **Transactions** and **Webhooks** tabs keep the previous dedicated views
  available for focused debugging.

### E2E tests

The dashboard has a Playwright E2E test that verifies the SPA loads in hardened
mode when the URL contains HTTP Basic Auth credentials:

```bash
cd web/dashboard
npm run test:e2e
```

Playwright starts a temporary hardened OpenMuara server via `go run`, waits for
`/healthz`, and asserts that no `fetch` credential errors are emitted. Run
`npx playwright install chromium` first if browsers are not installed.

## Admin API Notes

- `GET /_admin/ledger` returns a unified, paginated feed of transactions and
  webhook attempts: `{ limit, offset, total, results }`.
- `GET /_admin/transactions` and `GET /_admin/webhooks` return paginated
  envelopes: `{ limit, offset, results }`.
- Provider-specific webhook targets can be configured per provider under
  `webhook.targets.<provider>`.
