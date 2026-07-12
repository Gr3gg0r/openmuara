---
id: operations
title: OpenMuara Operations Guide
---

# OpenMuara Operations Guide

This guide covers running, observing, and maintaining OpenMuara in a local or team environment.

---

## Deployment

### Release binary

1. Download the archive for your platform from [GitHub Releases](https://github.com/openmuara/openmuara/releases).
2. Extract the `muara` binary.
3. Run `muara init` to create `.muara/config.yml`.
4. Run `muara start`.

### From source

```bash
go build -o bin/muara ./cmd/muara
./bin/muara init
./bin/muara start
```

### Docker

```bash
docker build -t openmuara:latest .
mkdir -p .muara
docker run --rm -v "$(pwd)/.muara:/app/.muara" openmuara:latest init
docker run --rm -p 127.0.0.1:9000:9000 -v "$(pwd)/.muara:/app/.muara" openmuara:latest start
```

### docker-compose

Use the included `docker-compose.yml`:

```bash
docker compose up --build
```

The compose file mounts `.muara/` into the container so configuration and the SQLite ledger persist across restarts.

---

## Configuration

Runtime config lives at `.muara/config.yml`. A template is available at [`muara.yml.example`](../muara.yml.example). Key sections:

| Section | Purpose |
|---------|---------|
| `server.host` / `server.port` | Bind address (`127.0.0.1:9000` by default) |
| `server.admin_port` | Optional second port for the admin dashboard (dual-port mode) |
| `server.public_base_url` | External URL for payment links when behind a reverse proxy |
| `admin` | Full-access dashboard credentials |
| `viewer` | Read-only dashboard credentials for testers |
| `log.level` | `debug`, `info`, `warn`, `error` |
| `persistence.type` | `sqlite` (default) or `memory` |
| `persistence.path` | SQLite file, default `.muara/data/ledger.db` |
| `providers.<name>` | Enable and configure each provider |
| `webhook.url` / `webhook.max_retries` | Outgoing webhook target and retry count |
| `cors` / `csrf` | Optional CORS origin list and CSRF toggle; see `internal/config` defaults |

Environment variables override YAML values using the `MUARA_` prefix and dot-to-underscore mapping, e.g. `MUARA_SERVER_PORT=8080`.

### Hosted deployment with a reverse proxy

When OpenMuara runs behind a reverse proxy or tunnel (Cloudflare Tunnel, nginx, Traefik), set `server.public_base_url` to the external URL testers use:

```yaml
server:
  host: 127.0.0.1
  port: 9000
  public_base_url: "https://muara.example.com"
```

This ensures provider payment links (Stripe Checkout, Fawry, Billplz, etc.) point to the public domain instead of the internal bind address.

The provider port must remain reachable by end-user browsers, because checkout pages are served from that port. The admin port can be restricted to internal networks if you use dual-port mode.

### Read-only tester access

To give non-developer testers access to the dashboard without letting them change configuration, enable the `viewer` account:

```yaml
admin:
  enabled: true
  username: admin
  password_hash: "$2a$10$..."

viewer:
  enabled: true
  username: viewer
  password_hash: "$2a$10$..."
```

Viewers can inspect the ledger, transactions, and webhook log. They cannot enable or disable providers, change webhook targets, or replay webhooks.

---

## Health checks

OpenMuara exposes two health endpoints on the main HTTP port:

| Endpoint | Purpose |
|----------|---------|
| `GET /healthz` | Liveness. Returns `200 OK` when the process is running. |
| `GET /readyz` | Readiness. Returns `200 OK` with enabled providers, or `503 Service Unavailable` if a provider failed to initialize. |

Configure load balancers or container orchestrators to probe these paths.

---

## Metrics

Prometheus metrics are exposed at `GET /metrics`.

### Available metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `openmuara_requests_total` | Counter | `method`, `path`, `status` | Total HTTP requests |
| `openmuara_request_duration_seconds` | Histogram | `method`, `path` | Request latency |
| `openmuara_webhook_attempts_total` | Counter | `provider`, `status` | Webhook delivery attempts |
| `openmuara_transactions_total` | Counter | `provider`, `status` | Transactions created or updated |

`GET /metrics` is not counted in `openmuara_requests_total`.

### Scraping configuration

```yaml
scrape_configs:
  - job_name: openmuara
    static_configs:
      - targets: ['localhost:9000']
    metrics_path: /metrics
```

### Alerting rules example

```yaml
groups:
  - name: openmuara
    rules:
      - alert: OpenMuaraDown
        expr: up{job="openmuara"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: OpenMuara instance is down

      - alert: OpenMuaraNotReady
        expr: openmuara_ready == 0
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: OpenMuara is not ready

      - alert: OpenMuaraHigh5xxRate
        expr: |
          sum(rate(openmuara_requests_total{status=~"5.."}[5m]))
          /
          sum(rate(openmuara_requests_total[5m])) > 0.05
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: High 5xx rate on OpenMuara

      - alert: OpenMuaraWebhookDeliveryFailures
        expr: rate(openmuara_webhook_attempts_total{status="failed"}[5m]) > 0
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: Webhook deliveries are failing

      - alert: OpenMuaraHighLatency
        expr: |
          histogram_quantile(0.99,
            rate(openmuara_request_duration_seconds_bucket[5m])
          ) > 1
        for: 3m
        labels:
          severity: warning
        annotations:
          summary: OpenMuara P99 latency is above 1s
```

> **Note:** OpenMuara does not expose an `openmuara_ready` gauge out of the box. Use a probe or a simple exporter that calls `/readyz` and emits `openmuara_ready` if you need it in alerts. The examples above assume such an exporter is in place.

---

## Logs

OpenMuara writes structured JSON logs to stdout using Go's `log/slog`. Each log line includes:

- `time`
- `level`
- `msg`
- `trace_id` (per-request)
- `method`, `path`, `status`, `duration_ms` (HTTP requests)

Set `log.level` to `debug` to see provider-specific details, signature verification steps, and dispatcher decisions.

### Log aggregation

Because logs are emitted on stdout, route them with your container runtime or process manager:

- **Docker:** `docker logs -f <container>` or a log driver such as `fluentd`, `json-file`, or `journald`.
- **systemd:** set `StandardOutput=journal` and query with `journalctl -u openmuara`.
- **Kubernetes:** containers logs are picked up automatically by node-level log collectors.

Forward logs to your aggregation stack (Loki, Elasticsearch, CloudWatch, etc.) using the shipper of your choice.

---

## Audit log

The audit log records security-relevant events in a structured, queryable store.

| Interface | Command / endpoint |
|-----------|--------------------|
| CLI | `muara audit list --limit 100 --since 2026-06-01T00:00:00Z` |
| API | `GET /_admin/audit?limit=100&offset=0` |
| Dashboard | `/_admin` |

Logged events include provider initialization, configuration reloads, charge/refund/scenario actions, and webhook delivery attempts.

---

## Backup and restore

All durable state is in `.muara/`:

- `.muara/config.yml` — runtime configuration
- `.muara/data/ledger.db` — SQLite transaction and audit log database
- `.muara/data/unified_matrix.json` — receipt lookup matrix (when used)

To back up:

```bash
cp -r .muara .muara-backup-$(date +%Y%m%d)
```

To restore, stop the server and replace `.muara/` with the backup. SQLite supports online backups with `.backup` as well:

```bash
sqlite3 .muara/data/ledger.db ".backup to /path/to/backup.db"
```

---

## Scaling limits

OpenMuara is designed as a local-first, single-process emulator.

- **Single-node only.** There is no clustering or replicated state.
- **Single SQLite writer.** Concurrent writes are serialized through one shared `*sql.DB` connection. Heavy concurrent write workloads may hit `database is locked`.
- **Memory-based webhook attempt store.** Webhook attempts live in memory and are lost on restart.
- **Auth is opt-in.** Enable `admin.enabled` and `hardened` for shared networks; use `viewer` for read-only tester access. See [Security](security) and [Hosted Testing](hosted-testing).

For larger or multi-user setups, consider running separate instances per developer or test environment.

---

## Security notes

- Admin routes, metrics, and the dashboard are unauthenticated by default. Enable `admin.enabled` for shared environments.
- Use the `viewer` role to give non-developer testers read-only dashboard access without letting them change configuration or replay webhooks.
- CSRF protection is enabled for `/_admin/*` mutations using a double-submit cookie (`openmuara_csrf`).
- CORS origins are configurable; the default is restrictive.
- Webhook signature verification uses the configured `webhook_secret`.
- Do not expose OpenMuara directly to the public internet. Run it behind a reverse proxy or tunnel and set `server.public_base_url`.
- See the [Hosted Testing Guide](hosted-testing) for Mailpit-style deployment patterns.

---

## Runbooks

See the following runbooks for common operational tasks:

- [Local Development](/runbooks/local-development) — build and run locally
- [Quality Gates](/runbooks/quality-gates) — test and lint workflow
- [On-Call](/runbooks/on-call) — alerts and first response
- [Debugging](/runbooks/debugging) — inspect state, replay webhooks, fix issues
